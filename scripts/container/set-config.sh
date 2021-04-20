ORG=doro
PROD_USER=$1
CHANNEL_NAME=$2

#change channel's config file format from json to pb
configtxlator proto_encode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config.json --type common.Config --output /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config.pb

#change modified channel's config file format from json to pb
configtxlator proto_encode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/modified_config.json --type common.Config --output /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/modified_config.pb

#extract the differences between existing channel's config and modified channel's config
configtxlator compute_update --channel_id ${CHANNEL_NAME} --original /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config.pb --updated /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig//modified_config.pb --output /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_update.pb

#change extracted config file formet from pb to json
configtxlator proto_decode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_update.pb --type common.ConfigUpdate | jq . >/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_update.json

#add header in extracted config file
echo '{"payload":{"header":{"channel_header":{"channel_id": "'${CHANNEL_NAME}'", "type":2}},"data":{"config_update":'$(cat /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_update.json)'}}}' | jq . >/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_update_in_envelope.json

#change final file format from json to pb
configtxlator proto_encode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_update_in_envelope.json --type common.Envelope --output /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_update_in_envelope.pb

#signature in final file
peer channel signconfigtx -f /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_update_in_envelope.pb

#  CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org${ORG}.com/peers/peer0.org${ORG}.com/msp/cacerts/ca-org${ORG}-com-7054.pem
#  CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org${ORG}.com/peers/peer0.org${ORG}.com/tls/server.key
#  CORE_PEER_LOCALMSPID=${ORG}MSP
#  CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org${ORG}.com/peers/peer0.org${ORG}.com/tls/server.crt
#  CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org${ORG}.com/users/Admin@org${ORG}.com/msp/
#  CORE_PEER_ADDRESS=peer0.org${ORG}.com:7051

#channel update with final pb file
peer channel update -f /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config_update_in_envelope.pb -o orderer0.orgorderer.com:37060 -c ${CHANNEL_NAME} --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orgorderer.com/msp/tlscacerts/ca-orgorderer-com-8054.pem

#delete created directory
rm -rf /opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig
