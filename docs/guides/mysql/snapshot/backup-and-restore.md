---
title: Instant Backup of MySQL
menu:
  docs_0.9.0-beta.0:
    identifier: my-backup-and-restore-snapshot
    name: Instant Backup
    parent: my-snapshot-mysql
    weight: 10
menu_name: docs_0.9.0-beta.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Database Snapshots

This tutorial will show you how to take snapshots of a KubeDB managed MySQL database.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

A `MySQL` database is needed to take snapshot for this tutorial. To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    1h
demo          Active    1m
kube-public   Active    1h
kube-system   Active    1h

$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.1/docs/examples/mysql/snapshot/demo-1.yaml
mysql "mysql-infant" created
```

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Instant Backups

You can easily take a snapshot of `MySQL` database by creating a `Snapshot` object. When a `Snapshot` object is created, KubeDB operator will launch a Job that runs the `mysql dump` command and uploads the output bson file to various cloud providers S3, GCS, Azure, OpenStack Swift and/or locally mounted volumes using [osm](https://github.com/appscode/osm).

In this tutorial, snapshots will be stored in a Google Cloud Storage (GCS) bucket. To do so, a secret is needed that has the following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic my-snap-secret -n demo \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "my-snap-secret" created
```

```yaml
$ kubectl get secret my-snap-secret -n demo -o yaml
apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2018-02-09T12:02:08Z
  name: my-snap-secret
  namespace: demo
  resourceVersion: "30349"
  selfLink: /api/v1/namespaces/demo/secrets/my-snap-secret
  uid: 0dccee80-0d91-11e8-9091-08002751ae8c
type: Opaque
```

To lean how to configure other storage destinations for Snapshots, please visit [here](/docs/concepts/snapshot.md). Now, create the Snapshot object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: snap-mysql-infant
  namespace: demo
  labels:
    kubedb.com/kind: MySQL
spec:
  databaseName: mysql-infant
  storageSecretName: my-snap-secret
  gcs:
    bucket: restic
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.1/docs/examples/mysql/snapshot/demo-2.yaml
snapshot "snap-mysql-infant" created

$ kubedb get snap -n demo
NAME                DATABASE          STATUS    AGE
snap-mysql-infant   my/mysql-infant   Running   22s
```

```yaml
$ kubedb get snap -n demo snap-mysql-infant -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-09T12:03:50Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: MySQL
    kubedb.com/name: mysql-infant
  name: snap-mysql-infant
  namespace: demo
  resourceVersion: "30488"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/snapshots/snap-mysql-infant
  uid: 4a507251-0d91-11e8-9091-08002751ae8c
spec:
  databaseName: mysql-infant
  gcs:
    bucket: restic
  storageSecretName: my-snap-secret
status:
  completionTime: 2018-02-09T12:04:52Z
  phase: Succeeded
  startTime: 2018-02-09T12:03:50Z
```

Here,

- `metadata.labels` should include the type of database `kubedb.com/kind: MySQL` whose snapshot will be taken.
- `spec.databaseName` points to the database whose snapshot is taken.
- `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.
- `spec.gcs.bucket` points to the bucket name used to store the snapshot data.

You can also run the `kubedb describe` command to see the recent snapshots taken for a database.

```console
$ kubedb describe my -n demo mysql-infant
Name:		mysql-infant
Namespace:	demo
StartTimestamp:	Fri, 09 Feb 2018 18:00:23 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			mysql-infant
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Fri, 09 Feb 2018 18:00:24 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		mysql-infant
  Type:		ClusterIP
  IP:		10.103.94.148
  Port:		db	3306/TCP

Database Secret:
  Name:	mysql-infant-auth
  Type:	Opaque
  Data
  ====
  password:	16 bytes
  user:		4 bytes

Snapshots:
  Name                Bucket      StartTime                         CompletionTime                    Phase
  ----                ------      ---------                         --------------                    -----
  snap-mysql-infant   gs:restic   Fri, 09 Feb 2018 18:03:50 +0600   Fri, 09 Feb 2018 18:04:52 +0600   Succeeded

Events:
  FirstSeen   LastSeen   Count     From                  Type       Reason               Message
  ---------   --------   -----     ----                  --------   ------               -------
  1m          1m         1         Job Controller        Normal     SuccessfulSnapshot   Successfully completed snapshot
  2m          2m         1         Snapshot Controller   Normal     Starting             Backup running
  5m          5m         1         MySQL operator        Normal     Successful           Successfully patched StatefulSet
  5m          5m         1         MySQL operator        Normal     Successful           Successfully patched MySQL
  5m          5m         1         MySQL operator        Normal     Successful           Successfully created StatefulSet
  5m          5m         1         MySQL operator        Normal     Successful           Successfully created MySQL
  5m          5m         1         MySQL operator        Normal     Successful           Successfully created Service
```

