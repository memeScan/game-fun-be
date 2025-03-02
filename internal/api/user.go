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
// @Router /user/login [post]
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
