package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"game-fun-be/internal/constants"
	"game-fun-be/internal/es"
	"game-fun-be/internal/es/query"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"
	"game-fun-be/internal/service"
	"log"
	"time"
)

var tokenInfoService = service.TokenInfoService{}
var redisKey = constants.RedisKeyHotTokens

func RefreshHotTokensHolderJob() {
	var req request.TickersRequest
	req.Limit = 30
	req.SortDirection = "DESC"
	req.SortedBy = "NATIVE_VOLUME_24H"

	// 1. 生成查询JSON
	queryJSON, err := query.MarketQuery(&req)
	if err != nil {
		log.Printf("Error generating market query: %v", err)
		return
	}

	// 2. 获取代币列表
	tokenAddresses, _, err := getTokenList(queryJSON)
	if err != nil {
		log.Printf("Error getting token list: %v", err)
		return
	}

	// 3. 设置任务参数
	taskName := "RefreshHoldersJob"
	lockKey := "lock:" + taskName
	ttl := time.Now().Add(2 * time.Hour).Unix()

	// 4. 调用RefreshHotTokensJob
	if err := RefreshHotTokensJob(taskName, lockKey, tokenAddresses, redisKey, ttl); err != nil {
		log.Printf("Error refreshing hot tokens: %v", err)
	}
}

// 获取代币列表
func getTokenList(queryJSON string) ([]string, []map[string]string, error) {
	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if err != nil {
		util.Log().Error("Error searching documents: %v", err)
		return nil, nil, err
	}

	aggregationResult, err := es.UnmarshalAggregationResult(result)
	if err != nil {
		util.Log().Error("Error unmarshaling aggregation result: %v", err)
		return nil, nil, err
	}

	if aggregationResult == nil || len(aggregationResult.Buckets) == 0 {
		util.Log().Info("No completed tokens found")
		return nil, nil, nil
	}

	var tokenAddresses []string
	var tokens []map[string]string

	for _, bucket := range aggregationResult.Buckets {
		var tokenTransaction response.TokenTransaction
		if len(bucket.LatestTransaction.Hits.Hits) > 0 &&
			bucket.LatestTransaction.Hits.Hits[0].Source != nil {
			if err := json.Unmarshal(bucket.LatestTransaction.Hits.Hits[0].Source, &tokenTransaction); err != nil {
				util.Log().Error("Error unmarshaling hit source: %v", err)
				continue
			}
			if tokenTransaction.ExtInfo == "" {
				continue
			}
		}

		tokenAddress := tokenTransaction.TokenAddress
		var tokenPoolAddress string
		if tokenTransaction.CreatedPlatformType != 1 {
			tokenPoolAddress = tokenTransaction.PoolAddress
		}

		tokens = append(tokens, map[string]string{
			"mints":         tokenAddress,
			"poolAddresses": tokenPoolAddress,
		})
		tokenAddresses = append(tokenAddresses, tokenAddress)
	}

	return tokenAddresses, tokens, nil
}

// 操作Redis
func processRedisOperations(redisKey string, tokenAddresses []string, tokenTTL int64) ([]string, error) {
	// 1. 检查 key 是否存在
	exists, err := redis.Exists(redisKey)
	if err != nil {
		util.Log().Error("检查 Redis key 是否存在失败: %v", err)
		return nil, err
	}

	// 2. 如果 key 不存在，先创建
	if !exists {
		util.Log().Info("Redis key 不存在，创建新的 Sorted Set")
		if err := redis.CreateSortedSet(redisKey, 0); err != nil {
			util.Log().Error("创建 Sorted Set 失败: %v", err)
			return nil, err
		}
	}

	// 3. 安全清理过期数据
	if err := redis.SafeCleanExpiredTokens(context.Background(), redisKey); err != nil {
		util.Log().Error("清理过期数据失败: %v", err)
	}

	// 4. 获取有效代币
	hotTokens, err := redis.GetValidTokens(context.Background(), redisKey)
	if err != nil {
		util.Log().Error("Error getting token info: %v", err)
		return nil, err
	}

	// 5. 如果需要初始化，直接使用新获取的代币列表
	if len(hotTokens) == 0 {
		hotTokens = tokenAddresses
		util.Log().Info("No completed tokens found in Redis, initialized")
	}

	return hotTokens, nil
}

