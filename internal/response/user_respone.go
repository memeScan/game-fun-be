package response

import (
	"time"
)

// LoginResult 登录结果数据
// @Description 登录成功后返回的具体数据
type LoginResponse struct {
	Token      string    `json:"token" example:"_vetZFMHLOlPWer8_uTrqe5_7gGGNPi-fvD--pomTdSd2kn3po2uBujC4o7raUiNHPpKJaS_BR3jmgOEdaAEFA=="` // 登录令牌
	User       UserInfo  `json:"user"`                                                                                                     // 用户信息
	ExpireTime time.Time `json:"expire_time" example:""`                                                                                   // 邀请码
}

// UserInfo 表示用户的基本信息。
// @Description 包含用户 ID、地址、昵称、头像和描述等信息。
type UserInfo struct {
	UserID          uint   `json:"user_id"`
	Address         string `json:"address"`
	TwitterId       string `json:"twitter_id"`
	TwitterUsername string `json:"twitter_username"`
	InvitationCode  string `json:"invite_code" example:""`
}
