package kafka

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"my-token-ai-be/internal/pkg/util"

	"github.com/IBM/sarama"
)

// MessageHandler 是处理消息的函数类型
type MessageHandler func([]byte, string) error

// BatchMessageHandler 是批量处理消息的函数类型，增加 topic、partition 和 goroutineID 参数
type BatchMessageHandler func(topic string, messages []sarama.ConsumerMessage, partition int32, goroutineID uint64) error

// TopicConsumer 管理多个 topic 的消费
type TopicConsumer struct {
	consumerGroup       sarama.ConsumerGroup
	stopChan            chan struct{}
	session             sarama.ConsumerGroupSession
	sessionMutex        sync.RWMutex
	batchSizes          map[string]int
	defaultBatchSize    int
	handlers            map[string]MessageHandler
	batchHandlers       map[string]BatchMessageHandler
	batchTimeouts       map[string]time.Duration
	defaultBatchTimeout time.Duration
	groupID             string
	processedOffsets    map[string]map[int32]int64 // topic -> partition -> last processed offset
	processMutex        sync.RWMutex
	minBatchSizes       map[string]int // 新增：每个topic的最小批量
	defaultMinBatchSize int            // 新增：默认最小批量
}

// NewTopicConsumer 创建一个新的 TopicConsumer
func NewTopicConsumer(groupID string) (*TopicConsumer, error) {
	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(os.Getenv("KAFKA_BROKERS"), ","), groupID, KafkaConfig)
	if err != nil {
		return nil, err
	}
	return &TopicConsumer{
		consumerGroup:       consumerGroup,
		groupID:             groupID,
		batchHandlers:       make(map[string]BatchMessageHandler),
		stopChan:            make(chan struct{}),
		handlers:            make(map[string]MessageHandler),
		batchTimeouts:       make(map[string]time.Duration),
		batchSizes:          make(map[string]int),
		defaultBatchSize:    100,
		defaultBatchTimeout: 5 * time.Second,
		processedOffsets:    make(map[string]map[int32]int64),
		minBatchSizes:       make(map[string]int),
		defaultMinBatchSize: 50,
	}, nil
}

// Setup 在新会话开始时运行
func (tc *TopicConsumer) Setup(session sarama.ConsumerGroupSession) error {
	tc.sessionMutex.Lock()
	tc.session = session
	tc.sessionMutex.Unlock()
	return nil
}

// AddHandler 为指定的 topic 添加处理函数
func (tc *TopicConsumer) AddHandler(topic string, handler MessageHandler) {
	if handler != nil {
		tc.handlers[topic] = handler
	}
}

// AddHandler 为指定的 topic 添加批处理函数
func (tc *TopicConsumer) AddBatchHandler(topic string, batchHandler BatchMessageHandler, maxBatchSize int, minBatchSize int, batchTimeout time.Duration) {
	if batchHandler != nil {
		tc.batchHandlers[topic] = batchHandler
	}
	if maxBatchSize > 0 {
		tc.batchSizes[topic] = maxBatchSize
	} else {
		tc.batchSizes[topic] = tc.defaultBatchSize
	}
	if minBatchSize > 0 {
		tc.minBatchSizes[topic] = minBatchSize
	} else {
		tc.minBatchSizes[topic] = tc.defaultMinBatchSize
	}
	if batchTimeout > 0 {
		tc.batchTimeouts[topic] = batchTimeout
	} else {
		tc.batchTimeouts[topic] = tc.defaultBatchTimeout
	}
}

