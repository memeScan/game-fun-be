package response

// TickerResponse 市场行情数据
// @Description 市场行情的具体数据
type TickersResponse struct {
	List    []TickerItem `json:"list"`                                              // 市场列表
	HasMore bool         `json:"has_more" example:"true"`                           // 是否有更多数据
	Cursor  string       `json:"cursor" example:"NTA0MTcuMzIwNzI5MDAwMDAwMDAwOjE="` // 分页游标
}

// TickerItem 单个市场行情数据
// @Description 单个市场的详细信息
type TickerItem struct {
	Market         Market         `json:"market"`          // 市场信息
	MarketMetadata MarketMetadata `json:"market_metadata"` // 市场元数据
	MarketTicker   MarketTicker   `json:"market_ticker"`   // 市场行情
}
