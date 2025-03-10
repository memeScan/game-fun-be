package services_test

import (
	"fmt"
	"game-fun-be/internal/model"
	"game-fun-be/internal/response"
	"game-fun-be/internal/service"
	"game-fun-be/test"
	"testing"

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

	// Create service
	pointsService := service.NewPointsServiceImpl(userRepo, pointsRepo)

	t.Run("get points success", func(t *testing.T) {
		// Test input
		userID := uint64(1)
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

	// Create service
	pointsService := service.NewPointsServiceImpl(userRepo, pointsRepo)

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

	// Create service
	pointsService := service.NewPointsServiceImpl(userRepo, pointsRepo)

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

func TestPointsService_PointsEstimated(t *testing.T) {
	// Create test repos
	userRepo := &model.UserInfoRepo{}
	pointsRepo := &model.PointRecordsRepo{}

	// Create service
	pointsService := service.NewPointsServiceImpl(userRepo, pointsRepo)

	t.Run("get points estimated success", func(t *testing.T) {
		// Test input
		userID := "1234"
		chainType := model.ChainType(1)

		// Call service
		resp := pointsService.PointsEstimated(userID, chainType)

		// Assert response
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)

		// Assert estimated points
		data, ok := resp.Data.(response.PointsEstimatedResponse)
		assert.True(t, ok)
		assert.Equal(t, "12862.90277", data.EstimatedPoints)
	})
}
