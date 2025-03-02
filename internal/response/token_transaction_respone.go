package response

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"
)

// PoolMarketInfo 交易池市场信息
type PoolMarketInfo struct {
	// 基础信息
	ID            int64  `json:"id"`
	PoolAddress   string `json:"address"`
	Creator       string `json:"creator"`
	PoolType      int    `json:"pool_type"`
	PoolTypeStr   string `json:"pool_type_str"`
	OpenTimestamp int64  `json:"open_timestamp"`
	Launchpad     string `json:"launchpad"`

	// 代币信息
	QuoteAddress string `json:"quote_address"`
	BaseAddress  string `json:"base_address"`
	QuoteSymbol  string `json:"quote_symbol"`
	BaseSymbol   string `json:"base_symbol"`

	// 储备金额
	QuoteReserve        string `json:"quote_reserve"`
	BaseReserve         string `json:"base_reserve"`
	InitialLiquidity    string `json:"initial_liquidity"`
	InitialQuoteReserve string `json:"initial_quote_reserve"`

	// 关联代币信息
	QuoteTokenInfo *TokenTransactionResponse `json:"quote_token_info"`
	BaseTokenInfo  *TokenTransactionResponse `json:"base_token_info"`
}

// TokenTransaction 代币交易信息
type TokenTransaction struct {
	// 池信息
	Price              float64 `json:"price,omitempty" default:"0"`
	NativePrice        float64 `json:"native_price,omitempty" default:"0"`
	PoolAddress        string  `json:"pool_address,omitempty" default:""`
	BaseAddress        string  `json:"base_address,omitempty" default:""`
	QuoteAddress       string  `json:"quote_address,omitempty" default:""`
	QuoteReserve       string  `json:"quote_reserve,omitempty" default:"0"`
	RealNativeReserves string  `json:"real_native_reserves,omitempty" default:"0"`
	Creator            string  `json:"creator,omitempty" default:""`
	PoolTypeStr        string  `json:"pool_type_str,omitempty" default:"unknown"`
	PoolType           int     `json:"pool_type,omitempty" default:"0"`
	QuoteSymbol        string  `json:"quote_symbol,omitempty" default:""`
	Launchpad          string  `json:"launchpad,omitempty" default:""`
	OpenTimestamp      string  `json:"block_time,omitempty" default:""`

	// 代币基本信息
	TokenAddress          string  `json:"token_address,omitempty" default:""`
	ExtInfo               string  `json:"ext_info,omitempty" default:"{}"`
	UserAddress           string  `json:"user_address,omitempty" default:""`
	TokenSupply           string  `json:"token_supply,omitempty" default:"0"`
	ChainType             int     `json:"chain_type,omitempty" default:"0"`
	PlatformType          int     `json:"platform_type,omitempty" default:"0"`
	Decimals              int     `json:"decimals,omitempty" default:"0"`
	TokenCreateTime       string  `json:"token_create_time,omitempty" default:""`
	MarketCap             float64 `json:"market_cap,omitempty" default:"0"`
	CreatedPlatformType   int     `json:"created_platform_type,omitempty" default:"0"`
	Progress              float64 `json:"progress,omitempty" default:"0"`
	IsComplete            bool    `json:"is_complete,omitempty" default:"false"`
	LatestTransactionTime string  `json:"transaction_time,omitempty" default:""`

	// 储备信息
	VirtualNativeReserves string  `json:"virtual_native_reserves,omitempty" default:"0"`
	VirtualTokenReserves  string  `json:"virtual_token_reserves,omitempty" default:"0"`
	InitialNativeReserve  float64 `json:"initial_pc_reserve,omitempty" default:"0"`
	InitialTokenReserve   float64 `json:"initial_token_reserve,omitempty" default:"0"`
	NativeReserveRate     float64 `json:"native_reserve_rate,omitempty" default:"0"`

	// 交易统计
	SwapsCount1m  int64   `json:"swaps_1m,omitempty" default:"0"`
	SwapsCount5m  int64   `json:"swaps_5m,omitempty" default:"0"`
	SwapsCount1h  int64   `json:"swaps_1h,omitempty" default:"0"`
	SwapsCount6h  int64   `json:"swaps_6h,omitempty" default:"0"`
	SwapsCount24h int64   `json:"swaps_24h,omitempty" default:"0"`
	Volume1m      float64 `json:"volume_1m,omitempty" default:"0"`
	Volume5m      float64 `json:"volume_5m,omitempty" default:"0"`
	Volume1h      float64 `json:"volume_1h,omitempty" default:"0"`
	Volume6h      float64 `json:"volume_6h,omitempty" default:"0"`
	Volume24h     float64 `json:"volume_24h,omitempty" default:"0"`
	BuyCount1m    int64   `json:"buy_count_1m,omitempty" default:"0"`
	BuyCount5m    int64   `json:"buy_count_5m,omitempty" default:"0"`
	BuyCount1h    int64   `json:"buy_count_1h,omitempty" default:"0"`
	BuyCount6h    int64   `json:"buy_count_6h,omitempty" default:"0"`
	BuyCount24h   int64   `json:"buy_count_24h,omitempty" default:"0"`
	SellCount1m   int64   `json:"sell_count_1m,omitempty" default:"0"`
	SellCount5m   int64   `json:"sell_count_5m,omitempty" default:"0"`
	SellCount1h   int64   `json:"sell_count_1h,omitempty" default:"0"`
	SellCount6h   int64   `json:"sell_count_6h,omitempty" default:"0"`
	SellCount24h  int64   `json:"sell_count_24h,omitempty" default:"0"`

	// 市场变化
	MarketCapTime        float64 `json:"market_cap_time,omitempty" default:"0"`
	MarketCapChange1m    float64 `json:"market_cap_change_1m,omitempty" default:"0"`
	MarketCapChange5m    float64 `json:"market_cap_change_5m,omitempty" default:"0"`
	MarketCapChange1h    float64 `json:"market_cap_change_1h,omitempty" default:"0"`
	PriceChangePercent1m float64 `json:"price_change_percent1m,omitempty" default:"0"`
	PriceChangePercent5m float64 `json:"price_change_percent5m,omitempty" default:"0"`
	PriceChangePercent1h float64 `json:"price_change_percent1h,omitempty" default:"0"`

	// 总计数据
	Swaps              int64   `json:"swaps,omitempty" default:"0"`
	Volume             float64 `json:"volume,omitempty" default:"0"`
	Buys               int64   `json:"buys,omitempty" default:"0"`
	Sells              int64   `json:"sells,omitempty" default:"0"`
	PriceChangePercent float64 `json:"price_change_percent,omitempty" default:"0"`
	CreatorBalanceRate string  `json:"creator_balance_rate,omitempty" default:"0"`

	// 代币标志信息
	Holder               int     `json:"holder,omitempty" default:"0"`
	DevStatus            int     `json:"dev_status,omitempty" default:"0"`
	CrownDuration        int     `json:"crown_duration,omitempty" default:"0"`
	RocketDuration       int     `json:"rocket_duration,omitempty" default:"0"`
	DevNativeTokenAmount float64 `json:"dev_native_token_amount,omitempty" default:"0"`
	DevTokenBurnAmount   float64 `json:"dev_token_burn_amount,omitempty" default:"0"`
	DevTokenBurnRatio    float64 `json:"dev_token_burn_ratio,omitempty" default:"0"`
	LpBurnPercentage     float64 `json:"burn_percentage,omitempty" default:"0"`
	Top10HolderRate      float64 `json:"top_10_percentage,omitempty" default:"0"`
	CommentCount         int     `json:"comment_count,omitempty" default:"0"`
	LastReply            int64   `json:"last_reply,omitempty" default:"0"`
	IsBurnedLp           bool    `json:"is_burned_lp,omitempty" default:"false"`
	IsDexAd              bool    `json:"is_dex_ad,omitempty" default:"false"`

	// 状态标志
	IsMedia           bool `json:"is_media,omitempty" default:"false"`
	DexscrAd          bool `json:"dexscr_ad,omitempty" default:"false"`
	DexscrUpdateLink  bool `json:"dexscr_update_link,omitempty" default:"false"`
	CtoFlag           bool `json:"cto_flag,omitempty" default:"false"`
	CreatorClose      bool `json:"creator_close,omitempty" default:"false"`
	TwitterChangeFlag bool `json:"twitter_change_flag,omitempty" default:"false"`
	TokenFlags        int  `json:"token_flags,omitempty" default:"0"`
	TotalHolders      bool `json:"total_holders,omitempty" default:"false"`
	MintAuthority     bool `json:"mint_authority,omitempty" default:"false"`
	FreezeAuthority   bool `json:"freeze_authority,omitempty" default:"false"`
}

