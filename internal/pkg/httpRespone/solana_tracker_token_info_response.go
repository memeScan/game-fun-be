package httpRespone

type TokenInfoResponse struct {
	Code    int    `json:"code"`
	Data    Token  `json:"data"`
	Message string `json:"message"`
}

type SolanaTrackerToken struct {
	TokenInfo TokenDetail `json:"token"`
	Pools     []Pool      `json:"pools"`
	Events    Events      `json:"events"`
	Risk      Risk        `json:"risk"`
	Buys      int         `json:"buys"`
	Sells     int         `json:"sells"`
	Txns      int         `json:"txns"`
	Holders   int         `json:"holders"`
}

type TokenDetail struct {
	Name            string                  `json:"name"`
	Symbol          string                  `json:"symbol"`
	Mint            string                  `json:"mint"`
	URI             string                  `json:"uri"`
	Decimals        int                     `json:"decimals"`
	Image           string                  `json:"image"`
	Description     string                  `json:"description"`
	Extensions      SolanaTrackerExtensions `json:"extensions"`
	Tags            []string                `json:"tags"`
	Creator         Creator                 `json:"creator"`
	HasFileMetaData bool                    `json:"hasFileMetaData"`
}

type SolanaTrackerExtensions struct {
	Twitter  string `json:"twitter"`
	Telegram string `json:"telegram"`
}

type Creator struct {
	Name string `json:"name"`
	Site string `json:"site"`
}

type Pool struct {
	Liquidity    Amounts  `json:"liquidity"`
	Price        Amounts  `json:"price"`
	TokenSupply  float64  `json:"tokenSupply"`
	LpBurn       int      `json:"lpBurn"`
	TokenAddress string   `json:"tokenAddress"`
	MarketCap    Amounts  `json:"marketCap"`
	Market       string   `json:"market"`
	QuoteToken   string   `json:"quoteToken"`
	Decimals     int      `json:"decimals"`
	Security     Security `json:"security"`
	LastUpdated  int64    `json:"lastUpdated"`
	CreatedAt    int64    `json:"createdAt"`
	PoolId       string   `json:"poolId"`
}

type Amounts struct {
	Quote float64 `json:"quote"`
	USD   float64 `json:"usd"`
}

type Security struct {
	FreezeAuthority string `json:"freezeAuthority"`
	MintAuthority   string `json:"mintAuthority"`
}

type Events struct {
	OneMin         PriceChange `json:"1m"`
	FiveMin        PriceChange `json:"5m"`
	FifteenMin     PriceChange `json:"15m"`
	ThirtyMin      PriceChange `json:"30m"`
	OneHour        PriceChange `json:"1h"`
	TwoHour        PriceChange `json:"2h"`
	ThreeHour      PriceChange `json:"3h"`
	FourHour       PriceChange `json:"4h"`
	FiveHour       PriceChange `json:"5h"`
	SixHour        PriceChange `json:"6h"`
	TwelveHour     PriceChange `json:"12h"`
	TwentyFourHour PriceChange `json:"24h"`
}

type PriceChange struct {
	PriceChangePercentage float64 `json:"priceChangePercentage"`
}

type Risk struct {
	Rugged bool       `json:"rugged"`
	Risks  []RiskItem `json:"risks"`
	Score  int        `json:"score"`
}

type RiskItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Score       int    `json:"score"`
}
