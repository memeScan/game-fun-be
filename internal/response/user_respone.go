package response

// LoginResult 登录结果数据
// @Description 登录成功后返回的具体数据
type LoginResponse struct {
	Token      string   `json:"token" example:"_vetZFMHLOlPWer8_uTrqe5_7gGGNPi-fvD--pomTdSd2kn3po2uBujC4o7raUiNHPpKJaS_BR3jmgOEdaAEFA=="` // 登录令牌
	User       UserInfo `json:"user"`                                                                                                     // 用户信息
	InviteCode string   `json:"invite_code" example:""`                                                                                   // 邀请码

}
