package httpRespone

// SolBalanceResponse represents the structure of the response from the sol-balance API
type SolBalanceResponse struct {
	Code    int                      `json:"code"`
	Data    []SolBalanceResponseData `json:"data"`
	Message string                   `json:"message"`
}

type SolBalanceResponseData struct {
	Address  string `json:"address"`
	Balance  string `json:"balance"`
	Decimals int    `json:"decimals"`
}
