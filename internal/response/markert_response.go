package response

import (
	"github.com/shopspring/decimal"
)

// Market 市场信息
// @Description 市场的基本信息
type Market struct {
	MarketID        int64           `json:"market_id" example:"1"`                                               // 市场ID
	Market          string          `json:"market" example:"B3o198GX45DZ9HeJZ56wuruyWWCfGVBSn5qC71q5AHZz"`       // 市场地址
	TokenMint       string          `json:"token_mint" example:"suprkbfvwpFZXzWaoKjTzGkW1nkvvwK9n2E6g1zyLFo"`    // 代币铸造地址
	Decimals        uint8           `json:"decimals" example:"6"`                                                // 代币铸造地址
	TokenVault      string          `json:"token_vault" example:"9SQy9T4pzadupDcWdYYgFmUgXLVydfQwgKZEbzeACNAQ"`  // 代币金库地址
	NativeVault     string          `json:"native_vault" example:"98J2RACW3wk8u36k1cPmDenzHdPjhcd5qXsm5g471kdh"` // 原生代币金库地址
	TokenName       string          `json:"token_name" example:"§uper Exchange"`                                 // 代币名称
	TokenSymbol     string          `json:"token_symbol" example:"SUPER"`                                        // 代币符号
	Creator         string          `json:"creator" example:"3P3PMv28AM7SvNkGmAdHCusPwahSB9QG9z5o4gTvmNBX"`      // 创建者地址
	URI             string          `json:"uri" example:"https://static.super.exchange/m/official/super.json"`   // URI
	Price           decimal.Decimal `json:"price" example:"0.000050417320729"`                                   // 当前价格
	Rank            int             `json:"rank" example:"1"`
	CreateTimestamp int64           `json:"create_timestamp" example:"1740033455"`
}

// MarketMetadata 市场元数据
// @Description 市场的元数据信息
type MarketMetadata struct {
	ImageURL    *string `json:"image_url" example:"https://static.super.exchange/i/official/super.png"` // 图片URL
	Description *string `json:"description" example:""`                                                 // 描述
	Twitter     *string `json:"twitter" example:"https://x.com/_superexchange"`                         // Twitter链接
	Website     *string `json:"website" example:"https://super.exchange"`                               // 网站链接
	Telegram    *string `json:"telegram" example:"https://t.me/SuperExchangeCommunity"`                 // Telegram链接
	Github      *string `json:"github"`
	Banner      *string `json:"banner" example:""` // 横幅图片
	Rules       *string `json:"rules" example:""`  // 规则
	Sort        *uint   `json:"sort" example:"0"`  // 排序
}

