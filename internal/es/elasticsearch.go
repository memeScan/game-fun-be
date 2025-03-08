package es

import (
	"context"
	"encoding/json"
	"fmt"
	"game-fun-be/internal/pkg/util"
	"io"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

var ESClient *elastic.Client

func Elasticsearch() {
	var err error
	ESClient, err = elastic.NewClient(
		elastic.SetURL(os.Getenv("ES_URL")),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetBasicAuth(
			os.Getenv("ES_USERNAME"),
			os.Getenv("ES_PASSWORD"),
		),
	)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	// 测试连接
	info, code, err := ESClient.Ping(os.Getenv("ES_URL")).Do(context.Background())
	if err != nil {
		log.Fatalf("Error pinging Elasticsearch: %s", err)
	}

	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	util.Log().Info("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	log.Println("Elasticsearch connected successfully")

	//预创建索引结构
	ctx := context.Background()
	exists, err := ESClient.IndexExists(ES_INDEX_TOKEN_TRANSACTIONS).Do(ctx)
	if err != nil {
		log.Fatalf("Error checking if index exists: %s", err)
	}

	if exists {
		util.Log().Info("Elasticsearch index '%s' already exists", ES_INDEX_TOKEN_TRANSACTIONS)
	} else {
		err = CreateTokenTransactionIndex(ES_INDEX_TOKEN_TRANSACTIONS)
		if err != nil {
			log.Fatalf("Error creating token transactions index: %s", err)
		}

		// 创建别名指向索引
		currentIndex := ES_INDEX_TOKEN_TRANSACTIONS
		if aliasErr := CreateAlias(ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, currentIndex); aliasErr != nil {
			util.Log().Error("创建指向索引的别名失败: %v", aliasErr)
		} else {
			util.Log().Info("成功创建指向索引的别名: %s", ES_INDEX_TOKEN_TRANSACTIONS_ALIAS)
		}
		util.Log().Info("Elasticsearch index '%s' created successfully", ES_INDEX_TOKEN_TRANSACTIONS)
	}

	// ...
}

// 这里可以添加其他 Elasticsearch 相关的函数，例如：

// IndexDocument 索引一个文档
func IndexDocument(index string, id string, document interface{}) error {
	_, err := ESClient.Index().
		Index(index).
		Id(id).
		BodyJson(document).
		Refresh("true").
		Do(context.Background())
	return err
}

// SearchDocuments 搜索文档
func SearchDocuments(index string, query string) ([]json.RawMessage, error) {
	searchResult, err := ESClient.Search().
		Index(index).
		Source(query). // 直接使用传入的查询字符串
		Do(context.Background())
	if err != nil {
		log.Println("Error searching documents: ", err)
		return nil, err
	}

	var results []json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}

// SearchDocuments 搜索文档
func SearchDocumentsV2(index string, query string) ([]json.RawMessage, error) {
	searchResult, err := ESClient.Search().
		Index(index).
		Source(query). // 直接使用传入的查询字符串
		Do(context.Background())
	if err != nil {
		log.Println("Error searching documents: ", err)
		return nil, err
	}

	var results []json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		results = append(results, json.RawMessage(hit.Id))
	}

	return results, nil
}

// DeleteDocument 删除一个文档
func DeleteDocument(index string, id string) error {
	_, err := ESClient.Delete().
		Index(index).
		Id(id).
		Do(context.Background())
	return err
}

// 可以根据需要添加更多函数
func CreateTokenTransactionIndex(indexName string) error {
	indexDefinition := `{
		"settings": {
			"number_of_shards": 4,
			"number_of_replicas": 1,
			"refresh_interval": "1s",
			"analysis": {
				"analyzer": {
					"default": {
						"type": "keyword"
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"transaction_hash": {
					"type": "text",
					"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
						}
					}
				},			
				"id": { 
            		 "type": "keyword", 
             		"ignore_above": 20 
     		 	},	
				"token_address": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"from_address": { "type": "keyword" },
				"to_address": { "type": "keyword" },
				"user_address": {
					"type": "text",
					"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
					}
				},
				"token_amount": { "type": "keyword", "ignore_above": 20 },
				"native_token_amount": { "type": "keyword", "ignore_above": 20 },
				"price": { "type": "double" },
				"native_price": { "type": "double" },
				"transaction_time": { "type": "date" },
				"create_time": { "type": "date" },
				"update_time": { "type": "date" },
				"chain_type": { "type": "integer" },
				"platform_type": { "type": "integer" },
				"created_platform_type": { "type": "integer" },
				"is_buy": { "type": "boolean" },
				"progress": { "type": "double" },
				"is_complete": { "type": "boolean" },	
				"decimals": { "type": "integer" },
				"virtual_native_reserves": { 
					"type": "keyword",
					"ignore_above": 20
				},
				"virtual_token_reserves": { 
					"type": "keyword",
					"ignore_above": 20
				},
				"real_native_reserves": { 
					"type": "keyword",
					"ignore_above": 20
				},
				"real_token_reserves": { 
					"type": "keyword",
					"ignore_above": 20
				},
				"token_create_time": { "type": "date" },
				"dev_native_token_amount": { "type": "double" },
				"holder": { "type": "integer" },
				"comment_count": { "type": "integer" },
				"market_cap": { "type": "double" },
				"token_supply": { 
					"type": "keyword",
					"ignore_above": 20
				},
				"dev_status": { "type": "integer" },
				"crown_duration": { "type": "integer" },
				"rocket_duration": { "type": "integer" },
				"is_media": { "type": "boolean" },
				"uri": { "type": "text" },
				"ext_info": {
					"type": "text",
					"index": false,
					"doc_values": false,
					"store": true
				},
				"pool_address": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"native_token_address": { "type": "text" },
				"block_time": { "type": "date" },
				"creator": {
					"type": "text",
					"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
					}
				},
				"initial_pc_reserve": { "type": "keyword", "ignore_above": 20 },
				"pc_symbol": { "type": "text" },
				"is_honeypot": { "type": "boolean" },
				"burn_status": { "type": "integer" },
				"top_10_holder_rate": { "type": "double" },
				"creator_token_status": { "type": "integer" },
				"liquidity": { "type": "double" },
				"token_flags": { "type": "integer" },
				"burn_percentage": { "type": "double" },
				"dev_burn_percentage": { "type": "double" },
				"dev_percentage": { "type": "double" },
				"top_10_percentage": { "type": "double" },
				"dev_token_amount": { 
					"type": "keyword",
					"ignore_above": 20
				},
				"token_creator": { "type": "keyword" },
				"token_name": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"symbol": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"is_burned_lp": {
					"type": "boolean"
				},
				"is_dex_ad": {
					"type": "boolean"
				}
			}
		}
	}`

	ctx := context.Background()
	exists, err := ESClient.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	createIndex, err := ESClient.CreateIndex(indexName).BodyString(indexDefinition).Do(ctx)
	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		return fmt.Errorf("index creation was not acknowledged")
	}

	return nil
}

