package kafka

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

	// 所有需要监听的 topics
	AllTopics = []string{
		TopicRayCreate,
		TopicRaySwap,
		TopicRayAddLiquidity,
		TopicRayRemoveLiquidity,
		TopicPumpCreate,
		TopicPumpTrade,
		TopicPumpComplete,
		TopicPumpSetParams,
		TopicUnknownToken,
	}
)
