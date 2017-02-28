#!/bin/bash

echo ""
echo "Backup process is going on..."
echo ""

Time=$(echo $(date +'%s'))
SNAPSHOT_FILE="$FOLDER_NAME"/"$SNAPSHOT_NAME-$Time"

mysqldump --all-databases > $SNAPSHOT_FILE

aws s3 cp $SNAPSHOT_FILE s3://"$BUCKET"/mysql-backup/"$SNAPSHOT_NAME-$Time"

echo ""