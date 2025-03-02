package api

import (
	"net/http"

	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"

	"github.com/gin-gonic/gin"
)

// GetSolPumpRank godoc
// @Summary Get SOL pump rank
// @Description Get the pump rank for SOL tokens
// @Tags pump
// @Accept json
// @Produce json
// @Param time query string true "Time range" Enums(1m, 5m, 1h, 6h, 24h)
// @Param limit query int true "Limit of results" minimum(1)
// @Param orderby query string false "Order by field" Enums(progress, created_timestamp, creator_balance, holder_count, swaps_1h, volume_1h, reply_count, usd_market_cap, last_trade_timestamp, koth_duration, time_since_koth, market_cap_1m, market_cap_5m)
// @Param direction query string false "Sort direction" Enums(asc, desc)
// @Param new_creation query bool false "Filter for new creations"
// @Param completing query bool false "Filter for completing projects"
// @Param pump query bool false "Filter for projects about to pump"
// @Param soaring query bool false "Filter for soaring projects"
// @Param filters[] query []string false "Additional filters"
// @Param min_created query string false "Minimum creation time"
// @Param max_created query string false "Maximum creation time"
// @Param min_holder_count query int false "Minimum holder count"
// @Param max_holder_count query int false "Maximum holder count"
// @Param min_swaps query int false "Minimum swaps"
// @Param max_swaps query int false "Maximum swaps"
// @Param min_marketcap query number false "Minimum market cap"
// @Param max_marketcap query number false "Maximum market cap"
// @Param min_volume query number false "Minimum volume"
// @Param max_volume query number false "Maximum volume"
// @Param min_reply query int false "Minimum reply count"
// @Param koth_duration query string false "Time to reach 1/2 crown"
// @Param time_since_koth query string false "Time to reach 2/2 rocket"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rank/sol/pump [get]
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

// GetSolDumpRank godoc
// @Summary Get SOL dump rank
// @Description Get the dump rank for SOL tokens
// @Tags pump
// @Accept json
// @Produce json
// @Param time query string true "Time range" Enums(1m, 5m, 1h, 6h, 24h)
// @Param limit query int true "Limit of results" minimum(1)
// @Param orderby query string false "Order by field" Enums(progress, created_timestamp, price_change_percent5m, price_change_percent1h, liquidity, creator_balance, holder_count, swaps_1h, volume_1h, reply_count, usd_market_cap, last_trade_timestamp, koth_duration, time_since_koth, market_cap_1m, market_cap_5m)
// @Param direction query string false "Sort direction" Enums(asc, desc)
// @Param completed query bool false "Filter for completed projects"
// @Param filters[] query []string false "Additional filters"
// @Param min_created query string false "Minimum creation time"
// @Param max_created query string false "Maximum creation time"
// @Param min_holder_count query int false "Minimum holder count"
// @Param max_holder_count query int false "Maximum holder count"
// @Param min_swaps query int false "Minimum swaps"
// @Param max_swaps query int false "Maximum swaps"
// @Param min_marketcap query number false "Minimum market cap"
// @Param max_marketcap query number false "Maximum market cap"
// @Param min_volume query number false "Minimum volume"
// @Param max_volume query number false "Maximum volume"
// @Param min_reply query int false "Minimum reply count"
// @Param koth_duration query string false "Time to reach 1/2 crown"
// @Param time_since_koth query string false "Time to reach 2/2 rocket"
// @Param min_init_liquidity query number false "Minimum initial liquidity"
// @Param max_init_liquidity query number false "Maximum initial liquidity"
// @Param min_quote_usd query number false "Minimum quote USD"
// @Param max_quote_usd query number false "Maximum quote USD"
// @Param min_liquidity query number false "Minimum liquidity"
// @Param max_liquidity query number false "Maximum liquidity"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rank/sol/new_pairs [get]
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

