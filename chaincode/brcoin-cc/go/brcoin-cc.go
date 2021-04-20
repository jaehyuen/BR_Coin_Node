package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid" 
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
)

type SmartContract struct {
	contractapi.Contract
}

type jsonResponse struct {
	Key           string `json:"key"`
	ResultFlag    bool   `json:"resultFlag"`
	ResultCode    string `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}


//return 데이터 생성 함수
func _makeJson(key string, resultFlag bool, resultCode string, resultMessage string) (res jsonResponse) {

	var response jsonResponse
	response.Key = key
	response.ResultFlag = resultFlag
	response.ResultCode = resultCode
	response.ResultMessage = resultMessage

	return response
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create brcoin-cc chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting brcoin-cc chaincode: %s", err.Error())
	}
}
