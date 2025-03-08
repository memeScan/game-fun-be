package model

import (
	"fmt"
	"game-fun-be/internal/constants"
	"game-fun-be/internal/redis"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TokenTransaction 代币交易记录模型
type TokenTransaction struct {
	ID                    uint64          `gorm:"column:id;primaryKey;autoIncrement;type:bigint unsigned" json:"id"`
	TransactionHash       string          `gorm:"column:transaction_hash;type:varchar(88)" json:"transaction_hash"`
	TransactionType       uint8           `gorm:"column:transaction_type;type:tinyint unsigned" json:"transaction_type"`
	MarketAddress         string          `gorm:"column:market_address;type:varchar(64)" json:"market_address"`
	PoolAddress           string          `gorm:"column:pool_address;type:varchar(64)" json:"pool_address"`
	TokenAddress          string          `gorm:"column:token_address;type:varchar(64);not null" json:"token_address"`
	NativeTokenAddress    string          `gorm:"column:native_token_address;type:varchar(64)" json:"native_token_address"`
	TransactionTime       time.Time       `gorm:"column:transaction_time;type:datetime" json:"transaction_time"`
	Block                 uint64          `gorm:"column:block;type:bigint unsigned;not null" json:"block"`
	PlatformType          uint8           `gorm:"column:platform_type;type:tinyint unsigned" json:"platform_type"`
	ChainType             uint8           `gorm:"column:chain_type;type:tinyint unsigned" json:"chain_type"`
	ProxyType             uint8           `gorm:"column:proxy_type;type:tinyint unsigned" json:"proxy_type"`
	NativeTokenAmount     uint64          `gorm:"column:native_token_amount;type:bigint unsigned;not null" json:"native_token_amount"`
	TokenAmount           uint64          `gorm:"column:token_amount;type:bigint unsigned;not null" json:"token_amount"`
	Decimals              uint8           `gorm:"column:decimals;type:tinyint unsigned;default:6" json:"decimals"`
	UserAddress           string          `gorm:"column:user_address;type:varchar(64);not null" json:"user_address"`
	VirtualNativeReserves uint64          `gorm:"column:virtual_native_reserves;type:bigint unsigned;not null" json:"virtual_native_reserves"`
	VirtualTokenReserves  uint64          `gorm:"column:virtual_token_reserves;type:bigint unsigned;not null" json:"virtual_token_reserves"`
	RealNativeReserves    uint64          `gorm:"column:real_native_reserves;type:bigint unsigned;default:0" json:"real_native_reserves"`
	RealTokenReserves     uint64          `gorm:"column:real_token_reserves;type:bigint unsigned;default:0" json:"real_token_reserves"`
	IsBuy                 bool            `gorm:"column:is_buy;type:tinyint(1);not null" json:"is_buy"`
	IsBuyback             bool            `gorm:"column:is_buyback;type:tinyint(1);not null" json:"is_buyback"`
	Progress              decimal.Decimal `gorm:"column:progress;type:decimal(5,2);not null" json:"progress"`
	IsComplete            bool            `gorm:"column:is_complete;type:tinyint(1);not null" json:"is_complete"`
	Price                 decimal.Decimal `gorm:"column:price;type:decimal(30,18)" json:"price"`
	NativePrice           decimal.Decimal `gorm:"column:native_price;type:decimal(30,18)" json:"native_price"`
	NativePriceUSD        decimal.Decimal `gorm:"column:native_price_usd;type:decimal(30,18)" json:"native_price_usd"`
	TransactionAmountUSD  decimal.Decimal `gorm:"column:transaction_amount_usd;type:decimal(30,18)" json:"transaction_amount_usd"`
	CreateTime            time.Time       `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime            time.Time       `gorm:"column:update_time;type:datetime" json:"update_time"`
}

// TableName 指定表名
func (TokenTransaction) TableName(date string) string {
	return "token_transaction_" + date
}

// CreateTokenTransaction 创建代币交易记录
func CreateTokenTransaction(tx *TokenTransaction, date string) error {
	tableName := getTableName(date)

	return DB.Table(tableName).Create(tx).Error
}

// GetTokenTransactionByHash 通过交易哈希获取代币交易记录
func GetTokenTransactionByHash(date string, transactionHash string, tokenAddress string) (*TokenTransaction, error) {
	var tx TokenTransaction
	err := DB.Table(getTableName(date)).Where("transaction_hash = ? and token_address = ?", transactionHash, tokenAddress).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// GetTokenTransactionByID 通过交易ID获取代币交易记录
func GetTokenTransactionByID(date string, transactionID uint64) (*TokenTransaction, error) {
	var tx TokenTransaction
	err := DB.Table(getTableName(date)).Where("id = ? ", transactionID).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// UpdateTokenTransaction 更新代币交易记录
func UpdateTokenTransaction(tx *TokenTransaction) error {
	return DB.Save(tx).Error
}

// ListTokenTransactions 列出代币交易记录
func ListTokenTransactions(limit, offset int) ([]TokenTransaction, error) {
	var txs []TokenTransaction
	err := DB.Limit(limit).Offset(offset).Find(&txs).Error
	return txs, err
}

// GetTokenTransactionsByUserAddress 通过用户地址获取代币交易记录
func GetTokenTransactionsByUserAddress(userAddress string, limit, offset int) ([]TokenTransaction, error) {
	var txs []TokenTransaction
	err := DB.Where("user_address = ?", userAddress).Limit(limit).Offset(offset).Find(&txs).Error
	return txs, err
}

// GetTokenTransactionsByChainAndToken 通过链类型和代币地址获取代币交易记录
func GetTokenTransactionsByChainAndToken(chainType uint8, tokenAddress string, limit, offset int) ([]TokenTransaction, error) {
	var txs []TokenTransaction
	err := DB.Where("chain_type = ? AND token_address = ?", chainType, tokenAddress).Limit(limit).Offset(offset).Find(&txs).Error
	return txs, err
}

// CreateTableForDate 创建指定日期的表
func CreateTableForDate(date string) error {
	tableName := (&TokenTransaction{}).TableName(date)
	return DB.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分布式id',
			transaction_hash VARCHAR(88) NOT NULL COMMENT '交易哈希',
			transaction_type TINYINT UNSIGNED DEFAULT NULL COMMENT '交易类型：1-买, 2-卖, 3-加池子, 4-减池子, 5-烧币',
			market_address VARCHAR(64) COMMENT '市场地址',
			pool_address VARCHAR(64) COMMENT '交易对池子地址',
			token_address VARCHAR(64) NOT NULL COMMENT '代币合约地址',
			native_token_address VARCHAR(64) COMMENT '原生定价代币的代币地址',
			transaction_time datetime DEFAULT NULL COMMENT '交易时间',
			block BIGINT UNSIGNED NOT NULL COMMENT '区块高度',
			platform_type TINYINT UNSIGNED DEFAULT NULL COMMENT '交易平台类型：1-Pump, 2-Raydium',
			chain_type TINYINT UNSIGNED DEFAULT NULL COMMENT '链类型：1-Solana, 2-Ethereum',
			proxy_type TINYINT UNSIGNED DEFAULT NULL COMMENT '代理类型',
			native_token_amount BIGINT UNSIGNED NOT NULL COMMENT '链上原生代币数量',
			token_amount BIGINT UNSIGNED NOT NULL COMMENT 'Token 数量',
			decimals TINYINT UNSIGNED DEFAULT 6 COMMENT '精度',
			user_address VARCHAR(64) NOT NULL COMMENT '用户地址',
			virtual_native_reserves BIGINT UNSIGNED NOT NULL COMMENT '虚拟流动性池子中原生代币的数量',
			virtual_token_reserves BIGINT UNSIGNED NOT NULL COMMENT '虚拟流动性池子中 token 的数量',
			real_native_reserves BIGINT UNSIGNED DEFAULT 0 COMMENT '流动性池子中原生代币的数量',
			real_token_reserves BIGINT UNSIGNED DEFAULT 0 COMMENT '流动性池子中 token 的数量',
			is_buy TINYINT(1) NOT NULL COMMENT '是否是买入',
			is_buyback TINYINT(1) NOT NULL COMMENT '是否是回购',
			progress DECIMAL(5, 2) NOT NULL COMMENT '当前内盘进度',
			is_complete TINYINT(1) NOT NULL COMMENT '内盘是否已完成',
			price DECIMAL(30, 18) DEFAULT NULL COMMENT 'token的u本位价格',
			native_price DECIMAL(30, 18) DEFAULT NULL COMMENT 'token的原生代币本位价格',
			native_price_usd DECIMAL(30, 18) DEFAULT NULL COMMENT '原生代币的u本位价格',
			transaction_amount_usd DECIMAL(30, 18) DEFAULT NULL COMMENT '交易额(USD)',
			create_time datetime DEFAULT NULL COMMENT '创建时间',
			update_time datetime DEFAULT NULL COMMENT '更新时间',
			PRIMARY KEY (id),
			KEY IDX_CHAIN_TYPE_TOKEN_ADDRESS (chain_type, token_address),
			KEY IDX_USER_ADDRESS (user_address),
			KEY IDX_HASH (transaction_hash),
			KEY IDX_PROXY (proxy_type)
		) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT='代币交易记录表'
	`, tableName)).Error
}

