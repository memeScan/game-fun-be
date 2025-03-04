package service

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/pkg/httpUtil"
	"my-token-ai-be/internal/response"
	"net/http"
)

type GlobalServiceImpl struct{}

func NewGlobalServiceImpl() *GlobalServiceImpl {
	return &GlobalServiceImpl{}
}

func (s *GlobalServiceImpl) NativeTokePrice(chainType model.ChainType) response.Response {
	solUsdPrice, err := getSolPrice()
	if err != nil {
		return response.Err(http.StatusInternalServerError, "failed to get sol price", err)
	}
	responseData := map[string]string{
		response.TokenPrices: solUsdPrice.StringFixed(8),
	}
	return response.Success(responseData)
}

func (s *GlobalServiceImpl) Balance(address string, chainType model.ChainType) response.Response {
	tokenBalances, err := httpUtil.GetTokenBalance([]string{address}, model.SolanaWrappedSOLAddress)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get balance", err)
	}
	balances := response.TokenBalance{
		Token:    address,
		Owner:    model.SolanaWrappedSOLAddress,
		Balance:  (*tokenBalances)[0].Balance,
		Decimals: (*tokenBalances)[0].Decimals,
	}
	return response.Success(balances)
}
