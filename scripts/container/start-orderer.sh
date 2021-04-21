touch /log/$HOST_NAME.log
chmod 666 /log/$HOST_NAME.log
orderer >> /log/$HOST_NAME.log 2>&1
