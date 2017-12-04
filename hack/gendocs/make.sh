#!/usr/bin/env bash

pushd $GOPATH/src/github.com/kubedb/cli/hack/gendocs
go run main.go

cd $GOPATH/src/github.com/kubedb/cli/docs/reference
sed -i 's/######\ Auto\ generated\ by.*//g' *
popd
