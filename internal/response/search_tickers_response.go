package response

// MarketResponse 定义市场响应数据结构
type MarketResponse struct {
	Market         Market         `json:"market"`          // 市场数据
	MarketMetadata MarketMetadata `json:"market_metadata"` // 市场元数据
	MarketTicker   MarketTicker   `json:"market_ticker"`   // 市场 Ticker 数据
}

// MarketListResponse 定义市场列表响应数据结构
type SearchTickerResponse struct {
	Posts   []interface{}    `json:"posts"`    // 帖子列表
	Users   []interface{}    `json:"users"`    // 用户列表
	Tickers []MarketResponse `json:"tickers"`  // 市场列表
	Cursor  string           `json:"cursor"`   // 游标
	HasMore bool             `json:"has_more"` // 是否有更多数据
}
