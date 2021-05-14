package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	"util"

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

// func TestCreateWallet(t *testing.T) {
// 	fmt.Println("[TestCreateWallet] start")
// 	address1 := checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("12341234")})
// 	checkQuery(t, stub, "queryWallet", address1)
// 	fmt.Println("[TestCreateWallet] fin")
// 	fmt.Println()

// }

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

func TestFull(t *testing.T) {
	stub.MockInit("1", [][]byte{[]byte("init")})

	var pubPEM = `-----BEGIN TEST PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCY7YuXgEeGlEn9EuIIBDKzV+IY\nr9WpaVHMNJcNuJMLAXEga/ACWaKSmDnWaziArjHYBjgXVqMlk2PImHKMaf2pv7ch\nhAhes1pBAzNoVai+zEqJCii+67Cy3ql7nLN3+aFBh457tiL3UcHG/pbryYKsizt5\nmC29e7if9uI6/xzj9QIDAQAB\n-----END TEST PUBLIC KEY-----\n`
	fmt.Println("[TestCreateTokenAndTransfer] start")
	//토큰 생성을 위한 테스트 지갑2개 생성
	address1 := checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte(pubPEM)})
	address2 := checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte(pubPEM)})
	checkQuery(t, stub, "queryWallet", address1)
	checkQuery(t, stub, "queryWallet", address2)
	checkQuery(t, stub, "queryWallet", "BRCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	fmt.Println("[address1] " + address1)
	fmt.Println("[address2] " + address2)
	checkQuery(t, stub, "queryAllTokens", "")

	//0번토큰(BRC2) 생성
	tokenData := "{\"owner\": \"" + address1 + "\", \"symbol\": \"BRC2\", \"totalSupply\": \"1000000\", \"name\": \"brcoin2\", \"information\": \"thisisbrcoin\", \"url\": \"https: //github.com/jaehyuen/BR_Coin_Node\",\"decimal\": 3, \"reserve\": []}"
	checkInvoke(t, stub, [][]byte{[]byte("createToken"), []byte(tokenData)})

	//1번토큰(AAA) 생성
	tokenData = "{\"owner\": \"" + address1 + "\", \"symbol\": \"AAA\", \"totalSupply\": \"1000\", \"name\": \"testaaa\", \"information\": \"a\", \"url\": \"https: //github.com/jaehyuen/BR_Coin_Node\",\"decimal\": 8, \"reserve\": []}"
	checkInvoke(t, stub, [][]byte{[]byte("createToken"), []byte(tokenData)})

	fmt.Println("[token create fin] ")
	// 1번지갑에서 2번지갑에 토큰 송금 1
	transferData := "{\"fromAddr\": \"" + address1 + "\",\"toAddr\": \"" + address2 + "\",\"amount\": \"3.312\",\"tokenId\": \"2\",\"unlockDate\": \"0\"}"
	checkInvoke(t, stub, [][]byte{[]byte("transfer"), []byte(transferData)})

	// 1번지갑에서 2번지갑에 토큰 송금 1
	transferData = "{\"fromAddr\": \"" + address1 + "\",\"toAddr\": \"" + address2 + "\",\"amount\": \"3.12345678\",\"tokenId\": \"2\",\"unlockDate\": \"1630623934\"}"
	checkInvoke(t, stub, [][]byte{[]byte("transfer"), []byte(transferData)})

	checkQuery(t, stub, "queryWallet", address1)
	checkQuery(t, stub, "queryWallet", address2)
	checkQuery(t, stub, "queryWallet", "BRCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	checkQuery(t, stub, "balanceOf", address1)
	checkQuery(t, stub, "balanceOf", address2)
	checkQuery(t, stub, "balanceOf", "BRCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	checkQuery(t, stub, "queryAllTokens", "")

	fmt.Println("[TestCreateTokenAndTransfer] fin")
	fmt.Println(time.Now().Unix())

}

func TestKey(t *testing.T) {
	// var encodedStr = "DK1kUGQjO/g5dD6TUWOo6ztmmA0sS4iGB2vgZm0l0rXY6T0L0zcijUQNA7axxz/8huavYXWNXKH7aop2kDGYWMlAIttYYtvwHx3i350HJxG1IQtn9bYc8TsWIWykUEUIaSZnCfercjYp1RXPnVZUxHvryuiPRLsGrE7fLmYPMxw="
	var singStr = "fiR1ZWjtF4iw2tuhfeGrA+Kj0KHkQQi1uusJfKLBSrCBqyb0+RVNkIDdZ/Vj/0H6bEwfHtTCZFIwIp2BRY2Kk9aUsDc4p0T/IOQ45gHPsJ4/ETy+3/tjndKdLRvQMcQPdM+tfxmDI+JOGFzj7NN3+Vmf3IRij1EVb1BhkJOzi/0="
	var pubPEM = `-----BEGIN TEST PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCY7YuXgEeGlEn9EuIIBDKzV+IY\nr9WpaVHMNJcNuJMLAXEga/ACWaKSmDnWaziArjHYBjgXVqMlk2PImHKMaf2pv7ch\nhAhes1pBAzNoVai+zEqJCii+67Cy3ql7nLN3+aFBh457tiL3UcHG/pbryYKsizt5\nmC29e7if9uI6/xzj9QIDAQAB\n-----END TEST PUBLIC KEY-----\n`

	util.SignatureVerification(pubPEM, singStr)
	// fmt.Println(err.Error())

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
