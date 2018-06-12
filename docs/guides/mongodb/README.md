---
title: MongoDB
menu:
  docs_0.8.0:
    identifier: mg-readme-mongodb
    name: MongoDB
    parent: mg-mongodb-guides
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
url: /docs/0.8.0/guides/mongodb/
aliases:
  - /docs/0.8.0/guides/mongodb/README/
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

## Supported MongoDB Features

|Features                                     | Availability |
|---------------------------------------------|:------------:|
|Clustering                                   | &#10007;     |
|Persistent Volume                            | &#10003;     |
|Instant Backup                               | &#10003;     |
|Scheduled Backup                             | &#10003;     |
|Initialize using Snapshot                    | &#10003;     |
|Initialize using Script (\*.js and/or \*.sh) | &#10003;     |
|Builtin Prometheus Discovery                 | &#10003;     |
|Using CoreOS Prometheus Operator             | &#10003;     |

<br/>

## Life Cycle of a MongoDB Object

<p align="center">
  <img alt="lifecycle"  src="/docs/images/mongodb/mgo-lifecycle.png" width="600" height="660">
</p>

<br/>

## Supported MongoDB Versions

| KubeDB Version | Mongo:3.4 | Mongo:3.6 |
|:--------------:|:---------:|:---------:|
| 0.1.0 - 0.7.0  | &#10007;  | &#10007;  |
| 0.8.0-beta.2   | &#10003;  | &#10003;  |
| 0.8.0-rc.0     | &#10003;  | &#10003;  |
| 0.8.0          | &#10003;  | &#10003;  |

## User Guide

- [Quickstart MongoDB](/docs/guides/mongodb/quickstart/quickstart.md) with KubeDB Operator.
- [Snapshot and Restore](/docs/guides/mongodb/snapshot/backup-and-restore.md) process of MongoDB databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mongodb/snapshot/scheduled-backup.md) of MongoDB databases using KubeDB.
- Initialize [MongoDB with Script](/docs/guides/mongodb/initialization/using-script.md).
- Initialize [MongoDB with Snapshot](/docs/guides/mongodb/initialization/using-snapshot.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mongodb/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mongodb/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mongodb/private-registry/using-private-registry.md) to deploy MongoDB with KubeDB.
- Use [kubedb cli](/docs/guides/mongodb/cli/cli.md) to manage databases like kubectl for Kubernetes.
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
