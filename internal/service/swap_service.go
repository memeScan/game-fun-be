package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"game-fun-be/internal/constants"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpRequest"
	"game-fun-be/internal/pkg/httpRespone"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"

	"github.com/IBM/sarama"
	"github.com/shopspring/decimal"
)

type SwapServiceImpl struct {
	userInfoRepo      *model.UserInfoRepo
	pointsServiceImpl *PointsServiceImpl
	kafka             sarama.SyncProducer
}

func NewSwapService(producer sarama.SyncProducer) *SwapServiceImpl {
	return &SwapServiceImpl{
		kafka: producer,
	}
}

func (s *SwapServiceImpl) GetSwapRoute(req request.SwapRouteRequest, chainType uint8) response.Response {
	startTime := time.Now().UnixNano()

	solPriceUSD, priceErr := getSolPrice()
	if priceErr != nil {
		return response.Err(http.StatusBadRequest, "price query failed", priceErr)
	}

	tokenDetail, poolDetail, errResp := s.getTokenAndPoolInfo(req.TokenAddress, chainType, req.PlatformType)
	if errResp != nil {
		return s.handleErrorResponse(errResp)
	}

	if req.SwapType == "buy" && chainType == 1 {
		solMultiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(model.SOL_DECIMALS)))
		req.InAmount = req.InAmount.Mul(solMultiplier)
	} else if req.SwapType == "sell" && chainType == 1 {
		tokenMultiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(tokenDetail.Decimals)))
		req.InAmount = req.InAmount.Mul(tokenMultiplier)
	}
	inAmountUint64 := req.InAmount.BigInt().Uint64()
	inAmountStr := strconv.FormatUint(inAmountUint64, 10)

	// Process Anti-MEV logic
	mev, jitotip, jitoOrderId, errResp := s.processAntiMev(req)
	if errResp != nil {
		return s.handleErrorResponse(errResp)
	}

	isCanBuyflag := true
	// pointsString := ""
	estimatedPoints := decimal.NewFromFloat(0.0)

	if req.PlatformType == "g_external" {
		isCanBuyflag = CheckBalanceSufficient(chainType, req.SwapType, req.PlatformType, req.InAmount, req.UserBalance, model.SOL_DECIMALS, req.PriorityFee)
	}
	// else if req.PlatformType == "g_points" {
	// 	pointsDecimal := decimal.NewFromFloat(req.Points) // 用户输入的 Points（代币数量）
	// 	// 计算 10^decimals，确保计算精度
	// 	pointsDecimalsFactor := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(tokenDetail.Decimals)))
	// 	// 去掉小数点部分，取整
	// 	pointsWithoutDecimal := pointsDecimal.Mul(pointsDecimalsFactor).Floor()
	// 	// 将去掉小数点后的结果转换为字符串
	// 	// pointsString = pointsWithoutDecimal.String()
	// 	// 计算代币的 USD 价值（代币数量 * 代币单价 * 10^精度）
	// 	tokenUsd := pointsDecimal.Mul(tokenDetail.Price)
	// 	// 计算需要多少 SOL（代币 USD 价值 / SOL 的 USD 价值）
	// 	req.InAmount = tokenUsd.Div(solPriceUSD)
	// 	req.InAmount = req.InAmount.Div(decimal.NewFromInt(2))
	// 	isCanBuyflag = CheckBalanceSufficient(chainType, req.SwapType, req.PlatformType, req.InAmount, req.UserBalance, model.SOL_DECIMALS, req.PriorityFee)
	// }

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
			swapStruct := s.buildGameFunGInstructionStruct(req, poolDetail, inAmountStr, mev, jitotip)
			if isCanBuyflag {
				return s.getGameFunGInstruction(swapStruct)
			}
			return &httpRespone.SwapTransactionResponse{
				Code:    2000,
				Message: "Your custom message here",
			}, nil
		},
		"g_points": func() (*httpRespone.SwapTransactionResponse, error) {
			// swapStruct := s.buildBuyGWithPointsStruct(req, pointsString)
			// if isCanBuyflag {
			// 	return s.getGetBuyGWithPointsInstruction(req.Points, swapStruct, startTime)
			// }
			return &httpRespone.SwapTransactionResponse{
				Code:    2000,
				Message: "Points rule is upgrading!",
			}, nil
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

	inDecimals := model.SOL_DECIMALS
	outDecimals := tokenDetail.Decimals
	outAmount, inAmountUSD, outAmountUSD, errResp := s.calculateSwapAmounts(req, tokenDetail, solPriceUSD, inDecimals, outDecimals)
	if errResp != nil {
		return s.handleErrorResponse(errResp)
	}

	// if req.PlatformType == "g_points" {
	// 	multiplier := decimal.NewFromFloat(math.Pow(10, float64(outDecimals)))
	// 	result := outAmount.Mul(multiplier)
	// 	uint64Result := result.IntPart()
	// 	if uint64Result < 0 {
	// 		uint64Result = 0
	// 	}
	// 	estimatedPointsUint, err := s.getVaultAndQuotaUsage(uint64(uint64Result))
	// 	if err != nil {
	// 	}
	// 	estimatedPointsRsp := divideByDecimals(estimatedPointsUint, outDecimals)
	// 	estimatedPoints = estimatedPointsRsp
	// }

	// Return the constructed response
	return ConstructSwapRouteResponse(req, swapTransaction, uint8(inDecimals), uint8(outDecimals), outAmount, inAmountUSD, outAmountUSD, startTime, jitoOrderId, estimatedPoints)
}

