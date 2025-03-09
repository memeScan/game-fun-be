package httpUtil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game-fun-be/internal/pkg/httpRespone"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// SOLANA 常量
const SOLANA = "solana"

// 全局变量，用于存储配置
var BaseURL = "https://public-api.birdeye.so/defi"

// sendRequest 统一请求封装，支持 GET 和 POST，支持传递链参数
func sendRequest(method, url string, bodyData interface{}, target interface{}, chain string) error {
	var body io.Reader

	// 如果是 POST 请求，处理 bodyData
	if method == http.MethodPost && bodyData != nil {
		jsonBody, err := json.Marshal(bodyData)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON body: %v", err)
		}
		body = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-API-KEY", os.Getenv("BIRDEYE_API_KEY"))
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-chain", chain)

	// 只有 POST 请求需要设置 Content-Type
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	// 使用全局 httpClient 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 处理非 200 响应
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// 解析 JSON 响应
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	return nil
}

// GetTokenHolders 获取代币持有者信息
func GetTokenHolders(tokenAddress string, offset, limit int, chain string) (*httpRespone.TokenHoldersResponse, error) {
	url := fmt.Sprintf("%s/v3/token/holder?address=%s&offset=%d&limit=%d", BaseURL, tokenAddress, offset, limit)
	var response httpRespone.TokenHoldersResponse
	if err := sendRequest(http.MethodGet, url, nil, &response, chain); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTokenMarketData 获取代币市场数据
func GetTokenMarketData(tokenAddress string, chain string) (*httpRespone.TokenMarketDataResponse, error) {
	url := fmt.Sprintf("%s/v3/token/market-data?address=%s", BaseURL, tokenAddress)
	var response httpRespone.TokenMarketDataResponse
	if err := sendRequest(http.MethodGet, url, nil, &response, chain); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTradeData 获取代币交易数据
func GetTradeData(tokenAddress string, chain string) (*httpRespone.TradeDataResponse, error) {
	url := fmt.Sprintf("%s/v3/token/trade-data/single?address=%s", BaseURL, tokenAddress)
	var response httpRespone.TradeDataResponse
	if err := sendRequest(http.MethodGet, url, nil, &response, chain); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTokenMetaData 获取多个代币的元数据
func GetTokenMetaData(tokenAddresses []string, chain string) (*httpRespone.TokenMetaDataResponse, error) {
	// 将 tokenAddresses 转换为逗号分隔的字符串
	addressList := strings.Join(tokenAddresses, "%2C")

	// 构建请求 URL
	url := fmt.Sprintf("%s/v3/token/meta-data/multiple?list_address=%s", BaseURL, addressList)
	var response httpRespone.TokenMetaDataResponse

	if err := sendRequest(http.MethodGet, url, nil, &response, chain); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTokenCreationInfo 获取代币创建信息
func GetTokenCreationInfo(tokenAddress string, chain string) (*httpRespone.TokenCreationInfoResponse, error) {
	url := fmt.Sprintf("%s/token_creation_info?address=%s", BaseURL, tokenAddress)
	log.Print(url)
	var response httpRespone.TokenCreationInfoResponse
	if err := sendRequest(http.MethodGet, url, nil, &response, chain); err != nil {
		return nil, err
	}
	return &response, nil
}
