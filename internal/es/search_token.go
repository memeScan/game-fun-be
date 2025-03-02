package es

import (
	"encoding/json"
	"errors"
	"strings"
)

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
