package service

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken はランダムな32byteのトークン文字列を生成する
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
