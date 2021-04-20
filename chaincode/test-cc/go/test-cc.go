package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid" //cid 테스트??
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
)

const CODE9999 = "ChaincodectxInterface Error - "
const CODE9991 = "Data already exist - "
const CODE9992 = "Data does not exist - "
const CODE9993 = "Empty key - "

const PRIVATE_COLLECCTION_NAME = "collectionTest"
const TRANSIENTKEY_PRIVATE = "priTest"

const IMPLICIT_COLLECCTION_NAME = "_implicit_org_apeerMSP"
const TRANSIENTKEY_IMPLICIT = "priImpTest"

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

type testDataWithKey struct {
	Key      string   `json:"key"`
	TestData testData `json:"testData"`
}

type testData struct {
	TestValue1 string `json:"testValue1"`
	TestValue2 string `json:"testValue2"`
	TestValue3 string `json:"testValue3"`
	TestValue4 string `json:"testValue4"`
}

type priTestDataWithKey struct {
	Key         string      `json:"key"`
	PriTestData priTestData `json:"priTestData"`
}

type priTestData struct {
	PriTestValue1 string `json:"priTestValue1"`
	PriTestValue2 string `json:"priTestValue2"`
	PriTestValue3 string `json:"priTestValue3"`
	PriTestValue4 string `json:"priTestValue4"`
}

type priImpTestDataWithKey struct {
	Key            string         `json:"key"`
	PriImpTestData priImpTestData `json:"priImpTestData"`
}

type priImpTestData struct {
	PriImpTestValue1 string `json:"priImpTestValue1"`
	PriImpTestValue2 string `json:"priImpTestValue2"`
	PriImpTestValue3 string `json:"priImpTestValue3"`
	PriImpTestValue4 string `json:"priImpTestValue4"`
}

