---
title: Postgres
menu:
  docs_0.8.0-beta.2:
    identifier: pg-readme-postgres
    name: Postgres
    parent: pg-postgres-guides
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
url: /docs/0.8.0-beta.2/guides/postgres/
aliases:
  - /docs/0.8.0-beta.2/guides/postgres/README/
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

## Supported PostgreSQL Features

|Features                                                | Availability |
|------------------------------------------------------- |:------------:|
|Clustering                                              | &#10003;     |
|Warm Standby                                            | &#10003;     |
|Hot Standby                                             | &#10003;     |
|Synchronous Replication                                 | &#10007;     |
|Streaming Replication                                   | &#10003;     |
|Continuous Archiving using `wal-g`                      | &#10003;     |
|Initialization from WAL archive                         | &#10003;     |
|Persistent Volume                                       | &#10003;     |
|Instant Backup                                          | &#10003;     |
|Scheduled Backup                                        | &#10003;     |
|Initialization from Snapshot                            | &#10003;     |
|Initialization using Script                             | &#10003;     |
|Builtin Prometheus Discovery                            | &#10003;     |
|Using CoreOS Prometheus Operator                        | &#10003;     |

<br/>

## Life Cycle of a PostgreSQL Object

<p align="center">
  <img alt="lifecycle"  src="/docs/images/postgres/lifecycle.png" width="581" height="362">
</p>

<br/>

## Supported PostgreSQL Versions

| KubeDB Version | Postgres:9.5 | Postgres:9.6 | Postgres:10.2 |
|----------------|:------------:|:------------:|:-------------:|
| 0.1.0 - 0.7.0  | &#10003;     | &#10007;     | &#10007;      |
| 0.8.0-beta.2   | &#10007;     | &#10003;     | &#10003;      |

<br/>

## External tools dependency

|Tool                                      |Version  |
|------------------------------------------|:-------:|
|[wal-g](https://github.com/wal-g/wal-g)   | v0.1.7  |
|[osm](https://github.com/appscode/osm)    | 0.6.2   |

## User Guide

- [Quickstart PostgreSQL](/docs/guides/postgres/quickstart/quickstart.md) with KubeDB Operator.
- Take [Instant Snapshot](/docs/guides/postgres/snapshot/instant_backup.md) of PostgreSQL database using KubeDB.
- [Schedule backup](/docs/guides/postgres/snapshot/scheduled_backup.md) of PostgreSQL database using KubeDB.
- Initialize [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Initialize [PostgreSQL with KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- [PostgreSQL Clustering](/docs/guides/postgres/clustering/ha_cluster.md) supported by KubeDB Postgres.
- [Streaming Replication](/docs/guides/postgres/clustering/streaming_replication.md) for PostgreSQL clustering.
- [Continuous Archiving](/docs/guides/postgres/snapshot/continuous_archiving.md) of Write-Ahead Log by `wal-g`.
- Monitor your PostgreSQL database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Use [private Docker registry](/docs/guides/postgres/private-registry/using-private-registry.md) to deploy PostgreSQL with KubeDB.
- Detail concepts of [Postgres object](/docs/concepts/databases/postgres.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).

