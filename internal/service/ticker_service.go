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
	"net/http"
	"reflect"
	"strconv"
	"time"

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
	var tickerResponse response.GetTickerResponse

	tokenInfo, err := s.tokenInfoRepo.GetTokenInfoByAddress(tokenAddress, uint8(chainType))
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token info by address", err)
	}
	if tokenInfo == nil {
		return response.Success(tickerResponse)
	}

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

	// 3. 获取市场信息
	marketTicker := s.MarketTicker(tokenAddress, chainType)
	if marketTicker.Code != http.StatusOK {
		return response.Err(http.StatusInternalServerError, "Failed to get market ticker", fmt.Errorf(marketTicker.Msg))
	}
	var marketData response.MarketTicker
	if marketTicker.Data != nil {
		var ok bool
		marketData, ok = marketTicker.Data.(response.MarketTicker)
		if !ok {
			return response.Err(http.StatusInternalServerError, "Failed to convert market ticker data to MarketTicker", fmt.Errorf("type assertion failed"))
		}
	}
	tickerResponse.MarketTicker = marketData

	return response.Success(tickerResponse)
}

func GetMarketData(tokenAddress, chainType string) (*httpRespone.TradeDataResponse, error) {
	var TradeDataResponse *httpRespone.TradeDataResponse
	marketDataKey := GetRedisKey(constants.TokenTradeData, tokenAddress)

	// 从 Redis 获取市场数据
	marketData, err := redis.Get(marketDataKey)
	if err != nil {
		util.Log().Error(fmt.Sprintf("Failed to get market data from Redis: %v", err))
	} else if marketData != "" {
		// 如果 Redis 中有数据，反序列化
		err = redis.Unmarshal(marketData, &TradeDataResponse)
		if err != nil {
			util.Log().Error(fmt.Sprintf("Failed to unmarshal market data: %v", err))
		} else {
			// 使用 Redis 中的数据
			return TradeDataResponse, nil
		}
	}

	// 如果 Redis 中没有数据，通过 HTTP 请求获取交易数据
	tokenTradeData, err := httpUtil.GetTradeData(tokenAddress, chainType)
	if err != nil {
		util.Log().Error(fmt.Sprintf("Failed to get trade data from HTTP: %v", err))
		return nil, err
	}

	// 检查 tokenTradeData 是否为 nil
	if tokenTradeData == nil {
		util.Log().Error("Trade data is nil")
		return nil, fmt.Errorf("trade data is nil")
	}

	// 检查 tokenTradeData.Data 是否为零值
	if reflect.ValueOf(tokenTradeData.Data).IsZero() {
		util.Log().Error("Trade data is empty")
		return nil, fmt.Errorf("trade data is empty")
	}

	// 将交易数据存储到 Redis
	err = redis.Set(marketDataKey, tokenTradeData, 3*time.Minute)
	if err != nil {
		util.Log().Error(fmt.Sprintf("Failed to set trade data to Redis: %v", err))
	}

	return tokenTradeData, nil
}

