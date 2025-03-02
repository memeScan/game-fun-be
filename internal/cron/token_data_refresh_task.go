package cron

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"my-token-ai-be/internal/constants"
	"my-token-ai-be/internal/es"
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/pkg/httpRespone"
	"my-token-ai-be/internal/pkg/httpUtil"
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/redis"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"
	"sync"
	"time"
)

var tokenInfoService = service.TokenInfoService{}
var redisKey = constants.RedisKeyHotTokens

func CompletedTokenDataRefreshTaskQuery() {
	var safety request.SolRankRequest
	safety.Time = "1h"
	safety.Limit = 50
	queryJSON, err := es.CompletedQuery(&safety)
	if err != nil {
		log.Printf("Error searching documents: %v", err)
	}
	taskName := "RefreshCompletedTokenDataJob"
	lockKey := "lock:" + taskName
	ttl := time.Now().Add(2 * time.Hour).Unix()
	RefreshHotTokensJob(taskName, lockKey, queryJSON, redisKey, ttl)
}

func SwapToken1mDataRefreshTaskQuery() error {
	var safety request.SolRankRequest
	safety.Time = "1m"
	safety.Limit = 20
	queryJSON, err := es.SolSwapQuery(&safety)
	if err != nil {
		log.Printf("Error searching documents: %v", err)
		return err
	}
	taskName := "RefreshSwapTokens1mJob"
	lockKey := "lock:" + taskName
	ttl := time.Now().Add(2 * time.Hour).Unix()
	RefreshHotTokensJob(taskName, lockKey, queryJSON, redisKey, ttl)
	return nil
}

func SwapToken5mDataRefreshTaskQuery() error {
	var safety request.SolRankRequest
	safety.Time = "5m"
	safety.Limit = 50
	queryJSON, err := es.SolSwapQuery(&safety)
	if err != nil {
		log.Printf("Error searching documents: %v", err)
		return err
	}
	taskName := "RefreshSwapTokens5mJob"
	lockKey := "lock:" + taskName
	ttl := time.Now().Add(2 * time.Hour).Unix()
	RefreshHotTokensJob(taskName, lockKey, queryJSON, redisKey, ttl)
	return nil
}

func SwapToken1hDataRefreshTaskQuery() error {
	var safety request.SolRankRequest
	safety.Time = "1h"
	safety.Limit = 50
	queryJSON, err := es.SolSwapQuery(&safety)
	if err != nil {
		log.Printf("Error searching documents: %v", err)
	}
	taskName := "RefreshSwapTokens1hJob"
	lockKey := "lock:" + taskName
	ttl := time.Now().Add(2 * time.Hour).Unix()
	RefreshHotTokensJob(taskName, lockKey, queryJSON, redisKey, ttl)
	return nil
}

func SwapToken6hDataRefreshTaskQuery() error {
	var safety request.SolRankRequest
	safety.Time = "6h"
	safety.Limit = 50
	queryJSON, err := es.SolSwapQuery(&safety)
	if err != nil {
		log.Printf("Error searching documents: %v", err)
		return err
	}
	taskName := "RefreshSwapTokens6hJob"
	lockKey := "lock:" + taskName
	ttl := time.Now().Add(2 * time.Hour).Unix()
	RefreshHotTokensJob(taskName, lockKey, queryJSON, redisKey, ttl)
	return nil
}

func SwapToken1dDataRefreshTaskQuery() error {
	var safety request.SolRankRequest
	safety.Time = "1d"
	safety.Limit = 50
	queryJSON, err := es.SolSwapQuery(&safety)
	if err != nil {
		log.Printf("Error searching documents: %v", err)
		return err
	}
	taskName := "RefreshSwapTokens1dJob"
	lockKey := "lock:" + taskName
	ttl := time.Now().Add(2 * time.Hour).Unix()
	RefreshHotTokensJob(taskName, lockKey, queryJSON, redisKey, ttl)
	return nil
}

