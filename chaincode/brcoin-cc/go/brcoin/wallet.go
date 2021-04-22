package brcoin

import (
	"fmt"
	"time"
	"util"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"structure"
)

func CreateWallet(stub shim.ChaincodeStubInterface, password string) (string, error) {
	var address string

	walletData := structure.BarakWallet{Regdate: time.Now().Unix(),
		Password: password,
		JobDate:  time.Now().Unix(),
		JobType:  "NewWallet",
		Nonce:    "util.MakeRandomString(40)",
		Balance:  []structure.BalanceInfo{structure.BalanceInfo{Balance: "0", TokenId: 0, UnlockDate: 0}}}

	var isSuccess = false

	// 새로운 지갑 주소만든다
	for !isSuccess {

		address = fmt.Sprintf("BRC%37s", util.MakeRandomString(37))
		_, err := stub.GetState(address)

		if err != nil {
			continue
		} else {
			isSuccess = true
			break
		}
	}

	if err := PutWallet(stub, address, walletData, "NewWallet", []string{address, password}); err != nil {
		return "", err
	}
	return address, nil
}
