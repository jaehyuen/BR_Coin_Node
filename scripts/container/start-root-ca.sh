#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

# Copy the root CA's signing certificate to the data directory to be used by others
fabric-ca-server init -b admin:adminpw 
mkdir -p /crypto-config/ca-certs
cp $FABRIC_CA_SERVER_HOME/ca-cert.pem /crypto-config/ca-certs/${FABRIC_CA_SERVER_CSR_CN}-cert.pem



# Add the custom orgs
for o in $PEER_ORGS; do
   aff=$aff"\n   $o: []"
done
aff="${aff#\\n   }"
sed -i "/affiliations:/a \\   $aff" \
   $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
   
# Start the root CA
touch /log/$FABRIC_CA_SERVER_CSR_HOSTS.log 
chmod 666 /log/$FABRIC_CA_SERVER_CSR_HOSTS.log 

fabric-ca-server start >> /log/$FABRIC_CA_SERVER_CSR_HOSTS.log 2>&1
