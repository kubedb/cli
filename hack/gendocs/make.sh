#!/usr/bin/env bash

pushd $GOPATH/src/github.com/k8sdb/cli/hack/gendocs
go run main.go
popd
