package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/pkg/httpRequest"
	"my-token-ai-be/internal/pkg/httpUtil"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
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

func CreateSwapPumpRequestBody(
	fromaddress string,
	tokenInAddress string,
	tokenOutAddress string,
	inAmount string,
	slippageBps string,
	tokenOutVirtualNativeReserves string,
	tokenOutVirtualTokenReserves string,
	feeLamports int64,
	feeBasisPoints int64,
	tokenOutTotalSupply string,
	feeRecipient string,
	mev bool,
	jitotip string,
) (httpRequest.SwapPumpStruct, error) {
	swapData := httpRequest.SwapPumpStruct{
		FromAddress:                 fromaddress,
		InputMint:                   tokenInAddress,
		OutputMint:                  tokenOutAddress,
		InAmount:                    inAmount,
		SlippageBps:                 slippageBps,
		VirtualSolReserves:          tokenOutVirtualTokenReserves,
		VirtualTokenReserves:        tokenOutVirtualNativeReserves,
		PriorityFee:                 feeLamports,
		FeeBasisPoints:              feeBasisPoints,
		TokenTotalSupply:            tokenOutTotalSupply,
		InitialRealTokenReserves:    model.PUMP_INITIAL_REAL_TOKEN_RESERVES,
		InitialVirtualSolReserves:   model.PUMP_INITIAL_VIRTUAL_SOL_RESERVES,
		InitialVirtualTokenReserves: model.PUMP_INITIAL_VIRTUAL_TOKEN_RESERVES,
		FeeRecipient:                feeRecipient,
		Mev:                         mev,
		Jitotip:                     jitotip,
	}
	// 变成json 字符串
	jsonData, err := json.Marshal(swapData)
	if err != nil {
		return httpRequest.SwapPumpStruct{}, fmt.Errorf("failed to marshal swap data: %w", err)
	}
	log.Println("swapData", string(jsonData))
	return swapData, nil
}

func CreateSwapRaydiumRequestBody(
	poolId string,
	marketId string,
	owner string,
	inputMint string,
	outputMint string,
	inAmount string,
	slippageBps string,
	poolPcReserve string,
	poolCoinReserve string,
	priorityFee string,
	poolPcAddress string,
	poolCoinAddress string,
	mev bool,
	jitotip string,
) (httpRequest.SwapRaydiumStruct, error) {
	swapData := httpRequest.SwapRaydiumStruct{
		PoolId:          poolId,
		MarketId:        marketId,
		Owner:           owner,
		InputMint:       inputMint,
		OutputMint:      outputMint,
		InAmount:        inAmount,
		SlippageBps:     slippageBps,
		PoolPcReserve:   poolPcReserve,
		PoolCoinReserve: poolCoinReserve,
		PriorityFee:     priorityFee,
		PoolPcAddress:   poolPcAddress,
		PoolCoinAddress: poolCoinAddress,
		Mev:             mev,
		Jitotip:         jitotip,
	}
	// 变成json 字符串
	jsonData, err := json.Marshal(swapData)
	if err != nil {
		return httpRequest.SwapRaydiumStruct{}, fmt.Errorf("failed to marshal swap data: %w", err)
	}
	log.Println("swapData", string(jsonData))
	return swapData, nil
}

const dateFormat = "20060102"

