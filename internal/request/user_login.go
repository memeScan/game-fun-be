package request

type UserLoginRequest struct {
	Address string `json:"address" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}