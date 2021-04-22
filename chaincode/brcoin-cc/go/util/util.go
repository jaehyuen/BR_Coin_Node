package util

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
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

// ParsePositive string is positive ?
func ParsePositive(s string) (decimal.Decimal, error) {

	//10진수인가?
	var d decimal.Decimal
	var err error

	if d, err = decimal.NewFromString(s); err != nil {
		fmt.Println("1")
		return d, errors.New("1101, " + s + " is not integer string")
	}
	if !d.IsPositive() {
		fmt.Println("3")
		return d, errors.New("1101, " + s + " is either 0 or negative.")
	}
	return d, nil
}
