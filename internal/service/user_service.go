package service

import (
	"game-fun-be/internal/auth"
	"game-fun-be/internal/constants"
	"game-fun-be/internal/model"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"

	"fmt"
	"net/http"
	"time"
)

type UserServiceImpl struct {
	userInfoRepo *model.UserInfoRepo
}

func NewUserServiceImpl(userInfoRepo *model.UserInfoRepo) *UserServiceImpl {
	return &UserServiceImpl{
		userInfoRepo: userInfoRepo,
	}
}

func (s UserServiceImpl) Login(req request.LoginRequest, chainType model.ChainType) response.Response {
	var loginResponse response.LoginResponse
	switch chainType {
	case model.ChainTypeSolana:

		userInfo, err := s.userInfoRepo.GetOrCreateUserByAddress(req.Address, uint8(chainType), req.InviteCode)
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to get or create user", err)
		}

		userIDStr := UintToString(userInfo.ID)

		userTokenKey := GetRedisKey(constants.UserTokenKeyFormat, userInfo.Address)

		token, exists, err := redis.GetToken(userTokenKey)
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to get token from Redis", err)
		}

		if exists {
			expireTime := time.Now().Add(model.TokenExpireDuration)
			loginResponse = buildLoginResponse(token, expireTime, userInfo)
			return response.Success(loginResponse)
		}

		message := fmt.Sprintf(model.LoginMessageTemplate, req.Timestamp)

		isValid, err := VerifySolanaSignature(req.Address, req.Signature, message)
		if err != nil || !isValid {
			return response.Err(response.CodeUnauthorized, "Signature verification failed", err)
		}

		token, expireTime, err := auth.GenerateJWT(userInfo.Address, userIDStr, model.TokenExpireDuration)
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to generate JWT", err)
		}

		err = redis.Set(userTokenKey, token, model.TokenExpireDuration)
		if err != nil {
			return response.Err(http.StatusInternalServerError, "Failed to store token in Redis", err)
		}

		loginResponse = buildLoginResponse(token, expireTime, userInfo)

	default:
		return response.Err(response.CodeUnauthorized, fmt.Sprintf("Unsupported chain type: %v", chainType), nil)
	}

	return response.Success(loginResponse)
}

func buildLoginResponse(token string, expireTime time.Time, userInfo *model.UserInfo) response.LoginResponse {
	return response.LoginResponse{
		Token:      token,
		ExpireTime: expireTime,
		User: response.UserInfo{
			UserID:          userInfo.ID,
			Address:         userInfo.Address,
			TwitterId:       userInfo.TwitterID,
			TwitterUsername: userInfo.TwitterUsername,
			InvitationCode:  userInfo.InvitationCode,
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
