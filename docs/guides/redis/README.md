> New to KubeDB? Please start [here](/docs/guides/README.md).

## Redis versions supported by KubeDB

| KubeDB Version | Redis:4 |
|:--:|:--:|
| 0.1.0 - 0.7.0 | &#10007; |
| 0.8.0-beta.0 | &#10003; |
| 0.8.0-beta.1 | &#10003; |

<br/>

## KubeDB Features and their availability for Redis

|Features |Availability|
|--|:--:|
|Clustering | &#10007; |
|Instant Backup | &#10007; |
|Scheduled Backup  | &#10007; |
|Initialize using Snapshot | &#10007; |
|Initialize using Script | &#10007; |
|Persistent Volume | &#10003; |
|out-of-the-box builtin-Prometheus Monitoring | &#10003; |
|out-of-the-box CoreOS-Prometheus-Operator Monitoring | &#10003; |

<br/>

## Life Cycle of Redis in KubeDB

<p align="center">
  <img alt="lifecycle"  src="/docs/images/redis/redis-lifecycle.png" width="600" height="373">
</p>

## User Guide

- [Quickstart Redis](/docs/guides/redis/quickstart/quickstart.md) with KubeDB Operator.
- Monitor your Redis database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/redis/monitoring/using-coreos-prometheus-operator.md).
- Monitor your Redis database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/redis/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/redis/private-registry/using-private-registry.md) to deploy Redis with KubeDB.
- Detail concepts of [Redis object](/docs/concepts/databases/redis.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
