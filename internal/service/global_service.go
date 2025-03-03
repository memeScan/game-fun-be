package service

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/pkg/httpUtil"
	"my-token-ai-be/internal/response"
	"net/http"
)

// GlobalService 定义接口
type GlobalService interface {
	SolUsdPrice() response.Response
	SolBalance(address string) response.Response
}

// GlobalServiceImpl 实现接口
type GlobalServiceImpl struct{}

// NewGlobalServiceService 创建服务实例
func NewGlobalServiceImpl() GlobalService {
	return &GlobalServiceImpl{}
}

// SolUsdPrice 获取 SOL 对 USD 的价格
func (s *GlobalServiceImpl) SolUsdPrice() response.Response {
	solUsdPrice, err := getSolPrice()
	if err != nil {
		return response.Err(http.StatusInternalServerError, "failed to get sol price", err)
	}
	responseData := map[string]string{
		response.TokenPrices: solUsdPrice.StringFixed(8),
	}
	return response.Success(responseData)
}

// SolBalance 获取 SOL 余额
func (s *GlobalServiceImpl) SolBalance(address string) response.Response {
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
