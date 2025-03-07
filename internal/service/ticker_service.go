package service

import (
	"fmt"
	"game-fun-be/internal/clickhouse"
	"game-fun-be/internal/constants"
	"game-fun-be/internal/es"
	"game-fun-be/internal/es/query"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"

	"encoding/json"
	"log"
	"math"
	"strconv"
	"time"

	"net/http"
)

type TickerServiceImpl struct {
	tokenInfoRepo            *model.TokenInfoRepo
	tokenMarketAnalyticsRepo *clickhouse.TokenMarketAnalyticsRepo
}

func NewTickerServiceImpl(tokenInfoRepo *model.TokenInfoRepo, tokenMarketAnalyticsRepo *clickhouse.TokenMarketAnalyticsRepo) *TickerServiceImpl {
	return &TickerServiceImpl{
		tokenInfoRepo:            tokenInfoRepo,
		tokenMarketAnalyticsRepo: tokenMarketAnalyticsRepo,
	}
}

func (s *TickerServiceImpl) Tickers(req request.TickersRequest, chainType model.ChainType) response.Response {
	TickersQuery, err := query.TickersQuery(&req)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to generate TickersQuery", err)
	}
	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, TickersQuery, es.UNIQUE_TOKENS)
	if err != nil || result == nil {
		status := http.StatusInternalServerError
		msg := "Failed to get pump rank"
		data := []response.TickersResponse{}
		if result == nil {
			status = http.StatusOK
			msg = "No data found"
			data := []response.TickersResponse{}
			return response.BuildResponse(data, status, msg, nil)
		}
		return response.BuildResponse(data, status, msg, err)
	}

	var tickersResponse response.TickersResponse

	return response.Success(tickersResponse)
}

func (s *TickerServiceImpl) TickerDetail(tokenAddress string, chainType model.ChainType) response.Response {

	tokenInfo, err := s.tokenInfoRepo.GetTokenInfoByAddress(tokenAddress, uint8(chainType))
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token info by address", err)
	}

	var tickerRespons response.GetTickerResponse

	tickerRespons.Market = response.Market{
		MarketID:  tokenInfo.ID,
		TokenMint: tokenInfo.TokenAddress,
		// 待确定
		// Market:          userInfo.MarketAddress,
		// TokenVault:      userInfo.TokenVault,
		// NativeVault:     userInfo.NativeVault,
		TokenName:       tokenInfo.TokenName,
		TokenSymbol:     tokenInfo.Symbol,
		Creator:         tokenInfo.Creator,
		URI:             tokenInfo.URI,
		Price:           tokenInfo.Price,
		Holders:         tokenInfo.Holder,
		CreateTimestamp: tokenInfo.TransactionTime.Unix(),
		Rank:            0,
	}

	var extInfo model.ExtInfo
	if err := UnmarshalJSON(tokenInfo.ExtInfo, &extInfo); err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to unmarshal ExtInfo", err)
	}

	tickerRespons.MarketMetadata = response.MarketMetadata{
		ImageURL:    &extInfo.Image,
		Description: &extInfo.Description,
		Twitter:     &extInfo.Twitter,
		Website:     &extInfo.Website,
		Telegram:    &extInfo.Telegram,
		Banner:      &extInfo.Banner,
		Rules:       &extInfo.Rules,
		Sort:        &extInfo.Sort,
	}

	tokenMarketAnalytics, err := s.tokenMarketAnalyticsRepo.GetTokenMarketAnalytics(tokenAddress, uint8(chainType))
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to unmarshal ExtInfo", err)
	}
	if tokenMarketAnalytics != nil {
		tickerRespons.MarketTicker = response.MarketTicker{
			// High24H:            fmt.Sprintf("%f", tokenMarketAnalytics.Price24H),
			// Low24H:             fmt.Sprintf("%f", tokenMarketAnalytics.Price24H),
			TokenVolume24H:    fmt.Sprintf("%f", tokenMarketAnalytics.TokenVolume24H),
			BuyTokenVolume24H: fmt.Sprintf("%f", tokenMarketAnalytics.BuyTokenVolume24H),
			// NativeVolume24H:    "0",
			// BuyNativeVolume24H: "0",
			PriceChange24H: fmt.Sprintf("%f", tokenMarketAnalytics.PriceChange24H),
			TxCount24H:     int(tokenMarketAnalytics.TxCount24H),
			BuyTxCount24H:  int(tokenMarketAnalytics.BuyTxCount24H),
			// High1H:             fmt.Sprintf("%f", tokenMarketAnalytics.Price1H),
			// Low1H:              fmt.Sprintf("%f", tokenMarketAnalytics.Price1H),
			TokenVolume1H:    fmt.Sprintf("%f", tokenMarketAnalytics.TokenVolume1H),
			BuyTokenVolume1H: fmt.Sprintf("%f", tokenMarketAnalytics.BuyTokenVolume1H),
			// NativeVolume1H:     "0",
			// BuyNativeVolume1H:  "0",
			PriceChange1H: fmt.Sprintf("%f", tokenMarketAnalytics.PriceChange1H),
			TxCount1H:     int(tokenMarketAnalytics.TxCount1H),
			BuyTxCount1H:  int(tokenMarketAnalytics.BuyTxCount1H),
			// High5M:             fmt.Sprintf("%f", tokenMarketAnalytics.Price5M),
			// Low5M:              fmt.Sprintf("%f", tokenMarketAnalytics.Price5M),
			TokenVolume5M:    fmt.Sprintf("%f", tokenMarketAnalytics.TokenVolume5M),
			BuyTokenVolume5M: fmt.Sprintf("%f", tokenMarketAnalytics.BuyTokenVolume5M),
			// NativeVolume5M:     "0",
			// BuyNativeVolume5M:  "0",
			PriceChange5M: fmt.Sprintf("%f", tokenMarketAnalytics.PriceChange5M),
			TxCount5M:     int(tokenMarketAnalytics.TxCount5M),
			BuyTxCount5M:  int(tokenMarketAnalytics.BuyTxCount5M),
			LastSwapAt:    tokenInfo.TransactionTime.Unix(),
			MarketCap:     tokenInfo.MarketCap.String(),
		}
	}

	return response.Success(tickerRespons)
}

