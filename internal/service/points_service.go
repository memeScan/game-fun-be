package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"game-fun-be/internal/conf"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/response"

	"github.com/IBM/sarama"
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

func NewPointsServiceImpl(userInfoRepo *model.UserInfoRepo,
	pointRecordsRecord *model.PointRecordsRepo,
	platformTokenStatisticRepo *model.PlatformTokenStatisticRepo,
) *PointsServiceImpl {
	return &PointsServiceImpl{
		userInfoRepo:               userInfoRepo,
		pointRecordsRecord:         pointRecordsRecord,
		platformTokenStatisticRepo: platformTokenStatisticRepo,
	}
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

	solPriceUSD, priceErr := getSolPrice()
	if priceErr != nil {
		return response.Err(http.StatusBadRequest, "price query failed", priceErr)
	}

	var pointsResponse response.PointsResponse = response.PointsResponse{
		TradingPoints:      formatPoints(userInfo.TradingPoints),
		InvitePoints:       formatPoints(userInfo.InvitePoints),
		AvailablePoints:    formatPoints(userInfo.AvailablePoints),
		AccumulatedPoints:  formatPoints(userInfo.TradingPoints + userInfo.InvitePoints),
		InviteRebate:       formatSolUsd(solPriceUSD.Mul(decimal.NewFromInt(int64(userInfo.InviteRebate)))),
		WithdrawableRebate: formatSolUsd(solPriceUSD.Mul(decimal.NewFromInt(int64(userInfo.WithdrawableRebate)))),
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

	solPriceUSD, priceErr := getSolPrice()
	if priceErr != nil {
		return response.Err(http.StatusBadRequest, "price query failed", priceErr)
	}

	// 构建响应数据
	details := make([]response.InvitedPointsDetail, len(records))

	for i, record := range records {
		details[i] = response.InvitedPointsDetail{
			Invitee:       userInfoMap[record.UserID].Address,
			InviteTime:    userInfoMap[record.UserID].CreateTime.Unix(),
			TradingPoints: formatPoints(record.Points),
			FeeRebate:     formatSolUsd(solPriceUSD.Mul(decimal.NewFromInt(int64(record.Rebate)))),
			UpdateTime:    time.Now().Unix(),
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

			point, onlineDayCount, err := s.CalculatePoint(transactionAmountDetailByTime.VaultAmount, detail.QuotaAmount, transactionAmountDetailByTime.QuotaTotalAmount)
			if err != nil {
				return err
			}

			pointCalculate := PointCalculate{
				StartTime:        transactionAmountDetailByTime.StartTime,
				EndTime:          transactionAmountDetailByTime.EndTime,
				TransactionTime:  detail.TransactionTime,
				QuotaTotalAmount: transactionAmountDetailByTime.QuotaTotalAmount,
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
				TokenAmount:       detail.QuotaAmount,
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

func (s *PointsServiceImpl) CalculatePointByDay(vaultAmount uint64) (float64, int, error) {
	onlineDate := os.Getenv("ONLINE_DATE")
	newCoeffientStr := os.Getenv("NEW_COEFFICIENT")

	newCoeffient, _ := strconv.ParseFloat(newCoeffientStr, 64)
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
	point := math.Pow(0.995/1.003, float64(onlineDayCount)) * float64(vaultAmount) / newCoeffient
	return point, onlineDayCount, nil
}

func (s *PointsServiceImpl) CalculatePoint(vaultAmount uint64, quotaAmount uint64, quotaTotalAmount uint64) (uint64, int, error) {
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
	pointByDay, _, _ := s.CalculatePointByDay(vaultAmount)
	point := pointByDay * float64(quotaAmount) / float64(quotaTotalAmount)
	return uint64(point), onlineDayCount, nil
}

func (s *PointsServiceImpl) PointsSave(address string, point uint64, hash string, transactionDetail string, tokenAmount uint64, baseTokenAmount uint64, tokenAddress string, amounts map[model.StatisticType]uint64) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		_, user, err := s.userInfoRepo.WithTx(tx).GetOrCreateUserByAddress(address, 1, "")
		if err != nil {
			return err // 400 Bad Request
		}

		baseFee := amounts[model.FeeAmount]

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
			inviteRebate := uint64(float64(baseFee) * 0.15)

			// 创建积分记录
			insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
				UserID:        inviter.ID,
				PointsChange:  invitePoints, // 积分变动
				PointsBalance: inviter.AvailablePoints + invitePoints,
				RecordType:    int8(model.Invite), // 积分类型
				RebateChange:  inviteRebate,
				RebateBalance: inviter.WithdrawableRebate + inviteRebate,
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
			userPoints[model.InviteRebate] = inviteRebate
			userPoints[model.WithdrawableRebate] = inviteRebate
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
			parentInviterRebate := uint64(float64(baseFee) * 0.05)
			// 创建积分记录
			insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
				UserID:        parentInviter.ID,
				PointsChange:  parentInviterPoints, // 积分变动
				PointsBalance: parentInviter.AvailablePoints + parentInviterPoints,
				RecordType:    int8(model.Invite), // 积分类型
				RebateChange:  parentInviterRebate,
				RebateBalance: parentInviter.WithdrawableRebate + parentInviterRebate,
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
			userPoints[model.InviteRebate] = parentInviterRebate
			userPoints[model.WithdrawableRebate] = parentInviterRebate

			// 更新用户积分
			if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(parentInviter.Address, userPoints); err != nil {
				return err
			}

		}

		return nil
	})
}

func (s *PointsServiceImpl) PointsEstimated(userID string, vaultAmount uint64, chainType model.ChainType) response.Response {
	points, _, _ := s.CalculatePointByDay(vaultAmount)
	pointsEstimatedResponse := response.PointsEstimatedResponse{
		EstimatedPoints: formatPoints(uint64(points)),
	}
	return response.Success(pointsEstimatedResponse)
}

func (s *PointsServiceImpl) IncrementStatisticsAndUpdateTime(address string, amounts map[model.StatisticType]uint64) error {
	return s.platformTokenStatisticRepo.IncrementStatisticsAndUpdateTime(address, amounts)
}

func (s *PointsServiceImpl) CheckRebate(address string, rebateAmount uint64) response.Response {
	user, err := s.userInfoRepo.GetUserByAddress(address, model.ChainTypeSolana.Uint8())
	if err != nil {
		return response.Err(http.StatusBadRequest, "用户不存在", err)
	}

	if user.WithdrawableRebate < rebateAmount {
		return response.Err(http.StatusBadRequest, "提现金额不足", errors.New("提现金额不足"))
	}

	return response.Success("有足够的提现金额")
}

func (s *PointsServiceImpl) SendClaimTransaction(address string) response.Response {
	user, err := s.userInfoRepo.GetUserByAddress(address, model.ChainTypeSolana.Uint8())
	if err != nil {
		return response.Err(http.StatusBadRequest, "用户不存在", err)
	}

	isTrue, err := s.userInfoRepo.DeductRebateWithOptimisticLock(uint64(user.ID), user.WithdrawableRebate)
	if err != nil {
		util.Log().Error("Failed to deduct points with optimistic lock: %v", err)
		return response.Err(http.StatusInternalServerError, "Failed to deduct points, please try again later", err)
	}
	if !isTrue {
		util.Log().Error("Optimistic lock failed, points deduction unsuccessful for user: %d", user.ID)
		return response.Err(http.StatusConflict, "Points deduction failed due to concurrent update, please try again", nil)
	}

	resp, err := httpUtil.SendClaimTransaction(address, user.WithdrawableRebate)
	if err != nil || resp == nil || resp.Code != 2000 {

		// 交易发送失败，恢复用户积分
		userInfoRepo := model.NewUserInfoRepo()
		if err := userInfoRepo.IncrementWithdrawableRebateByUserID(uint(user.ID), user.WithdrawableRebate); err != nil {
			util.Log().Error("Failed to restore points for user %d after transaction %s failed: %v",
				user.ID, resp.Data.Signature, err)
			return response.Err(http.StatusInternalServerError, "Failed to restore points, please try again later", err)
		}
		util.Log().Info("Transaction %s failed, restored %d points to user %d",
			resp.Data.Signature, user.WithdrawableRebate, user.ID)
		return response.Err(http.StatusInternalServerError, "Failed to get send transaction", err)
	}
	pointTxStatusMsg := model.PointTxStatusMessage{
		Signature: resp.Data.Signature,
		UserId:    uint(user.ID),
		Points:    0,
		Rebate:    user.WithdrawableRebate,
		TxType:    2,
	}

	msgBytes, err := json.Marshal(pointTxStatusMsg)
	if err != nil {
		util.Log().Error("Failed to marshal point transaction status message: %v", err)
	} else {
		// 发送消息到Kafka
		var topic string
		if conf.IsTest() {
			topic = "market.point.tx.status.test"
		} else {
			topic = "market.point.tx.status.prod"
		}

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(msgBytes),
		}
		fmt.Println(msg)
		// producer
		// _, _, err := &sarama.SendMessage(msg)
		// if err != nil {
		// 	util.Log().Error("Failed to send point transaction status message to Kafka: %v", err)
		// } else {
		// 	util.Log().Info("Sent point transaction status check message for transaction %s, user %d, points %d",
		// 		resp.Data.Signature, user.ID, user.WithdrawableRebate)
		// }
	}

	return response.Success(resp.Data)
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
