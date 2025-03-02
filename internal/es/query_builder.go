package es

import (
	"encoding/json"
)

// ESQueryBuilder 用于构建 Elasticsearch 查询的构建器
// 支持构建 bool 查询、聚合等复杂查询结构
type ESQueryBuilder struct {
	query        map[string]interface{} // 存储查询主体
	boolQuery    map[string]interface{} // 存储 bool 查询相关条件
	aggregations map[string]interface{} // 存储聚合查询
	currentAgg   string                 // 当前正在处理的聚合名称
	size         int                    // 查询返回的文档数量
}

// NewESQueryBuilder 创建并返回一个新的查询构建器实例
func NewESQueryBuilder() *ESQueryBuilder {
	return &ESQueryBuilder{
		query:        make(map[string]interface{}),
		boolQuery:    make(map[string]interface{}),
		aggregations: make(map[string]interface{}),
		size:         0,
	}
}

// SetSize 设置查询返回的文档数量
func (b *ESQueryBuilder) SetSize(size int) *ESQueryBuilder {
	b.size = size
	return b
}

// AddMustNotTerm 添加一个 must_not term 查询
// field: 字段名称
// value: 字段值
func (b *ESQueryBuilder) AddMustNotTerm(field, value string) *ESQueryBuilder {
	if b.boolQuery["must_not"] == nil {
		b.boolQuery["must_not"] = make(map[string]interface{})
	}
	b.boolQuery["must_not"].(map[string]interface{})["term"] = map[string]interface{}{
		field: value,
	}
	return b
}

// AddFilter 添加一个过滤条件
// filterType: 过滤类型（如 "term", "exists", "range" 等）
// field: 字段名称或完整的过滤条件
// value: 过滤值（可选）
func (b *ESQueryBuilder) AddFilter(filterType string, field interface{}, value ...interface{}) *ESQueryBuilder {
	if b.boolQuery["filter"] == nil {
		b.boolQuery["filter"] = make([]map[string]interface{}, 0)
	}

	filters := b.boolQuery["filter"].([]map[string]interface{})

	var filterMap map[string]interface{}

	switch filterType {
	case "term":
		filterMap = map[string]interface{}{
			"term": map[string]interface{}{
				field.(string): value[0],
			},
		}

	case "exists":
		filterMap = map[string]interface{}{
			"exists": map[string]interface{}{
				"field": field,
			},
		}

	case "range":
		if len(value) >= 2 {
			filterMap = map[string]interface{}{
				"range": map[string]interface{}{
					field.(string): map[string]interface{}{
						value[0].(string): value[1],
					},
				},
			}
		}

	default:
		// 处理自定义过滤条件
		if fieldMap, ok := field.(map[string]interface{}); ok {
			filterMap = map[string]interface{}{
				filterType: fieldMap,
			}
		}
	}

	if filterMap != nil {
		filters = append(filters, filterMap)
		b.boolQuery["filter"] = filters
	}

	return b
}

// SetTermsOrder 设置 terms 聚合的排序
// field: 排序字段
// order: 排序方式(asc/desc)
func (b *ESQueryBuilder) SetTermsOrder(field, order string) *ESQueryBuilder {
	if b.currentAgg == "" {
		return b
	}

	if agg, ok := b.aggregations[b.currentAgg].(map[string]interface{}); ok {
		if terms, ok := agg["terms"].(map[string]interface{}); ok {
			terms["order"] = map[string]interface{}{
				field: order,
			}
		}
	}
	return b
}

// AddFilterTerms 添加一个 terms 过滤条件
// field: 字段名称
// values: 过滤值列表
func (b *ESQueryBuilder) AddFilterTerms(field string, values interface{}) *ESQueryBuilder {
	if b.boolQuery["filter"] == nil {
		b.boolQuery["filter"] = make([]map[string]interface{}, 0)
	}

	filterClauses := b.boolQuery["filter"].([]map[string]interface{})
	filterClauses = append(filterClauses, map[string]interface{}{
		"terms": map[string]interface{}{
			field: values,
		},
	})

	b.boolQuery["filter"] = filterClauses
	return b
}

// AddShould 添加一个should查询条件列表
func (b *ESQueryBuilder) AddShould(shouldClauses []map[string]interface{}) *ESQueryBuilder {
	b.boolQuery["should"] = shouldClauses
	return b
}

// SetMinimumShouldMatch 设置minimum_should_match值
func (b *ESQueryBuilder) SetMinimumShouldMatch(value int) *ESQueryBuilder {
	b.boolQuery["minimum_should_match"] = value
	return b
}

// AddShouldTerms 添加一个 should terms 查询
// field: 字段名称
// values: 匹配值列表
func (b *ESQueryBuilder) AddShouldTerms(field string, values interface{}) *ESQueryBuilder {
	if b.boolQuery["should"] == nil {
		b.boolQuery["should"] = make([]map[string]interface{}, 0)
	}

	shouldClauses := b.boolQuery["should"].([]map[string]interface{})
	shouldClauses = append(shouldClauses, map[string]interface{}{
		"terms": map[string]interface{}{
			field: values,
		},
	})

	b.boolQuery["should"] = shouldClauses
	return b
}