// ConsumeTopics 开始消费指定的 topics
func (tc *TopicConsumer) ConsumeTopics(topics []string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-tc.stopChan:
			return nil
		default:
			err := tc.consumerGroup.Consume(ctx, topics, tc)
			if err != nil {
				util.Log().Error("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}
	}
}

// ConsumeClaim 消费消息的主循环
func (tc *TopicConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 修改初始日志格式
	util.Log().Info("=== Consumer Started ===\n"+
		"Topic:          %s\n"+
		"Partition:      %d\n"+
		"InitialOffset:  %d\n"+
		"Goroutine:      %d",
		claim.Topic(), claim.Partition(), claim.InitialOffset(), util.GetGoroutineID())

	batchMessages := make(map[string][]sarama.ConsumerMessage)
	timers := make(map[string]*time.Timer)
	messageCounter := 0 // 添加消息计数器

	// 定期打印状态的定时器
	statusTicker := time.NewTicker(20 * time.Second)
	defer statusTicker.Stop()

	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				util.Log().Info("=== Consumer Channel Closed ===\n"+
					"Topic:          %s\n"+
					"Partition:      %d\n"+
					"TotalProcessed: %d\n"+
					"Goroutine:      %d",
					claim.Topic(), claim.Partition(), messageCounter, util.GetGoroutineID())
				return nil
			}

			messageCounter++

			handler, handlerExists := tc.handlers[msg.Topic]
			batchHandler, batchHandlerExists := tc.batchHandlers[msg.Topic]

			if handlerExists && handler != nil {
				util.Log().Info("=== Immediate Processing ===\n"+
					"Topic:     %s\n"+
					"Partition: %d\n"+
					"Offset:    %d\n"+
					"Goroutine: %d",
					msg.Topic, msg.Partition, msg.Offset, util.GetGoroutineID())
				if err := handler(msg.Value, msg.Topic); err != nil {
					util.Log().Error("=== Message Handling Error ===\n"+
						"Topic:     %s\n"+
						"Partition: %d\n"+
						"Error:     %v",
						msg.Topic, msg.Partition, err)
				} else {
					session.MarkMessage(msg, "")
					session.Commit()
				}
			} else if batchHandlerExists && batchHandler != nil {
				// 批处理模式
				if _, exists := batchMessages[msg.Topic]; !exists {
					batchSize := tc.batchSizes[msg.Topic]
					batchMessages[msg.Topic] = make([]sarama.ConsumerMessage, 0, batchSize)
					timers[msg.Topic] = time.NewTimer(tc.batchTimeouts[msg.Topic])
				}

				// 如果当前批次超过限制，先处理它
				if len(batchMessages[msg.Topic]) >= tc.batchSizes[msg.Topic] {
					util.Log().Info("Processing batch: size=%d limit=%d",
						len(batchMessages[msg.Topic]), tc.batchSizes[msg.Topic])

					if err := tc.processBatchForTopic(msg.Topic, batchMessages[msg.Topic], session, claim.Partition(), util.GetGoroutineID()); err != nil {
						util.Log().Error("Failed to process batch, will retry: %v", err)
						continue // 保持失败时的重试逻辑
					}

					// 只有处理成功才清空批次
					batchMessages[msg.Topic] = batchMessages[msg.Topic][:0]
					timers[msg.Topic].Reset(tc.batchTimeouts[msg.Topic])
				}

				// 只有在批次未满时才添加新消息
				if len(batchMessages[msg.Topic]) < tc.batchSizes[msg.Topic] {
					batchMessages[msg.Topic] = append(batchMessages[msg.Topic], *msg)
				}
			}

		case <-statusTicker.C:
			currentOffset := claim.HighWaterMarkOffset()
			lastProcessedOffset := tc.getLastProcessedOffset(claim.Topic(), claim.Partition())
			realLag := currentOffset - lastProcessedOffset

			// 添加处理进度日志
			util.Log().Info("=== Processing Progress ===\n"+
				"Topic:              %s\n"+
				"Partition:          %d\n"+
				"Last Processed:     %d\n"+
				"Current HW:         %d\n"+
				"Real Lag:          %d\n"+
				"Batch Size:         %d\n"+
				"Messages in Batch:  %d",
				claim.Topic(),
				claim.Partition(),
				lastProcessedOffset,
				currentOffset,
				realLag,
				tc.batchSizes[claim.Topic()],
				len(batchMessages[claim.Topic()]))

			// 打印每个主题的批次状态并检查超时
			for topic, messages := range batchMessages {
				util.Log().Info("Batch Details:\n"+
					"Topic:              %s\n"+
					"Messages in Batch:  %d\n"+
					"Batch Size Limit:   %d\n"+
					"Batch Timeout:      %s\n"+
					"Has BatchHandler:   %v",

					topic,
					len(messages),
					tc.batchSizes[topic],
					tc.batchTimeouts[topic],
					tc.batchHandlers[topic] != nil)

				// 检查超时批次
				if len(messages) > 0 && len(messages) < tc.batchSizes[topic] {
					select {
					case <-timers[topic].C:
						if len(messages) >= tc.minBatchSizes[topic] {
							util.Log().Info("=== Processing Timeout Batch ===\n"+
								"Topic:     %s\n"+
								"Partition: %d\n"+
								"BatchSize: %d",
								topic, claim.Partition(), len(messages))

							if err := tc.processBatchForTopic(topic, messages, session, claim.Partition(), util.GetGoroutineID()); err != nil {
								util.Log().Error("Failed to process timeout batch, will retry: %v", err)
								continue
							}

							batchMessages[topic] = batchMessages[topic][:0]
							timers[topic] = time.NewTimer(tc.batchTimeouts[topic])
						} else {
							// 批次太小，继续等待
							util.Log().Info("Batch too small on timeout: topic=%s size=%d min_size=%d",
								topic, len(messages), tc.minBatchSizes[topic])
							timers[topic] = time.NewTimer(tc.batchTimeouts[topic])
						}
					default:
					}
				}
			}

		case <-session.Context().Done():
			// 处理剩余的批次消息
			for topic, messages := range batchMessages {
				if len(messages) > 0 {
					// 尝试处理剩余消息，但不重试
					if err := tc.processBatchForTopic(topic, messages, session, claim.Partition(), util.GetGoroutineID()); err != nil {
						util.Log().Error("Failed to process remaining messages before shutdown: %v", err)
					}
				}
			}
			util.Log().Info("=== Session Completed ===\n"+
				"Topic:          %s\n"+
				"Partition:      %d\n"+
				"TotalProcessed: %d",
				claim.Topic(), claim.Partition(), messageCounter)
			return nil
		}
	}
}

