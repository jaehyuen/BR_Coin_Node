module brcoin-cc

go 1.15

require (
	brcoin v0.0.0
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20210319203922-6b661064d4d9
	github.com/hyperledger/fabric-protos-go v0.0.0-20210318103044-13fdee960194
	structure v0.0.0
	util v0.0.0 // indirect
)

replace (
	brcoin v0.0.0 => ./brcoin
	structure v0.0.0 => ./structure
	util v0.0.0 => ./util
)
