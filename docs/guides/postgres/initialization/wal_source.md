---
title: Initialize Postgres from WAL
menu:
  docs_0.9.0:
    identifier: pg-wal-source-initialization
    name: From WAL
    parent: pg-initialization-postgres
    weight: 20
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).
> Don't know how to take continuous backup?  Check [tutorial](/docs/guides/postgres/snapshot/continuous_archiving.md) on Continuous Archiving.

# PostgreSQL Initialization from WAL files

KubeDB supports PostgreSQL database initialization. When you create a new Postgres object, you can provide existing WAL files to restore from by "replaying" the log entries.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace/demo created
```

## Prepare WAL Archive

We need a WAL archive to perform initialization. If you already don't have a WAL archive ready, create one by following the tutorial [here](/docs/guides/postgres/snapshot/continuous_archiving.md).

Let's populate the database so that we can verify that the initialized database has the same data. We will `exec` into the database pod and use `psql` command-line tool to create a table.

At first, find out the primary replica using the following command,

```console
$ kubectl get pods -n demo --selector="kubedb.com/name=wal-postgres","kubedb.com/role=primary"
NAME             READY     STATUS    RESTARTS   AGE
wal-postgres-0   1/1       Running   0          8m
```

Now, let's `exec` into the pod and create a table,

```console
$ kubectl exec -it -n demo wal-postgres-0 sh
# login as "postgres" superuser.
/ # psql -U postgres
psql (9.6.7)
Type "help" for help.

# list available databases
postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
(3 rows)

# connect to "postgres" database
postgres=# \c postgres
You are now connected to database "postgres" as user "postgres".

# create a table
postgres=# CREATE TABLE COMPANY( NAME TEXT NOT NULL, EMPLOYEE INT NOT NULL);
CREATE TABLE

# list tables
postgres=# \d
          List of relations
 Schema |  Name   | Type  |  Owner   
--------+---------+-------+----------
 public | company | table | postgres

# quit from the database
postgres=# \q

# exit from the pod
/ # exit
```

Now, we are ready to proceed for rest of the tutorial.

> Note: YAML files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create Postgres with WAL source

We can initialize a new database from this archived WAL files. We have to specify the archive backend in the `spec.init.postgresWAL` field of Postgres object.

Here, the YAML of Postgres object that we are going to create in this tutorial,

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: replay-postgres
  namespace: demo
spec:
  version: "9.6-v2"
  replicas: 2
  databaseSecret:
    secretName: wal-postgres-auth
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
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
  - `s3` points to s3 storage configuration.
  - `s3.bucket` points to the bucket name where archived WAL data is stored.
  - `s3.prefix` points to the path of archived WAL data.
  - `gcs` points to GCS storage configuration.
  - `gcs.bucket` points to the bucket name where archived WAL data is stored.
  - `gcs.prefix` points to the path of archived WAL data..

User can use either s3 or gcs WAL archiver for init.

**wal-g** receives archived WAL data from a folder called `/kubedb/{namespace}/{postgres-name}/archive/`.

Here, `{namespace}` & `{postgres-name}` indicates Postgres object whose WAL archived data will be replayed.

> Note: Postgres `replay-postgres` must have same superuser credentials as archived Postgres. In our case, it is `wal-postgres`.

Now, let's create the Postgres object that's YAML has shown above,

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/postgres/initialization/replay-postgres.yaml 
postgres.kubedb.com/replay-postgres created
```

This will create a new database and will initialize the database from the archived WAL files.

## Verify Initialization

Let's verify that the new database has been initialized successfully from the WAL archive. It must contain the table we have created for `wal-postgres` database.

We will `exec` into new database pod and use `psql` command-line tool to list tables of `postgres` database.

```console
$ kubectl exec -it -n demo replay-postgres-0 sh
# login as "postgres" superuser
/ # psql -U postgres
psql (9.6.7)
Type "help" for help.

# list available databases
postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
(3 rows)

# connect to "postgres" database
postgres=# \c postgres
You are now connected to database "postgres" as user "postgres".

# list tables
postgres=# \d
          List of relations
 Schema |  Name   | Type  |  Owner   
--------+---------+-------+----------
 public | company | table | postgres
(1 row)

# quit from the database
postgres=# \q

# exit from pod
/ # exit
```

So, we can see that our new database `replay-postgres` has been initialized successfully and contains the data we had inserted into `wal-postgres`.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo pg/replay-postgres -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo pg/replay-postgres

kubectl delete ns demo
```

Also cleanup the resources created for `wal-postgres` following the guide [here](/docs/guides/postgres/snapshot/continuous_archiving.md#cleaning-up).

## Next Steps

- Learn about initializing [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
