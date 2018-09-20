---
title: Instant Backup of MongoDB
menu:
  docs_0.8.0:
    identifier: mg-backup-and-restore-snapshot
    name: Instant Backup
    parent: mg-snapshot-mongodb
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Database Snapshots

This tutorial will show you how to take snapshots of a KubeDB managed MongoDB database.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

A `MongoDB` database is needed to take snapshot for this tutorial. To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    1h
demo          Active    1m
kube-public   Active    1h
kube-system   Active    1h

$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.0/docs/examples/mongodb/snapshot/demo-1.yaml
mongodb "mgo-infant" created
```

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Instant Backups

You can easily take a snapshot of `MongoDB` database by creating a `Snapshot` object. When a `Snapshot` object is created, KubeDB operator will launch a Job that runs the `mongodump` command and uploads the output bson file to various cloud providers S3, GCS, Azure, OpenStack Swift and/or locally mounted volumes using [osm](https://github.com/appscode/osm).

In this tutorial, snapshots will be stored in a Google Cloud Storage (GCS) bucket. To do so, a secret is needed that has the following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic mg-snap-secret -n demo \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "mg-snap-secret" created
```

```yaml
$ kubectl get secret mg-snap-secret -n demo -o yaml
apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2018-02-02T10:02:09Z
  name: mg-snap-secret
  namespace: demo
  resourceVersion: "48679"
  selfLink: /api/v1/namespaces/demo/secrets/mg-snap-secret
  uid: 220a7c60-0800-11e8-946f-080027c05a6e
type: Opaque
```

To lean how to configure other storage destinations for Snapshots, please visit [here](/docs/concepts/snapshot.md). Now, create the Snapshot object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: snapshot-infant
  namespace: demo
  labels:
    kubedb.com/kind: MongoDB
spec:
  databaseName: mgo-infant
  storageSecretName: mg-snap-secret
  gcs:
    bucket: restic
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.0/docs/examples/mongodb/snapshot/demo-2.yaml
snapshot "snapshot-infant" created

$ kubedb get snap -n demo
NAME              DATABASE        STATUS    AGE
snapshot-infant   mg/mgo-infant   Running   47s
```

```yaml
$ kubedb get snap -n demo snapshot-infant -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-02T10:05:36Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: MongoDB
    kubedb.com/name: mgo-infant
  name: snapshot-infant
  namespace: demo
  resourceVersion: "48991"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/snapshots/snapshot-infant
  uid: 9d4f37a0-0800-11e8-946f-080027c05a6e
spec:
  databaseName: mgo-infant
  gcs:
    bucket: restic
  storageSecretName: mg-snap-secret
status:
  completionTime: 2018-02-02T10:06:43Z
  phase: Succeeded
  startTime: 2018-02-02T10:05:37Z
```

Here,

- `metadata.labels` should include the type of database `kubedb.com/kind: MongoDB` whose snapshot will be taken.
- `spec.databaseName` points to the database whose snapshot is taken.
- `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.
- `spec.gcs.bucket` points to the bucket name used to store the snapshot data.

You can also run the `kubedb describe` command to see the recent snapshots taken for a database.

```console
$ kubedb describe mg -n demo mgo-infant
Name:		mgo-infant
Namespace:	demo
StartTimestamp:	Fri, 02 Feb 2018 16:04:50 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			mgo-infant
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Fri, 02 Feb 2018 16:04:56 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		mgo-infant
  Type:		ClusterIP
  IP:		10.99.34.23
  Port:		db	27017/TCP

Database Secret:
  Name:	mgo-infant-auth
  Type:	Opaque
  Data
  ====
  password:	16 bytes
  user:		4 bytes

Snapshots:
  Name              Bucket      StartTime                         CompletionTime                    Phase
  ----              ------      ---------                         --------------                    -----
  snapshot-infant   gs:restic   Fri, 02 Feb 2018 16:05:37 +0600   Fri, 02 Feb 2018 16:06:43 +0600   Succeeded

Events:
  FirstSeen   LastSeen   Count     From                  Type       Reason               Message
  ---------   --------   -----     ----                  --------   ------               -------
  8m          8m         1         Job Controller        Normal     SuccessfulSnapshot   Successfully completed snapshot
  9m          9m         1         Snapshot Controller   Normal     Starting             Backup running
  9m          9m         1         MongoDB operator      Normal     Successful           Successfully patched StatefulSet
  9m          9m         1         MongoDB operator      Normal     Successful           Successfully patched MongoDB
  9m          9m         1         MongoDB operator      Normal     Successful           Successfully created StatefulSet
  9m          9m         1         MongoDB operator      Normal     Successful           Successfully created MongoDB
  10m         10m        1         MongoDB operator      Normal     Successful           Successfully created Service
```

