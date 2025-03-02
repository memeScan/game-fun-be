package api

import (
	"my-token-ai-be/internal/response"
	"my-token-ai-be/internal/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetSolGasFee godoc
// @Summary Get Solana gas fee
// @Description Get the current gas fee for Solana blockchain
// @Tags on-chain-data
// @Accept json
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Produce json
// @Success 200 {object} response.Response{data=httpRespone.GasFee}
// @Failure 500 {object} response.Response
// @Router /chains/:chainType/gas_fee [get]
func GetSolGasFee(c *gin.Context) {
	chainType := c.Param("chainType")

	onChainDataService := service.NewOnChainDataService()
	res := onChainDataService.GetNativeTokenGasFee(chainType)

	c.JSON(res.Code, res)
}

// GetSolPrice godoc
// @Summary Get Solana real-time price
// @Description Get the current real-time price for Solana
// @Tags on-chain-data
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=float64}
// @Failure 500 {object} response.Response
// @Router /chains/:chainType/sol_price [get]
func GetSolPrice(c *gin.Context) {
	onChainDataService := service.NewOnChainDataService()
	res := onChainDataService.GetSolPrice()
	c.JSON(res.Code, res)
}

// @Summary Get token balance
// @Description Get the balance of a token for a given owner
// @Tags on-chain-data
// @Accept json
// @Produce json
// @Param chainType path string true "Type of the token (e.g., sol, eth)"
// @Param owner path string true "Owner of the token"
// @Param token path string true "Token address"
// @Success 200 {object} response.Response{data=response.TokenBalance}
// @Failure 500 {object} response.Response
// @Router /token_balance/{chainType}/{owner}/{token} [get]
func GetTokenBalance(c *gin.Context) {
	tokenBalance := service.NewOnChainDataService().GetTokenBalance(c.Param("chainType"), c.Param("owner"), c.Param("token"))
	c.JSON(tokenBalance.Code, tokenBalance)
}

// GetTokenKlines godoc
// @Summary Get token kline data
// @Description Get historical kline (candlestick) data for a specific token
// @Tags on-chain-data
// @Accept json
// @Produce json
// @Param klineType path string true "Kline type" Enums(kline, mcapkline)
// @Param chainType path string true "Chain type" Enums(sol, eth, btc)
// @Param tokenAddress path string true "Token address"
// @Param resolution query string true "Resolution (1S=1sec, 1=1min, 5=5min, 15=15min, 60=1hour, 240=4hour, 720=12hour, 1D=1day)"
// @Success 200 {object} response.Response{data=[]response.KlineData}
// @Failure 400 {object} response.Response{data=string} "Invalid parameters"
// @Failure 500 {object} response.Response{data=string} "Server error"
// @Router /tokens/{klineType}/{chainType}/{tokenAddress} [get]
func GetTokenKlines(c *gin.Context) {
	// klineType := c.Param("klineType") TODO: 后续再支持链类型和klineType
	// chainType := c.Param("chainType")
	tokenAddress := c.Param("tokenAddress")
	resolution := c.Query("resolution")

	// 解析时间参数
	startTs, err := strconv.ParseInt(c.Query("from"), 10, 64)
	if err != nil {
		c.JSON(400, response.Err(response.CodeParamErr, "Invalid from timestamp", err))
		return
	}

	endTs, err := strconv.ParseInt(c.Query("till"), 10, 64)
	if err != nil {
		c.JSON(400, response.Err(response.CodeParamErr, "Invalid till timestamp", err))
		return
	}

	start := time.Unix(startTs, 0)
	end := time.Unix(endTs, 0)

	// 验证时间范围
	if startTs >= endTs {
		c.JSON(400, response.Err(response.CodeParamErr, "Start time must be before end time", nil))
		return
	}

	// 验证 resolution 参数
	validResolutions := map[string]bool{
		"1S":  true, // 1 second
		"1":   true, // 1 minute
		"5":   true, // 5 minutes
		"15":  true, // 15 minutes
		"60":  true, // 1 hour
		"240": true, // 4 hours
		"720": true, // 12 hours
		"1D":  true, // 1 day
	}

	resolution = strings.ToUpper(resolution)
	if !validResolutions[resolution] {
		c.JSON(400, response.Err(response.CodeParamErr, "Invalid resolution", nil))
		return
	}

	// 调用 service 获取数据
	klineService := service.NewKlineService()
	klines, err := klineService.GetTokenKlines(tokenAddress, resolution, start, end)
	if err != nil {
		c.JSON(500, response.Err(response.CodeServerUnknown, err.Error(), err))
		return
	}

	// 空数据处理
	if len(klines) == 0 {
		c.JSON(200, response.BuildResponse([]response.KlineData{}, 200, "no data found", nil))
		return
	}

	decimals := uint8(6) // FIXME: 暂时取 6 ，后续从token_info表中获取

	// 转换为前端需要的格式
	klineDataList := response.BuildKlineDataList(klines, decimals)
	c.JSON(200, response.BuildResponse(klineDataList, 200, "success", nil))
}
