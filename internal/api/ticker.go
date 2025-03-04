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
// @Success 200 {object} response.Response{data=response.TickersResponse} "成功返回市场行情数据"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers [get]
func (t *TickersHandler) Tickers(c *gin.Context) {
	var req request.TickersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "Invalid request parameters", err))
		return
	}
	res := t.tickerService.Tickers(req)
	c.JSON(res.Code, res)
}

// TickerDetail 获取 Ticker 详情
// @Summary 获取 Ticker 详情
// @Description 根据 token_symbol 获取 Ticker 的详细信息
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param token_symbol path string true "代币符号" Enums(SUPER, BTC, ETH, USDT, BNB) example("SUPER")
// @Success 200 {object} response.Response{data=response.GetTickerResponse} "成功返回 Ticker 详情"
// @Failure 400 {object} response.Response "参数错误"
// @Router /tickers/{token_symbol} [get]
func (t *TickersHandler) TickerDetail(c *gin.Context) {
	tokenSymbol := c.Param("token_symbol")
	if tokenSymbol == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "token_symbol cannot be empty", errors.New("token_symbol is required")))
		return
	}
	res := t.tickerService.TickerDetail(tokenSymbol)
	c.JSON(res.Code, res)
}

// SwapHistories 获取交易历史记录
// @Summary 获取指定 Ticker 的交易历史记录
// @Description 根据 Ticker ID 获取交易历史记录
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param tickers_id path string true "Ticker ID" example("SUPER")
// @Success 200 {object} response.Response{data=response.SwapHistoriesResponse} "成功返回交易历史记录"
//
//	@Example ExampleResponse {
//	  "code": 200,
//	  "msg": "success",
//	  "data": {
//	    "transaction_histories": [
//	      {
//	        "market_id": 1,
//	        "is_buy": true,
//	        "payer": "BPeE27Y13ChxWHVceuv4JdWdZjh9tnc9xxgbJoxYbXLq",
//	        "recipient": "669VYcBRq51iQzFiPTQcsW2CsvLfHM9AwVmaoM1mAAR7",
//	        "signature": "3XTc73NfWb92PFv52ZDNSabNui4pyeWKaFB8cCjZR1SDWqgKDBaojWUCssRF5ijwvYDKF3k4MuGAakmDeDX9pwnd",
//	        "block_time": "2025-03-02T03:57:35Z",
//	        "index": 1,
//	        "slot": 324011096,
//	        "token_amount": "1987.143391",
//	        "native_amount": "0.100251255",
//	        "fee": "0.000501256",
//	        "total_native_amount": "0.100251255",
//	        "id": 749019,
//	        "payer_profile": {
//	          "user_id": "10000661",
//	          "avatar": "",
//	          "username": "",
//	          "nickname": "Zebra8fsmP"
//	        }
//	      }
//	    ],
//	    "has_more": false
//	  }
//	}
//
// @Failure 400 {object} response.Response "参数错误"
// @Router /tickers/swap_histories/{tickers_id} [get]
func (t *TickersHandler) SwapHistories(c *gin.Context) {
	tickersId := c.Param("tickers_id")
	if tickersId == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "tickers_id cannot be empty", errors.New("tickers_id is required")))
		return
	}
	// chainType := c.Param("chainType")
	// tokenAddress := c.Param("tokenAddress")
	// tokenInfoService := &service.TokenInfoService{}
	// orderBook := tokenInfoService.GetTokenOrderBook(tokenAddress, uint8(model.FromString(chainType)))
	// c.JSON(orderBook.Code, orderBook)

	res := t.tickerService.SwapHistories(tickersId)
	c.JSON(res.Code, res)
}

// TokenDistribution 获取代币分布信息
// @Summary 获取指定 Ticker 的代币分布信息
// @Description 根据 Ticker ID 获取代币持有者的分布信息
// @Tags 市场行情
// @Accept json
// @Produce json
// @Param tickers_id path string true "Ticker ID" example("SUPER")
// @Success 200 {object} response.Response{data=response.TokenDistributionResponse} "成功返回代币分布信息"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tickers/{tickers_id}/distribution [get]
func (t *TickersHandler) TokenDistribution(c *gin.Context) {
	tickersId := c.Param("tickers_id")
	if tickersId == "" {
		c.JSON(http.StatusBadRequest, response.Err(http.StatusBadRequest, "tickers_id cannot be empty", errors.New("tickers_id is required")))
		return
	}
	res := t.tickerService.TokenDistribution(tickersId)
	c.JSON(res.Code, res)
}

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
	cursor := c.Param("cursor")
	res := t.tickerService.SearchTickers(param, limit, cursor)
	c.JSON(res.Code, res)
}
