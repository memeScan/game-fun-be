package response

import "github.com/shopspring/decimal"

// TokenMarketAnalyticsResponse 代币市场分析数据响应结构
type TokenMarketAnalyticsResponse struct {
	TokenAddress   string          `json:"token_address"`
	CurrentPrice   float64         `json:"current_price"`
	PriceChange1m  float64         `json:"price_change_1m"`
	PriceChange5m  float64         `json:"price_change_5m"`
	PriceChange1h  float64         `json:"price_change_1h"`
	PriceChange24h float64         `json:"price_change_24h"`
	Volume1m       decimal.Decimal `json:"volume_1m"`
	Price          float64         `json:"Price"`

	Volume5m      decimal.Decimal `json:"volume_5m"`
	Volume1h      decimal.Decimal `json:"volume_1h"`
	Volume24h     decimal.Decimal `json:"volume_24h"`
	BuyCount1m    decimal.Decimal `json:"buy_count_1m"`
	BuyCount5m    decimal.Decimal `json:"buy_count_5m"`
	BuyCount1h    decimal.Decimal `json:"buy_count_1h"`
	BuyCount24h   decimal.Decimal `json:"buy_count_24h"`
	SellCount1m   decimal.Decimal `json:"sell_count_1m"`
	SellCount5m   decimal.Decimal `json:"sell_count_5m"`
	SellCount1h   decimal.Decimal `json:"sell_count_1h"`
	SellCount24h  decimal.Decimal `json:"sell_count_24h"`
	TotalCount1m  decimal.Decimal `json:"total_count_1m"`
	TotalCount5m  decimal.Decimal `json:"total_count_5m"`
	TotalCount1h  decimal.Decimal `json:"total_count_1h"`
	TotalCount24h decimal.Decimal `json:"total_count_24h"`
	BuyVolume1m   decimal.Decimal `json:"buy_volume_1m"`
	BuyVolume5m   decimal.Decimal `json:"buy_volume_5m"`
	BuyVolume1h   decimal.Decimal `json:"buy_volume_1h"`
	BuyVolume24h  decimal.Decimal `json:"buy_volume_24h"`
	SellVolume1m  decimal.Decimal `json:"sell_volume_1m"`
	SellVolume5m  decimal.Decimal `json:"sell_volume_5m"`
	SellVolume1h  decimal.Decimal `json:"sell_volume_1h"`
	SellVolume24h decimal.Decimal `json:"sell_volume_24h"`
	LastSwapAt    int64           `json:"last_swap_at"`
	MarketCap     string          `json:"market_cap"`
	Holders       int             `json:"holders" example:"4204"`
	Rank          int             `json:"rank" example:"1"`
}
