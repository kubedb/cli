#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

LIB_ROOT=$(dirname "${BASH_SOURCE}")/../../../..
source "$LIB_ROOT/hack/libbuild/common/lib.sh"
source "$LIB_ROOT/hack/libbuild/common/public_image.sh"

docker_names=( \
	"db" \
	"util" \
)

IMG=postgres
TAG=9.5-v3

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

docker_release() {
	for name in "${docker_names[@]}"
	do
        docker push appscode/$IMG:$TAG-$name
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

binary_repo $@
