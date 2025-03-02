package service

import (
	"encoding/json"
	"fmt"
	"my-token-ai-be/internal/auth"
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/redis"
	"my-token-ai-be/internal/response"
	"net/http"
	"regexp"
	"strings"
	"time"
	"my-token-ai-be/internal/request"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)


// Login 钱包登录		
func Login(address string, signature string) response.Response {

	if address == "" || signature == "" {
		return response.ParamErr("地址和签名不能为空", nil)
	}

	// Retrieve stored message from Redis
	key := fmt.Sprintf("%s%s", response.RedisKeyPrefixMessage, address)
	messageJSON, err := redis.GetAndDelete(key)
	if err != nil {
		util.Log().Error("Failed to retrieve message from Redis: " + err.Error())
		return response.Err(http.StatusInternalServerError, "Failed to retrieve message from Redis", err)
	}
		
	var messageMap map[string]interface{}

	// 解码 JSON 字符串
	err = json.Unmarshal([]byte(messageJSON), &messageMap)
	if err != nil {
		util.Log().Error("Error unmarshaling messageJSON: " + err.Error())
		return response.Err(http.StatusInternalServerError, "Failed to decode messageJSON", err)
	}

	if !VerifyExistingSignature(address, signature, messageJSON) {
		ProcessAuthLogUpdate(address,messageMap["nonce"].(string),signature, 2)
		return response.Err(http.StatusBadRequest, "签名验证失败", nil)
	}
	
	return processUserLogin(address, messageMap["nonce"].(string), signature)
}

// processUserLogin 处理用户登录逻辑
func processUserLogin(address string, messageNonce string, signature string) response.Response {
	// 获取或创建用户
	user, err := model.GetOrCreateUserByAddress(address)
	if err != nil {
		// 更新认证日志状态为失败
		ProcessAuthLogUpdate(address, messageNonce, signature, 2)
		return response.DBErr("用户创建或获取失败", err)
	}

	// 更新用户授权登录记录
	updateResponse := ProcessAuthLogUpdate(address, messageNonce, signature, 1)
	if updateResponse != nil {
		return response.Err(http.StatusInternalServerError, "更新用户授权登录记录失败", updateResponse)
	}

	// 生成并缓存 token
	token, err := auth.GenerateAndCacheJWT(user.Address)
	
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Token生成或缓存失败", err)
	}

	return response.Response{
		Code: http.StatusOK	,
		Data: map[string]interface{}{
			"token":   token,
			"address": user.Address,
		},
		Msg: "Login successful",
	}
}

// GetUserByAddress 根据地址获取用户信息
func GetUserByAddress(address string) (response.Response) {
	user, err := model.GetUserByAddress(address)
	if err != nil {
		return response.BuildResponse(nil, http.StatusInternalServerError, "获取用户信息失败", err)
	}
	return response.BuildResponse(user, http.StatusOK, "Success",nil)
}

// GetMessage 获取需要签名的消息
func GetMessage(address string) (response.Response) {
	isValid := isValidSolanaAddress(address)
	if !isValid {
		return response.Err(http.StatusInternalServerError, "Invalid Solana address", nil)
	}				
	nonce := util.RandStringRunes(8)
	now := time.Now().UTC()
	expirationTime := now.Add(30 * 24 * time.Hour) // 30天后过期
	timeFormat := "2006-01-02T15:04:05.000Z"
	message := request.Message{
		Domain:         "MytokenAI",
		Statement:      "MytokenAI wants you to sign in with your Solana account:" + address,
		URI:            "https://MytokenAI",
		Version:        "1",
		ChainID:        900,
		Nonce:          nonce,
		IssuedAt:       now.Format(timeFormat),
		ExpirationTime: expirationTime.Format(timeFormat),
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		util.Log().Error("Failed to marshal message: " + err.Error())
		return response.Err(http.StatusInternalServerError, "Failed to marshal message", err)
	}

	// 创建用户授权登录记录
	authLogResponse := ProcessAuthLogCreation(address, nonce)
	if authLogResponse.Code != 0 {
		return authLogResponse
	}

	// 将整个 message 存储到 Redis
	key := fmt.Sprintf("%s%s", response.RedisKeyPrefixMessage, address)
	err = redis.Set(key, string(messageJSON), 2*time.Hour)
	if err != nil {
		util.Log().Error("Failed to store message: " + err.Error())
		return response.Err(http.StatusInternalServerError, "Failed to store message", err)
	}

	return response.BuildResponse(message, http.StatusOK, "Success", nil)
}

// isValidSolanaAddress 验证是否为solana地址
func isValidSolanaAddress(address string) bool {
    // 检查地址是否为空
    if strings.TrimSpace(address) == "" {
        return false
    }

    // 检查地址长度，Solana 地址长度通常在 32 到 44 字符之间
    if len(address) < 32 || len(address) > 44 {
        return false
    }

    // 检查地址是否符合 Base58 编码格式
    base58Regex := regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]+$`)
    if !base58Regex.MatchString(address) {
        return false
    }

    // 尝试解码以确保它是有效的 Base58 编码
    decoded, err := base58.Decode(address)
    if err != nil {
        return false
    }

    // 检查解码后的字节长度，Solana 公钥应为 32 字节
    return len(decoded) == 32
}

// VerifyExistingSignature 验证签名
func VerifyExistingSignature(address string, signature string, decodedString string) bool {
	// Decode public key from Base58
	pubkey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		util.Log().Error("Invalid public key: " + err.Error())
		return false
	}

	// Decode signature from Base58
	sigBytes, err := base58.Decode(signature)
	if err != nil {
		util.Log().Error("Failed to decode signature: " + err.Error())
		return false
	}

	// Create signature object
	var sig solana.Signature
	copy(sig[:], sigBytes)

	// Verify signature
	isValid := pubkey.Verify([]byte(decodedString), sig)
	if !isValid {
		util.Log().Error("Signature verification failed")
	}

	return isValid
}
