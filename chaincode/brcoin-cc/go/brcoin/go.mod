module brcoin

go 1.15

require (
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20210319203922-6b661064d4d9
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/shopspring/decimal v1.2.0
	structure v0.0.0
	util v0.0.0
)

replace (
	structure v0.0.0 => ../structure
	util v0.0.0 => ../util
)
