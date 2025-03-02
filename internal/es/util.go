package es

import "encoding/json"

func addRangeFilter(filterClauses []map[string]interface{}, field string, operator string, value interface{}) []map[string]interface{} {
	filterClauses = append(filterClauses, map[string]interface{}{
		"range": map[string]interface{}{
			field: map[string]interface{}{operator: value},
		},
	})
	return filterClauses
}


// UnmarshalAggregationResult unmarshals the JSON result into an AggregationResult struct
func 	UnmarshalAggregationResult(result []byte) (*AggregationResult, error) {
	var aggregationResult AggregationResult
	if err := json.Unmarshal(result, &aggregationResult); err != nil {
		return nil, err
	}
	return &aggregationResult, nil
}