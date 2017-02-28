#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -x

RETVAL=0
IMG=postgres
TAG=9.4

docker_names=( \
	"base" \
	"base-arc" \
	"dump" \
	"st" \
	"ha" \
	"arc" \
)

build() {
	for name in "${docker_names[@]}"
	do
		cd $name
		docker build -t appscode/$IMG:$TAG-$name .
		cd ..
	done
}

docker_push() {
	for name in "${docker_names[@]}"
	do
		docker_up $IMG:$TAG-$name
	done
}

docker_check() {
	for i in "${docker_names[@]}"
	do
		echo "Chcking $IMG ..."
		name=$i-$(date +%s | sha256sum | base64 | head -c 8 ; echo)
		docker run -d -P -it --name=$name appscode/$IMG:$TAG-$i
		docker exec -it $name ps aux
		sleep 5
		docker exec -it $name ps aux
		docker stop $name && docker rm $name
	done
}

if [ $# -eq 0 ]; then
	build
	exit $RETVAL
fi

case "$1" in
	build)
		build
		;;
	push)
		docker_push
		;;
	check)
		docker_check
		;;
	*)	(10)
		echo $"Usage: $0 {build|push|check}"
		RETVAL=1
esac
exit $RETVAL