func (s *TickerServiceImpl) MarketTicker(tokenAddress string, chainType model.ChainType) response.Response {
	var marketTicker response.MarketTicker
	tokenMarketAnalytics, err := s.tokenMarketAnalyticsRepo.GetTokenMarketAnalytics(tokenAddress, uint8(chainType))
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to unmarshal ExtInfo", err)
	}

	if tokenMarketAnalytics != nil {
		marketTicker = response.MarketTicker{
			// High24H:            fmt.Sprintf("%f", tokenMarketAnalytics.Price24H),
			// Low24H:             fmt.Sprintf("%f", tokenMarketAnalytics.Price24H),
			TokenVolume24H:    fmt.Sprintf("%f", tokenMarketAnalytics.TokenVolume24H),
			BuyTokenVolume24H: fmt.Sprintf("%f", tokenMarketAnalytics.BuyTokenVolume24H),
			// NativeVolume24H:    "0",
			// BuyNativeVolume24H: "0",
			PriceChange24H: fmt.Sprintf("%f", tokenMarketAnalytics.PriceChange24H),
			TxCount24H:     int(tokenMarketAnalytics.TxCount24H),
			BuyTxCount24H:  int(tokenMarketAnalytics.BuyTxCount24H),
			// High1H:             fmt.Sprintf("%f", tokenMarketAnalytics.Price1H),
			// Low1H:              fmt.Sprintf("%f", tokenMarketAnalytics.Price1H),
			TokenVolume1H:    fmt.Sprintf("%f", tokenMarketAnalytics.TokenVolume1H),
			BuyTokenVolume1H: fmt.Sprintf("%f", tokenMarketAnalytics.BuyTokenVolume1H),
			// NativeVolume1H:     "0",
			// BuyNativeVolume1H:  "0",
			PriceChange1H: fmt.Sprintf("%f", tokenMarketAnalytics.PriceChange1H),
			TxCount1H:     int(tokenMarketAnalytics.TxCount1H),
			BuyTxCount1H:  int(tokenMarketAnalytics.BuyTxCount1H),
			// High5M:             fmt.Sprintf("%f", tokenMarketAnalytics.Price5M),
			// Low5M:              fmt.Sprintf("%f", tokenMarketAnalytics.Price5M),
			TokenVolume5M:    fmt.Sprintf("%f", tokenMarketAnalytics.TokenVolume5M),
			BuyTokenVolume5M: fmt.Sprintf("%f", tokenMarketAnalytics.BuyTokenVolume5M),
			// NativeVolume5M:     "0",
			// BuyNativeVolume5M:  "0",
			PriceChange5M: fmt.Sprintf("%f", tokenMarketAnalytics.PriceChange5M),
			TxCount5M:     int(tokenMarketAnalytics.TxCount5M),
			BuyTxCount5M:  int(tokenMarketAnalytics.BuyTxCount5M),
			// 查询es
			// LastSwapAt:    tokenInfo.TransactionTime.Unix(),
			// MarketCap:     tokenInfo.MarketCap.String(),
		}
	}
	return response.Success(marketTicker)
}

