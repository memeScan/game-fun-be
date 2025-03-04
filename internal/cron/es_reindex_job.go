package cron

import (
	"context"
	"fmt"
	"game-fun-be/internal/es"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/service"
	"time"

	"github.com/olivere/elastic/v7"
)

// 执行重新索引任务的包装函数
func ExecuteReindexJob() error {
	util.Log().Info("执行重新索引任务")

	// 1. 获取当前别名对应的实际索引名称
	currentIndex, err := es.GetAliasActualIndex(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS)
	if err != nil {
		return fmt.Errorf("获取当前索引名称失败: %v", err)
	}
	util.Log().Info("当前索引名称: %s", currentIndex)

	// 从当前索引名提取版本号并递增
	var currentVersion int
	_, err = fmt.Sscanf(currentIndex, "token_transactions_new_v%d", &currentVersion)
	if err != nil {
		return fmt.Errorf("解析当前索引版本号失败: %v", err)
	}

	// 创建新索引名称（版本号+1）
	newIndex := fmt.Sprintf("token_transactions_new_v%d", currentVersion+1)

	// 2. 创建新索引
	err = es.CreateTokenTransactionIndex(newIndex)
	if err != nil {
		return fmt.Errorf("创建新索引失败: %v", err)
	}
	util.Log().Info("创建新索引成功: %s", newIndex)

	// 3. 更新别名（原子操作：删除旧索引别名，添加新索引别名）
	err = es.UpdateAlias(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, currentIndex, newIndex)
	if err != nil {
		return fmt.Errorf("更新别名失败: %v", err)
	}
	util.Log().Info("更新别名成功，从 %s 切换到 %s", currentIndex, newIndex)

	// 4. 执行重建索引
	err = es.ReindexWithProgress(currentIndex, newIndex)
	if err != nil {
		util.Log().Error("重建索引失败，开始清理和恢复操作")

		// 1. 首先检查别名当前指向
		actualIndex, checkErr := es.GetAliasActualIndex(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS)
		if checkErr != nil {
			util.Log().Error("检查别名失败: %v", checkErr)
		}

		// 2. 根据别名指向决定恢复操作
		if actualIndex == newIndex {
			// 只有当别名确实指向新索引时，才需要切回旧索引
			if aliasErr := es.UpdateAlias(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, newIndex, currentIndex); aliasErr != nil {
				util.Log().Error("恢复别名到旧索引失败: %v", aliasErr)
			} else {
				util.Log().Info("成功将别名恢复到旧索引: %s", currentIndex)
			}
		} else if actualIndex == "" {
			// 如果别名不存在，直接创建指向旧索引的别名
			if aliasErr := es.CreateAlias(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, currentIndex); aliasErr != nil {
				util.Log().Error("创建指向旧索引的别名失败: %v", aliasErr)
			} else {
				util.Log().Info("成功创建指向旧索引的别名: %s", currentIndex)
			}
		} else {
			util.Log().Info("别名当前指向索引: %s，无需恢复", actualIndex)
		}

		// 3. 清理失败的新索引
		if deleteErr := es.DeleteIndex(newIndex); deleteErr != nil {
			util.Log().Error("清理失败的新索引失败: %v", deleteErr)
		} else {
			util.Log().Info("成功清理失败的新索引: %s", newIndex)
		}

		return fmt.Errorf("重建索引失败: %v", err)
	}

	// 重建成功后，删除旧索引
	if err := es.DeleteIndex(currentIndex); err != nil {
		util.Log().Error("删除旧索引失败: %v", err)
		// 不返回错误，因为重建已经成功，删除旧索引失败不应影响整体流程
	}
	util.Log().Info("重建索引完成，旧索引已删除")

	return nil
}

func SyncTokenInfoJob() error {
	// 偏移量根据es的文档数来
	total, err := es.ESClient.Count(es.ES_INDEX_TOKEN_INFO).Do(context.Background())
	if err != nil {
		util.Log().Error("获取代币信息总数失败: %v", err)
		return fmt.Errorf("failed to get total token infos: %w", err)
	}
	limit := 10000
	offset := int(total)
	totalProcessed := 0
	startTime := time.Now()
	batchSize := 500
	maxRetries := 3

	util.Log().Info("开始同步代币信息到 ES, 批次大小: %d", limit)

	tokenInfoService := service.TokenInfoService{}
	for {
		batchStartTime := time.Now()

		tokenInfos, err := tokenInfoService.ListTokenInfos(limit, offset)
		if err != nil {
			util.Log().Error("获取代币信息失败: %v", err)
			return fmt.Errorf("failed to get all token infos: %w", err)
		}

		if len(tokenInfos) == 0 {
			util.Log().Info("同步完成 - 总处理: %d, 总耗时: %v",
				totalProcessed, time.Since(startTime))
			break
		}

		util.Log().Info("处理批次 - 偏移量: %d, 数量: %d", offset, len(tokenInfos))

		for i := 0; i < len(tokenInfos); i += batchSize {
			end := i + batchSize
			if end > len(tokenInfos) {
				end = len(tokenInfos)
			}

			currentBatch := make([]*model.TokenInfo, len(tokenInfos[i:end]))
			for j, token := range tokenInfos[i:end] {
				currentBatch[j] = &token
			}

			failedTokens := make([]*model.TokenInfo, 0)
			successCount := processTokenBatch(currentBatch, &failedTokens)
			totalProcessed += successCount

			if len(failedTokens) > 0 {
				retryAndLogFailures(failedTokens, maxRetries, &totalProcessed)
			}

			time.Sleep(2 * time.Second) // 控制写入速率
		}

		offset += len(tokenInfos)
		util.Log().Info("批次完成 - 总处理: %d, 耗时: %v",
			totalProcessed, time.Since(batchStartTime))
	}

	return nil
}

