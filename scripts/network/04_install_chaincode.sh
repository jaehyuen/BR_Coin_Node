. scripts/network/conf

CHANNEL_NAME=$1
CHAINCODE_NAME=$2
CHAINCODE_VERSION=$3
CHAINCODE_PATH=github.com/chaincode/$CHAINCODE_NAME/go/
CHAINCODE_PACKAGE_FILE=${CHAINCODE_NAME}_v${CHAINCODE_VERSION}.tar.gz
CHAINCODE_PACKAGE_FILE_PATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/$CHAINCODE_PACKAGE_FILE



SEQ=`docker exec $ENV_STR_BRCOIN0 cli peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name $CHAINCODE_NAME  --cafile $CA_FILE`
SEQ=`echo ${SEQ%%, Endorsement*}`
SEQ=`echo ${SEQ#*Sequence: }`
SEQ=`expr $SEQ + 1`


###################################### BRCOIN0 ###########################################

echo "======================================="
echo "       start packaging chaincode       "
echo "======================================="

docker exec $ENV_STR_BRCOIN0 cli peer lifecycle chaincode package $CHAINCODE_PACKAGE_FILE_PATH --path $CHAINCODE_PATH --lang golang --label ${CHAINCODE_NAME}_${CHAINCODE_VERSION}

echo "======================================="
echo "  start buliding & install chaincode   "
echo "======================================="

docker exec $ENV_STR_BRCOIN0 cli peer lifecycle chaincode install $CHAINCODE_PACKAGE_FILE_PATH
docker exec $ENV_STR_BRCOIN1 cli peer lifecycle chaincode install $CHAINCODE_PACKAGE_FILE_PATH

echo "======================================="
echo " start checking installed chaincode id "
echo "======================================="

docker exec $ENV_STR_BRCOIN0 cli peer lifecycle chaincode queryinstalled

echo -e "Please Enter the chaincode ID labeled ${CHAINCODE_NAME}_${CHAINCODE_VERSION}"
read  CCID

echo "======================================="
echo " start approve a chaincode definition  "
echo "======================================="

docker exec $ENV_STR_BRCOIN0 cli peer lifecycle chaincode approveformyorg -o $ORDERER_ADDRESS0 --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --version $CHAINCODE_VERSION --package-id $CCID --sequence $SEQ --tls --cafile $CA_FILE 

echo "======================================="
echo "  start checking approved a chaincode  "
echo "======================================="

docker exec $ENV_STR_BRCOIN0 cli peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --version $CHAINCODE_VERSION --sequence $SEQ --tls --cafile $CA_FILE

echo "======================================="
echo "        start commit chaincode         "
echo "======================================="

docker exec $ENV_STR_BRCOIN0 cli peer lifecycle chaincode commit -o $ORDERER_ADDRESS0 --channelID $CHANNEL_NAME --name $CHAINCODE_NAME  --version $CHAINCODE_VERSION --sequence $SEQ --tls --cafile $CA_FILE 