func SearchTokenTransactionsWithAggs(indexName string, query string, unique string) (json.RawMessage, error) {

	_, err := ESClient.Refresh().Index(indexName).Do(context.Background())
	if err != nil {
		log.Println("Error refreshing index: ", err)
		return nil, err
	}
	searchResult, err := ESClient.Search().
		Index(indexName).
		Source(query).
		Do(context.Background())

	if err != nil {
		log.Println("Error searching documents: ", err)
		return nil, err
	}

	// 获取聚合结果
	if agg, found := searchResult.Aggregations[unique]; found {
		return agg, nil
	}
	return nil, nil
}

// BulkIndexDocuments 批量索引文档
func BulkIndexDocuments(index string, documents []map[string]interface{}) (*elastic.BulkResponse, error) {
	totalStart := time.Now()
	ctx := context.Background()

	processor, err := ESClient.BulkProcessor().
		Name("TokenTransactionProcessor").
		Workers(16).
		BulkActions(10000).
		BulkSize(50 * 1024 * 1024).
		FlushInterval(5 * time.Second).
		Stats(true).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create bulk processor: %w", err)
	}
	defer processor.Close()

	// 2. 批量添加文档
	addStart := time.Now()
	docSize := 0
	for _, doc := range documents {
		id, ok := doc["id"].(string)
		if !ok {
			return nil, fmt.Errorf("document missing id")
		}

		// 计算文档大小
		if docBytes, err := json.Marshal(doc); err == nil {
			docSize += len(docBytes)
		}

		req := elastic.NewBulkIndexRequest().
			Index(index).
			Id(id).
			Doc(doc)

		processor.Add(req)
	}
	addTime := time.Since(addStart)
	util.Log().Info("添加文档耗时: %v (文档数: %d, 总大小: %.2fMB, 平均: %.2fKB/doc)",
		addTime,
		len(documents),
		float64(docSize)/1024/1024,
		float64(docSize)/float64(len(documents))/1024)

	// 在获取统计之前，显式刷新
	if err := processor.Flush(); err != nil {
		util.Log().Error("刷新bulk processor失败: %v", err)
	}

	// 4. 获取处理统计
	stats := processor.Stats()
	if stats.Failed > 0 {
		util.Log().Error("ES批量写入部分失败 - 总数: %d, 成功: %d, 失败: %d",
			stats.Indexed,
			stats.Succeeded,
			stats.Failed)
	}

	// 5. 构造返回结果
	response := &elastic.BulkResponse{
		Took:   int(stats.Indexed),
		Errors: stats.Failed > 0,
		Items:  make([]map[string]*elastic.BulkResponseItem, stats.Succeeded),
	}

	// 6. 总耗时统计
	totalTime := time.Since(totalStart)
	util.Log().Info("ES批量写入总耗时: %v (处理器: %v, 添加: %v, 刷新: %v), QPS: %.2f/s",
		totalTime,
		time.Since(addStart),
		addTime,
		time.Since(addStart),
		float64(stats.Succeeded)/totalTime.Seconds())

	return response, nil
}

