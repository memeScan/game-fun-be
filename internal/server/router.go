package server

import (
	"game-fun-be/internal/api"
	"game-fun-be/internal/api/ws"
	"game-fun-be/internal/interceptor"
	"game-fun-be/internal/model"
	"game-fun-be/internal/service"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {

	globalService := service.NewGlobalServiceImpl()

	globalHandler := api.NewGlobalHandler(globalService)
	userInfoRepo := model.NewUserInfoRepo()
	userService := service.NewUserServiceImpl(userInfoRepo)
	userHandler := api.NewUserHandler(userService)

	tickerService := service.NewTickerServiceImpl()
	tickerHandler := api.NewTickersHandler(tickerService)

	tokenHoldingsService := service.NewTokenHoldingsServiceImpl()
	tokenHoldingsHandler := api.NewTokenHoldingsHandler(tokenHoldingsService)

	pointsService := service.NewPointsServiceImpl()
	pointsHandler := api.NewPointsHandler(pointsService)

	swapService := service.NewSwapService()
	swapHandler := api.NewSwapHandler(swapService)

	r := gin.New()

	// 基础中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Cache-Control",
			"X-Requested-With",
			"Referer",
			"Sec-Fetch-Mode",
			"Sec-Fetch-Site",
			"Sec-Fetch-Dest",
			"Sec-WebSocket-Key",
			"Sec-WebSocket-Version",
			"Sec-WebSocket-Extensions",
			"Sec-WebSocket-Protocol",
			"Upgrade",
			"Connection",
		},
		AllowCredentials: false,
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
		},
		AllowWildcard: true,
		MaxAge:        12 * time.Hour,
	}))

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// v1 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("users/:chain_type/login", userHandler.Login)
		v1.GET("tickers/:chain_type", tickerHandler.Tickers)
		v1.GET("tickers/:chain_type/:ticker_address", tickerHandler.TickerDetail)
		v1.GET("tickers/:chain_type/search", tickerHandler.SearchTickers)
		v1.GET("tickers/:chain_type/swap_histories/:ticker_address", tickerHandler.SwapHistories)
		v1.GET("tickers/:chain_type/token_distribution/:ticker_address", tickerHandler.TokenDistribution)
		v1.GET("token_holdings/:chain_type/:account", tokenHoldingsHandler.TokenHoldings)
		v1.GET("token_holdings/:chain_type/histories/:account", tokenHoldingsHandler.TokenHoldingsHistories)
		v1.GET("swap/:chainType/get_transaction", swapHandler.GetTransaction)
		v1.GET("swap/:chainType/send_transaction", swapHandler.SendTransaction)
		v1.GET("swap/:chainType/transaction_status", swapHandler.TransactionStatus)
		v1.GET("global/:chain_type/native_token_price", globalHandler.NativeTokePrice)

		auth := v1.Group("")
		// token登陆验证路由
		auth.Use(interceptor.AuthRequired())
		auth.GET("users/:chain_type/my_info", userHandler.MyInfo)
		auth.GET("users/:chain_type/invite/code", userHandler.InviteCode)
		auth.GET("points/:chain_type", pointsHandler.Points)
		auth.GET("points/:chain_type/detail", pointsHandler.PointsDetail)
		auth.GET("points/:chain_type/estimated", pointsHandler.PointsEstimated)
		auth.GET("global/:chain_type/balance", globalHandler.Balance)

		v1.GET("tokens/:klineType/:chainType/:tokenAddress", api.GetTokenKlines)

		// WebSocket 路由
		v1.GET("ws/kline/:tokenAddress", ws.HandleKlineWS)
		v1.GET("ws/market_analytics/:tokenAddress", ws.HandleMarketAnalyticsWS)
	}

	// 健康检查路由
	r.GET("/health", api.HealthCheck)
	// 工具路由
	r.POST("/tools/execute_reindex_job", api.ExecuteReindexJob)
	r.POST("/tools/reset_pool_info", api.ResetTokenPoolInfo)

	return r
}