// TokenTransactionResponse API响应用的代币交易信息
type TokenTransactionResponse struct {
	// 基础信息
	TokenAddress        string                 `json:"token_address"`
	UserAddress         string                 `json:"creator"`
	ExtInfo             map[string]interface{} `json:"ext_info"`
	ChainType           int                    `json:"chain_type"`
	PlatformType        int                    `json:"platform_type"`
	CreatedPlatformType int                    `json:"created_platform_type"`
	Decimals            int                    `json:"decimals"`
	TokenSupply         int64                  `json:"total_supply"`
	Progress            float64                `json:"progress"`
	IsBurnedLp          bool                   `json:"is_burned_lp"`

	// 时间戳
	TransactionTime int64 `json:"last_trade_timestamp"`
	TokenCreateTime int64 `json:"created_timestamp"`
	LastReply       int64 `json:"last_reply"`
	UpdatedAt       int64 `json:"updated_at"`
	OpenTimestamp   int64 `json:"open_timestamp"`

	// 价格和储备
	Price                 float64 `json:"price"`
	NativePrice           float64 `json:"native_price"`
	RealNativeReserves    float64 `json:"real_native_reserves"`
	VirtualNativeReserves float64 `json:"virtual_native_reserves"`
	VirtualTokenReserves  float64 `json:"virtual_token_reserves"`
	InitialNativeReserve  float64 `json:"initial_native_reserve"`
	InitialTokenReserve   float64 `json:"initial_token_reserve"`
	NativeReserveRate     float64 `json:"native_reserve_rate"`

	// 交易统计
	Volume1m      float64 `json:"volume_1m"`
	Volume5m      float64 `json:"volume_5m"`
	Volume1h      float64 `json:"volume_1h"`
	Volume6h      float64 `json:"volume_6h"`
	Volume24h     float64 `json:"volume_24h"`
	Volume        float64 `json:"volume"`
	SwapsCount1m  int64   `json:"swaps_1m"`
	SwapsCount5m  int64   `json:"swaps_5m"`
	SwapsCount1h  int64   `json:"swaps_1h"`
	SwapsCount6h  int64   `json:"swaps_6h"`
	SwapsCount24h int64   `json:"swaps_24h"`
	Swaps         int64   `json:"swaps"`
	BuyCount1m    int64   `json:"buy_count_1m"`
	BuyCount5m    int64   `json:"buy_count_5m"`
	BuyCount1h    int64   `json:"buy_count_1h"`
	BuyCount6h    int64   `json:"buy_count_6h"`
	BuyCount24h   int64   `json:"buy_count_24h"`
	SellCount1m   int64   `json:"sell_count_1m"`
	SellCount5m   int64   `json:"sell_count_5m"`
	SellCount1h   int64   `json:"sell_count_1h"`
	SellCount6h   int64   `json:"sell_count_6h"`
	SellCount24h  int64   `json:"sell_count_24h"`
	Buys          int64   `json:"buys"`
	Sells         int64   `json:"sells"`

	// 价格变化
	PriceChangePercent    float64 `json:"price_change_percent"`
	PriceChangePercent1m  float64 `json:"price_change_percent1m"`
	PriceChangePercent5m  float64 `json:"price_change_percent5m"`
	PriceChangePercent1h  float64 `json:"price_change_percent1h"`
	PriceChangePercent6h  float64 `json:"price_change_percent6h"`
	PriceChangePercent24h float64 `json:"price_change_percent24h"`

	// 市值相关
	MarketCap         float64 `json:"usd_market_cap"`
	MarketCapChange1m float64 `json:"market_cap_change_1m"`
	MarketCapChange5m float64 `json:"market_cap_change_5m"`

	// 代币状态
	Status               int     `json:"status"`
	DevStatus            int     `json:"creator_token_status"`
	CrownDuration        int     `json:"crown_duration"`
	RocketDuration       int     `json:"rocket_duration"`
	Holder               int     `json:"holder"`
	CommentCount         int     `json:"reply_count"`
	Top10HolderRate      float64 `json:"top_10_holder_rate"`
	DevTokenBurnAmount   float64 `json:"dev_token_burn_amount"`
	DevTokenBurnRatio    float64 `json:"dev_token_burn_ratio"`
	LpBurnPercentage     float64 `json:"lp_burn_percentage"`
	DevNativeTokenAmount float64 `json:"creator_balance"`
	CreatorBalanceRate   float64 `json:"creator_balance_rate"`
	RatTraderAmountRate  float64 `json:"rat_trader_amount_rate"`

	// 标志位
	TotalHolders      bool `json:"total_holders"`
	MintAuthority     bool `json:"mint_authority"`
	FreezeAuthority   bool `json:"freeze_authority"`
	IsMedia           bool `json:"is_media"`
	CtoFlag           bool `json:"cto_flag"`
	DexscrAd          bool `json:"dexscr_ad"`
	Completed         bool `json:"complete"`
	IsComplete        bool `json:"is_complete"`
	CreatorClose      bool `json:"creator_close"`
	DexscrUpdateLink  bool `json:"dexscr_update_link"`
	TwitterChangeFlag bool `json:"twitter_change_flag"`
}

