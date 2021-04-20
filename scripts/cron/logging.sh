#file formet with date
FILE_FORMET=$(date +"%Y-%m-%d")

#log files
LOGFILES=$(ls /logs/*/*.log)

PROD_USER=$1

for FILE in $LOGFILES; do

    if [ -z "$FILE" ]; then
        continue
    elif [[ "$FILE" == *"explorer"* ]]; then
        continue

    fi

    cp -r $FILE $FILE-$FILE_FORMET
    cat /dev/null >$FILE
    
done

chown -R $PROD_USER:$PROD_USER /logs