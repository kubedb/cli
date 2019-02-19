---
title: MySQL
menu:
  docs_0.9.0:
    identifier: my-readme-mysql
    name: MySQL
    parent: my-mysql-guides
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
url: /docs/0.9.0/guides/mysql/
aliases:
  - /docs/0.9.0/guides/mysql/README/
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

## Supported MySQL Features

|                        Features                         | Availability |
| ------------------------------------------------------- | :----------: |
| Clustering                                              |   &#10007;   |
| Persistent Volume                                       |   &#10003;   |
| Instant Backup                                          |   &#10003;   |
| Scheduled Backup                                        |   &#10003;   |
| Initialize using Snapshot                               |   &#10003;   |
| Initialize using Script (\*.sql, \*sql.gz and/or \*.sh) |   &#10003;   |
| Custom Configuration                                    |   &#10003;   |
| Using Custom docker image                               |   &#10003;   |
| Builtin Prometheus Discovery                            |   &#10003;   |
| Using CoreOS Prometheus Operator                        |   &#10003;   |

<br/>

## Life Cycle of a MySQL Object

<p align="center">
  <img alt="lifecycle"  src="/docs/images/mysql/mysql-lifecycle.png" >
</p>

<br/>

## Supported MySQL Versions

| KubeDB Version | MySQL:8.0 | MySQL:5.7 |
| :------------: | :-------: | :-------: |
| 0.1.0 - 0.7.0  | &#10007;  | &#10007;  |
|     0.8.0      | &#10003;  | &#10003;  |
|     0.9.0      | &#10003;  | &#10003;  |
|     0.10.0     | &#10003;  | &#10003;  |

## Supported MySQLVersion CRD

Here, &#10003; means supported and &#10007; means deprecated.

|  NAME  | VERSION | KubeDB: 0.9.0 | KubeDB: 0.9.0 |
| ------ | ------- | ------------- | ------------- |
| 5      | 5       | &#10007;      | &#10007;      |
| 5.7    | 5.7     | &#10007;      | &#10007;      |
| 8      | 8       | &#10007;      | &#10007;      |
| 8.0    | 8.0     | &#10007;      | &#10007;      |
| 5-v1   | 5       | &#10003;      | &#10007;      |
| 5.7-v1 | 5.7     | &#10003;      | &#10003;      |
| 8-v1   | 8       | &#10003;      | &#10007;      |
| 8.0-v1 | 8.0.3   | &#10003;      | &#10007;      |
| 8.0-v2 | 8.0.14  | &#10007;      | &#10003;      |
| 8.0.3  | 8.0.3   | &#10007;      | &#10003;      |
| 8.0-v2 | 8.0.14  | &#10007;      | &#10003;      |
| 8.0.14 | 8.0.14  | &#10007;      | &#10003;      |

## External tools dependency

|                  Tool                  | Version |
| -------------------------------------- | :-----: |
| [osm](https://github.com/appscode/osm) |  0.9.1  |

<br/>

## User Guide

- [Quickstart MySQL](/docs/guides/mysql/quickstart/quickstart.md) with KubeDB Operator.
- [Snapshot and Restore](/docs/guides/mysql/snapshot/backup-and-restore.md) process of MySQL databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mysql/snapshot/scheduled-backup.md) of MySQL databases using KubeDB.
- Initialize [MySQL with Script](/docs/guides/mysql/initialization/using-script.md).
- Initialize [MySQL with Snapshot](/docs/guides/mysql/initialization/using-snapshot.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Use [kubedb cli](/docs/guides/mysql/cli/cli.md) to manage databases like kubectl for Kubernetes.
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Detail concepts of [MySQLVersion object](/docs/concepts/catalog/mysql.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
