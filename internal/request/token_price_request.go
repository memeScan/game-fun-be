package request

// TokenPriceRequest represents the request structure for fetching token prices
type TokenPriceRequest struct {
	Addresses []string `json:"addresses" binding:"required"`
} 