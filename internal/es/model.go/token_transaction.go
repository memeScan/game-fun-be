package model

// TokenTransaction 代币交易信息
type TokenTransaction struct {
	// 池信息
	Price              float64 `json:"price,omitempty" default:"0"`
	NativePrice        float64 `json:"native_price,omitempty" default:"0"`
	NativePriceUSD     float64 `json:"native_price_usd,omitempty"`
	PoolAddress        string  `json:"pool_address,omitempty" default:""`
	BaseAddress        string  `json:"base_address,omitempty" default:""`
	QuoteAddress       string  `json:"quote_address,omitempty" default:""`
	QuoteReserve       string  `json:"quote_reserve,omitempty" default:"0"`
	RealNativeReserves string  `json:"real_native_reserves,omitempty" default:"0"`
	Creator            string  `json:"creator,omitempty" default:""`
	PoolTypeStr        string  `json:"pool_type_str,omitempty" default:"unknown"`
	QuoteSymbol        string  `json:"quote_symbol,omitempty" default:""`
	Launchpad          string  `json:"launchpad,omitempty" default:""`
	OpenTimestamp      string  `json:"block_time,omitempty" default:""`
	MarketAddress      string  `json:"market_address,omitempty"`
	NativeTokenAddress string  `json:"native_token_address,omitempty"`

	// 代币基本信息
	ChainType             int     `json:"chain_type,omitempty" default:"0"`
	PlatformType          int     `json:"platform_type,omitempty" default:"0"`
	Decimals              int     `json:"decimals,omitempty" default:"0"`
	CreatedPlatformType   int     `json:"created_platform_type,omitempty" default:"0"`
	TokenAddress          string  `json:"token_address,omitempty" default:""`
	ExtInfo               string  `json:"ext_info,omitempty" default:"{}"`
	UserAddress           string  `json:"user_address,omitempty" default:""`
	TokenSupply           string  `json:"token_supply,omitempty" default:"0"`
	TokenCreateTime       string  `json:"token_create_time,omitempty" default:""`
	LatestTransactionTime string  `json:"transaction_time,omitempty" default:""`
	Progress              float64 `json:"progress,omitempty" default:"0"`
	IsComplete            bool    `json:"is_complete,omitempty" default:"false"`

	// 储备信息
	VirtualNativeReserves string  `json:"virtual_native_reserves,omitempty" default:"0"`
	VirtualTokenReserves  string  `json:"virtual_token_reserves,omitempty" default:"0"`
	InitialNativeReserve  float64 `json:"initial_pc_reserve,omitempty" default:"0"`
	InitialTokenReserve   float64 `json:"initial_token_reserve,omitempty" default:"0"`
	NativeReserveRate     float64 `json:"native_reserve_rate,omitempty" default:"0"`

	// 交易统计
	SwapsCount1m  int64   `json:"swaps_1m,omitempty" default:"0"`
	SwapsCount5m  int64   `json:"swaps_5m,omitempty" default:"0"`
	SwapsCount1h  int64   `json:"swaps_1h,omitempty" default:"0"`
	SwapsCount6h  int64   `json:"swaps_6h,omitempty" default:"0"`
	SwapsCount24h int64   `json:"swaps_24h,omitempty" default:"0"`
	Volume1m      float64 `json:"volume_1m,omitempty" default:"0"`
	Volume5m      float64 `json:"volume_5m,omitempty" default:"0"`
	Volume1h      float64 `json:"volume_1h,omitempty" default:"0"`
	Volume6h      float64 `json:"volume_6h,omitempty" default:"0"`
	Volume24h     float64 `json:"volume_24h,omitempty" default:"0"`
	BuyCount1m    int64   `json:"buy_count_1m,omitempty" default:"0"`
	BuyCount5m    int64   `json:"buy_count_5m,omitempty" default:"0"`
	BuyCount1h    int64   `json:"buy_count_1h,omitempty" default:"0"`
	BuyCount6h    int64   `json:"buy_count_6h,omitempty" default:"0"`
	BuyCount24h   int64   `json:"buy_count_24h,omitempty" default:"0"`
	SellCount1m   int64   `json:"sell_count_1m,omitempty" default:"0"`
	SellCount5m   int64   `json:"sell_count_5m,omitempty" default:"0"`
	SellCount1h   int64   `json:"sell_count_1h,omitempty" default:"0"`
	SellCount6h   int64   `json:"sell_count_6h,omitempty" default:"0"`
	SellCount24h  int64   `json:"sell_count_24h,omitempty" default:"0"`

	// 市场变化
	MarketCap             float64 `json:"market_cap,omitempty" default:"0"`
	MarketCap1mAgo        float64 `json:"market_cap_1m_ago,omitempty" default:"0"`
	MarketCap5mAgo        float64 `json:"market_cap_5m_ago,omitempty" default:"0"`
	MarketCap1hAgo        float64 `json:"market_cap_1h_ago,omitempty" default:"0"`
	MarketCap24hAgo       float64 `json:"market_cap_24h_ago,omitempty" default:"0"`
	PriceChangePercent1m  float64 `json:"price_change_percent1m,omitempty" default:"0"`
	PriceChangePercent5m  float64 `json:"price_change_percent5m,omitempty" default:"0"`
	PriceChangePercent1h  float64 `json:"price_change_percent1h,omitempty" default:"0"`
	PriceChangePercent24h float64 `json:"price_change_percent24h,omitempty" default:"0"`

	// 总计数据
	Buys               int64   `json:"buys,omitempty" default:"0"`
	Sells              int64   `json:"sells,omitempty" default:"0"`
	Swaps              int64   `json:"swaps,omitempty" default:"0"`
	Volumes            float64 `json:"volumes,omitempty" default:"0"`
	PriceChangePercent float64 `json:"price_change_percent,omitempty" default:"0"`
	CreatorBalanceRate float64 `json:"creator_balance_rate,omitempty" default:"0"`

	// 代币标志信息
	Holder               int     `json:"holder,omitempty" default:"0"`
	DevStatus            int     `json:"dev_status,omitempty" default:"0"`
	CrownDuration        int     `json:"crown_duration,omitempty" default:"0"`
	RocketDuration       int     `json:"rocket_duration,omitempty" default:"0"`
	IsDexAd              bool    `json:"is_dex_ad,omitempty" default:"false"`
	IsBurnedLp           bool    `json:"is_burned_lp,omitempty" default:"false"`
	DevNativeTokenAmount float64 `json:"dev_native_token_amount,omitempty" default:"0"`
	DevTokenBurnAmount   float64 `json:"dev_token_burn_amount,omitempty" default:"0"`
	DevTokenBurnRatio    float64 `json:"dev_token_burn_ratio,omitempty" default:"0"`
	LpBurnPercentage     float64 `json:"burn_percentage,omitempty" default:"0"`
	Top10HolderRate      float64 `json:"top_10_percentage,omitempty" default:"0"`

	// 状态标志
	TokenFlags        int  `json:"token_flags,omitempty" default:"0"`
	IsMedia           bool `json:"is_media,omitempty" default:"false"`
	DexscrAd          bool `json:"dexscr_ad,omitempty" default:"false"`
	DexscrUpdateLink  bool `json:"dexscr_update_link,omitempty" default:"false"`
	CtoFlag           bool `json:"cto_flag,omitempty" default:"false"`
	CreatorClose      bool `json:"creator_close,omitempty" default:"false"`
	TwitterChangeFlag bool `json:"twitter_change_flag,omitempty" default:"false"`
	TotalHolders      bool `json:"total_holders,omitempty" default:"false"`
	MintAuthority     bool `json:"mint_authority,omitempty" default:"false"`
	FreezeAuthority   bool `json:"freeze_authority,omitempty" default:"false"`
}
