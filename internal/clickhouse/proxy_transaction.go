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

// InsertTransaction 插入单条交易数据
func InsertProxyTransaction(tx *ProxyTransaction) error {
	query := `
        INSERT INTO proxy_transaction_ck (
            transaction_hash, chain_type, proxy_type, user_address, token_address,
            pool_address, base_token_amount, quote_token_amount, base_token_reserve_amount,
            quote_token_reserve_amount, decimals, base_token_price, quote_token_price,
            transaction_type, is_burn, points_amount, feeQuote_amount, feeBase_amount,
            buybackFeeBase_amount, block_time, transaction_time
        ) VALUES (
            ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
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
	)
	if err != nil {
		return fmt.Errorf("insert transaction failed: %w", err)
	}
	return nil
}
