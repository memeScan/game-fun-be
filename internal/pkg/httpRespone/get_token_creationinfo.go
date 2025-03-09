package httpRespone

type TokenCreationInfoResponse struct {
	Success bool `json:"success"`
	Data    struct {
		TxHash         string `json:"txHash"`
		Slot           int64  `json:"slot"`
		TokenAddress   string `json:"tokenAddress"`
		Decimals       int    `json:"decimals"`
		Owner          string `json:"owner"`
		BlockUnixTime  int64  `json:"blockUnixTime"`
		BlockHumanTime string `json:"blockHumanTime"`
	} `json:"data"`
}