func SetupILMPolicy(client *elastic.Client) error {
	policy := map[string]interface{}{
		"policy": map[string]interface{}{
			"phases": map[string]interface{}{
				"hot": map[string]interface{}{
					"actions": map[string]interface{}{
						"rollover": map[string]interface{}{
							"max_age":  "7d",
							"max_size": "100gb",
						},
					},
				},
				"delete": map[string]interface{}{
					"min_age": "30d",
					"actions": map[string]interface{}{
						"delete": map[string]interface{}{},
					},
				},
			},
		},
	}

	_, err := client.XPackIlmPutLifecycle().Policy("token_tx_policy").BodyJson(policy).Do(context.Background())
	return err
}

// 档
func DeleteDocuments(indexName string, docIDs []string) (string, error) {
	bulkRequest := ESClient.Bulk()
	for _, id := range docIDs {
		req := elastic.NewBulkDeleteRequest().Index(indexName).Id(id)
		bulkRequest = bulkRequest.Add(req)
	}

	ctx := context.Background()
	response, err := bulkRequest.Do(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", response.Took), nil
}

// GetAliasActualIndex returns the actual index name for a given alias
func GetAliasActualIndex(aliasName string) (string, error) {
	resp, err := ESClient.Aliases().Index(aliasName).Do(context.Background())
	if err != nil {
		return "", err
	}

	// Get the first index name from the alias (assuming single index per alias)
	for indexName := range resp.Indices {
		return indexName, nil
	}

	return "", fmt.Errorf("no index found for alias: %s", aliasName)
}

func UpdateAlias(aliasName, oldIndex, newIndex string) error {
	_, err := ESClient.Alias().
		Remove(oldIndex, aliasName).
		Add(newIndex, aliasName).
		Do(context.Background())
	return err
}

func ReindexWithProgress(currentIndex, newIndex string) error {
	// 获取符合条件的文档总数（最近25小时的数据）
	rangeQuery := elastic.NewRangeQuery("transaction_time").
		Gte("now-23h").
		Lte("now")

	count, err := ESClient.Count(currentIndex).
		Query(rangeQuery).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("获取源索引文档数量失败: %v", err)
	}
	totalDocs := count

	body := map[string]interface{}{
		"source": map[string]interface{}{
			"index": currentIndex,
			"query": map[string]interface{}{
				"range": map[string]interface{}{
					"transaction_time": map[string]interface{}{
						"gte": "now-23h",
						"lte": "now",
					},
				},
			},
			"size": 300,
		},
		"dest": map[string]interface{}{
			"index":   newIndex,
			"op_type": "create",
		},
		"script": map[string]interface{}{
			"source": "ctx._id = ctx._source.id",
			"lang":   "painless",
		},
		"conflicts": "proceed",
	}

	res, err := ESClient.Reindex().
		Body(body).
		Refresh("true").
		RequestsPerSecond(-1).
		Slices(5).
		WaitForCompletion(false).
		Timeout("3h").
		Header("wait_for_completion_timeout", "3h").
		Header("keep_on_completion", "true").
		Header("keep_alive", "3h").
		DoAsync(context.Background())

	if err != nil {
		return err
	}

	util.Log().Info("开始重建索引，任务ID: %s", res.TaskId)

	taskID := res.TaskId
	for {
		time.Sleep(5 * time.Second)
		task, err := ESClient.TasksGetTask().TaskId(taskID).Do(context.Background())
		if err != nil {
			return err
		}

		if task.Error != nil {
			return fmt.Errorf("reindex task failed: %s", task.Error.Reason)
		}

		status := task.Task.Status.(map[string]interface{})

		// 添加更多状态监控
		created := int(status["created"].(float64))
		updated := int(status["updated"].(float64))
		deleted := int(status["deleted"].(float64))
		timeInSeconds := float64(task.Task.RunningTimeInNanos) / 1e9

		// 实时检查目标索引的文档数
		destCount, _ := ESClient.Count(newIndex).Do(context.Background())

		if task.Completed {
			util.Log().Info("索引重建完成 - 源索引总数: %d, 目标索引文档数: %d\n"+
				"已创建: %d, 已更新: %d, 已删除: %d\n"+
				"耗时: %.0f秒, 速度: %.0f文档/秒",
				totalDocs, destCount,
				created, updated, deleted,
				timeInSeconds,
				float64(created+updated)/timeInSeconds)
			break
		}

		util.Log().Info("索引重建进度 - 源索引总数: %d, 目标索引当前文档数: %d\n"+
			"已创建: %d, 已更新: %d, 已删除: %d\n"+
			"进度: %.1f%%, 耗时: %.0f秒, 速度: %.0f文档/秒",
			totalDocs, destCount,
			created, updated, deleted,
			float64(created+updated)/float64(totalDocs)*100,
			timeInSeconds,
			float64(created+updated)/timeInSeconds)
	}
	return nil
}

