package api

import (
	"strconv"

	"game-fun-be/internal/request"
	"game-fun-be/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminTokenConfigHandler struct {
	tokenconfigService *service.TokenConfigServiceImpl
}

func NewAdminTokenConfigHandler(tokenconfigService *service.TokenConfigServiceImpl) *AdminTokenConfigHandler {
	return &AdminTokenConfigHandler{tokenconfigService: tokenconfigService}
}

// GetAdminTokenConfigList 获取tokenconfig列表
// security:
//   - ApiKey: []
//
// @Summary 获取tokenconfig列表
// @Description 获取系统中所有tokenconfig的列表
// @Tags 管理员-令牌配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query string false "页码"
// @Param limit query string false "每页数量"
// @Success 200 {object} response.Response{data=model.TokenConfig} "成功返回tokenconfig列表"
// @Failure 401 {object} response.Response "未授权"
// @Router /admin/tokenconfigs/list [get]
func (a *AdminTokenConfigHandler) GetAdminTokenConfigList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	res := a.tokenconfigService.GetTokenConfigs(page, limit)
	c.JSON(res.Code, res)
}

// GetTokenConfig 获取tokenconfig详情
// security:
//   - ApiKey: []
//
// @Summary 获取tokenconfig详情
// @Description 根据ID获取tokenconfig的详细信息
// @Tags 管理员-令牌配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "TokenConfig ID"
// @Success 200 {object} response.Response{data=model.TokenConfig} "成功返回tokenconfig详情"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "未找到"
// @Router /admin/tokenconfigs/detail/{id} [get]
func (a *AdminTokenConfigHandler) GetTokenConfig(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "Invalid ID format"})
		return
	}

	res := a.tokenconfigService.GetTokenConfig(uint(idInt))
	c.JSON(res.Code, res)
}

// CreateTokenConfig 创建tokenconfig
// security:
//   - ApiKey: []
//
// @Summary 创建tokenconfig
// @Description 创建新的tokenconfig
// @Tags 管理员-令牌配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param tokenconfig body request.CreateTokenConfigRequest true "TokenConfig信息"
// @Success 200 {object} response.Response{data=model.TokenConfig} "成功创建tokenconfig"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /admin/tokenconfigs/create [post]
func (a *AdminTokenConfigHandler) CreateTokenConfig(c *gin.Context) {
	var request request.CreateTokenConfigRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	res := a.tokenconfigService.CreateTokenConfig(request.Name, request.Symbol, request.Address,
		request.EnableMining, request.MiningStartTime, request.MiningEndTime, request.IsListed, request.Description)
	c.JSON(res.Code, res)
}

// UpdateTokenConfig 更新tokenconfig
// security:
//   - ApiKey: []
//
// @Summary 更新tokenconfig
// @Description 根据ID更新tokenconfig的信息
// @Tags 管理员-令牌配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "TokenConfig ID"
// @Param tokenconfig body request.UpdateTokenConfigRequest true "TokenConfig信息"
// @Success 200 {object} response.Response{data=model.TokenConfig} "成功更新tokenconfig"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "未找到"
// @Router /admin/tokenconfigs/update/{id} [post]
func (a *AdminTokenConfigHandler) UpdateTokenConfig(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "Invalid ID format"})
		return
	}

	var request request.UpdateTokenConfigRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	res := a.tokenconfigService.UpdateTokenConfig(uint(idInt), request.Name, request.Symbol, request.Address,
		request.EnableMining, request.MiningStartTime, request.MiningEndTime, request.IsListed, request.Description)
	c.JSON(res.Code, res)
}

// DeleteTokenConfig 删除tokenconfig
// security:
//   - ApiKey: []
//
// @Summary 删除tokenconfig
// @Description 根据ID删除tokenconfig
// @Tags 管理员-令牌配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "TokenConfig ID"
// @Success 200 {object} response.Response "成功删除tokenconfig"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "未找到"
// @Router /admin/tokenconfigs/delete/{id} [get]
func (a *AdminTokenConfigHandler) DeleteTokenConfig(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "Invalid ID format"})
		return
	}

	res := a.tokenconfigService.DeleteTokenConfig(uint(idInt))
	c.JSON(res.Code, res)
}