// 定义排序字段常量
const (
	SortFieldCreatedTimestamp   = "created_timestamp"
	SortFieldHolderCount        = "holder_count"
	SortFieldLiquidity          = "liquidity"
	SortFieldReplyCount         = "reply_count"
	SortFieldLastTradeTimestamp = "last_trade_timestamp"
	SortFieldPrice              = "price"
	SortFieldPriceChange1m      = "price_change_percent1m"
	SortFieldPriceChange5m      = "price_change_percent5m"
	SortFieldPriceChange1h      = "price_change_percent1h"
	SortFieldMarketCap          = "usd_market_cap"
	SortFieldSwaps1h            = "swaps_1h"
	SortFieldVolume1h           = "volume_1h"
	SortFieldVolume             = "volume"
	SortFieldSwaps              = "swaps"
	SortFieldChange             = "change"
)

// SortPoolMarketInfos 通用排序方法
func SortPoolMarketInfos(responses []*PoolMarketInfo, field, direction string) {
	sort.Slice(responses, func(i, j int) bool {
		return comparePoolMarketInfos(responses[i], responses[j], field, direction == "asc")
	})
}

// comparePoolMarketInfos 比较两个 PoolMarketInfo 对象
func comparePoolMarketInfos(a, b *PoolMarketInfo, field string, isAsc bool) bool {
	// 空值检查
	if a == nil || b == nil ||
		a.BaseTokenInfo == nil || b.BaseTokenInfo == nil {
		return false
	}

	// 获取比较值
	var valA, valB any
	switch field {
	case SortFieldCreatedTimestamp:
		valA, valB = a.OpenTimestamp, b.OpenTimestamp
	case SortFieldHolderCount:
		valA, valB = a.BaseTokenInfo.Holder, b.BaseTokenInfo.Holder
	case SortFieldLiquidity:
		valA, valB = a.BaseTokenInfo.VirtualNativeReserves, b.BaseTokenInfo.VirtualNativeReserves
	case SortFieldReplyCount:
		valA, valB = a.BaseTokenInfo.CommentCount, b.BaseTokenInfo.CommentCount
	case SortFieldLastTradeTimestamp:
		valA, valB = a.BaseTokenInfo.TransactionTime, b.BaseTokenInfo.TransactionTime
	case SortFieldPrice:
		valA, valB = a.BaseTokenInfo.Price, b.BaseTokenInfo.Price
	case SortFieldPriceChange1m:
		valA, valB = a.BaseTokenInfo.PriceChangePercent1m, b.BaseTokenInfo.PriceChangePercent1m
	case SortFieldPriceChange5m:
		valA, valB = a.BaseTokenInfo.PriceChangePercent5m, b.BaseTokenInfo.PriceChangePercent5m
	case SortFieldPriceChange1h:
		valA, valB = a.BaseTokenInfo.PriceChangePercent1h, b.BaseTokenInfo.PriceChangePercent1h
	case SortFieldMarketCap:
		valA, valB = a.BaseTokenInfo.MarketCap, b.BaseTokenInfo.MarketCap
	case SortFieldSwaps1h:
		valA, valB = a.BaseTokenInfo.Swaps, b.BaseTokenInfo.Swaps
	case SortFieldVolume1h:
		valA, valB = a.BaseTokenInfo.Volume, b.BaseTokenInfo.Volume
	case SortFieldVolume:
		valA, valB = a.BaseTokenInfo.Volume, b.BaseTokenInfo.Volume
	case SortFieldSwaps:
		valA, valB = a.BaseTokenInfo.Swaps, b.BaseTokenInfo.Swaps
	case SortFieldChange:
		valA, valB = a.BaseTokenInfo.PriceChangePercent, b.BaseTokenInfo.PriceChangePercent
	default:
		return false
	}

	return compareValues(valA, valB, isAsc)
}

