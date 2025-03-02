
package clickhouse

type TokenMarketAnalytics struct {
    TokenAddress   string  `json:"token_address"`
    CurrentPrice   float64 `json:"current_price"`
    PriceChange1m  float64 `json:"price_change_1m"`
    PriceChange5m  float64 `json:"price_change_5m"`
    PriceChange1h  float64 `json:"price_change_1h"`
    PriceChange24h float64 `json:"price_change_24h"`
    Volume1m       float64 `json:"volume_1m"`
    Volume5m       float64 `json:"volume_5m"`
    Volume1h       float64 `json:"volume_1h"`
    Volume24h      float64 `json:"volume_24h"`
    BuyCount1m     uint64     `json:"buy_count_1m"`
    BuyCount5m     uint64     `json:"buy_count_5m"`
    BuyCount1h     uint64     `json:"buy_count_1h"`
    BuyCount24h    uint64     `json:"buy_count_24h"`
    SellCount1m    uint64     `json:"sell_count_1m"`
    SellCount5m    uint64     `json:"sell_count_5m"`
    SellCount1h    uint64     `json:"sell_count_1h"`
    SellCount24h   uint64     `json:"sell_count_24h"`
    TotalCount1m   uint64     `json:"total_count_1m"`
    TotalCount5m   uint64     `json:"total_count_5m"`
    TotalCount1h   uint64     `json:"total_count_1h"`
    TotalCount24h  uint64     `json:"total_count_24h"`
    BuyVolume1m    float64 `json:"buy_volume_1m"`
    BuyVolume5m    float64 `json:"buy_volume_5m"`
    BuyVolume1h    float64 `json:"buy_volume_1h"`
    BuyVolume24h   float64 `json:"buy_volume_24h"`
    SellVolume1m   float64 `json:"sell_volume_1m"`
    SellVolume5m   float64 `json:"sell_volume_5m"`
    SellVolume1h   float64 `json:"sell_volume_1h"`
    SellVolume24h  float64 `json:"sell_volume_24h"`
}