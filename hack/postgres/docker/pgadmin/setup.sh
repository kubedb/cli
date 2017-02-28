#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/github.com/appscode/container-datastore

source "$REPO_ROOT/hack/libbuild/common/lib.sh"
source "$REPO_ROOT/hack/libbuild/common/public_image.sh"

IMG=pgadmin
TAG=4-1.1

binary_repo $@
