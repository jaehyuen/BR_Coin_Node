. scripts/network/conf

PROD_USER=$(id -g)

CHANNEL_NAME=$1

#copy modified channel's config file in to cli container
docker cp ./${CHANNEL_NAME}_config.json cli:/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}_updateconfig/modified_config.json

#copy set-config script in to cli container
docker cp ./scripts/container/set-config.sh cli:/opt/gopath/src/github.com/hyperledger/fabric/peer

#run set-config script
docker exec ${ENV_STR_PEER0} cli /bin/bash /opt/gopath/src/github.com/hyperledger/fabric/peer/set-config.sh $PROD_USER $CHANNEL_NAME
# docker exec ${ENV_STR_ORDERER0} cli /bin/bash /opt/gopath/src/github.com/hyperledger/fabric/peer/set-config.sh $user $CHANNEL_NAME

#delete get-config script in cli container
docker exec cli rm -rf /opt/gopath/src/github.com/hyperledger/fabric/peer/set-config.sh

#delete modified channel's config file
rm -rf ./${CHANNEL_NAME}_config.json
