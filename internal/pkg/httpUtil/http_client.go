package httpUtil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"game-fun-be/internal/pkg/httpRequest"
	"game-fun-be/internal/pkg/httpRespone"
	"game-fun-be/internal/pkg/util"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"unicode/utf8"

	"game-fun-be/internal/pkg/metrics"

	"github.com/redis/go-redis/v9"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/unicode/norm"
)

var (
	// API 端点
	ApiTokenBalanceBatch      string
	ApiDexCheckBatch          string
	ApiSafetyCheckBatch       string
	ApiPriorityFee            string
	ApiPumpfunTrade           string
	ApiRaydiumTrade           string
	ApiTransactionSend        string
	ApiGameFunTransactionSend string
	ApiTransactionStatus      string
	ApiTokensBatch            string
	ApiSafetyCheckPool        string
	ApiTipFloor               string
	ApiPoolInfo               string
	ApiTokenFullInfo          string
	ApiBondingCurves          string
	ApiGetGameFunGInstruction string
	ApiBuyGWithPoints         string
	httpClient                *http.Client
	metricsClient             *metrics.MetricsHTTPClient
)

func init() {
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
		},
	}
}

func InitMetrics(redisClient *redis.Client) {
	metricsClient = metrics.NewMetricsHTTPClient(httpClient, redisClient)
}

func GetHTTPClient() *metrics.MetricsHTTPClient {
	if metricsClient != nil {
		return metricsClient
	}
	return nil
}

func InitAPI(endpoint *string) {
	ApiTokenBalanceBatch = *endpoint + "/block/api/v1/token-balance/batch"
	ApiDexCheckBatch = *endpoint + "/block/api/v1/dex-check-batch"
	ApiSafetyCheckBatch = *endpoint + "/block/api/v1/safety-check/authority"
	ApiSafetyCheckPool = *endpoint + "/block/api/v1/safety-check/pool"
	ApiPriorityFee = *endpoint + "/block/api/v1/get-priority-fee"
	ApiPumpfunTrade = *endpoint + "/block/api/v1/pumpfun/trade"
	ApiRaydiumTrade = *endpoint + "/block/api/v1/raydium/trade"
	ApiTransactionSend = *endpoint + "/block/api/v1/transaction/send"
	ApiTransactionStatus = *endpoint + "/block/api/v1/transaction/status"
	ApiTokensBatch = *endpoint + "/block/api/v1/tokens/batch"
	ApiTipFloor = *endpoint + "/block/api/v1/tip-floor"
	ApiPoolInfo = *endpoint + "/block/api/v1/pool-info/mints"
	ApiTokenFullInfo = *endpoint + "/block/api/v1/tokens/full-info"
	ApiBondingCurves = *endpoint + "/block/api/v1/bonding-curves"
	ApiGetGameFunGInstruction = *endpoint + "/block/api/v1/gamefun/swap-g-instruction"
	ApiBuyGWithPoints = *endpoint + "/block/api/v1/gamefun/buy-g-with-points"
	ApiGameFunTransactionSend = *endpoint + "/block/api/v1/gamefun/send-transaction"
}

// FetchURIWithRetry adds retry mechanism for fetching URI
func FetchURIWithRetry(uri string, maxRetries int) ([]byte, error) {
	if uri == "" {
		return nil, fmt.Errorf("empty URI provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if maxRetries <= 0 {
		maxRetries = 3
	}

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context deadline exceeded: %w", ctx.Err())
		default:
			body, err := fetchWithContext(ctx, uri)
			if err == nil {
				return body, nil
			}
			lastErr = err
			backoffDuration := getBackoffDuration(i)
			util.Log().Warning("Attempt %d failed to fetch URI %s: %v, retrying in %v",
				i+1, uri, err, backoffDuration)
			time.Sleep(backoffDuration)
		}
	}
	return nil, fmt.Errorf("failed to fetch URI after %d attempts: %w", maxRetries, lastErr)
}

