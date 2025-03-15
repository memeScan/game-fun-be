package service

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/response"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PointsServiceImpl struct {
	userInfoRepo               *model.UserInfoRepo
	pointRecordsRecord         *model.PointRecordsRepo
	platformTokenStatisticRepo *model.PlatformTokenStatisticRepo
}

type PointCalculate struct {
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	TransactionTime  time.Time `json:"transaction_time"`
	QuotaTotalAmount uint64    `json:"quota_total_amount"`
	VaultAmount      uint64    `json:"vault_amount"`
	OnlineDayCount   int       `json:"online_day"`
}

type TransactionAmountDetail struct {
	QuotaAmount     uint64
	TransactionHash string
	TransactionTime time.Time
}

type TransactionAmountDetailByTime struct {
	UserAddress              string
	QuotaTotalAmount         uint64
	VaultAmount              uint64
	TransactionAmountDetails []TransactionAmountDetail
	StartTime                time.Time
	EndTime                  time.Time
}

func NewPointsServiceImpl(userInfoRepo *model.UserInfoRepo, pointRecordsRecord *model.PointRecordsRepo, platformTokenStatisticRepo *model.PlatformTokenStatisticRepo) *PointsServiceImpl {
	return &PointsServiceImpl{userInfoRepo: userInfoRepo, pointRecordsRecord: pointRecordsRecord, platformTokenStatisticRepo: platformTokenStatisticRepo}
}

/**
* get points by user id
* @param userID uint64`json:"user_id"`
* @param chainType model.ChainType`json:"chain_type"`
 */
func (s *PointsServiceImpl) Points(userID uint, chainType model.ChainType) response.Response {
	userInfo, err := s.userInfoRepo.GetUserByUserID(userID)
	if err != nil {
		return response.DBErr("", err)
	}

	var pointsResponse response.PointsResponse = response.PointsResponse{
		TradingPoints:   formatPoints(userInfo.TradingPoints),
		InvitePoints:    formatPoints(userInfo.InvitePoints),
		AvailablePoints: formatPoints(userInfo.AvailablePoints),
	}
	return response.Success(pointsResponse)
}

/**
* get points detail by user id
 */
func (s *PointsServiceImpl) PointsDetail(userID uint64, cursor *uint, limit int, chainType model.ChainType) response.Response {
	util.Log().Info("PointsDetail: userID: %d, cursor: %d, limit: %d, chainType: %d", userID, cursor, limit, chainType)
	records, new_curs, has_more, err := s.pointRecordsRecord.GetPointRecordsByUserIDWithCursor(userID, cursor, limit)
	if err != nil {
		return response.DBErr("", err)
	}

	details := make([]response.PointsDetail, len(records))
	for i, record := range records {
		details[i] = response.PointsDetail{
			Points:    formatPoints(record.PointsChange),
			Amount:    formatPoints(record.TokenAmount),
			Timestamp: record.CreateTime.Unix(),
			Type:      record.RecordType, // 积分类型
		}
	}

	pointsDetailsResponse := response.PointsDetailsResponse{
		Details: details,
		Cursor:  new_curs,
		HasMore: has_more,
	}
	return response.Success(pointsDetailsResponse)
}

/**
*
 */
func (s *PointsServiceImpl) CreatePointRecord(wallet_address string, point uint64, hash string, transactionDetail string, record_type model.RecordType, tokenAmount uint64, nativeTokenAmount uint64, isAddPoints bool, tokenAddress string, amounts map[model.StatisticType]uint64) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		user, err := s.userInfoRepo.GetUserByAddress(wallet_address, model.ChainTypeSolana.Uint8())
		if user == nil || err != nil { // 用户不存在
			tx.Rollback()
			return err // 400 Bad Request
		}
		points := user.AvailablePoints
		// if isAddPoints {
		// 	points += point
		// } else {
		// 	points -= point
		// }

		// 创建积分记录
		insertErr := s.pointRecordsRecord.CreatePointRecord(&model.PointRecords{
			UserID:            user.ID,
			PointsChange:      point, // 积分变动
			PointsBalance:     points - point,
			RecordType:        int8(record_type), // 积分类型
			TransactionHash:   hash,
			TransactionDetail: transactionDetail,
			TokenAmount:       tokenAmount,
			NativeTokenAmount: nativeTokenAmount,
			CreateTime:        time.Now(),
			UpdateTime:        time.Now(),
		})
		if insertErr != nil {
			tx.Rollback()
			return insertErr
		}

		// userPoints := make(map[model.PointType]uint64)
		// userPoints[model.AvailablePoints] = point
		// userPoints[model.TradingPoints] = point

		// 更新统计数据
		if err := s.platformTokenStatisticRepo.WithTx(tx).IncrementStatisticsAndUpdateTime(tokenAddress, amounts); err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

