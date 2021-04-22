package main

import (
	"brcoin"
	"encoding/json"
	"fmt"
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
		return s.createToken(stub, args) //토큰(코인) 생성
	} else if function == "queryWallet" {
		return s.queryWallet(stub, args) //지갑 조회
	} else if function == "totalSupply" {
		return s.totalSupply(stub, args) // 토큰(코인의 총 발행량)
	} else if function == "balanceOf" {
		return s.balanceOf(stub, args) // 지갑에 있는 자산 조회
	} else if function == "transfer" {
		return s.transfer(stub, args) //토큰 송금
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

	var token *structure.Token

	if err := json.Unmarshal([]byte(args[0]), &token); err != nil {
		return shim.Error(brcoin.CODE0006 + " Token")

	}

	response, err := brcoin.CreateToken(stub, token)

	dataBytes, _ := json.Marshal(response)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(dataBytes)

}

func (s *SmartContract) createWallet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	response, err := brcoin.CreateWallet(stub, args[0])

	dataBytes, _ := json.Marshal(response)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(dataBytes)

}

func (s *SmartContract) queryWallet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//지갑 데이터 조회
	walletData, err := brcoin.GetWallet(stub, args[0])
	walletByte, _ := json.Marshal(walletData)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(walletByte)

}

func (s *SmartContract) totalSupply(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//해당 토큰의 총 발행량 조회
	totalSupplyData, err := brcoin.GetTotalSupply(stub, args[0])
	totalSupplyByte, _ := json.Marshal(totalSupplyData)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(totalSupplyByte)

}

func (s *SmartContract) balanceOf(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	return shim.Success(nil)

}
func (s *SmartContract) transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	return shim.Success(nil)

}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
