package response

// TickerData Ticker 详情数据
// @Description Ticker 的具体数据
type GetTickerResponse struct {
	Market         Market         `json:"market"`          // 市场信息
	MarketMetadata MarketMetadata `json:"market_metadata"` // 市场元数据
	MarketTicker   MarketTicker   `json:"market_ticker"`   // 市场行情
}
