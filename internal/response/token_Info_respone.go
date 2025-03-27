package response

import "github.com/shopspring/decimal"

// TokenInfoResponse represents the detailed information of a token.
// @Description Detailed information about a token
type TokenInfoResponse struct {
	Address              string          `json:"address"`                 // The unique address of the token
	Symbol               string          `json:"symbol"`                  // The symbol representing the token
	Name                 string          `json:"name"`                    // The name of the token
	Decimals             int             `json:"decimals"`                // The number of decimal places the token uses
	Creator              string          `json:"creator"`                 // The address of the token's creator
	DevNativeTokenAmount uint64          `json:"dev_native_token_amount"` // The amount of native tokens held by the developer
	Logo                 string          `json:"logo"`                    // URL to the token's logo image
	Website              string          `json:"website"`                 // The website of the token
	Twitter              string          `json:"twitter"`                 // The Twitter handle of the token
	Telegram             string          `json:"telegram"`                // The Telegram handle of the token
	CreateTimestamp      int64           `json:"create_timestamp"`        // The timestamp when the token was created
	BiggestPoolAddress   string          `json:"biggest_pool_address"`    // Address of the largest liquidity pool for the token
	OpenTimestamp        int64           `json:"open_timestamp"`          // The timestamp when the token was first opened for trading
	HolderCount          int             `json:"holder_count"`            // The number of holders of the token
	CirculatingSupply    uint64          `json:"circulating_supply"`      // The amount of tokens currently in circulation
	TotalSupply          uint64          `json:"total_supply"`            // The total supply of the token
	MaxSupply            uint64          `json:"max_supply"`              // The maximum supply of the token
	Top10Holdings        float64         `json:"top10_holdings"`          // The percentage of total supply held by the top 10 holders
	Liquidity            decimal.Decimal `json:"liquidity"`               // The liquidity of the token in the market
	Progress             decimal.Decimal `json:"progress"`                // The progress of the token
	Price                float64         `json:"price"`                   // The price of the token
	PriceChange          float64         `json:"price_change"`            // The price change of the token
	Volume               float64         `json:"volume"`                  // The volume of the tokenv
	IsComplete           bool            `json:"is_complete"`             // Flag indicating if the token is a common token
	// Security check data
	CtoFlag              bool    `json:"cto_flag"`            // Flag indicating if the token is a CTO
	DexscrAd             bool    `json:"dexscr_ad"`           // Advertisement status on DEX
	DexscrUpdateLink     bool    `json:"dexscr_update_link"`  // Update link status on DEX
	TwitterChangeFlag    bool    `json:"twitter_change_flag"` // Twitter change flag
	TotalHolders         bool    `json:"total_holders"`       // The total number of holders of the token
	MintAuthority        bool    `json:"mint_authority"`      // Indicates if the token has minting authority
	FreezeAuthority      bool    `json:"freeze_authority"`    // Indicates if the token has freezing authority
	MarketCap            float64 `json:"market_cap"`
	DexscrUpdateLinkTime int64   `json:"dexscr_update_link_time"` // Timestamp of the last update link on DEX
	CreatorBalanceRate   string  `json:"creator_balance_rate"`    // The percentage of total supply held by the creator
	RatTraderAmountRate  float64 `json:"rat_trader_amount_rate"`  // The percentage of total supply held by rat traders
	CreatorTokenBalance  string  `json:"creator_token_balance"`   // The balance of tokens held by the creator
	Top10HolderRate      float64 `json:"top_10_holder_rate"`      // The percentage of total supply held by the top 10 holders
	DevTokenBurnAmount   float64 `json:"dev_token_burn_amount"`   // The amount of tokens burned by the developer
	DevTokenBurnRatio    float64 `json:"dev_token_burn_ratio"`    // The ratio of tokens burned by the developer
	TokenBurnRatio       float64 `json:"token_burn_ratio"`        // The ratio of tokens burned by the token
	BurnPercentage       float64 `json:"burn_percentage"`         // The percentage of tokens burned by the token
	CreatedPlatformType  uint8   `json:"created_platform_type"`   // The platform type where the token was created

}

func BuildTokenInfoResponse(tokenInfo TokenInfoResponse) Response {
	return Response{
		Code: CodeSuccess,
		Data: tokenInfo,
		Msg:  "Success",
	}
}