// 批量更新代币信息
func batchUpdateTokenInfo(batch []string, tokens []map[string]string, tokenTTL int64) error {
	// 批量查询代币信息
	tokenInfoMap, err := tokenInfoService.GetTokenInfoMapByDB(batch, uint8(model.ChainTypeSolana))
	if err != nil {
		util.Log().Error("Error getting token info: %v", err)
		return err
	}

	var tokenInfos []*model.TokenInfo
	for _, tokenAddress := range batch {
		solanaTrackerToken, err := httpUtil.GetTokenInfoByAddress(tokenAddress)
		if err != nil {
			log.Printf("Failed to get token info for %s: %v", tokenAddress, err)
			continue
		}

		// 仅当 tokenInfoMap 中已存在该 tokenAddress 时，才更新数据
		if tokenInfo, exists := tokenInfoMap[tokenAddress]; exists {
			tokenInfo.Holder = solanaTrackerToken.Holders
			tokenInfos = append(tokenInfos, tokenInfo)
			log.Printf("Updated Token Info for %s: %+v", tokenAddress, tokenInfo)
		} else {
			log.Printf("Token %s not found in tokenInfoMap, skipping update", tokenAddress)
		}
	}

	batchUpdateResp := tokenInfoService.BatchUpdateTokenInfo(tokenInfos)
	if batchUpdateResp.Code != 0 {
		util.Log().Error("job 批量更新代币信息失败: %v", batchUpdateResp.Error)
		return fmt.Errorf("job 批量更新代币信息失败: %s", batchUpdateResp.Error)
	}

	util.Log().Info("job 批量更新代币信息成功: %d", len(tokenInfos))

	updatedTokens, ok := batchUpdateResp.Data.([]*model.TokenInfo)
	if !ok {
		util.Log().Error("Invalid type assertion for batchUpdateResp.Data")
		return fmt.Errorf("invalid type assertion for batchUpdateResp.Data")
	}

	updatedTokensString := make(map[string]int64)
	for _, token := range updatedTokens {
		if token != nil {
			updatedTokensString[token.TokenAddress] = int64(tokenTTL)
		}
	}

	if err := redis.SafeBatchAddTokens(context.Background(), redisKey, updatedTokensString, tokenTTL); err != nil {
		util.Log().Error("Error adding tokens to Redis: %v", err)
		return err
	}

	return nil
}

func RefreshHotTokensJob(taskName string, lockKey string, tokenAddresses []string, redisKey string, tokenTTL int64) error {
	// 获取锁，设置过期时间为 30 分钟
	locked, err := redis.SetNX(lockKey, "1", 30*time.Minute)
	if err != nil {
		util.Log().Error("%s 获取分布式锁失败: %v", taskName, err)
		return err
	}
	if !locked {
		util.Log().Info("%s 已有任务在执行中，跳过本次执行", taskName)
		return nil
	}

	start := time.Now()
	util.Log().Info("%s 开始执行", taskName)
	defer func() {
		redis.Delete(lockKey)
		util.Log().Info("%s 执行完成, 耗时 %v", taskName, time.Since(start))
	}()

	// 操作Redis
	hotTokens, err := processRedisOperations(redisKey, tokenAddresses, tokenTTL)
	if err != nil {
		return err
	}

	// 需要检测的代币集合
	var needCheckTokens []string
	if len(hotTokens) == 0 {
		needCheckTokens = tokenAddresses
	} else {
		needCheckTokens = Difference(tokenAddresses, hotTokens)
		if len(needCheckTokens) == 0 {
			util.Log().Info("No new hot tokens found")
			return nil
		}
	}

	// 批量处理代币信息
	batchSize := 20
	for i := 0; i < len(needCheckTokens); i += batchSize {
		end := i + batchSize
		if end > len(needCheckTokens) {
			end = len(needCheckTokens)
		}
		batch := needCheckTokens[i:end]

		util.Log().Info("Processing batch %d to %d of %d tokens", i, end, len(needCheckTokens))

		if err := batchUpdateTokenInfo(batch, nil, tokenTTL); err != nil {
			util.Log().Error("Error processing batch: %v", err)
			continue
		}
	}

	return nil
}

// Difference 返回 slice1 中有而 slice2 中没有的元素
func Difference(slice1, slice2 []string) []string {
	// 创建一个 map 用于存储 slice2 中的元素
	set := make(map[string]struct{}, len(slice2))
	for _, item := range slice2 {
		set[item] = struct{}{}
	}

	// 创建一个 slice 用于存储差集结果
	var diff []string

	// 遍历 slice1,将不在 slice2 中的元素添加到差集结果中
	for _, item := range slice1 {
		if _, found := set[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff
}

// Contains 检查切片中是否包含特定的元素
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
