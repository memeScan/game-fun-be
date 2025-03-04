package service

import (
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
)

type PointServiceImpl struct {
}

func NewPointService(txService *TransactionCkServiceImpl) *PointServiceImpl {
	return &PointServiceImpl{}
}

func (s *PointServiceImpl) CalculateVolumeStatistics() error {
	// 实现积分计算的具体逻辑
	// 例如：
	// 1. 获取需要计算积分的交易记录
	// 2. 计算积分
	// 3. 更新用户积分

	service := TransactionCkServiceImpl{}
	transactions, err := service.GetTokenTransactions("tokenAddress", 1)
	if err != nil {
		return fmt.Errorf("获取交易记录失败: %w", err)
	}

	pool, _ := ants.NewPool(8)
	defer pool.Release()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalVolume decimal.Decimal

	for _, tx := range transactions {
		wg.Add(1)

		err := pool.Submit(func() {
			defer wg.Done()

			// 计算单个交易的交易量
			volume := tx.TransactionAmountUSD

			// 累加到总交易量
			mu.Lock()
			totalVolume = totalVolume.Add(volume)
			mu.Unlock()
		})

		if err != nil {
			return fmt.Errorf("提交交易量计算任务失败: %w", err)
		}
	}

	// 等待所有计算任务完成
	wg.Wait()

	// 返回计算结果
	return nil
}
