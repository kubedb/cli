#!/usr/bin/env bash

pushd $GOPATH/src/github.com/k8sdb/cli/hack/gendocs
go run main.go

cd $GOPATH/src/github.com/k8sdb/cli/docs/user-guide/reference
sed -i 's/######\ Auto\ generated\ by.*//g' *
popd
