package response

// TokenHolding 定义代币持仓数据结构
type TokenHolding struct {
	TokenName    string `json:"token_name"`    // 代币名称
	Symbol       string `json:"symbol"`        // 代币符号
	Price        string `json:"price"`         // 当前价格
	ImageURI     string `json:"image_uri"`     // 代币图片 URI
	Balance      string `json:"balance"`       // 持仓数量
	TotalValue   string `json:"total_value"`   // 持仓总价值
	MarketCap    string `json:"market_cap"`    // 市值
	ID           int64  `json:"id"`            // 代币 ID
	HoldersCount int    `json:"holders_count"` // 持有者数量
	Profit       string `json:"profit"`        // 收益
	ProfitRate   string `json:"profit_rate"`   // 收益率
}

// TokenHoldingsResponse 定义代币持仓响应数据结构
type TokenHoldingsResponse struct {
	CurrentHolding       []TokenHolding `json:"current_holding"`        // 当前持仓
	HistoryTokenHoldings []TokenHolding `json:"history_token_holdings"` // 历史持仓
}
