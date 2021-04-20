. scripts/network/conf

docker exec $ENV_STR_BRCOIN0 cli peer channel join -b $CHANNEL_NAME.block
docker exec $ENV_STR_BRCOIN0 cli peer channel update -c $CHANNEL_NAME -o $ORDERER_ADDRESS0 -f ./channel-artifacts/$CHANNEL_NAME/${ORG}MSPanchors.tx --tls true --cafile $CA_FILE

echo "======================================="
echo "        End to join $CHANNEL_NAME      "
echo "       by peer0.orgbrcoin.com      "
echo "======================================="

docker exec $ENV_STR_BRCOIN1 cli peer channel join -b $CHANNEL_NAME.block

echo "======================================="
echo "        End to join $CHANNEL_NAME      "
echo "       by peer1.orgbrcoin.com      "
echo "======================================="





