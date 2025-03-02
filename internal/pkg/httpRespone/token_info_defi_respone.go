package httpRespone

type TokenInfoDefiResponse struct {
    Code    int      `json:"code"`
    Data    []Token  `json:"data"`
    Message string   `json:"message"`
}

type Token struct {
    Timestamp string `json:"timestamp"`
    Block     string `json:"block"`
    Mint      string `json:"mint"`
    Creator   string `json:"creator"`
    Name      string `json:"name"`
    Symbol    string `json:"symbol"`
    URI       string `json:"uri"`
    Signature string `json:"signature"`
	Decimals    int    `json:"decimals"`
	FreezeAuthority int `json:"freezeAuthority"`
	MintAuthority   int `json:"mintAuthority"`
	Supply          string `json:"supply"`
	IsInitialized   bool   `json:"isInitialized"`
}
