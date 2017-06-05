# kubedb

## Installing

Lets install `kubedb` CLI using `go get` from source code.

Following command will install the latest version of the library from master.

```bash
go get github.com/k8sdb/cli/...
```

## Usage

`kubedb` CLI is used to manipulate kubedb ThirdPartyResource objects.

### kubedb help

```bash
$ kubedb --help

kubedb CLI controls kubedb ThirdPartyResource objects.

Find more information at https://github.com/k8sdb/cli.

Basic Commands (Beginner):
  create      Create a resource by filename or stdin
  init        Create or upgrade unified operator

Basic Commands (Intermediate):
  get         Display one or many resources
  edit        Edit a resource on the server
  delete      Delete resources by filenames, stdin, resources and names, or by resources and label selector

Troubleshooting and Debugging Commands:
  describe    Show details of a specific resource or group of resources
  version     Prints binary version number.

Other Commands:
  help        Help about any command

Usage:
  kubedb [flags] [options]

Use "kubedb <command> --help" for more information about a given command.
```

---

We will go through each of these commands and will see how these commands interact with ThirdPartyResources for kubedb databases.

First of all, we need an unified [operator](https://github.com/k8sdb/operator) to handle ThirdPartyResources.

We can deploy this operator using `kubedb` CLI.

### kubedb init

The `kubedb init` command will start an unified operator for kubedb databases. This command can also be used to upgrade version of operator.

Following command will create a deployment with image `kubedb/operator:0.1.0` and a service in `default` namespace
```bash
$ kubedb init --namespace='default' --version='0.1.0'

Successfully created operator deployment.
Successfully created operator service.
```

Any existing operator can also be upgraded using this command.

```bash
$ kubedb init --version='0.2.0' --upgrade

Successfully upgraded operator deployment.
```

---

Now we can create a database supported by **kubedb** using this CLI.

Lets create a Postgres database.

### kubedb create

Following command will create a Postgres TPR as specified in `postgres.yaml`.
This will create an object in `default` namespace by default unless namespace is specified in `postgres.yaml`.

```bash
$ kubedb create -f postgres.yaml

postgres "postgres-demo" created
```
---

We can provide namespace as a flag `--namespace`.

kubedb CLI also supports _stdin_ to create objects.

```bash
cat postgres.yaml | kubedb create -f -
```

We can also provide folder path with `--recursive` flag.

Now lets describe `postgres-demo`

### kubedb describe

Following command will describe Postgres database object `postgres-demo` with relevant information.

```bash
$ kubedb describe pg postgres-demo

Name:		postgres-demo
Namespace:	default
StartTimestamp:	Mon, 05 Jun 2017 10:10:06 +0600
Status:		Running
No volumes.

StatefulSet:
  Name:			postgres-demo-pg
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Mon, 05 Jun 2017 10:10:14 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		postgres-demo
  Type:		ClusterIP
  IP:		10.0.0.36
  Port:		port	5432/TCP

Database Secret:
  Name:	postgres-demo-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	35 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                Type       Reason               Message
  ---------   --------   -----     ----                --------   ------               -------
  21m         21m        1         Postgres operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  21m         21m        1         Postgres operator   Normal     SuccessfulCreate     Successfully created Postgres
  29m         29m        1         Postgres operator   Normal     SuccessfulValidate   Successfully validate Postgres
  29m         29m        1         Postgres operator   Normal     Creating             Creating Kubernetes objects 37s         37s        1         Postgres operator   Normal     Creating             Creating Kubernetes object
```
---

The `kubedb describe` command provides basic information of TPR object, StatefulSet, Service, Secret and list of Snapshots.

If we have multiple objects, we can list them all.

### kubedb get

```bash
$ kubedb get postgres

NAME            STATUS    AGE
postgres-demo   Running   5h
postgres-dev    Running   4h
postgres-prod   Running   30m
postgres-qa     Running   2h
```

To get YAML of an object, we can provide `--output=yaml` flag

```bash
$ kubedb get postgres postgres-demo --output=yaml
```

This `--output` flag also takes `json|wide|name`

We can view all Postgres with their status with `--output=wide`.

```bash
$ kubedb get postgres -o wide
```

> We will see some other options later
