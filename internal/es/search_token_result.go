package es

import "encoding/json"

type TokenInfo struct {
	Block                int64       `json:"block"`
	BurnPercentage       float64     `json:"burn_percentage"`
	ChainType            int         `json:"chain_type"`
	CirculatingMarketCap float64     `json:"circulating_market_cap"`
	CirculatingSupply    json.Number `json:"circulating_supply"`
	CommentCount         int         `json:"comment_count"`
	CreateTime           string      `json:"create_time"`
	CreatedPlatformType  int         `json:"created_platform_type"`
	Creator              string      `json:"creator"`
	CrownDuration        int64       `json:"crown_duration"`
	Decimals             int         `json:"decimals"`
	DevBurnPercentage    float64     `json:"dev_burn_percentage"`
	DevNativeTokenAmount int64       `json:"dev_native_token_amount"`
	DevPercentage        float64     `json:"dev_percentage"`
	DevStatus            int         `json:"dev_status"`
	DevTokenAmount       int64       `json:"dev_token_amount"`
	ExtInfo              string      `json:"ext_info"`
	Holder               int64       `json:"holder"`
	ID                   int64       `json:"id"`
	IsComplete           bool        `json:"is_complete"`
	IsMedia              bool        `json:"is_media"`
	Liquidity            float64     `json:"liquidity"`
	MarketCap            float64     `json:"market_cap"`
	NativePrice          float64     `json:"native_price"`
	PoolAddress          string      `json:"pool_address"`
	Price                float64     `json:"price"`
	Progress             float64     `json:"progress"`
	RocketDuration       int64       `json:"rocket_duration"`
	Symbol               string      `json:"symbol"`
	TokenAddress         string      `json:"token_address"`
	TokenFlags           int         `json:"token_flags"`
	TokenName            string      `json:"token_name"`
	Top10Percentage      float64     `json:"top10_percentage"`
	TotalSupply          json.Number `json:"total_supply"`
	TransactionHash      string      `json:"transaction_hash"`
	TransactionTime      string      `json:"transaction_time"`
	UpdateTime           string      `json:"update_time"`
	URI                  string      `json:"uri"`
}