// compareValues 比较两个值
func compareValues(valA, valB any, isAsc bool) bool {
	switch v := valA.(type) {
	case int64:
		if isAsc {
			return v < valB.(int64)
		}
		return v > valB.(int64)
	case int:
		if isAsc {
			return v < valB.(int)
		}
		return v > valB.(int)
	case float64:
		if isAsc {
			return v < valB.(float64)
		}
		return v > valB.(float64)
	default:
		return false
	}
}

func ConvertAndSortSolPumpNewCreateTransactions(transactions []*TokenTransaction, field, direction string) []*TokenTransactionResponse {

	responses := make([]*TokenTransactionResponse, len(transactions))
	indices := make([]int, len(transactions))

	var wg sync.WaitGroup
	wg.Add(len(transactions))

	for i, transaction := range transactions {
		go func(i int, transaction *TokenTransaction) {
			defer wg.Done()
			if transaction == nil {
				return
			}

			parsedTransactionTime, err := time.Parse(time.RFC3339, transaction.LatestTransactionTime)
			if err != nil {
				fmt.Println("Error parsing transaction time:", err)
				// 设置一个默认时间
				parsedTransactionTime = time.Now()
			}

			parsedTokenCreateTime, err := time.Parse(time.RFC3339, transaction.TokenCreateTime)
			if err != nil {
				fmt.Println("Error parsing token create time:", err)
				parsedTokenCreateTime = time.Now()
			}

			extInfoMap := ConvertExtInfoToMap(transaction.ExtInfo)

			// 先判断 再转换 string to float64
			virtualTokenReserves, err := strconv.ParseFloat(transaction.VirtualTokenReserves, 64)
			if err != nil {
				virtualTokenReserves = 0
			}

			virtualNativeReserves, err := strconv.ParseFloat(transaction.VirtualNativeReserves, 64)
			if err != nil {
				virtualNativeReserves = 0
			}

			tokenSupply, err := strconv.ParseInt(transaction.TokenSupply, 10, 64)
			if err != nil {
				tokenSupply = 0
			}

			creatorBalanceRate, err := strconv.ParseFloat(transaction.CreatorBalanceRate, 64)
			if err != nil {
				creatorBalanceRate = 0
			}

			// 先判断 再转换 string to float64 再除于精度
			realNativeReserves, err := strconv.ParseFloat(transaction.RealNativeReserves, 64)
			if err != nil {
				realNativeReserves = 0
			}
			realNativeReserves = realNativeReserves / math.Pow(10, float64(SolDecimals))

			responses[i] = &TokenTransactionResponse{
				ExtInfo:               extInfoMap,
				TokenAddress:          transaction.TokenAddress,
				UserAddress:           transaction.UserAddress,
				Price:                 transaction.Price,
				NativePrice:           transaction.NativePrice,
				TransactionTime:       parsedTransactionTime.Unix(),
				ChainType:             transaction.ChainType,
				PlatformType:          transaction.PlatformType,
				Progress:              transaction.Progress,
				VirtualTokenReserves:  virtualTokenReserves,
				VirtualNativeReserves: virtualNativeReserves,
				InitialNativeReserve:  transaction.InitialNativeReserve,
				InitialTokenReserve:   transaction.InitialTokenReserve,
				TokenCreateTime:       parsedTokenCreateTime.Unix(),
				Decimals:              transaction.Decimals, // 待修改
				CreatedPlatformType:   transaction.CreatedPlatformType,
				DevNativeTokenAmount:  transaction.DevNativeTokenAmount,
				Holder:                transaction.Holder,
				CommentCount:          transaction.CommentCount,
				MarketCap:             transaction.MarketCap,
				RealNativeReserves:    realNativeReserves,
				TokenSupply:           tokenSupply,
				DevTokenBurnAmount:    transaction.DevTokenBurnAmount,
				DevTokenBurnRatio:     transaction.DevTokenBurnRatio,
				LpBurnPercentage:      transaction.LpBurnPercentage,
				DevStatus:             transaction.DevStatus,
				SwapsCount1h:          transaction.SwapsCount1h,
				Volume1h:              transaction.Volume1h,
				IsMedia:               transaction.IsMedia,
				Top10HolderRate:       transaction.Top10HolderRate,
				CreatorBalanceRate:    creatorBalanceRate,
				CtoFlag:               transaction.CtoFlag,
				DexscrAd:              transaction.DexscrAd,
				DexscrUpdateLink:      transaction.DexscrUpdateLink,
				TwitterChangeFlag:     transaction.TwitterChangeFlag,
				IsComplete:            transaction.IsComplete,
			}

			indices[i] = i
		}(i, transaction)
	}
	wg.Wait()

	sort.Slice(indices, func(i, j int) bool {
		switch field {
		case "swaps_1h":
			if direction == "asc" {
				return responses[indices[i]].SwapsCount1h < responses[indices[j]].SwapsCount1h
			} else {
				return responses[indices[i]].SwapsCount1h > responses[indices[j]].SwapsCount1h
			}
		case "volume_1h":
			if direction == "asc" {
				return responses[indices[i]].Volume1h < responses[indices[j]].Volume1h
			} else {
				return responses[indices[i]].Volume1h > responses[indices[j]].Volume1h
			}
		case "created_timestamp":
			if direction == "asc" {
				return responses[indices[i]].TokenCreateTime < responses[indices[j]].TokenCreateTime
			} else {
				return responses[indices[i]].TokenCreateTime > responses[indices[j]].TokenCreateTime
			}
		case "progress":
			if direction == "asc" {
				return responses[indices[i]].Progress < responses[indices[j]].Progress
			} else {
				return responses[indices[i]].Progress > responses[indices[j]].Progress
			}
		case "liquidity":
			if direction == "asc" {
				return responses[indices[i]].VirtualNativeReserves < responses[indices[j]].VirtualNativeReserves
			} else {
				return responses[indices[i]].VirtualNativeReserves > responses[indices[j]].VirtualNativeReserves
			}
		case "creator_balance":
			if direction == "asc" {
				return responses[indices[i]].RealNativeReserves < responses[indices[j]].RealNativeReserves
			} else {
				return responses[indices[i]].RealNativeReserves > responses[indices[j]].RealNativeReserves
			}
		case "holder_count":
			if direction == "asc" {
				return responses[indices[i]].Holder < responses[indices[j]].Holder
			} else {
				return responses[indices[i]].Holder > responses[indices[j]].Holder
			}
		case "reply_count":
			if direction == "asc" {
				return responses[indices[i]].CommentCount < responses[indices[j]].CommentCount
			} else {
				return responses[indices[i]].CommentCount > responses[indices[j]].CommentCount
			}
		case "last_trade_timestamp":
			if direction == "asc" {
				return responses[indices[i]].TransactionTime < responses[indices[j]].TransactionTime
			} else {
				return responses[indices[i]].TransactionTime > responses[indices[j]].TransactionTime
			}
		case "usd_market_cap":
			if direction == "asc" {
				return responses[indices[i]].MarketCap < responses[indices[j]].MarketCap
			} else {
				return responses[indices[i]].MarketCap > responses[indices[j]].MarketCap
			}
		default:
			return indices[i] < indices[j]
		}
	})

	if field == "" || direction == "" {
		return responses
	}

	sortedResponses := make([]*TokenTransactionResponse, len(responses))
	for i, index := range indices {
		sortedResponses[i] = responses[index]
	}

	return sortedResponses
}

