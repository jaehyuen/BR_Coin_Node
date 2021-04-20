#!/usr/bin/env bash

function printHelp() {

    echo
    echo
    echo "Usage: "
    echo "  network.sh -up <after start ca, setup container and copy crypto-config to other server, start fabric network> "
    echo "  network.sh -up [ca <ca, setup>] [o <orderer>] [p <peer>] [co <couchdb>] [ex <explorer>] [exdb <explorer db>]"
    echo "  network.sh -up [container_name <other containers>]"
    echo "  network.sh -down <all containers stop and network reset>"
    echo "  network.sh -down [ca <ca, setup>] [o <orderer>] [p <peer>] [co <couchdb>] [ex <explorer>] [exdb <explorer db>]"
    echo "  network.sh -down [container_name <other containers>]"
    echo
    echo "  network.sh -install [channel_name] [chaincode_name] [version] <install chaincode to specific version>"
    echo
    echo "  network.sh -addorg [channel_name] [org name] <add organization to the channel>"
    echo
    echo "  network.sh -removeorg [channel_name] [org name] <add organization to the channel>"
    echo
    echo "  network.sh -getconfig [channel_name] <get config in channel>"
    echo "  network.sh -setconfig [channel_name] <set config in channel after you modified existing config>"
    echo
    echo "start example: "
    echo "  network.sh -up "
    echo "  network.sh -up ca <ca, setup container start>"
    echo "  network.sh -up o <orderer container start>"
    echo "  network.sh -up p <peer container start>"
    echo "  network.sh -up cli <cli container start>"
    echo
    echo "shutdown example"
    echo "  network.sh -down "
    echo "  network.sh -down ca <ca, setup container stop>"
    echo "  network.sh -down o <orderer container stop>"
    echo "  network.sh -down p <peer container stop>"
    echo "  network.sh -down cli <cli container stop>"
    echo
    echo "install example"
    echo "  network.sh -install brcoin-channel brcoin-cc 3.2.11 <install to specific version>"
    echo
    echo "addorg example"
    echo "  network.sh -addorg brcoin-channel test"
    echo
    echo "removeorg example"
    echo "  network.sh -removeorg brcoin-channel test"
    echo
    echo "get config example"
    echo "  network.sh -getconfig brcoin-channel"
    echo
    echo "set config example"
    echo "  network.sh -setconfig brcoin-channel"
    echo
    echo "============"
    echo "   NOTICE   "
    echo "============"
    echo "* You must start first ca, setup container to enroll fabric certs"
    echo "* You must copy crypto-config dir to other servers"
    echo
    echo

}

function resetNetwork() {
    if (($EUID != 0)); then
        echo "This option must be run as root"
        exit 0
    fi

    docker stop $(docker ps --filter network=brcoin-network -aq)
    docker rm $(docker ps --filter network=brcoin-network -aq)
    docker rmi $(docker images dev* -aq)

    rm -rf ../data/production
    rm -rf ../data/couchdb

}

function removeOrg() {
    ch=$1
    removeorg=$2

    checkChannel $ch

    ./scripts/network/07_2_remove_orgs.sh $ch $removeorg
}

function addOrg() {
    ch=$1
    addorg=$2

    checkChannel $ch
    if [ -n "$3" ]; then
        anchor_port=$3
    else
        echo "You must input anchor peer port"
        exit 1

    fi

    echo $anchor_port
    ./scripts/network/07_1_add_orgs.sh $ch $addorg $anchor_port
}
function getConfig() {
    ch=$1
    checkChannel $ch

    ./scripts/network/08_1_get_chanel_config.sh $ch
}

function setConfig() {
    ch=$1
    checkChannel $ch

    ./scripts/network/08_2_set_chanel_config.sh $ch
}


