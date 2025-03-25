package httpRespone

type PNLResponse struct {
	Tokens   map[string]TokenPNL `json:"tokens"`
	Summary  Summary             `json:"summary"`
	PNLSince int64               `json:"pnl_since"`
}

type TokenPNL struct {
	Holding           float64 `json:"holding"`
	Held              float64 `json:"held"`
	Sold              float64 `json:"sold"`
	SoldUSD           float64 `json:"sold_usd"`
	Realized          float64 `json:"realized"`
	Unrealized        float64 `json:"unrealized"`
	Total             float64 `json:"total"`
	TotalSold         float64 `json:"total_sold"`
	TotalInvested     float64 `json:"total_invested"`
	AverageBuyAmount  float64 `json:"average_buy_amount"`
	CurrentValue      float64 `json:"current_value"`
	CostBasis         float64 `json:"cost_basis"`
	FirstBuyTime      int64   `json:"first_buy_time"`
	LastBuyTime       int64   `json:"last_buy_time"`
	LastSellTime      int64   `json:"last_sell_time"`
	LastTradeTime     int64   `json:"last_trade_time"`
	BuyTransactions   int     `json:"buy_transactions"`
	SellTransactions  int     `json:"sell_transactions"`
	TotalTransactions int     `json:"total_transactions"`
}

type Summary struct {
	Realized          float64 `json:"realized"`
	Unrealized        float64 `json:"unrealized"`
	Total             float64 `json:"total"`
	TotalInvested     float64 `json:"totalInvested"`
	TotalWins         int     `json:"totalWins"`
	TotalLosses       int     `json:"totalLosses"`
	AverageBuyAmount  float64 `json:"averageBuyAmount"`
	WinPercentage     float64 `json:"winPercentage"`
	LossPercentage    float64 `json:"lossPercentage"`
	NeutralPercentage float64 `json:"neutralPercentage"`
}
