#!/bin/bash

echo ""
echo "Backup process is going on..."
echo ""

Time=$(echo $(date +'%s'))
SNAPSHOT_FILE="$FOLDER_NAME"/"$SNAPSHOT_NAME-$Time"
mongodump --out "$SNAPSHOT_NAME-$Time"
tar  -cvzf "$SNAPSHOT_NAME-$Time".tgz "$SNAPSHOT_NAME-$Time/"*
rm -r "$SNAPSHOT_NAME-$Time"

aws s3 cp "$SNAPSHOT_NAME-$Time".tgz s3://"$BUCKET"/mongo-backup/"$SNAPSHOT_NAME-$Time"

rm "$SNAPSHOT_NAME-$Time".tgz
echo ""
