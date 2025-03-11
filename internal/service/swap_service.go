package service

import (
	"game-fun-be/internal/constants"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpRequest"
	"game-fun-be/internal/pkg/httpRespone"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"
	"log"

	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type SwapServiceImpl struct {
	userInfoRepo *model.UserInfoRepo
}

func NewSwapService() *SwapServiceImpl {
	return &SwapServiceImpl{}
}

func (s *SwapServiceImpl) GetSwapRoute(req request.SwapRouteRequest, chainType uint8) response.Response {

	startTime := time.Now()

	solPriceUSD, priceErr := getSolPrice()
	if priceErr != nil {
		return response.Err(http.StatusBadRequest, "price query failed", priceErr)
	}
	if req.SwapType == "buy" && chainType == 1 {
		solMultiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(model.SOL_DECIMALS)))
		req.InAmount = req.InAmount.Mul(solMultiplier)
	}
	inAmountUint64 := req.InAmount.BigInt().Uint64()
	inAmountStr := strconv.FormatUint(inAmountUint64, 10)

	// Get token and pool details
	tokenDetail, poolDetail, errResp := s.getTokenAndPoolInfo(req.TokenAddress, chainType)
	if errResp != nil {
		return s.handleErrorResponse(errResp)
	}

	// Process Anti-MEV logic
	mev, jitotip, jitoOrderId, errResp := s.processAntiMev(req)
	if errResp != nil {
		return s.handleErrorResponse(errResp)
	}

	// Create map for platform-specific functions
	platformHandlers := map[string]func() (*httpRespone.SwapTransactionResponse, error){
		"pump": func() (*httpRespone.SwapTransactionResponse, error) {
			swapStruct := s.buildSwapPumpStruct(req, tokenDetail, poolDetail, mev, jitotip)
			return s.getPumpFunTradeTx(swapStruct)
		},
		"raydium": func() (*httpRespone.SwapTransactionResponse, error) {
			swapStruct := s.buildSwapRaydiumStruct(req, tokenDetail, poolDetail, mev, jitotip)
			return s.getRaydiumTradeTx(swapStruct)
		},
		"g_external": func() (*httpRespone.SwapTransactionResponse, error) {
			swapStruct := s.buildGameFunGInstructionStruct(req, inAmountStr)
			return s.getGameFunGInstruction(swapStruct)
		},
		"g_points": func() (*httpRespone.SwapTransactionResponse, error) {
			pointsDecimal := decimal.NewFromFloat(req.Points) // 用户输入的 Points（代币数量）
			// 计算 10^decimals，确保计算精度
			pointsDecimalsFactor := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(tokenDetail.Decimals)))
			// 去掉小数点部分，取整
			pointsWithoutDecimal := pointsDecimal.Mul(pointsDecimalsFactor).Floor()
			// 将去掉小数点后的结果转换为字符串
			pointsString := pointsWithoutDecimal.String()
			// 计算代币的 USD 价值（代币数量 * 代币单价 * 10^精度）
			tokenAmount := pointsDecimal.Mul(tokenDetail.Price).Mul(pointsDecimalsFactor)
			// 计算 SOL 价格 * 10^SOL_DECIMALS
			solMultiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(model.SOL_DECIMALS)))
			solAmount := solPriceUSD.Mul(solMultiplier)
			// 计算需要多少 SOL（代币 USD 价值 / SOL 的 USD 价值）
			req.InAmount = tokenAmount.Div(solAmount)

			swapStruct := s.buildBuyGWithPointsStruct(req, pointsString)
			return s.getGetBuyGWithPointsInstruction(req.Points, swapStruct)
		},
	}

	// Handle platform-specific logic
	handler, exists := platformHandlers[req.PlatformType]
	if !exists {
		return response.Err(http.StatusBadRequest, "Unsupported platform type", errors.New("unsupported platform"))
	}

	swapTransaction, err := handler()
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to send swap request", err)
	}
	if swapTransaction.Code != 2000 {
		return response.Err(http.StatusInternalServerError, "Failed to send swap request", errors.New(swapTransaction.Message))
	}

	// Calculate amounts
	platform := model.CreatedPlatformType(tokenDetail.CreatedPlatformType)
	inDecimals := model.SOL_DECIMALS
	outDecimals := platform.GetDecimals()

	outAmount, inAmountUSD, outAmountUSD, errResp := s.calculateSwapAmounts(req, tokenDetail, solPriceUSD, inDecimals, outDecimals)
	if errResp != nil {
		return s.handleErrorResponse(errResp)
	}

	// Return the constructed response
	return ConstructSwapRouteResponse(req, swapTransaction, uint8(inDecimals), uint8(outDecimals), outAmount, inAmountUSD, outAmountUSD, startTime, jitoOrderId)
}

