#!/bin/bash

echo "Please enter backup time: "
echo "Example: 1445513365"
read Time

echo
echo "Restoring backup-$Time from S3..."
echo

aws s3 cp  s3://"$BUCKET"/mysql-backup/"$SNAPSHOT_NAME-$Time" /tmp/backup-temp/

mysql < /tmp/backup-temp/"$SNAPSHOT_NAME-$Time"

echo
echo "Restore complete..."
echo
