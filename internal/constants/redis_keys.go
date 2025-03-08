package constants

const (
	UserTokenKeyFormat = "auth:token"
	TokenMarketData    = "token:market:data"
	TokenMetaData      = "token:metadata"
	TokenDistribution  = "token:trade:distribution:data"
	TokenTradeData     = "token:trade:data"
	// RedisKeySolLatestPrice 存储最新 SOL 价格的 Redis 键
	RedisKeySolLatestPrice = "sol_latest_price"
	// SolPriceMultiplier SOL 价格的乘数
	SolPriceMultiplier = 100000000 // 10^8
	// RedisKeyHotTokens 热门代币集合的 Redis 键
	RedisKeyHotTokens = "hot_tokens_zset"
	// RedisKeySafetyCheck 安全检查的 Redis 键
	RedisKeySafetyCheck = "token:safety:"
	// RedisKeyTradingPool6hMarketInfo 交易池6小时市场信息
	RedisKeyTradingPool6hMarketInfo = "trading_pool_6h_market_info"
	// RedisKeyTradingPool24hMarketInfo 交易池24小时市场信息
	RedisKeyTradingPool24hMarketInfo = "trading_pool_24h_market_info"
	// RedisKeyCompletedTokens 已完成代币集合的 Redis 键
	RedisKeyCompletedTokens = "completed_tokens_zset"
	// SwapRedisKeyTokens1m 1分钟代币集合的 Redis 键
	SwapRedisKeyTokens1m = "swap_tokens_1m_zset"
	// SwapRedisKeyTokens5m 5分钟代币集合的 Redis 键
	SwapRedisKeyTokens5m = "swap_tokens_5m_zset"
	// SwapRedisKeyTokens1h 1小时代币集合的 Redis 键
	SwapRedisKeyTokens1h = "swap_tokens_1h_zset"
	// SwapRedisKeyTokens6h 6小时代币集合的 Redis 键
	SwapRedisKeyTokens6h = "swap_tokens_6h_zset"
	// SwapRedisKeyTokens1d 1天代币集合的 Redis 键
	SwapRedisKeyTokens1d = "swap_tokens_1d_zset"

	// RedisKeyPrefixTokenInfo 代币信息的 Redis key 前缀
	RedisKeyPrefixTokenInfo = "token:pump:info"

	// RedisKeyTokenTransactionID 代币交易记录的分布式ID生成器的Redis键
	RedisKeyTokenTransactionID = "token:transaction:id"
)