// Helper function to handle error responses
func (s *SwapServiceImpl) handleErrorResponse(errResp *response.Response) response.Response {
	return response.Err(errResp.Code, errResp.Msg, errors.New(errResp.Error))
}

func (s *SwapServiceImpl) getTokenAndPoolInfo(tokenAddress string, chainType uint8) (*model.TokenInfo, *model.TokenLiquidityPool, *response.Response) {
	tokenDetail, err := model.GetTokenInfoByAddress(tokenAddress, chainType)
	if tokenDetail == nil {
		return nil, nil, &response.Response{
			Code: http.StatusOK,
			Msg:  "token not found",
		}
	}
	if err != nil {
		return nil, nil, &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get token detail",
		}
	}

	poolDetail, err := QueryAndCheckPool(tokenAddress, chainType, 1)
	if poolDetail == nil {
		return nil, nil, &response.Response{
			Code: http.StatusOK,
			Msg:  "token not found",
		}
	}
	if err != nil {
		return nil, nil, &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get pool info",
		}
	}

	return tokenDetail, poolDetail, nil
}

func (s *SwapServiceImpl) processAntiMev(req request.SwapRouteRequest) (bool, string, string, *response.Response) {
	if !req.IsAntiMev {
		return false, "0", "", nil
	}

	tipFloorResponse, err := httpUtil.GetTipFloor(req.TokenOutAddress)
	if tipFloorResponse.Code != 2000 && err == nil {
		return false, "", "", &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get tip floor",
		}
	}

	tipFloor := tipFloorResponse.Data[0].EmaLandedTips50thPercentile * math.Pow(10, float64(response.SolDecimals))
	tipFloor = math.Floor(tipFloor)
	jitotip := strconv.FormatFloat(tipFloor, 'f', -1, 64)
	jitoOrderId := httpUtil.GenerateJitoOrderId(req.TokenOutAddress, req.TokenInAddress, req.InAmount.String(), jitotip)

	return true, jitotip, jitoOrderId, nil
}

func (s *SwapServiceImpl) buildGameFunGInstructionStruct(req request.SwapRouteRequest, inAmount string) httpRequest.SwapGInstructionStruct {
	return httpRequest.SwapGInstructionStruct{
		User: req.FromAddress,
		// User:        "GXL1pXLNKzFq7rzbFsGor6NaMsSMjoKhLqmxe8vsh7Gg",
		InputAmount: inAmount,
		InputMint:   req.TokenInAddress,
		// OutputMint:  req.TokenOutAddress,
		OutputMint:  "ZziTphJ4pYsbWZtpR8TaHy2xDqbNyf8yEp5d5jvpump",
		SlippageBps: req.Slippage,
		GMint:       "ZziTphJ4pYsbWZtpR8TaHy2xDqbNyf8yEp5d5jvpump",
		Amm:         "4ZaJqcDxgCCMpBL6TiAz6A8H8zQ6imas4eMs3Hk4ra52",
		Market:      "75dsjBhyyMsbEoqhgQGCdunYDrgmmPSmBDinnzqVL9Hv",
		GAmm:        "4ZaJqcDxgCCMpBL6TiAz6A8H8zQ6imas4eMs3Hk4ra52",
		GMarket:     "75dsjBhyyMsbEoqhgQGCdunYDrgmmPSmBDinnzqVL9Hv",
	}
}

