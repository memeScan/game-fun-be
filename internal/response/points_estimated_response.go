package response

// PointsEstimatedResponse 定义预估积分响应数据结构
type PointsEstimatedResponse struct {
	EstimatedPoints string `json:"estimated_points" example:"12862.90277"` // 预估积分
}
