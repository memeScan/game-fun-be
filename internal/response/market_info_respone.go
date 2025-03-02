package response

import "github.com/shopspring/decimal"

// MarketInfo represents the structure of market information
type MarketInfo struct {
	Address             string          `json:"address"`
	MarketId            string          `json:"market_id"`
	PoolAddress         string          `json:"pool_address"`
	QuoteAddress        string          `json:"quote_address"`
	QuoteSymbol         string          `json:"quote_symbol"`
	BaseSymbol          string          `json:"base_symbol"`
	BaseAddress         string          `json:"base_address"`
	Liquidity           float64         `json:"liquidity"`
	BaseReserve         float64         `json:"base_reserve"`
	MarketCap           float64         `json:"market_cap"`
	TotalSupply         float64         `json:"total_supply"`
	CirculatingSupply   float64         `json:"circulating_supply"`
	QuoteReserve        float64         `json:"quote_reserve"`
	InitialLiquidity    float64         `json:"initial_liquidity"`
	InitialBaseReserve  float64         `json:"initial_base_reserve"`
	InitialQuoteReserve float64         `json:"initial_quote_reserve"`
	CreationTimestamp   int64           `json:"creation_timestamp"`
	BaseReserveValue    float64         `json:"base_reserve_value"`
	QuoteReserveValue   float64         `json:"quote_reserve_value"`
	QuoteVaultAddress   string          `json:"quote_vault_address"`
	BaseVaultAddress    string          `json:"base_vault_address"`
	Creator             string          `json:"creator"`
	Progress            decimal.Decimal `json:"progress"`
}

func BuildMarketInfoResponse(marketInfo MarketInfo) Response {
	return Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: marketInfo,
	}
}
