package service

import (
	"fmt"
	"math"
	"my-token-ai-be/internal/es"
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

func processMarketCapChanges(tx *response.TokenTransaction, bucket es.Bucket) {

	// 按时间查看市值变化
	if len(bucket.MarketCapTime.LatestTransaction.Hits.Hits) > 0 {
		tx.MarketCapTime = tx.MarketCap - bucket.MarketCapTime.LatestTransaction.Hits.Hits[0].Source.MarketCap
		tx.PriceChangePercent = calculatePriceChangePercent(
			tx.MarketCap,
			bucket.MarketCapTime.LatestTransaction.Hits.Hits[0].Source.MarketCap,
		)
	}

	// 1分钟变化
	if len(bucket.MarketCap1m.LatestTransaction.Hits.Hits) > 0 {
		tx.MarketCapChange1m = tx.MarketCap - bucket.MarketCap1m.LatestTransaction.Hits.Hits[0].Source.MarketCap
		tx.PriceChangePercent1m = calculatePriceChangePercent(
			tx.MarketCap,
			bucket.MarketCap1m.LatestTransaction.Hits.Hits[0].Source.MarketCap,
		)
	}

	// 5分钟变化
	if len(bucket.MarketCap5m.LatestTransaction.Hits.Hits) > 0 {
		tx.MarketCapChange5m = tx.MarketCap - bucket.MarketCap5m.LatestTransaction.Hits.Hits[0].Source.MarketCap
		tx.PriceChangePercent5m = calculatePriceChangePercent(
			tx.MarketCap,
			bucket.MarketCap5m.LatestTransaction.Hits.Hits[0].Source.MarketCap,
		)
	}

	// 1小时变化
	if len(bucket.MarketCap1h.LatestTransaction.Hits.Hits) > 0 {
		tx.MarketCapChange1h = tx.MarketCap - bucket.MarketCap1h.LatestTransaction.Hits.Hits[0].Source.MarketCap
		tx.PriceChangePercent1h = calculatePriceChangePercent(
			tx.MarketCap,
			bucket.MarketCap1h.LatestTransaction.Hits.Hits[0].Source.MarketCap,
		)
	}
}

func calculatePriceChangePercent(latest, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return (latest - previous) / previous
}

func processTokenFlags(tokenTransaction *response.TokenTransaction) {
	var tokenInfo model.TokenInfo
	tokenInfo.SetFlag(tokenTransaction.TokenFlags)

	tokenTransaction.MintAuthority = tokenInfo.HasFlag(model.FLAG_MINT_AUTHORITY)
	tokenTransaction.FreezeAuthority = tokenInfo.HasFlag(model.FLAG_FREEZE_AUTHORITY)
	tokenTransaction.DexscrAd = tokenInfo.HasFlag(model.FLAG_DXSCR_AD)
	tokenTransaction.TwitterChangeFlag = tokenInfo.HasFlag(model.FLAG_TWITTER_CHANGE)
	tokenTransaction.DexscrUpdateLink = tokenInfo.HasFlag(model.FLAG_DEXSCR_UPDATE)
	tokenTransaction.CtoFlag = tokenInfo.HasFlag(model.FLAG_CTO)
}

func safeDiv(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func parseFloat64(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

func processTokenTransaction(tokenTransaction *response.TokenTransaction, solDecimals float64) {
	decimalsMultiplier := math.Pow(10, solDecimals)

	tokenTransaction.InitialNativeReserve = safeDiv(tokenTransaction.InitialNativeReserve, decimalsMultiplier)

	quoteReserve := parseFloat64(tokenTransaction.QuoteReserve)
	tokenTransaction.QuoteReserve = strconv.FormatFloat(
		safeDiv(quoteReserve, decimalsMultiplier),
		'f',
		-1,
		64,
	)

	tokenTransaction.InitialTokenReserve = safeDiv(tokenTransaction.InitialTokenReserve, decimalsMultiplier)

	virtualNativeReserves := parseFloat64(tokenTransaction.VirtualNativeReserves)
	tokenTransaction.VirtualNativeReserves = strconv.FormatFloat(
		safeDiv(virtualNativeReserves, decimalsMultiplier),
		'f',
		-1,
		64,
	)

	tokenTransaction.NativeReserveRate = safeDiv(
		virtualNativeReserves-float64(tokenTransaction.InitialNativeReserve),
		float64(tokenTransaction.InitialNativeReserve),
	)
}

func shouldIncludeTokenTransaction(transaction response.TokenTransaction, filters []string) bool {
	tokenInfo := model.TokenInfo{}
	tokenInfo.TokenFlags = transaction.TokenFlags

	if transaction.ExtInfo == "" {
		return false
	}

	for _, filter := range filters {
		switch filter {
		case "has_social":
			if !transaction.IsMedia {
				return false
			}
		case "creator_hold":
			if transaction.DevStatus != 1 {
				return false
			}
		case "creator_close":
			if transaction.DevStatus != 3 {
				return false
			}
		case "renounced":
			if tokenInfo.HasFlag(model.FLAG_MINT_AUTHORITY) {
				return false
			}
		case "frozen":
			if tokenInfo.HasFlag(model.FLAG_FREEZE_AUTHORITY) {
				return false
			}
		case "burn":
			if transaction.LpBurnPercentage < 50 {
				return false
			}
		case "distributed":
			if transaction.Top10HolderRate > 0.3 {
				return false
			}
		}
	}

	return true
}

// buildSolPumpRankQuery builds the query for SOL pump rank
func buildSolRankQuery(req *request.SolRankRequest) (string, error) {
	queryJSON := ""
	if req.NewCreation != nil && *req.NewCreation {
		query, err := es.NewCreateQuery(req)
		if err != nil {
			return "", err
		}
		queryJSON = query
	}
	if req.Completing != nil && *req.Completing {
		query, err := es.CompletingQuery(req)
		if err != nil {
			return "", err
		}
		queryJSON = query
	}
	if req.Soaring != nil && *req.Soaring {
		query, err := es.SoaringQuery(req)
		if err != nil {
			return "", err
		}
		queryJSON = query
	}
	if req.Completed != nil && *req.Completed {
		query, err := es.CompletedQuery(req)
		if err != nil {
			return "", err
		}
		queryJSON = query
	}
	return queryJSON, nil
}

// checkTransactionConditions 检查交易是否满足所有条件
func checkTransactionConditions(tx response.TokenTransaction, req *request.SolRankRequest) bool {

	// 检查持有者数量范围
	if req.MinHolderCount != nil && tx.Holder < *req.MinHolderCount {
		return false
	}
	if req.MaxHolderCount != nil && tx.Holder > *req.MaxHolderCount {
		return false
	}

	// 检查市值范围
	if req.MinMarketcap != 0 && tx.MarketCap < req.MinMarketcap {
		return false
	}
	if req.MaxMarketcap != 0 && tx.MarketCap > req.MaxMarketcap {
		return false
	}

	// 检查交易量范围
	if req.MinVolume != 0 && tx.Volume < req.MinVolume {
		return false
	}
	if req.MaxVolume != 0 && tx.Volume > req.MaxVolume {
		return false
	}

	if req.MinSwaps != nil && tx.Swaps < *req.MinSwaps {
		return false
	}
	if req.MaxSwaps != nil && tx.Swaps > *req.MaxSwaps {
		return false
	}

	// 检查交易次数范围
	if req.MinSwaps1h != nil && tx.Swaps < *req.MinSwaps1h {
		return false
	}
	if req.MaxSwaps1h != nil && tx.Swaps > *req.MaxSwaps1h {
		return false
	}

	// 如果需要检查quote usd，则获取sol价格
	if req.MinQuoteUsd != nil || req.MaxQuoteUsd != nil {
		solPrice, err := getSolPrice()
		if err != nil {
			return false
		}

		if req.MinQuoteUsd != nil {
			minQuoteUSD := decimal.NewFromFloat(*req.MinQuoteUsd)
			result := minQuoteUSD.Div(solPrice)
			virtualNativeReserves, err := decimal.NewFromString(tx.VirtualNativeReserves)
			if err != nil {
				return false
			}
			if virtualNativeReserves.LessThan(result) {
				return false
			}
		}
		if req.MaxQuoteUsd != nil {
			maxQuoteUSD := decimal.NewFromFloat(*req.MaxQuoteUsd)
			result := maxQuoteUSD.Div(solPrice)
			virtualNativeReserves, err := decimal.NewFromString(tx.VirtualNativeReserves)
			if err != nil {
				return false
			}
			if virtualNativeReserves.GreaterThan(result) {
				return false
			}
		}
	}

	// 所有条件都满足
	return true
}

// ParseTimeRange 将时间范围字符串转换为时间戳
func ParseTimeRange(timeStr string) (int64, error) {
	if timeStr == "" {
		return 0, fmt.Errorf("time string cannot be empty")
	}

	// 获取当前时间戳
	now := time.Now().Unix()

	// 解析时间单位和数值
	length := len(timeStr)
	if length < 2 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	// 获取数值和单位
	value, err := strconv.ParseInt(timeStr[:length-1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid time value: %s", timeStr)
	}
	unit := timeStr[length-1:]

	// 根据单位计算时间戳
	switch unit {
	case "m", "M": // 分钟
		return now - value*60, nil
	case "h", "H": // 小时
		return now - value*60*60, nil
	case "d", "D": // 天
		return now - value*24*60*60, nil
	case "w", "W": // 周
		return now - value*7*24*60*60, nil
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}
}
