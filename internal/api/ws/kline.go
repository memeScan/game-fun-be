package ws

import (
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/response"
	"game-fun-be/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需要proper检查
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	readWriteTimeout = 100 * time.Second
	heartbeatPeriod  = 30 * time.Second
)

// @Summary K线数据 WebSocket 接口
// @Description 通过 WebSocket 实时推送 K 线数据。连接成功后，服务器将每秒推送最新的K线数据。
// @Description 连接格式: ws://host/ws/kline/{tokenAddress}?resolution={resolution}
// @Description 示例: ws://localhost:8080/ws/kline/0x1234...?resolution=1S
// @Tags WebSocket
// @Accept  json
// @Produce  json
// @Param tokenAddress path string true "代币地址" example(0x1234567890abcdef1234567890abcdef12345678)
// @Param resolution query string true "K线周期" Enums(1S,1,5,15,60,240,720,1D) example(1S)
// @Success 101 {object} response.KlineData{timestamp=int64,open=string,high=string,low=string,close=string,volume=string} "WebSocket连接成功后的推送数据格式"
// @Router /ws/kline/{tokenAddress} [get]
func HandleKlineWS(c *gin.Context) {
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

	// 数据推送定时器
	pushTicker := time.NewTicker(getUpdateInterval(c.Query("resolution")))
	defer pushTicker.Stop()

	klineService := service.NewKlineService()
	var lastPrice string
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
	if klineData, err := getLatestKlineData(klineService, c.Param("tokenAddress"), decimals); err == nil {
		lastPrice = klineData.Close.String()
		if err := ws.WriteJSON(klineData); err != nil {
			return
		}
	}

	// 主循环
	for {
		select {
		case <-heartbeatTicker.C: // 心跳
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				util.Log().Error("WebSocket heartbeat failed: %v, tokenAddress: %s",
					err, c.Param("tokenAddress"))
				return
			}
			ws.SetWriteDeadline(time.Now().Add(readWriteTimeout))

		case <-pushTicker.C: // 数据推送
			klineData, err := getLatestKlineData(klineService, c.Param("tokenAddress"), decimals)
			if err != nil {
				continue
			}

			// 只在价格变化时推送
			currentPrice := klineData.Close.String()
			if currentPrice != lastPrice {
				if err := ws.WriteJSON(klineData); err != nil {
					util.Log().Error("WebSocket write failed: %v, tokenAddress: %s",
						err, c.Param("tokenAddress"))
					return
				}
				// 推送成功后才更新lastPrice
				lastPrice = currentPrice
			}
		}
	}
}

// 获取最新K线数据
func getLatestKlineData(klineService *service.KlineService, tokenAddress string,
	decimals uint8) (*response.KlineData, error) {

	latestKline, err := klineService.GetLatestKline(tokenAddress)
	if err != nil {
		util.Log().Error("获取K线数据失败: %v", err)
		return nil, err
	}

	klineData := response.BuildKlineData(*latestKline, decimals)
	return &klineData, nil
}

func getUpdateInterval(resolution string) time.Duration {
	switch resolution {
	case "1S":
		return time.Millisecond * 100
	case "1":
		return time.Millisecond * 100
	case "5":
		return time.Millisecond * 100
	case "15":
		return time.Millisecond * 200
	case "60":
		return time.Millisecond * 200
	case "240":
		return time.Millisecond * 200
	case "720":
		return time.Millisecond * 200
	case "1D":
		return time.Minute * 5
	default:
		return time.Second * 10
	}
}
