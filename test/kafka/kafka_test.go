package kafka_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"game-fun-be/internal/kafka"
	"game-fun-be/internal/model"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
)

func TestKafkaInitialization(t *testing.T) {
	// 设置测试环境
	t.Run("Initialize Kafka Producer", func(t *testing.T) {
		// 设置测试环境变量
		os.Setenv("KAFKA_BROKERS", "alikafka-post-public-intl-sg-jiy45rtfa0s-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-public-intl-sg-jiy45rtfa0s-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-public-intl-sg-jiy45rtfa0s-3-vpc.alikafka.aliyuncs.com:9092")

		// 初始化 Kafka
		kafka.Kafka()
		defer kafka.Close()

		// 验证 Producer 是否成功创建
		if kafka.KafkaProducer == nil {
			t.Error("KafkaProducer should not be nil")
		}
	})
}

func TestKafkaSendMessage(t *testing.T) {
	// 创建 mock producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	mock := mocks.NewSyncProducer(t, config)

	// 保存原始 producer
	originalProducer := kafka.KafkaProducer
	// 替换为 mock producer
	kafka.KafkaProducer = mock

	// 测试结束后恢复原始 producer
	defer func() {
		kafka.KafkaProducer = originalProducer
	}()

	t.Run("Send Message Successfully", func(t *testing.T) {
		// 设置预期
		mock.ExpectSendMessageAndSucceed()

		// 测试发送消息
		topic := "test-topic"
		message := []byte("test message")
		err := kafka.SendMessage(topic, message)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Send Message With Error", func(t *testing.T) {
		// 设置预期失败
		mock.ExpectSendMessageAndFail(sarama.ErrTopicAuthorizationFailed)

		// 测试发送消息
		topic := "test-topic"
		message := []byte("test message")
		err := kafka.SendMessage(topic, message)

		if err != sarama.ErrTopicAuthorizationFailed {
			t.Errorf("Expected error %v, got %v", sarama.ErrTopicAuthorizationFailed, err)
		}
	})
}

func TestKafkaIntegration(t *testing.T) {
	// 跳过集成测试，除非明确指定要运行
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Real Kafka Integration", func(t *testing.T) {
		// 设置测试环境变量
		// os.Setenv("KAFKA_BROKERS", "alikafka-post-public-intl-sg-jiy45rtfa0s-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-public-intl-sg-jiy45rtfa0s-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-public-intl-sg-jiy45rtfa0s-3-vpc.alikafka.aliyuncs.com:9092")
		os.Setenv("KAFKA_BROKERS", "192.168.31.107:9092,192.168.31.108:9092,192.168.31.109:9092")

		// 初始化 Kafka
		kafka.Kafka()
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic: %v", r)
			}
			kafka.Close() // 确保 Kafka 连接关闭
		}()

		// 准备批次消息
		batchSize := 1 // 批次大小
		messages := make([]*sarama.ProducerMessage, 0, batchSize)
		for i := 0; i < batchSize; i++ {
			// message := &sarama.ProducerMessage{
			// 	Topic: kafka.TopicGameOutTrade,
			// 	Value: sarama.ByteEncoder(fmt.Sprintf(`{
			// 		"timestamp": %d,
			// 		"block": 3130978492,
			// 		"signature": "2ZZf2p6RH4Zr4nKbGV9zmNkGjr3wEgpER9xxJD3WkfqVkCHzoW2HCZtg4JtKwqPmKBB5o6Um8cDjqxksVrBg8UUU",
			// 		"user": "SoLxyz987654321abc987654321abc987654321",
			// 		"poolAddress": "F3UWHvZWy41HpsbWfbBFbvAemYuCn8RJ2qvKfr2346Qf",
			// 		"isBuy": true,
			// 		"quoteToken": "8iFREvVdmLKxVeibpC5VLRr1S6X5dm7gYR3VCU1wpump",
			// 		"baseToken": "So11111111111111111111111111111111111111112",
			// 		"marketAddress": "",
			// 		"poolQuoteReserve": 18640745631097,
			// 		"poolBaseReserve": 539997130105,
			// 		"quoteAmount": 54610438,
			// 		"baseAmount": 1585960,
			// 		"decimals": 6,
			// 		"feeQuoteAmount": 546104,
			// 		"feeBaseAmount": 15859,
			// 		"buybackFeeBaseAmount": 7930,
			// 		"isBurn": false
			// 	}`, time.Now().Unix())),
			// }

			message := &sarama.ProducerMessage{
				Topic: kafka.TopicGameInTrade,
				Value: sarama.ByteEncoder(fmt.Sprintf(`{
            "timestamp": %d,
            "block": 3130978492,
            "signature": "2ZZf2p6RH4Zr4nKbGV9zmNkGjr3wEgpER9xxJD3WkfqVkCHzoW2HCZtg4JtKwqPmKBB5o6Um8cDjqxksVrBg8UUU",
            "user": "SoLxyz987654321abc987654321abc987654321",
            "isBuy": true,
            "quoteToken": "8iFREvVdmLKxVeibpC5VLRr1S6X5dm7gYR3VCU1wpump",
            "baseToken": "So11111111111111111111111111111111111111112",
            "quoteAmount": 54610438,
            "baseAmount": 1585960,
            "decimals": 6,
            "pointsAmount": 1000000,
            "feeBaseAmount": 15859
        }`, time.Now().Unix())),
			}
			messages = append(messages, message)
		}

		// 批次发送
		err := kafka.GetProducer().SendMessages(messages)
		if err != nil {
			t.Errorf("Failed to send messages: %v", err)
		} else {
			t.Logf("Successfully sent %d messages", len(messages))
		}
	})
}

