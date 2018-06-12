---
title: PostgreSQL Quickstart
menu:
  docs_0.8.0:
    identifier: pg-quickstart-quickstart
    name: Overview
    parent: pg-quickstart-postgres
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Running PostgreSQL

This tutorial will show you how to use KubeDB to run a PostgreSQL database.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/postgres/lifecycle.png" width="600" height="660">
</p>

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

This tutorial will also use a pgAdmin to connect and test PostgreSQL database, once it is running.

Run the following command to prepare your cluster for this tutorial

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/postgres/quickstart/pgadmin.yaml
deployment "pgadmin" created
service "pgadmin" created

$ kubectl get pods -n demo --watch
NAME                       READY     STATUS              RESTARTS   AGE
pgadmin-54688976f7-5rxfs   0/1       ContainerCreating   0          9s
pgadmin-54688976f7-5rxfs   1/1       Running   0         11s
^C⏎
```

Now, open pgAdmin in your browser by running `minikube service pgadmin -n demo`.

Or you can get the URL of Service `pgadmin` by running following command

```console
$ minikube service pgadmin -n demo --url
http://192.168.99.100:32231
```

To log into the pgAdmin, use username __`admin`__ and password __`admin`__.

> Note: Yaml files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

## Create a PostgreSQL database

KubeDB implements a Postgres CRD to define the specification of a PostgreSQL database.

Below is the Postgres object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: quick-postgres
  namespace: demo
spec:
  version: "9.6"
  doNotPause: true
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
```

Here,

- `spec.version` is the version of PostgreSQL database. In this tutorial, a PostgreSQL 9.6 database is created.
- `spec.doNotPause` prevents user from deleting this object if admission webhook is enabled.
- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. Since release 0.8.0, a storage spec is required for MySQL.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/postgres/quickstart/quick-postgres.yaml
postgres "quick-postgres" created
```

KubeDB operator watches for Postgres objects using Kubernetes api. When a Postgres object is created, KubeDB operator will create a new StatefulSet and two ClusterIP Service with the matching name.
KubeDB operator will also create a governing service for StatefulSet with the name `kubedb`, if one is not already present.

If RBAC is enabled in clusters, PostgreSQL specific RBAC permission is required. [Check here](/docs/guides/postgres/quickstart/rbac.md) for details.

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created.

```console
$ kubedb get pg -n demo quick-postgres -o wide
NAME             VERSION   STATUS    AGE
quick-postgres   9.6       Running   17m
```

Lets describe Postgres object `quick-postgres`

```console
$ kubedb describe pg -n demo quick-postgres
Name:           quick-postgres
Namespace:      demo
StartTimestamp: Thu, 08 Feb 2018 14:44:24 +0600
Status:         Running
Volume:
  StorageClass: standard
  Capacity:     50Mi
  Access Modes: RWO

StatefulSet:
  Name:                 quick-postgres
  Replicas:             1 current / 1 desired
  CreationTimestamp:    Thu, 08 Feb 2018 14:44:29 +0600
  Pods Status:          1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name: quick-postgres
  Type: ClusterIP
  IP:   10.104.43.68
  Port: api	5432/TCP

Service:
  Name: quick-postgres-replicas
  Type: ClusterIP
  IP:   10.96.98.122
  Port: api	5432/TCP

Database Secret:
  Name:	quick-postgres-auth
  Type:	Opaque
  Data
  ====
  POSTGRES_PASSWORD:    16 bytes

Topology:
  Type      Pod                StartTime                       Phase
  ----      ---                ---------                       -----
  primary   quick-postgres-0   2018-02-08 14:44:29 +0600 +06   Running

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                Type       Reason       Message
  ---------   --------   -----     ----                --------   ------       -------
  2m          2m         1         Postgres operator   Normal     Successful   Successfully patched StatefulSet
  2m          2m         1         Postgres operator   Normal     Successful   Successfully patched Postgres
  2m          2m         1         Postgres operator   Normal     Successful   Successfully created StatefulSet
  2m          2m         1         Postgres operator   Normal     Successful   Successfully created Postgres
  4m          4m         1         Postgres operator   Normal     Successful   Successfully created Service
  4m          4m         1         Postgres operator   Normal     Successful   Successfully created Service
```

```console
$ kubectl get service -n demo --selector=kubedb.com/kind=Postgres,kubedb.com/name=quick-postgres
NAME                        TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
quick-postgres              ClusterIP   10.104.43.68   <none>        5432/TCP   5m
quick-postgres-replicas     ClusterIP   10.96.98.122   <none>        5432/TCP   5m
```

Two services for each Postgres object.

- Service *`quick-postgres`* targets only one Pod which is acting as *primary* server
- Service *`quick-postgres-replicas`* targets all Pods created by StatefulSet

KubeDB supports PostgreSQL clustering where Pod can be either *primary* or *standby*.
To learn how to configure highly available PostgreSQL cluster, click [here](/docs/guides/postgres/clustering/ha_cluster.md).

Here, we create a PostgreSQL database with single node, *primary* only.

Please note that KubeDB operator has created a new Secret called `quick-postgres-auth` for storing the password for `postgres` superuser.

```yaml
apiVersion: v1
data:
  POSTGRES_PASSWORD: OXFyUktkNWFBT1JnSC1hVg==
