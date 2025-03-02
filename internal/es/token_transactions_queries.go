package es

import (
	"fmt"
	"my-token-ai-be/internal/request"
)

func NewPairRanksQuery(req *request.SolRankRequest) (string, error) {

	builder := NewESQueryBuilder().
		SetSize(0)

	// 只拿48小时的交易
	builder.AddRangeFilter("transaction_time", "gte", fmt.Sprintf("now-%s", req.Time))
	if req.DexScreenerSpent {
		builder.AddRangeFilter("block_time", "gte", fmt.Sprintf("now-%s", "24h"))
	} else {
		builder.AddRangeFilter("block_time", "gte", fmt.Sprintf("now-%s", "2h"))
	}

	// 设置基础聚合
	builder.AddTermsAggregation("unique_tokens", "token_address.keyword", req.Limit).
		SetTermsOrder("max_block_time", "desc")

	builder.AddTopHitsAggregation("latest_transaction", 1, "transaction_time", "desc", "unique_tokens")

	// 添加其他子聚合到 unique_tokens
	builder.
		AddMaxAggregation("max_block_time", "block_time", "unique_tokens").
		AddSwapCountAggregation("buys", "unique_tokens", true, req.Time).
		AddSwapCountAggregation("sells", "unique_tokens", false, req.Time).
		AddHolderBalanceAggregation("holder_count", "unique_tokens").
		AddVolumeAggregation("volume", "unique_tokens").
		AddMarketCapAggregation("market_cap_1m", "1m", "unique_tokens").
		AddMarketCapAggregation("market_cap_5m", "5m", "unique_tokens").
		AddMarketCapAggregation("market_cap_1h", "1h", "unique_tokens").
		AddLatestMarketCapAggregation("latest_market_cap", "unique_tokens").
		AddTopHoldersAggregation("holders", "unique_tokens").
		AddTopHoldersTotalPercentageAggregation("total_holders_percentage", "unique_tokens")

	if req.DexScreenerSpent {
		builder.AddMust("term", "is_dex_ad", true) // 或者使用 "terms" 如果需要多个值
	}

	// 处理平台过滤
	if len(req.Platforms) > 0 {
		shouldFilters := make([]map[string]interface{}, 0)

		for _, platform := range req.Platforms {
			switch platform {
			case "raydium":
				shouldFilters = append(shouldFilters, map[string]interface{}{
					"terms": map[string]interface{}{
						"platform_type": []int{2},
					},
				})
			case "pump":
				shouldFilters = append(shouldFilters, map[string]interface{}{
					"terms": map[string]interface{}{
						"created_platform_type": []int{1},
					},
				})
			case "mootshot":
				shouldFilters = append(shouldFilters, map[string]interface{}{
					"terms": map[string]interface{}{
						"created_platform_type": []int{2},
					},
				})
			}
		}

		if len(shouldFilters) > 0 {
			builder.AddShould(shouldFilters).SetMinimumShouldMatch(1)
		}
	} else {
		builder.AddShouldTerms("platform_type", []int{2}).
			AddFilterTerms("created_platform_type", []int{0, 1})
	}

	if req.MinCreated != nil {
		builder.AddRangeFilter("block_time", "gte", "now-"+*req.MinCreated)
	}
	if req.MaxCreated != nil {
		builder.AddRangeFilter("block_time", "lte", "now-"+*req.MaxCreated)
	}

	return builder.Build()
}
