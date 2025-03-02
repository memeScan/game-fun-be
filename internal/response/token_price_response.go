package response

import (
	"my-token-ai-be/internal/model"

	"github.com/shopspring/decimal"
)

// TokenPriceResponse represents the response structure for token prices
type TokenPriceResponse struct {
	Address string  `json:"address"`
	Price   float64 `json:"price"`
}

// BuildTokenPriceResponse constructs a response for a list of token prices
func BuildTokenPriceResponse(prices []model.TokenInfo) []TokenPriceResponse {
	var responses []TokenPriceResponse

	for _, info := range prices {
		// 将 uint64 转换为 decimal.Decimal，并考虑代币精度
		totalSupplyDecimal := decimal.NewFromInt(int64(info.TotalSupply))
		// 如果需要考虑代币精度
		actualTotalSupply := totalSupplyDecimal.Shift(-int32(info.Decimals))
		// 计算价格
		price := info.MarketCap.Div(actualTotalSupply).InexactFloat64()

		responses = append(responses, TokenPriceResponse{
			Address: info.TokenAddress,
			Price:   price,
		})
	}

	return responses
}
