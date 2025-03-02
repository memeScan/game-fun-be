// internal/service/kline.go
package service

import (
	"my-token-ai-be/internal/clickhouse"
	"my-token-ai-be/internal/model"
	"time"
)

type KlineService struct{}

func NewKlineService() *KlineService {
	return &KlineService{}
}

// ProcessTransaction 处理交易数据
func (s *KlineService) ProcessTransaction(tx *clickhouse.TokenTransactionCk) error {
	return clickhouse.InsertTransaction(tx)
}

// BatchProcessTransactions 批量处理交易数据
func (s *KlineService) BatchProcessTransactions(txs []*clickhouse.TokenTransactionCk) error {
	return clickhouse.BatchInsertTransactions(txs)
}

// GetTokenKlines 获取K线数据
func (s *KlineService) GetTokenKlines(tokenAddress string, interval string, start, end time.Time) ([]clickhouse.Kline, error) {
	return clickhouse.GetKlines(tokenAddress, interval, start, end)
}

// GetLatestKline 获取最新的K线数据
func (s *KlineService) GetLatestKline(tokenAddress string) (*clickhouse.Kline, error) {
	return clickhouse.GetLatestKline(tokenAddress)
}

// ConvertToTransactionCks 将 transactions 转换为 TokenTransactionCk 列表
func (s *KlineService) ConvertToTransactionCks(transactions []*model.TokenTransaction) []*clickhouse.TokenTransactionCk {
	txCks := make([]*clickhouse.TokenTransactionCk, 0, len(transactions))

	for _, tx := range transactions {
		txCk := &clickhouse.TokenTransactionCk{}

		// 设置基本信息
		txCk.TransactionHash = tx.TransactionHash
		txCk.ChainType = tx.ChainType
		txCk.UserAddress = tx.UserAddress
		txCk.TokenAddress = tx.TokenAddress
		txCk.PoolAddress = tx.PoolAddress

		// 设置数量
		baseAmount := tx.NativeTokenAmount
		quoteAmount := tx.TokenAmount
		txCk.BaseTokenAmount = baseAmount
		txCk.QuoteTokenAmount = quoteAmount
		txCk.Decimals = tx.Decimals

		// 设置价格
		txCk.BaseTokenPrice = tx.NativePriceUSD
		txCk.QuoteTokenPrice = tx.Price
		txCk.TransactionAmountUSD = tx.TransactionAmountUSD

		// 设置其他信息
		txCk.TransactionType = tx.TransactionType
		txCk.PlatformType = tx.PlatformType
		txCk.TransactionTime = tx.TransactionTime

		txCks = append(txCks, txCk)
	}

	return txCks
}
