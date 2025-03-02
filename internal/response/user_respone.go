package response

// LoginResult 登录结果数据
// @Description 登录成功后返回的具体数据
type LoginResponse struct {
	Token      string   `json:"token" example:"_vetZFMHLOlPWer8_uTrqe5_7gGGNPi-fvD--pomTdSd2kn3po2uBujC4o7raUiNHPpKJaS_BR3jmgOEdaAEFA=="` // 登录令牌
	User       UserInfo `json:"user"`                                                                                                     // 用户信息
	InviteCode string   `json:"invite_code" example:""`                                                                                   // 邀请码
}

// UserInfo 用户信息
// @Description 用户的基本信息
type UserInfo struct {
	UserID      string `json:"user_id" example:"10014855"`                                     // 用户ID
	Address     string `json:"address" example:"F59CSoJEmjDFQWZMjSjjvu6q7xV31p9rPzRynwrE71yk"` // 钱包地址
	Nickname    string `json:"nickname" example:"EagleI14Jv"`                                  // 昵称
	Avatar      string `json:"avatar" example:""`                                              // 头像
	Description string `json:"description" example:""`                                         // 个人描述
}
