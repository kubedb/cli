#!/usr/bin/env bash

pushd $GOPATH/src/kubedb.dev/cli/hack/gendocs
go run main.go

cd $GOPATH/src/kubedb.dev/cli/docs/reference
sed -i 's/######\ Auto\ generated\ by.*//g' *
popd
