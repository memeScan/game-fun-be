package httpRespone

type TransactionStatusResponse struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    TransactionStatusData `json:"data"`
}

type TransactionStatusData struct {
	Signature     string `json:"signature"`
	Status        string `json:"status"`
	Confirmations int    `json:"confirmations"`
	Error         string `json:"error"`
}
