package response

// SwapHistoriesResponse 表示交易历史记录的响应结构。
type SwapHistoriesResponse struct {
	TransactionHistories []TransactionHistory `json:"transaction_histories"` // 交易历史记录列表
	HasMore              bool                 `json:"has_more"`              // 是否还有更多数据
}

// TransactionHistory 表示单笔交易历史记录。
type TransactionHistory struct {
	MarketID          int          `json:"market_id"`           // 市场 ID
	IsBuy             bool         `json:"is_buy"`              // 是否为买入
	Payer             string       `json:"payer"`               // 支付方地址
	Recipient         string       `json:"recipient"`           // 接收方地址
	Signature         string       `json:"signature"`           // 交易签名
	BlockTime         string       `json:"block_time"`          // 区块时间
	Index             int          `json:"index"`               // 交易索引
	Slot              int          `json:"slot"`                // 区块槽位
	TokenAmount       string       `json:"token_amount"`        // 代币数量
	NativeAmount      string       `json:"native_amount"`       // 原生代币数量
	Fee               string       `json:"fee"`                 // 手续费
	TotalNativeAmount string       `json:"total_native_amount"` // 总原生代币数量
	ID                int          `json:"id"`                  // 交易 ID
	PayerProfile      PayerProfile `json:"payer_profile"`       // 支付方信息
}

// PayerProfile 表示支付方的信息。
type PayerProfile struct {
	UserID   string `json:"user_id"`  // 用户 ID
	Avatar   string `json:"avatar"`   // 用户头像
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 用户昵称
}
