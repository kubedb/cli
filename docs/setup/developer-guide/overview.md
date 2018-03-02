---
title: Developer Guide
menu:
  docs_0.8.0-beta.2:
    identifier: developer-guide-readme
    name: Overview
    parent: developer-guide
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: setup
aliases:
  - /docs/0.8.0-beta.2/setup/developer-guide/
---

# Development Guide

This document is intended to be the canonical source of truth for things like supported toolchain versions for building KubeDB.
If you find a requirement that this doc does not capture, please submit an issue on github.

This document is intended to be relative to the branch in which it is found. It is guaranteed that requirements will change over time
for the development branch, but release branches of KubeDB should not change.

## Build KubeDB

Some of the KubeDB development helper scripts rely on a fairly up-to-date GNU tools environment, so most recent Linux distros should
work just fine out-of-the-box.

### Setup GO

KubeDB is written in Google's GO programming language. Currently, KubeDB is developed and tested on **go 1.8.3**. If you haven't set up a GO
development environment, please follow [these instructions](https://golang.org/doc/code.html) to install GO.

### Code Organization

KubeDB codebase is across various repositories under github.com/kubedb organization. There are 5 categories of git repositories:

| Repository                             | Description                                                                                              |
|----------------------------------------|----------------------------------------------------------------------------------------------------------|
| https://github.com/kubedb/apimachinery | Contains api types, clientset and KubeDB framework interfaces.                                           |
| https://github.com/kubedb/xdb          | This repository contains operator for `db`, eg, https://github.com/kubedb/postgres                       |
| https://github.com/kubedb/db_exporter  | This repository contains Prometheus exporter for `db`, eg, https://github.com/kubedb/postgres_exporter . |
| https://github.com/kubedb/operator     | This repository contains the combined operator for all databases supported by KubeDB.                    |
| https://github.com/kubedb/cli          | This repository contains CLI for KubeDB.                                                                 |

For each of these repositories, you can get source code and build code using the following steps:

### Download Source

```console
$ go get github.com/kubedb/operator
$ cd $(go env GOPATH)/src/github.com/kubedb/operator
```

### Install Dev tools

To install various dev tools for KubeDB, run the following command:

```console
$ ./hack/builddeps.sh
```

### Build Binary

```console
$ ./hack/make.py
```

### Dependency management

For KubeDB original repositories, we use [Glide](https://github.com/Masterminds/glide) to manage dependencies. Dependencies are already checked in the `vendor` folder. If you want to update/add dependencies, run:
```console
$ glide slow
```

### Build Docker images

For unified operator or db specific operators, we support building Docker images. To build and push your custom Docker image, follow the steps below. To release a new version of KubeDB, please follow the [release guide](/docs/setup/developer-guide/release.md).

```console
# Build Docker image
$ ./hack/docker/operator/setup.sh; ./hack/docker/operator/setup.sh push

# Add docker tag for your repository
$ docker tag kubedb/operator:<tag> <image>:<tag>

# Push Image
$ docker push <image>:<tag>

# Example:
$ docker tag kubedb/operator:default aerokite/operator:default
$ docker push aerokite/operator:default
```

### Generate CLI Reference Docs

```console
$ cd ~/go/src/github.com/kubedb/cli
$ ./hack/gendocs/make.sh
```
