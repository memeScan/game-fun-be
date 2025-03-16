package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type ProxyTransaction struct {
	TransactionHash         string          `db:"transaction_hash"`
	ChainType               uint8           `db:"chain_type"`
	ProxyType               uint8           `db:"proxy_type"`
	UserAddress             string          `db:"user_address"`
	TokenAddress            string          `db:"token_address"`
	PoolAddress             string          `db:"pool_address"`
	BaseTokenAmount         uint64          `db:"base_token_amount"`
	QuoteTokenAmount        uint64          `db:"quote_token_amount"`
	BaseTokenReserveAmount  uint64          `db:"base_token_reserve_amount"`
	QuoteTokenReserveAmount uint64          `db:"quote_token_reserve_amount"`
	Decimals                uint8           `db:"decimals"`
	BaseTokenPrice          decimal.Decimal `db:"base_token_price"`
	QuoteTokenPrice         decimal.Decimal `db:"quote_token_price"`
	TransactionType         uint8           `db:"transaction_type"`
	IsBurn                  uint8           `db:"is_burn"`
	PointsAmount            uint64          `db:"points_amount"`
	FeeQuoteAmount          uint64          `db:"feeQuote_amount"`
	FeeBaseAmount           uint64          `db:"feeBase_amount"`
	BuybackFeeBaseAmount    uint64          `db:"buybackFeeBase_amount"`
	BlockTime               time.Time       `db:"block_time"`
	TransactionTime         time.Time       `db:"transaction_time"`
	CreateTime              time.Time       `db:"create_time"`
}

type QuoteVolumeAll struct {
	TimeBucket      time.Time `db:"time_bucket"`
	ChainType       uint8     `db:"chain_type"`
	UserAddress     string    `db:"user_address"`
	TokenAddress    string    `db:"token_address"`
	UserQuoteVolume uint64    `db:"user_quote_volume"`
}

// InsertTransaction 插入单条交易数据
func InsertProxyTransaction(tx *ProxyTransaction) error {
	query := `
        INSERT INTO proxy_transaction_ck_all (
            transaction_hash, chain_type, proxy_type, user_address, token_address,
            pool_address, base_token_amount, quote_token_amount, base_token_reserve_amount,
            quote_token_reserve_amount, decimals, base_token_price, quote_token_price,
            transaction_type, is_burn, points_amount, feeQuote_amount, feeBase_amount,
            buybackFeeBase_amount, block_time, transaction_time,create_time
        ) VALUES (
            ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?
        )
    `

	err := ClickHouseClient.Exec(context.Background(), query,
		tx.TransactionHash,
		tx.ChainType,
		tx.ProxyType,
		tx.UserAddress,
		tx.TokenAddress,
		tx.PoolAddress,
		tx.BaseTokenAmount,
		tx.QuoteTokenAmount,
		tx.BaseTokenReserveAmount,
		tx.QuoteTokenReserveAmount,
		tx.Decimals,
		tx.BaseTokenPrice,
		tx.QuoteTokenPrice,
		tx.TransactionType,
		tx.IsBurn,
		tx.PointsAmount,
		tx.FeeQuoteAmount,
		tx.FeeBaseAmount,
		tx.BuybackFeeBaseAmount,
		tx.BlockTime,
		tx.TransactionTime,
		tx.CreateTime,
	)
	if err != nil {
		return fmt.Errorf("insert transaction failed: %w", err)
	}
	return nil
}

// QueryProxyTransactionsByTime 根据时间范围查询交易数据
func QueryProxyTransactionsByTime(startTime, endTime time.Time, chainType uint8, tokenAddress string) ([]ProxyTransaction, error) {
	query := `
        SELECT 
            transaction_hash, chain_type, proxy_type, user_address, token_address,
            pool_address, base_token_amount, quote_token_amount, base_token_reserve_amount,
            quote_token_reserve_amount, decimals, base_token_price, quote_token_price,
            transaction_type, is_burn, points_amount, feeQuote_amount, feeBase_amount,
            buybackFeeBase_amount, block_time, transaction_time, create_time
        FROM proxy_transaction_ck_all
        WHERE transaction_time >= ? AND transaction_time < ?
        AND proxy_type = 1
        AND chain_type = ?
        AND token_address = ?
        ORDER BY transaction_time DESC
    `

	var transactions []ProxyTransaction
	rows, err := ClickHouseClient.Query(context.Background(), query,
		startTime,
		endTime,
		chainType,
		tokenAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("query transactions by time failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tx ProxyTransaction

		if err := rows.Scan(
			&tx.TransactionHash,
			&tx.ChainType,
			&tx.ProxyType,
			&tx.UserAddress,
			&tx.TokenAddress,
			&tx.PoolAddress,
			&tx.BaseTokenAmount,
			&tx.QuoteTokenAmount,
			&tx.BaseTokenReserveAmount,
			&tx.QuoteTokenReserveAmount,
			&tx.Decimals,
			&tx.BaseTokenPrice,
			&tx.QuoteTokenPrice,
			&tx.TransactionType,
			&tx.IsBurn,
			&tx.PointsAmount,
			&tx.FeeQuoteAmount,
			&tx.FeeBaseAmount,
			&tx.BuybackFeeBaseAmount,
			&tx.BlockTime,
			&tx.TransactionTime,
			&tx.CreateTime,
		); err != nil {
			return nil, fmt.Errorf("scan row failed: %w", err)
		}

		transactions = append(transactions, tx)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return transactions, nil
}