func (s *SwapServiceImpl) buildBuyGWithPointsStruct(req request.SwapRouteRequest, points string) httpRequest.BuyGWithPointsStruct {
	return httpRequest.BuyGWithPointsStruct{
		User: req.FromAddress,
		// User:      "GXL1pXLNKzFq7rzbFsGor6NaMsSMjoKhLqmxe8vsh7Gg",
		Points:    points,
		InputMint: req.TokenInAddress,
		// OutputMint:  req.TokenOutAddress,
		OutputMint:  "ZziTphJ4pYsbWZtpR8TaHy2xDqbNyf8yEp5d5jvpump",
		SlippageBps: 0,
		GMint:       "ZziTphJ4pYsbWZtpR8TaHy2xDqbNyf8yEp5d5jvpump",
		Amm:         "4ZaJqcDxgCCMpBL6TiAz6A8H8zQ6imas4eMs3Hk4ra52",
		Market:      "75dsjBhyyMsbEoqhgQGCdunYDrgmmPSmBDinnzqVL9Hv",
		GAmm:        "4ZaJqcDxgCCMpBL6TiAz6A8H8zQ6imas4eMs3Hk4ra52",
		GMarket:     "75dsjBhyyMsbEoqhgQGCdunYDrgmmPSmBDinnzqVL9Hv",
	}

}

func (s *SwapServiceImpl) buildSwapPumpStruct(req request.SwapRouteRequest, tokenDetail *model.TokenInfo, poolDetail *model.TokenLiquidityPool, mev bool, jitotip string) httpRequest.SwapPumpStruct {
	return httpRequest.SwapPumpStruct{
		FromAddress:                 req.FromAddress,
		InAmount:                    req.InAmount,
		InputMint:                   req.TokenInAddress,
		OutputMint:                  req.TokenOutAddress,
		SlippageBps:                 strconv.Itoa(req.Slippage),
		PriorityFee:                 req.PriorityFee,
		TokenTotalSupply:            strconv.FormatUint(tokenDetail.CirculatingSupply, 10),
		VirtualSolReserves:          strconv.FormatUint(poolDetail.PoolPcReserve, 10),
		VirtualTokenReserves:        strconv.FormatUint(poolDetail.PoolCoinReserve, 10),
		InitialRealTokenReserves:    model.PUMP_INITIAL_REAL_TOKEN_RESERVES,
		InitialVirtualSolReserves:   model.PUMP_INITIAL_VIRTUAL_SOL_RESERVES,
		InitialVirtualTokenReserves: model.PUMP_INITIAL_VIRTUAL_TOKEN_RESERVES,
		Mev:                         mev,
		Jitotip:                     jitotip,
	}
}

func (s *SwapServiceImpl) buildSwapRaydiumStruct(req request.SwapRouteRequest, tokenDetail *model.TokenInfo, poolDetail *model.TokenLiquidityPool, mev bool, jitotip string) httpRequest.SwapRaydiumStruct {
	return httpRequest.SwapRaydiumStruct{}
}

