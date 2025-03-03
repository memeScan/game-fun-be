package response

// InviteCodeResponse 定义邀请码和邀请数量数据结构
type InviteCodeResponse struct {
	InviteCode  string `json:"invite_code" example:"GT4L5B"` // 邀请码
	InviteCount int    `json:"invite_count" example:"0"`     // 邀请数量
}