// GetSolDumpRank godoc
// @Summary Get SOL dump rank
// @Description Get the dump rank for SOL tokens
// @Tags pump
// @Accept json
// @Produce json
// @Param time query string true "Time range" Enums(1m, 5m, 1h, 6h, 24h)
// @Param limit query int true "Limit of results" minimum(1)
// @Param orderby query string false "Order by field" Enums(progress, created_timestamp, price_change_percent5m, price_change_percent1h, liquidity, creator_balance, holder_count, swaps_1h, volume_1h, reply_count, usd_market_cap, last_trade_timestamp, koth_duration, time_since_koth, market_cap_1m, market_cap_5m)
// @Param direction query string false "Sort direction" Enums(asc, desc)
// @Param completed query bool false "Filter for completed projects"
// @Param filters[] query []string false "Additional filters"
// @Param min_created query string false "Minimum creation time"
// @Param max_created query string false "Maximum creation time"
// @Param min_holder_count query int false "Minimum holder count"
// @Param max_holder_count query int false "Maximum holder count"
// @Param min_swaps query int false "Minimum swaps"
// @Param max_swaps query int false "Maximum swaps"
// @Param min_marketcap query number false "Minimum market cap"
// @Param max_marketcap query number false "Maximum market cap"
// @Param min_volume query number false "Minimum volume"
// @Param max_volume query number false "Maximum volume"
// @Param min_reply query int false "Minimum reply count"
// @Param koth_duration query string false "Time to reach 1/2 crown"
// @Param time_since_koth query string false "Time to reach 2/2 rocket"
// @Param min_init_liquidity query number false "Minimum initial liquidity"
// @Param max_init_liquidity query number false "Maximum initial liquidity"
// @Param min_quote_usd query number false "Minimum quote USD"
// @Param max_quote_usd query number false "Maximum quote USD"
// @Param min_liquidity query number false "Minimum liquidity"
// @Param max_liquidity query number false "Maximum liquidity"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rank/sol/swap [get]
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

// GetNewPairRanks godoc
// @Summary Get new pair ranks
// @Description Get the new pair ranks
// @Tags pump
// @Accept json
// @Produce json
// @Param time query string true "Time range" Enums(1m, 5m, 1h, 6h, 24h)
// @Param limit query int true "Limit of results" minimum(1)
// @Param orderby query string false "Order by field" Enums(progress, created_timestamp, price_change_percent5m, change, creator_balance, holder_count, swaps, volume, reply_count, usd_market_cap, last_trade_timestamp, koth_duration, time_since_koth, market_cap_1m, market_cap_5m, price, price_change, price_change_percent1m, price_change_percent1h, volume, swaps, swaps_1h, volume_1h, liquidity)
// @Param direction query string false "Sort direction" Enums(asc, desc)
// @Param new_pool query bool false "Filter for new pools"
// @Param burnt query bool false "Filter for burnt pools"
// @Param dexscreener_spent query bool false "Filter for dexscreener spent pools"
// @Param filters[] query []string false "Additional filters"
// @Param min_quote_usd query int false "Minimum quote USD"
// @Param max_quote_usd query int false "Maximum quote USD"
// @Param min_marketcap query int false "Minimum market cap"
// @Param max_marketcap query int false "Maximum market cap"
// @Param min_volume query int false "Minimum volume"
// @Param max_volume query int false "Maximum volume"
// @Param min_swaps1h query int false "Minimum swaps 1h"
// @Param max_swaps1h query int false "Maximum swaps 1h"
// @Param min_holder_count query int false "Minimum holder count"
// @Param max_holder_count query int false "Maximum holder count"
// @Param min_created query string false "Minimum creation time"
// @Param max_created query string false "Maximum creation time"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rank/sol/new_pairs [get]
func GetNewPairRanks(c *gin.Context) {
	var req request.SolRankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(response.CodeParamErr, "Invalid request parameters", err))
		return
	}
	result := service.GetNewPairRanks(&req)
	c.JSON(result.Code, result)
}
