package server

import (
	"my-token-ai-be/internal/api"
	"my-token-ai-be/internal/api/ws"
	"my-token-ai-be/internal/interceptor"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
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
		v1.GET("user/message", api.GetMessage)
		v1.POST("user/wallet-login", api.WalletLogin)
		v1.GET(":chainType/:tradeType/get_swap_route", api.GetSwapRoute)
		v1.GET("token_launchpad_info/:chainType/:tokenAddress", api.GetTokenLaunchpadInfo)
		v1.GET("token_pool_info_sol/:chainType/:tokenAddress", api.GetMarketInfo)
		v1.GET("token_info/:chainType/:tokenAddress", api.GetTokenInfo)
		v1.POST("token_prices/:chainType", api.GetTokenPrices)
		v1.GET("rank/sol/pump", api.GetSolPumpRank)
		v1.GET("rank/sol/new_pairs", api.GetSolDumpRank)
		v1.GET("tokens/token_balance/:chainType/:owner/:token", api.GetTokenBalance)
		v1.GET("chains/:chainType/gas_fee", api.GetSolGasFee)
		v1.GET("chains/:chainType/native_token_price", api.GetSolPrice)
		v1.GET("token_market_analytics/:chainType/:tokenAddress", api.GetTokenMarketAnalytics)
		v1.GET("token_order_book/:chainType/:tokenAddress", api.GetTokenOrderBook)
		v1.GET("tokens/search/:chainType/:tokenAddress", api.Search)
		v1.GET("tokens/:klineType/:chainType/:tokenAddress", api.GetTokenKlines)
		v1.GET("transaction/sol/send_swap_transaction", api.SendSwapRequest)
		v1.GET("transaction/sol/get_swap_request_status", api.GetSwapRequestStatus)
		v1.GET("search_documents_job", api.SearchDocumentsJob)
		v1.GET("rank/sol/swap", api.GetSolSwapRank)
		v1.GET("token_base_info/:chainType/:tokenAddress", api.GetTokenBaseInfo)
		v1.GET("token_check_info/:chainType/:tokenAddress/:tokenPool", api.GetTokenCheckInfo)
		v1.GET("token_market_analytics_search/:chainType/:tokenAddress", api.GetTokenMarketQuery)
		v1.GET("pairs/:chainType/new_pair_ranks", api.GetNewPairRanks)
		// 代币信息同步
		v1.GET("tools/token_info_sync", api.TokenInfoSyncJob)
		auth := v1.Group("")
		auth.Use(interceptor.AuthRequired())
		auth.Use(interceptor.CurrentUser())
		{
			auth.GET("user/getuser", api.UserMe)
		}

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
