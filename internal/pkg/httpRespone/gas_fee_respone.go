package httpRespone

// PriorityFeeLevels represents the different levels of gas fees
type PriorityFeeLevels struct {
	Min      float64 `json:"min"`
	Low      float64 `json:"low"`
	Medium   float64 `json:"medium"`
	High     float64 `json:"high"`
	VeryHigh float64 `json:"veryHigh"`
	UnsafeMax float64 `json:"unsafeMax"`
}

// GasFeeResponse represents the response structure for gas fees
type GasFeeResponse struct {
	PriorityFeeLevels PriorityFeeLevels `json:"priorityFeeLevels"`
}

// GasFee represents the overall response structure
type GasFee struct {
	Code    int            `json:"code"`
	Data    GasFeeResponse `json:"data"`
	Message string         `json:"message"`
}

