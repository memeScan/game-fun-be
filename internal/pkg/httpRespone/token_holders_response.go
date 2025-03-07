package httpRespone

// TokenHoldersResponse 结构体
// @Description 代币持有者信息
type TokenHoldersResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Items []TokenHolder `json:"items"`
	} `json:"data"`
}

// TokenHolder 表示单个代币持有者
// @Description 代币持有者详情
type TokenHolder struct {
	// 持有代币的数量
	Amount string `json:"amount"`
	// 代币的小数位数
	Decimals int `json:"decimals"`
	// 代币的合约地址
	Mint string `json:"mint"`
	// 持有者钱包地址
	Owner string `json:"owner"`
	// 代币账户地址
	TokenAccount string `json:"token_account"`
	// 代币的 UI 显示数量
	UIAmount float64 `json:"ui_amount"`
}
