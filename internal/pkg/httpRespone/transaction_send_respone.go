package httpRespone

type TransactionSendResponse struct {
	Code    int          `json:"code"`
	Data    ResponseData `json:"data"`
	Message string       `json:"message"`
}

type ResponseData struct {
	Signature string `json:"signature"`
}
