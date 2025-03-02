package cron

import (
	"encoding/json"
	"my-token-ai-be/internal/es"
	"my-token-ai-be/internal/pkg/util"
)

func DeleteDocumentsJob() error {
	// 每次删除成功是10000条，
	deletedCount := 0
	i := 0
	// 循环查询删除
	for {
		// Step 1: 构建查询请求
		query := map[string]interface{}{
			"size":    8000,
			"_source": false, // 我们只需要文档 ID，不需要其他字段
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"range": map[string]interface{}{
								"token_create_time": map[string]interface{}{
									"lte": "now-1d",
								},
							},
						},
					},
				},
			},
			"sort": []map[string]interface{}{
				{
					"transaction_time": map[string]interface{}{
						"order": "desc",
					},
				},
			},
		}

		queryStr, err := json.Marshal(query)
		if err != nil {
			return err
		}

		util.Log().Info("查询条件： %s", string(queryStr))

		searchResult, err := es.SearchDocumentsV2(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, string(queryStr))
		if err != nil {
			return err
		}

		var docIDs []string

		for _, result := range searchResult {
			docIDs = append(docIDs, string(result))
		}

		// 如果没有更多文档，退出循环
		if len(docIDs) == 0 {
			util.Log().Info("No more documents to delete")
			break
		}
		_, err = es.DeleteDocuments(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, docIDs)
		if err != nil {
			return err
		}

		deletedCount += len(searchResult)

		i++
		util.Log().Info("删除第 %d 次", i)

	}

	util.Log().Info("删除成功！ 总共删除：%d 条", deletedCount)

	return nil
}
