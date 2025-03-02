package model

// TokenBaseInfo 代币基础信息
type TokenBaseInfo struct {
	Address string `json:"address"`
	// 池地址
	PoolAddress         string `json:"pool_address"`
	Symbol              string `json:"symbol"`
	URI                 string `json:"uri"`
	Name                string `json:"name"`
	Creator             string `json:"creator"`
	Decimals            int    `json:"decimals"`
	ChainType           uint8  `json:"chain_type"`
	Twitter             string `json:"twitter"`
	Website             string `json:"website"`
	Telegram            string `json:"telegram"`
	IsComplete          bool   `json:"is_complete"`
	CreatedPlatformType uint8  `json:"created_platform_type"`
}
