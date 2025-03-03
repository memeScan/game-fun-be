package service

import (
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
)

type TickerService interface {
	Tickers(req request.TickersRequest) response.Response
	GetTicker(tokenSymbol string) response.Response
	SwapHistories(tickersId string) response.Response
	TokenDistribution(tickersId string) response.Response
}

type TickerServiceImpl struct{}

func NewTickerService() TickerService {
	return &TickerServiceImpl{}
}

func (s *TickerServiceImpl) Tickers(req request.TickersRequest) response.Response {
	var tickersResponse response.TickersResponse

	return response.Response{
		Code: 200,
		Data: tickersResponse,
		Msg:  "success",
	}
}

func (s *TickerServiceImpl) GetTicker(tokenSymbol string) response.Response {
	var getTickerResponse response.GetTickerResponse

	return response.Response{
		Code: 200,
		Data: getTickerResponse,
		Msg:  "success",
	}
}

func (s *TickerServiceImpl) SwapHistories(tickersId string) response.Response {
	var swapHistoriesResponse response.SwapHistoriesResponse

	return response.Response{
		Code: 200,
		Data: swapHistoriesResponse,
		Msg:  "success",
	}
}

func (s *TickerServiceImpl) TokenDistribution(tickersId string) response.Response {
	var tokenDistributionResponse response.TokenDistributionResponse

	return response.Response{
		Code: 200,
		Data: tokenDistributionResponse,
		Msg:  "success",
	}
}
