package response

// MyInfoResponse 表示用户信息的响应结构。
// @Description 包含用户基本信息、社交统计、邀请码等详细信息。
type MyInfoResponse struct {
	User          UserInfo `json:"user"`           // 用户基本信息
	FollowerCount int      `json:"follower_count"` // 粉丝数量
	FanCount      int      `json:"fan_count"`      // 关注数量
	VoteCount     int      `json:"vote_count"`     // 投票数量
	MentionCount  int      `json:"mention_count"`  // 提及数量
	FollowStatus  string   `json:"follow_status"`  // 关注状态
	InviteCode    string   `json:"invite_code"`    // 邀请码
	HasBound      bool     `json:"has_bound"`      // 是否已绑定
}

// UserInfo 表示用户的基本信息。
// @Description 包含用户 ID、地址、昵称、头像和描述等信息。
type UserInfo struct {
	UserID      string `json:"user_id"`     // 用户 ID
	Address     string `json:"address"`     // 用户地址
	Nickname    string `json:"nickname"`    // 用户昵称
	Avatar      string `json:"avatar"`      // 用户头像
	Description string `json:"description"` // 用户描述
}
