package httpRespone

type SafetyData struct {
	Mint            string `json:"mint"`
	MintAuthority   int    `json:"mintAuthority"`
	FreezeAuthority int    `json:"freezeAuthority"`
}

type SafetyCheckResponse struct {
	Code    int          `json:"code"`
	Data    []SafetyData `json:"data"`
	Message string       `json:"message"`
}
