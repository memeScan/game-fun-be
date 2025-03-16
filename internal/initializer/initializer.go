package initializer

import (
	"game-fun-be/internal/clickhouse"
	"game-fun-be/internal/conf"
	"game-fun-be/internal/cron"
	"game-fun-be/internal/es"
	"game-fun-be/internal/kafka"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"os"

	"github.com/IBM/sarama"
)

// Setup 初始化所有组件
func Setup(env string) sarama.SyncProducer {
	// 设置环境变量
	conf.SetEnv(env)
	conf.InitGameConfig()

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

	// 初始化 Kafka 生产者
	producer := kafka.Kafka()

	// 如果不是 debug 环境
	if !conf.IsDebug() {
		// 通过环境变量控制:
		// kafka=1: 启动 Kafka 消费和定时任务
		// kafka=0 或未设置: 不启动任何服务
		if os.Getenv("kafka") == "1" {
			// 在消费 Kafka 时启动定时任务
			cron.InitCronJobs()
			// 启动 Kafka 消费
			go func() {
				if err := kafka.ConsumePumpfunTopics(); err != nil {
					util.Log().Error("Failed to consume Kafka topics: %v", err)
				}
			}()
		}
	}
	return producer
}
