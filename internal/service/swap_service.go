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
			swapStruct := s.buildGameFunGInstructionStruct(req)
			return s.getGameFunGInstruction(swapStruct)
		},
		"g_points": func() (*httpRespone.SwapTransactionResponse, error) {
			swapStruct := s.buildBuyGWithPointsStruct(req)
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

	// Calculate amounts
	platform := model.CreatedPlatformType(tokenDetail.CreatedPlatformType)
	inDecimals := model.SOL_DECIMALS
	outDecimals := platform.GetDecimals()

	outAmount, inAmountUSD, outAmountUSD, errResp := s.calculateSwapAmounts(req, tokenDetail, inDecimals, outDecimals)
	if errResp != nil {
		return s.handleErrorResponse(errResp)
	}

	// Return the constructed response
	return ConstructSwapRouteResponse(req, swapTransaction, inDecimals, outDecimals, outAmount, inAmountUSD, outAmountUSD, startTime, jitoOrderId)
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

func (s *SwapServiceImpl) buildGameFunGInstructionStruct(req request.SwapRouteRequest) httpRequest.SwapGInstructionStruct {
	return httpRequest.SwapGInstructionStruct{
		User:         req.FromAddress,
		InputAmount:  req.InAmount,
		InputMint:    req.TokenInAddress,
		OutputMint:   req.TokenOutAddress,
		SlippageBps:  req.Slippage,
		GMint:        "GMintDefault",
		Amm:          "AmmDefault",
		Market:       "MarketDefault",
		GAmm:         "GAmmDefault",
		GMarket:      "GMarketDefault",
		FeeRecipient: "FeeRecipientDefault",
	}
}

func (s *SwapServiceImpl) buildBuyGWithPointsStruct(req request.SwapRouteRequest) httpRequest.BuyGWithPointsStruct {
	return httpRequest.BuyGWithPointsStruct{
		User:         req.FromAddress,
		InputAmount:  req.InAmount,
		InputMint:    req.TokenInAddress,
		OutputMint:   req.TokenOutAddress,
		SlippageBps:  req.Slippage,
		GMint:        "GMintDefault",
		Amm:          "AmmDefault",
		Market:       "MarketDefault",
		GAmm:         "GAmmDefault",
		GMarket:      "GMarketDefault",
		FeeRecipient: "FeeRecipientDefault",
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

func (s *SwapServiceImpl) getGetBuyGWithPointsInstruction(points float64, swapStruct httpRequest.BuyGWithPointsStruct) (*httpRespone.SwapTransactionResponse, error) {

	resp, err := httpUtil.GetBuyGWithPointsInstruction(swapStruct)
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

	SwapGPointsKey := GetRedisKey(constants.SwapGPoints, swapTxResponse.Data.Base64SwapTransaction)
	multiplier := math.Pow10(model.PointsDecimal)
	scaledPoints := uint64(points * multiplier)
	err = redis.Set(SwapGPointsKey, scaledPoints, 5*time.Minute)
	if err != nil {
		util.Log().Error("Failed to set key in Redis: %v", err)
	}
	return swapTxResponse, nil
}

func (s *SwapServiceImpl) calculateSwapAmounts(
	req request.SwapRouteRequest,
	tokenDetail *model.TokenInfo,
	inDecimals, outDecimals uint8,
) (
	outAmount, inAmountUSD, outAmountUSD decimal.Decimal,
	err *response.Response,
) {

	inPriceUSD, priceErr := getSolPrice()

	if priceErr != nil {
		return decimal.Decimal{}, decimal.Decimal{}, decimal.Decimal{}, &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get SOL price",
		}
	}

	outPriceUSD := tokenDetail.Price

	if req.SwapType == "sell" {
		inDecimals, outDecimals = outDecimals, inDecimals
		inPriceUSD, outPriceUSD = outPriceUSD, inPriceUSD
	}

	multiplier := decimal.NewFromFloat(math.Pow(10, float64(inDecimals)))
	inAmount := req.InAmount.Mul(multiplier)
	inAmountUSD = inAmount.Mul(inPriceUSD)
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
				SwapTransaction:      swapResponse.Data.Base64SwapTransaction,
				LastValidBlockHeight: swapResponse.Data.LastValidBlockHeight,
				RecentBlockhash:      swapResponse.Data.SwapTransaction.Message.RecentBlockhash,
			},
			AmountInUSD:  amountInUSD,
			AmountOutUSD: amountOutUSD,
			JitoOrderID:  jitoOrderId,
		},
	}
	return swapRouteResponse
}

func (s *SwapServiceImpl) SendTransaction(userID string, swapTransaction string, isJito bool, platformType string) response.Response {
	isUsePoint := false
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
		userInfo, err := s.userInfoRepo.GetUserByUserID(userIDUint64)
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
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to send swap transaction", err)
	}
	defer resp.Body.Close() // 读取完毕后关闭，避免泄漏

	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to read response body", err)
	}

	return response.Success(string(readAll))
}

func (s *SwapServiceImpl) GetSwapStatusBySignature(swapTransaction string) response.Response {
	resp, err := httpUtil.GetSwapStatusBySignature(swapTransaction)
	if err != nil || resp == nil {
		return response.Err(http.StatusInternalServerError, "Failed to get swap request status", err)
	}

	defer resp.Body.Close()

	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "failed to read response body", err)
	}
	return response.Success(string(readAll))
}