func RefreshHotTokensJob(taskName string, lockKey string, queryJSON string, redisKey string, tokenTTL int64) error {

	// 1. 先检查 key 是否存在
	exists, err := redis.Exists(redisKey)
	if err != nil {
		util.Log().Error("检查 Redis key 是否存在失败: %v", err)
		return err
	}

	// 2. 如果 key 不存在，先创建
	if !exists {
		// key 的 ttl 0 表示永久存储
		keyTTL := 0
		util.Log().Info("Redis key 不存在，创建新的 Sorted Set")
		err = redis.CreateSortedSet(redisKey, int64(keyTTL))
		if err != nil {
			util.Log().Error("创建 Sorted Set 失败: %v", err)
			return err
		}
	}

	// 3. 安全清理过期数据
	err = redis.SafeCleanExpiredTokens(context.Background(), redisKey)
	if err != nil {
		util.Log().Error("清理过期数据失败: %v", err)
		// 不要直接返回，继续执行
	}

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

	// 2. 解析查询结果
	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if err != nil {
		log.Printf("Error searching documents: %v", err)
		return err
	}

	// 3. 解析查询结果
	aggregationResult, err := es.UnmarshalAggregationResult(result)

	if err != nil {
		log.Printf("Error unmarshaling aggregation result: %v", err)
		util.Log().Error("Error unmarshaling aggregation result: %v", err)
		return err
	}

	if aggregationResult == nil {
		log.Println("aggregationResult is nil")
		return nil
	}

	if len(aggregationResult.Buckets) == 0 {
		util.Log().Info("No completed tokens found")
		return nil
	}

	var tokenAddresses []string
	var tokens []map[string]string

	for _, bucket := range aggregationResult.Buckets {
		var tokenTransaction response.TokenTransaction
		if len(bucket.LatestTransaction.Hits.Hits) > 0 &&
			bucket.LatestTransaction.Hits.Hits[0].Source != nil {
			if err := json.Unmarshal(bucket.LatestTransaction.Hits.Hits[0].Source, &tokenTransaction); err != nil {
				log.Printf("Error unmarshaling hit source: %v", err)
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

	var hotTokens []string
	isInit := true

	hotTokens, err = redis.GetValidTokens(context.Background(), redisKey)
	if err != nil {
		util.Log().Error("Error getting token info: %v", err)
		return err
	}

	if len(hotTokens) == 0 {
		hotTokens = tokenAddresses
		isInit = false
		util.Log().Info("No completed tokens found in Redis, initialized")
	}

	// 需要检测的代币集合
	var needCheckTokens []string

	if !isInit {
		needCheckTokens = hotTokens
	} else {
		needCheckTokens = Difference(tokenAddresses, hotTokens)
		if len(needCheckTokens) == 0 {
			util.Log().Info("No new hot tokens found")
			return nil
		}
	}

	batchSize := 20
	for i := 0; i < len(needCheckTokens); i += batchSize {
		end := i + batchSize
		if end > len(needCheckTokens) {
			end = len(needCheckTokens)
		}
		batch := needCheckTokens[i:end]

		util.Log().Info("Processing batch %d to %d of %d tokens", i, end, len(needCheckTokens))

		// 批量查询代币信息
		tokenInfoMap, err := tokenInfoService.GetTokenInfoMapByDB(batch, uint8(model.ChainTypeSolana))
		if err != nil {
			util.Log().Error("Error getting token info: %v", err)
			continue
		}

		var tokenInfos []*model.TokenInfo
		// tokens 和 needCheckTokens 做差集，得到需要检测的代币地址池地址
		var needCheckTokenAddressPoolAddresses []map[string]string
		for _, token := range tokens {
			if Contains(batch, token["mints"]) {
				needCheckTokenAddressPoolAddresses = append(needCheckTokenAddressPoolAddresses, token)
			}
		}

		var wg sync.WaitGroup

		// 用于存储结果
		var safetyCheckData *[]httpRespone.SafetyCheckPoolData
		var dexCheckData *[]httpRespone.DexCheckData
		var safetyCheckErr error
		var dexCheckErr error

		// 启动 goroutine 获取安全检查数据
		wg.Add(1)
		go func() {
			defer wg.Done()
			safetyCheckData, safetyCheckErr = httpUtil.GetSafetyCheckPool(needCheckTokenAddressPoolAddresses)
			if safetyCheckErr != nil {
				util.Log().Error("Error getting safety check data: %v", safetyCheckErr)
			}
		}()

		// 启动 goroutine 获取 DEX 检查数据
		wg.Add(1)
		go func() {
			defer wg.Done()
			dexCheckData, dexCheckErr = httpUtil.GetDexCheck(batch)
			if dexCheckErr != nil {
				util.Log().Error("Error getting dex check data: %v", dexCheckErr)
			}
		}()

		wg.Wait()

		// 在处理 safetyCheckData 之前添加空指针检查
		if safetyCheckData != nil {
			for _, safetyData := range *safetyCheckData {
				tokenInfo := &model.TokenInfo{}
				tokenInfo.TokenAddress = safetyData.Mint
				if _, ok := tokenInfoMap[safetyData.Mint]; !ok {
					util.Log().Error("Token info not found for mint: %s", safetyData.Mint)
					continue
				}
				tokenInfo.ChainType = tokenInfoMap[safetyData.Mint].ChainType
				tokenInfo.TotalSupply = tokenInfoMap[safetyData.Mint].TotalSupply
				decimals := tokenInfoMap[safetyData.Mint].Decimals
				actualTotalSupply := float64(tokenInfo.TotalSupply) / math.Pow(10, float64(decimals))

				// 设置Top10持仓百分比
				if safetyData.Top10Holdings > 0 {
					percentage := float64(safetyData.Top10Holdings) / actualTotalSupply
					tokenInfo.Top10Percentage = percentage
					safetyData.Top10Holdings = percentage
				} else {
					tokenInfo.Top10Percentage = 0
				}
				if safetyData.LpBurnedPercentage > 0 {
					tokenInfo.BurnPercentage = safetyData.LpBurnedPercentage
					if safetyData.LpBurnedPercentage > 0.5 {
						tokenInfo.SetFlag(model.FLAG_BURNED_LP)
					}
				} else {
					tokenInfo.BurnPercentage = 0
				}
				if safetyData.Holders > 0 {
					tokenInfo.Holder = safetyData.Holders
					tokenInfos = append(tokenInfos, tokenInfo)
				}
				key := constants.RedisKeySafetyCheck + safetyData.Mint
				redis.Set(key, safetyData, 60*time.Minute)
			}
		}

		// 1. 检查 dexCheckData
		if dexCheckData != nil {
			for _, dexCheck := range *dexCheckData {
				for i, tokenInfo := range tokenInfos {
					if tokenInfo.TokenAddress == dexCheck.Address {
						if dexCheck.Websites != nil || dexCheck.Socials != nil {
							tokenInfos[i].SetFlag(model.FLAG_DEXSCR_UPDATE)
						}
						if dexCheck.Boosts != nil && dexCheck.Boosts.Active > 0 {
							tokenInfos[i].SetFlag(model.FLAG_DXSCR_AD)
						}
					}
				}
			}
		}

		batchUpdateResp := tokenInfoService.BatchUpdateTokenInfoV2(tokenInfos)

		if batchUpdateResp.Code != 0 {
			util.Log().Error("job 批量更新代币信息失败: %v", batchUpdateResp.Error)
			continue
		}
		util.Log().Info("job 批量更新代币信息成功: %d", len(tokenInfos))

		updatedTokens, ok := batchUpdateResp.Data.([]*model.TokenInfo) // 添加类型断言检查
		if !ok {
			util.Log().Error("Invalid type assertion for batchUpdateResp.Data")
			continue
		}
		updatedTokensString := make(map[string]int64)
		for _, token := range updatedTokens {
			if token != nil { // 检查 token 是否为 nil
				updatedTokensString[token.TokenAddress] = int64(tokenTTL)
			}
		}
		redis.SafeBatchAddTokens(context.Background(), redisKey, updatedTokensString, tokenTTL)

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
