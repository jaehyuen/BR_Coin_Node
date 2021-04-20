. scripts/network/conf


echo "======================================="
echo "     Start to create brcoin-channel     "
echo "======================================="
CHANNEL_NAME=brcoin-channel
docker exec $ENV_STR_BRCOIN0 cli peer channel create -o $ORDERER_ADDRESS0 -c $CHANNEL_NAME -f ./channel-artifacts/$CHANNEL_NAME/$CHANNEL_NAME.tx --tls true --cafile $CA_FILE

echo "======================================="
echo "         End to create channel         "
echo "======================================="
