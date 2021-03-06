package brcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"util"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"structure"
)

func PutWallet(stub shim.ChaincodeStubInterface, key string, walletData structure.BarakWallet, jobType string, jobArgs []string) error {

	var argsByte []byte
	var err error

	//지갑의 JobType을 파라미터로 받은 JobType으로 수정
	walletData.JobType = jobType
	//지갑의 JobDate를 현재시간으로 수정
	walletData.JobDate = time.Now().Unix()

	//배열이 nill이 아니거나 길이가 0보다 크면 arg를 byte 배열로 만든다
	if jobArgs != nil && len(jobArgs) > 0 {
		if argsByte, err = json.Marshal(jobArgs); err == nil {
			walletData.JobArgs = string(argsByte)
		}
	} else {
		walletData.JobArgs = ""
	}

	if err := PutState(stub, key, walletData); err != nil {
		fmt.Printf("[PutWallet Error] [%s] Error %s\n", key, err)
		return errors.New("8600,Hyperledger internal error - " + err.Error() + key)
	}
	return nil

}

func GetWallet(stub shim.ChaincodeStubInterface, address string) (structure.BarakWallet, error) {
	var walletData structure.BarakWallet

	//주소 유효성 검사
	if util.AddressValidation(address) {
		return walletData, errors.New(CODE0005 + " Address [" + address + "] is in the wrong format")
	}
	valueByte, err := stub.GetState(address)
	//오류 발생시 에러리턴
	if err != nil {
		return walletData, errors.New(CODE9999 + " Hyperledger internal error - " + err.Error())
	}

	//값이없으면 에러리턴
	if valueByte == nil {
		return walletData, errors.New(CODE0003 + " Can not find the address [" + address + "]")
	}

	//structure.BarakWallet 형식으로 Unmarshal
	if err = json.Unmarshal(valueByte, &walletData); err != nil {
		return walletData, errors.New(CODE0006 + " Address [" + address + "] is in the wrong data")
	}
	return walletData, nil
}

func GetToken(stub shim.ChaincodeStubInterface, tokenId string) (structure.Token, error) {

	var tokenData structure.Token

	valueByte, err := stub.GetState("TOKEN_DATA_" + tokenId)

	//오류 발생시 에러리턴
	if err != nil {
		return tokenData, errors.New(CODE9999 + " Hyperledger internal error - " + err.Error())
	}

	//값이없으면 에러리턴
	if valueByte == nil {
		return tokenData, errors.New(CODE0003 + " Can not find the tokenId [" + tokenId + "]")
	}

	//structure.Token 형식으로 Unmarshal
	if err = json.Unmarshal(valueByte, &tokenData); err != nil {
		return tokenData, errors.New(CODE0006 + " Address [" + tokenId + "] is in the wrong data")
	}

	return tokenData, nil
}

// 등록함수
func PutState(stub shim.ChaincodeStubInterface, key string, value interface{}) error {

	var err error
	var valueBytes []byte
	//public 데이터 변수 저장
	valueBytes, err = json.Marshal(value)

	fmt.Println(len(valueBytes))

	if err != nil {
		return errors.New("8600,Hyperledger internal error - " + err.Error() + key)
	}

	//public 데이터 등록
	if err := stub.PutState(key, valueBytes); err != nil {
		return errors.New("8600,Hyperledger internal error - " + err.Error() + key)
	}
	return nil

}

func InitBrcoin(stub shim.ChaincodeStubInterface) error {

	var err error
	var tokenByte []byte
	var address, totalSupply, symbol string

	symbol = "BRC"
	address = "BRCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	totalSupply = "100000000"

	tokenData := structure.Token{
		Owner:       address,
		Symbol:      symbol,
		CreateDate:  time.Now().Unix(),
		TotalSupply: totalSupply,
		TokeId:      0,
		Decimal:     8,
		JobDate:     time.Now().Unix(),
		JobType:     "CreateToken",
	}

	if tokenByte, err = json.Marshal(tokenData); err != nil {
		return errors.New(CODE0006 + " Invalid Data format")
	}

	if err = stub.PutState("TOKEN_DATA_0", tokenByte); err != nil {
		return err
	}

	if err = stub.PutState("TOKEN_MAX_NO", []byte(strconv.Itoa(0))); err != nil {
		return err
	}

	walletData := structure.BarakWallet{Regdate: time.Now().Unix(),
		PublicKey: "publicKey",
		JobDate:   time.Now().Unix(),
		JobType:   "NewWallet",
		Nonce:     "util.MakeRandomString(40)",
		Balance:   []structure.BalanceInfo{structure.BalanceInfo{Balance: totalSupply, TokenId: 0, UnlockDate: 0}}}

	address = "BRCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	if err := PutWallet(stub, address, walletData, "NewWallet", []string{address, "publicKey"}); err != nil {
		return err
	}
	return nil
}
