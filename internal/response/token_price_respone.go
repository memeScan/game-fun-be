package response

// Price represents the structure of the price response
type Price struct {
	Price float64 `json:"price"`
}

// BuildSolPriceResponse builds the response for Solana price
func BuildSolPriceResponse(price Price) Response {
	return Response{
		Code: CodeSuccess,
		Data: price,
		Msg:  "Success",
	}
} 