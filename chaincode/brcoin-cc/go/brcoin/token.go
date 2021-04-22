package brcoin

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	"util"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/shopspring/decimal"

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
	token.TokeId = currNo
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

func TransferToken(stub shim.ChaincodeStubInterface, transfer *structure.Transfer) (string, error) {

	var err error
	var fromWallet, toWallet structure.BarakWallet
	var iUnlockDate int64

	//지갑 주소 유효성 검사
	if util.AddressValidation(transfer.FromAddr) {
		return "", errors.New(CODE0005 + " Invalid from address")
	}
	if util.AddressValidation(transfer.ToAddr) {
		return "", errors.New(CODE0005 + " Invalid to address")
	}
	if transfer.FromAddr == transfer.ToAddr {
		return "", errors.New(CODE0005 + " From address and to address must be different values")
	}

	//unlock date 숫자로 변경(부호가 있는 숫자)
	if iUnlockDate, err = strconv.ParseInt(transfer.UnlockDate, 10, 64); err != nil {
		return "", errors.New("1102,Invalid unlock date")
	}

	//보내는 지갑 데이터 조회
	if fromWallet, err = GetWallet(stub, transfer.FromAddr); err != nil {
		return "", err
	}

	//받는 지갑 데이터 조회
	if toWallet, err = GetWallet(stub, transfer.ToAddr); err != nil {
		return "", err
	}

	//토큰 이동
	if err = MoveToken(stub, &fromWallet, &toWallet, transfer.TokenId, transfer.Amount, iUnlockDate); err != nil {
		if strings.Index(err.Error(), "5000,") == 0 {
			return "", errors.New("5001,The balance of fromuser is insufficient")
		}
		return "", err
	}

	args := []string{transfer.ToAddr, transfer.FromAddr, transfer.Amount, transfer.TokenId, transfer.UnlockDate}
	//보낸 지갑 정보 저장
	if err = PutWallet(stub, transfer.FromAddr, fromWallet, "transfer", args); err != nil {
		return "", err
	}
	//받은 지갑 정보 저장
	if err = PutWallet(stub, transfer.ToAddr, toWallet, "receive", args); err != nil {
		return "", err
	}

	return "", nil
}

// MoveToken 잔액을 다른 Wallet 로 이동
func MoveToken(stub shim.ChaincodeStubInterface, fromWallet *structure.BarakWallet, toWallet *structure.BarakWallet, TokenId string, amount string, unlockDate int64) error {
	var err error
	var subtractAmount, fromCoin, toCoin, addAmount decimal.Decimal
	var toIndex int
	var balanceTemp []structure.BalanceInfo
	var nowTime int64
	var tokenData structure.Token

	//현재시간과 거래 정지 시간 비교
	nowTime = time.Now().Unix()
	if unlockDate < nowTime {
		unlockDate = 0
	}

	//토큰 조회
	if tokenData, err = GetToken(stub, TokenId); err != nil {
		return errors.New(CODE9999 + " Hyperledger internal error - " + err.Error())
	}

	// 소수점 자리수 체크
	strTmp := strings.Split(amount, ".")
	if len(strTmp[len(strTmp)-1]) > tokenData.Decimal {

		return errors.New(CODE0007 + " The decimal places are bigger than " + strconv.Itoa(tokenData.Decimal))
	}

	if subtractAmount, err = util.ParsePositive(amount); err != nil {
		return errors.New(CODE0007 + " Amount must be an integer string")
	}
	// 추가, 삭제 코인량 초기화
	addAmount = subtractAmount

	isBalanceClean := false

	//보내는 지갑에서 해당 도큰 잔고 확인
	for index, element := range fromWallet.Balance {

		if element.TokenId != tokenData.TokeId {
			continue
		}

		//현재시간 보다 거래 정지 정지 시간이 더크면 다음 포문 실행
		if nowTime < element.UnlockDate {
			continue
		}

		//보낼수 았는 코인량
		if fromCoin, err = decimal.NewFromString(element.Balance); err != nil {
			continue
		}

		//fromCoin(보낼수 있는 코인량)이 subtractAmount(전송할 코인량) 비교
		if fromCoin.Cmp(subtractAmount) < 0 {
			//fromCoin(보낼수 있는 코인량)이 subtractAmount(전송할 코인량)작으면
			// subtractAmount - fromCoin 계산
			subtractAmount = subtractAmount.Sub(fromCoin).Round(int32(tokenData.Decimal))

			//Balance를 0으로 초기화
			fromWallet.Balance[index].Balance = "0"

			//유효한 토큰 아이디면 isBalanceClean true로 변경
			if tokenData.TokeId > 0 {
				isBalanceClean = true
			}
			continue

		} else {
			//fromCoin(보낼수 있는 코인량)이 subtractAmount(전송할 코인량)보다 크거나 같으면
			//Balance 값은 fromCoin - subtractAmount
			fromWallet.Balance[index].Balance = fromCoin.Sub(subtractAmount).Round(int32(tokenData.Decimal)).String()

			//subtractAmount 값은 0?
			subtractAmount = subtractAmount.Sub(subtractAmount).Round(int32(tokenData.Decimal))
			break
		}
	}

	//머지? 잔고가 텅텅
	if isBalanceClean {
		for _, element := range fromWallet.Balance {
			if element.TokenId > 0 && element.Balance == "0" {
				continue
			}
			balanceTemp = append(balanceTemp, element)
		}
		fromWallet.Balance = balanceTemp
	}

	//subtractAmount 가 0보다 클때 true (잔고 부족)
	if subtractAmount.IsPositive() {
		return errors.New(CODE0008)
	}

	//받을 지갑 잔고에 대한 인덱스
	toIndex = -1
	// 받는 코인을 0으로 초기화
	toCoin = decimal.Zero

	//코인을 받는 지갑의 잔고 포문
	for index, element := range toWallet.Balance {
		//토큰 아이디가 같으면
		if element.TokenId == tokenData.TokeId {
			//거래 정지시간  데이트가 같으면
			if element.UnlockDate == unlockDate {

				//현재 잔고를 toCoin에 대입
				toCoin, _ = decimal.NewFromString(element.Balance)

				//현재 인덱스로 설정
				toIndex = index
				break
			}
		}
	}

	//addAmount 만큼 더하고 소수점 자르기
	toCoin = toCoin.Add(addAmount).Round(int32(tokenData.Decimal))

	if toIndex == -1 {
		//지갑에 해당 잔고 (BalanceInfo)가 없을떄
		//iUnlockDate가 0보다 크면
		if unlockDate > 0 {
			//잔고를 추가
			toWallet.Balance = append(toWallet.Balance, structure.BalanceInfo{Balance: toCoin.String(), TokenId: tokenData.TokeId, UnlockDate: unlockDate})
		} else {
			toWallet.Balance = append(toWallet.Balance, structure.BalanceInfo{Balance: toCoin.String(), TokenId: tokenData.TokeId, UnlockDate: 0})
		}
	} else {

		//지갑에 해당 잔고 (BalanceInfo)가  있으면 toCoin, iUnlockDate 업데이트
		toWallet.Balance[toIndex].Balance = toCoin.String()
		if unlockDate > 0 {
			toWallet.Balance[toIndex].UnlockDate = unlockDate
		}
	}
	return nil
}