func (s *PointsServiceImpl) InvitedPointsDetail(userID uint64, cursor *uint, limit int, chainType model.ChainType) response.Response {
	// 获取用户信息
	users, new_curs, has_more, err := s.userInfoRepo.GetUsersByInviterId(userID, cursor, limit)
	if err != nil {
		return response.DBErr("获取用户信息失败", err)
	}

	fmt.Println(users)

	userIDs := make([]uint, len(users))

	userInfoMap := make(map[uint]*model.UserInfo)
	for _, user := range users {
		userInfoMap[user.ID] = user
		userIDs = append(userIDs, user.ID)
	}

	// 获取邀请积分记录
	records, err := s.pointRecordsRecord.InvitedPointsDetail(userIDs)
	if err != nil {
		return response.DBErr("获取邀请积分记录失败", err)
	}

	fmt.Println(records)

	// 构建响应数据
	details := make([]response.InvitedPointsDetail, len(records))

	for i, record := range records {
		details[i] = response.InvitedPointsDetail{
			Invitee:       userInfoMap[record.UserID].Address,
			InviteTime:    userInfoMap[record.UserID].CreateTime.Unix(),
			TradingPoints: formatPoints(record.Points),
		}
	}

	pointsTotalResponse := response.InvitedPointsTotalResponse{
		Details: details,
		HasMore: has_more,
		Cursor:  new_curs,
	}

	return response.Success(pointsTotalResponse)
}

