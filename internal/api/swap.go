package api

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/service"
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetSwapRoute godoc
// @Summary Get swap route
// @Description Get swap route for a specific token address
// @Tags swap
// @Accept json
// @Produce json
// @Param chainType path string true "Chain type" Enums(sol, eth, btc)
// @Param tradeType path string true "Trade type" Enums(tx)
// @Param token_in_chain query string true "Token in chain" Enums(sol, eth, btc)
// @Param token_out_chain query string true "Token out chain" Enums(sol, eth, btc)
// @Param from_address query string true "From address"
// @Param slippage query float64 true "Slippage"
// @Param token_in_address query string true "Token in address"
// @Param token_out_address query string true "Token out address"
// @Param in_amount query string true "In amount"
// @Param fee query float64 true "Fee"
// @Param is_anti_mev query bool false "Is anti MEV"
// @Param legacy query bool false "Legacy"
// @Param token_total_supply query int64 true "Token total supply"
// @Param initial_real_token_reserves query int64 true "Initial real token reserves"
// @Param initial_virtual_sol_reserves query int64 true "Initial virtual sol reserves"
// @Param initial_virtual_token_reserves query int64 true "Initial virtual token reserves"
// @Success 200 {object} response.Response{data=response.SwapRouteData}
// @Failure 400 {object} response.Response
// @Router /{chainType}/{tradeType}/get_swap_route [get]
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

// SendSwapRequest godoc
// @Summary Send swap request
// @Description Send swap request
// @Tags swap
// @Accept json
// @Produce json
// @Param transaction path string true "transaction"
// @Success 200 {object} response.Response{data=response.SwapTransactionResponse}
// @Failure 400 {object} response.Response
// @Router /transaction/sol/send_swap_transaction [get]
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

// GetSwapRequestStatus godoc
// @Summary Get swap request status``````
// @Description Get swap request status
// @Tags swap
// @Accept json
// @Produce json
// @Param swap_transaction path string true "swap_transaction"
// @Success 200 {object} response.Response{data=response.SwapTransactionResponse}
// @Failure 400 {object} response.Response
// @Router /transaction/sol/get_swap_request_status [get]
func GetSwapRequestStatus(c *gin.Context) {
	SwapTransaction := c.Query("swap_transaction")
	swapService := service.SwapService{}
	swapResponse := swapService.GetSwapRequestStatusBySignature(SwapTransaction)
	c.JSON(swapResponse.Code, swapResponse)
}
