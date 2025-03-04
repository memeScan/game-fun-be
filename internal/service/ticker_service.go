package service

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
)

type TickerServiceImpl struct {
}

func NewTickerServiceImpl() *TickerServiceImpl {
	return &TickerServiceImpl{}
}

func (s *TickerServiceImpl) Tickers(req request.TickersRequest, chainType model.ChainType) response.Response {
	var tickersResponse response.TickersResponse
	return response.Success(tickersResponse)
}

func (s *TickerServiceImpl) TickerDetail(tokenSymbol string, chainType model.ChainType) response.Response {
	var getTickerResponse response.GetTickerResponse
	return response.Success(getTickerResponse)
}

func (s *TickerServiceImpl) SwapHistories(tickersId string, chainType model.ChainType) response.Response {
	var swapHistoriesResponse response.SwapHistoriesResponse
	return response.Success(swapHistoriesResponse)

}

func (s *TickerServiceImpl) TokenDistribution(tickersId string, chainType model.ChainType) response.Response {
	var tokenDistributionResponse response.TokenDistributionResponse
	return response.Success(tokenDistributionResponse)
}

func (s *TickerServiceImpl) SearchTickers(param, limit, cursor string, chainType model.ChainType) response.Response {
	var sarchTickerResponse response.SearchTickerResponse
	return response.Success(sarchTickerResponse)
}