// 添加余额检测接口
func CheckBalanceSufficient(chainType uint8, swaType string, PlatformType string, tokneInAmount decimal.Decimal, tokneBalanceAmount decimal.Decimal, nativeTokenDecimals uint8, priorityFee float64) bool {
	needPayAmount := decimal.NewFromInt(0)

	priorityFeeDecimal := decimal.NewFromFloat(priorityFee)
	powerOfTen := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(nativeTokenDecimals)))
	scaledPriorityFee := priorityFeeDecimal.Mul(powerOfTen)
	if PlatformType == "g_external" {
		if swaType == "buy" {
			multiplier := decimal.NewFromFloat(0.005)
			powerOfTen := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(nativeTokenDecimals)))
			needPayAmount = multiplier.Mul(powerOfTen).Add(scaledPriorityFee).Add(tokneInAmount)
		} else if swaType == "sell" {
			multiplier := decimal.NewFromFloat(0.003)
			nativeTokenDecimalsDecimal := decimal.NewFromInt(int64(nativeTokenDecimals))
			needPayAmount = multiplier.Mul(nativeTokenDecimalsDecimal).Add(scaledPriorityFee)
		}
	} else if PlatformType == "g_points" {
		solMultiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(nativeTokenDecimals)))
		result := tokneInAmount.Mul(solMultiplier)
		multiplier := decimal.NewFromFloat(0.0025)
		powerOfTen := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(nativeTokenDecimals)))
		needPayAmount = multiplier.Mul(powerOfTen).Add(scaledPriorityFee).Add(result)
	}

	if tokneBalanceAmount.LessThan(needPayAmount) {
		return false
	}
	return true
}

// // GetVaultAndQuotaUsage 获取 Vault 余额和最近 10 分钟的额度使用量
// func (s *SwapServiceImpl) getVaultAndQuotaUsage(quotaTotalAmount uint64) (uint64, error) {
// 	var vaultBalance uint64        // 金库余额
// 	var quotaUsageLast10Min uint64 // 过去 10 分钟的额度使用量

// 	// 获取 Redis 中的 Vault 余额
// 	vaultStr, err := redis.Get(constants.RedisKeyVaultAmount)
// 	if err != nil {
// 		util.Log().Error("Failed to get vault balance from Redis:", err)
// 	} else {
// 		parsedValue, err := strconv.ParseUint(vaultStr, 10, 64)
// 		if err != nil {
// 			util.Log().Error("Failed to parse vault balance:", err)
// 		} else {
// 			vaultBalance = parsedValue
// 		}
// 	}

// 	// 获取 Redis 中最近 10 分钟的额度使用量
// 	quotaStr, err := redis.Get(constants.RedisKeyQuotaAmountLast10Min)
// 	if err != nil {
// 		util.Log().Error("Failed to get quota usage from Redis:", err)
// 	} else {
// 		parsedValue, err := strconv.ParseUint(quotaStr, 10, 64)
// 		if err != nil {
// 			util.Log().Error("Failed to parse quota usage:", err)
// 		} else {
// 			quotaUsageLast10Min = parsedValue
// 		}
// 	}

// 	// 计算积分
// 	point, _, err := s.pointsServiceImpl.CalculatePoint(vaultBalance, quotaUsageLast10Min, quotaTotalAmount)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return point, nil
// }

// func divideByDecimals(estimatedPointsUint uint64, outDecimals uint8) decimal.Decimal {
// 	estimatedPointsDecimal := decimal.NewFromInt(int64(estimatedPointsUint))
// 	divisor := decimal.New(10, 0).Pow(decimal.NewFromInt(int64(outDecimals)))
// 	result := estimatedPointsDecimal.Div(divisor)
// 	return result
// }

