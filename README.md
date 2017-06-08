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

* [kubedb init](docs/user-guide/task/init.md) to deploy unified operator.
* [kubedb create](docs/user-guide/task/create.md) to create a database object.
* [kubedb describe](docs/user-guide/task/describe.md) to describe a supported object.
* [kubedb get](docs/user-guide/task/get.md) to get/list supported object(s).
* [kubedb edit](docs/user-guide/task/edit.md) to edit supported object(s).
* [kubedb delete](docs/user-guide/task/delete.md) to delete supported object(s).


See step-by-step [tutorial](docs/user-guide/tutorial.md)
