package service

import (
	"fmt"
	"net/http"

	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/response"
)

type TokenHoldingsServiceImpl struct {
	// tokenInfoService *TokenInfoService
}

func NewTokenHoldingsServiceImpl() *TokenHoldingsServiceImpl {
	return &TokenHoldingsServiceImpl{}
}

func (s *TokenHoldingsServiceImpl) TokenHoldings(userAccount string, chainType model.ChainType) response.Response {
	tokenInfoService := &TokenInfoService{}

	walletPNL, err := httpUtil.GetWalletPNL(userAccount)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "failed to get wallet pnl", err)
	}

	tokenAddresses := make([]string, 0)
	for tokenAddress := range walletPNL.Tokens {
		tokenAddresses = append(tokenAddresses, tokenAddress)
	}

	tokenInfoMap, err := tokenInfoService.GetTokenInfoMapByDB(tokenAddresses, uint8(chainType))
	if err != nil {
		return response.Err(http.StatusInternalServerError, "failed to get token info map", err)
	}

	// fmt.Println(tokenInfoMap)

	currentHolding := make([]response.TokenHolding, 0)
	historyHolding := make([]response.TokenHolding, 0)

	for tokenAddress, tokenPNL := range walletPNL.Tokens {
		tokenInfo, ok := tokenInfoMap[tokenAddress]
		if !ok {
			continue
		}
		if tokenPNL.Sold-tokenPNL.Held >= 0 && tokenPNL.Holding == 0 {
			historyHolding = append(historyHolding, response.TokenHolding{
				TokenName:    tokenInfo.TokenName,
				Symbol:       tokenInfo.Symbol,
				Price:        tokenInfo.Price.String(),
				ImageURI:     tokenInfo.URI,
				Balance:      "0",
				TotalValue:   "",
				ID:           tokenInfo.ID,
				HoldersCount: tokenInfo.Holder,
				Profit:       fmt.Sprintf("%f", tokenPNL.Total),
				ProfitRate:   fmt.Sprintf("%f", tokenPNL.Total*100/tokenPNL.TotalInvested),
			})
		} else {
			currentHolding = append(currentHolding, response.TokenHolding{
				TokenName:    tokenInfo.TokenName,
				Symbol:       tokenInfo.Symbol,
				Price:        tokenInfo.Price.String(),
				ImageURI:     tokenInfo.URI,
				Balance:      fmt.Sprintf("%f", tokenPNL.Holding),
				TotalValue:   fmt.Sprintf("%f", tokenPNL.CurrentValue),
				ID:           tokenInfo.ID,
				HoldersCount: tokenInfo.Holder,
				Profit:       fmt.Sprintf("%f", tokenPNL.Total),
				ProfitRate:   fmt.Sprintf("%f", tokenPNL.Total*100/tokenPNL.TotalInvested),
			})
		}
	}

	var tokenHoldingsResponse response.TokenHoldingsResponse
	tokenHoldingsResponse.CurrentHolding = currentHolding
	tokenHoldingsResponse.HistoryTokenHoldings = historyHolding
	return response.Success(tokenHoldingsResponse)
}
