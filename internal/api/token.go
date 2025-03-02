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

// GetTokenInfo godoc
// @Summary Get token information
// @Description Retrieve detailed information for a specific token
// @Tags token
// @Accept json
// @Produce json
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Param tokenAddress path string true "Address of the token"
// @Success 200 {object} response.TokenInfoResponse "Detailed token information"
// @Failure 400 {object} response.Response{data=response.TokenInfoResponse}
// @Router /token_info/{chainType}/{tokenAddress} [get]
func GetTokenInfo(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	tokenInfo := tokenInfoService.GetTokenInfo(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(tokenInfo.Code, tokenInfo)
}

// 获取代币基础信息
// GetTokenInfo godoc
// @Summary Get token information
// @Description Retrieve detailed information for a specific token
// @Tags token
// @Accept json
// @Produce json
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Param tokenAddress path string true "Address of the token"
// @Success 200 {object} response.TokenInfoResponse "Detailed token information"
// @Failure 400 {object} response.Response{data=response.TokenInfoResponse}
// @Router /token_base_info/{chainType}/{tokenAddress} [get]
func GetTokenBaseInfo(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	tokenInfo := tokenInfoService.GetTokenBaseInfo(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(tokenInfo.Code, tokenInfo)
}

// TokenPrice godoc
// @Summary Get token prices
// @Description Retrieve prices for a list of tokens
// @Tags token
// @Accept json
// @Produce json
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Param request body request.TokenPriceRequest true "Array of token addresses"
// @Success 200 {object} response.Response{data=[]response.TokenPriceResponse}
// @Failure 400 {object} response.Response
// @Router /token_prices/{chainType} [post]
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

// GetTokenLaunchpadInfo godoc
// @Summary Get token launchpad info
// @Description Retrieve detailed information for a specific token
// @Tags token
// @Accept json
// @Produce json
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Param tokenAddress path string true "Address of the token"
// @Success 200 {object} response.Response{data=response.TokenLaunchpadInfo}
// @Failure 400 {object} response.Response
// @Router /token_launchpad_info/{chainType}/{tokenAddress} [get]
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

// GetTokenMarketAnalytics godoc
// @Summary Get token market analytics
// @Description Retrieve market analytics for a specific token
// @Tags token
// @Accept json
// @Produce json
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Param tokenAddress path string true "Address of the token"
// @Success 200 {object} response.Response{data=response.TokenMarketAnalyticsResponse}
// @Failure 400 {object} response.Response
// @Router /token_market_analytics/{chainType}/{tokenAddress} [get]
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

// GetTokenOrderBook godoc
// @Summary Get token order book
// @Description Retrieve order book information for a specific token
// @Tags token
// @Accept json
// @Produce json
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Param tokenAddress path string true "Address of the token"
// @Success 200 {object} response.Response{data=response.TokenOrderBookItem}
// @Failure 400 {object} response.Response
// @Router /token_order_book/{chainType}/{tokenAddress} [get]
func GetTokenOrderBook(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	orderBook := tokenInfoService.GetTokenOrderBook(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(orderBook.Code, orderBook)
}

// 代币检测接口
// GetTokenCheckInfo godoc
// @Summary Get token check info
// @Description Retrieve check info for a specific token
// @Tags token
// @Accept json
// @Produce json
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Param tokenAddress path string true "Address of the token"
// @Param tokenPool path string true "Pool of the token"
// @Success 200 {object} response.Response{data=response.TokenCheckInfo}
// @Failure 400 {object} response.Response
// @Router /token_check_info/{chainType}/{tokenAddress}/{tokenPool} [get]
func GetTokenCheckInfo(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenPool := c.Param("tokenPool")

	tokenInfoService := &service.TokenInfoService{}
	detection := tokenInfoService.GetTokenCheckInfo(tokenAddress, uint8(model.FromString(chainType)), tokenPool)
	c.JSON(detection.Code, detection)
}

// Search godoc
// @Summary Search for tokens
// @Description Search for tokens
// @Tags search
// @Accept json
// @Produce json
// @Param query path string true "Query"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /tokens/{chainType}/{tokenAddress}/search [get]
func Search(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	response := tokenInfoService.SearchToken(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(response.Code, response)
}

// 代币市场	查询
// GetTokenMarketQuery godoc
// @Summary Get token market query
// @Description Retrieve market query for a specific token
// @Tags token
// @Accept json
// @Produce json
// @Param tokenAddress path string true "Address of the token"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /token_market_analytics_search/{chainType}/{createdPlatformType/{tokenAddress} [get]
func GetTokenMarketQuery(c *gin.Context) {
	tokenAddress := c.Param("tokenAddress")
	chainType := c.Param("chainType")
	tokenInfoService := &service.TokenInfoService{}
	response := tokenInfoService.GetTokenMarketQuery(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(response.Code, response)
}
