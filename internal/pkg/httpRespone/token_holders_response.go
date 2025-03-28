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

// Value 结构体表示代币持有者的资产价值
type Value struct {
	Quote float64 `json:"quote"`
	USD   float64 `json:"usd"`
}

// Holder 结构体表示代币持有者的信息
type Top20Holders struct {
	Address    string  `json:"address"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
	Value      Value   `json:"value"`
}
