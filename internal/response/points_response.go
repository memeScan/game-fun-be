package response

// PointsResponse 定义积分数据结构
type PointsResponse struct {
	TradingPoints   string `json:"trading_points" example:"0.147938"`   // 交易积分
	InvitePoints    string `json:"invite_points" example:"0"`           // 邀请积分
	AvailablePoints string `json:"available_points" example:"0.147938"` // 可用积分
}
