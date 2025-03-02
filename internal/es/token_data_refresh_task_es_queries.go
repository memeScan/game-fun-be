package es

import "encoding/json"

// CompletedTokenDataRefreshTaskQuery generates a query for completed token data refresh tasks.
func CompletedTokenDataRefreshTaskQuery(limit int) (string, error) {

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
								"gte": "now-1h",
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
					"size":  limit,
					"order": map[string]interface{}{
						"max_token_create_time": "desc",
					},
				},
				"aggs": map[string]interface{}{
					"latest_transaction": map[string]interface{}{
						"top_hits": map[string]interface{}{
							"size": 1,
							"sort": []map[string]interface{}{
								{},
							},
							"_source": []string{"token_address", "pool_address"},
						},
					},
					"max_token_create_time": map[string]interface{}{
						"min": map[string]interface{}{
							"field": "block_time",
						},
					},
				},
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	return string(queryJSON), nil
}

// SwapsTokenDataRefreshTaskQuery generates a query for swaps token data refresh tasks.
func SwapsTokenDataRefreshTaskQuery(limit int, time string) (string, error) {
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"is_complete": true,
						},
					},
					{
						"range": map[string]interface{}{
							"transaction_time": map[string]interface{}{
								"gte": "now-" + time,
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
					"order": map[string]interface{}{
						"swaps": "desc",
					},
					"size": limit,
				},
				"aggs": map[string]interface{}{
					"latest_transaction": map[string]interface{}{
						"top_hits": map[string]interface{}{
							"_source": []string{"token_address", "pool_address"},
							"size":    1,
							"sort": []map[string]interface{}{
								{
									"transaction_time": map[string]interface{}{
										"order": "desc",
									},
								},
							},
						},
					},
					"swaps": map[string]interface{}{
						"filter": map[string]interface{}{
							"range": map[string]interface{}{
								"transaction_time": map[string]interface{}{
									"gte": "now-" + time,
								},
							},
						},
						"aggs": map[string]interface{}{
							"transaction_count": map[string]interface{}{
								"value_count": map[string]interface{}{
									"field": "transaction_hash.keyword",
								},
							},
						},
					},
				},
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	return string(queryJSON), nil
}
