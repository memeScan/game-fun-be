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
// @Param trade_type path string true "交易类型（buy 或 sell）"
// @Param token_address query string true "代币地址" example:"So11111111111111111111111111111111111111112"
// @Param from_address query string true "发送地址（用户地址）" example:"CN8R1aHNWLAZm99ymTCd3asErjc2fhe5471cRXs7nJ3m"
// @Param token_in_address query string true "输入代币地址" example:"So11111111111111111111111111111111111111112"
// @Param token_out_address query string true "输出代币地址" example:"FfYhzJ7j3rrs4m4i1wKy5Bz5aYW8mKEGq2rxChU3pump"
// @Param token_in_chain query string true "输入代币所在链（sol、eth、bsc）" example:"sol"
// @Param token_out_chain query string true "输出代币所在链（sol、eth、bsc）" example:"sol"
// @Param in_amount query string true "输入金额（单位：最小代币单位，如 lamports、wei）" example:"100000000"
// @Param priorityFee query int false "交易优先费（单位：最小代币单位，如 lamports）" example:200000000
// @Param slippage query int true "滑点（100 * 100 代表 1%）" example:10000
// @Param is_anti_mev query bool false "是否启用 Anti-MEV（默认 false）" example:false
// @Param legacy query bool false "是否使用 Legacy 交易模式（默认 false）" example:false
// @Param swap_type query string false "交易方向（buy 或 sell，可选）" example:"buy"
// @Param points query int false "积分数（用于 g_points 交易类型）" example:200000000
// @Param platform_type query string true "交易平台类型（pump、raydium、game、g_external、g_points）" example:"pump"
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
