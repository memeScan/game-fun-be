package api

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/service"

	"github.com/gin-gonic/gin"
)

func GetMarketInfo(c *gin.Context) {
	tokenType := c.Param("chainType")
	address := c.Param("tokenAddress")
	marketInfo := service.GetMarketInfo(address, model.FromString(tokenType))
	c.JSON(marketInfo.Code, marketInfo)
}