type jsonResponse struct {
	Key           string `json:"key"`
	ResultFlag    bool   `json:"resultFlag"`
	ResultCode    string `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

/*
 * 테스트 인보크
 * args[0]: []testDataWithKey
 * detail args[0]: [{Key: ... , TestData:{"TestValue1":"1","TestValue2":"00","TestValue3":"zz","TestValue4":"00"}}]
 */

func (s *SmartContract) InvokeTest(ctx contractapi.TransactionContextInterface, tdkList []testDataWithKey) (string, error) {

	var response []jsonResponse
	var isSuccess = false

	var creator, _ = ctx.GetStub().GetCreator()
	log.Println("ctx.GetStub().GetCreator() --> creator")
	log.Println(string(creator))

	var id, _ = cid.GetID(ctx.GetStub())
	log.Println("cid.GetID(ctx.GetStub()) --> id")
	log.Println(id)

	var mspid, _ = cid.GetMSPID(ctx.GetStub())
	log.Println("cid.GetMSPID(ctx.GetStub()) --> mspid")
	log.Println(mspid)

	// var cert, err := cid.GetX509Certificate(ctx.GetStub())
	// log.Println("cid.GetX509Certificate(ctx.GetStub()) --> cert")
	// log.Println(cert)

	// 인증서별 인보크
	for i := 0; i < len(tdkList); i++ {

		{
			if res, ok := _writeData(ctx, tdkList[i].Key, &tdkList[i].TestData); !ok {
				response = append(response, res)

			} else {
				response = append(response, res)
				isSuccess = true
			}
		}
	}

	log.Println("%+v", response)

	var dataBytes []byte
	dataBytes, _ = json.Marshal(response)
	log.Println(string(dataBytes))

	if isSuccess {
		return string(dataBytes), nil
	} else {
		return "", errors.New(string(dataBytes))
	}
}

/*
 * 테스트 조회
 * args[0]: []testDataWithKey
 * detail args[0]: [{Key: ...}, {Key: ...}]
 */

func (s *SmartContract) QueryTest(ctx contractapi.TransactionContextInterface, tdkList []testDataWithKey) (string, error) {

	var response []jsonResponse
	var isSuccess = false

	// 인증서별 쿼리
	for i := 0; i < len(tdkList); i++ {

		testAsBytes, err := ctx.GetStub().GetState(tdkList[i].Key)
		//public 데이터 조회

		if err != nil {
			log.Println("Fail query ()")
			response = append(response, _makeJson(tdkList[i].Key, false, "9999", CODE9999+"GetState() : "+err.Error()))
			continue
		}

		if testAsBytes == nil {
			log.Println("Fail query ()")
			response = append(response, _makeJson(tdkList[i].Key, false, "9992", CODE9992+"key : "+tdkList[i].Key))
			continue

		} else {
			log.Println("Success query ()")
			response = append(response, _makeJson(tdkList[i].Key, true, "0000", string(testAsBytes)))
			isSuccess = true
		}

	}

	var dataBytes []byte
	dataBytes, _ = json.Marshal(response)

	if isSuccess {
		return string(dataBytes), nil
	} else {
		return "", errors.New(string(dataBytes))
	}
}

/*
 * 테스트 PDC 인보크 (Collection Config)
 * transient: priTestDataWithKey
 * detail transient: {Key: ... , TestData:{"PriTestValue1":"1","PriTestValue2":"00","PriTestValue3":"zz","PriTestValue4":"00"}} by base64
 */

func (s *SmartContract) PriInvokeTest(ctx contractapi.TransactionContextInterface) (string, error) {

	var transMap map[string][]byte
	var err error

	var response []jsonResponse
	var priTdk priTestDataWithKey
	var dataBytes []byte

	if transMap, err = ctx.GetStub().GetTransient(); err != nil {
		response = append(response, _makeJson("", false, "9999", CODE9999+"GetTransient() : "+err.Error()))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//Transient에서 가져온 값의 key가 존재하는지 확인합니다.
	if _, ok := transMap[TRANSIENTKEY_PRIVATE]; !ok {
		response = append(response, _makeJson("", false, "9999", CODE9999+"UserData must be a key in the transient map"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//privPersonal 키에 속한 value가 존재하는지 확인합니다.
	if len(transMap[TRANSIENTKEY_PRIVATE]) == 0 {
		response = append(response, _makeJson("", false, "9999", CODE9999+"Value in the transient map must be a non-empty Json string"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//고객(private) 데이터 unmarshal
	if err = json.Unmarshal(transMap[TRANSIENTKEY_PRIVATE], &priTdk); err != nil {
		response = append(response, _makeJson("", false, "9999", CODE9999+"Value in the transient map must be a non-empty Json string"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	if res, ok := _writeData(ctx, priTdk.Key, &priTdk.PriTestData); !ok {
		response = append(response, res)
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	} else {
		response = append(response, res)
	}

	dataBytes, _ = json.Marshal(response)

	return string(dataBytes), nil

}

/*
 * 테스트 PDC 조회 (Collection Config)
 * transient: priTestDataWithKey
 * detail transient: {Key: ...}
 */

func (s *SmartContract) PriQueryTest(ctx contractapi.TransactionContextInterface) (string, error) {

	var transMap map[string][]byte
	var err error

	var response []jsonResponse
	var dataBytes []byte
	var priTdk priTestDataWithKey

	if transMap, err = ctx.GetStub().GetTransient(); err != nil {
		response = append(response, _makeJson("", false, "9999", CODE9999+"GetTransient() : "+err.Error()))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//Transient에서 가져온 값의 key가 존재하는지 확인합니다.
	if _, ok := transMap[TRANSIENTKEY_PRIVATE]; !ok {
		response = append(response, _makeJson("", false, "9999", CODE9999+"UserData must be a key in the transient map"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//privPersonal 키에 속한 value가 존재하는지 확인합니다.
	if len(transMap[TRANSIENTKEY_PRIVATE]) == 0 {
		response = append(response, _makeJson("", false, "9999", CODE9999+"Value in the transient map must be a non-empty Json string"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//고객(private) 데이터 unmarshal
	if err = json.Unmarshal(transMap[TRANSIENTKEY_PRIVATE], &priTdk); err != nil {
		response = append(response, _makeJson("", false, "9999", CODE9999+"Value in the transient map must be a non-empty Json string"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	testAsBytes, err := ctx.GetStub().GetPrivateData(PRIVATE_COLLECCTION_NAME, priTdk.Key)

	if err != nil {
		log.Println("Fail query ()")
		response = append(response, _makeJson(priTdk.Key, false, "9999", CODE9999+"GetState() : "+err.Error()))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	if testAsBytes == nil {
		log.Println("Fail query ()")
		response = append(response, _makeJson(priTdk.Key, false, "9992", CODE9992+"key : "+priTdk.Key))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	} else {
		log.Println("Success query ()")
		response = append(response, _makeJson(priTdk.Key, true, "0000", string(testAsBytes)))
	}

	dataBytes, _ = json.Marshal(response)

	return string(dataBytes), nil

}

/*
 * 테스트 PDC 인보크 (Implicit Collection)
 * transient: priImpTestDataWithKey
 * detail transient: {Key: ... , TestData:{"PriImpTestValue1":"1","PriImpTestValue2":"00","PriImpTestValue3":"zz","PriImpTestValue4":"00"}} by base64
 */

func (s *SmartContract) PriImpInvokeTest(ctx contractapi.TransactionContextInterface) (string, error) {

	var transMap map[string][]byte
	var err error

	var response []jsonResponse
	var priImpTdk priImpTestDataWithKey
	var dataBytes []byte

	if transMap, err = ctx.GetStub().GetTransient(); err != nil {
		response = append(response, _makeJson("", false, "9999", CODE9999+"GetTransient() : "+err.Error()))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//Transient에서 가져온 값의 key가 존재하는지 확인합니다.
	if _, ok := transMap[TRANSIENTKEY_IMPLICIT]; !ok {
		response = append(response, _makeJson("", false, "9999", CODE9999+"UserData must be a key in the transient map"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//privPersonal 키에 속한 value가 존재하는지 확인합니다.
	if len(transMap[TRANSIENTKEY_IMPLICIT]) == 0 {
		response = append(response, _makeJson("", false, "9999", CODE9999+"Value in the transient map must be a non-empty Json string"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//고객(private) 데이터 unmarshal
	if err = json.Unmarshal(transMap[TRANSIENTKEY_IMPLICIT], &priImpTdk); err != nil {
		response = append(response, _makeJson("", false, "9999", CODE9999+"Value in the transient map must be a non-empty Json string"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	if res, ok := _writeData(ctx, priImpTdk.Key, &priImpTdk.PriImpTestData); !ok {
		response = append(response, res)
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	} else {
		response = append(response, res)
	}

	dataBytes, _ = json.Marshal(response)

	return string(dataBytes), nil

}

/*
 * 테스트 PDC 조회 (Implicit Collection)
 * transient: priImpTestDataWithKey
 * detail transient: {Key: ...}
 */

func (s *SmartContract) PriImpQueryTest(ctx contractapi.TransactionContextInterface) (string, error) {

	var transMap map[string][]byte
	var err error

	var response []jsonResponse
	var dataBytes []byte
	var priImpTdk priImpTestDataWithKey

	if transMap, err = ctx.GetStub().GetTransient(); err != nil {
		response = append(response, _makeJson("", false, "9999", CODE9999+"GetTransient() : "+err.Error()))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//Transient에서 가져온 값의 key가 존재하는지 확인합니다.
	if _, ok := transMap[TRANSIENTKEY_IMPLICIT]; !ok {
		response = append(response, _makeJson("", false, "9999", CODE9999+"UserData must be a key in the transient map"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//privPersonal 키에 속한 value가 존재하는지 확인합니다.
	if len(transMap[TRANSIENTKEY_IMPLICIT]) == 0 {
		response = append(response, _makeJson("", false, "9999", CODE9999+"Value in the transient map must be a non-empty Json string"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	//고객(private) 데이터 unmarshal
	if err = json.Unmarshal(transMap[TRANSIENTKEY_IMPLICIT], &priImpTdk); err != nil {
		response = append(response, _makeJson("", false, "9999", CODE9999+"Value in the transient map must be a non-empty Json string"))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	testAsBytes, err := ctx.GetStub().GetPrivateData(PRIVATE_COLLECCTION_NAME, priImpTdk.Key)

	if err != nil {
		log.Println("Fail query ()")
		response = append(response, _makeJson(priImpTdk.Key, false, "9999", CODE9999+"GetState() : "+err.Error()))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	}

	if testAsBytes == nil {
		log.Println("Fail query ()")
		response = append(response, _makeJson(priImpTdk.Key, false, "9992", CODE9992+"key : "+priImpTdk.Key))
		dataBytes, _ = json.Marshal(response)
		return "", errors.New(string(dataBytes))

	} else {
		log.Println("Success query ()")
		response = append(response, _makeJson(priImpTdk.Key, true, "0000", string(testAsBytes)))
	}

	dataBytes, _ = json.Marshal(response)

	return string(dataBytes), nil

}

// 등록함수
func _writeData(ctx contractapi.TransactionContextInterface, key string, ledger interface{}) (res jsonResponse, ok bool) {
	var err error

	switch ledger.(type) {
	case *testData:

		var testDataBytes []byte
		testLedger := ledger.(*testData)
		//public 데이터 변수 저장
		testDataBytes, err = json.Marshal(testLedger)
		if err != nil {
			log.Println("Fail _writeData ()")
			return _makeJson(key, false, "9999", CODE9999+"_insertCertData() : "+err.Error()), false
		}

		//public 데이터 등록
		err = ctx.GetStub().PutState(key, testDataBytes)
		if err != nil {
			log.Println("Fail _writeData ()")
			return _makeJson(key, false, "9999", CODE9999+"_writeCertData() : "+err.Error()), false
		}
		log.Println("Success _writeData ()")
		return _makeJson(key, true, "0000", "Success"), true

	case *priTestData:

		var priTestDataBytes []byte
		priTestLedger := ledger.(*priTestData)

		//public 데이터 변수 저장
		priTestDataBytes, err = json.Marshal(priTestLedger)
		if err != nil {
			log.Println("Fail _writeData ()")
			return _makeJson(key, false, "9999", CODE9999+"_insertCertData() : "+err.Error()), false
		}

		//public 데이터 등록
		err = ctx.GetStub().PutPrivateData(PRIVATE_COLLECCTION_NAME, key, priTestDataBytes)
		if err != nil {
			log.Println("Fail _writeData ()")
			return _makeJson(key, false, "9999", CODE9999+"_writeCertData() : "+err.Error()), false
		}
		log.Println("Success _writeData ()")
		return _makeJson(key, true, "0000", "Success"), true

	case *priImpTestData:

		var priImpTestDataBytes []byte
		priImpTestLedger := ledger.(*priImpTestData)

		//public 데이터 변수 저장
		priImpTestDataBytes, err = json.Marshal(priImpTestLedger)
		if err != nil {
			log.Println("Fail _writeData ()")
			return _makeJson(key, false, "9999", CODE9999+"_insertCertData() : "+err.Error()), false
		}

		//public 데이터 등록
		err = ctx.GetStub().PutPrivateData(IMPLICIT_COLLECCTION_NAME, key, priImpTestDataBytes)
		if err != nil {
			log.Println("Fail _writeData ()")
			return _makeJson(key, false, "9999", CODE9999+"_writeCertData() : "+err.Error()), false
		}
		log.Println("Success _writeData ()")
		return _makeJson(key, true, "0000", "Success"), true

	default:
		return _makeJson(key, false, "9999", CODE9999+"_writeCertData() : "), false
	}

}

//return 데이터 생성 함수
func _makeJson(key string, resultFlag bool, resultCode string, resultMessage string) (res jsonResponse) {

	var response jsonResponse
	response.Key = key
	response.ResultFlag = resultFlag
	response.ResultCode = resultCode
	response.ResultMessage = resultMessage

	return response
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create doro-cc chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting doro-cc chaincode: %s", err.Error())
	}
}
