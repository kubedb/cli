> New to KubeDB? Please start [here](/docs/tutorial.md).

# Installation Guide

## Install KubeDB CLI
KubeDB provides a CLI to work with database objects. Download pre-built binaries from [k8sdb/cli Github releases](https://github.com/k8sdb/cli/releases) and put the binary to some directory in your `PATH`. To install on Linux 64-bit and MacOS 64-bit you can run the following commands:

```sh
# Linux amd 64-bit
wget -O kubedb https://github.com/k8sdb/cli/releases/download/0.2.0/kubedb-linux-amd64 \
  && chmod +x kubedb \
  && sudo mv kubedb /usr/local/bin/

# Mac 64-bit
wget -O kubedb https://github.com/k8sdb/cli/releases/download/0.2.0/kubedb-darwin-amd64 \
  && chmod +x kubedb \
  && sudo mv kubedb /usr/local/bin/
```

If you prefer to install KubeDB cli from source code, you will need to set up a GO development environment following [these instructions](https://golang.org/doc/code.html). Then, install `kubedb` CLI using `go get` from source code.

```bash
go get github.com/k8sdb/cli/...
```

Please note that this will install KubeDB cli from master branch which might include breaking and/or undocumented changes.

## Install KubeDB Operator
To use KubeDB, you will need to install KubeDB [operator](https://github.com/k8sdb/operator).  `kubedb init` command will deploy operator for kubedb databases. 

```sh
$ kubedb init

Successfully created operator deployment.
Successfully created operator service.
```

## Verify installation
To check if KubeDB operator pods have started, run the following command:
```sh
$ kubectl get pods --all-namespaces -l app=kubedb --watch
```

Once the operator pods are running, you can cancel the above command by typing `Ctrl+C`.

Now, to confirm TPR groups have been registered by the operator, run the following command:
```sh
$ kubectl get thirdpartyresources -l app=kubedb
```

Now, you are ready to [create your first database](/docs/tutorial.md) using KubeDB.

## Upgrade KubeDB
To upgrade KubeDB cli, just replace the old cli with the new version.

`kubedb init` command can be used to upgrade operator. Re-run the `kubedb init` command with `--upgrade flag` to upgrade operator.

```sh
$ kubedb init --version='0.2.0' --upgrade

Successfully upgraded operator deployment.
```
To learn about various options of `init` command, please visit [here](/docs/reference/kubedb_init.md).
