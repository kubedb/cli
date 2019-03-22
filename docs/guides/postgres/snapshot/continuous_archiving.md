---
title: Continuous Archiving of PostgreSQL
menu:
  docs_0.11.0:
    identifier: pg-continuous-archiving-snapshot
    name: WAL Archiving
    parent: pg-snapshot-postgres
    weight: 20
menu_name: docs_0.11.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Continuous Archiving with WAL-G

KubeDB also supports continuous archiving of PostgreSQL using [WAL-G ](https://github.com/wal-g/wal-g). Users can now use any one of `s3`, `gcs`, `azure`, or `swift` as cloud storage destination. 

**What is this Continuous Archiving**

PostgreSQL maintains a write ahead log (WAL) in the `pg_xlog/` subdirectory of the cluster's data directory.  The existence of the log makes it possible to use a third strategy for backing up databases and if recovery is needed, restore from the backed-up WAL files to bring the system back to last known state.

To know more about continuous archiving, please refer to the [ofiicial postgres document](https://www.postgresql.org/docs/10/continuous-archiving.html) on this topic.

**Continuous Archiving Setup**

Following additional parameters are set in `postgresql.conf` for *primary* server

```console
archive_command = 'wal-g wal-push %p'
archive_timeout = 60
```

Here, these commands are used to push files to the cloud.



**List of supported Cloud Destination for PostgresVersion CRDs**

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

Users can use supported cloud destinations to backup  WAL files to restore whaen needed. KubeDB currently supports *archiving to* [Amazon S3](/docs/guides/postgres/snapshot/archiving_to_s3.md), [Google Cloud Storage](/docs/guides/postgres/snapshot/archiving_to_gcs.md), [Azure Storage](/docs/guides/postgres/snapshot/archiving_to_azure.md), and [OpenStack Object Storage (Swift)](/docs/guides/postgres/snapshot/archiving_to_swift.md).

## Next Steps

- Learn about archiving to [Amazon S3](/docs/guides/postgres/snapshot/archiving_to_s3.md).
- Learn about archiving to [Google Cloud Storage](/docs/guides/postgres/snapshot/archiving_to_gcs.md).
- Learn about archiving to [Azure Storage](/docs/guides/postgres/snapshot/archiving_to_azure.md).
- Learn about archiving to [OpenStack Object Storage (Swift)](/docs/guides/postgres/snapshot/archiving_to_swift.md).
