package service

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
)

type UserServiceImpl struct{}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{}
}

func (s *UserServiceImpl) Login(req request.LoginRequest, chainType model.ChainType) response.Response {
	var loginResponse response.LoginResponse
	return response.Success(loginResponse)
}

func (s *UserServiceImpl) MyInfo(userID string, chainType model.ChainType) response.Response {
	var myInfoResponse response.MyInfoResponse
	return response.Success(myInfoResponse)
}

func (s *UserServiceImpl) GetCode(userID string, chainType model.ChainType) response.Response {
	var inviteCodeResponse response.InviteCodeResponse
	return response.Success(inviteCodeResponse)
}
