package model

import (
	"fmt"
	"game-fun-be/internal/pkg/util"
	"strings"
	"time"

	"gorm.io/gorm"
)

// TokenTxIndex 代币交易索引模型
type TokenTxIndex struct {
	TransactionHash string    `gorm:"column:transaction_hash;primaryKey;type:varchar(88)" json:"transaction_hash"`
	TokenAddress    string    `gorm:"column:token_address;type:varchar(64);not null" json:"token_address"`
	TransactionDate time.Time `gorm:"column:transaction_date;type:date;not null" json:"transaction_date"`
	ChainType       uint8     `gorm:"column:chain_type;type:tinyint;not null" json:"chain_type"`
}

// TableName 返回表名
func (TokenTxIndex) TableName() string {
	return "token_tx_index"
}

// GetShardTableName 返回分片表名
func (t *TokenTxIndex) GetShardTableName() string {
	shardIndex := util.HashString(t.TokenAddress) % 128
	return fmt.Sprintf("token_tx_index_%03d", shardIndex)
}

// CreateTokenTxIndex 创建代币交易索引记录
func CreateTokenTxIndex(index *TokenTxIndex) error {
	return DB.Table(index.GetShardTableName()).Create(index).Error
}

// UpdateTokenTxIndex 更新代币交易索引记录
func UpdateTokenTxIndex(index *TokenTxIndex) error {
	return DB.Table(index.GetShardTableName()).Save(index).Error
}

// DeleteTokenTxIndex 删除代币交易索引记录
func DeleteTokenTxIndex(hash string, tokenAddress string) error {
	index := &TokenTxIndex{TransactionHash: hash, TokenAddress: tokenAddress}
	return DB.Table(index.GetShardTableName()).Where("transaction_hash = ?", hash).Delete(index).Error
}

// ListTokenTxIndices 列出代币交易索引记录
func ListTokenTxIndices(chainType uint8, tokenAddress string, limit, offset int) ([]TokenTxIndex, error) {
	var indices []TokenTxIndex
	index := &TokenTxIndex{TokenAddress: tokenAddress}
	err := DB.Table(index.GetShardTableName()).
		Where("chain_type = ? AND token_address = ?", chainType, tokenAddress).
		Limit(limit).Offset(offset).
		Find(&indices).Error
	return indices, err
}

// CreateTokenTxIndexTables 创建代币交易索引表
func CreateTokenTxIndexTables() error {
	for i := 0; i < 128; i++ {
		tableName := fmt.Sprintf("token_tx_index_%03d", i)
		err := DB.Exec(fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				transaction_hash VARCHAR(88) NOT NULL COMMENT '交易哈希',
				token_address VARCHAR(64) NOT NULL COMMENT '代币合约地址',
				transaction_date DATE NOT NULL COMMENT '交易日期',
				chain_type TINYINT NOT NULL COMMENT '链类型：1-Solana, 2-Ethereum',
				PRIMARY KEY (transaction_hash),
				KEY idx_chain_token (chain_type, token_address)
			) COMMENT = '代币交易索引表'
		`, tableName)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchCreateTokenTxIndexes 批量创建代币交易索引，使用 INSERT IGNORE
func BatchCreateTokenTxIndexes(indexes []*TokenTxIndex) error {
	if len(indexes) == 0 {
		return nil
	}

	// 按照 tokenAddress 取模分组
	groupedIndexes := make(map[int][]*TokenTxIndex)
	for _, index := range indexes {
		mod := int(util.HashString(index.TokenAddress) % 128)
		groupedIndexes[mod] = append(groupedIndexes[mod], index)
	}

	// 对每个分组执行批量插入
	for mod, groupIndexes := range groupedIndexes {
		tableName := fmt.Sprintf("token_tx_index_%03d", mod)

		// 使用占位符
		placeholder := "(?,?,?,?)" // 4个字段
		placeholders := make([]string, len(groupIndexes))
		for i := range placeholders {
			placeholders[i] = placeholder
		}

		// 构建参数数组，预分配容量
		valueArgs := make([]interface{}, 0, len(groupIndexes)*4) // 4个字段
		for _, index := range groupIndexes {
			valueArgs = append(valueArgs,
				index.TransactionHash,
				index.TokenAddress,
				index.ChainType,
				index.TransactionDate,
			)
		}

		// 构建 SQL
		query := fmt.Sprintf(`INSERT IGNORE INTO %s (
			transaction_hash,
			token_address,
			chain_type,
			transaction_date
		) VALUES %s`, tableName, strings.Join(placeholders, ","))

		// 执行批量插入，使用参数绑定
		err := DB.Session(&gorm.Session{
			PrepareStmt: false,
		}).Exec(query, valueArgs...).Error

		if err != nil {
			return err
		}
	}

	return nil
}

// GetLatestTokenTxIndexByTokenAddress 根据tokenAddress获取最新的交易索引
func GetLatestTokenTxIndexByTokenAddress(tokenAddress string, chainType uint8) (*TokenTxIndex, error) {
	var index TokenTxIndex
	index.TokenAddress = tokenAddress
	err := DB.Table(index.GetShardTableName()).
		Where("token_address = ? AND chain_type = ?", tokenAddress, chainType).
		Order("transaction_date DESC").
		First(&index).Error
	return &index, err
}
