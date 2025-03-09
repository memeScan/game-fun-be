package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
)

type TokenMarketAnalyticsRepo struct{}

func NewTokenMarketAnalyticsRepo() *TokenMarketAnalyticsRepo {
	return &TokenMarketAnalyticsRepo{}
}

type TokenMarketAnalytics struct {
	ChainType         uint8   `json:"chain_type"`
	TokenAddress      string  `json:"token_address"`
	TxCount5M         int     `json:"tx_count_5m"`
	TxCount1H         int     `json:"tx_count_1h"`
	TxCount24H        int     `json:"tx_count_24h"`
	BuyTxCount5M      uint64  `json:"buy_tx_count_5m"`
	BuyTxCount1H      uint64  `json:"buy_tx_count_1h"`
	BuyTxCount24H     uint64  `json:"buy_tx_count_24h"`
	TokenVolume5M     float64 `json:"token_volume_5m"`
	TokenVolume1H     float64 `json:"token_volume_1h"`
	TokenVolume24H    float64 `json:"token_volume_24h"`
	BuyTokenVolume5M  float64 `json:"buy_token_volume_5m"`
	BuyTokenVolume1H  float64 `json:"buy_token_volume_1h"`
	BuyTokenVolume24H float64 `json:"buy_token_volume_24h"`
	Price5M           float64 `json:"price_5m"`
	Price1H           float64 `json:"price_1h"`
	Price24H          float64 `json:"price_24h"`
	CurrentPrice      float64 `json:"current_price"`
	PriceChange5M     float64 `json:"price_change_5m"`
	PriceChange1H     float64 `json:"price_change_1h"`
	PriceChange24H    float64 `json:"price_change_24h"`
}

func (t *TokenMarketAnalyticsRepo) GetTokenMarketAnalytics(tokenAddress string, chainType uint8) (*TokenMarketAnalytics, error) {
	query := `
		SELECT 
			chain_type,
			token_address,
			TxCount5M, TxCount1H, TxCount24H,
			BuyTxCount5M, BuyTxCount1H, BuyTxCount24H,
            TokenVolume5M, TokenVolume1H, TokenVolume24H
			BuyTokenVolume5M, BuyTokenVolume1H, BuyTokenVolume24H,
			Price5M, Price1H, Price24H, CurrentPrice,
			(CurrentPrice - Price5M) / Price5M * 100 AS PriceChange5M,
			(CurrentPrice - Price1H) / Price1H * 100 AS PriceChange1H,
			(CurrentPrice - Price24H) / Price24H * 100 AS PriceChange24H
		FROM token_transaction_stats_mv
		WHERE chain_type = ? 
		AND token_address = ?
		ORDER BY timestamp DESC
		LIMIT 1;
	`

	var analytics TokenMarketAnalytics
	row := ClickHouseClient.QueryRow(context.Background(), query, chainType, tokenAddress)
	err := row.Scan(
		&analytics.ChainType,
		&analytics.TokenAddress,
		&analytics.TxCount5M,
		&analytics.TxCount1H,
		&analytics.TxCount24H,
		&analytics.BuyTxCount5M,
		&analytics.BuyTxCount1H,
		&analytics.BuyTxCount24H,
		&analytics.TokenVolume5M,
		&analytics.TokenVolume1H,
		&analytics.TokenVolume24H,
		&analytics.BuyTokenVolume5M,
		&analytics.BuyTokenVolume1H,
		&analytics.BuyTokenVolume24H,
		&analytics.Price5M,
		&analytics.Price1H,
		&analytics.Price24H,
		&analytics.CurrentPrice,
		&analytics.PriceChange5M,
		&analytics.PriceChange1H,
		&analytics.PriceChange24H,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query token_transaction_stats_mv: %v", err)
	}

	return &analytics, nil
}
