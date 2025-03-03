package api

import (
	"my-token-ai-be/internal/model"

	"my-token-ai-be/internal/service"

	"github.com/gin-gonic/gin"
)

func GetTokenOrderBook(c *gin.Context) {
	chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	tokenInfoService := &service.TokenInfoService{}
	orderBook := tokenInfoService.GetTokenOrderBook(tokenAddress, uint8(model.FromString(chainType)))
	c.JSON(orderBook.Code, orderBook)
}