Once the snapshot Job is complete, you should see the output of the `mongodump` command stored in the GCS bucket.

![snapshot-console](/docs/images/mongodb/mgo-xyz-snapshot.png)

From the above image, you can see that the snapshot output is stored in a folder called `{bucket}/kubedb/{namespace}/{mongodb-object}/{snapshot}/`.

## Restore from Snapshot

You can create a new database from a previously taken Snapshot. Specify the Snapshot name in the `spec.init.snapshotSource` field of a new MongoDB object. See the example `mgo-recovered` object below:

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-recovered
  namespace: demo
spec:
  version: "3.4"
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    snapshotSource:
      name: snapshot-infant
      namespace: demo
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.0/docs/examples/mongodb/snapshot/demo-3.yaml
mongodb "mgo-recovered" created
```

Here,

- `spec.init.snapshotSource.name` refers to a Snapshot object for a MongoDB database in the same namespaces as this new `mgo-recovered` MongoDB object.

Now, wait several seconds. KubeDB operator will create a new StatefulSet. Then KubeDB operator launches a Kubernetes Job to initialize the new database using the data from `snapshot-infant` Snapshot.

```console
$ kubedb get mg -n demo
NAME            STATUS    AGE
mgo-infant      Running   23m
mgo-recovered   Running   4m

$ kubedb describe mg -n demo mgo-recovered
Name:		mgo-recovered
Namespace:	demo
StartTimestamp:	Fri, 02 Feb 2018 16:24:23 +0600
Status:		Running
Annotations:	kubedb.com/initialized=
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			mgo-recovered
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Fri, 02 Feb 2018 16:24:36 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		mgo-recovered
  Type:		ClusterIP
  IP:		10.107.157.253
  Port:		db	27017/TCP

Database Secret:
  Name:	mgo-recovered-auth
  Type:	Opaque
  Data
  ====
  user:		4 bytes
  password:	16 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From               Type       Reason                 Message
  ---------   --------   -----     ----               --------   ------                 -------
  3m          3m         1         MongoDB operator   Normal     Successful             Successfully patched StatefulSet
  3m          3m         1         MongoDB operator   Normal     Successful             Successfully patched MongoDB
  3m          3m         1         Job Controller     Normal     SuccessfulInitialize   Successfully completed initialization
  4m          4m         1         MongoDB operator   Normal     Successful             Successfully patched StatefulSet
  4m          4m         1         MongoDB operator   Normal     Successful             Successfully patched MongoDB
  4m          4m         1         MongoDB operator   Normal     Initializing           Initializing from Snapshot: "snapshot-infant"
  4m          4m         1         MongoDB operator   Normal     Successful             Successfully created StatefulSet
  4m          4m         1         MongoDB operator   Normal     Successful             Successfully created MongoDB
  4m          4m         1         MongoDB operator   Normal     Successful             Successfully created Service
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo mg/mgo-infant mg/mgo-recovered -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo mg/mgo-infant mg/mgo-recovered

$ kubectl patch -n demo drmn/mgo-infant drmn/mgo-recovered -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/mgo-infant drmn/mgo-recovered

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Take [Scheduled Snapshot](/docs/guides/mongodb/snapshot/scheduled-backup.md) of MongoDB databases using KubeDB.
- Initialize [MongoDB with Script](/docs/guides/mongodb/initialization/using-script.md).
- Initialize [MongoDB with Snapshot](/docs/guides/mongodb/initialization/using-snapshot.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mongodb/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mongodb/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mongodb/private-registry/using-private-registry.md) to deploy MongoDB with KubeDB.
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
