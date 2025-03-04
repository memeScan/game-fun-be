package service

import (
	"my-token-ai-be/internal/response"
)

type TokenHoldingsServiceImpl struct{}

func NewTokenHoldingsServiceImpl() *TokenHoldingsServiceImpl {
	return &TokenHoldingsServiceImpl{}
}

func (s *TokenHoldingsServiceImpl) TokenHoldings(userAccount, targetAccount, allowZeroBalance string) response.Response {
	var tokenHoldingsResponse response.TokenHoldingsResponse
	return response.Success(tokenHoldingsResponse)
}

func (s *TokenHoldingsServiceImpl) TokenHoldingsHistories(userAccount, page, limit string) response.Response {
	var tokenHoldingHistoriesResponse response.TokenHoldingHistoriesResponse
	return response.Success(tokenHoldingHistoriesResponse)
}