func ConvertAndSortSolPumpCompletingTransactions(transactions []*TokenTransaction, field, direction string) []*TokenTransactionResponse {
	// 创建一个切片来存储响应
	responses := make([]*TokenTransactionResponse, len(transactions))
	indices := make([]int, len(transactions))

	var wg sync.WaitGroup
	wg.Add(len(transactions))

	// 遍历 transactions，进行赋值
	for i, transaction := range transactions {

		go func(i int, transaction *TokenTransaction) {
			defer wg.Done()

			parsedTransactionTime, err := time.Parse(time.RFC3339, transaction.LatestTransactionTime)
			if err != nil {
				fmt.Println("Error parsing transaction time:", err)
				parsedTransactionTime = time.Now()
			}

			parsedTokenCreateTime, err := time.Parse(time.RFC3339, transaction.TokenCreateTime)
			if err != nil {
				fmt.Println("Error parsing token create time:", err)
				parsedTokenCreateTime = time.Now()
			}

			extInfoMap := ConvertExtInfoToMap(transaction.ExtInfo)

			// 先判断 再转换 string to float64
			virtualTokenReserves, err := strconv.ParseFloat(transaction.VirtualTokenReserves, 64)
			if err != nil {
				virtualTokenReserves = 0
			}

			virtualNativeReserves, err := strconv.ParseFloat(transaction.VirtualNativeReserves, 64)
			if err != nil {
				virtualNativeReserves = 0
			}

			tokenSupply, err := strconv.ParseInt(transaction.TokenSupply, 10, 64)
			if err != nil {
				tokenSupply = 0
			}

			creatorBalanceRate, err := strconv.ParseFloat(transaction.CreatorBalanceRate, 64)
			if err != nil {
				creatorBalanceRate = 0
			}

			responses[i] = &TokenTransactionResponse{
				ExtInfo:               extInfoMap,
				TokenAddress:          transaction.TokenAddress,
				UserAddress:           transaction.UserAddress,
				Price:                 transaction.Price,
				NativePrice:           transaction.NativePrice,
				TransactionTime:       parsedTransactionTime.Unix(),
				ChainType:             transaction.ChainType,
				PlatformType:          transaction.PlatformType,
				Progress:              transaction.Progress,
				VirtualTokenReserves:  virtualTokenReserves,
				VirtualNativeReserves: virtualNativeReserves,
				InitialNativeReserve:  transaction.InitialNativeReserve,
				InitialTokenReserve:   transaction.InitialTokenReserve,
				Decimals:              transaction.Decimals,
				TokenCreateTime:       parsedTokenCreateTime.Unix(),
				CrownDuration:         transaction.RocketDuration,
				RocketDuration:        transaction.CrownDuration,
				DevNativeTokenAmount:  transaction.DevNativeTokenAmount,
				Holder:                transaction.Holder,
				CommentCount:          transaction.CommentCount,
				CreatedPlatformType:   transaction.CreatedPlatformType,
				MarketCap:             transaction.MarketCap,
				TokenSupply:           tokenSupply,
				DevTokenBurnAmount:    transaction.DevTokenBurnAmount,
				DevTokenBurnRatio:     transaction.DevTokenBurnRatio,
				LpBurnPercentage:      transaction.LpBurnPercentage,
				IsBurnedLp:            transaction.IsBurnedLp,
				DevStatus:             transaction.DevStatus,
				Volume1h:              transaction.Volume1h,
				SwapsCount1h:          transaction.SwapsCount1h,
				IsMedia:               transaction.IsMedia,
				Top10HolderRate:       transaction.Top10HolderRate,
				CreatorBalanceRate:    creatorBalanceRate,
				CtoFlag:               transaction.CtoFlag,
				DexscrAd:              transaction.DexscrAd,
				DexscrUpdateLink:      transaction.DexscrUpdateLink,
				TwitterChangeFlag:     transaction.TwitterChangeFlag,
				IsComplete:            transaction.IsComplete,
			}
			indices[i] = i
		}(i, transaction)
	}
	wg.Wait()

	sort.Slice(indices, func(i, j int) bool {
		switch field {
		case "swaps_1h":
			if direction == "asc" {
				return responses[indices[i]].SwapsCount1h < responses[indices[j]].SwapsCount1h
			} else {
				return responses[indices[i]].SwapsCount1h > responses[indices[j]].SwapsCount1h
			}
		case "volume_1h":
			if direction == "asc" {
				return responses[indices[i]].Volume1h < responses[indices[j]].Volume1h
			} else {
				return responses[indices[i]].Volume1h > responses[indices[j]].Volume1h
			}
		case "created_timestamp":
			if direction == "asc" {
				return responses[indices[i]].TokenCreateTime < responses[indices[j]].TokenCreateTime
			} else {
				return responses[indices[i]].TokenCreateTime > responses[indices[j]].TokenCreateTime
			}
		case "progress":
			if direction == "asc" {
				return responses[indices[i]].Progress < responses[indices[j]].Progress
			} else {
				return responses[indices[i]].Progress > responses[indices[j]].Progress
			}
		case "liquidity":
			if direction == "asc" {
				return responses[indices[i]].VirtualNativeReserves < responses[indices[j]].VirtualNativeReserves
			} else {
				return responses[indices[i]].VirtualNativeReserves > responses[indices[j]].VirtualNativeReserves
			}
		case "holder_count":
			if direction == "asc" {
				return responses[indices[i]].Holder < responses[indices[j]].Holder
			} else {
				return responses[indices[i]].Holder > responses[indices[j]].Holder
			}
		case "reply_count":
			if direction == "asc" {
				return responses[indices[i]].CommentCount < responses[indices[j]].CommentCount
			} else {
				return responses[indices[i]].CommentCount > responses[indices[j]].CommentCount
			}
		case "last_trade_timestamp":
			if direction == "asc" {
				return responses[indices[i]].TransactionTime < responses[indices[j]].TransactionTime
			} else {
				return responses[indices[i]].TransactionTime > responses[indices[j]].TransactionTime
			}
		case "koth_duration":
			if direction == "asc" {
				return responses[indices[i]].CrownDuration > responses[indices[j]].CrownDuration
			} else {
				return responses[indices[i]].CrownDuration < responses[indices[j]].CrownDuration
			}
		case "time_since_koth":
			if direction == "asc" {
				return responses[indices[i]].RocketDuration > responses[indices[j]].RocketDuration
			} else {
				return responses[indices[i]].RocketDuration < responses[indices[j]].RocketDuration
			}
		case "usd_market_cap":
			if direction == "asc" {
				return responses[indices[i]].MarketCap < responses[indices[j]].MarketCap
			} else {
				return responses[indices[i]].MarketCap > responses[indices[j]].MarketCap
			}
		default:
			// 默认根据progress排序
			return responses[indices[i]].Progress > responses[indices[j]].Progress
		}
	})

	// 创建一个新的切片以返回排序后的结果
	sortedResponses := make([]*TokenTransactionResponse, len(responses))
	for i, index := range indices {
		sortedResponses[i] = responses[index]
	}

	return sortedResponses
}

