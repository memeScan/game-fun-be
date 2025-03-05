package request

// LoginRequest 用户登录请求参数
// @Description 用户登录时提交的请求参数
type LoginRequest struct {
	// 用户钱包地址
	Address string `json:"address" example:"F59CSoJEmjDFQWZMjSjjvu6q7xV31p9rPzRynwrE71yk"`
	// 评论ID（可选）
	CommentID string `json:"comment_id" example:""`
	// 邀请码（可选）
	InviteCode string `json:"invite_code" example:"INVITE123"`
	// 登录消息
	Message string `json:"message" example:"Sign in to the super.exchange\r\n\r\nTimestamp: 1740885327"`
	// 帖子ID（可选）
	PostID string `json:"post_id" example:""`
	// 签名信息
	Signature string `json:"signature" example:"Pli28P7PHx6Mzh+RwRTcqtCSuNs2qlncAqC4tK9PFN/CLnvZ2Gm8koqvkLPEZNoiVJOX71hCBaQV9NUKVX8KDg=="`
	// 时间戳
	Timestamp string `json:"timestamp" example:"1740885327"`
}
