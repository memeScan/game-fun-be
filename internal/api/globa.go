package api

import (
	"game-fun-be/internal/service"

	"github.com/gin-gonic/gin"
)

type GlobalHandler struct {
	globalService *service.GlobalServiceImpl
}

func NewGlobalHandler(globalService *service.GlobalServiceImpl) *GlobalHandler {
	return &GlobalHandler{globalService: globalService}
}

// SolUsdPrice 获取 SOL 对 USD 的价格
// @Summary 获取 SOL 对 USD 的当前价格
// @Description 返回 SOL 对 USD 的当前价格，保留 8 位小数
// @Tags 全局数据
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]string} "成功返回 SOL 价格"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /global/sol_usd_price [get]
func (g *GlobalHandler) SolUsdPrice(c *gin.Context) {
	res := g.globalService.SolUsdPrice()
	c.JSON(res.Code, res)
}

// SolBalance 获取 SOL 余额
// @Summary 获取当前用户的 SOL 余额
// @Description 根据 JWT Token 获取当前用户的 SOL 余额
// @Tags 全局数据
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.TokenBalance} "成功返回 SOL 余额"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /global/sol_balance [get]
func (g *GlobalHandler) SolBalance(c *gin.Context) {
	address, errResp := GetAddressFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := g.globalService.SolBalance(address)
	c.JSON(res.Code, res)
}
