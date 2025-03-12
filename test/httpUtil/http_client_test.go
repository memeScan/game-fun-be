package httpUtil_test

import (
	"game-fun-be/internal/pkg/httpRequest"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// 测试前的设置
	setup()

	// 运行测试
	code := m.Run()

	// 测试后的清理
	teardown()

	os.Exit(code)
}

func setup() {
	// 1. 加载测试环境的 .env 文件
	loadTestEnv()

	// 2. 设置环境变量
	os.Setenv("APP_ENV", "test")
	os.Setenv("LOG_LEVEL", "debug")

	endpoint := os.Getenv("BLOCKCHAIN_API_ENDPOINT")
	httpUtil.InitAPI(&endpoint)
	redis.Redis()
	httpUtil.InitMetrics(redis.RedisClient)

	// 3. 初始化日志
	util.BuildLogger("debug")

	util.Log().Info("Test environment setup completed")
}

func teardown() {
	util.Log().Info("Test environment cleanup completed")
}

// 可选：每个测试用例的设置
func setupTest(t *testing.T) {
	t.Helper()
	util.Log().Info("Starting test: %s", t.Name())
}

func TestGetPythResponse(t *testing.T) {
	setupTest(t)
	resp, err := httpUtil.GetPythResponse()
	if err != nil {
		t.Fatalf("Failed to get Pyth response: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	pythResp := (*resp)[0]
	t.Logf("Received Pyth Response: %+v", pythResp)
}

func TestGetSolPrice(t *testing.T) {
	setupTest(t)
	price, err := httpUtil.GetSolPrice()

	t.Logf("Current SOL price: $%.2f", price)

	if err != nil {
		t.Fatalf("Failed to get SOL price: %v", err)
	}

	t.Logf("Current SOL price: $%.2f", price)

	if price <= 0 {
		t.Error("Expected positive SOL price")
	}
}

func TestGetTokenBalance(t *testing.T) {
	setupTest(t)
	accounts := []string{"5YNmS1R9nNSCDzb5a7mMJ1dwK9uHeAAF4CmPEwKgVWr8", "AF1pcxhN9rYVg2W3J3gYc6W9aWvyp9PRY2qgQiFyANr4"}
	mintTokenAddress := "2d7zw2qCXUrXdiGvTyQwNH99a5vACKAy4CzpD8Y4DLVo"

	resp, err := httpUtil.GetTokenBalance(accounts, mintTokenAddress)
	if err != nil {
		t.Fatalf("Failed to get token balance: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	t.Logf("Received Token Balance Response: %+v", *resp)
}

func TestGetDexCheck(t *testing.T) {
	setupTest(t)
	tokenAddresses := []string{"5YNmS1R9nNSCDzb5a7mMJ1dwK9uHeAAF4CmPEwKgVWr8", "AF1pcxhN9rYVg2W3J3gYc6W9aWvyp9PRY2qgQiFyANr4"}

	resp, err := httpUtil.GetDexCheck(tokenAddresses)
	if err != nil {
		t.Fatalf("Failed to get dex check data: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	t.Logf("Received Dex Check Response: %+v", *resp)
}

func TestGetSafetyCheckData(t *testing.T) {
	setupTest(t)
	tokenAddresses := []string{"2XeFanPEcSCNNwYcGgzqk8NR8KuJEdkM2jZ3wrGqJK5p"}

	resp, err := httpUtil.GetSafetyCheckData(tokenAddresses)
	if err != nil {
		t.Fatalf("Failed to get safety check data: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	t.Logf("Received Safety Check Response: %+v", *resp)
}

func TestGetPriorityFee(t *testing.T) {
	setupTest(t)
	resp, err := httpUtil.GetPriorityFee()
	if err != nil {
		t.Fatalf("Failed to get priority fee: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	t.Logf("Received Priority Fee Response: %+v", *resp)
}

func TestSendSwapRequest(t *testing.T) {
	setupTest(t)
	swapStruct := httpRequest.SwapPumpStruct{
		// Fill in the required fields for testing
	}

	resp, err := httpUtil.SendSwapRequest(swapStruct)
	if err != nil {
		t.Fatalf("Failed to send swap request: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	t.Logf("Received Swap Request Response: %+v", resp)
}

func TestSendSwapTransaction(t *testing.T) {
	setupTest(t)
	swapTransaction := "dummy_swap_transaction"

	resp, err := httpUtil.SendTransaction(swapTransaction, false)
	if err != nil {
		t.Fatalf("Failed to send swap transaction: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	t.Logf("Received Swap Transaction Response: %+v", resp)
}

func TestGetSwapRequestStatusBySignature(t *testing.T) {
	setupTest(t)
	signature := "dummy_signature"

	resp, err := httpUtil.GetSwapStatusBySignature(signature)
	if err != nil {
		t.Fatalf("Failed to get swap request status: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	t.Logf("Received Swap Request Status Response: %+v", resp)
}

func TestGetTokenInfoDefi(t *testing.T) {
	setupTest(t)
	tokenAddress := "4xA8psLrTiifGGKrqLs7k5uo3boyYQZpJGxk1tJBsMeV"
	chainType := uint8(1)

	resp, err := httpUtil.GetTokenInfoDefi([]string{tokenAddress}, chainType)
	if err != nil {
		t.Fatalf("Failed to get token info defi: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	t.Logf("Received Token Info Defi Response: %+v", resp)
}

// func TestGetTokenInfos(t *testing.T) {
// 	tokenAddresses := []string{"BbqyciFUdBgJ8RGuR9zx7V7bdGc1w9HiKXTiMDoDpump", "9jBkMwUXzqfDfWrzwdscqMXL4cdGrFRgjq6gHyTPpump"}
// 	chainType := uint8(1)

// 	resp, err := httpUtil.GetTokenInfos(tokenAddresses, chainType)
// 	if err != nil {
// 		t.Fatalf("Failed to get token infos: %v", err)
// 	}

// 	if resp == nil {
// 		t.Fatal("Expected response, got nil")
// 	}

// 	t.Logf("Received Token Infos Response: %+v", resp)
// }

// 加载测试环境的 .env 文件
func loadTestEnv() {
	// 获取测试目录路径
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)

	// 构建 .env.test 文件路径
	envFile := filepath.Join(testDir, "..", ".env.test")

	// 加载环境变量文件
	if err := godotenv.Load(envFile); err != nil {
		util.Log().Warning("Error loading test env file: %v", err)
	}
}

func TestGetTokenFullInfo(t *testing.T) {
	setupTest(t)
	tokenAddresses := []string{"9yXuW9iu9YYfHhMj2wQUJCPwZbVEZ7JwRrby1rTSpump"}
	chainType := "sol"

	resp, err := httpUtil.GetTokenFullInfo(tokenAddresses, chainType)
	if err != nil {
		t.Fatalf("Failed to get token full info: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.Code != 2000 {
		t.Errorf("Expected code 2000, got %d", resp.Code)
	}

	if len(resp.Data) > 0 {
		tokenInfo := resp.Data[0]
		t.Logf("Token Info: Name=%s, Symbol=%s, Decimals=%d",
			tokenInfo.Name, tokenInfo.Symbol, tokenInfo.Decimals)
	}
}

func TestGetPoolInfo(t *testing.T) {
	setupTest(t)
	tokenAddresses := []string{"3wnN4QSD9aQ89bJunXnPWqeG9sgMY1gGoPwoXrFsRY4e"}

	resp, err := httpUtil.GetPoolInfo(tokenAddresses)
	if err != nil {
		t.Fatalf("Failed to get pool info: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if len(*resp) > 0 {
		poolInfo := (*resp)[0]
		t.Logf("Pool Info: BaseMint=%s, QuoteMint=%s",
			poolInfo.Mint, poolInfo.Mint)
	}
}

func TestGetBondingCurves(t *testing.T) {
	setupTest(t)
	addresses := []string{"BkKVKe1zwcs6PqSszd3R1ipadML6bR8eww8f9xjVZYVE"}

	resp, err := httpUtil.GetBondingCurves(addresses)
	if err != nil {
		t.Fatalf("Failed to get bonding curves: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.Code != 2000 {
		t.Errorf("Expected code 2000, got %d", resp.Code)
	}

	if len(resp.Data) > 0 {
		curveInfo := resp.Data[0]
		t.Logf("Bonding Curve Info: Address=%s, TokenSupply=%s, Complete=%v",
			curveInfo.Address,
			curveInfo.Data.TokenTotalSupply,
			curveInfo.Data.Complete)
	}
}

// func TestGetPriorityFee(t *testing.T) {
//     setupTest(t)

//     // 调用 GetPriorityFee 方法
//     gasFeeResponse, err := httpUtil.GetPriorityFee()

//     // 检查是否有错误
//     if err != nil {
//         t.Fatalf("Failed to get priority fee: %v", err)
//     }

//     // 检查返回值是否为空
//     if gasFeeResponse == nil {
//         t.Fatal("Expected response, got nil")
//     }

//     // 打印返回的 GasFeeResponse 信息
//     t.Logf("Received Priority Fee Response: Low=%s, High=%s, Medium=%s",
//         gasFeeResponse.PriorityFeeLevels.Low, gasFeeResponse.PriorityFeeLevels.High, gasFeeResponse.PriorityFeeLevels.Medium)

//     // 检查 BaseFee 和 PriorityFee 是否为空
//     if gasFeeResponse.PriorityFeeLevels.Min == 0 {
//         t.Error("Expected non-empty BaseFee")
//     }
//     if gasFeeResponse.PriorityFeeLevels.Medium == 0 {
//         t.Error("Expected non-empty PriorityFee")
//     }
// }
