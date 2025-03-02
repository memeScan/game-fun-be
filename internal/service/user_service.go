package service

import (
	"my-token-ai-be/internal/request"
	"my-token-ai-be/internal/response"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Login(req request.LoginRequest) response.Response {
	var LoginResponse response.LoginResponse

	return response.Response{
		Code: 200,
		Data: LoginResponse,
		Msg:  "success",
	}
}

func (s *UserService) MyInfo(userID string) response.Response {
	var myInfoResponse response.MyInfoResponse

	return response.Response{
		Code: 200,
		Data: myInfoResponse,
		Msg:  "success",
	}
}
