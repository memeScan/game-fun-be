package api

import (
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenHoldingsHandler struct {
	tokenHoldingsService service.TokenHoldingsService
}

func NewTokenHoldingsHandler(tokenHoldingsService service.TokenHoldingsService) *TokenHoldingsHandler {
	return &TokenHoldingsHandler{tokenHoldingsService: tokenHoldingsService}
}

func (h *TokenHoldingsHandler) TokenHoldings(c *gin.Context) {
	userAccount := c.Param("account")
	if userAccount == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "User account parameter is required", nil))
		return
	}
	targetAccount := c.Query("account")
	if targetAccount == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "Target account parameter is required", nil))
		return
	}
	allowZeroBalance := c.DefaultQuery("allow_zero_balance", "false")
	res := h.tokenHoldingsService.TokenHoldings(userAccount, targetAccount, allowZeroBalance)
	c.JSON(res.Code, res)
}

func (h *TokenHoldingsHandler) TokenHoldingsHistories(c *gin.Context) {
	userAccount := c.Param("account")
	if userAccount == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "User account parameter is required", nil))
		return
	}
	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "20")
	res := h.tokenHoldingsService.TokenHoldingsHistories(userAccount, page, limit)
	c.JSON(res.Code, res)
}
