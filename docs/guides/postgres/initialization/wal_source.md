---
title: Initialize Postgres from WAL
menu:
  docs_0.11.0:
    identifier: pg-wal-source-initialization
    name: From WAL
    parent: pg-initialization-postgres
    weight: 20
menu_name: docs_0.11.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).
> Don't know how to take continuous backup?  Check this [tutorial](/docs/guides/postgres/snapshot/continuous_archiving.md) on Continuous Archiving.

# PostgreSQL Initialization from WAL files

KubeDB supports PostgreSQL database initialization. When you create a new Postgres object, you can provide existing WAL files to restore from by "replaying" the log entries. Users can now restore from any one of `s3`, `gcs`, `azure`, or `swift` as cloud backup provider.

**What is Continuous Archiving**

PostgreSQL maintains a write ahead log (WAL) in the `pg_xlog/` subdirectory of the cluster's data directory.  The existence of the log makes it possible to restore from the backed-up WAL files to bring the system back to a last known state.

To know more about continuous archiving, please refer to the [ofiicial postgres document](https://www.postgresql.org/docs/10/continuous-archiving.html) on this topic.

**List of supported Cloud Providers for PostgresVersion CRDs**

|   Name   | Version |  S3  | GCS  | Azure | Swift |
| :------: | :-----: | :--: | :--: | :---: | :---: |
|  9.6-v2  |   9.6   |  ✓   |  ✓   |   ✗   |   ✗   |
| 9.6.7-v2 |  9.6.7  |  ✓   |  ✓   |   ✗   |   ✗   |
| 10.2-v2  |  10.2   |  ✓   |  ✓   |   ✗   |   ✗   |
|   10.6   |  10.6   |  ✓   |  ✓   |   ✗   |   ✗   |
|   11.1   |  11.1   |  ✓   |  ✓   |   ✗   |   ✗   |
|  9.6-v3  |   9.6   |  ✓   |  ✓   |   ✓   |   ✓   |
| 9.6.7-v3 |  9.6.7  |  ✓   |  ✓   |   ✓   |   ✓   |
| 10.2-v3  |  10.2   |  ✓   |  ✓   |   ✓   |   ✓   |
| 10.6-v1  |  10.6   |  ✓   |  ✓   |   ✓   |   ✓   |
| 11.1-v1  |  11.1   |  ✓   |  ✓   |   ✓   |   ✓   |

## Next Steps

- Learn about restoring from [Amazon S3](/docs/guides/postgres/initialization/replay_from_s3.md).
- Learn about restoring from [Google Cloud Storage](/docs/guides/postgres/initialization/replay_from_gcs.md).
- Learn about restoring from [Azure Storage](/docs/guides/postgres/initialization/replay_from_azure.md).
- Learn about restoring from [OpenStack Object Storage (Swift)](/docs/guides/postgres/initialization/replay_from_swift.md).
- Learn about initializing [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
