package request

type RebateClaimRequest struct {
	Address         string `form:"address" binding:"required" example:"CN8R1aHNWLAZm99ymTCd3asErjc2fhe5471cRXs7nJ3m"`
	RebateAmount    uint64 `form:"rebate_amount" binding:"required" example:"200000000"`
	SwapTransaction string `json:"swap_transaction" binding:"required"`
}