func TestSolBalanceAPI(t *testing.T) {
	// 定义要请求的 URL
	url := "http://172.20.8.16:3001/api/v1/sol-balance?address=Huy4cz1yTxS6GrGMN7Q5acQ7ws3PsHsc886i4iS2pump"

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// 输出响应内容
	t.Logf("Response Body: %s", string(body))
}

func TestRaydiumCreateMessage(t *testing.T) {
	// 跳过集成测试，除非明确指定要运行
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Send Raydium Create Message", func(t *testing.T) {
		// 设置测试环境变量
		os.Setenv("KAFKA_BROKERS", "alikafka-post-public-intl-sg-jiy45rtfa0s-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-public-intl-sg-jiy45rtfa0s-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-public-intl-sg-jiy45rtfa0s-3-vpc.alikafka.aliyuncs.com:9092")

		// 初始化 Kafka
		kafka.Kafka()
		defer kafka.Close()

		// 测试发送实际消息
		topic := kafka.TopicRayCreate
		message := []byte(`{
			"timeStamp": 1731594912,
			"signature": "2c8FviquiAbHSh4rDK2jgJanYzwUEFchCvBBdG7nEC5bkLYUQ2AFhQ7wkPRYQoD3WAxvqpyehWWaUt5YtxAjyXaQ",
			"marketAddress": "MarketAddressExample1234567890",
			"poolAddress": "8nsjiwgZGpqMQ4n3fSWcEdMoQfMaAqxBFTkaGDtzeD4J",
			"poolState": 1,
			"pcAddress": "CUCfqECyNKLe8dzxBQtT8vzHXpDryT72c8NtEeoN7WmS",
			"coinAddress": "So11111111111111111111111111111111111111112",
			"changePoolPcAmount": "50000000",
			"changePoolCoinAmount": "198844097903",
			"poolPcReserve": "150000000",
			"poolCoinReserve": "796532293710",
			"user": "EYANY4XNWRcx3YBhFygQLo3UAzGnXEWBskZMctyuxyFG",
			"block": 301364918
		}`)

		err := kafka.SendMessage(topic, message)
		if err != nil {
			t.Errorf("Failed to send message: %v", err)
		}
		// 给消息发送一些时间
		time.Sleep(time.Second)
	})
}

func TestKafkaLoopIntegration(t *testing.T) {
	// 跳过集成测试，除非明确指定要运行
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	loopCount := 100 // 设置循环次数

	for i := 0; i < loopCount; i++ {
		t.Run(fmt.Sprintf("Integration Test Loop %d", i+1), func(t *testing.T) {
			TestKafkaIntegration(t)
		})
	}
}

func TestPointTxStatusHandler(t *testing.T) {
	// 跳过集成测试，除非明确指定要运行
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 设置环境变量，配置Kafka连接
	os.Setenv("KAFKA_BROKERS", "alikafka-post-public-intl-sg-jiy45rtfa0s-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-public-intl-sg-jiy45rtfa0s-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-public-intl-sg-jiy45rtfa0s-3-vpc.alikafka.aliyuncs.com:9092")

	// 初始化 Kafka
	kafka.Kafka()
	defer kafka.Close()

	// 创建测试用的消息
	testSignature := "2NbG9aAP8DdTejHSjrDDNDAAdgJ5fcnwwGfAA12W7UMUmqgLqoBSRX2ukD62KA91L1xb4L7jgD9uxAF9X5ijtWpJ"
	testUserId := uint(3)
	testPoints := uint64(1000)

	statusMsg := model.PointTxStatusMessage{
		Signature: testSignature,
		UserId:    testUserId,
		Points:    testPoints,
	}

	// 序列化消息
	msgBytes, err := json.Marshal(statusMsg)
	if err != nil {
		t.Fatalf("Failed to marshal test message: %v", err)
	}

	// 使用Kafka发送消息
	topic := kafka.TopicPointTxStatus
	t.Logf("Sending test message to topic %s: %s", topic, string(msgBytes))

	err = kafka.SendMessage(topic, msgBytes)
	if err != nil {
		t.Errorf("Failed to send message: %v", err)
		return
	}

	// 等待消息处理完成
	t.Log("Message sent, waiting for processing...")
	time.Sleep(5 * time.Second)

	// 注意：在实际测试中，你应该验证用户积分是否按预期更新
	// 例如，如果模拟交易失败的情况，应该检查用户积分是否被恢复
	t.Log("Test completed - check logs for processing results")

}
