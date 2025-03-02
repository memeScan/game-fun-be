package api

import (
	"my-token-ai-be/internal/cron"
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"

	"strconv"

	"github.com/gin-gonic/gin"
)

// ExecuteJob godoc
// @Summary Execute a job
// @Description Trigger a job execution
// @Tags tool
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=string}
// @Failure 500 {object} response.Response
// @Router /tools/execute_reindex_job [post]
func ExecuteReindexJob(c *gin.Context) {
	// 异步执行重建索引任务
	go func() {
		if err := cron.ExecuteReindexJob(); err != nil {
			util.Log().Error("Reindex job failed: %v", err)
		}
	}()

	// 立即返回响应
	c.JSON(200, response.Response{
		Code: 200,
		Data: "reindex job started",
	})
}

// TokenInfoSync godoc
// @Summary Sync token info data
// @Description Synchronize token information data from source to elasticsearch
// @Tags tool
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=string}
// @Failure 500 {object} response.Response
// @Router /tools/token_info_sync [post]
func TokenInfoSyncJob(c *gin.Context) {
	// 异步 立即返回
	go func() {
		err := cron.SyncTokenInfoJob()
		if err != nil {
			util.Log().Error("TokenInfoSyncJob failed: %v", err)
			c.JSON(500, response.Response{Code: 500, Msg: err.Error()})
		}
	}()

	c.JSON(200, response.Response{Code: 200, Data: "token info sync job started"})
}

// ResetTokenPoolInfo 重置所有代币的流动池相关信息
// @Summary 重置所有代币的流动池相关信息
// @Description 异步重置所有代币的流动性、价格等池子相关信息
// @Tags Tool
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=string}
// @Failure 500 {object} response.Response
// @Router /tools/reset_pool_info [post]
func ResetTokenPoolInfo(c *gin.Context) {
	// 获取 startID 参数
	startID, _ := strconv.ParseInt(c.DefaultQuery("start_id", "0"), 10, 64)

	// 异步执行，立即返回
	go func() {
		toolService := service.NewToolService()
		resp := toolService.ResetPoolInfo(startID)
		if resp.Code != 0 {
			util.Log().Error("ResetPoolInfo failed: %v", resp.Error)
		} else {
			util.Log().Info("ResetPoolInfo completed successfully, processed %v records",
				resp.Data.(map[string]interface{})["total_processed"])
		}
	}()

	c.JSON(200, response.Response{
		Code: 0,
		Msg:  "reset pool info job started",
		Data: "任务已开始执行，请查看日志了解进度",
	})
}
