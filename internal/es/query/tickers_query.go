package query

import (
	"game-fun-be/internal/request"

	"encoding/json"
	"errors"
	"strings"
)

func MarketQuery(req *request.TickersRequest) (string, error) {

	query := map[string]interface{}{
		"size": 0,
		// "query": map[string]interface{}{
		// 	"bool": map[string]interface{}{
		// 		"must": []map[string]interface{}{
		// 			{
		// 				"exists": map[string]interface{}{
		// 					"field": "ext_info",
		// 				},
		// 			},
		// 		},
		// 		"must_not": []map[string]interface{}{
		// 			{
		// 				"term": map[string]interface{}{
		// 					"ext_info.keyword": "",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		"aggs": map[string]interface{}{
			"unique_tokens": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "token_address.keyword",
					"size":  req.Limit,
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
					"swaps_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1d", // 过去 24 小时
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_swaps": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword", // 直接统计文档数量，更高效
								},
							},
						},
					},
					"sell_count_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1d",
											},
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
					"volume_1h": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1h",
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
					"volume_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1d",
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
					"last_transaction_24_price": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1d",
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
					"bucket_sort": map[string]interface{}{
						"bucket_sort": map[string]interface{}{
							"sort": []map[string]interface{}{
								{
									"volume_24h.total_volume": map[string]interface{}{
										"order": "desc",
									},
								},
							},
							"from": 0,
							"size": req.Limit,
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

func NewPairQuery(req *request.TickersRequest) (string, error) {

	query := map[string]interface{}{
		"size": 0,
		// "query": map[string]interface{}{
		// 	"bool": map[string]interface{}{
		// 		"must": []map[string]interface{}{
		// 			{
		// 				"exists": map[string]interface{}{
		// 					"field": "ext_info",
		// 				},
		// 			},
		// 		},
		// 		"must_not": []map[string]interface{}{
		// 			{
		// 				"term": map[string]interface{}{
		// 					"ext_info.keyword": "",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		"aggs": map[string]interface{}{
			"unique_tokens": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "token_address.keyword",
					"size":  req.Limit,
					"order": map[string]interface{}{
						"max_token_create_time": "desc",
					},
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
					"swaps_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1d", // 过去 24 小时
								},
							},
						},
						"aggs": map[string]interface{}{
							"total_swaps": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword", // 直接统计文档数量，更高效
								},
							},
						},
					},
					"sell_count_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{
										"term": map[string]interface{}{
											"is_buy": false,
										},
									},
									{
										"range": map[string]interface{}{
											"transaction_time": map[string]interface{}{
												"gte": "now-1d",
											},
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
					"volume_1h": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1h",
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
					"volume_24h": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1d",
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
									"gte": "now-1d",
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
					"max_token_create_time": map[string]interface{}{
						"max": map[string]interface{}{
							"field": "token_create_time",
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

func SearchToken(tokenAddresses []string, chainType uint8) (string, error) {
	query := map[string]interface{}{
		"size": 0, // 不需要返回具体的文档，只需要聚合结果
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"transaction_time": map[string]interface{}{
								"gte": "now-24h",
							},
						},
					},
				},
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"chain_type": chainType,
						},
					},
					{
						"bool": map[string]interface{}{
							"should": []map[string]interface{}{
								{
									"terms": map[string]interface{}{
										"token_address": tokenAddresses,
									},
								},
							},
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
					"volume": map[string]interface{}{
						"sum": map[string]interface{}{
							"script": map[string]interface{}{
								"source": "doc['token_amount'].size() > 0 ? Double.parseDouble(doc['token_amount'].value) : 0",
							},
						},
					},
					"last_transaction_24h_price": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1d",
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

func SearchTokenBySymbol(token string, chainType uint8, isTokenAddress bool) (string, error) {

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"chain_type": chainType,
						},
					},
				},
			},
		},
		"size": 2000,
	}

	queryMap, ok := query["query"].(map[string]interface{})
	if !ok {
		return "", errors.New("query is not a map")
	}

	// 确保 bool 查询存在
	boolQuery, ok := queryMap["bool"].(map[string]interface{})
	if !ok {
		return "", errors.New("bool query is not a map")
	}

	if isTokenAddress {
		boolQuery["must"] = []map[string]interface{}{
			{
				"match": map[string]interface{}{
					"token_address.keyword": map[string]interface{}{
						"query": token,
					},
				},
			},
		}
	} else {
		boolQuery["must"] = []map[string]interface{}{
			{
				"wildcard": map[string]interface{}{
					"symbol.lowercase": strings.ToLower(token) + "*",
				},
			},
		}
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	return string(queryBytes), nil
}
