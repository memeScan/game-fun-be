package api

import (
	"fmt"
	"game-fun-be/internal/service"

	"strconv"

	"github.com/gin-gonic/gin"
)

type PointsHandler struct {
	pointsService *service.PointsServiceImpl
}

func NewPointsHandler(pointsService *service.PointsServiceImpl) *PointsHandler {
	return &PointsHandler{pointsService: pointsService}
}

// Points 获取用户积分数据
// @Summary 获取用户积分数据
// @Description 根据链类型和用户 ID 获取用户的交易积分、邀请积分和可用积分。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
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
// @Summary 获取用户积分明细
// @Description 根据链类型和用户 ID 获取用户的积分明细数据。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param page query string true "分页页码"
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

	limit := c.GetInt("limit")
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

// PointsDetail 获取用户邀请明细
// @Summary 获取用户邀请明细
// @Description 根据链类型和用户 ID 获取用户的积分明细数据。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param page query string true "分页页码"
// @Param limit query string true "每页数量"
// @Success 200 {object} response.Response{data=response.PointsDetailsResponse} "成功返回用户积分明细"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /points/{chain_type}/detail [get]
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

	limit := c.GetInt("limit")
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
// @Summary 获取用户预估积分数据
// @Description 根据链类型和用户 ID 获取用户的预估积分数据。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 积分
// @Accept json
// @Produce json
// @Security ApiKeyAuth
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
	res := p.pointsService.PointsEstimated(userID, chainType)
	c.JSON(res.Code, res)
}
