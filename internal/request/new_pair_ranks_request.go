package request

// NewPairRanksRequest represents the structure of the request to the API.
// @Description NewPairRanksRequest represents the structure of the request to the API.
type NewPairRanksRequest struct {
	Limit            int      `form:"limit"`          // 使用 form 标签
	From             int      `form:"from,omitempty"` // 新增的字段
	Time             string   `form:"time"`
	OrderBy          string   `form:"orderby" binding:"omitempty,oneof=progress created_timestamp price_change_percent5m change creator_balance holder_count swaps volume reply_count usd_market_cap last_trade_timestamp koth_duration time_since_koth market_cap_1m market_cap_5m price price_change price_change_percent1m price_change_percent1h volume swaps swaps_1h volume_1h liquidity"`
	Direction        string   `form:"direction" binding:"omitempty,oneof=desc asc"`
	NewPool          bool     `form:"new_pool"`
	Burnt            bool     `form:"burnt"`
	DexScreenerSpent bool     `form:"dexscreener_spent"`
	Platforms        []string `form:"platforms"`
	Filters          []string `form:"filters"`
	MinQuoteUSD      float64  `form:"min_quote_usd,omitempty"`
	MaxQuoteUSD      float64  `form:"max_quote_usd,omitempty"`
	MinMarketCap     float64  `form:"min_marketcap,omitempty"`
	MaxMarketCap     float64  `form:"max_marketcap,omitempty"`
	MinVolume        float64  `form:"min_volume,omitempty"`
	MaxVolume        float64  `form:"max_volume,omitempty"`
	MinSwaps1h       int64    `form:"min_swaps1h,omitempty"`
	MaxSwaps1h       int64    `form:"max_swaps1h,omitempty"`
	MinHolderCount   int      `form:"min_holder_count,omitempty"`
	MaxHolderCount   int      `form:"max_holder_count,omitempty"`
	MinCreated       string   `form:"min_created,omitempty"`
	MaxCreated       string   `form:"max_created,omitempty"`
}
