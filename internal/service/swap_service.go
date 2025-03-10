package service

import (
	"encoding/json"
	"fmt"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpRequest"
	"game-fun-be/internal/pkg/httpRespone"

	"errors"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type SwapServiceImpl struct {
}

func NewSwapService() *SwapServiceImpl {
	return &SwapServiceImpl{}
}

func (s *SwapServiceImpl) GetSwapRoute(req request.SwapRouteRequest, chainType uint8) response.Response {

	startTime := time.Now()

	tokenDetail, poolDetail, errResp := s.getTokenAndPoolInfo(req.TokenAddress, chainType)
	if errResp != nil {
		return response.Err(errResp.Code, errResp.Msg, errors.New(errResp.Error))
	}

	mev, jitotip, jitoOrderId, errResp := s.processAntiMev(req)
	if errResp != nil {
		return response.Err(errResp.Code, errResp.Msg, errors.New(errResp.Error))
	}

	var swapTransaction *httpRespone.SwapTransactionResponse

	if req.PlatformType == "pump" {
		swapStruct := s.buildSwapPumpStruct(req, tokenDetail, poolDetail, mev, jitotip)
		swapTransactionResponse, err := s.sendSwapRequest(swapStruct)
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to send swap request", err)
		}
		swapTransaction = swapTransactionResponse
	}
	if req.PlatformType == "raydium" {
		swapStruct := s.buildSwapPumpStruct(req, tokenDetail, poolDetail, mev, jitotip)
		swapTransactionResponse, err := s.sendSwapRequest(swapStruct)
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to send swap request", err)
		}
		swapTransaction = swapTransactionResponse
	}

	platform := model.CreatedPlatformType(tokenDetail.CreatedPlatformType)
	inDecimals := model.SOL_DECIMALS
	outDecimals := platform.GetDecimals()

	outAmount, inAmountUSD, outAmountUSD, errResp := s.calculateSwapAmounts(req, tokenDetail, inDecimals, outDecimals)
	if errResp != nil {
		return response.Err(errResp.Code, errResp.Msg, errors.New(errResp.Error))
	}

	return ConstructSwapRouteResponse(req, swapTransaction, inDecimals, outDecimals, outAmount, inAmountUSD, outAmountUSD, startTime, jitoOrderId)
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

func (s *SwapServiceImpl) buildSwapPumpStruct(req request.SwapRouteRequest, tokenDetail *model.TokenInfo, poolDetail *model.TokenLiquidityPool, mev bool, jitotip string) httpRequest.SwapPumpStruct {
	return httpRequest.SwapPumpStruct{
		FromAddress:                 req.FromAddress,
		InAmount:                    req.InAmount,
		InputMint:                   req.TokenInAddress,
		OutputMint:                  req.TokenOutAddress,
		SlippageBps:                 req.Slippage,
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

func (s *SwapServiceImpl) sendSwapRequest(swapStruct httpRequest.SwapPumpStruct) (*httpRespone.SwapTransactionResponse, error) {

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
				SlippageBps:          req.Slippage,
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

func (s *SwapServiceImpl) SendTransaction(swapTransaction string, isJito bool) response.Response {
	resp, err := httpUtil.SendTransaction(swapTransaction, isJito)
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
