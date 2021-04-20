. scripts/network/conf

CHANNEL=$1
REMOVE_ORG=$2

#copy remove-org script in to cli container
docker cp ./scripts/container/remove-org.sh cli:/opt/gopath/src/github.com/hyperledger/fabric/peer

#run remove-org script
docker exec ${ENV_STR_PEER0} cli /bin/bash /opt/gopath/src/github.com/hyperledger/fabric/peer/remove-org.sh $CHANNEL $REMOVE_ORG

#delete remove-org script in cli container
docker exec cli rm -rf /opt/gopath/src/github.com/hyperledger/fabric/peer/remove-org.sh
