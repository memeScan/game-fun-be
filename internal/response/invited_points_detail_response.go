package response

type InvitedPointsDetail struct {
	Invitee       string `json:"invitee"`
	InviteTime    int64  `json:"invite_time"`
	TradingPoints string `json:"trading_points"`
}

type InvitedPointsTotalResponse struct {
	Details []InvitedPointsDetail `json:"details"`
	HasMore bool                  `json:"has_more"`
	Cursor  *uint                 `json:"cursor"`
}