// MarketTicker 市场行情
// @Description 市场的行情数据
type MarketTicker struct {
	PriceChange24H       string `json:"price_change_24h" example:"0.224965"`
	TxCount24H           int    `json:"tx_count_24h" example:"7590"`
	BuyTxCount24H        int    `json:"buy_tx_count_24h" example:"7566"`
	SellTxCount24H       int    `json:"Sell_tx_count_24h" example:"7566"`
	TokenVolume24H       string `json:"token_volume_24h" example:"15450301.872331"`
	TokenVolume24HUsd    string `json:"token_volume_24h_usd" example:"15450301.872331"`
	BuyTokenVolume24H    string `json:"buy_token_volume_24h" example:"15435782.918031"`
	BuyTokenVolume24Usd  string `json:"buy_token_volume_24h_usd" example:"15435782.918031"`
	SellTokenVolume24H   string `json:"sell_token_volume_24h" example:"15435782.918031"`
	SellTokenVolume24Usd string `json:"sell_token_volume_24h_usd" example:"15435782.918031"`

	PriceChange1H       string `json:"price_change_1h" example:"0.005103"`
	TxCount1H           int    `json:"tx_count_1h" example:"271"`
	BuyTxCount1H        int    `json:"buy_tx_count_1h" example:"271"`
	SellTxCount1H       int    `json:"sell_tx_count_1h" example:"271"`
	TokenVolume1H       string `json:"token_volume_1h" example:"379231.972430"`
	TokenVolume1HUsd    string `json:"token_volume_1h_usd" example:"379231.972430"`
	BuyTokenVolume1H    string `json:"buy_token_volume_1h" example:"379231.972430"`
	BuyTokenVolume1Usd  string `json:"buy_token_volume_1h_usd" example:"379231.972430"`
	SellTokenVolume1H   string `json:"sell_token_volume_1h" example:"379231.972430"`
	SellTokenVolume1Usd string `json:"sell_token_volume_1h_usd" example:"379231.972430"`

	PriceChange30M       string `json:"price_change_30m" example:"0.000836"`
	TxCount30M           int    `json:"tx_count_30m" example:"25"`
	BuyTxCount30M        int    `json:"buy_tx_count_30m" example:"25"`
	SellTxCount30M       int    `json:"sell_tx_count_30m" example:"25"`
	TokenVolume30M       string `json:"token_volume_30m" example:"62269.324250"`
	TokenVolume30MUsd    string `json:"token_volume_30m_usd" example:"62269.324250"`
	BuyTokenVolume30M    string `json:"buy_token_volume_30m" example:"62269.324250"`
	BuyTokenVolume30Usd  string `json:"buy_token_volume_30m_usd" example:"62269.324250"`
	SellTokenVolume30M   string `json:"sell_token_volume_30m" example:"62269.324250"`
	SellTokenVolume30Usd string `json:"sell_token_volume_30m_usd" example:"62269.324250"`

	PriceChange5M       string `json:"price_change_5m" example:"0.000836"`
	TxCount5M           int    `json:"tx_count_5m" example:"25"`
	BuyTxCount5M        int    `json:"buy_tx_count_5m" example:"25"`
	SellTxCount5M       int    `json:"sell_tx_count_5m" example:"25"`
	TokenVolume5M       string `json:"token_volume_5m" example:"62269.324250"`
	TokenVolume5MUsd    string `json:"token_volume_5m_usd" example:"62269.324250"`
	BuyTokenVolume5M    string `json:"buy_token_volume_5m" example:"62269.324250"`
	BuyTokenVolume5Usd  string `json:"buy_token_volume_5m_usd" example:"62269.324250"`
	SellTokenVolume5M   string `json:"sell_token_volume_5m" example:"62269.324250"`
	SellTokenVolume5Usd string `json:"sell_token_volume_5m_usd" example:"62269.324250"`

	LastSwapAt int64  `json:"last_swap_at" example:"1740887399"`
	MarketCap  string `json:"market_cap" example:"50417.320729000000000"`
	Holders    int    `json:"holders" example:"4204"`

	// High1H             string `json:"high_1h" example:"0.000050417320729"`            // 1小时最高价
	// Low1H              string `json:"low_1h" example:"0.000050161350930"`             // 1小时最低价
	// NativeVolume1H     string `json:"native_volume_1h" example:"19.071274881"`        // 1小时原生代币交易量
	// BuyNativeVolume1H  string `json:"buy_native_volume_1h" example:"19.071274881"`    // 1小时买入原生代币交易量
	// High5M             string `json:"high_5m" example:"0.000050417320729"`            // 5分钟最高价
	// Low5M              string `json:"low_5m" example:"0.000050375183502"`             // 5分钟最低价
	// NativeVolume5M     string `json:"native_volume_5m" example:"3.138140370"`         // 5分钟原生代币交易量
	// BuyNativeVolume5M  string `json:"buy_native_volume_5m" example:"3.138140370"`     // 5分钟买入原生代币交易量
	// High24H            string `json:"high_24h" example:"0.000050417320729"`           // 24小时最高价
	// Low24H             string `json:"low_24h" example:"0.000041158184158"`            // 24小时最低价
	// NativeVolume24H    string `json:"native_volume_24h" example:"704.558749824"`      // 24小时原生代币交易量
	// BuyNativeVolume24H string `json:"buy_native_volume_24h" example:"703.884254256"`  // 24小时买入原生代币交易量

}
