> New to KubeDB? Please start [here](/docs/guides/README.md).

## PostgreSQL versions supported by KubeDB

| KubeDB Version | Postgres:9.5 | Postgres:9.6 | Postgres:10.2 |
|----------------|:------------:|:------------:|:-------------:|
| 0.1.0 - 0.7.0  | &#10003;     | &#10007;     | &#10007;      |
| 0.8.0-beta.0   | &#10007;     | &#10003;     | &#10007;      |
| 0.8.0-beta.1   | &#10007;     | &#10003;     | &#10003;      |

<br/>

## KubeDB features and their availability for Postgres

|Features                                               |Availability|
|-------------------------------------------------------|:----------:|
|Persistent Volume                                      | &#10003;   |
|Instant Backup                                         | &#10003;   |
|Scheduled Backup                                       | &#10003;   |
|Initialization from Snapshot                           | &#10003;   |
|Initialization using Script                            | &#10003;   |
|out-of-the-box builtin-Prometheus Monitoring           | &#10003;   |
|out-of-the-box CoreOS-Prometheus-Operator Monitoring   | &#10003;   |

## PostgreSQL features and their availability in KubeDB Postgres

|Features                                               |Availability|
|-------------------------------------------------------|:----------:|
|Clustering                                             | &#10003;   |
|Warm Standby                                           | &#10003;   |
|Hot Standby                                            | &#10003;   |
|Synchronous Replication                                | &#10007;   |
|Streaming Replication                                  | &#10003;   |
|Continuous Archiving using `wal-g`                     | &#10003;   |
|Initialization from WAL archive                        | &#10003;   |

<br/>

## Life Cycle of Postgres in KubeDB

<p align="center">
  <img alt="lifecycle"  src="/docs/images/postgres/lifecycle.png" width="581" height="362">
</p>

## User Guide

- [Quickstart PostgreSQL](/docs/guides/postgres/quickstart/quickstart.md) with KubeDB Operator.
- Take [Instant Snapshot](/docs/guides/postgres/snapshot/instant_backup.md) of PostgreSQL database using KubeDB.
- [Schedule backup](/docs/guides/postgres/snapshot/scheduled_backup.md) of PostgreSQL database using KubeDB.
- Initialize [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Initialize [PostgreSQL with KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- [PostgreSQL Clustering](/docs/guides/postgres/clustering/ha_cluster.md) supported by KubeDB Postgres.
- [Streaming Replication](/docs/guides/postgres/clustering/streaming_replication.md) for PostgreSQL clustering.
- [Continuous Archiving](/docs/guides/postgres/snapshot/continuous_archiving.md) of Write-Ahead Log by `wal-g`.
- Monitor your PostgreSQL database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/postgres/monitoring/using_builtin_prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using_coreos_prometheus_operator.md).
- Use [private Docker registry](/docs/guides/postgres/private-registry/using-private-registry.md) to deploy Postgres with KubeDB.
- Detail concepts of [Postgres object](/docs/concepts/databases/postgres.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).

