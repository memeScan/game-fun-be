package es

import (
	"encoding/json"
)

// 根据token_address查询到代币列表
func TokenAnalyticsQuery(tokenAddress string, chainType uint8) (string, error) {

	query := map[string]interface{}{
		"size": 0, // 不需要返回具体的文档，只需要聚合结果
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{},
				"filter": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"token_address.keyword": tokenAddress,
						},
					},
					{
						"term": map[string]interface{}{
							"chain_type": chainType,
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"unique_tokens": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "token_address.keyword",
				},
				"aggs": map[string]interface{}{

					"latest_transaction": map[string]interface{}{
						"top_hits": map[string]interface{}{
							"size": 1,
							"sort": []map[string]interface{}{
								{
									"transaction_time": map[string]interface{}{
										"order": "desc",
									},
								},
							},
						},
					},
					"last_transaction_1h_price": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1h",
								},
							},
						},
						"aggs": map[string]interface{}{
							"latest": map[string]interface{}{
								"top_hits": map[string]interface{}{
									"size": 1,
									"sort": []map[string]interface{}{
										{
											"transaction_time": map[string]interface{}{
												"order": "asc",
											},
										},
									},
									"_source": map[string]interface{}{
										"includes": []string{"price"}, // 只返回价格字段
									},
								},
							},
						},
					},
					"last_transaction_4h_price": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-4h",
								},
							},
						},
						"aggs": map[string]interface{}{
							"latest": map[string]interface{}{
								"top_hits": map[string]interface{}{
									"size": 1,
									"sort": []map[string]interface{}{
										{
											"transaction_time": map[string]interface{}{
												"order": "asc",
											},
										},
									},
									"_source": map[string]interface{}{
										"includes": []string{"price"}, // 只返回价格字段
									},
								},
							},
						},
					},
				},
			},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return "", err
	}	

	return string(queryBytes), nil
}