package main

import (
	"brcoin"
	"encoding/json"
	"fmt"
	"log"
	"structure"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type SmartContract struct {
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()
	if function == "createWallet" {
		return s.createWallet(stub, args) //지갑 생성
	} else if function == "createToken" {
		return s.createToken(stub, args) //토큰 생성
	} else if function == "queryWallet" {
		return s.queryWallet(stub, args) //지갑 조회
	}

	return shim.Error(brcoin.CODE9999 + " Invalid Smart Contract function name")
}

/*
 * 토큰(코인) 셍상
 * args[0]: []structure.Token
 * detail args[0]:
 * {
 * "owner": "지갑주소",
 * "symbol": "토큰 심볼(BRC)",
 * "totalsupply": 토큰발행량(int),
 * "name": "토큰이름",
 * "information": "토큰 정보(nullable)",
 * "url": "토큰 관련 url(nullable)",
 * "decimal": 토큰총발행량(int),
 * "reserve": [     토큰 예약 리스트
 *   {
 *     "address": "지갑주소",
 *     "value": "토큰 량",
 *     "unlockdate": "거래 제한 날짜"
 *   }
 *  ]
 * }
 *
 */

func (s *SmartContract) createToken(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var dataBytes []byte
	var token *structure.Token
	var err error
	var response string

	if err := json.Unmarshal([]byte(args[0]), &token); err != nil {
		return shim.Error(brcoin.CODE0006 + " Token")

	}

	response, err = brcoin.CreateToken(stub, token)

	dataBytes, _ = json.Marshal(response)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(dataBytes)

}

func (s *SmartContract) createWallet(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var dataBytes []byte
	var err error
	var response string

	var password = args[0]

	response, err = brcoin.CreateWallet(stub, password)

	dataBytes, _ = json.Marshal(response)
	log.Println(string(dataBytes))

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(dataBytes)

}

func (s *SmartContract) queryWallet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var walletByte []byte
	// fmt.Println("[queryWallet] wallet address is : ", args[0])

	//public 데이터 조회
	walletData, err := brcoin.GetWallet(stub, args[0])
	walletByte, _ = json.Marshal(walletData)
	// fmt.Println("[queryWallet] : ", string(walletByte))

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(walletByte)

}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
