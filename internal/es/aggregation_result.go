package es

import "encoding/json"

// 主要聚合结构体
type AggregationResult struct {
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`

	// 持有者统计
	HolderCount            HolderCountStats       `json:"holder_count"`
	Holders                TopHoldersAggregation  `json:"holders"`
	TotalHoldersPercentage TotalHoldersPercentage `json:"total_holders_percentage"`

	// 市值相关
	LatestTransactionMarketCap MarketCapTimeframe `json:"latest_transaction_market_cap"`
	MarketCap1m                MarketCapTimeframe `json:"market_cap_1m"`
	MarketCap5m                MarketCapTimeframe `json:"market_cap_5m"`
	MarketCap1h                MarketCapTimeframe `json:"market_cap_1h"`
	MarketCap6h                MarketCapTimeframe `json:"market_cap_6h"`
	MarketCap24h               MarketCapTimeframe `json:"market_cap_24h"`
	MarketCapTime              MarketCapTimeframe `json:"market_cap_time"`
	LatestMarketCap            struct {
		Hits struct {
			Hits []struct {
				Source MarketCapSource `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	} `json:"latest_market_cap"`

	// 交易相关
	Swaps    SwapStats `json:"swaps"`
	Swaps1m  SwapStats `json:"swaps_1m"`
	Swaps5m  SwapStats `json:"swaps_5m"`
	Swaps1h  SwapStats `json:"swaps_1h"`
	Swaps6h  SwapStats `json:"swaps_6h"`
	Swaps24h SwapStats `json:"swaps_24h"`

	// 交易量相关
	Volume    Volume      `json:"volume"`
	Volume1m  VolumeStats `json:"volume_1m"`
	Volume5m  VolumeStats `json:"volume_5m"`
	Volume1h  VolumeStats `json:"volume_1h"`
	Volume6h  VolumeStats `json:"volume_6h"`
	Volume24h VolumeStats `json:"volume_24h"`

	// 买入量相关
	BuyVolume1m  VolumeStats `json:"buy_volume_1m"`
	BuyVolume5m  VolumeStats `json:"buy_volume_5m"`
	BuyVolume1h  VolumeStats `json:"buy_volume_1h"`
	BuyVolume24h VolumeStats `json:"buy_volume_24h"`

	// 卖出量相关
	SellVolume1m  VolumeStats `json:"sell_volume_1m"`
	SellVolume5m  VolumeStats `json:"sell_volume_5m"`
	SellVolume1h  VolumeStats `json:"sell_volume_1h"`
	SellVolume24h VolumeStats `json:"sell_volume_24h"`

	// 买入统计
	Buys        BuyStats `json:"Buys"`
	BuyCount1m  BuyStats `json:"buy_count_1m"`
	BuyCount5m  BuyStats `json:"buy_count_5m"`
	BuyCount1h  BuyStats `json:"buy_count_1h"`
	BuyCount6h  BuyStats `json:"buy_count_6h"`
	BuyCount24h BuyStats `json:"buy_count_24h"`

	// 卖出统计
	Sells        SellStats `json:"Sells"`
	SellCount1m  SellStats `json:"sell_count_1m"`
	SellCount5m  SellStats `json:"sell_count_5m"`
	SellCount1h  SellStats `json:"sell_count_1h"`
	SellCount6h  SellStats `json:"sell_count_6h"`
	SellCount24h SellStats `json:"sell_count_24h"`

	// 交易记录
	LastTransaction1h TransactionTimeframe `json:"last_transaction_1h"`
	LastTransaction4h TransactionTimeframe `json:"last_transaction_4h"`
	LatestTransaction struct {
		Hits struct {
			BaseHits
			Hits []RawMessageHit `json:"hits"`
		} `json:"hits"`
	} `json:"latest_transaction"`

	// 价格相关
	LastTransactionPrice    PriceTimeframe `json:"last_transaction_price"`
	LastTransaction1mPrice  PriceTimeframe `json:"last_transaction_1m_price"`
	LastTransaction5mPrice  PriceTimeframe `json:"last_transaction_5m_price"`
	LastTransaction1hPrice  PriceTimeframe `json:"last_transaction_1h_price"`
	LastTransaction4hPrice  PriceTimeframe `json:"last_transaction_4h_price"`
	LastTransaction24hPrice PriceTimeframe `json:"last_transaction_24h_price"`
}

// 基础结构
type BaseTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type BaseHits struct {
	Total    BaseTotal   `json:"total"`
	MaxScore interface{} `json:"max_score"`
}

type BaseHit struct {
	Index string        `json:"_index"`
	Type  string        `json:"_type"`
	ID    string        `json:"_id"`
	Score interface{}   `json:"_score"`
	Sort  []interface{} `json:"sort"`
}

// Source 结构体
type MarketCapSource struct {
	MarketCap float64 `json:"market_cap"`
}

type PriceSource struct {
	Price           float64 `json:"price"`
	Decimals        int     `json:"decimals,omitempty"`
	NativePrice     float64 `json:"native_price,omitempty"`
	TransactionTime int64   `json:"native_price,transaction_time"`
}

type MarketCapHit struct {
	BaseHit
	Source MarketCapSource `json:"_source"`
}

type PriceHit struct {
	BaseHit
	Source PriceSource `json:"_source"`
}

type RawMessageHit struct {
	BaseHit
	Source json.RawMessage `json:"_source"`
}

// 统计结构体
type Volume struct {
	Value float64 `json:"value"`
}

type VolumeStats struct {
	DocCount    int `json:"doc_count"`
	TotalVolume struct {
		Value float64 `json:"value"`
	} `json:"total_volume"`
}

type SwapStats struct {
	DocCount         int `json:"doc_count"`
	TransactionCount struct {
		Value int64 `json:"value"`
	} `json:"transaction_count"`
}

type BuyStats struct {
	DocCount  int `json:"doc_count"`
	BuyVolume struct {
		Value int64 `json:"value"`
	} `json:"buy_volume"`
}

type SellStats struct {
	DocCount   int `json:"doc_count"`
	SellVolume struct {
		Value int64 `json:"value"`
	} `json:"sell_volume"`
}

type HolderCountStats struct {
	DocCount    int `json:"doc_count"`
	UniqueUsers struct {
		Value int `json:"value"`
	} `json:"unique_users"`
}

// 时间段结构体
type MarketCapTimeframe struct {
	DocCount          int `json:"doc_count"`
	LatestTransaction struct {
		Hits struct {
			BaseHits
			Hits []MarketCapHit `json:"hits"`
		} `json:"hits"`
	} `json:"latest_transaction"`
}

type TransactionTimeframe struct {
	Hits struct {
		Hits []struct {
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type PriceTimeframe struct {
	DocCount int `json:"doc_count"`
	Latest   struct {
		Hits struct {
			BaseHits
			Hits []PriceHit `json:"hits"`
		} `json:"hits"`
	} `json:"latest"`
}

type TopHoldersAggregation struct {
	Buckets []HolderStats `json:"buckets"`
}

type TotalHoldersPercentage struct {
	Value float64 `json:"value"`
}

type HolderStats struct {
	DocCount   int    `json:"doc_count"`
	Key        string `json:"key"` // user_address
	BalanceSum struct {
		Value float64 `json:"value"`
	} `json:"balance_sum"`
	TokenSupply struct {
		Value float64 `json:"value"`
	} `json:"token_supply"`
	HolderPercentage struct {
		Value float64 `json:"value"`
	} `json:"holder_percentage"`
}
