package api

import (
	"game-fun-be/internal/response"
	"game-fun-be/internal/service"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GlobalHandler struct {
	globalService *service.GlobalServiceImpl
}

func NewGlobalHandler(globalService *service.GlobalServiceImpl) *GlobalHandler {
	return &GlobalHandler{globalService: globalService}
}

// SolUsdPrice 获取链的原生代币价格
// @Summary 获取链的原生代币价格
// @Description 根据链类型获取原生代币的 USD 价格，保留 8 位小数
// @Tags 全局数据
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Success 200 {object} response.Response{data=map[string]string} "成功返回原生代币价格"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /global/{chain_type}/native_token_price [get]
func (g *GlobalHandler) NativeTokePrice(c *gin.Context) {
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := g.globalService.NativeTokePrice(chainType)
	c.JSON(res.Code, res)
}

// Balance 获取链的原生代币钱包余额
// security:
//   - Bearer: []
//
// @Summary 获取链的原生代币钱包余额
// @Description 根据链类型和用户地址获取用户的钱包原生代币余额。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 全局数据
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Success 200 {object} response.Response{data=response.TokenBalance} "成功返回原生代币余额"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /global/{chain_type}/native_balance [get]
func (g *GlobalHandler) NativeBalance(c *gin.Context) {
	userAddress, errResp := GetAddressFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := g.globalService.NativeBalance(userAddress, chainType)
	c.JSON(res.Code, res)
}

// Balance 获取链的原生代币钱包余额
// security:
//   - Bearer: []
//
// @Summary 获取链的原生代币钱包余额
// @Description 根据链类型和用户地址获取用户的钱包原生代币余额。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 全局数据
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param ticker_address path string true "代币地址"
// @Success 200 {object} response.Response{data=response.TokenBalance} "成功返回原生代币余额"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /global/{chain_type}/ticker_balance/{ticker_address} [get]
func (g *GlobalHandler) TickerBalance(c *gin.Context) {
	address, errResp := GetAddressFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	tickerAddress := c.Param("ticker_address")
	if tickerAddress == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "ticker_address cannot be empty", errors.New("ticker_address is required")))
		return
	}
	res := g.globalService.TickerBalance(address, tickerAddress, chainType)
	c.JSON(res.Code, res)
}
