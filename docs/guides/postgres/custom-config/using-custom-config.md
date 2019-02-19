---
title: Run PostgreSQL with Custom Configuration
menu:
  docs_0.9.0:
    identifier: pg-custom-config-quickstart
    name: Quickstart
    parent: pg-custom-config
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Custom Configuration File

KubeDB supports providing custom configuration for PostgreSQL. This tutorial will show you how to use KubeDB to run PostgreSQL database with custom configuration.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace/demo created
```

> Note: YAML files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Overview

PostgreSQL allows to configure database via **Configuration File**, **SQL** and **Shell**. The most common way is to edit configuration file `postgresql.conf`. When PostgreSQL docker image starts, it uses the configuration specified in `postgresql.conf` file. This file can have `include` directive which allows to include configuration from other files. One of these `include` directives is `include_if_exists` which accept a file reference. If the referenced file exists, it includes configuration from the file. Otherwise, it uses default configuration. KubeDB takes advantage of this feature to allow users to provide their custom configuration. To know more about configuring PostgreSQL see [here](https://www.postgresql.org/docs/current/static/runtime-config.html).

At first, you have to create a config file named `user.conf` with your desired configuration. Then you have to put this file into a [volume](https://kubernetes.io/docs/concepts/storage/volumes/). You have to specify this volume in `spec.configSource` section while creating Postgres crd. KubeDB will mount this volume into `/etc/config/` directory of the database pod which will be referenced by `include_if_exists` directive.

In this tutorial, we will configure `max_connections` and `shared_buffers` via a custom config file. We will use configMap as volume source.

## Custom Configuration

At first, let's create `user.conf` file setting `max_connections` and `shared_buffers` parameters.

```ini
$ cat user.conf
max_connections=300
shared_buffers=256MB
```

> Note that config file name must be `user.conf`

Now, create a configMap with this configuration file.

```console
$ kubectl create configmap -n demo pg-custom-config --from-literal=user.conf="$(curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/postgres/custom-config/user.conf)"
configmap/pg-custom-config created
```

Verify the config map has the configuration file.

```yaml
$ kubectl get configmap -n demo pg-custom-config -o yaml
apiVersion: v1
data:
  user.conf: |-
    max_connections=300
    shared_buffers=256MB
kind: ConfigMap
metadata:
  creationTimestamp: "2019-02-07T12:08:26Z"
  name: pg-custom-config
  namespace: demo
  resourceVersion: "44214"
  selfLink: /api/v1/namespaces/demo/configmaps/pg-custom-config
  uid: 131b321f-2ad1-11e9-9d44-080027154f61
```

Now, create Postgres crd specifying `spec.configSource` field.

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/postgres/custom-config/pg-custom-config.yaml 
postgres.kubedb.com/custom-postgres created
```

Below is the YAML for the Postgres crd we just created.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: custom-postgres
  namespace: demo
spec:
  version: "9.6-v2"
  configSource:
    configMap:
      name: pg-custom-config
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

Now, wait a few minutes. KubeDB operator will create necessary PVC, statefulset, services, secret etc. If everything goes well, we will see that a pod with the name `custom-postgres-0` has been created.

Check that the statefulset's pod is running

```console
$ kubectl get pod -n demo custom-postgres-0
NAME                READY     STATUS    RESTARTS   AGE
custom-postgres-0   1/1       Running   0          14m
```

Check the pod's log to see if the database is ready

```console
$ kubectl logs -f -n demo custom-postgres-0
I0705 12:05:51.697190       1 logs.go:19] FLAG: --alsologtostderr="false"
I0705 12:05:51.717485       1 logs.go:19] FLAG: --enable-analytics="true"
I0705 12:05:51.717543       1 logs.go:19] FLAG: --help="false"
I0705 12:05:51.717558       1 logs.go:19] FLAG: --log_backtrace_at=":0"
I0705 12:05:51.717566       1 logs.go:19] FLAG: --log_dir=""
I0705 12:05:51.717573       1 logs.go:19] FLAG: --logtostderr="false"
I0705 12:05:51.717581       1 logs.go:19] FLAG: --stderrthreshold="0"
I0705 12:05:51.717589       1 logs.go:19] FLAG: --v="0"
I0705 12:05:51.717597       1 logs.go:19] FLAG: --vmodule=""
We want "custom-postgres-0" as our leader
I0705 12:05:52.753464       1 leaderelection.go:175] attempting to acquire leader lease  demo/custom-postgres-leader-lock...
I0705 12:05:52.822093       1 leaderelection.go:184] successfully acquired lease demo/custom-postgres-leader-lock
Got leadership, now do your jobs
Running as Primary
sh: locale: not found

WARNING: enabling "trust" authentication for local connections
You can change this by editing pg_hba.conf or using the option -A, or
--auth-local and --auth-host, the next time you run initdb.
ALTER ROLE
/scripts/primary/start.sh: ignoring /var/initdb/*

LOG:  database system was shut down at 2018-07-05 12:07:51 UTC
LOG:  MultiXact member wraparound protections are now enabled
LOG:  database system is ready to accept connections
LOG:  autovacuum launcher started
```

Once we see `LOG:  database system is ready to accept connections` in the log, the database is ready.

Now, we will check if the database has started with the custom configuration we have provided. We will `exec` into the pod and use [SHOW](https://www.postgresql.org/docs/9.6/static/sql-show.html) query to check the run-time parameters.

```console
 $ kubectl exec -it -n demo custom-postgres-0 sh
 / #
 ## login as user "postgres". no authentication required from inside the pod because it is using trust authentication local connection.
/ # psql -U postgres
psql (9.6.7)
Type "help" for help.

## query for "max_connections"
postgres=# SHOW max_connections;
 max_connections
-----------------
 300
(1 row)

## query for "shared_buffers"
postgres=# SHOW shared_buffers;
 shared_buffers
----------------
 256MB
(1 row)

## log out from database
postgres=# \q
/ #

```

You can also connect to this database from pgAdmin and use following SQL query to check these configuration.

```sql
SELECT name,setting
FROM pg_settings
WHERE name='max_connections' OR name='shared_buffers';
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo pg/custom-postgres -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo pg/custom-postgres

kubectl delete -n demo configmap pg-custom-config
kubectl delete ns demo
```

If you would like to uninstall KubeDB operator, please follow the steps [here](/docs/setup/uninstall.md).

## Next Steps

- Learn about [taking instant backup](/docs/guides/postgres/snapshot/instant_backup.md) of PostgreSQL database using KubeDB Snapshot.
- Learn how to [schedule backup](/docs/guides/postgres/snapshot/scheduled_backup.md)  of PostgreSQL database.
- Learn about initializing [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Learn about initializing [PostgreSQL from KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