func (s *SwapServiceImpl) GetPumpSwapRoute(chainType model.ChainType, tradeType string, req request.SwapRouteRequest) *response.Response {

	startTime := time.Now()

	feeLamports := int64(req.Fee * 1000000000)
	slippageInt := req.Slippage * 100

	// 转成字符串
	feeLamportsStr := strconv.FormatInt(feeLamports, 10)
	slippageIntStr := strconv.FormatInt(slippageInt, 10)

	inDecimals := 0

	inPriceUSD := decimal.NewFromInt(1)
	outPriceUSD := decimal.NewFromInt(1)

	mev := false
	jitotip := "0"
	jitoOrderId := ""
	if req.IsAntiMev {
		mev = true
		// 调用链端接口拿小费
		tipFloorResponse, err := httpUtil.GetTipFloor(req.TokenOutAddress)
		if tipFloorResponse.Code != 2000 && err == nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to get tip floor",
			}
		}

		tipFloorResponse.Data[0].EmaLandedTips50thPercentile = tipFloorResponse.Data[0].EmaLandedTips50thPercentile * math.Pow(10, float64(response.SolDecimals))
		tipFloorResponse.Data[0].EmaLandedTips50thPercentile = math.Floor(tipFloorResponse.Data[0].EmaLandedTips50thPercentile)
		jitotip = strconv.FormatFloat(tipFloorResponse.Data[0].EmaLandedTips50thPercentile, 'f', -1, 64)

		// 用规则 生成jito order id
		jitoOrderId = httpUtil.GenerateJitoOrderId(req.TokenOutAddress, req.TokenInAddress, req.InAmount, jitotip)
	}

	var token = req.TokenOutAddress
	IsBuy := true
	if req.TokenInAddress != "So11111111111111111111111111111111111111112" {
		token = req.TokenInAddress
		IsBuy = false
	}

	swapType := 0

	tokenDetail, err := model.GetTokenInfoByAddress(token, uint8(chainType))
	if tokenDetail == nil {
		return &response.Response{
			Code: http.StatusNotFound,
			Msg:  "token not found",
		}
	}
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get token detail",
		}
	}

	platformType := uint8(1)
	if tokenDetail.IsComplete {
		platformType = 2
	}

	pool, err := QueryAndCheckPool(token, uint8(chainType), platformType)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get pool info",
		}
	}

	poolPcReserve := strconv.FormatUint(pool.PoolPcReserve, 10)
	poolCoinReserve := strconv.FormatUint(pool.PoolCoinReserve, 10)
	tokenTotalSupply := strconv.FormatUint(tokenDetail.TotalSupply, 10)

	if tokenDetail.IsComplete {
		swapType = 1
	}

	solPrice, err := getSolPrice()
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get sol price",
		}
	}

	if IsBuy {
		inDecimals = response.SolDecimals
		inPriceUSD = solPrice
		outPriceUSD = decimal.NewFromFloat(tokenDetail.Price.InexactFloat64())
	} else {
		inDecimals = int(tokenDetail.Decimals)
		inPriceUSD = decimal.NewFromFloat(tokenDetail.Price.InexactFloat64())
		outPriceUSD = solPrice
	}

	inAmount, err := strconv.ParseFloat(req.InAmount, 64)
	if err != nil {
		return &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "invalid in amount",
		}
	}

	inAmountUSD := decimal.NewFromFloat(inAmount).Mul(inPriceUSD)
	outAmount := inAmountUSD.Div(outPriceUSD)
	outAmountUSD := outAmount.Mul(outPriceUSD)

	// 乘以精度
	inAmountDecimal, err := decimal.NewFromString(req.InAmount)
	if err != nil {
		return &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "invalid in amount",
		}
	}

	req.InAmount = inAmountDecimal.
		Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(inDecimals)))).
		Truncate(0). // 去掉小数部分
		String()     // 转换成字符串

	var swapTransactionResponse *response.SwapTransactionResponse
	switch swapType {
	case 0:
		swapStruct, err := CreateSwapPumpRequestBody(
			req.FromAddress,
			req.TokenInAddress,
			req.TokenOutAddress,
			req.InAmount,
			slippageIntStr,
			poolCoinReserve,
			poolPcReserve,
			feeLamports,
			100,
			tokenTotalSupply,
			"CebN5WGQ4jvEPvsVU4EoHEpgzq1VV7AbicfhtW4xC9iM",
			mev,
			jitotip,
		)
		log.Println("swapStruct", swapStruct)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to create swap request body",
			}
		}

		resp, err := httpUtil.SendSwapRequest(swapStruct)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to send swap request",
			}
		}
		defer resp.Body.Close()

		globalServiceImpl := NewGlobalServiceImpl()
		res := globalServiceImpl.UsdPrice(model.ChainTypeSolana)
		if res.Error != "" {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to get sol price",
			}
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to read response body",
			}
		}

		swapTransactionResponse, err = ProcessSwapTransactionResponse(string(respBody), req)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to process response",
			}
		}

	case 1:
		swapStruct, err := CreateSwapRaydiumRequestBody(
			pool.PoolAddress,
			pool.MarketAddress,
			req.FromAddress,
			req.TokenInAddress,
			req.TokenOutAddress,
			req.InAmount,
			slippageIntStr,
			poolPcReserve,
			poolCoinReserve,
			feeLamportsStr,
			pool.PcAddress,
			pool.CoinAddress,
			mev,
			jitotip,
		)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to create swap request body",
			}
		}

		resp, err := httpUtil.SendRaydiumTradeRequest(swapStruct)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to send raydium trade request",
			}
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to read response body",
			}
		}

		swapTransactionResponse, err = ProcessSwapTransactionResponse(string(respBody), req)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "failed to process response",
			}
		}

	default:
		return &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "invalid swap type",
		}
	}

	return ConstructSwapRouteResponse(req, int64(feeLamports), slippageInt, swapTransactionResponse, int(inDecimals), int(tokenDetail.Decimals), outAmount.String(), inAmountUSD.String(), outAmountUSD.String(), startTime, jitoOrderId)
}

