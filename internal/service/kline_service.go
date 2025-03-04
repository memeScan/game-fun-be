// internal/service/kline.go
package service

import (
	"game-fun-be/internal/clickhouse"
	"time"
)

type KlineService struct{}

func NewKlineService() *KlineService {
	return &KlineService{}
}

// GetTokenKlines 获取K线数据
func (s *KlineService) GetTokenKlines(tokenAddress string, interval string, start, end time.Time) ([]clickhouse.Kline, error) {
	return clickhouse.GetKlines(tokenAddress, interval, start, end)
}

// GetLatestKline 获取最新的K线数据
func (s *KlineService) GetLatestKline(tokenAddress string) (*clickhouse.Kline, error) {
	return clickhouse.GetLatestKline(tokenAddress)
}
