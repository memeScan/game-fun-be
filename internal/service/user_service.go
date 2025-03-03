package service

import (
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
)

type UserService interface {
	Login(req request.LoginRequest) response.Response
	MyInfo(userID string) response.Response
	GetCode(userID string) response.Response
}

type UserServiceImpl struct{}

func NewUserServiceImpl() UserService {
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