// getTableName 获取表名
func getTableName(date string) string {
	return (&TokenTransaction{}).TableName(date)
}

// BatchCreateTokenTransactions 批量创建代币交易记录
func BatchCreateTokenTransactions(txs []*TokenTransaction, date string) error {
	if len(txs) == 0 {
		return nil
	}

	tableName := getTableName(date)

	// 从Redis获取一批分布式ID
	startID, _, err := redis.GetBatchIDs(constants.RedisKeyTokenTransactionID, int64(len(txs)))
	if err != nil {
		return fmt.Errorf("获取分布式ID失败: %v", err)
	}

	// 为每个交易分配ID
	currentID := startID
	for _, tx := range txs {
		tx.ID = uint64(currentID)
		currentID++
	}

	// 用占位符和参数绑定
	placeholder := "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	placeholders := make([]string, len(txs))
	for i := range placeholders {
		placeholders[i] = placeholder
	}

	// 构建参数数组
	valueArgs := make([]interface{}, 0, len(txs)*30)
	for _, tx := range txs {
		valueArgs = append(valueArgs,
			tx.ID,
			tx.TransactionHash,
			tx.TransactionType,
			tx.MarketAddress,
			tx.PoolAddress,
			tx.TokenAddress,
			tx.NativeTokenAddress,
			tx.TransactionTime,
			tx.Block,
			tx.PlatformType,
			tx.ChainType,
			tx.ProxyType,
			tx.NativeTokenAmount,
			tx.TokenAmount,
			tx.Decimals,
			tx.UserAddress,
			tx.VirtualNativeReserves,
			tx.VirtualTokenReserves,
			tx.RealNativeReserves,
			tx.RealTokenReserves,
			tx.IsBuy,
			tx.IsBuyback,
			tx.Progress,
			tx.IsComplete,
			tx.Price,
			tx.NativePrice,
			tx.NativePriceUSD,
			tx.TransactionAmountUSD,
			tx.CreateTime,
			tx.UpdateTime,
		)
	}

	query := fmt.Sprintf(`INSERT IGNORE INTO %s (
		id,
		transaction_hash,
		transaction_type,
		market_address,
		pool_address,
		token_address,
		native_token_address,
		transaction_time,
		block,
		platform_type,
		chain_type,
		proxy_type,
		native_token_amount,
		token_amount,
		decimals,
		user_address,
		virtual_native_reserves,
		virtual_token_reserves,
		real_native_reserves,
		real_token_reserves,
		is_buy,
		is_buyback,
		progress,
		is_complete,
		price,
		native_price,
		native_price_usd,
		transaction_amount_usd,
		create_time,
		update_time
	) VALUES %s`, tableName, strings.Join(placeholders, ","))

	// 保持原有的 Session 配置
	return DB.Session(&gorm.Session{
		PrepareStmt: false,
	}).Exec(query, valueArgs...).Error
}
