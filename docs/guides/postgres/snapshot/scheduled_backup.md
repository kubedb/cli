> Don't know how backup works?  Check [tutorial](/docs/guides/postgres/snapshot/instant_backup.md) on Instant Backup.

## Scheduled Backup

KubeDB supports taking periodic backups for PostgreSQL database.

To enable this, you need to add BackupSchedule in Postgres `spec`.

### Create Postgres with BackupSchedule

Below is the Postgres object with BackupSchedule field.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: scheduled-pg
  namespace: demo
spec:
  version: 9.6
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
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/postgres/snapshot/scheduled-pg.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/postgres/snapshot/scheduled-pg.yaml"
postgres "scheduled-pg" created
```

When Postgres is successfully created, KubeDB operator creates a Snapshot object immediately and registers to create a new Snapshot object on each tick of the cron expression.

```console
$ kubedb get snap -n demo --selector="kubedb.com/kind=Postgres,kubedb.com/name=scheduled-pg"
NAME                           DATABASE          STATUS      AGE
scheduled-pg-20180208-105341   pg/scheduled-pg   Succeeded   32s
```

### Update Postgres to enable periodic backups

If you already have a running Postgres, you can enable periodic backups by adding BackupSchedule.

Edit the Postgres `script-postgres` to add following `spec.backupSchedule` section.

```yaml
$ kubedb edit pg script-postgres -n demo
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

### Update Postgres to disable periodic backups

If you already have a running Postgres that takes backup periodically, you can disable that by removing BackupSchedule field.
Edit your Postgres object and remove BackupSchedule. This will stop taking future backups for this schedule.

## Next Steps
- Learn about [taking instant backup](/docs/guides/postgres/snapshot/instant_backup.md) of PostgreSQL database using KubeDB Snapshot.
- Learn about initializing [PostgreSQL from KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- Setup [Continuous Archiving](/docs/guides/postgres/snapshot/continuous_archiving.md) in PostgreSQL using `wal-g`
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
