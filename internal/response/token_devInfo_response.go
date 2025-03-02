package response
type TokenDevInfoResponse struct {
	Address              string  `json:"address"`
	Creator              string  `json:"creator_address"`
	CreatorTokenBalance  string  `json:"creator_token_balance"`
	CreatorTokenStatus   string  `json:"creator_token_status"`
}