func ConvertAndSortSolPumpSoaringTransactions(transactions []*TokenTransaction, field, direction string) []*TokenTransactionResponse {

	responses := make([]*TokenTransactionResponse, len(transactions))
	indices := make([]int, len(transactions))

	var wg sync.WaitGroup
	wg.Add(len(transactions))

	for i, transaction := range transactions {
		go func(i int, transaction *TokenTransaction) {
			defer wg.Done()
			if transaction == nil {
				return
			}

			parsedTransactionTime, err := time.Parse(time.RFC3339, transaction.LatestTransactionTime)
			if err != nil {
				fmt.Println("Error parsing transaction time:", err)
				parsedTransactionTime = time.Now()
			}

			parsedTokenCreateTime, err := time.Parse(time.RFC3339, transaction.TokenCreateTime)
			if err != nil {
				fmt.Println("Error parsing token create time:", err)
				parsedTokenCreateTime = time.Now()
			}

			extInfoMap := ConvertExtInfoToMap(transaction.ExtInfo)
			// 先判断 再转换 string to float64
			virtualTokenReserves, err := strconv.ParseFloat(transaction.VirtualTokenReserves, 64)
			if err != nil {
				virtualTokenReserves = 0
			}

			virtualNativeReserves, err := strconv.ParseFloat(transaction.VirtualNativeReserves, 64)
			if err != nil {
				virtualNativeReserves = 0
			}

			tokenSupply, err := strconv.ParseInt(transaction.TokenSupply, 10, 64)
			if err != nil {
				tokenSupply = 0
			}

			creatorBalanceRate, err := strconv.ParseFloat(transaction.CreatorBalanceRate, 64)
			if err != nil {
				creatorBalanceRate = 0
			}

			responses[i] = &TokenTransactionResponse{
				ExtInfo:               extInfoMap,
				TokenAddress:          transaction.TokenAddress,
				UserAddress:           transaction.UserAddress,
				Price:                 transaction.Price,
				NativePrice:           transaction.NativePrice,
				TransactionTime:       parsedTransactionTime.Unix(),
				ChainType:             transaction.ChainType,
				PlatformType:          transaction.PlatformType,
				Progress:              transaction.Progress,
				VirtualTokenReserves:  virtualTokenReserves,
				VirtualNativeReserves: virtualNativeReserves,
				InitialNativeReserve:  transaction.InitialNativeReserve,
				InitialTokenReserve:   transaction.InitialTokenReserve,
				TokenCreateTime:       parsedTokenCreateTime.Unix(),
				DevNativeTokenAmount:  transaction.DevNativeTokenAmount,
				Decimals:              transaction.Decimals,
				CreatedPlatformType:   transaction.CreatedPlatformType,
				// DevTokenAmount:        transaction.DevTokenAmount,
				Holder:               transaction.Holder,
				CommentCount:         transaction.CommentCount,
				MarketCap:            transaction.MarketCap,
				TokenSupply:          tokenSupply,
				DevTokenBurnAmount:   transaction.DevTokenBurnAmount,
				DevTokenBurnRatio:    transaction.DevTokenBurnRatio,
				LpBurnPercentage:     transaction.LpBurnPercentage,
				IsBurnedLp:           transaction.IsBurnedLp,
				DevStatus:            transaction.DevStatus,
				Volume1h:             transaction.Volume,
				SwapsCount1h:         transaction.Swaps,
				Swaps:                transaction.Swaps,
				Volume:               transaction.Volume,
				Buys:                 transaction.Buys,
				Sells:                transaction.Sells,
				PriceChangePercent:   transaction.PriceChangePercent,
				MarketCapChange1m:    transaction.MarketCapChange1m,
				MarketCapChange5m:    transaction.MarketCapChange5m,
				PriceChangePercent1m: transaction.PriceChangePercent1m,
				PriceChangePercent5m: transaction.PriceChangePercent5m,
				IsMedia:              transaction.IsMedia,
				Top10HolderRate:      transaction.Top10HolderRate,
				CreatorBalanceRate:   creatorBalanceRate,
				CtoFlag:              transaction.CtoFlag,
				DexscrAd:             transaction.DexscrAd,
				DexscrUpdateLink:     transaction.DexscrUpdateLink,
				TwitterChangeFlag:    transaction.TwitterChangeFlag,
				IsComplete:           transaction.IsComplete,
			}

			indices[i] = i
		}(i, transaction)
	}
	wg.Wait()

	sort.Slice(indices, func(i, j int) bool {
		switch field {
		case "created_timestamp":
			if direction == "asc" {
				return responses[indices[i]].TokenCreateTime < responses[indices[j]].TokenCreateTime
			} else {
				return responses[indices[i]].TokenCreateTime > responses[indices[j]].TokenCreateTime
			}
		case "progress":
			if direction == "asc" {
				return responses[indices[i]].Progress < responses[indices[j]].Progress
			} else {
				return responses[indices[i]].Progress > responses[indices[j]].Progress
			}
		case "holder_count":
			if direction == "asc" {
				return responses[indices[i]].Holder < responses[indices[j]].Holder
			} else {
				return responses[indices[i]].Holder > responses[indices[j]].Holder
			}
		case "swaps_1h":
			if direction == "asc" {
				return responses[indices[i]].SwapsCount1h < responses[indices[j]].SwapsCount1h
			} else {
				return responses[indices[i]].SwapsCount1h > responses[indices[j]].SwapsCount1h
			}
		case "volume_1h":
			if direction == "asc" {
				return responses[indices[i]].Volume1h < responses[indices[j]].Volume1h
			} else {
				return responses[indices[i]].Volume1h > responses[indices[j]].Volume1h
			}
		case "liquidity":
			if direction == "asc" {
				return responses[indices[i]].VirtualNativeReserves < responses[indices[j]].VirtualNativeReserves
			} else {
				return responses[indices[i]].VirtualNativeReserves > responses[indices[j]].VirtualNativeReserves
			}
		case "reply_count":
			if direction == "asc" {
				return responses[indices[i]].CommentCount < responses[indices[j]].CommentCount
			} else {
				return responses[indices[i]].CommentCount > responses[indices[j]].CommentCount
			}
		case "last_trade_timestamp":
			if direction == "asc" {
				return responses[indices[i]].TransactionTime < responses[indices[j]].TransactionTime
			} else {
				return responses[indices[i]].TransactionTime > responses[indices[j]].TransactionTime
			}
		case "market_cap_1m":
			if direction == "asc" {
				return responses[indices[i]].MarketCapChange1m < responses[indices[j]].MarketCapChange1m
			} else {
				return responses[indices[i]].MarketCapChange1m > responses[indices[j]].MarketCapChange1m
			}
		case "market_cap_5m":
			if direction == "asc" {
				return responses[indices[i]].MarketCapChange5m < responses[indices[j]].MarketCapChange5m
			} else {
				return responses[indices[i]].MarketCapChange5m > responses[indices[j]].MarketCapChange5m
			}
		case "usd_market_cap":
			if direction == "asc" {
				return responses[indices[i]].MarketCap < responses[indices[j]].MarketCap
			} else {
				return responses[indices[i]].MarketCap > responses[indices[j]].MarketCap
			}
		case "swaps":
			if direction == "asc" {
				return responses[indices[i]].Swaps < responses[indices[j]].Swaps
			} else {
				return responses[indices[i]].Swaps > responses[indices[j]].Swaps
			}
		case "change":
			if direction == "asc" {
				return responses[indices[i]].PriceChangePercent < responses[indices[j]].PriceChangePercent
			} else {
				return responses[indices[i]].PriceChangePercent > responses[indices[j]].PriceChangePercent
			}
		case "price":
			if direction == "asc" {
				return responses[indices[i]].Price < responses[indices[j]].Price
			} else {
				return responses[indices[i]].Price > responses[indices[j]].Price
			}
		default:
			return indices[i] < indices[j]
		}
	})

	if field == "" || direction == "" {
		return responses
	}

	sortedResponses := make([]*TokenTransactionResponse, len(responses))
	for i, index := range indices {
		sortedResponses[i] = responses[index]
	}

	return sortedResponses
}

