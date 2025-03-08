package api

import (
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"
	"game-fun-be/internal/service"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TickersHandler struct {
	tickerService *service.TickerServiceImpl
}

func NewTickersHandler(tickerService *service.TickerServiceImpl) *TickersHandler {
	return &TickersHandler{tickerService: tickerService}
}

// Tickers 获取市场行情
// @Summary 获取市场行情数据
// @Description 根据链类型、排序、分页和搜索条件获取市场行情数据。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param sorted_by query string false "排序字段，支持以下值：MARKET_CAP, PRICE_CHANGE_5M, PRICE_CHANGE_1H, PRICE_CHANGE_24H, NATIVE_VOLUME_1H, NATIVE_VOLUME_24H, TX_COUNT_24H, HOLDERS, INITIALIZE_AT, Links" Enums(MARKET_CAP, PRICE_CHANGE_5M, PRICE_CHANGE_1H, PRICE_CHANGE_24H, NATIVE_VOLUME_1H, NATIVE_VOLUME_24H, TX_COUNT_24H, HOLDERS, INITIALIZE_AT, Links) example("INITIALIZE_AT")
// @Param sort_direction query string false "排序方向，支持以下值：DESC, ASC" Enums(DESC, ASC) example("DESC")
// @Param page_cursor query string false "分页游标，用于分页查询" example("")
// @Param limit query int false "每页返回的数据条数" example(50)
// @Param search query string false "搜索关键字，用于筛选数据" example("")
// @Param new_pairs_resolution query string false "新交易对的时间分辨率，例如 1D（1 天）" Enums(1D, 1H, 1M) example("1D")
// @Success 200 {object} response.Response{data=response.TickersResponse} "成功返回市场行情数据"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers/{chain_type} [get]
func (t *TickersHandler) Tickers(c *gin.Context) {
	var req request.TickersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "Invalid request parameters", err))
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := t.tickerService.Tickers(req, chainType)
	c.JSON(res.Code, res)
}

// TickerDetail 获取 Ticker 详情
// @Summary 获取 Ticker 详情
// @Description 根据链类型和代币地址获取 Ticker 的详细信息。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param ticker_address path string true "代币地址"
// @Success 200 {object} response.Response{data=response.GetTickerResponse} "成功返回 Ticker 详情"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers/{chain_type}/detail/{ticker_address} [get]
func (t *TickersHandler) TickerDetail(c *gin.Context) {
	tickerAddress := c.Param("ticker_address")
	if tickerAddress == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "ticker_address cannot be empty", errors.New("ticker_address is required")))
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := t.tickerService.TickerDetail(tickerAddress, chainType)
	c.JSON(res.Code, res)
}

// TickerDetail 获取 Ticker 详情
// @Summary 获取 Ticker 详情
// @Description 根据链类型和代币地址获取 Ticker 的详细信息。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param ticker_address path string true "代币地址"
// @Success 200 {object} response.Response{data=response.GetTickerResponse} "成功返回 Ticker 详情"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers/{chain_type}/market/{ticker_address} [get]
func (t *TickersHandler) MarketTicker(c *gin.Context) {
	tickerAddress := c.Param("ticker_address")
	if tickerAddress == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "ticker_address cannot be empty", errors.New("ticker_address is required")))
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := t.tickerService.MarketTicker(tickerAddress, chainType)
	c.JSON(res.Code, res)
}

// SwapHistories 获取交易历史记录
// @Summary 获取指定 Ticker 的交易历史记录
// @Description 根据链类型和代币地址获取交易历史记录。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param ticker_address path string true "代币地址"
// @Success 200 {object} response.Response{data=response.SwapHistoriesResponse} "成功返回交易历史记录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers/{chain_type}/swap_histories/{ticker_address} [get]
func (t *TickersHandler) SwapHistories(c *gin.Context) {
	tickerAddress := c.Param("ticker_address")
	if tickerAddress == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "ticker_address cannot be empty", errors.New("ticker_address is required")))
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}

	res := t.tickerService.SwapHistories(tickerAddress, chainType)
	c.JSON(res.Code, res)
}

// TokenDistribution 获取代币分布信息
// @Summary 获取指定 Ticker 的代币分布信息
// @Description 根据链类型和代币地址获取代币持有者的分布信息。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param ticker_address path string true "代币地址"
// @Success 200 {object} response.Response{data=response.TokenDistributionResponse} "成功返回代币分布信息"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers/{chain_type}/token_distribution/{ticker_address} [get]
func (t *TickersHandler) TokenDistribution(c *gin.Context) {
	tickerAddress := c.Param("ticker_address")
	if tickerAddress == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "ticker_address cannot be empty", errors.New("ticker_address is required")))
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	res := t.tickerService.TokenDistribution(tickerAddress, chainType)
	c.JSON(res.Code, res)
}

// SearchTickers 根据条件搜索 Tickers
// @Summary 搜索 Tickers
// @Description 根据链类型、搜索参数、分页参数等条件搜索 Tickers。支持的链类型：sol（Solana）、eth（Ethereum）、bsc（Binance Smart Chain）。
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param chain_type path string true "链类型（sol、eth、bsc）"
// @Param param path string true "搜索参数（如代币名称或地址）"
// @Param limit query string true "分页大小"
// @Param cursor query string false "分页游标"
// @Success 200 {object} response.SearchTickerResponse "成功返回 Tickers 列表"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers/{chain_type}/search [get]
func (t *TickersHandler) SearchTickers(c *gin.Context) {
	param := c.Param("param")
	if param == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "tickers_id cannot be empty", errors.New("tickers_id is required")))
		return
	}
	limit, errResp := GetLimit(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	chainType, errResp := ParseChainTypeWithResponse(c)
	if errResp != nil {
		c.JSON(errResp.Code, errResp)
		return
	}
	cursor := c.Param("cursor")
	res := t.tickerService.SearchTickers(param, limit, cursor, chainType)
	c.JSON(res.Code, res)
}
