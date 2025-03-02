package service

import (
	"math"
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/response"
	"net/http"
)

func GetMarketInfo(tokenAddress string, tokenType model.ChainType) response.Response {

	poolAddress := ""
	baseAddress := ""
	quoteAddress := ""
	baseSymbol := ""
	creator := ""
	creationTimestamp := int64(0)

	tokenInfo, err := model.GetTokenInfoByAddress(tokenAddress, uint8(tokenType))
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "Failed to get token info", err)
	}
	if tokenInfo == nil {
		return response.BuildResponse(nil, http.StatusOK, "Token info not found", nil)
	}
	baseSymbol = tokenInfo.Symbol

	platformType := uint8(1)
	if tokenInfo.IsComplete {
		platformType = 2
	}

	poolInfo, err := QueryAndCheckPool(tokenInfo.TokenAddress, uint8(tokenType), platformType)
	if err != nil {
		util.Log().Error("Failed to get pool info", "pool_address", tokenInfo.PoolAddress, "error", err)
		return response.BuildResponse(nil, http.StatusInternalServerError, "Failed to get pool info", err)
	} else if poolInfo == nil {
		return response.BuildResponse(nil, http.StatusOK, "Pool info not found", nil)
	} else {
		poolAddress = poolInfo.PoolAddress
		quoteAddress = poolInfo.CoinAddress
		baseAddress = poolInfo.PcAddress
		creator = poolInfo.UserAddress
		creationTimestamp = poolInfo.BlockTime.Unix()
	}

	solPrice, err := getSolPrice()
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "failed to get sol price", err)
	}

	solPriceFloat, _ := solPrice.Float64()
	tokenPriceFloat, _ := tokenInfo.Price.Float64()

	baseReserveValue := float64(poolInfo.PoolPcReserve) / math.Pow(10, response.SolDecimals) * solPriceFloat
	quoteReserveValue := float64(poolInfo.PoolCoinReserve) * tokenPriceFloat / math.Pow(10, float64(tokenInfo.Decimals))
	baseReserve := float64(poolInfo.PoolPcReserve) / math.Pow(10, response.SolDecimals)
	quoteReserve := float64(poolInfo.PoolCoinReserve) / math.Pow(10, float64(tokenInfo.Decimals))
	totalSupply := float64(tokenInfo.TotalSupply) / math.Pow(10, float64(tokenInfo.Decimals))
	circulatingSupply := float64(tokenInfo.CirculatingSupply) / math.Pow(10, float64(tokenInfo.Decimals))

	initialBaseReserve := float64(poolInfo.InitialPcReserve) / math.Pow(10, response.SolDecimals)
	initialQuoteReserve := float64(poolInfo.InitialCoinReserve) / math.Pow(10, float64(tokenInfo.Decimals))

	if !tokenInfo.IsComplete {
		baseReserveValue = float64(poolInfo.RealNativeReserves) / math.Pow(10, response.SolDecimals) * solPriceFloat
		// quoteReserveValue = float64(poolInfo.RealTokenReserves) * tokenPriceFloat / math.Pow(10, float64(tokenInfo.Decimals))
		baseReserve = float64(poolInfo.RealNativeReserves) / math.Pow(10, response.SolDecimals)
		// quoteReserve = float64(poolInfo.RealTokenReserves) / math.Pow(10, float64(tokenInfo.Decimals))
	}

	marketInfo := response.MarketInfo{
		Address:             tokenAddress,
		MarketId:            poolInfo.MarketAddress,
		PoolAddress:         poolAddress,
		QuoteAddress:        quoteAddress,
		QuoteSymbol:         baseSymbol,
		BaseSymbol:          "SOL",
		BaseAddress:         baseAddress,
		MarketCap:           totalSupply * tokenPriceFloat,
		TotalSupply:         totalSupply,
		CirculatingSupply:   circulatingSupply,
		Liquidity:           baseReserve * 2 * solPriceFloat,
		BaseReserve:         baseReserve,
		QuoteReserve:        quoteReserve,
		BaseReserveValue:    baseReserveValue,
		QuoteReserveValue:   quoteReserveValue,
		QuoteVaultAddress:   tokenAddress,
		BaseVaultAddress:    tokenAddress,
		Creator:             creator,
		CreationTimestamp:   creationTimestamp,
		Progress:            tokenInfo.Progress,
		InitialBaseReserve:  initialBaseReserve,
		InitialQuoteReserve: initialQuoteReserve,
	}

	return response.BuildResponse(marketInfo, http.StatusOK, "Success", nil)
}
