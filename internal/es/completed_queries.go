package es

import (
	"encoding/json"
	"my-token-ai-be/internal/request"
)

func CompletedQuery(req *request.SolRankRequest) (string, error) {

	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"created_platform_type": "1",
						},
					},
					{
						"term": map[string]interface{}{
							"is_complete": true,
						},
					},
					{
						"range": map[string]interface{}{
							"transaction_time": map[string]interface{}{
								"gte": "now-" + req.Time,
							},
						},
					},
					{
						"range": map[string]interface{}{
							"block_time": map[string]interface{}{
								"gte": "now-24h",
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
					"volume": map[string]interface{}{
						"sum": map[string]interface{}{
							"script": map[string]interface{}{
								"source": "doc['token_amount'].size() > 0 ? Double.parseDouble(doc['token_amount'].value) : 0",
							},
						},
					},
					"buys": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
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
					"sells": map[string]interface{}{
						"filter": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
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
					"max_token_create_time": map[string]interface{}{
						"max": map[string]interface{}{
							"field": "block_time",
						},
					},
					"market_cap_1m": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1m",
								},
							},
						},
						"aggs": map[string]interface{}{
							"latest_transaction": map[string]interface{}{
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
										"includes": []string{"market_cap"},
									},
								},
							},
						},
					},
					"market_cap_5m": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-5m",
								},
							},
						},
						"aggs": map[string]interface{}{
							"latest_transaction": map[string]interface{}{
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
										"includes": []string{"market_cap"},
									},
								},
							},
						},
					},
					"market_cap_1h": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-1h",
								},
							},
						},
						"aggs": map[string]interface{}{
							"latest_transaction": map[string]interface{}{
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
										"includes": []string{"market_cap"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	boolQuery := query["query"].(map[string]interface{})["bool"].(map[string]interface{})

	filterClauses := boolQuery["filter"].([]map[string]interface{})

	if req.MinCreated != nil {
		filterClauses = addRangeFilter(filterClauses, "block_time", "lte", "now-"+*req.MinCreated)
	}
	if req.MaxCreated != nil {
		filterClauses = addRangeFilter(filterClauses, "block_time", "gte", "now-"+*req.MaxCreated)
	}

	boolQuery["filter"] = filterClauses

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	return string(queryJSON), nil
}
