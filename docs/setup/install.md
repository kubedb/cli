---
title: Install
menu:
  docs_0.8.0-beta.2:
    identifier: install-kubedb
    name: Install
    parent: setup
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: setup
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Installation Guide

## Install KubeDB CLI

KubeDB provides a CLI to work with database objects. Download pre-built binaries from [kubedb/cli Github releases](https://github.com/kubedb/cli/releases) and put the binary to some directory in your `PATH`. To install on Linux 64-bit and MacOS 64-bit you can run the following commands:

```console
# Linux amd 64-bit
wget -O kubedb https://github.com/kubedb/cli/releases/download/0.8.0-beta.2/kubedb-linux-amd64 \
  && chmod +x kubedb \
  && sudo mv kubedb /usr/local/bin/

# Mac 64-bit
wget -O kubedb https://github.com/kubedb/cli/releases/download/0.8.0-beta.2/kubedb-darwin-amd64 \
  && chmod +x kubedb \
  && sudo mv kubedb /usr/local/bin/
```

If you prefer to install KubeDB cli from source code, you will need to set up a GO development environment following [these instructions](https://golang.org/doc/code.html). Then, install `kubedb` CLI using `go get` from source code.

```bash
go get github.com/kubedb/cli/...
```

Please note that this will install KubeDB cli from master branch which might include breaking and/or undocumented changes.

## Install KubeDB Operator

To use `kubedb`, you will need to install KubeDB [operator](https://github.com/kubedb/operator).

### Using YAML

KubeDB can be installed via installer script included in the [/hack/deploy](https://github.com/kubedb/cli/tree/0.8.0-beta.2/hack/deploy) folder.

```console
$ curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/hack/deploy/kubedb.sh | bash -s -- -h
kubedb.sh - install kubedb operator

kubedb.sh [options]

options:
-h, --help                         show brief help
-n, --namespace=NAMESPACE          specify namespace (default: kube-system)
    --rbac                         create RBAC roles and bindings
    --docker-registry              docker registry used to pull kubedb images (default: appscode)
    --image-pull-secret            name of secret used to pull kubedb operator images
    --run-on-master                run kubedb operator on master
    --enable-admission-webhook     configure admission webhook for kubedb CRDs
    --uninstall                    uninstall kubedb

# install without RBAC roles
$ curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/hack/deploy/kubedb.sh \
    | bash

# Install with RBAC roles
$ curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/hack/deploy/kubedb.sh \
    | bash -s -- --rbac
```

If you would like to run KubeDB operator pod in `master` instances, pass the `--run-on-master` flag:

```console
$ curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/hack/deploy/kubedb.sh \
    | bash -s -- --run-on-master [--rbac]
```

KubeDB operator will be installed in a `kube-system` namespace by default. If you would like to run KubeDB operator pod in `kubedb` namespace, pass the `--namespace=kubedb` flag:

```console
$ kubectl create namespace kubedb
$ curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/hack/deploy/kubedb.sh \
    | bash -s -- --namespace=kubedb [--run-on-master] [--rbac]
```

If you are using a private Docker registry, you need to pull required images from KubeDB's [Docker Hub account](https://hub.docker.com/r/kubedb/).

To pass the address of your private registry and optionally a image pull secret use flags `--docker-registry` and `--image-pull-secret` respectively.

```console
$ kubectl create namespace kubedb
$ curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/hack/deploy/kubedb.sh \
    | bash -s -- --docker-registry=MY_REGISTRY [--image-pull-secret=SECRET_NAME] [--rbac]
```

KubeDB implements a [validating admission webhook](https://kubernetes.io/docs/admin/admission-controllers/#validatingadmissionwebhook-alpha-in-18-beta-in-19) to validate KubeDB CRDs. This is enabled by default for Kubernetes 1.9.0 or later releases. To disable this feature, pass the `--enable-admission-webhook=false` flag.

```console
$ curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/hack/deploy/kubedb.sh \
    | bash -s -- --enable-admission-webhook [--rbac]
```

### Using Helm

KubeDB can be installed via [Helm](https://helm.sh/) using the [chart](https://github.com/kubedb/cli/tree/master/chart/stable/kubedb) included in this repository. To install the chart with the release name `my-release`:

```console
# Mac OSX amd64:
curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-darwin-amd64 \
  && chmod +x onessl \
  && sudo mv onessl /usr/local/bin/

# Linux amd64:
curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-linux-amd64 \
  && chmod +x onessl \
  && sudo mv onessl /usr/local/bin/

# Linux arm64:
curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-linux-arm64 \
  && chmod +x onessl \
  && sudo mv onessl /usr/local/bin/

# Kubernetes 1.8.x
$ helm repo update
$ helm install stable/kubedb --name my-release

# Kubernetes 1.9.0 or later
$ helm repo update
$ helm install stable/kubedb --name my-release \
  --set apiserver.ca="$(onessl get kube-ca)"
```

To see the detailed configuration options, visit [here](https://github.com/kubedb/cli/tree/master/chart/stable/kubedb).

## Verify installation

To check if KubeDB operator pods have started, run the following command:

```console
$ kubectl get pods --all-namespaces -l app=kubedb --watch
```

Once the operator pods are running, you can cancel the above command by typing `Ctrl+C`.

Now, to confirm CRD groups have been registered by the operator, run the following command:

```console
$ kubectl get crd -l app=kubedb
```

Now, you are ready to [create your first database](/docs/guides/README.md) using KubeDB.

## Configuring RBAC

KubeDB installer will create 3 user facing cluster roles:

| ClusterRole       | Aggregates To | Desription |
| ----------------- | --------------| ---------- |
| kubedb:core:admin | admin         | Allows edit access to all `KubeDB` CRDs, intended to be granted within a namespace using a RoleBinding. This grants ability to wipeout dormant database and delete their record. |
| kubedb:core:edit  | edit          | Allows edit access to all `KubeDB` CRDs except `DormantDatabase` CRD, intended to be granted within a namespace using a RoleBinding. |
| kubedb:core:view  | view          | Allows read-only access to `KubeDB` CRDs, intended to be granted within a namespace using a RoleBinding. |

These user facing roles supports [ClusterRole Aggregation](https://kubernetes.io/docs/admin/authorization/rbac/#aggregated-clusterroles) feature in Kubernetes 1.9 or later clusters.

## Upgrade KubeDB

To upgrade KubeDB cli, just replace the old cli with the new version. To upgrade KubeDB operator, please follow the instruction for the corresponding release.