func (s *TickerServiceImpl) MarketTicker(tokenAddress string, chainType model.ChainType) response.Response {
	var analytics response.TokenMarketAnalyticsResponse

	queryJSON, err := query.TokenMarketAnalyticsQuery(tokenAddress, uint8(chainType))
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token creation info from API", err)
	}

	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if result == nil {
		return response.Success(analytics)
	}
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token transaction data from Elasticsearch", err)
	}

	aggregationResult, err := es.UnmarshalAggregationResult(result)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token creation info from API", err)
	}

	if len(aggregationResult.Buckets) == 0 {
		return response.Success(analytics)
	}

	tokenHolders := 0
	tradeDataResponse, err := GetMarketData(tokenAddress, chainType.ToString())
	if err == nil {
		tokenHolders = tradeDataResponse.Data.Holder
	}

	buyVolume1m := decimal.NewFromInt(0)
	sellVolume1m := decimal.NewFromInt(0)
	buyVolume5m := decimal.NewFromInt(0)
	sellVolume5m := decimal.NewFromInt(0)
	buyVolume1h := decimal.NewFromInt(0)
	sellVolume1h := decimal.NewFromInt(0)
	buyVolume24h := decimal.NewFromInt(0)
	sellVolume24h := decimal.NewFromInt(0)

	buyCount1m := decimal.NewFromInt(0)
	sellCount1m := decimal.NewFromInt(0)
	buyCount5m := decimal.NewFromInt(0)
	sellCount5m := decimal.NewFromInt(0)
	buyCount1h := decimal.NewFromInt(0)
	sellCount1h := decimal.NewFromInt(0)
	buyCount24h := decimal.NewFromInt(0)
	sellCount24h := decimal.NewFromInt(0)

	price := float64(0)
	lastSwapAt := ""

	price1m := float64(0)
	price5m := float64(0)
	price1h := float64(0)
	price24h := float64(0)
	marketCap := float64(0)
	nativePrice := float64(0)
	solPrice := float64(0)
	priceChange1m := float64(0)
	priceChange5m := float64(0)
	priceChange1h := float64(0)
	priceChange24h := float64(0)
	var decimals int

	for _, bucket := range aggregationResult.Buckets {

		if len(bucket.LastTransactionPrice.Latest.Hits.Hits) > 0 {
			price = bucket.LastTransactionPrice.Latest.Hits.Hits[0].Source.Price
			decimals = bucket.LastTransactionPrice.Latest.Hits.Hits[0].Source.Decimals
			nativePrice = bucket.LastTransactionPrice.Latest.Hits.Hits[0].Source.NativePrice
			lastSwapAt = bucket.LastTransactionPrice.Latest.Hits.Hits[0].Source.TransactionTime
			marketCap = bucket.LastTransactionPrice.Latest.Hits.Hits[0].Source.MarketCap

			// sol 的价格
			solPrice = price / nativePrice
			decimals = response.SolDecimals

		} else {
			decimals = 0
			price = 0
			continue
		}

		if len(bucket.LastTransaction1mPrice.Latest.Hits.Hits) > 0 {
			price1m = bucket.LastTransaction1mPrice.Latest.Hits.Hits[0].Source.Price
		} else {
			price1m = 0
		}
		if len(bucket.LastTransaction5mPrice.Latest.Hits.Hits) > 0 {
			price5m = bucket.LastTransaction5mPrice.Latest.Hits.Hits[0].Source.Price
		}
		if len(bucket.LastTransaction1hPrice.Latest.Hits.Hits) > 0 {
			price1h = bucket.LastTransaction1hPrice.Latest.Hits.Hits[0].Source.Price
		}
		if len(bucket.LastTransaction24hPrice.Latest.Hits.Hits) > 0 {
			price24h = bucket.LastTransaction24hPrice.Latest.Hits.Hits[0].Source.Price
		}

		priceChange1m = calculatePriceChange(price, price1m)
		priceChange5m = calculatePriceChange(price, price5m)
		priceChange1h = calculatePriceChange(price, price1h)
		priceChange24h = calculatePriceChange(price, price24h)

		if bucket.BuyVolume1m.TotalVolume.Value > 0 {
			buyVolume1m, _ = processVolume(bucket.BuyVolume1m.TotalVolume.Value, solPrice, decimals)
		}
		if bucket.SellVolume1m.TotalVolume.Value > 0 {
			sellVolume1m, _ = processVolume(bucket.SellVolume1m.TotalVolume.Value, solPrice, decimals)
		}
		if bucket.BuyVolume5m.TotalVolume.Value > 0 {
			buyVolume5m, _ = processVolume(bucket.BuyVolume5m.TotalVolume.Value, solPrice, decimals)
		}
		if bucket.SellVolume5m.TotalVolume.Value > 0 {
			sellVolume5m, _ = processVolume(bucket.SellVolume5m.TotalVolume.Value, solPrice, decimals)
		}
		if bucket.BuyVolume1h.TotalVolume.Value > 0 {
			buyVolume1h, _ = processVolume(bucket.BuyVolume1h.TotalVolume.Value, solPrice, decimals)
		}
		if bucket.SellVolume1h.TotalVolume.Value > 0 {
			sellVolume1h, _ = processVolume(bucket.SellVolume1h.TotalVolume.Value, solPrice, decimals)
		}
		if bucket.BuyVolume24h.TotalVolume.Value > 0 {
			buyVolume24h, _ = processVolume(bucket.BuyVolume24h.TotalVolume.Value, solPrice, decimals)
		}
		if bucket.SellVolume24h.TotalVolume.Value > 0 {
			sellVolume24h, _ = processVolume(bucket.SellVolume24h.TotalVolume.Value, solPrice, decimals)
		}

		if bucket.BuyCount1m.BuyVolume.Value > 0 {
			buyCount1m = decimal.NewFromInt(bucket.BuyCount1m.BuyVolume.Value)
		}
		if bucket.SellCount1m.SellVolume.Value > 0 {
			sellCount1m = decimal.NewFromInt(bucket.SellCount1m.SellVolume.Value)
		}
		if bucket.BuyCount5m.BuyVolume.Value > 0 {
			buyCount5m = decimal.NewFromInt(bucket.BuyCount5m.BuyVolume.Value)
		}
		if bucket.SellCount5m.SellVolume.Value > 0 {
			sellCount5m = decimal.NewFromInt(bucket.SellCount5m.SellVolume.Value)
		}
		if bucket.BuyCount1h.BuyVolume.Value > 0 {
			buyCount1h = decimal.NewFromInt(bucket.BuyCount1h.BuyVolume.Value)
		}
		if bucket.SellCount1h.SellVolume.Value > 0 {
			sellCount1h = decimal.NewFromInt(bucket.SellCount1h.SellVolume.Value)
		}
		if bucket.BuyCount24h.BuyVolume.Value > 0 {
			buyCount24h = decimal.NewFromInt(bucket.BuyCount24h.BuyVolume.Value)
		}
		if bucket.SellCount24h.SellVolume.Value > 0 {
			sellCount24h = decimal.NewFromInt(bucket.SellCount24h.SellVolume.Value)
		}
	}

	timestamp, err := StringToTimestamp(lastSwapAt, ISO8601Layout)
	if err != nil {
		fmt.Println("Error:", err)
		timestamp = 0
	}
	analytics.Holders = tokenHolders
	analytics.TokenAddress = tokenAddress
	analytics.BuyVolume1m = buyVolume1m
	analytics.SellVolume1m = sellVolume1m
	analytics.BuyVolume5m = buyVolume5m
	analytics.SellVolume5m = sellVolume5m
	analytics.BuyVolume1h = buyVolume1h
	analytics.SellVolume1h = sellVolume1h
	analytics.BuyVolume24h = buyVolume24h
	analytics.SellVolume24h = sellVolume24h
	analytics.BuyCount1m = buyCount1m
	analytics.BuyCount5m = buyCount5m
	analytics.BuyCount1h = buyCount1h
	analytics.BuyCount24h = buyCount24h
	analytics.SellCount1m = sellCount1m
	analytics.SellCount5m = sellCount5m
	analytics.SellCount1h = sellCount1h
	analytics.SellCount24h = sellCount24h
	analytics.PriceChange1m = priceChange1m
	analytics.PriceChange5m = priceChange5m
	analytics.PriceChange1h = priceChange1h
	analytics.PriceChange24h = priceChange24h
	analytics.CurrentPrice = price
	analytics.Volume1m = analytics.BuyVolume1m.Add(analytics.SellVolume1m)
	analytics.Volume5m = analytics.BuyVolume5m.Add(analytics.SellVolume5m)
	analytics.Volume1h = analytics.BuyVolume1h.Add(analytics.SellVolume1h)
	analytics.Volume24h = analytics.BuyVolume24h.Add(analytics.SellVolume24h)
	analytics.TotalCount1m = analytics.BuyCount1m.Add(analytics.SellCount1m)
	analytics.TotalCount5m = analytics.BuyCount5m.Add(analytics.SellCount5m)
	analytics.TotalCount1h = analytics.BuyCount1h.Add(analytics.SellCount1h)
	analytics.TotalCount24h = analytics.BuyCount24h.Add(analytics.SellCount24h)
	analytics.MarketCap = strconv.FormatFloat(marketCap, 'f', -1, 64)
	analytics.LastSwapAt = timestamp
	marketTicker := populateMarketTicker(analytics)

	// 4. 返回成功响应
	return response.Success(marketTicker)
}

