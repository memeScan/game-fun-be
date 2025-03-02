package response

// TokenCheckPool represents the structure of the token check pool response
type TokenCheckPool struct {
	TokenAddress       string  `json:"token_address"`
	PoolAddress        string  `json:"pool_address"`
	DevStatus          uint8   `json:"dev_status"`           // Indicates if the token is a CTO
	CtoFlag            bool    `json:"cto_flag"`             // Flag indicating if the token is a CTO
	DexscrAd           bool    `json:"dexscr_ad"`            // Advertisement status on DEX
	DexscrUpdateLink   bool    `json:"dexscr_update_link"`   // Update link status on DEX
	TwitterChangeFlag  bool    `json:"twitter_change_flag"`  // Twitter change flag
	MintAuthority      bool    `json:"mint_authority"`       // Indicates if the token has minting authority
	FreezeAuthority    bool    `json:"freeze_authority"`     // Indicates if the token has freezing authority
	Top10Holders       bool    `json:"top_10_holders"`       // Indicates if the token has top 10 holders
	Holders            int     `json:"holders"`              // The total number of holders of the token
	Top10HolderRate    float64 `json:"top_10_holder_rate"`   // The percentage of total supply held by the top 10 holders
	IsBurnedLp         bool    `json:"is_burned_lp"`         // Indicates if the token has burned LP
	LpBurnedPercentage float64 `json:"lp_burned_percentage"` // The percentage of LP tokens burned

}
