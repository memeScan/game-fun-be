package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type TokenTransactionCk struct {
	TransactionID        uint64          `db:"transaction_id"`
	TransactionHash      string          `db:"transaction_hash"`
	ChainType            uint8           `db:"chain_type"`
	UserAddress          string          `db:"user_address"`
	TokenAddress         string          `db:"token_address"`
	PoolAddress          string          `db:"pool_address"`
	BaseTokenAmount      uint64          `db:"base_token_amount"`
	QuoteTokenAmount     uint64          `db:"quote_token_amount"`
	Decimals             uint8           `db:"decimals"`
	BaseTokenPrice       decimal.Decimal `db:"base_token_price"`
	QuoteTokenPrice      decimal.Decimal `db:"quote_token_price"`
	TransactionAmountUSD decimal.Decimal `db:"transaction_amount_usd"`
	TransactionType      uint8           `db:"transaction_type"`
	PlatformType         uint8           `db:"platform_type"`
	IsBuyback            bool            `db:"is_buyback"`
	TransactionTime      time.Time       `db:"transaction_time"`
	CreateTime           time.Time       `db:"create_time"`
}

type Kline struct {
	TokenAddress      string          `db:"token_address"`
	IntervalTimestamp time.Time       `db:"interval_timestamp"`
	OpenPrice         decimal.Decimal `db:"open_price"`
	HighPrice         decimal.Decimal `db:"high_price"`
	LowPrice          decimal.Decimal `db:"low_price"`
	ClosePrice        decimal.Decimal `db:"close_price"`
	Volume            uint64          `db:"volume"`
	TradesCount       uint64          `db:"trades_count"`
	BuyCount          uint64          `db:"buy_count"`
	SellCount         uint64          `db:"sell_count"`
	BuyVolume         uint64          `db:"buy_volume"`
	SellVolume        uint64          `db:"sell_volume"`
	BaseVolume        uint64          `db:"base_volume"`
}

// InsertTransaction 插入单条交易数据
func InsertTransaction(tx *TokenTransactionCk) error {
	query := `
        INSERT INTO token_transaction_ck_new (
            transaction_id, transaction_hash, chain_type, user_address, token_address,
            pool_address, base_token_amount, quote_token_amount, decimals,
            base_token_price, quote_token_price, transaction_amount_usd,
            transaction_type, platform_type, is_buyback, transaction_time
        ) VALUES (
            ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
        )
    `

	err := ClickHouseClient.Exec(context.Background(), query,
		tx.TransactionID,
		tx.TransactionHash,
		tx.ChainType,
		tx.UserAddress,
		tx.TokenAddress,
		tx.PoolAddress,
		tx.BaseTokenAmount,
		tx.QuoteTokenAmount,
		tx.Decimals,
		tx.BaseTokenPrice,
		tx.QuoteTokenPrice,
		tx.TransactionAmountUSD,
		tx.TransactionType,
		tx.PlatformType,
		tx.IsBuyback,
		tx.TransactionTime,
	)

	if err != nil {
		return fmt.Errorf("insert transaction failed: %w", err)
	}
	return nil
}

