package httpRespone

type PoolInfoResponse2 struct {
	Code    int         `json:"code"`
	Data    []PoolItem2 `json:"data"`
	Message string      `json:"message"`
}

type PoolItem2 struct {
	PoolAddress string    `json:"poolAddress"`
	Data        PoolInfo2 `json:"data"`
}

type PoolInfo2 struct {
	Status                 int64  `json:"status"`
	Nonce                  int64  `json:"nonce"`
	MaxOrder               int64  `json:"maxOrder"`
	Depth                  int64  `json:"depth"`
	BaseDecimal            int64  `json:"baseDecimal"`
	QuoteDecimal           int64  `json:"quoteDecimal"`
	State                  int64  `json:"state"`
	ResetFlag              int64  `json:"resetFlag"`
	MinSize                int64  `json:"minSize"`
	VolMaxCutRatio         int64  `json:"volMaxCutRatio"`
	AmountWaveRatio        int64  `json:"amountWaveRatio"`
	BaseLotSize            int64  `json:"baseLotSize"`
	QuoteLotSize           int64  `json:"quoteLotSize"`
	MinPriceMultiplier     int64  `json:"minPriceMultiplier"`
	MaxPriceMultiplier     int64  `json:"maxPriceMultiplier"`
	SystemDecimalValue     int64  `json:"systemDecimalValue"`
	MinSeparateNumerator   int64  `json:"minSeparateNumerator"`
	MinSeparateDenominator int64  `json:"minSeparateDenominator"`
	TradeFeeNumerator      int64  `json:"tradeFeeNumerator"`
	TradeFeeDenominator    int64  `json:"tradeFeeDenominator"`
	PnlNumerator           int64  `json:"pnlNumerator"`
	PnlDenominator         int64  `json:"pnlDenominator"`
	SwapFeeNumerator       int64  `json:"swapFeeNumerator"`
	SwapFeeDenominator     int64  `json:"swapFeeDenominator"`
	BaseNeedTakePnl        int64  `json:"baseNeedTakePnl"`
	QuoteNeedTakePnl       int64  `json:"quoteNeedTakePnl"`
	QuoteTotalPnl          int64  `json:"quoteTotalPnl"`
	BaseTotalPnl           int64  `json:"baseTotalPnl"`
	PoolOpenTime           int64  `json:"poolOpenTime"`
	PunishPcAmount         int64  `json:"punishPcAmount"`
	PunishCoinAmount       int64  `json:"punishCoinAmount"`
	OrderbookToInitTime    int64  `json:"orderbookToInitTime"`
	SwapBaseInAmount       int64  `json:"swapBaseInAmount"`
	SwapQuoteOutAmount     int64  `json:"swapQuoteOutAmount"`
	SwapBase2QuoteFee      int64  `json:"swapBase2QuoteFee"`
	SwapQuoteInAmount      int64  `json:"swapQuoteInAmount"`
	SwapBaseOutAmount      int64  `json:"swapBaseOutAmount"`
	SwapQuote2BaseFee      int64  `json:"swapQuote2BaseFee"`
	BaseVault              string `json:"baseVault"`
	QuoteVault             string `json:"quoteVault"`
	BaseMint               string `json:"baseMint"`
	QuoteMint              string `json:"quoteMint"`
	LpMint                 string `json:"lpMint"`
	OpenOrders             string `json:"openOrders"`
	MarketId               string `json:"marketId"`
	MarketProgramId        string `json:"marketProgramId"`
	TargetOrders           string `json:"targetOrders"`
	WithdrawQueue          string `json:"withdrawQueue"`
	LpVault                string `json:"lpVault"`
	Owner                  string `json:"owner"`
	LpReserve              int64  `json:"lpReserve"`
	Padding                []int  `json:"padding"`
	BaseReserve            int64  `json:"baseReserve"`
	QuoteReserve           int64  `json:"quoteReserve"`
}
