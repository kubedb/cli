#!/bin/bash

echo "Please enter backup time: "
echo "Example: 1445513365"
read Time

echo
echo "Restoring backup-$Time from S3..."
echo

aws s3 cp  s3://"$BUCKET"/mongo-backup/"$SNAPSHOT_NAME-$Time" .

mkdir "$FOLDER_NAME/$SNAPSHOT_NAME-$Time"
tar -xvzf "$SNAPSHOT_NAME-$Time" -C "$FOLDER_NAME/"
rm "$SNAPSHOT_NAME-$Time"

mongorestore "$FOLDER_NAME/$SNAPSHOT_NAME-$Time"

echo
echo "Restore complete..."
echo