func (s *PointsServiceImpl) SavePointsEveryTimeBucket(transactionAmountDetailByTime TransactionAmountDetailByTime) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		_, user, err := s.userInfoRepo.WithTx(tx).GetOrCreateUserByAddress(transactionAmountDetailByTime.UserAddress, 1, "")
		if err != nil {
			return err // 400 Bad Request
		}

		records := make([]*model.PointRecords, len(transactionAmountDetailByTime.TransactionAmountDetails))
		userPoints := uint64(0)
		for i, detail := range transactionAmountDetailByTime.TransactionAmountDetails {

			point, onlineDayCount, err := s.CalculatePont(detail.QuotaAmount, transactionAmountDetailByTime.QuotaTotalAmount)
			if err != nil {
				return err
			}

			pointCalculate := PointCalculate{
				StartTime:        transactionAmountDetailByTime.StartTime,
				EndTime:          transactionAmountDetailByTime.EndTime,
				TransactionTime:  detail.TransactionTime,
				QuotaTotalAmount: detail.QuotaAmount,
				VaultAmount:      transactionAmountDetailByTime.VaultAmount,
				OnlineDayCount:   onlineDayCount,
			}
			data, err := json.Marshal(pointCalculate)
			if err != nil {
				data = []byte{}
			}

			// point := math.Pow(0.995/1.003, float64(transactionAmountDetailByTime.onlineDayCountonlineDayCount)) * float64(detail.QuotaAmount) / float64(transactionAmountDetailByTime.QuotaTotalAmount*7220)
			records[i] = &model.PointRecords{
				UserID:            user.ID,
				PointsChange:      uint64(point),
				RecordType:        int8(model.Trading),
				TransactionHash:   detail.TransactionHash,
				TransactionDetail: string(data),
				CreateTime:        time.Now(),
				UpdateTime:        time.Now(),
			}
			userPoints += uint64(point)
		}
		// 创建积分记录
		insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecords(records)
		if insertErr != nil {
			return insertErr
		}

		userPointsMap := make(map[model.PointType]uint64)
		userPointsMap[model.AvailablePoints] = userPoints
		userPointsMap[model.TradingPoints] = userPoints

		// 更新用户积分
		if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(transactionAmountDetailByTime.UserAddress, userPointsMap); err != nil {
			tx.Rollback()
			return err
		}

		// // 更新统计数据
		// if err := s.platformTokenStatisticRepo.WithTx(tx).IncrementStatisticsAndUpdateTime(transactionAmountDetailByTime.TokenAddress, transactionAmountDetailByTime.Amounts); err != nil {
		// 	tx.Rollback()
		// 	return err
		// }

		if user.InviterID != 0 {

			inviter, err := s.userInfoRepo.WithTx(tx).GetUserByUserID(user.InviterID)
			if inviter == nil || err != nil { // 用户不存在
				tx.Rollback()
				return err // 400 Bad Request
			}

			invitePoints := uint64(float64(userPoints) * 0.15)

			// 创建积分记录
			insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
				UserID:       inviter.ID,
				PointsChange: invitePoints,       // 积分变动
				RecordType:   int8(model.Invite), // 积分类型
				InviteeID:    user.ID,
				CreateTime:   time.Now(),
				UpdateTime:   time.Now(),
			})
			if insertErr != nil {
				tx.Rollback()

				return insertErr
			}

			userPointsMap := make(map[model.PointType]uint64)
			userPointsMap[model.AvailablePoints] = invitePoints
			userPointsMap[model.InvitePoints] = invitePoints

			// 更新用户积分
			if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(inviter.Address, userPointsMap); err != nil {
				tx.Rollback()
				return err
			}

		}

		if user.ParentInviteId != 0 {

			parentInviter, err := s.userInfoRepo.WithTx(tx).GetUserByUserID(user.ParentInviteId)
			if parentInviter == nil || err != nil { // 用户不存在
				return err // 400 Bad Request
			}

			parentInviterPoints := uint64(float64(userPoints) * 0.05)

			// 创建积分记录
			insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
				UserID:       parentInviter.ID,
				PointsChange: parentInviterPoints, // 积分变动
				RecordType:   int8(model.Invite),  // 积分类型
				InviteeID:    user.InviterID,
				CreateTime:   time.Now(),
				UpdateTime:   time.Now(),
			})
			if insertErr != nil {
				return insertErr
			}

			userPointsMap := make(map[model.PointType]uint64)
			userPointsMap[model.AvailablePoints] = parentInviterPoints
			userPointsMap[model.InvitePoints] = parentInviterPoints

			// 更新用户积分
			if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(parentInviter.Address, userPointsMap); err != nil {
				return err
			}

		}

		return nil
	})
}

func (s *PointsServiceImpl) CalculatePont(quotaAmount uint64, quotaTotalAmount uint64) (uint64, int, error) {
	onlineDate := os.Getenv("ONLINE_DATE")
	// 计算上线天数
	onlineDayCount := 1
	if onlineDate != "" {
		if onlineTime, err := time.Parse("20060102", onlineDate); err == nil {
			onlineDayCount = int(time.Now().Sub(onlineTime).Hours()/24) + 1
			if onlineDayCount < 1 {
				onlineDayCount = 1
			}
		} else {
			util.Log().Error("Failed to parse ONLINE_DATE: %v", err)
		}
	}
	point := math.Pow(0.995/1.003, float64(onlineDayCount)) * float64(quotaAmount) / float64(quotaTotalAmount*7220)
	return uint64(point), onlineDayCount, nil
}

