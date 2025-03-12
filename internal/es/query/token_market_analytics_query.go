package query

import (
	"encoding/json"
)

// 根据token_address查询到代币列表
func TokenMarketAnalyticsQuery(tokenAddress string, chainType uint8) (string, error) {

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
					"buy_count_1m": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1m",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": true,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"buy_volume": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword",
								},
							},
						},
					},
					"buy_count_5m": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-5m",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": true,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"buy_volume": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword",
								},
							},
						},
					},
					"buy_count_1h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1h",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": true,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"buy_volume": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword",
								},
							},
						},
					},
					"buy_count_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-24h",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": true,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"buy_volume": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword",
								},
							},
						},
					},
					"sell_count_1m": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1m",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"sell_volume": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword",
								},
							},
						},
					},
					"sell_count_5m": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-5m",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"sell_volume": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword",
								},
							},
						},
					},
					"sell_count_1h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1h",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"sell_volume": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword",
								},
							},
						},
					},
					"sell_count_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-24h",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"sell_volume": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword",
								},
							},
						},
					},
					"buy_volume_1m": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1m",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": true,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_volume": map[string]interface{}{
								"sum": map[string]interface{}{
									"script": map[string]interface{}{
										"source": "doc['native_token_amount'].size() > 0 ? Double.parseDouble(doc['native_token_amount'].value) : 0",
									},
								},
							},
						},
					},
					"buy_volume_5m": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-5m",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": true,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_volume": map[string]interface{}{
								"sum": map[string]interface{}{
									"script": map[string]interface{}{
										"source": "doc['native_token_amount'].size() > 0 ? Double.parseDouble(doc['native_token_amount'].value) : 0",
									},
								},
							},
						},
					},
					"buy_volume_1h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1h",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": true,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_volume": map[string]interface{}{
								"sum": map[string]interface{}{
									"script": map[string]interface{}{
										"source": "doc['native_token_amount'].size() > 0 ? Double.parseDouble(doc['native_token_amount'].value) : 0",
									},
								},
							},
						},
					},
					"buy_volume_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-24h",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": true,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_volume": map[string]interface{}{
								"sum": map[string]interface{}{
									"script": map[string]interface{}{
										"source": "doc['native_token_amount'].size() > 0 ? Double.parseDouble(doc['native_token_amount'].value) : 0",
									},
								},
							},
						},
					},
					"sell_volume_1m": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1m",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_volume": map[string]interface{}{
								"sum": map[string]interface{}{
									"script": map[string]interface{}{
										"source": "doc['native_token_amount'].size() > 0 ? Double.parseDouble(doc['native_token_amount'].value) : 0",
									},
								},
							},
						},
					},
					"sell_volume_5m": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-5m",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_volume": map[string]interface{}{
								"sum": map[string]interface{}{
									"script": map[string]interface{}{
										"source": "doc['native_token_amount'].size() > 0 ? Double.parseDouble(doc['native_token_amount'].value) : 0",
									},
								},
							},
						},
					},
					"sell_volume_1h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1h",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_volume": map[string]interface{}{
								"sum": map[string]interface{}{
									"script": map[string]interface{}{
										"source": "doc['native_token_amount'].size() > 0 ? Double.parseDouble(doc['native_token_amount'].value) : 0",
									},
								},
							},
						},
					},
					"sell_volume_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-24h",
											},
										},
									},
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_volume": map[string]interface{}{
								"sum": map[string]interface{}{
									"script": map[string]interface{}{
										"source": "doc['native_token_amount'].size() > 0 ? Double.parseDouble(doc['native_token_amount'].value) : 0",
									},
								},
							},
						},
					},

					// 最新的交易价格
					"last_transaction_price": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{},
						},
						"aggs": map[string]interface{}{
							"latest": map[string]interface{}{
								"top_hits": map[string]interface{}{
									"size": 1,
									"sort": []map[string]interface{}{
										{
											"transaction_time": map[string]interface{}{
												"order": "desc",
											},
										},
									},
									"_source": map[string]interface{}{
										"includes": []string{"price", "decimals", "native_price", "transaction_time", "market_cap"}, // 只返回价格字段
									},
								},
							},
						},
					},
					"last_transaction_1m_price": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1m",
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
					"last_transaction_5m_price": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-5m",
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
					"last_transaction_24h_price": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-24h",
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