// Helper function to handle error responses
func (s *SwapServiceImpl) handleErrorResponse(errResp *response.Response) response.Response {
	return response.Err(errResp.Code, errResp.Msg, errors.New(errResp.Error))
}

func (s *SwapServiceImpl) getTokenAndPoolInfo(tokenAddress string, chainType uint8, platformType string) (*model.TokenInfo, *model.TokenLiquidityPool, *response.Response) {
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

	poolDetail, err := QueryAndCheckPool(tokenAddress, chainType, 2)
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

	var tipFloor float64 = 10000 // 默认值

	if len(tipFloorResponse.Data) > 0 {
		tipFloor = tipFloorResponse.Data[0].EmaLandedTips50thPercentile * math.Pow(10, float64(response.SolDecimals))
		tipFloor = math.Floor(tipFloor)
	}
	jitotip := strconv.FormatFloat(tipFloor, 'f', -1, 64)
	jitoOrderId := httpUtil.GenerateJitoOrderId(req.TokenOutAddress, req.TokenInAddress, req.InAmount.String(), jitotip)

	return true, jitotip, jitoOrderId, nil
}

func (s *SwapServiceImpl) buildGameFunGInstructionStruct(req request.SwapRouteRequest, poolDetail *model.TokenLiquidityPool, inAmount string, mev bool, jitotip string) httpRequest.SwapGInstructionStruct {
	return httpRequest.SwapGInstructionStruct{
		User:            req.FromAddress,
		InputAmount:     inAmount,
		InputMint:       req.TokenInAddress,
		OutputMint:      req.TokenOutAddress,
		SlippageBps:     req.Slippage,
		PoolPcAddress:   poolDetail.PcAddress,
		PoolCoinAddress: poolDetail.CoinAddress,
		PoolPcReserve:   strconv.FormatUint(poolDetail.PoolPcReserve, 10),
		PoolCoinReserve: strconv.FormatUint(poolDetail.PoolCoinReserve, 10),
		GMint:           os.Getenv("GMINT_ADDRESS"),
		Amm:             os.Getenv("AMM_ADDRESS"),
		Market:          os.Getenv("MARKET_ADDRESS"),
		GAmm:            os.Getenv("GAMM_ADDRESS"),
		GMarket:         os.Getenv("GMARKET_ADDRESS"),
		Mev:             mev,
		Jitotip:         jitotip,
	}
}

