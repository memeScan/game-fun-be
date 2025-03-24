package services_test

import (
	"fmt"
	"testing"

	"game-fun-be/internal/cron"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/response"
	"game-fun-be/internal/service"
	"game-fun-be/test"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Use the common test setup from test package
	test.TestSetup()
	m.Run()
}

func TestPointsService_Points(t *testing.T) {
	// Create test repos
	userRepo := &model.UserInfoRepo{}
	pointsRepo := &model.PointRecordsRepo{}
	platformTokenStatisticRepo := &model.PlatformTokenStatisticRepo{}

	// Create service
	pointsService := service.NewPointsServiceImpl(userRepo, pointsRepo, platformTokenStatisticRepo)

	t.Run("get points success", func(t *testing.T) {
		// Test input
		userID := uint(1)
		chainType := model.ChainType(1)

		// Call service
		resp := pointsService.Points(userID, chainType)
		fmt.Println(resp)

		// Assert response
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
	})
}

func TestPointsService_InvitedPointsDetail(t *testing.T) {
	// Create test repos
	userRepo := &model.UserInfoRepo{}
	pointsRepo := &model.PointRecordsRepo{}

	platformTokenStatisticRepo := &model.PlatformTokenStatisticRepo{}

	// Create service
	pointsService := service.NewPointsServiceImpl(userRepo, pointsRepo, platformTokenStatisticRepo)

	t.Run("get points success", func(t *testing.T) {
		// Test input
		userID := uint64(1)
		// var cursor uint = 1
		limit := 3
		chainType := model.ChainType(1)

		// Call service
		resp := pointsService.InvitedPointsDetail(userID, nil, limit, chainType)
		fmt.Println(resp)

		// Assert response
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
	})
}

func TestPointsService_PointsDetail(t *testing.T) {
	// Create test repos
	userRepo := &model.UserInfoRepo{}
	pointsRepo := &model.PointRecordsRepo{}

	platformTokenStatisticRepo := &model.PlatformTokenStatisticRepo{}

	// Create service
	pointsService := service.NewPointsServiceImpl(userRepo, pointsRepo, platformTokenStatisticRepo)

	t.Run("get points detail success", func(t *testing.T) {
		// Test input
		userID := uint64(1)
		var cursor uint = 1
		limit := 2
		chainType := model.ChainType(1)

		// Call service
		resp := pointsService.PointsDetail(userID, &cursor, limit, chainType)
		fmt.Println(resp)

		// Assert response
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
	})
}

func TestPointsService_PointsSave(t *testing.T) {
	// globalService := service.NewGlobalServiceImpl()

	t.Run("save points success", func(t *testing.T) {
		// Test input
		// vaultAddress := os.Getenv("VAULT_ADDRESS")
		// tokenAddress := os.Getenv("TOKEN_ADDRESS")
		// resp := globalService.TickerBalance(vaultAddress, tokenAddress, model.ChainTypeSolana)
		// t.Logf("resp: %v", resp)
		// util.Log().Info("resp: %v", resp)
		// fmt.Println(resp)
		cron.ExecutePointJob()
	})
}

func TestPointsService_PlatformTokenQuery(t *testing.T) {
	// Create test repos
	platformTokenStatisticRepo := &model.PlatformTokenStatisticRepo{}

	// Create service
	platformTokenStatisticServiceImpl := service.NewPlatformTokenStatisticServiceImpl(platformTokenStatisticRepo)

	t.Run("get estimated points success", func(t *testing.T) {
		// Test input
		token_address := "8iFREvVdmLKxVeibpC5VLRr1S6X5dm7gYR3VCU1wpump"
		chainType := model.ChainType(1)

		// Call service
		resp := platformTokenStatisticServiceImpl.GetTokenPointsStatistic(token_address, uint8(chainType))
		fmt.Println(resp)

		// Assert response
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)

		// Assert response data
		data, ok := resp.Data.(response.PointsEstimatedResponse)
		assert.True(t, ok)
		assert.NotEmpty(t, data.EstimatedPoints)
	})
}

func Test_GetTokenInfoByAddress(t *testing.T) {
	tokenAddress := "8iFREvVdmLKxVeibpC5VLRr1S6X5dm7gYR3VCU1wpump"
	resp, err := httpUtil.GetTokenInfoByAddress(tokenAddress)
	if err != nil {
		t.Fatalf("failed to get token info by address: %v", err)
	}
	fmt.Println(resp)
}
