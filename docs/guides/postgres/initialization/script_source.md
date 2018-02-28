---
title: Initialize Postgres using Script Source
menu:
  docs_0.8.0-beta.2:
    identifier: pg-script-source-initialization
    name: Using Script
    parent: pg-initialization-postgres
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---

> New to KubeDB Postgres?  Quick start [here](/docs/guides/postgres/quickstart/quickstart.md).

# Initialize PostgreSQL with Script

KubeDB supports PostgreSQL database initialization.

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

This tutorial will show you how to use KubeDB to initialize a PostgreSQL database with supported script.

## Create PostgreSQL with script source

PostgreSQL database can be initialized by scripts provided to it.

Following YAML describes a Postgres object that holds VolumeSource of initialization scripts.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: script-postgres
  namespace: demo
spec:
  version: 9.6
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    scriptSource:
      gitRepo:
        repository: "https://github.com/kubedb/postgres-init-scripts.git"
        directory: "."
```

Here,

 -  `init.scriptSource` specifies scripts used to initialize the database when it is being created.

VolumeSource provided in `init.scriptSource` will be mounted in Pod and will be executed while creating PostgreSQL.

In this tutorial, `data.sql` script from the git repository `https://github.com/kubedb/postgres-init-scripts.git` is used to create a TABLE `dashboard` in `data` Schema.

> Note: PostgreSQL supports initialization with `.sh`, `.sql` and `.sql.gz` files.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/postgres/initialization/script-postgres.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/postgres/initialization/script-postgres.yaml"
postgres "script-postgres" created
```

```console
$ kubedb describe pg -n demo script-postgres -S=false -W=false
Name:           script-postgres
Namespace:      demo
StartTimestamp: Thu, 08 Feb 2018 15:55:11 +0600
Status:         Running
Init:
  scriptSource:
    Type:       GitRepo (a volume that is pulled from git when the pod is created)
    Repository: https://github.com/kubedb/postgres-init-scripts.git
    Directory:  .
Volume:
  StorageClass: standard
  Capacity:     50Mi
  Access Modes: RWO
StatefulSet:    script-postgres
Service:        script-postgres, script-postgres-replicas
Secrets:        script-postgres-auth

Topology:
  Type      Pod                 StartTime                       Phase
  ----      ---                 ---------                       -----
  primary   script-postgres-0   2018-02-08 15:55:29 +0600 +06   Running

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                Type       Reason       Message
  ---------   --------   -----     ----                --------   ------       -------
  4m          4m         1         Postgres operator   Normal     Successful   Successfully patched StatefulSet
  4m          4m         1         Postgres operator   Normal     Successful   Successfully patched Postgres
  4m          4m         1         Postgres operator   Normal     Successful   Successfully created StatefulSet
  4m          4m         1         Postgres operator   Normal     Successful   Successfully created Postgres
  5m          5m         1         Postgres operator   Normal     Successful   Successfully created Service
  5m          5m         1         Postgres operator   Normal     Successful   Successfully created Service
```

Now lets connect to our Postgres `script-postgres`  using pgAdmin we have installed in [quickstart](/docs/guides/postgres/quickstart.md#before-you-begin) tutorial.

Connection information:

- address: use Service `script-postgres.demo`
- port: `5432`
- database: `postgres`
- username: `postgres`

Run following command to get `postgres` superuser password

    $ kubectl get secrets -n demo script-postgres-auth -o jsonpath='{.data.\POSTGRES_PASSWORD}' | base64 -d

In PostgreSQL, run following query to check `pg_catalog.pg_tables` to confirm initialization.

```console
select * from pg_catalog.pg_tables where schemaname = 'data';
```

 schemaname | tablename | tableowner | hasindexes | hasrules | hastriggers | rowsecurity
------------|-----------|------------|------------|----------|-------------|-------------
 data       | dashboard | postgres   | t          | f        | f           | f

We can see TABLE `dashboard` in `data` Schema which is created for initialization.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete pg,drmn,snap -n demo --all --force
$ kubectl delete ns demo
```

## Next Steps

- Learn about [taking instant backup](/docs/guides/postgres/snapshot/instant_backup.md) of PostgreSQL database using KubeDB Snapshot.
- Learn about initializing [PostgreSQL from KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
