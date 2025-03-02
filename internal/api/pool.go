package api

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/service"

	"github.com/gin-gonic/gin"
)

// GetTokenPoolInfo godoc
// @Summary Get token pool information
// @Description Retrieve pool information for a specific token
// @Tags pool
// @Accept json
// @Produce json
// @Param chainType path string true "Type of the token (e.g., sol)"
// @Param tokenAddress path string true "Address of the token"
// @Success 200 {object} response.Response{data=response.MarketInfo}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /token_pool_info_sol/{chainType}/{tokenAddress} [get]
func GetMarketInfo(c *gin.Context) {
	tokenType := c.Param("chainType")
	address := c.Param("tokenAddress")
	marketInfo := service.GetMarketInfo(address, model.FromString(tokenType))
	c.JSON(marketInfo.Code, marketInfo)
}