func (s *SwapServiceImpl) buildBuyGWithPointsStruct(req request.SwapRouteRequest, points string) httpRequest.BuyGWithPointsStruct {
	return httpRequest.BuyGWithPointsStruct{
		User:        req.FromAddress,
		Points:      points,
		InputMint:   req.TokenInAddress,
		OutputMint:  req.TokenOutAddress,
		SlippageBps: 0,
		GMint:       os.Getenv("GMINT_ADDRESS"),
		Amm:         os.Getenv("AMM_ADDRESS"),
		Market:      os.Getenv("MARKET_ADDRESS"),
		GAmm:        os.Getenv("GAMM_ADDRESS"),
		GMarket:     os.Getenv("GMARKET_ADDRESS"),
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

func (s *SwapServiceImpl) getGetBuyGWithPointsInstruction(points float64, swapStruct httpRequest.BuyGWithPointsStruct, startTime int64) (*httpRespone.SwapTransactionResponse, error) {
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
	startTimeStr := strconv.FormatInt(startTime, 10)
	SwapGPointsKey := GetRedisKey(constants.SwapGPoints, swapStruct.User, startTimeStr)
	multiplier := math.Pow10(model.PointsDecimal)
	scaledPoints := uint64(points * multiplier)

	// 将 scaledPoints 存入 Redis
	err = redis.Set(SwapGPointsKey, scaledPoints, 5*time.Minute)
	if err != nil {
		util.Log().Error("Failed to set key in Redis: %v", err)
	}

	return &swapTxResponse, nil
}

func (s *SwapServiceImpl) CheckRebate(address string, rebateAmount uint64) response.Response {
	startTime := time.Now().UnixNano()
	user, err := s.userInfoRepo.GetUserByAddress(address, model.ChainTypeSolana.Uint8())
	if err != nil {
		return response.Err(http.StatusBadRequest, "用户不存在", err)
	}

	if user.WithdrawableRebate < rebateAmount {
		return response.Err(http.StatusBadRequest, "提现金额不足", errors.New("提现金额不足"))
	}

	swapStruct := httpRequest.ClaimRebateStruct{
		User:   address,
		Amount: strconv.FormatUint(rebateAmount, 10),
	}
	resp, err := s.getClaimRebateInstruction(swapStruct, startTime)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get claim rebate instruction", err)
	}

	return response.Success(resp.Data)
}

func (s *SwapServiceImpl) getClaimRebateInstruction(swapStruct httpRequest.ClaimRebateStruct, startTime int64) (*httpRespone.SwapTransactionResponse, error) {
	// 使用 GetBuyGWithPointsInstruction 函数发起请求
	resp, err := httpUtil.GetClaimRebateInstruction(swapStruct)
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
	// startTimeStr := strconv.FormatInt(startTime, 10)
	// SwapGPointsKey := GetRedisKey(constants.SwapGPoints, swapStruct.User, startTimeStr)
	// multiplier := math.Pow10(model.PointsDecimal)
	// scaledPoints := uint64(points * multiplier)

	// // 将 scaledPoints 存入 Redis
	// err = redis.Set(SwapGPointsKey, scaledPoints, 5*time.Minute)
	// if err != nil {
	// 	util.Log().Error("Failed to set key in Redis: %v", err)
	// }

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
	outAmount = inAmountUSD.Div(outPriceUSD)
	outAmountUSD = outAmount.Mul(outPriceUSD)

	if req.PlatformType == "g_points" {
		divisor := decimal.NewFromInt(2)
		inAmountUSD = inAmountUSD.Div(divisor)
		outAmount = decimal.NewFromFloat(req.Points)
	}

	return outAmount, inAmountUSD, outAmountUSD, nil
}

func ConstructSwapRouteResponse(req request.SwapRouteRequest, swapResponse *httpRespone.SwapTransactionResponse, inDecimals, outDecimals uint8, amountOut, amountInUSD, amountOutUSD decimal.Decimal, startTime int64, jitoOrderId string, estimatedPoints decimal.Decimal) response.Response {
	swapRouteResponse := response.Response{
		Code: http.StatusOK,
		Msg:  swapResponse.Message,
		Data: response.SwapRouteData{
			Quote: response.Quote{
				InputMint:            req.TokenInAddress,
				InAmount:             req.InAmount,
				OutAmount:            amountOut,
				InDecimals:           inDecimals,
				OutDecimals:          outDecimals,
				OutputMint:           req.TokenOutAddress,
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
				TimeTaken: strconv.FormatInt(startTime, 10),
			},
			RawTx: response.RawTx{
				SwapTransaction: swapResponse.Data,
				// LastValidBlockHeight: swapResponse.Data.,
				// RecentBlockhash: swapResponse.Data.SwapTransaction.Message.RecentBlockhash,
			},
			PlatformType:    req.PlatformType,
			AmountInUSD:     amountInUSD,
			AmountOutUSD:    amountOutUSD,
			JitoOrderID:     jitoOrderId,
			EstimatedPoints: estimatedPoints,
		},
	}
	return swapRouteResponse
}

func (s *SwapServiceImpl) SendTransaction(userID string, userAddress string, swapRequest request.SwapRequest) response.Response {
	isUsePoint := false
	userIDUint64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		util.Log().Error("Failed to convert userID to uint64: %v", err)
		return response.Err(http.StatusInternalServerError, "Invalid user ID", err)
	}
	points := uint64(0)
	if swapRequest.PlatformType == "g_points" {
		SwapGPointsKey := GetRedisKey(constants.SwapGPoints, userAddress, swapRequest.StartTime)
		redisValue, err := redis.Get(SwapGPointsKey)
		if err != nil {
			util.Log().Error("Failed to get key from Redis: %v", err)
			return response.Err(http.StatusInternalServerError, "Failed to retrieve points from Redis", err)
		}
		if redisValue == "" {
			util.Log().Error("Key not found in Redis: %s", SwapGPointsKey)
			return response.Err(http.StatusNotFound, "Transaction expired, please initiate the transaction again!", nil)
		}
		points, err = strconv.ParseUint(redisValue, 10, 64)
		if err != nil {
			util.Log().Error("Failed to convert Redis value to uint64: %v", err)
			return response.Err(http.StatusInternalServerError, "Failed to convert Redis points value to uint64", err)
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

	resp, err := httpUtil.SendGameFunTransaction(swapRequest.SwapTransaction, swapRequest.IsAntiMEV, isUsePoint)

	if err != nil || resp == nil || resp.Code != 2000 {
		if swapRequest.PlatformType == "g_points" {
			// 交易发送失败，恢复用户积分
			userInfoRepo := model.NewUserInfoRepo()
			if err := userInfoRepo.IncrementAvailablePointsByUserID(uint(userIDUint64), points); err != nil {
				util.Log().Error("Failed to restore points for user %d after transaction %s failed: %v",
					userIDUint64, swapRequest.SwapTransaction, err)
				return response.Err(http.StatusInternalServerError, "Failed to restore points, please try again later", err)
			}
			util.Log().Info("Transaction %s failed, restored %d points to user %d",
				swapRequest.SwapTransaction, points, userIDUint64)
		}
		return response.Err(http.StatusInternalServerError, "Failed to get send transaction", err)
	}
	if swapRequest.PlatformType == "g_points" {
		// 交易发送成功，发送Kafka积分交易检测消息
		pointTxStatusMsg := model.PointTxStatusMessage{
			Signature: resp.Data.Signature,
			UserId:    uint(userIDUint64),
			Points:    points,
			TxType:    1,
		}

		msgBytes, err := json.Marshal(pointTxStatusMsg)
		if err != nil {
			util.Log().Error("Failed to marshal point transaction status message: %v", err)
		} else {
			// 发送消息到Kafka
			if err := s.SendMessage(constants.TopicPointTxStatus, msgBytes); err != nil {
				util.Log().Error("Failed to send point transaction status message to Kafka: %v", err)
			} else {
				util.Log().Info("Sent point transaction status check message for transaction %s, user %d, points %d",
					swapRequest.SwapTransaction, uint(userIDUint64), points)
			}
		}
	}
	return response.Success(resp.Data)
}

// SendMessage 发送消息到指定的 topic
func (s *SwapServiceImpl) SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := s.kafka.SendMessage(msg)
	return err
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

// SendClaimTransaction 发送积分提现交易
func (s *SwapServiceImpl) SendClaimTransaction(req request.RebateClaimRequest) response.Response {
	user, err := s.userInfoRepo.GetUserByAddress(req.Address, model.ChainTypeSolana.Uint8())
	if err != nil {
		return response.Err(http.StatusBadRequest, "用户不存在", err)
	}

	isTrue, err := s.userInfoRepo.DeductRebateWithOptimisticLock(uint64(user.ID), req.RebateAmount)
	if err != nil {
		util.Log().Error("Failed to deduct points with optimistic lock: %v", err)
		return response.Err(http.StatusInternalServerError, "Failed to deduct points, please try again later", err)
	}
	if !isTrue {
		util.Log().Error("Optimistic lock failed, points deduction unsuccessful for user: %d", user.ID)
		return response.Err(http.StatusConflict, "Points deduction failed due to concurrent update, please try again", nil)
	}

	resp, err := httpUtil.SendClaimTransaction(req.Address, req.SwapTransaction, req.RebateAmount)
	if err != nil || resp == nil || resp.Code != 2000 {

		// 交易发送失败，恢复用户积分
		userInfoRepo := model.NewUserInfoRepo()
		if err := userInfoRepo.IncrementWithdrawableRebateByUserID(uint(user.ID), user.WithdrawableRebate); err != nil {
			util.Log().Error("Failed to restore points for user %d after transaction %s failed: %v",
				user.ID, resp.Data.Signature, err)
			return response.Err(http.StatusInternalServerError, "Failed to restore points, please try again later", err)
		}
		util.Log().Info("Transaction %s failed, restored %d points to user %d",
			resp.Data.Signature, user.WithdrawableRebate, user.ID)
		return response.Err(http.StatusInternalServerError, "Failed to get send transaction", err)
	}
	pointTxStatusMsg := model.PointTxStatusMessage{
		Signature: resp.Data.Signature,
		UserId:    uint(user.ID),
		Points:    0,
		Rebate:    user.WithdrawableRebate,
		TxType:    2,
	}

	msgBytes, err := json.Marshal(pointTxStatusMsg)
	if err != nil {
		util.Log().Error("Failed to marshal point transaction status message: %v", err)
	} else {
		// 发送消息到Kafka
		if err := s.SendMessage(constants.TopicPointTxStatus, msgBytes); err != nil {
			util.Log().Error("Failed to send point transaction status message to Kafka: %v", err)
		} else {
			util.Log().Info("Sent point transaction status check message for transaction %s, user %d, points %d",
				resp.Data.Signature, user.ID, user.WithdrawableRebate)
		}
	}

	return response.Success(resp.Data)
}
