package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/request"
	"game-fun-be/internal/service"

	"github.com/gin-gonic/gin"
)

type PointsHandler struct {
	pointsService *service.PointsServiceImpl
	globalService *service.GlobalServiceImpl
}

func NewPointsHandler(pointsService *service.PointsServiceImpl, globalService *service.GlobalServiceImpl) *PointsHandler {
	return &PointsHandler{pointsService: pointsService, globalService: globalService}
}

// Points 获取用户积分数据
// security:
//   - Bearer: []
//
// @Summary 获取用户积分数据
// @Description 根据链类型和用户 ID 获取用户的交易积分、邀请积分和可用积分。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 积分
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Success 200 {object} response.Response{data=response.PointsResponse} "成功返回用户积分数据"
// @Failure 401 {object} response.Response "未授权"
// @Router /points/{chain_type} [get]
func (p *PointsHandler) Points(c *gin.Context) {
	userID, errResp := GetUserIDFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	userIDInt, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}
	res := p.pointsService.Points(uint(userIDInt), chainType)
	c.JSON(res.Code, res)
}

// PointsDetail 获取用户积分明细
// security:
//   - Bearer: []
//
// @Summary 获取用户积分明细
// @Description 根据链类型和用户 ID 获取用户的积分明细数据。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 积分
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param cursor query string false "游标"
// @Param limit query string true "每页数量"
// @Success 200 {object} response.Response{data=response.PointsDetailsResponse} "成功返回用户积分明细"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /points/{chain_type}/detail [get]
func (p *PointsHandler) PointsDetail(c *gin.Context) {
	userID, errResp := GetUserIDFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	cursorStr := c.Query("cursor")
	var cursor *uint
	if cursorStr != "" {
		if cursorVal, err := strconv.ParseUint(cursorStr, 10, 64); err == nil {
			cursorUint := uint(cursorVal)
			cursor = &cursorUint
		}
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	userIDInt, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}
	res := p.pointsService.PointsDetail(userIDInt, cursor, limit, chainType)
	c.JSON(res.Code, res)
}

// InvitedPointsDetail 获取用户邀请明细
// @Summary 获取用户邀请明细
// @Description 根据链类型和用户 ID 获取用户的积分统计数据。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param cursor query string false "游标"
// @Param limit query string true "每页数量"
// @Success 200 {object} response.Response{data=response.InvitedPointsTotalResponse} "成功返回用户积分明细"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /points/{chain_type}/invite/detail [get]
func (p *PointsHandler) InvitedPointsDetail(c *gin.Context) {
	userID, errResp := GetUserIDFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	cursorStr := c.Query("cursor")
	var cursor *uint
	if cursorStr != "" {
		if cursorVal, err := strconv.ParseUint(cursorStr, 10, 64); err == nil {
			cursorUint := uint(cursorVal)
			cursor = &cursorUint
		}
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	userIDInt, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}
	res := p.pointsService.InvitedPointsDetail(userIDInt, cursor, limit, chainType)
	c.JSON(res.Code, res)
}

// PointsEstimated 获取用户预估积分数据
// security:
//   - Bearer: []
//
// @Summary 获取用户预估积分数据
// @Description 根据链类型和用户 ID 获取用户的预估积分数据。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 积分
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Success 200 {object} response.Response{data=response.PointsEstimatedResponse} "成功返回用户预估积分数据"
// @Failure 401 {object} response.Response "未授权"
// @Router /points/{chain_type}/estimated [get]
func (p *PointsHandler) PointsEstimated(c *gin.Context) {
	userID, errResp := GetUserIDFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}

	tokenAddress := os.Getenv("TOKEN_ADDRESS")
	vaultAddress := os.Getenv("VAULT_ADDRESS")

	tokenBalances, _ := httpUtil.GetTokenBalance(vaultAddress, tokenAddress)

	balanceStr := tokenBalances.Data.Balance
	vaultAmount, _ := strconv.ParseUint(balanceStr, 0, 64)

	res := p.pointsService.PointsEstimated(userID, vaultAmount, chainType)
	c.JSON(res.Code, res)
}

// GetSwapRoute 获取 Swap 路由
// security:
//   - Bearer: []
//
// @Summary 获取 Swap 路由
// @Description 根据链类型、交易类型和请求参数获取 Swap 路由。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags Swap
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param address query string true "钱包地址"
// @Param rebate_amount query string true "返利金额"
// @Success 200 {object} response.Response "成功返回 Swap 路由信息"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /swap/{chain_type}/get_transaction [get]
func (p *PointsHandler) GetTransaction(c *gin.Context) {
	var req request.RebateClaimRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := p.pointsService.CheckRebate(req.Address, req.RebateAmount)
	c.JSON(res.Code, res)
}

// SendSwapRequest 发送 Swap 请求
// security:
//   - Bearer: []
//
// @Summary 发送 Swap 请求
// @Description 根据链类型和 Swap 交易数据发送 Swap 请求。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags Swap
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Success 200 {object} response.Response "成功返回交易结果"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /swap/{chain_type}/send_transaction [post]
func (p *PointsHandler) SendTransaction(c *gin.Context) {
	userAddress, errResp := GetAddressFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	// 定义请求体结构
	var req request.RebateClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	res := p.pointsService.SendClaimTransaction(userAddress)
	c.JSON(res.Code, res)
}
