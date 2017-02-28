#!/bin/bash

exec 1> >(logger -s -p daemon.info -t es)
exec 2> >(logger -s -p daemon.error -t es)

RETVAL=0

backup(){
  # 1 - host
  # 2 - snap-name
  # 3 - index

  path=/var/dump-backup/$2/$3
  mkdir -p $path
  cd $path
  rm -rf $path/*

  elasticdump --quiet --input http://$1:9200/$3 --output $3.mapping.json --type mapping
  retval=$?
  if [ "$retval" -ne 0 ]; then
    echo "Fail to dump mapping for $3"
    exit 1
  fi

  elasticdump --quiet --input http://$1:9200/$3 --output $3.analyzer.json --type analyzer
  retval=$?
  if [ "$retval" -ne 0 ]; then
    echo "Fail to dump analyzer for $3"
    exit 1
  fi

  elasticdump --quiet --input http://$1:9200/$3 --output $3.json --type data
  retval=$?
  if [ "$retval" -ne 0 ]; then
    echo "Fail to dump data for $3"
    exit 1
  fi

  echo "Successfully dump for $3"
}

restore(){
  # 1 - host
  # 2 - snap-name
  # 3 - index
  path=/var/dump-restore/$2/$3
  cd $path

  elasticdump --quiet --input $3.analyzer.json --output http://$1:9200/$3 --type analyzer
  retval=$?
  if [ "$retval" -ne 0 ]; then
    echo "Fail to restore analyzer for $3"
    exit 1
  fi

  elasticdump --quiet --input $3.mapping.json --output http://$1:9200/$3 --type mapping
  retval=$?
  if [ "$retval" -ne 0 ]; then
    echo "Fail to restore mapping for $3"
    exit 1
  fi


  elasticdump --quiet --input $3.json --output http://$1:9200/$3 --type data
  retval=$?
  if [ "$retval" -ne 0 ]; then
    echo "Fail to restore data for $3"
    exit 1
  fi

  echo "Successfully restore for $3"
}

push() {
  # 1 - cloud
  # 2 - bucket
  # 3 - database
  # 4 - snap-name

  path=/var/dump-backup/$4

  if [ "$1" = 'gce' ]; then
    gsutil -m cp -r $path gs://$2/$3/$4
    retval=$?
    if [ "$retval" -ne 0 ]; then
        exit 1
    fi
  fi

  if [ "$1" = 'aws' ]; then
    region=$(aws s3api get-bucket-location --bucket=$2 --output=text)
    if [ $region = "None" ]; then
        aws s3 cp --recursive $path s3://$2/$3/$4
    else
        aws s3 cp --region $region --recursive $path s3://$2/$3/$4
    fi
    retval=$?
    if [ "$retval" -ne 0 ]; then
        exit 1
      fi
  fi
  exit 0
}


pull() {
  # 1 - cloud
  # 2 - bucket
  # 3 - database
  # 4 - snap-name

  path=/var/dump-restore
  mkdir -p $path
  cd $path
  rm -rf *

  if [ "$1" = 'gce' ]; then
      gsutil -m cp -r gs://$2/$3/$4 .
      retval=$?
      if [ "$retval" -ne 0 ]; then
          exit 1
      fi
  fi

  if [ "$1" = 'aws' ]; then
      region=$(aws s3api get-bucket-location --bucket=$2 --output=text)
      if [ $region = "None" ]; then
            aws s3 cp --recursive s3://$2/$3/$4 $path
      else
            aws s3 cp --region $region --recursive s3://$2/$3/$4 $path
      fi
      retval=$?
      if [ "$retval" -ne 0 ]; then
          exit 1
      fi
  fi
}


process=$1
shift
case "$process" in
	backup)
		backup "$@"
		;;
	restore)
		restore "$@"
		;;
	push)
	  push "$@"
	  ;;
	pull)
		pull "$@"
		;;
	*)	(10)
		echo $"Unknown process!"
		RETVAL=1
esac
exit $RETVAL
