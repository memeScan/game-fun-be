package metrics

import (
	"context"
	"fmt"
	"game-fun-be/internal/pkg/util"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

// MetricsHTTPClient 装饰器结构体
type MetricsHTTPClient struct {
	client      *http.Client
	redisClient *redis.Client
}

// NewMetricsHTTPClient 创建新的装饰器客户端
func NewMetricsHTTPClient(client *http.Client, redisClient *redis.Client) *MetricsHTTPClient {
	return &MetricsHTTPClient{
		client:      client,
		redisClient: redisClient,
	}
}

// Do 实现 http.Client 的 Do 方法
func (m *MetricsHTTPClient) Do(req *http.Request) (*http.Response, error) {
	startTime := time.Now()

	// 执行原始请求
	resp, err := m.client.Do(req)

	// 记录指标 - 使用完整 URL
	m.recordMetrics(req.URL.String(), time.Since(startTime), req.Method)

	return resp, err
}

// Get 实现 GET 方法
func (m *MetricsHTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return m.Do(req)
}

// Post 实现 POST 方法
func (m *MetricsHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return m.Do(req)
}

func (m *MetricsHTTPClient) recordMetrics(fullURL string, duration time.Duration, method string) {
	var baseURL string
	if method == http.MethodGet {
		// 只对 GET 请求解析 URL
		parsedURL, err := url.Parse(fullURL)
		if err != nil {
			util.Log().Error("Failed to parse URL: %v", err)
			return
		}
		baseURL = fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
	} else {
		// POST 等其他请求直接使用原始 URL
		baseURL = fullURL
	}

	pipe := m.redisClient.Pipeline()
	now := time.Now()

	// 使用 baseURL 替代完整 URL
	dailyKey := fmt.Sprintf("api_calls:daily:%s", now.Format("2006-01-02"))
	pipe.ZIncrBy(context.Background(), dailyKey, 1, baseURL)

	pipe.ZIncrBy(context.Background(), "api_calls:total", 1, baseURL)

	latencyKey := fmt.Sprintf("api_latency:%s", now.Format("2006-01-02"))
	pipe.ZIncrBy(context.Background(), latencyKey, float64(duration.Milliseconds()), baseURL)

	if _, err := pipe.Exec(context.Background()); err != nil {
		util.Log().Error("Failed to record API metrics: %v", err)
	}
}
