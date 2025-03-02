package api

import (
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"
	"net/http"
	"my-token-ai-be/internal/request"
	"github.com/gin-gonic/gin"
)

// GetAddressFromContext 从 Gin 上下文中获取地址
func GetAddressFromContext(c *gin.Context) string {
	address, _ := c.Get("address")
	return address.(string)
}

// UserMe godoc
// @Summary Get user details
// @Description Get details of the currently authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/me [get]
func UserMe(c *gin.Context) {
	address := GetAddressFromContext(c)
	// 使用 address 获取用户信息
	respone := service.GetUserByAddress(address)
	c.JSON(respone.Code, respone)
}

// GetMessage godoc
// @Summary Get message for wallet authentication
// @Description Get a message that needs to be signed for wallet authentication
// @Tags authentication
// @Accept json
// @Produce json
// @Param address query string true "Solana wallet address"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.Response
// @Router /user/message [get]
func GetMessage(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "Address is required", nil))
		return
	}

	authMessage := service.GetMessage(address)
	c.JSON(authMessage.Code, authMessage)
}

// WalletLogin godoc
// @Summary Wallet login
// @Description Authenticate user with wallet signature
// @Tags authentication
// @Accept json
// @Produce json
// @Param input body service.WalletAuthService true "Wallet authentication input"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/wallet-login [post]
func WalletLogin(c *gin.Context) {
	var req request.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "Invalid request parameters", err))
		return	
	}
	res := service.Login(req.Address, req.Signature)
	c.JSON(res.Code, res)
}
