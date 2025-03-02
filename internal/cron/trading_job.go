package cron

import (
	"my-token-ai-be/internal/constants"
	"my-token-ai-be/internal/es"
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/redis"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/service"
	"time"
)

func Trading6hJob() {
	var safety request.SolRankRequest
	safety.Time = "6h"
	safety.Limit = 50

	queryJSON, err := es.SolSwapQuery(&safety)
	if err != nil {
		util.Log().Error("failed to get swap rank data for 6h")
		return
	}

	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if err != nil || result == nil {
		util.Log().Error("failed to get swap rank data for 6h")
		return
	}

	transactions, err := service.ProcessAggregationResult(result, safety.Filters, &safety)
	if err != nil {
		util.Log().Error("failed to process swap rank data for 6h")
		return
	}

	redis.Set(constants.RedisKeyTradingPool6hMarketInfo, transactions, 1*time.Hour)

}

func Trading24hJob() {
	var safety request.SolRankRequest
	safety.Time = "24h"
	safety.Limit = 50

	queryJSON, err := es.SolSwapQuery(&safety)
	if err != nil {
		util.Log().Error("failed to get swap rank data for 24h")
		return
	}

	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if err != nil || result == nil {
		util.Log().Error("failed to get swap rank data for 24h")
		return
	}

	transactions, err := service.ProcessAggregationResult(result, safety.Filters, &safety)
	if err != nil {
		util.Log().Error("failed to process swap rank data for 24h")
		return
	}

	redis.Set(constants.RedisKeyTradingPool24hMarketInfo, transactions, 1*time.Hour)
}
