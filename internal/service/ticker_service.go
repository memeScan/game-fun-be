package service

import (
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
)

type TickerService interface {
	Tickers(req request.TickersRequest) response.Response
	TickerDetail(tokenSymbol string) response.Response
	SwapHistories(tickersId string) response.Response
	TokenDistribution(tickersId string) response.Response
	SearchTickers(param, limit, cursor string) response.Response
}

type TickerServiceImpl struct{}

func NewTickerServiceImpl() TickerService {
	return &TickerServiceImpl{}
}

func (s *TickerServiceImpl) Tickers(req request.TickersRequest) response.Response {
	var tickersResponse response.TickersResponse
	return response.Success(tickersResponse)
}

func (s *TickerServiceImpl) TickerDetail(tokenSymbol string) response.Response {
	var getTickerResponse response.GetTickerResponse
	return response.Success(getTickerResponse)
}

func (s *TickerServiceImpl) SwapHistories(tickersId string) response.Response {
	var swapHistoriesResponse response.SwapHistoriesResponse
	return response.Success(swapHistoriesResponse)

}

func (s *TickerServiceImpl) TokenDistribution(tickersId string) response.Response {
	var tokenDistributionResponse response.TokenDistributionResponse
	return response.Success(tokenDistributionResponse)
}

func (s *TickerServiceImpl) SearchTickers(param, limit, cursor string) response.Response {
	var sarchTickerResponse response.SearchTickerResponse
	return response.Success(sarchTickerResponse)
}
