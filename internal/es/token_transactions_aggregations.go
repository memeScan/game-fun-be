package es

import (
	"fmt"
)

// AddMarketCapAggregation 添加市值聚合
func (b *ESQueryBuilder) AddMarketCapAggregation(name string, timeRange string, parentAgg string) *ESQueryBuilder {
	marketCapAgg := map[string]interface{}{
		"filter": map[string]interface{}{
			"range": map[string]interface{}{
				"transaction_time": map[string]interface{}{
					"gte": fmt.Sprintf("now-%s", timeRange),
				},
			},
		},
		"aggs": map[string]interface{}{
			"latest_transaction": map[string]interface{}{
				"top_hits": map[string]interface{}{
					"_source": map[string]interface{}{
						"includes": []string{"market_cap"},
					},
					"size": 1,
					"sort": []map[string]interface{}{
						{
							"transaction_time": map[string]interface{}{
								"order": "asc",
							},
						},
					},
				},
			},
		},
	}

	if parentAgg != "" {
		if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
			if parent["aggs"] == nil {
				parent["aggs"] = make(map[string]interface{})
			}
			parent["aggs"].(map[string]interface{})[name] = marketCapAgg
		}
	} else {
		b.aggregations[name] = marketCapAgg
	}

	return b
}

// AddLatestMarketCapAggregation 添加获取最新市值的聚合
// name: 聚合名称
// parentAgg: 父聚合名称（可选，如果为空则添加到根级别）
func (b *ESQueryBuilder) AddLatestMarketCapAggregation(name string, parentAgg string) *ESQueryBuilder {
	latestMarketCapAgg := map[string]interface{}{
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
				"includes": []string{"market_cap"},
			},
		},
	}

	if parentAgg != "" {
		if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
			if parent["aggs"] == nil {
				parent["aggs"] = make(map[string]interface{})
			}
			parent["aggs"].(map[string]interface{})[name] = latestMarketCapAgg
		}
	} else {
		b.aggregations[name] = latestMarketCapAgg
	}

	return b
}

// AddVolumeAggregation 添加交易量聚合
// name: 聚合名称
// parentAgg: 父聚合名称(通常是 "unique_tokens")
func (b *ESQueryBuilder) AddVolumeAggregation(name string, parentAgg string) *ESQueryBuilder {
	volumeAgg := map[string]interface{}{
		"sum": map[string]interface{}{
			"script": map[string]interface{}{
				"source": "doc['token_amount'].size() > 0 ? Double.parseDouble(doc['token_amount'].value) : 0",
			},
		},
	}

	// 添加到父聚合
	if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
		if parent["aggs"] == nil {
			parent["aggs"] = make(map[string]interface{})
		}
		parent["aggs"].(map[string]interface{})[name] = volumeAgg
	}

	return b
}

// AddMaxAggregation 添加一个获取字段最大值的聚合
// name: 聚合名称（可选，如果为空则使用默认名称）
// field: 需要获取最大值的字段名
// parentAgg: 父聚合名称（可选，如果为空则添加到根级别）
func (b *ESQueryBuilder) AddMaxAggregation(name string, field string, parentAgg string) *ESQueryBuilder {
	if name == "" {
		name = fmt.Sprintf("max_%s", field)
	}

	maxAgg := map[string]interface{}{
		"max": map[string]interface{}{
			"field": field,
		},
	}

	if parentAgg != "" {
		if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
			if parent["aggs"] == nil {
				parent["aggs"] = make(map[string]interface{})
			}
			parent["aggs"].(map[string]interface{})[name] = maxAgg
		}
	} else {
		b.aggregations[name] = maxAgg
	}

	return b
}

// AddBuysAggregation 添加买入聚合
// name: 聚合名称
// parentAgg: 父聚合名称（可选，如果为空则添加到根级别）
// isBuy: 是否买入
// timeRange: 时间范围（例如：1h, 24h）
func (b *ESQueryBuilder) AddSwapCountAggregation(name string, parentAgg string, isBuy bool, timeRange string) *ESQueryBuilder {
	value_name := "buy_volume"
	if !isBuy {
		value_name = "sell_volume"
	}
	// 构建基础聚合
	agg := map[string]interface{}{
		"filter": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"transaction_time": map[string]interface{}{
								"gte": fmt.Sprintf("now-%s", timeRange),
							},
						},
					},
					{
						"term": map[string]interface{}{
							"is_buy": isBuy,
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			value_name: map[string]interface{}{
				"value_count": map[string]interface{}{
					"field": "transaction_hash.keyword",
				},
			},
		},
	}

	// 添加到父聚合或根级别
	if parentAgg != "" {
		if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
			if _, exists := parent["aggs"]; !exists {
				parent["aggs"] = make(map[string]interface{})
			}
			parent["aggs"].(map[string]interface{})[name] = agg
		}
	} else {
		b.aggregations[name] = agg
	}

	return b
}

