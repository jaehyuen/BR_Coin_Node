PROD_USER=$(id -g)
docker exec cli /bin/bash /scripts/logging.sh $PROD_USER

cd /svc/nhblock/logs
find ./  -name "*log*" -mtime +5 -exec gzip {} \;
find ./  -name "access*" -mtime +5 -exec gzip {} \;
find ./  -name "error*" -mtime +5 -exec gzip {} \;
find ./  -name "*.gz" -mtime +31 -exec rm {} \;