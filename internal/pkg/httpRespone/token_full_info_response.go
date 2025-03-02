package httpRespone

type TokenFullInfoResponse struct {
	Code    int         `json:"code"`
	Data    []TokenInfo `json:"data"`
	Message string      `json:"message"`
}

type TokenInfo struct {
	Decimals            int64  `json:"decimals"`
	FreezeAuthority     int64  `json:"freezeAuthority"`
	MintAuthority       int64  `json:"mintAuthority"`
	Supply              string `json:"supply"`
	Timestamp           int64  `json:"timestamp"`
	Block               int64  `json:"block"`
	Mint                string `json:"mint"`
	Creator             string `json:"creator"`
	Name                string `json:"name"`
	Symbol              string `json:"symbol"`
	URI                 string `json:"uri"`
	Signature           string `json:"signature"`
	PlatformType        int64  `json:"platformType"`
	BondingCurveAddress string `json:"bondingCurveAddress"`
}
