> New to KubeDB? Please start [here](/docs/guides/README.md).

## Elasticsearch versions supported by KubeDB

| KubeDB Version | Elasticsearch:2.3 | Elasticsearch:5.6 |
|----------------|:------------:|:------------:|
| 0.1.0 - 0.7.0  | &#10003;     | &#10007;     |
| 0.8.0-beta.0   | &#10007;     | &#10003;     |
| 0.8.0-beta.1   | &#10007;     | &#10003;     |

<br/>

## KubeDB features and their availability for Elasticsearch

|Features                                               |Availability|
|-------------------------------------------------------|:----------:|
|Persistent Volume                                      | &#10003;   |
|Instant Backup                                         | &#10003;   |
|Scheduled Backup                                       | &#10003;   |
|Initialization from Snapshot                           | &#10003;   |
|out-of-the-box builtin-Prometheus Monitoring           | &#10003;   |
|out-of-the-box CoreOS-Prometheus-Operator Monitoring   | &#10003;   |

## Elasticsearch features and their availability in KubeDB Elasticsearch

|Features                                               |Availability|
|-------------------------------------------------------|:----------:|
|Clustering                                             | &#10003;   |
|Authentication                                         | &#10003;   |
|Authorization                                          | &#10003;   |
|TLS certificates                                       | &#10003;   |

<br/>

## External tools dependency

|Tool                                                               |Version   |
|-------------------------------------------------------------------|:--------:|
|[searchguard](https://github.com/floragunncom/search-guard)        | 5.6.4-18 |
|[elasticdump](https://github.com/taskrabbit/elasticsearch-dump/)   | 3.3.1    |
|[osm](https://github.com/appscode/osm)                             | 0.6.2    |

<br/>

## Life Cycle of Elasticsearch in KubeDB

<p align="center">
  <img alt="lifecycle"  src="/docs/images/elasticsearch/lifecycle.png" width="581" height="362">
</p>

## User Guide

- [Quickstart Elasticsearch](/docs/guides/elasticsearch/quickstart/quickstart.md) with KubeDB Operator.
- [Take instant backup](/docs/guides/elasticsearch/snapshot/instant_backup.md) of Elasticsearch database using KubeDB.
- [Schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Initialize [Elasticsearch with Snapshot](/docs/guides/elasticsearch/initialization/snapshot_source.md).
- [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md) supported by KubeDB
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using_builtin_prometheus.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using_coreos_prometheus_operator.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Detail concepts of [Elasticsearch object](/docs/concepts/databases/elasticsearch.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).

