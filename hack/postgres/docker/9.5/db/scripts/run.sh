#!/bin/bash

exec 1> >(logger -s -p daemon.info -t pg)
exec 2> >(logger -s -p daemon.error -t pg)

RETVAL=0

mode=$1
shift
case "$mode" in
	basic)
		/scripts/basic/basic.sh "$@"
		;;
	*)	(10)
		echo $"Unknown mode!"
		RETVAL=1
esac
exit $RETVAL