func fetchWithContext(ctx context.Context, uri string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 首先尝试 Windows-1252 解码
	decoder := charmap.Windows1252.NewDecoder()
	utf8Body, err := decoder.Bytes(body)
	if err != nil {
		// 如果 Windows-1252 转换失败，使用原始数据继续处理
		utf8Body = body
	}

	// 然后进行 Unicode 规范化
	normalized := norm.NFKC.Bytes(utf8Body)

	// 验证最终结果是否为有效的 UTF-8
	if !utf8.Valid(normalized) {
		return body, nil // 如果所有处理都失败，返回原始数据
	}

	return normalized, nil
}

func getBackoffDuration(attempt int) time.Duration {
	backoff := time.Duration(1<<uint(attempt)) * time.Second
	maxBackoff := 30 * time.Second
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	return backoff
}

func GetRaydiumTradeTx(tradeStruct httpRequest.SwapRaydiumStruct) (*http.Response, error) {
	url := ApiRaydiumTrade
	resp, err := postRequest(url, tradeStruct, false)
	if err != nil {
		return nil, fmt.Errorf("failed to send raydium trade request: %w", err)
	}
	return resp, nil
}

func GetPumpFunTradeTx(swapStruct httpRequest.SwapPumpStruct) (*http.Response, error) {
	url := ApiPumpfunTrade
	resp, err := postRequest(url, swapStruct, false)
	if err != nil {
		return nil, fmt.Errorf("failed to send swap request: %w", err)
	}
	return resp, nil
}

func GetGameFunGInstruction(swapGInstructionStruct httpRequest.SwapGInstructionStruct) (*http.Response, error) {
	url := ApiGetGameFunGInstruction
	resp, err := postRequest(url, swapGInstructionStruct, false)
	if err != nil {
		return nil, fmt.Errorf("failed to send swap request: %w", err)
	}
	return resp, nil
}

func GetBuyGWithPointsInstruction(buyGWithPointsStruct httpRequest.BuyGWithPointsStruct) (*http.Response, error) {
	url := ApiBuyGWithPoints
	resp, err := postRequest(url, buyGWithPointsStruct, false)
	if err != nil {
		return nil, fmt.Errorf("failed to send swap request: %w", err)
	}
	return resp, nil
}

