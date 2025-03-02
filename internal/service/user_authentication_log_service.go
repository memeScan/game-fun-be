package service

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/response"
	"net/http"
	"time"
)


// CreateAuthLog 创建认证日志
func CreateAuthLog(address, messageNonce string) (*model.UserAuthenticationLog, error) {
	authLog := &model.UserAuthenticationLog{
		Address:      address,
		MessageNonce: messageNonce,
		Status:       0, // 未签名状态
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	
	err := model.CreateUserAuthenticationLog(authLog)
	if err != nil {
		return nil, err
	}
	
	return authLog, nil
}

func UpdateAuthLog(address string, messageNonce string, signature string, status int8) (*model.UserAuthenticationLog, error) {

    authLog, err := model.GetUserAuthenticationLogByAddress(address, messageNonce)
    if err != nil {
        return nil, err
    }

    authLog.Signature = signature
    authLog.Status = status
    authLog.UpdateTime = time.Now()

    err = model.UpdateUserAuthenticationLog(authLog)
    if err != nil {
        return nil, err
    }

    return authLog, nil
}

// GetLatestAuthLog 获取最新的认证日志
func GetLatestAuthLog(id uint) (*model.UserAuthenticationLog, error) {
	return model.GetUserAuthenticationLogID(id)
}

// ProcessAuthLogCreation 处理认证日志创建
func ProcessAuthLogCreation(address, messageNonce string) response.Response {
	authLog, err := CreateAuthLog(address, messageNonce)
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to create authentication log", err)
	}
	
	return response.Response{
		Code: 0,
		Data: authLog,
		Msg:  "Authentication log created successfully",
	}
}

// ProcessAuthLogUpdate 处理认证日志更新
func ProcessAuthLogUpdate(address string, messageNonce string, signature string, status int8) error {
	_, err := UpdateAuthLog(address, messageNonce, signature, status)
	if err != nil {
		return err
	}
	
	return nil
}
