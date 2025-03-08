package service

import (
	"fmt"
	"game-fun-be/internal/clickhouse"
	"game-fun-be/internal/constants"
	"game-fun-be/internal/es"
	"game-fun-be/internal/es/query"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpRespone"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"

	"log"
	"math"
	"strconv"
	"time"

	"net/http"

	"github.com/shopspring/decimal"
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

	if tokenInfo != nil {
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
	} else {

		redisKey := GetRedisKey(constants.TokenMetaData, tokenAddress)
		var tokenMetaData httpRespone.TokenMetaData // 改为非指针类型

		// 1. 优先从 Redis 获取数据
		value, err := redis.Get(redisKey)
		if err != nil {
			util.Log().Error("Failed to get token meta data from Redis: %v", err)
		} else if value != "" {
			if err := redis.Unmarshal(value, &tokenMetaData); err != nil {
				util.Log().Error("Failed to unmarshal token meta data: %v", err)
			}
		}

		// 2. 如果 Redis 中没有数据，调用 API 获取数据
		if tokenMetaData.Address == "" { // 检查是否为空
			tokenMetaDatas, err := httpUtil.GetTokenMetaData([]string{tokenAddress}, chainType.ToString())
			if err != nil {
				return response.Err(http.StatusInternalServerError, "Failed to get token meta data", err)
			}

			// 3. 根据 tokenAddress 在 tokenMetaDatas.Data 中找到对应的代币信息
			var exists bool
			tokenMetaData, exists = tokenMetaDatas.Data[tokenAddress]
			if !exists {
				return response.Err(http.StatusNotFound, "Token meta data not found", fmt.Errorf("token address %s not found in response", tokenAddress))
			}

			// 4. 如果 API 返回的数据有效，缓存到 Redis
			if err := redis.Set(redisKey, tokenMetaData); err != nil {
				util.Log().Error("Failed to set token meta data in Redis: %v", err)
			}
		}

		// 5. 填充 MarketMetadata
		tickerRespons.MarketMetadata = response.MarketMetadata{
			ImageURL:    &tokenMetaData.LogoURI,
			Description: tokenMetaData.Extensions.Description,
			Twitter:     tokenMetaData.Extensions.Twitter,
			Website:     tokenMetaData.Extensions.Website,
			Telegram:    tokenMetaData.Extensions.Telegram,
			Github:      tokenMetaData.Extensions.Github,
			Banner:      nil,
			Rules:       nil,
			Sort:        nil,
		}

		tokenMarketDataRes, err := GetOrFetchTokenMarketData(tokenAddress, chainType)
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to get token market data", err)
		}

		tickerRespons.Market = response.Market{
			MarketID:  0,
			TokenMint: tokenMetaData.Address,
			// 待确定
			// Market:          userInfo.MarketAddress,
			// TokenVault:      userInfo.TokenVault,
			// NativeVault:     userInfo.NativeVault,
			TokenName:   tokenMetaData.Name,
			TokenSymbol: tokenMetaData.Symbol,
			// Creator:     tokenInfo.Creator,
			URI:   tokenMetaData.LogoURI,
			Price: decimal.NewFromFloat(tokenMarketDataRes.Data.Price),
			// CreateTimestamp: ,
			// Rank:            0,
		}
	}

	// // 从clickhoues 获取市场信息
	// tokenMarketAnalytics, err := s.tokenMarketAnalyticsRepo.GetTokenMarketAnalytics(tokenAddress, uint8(chainType))
	// if err != nil {
	// 	return response.Err(http.StatusInternalServerError, "Failed to unmarshal ExtInfo", err)
	// }

	marketTicker := s.MarketTicker(tokenAddress, chainType)
	if marketTicker.Code != http.StatusOK {
		return response.Err(http.StatusInternalServerError, "failed to get market ticker", fmt.Errorf(marketTicker.Msg))
	}
	marketData, ok := marketTicker.Data.(response.MarketTicker)
	if !ok {
		return response.Err(http.StatusInternalServerError, "failed to convert market ticker data to MarketTicker", fmt.Errorf("type assertion failed"))
	}

	tickerRespons.MarketTicker = marketData

	return response.Success(tickerRespons)
}

