package response

// WalletHolding 定义代币持仓数据结构
type WalletHolding struct {
	TokenName       string `json:"token_name"`        // 代币名称
	Symbol          string `json:"symbol"`            // 代币符号
	Price           string `json:"price"`             // 当前价格
	ImageURI        string `json:"image_uri"`         // 代币图片 URI
	Balance         string `json:"balance"`           // 持仓数量
	TotalValue      string `json:"total_value"`       // 持仓总价值
	ID              int    `json:"id"`                // 代币 ID
	MarketID        int    `json:"market_id"`         // 市场 ID
	HoldersCount    int    `json:"holders_count"`     // 持有者数量
	FilledPrice     string `json:"filled_price"`      // 成交价格
	RealizedPNL     string `json:"realized_pnl"`      // 已实现盈亏
	MarketAddress   string `json:"market_address"`    // 市场地址
	TotalBuy        string `json:"total_buy"`         // 总买入数量
	TotalBuyNative  string `json:"total_buy_native"`  // 总买入原生代币数量
	TotalSell       string `json:"total_sell"`        // 总卖出数量
	TotalSellNative string `json:"total_sell_native"` // 总卖出原生代币数量
}

// TokenHoldingsResponse 定义代币持仓响应数据结构
type WalletHoldingsResponse struct {
	TokenHoldings []WalletHolding `json:"token_holdings"` // 代币持仓列表
}
