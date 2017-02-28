#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/github.com/appscode/container-datastore

source "$REPO_ROOT/hack/libbuild/common/lib.sh"
source "$REPO_ROOT/hack/libbuild/common/public_image.sh"

IMG=mariadb
TAG=10.1
EXTRA_DOCKER_OPTS="-v $PWD/example.admin:/srv/mysql/secrets/.admin"

binary_repo $@
