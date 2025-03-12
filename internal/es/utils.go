package es

import (
	"encoding/json"
)

// UnmarshalAggregationResult unmarshals the JSON result into an AggregationResult struct
func UnmarshalAggregationResult(result []byte) (*AggregationResult, error) {
	var aggregationResult AggregationResult
	if err := json.Unmarshal(result, &aggregationResult); err != nil {
		return nil, err
	}
	return &aggregationResult, nil
}
