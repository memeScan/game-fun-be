package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gagliardetto/solana-go"
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
	sign := "+"
	if value < 0 {
		sign = "-"
		value = -value // 取绝对值
	}
	return fmt.Sprintf("%s%.2f%%", sign, value)
}
