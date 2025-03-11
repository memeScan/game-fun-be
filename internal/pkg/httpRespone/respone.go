package httpRespone

type ApiResponse struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
}

// ApiResponse 结构体
type SendResponse struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

// Data 结构体
type Data struct {
	Signature string `json:"signature"`
}
