---
title: Using Postgres Streaming Replication
menu:
  docs_0.9.0:
    identifier: pg-streaming-replication-clustering
    name: Streaming Replication
    parent: pg-clustering-postgres
    weight: 15
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Streaming Replication

Streaming Replication provides *asynchronous* replication to one or more *standby* servers.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace/demo created
```

> Note: YAML files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create PostgreSQL with Streaming replication

The example below demonstrates KubeDB PostgreSQL for Streaming Replication

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: ha-postgres
  namespace: demo
spec:
  version: "9.6-v2"
  replicas: 3
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

In this examples:

- This `Postgres` object creates three PostgreSQL servers, indicated by the **`replicas`** field.
- One server will be *primary* and two others will be *warm standby* servers, default of **`spec.standbyMode`**

### What is Streaming Replication

Streaming Replication allows a *standby* server to stay more up-to-date by shipping and applying the [WAL XLOG](http://www.postgresql.org/docs/9.6/static/wal.html)
records continuously. The *standby* connects to the *primary*, which streams WAL records to the *standby* as they're generated, without waiting for the WAL file to be filled.

Streaming Replication is **asynchronous** by default. As a result, there is a small delay between committing a transaction in the *primary* and the changes becoming visible in the *standby*.

### Streaming Replication setup

Following parameters are set in `postgresql.conf` for both *primary* and *standby* server

```console
wal_level = replica
max_wal_senders = 99
wal_keep_segments = 32
```

Here,

- _wal_keep_segments_ specifies the minimum number of past log file segments kept in the pg_xlog directory.

And followings are in `recovery.conf` for *standby* server

```console
standby_mode = on
trigger_file = '/tmp/pg-failover-trigger'
recovery_target_timeline = 'latest'
primary_conninfo = 'application_name=$HOSTNAME host=$PRIMARY_HOST'
```

Here,

- _trigger_file_ is created to trigger a *standby* to take over as *primary* server.
- *$PRIMARY_HOST* holds the Kubernetes Service name that targets *primary* server

Now create this Postgres object with Streaming Replication support

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/postgres/clustering/ha-postgres.yaml
postgres.kubedb.com/ha-postgres created
```

KubeDB operator creates three Pod as PostgreSQL server.

```console
$ kubectl get pods -n demo --selector="kubedb.com/name=ha-postgres" --show-labels
NAME            READY   STATUS    RESTARTS   AGE   LABELS
ha-postgres-0   1/1     Running   0          20s   controller-revision-hash=ha-postgres-6b7998ccfd,kubedb.com/kind=Postgres,kubedb.com/name=ha-postgres,kubedb.com/role=primary,statefulset.kubernetes.io/pod-name=ha-postgres-0
ha-postgres-1   1/1     Running   0          16s   controller-revision-hash=ha-postgres-6b7998ccfd,kubedb.com/kind=Postgres,kubedb.com/name=ha-postgres,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=ha-postgres-1
ha-postgres-2   1/1     Running   0          10s   controller-revision-hash=ha-postgres-6b7998ccfd,kubedb.com/kind=Postgres,kubedb.com/name=ha-postgres,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=ha-postgres-2
```

Here,

- Pod `ha-postgres-0` is serving as *primary* server, indicated by label `kubedb.com/role=primary`
- Pod `ha-postgres-1` & `ha-postgres-2` both are serving as *standby* server, indicated by label `kubedb.com/role=replica`

And two services for Postgres `ha-postgres` are created.

```console
$ kubectl get svc -n demo --selector="kubedb.com/name=ha-postgres"
NAME                   TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
ha-postgres            ClusterIP   10.102.19.49   <none>        5432/TCP   4m
ha-postgres-replicas   ClusterIP   10.97.36.117   <none>        5432/TCP   4m
```

```console
$ kubectl get svc -n demo --selector="kubedb.com/name=ha-postgres" -o=custom-columns=NAME:.metadata.name,SELECTOR:.spec.selector
NAME                   SELECTOR
ha-postgres            map[kubedb.com/kind:Postgres kubedb.com/name:ha-postgres kubedb.com/role:primary]
ha-postgres-replicas   map[kubedb.com/kind:Postgres kubedb.com/name:ha-postgres]
```

Here,

- Service `ha-postgres` targets Pod `ha-postgres-0`, which is *primary* server, by selector `kubedb.com/kind=Postgres,kubedb.com/name=ha-postgres,kubedb.com/role=primary`.
- Service `ha-postgres-replicas` targets all Pods (*`ha-postgres-0`*, *`ha-postgres-1`* and *`ha-postgres-2`*) with label `kubedb.com/kind=Postgres,kubedb.com/name=ha-postgres`.

>These *standby* servers are asynchronous *warm standby* server. That means, you can only connect to *primary* sever.

