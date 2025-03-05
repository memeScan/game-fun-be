package service

import (
	"game-fun-be/internal/constants"

	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/gagliardetto/solana-go"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func VerifySolanaSignature(address, signature, message string) (bool, error) {
	pubkey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return false, errors.New("无效的公钥地址")
	}

	sigBytes, err := solana.SignatureFromBase58(signature)
	if err != nil {
		return false, errors.New("无效的签名格式")
	}

	isValid := pubkey.Verify([]byte(message), sigBytes)
	if !isValid {
		return false, errors.New("签名验证失败")
	}

	return true, nil
}

func UintToString(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

func GetUserTokenKey(userAddress string) string {
	return constants.UserTokenKeyFormat + userAddress
}

func GenerateInviteCode(address string) string {
	hash := sha256.Sum256([]byte(address))
	hashHex := hex.EncodeToString(hash[:])

	var codeBuilder strings.Builder
	codeBuilder.Grow(6)

	for i := 0; i < len(hashHex) && codeBuilder.Len() < 6; i++ {
		char := hashHex[i]
		index := int(char) % 62
		codeBuilder.WriteByte(base62Chars[index])
	}

	for codeBuilder.Len() < 6 {
		codeBuilder.WriteByte(base62Chars[0])
	}

	return codeBuilder.String()
}