func processTokenBatch(tokens []*model.TokenInfo, failedTokens *[]*model.TokenInfo) int {
	if len(tokens) == 0 {
		return 0
	}

	// 检查已存在文档
	existingDocs := make(map[string]bool)
	var searchQueries []elastic.Query

	for _, tokenInfo := range tokens {
		docID := fmt.Sprintf("%s_%d", tokenInfo.TokenAddress, tokenInfo.ChainType)
		searchQueries = append(searchQueries, elastic.NewTermQuery("_id", docID))
	}

	existsQuery := elastic.NewBoolQuery().Should(searchQueries...).MinimumShouldMatch("1")
	searchResult, err := es.ESClient.Search().
		Index("token_info_v1").
		Query(existsQuery).
		Size(len(tokens)).
		FetchSource(false).
		Do(context.Background())

	if err != nil {
		util.Log().Error("检查文档存在性失败: %v", err)
	} else {
		for _, hit := range searchResult.Hits.Hits {
			existingDocs[hit.Id] = true
		}
	}

	bulkRequest := es.ESClient.Bulk().Index("token_info_v1")
	tokenInfoService := service.TokenInfoService{}
	failedInBulk := make(map[string]*model.TokenInfo)
	skippedCount := 0

	for _, tokenInfo := range tokens {
		docID := fmt.Sprintf("%s_%d", tokenInfo.TokenAddress, tokenInfo.ChainType)

		if existingDocs[docID] {
			skippedCount++
			continue
		}

		doc, err := tokenInfoService.TokenInfoToESDoc(tokenInfo)
		if err != nil {
			util.Log().Error("转换文档失败: %s_%d: %v",
				tokenInfo.TokenAddress, tokenInfo.ChainType, err)
			*failedTokens = append(*failedTokens, tokenInfo)
			continue
		}

		failedInBulk[docID] = tokenInfo
		req := elastic.NewBulkIndexRequest().
			Index("token_info_v1").
			Id(docID).
			Doc(doc)

		bulkRequest = bulkRequest.Add(req)
	}

	if bulkRequest.NumberOfActions() == 0 {
		if skippedCount > 0 {
			util.Log().Info("跳过 %d 个已存在文档", skippedCount)
		}
		return skippedCount
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := bulkRequest.Do(ctx)
	if err != nil {
		util.Log().Error("批量索引失败: %v", err)
		*failedTokens = append(*failedTokens, tokens...)
		return skippedCount
	}

	if resp.Errors {
		for _, item := range resp.Failed() {
			if token, exists := failedInBulk[item.Id]; exists {
				*failedTokens = append(*failedTokens, token)
			}
		}
		util.Log().Error("部分文档索引失败 - 总数: %d, 失败: %d",
			bulkRequest.NumberOfActions(), len(resp.Failed()))
	}

	return len(tokens) - len(*failedTokens)
}

func retryAndLogFailures(failedTokens []*model.TokenInfo, maxRetries int, totalProcessed *int) {
	remainingTokens := failedTokens

	for retry := 1; retry < maxRetries && len(remainingTokens) > 0; retry++ {
		sleepTime := time.Duration(retry*2) * time.Second
		util.Log().Info("重试 %d 个失败文档 (第 %d/%d 次)",
			len(remainingTokens), retry+1, maxRetries)
		time.Sleep(sleepTime)

		var retryFailedTokens []*model.TokenInfo
		retrySuccess := processTokenBatch(remainingTokens, &retryFailedTokens)
		*totalProcessed += retrySuccess
		remainingTokens = retryFailedTokens

		if len(remainingTokens) == 0 {
			util.Log().Info("重试成功")
			break
		}
	}

	if len(remainingTokens) > 0 {
		util.Log().Error("最终失败 %d 个文档", len(remainingTokens))
		for _, token := range remainingTokens {
			util.Log().Error("失败文档: %s_%d",
				token.TokenAddress, token.ChainType)
		}
	}
}