func postRequest(url string, body interface{}, useMetrics bool) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	var resp *http.Response
	if useMetrics {
		resp, err = GetHTTPClient().Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	} else {
		resp, err = httpClient.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

func GetTokenBalance(accounts []string, mintTokenAddress string) (*[]httpRespone.SolBalanceResponseData, error) {
	url := ApiTokenBalanceBatch
	body := map[string]interface{}{
		"addresses": accounts,
		"mint":      mintTokenAddress,
	}
	resp, err := postRequest(url, body, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get token balance: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var balanceResponse httpRespone.SolBalanceResponse
	if err := json.Unmarshal([]byte(bodyBytes), &balanceResponse); err != nil {
		log.Fatalf("Failed to decode dex check response: %v", err)
	}

	return &balanceResponse.Data, nil
}

func GetDexCheck(tokenAddresses []string) (*[]httpRespone.DexCheckData, error) {
	url := ApiDexCheckBatch
	body := map[string]interface{}{
		"tokenAddresses": tokenAddresses,
	}
	resp, err := postRequest(url, body, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get dex check data: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	var dexscrCheck httpRespone.DexCheckResponse
	if err := json.Unmarshal([]byte(bodyBytes), &dexscrCheck); err != nil {
		log.Fatalf("Failed to decode dex check response: %v", err)
	}

	return &dexscrCheck.Data, nil
}

func GetSafetyCheckData(tokenAddresses []string) (*[]httpRespone.SafetyData, error) {

	url := ApiSafetyCheckBatch

	body := map[string]interface{}{
		"addresses": tokenAddresses,
	}

	resp, err := postRequest(url, body, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get safety check data: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var safetyCheckResponse httpRespone.SafetyCheckResponse
	if err := json.Unmarshal([]byte(bodyBytes), &safetyCheckResponse); err != nil {
		log.Fatalf("Failed to decode safety check response: %v", err)
	}

	return &safetyCheckResponse.Data, nil
}

func GetPriorityFee() (*httpRespone.GasFeeResponse, error) {
	url := ApiPriorityFee

	resp, err := GetHTTPClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get priority fee: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read priority fee response body: %w", err)
	}

	var gasFee httpRespone.GasFee
	if err := json.Unmarshal(body, &gasFee); err != nil {
		return nil, fmt.Errorf("failed to unmarshal priority fee response: %w", err)
	}

	return &gasFee.Data, nil
}

func SendSwapRequest(swapStruct httpRequest.SwapPumpStruct) (*http.Response, error) {
	url := ApiPumpfunTrade

	resp, err := postRequest(url, swapStruct, false)
	if err != nil {
		return nil, fmt.Errorf("failed to send swap request: %w", err)
	}
	return resp, nil
}

func GetPythResponse() (*[]httpRespone.PythResponse, error) {
	now := time.Now().Unix()
	url := fmt.Sprintf("https://benchmarks.pyth.network/v1/shims/tradingview/history?symbol=Crypto.SOL/USD&resolution=1&from=%d&to=%d", now-60, now)

	resp, err := GetHTTPClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get pyth response: %w", err)
	}
	defer resp.Body.Close()

	var pythResp httpRespone.PythResponse
	if err := json.NewDecoder(resp.Body).Decode(&pythResp); err != nil {
		return nil, fmt.Errorf("failed to decode pyth response: %w", err)
	}

	return &[]httpRespone.PythResponse{pythResp}, nil
}

func SendTransaction(swapTransaction string, isJito bool) (*http.Response, error) {
	url := ApiTransactionSend
	body := map[string]interface{}{
		"signedTransaction": swapTransaction,
		"mev":               isJito,
	}

	resp, err := postRequest(url, body, false)
	if err != nil {
		return nil, fmt.Errorf("failed to send swap transaction: %w", err)
	}
	return resp, nil
}

func SendGameFunTransaction(swapTransaction string, isJito bool, isUsePoint bool) (*httpRespone.SendResponse, error) {
	url := ApiGameFunTransactionSend
	body := map[string]interface{}{
		"signedTransaction": swapTransaction,
		// "mev":               isJito,
		"usePoint": isUsePoint,
	}

	resp, err := postRequest(url, body, false)
	if err != nil {
		return nil, fmt.Errorf("failed to send swap transaction: %w", err)
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read priority fee response body: %w", err)
	}

	var apiResp *httpRespone.SendResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal priority fee response: %w", err)
	}

	return apiResp, nil
}

func GetSwapStatusBySignature(signature string) (*httpRespone.ApiResponse, error) {
	url := fmt.Sprintf("%s?signature=%s", ApiTransactionStatus, signature)

	resp, err := GetHTTPClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get priority fee: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read priority fee response body: %w", err)
	}

	var apiResp *httpRespone.ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal priority fee response: %w", err)
	}

	return apiResp, nil
}

