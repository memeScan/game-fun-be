package service

import (
	"fmt"
	"game-fun-be/internal/clickhouse"
	"game-fun-be/internal/model"
	"game-fun-be/internal/response"

	"github.com/shopspring/decimal"
)

type TransactionCkServiceImpl struct {
}

func NewTransactionCkService() *TransactionCkServiceImpl {
	return &TransactionCkServiceImpl{}
}

func (service *TransactionCkServiceImpl) GetTokenOrderBook(tokenAddress string, chainType uint8) response.Response {
	transactions, err := clickhouse.GetTokenTransactions(tokenAddress, 100)
	if err != nil {
		return response.Err(response.CodeDBError, "failed to get token order book", err)
	}

	convertedTransactions := make([]response.TokenOrderBookItem, len(transactions))
	for i, tx := range transactions {
		// 计算基础代币数量
		baseAmount := decimal.NewFromInt(int64(tx.BaseTokenAmount))
		if tx.Decimals > 0 {
			baseAmount = baseAmount.Shift(-9)
		}

		// 计算报价代币数量
		quoteAmount := decimal.NewFromInt(int64(tx.QuoteTokenAmount))
		if tx.Decimals > 0 {
			quoteAmount = quoteAmount.Shift(-int32(tx.Decimals))
		}

		// 计算 USD 金额
		var usdAmount decimal.Decimal
		if tx.TransactionType == 1 { // 1 是买入
			usdAmount = quoteAmount.Mul(tx.QuoteTokenPrice)
		} else {
			usdAmount = baseAmount.Mul(tx.BaseTokenPrice)
		}

		// 创建一个新的 TokenOrderBookItem
		item := response.TokenOrderBookItem{}

		// 逐个赋值
		item.TransactionHash = tx.TransactionHash
		item.ChainType = tx.ChainType
		item.UserAddress = tx.UserAddress
		item.TokenAddress = tx.TokenAddress
		item.PoolAddress = tx.PoolAddress
		item.BaseTokenAmount = baseAmount
		item.QuoteTokenAmount = quoteAmount
		item.BaseTokenPrice = tx.BaseTokenPrice
		item.QuoteTokenPrice = tx.QuoteTokenPrice
		item.TransactionType = tx.TransactionType
		item.PlatformType = tx.PlatformType
		item.TransactionTime = tx.TransactionTime.Unix()
		item.UsdAmount = usdAmount

		// 将item赋值给切片
		convertedTransactions[i] = item
	}

	return response.BuildTokenOrderBookResponse(convertedTransactions)
}

func (s *TransactionCkServiceImpl) GetTokenTransactions(tokenAddress string, limit int) ([]clickhouse.TokenTransactionCk, error) {
	transactions, err := clickhouse.GetTokenTransactions(tokenAddress, limit)
	if err != nil {
		return nil, fmt.Errorf("获取交易记录失败: %w", err)
	}
	return transactions, nil
}

// ConvertToTransactionCks 将 transactions 转换为 TokenTransactionCk 列表
func (s *TransactionCkServiceImpl) ConvertToTransactionCks(transactions []*model.TokenTransaction) []*clickhouse.TokenTransactionCk {
	txCks := make([]*clickhouse.TokenTransactionCk, 0, len(transactions))

	for _, tx := range transactions {
		txCk := &clickhouse.TokenTransactionCk{}

		// 设置基本信息
		txCk.TransactionID = tx.ID
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

// BatchProcessTransactions 批量处理交易数据
func (s *TransactionCkServiceImpl) BatchProcessTransactions(txs []*clickhouse.TokenTransactionCk) error {
	return clickhouse.BatchInsertTransactions(txs)
}
