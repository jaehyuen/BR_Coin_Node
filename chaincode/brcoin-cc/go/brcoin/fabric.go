package brcoin

import (
	"encoding/json"
	"errors"
	"fmt" //cid 테스트??
	"time"

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

// 등록함수
func PutState(stub shim.ChaincodeStubInterface, key string, value interface{}) error {

	var err error
	var valueBytes []byte
	//public 데이터 변수 저장
	valueBytes, err = json.Marshal(value)

	if err != nil {
		return errors.New("8600,Hyperledger internal error - " + err.Error() + key)
	}

	//public 데이터 등록
	if err := stub.PutState(key, valueBytes); err != nil {
		return errors.New("8600,Hyperledger internal error - " + err.Error() + key)
	}
	return nil

}
