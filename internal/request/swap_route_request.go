package request

type SwapRouteRequest struct {
	TokenInChain    string  `form:"token_in_chain" binding:"required" example:"sol"`
	TokenOutChain   string  `form:"token_out_chain" binding:"required" example:"sol"`
	FromAddress     string  `form:"from_address" binding:"required" example:"CN8R1aHNWLAZm99ymTCd3asErjc2fhe5471cRXs7nJ3m"`
	Slippage        int64 `form:"slippage" binding:"required" example:"30"`
	TokenInAddress  string  `form:"token_in_address" binding:"required" example:"So11111111111111111111111111111111111111112"`
	TokenOutAddress string  `form:"token_out_address" binding:"required" example:"FfYhzJ7j3rrs4m4i1wKy5Bz5aYW8mKEGq2rxChU3pump"`
	InAmount        string   `form:"in_amount" binding:"required" example:"100000000"`
	Fee             float64  `form:"fee" binding:"required" example:"200"`
	IsAntiMev       bool    `form:"is_anti_mev"`
	Legacy          bool    `form:"legacy"`
	} 

	