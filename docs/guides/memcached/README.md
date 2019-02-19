---
title: Memcached
menu:
  docs_0.9.0:
    identifier: mc-readme-memcached
    name: Memcached
    parent: mc-memcached-guides
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
url: /docs/0.9.0/guides/memcached/
aliases:
  - /docs/0.9.0/guides/memcached/README/
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

## Supported Memcached Features

|             Features             | Availability |
| -------------------------------- | :----------: |
| Clustering                       |   &#10007;   |
| Persistent Volume                |   &#10007;   |
| Instant Backup                   |   &#10007;   |
| Scheduled Backup                 |   &#10007;   |
| Initialize using Snapshot        |   &#10007;   |
| Initialize using Script          |   &#10007;   |
| Custom Configuration             |   &#10003;   |
| Using Custom docker image        |   &#10003;   |
| Builtin Prometheus Discovery     |   &#10003;   |
| Using CoreOS Prometheus Operator |   &#10003;   |

<br/>

## Life Cycle of a Memcached Object

<p align="center">
  <img alt="lifecycle"  src="/docs/images/memcached/memcached-lifecycle.png">
</p>

<br/>

## Supported Memcached Versions

| KubeDB Version | Memcached:1.5.4 |
| :------------: | :-------------: |
| 0.1.0 - 0.7.0  |    &#10007;     |
|     0.8.0      |    &#10003;     |
|     0.9.0      |    &#10003;     |
|     0.10.0     |    &#10003;     |

## Supported MemcachedVersion CRD

Here, &#10003; means supported and &#10007; means deprecated.

|   NAME   | VERSION | KubeDB: 0.9.0 | KubeDB: 0.10.0 |
| :------: | :-----: | :-----------: | :------------: |
|   1.5    |   1.5   |   &#10007;    |    &#10007;    |
|  1.5-v1  |   1.5   |   &#10003;    |    &#10003;    |
|  1.5.4   |  1.5.4  |   &#10007;    |    &#10007;    |
| 1.5.4-v1 |  1.5.4  |   &#10003;    |    &#10003;    |

## User Guide

- [Quickstart Memcached](/docs/guides/memcached/quickstart/quickstart.md) with KubeDB Operator.
- Monitor your Memcached server with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/memcached/monitoring/using-coreos-prometheus-operator.md).
- Monitor your Memcached server with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/memcached/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/memcached/private-registry/using-private-registry.md) to deploy Memcached with KubeDB.
- Use [kubedb cli](/docs/guides/memcached/cli/cli.md) to manage databases like kubectl for Kubernetes.
- Detail concepts of [Memcached object](/docs/concepts/databases/memcached.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
