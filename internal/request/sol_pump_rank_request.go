package request

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// SolRankRequest represents the request parameters for GetSolPumpRank API
type SolRankRequest struct {
	Time             string   `form:"time" binding:"required,oneof=1m 5m 1h 6h 24h"`
	Limit            int      `form:"limit" binding:"required,min=1"`
	From             int      `form:"from,omitempty"` // 新增字段
	OrderBy          string   `form:"orderby" binding:"omitempty,oneof=progress created_timestamp price_change_percent5m change creator_balance holder_count swaps volume reply_count usd_market_cap last_trade_timestamp koth_duration time_since_koth market_cap_1m market_cap_5m price price_change price_change_percent1m price_change_percent1h volume swaps swaps_1h volume_1h liquidity"`
	Direction        string   `form:"direction" binding:"omitempty,oneof=desc asc"`
	NewCreation      *bool    `form:"new_creation" binding:"omitempty"`
	Completing       *bool    `form:"completing" binding:"omitempty"`
	Completed        *bool    `form:"completed" binding:"omitempty"`
	Soaring          *bool    `form:"soaring" binding:"omitempty"`
	NewPool          bool     `form:"new_pool"`
	Burnt            bool     `form:"burnt"`
	DexScreenerSpent bool     `form:"dexscreener_spent"`
	Filters          []string `form:"filters[]" binding:"omitempty"`
	Platforms        []string `form:"platforms[]" binding:"omitempty"`
	MinCreated       *string  `form:"min_created" binding:"omitempty"`
	MaxCreated       *string  `form:"max_created" binding:"omitempty"`
	MinHolderCount   *int     `form:"min_holder_count" binding:"omitempty"`
	MaxHolderCount   *int     `form:"max_holder_count" binding:"omitempty"`
	MinSwaps         *int64   `form:"min_swaps" binding:"omitempty"`
	MaxSwaps         *int64   `form:"max_swaps" binding:"omitempty"`
	MinSwaps1h       *int64   `form:"min_swaps1h,omitempty"`
	MaxSwaps1h       *int64   `form:"max_swaps1h,omitempty"`
	MinMarketcap     float64  `form:"min_marketcap" binding:"omitempty"`
	MaxMarketcap     float64  `form:"max_marketcap" binding:"omitempty"`
	MinVolume        float64  `form:"min_volume" binding:"omitempty"`
	MaxVolume        float64  `form:"max_volume" binding:"omitempty"`
	MinReply         int      `form:"min_reply" binding:"omitempty"`
	KothDuration     string   `form:"koth_duration" binding:"omitempty"`
	TimeSinceKoth    string   `form:"time_since_koth" binding:"omitempty"`
	MinInitLiquidity float64  `form:"min_init_liquidity" binding:"omitempty"`
	MaxInitLiquidity float64  `form:"max_init_liquidity" binding:"omitempty"`
	MinQuoteUsd      *float64 `form:"min_quote_usd" binding:"omitempty"`
	MaxQuoteUsd      *float64 `form:"max_quote_usd" binding:"omitempty"`
	MinProgress      float64  `form:"min_progress" binding:"omitempty"`
}

// Validate performs additional validation on the request
func (r *SolRankRequest) Validate() error {

	// Ensure direction is specified if orderby is provided
	if r.Direction == "" {
		r.Direction = "desc"
	}

	// Validate created time range
	if r.MinCreated != nil && r.MaxCreated != nil {
		minMinutes, errMin := ConvertRelativeTimeToMinutes(*r.MinCreated)
		maxMinutes, errMax := ConvertRelativeTimeToMinutes(*r.MaxCreated)

		if errMin != nil || errMax != nil {
			return errors.New("min_created and max_created must be in valid format")
		}

		if minMinutes > maxMinutes {
			return errors.New("min_created must be less than or equal to max_created")
		}
	}

	// Validate holder count range
	if r.MinHolderCount != nil && r.MaxHolderCount != nil && *r.MinHolderCount > *r.MaxHolderCount {
		return errors.New("min_holder_count must be less than or equal to max_holder_count")
	}

	// Validate swaps range
	if r.MinSwaps != nil && r.MaxSwaps != nil && *r.MinSwaps > *r.MaxSwaps {
		return errors.New("min_swaps must be less than or equal to max_swaps")
	}

	// Validate volume range
	if r.MinVolume != 0 && r.MaxVolume != 0 && r.MinVolume > r.MaxVolume {
		return errors.New("min_volume must be less than or equal to max_volume")
	}

	// Validate market cap range
	if r.MinMarketcap != 0 && r.MaxMarketcap != 0 && r.MinMarketcap > r.MaxMarketcap {
		return errors.New("min_marketcap must be less than or equal to max_marketcap")
	}

	// Validate initial liquidity range
	if r.MinInitLiquidity != 0 && r.MaxInitLiquidity != 0 && r.MinInitLiquidity > r.MaxInitLiquidity {
		return errors.New("min_init_liquidity must be less than or equal to max_init_liquidity")
	}

	// Validate quote USD range
	if r.MinQuoteUsd != nil && r.MaxQuoteUsd != nil && *r.MinQuoteUsd > *r.MaxQuoteUsd {
		return errors.New("min_quote_usd must be less than or equal to max_quote_usd")
	}

	return nil
}

// ConvertRelativeTimeToMinutes converts a relative time string to minutes
func ConvertRelativeTimeToMinutes(relativeTime string) (int, error) {
	re := regexp.MustCompile(`(\d+)\s*([smh])`)
	matches := re.FindStringSubmatch(relativeTime)

	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid relative time format")
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	unit := matches[2]

	var minutes int
	switch unit {
	case "s":
		minutes = value / 60
	case "m":
		minutes = value
	case "h":
		minutes = value * 60
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}

	return minutes, nil
}
