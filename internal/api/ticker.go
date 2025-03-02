package api

import (
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Tickers 获取市场行情
// @Summary 获取市场行情数据
// @Description 根据排序、分页和搜索条件获取市场行情数据
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param sorted_by query string false "排序字段，支持以下值：MARKET_CAP, PRICE_CHANGE_5M, PRICE_CHANGE_1H, PRICE_CHANGE_24H, NATIVE_VOLUME_1H, NATIVE_VOLUME_24H, TX_COUNT_24H, HOLDERS, INITIALIZE_AT, Links" Enums(MARKET_CAP, PRICE_CHANGE_5M, PRICE_CHANGE_1H, PRICE_CHANGE_24H, NATIVE_VOLUME_1H, NATIVE_VOLUME_24H, TX_COUNT_24H, HOLDERS, INITIALIZE_AT, Links) example("INITIALIZE_AT")
// @Param sort_direction query string false "排序方向，支持以下值：DESC, ASC" Enums(DESC, ASC) example("DESC")
// @Param page_cursor query string false "分页游标，用于分页查询" example("")
// @Param limit query int false "每页返回的数据条数" example(50)
// @Param search query string false "搜索关键字，用于筛选数据" example("")
// @Param new_pairs_resolution query string false "新交易对的时间分辨率，例如 1D（1 天）" Enums(1D, 1H, 1M) example("1D")
// @Success 200 {object} response.TickersResponse "成功返回市场行情数据"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers [get]
func Tickers(c *gin.Context) {
	var req request.TickersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "Invalid request parameters", err))
		return
	}
	tickerService := service.NewTickerService()
	response := tickerService.Tickers(req)
	c.JSON(200, response)
}

// GetTicker 获取 Ticker 详情
// @Summary 获取 Ticker 详情
// @Description 根据 token_symbol 获取 Ticker 的详细信息
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param token_symbol path string true "代币符号" Enums(SUPER, BTC, ETH, USDT, BNB) example("SUPER")
// @Success 200 {object} response.GetTickerResponse "成功返回 Ticker 详情"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers/{token_symbol} [get]
func GetTicker(c *gin.Context) {
	tokenSymbol := c.Param("token_symbol")
	if tokenSymbol == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "token_symbol cannot be empty", errors.New("token_symbol is required")))
		return
	}
	tickerService := service.NewTickerService()
	response := tickerService.GetTicker(tokenSymbol)
	c.JSON(200, response)
}