// BatchInsertTransactions 批量插入交易数据
func BatchInsertTransactions(txs []*TokenTransactionCk) error {
	batch, err := ClickHouseClient.PrepareBatch(context.Background(), `
        INSERT INTO token_transaction_ck_new (
            transaction_id, transaction_hash, chain_type, user_address, token_address,
            pool_address, base_token_amount, quote_token_amount, decimals,
            base_token_price, quote_token_price, transaction_amount_usd,
            transaction_type, platform_type, is_buyback, transaction_time
        )
    `)
	if err != nil {
		return fmt.Errorf("prepare batch failed: %w", err)
	}
	defer batch.Abort()

	for _, tx := range txs {
		err := batch.Append(
			tx.TransactionID,
			tx.TransactionHash,
			tx.ChainType,
			tx.UserAddress,
			tx.TokenAddress,
			tx.PoolAddress,
			tx.BaseTokenAmount,
			tx.QuoteTokenAmount,
			tx.Decimals,
			tx.BaseTokenPrice,
			tx.QuoteTokenPrice,
			tx.TransactionAmountUSD,
			tx.TransactionType,
			tx.PlatformType,
			tx.IsBuyback,
			tx.TransactionTime,
		)
		if err != nil {
			return fmt.Errorf("append to batch failed: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("send batch failed: %w", err)
	}
	return nil
}

// GetKlines 查询K数据
func GetKlines(tokenAddress string, interval string, start, end time.Time) ([]Kline, error) {

	rows, err := ClickHouseClient.Query(context.Background(), `
        SELECT
            token_address,
            toStartOfInterval(timestamp, INTERVAL ? second) as interval_timestamp,
            any(open_price) as open_price,
            max(high_price) as high_price,
            min(low_price) as low_price,
            anyLast(open_price) OVER (
                PARTITION BY token_address 
                ORDER BY interval_timestamp 
                ROWS BETWEEN 1 FOLLOWING AND 1 FOLLOWING
            ) as close_price,
            sum(volume) as volume,
            sum(trades_count) as trades_count,
            sum(buy_count) as buy_count,
            sum(sell_count) as sell_count,
            sum(buy_volume) as buy_volume,
            sum(sell_volume) as sell_volume,
            sum(base_volume) as base_volume
        FROM token_kline_1s
        WHERE token_address = ?
          AND timestamp >= ?
          AND timestamp < ?
        GROUP BY token_address, interval_timestamp
        ORDER BY token_address, interval_timestamp;
    `, getIntervalSeconds(interval), tokenAddress, start, end)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var klines []Kline
	for rows.Next() {
		var k Kline

		err := rows.Scan(
			&k.TokenAddress,
			&k.IntervalTimestamp,
			&k.OpenPrice,
			&k.HighPrice,
			&k.LowPrice,
			&k.ClosePrice,
			&k.Volume,
			&k.TradesCount,
			&k.BuyCount,
			&k.SellCount,
			&k.BuyVolume,
			&k.SellVolume,
			&k.BaseVolume,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		klines = append(klines, k)
	}

	if len(klines) > 0 {
		lastKline := &klines[len(klines)-1]
		if lastKline.ClosePrice.IsZero() {
			lastKline.ClosePrice = lastKline.OpenPrice
		}
	}

	return klines, nil
}

func getIntervalSeconds(interval string) int {
	return map[string]int{
		"1S":  1,     // 1 second
		"1":   60,    // 1 minute
		"5":   300,   // 5 minutes
		"15":  900,   // 15 minutes
		"60":  3600,  // 1 hour
		"240": 14400, // 4 hours
		"720": 43200, // 12 hours
		"1D":  86400, // 1 day
	}[interval]
}

// GetTokenTransactions retrieves the latest token transactions from ClickHouse
func GetTokenTransactions(tokenAddress string, limit int) ([]TokenTransactionCk, error) {
	query := `
        SELECT *
        FROM token_transaction_ck_new 
        WHERE token_address = ?
        ORDER BY transaction_time DESC , user_address DESC
        LIMIT ?
    `

	var transactions []TokenTransactionCk
	rows, err := ClickHouseClient.Query(context.Background(), query, tokenAddress, limit)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tx TokenTransactionCk
		var isBuyback uint8 // Use uint8 as intermediate type for boolean field

		if err := rows.Scan(
			&tx.TransactionID,
			&tx.TransactionHash,
			&tx.ChainType,
			&tx.UserAddress,
			&tx.TokenAddress,
			&tx.PoolAddress,
			&tx.BaseTokenAmount,
			&tx.QuoteTokenAmount,
			&tx.Decimals,
			&tx.BaseTokenPrice,
			&tx.QuoteTokenPrice,
			&tx.TransactionAmountUSD,
			&tx.TransactionType,
			&tx.PlatformType,
			&isBuyback, // Scan into uint8 variable
			&tx.TransactionTime,
			&tx.CreateTime,
		); err != nil {
			return nil, fmt.Errorf("scan row failed: %w", err)
		}

		// Explicitly convert the uint8 to bool (only true if value is 1)
		tx.IsBuyback = isBuyback == 1

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// GetLatestKline 获取最新的K线数据
func GetLatestKline(tokenAddress string) (*Kline, error) {
	query := `
        SELECT
            token_address,
            timestamp as interval_timestamp,
            open_price,
            high_price,
            low_price,
            close_price,
            volume,
            trades_count,
            buy_count,
            sell_count,
            buy_volume,
            sell_volume,
            base_volume
        FROM token_kline_1s
        WHERE token_address = ?
            AND timestamp >= now() - INTERVAL 12 HOUR
        ORDER BY timestamp DESC
        LIMIT 1
    `

	var kline Kline
	row := ClickHouseClient.QueryRow(context.Background(), query, tokenAddress)
	err := row.Scan(
		&kline.TokenAddress,
		&kline.IntervalTimestamp,
		&kline.OpenPrice,
		&kline.HighPrice,
		&kline.LowPrice,
		&kline.ClosePrice,
		&kline.Volume,
		&kline.TradesCount,
		&kline.BuyCount,
		&kline.SellCount,
		&kline.BuyVolume,
		&kline.SellVolume,
		&kline.BaseVolume,
	)
	if err != nil {
		return nil, fmt.Errorf("get latest kline failed: %w", err)
	}

	return &kline, nil
}