func ProcessSwapTransactionResponse(readAllStr string, req request.SwapRouteRequest) (*response.SwapTransactionResponse, error) {
	var response response.SwapTransactionResponse
	err := json.Unmarshal([]byte(readAllStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return &response, nil
}

func ConstructSwapRouteResponse(req request.SwapRouteRequest, feeLamports int64, slippageInt int64, swapTransactionResponse *response.SwapTransactionResponse, inDecimals, outDecimals int, amountOut string, amountInUSD string, amountOutUSD string, startTime time.Time, jitoOrderId string) *response.Response {
	if swapTransactionResponse == nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "swapTransactionResponse is nil",
		}
	}

	swapRouteResponse := &response.Response{
		Code: http.StatusOK,
		Msg:  swapTransactionResponse.Message,
		Data: response.SwapRouteData{
			Quote: response.Quote{
				InputMint:            req.TokenInAddress,
				InAmount:             req.InAmount,
				InDecimals:           inDecimals,
				OutDecimals:          outDecimals,
				OutputMint:           req.TokenOutAddress,
				OutAmount:            amountOut,
				OtherAmountThreshold: decimal.NewFromInt(0).String(),
				SwapMode:             "ExactIn",
				SlippageBps:          slippageInt,
				PlatformFee:          feeLamports,
				PriceImpactPct:       "0",
				RoutePlan: []response.RoutePlan{
					{
						SwapInfo: response.SwapInfo{
							AmmKey:     "Pump",
							Label:      "Pump",
							InputMint:  req.TokenInAddress,
							OutputMint: req.TokenOutAddress,
							InAmount:   req.InAmount,
							OutAmount:  amountOut,
							FeeAmount:  feeLamports,
							FeeMint:    "So11111111111111111111111111111111111111112",
						},
						Percent: 100,
					},
				},
				TimeTaken: time.Since(startTime).Seconds(),
			},
			RawTx: response.RawTx{
				SwapTransaction:           swapTransactionResponse.Data.Base64SwapTransaction,
				LastValidBlockHeight:      swapTransactionResponse.Data.LastValidBlockHeight,
				PrioritizationFeeLamports: 0,
				RecentBlockhash:           swapTransactionResponse.Data.SwapTransaction.Message.RecentBlockhash,
			},
			AmountInUSD:  amountInUSD,
			AmountOutUSD: amountOutUSD,
			JitoOrderID:  jitoOrderId,
		},
	}
	return swapRouteResponse
}

func (s *SwapServiceImpl) SendSwapRequest(swapTransaction string, isJito bool) *response.Response {
	resp, err := httpUtil.SendSwapTransaction(swapTransaction, isJito)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get swap request status",
		}
	}

	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to read response body",
		}
	}
	return &response.Response{
		Code: http.StatusOK,
		Msg:  string(readAll),
	}
}

func (s *SwapServiceImpl) GetSwapRequestStatusBySignature(SwapTransaction string) *response.Response {
	resp, err := httpUtil.GetSwapRequestStatusBySignature(SwapTransaction)
	if resp == nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get swap request status",
		}
	}
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get swap request status",
		}
	}

	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to read response body",
		}
	}

	return &response.Response{
		Code: resp.StatusCode,
		Msg:  string(readAll),
	}
}
