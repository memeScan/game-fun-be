package request

type SwapRequest struct {
	SwapTransaction string `json:"swap_transaction" binding:"required"`
	PlatformType    string `json:"platform_type" binding:"required"`
	IsAntiMEV       bool   `json:"is_anti_mev"`
	StartTime       string `json:"start_time"`
}
