#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
#
# The remainder of this file contains variables which typically would not be changed.
#
# All org names
ORGS="$ORDERER_ORGS $PEER_ORGS"

# initOrgVars <ORG>
function initOrgVars() {
    if [ $# -ne 1 ]; then
        echo "Usage: initOrgVars <ORG>"
        exit 1
    fi
    ORG=$1

    # Root CA admin identity
    ROOT_CA_ADMIN_USER=admin
    ROOT_CA_ADMIN_PASS=adminpw
    ROOT_CA_ADMIN_USER_PASS=${ROOT_CA_ADMIN_USER}:${ROOT_CA_ADMIN_PASS}
    # Admin identity for the org
    ADMIN_NAME=admin-${ORG}
    ADMIN_PASS=adminpw
    # Typical user identity for the org
    USER_NAME=user1
    USER_PASS=${USER_NAME}pw

    ROOT_CA_CERTFILE=/crypto-config/ca-certs/ca.org${ORG}.com-cert.pem
    ORG_MSP_ID=${ORG}MSP
    ORG_MSP_DIR=/root/orgs/${ORG}/msp
    ORG_ADMIN_CERT=${ORG_MSP_DIR}/admincerts/cert.pem
    ORG_ADMIN_HOME=/root/orgs/${ORG}/admin

    CA_NAME=ca-${ORG,,}
    PORT="CA_SERVER_PORT_${ORG}"
    eval PORT='$'$PORT

    CA_HOST_PORT=ca.org${ORG}.com:${PORT}
    CA_CHAINFILE=$ROOT_CA_CERTFILE
    CA_ADMIN_USER_PASS=${ROOT_CA_ADMIN_USER}:${ROOT_CA_ADMIN_PASS}
}
function initOrdererVars() {
    if [ $# -ne 2 ]; then
        echo "Usage: initOrdererVars <ORG> <NUM>"
        exit 1
    fi
    initOrgVars $1
    NUM=$2
    CA_HOST_PORT=ca.org${ORG}.com:${PORT}
    ORDERER_HOST=orderer${NUM}.org${ORG}.com
    ORDERER_NAME=orderer${NUM}.org${ORG}.com
    ORDERER_PASS=orderer${NUM}pw
    ENROLLMENT_URL=http://${ORDERER_NAME}:${ORDERER_PASS=}@${CA_HOST_PORT}
    ORDERER_PORT="ORDERER_PORT_${ORG}${NUM}"
    eval ORDERER_PORT='$'$ORDERER_PORT
    TLSDIR=/crypto-config/ordererOrganizations/${ORDERER_HOST:9}/orderers/$ORDERER_HOST/tls
    MSPDIR=/crypto-config/ordererOrganizations/${ORDERER_HOST:9}/orderers/$ORDERER_HOST/msp

}

# initPeerVars <ORG> <NUM>
function initPeerVars() {
    if [ $# -ne 2 ]; then
        echo "Usage: initPeerVars <ORG> <NUM>: $*"
        exit 1
    fi
    initOrgVars $1
    NUM=$2
    CA_HOST_PORT=ca.org${ORG}.com:${PORT}
    PEER_HOST=peer${NUM}.org${ORG}.com
    PEER_NAME=peer${NUM}.org${ORG}.com
    PEER_PASS=peer${NUM}_${ORG}pw
    PEER_PORT="PEER_PORT_${ORG}${NUM}"
    eval PEER_PORT='$'$PEER_PORT

    ENROLLMENT_URL=http://${PEER_NAME}:${PEER_PASS}@${CA_HOST_PORT}
    TLSDIR=/crypto-config/peerOrganizations/${PEER_HOST:6}/peers/$PEER_HOST/tls
    MSPDIR=/crypto-config/peerOrganizations/${PEER_HOST:6}/peers/$PEER_HOST/msp

}

# Wait for one or more files to exist
# Usage: dowait <what> <timeoutInSecs> <file> [<file> ...]
function dowait() {
    if [ $# -lt 3 ]; then
        echo "Usage: dowait: $*"
        exit 1
    fi
    local what=$1
    local secs=$2
    shift 2
    local starttime=$(date +%s)
    for file in $*; do
        until [ -f $file ]; do
            echo "Waiting for $what ..."
            sleep 1
            if [ "$(($(date +%s) - starttime))" -gt "$secs" ]; then
                echo "Failed waiting for $what ($file not found);"
                exit 1
            fi
        done
    done
    echo ""
}

# Wait for a process to begin to listen on a particular host and port
# Usage: waitPort <what> <timeoutInSecs> <errorLogFile> <host> <port>
function waitPort() {
    set +e
    local what=$1
    local secs=$2
    input=$3
    local host=${input%%:*}
    local port=${input##*:}
    nc -z $host $port >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo "Waiting for $what ..."
        local starttime=$(date +%s)
        while true; do
            sleep 1
            nc -z $host $port >/dev/null 2>&1
            if [ $? -eq 0 ]; then
                break
            fi
            if [ "$(($(date +%s) - starttime))" -gt "$secs" ]; then
                echo "Failed waiting for $what"
                exit 1
            fi
            echo -n "."
        done
        echo ""
    fi
    set -e
}

# Create the TLS directories of the MSP folder if they don't exist.
# The fabric-ca-client should do this.
function finishMSPSetup() {
    if [ $# -ne 1 ]; then
        echo "Usage: finishMSPSetup <targetMSPDIR>"
        exit 1
    fi
    if [ ! -d $1/tlscacerts ]; then
        mkdir $1/tlscacerts
        cp $1/cacerts/* $1/tlscacerts
        if [ -d $1/intermediatecerts ]; then
            mkdir $1/tlsintermediatecerts
            cp $1/intermediatecerts/* $1/tlsintermediatecerts
        fi
    fi
}

# Copy the org's admin cert into some target MSP directory
# This is only required if ADMINCERTS is enabled.
function copyAdminCert() {
    if [ $# -ne 2 ]; then
        echo "Usage: copyAdminCert <targetMSPDIR>"
        exit 1
    fi

    dstDir=$1/admincerts
    mkdir -p $dstDir
    dowait "$ORG administator to enroll" 60 $ORG_ADMIN_CERT
    cp $ORG_ADMIN_CERT $dstDir
}

# Switch to the current org's admin identity.  Enroll if not previously enrolled.
function switchToAdminIdentity() {
    if [ ! -d $ORG_ADMIN_HOME ]; then
        dowait "$CA_NAME to start" 60 $CA_CHAINFILE
        echo "Enrolling admin '$ADMIN_NAME' with $CA_HOST_PORT ..."
        export FABRIC_CA_CLIENT_HOME=$ORG_ADMIN_HOME
        export FABRIC_CA_CLIENT_TLS_CERTFILES=$CA_CHAINFILE
        fabric-ca-client enroll -d -u http://$ADMIN_NAME:$ADMIN_PASS@$CA_HOST_PORT
        # If admincerts are required in the MSP, copy the cert there now and to my local MSP also
        if [ $ADMINCERTS ]; then
            mkdir -p $(dirname "${ORG_ADMIN_CERT}")
            cp $ORG_ADMIN_HOME/msp/signcerts/* $ORG_ADMIN_CERT
            mkdir -p $ORG_ADMIN_HOME/msp/admincerts
            cp $ORG_ADMIN_HOME/msp/signcerts/* $ORG_ADMIN_HOME/msp/admincerts
            # local copy
            mkdir -p /crypto-config/admincerts
            cp $ORG_ADMIN_HOME/msp/signcerts/* /crypto-config/admincerts/${ORG}_admin_cert.pem
        fi
    fi
    export CORE_PEER_MSPCONFIGPATH=$ORG_ADMIN_HOME/msp
}