Once the snapshot Job is complete, you should see the output of the `mysql dump` command stored in the GCS bucket.

![snapshot-console](/docs/images/mysql/m1-xyz-snapshot.png)

From the above image, you can see that the snapshot output is stored in a folder called `{bucket}/kubedb/{namespace}/{mysql-object}/{snapshot}/`.

## Restore from Snapshot

You can create a new database from a previously taken Snapshot. Specify the Snapshot name in the `spec.init.snapshotSource` field of a new MySQL object. See the example `mysql-recovered` object below:

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-recovered
  namespace: demo
spec:
  version: "8.0"
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    snapshotSource:
      name: snap-mysql-infant
      namespace: demo
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.1/docs/examples/mysql/snapshot/demo-3.yaml
mysql "mysql-recovered" created
```

Here,

- `spec.init.snapshotSource.name` refers to a Snapshot object for a MySQL database in the same namespaces as this new `mysql-recovered` MySQL object.

Now, wait several seconds. KubeDB operator will create a new StatefulSet. Then KubeDB operator launches a Kubernetes Job to initialize the new database using the data from `snap-mysql-infant` Snapshot.

```console
$ kubedb get my -n demo
NAME              STATUS         AGE
mysql-infant      Running        8m
mysql-recovered   Initializing   21s

$ kubedb get my -n demo
$ NAME              STATUS    AGE
mysql-infant      Running     14m
mysql-recovered   Running     6m

$ kubedb describe my -n demo mysql-recovered
Name:		mysql-recovered
Namespace:	demo
StartTimestamp:	Fri, 09 Feb 2018 18:08:23 +0600
Status:		Running
Annotations:	kubedb.com/initialized=
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			mysql-recovered
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Fri, 09 Feb 2018 18:08:25 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		mysql-recovered
  Type:		ClusterIP
  IP:		10.105.110.215
  Port:		db	3306/TCP

Database Secret:
  Name:	mysql-recovered-auth
  Type:	Opaque
  Data
  ====
  user:		4 bytes
  password:	16 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From             Type       Reason               Message
  ---------   --------   -----     ----             --------   ------               -------
  5m          5m         1         Job Controller   Normal     SuccessfulSnapshot   Successfully completed initialization
  10m         10m        1         MySQL operator   Normal     Successful           Successfully patched StatefulSet
  10m         10m        1         MySQL operator   Normal     Successful           Successfully patched MySQL
  10m         10m        1         MySQL operator   Normal     Initializing         Initializing from Snapshot: "snap-mysql-infant"
  10m         10m        1         MySQL operator   Normal     Successful           Successfully created StatefulSet
  10m         10m        1         MySQL operator   Normal     Successful           Successfully created MySQL
  10m         10m        1         MySQL operator   Normal     Successful           Successfully created Service
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo mysql/mysql-infant mysql/mysql-recovered -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo mysql/mysql-infant mysql/mysql-recovered

$ kubectl patch -n demo drmn/mysql-infant drmn/mysql-recovered -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/mysql-infant drmn/mysql-recovered

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Take [Scheduled Snapshot](/docs/guides/mysql/snapshot/scheduled-backup.md) of MySQL databases using KubeDB.
- Initialize [MySQL with Script](/docs/guides/mysql/initialization/using-script.md).
- Initialize [MySQL with Snapshot](/docs/guides/mysql/initialization/using-snapshot.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