func (s *TickerServiceImpl) MarketTicker(tokenAddress string, chainType model.ChainType) response.Response {
	// // 从clickhoues 获取
	// tokenMarketAnalytics, err := s.tokenMarketAnalyticsRepo.GetTokenMarketAnalytics(tokenAddress, chainType.Uint8())
	// if err != nil {
	// 	return response.Err(http.StatusInternalServerError, "Failed to unmarshal ExtInfo", err)
	// }

	redisKey := GetRedisKey(constants.TokenTradeData, tokenAddress)

	var marketTicker response.MarketTicker
	var tradeData httpRespone.TradeData

	// 1. 优先从 Redis 获取数据
	value, err := redis.Get(redisKey)
	if err != nil {
		util.Log().Error("Failed to get data from Redis: %v", err)
	} else if value != "" {
		if err := redis.Unmarshal(value, &tradeData); err != nil {
			util.Log().Error("Failed to unmarshal trade data: %v", err)
		} else {
			marketTicker = populateMarketTicker(tradeData)
		}
	} else {
		// 2. 如果 Redis 中没有数据，调用 API 获取数据
		tokenMarketDataRes, err := httpUtil.GetTradeData(tokenAddress, chainType.ToString())
		if err != nil {
			util.Log().Error("Failed to get trade data for token %s on chain %s: %v", tokenAddress, chainType.ToString(), err)
			return response.Err(http.StatusInternalServerError, "failed to get trade data: %w", err)
		}

		// 3. 如果 API 返回的数据有效，更新 marketTicker 并缓存到 Redis
		if tokenMarketDataRes != nil {
			marketTicker = populateMarketTicker(tokenMarketDataRes.Data)

			// 缓存数据到 Redis，设置过期时间为 10 分钟
			if err := redis.Set(redisKey, tokenMarketDataRes.Data, 20*time.Minute); err != nil {
				util.Log().Error("Failed to set data in Redis: %v", err)
			}
		} else {
			// 如果 API 返回的数据为空，返回错误
			return response.Err(http.StatusInternalServerError, "API returned empty data", fmt.Errorf("API returned empty data"))
		}
	}

	// 4. 返回成功响应
	return response.Success(marketTicker)
}

// populateMarketTicker 将 TradeData 转换为 MarketTicker
func populateMarketTicker(tradeData httpRespone.TradeData) response.MarketTicker {
	return response.MarketTicker{
		Holders:              tradeData.Holder,
		TxCount24H:           int(tradeData.Trade24h),
		BuyTxCount24H:        int(tradeData.Buy24h),
		SellTxCount24H:       int(tradeData.Sell24h),
		TokenVolume24H:       fmt.Sprintf("%f", tradeData.Volume24h),
		TokenVolume24HUsd:    fmt.Sprintf("%f", tradeData.Volume24hUSD),
		BuyTokenVolume24H:    fmt.Sprintf("%f", tradeData.VolumeBuy24h),
		BuyTokenVolume24Usd:  fmt.Sprintf("%f", tradeData.VolumeBuy24hUSD),
		SellTokenVolume24H:   fmt.Sprintf("%f", tradeData.VolumeSell24h),
		SellTokenVolume24Usd: fmt.Sprintf("%f", tradeData.VolumeSell24hUSD),
		PriceChange24H:       fmt.Sprintf("%f", tradeData.PriceChange24hPercent),
		TxCount1H:            int(tradeData.Trade1h),
		BuyTxCount1H:         int(tradeData.Buy1h),
		SellTxCount1H:        int(tradeData.Sell1h),
		TokenVolume1H:        fmt.Sprintf("%f", tradeData.Volume1h),
		TokenVolume1HUsd:     fmt.Sprintf("%f", tradeData.Volume1hUSD),
		BuyTokenVolume1H:     fmt.Sprintf("%f", tradeData.VolumeBuy1h),
		BuyTokenVolume1Usd:   fmt.Sprintf("%f", tradeData.VolumeBuy1hUSD),
		SellTokenVolume1H:    fmt.Sprintf("%f", tradeData.VolumeSell1h),
		SellTokenVolume1Usd:  fmt.Sprintf("%f", tradeData.VolumeSell1hUSD),
		PriceChange1H:        fmt.Sprintf("%f", tradeData.PriceChange1hPercent),
		TxCount30M:           int(tradeData.Trade30m),
		BuyTxCount30M:        int(tradeData.Buy30m),
		SellTxCount30M:       int(tradeData.Sell30m),
		TokenVolume30M:       fmt.Sprintf("%f", tradeData.Volume30m),
		TokenVolume30MUsd:    fmt.Sprintf("%f", tradeData.Volume30mUsd),
		BuyTokenVolume30M:    fmt.Sprintf("%f", tradeData.VolumeBuy30m),
		BuyTokenVolume30Usd:  fmt.Sprintf("%f", tradeData.VolumeBuy30mUsd),
		SellTokenVolume30M:   fmt.Sprintf("%f", tradeData.VolumeSell30m),
		SellTokenVolume30Usd: fmt.Sprintf("%f", tradeData.VolumeSell30mUsd),
		PriceChange30M:       fmt.Sprintf("%f", tradeData.PriceChange30mPercent),
		LastSwapAt:           tradeData.LastTradeUnixTime,
	}
}