func GetTokenInfoDefi(tokenAddresses []string, chainType uint8) ([]httpRespone.Token, error) {
	url := ApiTokensBatch

	body := map[string]interface{}{
		"chain_type":      chainType,
		"token_addresses": tokenAddresses,
	}

	resp, err := postRequest(url, body, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get token info defi: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var tokenInfoDefiResponse httpRespone.TokenInfoDefiResponse
	if err := json.Unmarshal(bodyBytes, &tokenInfoDefiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token info defi response: %w", err)
	}
	return tokenInfoDefiResponse.Data, nil
}

func GetSolPrice() (float64, error) {
	solPriceID := "0xef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d"
	url := fmt.Sprintf("https://hermes.pyth.network/v2/updates/price/latest?ids[]=%s", solPriceID)

	resp, err := GetHTTPClient().Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch SOL price: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Parsed []struct {
			Price struct {
				Price       string `json:"price"`
				Conf        string `json:"conf"`
				Expo        int    `json:"expo"`
				PublishTime int64  `json:"publish_time"`
			} `json:"price"`
		} `json:"parsed"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Parsed) == 0 {
		return 0, fmt.Errorf("no price data received")
	}

	// 将价格字符串转为 float64
	price, err := strconv.ParseFloat(response.Parsed[0].Price.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	// 应用指数
	expo := response.Parsed[0].Price.Expo
	actualPrice := price * math.Pow10(expo)

	return actualPrice, nil
}

func GetSafetyCheckPool(tokenAddressPoolAddresses []map[string]string) (*[]httpRespone.SafetyCheckPoolData, error) {
	url := ApiSafetyCheckPool
	var request struct {
		Data []map[string]string `json:"data"`
	}
	request.Data = tokenAddressPoolAddresses
	resp, err := postRequest(url, request, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get safety check pool: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var safetyCheckResponse httpRespone.SafetyCheckPoolResponse
	if err := json.Unmarshal(bodyBytes, &safetyCheckResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal safety check response: %w", err)
	}
	return &safetyCheckResponse.Data, nil
}

func GetTipFloor(tokenAddress string) (*httpRespone.TipFloorResponse, error) {
	url := ApiTipFloor
	resp, err := GetHTTPClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get tip floor: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var tipFloorResponse httpRespone.TipFloorResponse
	if err := json.Unmarshal(bodyBytes, &tipFloorResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tip floor response: %w", err)
	}

	return &tipFloorResponse, nil
}

func GetPoolInfo(tokenAddresses []string) (*[]httpRespone.MintData, error) {

	url := ApiPoolInfo + "?mints=" + strings.Join(tokenAddresses, ",")

	resp, err := GetHTTPClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool info: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var poolInfoResponse httpRespone.PoolInfoResponse
	if err := json.Unmarshal(bodyBytes, &poolInfoResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pool info response: %w", err)
	}

	return &poolInfoResponse.Data, nil
}
func GetPoolInfoByPoolAddress(poolAddresses []string) (*[]httpRespone.PoolItem2, error) {

	url := ApiPoolInfo + "/address?address=" + strings.Join(poolAddresses, ",")

	resp, err := GetHTTPClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool info: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var poolInfoResponse httpRespone.PoolInfoResponse2
	if err := json.Unmarshal(bodyBytes, &poolInfoResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pool info response: %w", err)
	}

	return &poolInfoResponse.Data, nil
}

func GetTokenFullInfo(tokenAddresses []string, chainType string) (*httpRespone.TokenFullInfoResponse, error) {
	url := ApiTokenFullInfo

	body := map[string]interface{}{
		"chain_type":      chainType,
		"token_addresses": tokenAddresses,
	}

	resp, err := postRequest(url, body, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get token full info: %w", err)
	}

	var result httpRespone.TokenFullInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	defer resp.Body.Close()

	return &result, nil
}

func GetBondingCurves(addresses []string) (*httpRespone.BondingCurvesResponse, error) {
	url := ApiBondingCurves

	// 构造请求体
	body := map[string]interface{}{
		"addresses": addresses,
	}

	resp, err := postRequest(url, body, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get bonding curves: %w", err)
	}
	defer resp.Body.Close()

	var result httpRespone.BondingCurvesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode bonding curves response: %w", err)
	}

	return &result, nil
}

func GenerateJitoOrderId(tokenOutAddress string, tokenInAddress string, inAmount string, jitoTip string) string {
	return fmt.Sprintf("%s-%s-%s-%s", tokenOutAddress, tokenInAddress, inAmount, jitoTip)
}