// ReindexWithProgressBatch 使用批处理方式重建索引
func ReindexWithProgressBatch(sourceIndex, destIndex string, batchSize int) error {
	rangeQuery := elastic.NewRangeQuery("transaction_time").
		Gte("now-25h").
		Lte("now")

	count, err := ESClient.Count(sourceIndex).
		Query(rangeQuery).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("count documents failed: %w", err)
	}
	util.Log().Info("开始重建索引，总文档数: %d", count)

	ctx := context.Background()
	scroll := ESClient.Scroll(sourceIndex).
		Query(rangeQuery).
		Size(batchSize)

	defer scroll.Clear(ctx) // 确保清理 scroll

	bulkService := ESClient.Bulk().Refresh("true")
	total := 0
	totalStart := time.Now()
	batchStart := time.Now()
	for {
		results, err := scroll.Do(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("scroll failed: %w", err)
		}

		for _, hit := range results.Hits.Hits {
			var source map[string]interface{}
			if err := json.Unmarshal(hit.Source, &source); err != nil {
				return fmt.Errorf("unmarshal document failed: %w", err)
			}

			// 安全地获取
			docID, ok := source["id"].(string)
			if !ok {
				util.Log().Info("document missing id or invalid type: %v", source)
				continue
			}

			bulkReq := elastic.NewBulkIndexRequest().
				Index(destIndex).
				Id(docID).
				OpType("create").
				Doc(source)
			bulkService.Add(bulkReq)
		}

		if bulkService.NumberOfActions() >= batchSize {
			resp, err := bulkService.Do(ctx)
			if err != nil {
				return fmt.Errorf("bulk index failed: %w", err)
			}

			if resp.Errors {
				// 记录错误但继续处理
				for _, item := range resp.Failed() {
					util.Log().Error("bulk item failed: ID=%s, Error=%s", item.Id, item.Error.Reason)
				}
			}

			batchDuration := time.Since(batchStart).Seconds()
			totalDuration := time.Since(totalStart).Seconds()
			total += bulkService.NumberOfActions()
			progress := float64(total) / float64(count) * 100

			util.Log().Info("进度: %.2f%%, 已处理: %d/%d, 本批次耗时: %.2f秒, 总耗时: %.2f秒",
				progress, total, count, batchDuration, totalDuration)

			bulkService = ESClient.Bulk().Refresh("true")
			batchStart = time.Now()
		}
	}

	// 处理剩余档
	if bulkService.NumberOfActions() > 0 {
		resp, err := bulkService.Do(ctx)
		if err != nil {
			return fmt.Errorf("final bulk index failed: %w", err)
		}
		if resp.Errors {
			for _, item := range resp.Failed() {
				util.Log().Error("bulk item failed: ID=%s, Error=%s", item.Id, item.Error.Reason)
			}
		}
		total += bulkService.NumberOfActions()
	}

	totalDuration := time.Since(totalStart).Seconds()
	util.Log().Info("索引重建完成，总文档数: %d, 总耗时: %.2f秒", total, totalDuration)
	return nil
}

// DeleteIndex deletes an Elasticsearch index by name
func DeleteIndex(indexName string) error {
	resp, err := ESClient.DeleteIndex(indexName).Do(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete index %s: %v", indexName, err)
	}
	if !resp.Acknowledged {
		return fmt.Errorf("delete index %s not acknowledged", indexName)
	}
	return nil
}

func CreateAlias(aliasName, indexName string) error {
	_, err := ESClient.Alias().
		Add(indexName, aliasName).
		Do(context.Background())
	return err
}

// convertToISODuration 将简单时间格式转换为 ISO-8601 格式
func convertToISODuration(duration string) string {
	// 解析数字和单位
	var number int
	var unit string
	fmt.Sscanf(duration, "%d%s", &number, &unit)

	switch unit {
	case "m", "min":
		return fmt.Sprintf("PT%dM", number)
	case "h", "hour":
		return fmt.Sprintf("PT%dH", number)
	case "d", "day":
		return fmt.Sprintf("P%dD", number)
	case "w", "week":
		return fmt.Sprintf("P%dW", number)
	case "M", "month":
		return fmt.Sprintf("P%dM", number)
	case "y", "year":
		return fmt.Sprintf("P%dY", number)
	default:
		// 如果无法识别，返回原始值
		return duration
	}
}
