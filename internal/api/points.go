package api

import (
	"game-fun-be/internal/service"

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
// @Description 根据用户 ID 获取用户的交易积分、邀请积分和可用积分
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.PointsResponse} "成功返回用户积分数据"
// @Failure 401 {object} response.Response "未授权"
// @Router /points [get]
func (p *PointsHandler) Points(c *gin.Context) {
	userID, errResp := GetUserIDFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := p.pointsService.Points(userID)
	c.JSON(res.Code, res)
}

// PointsDetail 获取用户积分明细
// @Summary 获取用户积分明细
// @Description 根据用户 ID 获取用户的积分明细数据
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param userID path string true "用户 ID"
// @Param page query string true "分页页码"
// @Param limit query string true "每页数量"
// @Success 200 {object} response.Response{data=response.PointsDetailsResponse} "成功返回用户积分明细"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /points/detail [get]
func (p *PointsHandler) PointsDetail(c *gin.Context) {
	userID, errResp := GetUserIDFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	page, limit, errResp := GetPageAndLimit(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := p.pointsService.PointsDetail(userID, page, limit)
	c.JSON(res.Code, res)
}

// PointsEstimated 获取用户预估积分数据
// @Summary 获取用户预估积分数据
// @Description 根据用户 ID 获取用户的预估积分数据
// @Tags 积分
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.PointsEstimatedResponse} "成功返回用户预估积分数据"
// @Failure 401 {object} response.Response "未授权"
// @Router /points/estimated [get]
func (p *PointsHandler) PointsEstimated(c *gin.Context) {
	userID, errResp := GetUserIDFromContext(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := p.pointsService.PointsEstimated(userID)
	c.JSON(res.Code, res)
}