func (s *TickerServiceImpl) SwapHistories(tickersId string, chainType model.ChainType) response.Response {
	var swapHistoriesResponse response.SwapHistoriesResponse
	return response.Success(swapHistoriesResponse)
}

func (s *TickerServiceImpl) TokenDistribution(tokenAddress string, chainType model.ChainType) response.Response {
	redisKey := GetRedisKey(constants.TokenDistribution, tokenAddress)
	tokenDistribution, err := redis.Get(redisKey)
	if err != nil {
		util.Log().Error("Failed to get data from Redis: %v\n", err)
	}
	if tokenDistribution != "" {
		var tokenDistributionResponse response.TokenDistributionResponse
		if err := json.Unmarshal([]byte(tokenDistribution), &tokenDistributionResponse); err == nil {
			return response.Success(tokenDistributionResponse)
		}
		util.Log().Error("Failed to unmarshal token distribution data: %v\n", err)
	}
	BirdeyeClient := httpUtil.NewBirdeyeClient(httpUtil.SOLANA)
	tokenMarketDataRes, err := BirdeyeClient.GetTokenMarketData(tokenAddress)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token market data", err)
	}
	CirculatingSupply := tokenMarketDataRes.Data.CirculatingSupply

	tokenHoldersRes, err := BirdeyeClient.GetTokenHolders(tokenAddress, 0, 20)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token holders", err)
	}
	if !tokenHoldersRes.Success {
		return response.Err(http.StatusInternalServerError, "Failed to fetch token holders data", nil)
	}
	var tokenHolders []response.TokenHolder
	for _, holder := range tokenHoldersRes.Data.Items {
		var tokenHolder response.TokenHolder
		tokenHolder.Account = holder.Owner
		amount, err := strconv.ParseFloat(holder.Amount, 64)
		if err != nil {
			log.Printf("Failed to parse amount for holder %s: %v\n", holder.Owner, err)
			continue
		}
		percentage := (amount / math.Pow(10, float64(holder.Decimals))) / CirculatingSupply * 100

		tokenHolder.Percentage = strconv.FormatFloat(percentage, 'f', 2, 64)
		tokenHolder.IsAssociatedBondingCurve = false
		tokenHolder.UserProfile = nil
		tokenHolder.Amount = holder.Amount
		tokenHolder.UIAmount = holder.UIAmount
		var moderator response.Moderator
		moderator.BannedModID = 0
		moderator.Status = "NORMAL"
		moderator.Banned = false
		tokenHolder.Moderator = moderator
		tokenHolder.IsCommunityVault = false
		tokenHolder.IsBlackHole = false
		tokenHolders = append(tokenHolders, tokenHolder)
	}
	var tokenDistributionResponse response.TokenDistributionResponse
	tokenDistributionResponse.TokenHolders = tokenHolders

	if err := redis.Set(redisKey, tokenDistributionResponse, 10*time.Minute); err != nil {
		util.Log().Error("Failed to set data in Redis: %v\n", err)
	}
	return response.Success(tokenDistributionResponse)
}

func (s *TickerServiceImpl) SearchTickers(param, limit, cursor string, chainType model.ChainType) response.Response {
	var sarchTickerResponse response.SearchTickerResponse
	return response.Success(sarchTickerResponse)
}
