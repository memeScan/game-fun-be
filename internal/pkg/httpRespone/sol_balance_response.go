package httpRespone

// SolBalanceResponse represents the structure of the response from the sol-balance API
type SolBalanceBatchResponse struct {
	Code    int                           `json:"code"`
	Data    []SolBalanceBatchResponseData `json:"data"`
	Message string                        `json:"message"`
}

type SolBalanceBatchResponseData struct {
	Address  string `json:"address"`
	Balance  string `json:"balance"`
	Decimals int    `json:"decimals"`
}