func ConverSolRaydiumTransactions(transactions []*TokenTransaction, field, direction string) []*PoolMarketInfo {

	responses := make([]*PoolMarketInfo, len(transactions))
	var wg sync.WaitGroup
	wg.Add(len(transactions))

	for i, transaction := range transactions {

		i, transaction := i, transaction
		go func() {
			defer wg.Done()
			if transaction == nil {
				return
			}

			parsedTokenCreateTime, err := time.Parse(time.RFC3339, transaction.OpenTimestamp)
			if err != nil {
				fmt.Println("Error parsing token create time:", err)
				return
			}

			// 转换为string
			initialNativeReserve := strconv.FormatFloat(transaction.InitialNativeReserve, 'f', -1, 64)
			initialTokenReserve := strconv.FormatFloat(transaction.InitialTokenReserve, 'f', -1, 64)

			// Construct PoolMarketInfo from TokenTransaction
			responses[i] = &PoolMarketInfo{
				PoolAddress:         transaction.PoolAddress,
				BaseAddress:         transaction.BaseAddress,
				QuoteAddress:        transaction.QuoteAddress,
				QuoteReserve:        transaction.QuoteReserve,
				InitialLiquidity:    initialNativeReserve,
				InitialQuoteReserve: initialTokenReserve,
				Creator:             transaction.Creator,
				PoolTypeStr:         transaction.PoolTypeStr,
				PoolType:            transaction.PoolType,
				QuoteSymbol:         transaction.QuoteSymbol,
				BaseTokenInfo:       ConvertSingleSolSwapTransaction(transaction),
				OpenTimestamp:       parsedTokenCreateTime.Unix(),
				Launchpad:           transaction.Launchpad,
			}
		}()
	}
	wg.Wait()

	// 使用通用排序方法
	if field != "" && direction != "" {
		SortPoolMarketInfos(responses, field, direction)
	}

	return responses
}