kind: Secret
metadata:
  creationTimestamp: 2018-02-08T08:44:29Z
  labels:
    kubedb.com/kind: Postgres
    kubedb.com/name: quick-postgres
  name: quick-postgres-auth
  namespace: demo
  resourceVersion: "19311"
  selfLink: /api/v1/namespaces/demo/secrets/quick-postgres-auth
  uid: 46ac9a71-0cac-11e8-b21d-0800273fbab1
type: Opaque
```

This Secret contains `postgres` superuser password as `POSTGRES_PASSWORD` key.

> Note: Auth Secret name format: `{postgres-name}-auth`

Now, you can connect to this database from the pgAdmin dashboard using Service `quick-postgres.demo` and `postgres` superuser password .

Connection information:

- address: you can use any of these
  - Service `quick-postgres.demo`
  - Pod IP (`$ kubectl get pods quick-postgres-0 -n demo -o yaml | grep podIP`)
- port: `5432`
- database: `postgres`
- username: `postgres`

Run following command to get `postgres` superuser password

    $ kubectl get secrets -n demo quick-postgres-auth -o jsonpath='{.data.\POSTGRES_PASSWORD}' | base64 -d

<p align="center">
  <kbd>
    <img alt="quick-postgres"  src="/docs/images/postgres/quick-postgres.gif">
  </kbd>
</p>

## Pause Database

KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `doNotPause` feature.
If admission webhook is enabled, It prevents user from deleting the database as long as the `spec.doNotPause` is set `true`.

In this tutorial, Postgres `quick-postgres` is created with `spec.doNotPause: true`. So, if you delete this Postgres object, admission webhook will nullify the delete operation.

```console
$ kubedb delete pg -n demo quick-postgres
error: Postgres "quick-postgres " can't be paused. To continue delete, unset spec.doNotPause and retry.
```

To continue with this tutorial, unset `spec.doNotPause` by updating Postgres object

```console
$ kubedb edit pg -n demo quick-postgres
spec:
  doNotPause: false
```

Now, if you delete the Postgres object, KubeDB operator will create a matching DormantDatabase object. KubeDB operator watches for DormantDatabase objects and it will take necessary steps when a DormantDatabase object is created.

KubeDB operator will delete the StatefulSet and its Pods, but leaves the Secret, PVCs unchanged.

```console
$ kubedb delete pg -n demo quick-postgres
postgres "quick-postgres" deleted
```

Check DormantDatabase entry

```console
$ kubedb get drmn -n demo quick-postgres
NAME             STATUS    AGE
quick-postgres   Paused    19s
```

In KubeDB parlance, we say that Postgres `quick-postgres`  has entered into dormant state.

Lets see, what we have in this DormantDatabase object

```yaml
$ kubedb get drmn -n demo quick-postgres -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-08T09:33:22Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: Postgres
  name: quick-postgres
  namespace: demo
  resourceVersion: "24091"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/quick-postgres
  uid: 1ab5628d-0cb3-11e8-b21d-0800273fbab1
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: quick-postgres
      namespace: demo
    spec:
      postgres:
        databaseSecret:
          secretName: quick-postgres-auth
        resources: {}
        storage:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 50Mi
          storageClassName: standard
        version: "9.6"
status:
  creationTime: 2018-02-08T09:33:22Z
  pausingTime: 2018-02-08T09:33:24Z
  phase: Paused
```

Here,

- `spec.origin` contains original Postgres object.
- `status.phase` points to the current database state `Paused`.

## Resume DormantDatabase

To resume the database from the dormant state, create same Postgres object with same Spec.

In this tutorial, the DormantDatabase `quick-postgres` can be resumed by creating original Postgres object.

The below command will resume the DormantDatabase `quick-postgres`

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/postgres/quickstart/quick-postgres.yaml
postgres "quick-postgres" created
```

## WipeOut DormantDatabase

You can wipe out a DormantDatabase while deleting the objet by setting `spec.wipeOut` to true. KubeDB operator will delete any relevant resources of this `Elasticsearch` database (i.e, PVCs, Secrets, Snapshots). It will also delete snapshot data stored in the Cloud Storage buckets.

```yaml
$ kubedb edit drmn -n demo quick-postgres
spec:
  wipeOut: true
```

If `spec.wipeOut` is not set to true while deleting the `dormantdatabase` object, then only this object will be deleted and `kubedb-operator` won't delete related Secrets, PVCs and Snapshots. So, user still can access the stored data in the cloud storage buckets as well as PVCs.

## Delete DormantDatabase

As it is already discussed above, `DormantDatabase` can be deleted with or without wiping out the resources. To delete the `dormantdatabase`,

```console
$ kubedb delete drmn -n demo quick-postgres
dormantdatabase "quick-postgres" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo pg/quick-postgres -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo pg/quick-postgres

$ kubectl patch -n demo drmn/quick-postgres -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/quick-postgres

$ kubectl delete ns demo
```

## Next Steps

- Learn about [taking instant backup](/docs/guides/postgres/snapshot/instant_backup.md) of PostgreSQL database using KubeDB Snapshot.
- Learn how to [schedule backup](/docs/guides/postgres/snapshot/scheduled_backup.md)  of PostgreSQL database.
- Learn about initializing [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Learn about initializing [PostgreSQL from KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Detail concepts of [Postgres object](/docs/concepts/databases/postgres.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Use [private Docker registry](/docs/guides/postgres/private-registry/using-private-registry.md) to deploy PostgreSQL with KubeDB.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
