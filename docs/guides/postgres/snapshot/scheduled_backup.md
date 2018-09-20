---
title: Scheduled Backup of PostgreSQL
menu:
  docs_0.8.0:
    identifier: pg-scheduled-backup-snapshot
    name: Scheduled Backup
    parent: pg-snapshot-postgres
    weight: 15
menu_name: docs_0.8.0
section_menu_id: guides
---
> Don't know how backup works?  Check [tutorial](/docs/guides/postgres/snapshot/instant_backup.md) on Instant Backup.

# Database Scheduled Snapshots

KubeDB supports taking periodic backups for PostgreSQL database.

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

## Create Postgres with BackupSchedule

KubeDB supports taking periodic backups for a database using a [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26).
KubeDB operator will launch a Job periodically that takes backup and uploads the output files to various cloud providers S3, GCS, Azure,
OpenStack Swift and/or locally mounted volumes using [osm](https://github.com/appscode/osm).

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

Below is the Postgres object with BackupSchedule field.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: scheduled-pg
  namespace: demo
spec:
  version: "9.6"
  replicas: 3
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
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

> Note: Secret object must be in the same namespace as Postgres, `scheduled-pg`, in this case.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.0/docs/examples/postgres/snapshot/scheduled-pg.yaml
postgres "scheduled-pg" created
```

When PostgreSQL is successfully created, KubeDB operator creates a Snapshot object immediately and registers to create a new Snapshot object on each tick of the cron expression.

```console
$ kubedb get snap -n demo --selector="kubedb.com/kind=Postgres,kubedb.com/name=scheduled-pg"
NAME                           DATABASE          STATUS      AGE
scheduled-pg-20180208-105341   pg/scheduled-pg   Succeeded   32s
```

## Update Postgres to disable periodic backups

If you already have a running PostgreSQL that takes backup periodically, you can disable that by removing BackupSchedule field.

Edit your Postgres object and remove BackupSchedule. This will stop taking future backups for this schedule.

```console
$ kubedb edit pg -n demo scheduled-pg
spec:
#  backupSchedule:
#    cronExpression: '@every 6h'
#    storageSecretName: gcs-secret
#    gcs:
#      bucket: kubedb
```

## Update Postgres to enable periodic backups

If you already have a running Postgres, you can enable periodic backups by adding BackupSchedule.

Edit the Postgres `scheduled-pg` to add following `spec.backupSchedule` section.

```console
$ kubedb edit pg scheduled-pg -n demo
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb
```

Once the `spec.backupSchedule` is added, KubeDB operator creates a Snapshot object immediately and registers to create a new Snapshot object on each tick of the cron expression.

```console
$ kubedb get snap -n demo --selector="kubedb.com/kind=Postgres,kubedb.com/name=script-postgres"
NAME                              DATABASE             STATUS      AGE
instant-snapshot                  pg/script-postgres   Succeeded   30m
script-postgres-20180208-105625   pg/script-postgres   Succeeded   1m
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo pg/scheduled-pg -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo pg/scheduled-pg

$ kubectl patch -n demo drmn/scheduled-pg -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/scheduled-pg

$ kubectl delete ns demo
```

## Next Steps

- Learn about [taking instant backup](/docs/guides/postgres/snapshot/instant_backup.md) of PostgreSQL database using KubeDB Snapshot.
- Learn about initializing [PostgreSQL from KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- Setup [Continuous Archiving](/docs/guides/postgres/snapshot/continuous_archiving.md) in PostgreSQL using `wal-g`
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
