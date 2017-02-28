#!/bin/bash

backup() {
    SNAPSHOT_PATH="/var/influxdb-backup/$3"
    influxd backup -host "$1" -database "$2" $SNAPSHOT_PATH
    retval=$?
    if [ "$retval" -ne 0 ]; then
    exit 1
    fi
    echo "$2"  >> "$SNAPSHOT_PATH/db-list.txt"
    exit 0
}

push() {
    SNAPSHOT_PATH="/var/influxdb-backup/$4"
    if [ "$1" = 'gce' ]; then
        chmod +x /root/.boto
        gsutil -m cp -r $SNAPSHOT_PATH gs://"$2"/"$3"/"$4"
        retval=$?
        if [ "$retval" -ne 0 ]; then
        exit 1
        fi
        echo "$2"  >> "$SNAPSHOT_PATH/db-list.txt"
        exit 0
    fi

    if [ "$1" = 'aws' ]; then
        chmod +x /root/.aws/credentials
        region=$(aws s3api get-bucket-location --bucket=$2 --output=text)
        if [ $region = "None" ]; then
            aws s3 cp --recursive $SNAPSHOT_PATH s3://"$2"/"$3"/"$4"
        else
            aws s3 cp --region $region --recursive $SNAPSHOT_PATH s3://"$2"/"$3"/"$4"
        fi
        retval=$?
        if [ "$retval" -ne 0 ]; then
        exit 1
        fi
        echo "$2"  >> "$SNAPSHOT_PATH/db-list.txt"
        exit 0
    fi
}


if [ "$1" == "backup" ]; then
    backup $2 $3 $4
fi

if [ "$1" == "push" ]; then
    push $2 $3 $4 $5
fi