// AddRangeFilter 添加一个范围过滤条件
// field: 字段名称
// operator: 操作符(gt/gte/lt/lte)
// value: 比较值
func (b *ESQueryBuilder) AddRangeFilter(field, operator string, value interface{}) *ESQueryBuilder {
	filters := b.boolQuery["filter"]
	if filters == nil {
		filters = make([]map[string]interface{}, 0)
	}

	filterSlice := filters.([]map[string]interface{})
	filterSlice = append(filterSlice, map[string]interface{}{
		"range": map[string]interface{}{
			field: map[string]interface{}{
				operator: value,
			},
		},
	})

	b.boolQuery["filter"] = filterSlice
	return b
}

// AddMust 添加一个 must 查询条件
func (b *ESQueryBuilder) AddMust(queryType string, field string, value interface{}) *ESQueryBuilder {
	if b.boolQuery["must"] == nil {
		b.boolQuery["must"] = make([]map[string]interface{}, 0)
	}

	mustClauses := b.boolQuery["must"].([]map[string]interface{})

	var clause map[string]interface{}
	switch queryType {
	case "term":
		clause = map[string]interface{}{
			"term": map[string]interface{}{
				field: value,
			},
		}
	case "terms":
		clause = map[string]interface{}{
			"terms": map[string]interface{}{
				field: value,
			},
		}
	}

	if clause != nil {
		mustClauses = append(mustClauses, clause)
		b.boolQuery["must"] = mustClauses
	}

	return b
}

// AddTermsAggregation 添加一个 terms 聚合
// name: 聚合名称
// field: 聚合字段
// size: 返回的桶数量
func (b *ESQueryBuilder) AddTermsAggregation(name, field string, size int) *ESQueryBuilder {
	b.currentAgg = name // 设置当前聚合名称
	b.aggregations[name] = map[string]interface{}{
		"terms": map[string]interface{}{
			"field": field,
			"size":  size,
		},
	}
	return b
}

// AddTopHitsAggregation 添加一个 top_hits 聚合
// name: 聚合名称
// size: 返回的文档数量
// sortField: 排序字段
// sortOrder: 排序方式(asc/desc)
// parentAgg: 父聚合名称（可选）
func (b *ESQueryBuilder) AddTopHitsAggregation(name string, size int, sortField string, sortOrder string, parentAgg string) *ESQueryBuilder {
	topHits := map[string]interface{}{
		"top_hits": map[string]interface{}{
			"size": size,
			"sort": []map[string]interface{}{
				{
					sortField: map[string]interface{}{
						"order": sortOrder,
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
			parent["aggs"].(map[string]interface{})[name] = topHits
		}
	} else {
		b.aggregations[name] = topHits
	}

	return b
}

// SetSource 设置要返回的字段
func (b *ESQueryBuilder) SetSource(includes []string, excludes ...[]string) *ESQueryBuilder {
	source := map[string]interface{}{
		"includes": includes,
	}

	// 如果提供了 excludes
	if len(excludes) > 0 {
		source["excludes"] = excludes[0]
	}

	// 找到最近添加的 top_hits 聚合并设置 _source
	if parent, ok := b.aggregations[b.currentAgg].(map[string]interface{}); ok {
		if aggs, ok := parent["aggs"].(map[string]interface{}); ok {
			for _, agg := range aggs {
				if aggMap, ok := agg.(map[string]interface{}); ok {
					if topHits, ok := aggMap["top_hits"].(map[string]interface{}); ok {
						topHits["_source"] = source
					}
				}
			}
		}
	}

	return b
}

// Build 构建并返回最终的查询 JSON 字符串
// 返回序列化后的查询 JSON 和可能的错误
func (b *ESQueryBuilder) Build() (string, error) {
	finalQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": b.boolQuery,
		},
	}

	finalQuery["size"] = b.size

	if len(b.aggregations) > 0 {
		finalQuery["aggs"] = b.aggregations
	}

	queryBytes, err := json.Marshal(finalQuery)
	if err != nil {
		return "", err
	}

	return string(queryBytes), nil
}

// AddBucketSelector 添加一个桶选择器过滤器
// name: 过滤器名称
// bucketsPath: 引用其他聚合的路径映射
// script: 过滤脚本
// parentAgg: 父聚合名称
func (b *ESQueryBuilder) AddBucketSelector(name string, bucketsPath map[string]string, script string, parentAgg string) *ESQueryBuilder {
	bucketSelector := map[string]interface{}{
		"bucket_selector": map[string]interface{}{
			"buckets_path": bucketsPath,
			"script":       script,
		},
	}

	// 将 bucket_selector 添加为子聚合
	if parentAgg != "" {
		if parent, ok := b.aggregations[parentAgg].(map[string]interface{}); ok {
			if parent["aggs"] == nil {
				parent["aggs"] = make(map[string]interface{})
			}
			parent["aggs"].(map[string]interface{})[name] = bucketSelector
		}
	}

	return b
}
