> New to KubeDB Postgres?  Quick start [here](/docs/guides/postgres/quickstart.md).

> Don't know how to take continuous backup?  Check [tutorial](/docs/guides/postgres/snapshot/continuous_archiving.md) on Continuous Archiving.

# Postgres Initialization

KubeDB supports PostgreSQL database initialization. When you create a new Postgres object, you can provide existing WAL files to restore from by "replaying" the log entries.

## PostgreSQL WAL Source

You can create a new database from archived WAL files using [wal-g ](https://github.com/wal-g/wal-g).

Specify storage backend in the `spec.init.postgresWAL` field of a new Postgres object. Add following additional information in `spec` of a new Postgres:

See the example Postgres object below

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: replay-postgres
  namespace: demo
spec:
  version: 9.6
  replicas: 2
  databaseSecret:
    secretName: wal-postgres-auth
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  archiver:
    storage:
      storageSecretName: s3-secret
      s3:
        bucket: kubedb
  init:
    postgresWAL:
      storageSecretName: s3-secret
      s3:
        bucket: kubedb
        prefix: 'kubedb/demo/wal-postgres/archive'
```

Here,

- `spec.init.postgresWAL` specifies storage information that will be used by `wal-g`
	- `storageSecretName` points to the Secret containing the credentials for cloud storage destination.
	- `s3.bucket` points to the bucket name used to store continuous archiving data.
	- `s3.prefix` points to the path where archived WAL data is stored. Prefix format: `/kubedb/{namespace}/{postgres-name}/archive/`.
	Here `{namespace}` & `{postgres-name}` indicates Postgres object from where this archived WAL data is generated.

> Note: Postgres `replay-postgres` must have same `postgres` superuser password as Postgres `wal-postgres`.

[//]: # (Describe authentication part. This should match with existing one)

Now create this Postgres

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/postgres/initialization/replay-postgres.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/postgres/initialization/replay-postgres.yaml"
postgres "replay-postgres" created
```
This will create a new database with existing _basebackup_ and will restore from archived _wal_ files.

When this database is ready, **wal-g** takes a _basebackup_ and uploads it to cloud storage defined by storage backend in `spec.archiver`.

## Next Steps
- Learn about initializing [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using_builtin_prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using_coreos_prometheus_operator.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
