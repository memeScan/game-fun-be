package api

import (
	"my-token-ai-be/internal/cron"
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTokenInfo(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	tokenInfo := tokenInfoService.GetTokenInfo(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(tokenInfo.Code, tokenInfo)
}

func GetTokenBaseInfo(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	tokenInfo := tokenInfoService.GetTokenBaseInfo(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(tokenInfo.Code, tokenInfo)
}

func GetTokenPrices(c *gin.Context) {
	chainType := c.Param("chainType")
	var req request.TokenPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "Invalid request parameters", err))
		return
	}
	tokenInfoService := &service.TokenInfoService{}
	prices := tokenInfoService.GetTokenPrices(req.Addresses, model.FromString(chainType))
	c.JSON(prices.Code, prices)
}

func GetTokenLaunchpadInfo(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	info := tokenInfoService.GetTokenLaunchpadInfo(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(info.Code, info)
}

func SearchDocumentsJob(c *gin.Context) {
	jobType := c.Query("type")
	if jobType == "hot5m" {
		err := cron.SwapToken5mDataRefreshTaskQuery()
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Err(http.StatusInternalServerError, "Failed to refresh hot tokens", err))
			return
		}
	} else if jobType == "hot1h" {
		err := cron.SwapToken1hDataRefreshTaskQuery()
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Err(http.StatusInternalServerError, "Failed to refresh hot tokens", err))
			return
		}
	} else if jobType == "delete" {
		cron.DeleteDocumentsJob()
	}
	c.JSON(http.StatusOK, nil)
}

func GetTokenMarketAnalytics(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	analytics := tokenInfoService.GetTokenMarketAnalytics(tokenAddress, uint8(model.FromString(chainType)))
	if analytics.Code != http.StatusOK {
		c.JSON(analytics.Code, analytics)
		return
	}
	c.JSON(http.StatusOK, analytics)
}

func GetTokenOrderBook(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	orderBook := tokenInfoService.GetTokenOrderBook(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(orderBook.Code, orderBook)
}

func GetTokenCheckInfo(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenPool := c.Param("tokenPool")

	tokenInfoService := &service.TokenInfoService{}
	detection := tokenInfoService.GetTokenCheckInfo(tokenAddress, uint8(model.FromString(chainType)), tokenPool)
	c.JSON(detection.Code, detection)
}

func Search(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	response := tokenInfoService.SearchToken(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(response.Code, response)
}

func GetTokenMarketQuery(c *gin.Context) {
	tokenAddress := c.Param("tokenAddress")
	chainType := c.Param("chainType")
	tokenInfoService := &service.TokenInfoService{}
	response := tokenInfoService.GetTokenMarketQuery(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(response.Code, response)
}
