package server

import (
	"time"

	"game-fun-be/internal/api"
	apiadmin "game-fun-be/internal/api/admin"
	"game-fun-be/internal/api/ws"
	"game-fun-be/internal/clickhouse"
	"game-fun-be/internal/interceptor"
	"game-fun-be/internal/model"
	"game-fun-be/internal/service"

	"github.com/IBM/sarama"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter 路由配置
func NewRouter(producer sarama.SyncProducer) *gin.Engine {
	globalService := service.NewGlobalServiceImpl()
	globalHandler := api.NewGlobalHandler(globalService)

	userInfoRepo := model.NewUserInfoRepo()
	userService := service.NewUserServiceImpl(userInfoRepo)
	userHandler := api.NewUserHandler(userService)
	platformTokenStatisticRepo := model.NewPlatformTokenStatisticRepo()

	tickerInfoRepo := model.NewTokenInfoRepo()
	TokenMarketAnalyticsRepo := clickhouse.NewTokenMarketAnalyticsRepo()
	tickerService := service.NewTickerServiceImpl(tickerInfoRepo, TokenMarketAnalyticsRepo)
	platformTokenStatisticService := service.NewPlatformTokenStatisticServiceImpl(platformTokenStatisticRepo)
	tickerHandler := api.NewTickersHandler(tickerService, platformTokenStatisticService)

	tokenHoldingsService := service.NewTokenHoldingsServiceImpl()
	tokenHoldingsHandler := api.NewTokenHoldingsHandler(tokenHoldingsService)

	pointRecordsRepo := model.NewPointRecordsRepo()

	pointsService := service.NewPointsServiceImpl(userInfoRepo, pointRecordsRepo, platformTokenStatisticRepo)
	pointsHandler := api.NewPointsHandler(pointsService, globalService)

	swapService := service.NewSwapService(producer)
	swapHandler := api.NewSwapHandler(swapService)
	tokenConfigService := service.NewTokenConfigServiceImpl()
	tokenConfigHandler := apiadmin.NewAdminTokenConfigHandler(tokenConfigService)

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
		v1.GET("tickers/:chain_type/detail/:ticker_address", tickerHandler.TickerDetail)
		v1.GET("tickers/:chain_type/statistic/:ticker_address", tickerHandler.TickerPointsStatistic)
		v1.GET("tickers/:chain_type/market/:ticker_address", tickerHandler.MarketTicker)
		v1.GET("tickers/:chain_type/search", tickerHandler.SearchTickers)
		v1.GET("tickers/:chain_type/swap_histories/:ticker_address", tickerHandler.SwapHistories)
		v1.GET("tickers/:chain_type/token_distribution/:ticker_address", tickerHandler.TokenDistribution)
		v1.GET("klines/:klineType/:chainType/:tokenAddress", tickerHandler.GetTokenKlines)
		v1.GET("global/:chain_type/native_token_price", globalHandler.NativeTokePrice)

		auth := v1.Group("")
		// token登陆验证路由
		auth.Use(interceptor.BearerAuth())
		auth.GET("users/:chain_type/my_info", userHandler.MyInfo)
		auth.GET("users/:chain_type/invite/code", userHandler.InviteCode)
		auth.GET("points/:chain_type", pointsHandler.Points)
		auth.GET("points/:chain_type/detail", pointsHandler.PointsDetail)
		auth.GET("points/:chain_type/invite/detail", pointsHandler.InvitedPointsDetail)
		auth.GET("points/:chain_type/estimated", pointsHandler.PointsEstimated)
		auth.GET("global/:chain_type/native_balance", globalHandler.NativeBalance)
		auth.GET("global/:chain_type/ticker_balance/:ticker_address", globalHandler.TickerBalance)
		auth.GET("swap/:chain_type/get_transaction", swapHandler.GetTransaction)
		auth.POST("swap/:chain_type/send_transaction", swapHandler.SendTransaction)
		auth.GET("swap/:chain_type/transaction_status", swapHandler.TransactionStatus)
		auth.GET("token_holdings/:chain_type/:account", tokenHoldingsHandler.TokenHoldings)

		// WebSocket 路由
		v1.GET("ws/kline/:tokenAddress", ws.HandleKlineWS)
		v1.GET("ws/market_analytics/:tokenAddress", ws.HandleMarketAnalyticsWS)
	}

	// 健康检查路由
	r.GET("/health", api.HealthCheck)
	// 工具路由
	r.GET("/tools/execute_reindex_job", api.ExecuteReindexJob)
	r.POST("/tools/reset_pool_info", api.ResetTokenPoolInfo)

	// 仅限管理员的API路由组，需要API Key认证
	adminAPI := r.Group("/admin")
	// 应用API Key认证中间件仅到admin路由组
	adminAPI.Use(interceptor.ApiKeyAuth())
	{
		// 获取token列表
		adminAPI.GET("/tokenconfigs/list", tokenConfigHandler.GetAdminTokenConfigList)
		// 获取token详情
		adminAPI.GET("/tokenconfigs/detail/:id", tokenConfigHandler.GetTokenConfig)
		// 创建token
		adminAPI.POST("/tokenconfigs/create", tokenConfigHandler.CreateTokenConfig)
		// 更新token
		adminAPI.POST("/tokenconfigs/update/:id", tokenConfigHandler.UpdateTokenConfig)
		// 删除token
		adminAPI.GET("/tokenconfigs/delete/:id", tokenConfigHandler.DeleteTokenConfig)
	}

	return r
}