func (s *SwapServiceImpl) getRaydiumTradeTx(swapStruct httpRequest.SwapRaydiumStruct) (*httpRespone.SwapTransactionResponse, error) {

	resp, err := httpUtil.GetRaydiumTradeTx(swapStruct)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var swapTxResponse *httpRespone.SwapTransactionResponse
	if err := json.Unmarshal(respBody, &swapTxResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return swapTxResponse, nil
}

func (s *SwapServiceImpl) getPumpFunTradeTx(swapStruct httpRequest.SwapPumpStruct) (*httpRespone.SwapTransactionResponse, error) {

	resp, err := httpUtil.GetPumpFunTradeTx(swapStruct)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var swapTxResponse *httpRespone.SwapTransactionResponse
	if err := json.Unmarshal(respBody, &swapTxResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return swapTxResponse, nil
}

func (s *SwapServiceImpl) getGameFunGInstruction(swapStruct httpRequest.SwapGInstructionStruct) (*httpRespone.SwapTransactionResponse, error) {
	log.Print(swapStruct)
	resp, err := httpUtil.GetGameFunGInstruction(swapStruct)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var swapTxResponse *httpRespone.SwapTransactionResponse
	if err := json.Unmarshal(respBody, &swapTxResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return swapTxResponse, nil
}

func (s *SwapServiceImpl) getGetBuyGWithPointsInstruction(points float64, swapStruct httpRequest.BuyGWithPointsStruct) (*httpRespone.SwapTransactionResponse, error) {
	// 使用 GetBuyGWithPointsInstruction 函数发起请求
	resp, err := httpUtil.GetBuyGWithPointsInstruction(swapStruct)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析 JSON 响应体
	var swapTxResponse httpRespone.SwapTransactionResponse
	if err := json.Unmarshal(respBody, &swapTxResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// 获取 Redis 键并设置值
	SwapGPointsKey := GetRedisKey(constants.SwapGPoints, swapTxResponse.Data)
	multiplier := math.Pow10(model.PointsDecimal)
	scaledPoints := uint64(points * multiplier)

	// 将 scaledPoints 存入 Redis
	err = redis.Set(SwapGPointsKey, scaledPoints, 5*time.Minute)
	if err != nil {
		util.Log().Error("Failed to set key in Redis: %v", err)
	}

	return &swapTxResponse, nil
}

func (s *SwapServiceImpl) calculateSwapAmounts(
	req request.SwapRouteRequest,
	tokenDetail *model.TokenInfo,
	inPriceUSD decimal.Decimal,
	inDecimals, outDecimals uint8,
) (
	outAmount, inAmountUSD, outAmountUSD decimal.Decimal,
	err *response.Response,
) {
	outPriceUSD := tokenDetail.Price
	if req.SwapType == "sell" {
		inDecimals, outDecimals = outDecimals, inDecimals
		inPriceUSD, outPriceUSD = outPriceUSD, inPriceUSD
	}

	// 计算 10^inDecimals
	decimalsFactor := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(inDecimals)))

	// 计算 inAmountUSD
	inAmountUSD = req.InAmount.Mul(inPriceUSD).Div(decimalsFactor)

	if req.PlatformType == "g_points" {
		divisor := decimal.NewFromInt(2)
		inAmountUSD = inAmountUSD.Div(divisor)
	}
	outAmount = inAmountUSD.Div(outPriceUSD)
	outAmountUSD = outAmount.Mul(outPriceUSD)

	return outAmount, inAmountUSD, outAmountUSD, nil
}

func ConstructSwapRouteResponse(req request.SwapRouteRequest, swapResponse *httpRespone.SwapTransactionResponse, inDecimals, outDecimals uint8, amountOut, amountInUSD, amountOutUSD decimal.Decimal, startTime time.Time, jitoOrderId string) response.Response {

	swapRouteResponse := response.Response{
		Code: http.StatusOK,
		Msg:  swapResponse.Message,
		Data: response.SwapRouteData{
			Quote: response.Quote{
				InputMint:            req.TokenInAddress,
				InAmount:             req.InAmount,
				InDecimals:           inDecimals,
				OutDecimals:          outDecimals,
				OutputMint:           req.TokenOutAddress,
				OutAmount:            amountOut,
				OtherAmountThreshold: decimal.NewFromInt(0).String(),
				SlippageBps:          strconv.Itoa(req.Slippage),
				PlatformFee:          0,
				RoutePlan: []response.RoutePlan{
					{
						SwapInfo: response.SwapInfo{
							AmmKey:     "Pump",
							Label:      "Pump",
							InputMint:  req.TokenInAddress,
							OutputMint: req.TokenOutAddress,
							InAmount:   req.InAmount,
							OutAmount:  amountOut,
							FeeAmount:  req.PriorityFee,
							FeeMint:    "So11111111111111111111111111111111111111112",
						},
						Percent: 0,
					},
				},
				TimeTaken: time.Since(startTime).Seconds(),
			},
			RawTx: response.RawTx{
				SwapTransaction: swapResponse.Data,
				// LastValidBlockHeight: swapResponse.Data.,
				// RecentBlockhash: swapResponse.Data.SwapTransaction.Message.RecentBlockhash,
			},
			PlatformType: req.PlatformType,
			AmountInUSD:  amountInUSD,
			AmountOutUSD: amountOutUSD,
			JitoOrderID:  jitoOrderId,
		},
	}
	return swapRouteResponse
}

func (s *SwapServiceImpl) SendTransaction(userID string, swapTransaction string, isJito bool, platformType string) response.Response {
	isUsePoint := false
	log.Print(swapTransaction)
	if platformType == "g_points" {
		SwapGPointsKey := GetRedisKey(constants.SwapGPoints, swapTransaction)
		redisValue, err := redis.Get(SwapGPointsKey)
		if err != nil {
			util.Log().Error("Failed to get key from Redis: %v", err)
			return response.Err(http.StatusInternalServerError, "Failed to retrieve points from Redis", err)
		}
		if redisValue == "" {
			util.Log().Error("Key not found in Redis: %s", SwapGPointsKey)
			return response.Err(http.StatusNotFound, "Transaction expired, please initiate the transaction again!", nil)
		}
		points, err := strconv.ParseUint(redisValue, 10, 64)
		if err != nil {
			util.Log().Error("Failed to convert Redis value to uint64: %v", err)
			return response.Err(http.StatusInternalServerError, "Failed to convert Redis points value to uint64", err)
		}
		userIDUint64, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			util.Log().Error("Failed to convert userID to uint64: %v", err)
			return response.Err(http.StatusInternalServerError, "Invalid user ID", err)
		}
		userInfo, err := s.userInfoRepo.GetUserByUserID(uint(userIDUint64))
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Unable to retrieve user information, transaction failed!", err)
		}
		if points > userInfo.AvailablePoints {
			return response.Err(http.StatusInternalServerError, "Your available points are insufficient, transaction failed!", err)
		}
		isTrue, err := s.userInfoRepo.DeductPointsWithOptimisticLock(userIDUint64, points)
		if err != nil {
			util.Log().Error("Failed to deduct points with optimistic lock: %v", err)
			return response.Err(http.StatusInternalServerError, "Failed to deduct points, please try again later", err)
		}
		if !isTrue {
			util.Log().Error("Optimistic lock failed, points deduction unsuccessful for user: %d", userIDUint64)
			return response.Err(http.StatusConflict, "Points deduction failed due to concurrent update, please try again", nil)
		}
		isUsePoint = true
	}

	resp, err := httpUtil.SendGameFunTransaction(swapTransaction, isJito, isUsePoint)
	if err != nil || resp == nil || resp.Code != 2000 {
		return response.Err(http.StatusInternalServerError, "Failed to get send transaction", err)
	}
	return response.Success(resp.Data)
}

func (s *SwapServiceImpl) GetSwapStatusBySignature(swapTransaction string) response.Response {

	resp, err := httpUtil.GetSwapStatusBySignature(swapTransaction)
	if err != nil || resp == nil || resp.Code != 2000 {
		return response.Err(http.StatusInternalServerError, "Failed to get swap request status", err)
	}

	status := 0

	if resp.Data == "success" {
		status = 1
	} else if resp.Data == "failed" {
		status = 2
	} else if resp.Data == "processing" {
		status = 0
	}

	return response.Success(
		map[string]interface{}{
			"status": status,
		},
	)

}
