ORG=doro
REMOVE_ORG=$2
CHANNEL_NAME=$1

#create directory to config remove org
mkdir /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config

#get channel's config block
peer channel fetch config /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_block.pb -o orderer0.orgorderer.com:37060 -c $CHANNEL_NAME --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orgorderer.com/msp/tlscacerts/ca-orgorderer-com-8054.pem

#change channel's config file format from pb to json
configtxlator proto_decode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_block.pb --type common.Block | jq .data.data[0].payload.data.config >/opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config.json

#remove org in channel' config
jq 'del(.channel_group.groups.Application.groups.'$REMOVE_ORG')' /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config.json >/opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/remove_$REMOVE_ORG.json

#change channel's config file format from json to pb
configtxlator proto_encode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config.json --type common.Config --output /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config.pb

#change removed org's config file format from json to pb
configtxlator proto_encode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/remove_$REMOVE_ORG.json --type common.Config --output /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/remove_$REMOVE_ORG.pb

#extract the differences between existing channel's config and removed org channel's config
configtxlator compute_update --channel_id $CHANNEL_NAME --original /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config.pb --updated /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/remove_$REMOVE_ORG.pb --output /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_update.pb

#change extracted config file formet from pb to json
configtxlator proto_decode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_update.pb --type common.ConfigUpdate | jq . >/opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_update.json

#add header in extracted config file
echo '{"payload":{"header":{"channel_header":{"channel_id": "'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_update.json)'}}}' | jq . >/opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_update_in_envelope.json

#change final file format from json to pb
configtxlator proto_encode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_update_in_envelope.json --type common.Envelope --output /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_update_in_envelope.pb

#channel update with final pb file
peer channel update -f /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config/config_update_in_envelope.pb -o orderer0.orgorderer.com:37060 -c $CHANNEL_NAME --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orgorderer.com/msp/tlscacerts/ca-orgorderer-com-8054.pem

#delete created directory
rm -rf /opt/gopath/src/github.com/hyperledger/fabric/peer/remove_config
