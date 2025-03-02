package httpRespone

type BondingCurvesResponse struct {
	Code    int                `json:"code"`
	Data    []BondingCurveInfo `json:"data"`
	Message string             `json:"message"`
}

type BondingCurveInfo struct {
	Address string      `json:"address"`
	Error   interface{} `json:"error"`
	Data    CurveData   `json:"data"`
}

type CurveData struct {
	Discriminator        string `json:"discriminator"`
	VirtualTokenReserves string `json:"virtualTokenReserves"`
	VirtualSolReserves   string `json:"virtualSolReserves"`
	RealTokenReserves    string `json:"realTokenReserves"`
	RealSolReserves      string `json:"realSolReserves"`
	TokenTotalSupply     string `json:"tokenTotalSupply"`
	Complete             bool   `json:"complete"`
}
