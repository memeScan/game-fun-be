package httpRespone

type TipFloorResponse struct {
	Code    int64 `json:"code"`
	Data    []struct {
		Time string `json:"time"`
		LandedTips25thPercentile float64 `json:"landed_tips_25th_percentile"`
		LandedTips50thPercentile float64 `json:"landed_tips_50th_percentile"`
		LandedTips75thPercentile float64 `json:"landed_tips_75th_percentile"`	
		LandedTips95thPercentile float64 `json:"landed_tips_95th_percentile"`
		LandedTips99thPercentile float64 `json:"landed_tips_99th_percentile"`
		EmaLandedTips50thPercentile float64 `json:"ema_landed_tips_50th_percentile"`
	} `json:"data"`
	Message string `json:"message"`
}
	