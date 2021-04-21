package util

import (
	"crypto/rand"
	"math/big"
	"strings"
)

//랜덤 스트링 생성
func MakeRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// 지갑주소 유효성 검사
func AddressValidation(address string) bool {
	if len(strings.TrimSpace(address)) != 40 {

		return true
	}
	if address[:3] != "BRC" {

		return true
	}
	return false
}
