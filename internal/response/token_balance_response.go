package response

// TokenBalance represents the structure of the token balance response
type TokenBalance struct {
	Token    string `json:"token"`
	Owner    string `json:"owner"`
	Balance  string `json:"balance"`
	Decimals int    `json:"decimals"`
}
