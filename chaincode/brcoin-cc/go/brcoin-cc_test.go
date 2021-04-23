package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

var scc *SmartContract
var stub *shimtest.MockStub

func setup() {
	scc = new(SmartContract)
	stub = shimtest.NewMockStub("unittest", scc)
}
func shutdown() {
	scc = nil
	stub = nil
}

func Reset() {
	shutdown()
	setup()
}

func checkQuery(t *testing.T, stub *shimtest.MockStub, function string, param string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(param)})

	if res.Status != shim.OK {
		fmt.Println("Query", param, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", param, "failed to get value")
		t.FailNow()
	}

	fmt.Println("[Query] result : ", string(res.Payload))

}

func checkInvoke(t *testing.T, stub *shimtest.MockStub, args [][]byte) string {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}

	var test, _ = strconv.Unquote(string(res.Payload))
	return test

}
func TestCreateWallet(t *testing.T) {
	fmt.Println("[TestCreateWallet] start")
	address1 := checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("12341234")})
	checkQuery(t, stub, "queryWallet", address1)
	fmt.Println("[TestCreateWallet] fin")
	fmt.Println()

}

// func TestCreateToken(t *testing.T) {

// 	fmt.Println("[TestCreateToken] start")
// 	//토큰 생성을 위한 테스트 지갑 생성
// 	address1 := checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("12341234")})

// 	//0번토큰(BRC) 생성
// 	tokenData := "{\"owner\": \"" + address1 + "\", \"symbol\": \"BRC\", \"totalSupply\": \"1000000\", \"name\": \"brcoin\", \"information\": \"thisisbrcoin\", \"url\": \"https: //github.com/jaehyuen/BR_Coin_Node\",	\"decimal\": 3, \"reserve\": []}"
// 	checkInvoke(t, stub, [][]byte{[]byte("createToken"), []byte(tokenData)})

// 	//1번토큰(AAA) 생성
// 	tokenData = "{\"owner\": \"" + address1 + "\", \"symbol\": \"AAA\", \"totalSupply\": \"1000\", \"name\": \"testaaa\", \"information\": \"a\", \"url\": \"https: //github.com/jaehyuen/BR_Coin_Node\",\"decimal\": 8, \"reserve\": []}"
// 	checkInvoke(t, stub, [][]byte{[]byte("createToken"), []byte(tokenData)})
// 	checkQuery(t, stub, "queryWallet", address1)

// 	fmt.Println("[TestCreateToken] fin")
// 	fmt.Println()

// }

func TestCreateTokenAndTransfer(t *testing.T) {

	fmt.Println("[TestCreateTokenAndTransfer] start")
	//토큰 생성을 위한 테스트 지갑2개 생성
	address1 := checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("12341234")})
	address2 := checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("12341234")})

	//0번토큰(BRC) 생성
	tokenData := "{\"owner\": \"" + address1 + "\", \"symbol\": \"BRC\", \"totalSupply\": \"1000000\", \"name\": \"brcoin\", \"information\": \"thisisbrcoin\", \"url\": \"https: //github.com/jaehyuen/BR_Coin_Node\",	\"decimal\": 3, \"reserve\": []}"
	checkInvoke(t, stub, [][]byte{[]byte("createToken"), []byte(tokenData)})

	//1번토큰(AAA) 생성
	tokenData = "{\"owner\": \"" + address1 + "\", \"symbol\": \"AAA\", \"totalSupply\": \"1000\", \"name\": \"testaaa\", \"information\": \"a\", \"url\": \"https: //github.com/jaehyuen/BR_Coin_Node\",\"decimal\": 8, \"reserve\": []}"
	checkInvoke(t, stub, [][]byte{[]byte("createToken"), []byte(tokenData)})

	// 1번지갑에서 2번지갑에 토큰 송금 1
	transferData := "{\"fromAddr\": \"" + address1 + "\",\"toAddr\": \"" + address2 + "\",\"amount\": \"3.312\",\"tokenId\": \"1\",\"unlockDate\": \"0\"}"
	checkInvoke(t, stub, [][]byte{[]byte("transfer"), []byte(transferData)})

	// 1번지갑에서 2번지갑에 토큰 송금 1
	transferData = "{\"fromAddr\": \"" + address1 + "\",\"toAddr\": \"" + address2 + "\",\"amount\": \"3.12345678\",\"tokenId\": \"1\",\"unlockDate\": \"1617081582\"}"
	checkInvoke(t, stub, [][]byte{[]byte("transfer"), []byte(transferData)})

	checkQuery(t, stub, "queryWallet", address1)
	checkQuery(t, stub, "queryWallet", address2)

	checkQuery(t, stub, "balanceOf", address1)
	checkQuery(t, stub, "balanceOf", address2)

	fmt.Println("[TestCreateTokenAndTransfer] fin")
	fmt.Println(time.Now().Unix())

}

// func TestGetTotalSupply(t *testing.T) {

// 	fmt.Println("[TestGetTotalSupply] start")
// 	checkQuery(t, stub, "totalSupply", "1")
// 	checkQuery(t, stub, "totalSupply", "0")
// 	fmt.Println("[TestGetTotalSupply] fin")
// 	fmt.Println()

// }

func TestMain(m *testing.M) {
	setup()

	m.Run()
	shutdown()
}
