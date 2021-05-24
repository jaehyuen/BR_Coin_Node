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

	fmt.Println("initinitinitinitinitinitinitinitinitinitinit")
	brcoin.InitBrcoin(stub)

	if err := brcoin.InitBrcoin(stub); err != nil {
		return shim.Error(brcoin.CODE9999 + " Init")

	}
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
		return s.transfer(stub, args) //토큰(코인) 송금
	} else if function == "mint" {
		return s.mint(stub, args) //토큰(코인) 추가 발행
	} else if function == "burn" {
		return s.burn(stub, args) //토큰(코인) 소각
	} else if function == "queryAllTokens" {
		return s.queryAllTokens(stub, args) //토큰(코인) 소각
	} else if function == "init" {
		return s.Init(stub) //토큰(코인) 소각
	}

	return shim.Error(brcoin.CODE9999 + " Invalid Smart Contract function name")
}

/*
 * 토큰(코인) 생성
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

	err := brcoin.CreateToken(stub, token)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

/*
 * 지갑 생성
 * args[0]: 지갑 public key
 *
 * return string 지갑 id
 */

func (s *SmartContract) createWallet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	response, err := brcoin.CreateWallet(stub, args[0])

	dataBytes, _ := json.Marshal(response)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(dataBytes)

}

/*
 * 지갑 조회
 * args[0]: 지갑 주소
 *
 * return structure.Token 지갑정보 json
 */

func (s *SmartContract) queryWallet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//지갑 데이터 조회
	walletData, err := brcoin.GetWallet(stub, args[0])
	walletByte, _ := json.Marshal(walletData)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(walletByte)

}

/*
 * 토큰(코인) 총 발행량 조회
 * args[0]: tokenId
 *
 * return string 해당 토큰의 총 발행량
 */

func (s *SmartContract) totalSupply(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//해당 토큰의 총 발행량 조회
	totalSupplyData, err := brcoin.GetToken(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	totalSupplyByte, _ := json.Marshal(totalSupplyData.TotalSupply)

	return shim.Success(totalSupplyByte)

}

/*
 * 잔고 조회
 * args[0]: 지갑 주소
 *
 * return []structure.BalanceInfo 지갑정보 jsonArray
 */

func (s *SmartContract) balanceOf(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	walletData, err := brcoin.GetWallet(stub, args[0])
	balanceByte, _ := json.Marshal(walletData.Balance)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(balanceByte)

}

/*
 * 토큰(코인) 송금
 * args[0]: []structure.Transfer
 * detail args[0]:
 * {
 *  "fromAddr": "보내는 지갑주소",
 *  "toAddr": "받는 지갑주소",
 *  "amount": "보내는 토큰(코인 량)",
 *  "tokenId": "토큰(코인) id",
 *  "unlockDate": "거래 금지 날짜(타임스탬프)"
 * }
 *
 */
func (s *SmartContract) transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var transfer *structure.Transfer

	if err := json.Unmarshal([]byte(args[0]), &transfer); err != nil {
		fmt.Println(err.Error())
		return shim.Error(brcoin.CODE0006 + " transfer")

	}

	err := brcoin.TransferToken(stub, transfer)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

/*
* 토큰(코인) 추가 발행
 */
func (s *SmartContract) mint(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//지갑 데이터 조회
	walletData, err := brcoin.GetWallet(stub, args[0])
	walletByte, _ := json.Marshal(walletData)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(walletByte)

}

/*
* 토큰(코인) 소각
 */
func (s *SmartContract) burn(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//지갑 데이터 조회
	walletData, err := brcoin.GetWallet(stub, args[0])
	walletByte, _ := json.Marshal(walletData)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(walletByte)

}

func (s *SmartContract) queryAllTokens(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//지갑 데이터 조회
	tokenData, err := brcoin.FindAllTokens(stub)
	tokenByte, _ := json.Marshal(tokenData)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(tokenByte)

}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
