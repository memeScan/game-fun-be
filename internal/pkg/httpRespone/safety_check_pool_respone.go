package httpRespone

type SafetyCheckPoolData struct {
	Mint               string  `json:"mint"`
	PoolAddress        string  `json:"poolAddress"`
	Top10Holdings      float64 `json:"top10Holdings"`
	LpBurnedPercentage float64 `json:"lpBurnedPercentage"`
	Holders            int     `json:"holders"`
}

type SafetyCheckPoolResponse struct {
	Code    int                   `json:"code"`
	Data    []SafetyCheckPoolData `json:"data"`
	Message string                `json:"message"`
}