function checkChannel() {

    flag=false
    if [ $# -ne 1 ]; then
        echo "You must input channel name"
        exit 1

    fi

    docker exec -it cli peer channel list >./scripts/channel.txt
    sed -i '1,2d' ./scripts/channel.txt
    tr -d '\r' <./scripts/channel.txt >./scripts/channel_list.txt

    rm -rf ./scripts/channel.txt
    channel_list=$(cat ./scripts/channel_list.txt)

    for channel in $channel_list; do

        if [[ "$channel" == "$1" ]]; then
            flag=true
            break
        else
            echo
        fi
    done

    if [ "$flag" == "false" ]; then
        echo "Please check channel name it isn't exist channel : $1"
        exit 1
    fi

}

function installChaincode() {

    ch=$1
    cc=$2
    # ccver_list=()
    # checkChannel $ch
    # # checkChaincode $ch $cc

    # ccver=$(checkChaincode $ch $cc)
    # ccver=${ccver#*Version:}
    # ccver=${ccver%%,*}

    # lists=$(echo $ccver | tr "." "\n")
    # i=0

    # for list in $lists; do
    #     ccver_list[$i]=$(echo $list)
    #     i=$i+1
    # done

    # if [ -n "$3" ]; then
    #     newccver=$3
    # else
    #     echo "You must input chaincode version"
    #     exit 1

    # fi

    ./scripts/network/04_install_chaincode.sh $ch $cc $newccver

    echo "Succcess install in $ch chaincode : $cc, version: $newccver"

}

function checkChaincode() {

    if [ $# -ne 2 ]; then
        echo "You must input chaincode name"
        exit 1

    fi

    docker exec -it cli peer chaincode list --instantiated -C $1 >./scripts/chaincodes.txt
    sed -i '1,1d' ./scripts/chaincodes.txt
    tr -d '\r' <./scripts/chaincodes.txt >./scripts/chaincode.txt

    rm -rf ./scripts/chaincodes.txt
    chaincode=$(cat ./scripts/chaincode.txt)

    if [[ "$chaincode" == *"$2"* ]]; then
        echo
    else
        echo "Please check chaincode name it isn't exist chaincode in channel $1 : $2"
        exit 1
    fi
    echo $chaincode
}

function startDocker() {

    if [ "$2" == "ca" ]; then
        container_name="ca.orgbrcoin.com ca.orgorderer.com setup"
    elif [ "$2" == "o" ]; then
        container_name="orderer0.orgorderer.com orderer1.orgorderer.com"
    elif [ "$2" == "p" ]; then
        container_name="peer0.orgbrcoin.com peer1.orgbrcoin.com"
    elif [ "$2" == "co" ]; then
        container_name="couchdb0.orgbrcoin.com couchdb1.orgbrcoin.com"
    else
        container_name=$2
    fi

    if [ $# -ne 2 ]; then
        echo "start brcoin network"
        ./scripts/network/01_start_docker_container_server.sh
        sleep 20
        ./scripts/network/02_channel_create.sh
        ./scripts/network/03_channel_join_and_anchor.sh
        exit 1

    elif [ $# -ne 3 ]; then
        docker-compose -f ./compose-files/docker-compose.yaml up -d $container_name

        exit 1
    fi

}

function stopDocker() {

    if [ "$2" == "ca" ]; then
        container_name="ca.orgbrcoin.com ca.orgorderer.com setup"
    elif [ "$2" == "o" ]; then
        container_name="orderer0.orgorderer.com orderer1.orgorderer.com"
    elif [ "$2" == "p" ]; then
        container_name="peer0.orgbrcoin.com peer1.orgbrcoin.com"
    elif [ "$2" == "co" ]; then
        container_name="couchdb0.orgbrcoin.com couchdb1.orgbrcoin.com"
    else
        container_name=$2
    fi

    if [ $# -ne 2 ]; then

        while true; do
            read -p "Do you wish to reset network? y/n  " yn
            case $yn in
            [Yy]*)
                ./scripts/network/99_stop_docker_container_server.sh
                break
                ;;
            [Nn]*) exit ;;
            *) echo "Please answer yes or no." ;;
            esac
        done

        exit 1
    elif [ $# -ne 3 ]; then
        docker stop $container_name
        exit 1
    fi
}
export PROD_USER=$(id -g)
. scripts/network/conf
if [ "$1" == "-up" ]; then
    startDocker $1 $2
elif [ "$1" == "-down" ]; then
    stopDocker $1 $2
elif [ "$1" == "-install" ]; then
    installChaincode $2 $3 $4
elif [ "$1" == "-getconfig" ]; then
    getConfig $2
elif [ "$1" == "-setconfig" ]; then
    setConfig $2
elif [ "$1" == "-addorg" ]; then
    addOrg $2 $3 $4
elif [ "$1" == "-removeorg" ]; then
    removeOrg $2 $3
elif [ "$1" == "-reset" ]; then
    resetNetwork
elif [ "$1" == "-h" ]; then
    printHelp
else
    printHelp
    exit 1
fi
