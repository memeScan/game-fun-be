package response

type InvitedPointsDetail struct {
	Invitee       string `json:"invitee"`
	InviteTime    int64  `json:"invite_time"`
	TradingPoints string `json:"trading_points"`
	FeeRebate     string `json:"fee_rebate"`
	UpdateTime    int64  `json:"update_time"`
}

type InvitedPointsTotalResponse struct {
	Details []InvitedPointsDetail `json:"details"`
	HasMore bool                  `json:"has_more"`
	Cursor  *uint                 `json:"cursor"`
}