func (s *TickerServiceImpl) SwapHistories(tickersId string, chainType model.ChainType) response.Response {
	var swapHistoriesResponse response.SwapHistoriesResponse
	return response.Success(swapHistoriesResponse)
}

func (s *TickerServiceImpl) TokenDistribution(tokenAddress string, chainType model.ChainType) response.Response {
	redisKey := GetRedisKey(constants.TokenDistribution, tokenAddress)

	var tokenDistributionResponse response.TokenDistributionResponse
	value, err := redis.Get(redisKey)
	if err != nil {
		util.Log().Error("Failed to get token distribution data from Redis: %v", err)
	} else if value != "" {
		if err := redis.Unmarshal(value, &tokenDistributionResponse); err != nil {
			util.Log().Error("Failed to unmarshal token distribution data: %v", err)
		} else if len(tokenDistributionResponse.TokenHolders) > 0 {
			return response.Success(tokenDistributionResponse)
		}
	}

	tokenMarketDataRes, err := GetOrFetchTokenMarketData(tokenAddress, chainType)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token market data", err)
	}

	CirculatingSupply := tokenMarketDataRes.Data.CirculatingSupply

	tokenHoldersRes, err := httpUtil.GetTokenHolders(tokenAddress, 0, 20, chainType.ToString())
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
	tokenDistributionResponse.TokenHolders = tokenHolders

	if err := redis.Set(redisKey, tokenDistributionResponse, 20*time.Minute); err != nil {
		log.Printf("Failed to set data in Redis: %v\n", err)
	}
	return response.Success(tokenDistributionResponse)
}

func (s *TickerServiceImpl) SearchTickers(param, limit, cursor string, chainType model.ChainType) response.Response {
	var sarchTickerResponse response.SearchTickerResponse
	return response.Success(sarchTickerResponse)
}

func GetOrFetchTokenMarketData(tokenAddress string, chainType model.ChainType) (*httpRespone.TokenMarketDataResponse, error) {
	redisKey := GetRedisKey(constants.TokenMarketData, tokenAddress)
	var tokenMarketDataRes *httpRespone.TokenMarketDataResponse

	// 1. 优先从 Redis 获取数据
	value, err := redis.Get(redisKey)
	if err != nil {
		util.Log().Error("Failed to get token market data from Redis: %v", err)
	} else if value != "" {
		if err := redis.Unmarshal(value, &tokenMarketDataRes); err != nil {
			util.Log().Error("Failed to unmarshal token market data: %v", err)
		}
	}

	// 2. 如果 Redis 中没有数据，调用 API 获取数据
	if tokenMarketDataRes == nil {
		tokenMarketDataRes, err = httpUtil.GetTokenMarketData(tokenAddress, chainType.ToString())
		if err != nil {
			return nil, fmt.Errorf("failed to get token market data: %w", err)
		}

		// 3. 如果 API 返回的数据有效，缓存到 Redis
		if tokenMarketDataRes != nil {
			if err := redis.Set(redisKey, tokenMarketDataRes, 10*time.Minute); err != nil {
				util.Log().Error("Failed to set token market data in Redis: %v", err)
			}
		}
	}

	return tokenMarketDataRes, nil
}
