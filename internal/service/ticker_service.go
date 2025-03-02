package service

import (
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
)

type TickerService struct{}

func NewTickerService() *TickerService {
	return &TickerService{}
}

func (s *TickerService) Tickers(req request.TickersRequest) response.Response {
	var tickersResponse response.TickersResponse

	return response.Response{
		Code: 200,
		Data: tickersResponse,
		Msg:  "success",
	}
}

func (s *TickerService) GetTicker(tokenSymbol string) response.Response {
	var getTickerResponse response.GetTickerResponse

	return response.Response{
		Code: 200,
		Data: getTickerResponse,
		Msg:  "success",
	}
}
