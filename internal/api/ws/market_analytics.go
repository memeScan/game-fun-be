package ws

import (
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/service"
	"my-token-ai-be/internal/response"
	"github.com/gin-gonic/gin"
	"time"
)

func HandleMarketAnalyticsWS(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		util.Log().Error("WebSocket upgrade failed: %v", err)
		return
	}
	defer func() {
		ws.Close()
		util.Log().Info("WebSocket connection closed: %s", c.Param("tokenAddress"))
	}()

	// 设置写入超时
	ws.SetWriteDeadline(time.Now().Add(readWriteTimeout))

	// 心跳检测定时器
	heartbeatTicker := time.NewTicker(heartbeatPeriod)
	defer heartbeatTicker.Stop()

	klineService := service.NewKlineService()

	decimals := uint8(6)

	// 等待客户端发送订阅确认消息
	_, message, err := ws.ReadMessage()
	if err != nil {
		util.Log().Error("Failed to read subscribe message: %v", err)
		return
	}

	// 验证订阅消息
	if string(message) != "subscribe" {
		util.Log().Error("Invalid subscribe message: %s", message)
		return
	}

	// 初始推送
	if klineData, err := getMarketAnalyticsData(klineService, c.Param("tokenAddress"), decimals); err == nil {
		if err := ws.WriteJSON(klineData); err != nil {
			return
		}
	}

}

// 获取k线数据
func getMarketAnalyticsData(klineService *service.KlineService, tokenAddress string,
	decimals uint8) (*response.KlineData, error) {

	interval := "60"
	// 获取过去1天的数据	
	start := time.Now().AddDate(0, -1, 0)
	// 获取当前时间
	end := time.Now()

	latestKlines, err := klineService.GetTokenKlines(tokenAddress, interval, start, end)
	if err != nil {
		util.Log().Error("获取K线数据失败: %v", err)
		return nil, err
	}

	// 拿到了分钟级别的
	lastPrice := ""
	// 1分钟的总购买量
	totalBuyVolume := uint64(0)	
	// 1分钟的总销售量
	totalSellVolume := uint64(0)
	// 1分钟的买笔数
	totalBuyCount1m := uint64(0)
	// 1分钟的卖笔数
	totalSellCount1m := uint64(0)
	// 5分钟的总购买量
	totalBuyVolume5m := uint64(0)
	// 5分钟的总销售量
	totalSellVolume5m := uint64(0)
	// 5分钟的买笔数
	totalBuyCount5m := uint64(0)
	// 5分钟的卖笔数
	totalSellCount5m := uint64(0)
	// 1小时的总购买量
	totalBuyVolume1h := uint64(0)
	// 1小时的总销售量
	totalSellVolume1h := uint64(0)
	// 1小时的买笔数
	totalBuyCount1h := uint64(0)
	// 1小时的卖笔数
	totalSellCount1h := uint64(0)
	// 24小时的总购买量
	totalBuyVolume24h := uint64(0)
	// 24小时的总销售量
	totalSellVolume24h := uint64(0)
	// 24小时的买笔数
	totalBuyCount24h := uint64(0)
	// 24小时的卖笔数
	totalSellCount24h := uint64(0)

	// latestKlines 遍历 第一个拿到的关盘价格为最新价格
	for _, kline := range latestKlines {
		// 只有第一个遍历到的关盘价格为最新价格
		if lastPrice == "" {
			lastPrice = kline.ClosePrice.String()
		}
		
		if kline.IntervalTimestamp.Unix() / 60 == 0 {
			// 1分钟的总购买量
			totalBuyVolume = totalBuyVolume + kline.BuyVolume
			// 1分钟的总销售量
			totalSellVolume = totalSellVolume + kline.SellVolume

		// 1分钟的买笔数
		totalBuyCount1m += kline.BuyCount
		// 1分钟的卖笔数
		totalSellCount1m += kline.SellCount

		} else if kline.IntervalTimestamp.Unix() / 300 == 0 {
		// 5分钟的总购买量
		totalBuyVolume5m = totalBuyVolume5m + kline.BuyVolume
		// 5分钟的总销售量
		totalSellVolume5m = totalSellVolume5m + kline.SellVolume
		// 5分钟的买笔数
		totalBuyCount5m += kline.BuyCount
		// 5分钟的卖笔数
		totalSellCount5m += kline.SellCount

		} else if kline.IntervalTimestamp.Unix() / 3600 == 0 {
		// 1小时的总购买量
		totalBuyVolume1h = totalBuyVolume1h + kline.BuyVolume
		// 1小时的总销售量
		totalSellVolume1h = totalSellVolume1h + kline.SellVolume
		// 1小时的买笔数
		totalBuyCount1h += kline.BuyCount
		// 1小时的卖笔数
		totalSellCount1h += kline.SellCount
		} else {
			// 24小时的总购买量
			totalBuyVolume24h = totalBuyVolume24h + kline.BuyVolume
			// 24小时的总销售量
			totalSellVolume24h = totalSellVolume24h + kline.SellVolume
			// 24小时的买笔数
			totalBuyCount24h += kline.BuyCount
			// 24小时的卖笔数
			totalSellCount24h += kline.SellCount
		}

	}
	// 5分钟的总购买量m
	totalBuyVolume5m = totalBuyVolume5m + totalBuyVolume
	// 5分钟的总销售量
	totalSellVolume5m = totalSellVolume5m + totalSellVolume
	// 1小时的总购买量
	totalBuyVolume1h = totalBuyVolume1h + totalBuyVolume5m
	// 1小时的总销售量
	totalSellVolume1h = totalSellVolume1h + totalSellVolume5m
	// 24小时的总购买量
	totalBuyVolume24h = totalBuyVolume24h + totalBuyVolume1h
	// 24小时的总销售量
	totalSellVolume24h = totalSellVolume24h + totalSellVolume1h

	// 5分钟的购买笔数
	totalBuyCount5m = totalBuyCount5m + totalBuyCount1m
	// 5分钟的卖笔数
	totalSellCount5m = totalSellCount5m + totalSellCount1m
	// 1小时的购买笔数
	totalBuyCount1h = totalBuyCount1h + totalBuyCount5m
	// 1小时的卖笔数
	totalSellCount1h = totalSellCount1h + totalSellCount5m
	// 24小时的购买笔数
	totalBuyCount24h = totalBuyCount24h + totalBuyCount1h
	// 24小时的卖笔数
	totalSellCount24h = totalSellCount24h + totalSellCount1h

	return nil, nil
}

