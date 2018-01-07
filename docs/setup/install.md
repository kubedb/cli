---
title: Install
menu:
  docs_0.8.0-beta.0:
    identifier: install-kubedb
    name: Install
    parent: setup
    weight: 10
menu_name: docs_0.8.0-beta.0
section_menu_id: setup
---

> New to KubeDB? Please start [here](/docs/guides/README.md).

# Installation Guide

## Install KubeDB CLI
KubeDB provides a CLI to work with database objects. Download pre-built binaries from [kubedb/cli Github releases](https://github.com/kubedb/cli/releases) and put the binary to some directory in your `PATH`. To install on Linux 64-bit and MacOS 64-bit you can run the following commands:

```console
# Linux amd 64-bit
wget -O kubedb https://github.com/kubedb/cli/releases/download/0.8.0-beta.0/kubedb-linux-amd64 \
  && chmod +x kubedb \
  && sudo mv kubedb /usr/local/bin/

# Mac 64-bit
wget -O kubedb https://github.com/kubedb/cli/releases/download/0.8.0-beta.0/kubedb-darwin-amd64 \
  && chmod +x kubedb \
  && sudo mv kubedb /usr/local/bin/
```

If you prefer to install KubeDB cli from source code, you will need to set up a GO development environment following [these instructions](https://golang.org/doc/code.html). Then, install `kubedb` CLI using `go get` from source code.

```bash
go get github.com/kubedb/cli/...
```

Please note that this will install KubeDB cli from master branch which might include breaking and/or undocumented changes.

## Install KubeDB Operator
To use `kubedb`, you will need to install KubeDB [operator](https://github.com/kubedb/operator).  `kubedb init` command will deploy operator for KubeDB databases.

```console
$ kubedb init

Successfully created operator deployment.
Successfully created operator service.
```

`KubeDB` operators supports RBAC v1beta1 authorization api for Kubernetes clusters. To use `KubeDB` in an [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/) enabled cluster, run the `kubedb init` command with `--rbac` flag.
```console
$ kubedb init --rbac
Successfully created cluster role.
Successfully created service account.
Successfully created cluster role bindings.
Successfully created operator deployment.
Successfully created operator service.
```

This will create a `kubedb-operator` ClusterRole, ClusterRoleBinding and ServiceAccount to run `KubeDB` operator pods. With RBAC enabled, a separate set of ClusterRole, ClusterRoleBinding and ServiceAccount will be created for each KubeDB database. To learn more, please visit [here](/docs/guides/rbac.md).

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

## Upgrade KubeDB
To upgrade KubeDB cli, just replace the old cli with the new version.

`kubedb init` command can be used to upgrade operator. Re-run the `kubedb init` command with `--upgrade` flag to upgrade operator.

```console
$ kubedb init --version='0.8.0-beta.0' --upgrade

Successfully upgraded operator deployment.
```
To learn about various options of `init` command, please visit [here](/docs/reference/kubedb_init.md).
