#!/usr/bin/env bash

echo "
*************************

copy reenrolli

*************************
"
set -e

cp -r ./crypto ./reenroll_crypto-config

sleep 2

ORDERER_ORG="orderer"

echo "***************************************************"
echo "###############  ORDERER ADMIN START ###################"
echo "***************************************************"
ORDERER_ADMIN_NAME=admin-orderer
ORDERER_ADMIN_PASS=adminpw
ORDERER_PORT=8054
ORDERER_CA_HOST=ca.orgorderer.com



mv $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/users/Admin@orgorderer.com/msp/keystore/*_sk $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/users/Admin@orgorderer.com/msp/keystore/server.key
./channel-artifacts/bin/fabric-ca-client reenroll -d -u http://${ORDERER_ADMIN_NAME}:${ORDERER_ADMIN_PASS}@${ORDERER_CA_HOST}:${ORDERER_PORT} -M $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/users/Admin@orgorderer.com/msp

rm -rf $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/users/Admin@orgorderer.com/msp/keystore/server.key

echo "***************************************************"
echo "###############  ORDERER ADMIN END ###################"
echo "***************************************************"

echo "***************************************************"
echo "###############  ORDERER ORG START ###################"
echo "***************************************************"

for ((i = 0; i < 3; i++)); do
   PORT=8054
   mv $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/msp/keystore/*_sk $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/msp/keystore/server.key

   ./channel-artifacts/bin/fabric-ca-client reenroll -d -u http://orderer${i}_${ORDERER_ORG}:orderer${i}_${ORDERER_ORG}pw@${ORDERER_CA_HOST}:${PORT} -M $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/msp

   mkdir -p $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/signcerts
   mkdir -p $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/keystore
   mv $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/server.crt $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/signcerts/cert.pem
   mv $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/server.key $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/keystore/

   ./channel-artifacts/bin/fabric-ca-client reenroll -d -u http://orderer${i}_${ORDERER_ORG}:orderer${i}_${ORDERER_ORG}pw@${ORDERER_CA_HOST}:${PORT} -M $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls --enrollment.profile tls --csr.hosts orderer${i}.org${ORDERER_ORG}.com

   rm -rf $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/msp/keystore/server.key
   mv $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/keystore/*_sk $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/server.key
   mv $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/signcerts/cert.pem $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls/server.crt

   cp $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/msp/signcerts/cert.pem $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/msp/admincerts/${ORDERER_ORG}_admin_cert.pem

   cd $PWD/reenroll_crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer${i}.org${ORDERER_ORG}.com/tls
   find . ! -path "./server.*" | cut -d "." -f2 | cut -d "/" -f2 | xargs rm -rf {} \;
   cd ../../../../../../

done
