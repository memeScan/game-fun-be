package service

import (
	"my-token-ai-be/internal/pkg/httpUtil"
	"my-token-ai-be/internal/response"
	"net/http"
	"os"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
)

type OnChainDataService interface {
	GetNativeTokenGasFee(chainType string) response.Response
	GetSolPrice() (decimal.Decimal, error)
	GetTokenBalance(chainType, owner, token string) (response.TokenBalance, error)
}

type onChainDataService struct {
	client *rpc.Client
}

func NewOnChainDataService() *onChainDataService {
	return &onChainDataService{
		client: rpc.New(os.Getenv("SOLANA_RPC_URL")),
	}
}

func (s *onChainDataService) GetNativeTokenGasFee(chainType string) response.Response {
	gasFee, err := httpUtil.GetPriorityFee()
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "failed to get gas fee", err)
	}
	return response.BuildResponse(gasFee.PriorityFeeLevels, http.StatusOK, "success", nil)
}

func (s *onChainDataService) GetSolPrice() response.Response {
	solPrice, err := getSolPrice()
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "failed to get sol price", err)
	}
	return response.BuildResponse(solPrice, http.StatusOK, "success", nil)
}

func (s *onChainDataService) GetTokenBalance(chainType, owner, token string) response.Response {

	tokenBalances, err := httpUtil.GetTokenBalance([]string{owner}, token)
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "failed to get balance", err)
	}

	balances := response.TokenBalance{
		Token:    token,
		Owner:    owner,
		Balance:  (*tokenBalances)[0].Balance,
		Decimals: (*tokenBalances)[0].Decimals,
	}

	return response.BuildResponse(balances, http.StatusOK, "success", nil)
}
