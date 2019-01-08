---
title: Scheduled Backup of Elasticsearch
menu:
  docs_0.9.0:
    identifier: es-scheduled-backup-snapshot
    name: Scheduled Backup
    parent: es-snapshot-elasticsearch
    weight: 15
menu_name: docs_0.9.0
section_menu_id: guides
---
> Don't know how backup works?  Check [tutorial](/docs/guides/elasticsearch/snapshot/instant_backup.md) on Instant Backup.

# Database Scheduled Snapshots

KubeDB supports taking periodic backups for Elasticsearch database.

## Before you begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: YAML files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create Elasticsearch with BackupSchedule

KubeDB supports taking periodic backups for a database using a [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). KubeDB operator will launch a Job periodically that takes backup and uploads the output files to various cloud providers S3, GCS, Azure, OpenStack Swift and/or locally mounted volumes using [osm](https://github.com/appscode/osm).

In this tutorial, snapshots will be stored in a Google Cloud Storage (GCS) bucket. To do so, a secret is needed that has the following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret -n demo generic gcs-secret \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "gcs-secret" created
```

To learn how to configure other storage destinations for Snapshots, please [visit here](/docs/concepts/snapshot.md).

Below is the Elasticsearch object with BackupSchedule field.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: scheduled-es
  namespace: demo
spec:
  version: "6.3-v1"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb
```

Here,

- [`cronExpression`](https://github.com/robfig/cron/blob/v2/doc.go) represents a set of times or interval when a single backup will be created.
- `storageSecretName` points to the Secret containing the credentials for snapshot storage destination.
- `gcs.bucket` points to the bucket name used to store the snapshot data

> Note: Secret object must be in the same namespace as Elasticsearch, `scheduled-es`, in this case.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/snapshot/scheduled-es.yaml
elasticsearch.kubedb.com/scheduled-es created
```

When Elasticsearch starts running successfully, KubeDB operator creates a Snapshot object immediately and registers to create a new Snapshot object on each tick of the cron expression.

```console
$ kubectl get snap -n demo --selector="kubedb.com/kind=Elasticsearch,kubedb.com/name=scheduled-es"
NAME                           DATABASENAME   STATUS    AGE
scheduled-es-20181005-120106   scheduled-es   Running   3s
```

## Update Elasticsearch to disable periodic backups

If you already have a running Elasticsearch that takes backup periodically, you can disable that by removing BackupSchedule field.

Edit your Elasticsearch object and remove BackupSchedule. This will stop taking future backups for this schedule.

```console
$ kubectl edit es -n demo scheduled-es
spec:
#  backupSchedule:
#    cronExpression: '@every 6h'
#    gcs:
#      bucket: kubedb
#    storageSecretName: gcs-secret
```

## Update Elasticsearch to enable periodic backups

If you already have a running Elasticsearch, you can enable periodic backups by adding BackupSchedule.

Edit the Elasticsearch `scheduled-es` to add following `spec.backupSchedule` section.

```yaml
$ kubectl edit es scheduled-es -n demo
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb
```

Once the `spec.backupSchedule` is added, KubeDB operator creates a Snapshot object immediately and registers to create a new Snapshot object on each tick of the cron expression.

```console
$ kubectl get snap -n demo --selector=kubedb.com/kind=Elasticsearch,kubedb.com/name=scheduled-es
NAME                           DATABASE          STATUS      AGE
scheduled-es-20180214-095019   es/scheduled-es   Succeeded   17m
scheduled-es-20180214-100711   es/scheduled-es   Running     9s
```

## Customizing `backupSchedule`

You can customize pod template spec and volume claim spec for the backup jobs by customizing `backupSchedule` section.

Some common customization sample is shown below.

**Specify PVC Template:**

Backup job needs a temporary storage to hold `dump` files before it can be uploaded to cloud backend. By default, KubeDB reads storage specification from `spec.storage` section of database crd and creates PVC with similar specification for backup job. However, if you want to specify custom PVC template, you can do it through `spec.backupSchedule.podVolumeClaimSpec` field. This is particularly helpful when you want to use different `storageclass` for backup job than the database.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: scheduled-es
  namespace: demo
spec:
  version: "6.3-v1"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb-dev
    podVolumeClaimSpec:
      storageClassName: "standard"
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi # make sure size is larger or equal than your database size
```

**Specify Resources for Backup Job:**

You can specify resources for backup job through `spec.backupSchedule.podTemplate.spec.resources` field.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: scheduled-es
  namespace: demo
spec:
  version: "6.3-v1"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb-dev
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

**Provide Annotation for Backup Job:**

If you need to add some annotations to backup job, you can specify this in `spec.backupSchedule.podTemplate.controller.annotations`. You can also specify annotation for the pod created by backup job through `spec.backupSchedule.podTemplate.annotations` field.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: scheduled-es
  namespace: demo
spec:
  version: "6.3-v1"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb-dev
    podTemplate:
      annotations:
        passMe: ToBackupJobPod
      controller:
        annotations:
          passMe: ToBackupJob
```

**Pass Arguments to Backup Job:**

KubeDB also allows to pass extra arguments for backup job. You can provide these arguments through `spec.backupSchedule.podTemplate.spec.args` field of Snapshot crd.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: scheduled-es
  namespace: demo
spec:
  version: "6.3-v1"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb-dev
    podTemplate:
      spec:
        args:
        - --extra-args-to-backup-command
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/scheduled-es -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo es/scheduled-es

$ kubectl delete ns demo
```

## Next Steps

- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Learn about initializing [Elasticsearch with Snapshot](/docs/guides/elasticsearch/initialization/snapshot_source.md).
- Learn how to configure [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
