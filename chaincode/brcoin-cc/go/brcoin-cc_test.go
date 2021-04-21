package main

import (
	"fmt"
	"strconv"
	"testing"

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

func checkQuery(t *testing.T, stub *shimtest.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("queryWallet"), []byte(name)})

	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
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

	address1 := checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("12341234")})
	checkQuery(t, stub, address1, "zz")

}

func TestCreateToken(t *testing.T) {

	//토큰 생성을 위한 테스트 지갑 생성
	address1 := checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("12341234")})

	//0번토큰(BRC) 생성
	tokenData := "{\"owner\": \"" + address1 + "\", \"symbol\": \"BRC\", \"totalsupply\": \"1000000\", \"name\": \"brcoin\", \"information\": \"thisisbrcoin\", \"url\": \"https: //github.com/jaehyuen/BR_Coin_Node\",	\"decimal\": 3, \"reserve\": []}"
	checkInvoke(t, stub, [][]byte{[]byte("createToken"), []byte(tokenData)})

	//1번토큰(AAA) 생성
	tokenData = "{\"owner\": \"" + address1 + "\", \"symbol\": \"AAA\", \"totalsupply\": \"1000\", \"name\": \"testaaa\", \"information\": \"a\", \"url\": \"https: //github.com/jaehyuen/BR_Coin_Node\",	\"decimal\": 8, \"reserve\": []}"
	checkInvoke(t, stub, [][]byte{[]byte("createToken"), []byte(tokenData)})
	checkQuery(t, stub, address1, "zz")

}

func TestMain(m *testing.M) {
	setup()

	m.Run()
	shutdown()
}
