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

// BatchCreateIndexFromTransactions 批量创建交易索引
func (service *TokenTxIndexService) BatchCreateIndexFromTransactions(transactions []*model.TokenTransaction) response.Response {
	indexes := make([]*model.TokenTxIndex, 0, len(transactions))
	for _, tx := range transactions {
		index := &model.TokenTxIndex{
			TransactionID:   tx.ID,
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
