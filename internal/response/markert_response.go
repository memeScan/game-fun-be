package response

// Market 市场信息
// @Description 市场的基本信息
type Market struct {
	MarketID        int    `json:"market_id" example:"1"`                                               // 市场ID
	Market          string `json:"market" example:"B3o198GX45DZ9HeJZ56wuruyWWCfGVBSn5qC71q5AHZz"`       // 市场地址
	TokenMint       string `json:"token_mint" example:"suprkbfvwpFZXzWaoKjTzGkW1nkvvwK9n2E6g1zyLFo"`    // 代币铸造地址
	TokenVault      string `json:"token_vault" example:"9SQy9T4pzadupDcWdYYgFmUgXLVydfQwgKZEbzeACNAQ"`  // 代币金库地址
	NativeVault     string `json:"native_vault" example:"98J2RACW3wk8u36k1cPmDenzHdPjhcd5qXsm5g471kdh"` // 原生代币金库地址
	TokenName       string `json:"token_name" example:"§uper Exchange"`                                 // 代币名称
	TokenSymbol     string `json:"token_symbol" example:"SUPER"`                                        // 代币符号
	Creator         string `json:"creator" example:"3P3PMv28AM7SvNkGmAdHCusPwahSB9QG9z5o4gTvmNBX"`      // 创建者地址
	URI             string `json:"uri" example:"https://static.super.exchange/m/official/super.json"`   // URI
	Price           string `json:"price" example:"0.000050417320729"`                                   // 当前价格
	Holders         int    `json:"holders" example:"4204"`                                              // 持有者数量
	CreateTimestamp int64  `json:"create_timestamp" example:"1740033455"`                               // 创建时间戳
}

// MarketMetadata 市场元数据
// @Description 市场的元数据信息
type MarketMetadata struct {
	ImageURL    string  `json:"image_url" example:"https://static.super.exchange/i/official/super.png"` // 图片URL
	Description string  `json:"description" example:""`                                                 // 描述
	Twitter     string  `json:"twitter" example:"https://x.com/_superexchange"`                         // Twitter链接
	Website     string  `json:"website" example:"https://super.exchange"`                               // 网站链接
	Telegram    string  `json:"telegram" example:"https://t.me/SuperExchangeCommunity"`                 // Telegram链接
	Banner      string  `json:"banner" example:""`                                                      // 横幅图片
	Rules       *string `json:"rules" example:""`                                                       // 规则
	Sort        int     `json:"sort" example:"0"`                                                       // 排序
}

// MarketTicker 市场行情
// @Description 市场的行情数据
type MarketTicker struct {
	High24H            string `json:"high_24h" example:"0.000050417320729"`           // 24小时最高价
	Low24H             string `json:"low_24h" example:"0.000041158184158"`            // 24小时最低价
	TokenVolume24H     string `json:"token_volume_24h" example:"15450301.872331"`     // 24小时代币交易量
	BuyTokenVolume24H  string `json:"buy_token_volume_24h" example:"15435782.918031"` // 24小时买入代币交易量
	NativeVolume24H    string `json:"native_volume_24h" example:"704.558749824"`      // 24小时原生代币交易量
	BuyNativeVolume24H string `json:"buy_native_volume_24h" example:"703.884254256"`  // 24小时买入原生代币交易量
	PriceChange24H     string `json:"price_change_24h" example:"0.224965"`            // 24小时价格变化
	TxCount24H         int    `json:"tx_count_24h" example:"7590"`                    // 24小时交易次数
	BuyTxCount24H      int    `json:"buy_tx_count_24h" example:"7566"`                // 24小时买入交易次数
	High1H             string `json:"high_1h" example:"0.000050417320729"`            // 1小时最高价
	Low1H              string `json:"low_1h" example:"0.000050161350930"`             // 1小时最低价
	TokenVolume1H      string `json:"token_volume_1h" example:"379231.972430"`        // 1小时代币交易量
	BuyTokenVolume1H   string `json:"buy_token_volume_1h" example:"379231.972430"`    // 1小时买入代币交易量
	NativeVolume1H     string `json:"native_volume_1h" example:"19.071274881"`        // 1小时原生代币交易量
	BuyNativeVolume1H  string `json:"buy_native_volume_1h" example:"19.071274881"`    // 1小时买入原生代币交易量
	PriceChange1H      string `json:"price_change_1h" example:"0.005103"`             // 1小时价格变化
	TxCount1H          int    `json:"tx_count_1h" example:"271"`                      // 1小时交易次数
	BuyTxCount1H       int    `json:"buy_tx_count_1h" example:"271"`                  // 1小时买入交易次数
	High5M             string `json:"high_5m" example:"0.000050417320729"`            // 5分钟最高价
	Low5M              string `json:"low_5m" example:"0.000050375183502"`             // 5分钟最低价
	TokenVolume5M      string `json:"token_volume_5m" example:"62269.324250"`         // 5分钟代币交易量
	BuyTokenVolume5M   string `json:"buy_token_volume_5m" example:"62269.324250"`     // 5分钟买入代币交易量
	NativeVolume5M     string `json:"native_volume_5m" example:"3.138140370"`         // 5分钟原生代币交易量
	BuyNativeVolume5M  string `json:"buy_native_volume_5m" example:"3.138140370"`     // 5分钟买入原生代币交易量
	PriceChange5M      string `json:"price_change_5m" example:"0.000836"`             // 5分钟价格变化
	TxCount5M          int    `json:"tx_count_5m" example:"25"`                       // 5分钟交易次数
	BuyTxCount5M       int    `json:"buy_tx_count_5m" example:"25"`                   // 5分钟买入交易次数
	LastSwapAt         int64  `json:"last_swap_at" example:"1740887399"`              // 最后交易时间戳
	MarketCap          string `json:"market_cap" example:"50417.320729000000000"`     // 市值
	Rank               int    `json:"rank" example:"1"`                               // 排名
}