func ConvertSingleSolSwapTransaction(transaction *TokenTransaction) *TokenTransactionResponse {
	if transaction == nil {
		return nil
	}
	parsedTransactionTime, err := time.Parse(time.RFC3339, transaction.LatestTransactionTime)
	if err != nil {
		fmt.Println("Error parsing transaction time:", err)
		return nil
	}

	parsedTokenCreateTime, err := time.Parse(time.RFC3339, transaction.OpenTimestamp)
	if err != nil {
		fmt.Println("Error parsing token create time:", err)
		return nil
	}

	extInfoMap := ConvertExtInfoToMap(transaction.ExtInfo)

	// 先判断 再转换 string to float64
	virtualTokenReserves, err := strconv.ParseFloat(transaction.VirtualTokenReserves, 64)
	if err != nil {
		virtualTokenReserves = 0
	}

	virtualNativeReserves, err := strconv.ParseFloat(transaction.VirtualNativeReserves, 64)
	if err != nil {
		virtualNativeReserves = 0
	}

	tokenSupply, err := strconv.ParseInt(transaction.TokenSupply, 10, 64)
	if err != nil {
		tokenSupply = 0
	}

	creatorBalanceRate, err := strconv.ParseFloat(transaction.CreatorBalanceRate, 64)
	if err != nil {
		creatorBalanceRate = 0
	}

	return &TokenTransactionResponse{
		ExtInfo:               extInfoMap,
		TokenAddress:          transaction.TokenAddress,
		UserAddress:           transaction.UserAddress,
		Price:                 transaction.Price,
		NativePrice:           transaction.NativePrice,
		TransactionTime:       parsedTransactionTime.Unix(),
		ChainType:             transaction.ChainType,
		PlatformType:          transaction.PlatformType,
		Progress:              transaction.Progress,
		VirtualTokenReserves:  virtualTokenReserves,
		VirtualNativeReserves: virtualNativeReserves,
		InitialNativeReserve:  transaction.InitialNativeReserve,
		InitialTokenReserve:   transaction.InitialTokenReserve,
		TokenCreateTime:       parsedTokenCreateTime.Unix(),
		DevNativeTokenAmount:  transaction.DevNativeTokenAmount,
		Holder:                transaction.Holder,
		CommentCount:          transaction.CommentCount,
		MarketCap:             transaction.MarketCap,
		TokenSupply:           tokenSupply,
		Decimals:              transaction.Decimals,
		CreatedPlatformType:   transaction.CreatedPlatformType,
		DevTokenBurnAmount:    transaction.DevTokenBurnAmount,
		DevTokenBurnRatio:     transaction.DevTokenBurnRatio,
		LpBurnPercentage:      transaction.LpBurnPercentage,
		IsBurnedLp:            transaction.IsBurnedLp,
		DevStatus:             transaction.DevStatus,
		Volume1h:              transaction.Volume1h,
		SwapsCount1h:          transaction.SwapsCount1h,
		Swaps:                 transaction.Swaps,
		Volume:                transaction.Volume,
		Buys:                  transaction.Buys,
		Sells:                 transaction.Sells,
		PriceChangePercent:    transaction.PriceChangePercent,
		PriceChangePercent1m:  transaction.PriceChangePercent1m,
		PriceChangePercent5m:  transaction.PriceChangePercent5m,
		PriceChangePercent1h:  transaction.PriceChangePercent1h,
		IsMedia:               transaction.IsMedia,
		Top10HolderRate:       transaction.Top10HolderRate,
		CreatorBalanceRate:    creatorBalanceRate,
		CtoFlag:               transaction.CtoFlag,
		DexscrAd:              transaction.DexscrAd,
		DexscrUpdateLink:      transaction.DexscrUpdateLink,
		TwitterChangeFlag:     transaction.TwitterChangeFlag,
		IsComplete:            transaction.IsComplete,
	}
}

func ConvertExtInfoToMap(extInfoStr string) map[string]interface{} {
	// 1. 空值检查
	if extInfoStr == "" || extInfoStr == "{}" {
		return map[string]interface{}{
			"name":        "Unknown",
			"symbol":      "Unknown",
			"description": "No description",
			"image":       "",
			"showName":    false,
			"createdOn":   "",
			"twitter":     "",
			"website":     "",
			"telegram":    "",
		}
	}

	// 2. 尝试解析
	var extInfoMap map[string]interface{}
	if err := json.Unmarshal([]byte(extInfoStr), &extInfoMap); err != nil {
		log.Printf("Error decoding ExtInfo: %v, data: %s", err, extInfoStr)
		return map[string]interface{}{
			"name":        "Unknown",
			"symbol":      "Unknown",
			"description": "No description",
			"image":       "",
			"showName":    false,
			"createdOn":   "",
			"twitter":     "",
			"website":     "",
			"telegram":    "",
		}
	}

	// 3. 如果解析成功，直接返回
	return extInfoMap
}
