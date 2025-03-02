package service

import (
	"encoding/json"
	"fmt"
	"math"
	"my-token-ai-be/internal/constants"
	"my-token-ai-be/internal/es"
	"my-token-ai-be/internal/pkg/httpRespone"
	"my-token-ai-be/internal/redis"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
)

type TransactionProcessor struct {
	aggregationResult *es.AggregationResult
	filters           []string
	rankRequest       *request.SolRankRequest
}

// 处理聚合结果的通用方法
func ProcessAggregationResult(result json.RawMessage, filters []string, rankRequest *request.SolRankRequest) ([]*response.TokenTransaction, error) {
	if result == nil {
		return nil, fmt.Errorf("no data found")
	}

	aggregationResult, err := es.UnmarshalAggregationResult([]byte(result))
	if err != nil {
		return nil, err
	}

	if len(aggregationResult.Buckets) == 0 {
		return nil, nil
	}

	processor := &TransactionProcessor{
		aggregationResult: aggregationResult,
		filters:           filters,
		rankRequest:       rankRequest,
	}

	return processor.processTransactions()
}

func (p *TransactionProcessor) processTransactions() ([]*response.TokenTransaction, error) {
	var tokenTransactions []*response.TokenTransaction

	for _, bucket := range p.aggregationResult.Buckets {
		tx, err := p.processSingleTransaction(bucket)
		if err != nil {
			continue
		}
		if tx != nil {
			tokenTransactions = append(tokenTransactions, tx)
		}
	}

	return tokenTransactions, nil
}

func (p *TransactionProcessor) processSingleTransaction(bucket es.Bucket) (*response.TokenTransaction, error) {
	if len(bucket.LatestTransaction.Hits.Hits) == 0 {
		return nil, nil
	}

	var tx response.TokenTransaction
	if err := json.Unmarshal(bucket.LatestTransaction.Hits.Hits[0].Source, &tx); err != nil {
		return nil, err
	}

	// 处理基础数据
	p.processBasicInfo(&tx, bucket)

	// 处理市值变化
	processMarketCapChanges(&tx, bucket)

	// 处理标志位
	processTokenFlags(&tx)

	if !shouldIncludeTokenTransaction(tx, p.filters) {
		return nil, nil
	}

	if !checkTransactionConditions(tx, p.rankRequest) {
		return nil, nil
	}

	return &tx, nil
}

func (p *TransactionProcessor) processBasicInfo(tx *response.TokenTransaction, bucket es.Bucket) {

	// 处理交易量
	if bucket.Volume.Value != 0 {
		tx.Volume = bucket.Volume.Value / math.Pow(10, float64(response.PumpDecimals)) * tx.Price
	}

	if bucket.SellCount1h.SellVolume.Value != 0 {
		tx.SellCount1h = bucket.SellCount1h.SellVolume.Value
	}
	if bucket.BuyCount1h.BuyVolume.Value != 0 {
		tx.BuyCount1h = bucket.BuyCount1h.BuyVolume.Value
	}

	// 处理买卖统计
	tx.Buys = bucket.Buys.BuyVolume.Value
	tx.Sells = bucket.Sells.SellVolume.Value
	tx.Swaps = tx.Buys + tx.Sells
	tx.SwapsCount1h = tx.Swaps
	tx.Volume1h = tx.Volume

	// 处理持有者数量
	if bucket.HolderCount.UniqueUsers.Value != 0 {
		tx.Holder = bucket.HolderCount.UniqueUsers.Value
	}

	if bucket.TotalHoldersPercentage.Value != 0 {
		tx.Top10HolderRate = bucket.TotalHoldersPercentage.Value
	}

	tx = updateHolderFromRedis(tx)

	// 处理流动性 除以 sol的 decimals 次方
	processTokenTransaction(tx, response.SolDecimals)

}

func updateHolderFromRedis(tx *response.TokenTransaction) *response.TokenTransaction {
	if tx.Holder != 0 {
		return tx
	}

	// 查 redis
	key := constants.RedisKeySafetyCheck + tx.TokenAddress
	redisData, err := redis.Get(key)
	if err != nil || redisData == "" {
		return tx
	}

	// redisData 转换成 model.SafetyCheckData
	safetyCheckData := &httpRespone.SafetyCheckPoolData{}
	if err := json.Unmarshal([]byte(redisData), safetyCheckData); err != nil {
		return tx
	}

	tx.Holder = safetyCheckData.Holders
	tx.Top10HolderRate = safetyCheckData.Top10Holdings
	return tx
}