// AddHolderBalanceAggregation 添加持有者余额和数量聚合
// name: 聚合名称
// parentAgg: 父聚合名称（可选，如果为空则添加到根级别）
func (b *ESQueryBuilder) AddHolderBalanceAggregation(name string, parentAgg string) *ESQueryBuilder {
	// 如果名称为空，使用默认名称
	if name == "" {
		name = "holder_stats"
	}

	// 构建持有者统计聚合
	holderStatsAgg := map[string]interface{}{
		"filter": map[string]interface{}{
			"term": map[string]interface{}{
				"is_buy": true,
			},
		},
		"aggs": map[string]interface{}{
			"holders": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "user_address.keyword",
					"size":  0, // 获取所有用户
				},
				"aggs": map[string]interface{}{
					"token_balance": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "token_amount",
						},
					},
					"filter_positive": map[string]interface{}{
						"bucket_selector": map[string]interface{}{
							"buckets_path": map[string]string{
								"balance": "token_balance",
							},
							"script": "params.balance >= 1",
						},
					},
				},
			},
		},
	}

	// 添加到父聚合或根级别
	if parentAgg != "" {
		if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
			if parent["aggs"] == nil {
				parent["aggs"] = make(map[string]interface{})
			}
			aggs := parent["aggs"].(map[string]interface{})

			// 添加持有者余额聚合
			aggs[name] = holderStatsAgg

			// 添加持有者数量聚合
			aggs["holder_count"] = map[string]interface{}{
				"filter": map[string]interface{}{
					"term": map[string]interface{}{
						"is_buy": true,
					},
				},
				"aggs": map[string]interface{}{
					"unique_users": map[string]interface{}{
						"cardinality": map[string]interface{}{
							"field": "user_address.keyword",
						},
					},
				},
			}
		}
	} else {
		b.aggregations[name] = holderStatsAgg
		b.aggregations["holder_count"] = map[string]interface{}{
			"filter": map[string]interface{}{
				"term": map[string]interface{}{
					"is_buy": true,
				},
			},
			"aggs": map[string]interface{}{
				"unique_users": map[string]interface{}{
					"cardinality": map[string]interface{}{
						"field": "user_address.keyword",
					},
				},
			},
		}
	}

	return b
}

func (b *ESQueryBuilder) AddLatestStatusFilterAggregation(name string, parentAgg string, isBurnt bool, isDexScreenerSpent bool) *ESQueryBuilder {
	if name == "" {
		name = "status_filter"
	}

	// 构建过滤条件
	mustConditions := make([]map[string]interface{}, 0)

	if isBurnt {
		mustConditions = append(mustConditions, map[string]interface{}{
			"term": map[string]interface{}{
				"is_burned_lp": true,
			},
		})
	}

	if isDexScreenerSpent {
		mustConditions = append(mustConditions, map[string]interface{}{
			"term": map[string]interface{}{
				"is_dex_ad": true,
			},
		})
	}

	// 创建 filter 聚合
	filterAgg := map[string]interface{}{
		"filter": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
	}

	// 添加到父聚合
	if parentAgg != "" {
		if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
			if parent["aggs"] == nil {
				parent["aggs"] = make(map[string]interface{})
			}
			parent["aggs"].(map[string]interface{})[name] = filterAgg
		}
	} else {
		b.aggregations[name] = filterAgg
	}

	return b
}

// AddTopHoldersAggregation 添加前10大持仓者聚合
func (b *ESQueryBuilder) AddTopHoldersAggregation(name string, parentAgg string) *ESQueryBuilder {
	if name == "" {
		name = "holders"
	}

	// 构建持仓者聚合
	holdersAgg := map[string]interface{}{
		"terms": map[string]interface{}{
			"field": "user_address.keyword",
			"size":  10,
			"order": map[string]interface{}{
				"balance_sum": "desc",
			},
		},
		"aggs": map[string]interface{}{
			"balance_sum": map[string]interface{}{
				"sum": map[string]interface{}{
					"script": map[string]interface{}{
						"source": "doc['is_buy'].value ? Double.parseDouble(doc['token_amount'].value) : -Double.parseDouble(doc['token_amount'].value)",
					},
				},
			},
			"token_supply": map[string]interface{}{
				"max": map[string]interface{}{
					"script": map[string]interface{}{
						"source": "Double.parseDouble(doc['token_supply'].value)",
					},
				},
			},
			"holder_percentage": map[string]interface{}{
				"bucket_script": map[string]interface{}{
					"buckets_path": map[string]interface{}{
						"balance": "balance_sum",
						"supply":  "token_supply",
					},
					"script": "params.balance / params.supply",
				},
			},
		},
	}

	// 添加到父聚合
	if parentAgg != "" {
		if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
			if parent["aggs"] == nil {
				parent["aggs"] = make(map[string]interface{})
			}
			parent["aggs"].(map[string]interface{})[name] = holdersAgg
		}
	} else {
		b.aggregations[name] = holdersAgg
	}

	return b
}

// AddTopHoldersTotalPercentageAggregation 添加前10大持仓总占比聚合
func (b *ESQueryBuilder) AddTopHoldersTotalPercentageAggregation(name string, parentAgg string) *ESQueryBuilder {
	if name == "" {
		name = "total_holders_percentage"
	}

	// 构建总占比聚合
	totalPercentageAgg := map[string]interface{}{
		"sum_bucket": map[string]interface{}{
			"buckets_path": "holders>holder_percentage", // 直接使用固定路径，因为它总是依赖于holders聚合
		},
	}

	// 添加到父聚合
	if parentAgg != "" {
		if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
			if parent["aggs"] == nil {
				parent["aggs"] = make(map[string]interface{})
			}
			parent["aggs"].(map[string]interface{})[name] = totalPercentageAgg
		}
	} else {
		b.aggregations[name] = totalPercentageAgg
	}

	return b
}
