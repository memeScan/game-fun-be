package api

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSwapRoute(c *gin.Context) {
	var req request.SwapRouteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chainType := c.Param("chainType")
	tradeType := c.Param("tradeType")
	swapService := service.SwapService{}
	swapRoute := swapService.GetPumpSwapRoute(model.FromString(chainType), tradeType, req)
	c.JSON(swapRoute.Code, swapRoute)
}

func SendSwapRequest(c *gin.Context) {
	swapTransaction := c.Query("swap_transaction")
	isJito := c.Query("is_anti_mev")
	isJitoBool := false
	if isJito == "true" {
		isJitoBool = true
	}
	swapService := service.SwapService{}
	swapResponse := swapService.SendSwapRequest(swapTransaction, isJitoBool)
	c.JSON(swapResponse.Code, swapResponse)
}

func GetSwapRequestStatus(c *gin.Context) {
	SwapTransaction := c.Query("swap_transaction")
	swapService := service.SwapService{}
	swapResponse := swapService.GetSwapRequestStatusBySignature(SwapTransaction)
	c.JSON(swapResponse.Code, swapResponse)
}
