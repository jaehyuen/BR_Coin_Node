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
	res := stub.MockInvoke("1", [][]byte{[]byte("query"), []byte(name)})

	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}

	fmt.Println("[checkQuery] result : ", string(res.Payload))

}

func checkInvoke(t *testing.T, stub *shimtest.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}

	var test, _ = strconv.Unquote(string(res.Payload))

	checkQuery(t, stub, test, "zz")
}
func TestCreateWallet(t *testing.T) {

	checkInvoke(t, stub, [][]byte{[]byte("createWallet"), []byte("12341234")})

}

func TestMain(m *testing.M) {
	setup()

	m.Run()
	shutdown()
}
