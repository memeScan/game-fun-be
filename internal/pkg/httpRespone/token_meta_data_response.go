package httpRespone

type TokenMetaDataResponse struct {
	Data    map[string]TokenMetaData `json:"data"`
	Success bool                     `json:"success"`
}

type TokenMetaData struct {
	Address    string      `json:"address"`
	Symbol     string      `json:"symbol"`
	Name       string      `json:"name"`
	Decimals   uint8       `json:"decimals"`
	Extensions *Extensions `json:"extensions"`
	LogoURI    string      `json:"logo_uri"`
}

type Extensions struct {
	CoingeckoID *string `json:"coingecko_id"`
	Website     *string `json:"website"`
	Twitter     *string `json:"twitter"`
	Discord     *string `json:"discord"`
	Medium      *string `json:"medium"`
	Telegram    *string `json:"telegram"` // 可能为 null，使用指针类型
	Description *string `json:"description"`
	SerumV3USDC *string `json:"serum_v3_usdc"`
	SerumV3USDT *string `json:"serum_v3_usdt"`
	Github      *string `json:"github"`
}
