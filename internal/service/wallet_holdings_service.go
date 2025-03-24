package service

import (
	"game-fun-be/internal/model"
	"game-fun-be/internal/response"
)

type WalletHoldingsServiceImpl struct{}

func NewWalletHoldingsServiceImpl() *WalletHoldingsServiceImpl {
	return &WalletHoldingsServiceImpl{}
}

func (s *TokenHoldingsServiceImpl) WalletHoldings(userAccount, targetAccount, allowZeroBalance string, chainType model.ChainType) response.Response {
	var tokenHoldingsResponse response.TokenHoldingsResponse
	return response.Success(tokenHoldingsResponse)
}

func (s *TokenHoldingsServiceImpl) WalletHoldingsHistories(userAccount, page, limit string, chainType model.ChainType) response.Response {
	var tokenHoldingHistoriesResponse response.TokenHoldingHistoriesResponse
	return response.Success(tokenHoldingHistoriesResponse)
}