// populateMarketTicker 将 TradeData 转换为 MarketTicker
func populateMarketTicker(analytics response.TokenMarketAnalyticsResponse) response.MarketTicker {
	priceChange24hPercentStr := FormatPercent(analytics.PriceChange24h)
	priceChange1hPercentStr := FormatPercent(analytics.PriceChange1h)
	priceChange5mPercentStr := FormatPercent(analytics.PriceChange5m)
	txCount24H := ConvertDecimalToInt(analytics.TotalCount24h, false)
	buyCount24h := ConvertDecimalToInt(analytics.BuyCount24h, false)
	sellCount24h := ConvertDecimalToInt(analytics.SellCount24h, false)
	txCount1H := ConvertDecimalToInt(analytics.TotalCount1h, false)
	buyCount1h := ConvertDecimalToInt(analytics.BuyCount1h, false)
	sellCount1h := ConvertDecimalToInt(analytics.SellCount1h, false)
	txCount5m := ConvertDecimalToInt(analytics.TotalCount5m, false)
	buyCount5m := ConvertDecimalToInt(analytics.BuyCount5m, false)
	sellCount5m := ConvertDecimalToInt(analytics.SellCount5m, false)
	currentPrice := decimal.NewFromFloat(analytics.CurrentPrice)
	roundTo9 := func(d decimal.Decimal) decimal.Decimal {
		return d.Round(9)
	}
	return response.MarketTicker{
		TxCount24H:           txCount24H,
		BuyTxCount24H:        buyCount24h,
		SellTxCount24H:       sellCount24h,
		TokenVolume24H:       analytics.Volume24h.String(),
		BuyTokenVolume24H:    analytics.BuyVolume24h.String(),
		SellTokenVolume24H:   analytics.SellVolume24h.String(),
		PriceChange24H:       priceChange24hPercentStr,
		TxCount1H:            txCount1H,
		BuyTxCount1H:         buyCount1h,
		SellTxCount1H:        sellCount1h,
		TokenVolume1H:        analytics.Volume1h.String(),
		BuyTokenVolume1H:     analytics.BuyVolume1h.String(),
		SellTokenVolume1H:    analytics.SellVolume1h.String(),
		PriceChange1H:        priceChange1hPercentStr,
		TxCount5M:            txCount5m,
		BuyTxCount5M:         buyCount5m,
		SellTxCount5M:        sellCount5m,
		TokenVolume5M:        analytics.Volume5m.String(),
		BuyTokenVolume5M:     analytics.BuyVolume5m.String(),
		SellTokenVolume5M:    analytics.SellVolume5m.String(),
		PriceChange5M:        priceChange5mPercentStr,
		TokenVolume24HUsd:    roundTo9(currentPrice.Mul(analytics.Volume24h)).String(),
		BuyTokenVolume24Usd:  roundTo9(currentPrice.Mul(analytics.BuyVolume24h)).String(),
		SellTokenVolume24Usd: roundTo9(currentPrice.Mul(analytics.SellVolume24h)).String(),
		TokenVolume1HUsd:     roundTo9(currentPrice.Mul(analytics.Volume1h)).String(),
		BuyTokenVolume1Usd:   roundTo9(currentPrice.Mul(analytics.BuyVolume1h)).String(),
		SellTokenVolume1Usd:  roundTo9(currentPrice.Mul(analytics.SellVolume1h)).String(),
		TokenVolume5MUsd:     roundTo9(currentPrice.Mul(analytics.Volume5m)).String(),
		BuyTokenVolume5Usd:   roundTo9(currentPrice.Mul(analytics.BuyVolume5m)).String(),
		SellTokenVolume5Usd:  roundTo9(currentPrice.Mul(analytics.SellVolume5m)).String(),
		LastSwapAt:           analytics.LastSwapAt,
		MarketCap:            analytics.MarketCap,
		Holders:              analytics.Holders,
		Rank:                 1,
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
		TransactionType := item.TransactionType
		if item.IsBuyback {
			TransactionType = 3
		}

		// Create a new transaction history
		history := response.TransactionHistory{
			TradeType:    TransactionType,
			Payer:        item.UserAddress,
			Signature:    item.TransactionHash,
			BlockTime:    item.TransactionTime,
			TokenAmount:  item.QuoteTokenAmount.String(),
			NativeAmount: item.BaseTokenAmount.String(),
			TokenPrice:   item.QuoteTokenPrice.String(),
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

	tokenMarketDataRes, err := GetOrFetchTokenMarketData(tokenAddress, chainType.ToString())
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
	if len(tokenHoldersRes.Data.Items) == 0 {
		return response.Success("No token holders found")
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

func GetOrFetchTokenMarketData(tokenAddress string, chainType string) (*httpRespone.TokenMarketDataResponse, error) {
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
		tokenMarketDataRes, err = httpUtil.GetTokenMarketData(tokenAddress, chainType)
		if err != nil {
			return nil, fmt.Errorf("failed to get token market data: %w", err)
		}

		// 3. 如果 API 返回的数据有效，缓存到 Redis
		if tokenMarketDataRes != nil {
			if err := redis.Set(redisKey, tokenMarketDataRes, 24*time.Hour); err != nil {
				util.Log().Error("Failed to set token market data in Redis: %v", err)
			}
		}
	}

	return tokenMarketDataRes, nil
}
