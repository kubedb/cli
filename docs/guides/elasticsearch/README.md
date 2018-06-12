---
title: Elasticsearch
menu:
  docs_0.8.0:
    identifier: es-readme-elasticsearch
    name: Elasticsearch
    parent: es-elasticsearch-guides
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
url: /docs/0.8.0/guides/elasticsearch/
aliases:
  - /docs/0.8.0/guides/elasticsearch/README/
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

## Supported Elasticsearch Features

|Features                                               |Availability|
|-------------------------------------------------------|:----------:|
|Clustering                                             | &#10003;   |
|Authentication (using [Search Guard](https://github.com/floragunncom/search-guard))   | &#10003;   |
|Authorization (using [Search Guard](https://github.com/floragunncom/search-guard))    | &#10003;   |
|TLS certificates (using [Search Guard](https://github.com/floragunncom/search-guard)) | &#10003;   |
|Persistent Volume                                      | &#10003;   |
|Instant Backup                                         | &#10003;   |
|Scheduled Backup                                       | &#10003;   |
|Initialization from Script                             | &#10007;   |
|Initialization from Snapshot                           | &#10003;   |
|Builtin Prometheus Discovery                           | &#10003;   |
|Using CoreOS Prometheus Operator                       | &#10003;   |

<br/>

## Life Cycle of an Elasticsearch Object

<p align="center">
  <img alt="lifecycle"  src="/docs/images/elasticsearch/lifecycle.png" width="600" height="660">
</p>

<br/>

## Supported Elasticsearch Version

| KubeDB Version | Elasticsearch:2.3 | Elasticsearch:5.6 |
|----------------|:------------:|:------------:|
| 0.1.0 - 0.7.0  | &#10003;     | &#10007;     |
| 0.8.0-beta.2   | &#10007;     | &#10003;     |
| 0.8.0   | &#10007;     | &#10003;     |

## External tools dependency

|Tool                                                               |Version   |
|-------------------------------------------------------------------|:--------:|
|[Search Guard](https://github.com/floragunncom/search-guard)       | 5.6.4-18 |
|[Elasticdump](https://github.com/taskrabbit/elasticsearch-dump/)   | 3.3.1    |
|[osm](https://github.com/appscode/osm)                             | 0.7.0    |

<br/>

## User Guide

- [Quickstart Elasticsearch](/docs/guides/elasticsearch/quickstart/quickstart.md) with KubeDB Operator.
- [Take instant backup](/docs/guides/elasticsearch/snapshot/instant_backup.md) of Elasticsearch database using KubeDB.
- [Schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Initialize [Elasticsearch with Snapshot](/docs/guides/elasticsearch/initialization/snapshot_source.md).
- [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md) supported by KubeDB
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Use [kubedb cli](/docs/guides/elasticsearch/cli/cli.md) to manage databases like kubectl for Kubernetes.
- Detail concepts of [Elasticsearch object](/docs/concepts/databases/elasticsearch.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).

