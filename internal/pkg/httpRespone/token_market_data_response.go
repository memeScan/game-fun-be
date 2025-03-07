package httpRespone

// TokenMarketDataResponse 结构体
// @Description 代币市场数据
type TokenMarketDataResponse struct {
	Data struct {
		Address           string  `json:"address"`
		Liquidity         float64 `json:"liquidity"`
		Price             float64 `json:"price"`
		TotalSupply       float64 `json:"total_supply"`
		CirculatingSupply float64 `json:"circulating_supply"`
		Fdv               float64 `json:"fdv"`
		MarketCap         float64 `json:"market_cap"`
	} `json:"data"`
	Success bool `json:"success"`
}
