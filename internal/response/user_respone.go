package response

// LoginResult 登录结果数据
// @Description 登录成功后返回的具体数据
type LoginResponse struct {
	Token      string   `json:"token" example:"_vetZFMHLOlPWer8_uTrqe5_7gGGNPi-fvD--pomTdSd2kn3po2uBujC4o7raUiNHPpKJaS_BR3jmgOEdaAEFA=="` // 登录令牌
	User       UserInfo `json:"user"`                                                                                                     // 用户信息
	InviteCode string   `json:"invite_code" example:""`                                                                                   // 邀请码

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
