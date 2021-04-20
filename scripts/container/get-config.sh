ORG=doro
PROD_USER=$1
CHANNEL_NAME=$2

#create directory to update config
mkdir /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig

#get channel's config block
peer channel fetch config /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_block.pb -o orderer0.orgorderer.com:37060 -c ${CHANNEL_NAME} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orgorderer.com/msp/tlscacerts/ca-orgorderer-com-8054.pem

#change channel's config file format from pb to json
configtxlator proto_decode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_block.pb --type common.Block | jq .data.data[0].payload.data.config >/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config.json
chown $PROD_USER:$PROD_USER /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config.json
