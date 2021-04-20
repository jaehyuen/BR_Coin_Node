touch /log/$CORE_PEER_ID.log
chmod 666 /log/$CORE_PEER_ID.log

peer node start >> /log/$CORE_PEER_ID.log 2>&1
