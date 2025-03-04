package service

import (
	"game-fun-be/internal/response"
	"net/http"
)

type ResponseBuilder struct {
	transactions []*response.TokenTransaction
	orderBy      string
	direction    string
}

func NewResponseBuilder(txs []*response.TokenTransaction, orderBy, direction string) *ResponseBuilder {
	return &ResponseBuilder{
		transactions: txs,
		orderBy:      orderBy,
		direction:    direction,
	}
}

func (rb *ResponseBuilder) BuildPumpResponse() response.Response {
	if rb.transactions == nil {
		return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusNotFound, "No data found", nil)
	}

	responses := response.ConvertAndSortSolPumpSoaringTransactions(
		rb.transactions,
		rb.orderBy,
		rb.direction,
	)

	return response.BuildResponse(responses, http.StatusOK, "Success", nil)
}

func (rb *ResponseBuilder) BuildRaydiumResponse() response.Response {
	if rb.transactions == nil {
		return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusOK, "No data found", nil)
	}

	responses := response.ConverSolRaydiumTransactions(
		rb.transactions,
		rb.orderBy,
		rb.direction,
	)

	return response.BuildResponse(responses, http.StatusOK, "Success", nil)
}

func (rb *ResponseBuilder) BuildNewCreateResponse() response.Response {
	if rb.transactions == nil {
		return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusOK, "No data found", nil)
	}
	responses := response.ConvertAndSortSolPumpNewCreateTransactions(
		rb.transactions,
		rb.orderBy,
		rb.direction,
	)
	return response.BuildResponse(responses, http.StatusOK, "Success", nil)
}

func (rb *ResponseBuilder) BuildCompletingResponse() response.Response {
	if rb.transactions == nil {
		return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusOK, "No data found", nil)
	}
	responses := response.ConvertAndSortSolPumpCompletingTransactions(
		rb.transactions,
		rb.orderBy,
		rb.direction,
	)
	return response.BuildResponse(responses, http.StatusOK, "Success", nil)
}

func (rb *ResponseBuilder) BuildSoaringResponse() response.Response {
	if rb.transactions == nil {
		return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusOK, "No data found", nil)
	}
	responses := response.ConvertAndSortSolPumpSoaringTransactions(
		rb.transactions,
		rb.orderBy,
		rb.direction,
	)
	return response.BuildResponse(responses, http.StatusOK, "Success", nil)
}

func (rb *ResponseBuilder) BuildSwapResponse() response.Response {
	if rb.transactions == nil {
		return response.BuildResponse([]response.PoolMarketInfo{}, http.StatusOK, "No data found", nil)
	}
	responses := response.ConverSolRaydiumTransactions(
		rb.transactions,
		rb.orderBy,
		rb.direction,
	)
	return response.BuildResponse(responses, http.StatusOK, "Success", nil)
}
