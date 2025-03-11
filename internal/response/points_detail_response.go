package response

// PointsDetail 定义积分明细数据结构
type PointsDetail struct {
	Points    string `json:"points" example:"0.000182"`      // 积分
	Amount    string `json:"amount" example:"0.000182"`      // token数量
	Timestamp int64  `json:"timestamp" example:"1740840671"` // 时间戳
	Type      int8   `json:"type" example:"1"`               // 类型
}

// PointsDetailsData 定义积分明细响应数据结构
type PointsDetailsResponse struct {
	Details []PointsDetail `json:"details"`  // 积分明细列表
	HasMore bool           `json:"has_more"` // 是否还有更多
	Cursor  *uint          `json:"cursor"`   // 游标

}
