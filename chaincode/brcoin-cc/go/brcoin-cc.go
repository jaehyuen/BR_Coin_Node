package main

import (
	"brcoin"
	"encoding/json"
	"fmt"
	"log"

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
	} else if function == "query" {
		return s.query(stub, args) //Tsa토큰 조회 (DocuSeq)
	}

	return shim.Error(brcoin.CODE9999 + " Invalid Smart Contract function name")
}

func (s *SmartContract) createToken(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	// var dataBytes []byte
	// var token *structure.Token

	// if err := json.Unmarshal([]byte(args[0]), token); err != nil {
	// 	return shim.Error(brcoin.CODE9999 + " Invalid Smart Contract function name")

	// }

	// response, err = brcoin.CreateToken(ctx, token)

	// dataBytes, _ = json.Marshal(response)
	// log.Println(string(dataBytes))

	// if err != nil {
	// 	return "", errors.New(string(dataBytes))
	// }
	return shim.Success(nil)

}

func (s *SmartContract) createWallet(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("[createWallet] result : ")
	var dataBytes []byte
	var err error
	var response string

	var password = args[0]

	response, err = brcoin.CreateWallet(stub, password)

	dataBytes, _ = json.Marshal(response)
	log.Println(string(dataBytes))

	if err != nil {
		return shim.Error(string(dataBytes))
	}
	return shim.Success(dataBytes)

}

func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var tsaAsBytes []byte
	fmt.Println("[query] wallet address is : ", args[0])

	//public 데이터 조회
	tsaAsBytes, _ = stub.GetState(args[0])

	fmt.Println("[query] : ", string(tsaAsBytes))

	return shim.Success(tsaAsBytes)

}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
