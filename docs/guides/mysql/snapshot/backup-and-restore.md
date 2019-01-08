---
title: Instant Backup of MySQL
menu:
  docs_0.9.0:
    identifier: my-backup-and-restore-snapshot
    name: Instant Backup
    parent: my-snapshot-mysql
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Database Snapshots

This tutorial will show you how to take snapshots of a KubeDB managed MySQL database.

> Note: The yaml files used in this tutorial are stored in [docs/examples/mysql](https://github.com/kubedb/cli/tree/master/docs/examples/mysql) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/) is required to run KubeDB. Check the available StorageClass in cluster.

  ```console
  $ kubectl get storageclasses
  NAME                 PROVISIONER                AGE
  standard (default)   k8s.io/minikube-hostpath   4h
  ```

- A `MySQL` database is needed to take snapshot for this tutorial. To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace "demo" created
  
  $ kubectl get ns
  NAME          STATUS    AGE
  demo          Active    1m
  
  $ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/snapshot/demo-1.yaml
  mysql.kubedb.com/mysql-infant created
  ```

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
secret/my-snap-secret created
```

```yaml
$ kubectl get secret my-snap-secret -n demo -o yaml
apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdX....1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInN...9tIgp9Cg==
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
    bucket: kubedb
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/snapshot/demo-2.yaml
snapshot.kubedb.com/snap-mysql-infant created

$ kubedb get snap -n demo
NAME                DATABASENAME   STATUS    AGE
snap-mysql-infant   mysql-infant   Running   13s
```

```yaml
$ kubedb get snap -n demo snap-mysql-infant -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  creationTimestamp: 2018-09-27T06:12:37Z
  finalizers:
  - kubedb.com
  generation: 1
  labels:
    kubedb.com/kind: MySQL
    kubedb.com/name: mysql-infant
  name: snap-mysql-infant
  namespace: demo
  resourceVersion: "1754"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/snapshots/snap-mysql-infant
  uid: 54efc1fe-c21c-11e8-850e-080027517bbf
spec:
  databaseName: mysql-infant
  gcs:
    bucket: kubedb
  storageSecretName: my-snap-secret
status:
  completionTime: 2018-09-27T06:18:41Z
  phase: Succeeded
  startTime: 2018-09-27T06:12:38Z
```

Here,

- `metadata.labels` should include the type of database `kubedb.com/kind: MySQL` whose snapshot will be taken.
- `spec.databaseName` points to the database whose snapshot is taken.
- `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.
- `spec.gcs.bucket` points to the bucket name used to store the snapshot data.

You can also run the `kubedb describe` command to see the recent snapshots taken for a database.

```console
$ kubedb describe my -n demo mysql-infant
Name:               mysql-infant
Namespace:          demo
CreationTimestamp:  Thu, 27 Sep 2018 12:12:10 +0600
Labels:             <none>
Annotations:        <none>
Replicas:           1  total
Status:             Running
  StorageType:      Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:          
  Name:               mysql-infant
  CreationTimestamp:  Thu, 27 Sep 2018 12:12:11 +0600
  Labels:               kubedb.com/kind=MySQL
                        kubedb.com/name=mysql-infant
  Annotations:        <none>
  Replicas:           824641842156 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         mysql-infant
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-infant
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.109.47.223
  Port:         db  3306/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.5:3306

Database Secret:
  Name:         mysql-infant-auth
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-infant
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  password:  16 bytes
  user:      4 bytes

Snapshots:
  Name               Bucket     StartTime                        CompletionTime                   Phase
  ----               ------     ---------                        --------------                   -----
  snap-mysql-infant  gs:kubedb  Thu, 27 Sep 2018 12:12:38 +0600  Thu, 27 Sep 2018 12:18:41 +0600  Succeeded

Events:
  Type    Reason              Age   From            Message
  ----    ------              ----  ----            -------
  Normal  Successful          17m   MySQL operator  Successfully created Service
  Normal  Starting            17m   Job Controller  Backup running
  Normal  Successful          14m   MySQL operator  Successfully created StatefulSet
  Normal  Successful          14m   MySQL operator  Successfully created MySQL
  Normal  Successful          14m   MySQL operator  Successfully patched StatefulSet
  Normal  Successful          14m   MySQL operator  Successfully patched MySQL
  Normal  Successful          14m   MySQL operator  Successfully patched StatefulSet
  Normal  Successful          14m   MySQL operator  Successfully patched MySQL
  Normal  SuccessfulSnapshot  11m   Job Controller  Successfully completed snapshot
```

Once the snapshot Job is complete, you should see the output of the `mysql dump` command stored in the GCS bucket.

![snapshot-console](/docs/images/mysql/m1-xyz-snapshot.png)

From the above image, you can see that the snapshot output is stored in a folder called `{bucket}/kubedb/{namespace}/{mysql-object}/{snapshot}/`.

## Restore from Snapshot

You can create a new database from a previously taken Snapshot. Specify the Snapshot name in the `spec.init.snapshotSource` field of a new MySQL object. See the example `mysql-recovered` object below:

> Note: MySQL `mysql-recovered` must have same superuser credentials as MySQL `mysql-infant`.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-recovered
  namespace: demo
spec:
  version: "8.0-v1"
  databaseSecret:
    secretName: mysql-infant-auth
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  init:
    snapshotSource:
      name: snap-mysql-infant
      namespace: demo
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/snapshot/demo-3.yaml
mysql.kubedb.com/mysql-recovered created
```

Here,

- `spec.init.snapshotSource.name` refers to a Snapshot object for a MySQL database in the same namespaces as this new `mysql-recovered` MySQL object.

Now, wait several seconds. KubeDB operator will create a new StatefulSet. Then KubeDB operator launches a Kubernetes Job to initialize the new database using the data from `snap-mysql-infant` Snapshot.

```console
$ kubedb get my -n demo
NAME              VERSION   STATUS         AGE
mysql-infant      8.0-v1    Running        27m
mysql-recovered   8.0-v1    Initializing   5m

$ kubedb get my -n demo
NAME              VERSION   STATUS    AGE
mysql-infant      8.0-v1    Running   31m
mysql-recovered   8.0-v1    Running   9m

$ kubedb describe my -n demo mysql-recovered
Name:               mysql-recovered
Namespace:          demo
CreationTimestamp:  Thu, 27 Sep 2018 12:34:07 +0600
Labels:             <none>
Annotations:        kubedb.com/initialized=
Replicas:           1  total
Status:             Running
  StorageType:      Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:          
  Name:               mysql-recovered
  CreationTimestamp:  Thu, 27 Sep 2018 12:34:09 +0600
  Labels:               kubedb.com/kind=MySQL
                        kubedb.com/name=mysql-recovered
  Annotations:        <none>
  Replicas:           824640109500 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         mysql-recovered
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-recovered
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.99.66.59
  Port:         db  3306/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.6:3306

Database Secret:
  Name:         mysql-infant-auth
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-infant
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  password:  16 bytes
  user:      4 bytes

No Snapshots.

Events:
  Type    Reason                Age   From            Message
  ----    ------                ----  ----            -------
  Normal  Successful            9m    MySQL operator  Successfully created Service
  Normal  Successful            9m    MySQL operator  Successfully created MySQL
  Normal  Successful            9m    MySQL operator  Successfully created StatefulSet
  Normal  Initializing          9m    MySQL operator  Initializing from Snapshot: "snap-mysql-infant"
  Normal  Successful            8m    MySQL operator  Successfully patched StatefulSet
  Normal  Successful            8m    MySQL operator  Successfully patched MySQL
  Normal  SuccessfulInitialize  3m    Job Controller  Successfully completed initialization
```

## Customizing Snapshot

You can customize pod template spec and volume claim spec for the backup and restore jobs. For details options read [this doc](/docs/concepts/snapshot.md).

Some common customization sample is shown below.

**Specify PVC Template:**

Backup and recovery job needs a temporary storage to hold `dump` files before it can be uploaded to cloud backend or inserted into database. By default, KubeDB reads storage specification from `spec.storage` section of database crd and creates PVC with similar specification for backup or recovery job. However, if you want to specify custom PVC template, you can do it through `spec.podVolumeClaimSpec` field of Snapshot crd. This is particularly helpful when you want to use different `storageclass` for backup or recovery job than the database.

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
    bucket: kubedb
  podVolumeClaimSpec:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi # make sure size is larger or equal than your database size
```

**Specify Resources for Backup/Recovery Job:**

You can specify resources for backup or recovery job through `spec.podTemplate.spec.resources` field.

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
    bucket: kubedb
  podTemplate:
    spec:
      resources:
        requests:
          memory: "64Mi"
          cpu: "250m"
        limits:
          memory: "128Mi"
          cpu: "500m"
```

**Provide Annotation for Backup/Recovery Job:**

If you need to add some annotations to backup or recovery job, you can specify this in `spec.podTemplate.controller.annotations`. You can also specify annotation for the pod created by backup or recovery job through `spec.podTemplate.annotations` field.

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
    bucket: kubedb
  podTemplate:
    annotations:
      passMe: ToBackupJobPod
    controller:
      annotations:
        passMe: ToBackupJob
```

**Pass Arguments to Backup/Recovery Job:**

KubeDB also allows to pass extra arguments for backup or recovery job. You can provide these arguments through `spec.podTemplate.spec.args` field of Snapshot crd.

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
    bucket: kubedb
  podTemplate:
    spec:
      args:
      - --extra-args-to-backup-command
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mysql/mysql-infant mysql/mysql-recovered -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mysql/mysql-infant mysql/mysql-recovered

kubectl patch -n demo drmn/mysql-infant drmn/mysql-recovered -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mysql-infant drmn/mysql-recovered

kubectl delete ns demo
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
- Detail concepts of [MySQLVersion object](/docs/concepts/catalog/mysql.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
