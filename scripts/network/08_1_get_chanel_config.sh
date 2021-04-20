. scripts/network/conf

PROD_USER=$(id -g)

CHANNEL_NAME=$1

#copy get-config script in to cli container
docker cp ./scripts/container/get-config.sh cli:/opt/gopath/src/github.com/hyperledger/fabric/peer

#run get-config script
docker exec ${ENV_STR_PEER0} cli /bin/bash /opt/gopath/src/github.com/hyperledger/fabric/peer/get-config.sh $PROD_USER $CHANNEL_NAME

#get channel's config file in cli container
docker cp cli:/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/config.json ./${CHANNEL_NAME}_config.json

#delete get-config script in cli container
docker exec cli rm -rf /opt/gopath/src/github.com/hyperledger/fabric/peer/get-config.sh
