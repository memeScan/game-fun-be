package service

import (
	"encoding/json"
	"fmt"
	"my-token-ai-be/internal/es"
	"my-token-ai-be/internal/redis"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
	"net/http"
)

func GetSolPumpRank(req *request.SolRankRequest) response.Response {
	// 1. 构建查询
	queryJSON, err := buildSolRankQuery(req)
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "Failed to build query", err)
	}
	// 2. 执行搜索
	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if err != nil || result == nil {
		status := http.StatusInternalServerError
		msg := "Failed to get pump rank"
		data := []response.PoolMarketInfo{}
		if result == nil {
			status = http.StatusOK
			msg = "No data found"
			data := []response.PoolMarketInfo{}
			return response.BuildResponse(data, status, msg, nil)
		}
		return response.BuildResponse(data, status, msg, err)
	}
	// 3. 处理聚合结果
	transactions, err := ProcessAggregationResult(result, req.Filters, req)
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "Failed to process results", err)
	}

	// 4. 根据不同类型构建响应
	responseBuilder := NewResponseBuilder(transactions, req.OrderBy, req.Direction)

	if req.NewCreation != nil && *req.NewCreation {
		return responseBuilder.BuildNewCreateResponse()
	}
	if req.Completing != nil && *req.Completing {
		return responseBuilder.BuildCompletingResponse()
	}
	if req.Soaring != nil && *req.Soaring {
		return responseBuilder.BuildSoaringResponse()
	}

	return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusInternalServerError, "Invalid request type", nil)
}

func GetSolRaydiumRank(req *request.SolRankRequest) response.Response {
	// 1. 构建查询
	queryJSON, err := buildSolRankQuery(req)
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "Failed to build query", err)
	}

	// 2. 执行搜索
	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if err != nil || result == nil {
		status := http.StatusInternalServerError
		msg := "Failed to get pump rank"
		data := []response.PoolMarketInfo{}
		if result == nil {
			status = http.StatusOK
			msg = "No data found"
			data := []response.PoolMarketInfo{}
			return response.BuildResponse(data, status, msg, nil)
		}
		return response.BuildResponse(data, status, msg, err)
	}

	// 3. 处理聚合结果
	transactions, err := ProcessAggregationResult(result, req.Filters, req)
	if err != nil {
		return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusInternalServerError, "Failed to process results", err)
	}

	// 4. 构建响应
	responseBuilder := NewResponseBuilder(transactions, req.OrderBy, req.Direction)
	return responseBuilder.BuildRaydiumResponse()
}

func GetSolSwapRank(req *request.SolRankRequest) response.Response {

	// 1. 构建查询
	queryJSON, err := es.SolSwapQuery(req)
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "Failed to build query", err)
	}

	// 2. 执行搜索
	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if err != nil || result == nil {
		status := http.StatusInternalServerError
		msg := "Failed to get pump rank"
		data := []response.PoolMarketInfo{}
		if result == nil {
			status = http.StatusOK
			msg = "No data found"
			data := []response.PoolMarketInfo{}
			return response.BuildResponse(data, status, msg, nil)
		}
		return response.BuildResponse(data, status, msg, err)
	}
	// 3. 处理聚合结果
	transactions, err := ProcessAggregationResult(result, req.Filters, req)
	if err != nil {
		return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusInternalServerError, "Failed to process results", err)
	}

	// 4. 构建响应
	responseBuilder := NewResponseBuilder(transactions, req.OrderBy, req.Direction)
	return responseBuilder.BuildSwapResponse()
}

func GetNewPairRanks(req *request.SolRankRequest) response.Response {
	// 1. 构建查询
	queryJSON, err := es.NewPairRanksQuery(req)
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "Failed to build query", err)
	}

	// 2. 执行搜索
	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if err != nil || result == nil {
		status := http.StatusInternalServerError
		msg := "Failed to get pump rank"
		data := []response.PoolMarketInfo{}
		if result == nil {
			status = http.StatusOK
			msg = "No data found"
			data := []response.PoolMarketInfo{}
			return response.BuildResponse(data, status, msg, nil)
		}
		return response.BuildResponse(data, status, msg, err)
	}
	// 3. 处理聚合结果
	transactions, err := ProcessAggregationResult(result, req.Filters, req)
	if err != nil {
		return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusInternalServerError, "Failed to process results", err)
	}
	// 4. 构建响应
	responseBuilder := NewResponseBuilder(transactions, req.OrderBy, req.Direction)
	return responseBuilder.BuildRaydiumResponse()
}

// 通用的获取函数
func GetTokenTransactionsFromRedis(key string) ([]response.TokenTransaction, error) {
	data, err := redis.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from redis: %v", err)
	}

	var transactions []response.TokenTransaction
	if err := json.Unmarshal([]byte(data), &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return transactions, nil
}