func (s *PointsServiceImpl) PointsSave(address string, point uint64, hash string, transactionDetail string, tokenAmount uint64, baseTokenAmount uint64, tokenAddress string, amounts map[model.StatisticType]uint64) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		_, user, err := s.userInfoRepo.WithTx(tx).GetOrCreateUserByAddress(address, 1, "")
		if err != nil {
			return err // 400 Bad Request
		}

		// 创建积分记录
		insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
			UserID:            user.ID,
			PointsChange:      point, // 积分变动
			PointsBalance:     user.AvailablePoints + point,
			RecordType:        int8(model.Trading), // 积分类型
			TransactionHash:   hash,
			TransactionDetail: transactionDetail,
			TokenAmount:       tokenAmount,
			NativeTokenAmount: baseTokenAmount,
			CreateTime:        time.Now(),
			UpdateTime:        time.Now(),
		})
		if insertErr != nil {
			return insertErr
		}

		userPoints := make(map[model.PointType]uint64)
		userPoints[model.AvailablePoints] = point
		userPoints[model.TradingPoints] = point

		// 更新用户积分
		if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(address, userPoints); err != nil {
			tx.Rollback()
			return err
		}

		// 更新统计数据
		if err := s.platformTokenStatisticRepo.WithTx(tx).IncrementStatisticsAndUpdateTime(tokenAddress, amounts); err != nil {
			tx.Rollback()
			return err
		}

		if user.InviterID != 0 {

			inviter, err := s.userInfoRepo.WithTx(tx).GetUserByUserID(user.InviterID)
			if inviter == nil || err != nil { // 用户不存在
				tx.Rollback()
				return err // 400 Bad Request
			}

			invitePoints := uint64(float64(point) * 0.15)

			// 创建积分记录
			insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
				UserID:        inviter.ID,
				PointsChange:  invitePoints, // 积分变动
				PointsBalance: inviter.AvailablePoints + invitePoints,
				RecordType:    int8(model.Invite), // 积分类型
				InviteeID:     user.ID,
				CreateTime:    time.Now(),
				UpdateTime:    time.Now(),
			})
			if insertErr != nil {
				tx.Rollback()

				return insertErr
			}

			userPoints := make(map[model.PointType]uint64)
			userPoints[model.AvailablePoints] = invitePoints
			userPoints[model.InvitePoints] = invitePoints

			// 更新用户积分
			if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(inviter.Address, userPoints); err != nil {
				tx.Rollback()
				return err
			}

		}

		if user.ParentInviteId != 0 {

			parentInviter, err := s.userInfoRepo.WithTx(tx).GetUserByUserID(user.ParentInviteId)
			if parentInviter == nil || err != nil { // 用户不存在
				return err // 400 Bad Request
			}

			parentInviterPoints := uint64(float64(point) * 0.05)

			// 创建积分记录
			insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
				UserID:        parentInviter.ID,
				PointsChange:  parentInviterPoints, // 积分变动
				PointsBalance: parentInviter.AvailablePoints + parentInviterPoints,
				RecordType:    int8(model.Invite), // 积分类型
				InviteeID:     user.InviterID,
				CreateTime:    time.Now(),
				UpdateTime:    time.Now(),
			})
			if insertErr != nil {
				return insertErr
			}

			userPoints := make(map[model.PointType]uint64)
			userPoints[model.AvailablePoints] = parentInviterPoints
			userPoints[model.InvitePoints] = parentInviterPoints

			// 更新用户积分
			if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(parentInviter.Address, userPoints); err != nil {
				return err
			}

		}

		return nil
	})
}

func (s *PointsServiceImpl) PointsEstimated(userID string, chainType model.ChainType) response.Response {
	pointsEstimatedResponse := response.PointsEstimatedResponse{
		EstimatedPoints: "12862.90277",
	}
	return response.Success(pointsEstimatedResponse)
}

func formatPoints(points uint64) string {
	return fmt.Sprintf("%.6f", float64(points)/1e6)
}

func formatSol(points uint64) string {
	return fmt.Sprintf("%.9f", float64(points)/1e9)
}

func formatSolUsd(solUsdPrice decimal.Decimal) string {
	price := solUsdPrice.Shift(-9)
	// Use 6 decimal places instead if needed for crypto prices
	return price.Round(6).StringFixed(6)
}
