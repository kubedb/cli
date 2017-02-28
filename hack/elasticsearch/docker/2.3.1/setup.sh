#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/github.com/appscode/kubedb

source "$REPO_ROOT/hack/libbuild/common/lib.sh"
source "$REPO_ROOT/hack/libbuild/common/public_image.sh"

IMG=elasticsearch
TAG=2.3.1-v2

build() {
  cp -r ../lib/* .
	gsutil cp gs://appscode-dev/binaries/elasticsearch_discovery/0.1/elasticsearch_discovery-linux-amd64 elasticsearch_discovery
	chmod 755 elasticsearch_discovery
	local cmd="docker build -t appscode/$IMG:$TAG ."
	echo $cmd; $cmd
	rm elasticsearch_discovery
}

binary_repo $@
