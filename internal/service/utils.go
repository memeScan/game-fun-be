package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"
)

func VerifySolanaSignature(address, signature, message string) (bool, error) {
	pubkey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return false, errors.New("invalid public key address")
	}

	sigBytes, err := solana.SignatureFromBase58(signature)
	if err != nil {
		return false, errors.New("invalid signature format")
	}

	isValid := pubkey.Verify([]byte(message), sigBytes)
	if !isValid {
		return false, errors.New("signature verification failed")
	}

	return true, nil
}

func UintToString(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

func GetRedisKey(prefix string, parts ...string) string {
	key := prefix
	for _, part := range parts {
		if part != "" {
			key += ":" + part
		}
	}
	return key
}

func UnmarshalJSON(jsonStr string, target interface{}) error {
	if jsonStr == "" {
		return nil
	}

	err := json.Unmarshal([]byte(jsonStr), target)
	if err != nil {
		return err
	}

	return nil
}

func FormatPercent(value float64) string {
	// 如果 value 是 0.00，直接返回 "0"
	if value == 0.00 {
		return "0"
	}

	sign := "+"
	if value < 0 {
		sign = "-"
		value = -value // 取绝对值
	}
	return fmt.Sprintf("%s%.2f", sign, value)
}

// ConvertDecimalToInt 将 decimal.Decimal 转换为 int，支持四舍五入
func ConvertDecimalToInt(value decimal.Decimal, round bool) int {
	if round {
		// 如果需要四舍五入，先调用 Round(0) 方法
		value = value.Round(0)
	}
	// 转换为 int
	return int(value.IntPart())
}

func safeNewFromFloat(value float64) (decimal.Decimal, error) {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return decimal.Decimal{}, fmt.Errorf("invalid value: %v", value)
	}
	return decimal.NewFromFloat(value), nil
}

func processVolume(value float64, solPrice float64, decimals int) (decimal.Decimal, error) {
	volume, err := safeNewFromFloat(value)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("invalid volume value: %v", err)
	}

	solPriceDec, err := safeNewFromFloat(solPrice)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("invalid solPrice: %v", err)
	}

	return volume.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))).Mul(solPriceDec), nil
}

// 定义时间布局常量
const (
	ISO8601Layout = "2006-01-02T15:04:05Z"
)

// StringToTimestamp 将时间字符串转换为时间戳
func StringToTimestamp(timeStr string, layout string) (int64, error) {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time: %v", err)
	}
	return t.Unix(), nil
}

// 定义一个方法，用于保留两位小数
func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}

// 计算价格变化并保留两位小数
func calculatePriceChange(currentPrice, previousPrice float64) float64 {
	if currentPrice != 0 && previousPrice != 0 {
		return roundToTwoDecimalPlaces((currentPrice - previousPrice) / previousPrice)
	}
	return 0
}
