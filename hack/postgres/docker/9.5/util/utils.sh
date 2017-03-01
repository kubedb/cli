#!/bin/bash

dump() {
    cd /var/dump-backup
    # Wait for postgres to start
	# ref: http://unix.stackexchange.com/a/5279
	while ! nc -q 1 $1 5432 </dev/null; do echo "Waiting... Master pod is not ready yet"; sleep 5; done
	PGPASSWORD=$3 pg_dumpall -U $2 -h $1 > dumpfile.sql
	retval=$?
	if [ "$retval" -ne 0 ]; then
	    exit 1
	fi
    exit 0
}

restore() {
    cd /var/dump-restore
    # Wait for postgres to start
    # ref: http://unix.stackexchange.com/a/5279
    while ! nc -q 1 $1 5432 </dev/null; do echo "Waiting... Master pod is not ready yet"; sleep 5; done
    PGPASSWORD=$3 psql -U $2 -h $1  -f dumpfile.sql postgres
    retval=$?
    if [ "$retval" -ne 0 ]; then
        exit 1
    fi
    exit 0
}

pull() {
    cd /var/dump-restore

    if [ "$1" = 'gce' ]; then
        gsutil cp  gs://$2/$3/$4.sql dumpfile.sql
        retval=$?
        if [ "$retval" -ne 0 ]; then
            exit 1
        fi
        exit 0
    fi

    if [ "$1" = 'aws' ]; then
        region=$(aws s3api get-bucket-location --bucket=$2 --output=text)
        if [ $region = "None" ]; then
            aws s3 cp s3://$2/$3/$4.sql dumpfile.sql
        else
            aws s3 cp --region $region s3://$2/$3/$4.sql dumpfile.sql
        fi
        retval=$?
        if [ "$retval" -ne 0 ]; then
            exit 1
        fi
        exit 0
    fi
}

push() {
    cd /var/dump-backup

    if [ "$1" = 'gce' ]; then
        gsutil cp dumpfile.sql gs://$2/$3/$4.sql
        retval=$?
        if [ "$retval" -ne 0 ]; then
            exit 1
        fi
        exit 0
    fi

    if [ "$1" = 'aws' ]; then
        region=$(aws s3api get-bucket-location --bucket=$2 --output=text)
        if [ $region = "None" ]; then
            aws s3 cp dumpfile.sql s3://$2/$3/$4.sql
        else
            aws s3 cp --region $region dumpfile.sql s3://$2/$3/$4.sql
        fi
        retval=$?
        if [ "$retval" -ne 0 ]; then
            exit 1
        fi
        exit 0
    fi
}


if [ "$1" == "dump" ]; then
    dump $2 $3 $4
fi

if [ "$1" == "restore" ]; then
    restore $2 $3 $4
fi

if [ "$1" == "push" ]; then
    push $2 $3 $4 $5
fi

if [ "$1" == "pull" ]; then
    pull $2 $3 $4 $5
fi
