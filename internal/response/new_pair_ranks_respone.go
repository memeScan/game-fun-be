package response

// NewPairRanksResponse 表示来自 API 的响应结构。
type NewPairRanksResponse struct {
	ID                  int            `json:"id"`                     // ID
	Address             string         `json:"address"`                // 代币地址
	BaseAddress         string         `json:"base_address"`          // 基础代币地址
	QuoteAddress        string         `json:"quote_address"`         // 报价代币地址
	QuoteReserve        string         `json:"quote_reserve"`         // 报价储备
	InitialLiquidity    string         `json:"initial_liquidity"`     // 初始流动性
	InitialQuoteReserve  string        `json:"initial_quote_reserve"`  // 初始报价储备
	Creator             string         `json:"creator"`                // 创建者地址
	PoolTypeStr         string         `json:"pool_type_str"`         // 池类型字符串（如 Raydium、Pump）
	PoolType            int            `json:"pool_type"`             // 池类型（整数表示）
	QuoteSymbol         string         `json:"quote_symbol"`          // 报价代币符号（如 SOL）
	BaseTokenInfo       BaseTokenInfo   `json:"base_token_info"`       // 基础代币信息
	OpenTimestamp       int64          `json:"open_timestamp"`        // 开放时间戳
	Launchpad           string         `json:"launchpad"`             // 启动平台（如 Pump.fun）
}

// BaseTokenInfo 表示基础代币的信息。
type BaseTokenInfo struct {
	Symbol                  string            `json:"symbol"`                     // 代币符号
	Name                    string            `json:"name"`                       // 代币名称
	Logo                    string            `json:"logo"`                       // 代币徽标 URL
	TotalSupply             int64             `json:"total_supply"`               // 总供应量
	Price                   string            `json:"price"`                      // 当前价格
	HolderCount             int               `json:"holder_count"`               // 持有者数量
	LaunchpadStatus         int               `json:"launchpad_status"`           // 启动平台状态（如是否已启动）
	PriceChangePercent1m    string            `json:"price_change_percent1m"`     // 1个月价格变化百分比
	PriceChangePercent5m    string            `json:"price_change_percent5m"`     // 5分钟价格变化百分比
	PriceChangePercent1h    string            `json:"price_change_percent1h"`     // 1小时价格变化百分比
	BurnRatio               string            `json:"burn_ratio"`                 // 销毁比例
	BurnStatus              string            `json:"burn_status"`                // 销毁状态（如已销毁）
	IsShowAlert             bool              `json:"is_show_alert"`              // 是否显示警报
	HotLevel                int               `json:"hot_level"`                  // 热度等级
	Liquidity               string            `json:"liquidity"`                  // 流动性
	Top10HolderRate         float64           `json:"top_10_holder_rate"`         // 前10持有者比例
	RenouncedMint           int               `json:"renounced_mint"`             // 是否放弃铸造（1表示是）
	RenouncedFreezeAccount   int              `json:"renounced_freeze_account"`   // 是否放弃冻结账户（1表示是）
	SocialLinks             SocialLinks       `json:"social_links"`               // 社交链接
	DevTokenBurnAmount      string            `json:"dev_token_burn_amount"`      // 开发者代币销毁数量
	DevTokenBurnRatio       float64           `json:"dev_token_burn_ratio"`       // 开发者代币销毁比例
	DexscrUpdateLink        int               `json:"dexscr_update_link"`         // DEX 更新链接
	CtoFlag                 int               `json:"cto_flag"`                   // CTO 标志
	MarketCap               string            `json:"market_cap"`                 // 市值
	CreatorBalanceRate      float64           `json:"creator_balance_rate"`       // 创建者余额比例
	CreatorTokenStatus      string            `json:"creator_token_status"`       // 创建者代币状态
	RatTraderAmountRate     float64           `json:"rat_trader_amount_rate"`     // RAT 交易者数量比例
	BluechipOwnerPercentage  float64          `json:"bluechip_owner_percentage"`   // 蓝筹股持有者比例
	SmartDegenCount         int               `json:"smart_degen_count"`          // 智能 Degen 数量
	RenownedCount           int               `json:"renowned_count"`             // 知名度数量
	Volume                  string            `json:"volume"`                     // 交易量
	Swaps                   int               `json:"swaps"`                      // 交换次数
	Buys                    int               `json:"buys"`                       // 买入次数
	Sells                   int               `json:"sells"`                      // 卖出次数
	BuyTax                  *float64          `json:"buy_tax,omitempty"`          // 买入税（可选）
	SellTax                 *float64          `json:"sell_tax,omitempty"`         // 卖出税（可选）
	IsHoneypot              *bool             `json:"is_honeypot,omitempty"`      // 是否为蜜罐（可选）
	Renounced               *bool             `json:"renounced,omitempty"`        // 是否放弃（可选）
	DexscrAd                int               `json:"dexscr_ad"`                  // DEX 广告
	TwitterChangeFlag       int               `json:"twitter_change_flag"`        // 推特变化标志
	Address                 string            `json:"address"`                    // 代币地址
}

// SocialLinks 表示代币的社交媒体链接。
type SocialLinks struct {
	TwitterUsername string `json:"twitter_username"` // 推特用户名
	Website         *string `json:"website,omitempty"` // 网站（可选）
	Telegram        *string `json:"telegram,omitempty"` // 电报（可选）
}