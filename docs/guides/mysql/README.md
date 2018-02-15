> New to KubeDB? Please start [here](/docs/guides/README.md).

## MySQL versions supported by KubeDB

| KubeDB Version | MySQL:8.0 |
|:--:|:--:|
| 0.1.0 - 0.7.0 | &#10007; |
| 0.8.0-beta.0 | &#10003; |
| 0.8.0-beta.1 | &#10003; |

<br/>

## KubeDB Features and their availability for MySQL

|Features |Availability|
|--|:--:|
|Clustering | &#10007; |
|Persistent Volume | &#10003; |
|Instant Backup | &#10003; |
|Scheduled Backup  | &#10003; |
|Initialize using Snapshot | &#10003; |
|Initialize using Script (\*.sql, \*sql.gz and/or \*.sh) | &#10003; |
|out-of-the-box builtin-Prometheus Monitoring | &#10003; |
|out-of-the-box CoreOS-Prometheus-Operator Monitoring | &#10003; |

<br/>

## Life Cycle of MySQL in KubeDB

<p align="center">
  <img alt="lifecycle"  src="/docs/images/mysql/mysql-lifecycle.png" width="600" height="373">
</p>

## User Guide

- [Quickstart MySQL](/docs/guides/mysql/quickstart/quickstart.md) with KubeDB Operator.
- [Snapshot and Restore](/docs/guides/mysql/snapshot/backup-and-restore.md) process of MySQL databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mysql/snapshot/scheduled-backup.md) of MySQL databases using KubeDB.
- Initialize [MySQL with Script](/docs/guides/mysql/initialization/using-script.md).
- Initialize [MySQL with Snapshot](/docs/guides/mysql/initialization/using-snapshot.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
