package response

type TokenAnalytisResponse struct {
	TokenAddress string `json:"token_address"`
	Price float64 `json:"price"`
	CreatedPlatformType uint8 `json:"created_platform_type"` // The platform type where the token was created
	IsComplete bool `json:"is_complete"` // Flag indicating if the token is a common token
	// 1h价格变化
	PriceChange1h float64 `json:"price_change_1h"`
	// 4h价格变化
	PriceChange4h float64 `json:"price_change_4h"`
}

