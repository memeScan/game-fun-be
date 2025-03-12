package kafka

import (
	"log"
	"os"
	"strings"
	"time"

	"game-fun-be/internal/pkg/util"

	"github.com/IBM/sarama"
)

var (
	KafkaProducer sarama.SyncProducer
	KafkaConfig   *sarama.Config
)

func Kafka() sarama.SyncProducer {
	// 从环境变量获取 Kafka 配置
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	// 使用优化后的配置替换默认配置
	KafkaConfig = optimizeKafkaConfig()

	// 初始化同步生产者
	producer, err := sarama.NewSyncProducer(brokers, KafkaConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %s", err)
	}
	KafkaProducer = producer

	// 获取 Kafka 版本信息
	version := KafkaConfig.Version
	util.Log().Info("Kafka version: %v", version)
	util.Log().Info("Kafka connected successfully to: %v", brokers)
	return producer
}

// SendMessage 发送消息到指定的 topic
func SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := KafkaProducer.SendMessage(msg)
	return err
}

// Close 关闭 Kafka 连接
func Close() {
	if KafkaProducer != nil {
		KafkaProducer.Close()
	}
}

func optimizeKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()

	// 从环境变量获取 Client ID，如果没有设置则使用默认值
	clientID := os.Getenv("KAFKA_CLIENT_ID")
	if clientID == "" {
		clientID = "my-token-ai-consumer-1"
	}
	config.ClientID = clientID

	// 添加同步生产者必需的配置
	config.Producer.Return.Successes = true          // 必须设置为 true 才能用于同步生产者
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	// 生产者超时设置
	config.Producer.Timeout = 10 * time.Second             // 生产者超时
	config.Producer.Retry.Max = 3                          // 最大重试次数
	config.Producer.Retry.Backoff = 100 * time.Millisecond // 重试间隔

	// 消费者基础配置
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = false

	// Channel 配置
	// config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange() // 确保多个消费者分摊分区

	config.ChannelBufferSize = 1024 // 从 8192 降到 1024
	// 12个分区 * 1024 = 最多可以缓存12288条消息
	// 小缓冲区可以更好地处理流量峰值和消息积压

	// Fetch 配置 - 针对15批消息优化（考虑12个分区）
	config.Consumer.Fetch.Min = 32 * 1024       // 降到 32KB
	config.Consumer.Fetch.Default = 1024 * 1024 // 降到 1MB
	config.Consumer.Fetch.Max = 2 * 1024 * 1024 // 降到 2MB

	// 时间配置
	config.Consumer.MaxWaitTime = 500 * time.Millisecond
	// 当没有足够消息时，等待积累新消息的最大时间
	// 降低到250ms可以提高实时性，因为我们已经设置了较大的Fetch大小(1.5MB)
	// 不需要等太久来积累消息

	config.Consumer.MaxProcessingTime = 5 * time.Second // 从 10s 降到 5s
	// 处理一批消息的最大允许时间，超过这个时间会触发rebalance
	// 增加到20s，因为：
	// - 12个分区同时处理可能需要更多时间
	// - 每个fetch最多1.5MB数据(15批消息)需要足够处理时间
	// - 避免因临时性的处理延迟导致不必要的rebalance

	// 会话配置
	config.Consumer.Group.Session.Timeout = 20 * time.Second // 从 30s 降到 20s
	// 消费者组成员被认为死亡前的最大时间
	// 45s给予足够的时间处理消息和网络波动
	// 如果超过这个时间没有心跳，会触发rebalance

	config.Consumer.Group.Heartbeat.Interval = 6 * time.Second // 从 10s 降到 6s
	// 向coordinator发送心跳的间隔时间
	// 设置为session timeout的1/3是最佳实践
	// 确保有足够的重试机会，避免误判为死亡

	config.Consumer.Group.Rebalance.Timeout = 30 * time.Second // 从 60s 降到 30s
	// rebalance过程的最大允许时间
	// 60s足够12个分区完成重新分配
	// 特别是在处理大量数据时，需要足够时间完成收尾工作

	// 性能优化
	config.Net.MaxOpenRequests = 15 // 从 30 降到 15
	// 限制每个broker连接的最大并发请求数
	// 设置为15是因为：
	// - 12个分区需要同时发送fetch请求
	// - 额外预留3个请求用于其他操作（如心跳、提交offset等）
	// - 确保每个分区都能及时获取数据，不会因为请求限制而等待

	config.Net.KeepAlive = 60 * time.Second
	// TCP keepalive 时间
	// 增加到60s以减少连接重建的频率
	// 特别是在网络稳定的环境中，可以维持更长的连接时间

	// 网络相关配置调整
	config.Net.DialTimeout = 30 * time.Second  // 连接超时时间
	config.Net.ReadTimeout = 30 * time.Second  // 读取超时时间
	config.Net.WriteTimeout = 30 * time.Second // 写入超时时间

	// 5. 添加流控制
	config.Net.SASL.Enable = false // 如果不需要认证，禁用 SASL
	config.Net.TLS.Enable = false  // 如果不需要 TLS，��用它

	// 重试策略
	config.Metadata.Retry.Max = 3                   // 元数据重试次数
	config.Metadata.Retry.Backoff = 5 * time.Second // 重试间隔

	return config
}

// GetProducer returns the global Kafka producer instance
func GetProducer() sarama.SyncProducer {
	return KafkaProducer
}
