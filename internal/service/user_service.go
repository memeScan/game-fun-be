package service

import (
	"game-fun-be/internal/request"
	"game-fun-be/internal/response"
)

type UserServiceImpl struct{}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{}
}

func (s *UserServiceImpl) Login(req request.LoginRequest) response.Response {
	var loginResponse response.LoginResponse
	return response.Success(loginResponse)
}

func (s *UserServiceImpl) MyInfo(userID string) response.Response {
	var myInfoResponse response.MyInfoResponse
	return response.Success(myInfoResponse)
}

func (s *UserServiceImpl) GetCode(userID string) response.Response {
	var inviteCodeResponse response.InviteCodeResponse
	return response.Success(inviteCodeResponse)
}
