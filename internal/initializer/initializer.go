package initializer

import (
	"my-token-ai-be/internal/clickhouse"
	"my-token-ai-be/internal/conf"
	"my-token-ai-be/internal/cron"
	"my-token-ai-be/internal/es"
	"my-token-ai-be/internal/kafka"
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/pkg/httpUtil"
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/redis"
	"os"
)

// Setup 初始化所有组件
func Setup(env string) {
	// 设置环境变量
	conf.SetEnv(env)

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// 初始化数据库、redis、elasticsearch
	model.Database(os.Getenv("MYSQL_DSN"))
	redis.Redis()
	es.Elasticsearch()
	clickhouse.ClickHouse()

	endpoint := os.Getenv("BLOCKCHAIN_API_ENDPOINT")
	httpUtil.InitAPI(&endpoint)
	httpUtil.InitMetrics(redis.RedisClient)

	// 如果不是 debug 环境
	if !conf.IsDebug() {
		// 通过环境变量控制:
		// kafka=1: 启动 Kafka 消费和定时任务
		// kafka=0 或未设置: 不启动任何服务
		if os.Getenv("kafka") == "1" {
			// 在消费 Kafka 时启动定时任务
			cron.InitCronJobs()
			// 启动 Kafka 消费
			kafka.Kafka()
			go func() {
				if err := kafka.ConsumePumpfunTopics(); err != nil {
					util.Log().Error("Failed to consume Kafka topics: %v", err)
				}
			}()
		}
	}
}
