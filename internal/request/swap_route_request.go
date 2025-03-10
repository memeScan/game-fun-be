package request

import "github.com/shopspring/decimal"

type SwapRouteRequest struct {
	TokenAddress    string          `form:"token_address" binding:"required" example:"So11111111111111111111111111111111111111112"`
	FromAddress     string          `form:"from_address" binding:"required" example:"CN8R1aHNWLAZm99ymTCd3asErjc2fhe5471cRXs7nJ3m"`
	TokenInAddress  string          `form:"token_in_address" binding:"required" example:"So11111111111111111111111111111111111111112"`
	TokenOutAddress string          `form:"token_out_address" binding:"required" example:"FfYhzJ7j3rrs4m4i1wKy5Bz5aYW8mKEGq2rxChU3pump"`
	TokenInChain    string          `form:"token_in_chain" binding:"required" example:"sol"`
	TokenOutChain   string          `form:"token_out_chain" binding:"required" example:"sol"`
	InAmount        decimal.Decimal `form:"in_amount" binding:"required" example:"100000000"`
	PriorityFee     uint64          `form:"priorityFee" binding:"required" example:"200000000"`
	Slippage        string          `form:"slippage" binding:"required" example:"100 * 100"`
	IsAntiMev       bool            `form:"is_anti_mev"`
	Legacy          bool            `form:"legacy"`
	SwapType        string          `form:"swap_type" binding:"omitempty,oneof=buy sell"`
	Points          uint64          `form:"points" binding:"omitempty" example:"200000000"`
	PlatformType    string          `form:"platform_type" binding:"required,oneof=pump raydium game g_external g_points" example:"pump"`
}
