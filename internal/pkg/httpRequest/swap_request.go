package httpRequest

import "github.com/shopspring/decimal"

type SwapPumpStruct struct {
	FromAddress                 string          `json:"fromAddress"`
	InputMint                   string          `json:"inputMint"`
	OutputMint                  string          `json:"outputMint"`
	InAmount                    decimal.Decimal `json:"inAmount"`
	SlippageBps                 string          `json:"slippageBps"`
	VirtualSolReserves          string          `json:"virtualSolReserves"`
	VirtualTokenReserves        string          `json:"virtualTokenReserves"`
	PriorityFee                 float64         `json:"priorityFee"`
	TokenTotalSupply            string          `json:"tokenTotalSupply"`
	InitialRealTokenReserves    string          `json:"initialRealTokenReserves"`
	InitialVirtualSolReserves   string          `json:"initialVirtualSolReserves"`
	InitialVirtualTokenReserves string          `json:"initialVirtualTokenReserves"`
	Mev                         bool            `json:"mev"`
	Jitotip                     string          `json:"jitotip"`
}

type SwapRaydiumStruct struct {
	PoolId          string `json:"poolId"`
	MarketId        string `json:"marketId"`
	Owner           string `json:"owner"`
	InputMint       string `json:"inputMint"`
	OutputMint      string `json:"outputMint"`
	InAmount        string `json:"inAmount"`
	SlippageBps     string `json:"slippageBps"`
	PoolPcReserve   string `json:"poolPcReserve"`
	PoolCoinReserve string `json:"poolCoinReserve"`
	PriorityFee     string `json:"priorityFee"`
	PoolPcAddress   string `json:"poolPcAddress"`
	PoolCoinAddress string `json:"poolCoinAddress"`
	Mev             bool   `json:"mev"`
	Jitotip         string `json:"jitotip"`
}

type SwapGInstructionStruct struct {
	User            string `json:"user"`
	InputAmount     string `json:"inputAmount"`
	InputMint       string `json:"inputMint"`
	OutputMint      string `json:"outputMint"`
	PoolPcAddress   string `json:"poolPcAddress"`
	PoolCoinAddress string `json:"poolCoinAddress"`
	PoolPcReserve   string `json:"poolPcReserve"`
	PoolCoinReserve string `json:"poolCoinReserve"`
	SlippageBps     int    `json:"slippageBps"`
	GMint           string `json:"gMint"`
	Amm             string `json:"amm"`
	Market          string `json:"market"`
	GAmm            string `json:"gAmm"`
	GMarket         string `json:"gMarket"`
	Mev             bool   `json:"mev"`
	Jitotip         string `json:"jitotip"`
}

type BuyGWithPointsStruct struct {
	User        string `json:"user"`
	Points      string `json:"points"`
	InputMint   string `json:"inputMint"`
	OutputMint  string `json:"outputMint"`
	SlippageBps int    `json:"slippageBps"`
	GMint       string `json:"gMint"`
	Amm         string `json:"amm"`
	Market      string `json:"market"`
	GAmm        string `json:"gAmm"`
	GMarket     string `json:"gMarket"`
}

type ClaimRebateStruct struct {
	User   string `json:"user"`
	Amount string `json:"amount"`
}
