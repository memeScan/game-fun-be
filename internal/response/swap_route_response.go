package response

type SwapRouteData struct {
	Quote        Quote       `json:"quote"`
	RawTx        RawTx       `json:"raw_tx"`
	AmountInUSD  string `json:"amountInUSD"`
	AmountOutUSD string `json:"amountOutUSD"`
	JitoOrderID  interface{} `json:"jito_order_id"`
}

type Quote struct {
	InputMint            string      `json:"inputMint"`
	InAmount             string      `json:"inAmount"`
	InDecimals           int         `json:"inDecimals"`
	OutDecimals          int         `json:"outDecimals"`
	OutputMint           string      `json:"outputMint"`
	OutAmount            string      `json:"outAmount"`
	OtherAmountThreshold string      `json:"otherAmountThreshold"`
	SwapMode             string      `json:"swapMode"`
	SlippageBps          int64       `json:"slippageBps"`
	PlatformFee          int64       `json:"platformFee"`
	PriceImpactPct       string      `json:"priceImpactPct"`
	RoutePlan            []RoutePlan `json:"routePlan"`
	TimeTaken            float64     `json:"timeTaken"`
}

type RoutePlan struct {
	SwapInfo SwapInfo `json:"swapInfo"`
	Percent  int      `json:"percent"`
}

type SwapInfo struct {
	AmmKey     string `json:"ammKey"`
	Label      string `json:"label"`
	InputMint  string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	InAmount   string `json:"inAmount"`
	OutAmount  string `json:"outAmount"`
	FeeAmount  int64 `json:"feeAmount"`
	FeeMint    string `json:"feeMint"`
}

type RawTx struct {
	SwapTransaction           string `json:"swapTransaction"`
	LastValidBlockHeight      int    `json:"lastValidBlockHeight"`
	PrioritizationFeeLamports int    `json:"prioritizationFeeLamports"`
	RecentBlockhash           string `json:"recentBlockhash"`
}
