package api

import (
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

// SolUsdPrice 获取 SOL 对 USD 的价格
// @Summary 获取 SOL 对 USD 的当前价格
// @Description 返回 SOL 对 USD 的当前价格，保留 8 位小数
// @Tags 全局数据
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]string} "成功返回 SOL 价格"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /globa/sol_usd_price [get]
func SolUsdPrice(c *gin.Context) {
	globaService := service.NewGlobalServiceImpl()
	res := globaService.SolUsdPrice()
	c.JSON(200, res)
}

// SolBalance 获取 SOL 余额
// @Summary 获取当前用户的 SOL 余额
// @Description 根据 JWT Token 获取当前用户的 SOL 余额
// @Tags 全局数据
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.TokenBalance} "成功返回 SOL 余额"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /globa/sol_balance [get]
func SolBalance(c *gin.Context) {
	address, exists := c.Get("address")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Err(http.StatusUnauthorized, "Address not found in context", nil))
		return
	}
	addressStr, ok := address.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Err(http.StatusUnauthorized, "Invalid address type in context", nil))
		return
	}
	globaService := service.NewGlobalServiceImpl()
	res := globaService.SolBalance(addressStr)
	c.JSON(res.Code, res)
}
