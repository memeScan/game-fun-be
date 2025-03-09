package api

import (
	"game-fun-be/internal/request"
	"game-fun-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SwapHandler struct {
	swapService *service.SwapServiceImpl
}

func NewSwapHandler(swapService *service.SwapServiceImpl) *SwapHandler {
	return &SwapHandler{swapService: swapService}
}

// GetSwapRoute 获取 Swap 路由
// @Summary 获取 Swap 路由
// @Description 根据链类型、交易类型和请求参数获取 Swap 路由。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags Swap
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param tradeType path string true "交易类型（如 buy、sell）"
// @Param inAmount query string true "输入金额"
// @Param tokenInAddress query string true "输入代币地址"
// @Param tokenOutAddress query string true "输出代币地址"
// @Param fromAddress query string true "发送地址"
// @Param slippage query float64 true "滑点（百分比）"
// @Param fee query float64 true "手续费（SOL）"
// @Param isAntiMev query bool false "是否启用 Anti-MEV（默认为 false）"
// @Success 200 {object} response.Response "成功返回 Swap 路由信息"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /swap/{chain_type}/get_transaction [get]
func (s *SwapHandler) GetTransaction(c *gin.Context) {
	var req request.SwapRouteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := s.swapService.GetSwapRoute(req, chainType.Uint8())
	c.JSON(res.Code, res)
}

// SendSwapRequest 发送 Swap 请求
// @Summary 发送 Swap 请求
// @Description 根据链类型和 Swap 交易数据发送 Swap 请求。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags Swap
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param swap_transaction query string true "Swap 交易数据（Base64 编码）"
// @Param is_anti_mev query bool false "是否启用 Anti-MEV（默认为 false）"
// @Success 200 {object} response.Response "成功返回交易结果"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /swap/{chain_type}/send_transaction [post]
func (s *SwapHandler) SendTransaction(c *gin.Context) {
	swapTransaction := c.Query("swap_transaction")
	isJito := c.Query("is_anti_mev")
	isJitoBool := false
	if isJito == "true" {
		isJitoBool = true
	}
	res := s.swapService.SendTransaction(swapTransaction, isJitoBool)
	c.JSON(res.Code, res)
}

// GetSwapRequestStatus 获取 Swap 请求状态
// @Summary 获取 Swap 请求状态
// @Description 根据链类型和交易签名获取 Swap 请求状态。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags Swap
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param swap_transaction query string true "交易签名"
// @Success 200 {object} response.Response "成功返回交易状态"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /swap/{chain_type}/transaction_status [get]
func (s *SwapHandler) TransactionStatus(c *gin.Context) {
	swapTransaction := c.Query("swap_transaction")
	res := s.swapService.GetSwapStatusBySignature(swapTransaction)
	c.JSON(res.Code, res)
}
