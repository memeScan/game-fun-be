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

	"encoding/json"
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
	// 1. 获取代币信息
	tokenInfo, err := s.getTokenInfo(tokenAddress, chainType)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token info", err)
	}

	// 2. 填充响应数据
	var tickerResponse response.GetTickerResponse
	if tokenInfo != nil {
		tickerResponse.Market = response.Market{
			MarketID:        tokenInfo.ID,
			TokenMint:       tokenInfo.TokenAddress,
			Decimals:        tokenInfo.Decimals,
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

		tickerResponse.MarketMetadata = response.MarketMetadata{
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
		// 如果 MySQL 中没有数据，从 API 获取并填充
		tokenMetaData, err := s.getTokenMetaDataFromAPI(tokenAddress, chainType)
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to get token meta data from API", err)
		}

		tokenMarketDataRes, err := GetOrFetchTokenMarketData(tokenAddress, chainType)
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to get token market data", err)
		}

		tickerResponse.Market = response.Market{
			MarketID:    0,
			TokenMint:   tokenMetaData.Address,
			Decimals:    tokenMetaData.Decimals,
			TokenName:   tokenMetaData.Name,
			TokenSymbol: tokenMetaData.Symbol,

			URI:             tokenMetaData.LogoURI,
			Price:           decimal.NewFromFloat(tokenMarketDataRes.Data.Price),
			CreateTimestamp: 0, // 需要从 API 或其他地方获取
			Rank:            0,
		}

		tickerResponse.MarketMetadata = response.MarketMetadata{
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

		// 插入到mysql
		var token model.TokenInfo
		token.TokenAddress = tokenMetaData.Address
		token.TokenName = tokenMetaData.Name
		token.Symbol = tokenMetaData.Symbol
		extInfoJSON, err := json.Marshal(tokenMetaData.Extensions)
		if err != nil {
			util.Log().Error("Failed to marshal Extensions to JSON: %v", err)
			return response.Err(http.StatusInternalServerError, "Failed to marshal Extensions to JSON", err)
		}
		token.ExtInfo = string(extInfoJSON)
	}

	// 3. 获取市场信息
	marketTicker := s.MarketTicker(tokenAddress, chainType)
	if marketTicker.Code != http.StatusOK {
		return response.Err(http.StatusInternalServerError, "Failed to get market ticker", fmt.Errorf(marketTicker.Msg))
	}
	marketData, ok := marketTicker.Data.(response.MarketTicker)
	if !ok {
		return response.Err(http.StatusInternalServerError, "Failed to convert market ticker data to MarketTicker", fmt.Errorf("type assertion failed"))
	}

	tickerResponse.MarketTicker = marketData

	// 4. 返回响应
	return response.Success(tickerResponse)
}

// getTokenInfo 获取代币信息（优先从 Redis 和 MySQL 获取）
func (s *TickerServiceImpl) getTokenInfo(tokenAddress string, chainType model.ChainType) (*model.TokenInfo, error) {
	// 1. 从 Redis 获取
	redisKey := GetRedisKey(constants.TokenMetaData, tokenAddress)
	value, err := redis.Get(redisKey)
	if err == nil && value != "" {
		var tokenInfo model.TokenInfo
		if err := redis.Unmarshal(value, &tokenInfo); err == nil {
			return &tokenInfo, nil
		}
		util.Log().Error("Failed to unmarshal token info from Redis: %v", err)
	}

	// 2. 从 MySQL 获取
	tokenInfo, err := s.tokenInfoRepo.GetTokenInfoByAddress(tokenAddress, uint8(chainType))
	if err != nil {
		return nil, fmt.Errorf("failed to get token info from MySQL: %v", err)
	}
	if tokenInfo != nil {
		// 将数据缓存到 Redis
		if err := redis.Set(redisKey, tokenInfo); err != nil {
			util.Log().Error("Failed to set token info in Redis: %v", err)
		}
		return tokenInfo, nil
	}

	// 3. 如果 Redis 和 MySQL 都没有数据，返回 nil
	return nil, nil
}

// getTokenMetaDataFromAPI 从 API 获取代币元数据
func (s *TickerServiceImpl) getTokenMetaDataFromAPI(tokenAddress string, chainType model.ChainType) (httpRespone.TokenMetaData, error) {
	tokenMetaDatas, err := httpUtil.GetTokenMetaData([]string{tokenAddress}, chainType.ToString())
	if err != nil {
		return httpRespone.TokenMetaData{}, fmt.Errorf("failed to get token meta data from API: %v", err)
	}

	tokenMetaData, exists := tokenMetaDatas.Data[tokenAddress]
	if !exists {
		return httpRespone.TokenMetaData{}, fmt.Errorf("token address %s not found in API response", tokenAddress)
	}

	// 将数据缓存到 Redis
	redisKey := GetRedisKey(constants.TokenMetaData, tokenAddress)
	if err := redis.Set(redisKey, tokenMetaData); err != nil {
		util.Log().Error("Failed to set token meta data in Redis: %v", err)
	}

	return tokenMetaData, nil
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

	// Get token transactions from ClickHouse
	service := TransactionCkServiceImpl{}
	resp := service.GetTokenOrderBook(tickersId, uint8(chainType))

	// Check if there was an error
	if resp.Code != response.CodeSuccess {
		return resp
	}

	// Convert the token order book items to transaction histories
	items, ok := resp.Data.([]response.TokenOrderBookItem)
	if !ok {
		return response.Err(response.CodeServerUnknown, "Failed to convert token order book data", nil)
	}

	// Create transaction histories
	transactionHistories := make([]response.TransactionHistory, 0, len(items))
	for _, item := range items {
		// Convert transaction type to isBuy (1 is buy, 2 is sell)
		isBuy := item.TransactionType == 1

		// Format the transaction time
		blockTime := time.Unix(item.TransactionTime, 0).Format(time.RFC3339)

		// Create a new transaction history
		history := response.TransactionHistory{
			IsBuy:        isBuy,
			Payer:        item.UserAddress,
			Signature:    item.TransactionHash,
			BlockTime:    blockTime,
			TokenAmount:  item.QuoteTokenAmount.String(),
			NativeAmount: item.BaseTokenAmount.String(),
		}

		transactionHistories = append(transactionHistories, history)
	}

	// Set the response data
	swapHistoriesResponse.TransactionHistories = transactionHistories
	swapHistoriesResponse.HasMore = len(transactionHistories) >= 100 // Assuming limit is 100

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
