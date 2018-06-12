---
title: Redis
menu:
  docs_0.8.0:
    identifier: rd-readme-redis
    name: Redis
    parent: rd-redis-guides
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
url: /docs/0.8.0/guides/redis/
aliases:
  - /docs/0.8.0/guides/redis/README/
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

## Supported Redis Features

|Features                                                | Availability |
|--------------------------------------------------------|:------------:|
|Clustering                                              | &#10007;     |
|Instant Backup                                          | &#10007;     |
|Scheduled Backup                                        | &#10007;     |
|Initialize using Snapshot                               | &#10007;     |
|Initialize using Script                                 | &#10007;     |
|Persistent Volume                                       | &#10003;     |
|Builtin Prometheus Discovery                            | &#10003;     |
|Using CoreOS Prometheus Operator                        | &#10003;     |

<br/>

## Supported Redis Versions

| KubeDB Version | Redis:4  |
|:--------------:|:--------:|
| 0.1.0 - 0.7.0  | &#10007; |
| 0.8.0-beta.2   | &#10003; |
| 0.8.0   | &#10003; |

<br/>

## Life Cycle of a Redis Object

<p align="center">
  <img alt="lifecycle"  src="/docs/images/redis/redis-lifecycle.png" width="600" height="660">
</p>

## User Guide

- [Quickstart Redis](/docs/guides/redis/quickstart/quickstart.md) with KubeDB Operator.
- Monitor your Redis database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/redis/monitoring/using-coreos-prometheus-operator.md).
- Monitor your Redis database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/redis/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/redis/private-registry/using-private-registry.md) to deploy Redis with KubeDB.
- Use [kubedb cli](/docs/guides/redis/cli/cli.md) to manage databases like kubectl for Kubernetes.
- Detail concepts of [Redis object](/docs/concepts/databases/redis.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