// 处理单个主题的批次消息
func (tc *TopicConsumer) processBatchForTopic(topic string, messages []sarama.ConsumerMessage, session sarama.ConsumerGroupSession, partition int32, goroutineID uint64) error {
	// 检查最小批量
	minSize := tc.minBatchSizes[topic]
	if len(messages) < minSize {
		util.Log().Info("Batch too small, skipping processing: topic=%s size=%d min_size=%d",
			topic, len(messages), minSize)
		return nil
	}

	startTime := time.Now()
	firstOffset := messages[0].Offset
	lastOffset := messages[len(messages)-1].Offset

	util.Log().Info("Starting batch process: topic=%s partition=%d size=%d offset_range=%d-%d",
		topic, partition, len(messages), firstOffset, lastOffset)

	batchHandler := tc.batchHandlers[topic]
	if batchHandler != nil {
		maxRetries := 3
		for retry := 0; retry < maxRetries; retry++ {
			err := batchHandler(topic, messages, partition, goroutineID)
			if err != nil {
				util.Log().Error("=== Batch Processing Error (Attempt %d/%d) ===\n"+
					"Topic:     %s\n"+
					"Partition: %d\n"+
					"Goroutine: %d\n"+
					"Error:     %v",
					retry+1, maxRetries, topic, partition, goroutineID, err)

				if retry == maxRetries-1 {
					util.Log().Error("=== Max Retries Reached ===")
					return err // 不提交消息，让消费者重试
				}

				time.Sleep(time.Second * time.Duration(retry+1))
				continue
			}

			// 处理成功后，逐条更新进度并标记
			for _, msg := range messages {
				session.MarkMessage(&msg, "")
				tc.recordProgress(topic, partition, msg.Offset)
			}
			session.Commit()

			util.Log().Info("Completed batch process: topic=%s partition=%d duration=%v messages=%d",
				topic, partition, time.Since(startTime), len(messages))
			return nil
		}
	}
	return nil
}

// Cleanup 在会话结束时运行
func (tc *TopicConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	tc.sessionMutex.Lock()
	tc.session = nil
	tc.sessionMutex.Unlock()
	return nil
}

// 关闭消费者
func (tc *TopicConsumer) Close() error {
	close(tc.stopChan) // 停止 goroutine
	return tc.consumerGroup.Close()
}

// 在处理消息时记录进度
func (tc *TopicConsumer) recordProgress(topic string, partition int32, offset int64) {
	tc.processMutex.Lock()
	defer tc.processMutex.Unlock()

	if tc.processedOffsets[topic] == nil {
		tc.processedOffsets[topic] = make(map[int32]int64)
	}
	tc.processedOffsets[topic][partition] = offset
}

// 修改状态检查逻辑
func (tc *TopicConsumer) getLastProcessedOffset(topic string, partition int32) int64 {
	tc.processMutex.RLock()
	defer tc.processMutex.RUnlock()

	if tc.processedOffsets[topic] == nil {
		return 0
	}
	return tc.processedOffsets[topic][partition]
}
