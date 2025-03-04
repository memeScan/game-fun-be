package service

import (
	"context"
	"fmt"
	"game-fun-be/internal/model"
	"game-fun-be/internal/response"
	"time"

	"github.com/olivere/elastic/v7"
)

var esClient *elastic.Client

// TokenTxIndexService 代币交易索引服务
type TokenTxIndexService struct{}

// 初始化函数
func InitTokenTxIndexService(client *elastic.Client) {
	esClient = client
}

// CreateTokenTxIndex 创建代币交易索引记录
func (service *TokenTxIndexService) CreateTokenTxIndex(index *model.TokenTxIndex) (*model.TokenTxIndex, error) {
	err := model.CreateTokenTxIndex(index)
	if err != nil {
		return nil, err
	}
	return index, nil
}

// UpdateTokenTxIndex 更新代币交易索引记录
func (service *TokenTxIndexService) UpdateTokenTxIndex(index *model.TokenTxIndex) error {
	return model.UpdateTokenTxIndex(index)
}

// DeleteTokenTxIndex 删除代币交易索引记录
func (service *TokenTxIndexService) DeleteTokenTxIndex(hash string, tokenAddress string) error {
	return model.DeleteTokenTxIndex(hash, tokenAddress)
}

// ListTokenTxIndices 列出代币交易索引记录
func (service *TokenTxIndexService) ListTokenTxIndices(chainType uint8, tokenAddress string, limit, offset int) ([]model.TokenTxIndex, error) {
	return model.ListTokenTxIndices(chainType, tokenAddress, limit, offset)
}

// ProcessTokenTxIndexCreation 处理代币交易索引记录创建
func (service *TokenTxIndexService) ProcessTokenTxIndexCreation(index *model.TokenTxIndex) response.Response {
	createdIndex, err := service.CreateTokenTxIndex(index)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to create token transaction index", err)
	}

	return response.Response{
		Code: 0,
		Data: createdIndex,
		Msg:  "Token transaction index created successfully",
	}
}

// ProcessTokenTxIndexUpdate 处理代币交易索引记录更新
func (service *TokenTxIndexService) ProcessTokenTxIndexUpdate(index *model.TokenTxIndex) response.Response {
	err := service.UpdateTokenTxIndex(index)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to update token transaction index", err)
	}

	return response.Response{
		Code: 0,
		Data: index,
		Msg:  "Token transaction index updated successfully",
	}
}

// ProcessTokenTxIndexDeletion 处理代币交易索引记录删除
func (service *TokenTxIndexService) ProcessTokenTxIndexDeletion(hash string, tokenAddress string) response.Response {
	err := service.DeleteTokenTxIndex(hash, tokenAddress)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to delete token transaction index", err)
	}

	return response.Response{
		Code: 0,
		Msg:  "Token transaction index deleted successfully",
	}
}

// ProcessTokenTxIndexQuery 处理代币交易索引记录查询
func (service *TokenTxIndexService) ProcessTokenTxIndexQuery(chainType uint8, tokenAddress string, limit, offset int) response.Response {
	indices, err := service.ListTokenTxIndices(chainType, tokenAddress, limit, offset)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to query token transaction indices", err)
	}

	return response.Response{
		Code: 0,
		Data: indices,
		Msg:  "Token transaction indices queried successfully",
	}
}

// ConvertTokenTransactionToIndex 将 TokenTransaction 转换为 TokenTxIndex
func (service *TokenTxIndexService) ConvertTokenTransactionToIndex(tx *model.TokenTransaction) *model.TokenTxIndex {
	return &model.TokenTxIndex{
		TransactionHash: tx.TransactionHash,
		TokenAddress:    tx.TokenAddress,
		TransactionDate: tx.TransactionTime.Truncate(24 * time.Hour), // 只保留日期部分
		ChainType:       tx.ChainType,
	}
}

// CreateIndexFromTransaction 从 TokenTransaction 创建 TokenTxIndex 并保存到数据库
func (service *TokenTxIndexService) CreateIndexFromTransaction(tx *model.TokenTransaction) response.Response {
	// 转换 TokenTransaction 为 TokenTxIndex
	index := service.ConvertTokenTransactionToIndex(tx)

	// 调用 ProcessTokenTxIndexCreation 创建索引记录
	return service.ProcessTokenTxIndexCreation(index)
}

// BatchCreateIndexFromTransactions 批量创建交易索引
func (service *TokenTxIndexService) BatchCreateIndexFromTransactions(transactions []*model.TokenTransaction) response.Response {
	indexes := make([]*model.TokenTxIndex, 0, len(transactions))
	for _, tx := range transactions {
		index := &model.TokenTxIndex{
			TransactionHash: tx.TransactionHash,
			TokenAddress:    tx.TokenAddress,
			ChainType:       tx.ChainType,
			TransactionDate: tx.TransactionTime.Truncate(24 * time.Hour), // 只保留日期部分
		}
		indexes = append(indexes, index)
	}

	err := model.BatchCreateTokenTxIndexes(indexes)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to create batch token transaction indexes", err)
	}

	return response.Response{
		Code: 0,
		Msg:  fmt.Sprintf("%d token transaction indexes created successfully", len(indexes)),
	}
}

func BatchIndexTokenTx(txs []*model.TokenTransaction) error {
	bulk := esClient.Bulk().Index("token_tx")
	for _, tx := range txs {
		doc := elastic.NewBulkIndexRequest().Doc(tx)
		bulk.Add(doc)
	}
	_, err := bulk.Do(context.Background())
	return err
}
