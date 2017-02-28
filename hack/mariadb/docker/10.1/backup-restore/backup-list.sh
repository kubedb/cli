#!/bin/bash

aws s3 ls s3://"$BUCKET"/mysql-backup/ >> /tmp/list

echo ""
echo "Backup List!!!"
echo ""

echo "backup-xxxxxxxxxx :         TIME"
echo "-----------------   -------------------"
while read line; do
     IFS=' ' read -ra val <<< "$line"
     echo "${val[3]} : ${val[0]} ${val[1]}"
done < /tmp/list

rm /tmp/list
