#!/bin/bash

function main() {
    echo "[main] Start main function"
    registerIdentities
    getCACerts
    enrollOrderer
    enrollPeer
    makeConfigTxYaml
    generateChannelArtifacts
    echo "[main] Finish main function"

    chown -R $PROD_USER:$PROD_USER /crypto-config
    chown -R $PROD_USER:$PROD_USER /root

}

# Enroll the CA administrator
function enrollCAAdmin() {
    waitPort "$CA_NAME to start" 20 $CA_HOST_PORT
    echo "Enrolling with $CA_NAME as bootstrap identity ..."
    export FABRIC_CA_CLIENT_HOME=$HOME/cas/$CA_NAME
    export FABRIC_CA_CLIENT_TLS_CERTFILES=$CA_CHAINFILE
    fabric-ca-client enroll -d -u http://$CA_ADMIN_USER_PASS@$CA_HOST_PORT
}

function registerIdentities() {
    echo "[registerIdentities] Start registerIdentities function"
    registerOrdererIdentities
    registerPeerIdentities
    echo "[registerIdentities] Finish registerIdentities function"
}

function registerOrdererIdentities() {
    echo "[registerOrdererIdentities] Start registerOrdererIdentities function"

    for ORG in $ORDERER_ORGS; do
        initOrgVars $ORG
        enrollCAAdmin
        local COUNT=1
        while [[ "$COUNT" -le $NUM_ORDERERS ]]; do
            initOrdererVars $ORG $((COUNT - 1))
            echo "[registerOrdererIdentities] Registering orderer $ORDERER_NAME"
            fabric-ca-client register -d --id.name $ORDERER_NAME --id.secret $ORDERER_PASS --id.type orderer
            COUNT=$((COUNT + 1))
        done
        echo "[registerOrdererIdentities] Registering $ORG admin identity"
        fabric-ca-client register -d --id.name $ADMIN_NAME --id.secret $ADMIN_PASS --id.attrs "admin=true:ecert"
    done

    echo "[registerOrdererIdentities] Finish registerOrdererIdentities function"
}

function registerPeerIdentities() {
    echo "[registerPeerIdentities] Start registerPeerIdentities function"
    for ORG in $PEER_ORGS; do
        initOrgVars $ORG
        enrollCAAdmin
        local COUNT=1
        while [[ "$COUNT" -le $NUM_PEERS ]]; do
            initPeerVars $ORG $((COUNT - 1))
            echo "[registerPeerIdentities] Registering peer $ORDERER_NAME"
            fabric-ca-client register -d --id.name $PEER_NAME --id.secret $PEER_PASS --id.type peer
            COUNT=$((COUNT + 1))
        done
        echo "[registerPeerIdentities] Registering $ORG admin identity"
        fabric-ca-client register -d --id.name $ADMIN_NAME --id.secret $ADMIN_PASS --id.attrs "hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert,abac.init=true:ecert"

    done
    echo "[registerPeerIdentities] Finish registerPeerIdentities function"

}

function getCACerts() {
    echo "[getCACerts] Start getCACerts function"
    for ORG in $ORGS; do
        initOrgVars $ORG
        echo "[getCACerts] Getting CA certs for organization $ORG and storing in $ORG_MSP_DIR"
        export FABRIC_CA_CLIENT_TLS_CERTFILES=$CA_CHAINFILE
        fabric-ca-client getcacert -d -u http://$CA_HOST_PORT -M $ORG_MSP_DIR
        finishMSPSetup $ORG_MSP_DIR

        if [ $ADMINCERTS ]; then
            switchToAdminIdentity
        fi
    done
    echo "[getCACerts] Finish getCACerts function"
}

