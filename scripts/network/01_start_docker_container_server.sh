echo "======================================="
echo "            start ca docker            "
echo "======================================="

docker-compose -f ./compose-files/docker-compose.yaml up -d ca.orgbrcoin.com ca.orgorderer.com 
sleep 2

echo "======================================="
echo "          start setup docker           "
echo "======================================="

docker-compose -f ./compose-files/docker-compose.yaml up -d setup

sleep 3

echo "======================================="
echo "           start peer docker           "
echo "======================================="

sleep 3
docker-compose -f ./compose-files/docker-compose.yaml up -d orderer0.orgorderer.com orderer1.orgorderer.com peer0.orgbrcoin.com peer1.orgbrcoin.com 

echo "======================================="
echo "           start other docker          "
echo "======================================="

docker-compose -f ./compose-files/docker-compose.yaml up -d cli

sleep 2

echo "                                       "
echo "                                       "
echo "                                       "
echo "                                       "

docker ps

echo "                                       "
echo "                                       "
echo "                                       "
echo "                                       "
