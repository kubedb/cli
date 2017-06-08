# kubedb

## Installing

Lets install `kubedb` CLI using `go get` from source code.

Following command will install the latest version of the library from master.

```bash
go get github.com/k8sdb/cli/...
```

## Usage

`kubedb` CLI is used to manipulate kubedb ThirdPartyResource objects.

We will go through each of the commands and will see how these commands interact with TPR objects for kubedb databases.

* [kubedb init](user-guide/task/init.md) to deploy unified operator.
* [kubedb create](user-guide/task/create.md) to create a database object.
* [kubedb describe](user-guide/task/describe.md) to describe a supported object.
* [kubedb get](user-guide/task/get.md) to get/list supported object(s).
* [kubedb edit](user-guide/task/edit.md) to edit supported object(s).
* [kubedb delete](user-guide/task/delete.md) to delete supported object(s).


See step-by-step [tutorial](user-guide/tutorial.md)
