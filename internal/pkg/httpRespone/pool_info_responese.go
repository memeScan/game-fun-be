package httpRespone

type PoolInfoResponse struct {
	Code    int        `json:"code"`
	Data    []MintData `json:"data"`
	Message string     `json:"message"`
}

type MintData struct {
	Mint string   `json:"mint"`
	Data PoolData `json:"data"`
}

type PoolData struct {
	PoolAddress    string   `json:"poolAddress"`
	ReturnPoolData PoolInfo `json:"returnPoolData"`
}

type PoolInfo struct {
	BaseDecimal  int64 `json:"-"`
	QuoteDecimal int64 `json:"-"`
	// PoolOpenTime        string `json:"poolOpenTime"`
	OrderbookToInitTime int64  `json:"-"`
	BaseVault           string `json:"baseVault"`
	QuoteVault          string `json:"quoteVault"`
	BaseMint            string `json:"baseMint"`
	QuoteMint           string `json:"quoteMint"`
	LpMint              string `json:"lpMint"`
	OpenOrders          string `json:"openOrders"`
	MarketId            string `json:"marketId"`
	MarketProgramId     string `json:"marketProgramId"`
	LpVault             string `json:"lpVault"`
	Owner               string `json:"owner"`
	LpReserve           int64  `json:"-"`
	BaseReserve         string `json:"baseReserve"`
	QuoteReserve        string `json:"quoteReserve"`
}
