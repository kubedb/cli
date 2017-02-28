#!/bin/bash

push() {
    SNAPSHOT_PATH="/mount/snapshots/$2"
    if [ "$1" = 'gce' ]; then
        gsutil -m cp -r $SNAPSHOT_PATH gs://"$2"/"$3"/"$4"
        retval=$?
        if [ "$retval" -ne 0 ]; then
            exit 1
        fi
        exit 0
    fi

    if [ "$1" = 'aws' ]; then
        aws s3 cp --recursive $SNAPSHOT_PATH s3://"$2"/"$3"/"$4"
        retval=$?
        if [ "$retval" -ne 0 ]; then
            exit 1
        fi
        exit 0
    fi
}

pull() {
    SNAPSHOT_PATH="/mount/snapshots/$2"
    if [ "$1" = 'gce' ]; then
        gsutil -m cp -r gs://"$2"/"$3"/"$4"/* $SNAPSHOT_PATH
        retval=$?
        if [ "$retval" -ne 0 ]; then
            exit 1
        fi
    fi
    if [ "$1" = 'aws' ]; then
	      aws s3 cp --recursive s3://"$2"/"$3"/"$4" $SNAPSHOT_PATH
	      retval=$?
        if [ "$retval" -ne 0 ]; then
            exit 1
        fi
    fi
    chown -R es:es $SNAPSHOT_PATH
    exit 0
}

if [ "$1" == "push" ]; then
    push $2 $3 $4 $5
fi

if [ "$1" == "pull" ]; then
    pull $2 $3 $4 $5
fi
