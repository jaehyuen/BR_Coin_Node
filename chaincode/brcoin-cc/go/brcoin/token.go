package brcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"structure"
)

func CreateToken(stub shim.ChaincodeStubInterface, token *structure.Token) (string, error) {

	var err error
	var currNo int
	var valueByte, tokenByte []byte
	var reserveInfo structure.TokenReserve
	var reserveWallet structure.BarakWallet

	// data check
	if len(token.Symbol) < 1 {
		return "", errors.New(CODE0005 + " Symbol is empty")
	}

	if len(token.Name) < 1 {
		return "", errors.New(CODE0005 + " Name is empty")
	}

	if token.Decimal < 0 {
		return "", errors.New(CODE0005 + " The decimal number must be bigger then 0")
	}

	if token.Decimal > 8 {
		return "", errors.New(CODE0005 + " The decimal number must be less than 8")
	}

	if valueByte, err = stub.GetState("TOKEN_MAX_NO"); err != nil {
		return "", errors.New(CODE9999 + " Hyperledger internal error - " + err.Error())
	}

	if valueByte == nil {

		currNo = 0
	} else {
		currNo64, _ := strconv.ParseInt(string(valueByte), 10, 32)
		currNo = int(currNo64)
		currNo = currNo + 1
	}
	fmt.Println(currNo)
	token.Token = currNo
	token.JobDate = time.Now().Unix()
	token.CreateDate = time.Now().Unix()
	token.JobType = "CreateToken"

	if tokenByte, err = json.Marshal(token); err != nil {
		return "", errors.New(CODE0006 + " Invalid Data format")
	}

	token.JobArgs = string(tokenByte)

	if tokenByte, err = json.Marshal(token); err != nil {
		return "", errors.New(CODE0006 + " Invalid Data format")
	}

	if err = stub.PutState("TOKEN_DATA_"+strconv.Itoa(currNo), tokenByte); err != nil {
		return "", err
	}

	if len(token.Reserve) == 0 {

		if reserveWallet, err = GetWallet(stub, token.Owner); err != nil {
			return "", errors.New(CODE0003 + " Token reserve address " + token.Owner + " not found")
		}
		if currNo == 0 {
			reserveWallet.Balance[0].Balance = token.TotalSupply
		} else {
			reserveWallet.Balance = append(reserveWallet.Balance, structure.BalanceInfo{Balance: token.TotalSupply, TokenId: currNo, UnlockDate: 0})
		}

		if err = PutWallet(stub, token.Owner, reserveWallet, "TokenReserve", []string{token.Owner, token.Owner, token.TotalSupply, strconv.Itoa(currNo)}); err != nil {
			return "", err
		}
	} else {

		for _, reserveInfo = range token.Reserve {
			if reserveWallet, err = GetWallet(stub, reserveInfo.Address); err != nil {
				return "", errors.New(CODE0003 + " Token reserve address " + reserveInfo.Address + " not found")
			}
			if currNo == 0 {
				reserveWallet.Balance[0].Balance = reserveInfo.Value
			} else {
				reserveWallet.Balance = append(reserveWallet.Balance, structure.BalanceInfo{Balance: reserveInfo.Value, TokenId: currNo, UnlockDate: reserveInfo.UnlockDate})
			}

			if err = PutWallet(stub, reserveInfo.Address, reserveWallet, "TokenReserve", []string{token.Owner, reserveInfo.Address, reserveInfo.Value, strconv.Itoa(currNo)}); err != nil {
				return "", err
			}
		}
	}

	if err = stub.PutState("TOKEN_MAX_NO", []byte(strconv.Itoa(currNo))); err != nil {
		return "", err
	}
	return "", nil
}
