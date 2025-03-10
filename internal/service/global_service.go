package service

import (
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/response"
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

func (s *GlobalServiceImpl) NativeBalance(userAddress string, chainType model.ChainType) response.Response {
	tokenBalances, err := httpUtil.GetTokenBalance([]string{userAddress}, model.SolanaWrappedSOLAddress)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get balance", err)
	}
	balances := response.TokenBalance{
		Token:    model.SolanaWrappedSOLAddress,
		Owner:    userAddress,
		Balance:  (*tokenBalances)[0].Balance,
		Decimals: (*tokenBalances)[0].Decimals,
	}
	return response.Success(balances)
}

func (s *GlobalServiceImpl) TickerBalance(userAddress string, tokenAddress string, chainType model.ChainType) response.Response {
	tokenBalances, err := httpUtil.GetTokenBalance([]string{userAddress}, tokenAddress)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get balance", err)
	}
	balances := response.TokenBalance{
		Token:    tokenAddress,
		Owner:    userAddress,
		Balance:  (*tokenBalances)[0].Balance,
		Decimals: (*tokenBalances)[0].Decimals,
	}
	return response.Success(balances)
}
