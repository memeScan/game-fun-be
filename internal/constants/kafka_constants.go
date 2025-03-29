package constants

import "os"

// Kafka 消费组 ID
const (
	KafkaGroupDexProcessor = "dex_processor"
)

// 基础 topic 前缀常量
const (
	RayNewPoolPrefix         = "market.raydium.newpool."
	RaySwapPrefix            = "market.raydium.swap."
	RayAddLiquidityPrefix    = "market.raydium.addliquidity."
	RayRemoveLiquidityPrefix = "market.raydium.removeliquidity."

	PumpCreatePrefix    = "market.pump.create."
	PumpTradePrefix     = "market.pump.trade."
	PumpCompletePrefix  = "market.pump.complete."
	PumpSetParamsPrefix = "market.pump.setparams."

	UnknownTokenPrefix = "market.unknown.token."

	GameOutTradePrefix  = "market.game.out.trade."
	GameInTradePrefix   = "market.game.in.trade."
	PointTxStatusPrefix = "market.point.tx.status." // 积分交易链上状态检测

	// PumpSwap DEX 相关
	PumpAmmNewPoolPrefix         = "market.pump.amm.newpool."
	PumpAmmSwapPrefix            = "market.pump.amm.swap."
	PumpAmmAddLiquidityPrefix    = "market.pump.amm.addliquidity."
	PumpAmmRemoveLiquidityPrefix = "market.pump.amm.removeliquidity."
)

var (
	// 获取环境后缀
	envSuffix = func() string {
		env := os.Getenv("APP_ENV")
		if env == "release" {
			return "prod"
		}
		return "test" // 默认环境和其他环境都返回test
	}()

	// Raydium topics - 使用前缀常量
	TopicRayCreate          = RayNewPoolPrefix + envSuffix
	TopicRaySwap            = RaySwapPrefix + envSuffix
	TopicRayAddLiquidity    = RayAddLiquidityPrefix + envSuffix
	TopicRayRemoveLiquidity = RayRemoveLiquidityPrefix + envSuffix

	TopicPumpCreate    = PumpCreatePrefix + envSuffix
	TopicPumpTrade     = PumpTradePrefix + envSuffix
	TopicPumpComplete  = PumpCompletePrefix + envSuffix
	TopicPumpSetParams = PumpSetParamsPrefix + envSuffix

	TopicUnknownToken = UnknownTokenPrefix + envSuffix

	// Game trading topics
	TopicGameOutTrade = GameOutTradePrefix + envSuffix // 代理合约外盘买卖事件
	TopicGameInTrade  = GameInTradePrefix + envSuffix  // 代理合约内盘买事件（积分兑换买）

	// 积分交易状态检测
	TopicPointTxStatus = PointTxStatusPrefix + envSuffix // 积分交易链上状态检测

	// PumpSwap DEX 相关
	TopicPumpAmmNewPool         = PumpAmmNewPoolPrefix + envSuffix
	TopicPumpAmmSwap            = PumpAmmSwapPrefix + envSuffix
	TopicPumpAmmAddLiquidity    = PumpAmmAddLiquidityPrefix + envSuffix
	TopicPumpAmmRemoveLiquidity = PumpAmmRemoveLiquidityPrefix + envSuffix

	// 所有需要监听的 topics
	AllTopics = []string{
		TopicRayCreate,
		TopicRaySwap,
		TopicRayAddLiquidity,
		TopicRayRemoveLiquidity,
		// TopicPumpCreate,
		// TopicPumpTrade,
		// TopicPumpComplete,
		// TopicPumpSetParams,
		TopicUnknownToken,
		TopicGameOutTrade,
		TopicGameInTrade,
		TopicPointTxStatus,
		TopicPumpAmmNewPool,
		TopicPumpAmmSwap,
		TopicPumpAmmAddLiquidity,
		TopicPumpAmmRemoveLiquidity,
	}
)