function enrollOrderer() {

    echo "[enrollOrderer] Start enrollOrderer function"
    for ORG in $ORDERER_ORGS; do
        for ((i = 0; i < $NUM_ORDERERS; i++)); do
            initOrdererVars $ORG $i

            echo "[enrollOrderer] enrolling orderer $ORDERER_NAME"
            fabric-ca-client enroll -d --enrollment.profile tls -u $ENROLLMENT_URL -M /tmp/$ORG/peer$i/tls --csr.hosts $ORDERER_HOST

            mkdir -p $TLSDIR
            cp /tmp/$ORG/peer$i/tls/signcerts/* $TLSDIR/server.crt
            cp /tmp/$ORG/peer$i/tls/keystore/* $TLSDIR/server.key

            fabric-ca-client enroll -d -u $ENROLLMENT_URL -M $MSPDIR

            mkdir -p /crypto-config/ordererOrganizations/${ORDERER_HOST:9}/users/Admin@${ORDERER_HOST:9}/
            mkdir -p /crypto-config/ordererOrganizations/${ORDERER_HOST:9}/msp
            cp -r /root/orgs/$ORG/msp/ /crypto-config/ordererOrganizations/${ORDERER_HOST:9}/

            cp -r /root/orgs/$ORG/admin/msp/ /crypto-config/ordererOrganizations/${ORDERER_HOST:9}/users/Admin@${ORDERER_HOST:9}/

            finishMSPSetup $MSPDIR
            copyAdminCert $MSPDIR $ORDERER_HOST
        done
    done
    echo "[enrollOrderer] Finish enrollOrderer function"
}

function enrollPeer() {

    echo "[enrollPeer] Start enrollPeer function"
    for ORG in $PEER_ORGS; do
        for ((i = 0; i < $NUM_PEERS; i++)); do
            initPeerVars $ORG $i
            echo "[enrollPeer] Registering peer $ORDERER_NAME"
            fabric-ca-client enroll -d --enrollment.profile tls -u $ENROLLMENT_URL -M /tmp/$ORG/peer$i/tls --csr.hosts $PEER_HOST

            mkdir -p $TLSDIR
            cp /tmp/$ORG/peer$i/tls/signcerts/* $TLSDIR/server.crt
            cp /tmp/$ORG/peer$i/tls/keystore/* $TLSDIR/server.key

            fabric-ca-client enroll -d -u $ENROLLMENT_URL -M $MSPDIR

            mkdir -p /crypto-config/peerOrganizations/${PEER_HOST:6}/users/Admin@${PEER_HOST:6}
            mkdir -p /crypto-config/peerOrganizations/${PEER_HOST:6}/msp
            cp -r /root/orgs/$ORG/msp/ /crypto-config/peerOrganizations/${PEER_HOST:6}/
            cp -r /root/orgs/$ORG/admin/msp/ /crypto-config/peerOrganizations/${PEER_HOST:6}/users/Admin@${PEER_HOST:6}/

            finishMSPSetup $MSPDIR
            copyAdminCert $MSPDIR $PEER_HOST
        done
    done
    echo "[enrollPeer] Finish enrollPeer function"
}

function makeConfigTxYaml() {
    {
        echo "
Organizations:"

        for ORG in $ORDERER_ORGS; do
            printOrdererOrg $ORG
        done

        for ORG in $PEER_ORGS; do
            printPeerOrg $ORG 0
        done

        echo " 
Capabilities:
    Channel: &ChannelCapabilities
        V2_0: true
        V1_3: false
        V1_1: false

    Orderer: &OrdererCapabilities
        V2_0: true
        V1_1: false

    Application: &ApplicationCapabilities
        V2_0: true
        V1_3: false
        V1_2: false
        V1_1: false

Application: &ApplicationDefaults
    Organizations:
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: \"ANY Readers\"
        Writers:
            Type: ImplicitMeta
            Rule: \"ANY Writers\"
        Admins:
            Type: Signature
            Rule: \"OR('apeerMSP.admin')\"
        LifecycleEndorsement:
            Type: ImplicitMeta
            Rule: \"ANY Endorsement\"
        Endorsement:
            Type: ImplicitMeta
            Rule: \"ANY Endorsement\"
            

    Capabilities:
        <<: *ApplicationCapabilities

Orderer: &OrdererDefaults

    BatchTimeout: 1s

    BatchSize:
        MaxMessageCount: 20
        AbsoluteMaxBytes: 80 KB
        PreferredMaxBytes: 20 KB

    Organizations:
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: \"ANY Readers\"
        Writers:
            Type: ImplicitMeta
            Rule: \"ANY Writers\"
        Admins:
            Type: Signature
            Rule: \"OR('apeerMSP.admin')\"

        BlockValidation:
            Type: ImplicitMeta
            Rule: \"ANY Writers\"

Channel: &ChannelDefaults
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: \"ANY Readers\"
        Writers:
            Type: ImplicitMeta
            Rule: \"ANY Writers\"
        Admins:
            Type: Signature
            Rule: \"OR('apeerMSP.admin')\"

    Capabilities:
        <<: *ChannelCapabilities


Profiles:"
        echo "
    OneOrgChannel:
        Consortium: SampleConsortium
        <<: *ChannelDefaults
        Application:
            <<: *ApplicationDefaults
            Organizations:"
        for ORG in $PEER_ORGS; do
            echo "                - *${ORG}"
        done

        echo "

    EtcdRaftNetwork:
        <<: *ChannelDefaults
        Capabilities:
            <<: *ChannelCapabilities
        Orderer:
            <<: *OrdererDefaults
            OrdererType: etcdraft
            EtcdRaft:
                Consenters:"
        for ORG in $ORDERER_ORGS; do
            for ((i = 0; i < $NUM_ORDERERS; i++)); do
                initOrdererVars $ORG $i
                echo "                - Host: ${ORDERER_HOST}
                  Port: ${ORDERER_PORT}
                  ClientTLSCert: ${TLSDIR}/server.crt
                  ServerTLSCert: ${TLSDIR}/server.crt"
            done
        done
        # echo " 
        #     Addresses:"
        # for ORG in $ORDERER_ORGS; do
        #     for ((i = 0; i < $NUM_ORDERERS; i++)); do
        #         initOrdererVars $ORG $i
        #         echo "                - ${ORDERER_HOST}:${ORDERER_PORT}"
        #     done
        # done
        echo "
                Options:
                    SnapshotIntervalSize: 20 MB"
        echo "
            Organizations:"
        for ORG in $ORDERER_ORGS; do
            echo "            - *${ORG}"
        done
        echo "
            Capabilities:
                <<: *OrdererCapabilities
        Application:
            <<: *ApplicationDefaults
            Organizations:"
        for ORG in $ORDERER_ORGS; do
            echo "            - <<: *${ORG}"
        done
        echo "
        Consortiums:
            SampleConsortium:
                Organizations:"
        for ORG in $PEER_ORGS; do
            echo "                - *${ORG}"
        done

    } >/root/data/configtx.yaml
}

function printOrg() {
    echo "
  - &$ORG
    Name: $ORG
    ID: $ORG_MSP_ID
    MSPDir: $ORG_MSP_DIR
    Policies:
        Readers:
            Type: Signature
            Rule: \"OR('$ORG_MSP_ID.member')\"
        Writers:
            Type: Signature
            Rule: \"OR('$ORG_MSP_ID.member')\"
        Admins:
            Type: Signature
            Rule: \"OR('$ORG_MSP_ID.admin')\"
        Endorsement:
            Type: Signature
            Rule: \"OR('$ORG_MSP_ID.member')\""
}

function printOrdererOrg() {
    initOrgVars $1
    printOrg
    echo "
    OrdererEndpoints:"
    for ((i = 0; i < $NUM_ORDERERS; i++)); do
        initOrdererVars $1 $i
        echo "         - ${ORDERER_HOST}:${ORDERER_PORT}"
    done
}

function printPeerOrg() {
    initPeerVars $1 $2
    printOrg
    echo "
    AnchorPeers:
       - Host: $PEER_HOST
         Port: $PEER_PORT
         "

}

function generateChannelArtifacts() {

    echo "[generateChannelArtifacts] Start generateChannelArtifacts function"

    export FABRIC_CFG_PATH=/root/data/
    for CHANNEL_NAME in $CHANNEL_NAMES; do
        mkdir -p /root/data/$CHANNEL_NAME

        echo "[generateChannelArtifacts] Generating channel configuration transaction"
        /root/data/bin/configtxgen -profile OneOrgChannel -outputCreateChannelTx /root/data/$CHANNEL_NAME/$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
        if [ "$?" -ne 0 ]; then
            echo "[generateChannelArtifacts] Failed to generate channel configuration transaction"
            exit 1
        fi

        for ORG in $PEER_ORGS; do
            initOrgVars $ORG
            echo "[generateChannelArtifacts] Generating anchor peer update transaction for $ORG at $ANCHOR_TX_FILE"
            /root/data/bin/configtxgen -profile OneOrgChannel -outputAnchorPeersUpdate /root/data/$CHANNEL_NAME/${ORG}MSPanchors.tx \
                -channelID $CHANNEL_NAME -asOrg $ORG
            if [ "$?" -ne 0 ]; then
                echo "[generateChannelArtifacts] Failed to generate anchor peer update for $ORG"
                exit 1
            fi
        done
    done

    /root/data/bin/configtxgen -profile EtcdRaftNetwork -outputBlock /root/data/genesis.block -channelID testchainid

    sleep 2
    echo "[generateChannelArtifacts] Finish generateChannelArtifacts function"

}

set -e

SDIR=$(dirname "$0")
source $SDIR/env.sh

main
