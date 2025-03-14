package api

import (
	"game-fun-be/internal/response"
	"game-fun-be/internal/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenHoldingsHandler struct {
	tokenHoldingsService *service.TokenHoldingsServiceImpl
}

func NewTokenHoldingsHandler(tokenHoldingsService *service.TokenHoldingsServiceImpl) *TokenHoldingsHandler {
	return &TokenHoldingsHandler{tokenHoldingsService: tokenHoldingsService}
}

// TokenHoldings 获取代币持仓数据
// security:
//   - Bearer: []
//
// @Summary 获取代币持仓数据
// @Description 根据链类型、用户账户和目标账户获取代币持仓数据。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 代币持仓
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param account path string true "用户账户"
// @Param account query string true "目标账户"
// @Param allow_zero_balance query string false "是否包含零余额" default(false)
// @Success 200 {object} response.Response{data=response.TokenHoldingsResponse} "成功返回代币持仓数据"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /token_holdings/{chain_type}/{account} [get]
func (h *TokenHoldingsHandler) TokenHoldings(c *gin.Context) {
	userAccount := c.Param("account")
	if userAccount == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "User account parameter is required", nil))
		return
	}
	targetAccount := c.Query("account")
	if targetAccount == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "Target account parameter is required", nil))
		return
	}
	allowZeroBalance := c.DefaultQuery("allow_zero_balance", "false")
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := h.tokenHoldingsService.TokenHoldings(userAccount, targetAccount, allowZeroBalance, chainType)
	c.JSON(res.Code, res)
}

// TokenHoldingsHistories 获取代币持仓历史数据
// security:
//   - Bearer: []
//
// @Summary 获取代币持仓历史数据
// @Description 根据链类型和用户账户获取代币持仓历史数据。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 代币持仓
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param account path string true "用户账户"
// @Param page query string false "分页页码" default(0)
// @Param limit query string false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.TokenHoldingHistoriesResponse} "成功返回代币持仓历史数据"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /token_holdings/{chain_type}/histories/{account} [get]
func (h *TokenHoldingsHandler) TokenHoldingsHistories(c *gin.Context) {
	userAccount := c.Param("account")
	if userAccount == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "User account parameter is required", nil))
		return
	}
	page, limit, errResp := GetPageAndLimit(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := h.tokenHoldingsService.TokenHoldingsHistories(userAccount, page, limit, chainType)
	c.JSON(res.Code, res)
}
