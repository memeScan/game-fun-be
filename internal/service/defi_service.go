package service

import (
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/response"
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
