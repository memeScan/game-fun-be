package service

import (
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
)

type UserService interface {
	Login(req request.LoginRequest) response.Response
	MyInfo(userID string) response.Response
}

type UserServiceImpl struct{}

func NewUserService() UserService {
	return &UserServiceImpl{}
}

func (s *UserServiceImpl) Login(req request.LoginRequest) response.Response {
	var loginResponse response.LoginResponse

	return response.Response{
		Code: 200,
		Data: loginResponse,
		Msg:  "success",
	}
}

func (s *UserServiceImpl) MyInfo(userID string) response.Response {
	var myInfoResponse response.MyInfoResponse
	return response.Response{
		Code: 200,
		Data: myInfoResponse,
		Msg:  "success",
	}
}
