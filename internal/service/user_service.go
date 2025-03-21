package service

import (
	"game-fun-be/internal/auth"
	"game-fun-be/internal/constants"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"

	"fmt"
	"net/http"
	"strconv"
	"time"
)

type UserServiceImpl struct {
	userInfoRepo              *model.UserInfoRepo
	UserAuthenticationLogRepo *model.UserAuthenticationLogRepo
}

func NewUserServiceImpl(userInfoRepo *model.UserInfoRepo) *UserServiceImpl {
	return &UserServiceImpl{
		userInfoRepo: userInfoRepo,
	}
}

func (s UserServiceImpl) Login(req request.LoginRequest, chainType model.ChainType) response.Response {

	// 判断签名是否已经使用过
	isUsed, err := s.UserAuthenticationLogRepo.IsSignatureUsed(req.Address, req.Signature)
	if err != nil {
		s.insertAuthLog(req, 0, "failed to check signature")
		return response.Err(http.StatusInternalServerError, "failed to check signature", err)
	}
	if isUsed {
		s.insertAuthLog(req, 0, "signature already used")
		return response.Err(http.StatusBadRequest, "signature already used, please sign again", nil)
	}

	// 判断签名有效期是否过期
	timestampInt, err := strconv.ParseInt(req.Timestamp, 10, 64)
	if err != nil {
		s.insertAuthLog(req, 0, "invalid timestamp format")
		return response.Err(http.StatusBadRequest, "Invalid timestamp format", err)
	}
	timestamp := time.Unix(timestampInt, 0)
	if time.Since(timestamp) > 3*time.Minute {
		s.insertAuthLog(req, 0, "login timeout")
		return response.Err(http.StatusBadRequest, "Login timeout, please try again", nil)
	}

	// 验证签名
	message := fmt.Sprintf(model.LoginMessageTemplate, req.Timestamp)
	isValid, err := VerifySolanaSignature(req.Address, req.Signature, message)
	if err != nil || !isValid {
		s.insertAuthLog(req, 0, "signature verification failed")
		return response.Err(response.CodeUnauthorized, "Signature verification failed", err)
	}

	var loginResponse response.LoginResponse
	switch chainType {
	case model.ChainTypeSolana:
		loginType, userInfo, err := s.userInfoRepo.GetOrCreateUserByAddress(req.Address, uint8(chainType), req.InviteCode)
		if err != nil {
			s.insertAuthLog(req, 0, "failed to get or create user")
			return response.Err(http.StatusInternalServerError, "Failed to get or create user", err)
		}

		userIDStr := UintToString(userInfo.ID)

		userTokenKey := GetRedisKey(constants.UserTokenKeyFormat, userInfo.Address)

		token, exists, err := redis.GetToken(userTokenKey)
		if err != nil {
			s.insertAuthLog(req, 0, "failed to get token from Redis")
			return response.Err(http.StatusInternalServerError, "Failed to get token from Redis", err)
		}

		if exists {
			expireTime := time.Now().Add(model.TokenExpireDuration)
			loginResponse = buildLoginResponse(token, expireTime, userInfo, loginType)
			return response.Success(loginResponse)
		}

		token, expireTime, err := auth.GenerateJWT(userInfo.Address, userIDStr, model.TokenExpireDuration)
		if err != nil {
			s.insertAuthLog(req, 0, "failed to generate JWT")
			return response.Err(http.StatusInternalServerError, "Failed to generate JWT", err)
		}

		err = redis.Set(userTokenKey, token, model.TokenExpireDuration)
		if err != nil {
			s.insertAuthLog(req, 0, "failed to store token in Redis")
			return response.Err(http.StatusInternalServerError, "Failed to store token in Redis", err)
		}

		if err := s.insertAuthLog(req, 1, "login successful"); err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to create authentication log", err)
		}

		loginResponse = buildLoginResponse(token, expireTime, userInfo, loginType)

	default:
		s.insertAuthLog(req, 0, fmt.Sprintf("unsupported chain type: %v", chainType))
		return response.Err(response.CodeUnauthorized, fmt.Sprintf("Unsupported chain type: %v", chainType), nil)
	}

	return response.Success(loginResponse)
}

func (s UserServiceImpl) insertAuthLog(req request.LoginRequest, status int8, message string) error {
	parsedTimestamp, err := strconv.ParseInt(req.Timestamp, 10, 64)
	if err != nil {
		util.Log().Error("Error parsing timestamp: %v", err)
	}
	timestamp := time.Unix(parsedTimestamp, 0)
	authLog := &model.UserAuthenticationLog{
		Address:       req.Address,
		Message:       req.Timestamp,
		Signature:     req.Signature,
		Status:        status,
		SignatureTime: timestamp,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	if err != nil {
		authLog.Message = fmt.Sprintf("%s: %v", message, err)
	} else {
		authLog.Message = message
	}
	return s.UserAuthenticationLogRepo.CreateUserAuthenticationLog(authLog)
}

func buildLoginResponse(token string, expireTime time.Time, userInfo *model.UserInfo, loginType uint8) response.LoginResponse {
	return response.LoginResponse{
		Token:      token,
		ExpireTime: expireTime,
		User: response.UserInfo{
			UserID:          userInfo.ID,
			Address:         userInfo.Address,
			TwitterId:       userInfo.TwitterID,
			TwitterUsername: userInfo.TwitterUsername,
			InvitationCode:  userInfo.InvitationCode,
			Type:            loginType,
		},
	}
}

func (s *UserServiceImpl) MyInfo(userAddress string, chainType model.ChainType) response.Response {
	var myInfoResponse response.MyInfoResponse

	user, err := s.userInfoRepo.GetUserByAddress(userAddress, uint8(chainType))
	if err != nil {
		return response.Err(http.StatusNotFound, "User not found", err)
	}

	myInfoResponse.User = response.UserInfo{
		UserID:          user.ID,
		Address:         user.Address,
		TwitterId:       user.TwitterID,
		TwitterUsername: user.TwitterUsername,
		InvitationCode:  user.InvitationCode,
	}

	myInfoResponse.FollowerCount = 0
	myInfoResponse.FanCount = 0
	myInfoResponse.VoteCount = 0
	myInfoResponse.MentionCount = 0

	myInfoResponse.FollowStatus = "not_followed"
	myInfoResponse.InviteCode = user.InvitationCode
	myInfoResponse.InviterID = user.InviterID
	myInfoResponse.ParentInviteId = user.ParentInviteId
	myInfoResponse.HasBound = user.TwitterID != ""

	return response.Success(myInfoResponse)
}

func (s *UserServiceImpl) GetCode(userAddress string, chainType model.ChainType) response.Response {
	var inviteCodeResponse response.InviteCodeResponse

	userInfo, inviteCount, err := s.userInfoRepo.GetInviteCodeAndCount(userAddress, uint8(chainType))
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get invite code and count", err)
	}

	inviteCodeResponse.InviteCode = userInfo.InvitationCode
	inviteCodeResponse.InviteCount = inviteCount

	return response.Success(inviteCodeResponse)
}

func (s *UserServiceImpl) DeductPoints(userID uint64, amount uint64) (bool, error) {
	return s.userInfoRepo.DeductPointsWithOptimisticLock(userID, amount)
}
