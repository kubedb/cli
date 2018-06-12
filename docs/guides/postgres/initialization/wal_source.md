---
title: Initialize Postgres from WAL
menu:
  docs_0.8.0:
    identifier: pg-wal-source-initialization
    name: From WAL
    parent: pg-initialization-postgres
    weight: 20
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

> Don't know how to take continuous backup?  Check [tutorial](/docs/guides/postgres/snapshot/continuous_archiving.md) on Continuous Archiving.

# PostgreSQL Initialization

KubeDB supports PostgreSQL database initialization. When you create a new Postgres object, you can provide existing WAL files to restore from by "replaying" the log entries.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: Yaml files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

## Create PostgreSQL with WAL Source

You can create a new database from archived WAL files using [wal-g ](https://github.com/wal-g/wal-g).

Specify storage backend in the `spec.init.postgresWAL` field of a new Postgres object.

See the example Postgres object below

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: replay-postgres
  namespace: demo
spec:
  version: "9.6"
  replicas: 2
  databaseSecret:
    secretName: wal-postgres-auth
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  archiver:
    storage:
      storageSecretName: s3-secret
      s3:
        bucket: kubedb
  init:
    postgresWAL:
      storageSecretName: s3-secret
      s3:
        bucket: kubedb
        prefix: 'kubedb/demo/wal-postgres/archive'
```

Here,

- `spec.init.postgresWAL` specifies storage information that will be used by `wal-g`
  - `storageSecretName` points to the Secret containing the credentials for cloud storage destination.
  - `s3.bucket` points to the bucket name used to store continuous archiving data.
  - `s3.prefix` points to the path where archived WAL data is stored.

**wal-g** receives archived WAL data from a folder called `/kubedb/{namespace}/{postgres-name}/archive/`.

Here, `{namespace}` & `{postgres-name}` indicates Postgres object whose WAL archived data will be replayed.

> Note: Postgres `replay-postgres` must have same `postgres` superuser password as Postgres `wal-postgres`.

[//]: # (Describe authentication part. This should match with existing one)

Now create this Postgres

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/postgres/initialization/replay-postgres.yaml
postgres "replay-postgres" created
```

This will create a new database with existing _basebackup_ and will restore from archived _wal_ files.

When this database is ready, **wal-g** takes a _basebackup_ and uploads it to cloud storage defined by storage backend in `spec.archiver`.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo pg/replay-postgres -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo pg/replay-postgres

$ kubectl patch -n demo drmn/replay-postgres -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/replay-postgres

$ kubectl delete ns demo
```

## Next Steps

- Learn about initializing [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
