package response

import (
	"github.com/shopspring/decimal"
)

// TokenOrderBookItem represents a single transaction record in the order book.
// @Description Single transaction record in the token order book
type TokenOrderBookItem struct {
	TransactionHash  string          `json:"transaction_hash"`   // The unique hash of the transaction
	ChainType        uint8           `json:"chain_type"`         // The type of blockchain network
	UserAddress      string          `json:"user_address"`       // The address of the user who made the transaction
	TokenAddress     string          `json:"token_address"`      // The address of the token being traded
	PoolAddress      string          `json:"pool_address"`       // The address of the liquidity pool
	BaseTokenAmount  decimal.Decimal `json:"base_token_amount"`  // The amount of base tokens in the transaction
	QuoteTokenAmount decimal.Decimal `json:"quote_token_amount"` // The amount of quote tokens in the transaction
	BaseTokenPrice   decimal.Decimal `json:"base_token_price"`   // The price of the base token
	QuoteTokenPrice  decimal.Decimal `json:"quote_token_price"`  // The price of the quote token
	TransactionType  uint8           `json:"transaction_type"`   // The type of transaction (buy/sell)
	PlatformType     uint8           `json:"platform_type"`      // The type of platform where the transaction occurred
	TransactionTime  int64           `json:"transaction_time"`   // The timestamp of the transaction in unix seconds
	UsdAmount        decimal.Decimal `json:"usd_amount"`         // The transaction amount in USD
}

// BuildTokenOrderBookResponse creates a standardized response for token order book data
func BuildTokenOrderBookResponse(transactions []TokenOrderBookItem) Response {
	return Response{
		Code: CodeSuccess,
		Data: transactions,
		Msg:  "Success",
	}
}
