package util

import (
	// "brcoin"

	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

//랜덤 스트링 생성
func MakeRandomString(n int, test string) string {
	// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(test))))
		b[i] = test[n.Int64()]
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

func SignatureVerification(stringPublicKey string, signature string) (bool, error) {

	stringPublicKey = strings.Replace(stringPublicKey, `\n`, "\n", -1)

	block, _ := pem.Decode([]byte(stringPublicKey))
	if block == nil {
		return false, errors.New("failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, errors.New("failed to parse DER encoded public key: " + err.Error())
	}

	pubkey := publicKey.(*rsa.PublicKey)

	message := []byte("")

	data, err := base64.StdEncoding.DecodeString(signature)

	if err != nil {
		return false, errors.New(err.Error())
	}
	hashed := sha256.Sum256(message)

	//서명 검증
	err = rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hashed[:], data)

	//서명 검증 실패
	if err != nil {
		fmt.Printf(err.Error())
		return false, errors.New("Error from verification: " + err.Error())

	}

	return true, nil
}
