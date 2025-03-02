package api

import (
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
// @Summary 用户钱包登录
// @Description 通过钱包地址和签名进行登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param login body request.LoginRequest true "登录请求参数"
// @Success 200 {object} response.LoginResponse "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/login [post]
func Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Err(http.StatusBadRequest, "Invalid request parameters", err))
		return
	}
	userService := service.NewUserService()
	response := userService.Login(req)
	c.JSON(200, response)
}

// MyInfo 获取用户信息
// @Summary 获取当前用户信息
// @Description 根据 JWT Token 获取当前用户的详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.MyInfoResponse} "成功返回用户信息"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/my_info [get]
func MyInfo(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, response.Err(http.StatusUnauthorized, "Authorization header is required", nil))
		return
	}
	userService := service.NewUserService()
	response := userService.MyInfo("userID")
	c.JSON(200, response)
}
