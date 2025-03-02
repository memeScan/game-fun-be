package api

import (
	"net/http"

	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"

	"github.com/gin-gonic/gin"
)

func GetSolPumpRank(c *gin.Context) {
	var req request.SolRankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(response.CodeParamErr, "Invalid request parameters", err))
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(response.CodeParamErr, "Invalid request parameters", err))
		return
	}
	result := service.GetSolPumpRank(&req)
	c.JSON(result.Code, result)
}

func GetSolDumpRank(c *gin.Context) {
	var req request.SolRankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(response.CodeParamErr, "Invalid request parameters", err))
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(response.CodeParamErr, "Invalid request parameters", err))
		return
	}
	result := service.GetSolRaydiumRank(&req)
	c.JSON(result.Code, result)
}

func GetSolSwapRank(c *gin.Context) {
	var req request.SolRankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(response.CodeParamErr, "Invalid request parameters", err))
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(response.CodeParamErr, "Invalid request parameters", err))
		return
	}
	result := service.GetSolSwapRank(&req)
	c.JSON(result.Code, result)
}

func GetNewPairRanks(c *gin.Context) {
	var req request.SolRankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(response.CodeParamErr, "Invalid request parameters", err))
		return
	}
	result := service.GetNewPairRanks(&req)
	c.JSON(result.Code, result)
}