Now connect to this *primary* server Pod `ha-postgres-0` using pgAdmin installed in [quickstart](/docs/guides/postgres/quickstart/quickstart.md#before-you-begin) tutorial.

**Connection information:**

- Host name/address: you can use any of these
  - Service: `ha-postgres.demo`
  - Pod IP: (`$kubectl get pods ha-postgres-0 -n demo -o yaml | grep podIP`)
- Port: `5432`
- Maintenance database: `postgres`
- Username: Run following command to get *username*,

  ```console
  $ kubectl get secrets -n demo ha-postgres-auth -o jsonpath='{.data.\POSTGRES_USER}' | base64 -d
  postgres
  ```

- Password: Run the following command to get *password*,

  ```console
  $ kubectl get secrets -n demo ha-postgres-auth -o jsonpath='{.data.\POSTGRES_PASSWORD}' | base64 -d
  MHRrOcuyddfh3YpU
  ```

You can check `pg_stat_replication` information to know who is currently streaming from *primary*.

```console
postgres=# select * from pg_stat_replication;
```

 pid | usesysid | usename  | application_name | client_addr | client_port |         backend_start         |   state   | sent_location | write_location | flush_location | replay_location | sync_priority | sync_state
-----|----------|----------|------------------|-------------|-------------|-------------------------------|-----------|---------------|----------------|----------------|-----------------|---------------|------------
  89 |       10 | postgres | ha-postgres-2    | 172.17.0.8  |       35306 | 2018-02-09 04:27:11.674828+00 | streaming | 0/5000060     | 0/5000060      | 0/5000060      | 0/5000060       |             0 | async
  90 |       10 | postgres | ha-postgres-1    | 172.17.0.7  |       42400 | 2018-02-09 04:27:13.716104+00 | streaming | 0/5000060     | 0/5000060      | 0/5000060      | 0/5000060       |             0 | async

Here, both `ha-postgres-1` and `ha-postgres-2` are streaming asynchronously from *primary* server.

### Lease Duration

Get the postgres CRD at this point.

```yaml
$ kubectl get pg -n demo   ha-postgres -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  creationTimestamp: "2019-02-07T12:14:05Z"
  finalizers:
  - kubedb.com
  generation: 2
  name: ha-postgres
  namespace: demo
  resourceVersion: "44966"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/postgreses/ha-postgres
  uid: dcf6d96a-2ad1-11e9-9d44-080027154f61
spec:
  databaseSecret:
    secretName: ha-postgres-auth
  leaderElection:
    leaseDurationSeconds: 15
    renewDeadlineSeconds: 10
    retryPeriodSeconds: 2
  podTemplate:
    controller: {}
    metadata: {}
    spec:
      resources: {}
  replicas: 3
  serviceTemplate:
    metadata: {}
    spec: {}
  storage:
    accessModes:
    - ReadWriteOnce
    dataSource: null
    resources:
      requests:
        storage: 1Gi
    storageClassName: standard
  storageType: Durable
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
  version: 9.6-v2
status:
  observedGeneration: 2$4213139756412538772
  phase: Running
```

There are three fields in Postgres CRD's `spec.LeaderElection`. These values defines how fast the leader election can happen.

- leaseDurationSeconds: LeaseDuration is the duration in second that non-leader candidates will wait to force acquire leadership. This is measured against time of last observed ack. Default 15
- renewDeadlineSeconds: RenewDeadline is the duration in second that the acting master will retry refreshing leadership before giving up. Normally, LeaseDuration * 2 / 3. Default 10
- retryPeriodSeconds: RetryPeriod is the duration in second the LeaderElector clients should wait between tries of actions. Normally, LeaseDuration / 3. Default 2

If the Cluster machine is powerful, user can reduce the times. But, Do not make it so little, in that case Postgres will restarts very often. 

### Automatic failover

If *primary* server fails, another *standby* server will take over and serve as *primary*.

Delete Pod `ha-postgres-0` to see the failover behavior.

```console
kubectl delete pod -n demo ha-postgres-0
```

```console
$ kubectl get pods -n demo --selector="kubedb.com/name=ha-postgres" --show-labels
NAME            READY     STATUS    RESTARTS   AGE       LABELS
ha-postgres-0   1/1       Running   0          10s       controller-revision-hash=ha-postgres-b8b4b5fc4,kubedb.com/kind=Postgres,kubedb.com/name=ha-postgres,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=ha-postgres-0
ha-postgres-1   1/1       Running   0          52m       controller-revision-hash=ha-postgres-b8b4b5fc4,kubedb.com/kind=Postgres,kubedb.com/name=ha-postgres,kubedb.com/role=primary,statefulset.kubernetes.io/pod-name=ha-postgres-1
ha-postgres-2   1/1       Running   0          51m       controller-revision-hash=ha-postgres-b8b4b5fc4,kubedb.com/kind=Postgres,kubedb.com/name=ha-postgres,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=ha-postgres-2
```

Here,

- Pod `ha-postgres-1` is now serving as *primary* server
- Pod `ha-postgres-0` and `ha-postgres-2` both are serving as *standby* server

And result from `pg_stat_replication`

```console
postgres=# select * from pg_stat_replication;
```

 pid | usesysid | usename  | application_name | client_addr | client_port |         backend_start         |   state   | sent_location | write_location | flush_location | replay_location | sync_priority | sync_state
-----|----------|----------|------------------|-------------|-------------|-------------------------------|-----------|---------------|----------------|----------------|-----------------|---------------|------------
  57 |       10 | postgres | ha-postgres-0    | 172.17.0.6  |       52730 | 2018-02-09 04:33:06.051716|00 | streaming | 0/7000060     | 0/7000060      | 0/7000060      | 0/7000060       |             0 | async
  58 |       10 | postgres | ha-postgres-2    | 172.17.0.8  |       42824 | 2018-02-09 04:33:09.762168|00 | streaming | 0/7000060     | 0/7000060      | 0/7000060      | 0/7000060       |             0 | async

You can see here, now `ha-postgres-0` and `ha-postgres-2` are streaming asynchronously from `ha-postgres-1`, our *primary* server.

<p align="center">
  <kbd>
    <img alt="recovered-postgres"  src="/docs/images/postgres/ha-postgres.gif">
  </kbd>
</p>

[//]: # (If you want to know how this failover process works, [read here])

## Streaming Replication with `hot standby`

Streaming Replication also works with one or more *hot standby* servers.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: hot-postgres
  namespace: demo
spec:
  version: "9.6-v2"
  replicas: 3
  standbyMode: Hot
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

In this examples:

- This `Postgres` object creates three PostgreSQL servers, indicated by the **`replicas`** field.
- One server will be *primary* and two others will be *hot standby* servers, as instructed by **`spec.standbyMode`**

### `hot standby` setup

Following parameters are set in `postgresql.conf` for *standby* server

```console
hot_standby = on
```

Here,

- _hot_standby_ specifies that *standby* server will act as *hot standby*.

Now create this Postgres object

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/postgres/clustering/hot-postgres.yaml
postgres "hot-postgres" created
```

KubeDB operator creates three Pod as PostgreSQL server.

```console
$ kubectl get pods -n demo --selector="kubedb.com/name=hot-postgres" --show-labels
NAME             READY     STATUS    RESTARTS   AGE       LABELS
hot-postgres-0   1/1       Running   0          1m        controller-revision-hash=hot-postgres-6c48cfb5bb,kubedb.com/kind=Postgres,kubedb.com/name=hot-postgres,kubedb.com/role=primary,statefulset.kubernetes.io/pod-name=hot-postgres-0
hot-postgres-1   1/1       Running   0          1m        controller-revision-hash=hot-postgres-6c48cfb5bb,kubedb.com/kind=Postgres,kubedb.com/name=hot-postgres,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=hot-postgres-1
hot-postgres-2   1/1       Running   0          48s       controller-revision-hash=hot-postgres-6c48cfb5bb,kubedb.com/kind=Postgres,kubedb.com/name=hot-postgres,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=hot-postgres-2
```

Here,

- Pod `hot-postgres-0` is serving as *primary* server, indicated by label `kubedb.com/role=primary`
- Pod `hot-postgres-1` & `hot-postgres-2` both are serving as *standby* server, indicated by label `kubedb.com/role=replica`

> These *standby* servers are asynchronous *hot standby* servers.

That means, you can connect to both *primary* and *standby* sever. But these *hot standby* servers only accept read-only queries.

Now connect to one of our *hot standby* servers Pod `hot-postgres-2` using pgAdmin installed in [quickstart](/docs/guides/postgres/quickstart/quickstart.md#before-you-begin) tutorial.

**Connection information:**

- Host name/address: you can use any of these
  - Service: `hot-postgres-replicas.demo`
  - Pod IP: (`$kubectl get pods hot-postgres-2 -n demo -o yaml | grep podIP`)
- Port: `5432`
- Maintenance database: `postgres`
- Username: Run following command to get *username*,

  ```console
  $ kubectl get secrets -n demo hot-postgres-auth -o jsonpath='{.data.\POSTGRES_USER}' | base64 -d
  postgres
  ```

- Password: Run the following command to get *password*,

  ```console
  $ kubectl get secrets -n demo hot-postgres-auth -o jsonpath='{.data.\POSTGRES_PASSWORD}' | base64 -d
  ZZgjjQMUdKJYy1W9
  ```

Try to create a database (write operation)

```console
postgres=# CREATE DATABASE standby;
ERROR:  cannot execute CREATE DATABASE in a read-only transaction
```

Failed to execute write operation. But it can execute following read query

```console
postgres=# select pg_last_xlog_receive_location();
 pg_last_xlog_receive_location
-------------------------------
 0/7000220
```

So, you can see here that you can connect to *hot standby* and it only accepts read-only queries.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo pg/ha-postgres pg/hot-postgres -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo pg/ha-postgres pg/hot-postgres

$ kubectl delete ns demo
```

## Next Steps

- Setup [Continuous Archiving](/docs/guides/postgres/snapshot/continuous_archiving.md) in PostgreSQL using `wal-g`
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
