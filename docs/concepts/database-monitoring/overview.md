---
title: Database Monitoring Overview
description: Database Monitoring Overview
menu:
  docs_0.11.0:
    identifier: database-monitoring-overview
    name: Overview
    parent: database-monitoring
    weight: 10
menu_name: docs_0.11.0
section_menu_id: concepts
---

# Monitoring Database with KubeDB

KubeDB has native support for monitoring via [Prometheus](https://prometheus.io/). You can use builtin [Prometheus](https://github.com/prometheus/prometheus) scrapper or [CoreOS Prometheus Operator](https://github.com/coreos/prometheus-operator) to monitor KubeDB managed databases. This tutorial will show you how database monitoring works with KubeDB and how to configure Database crd to enable monitoring.

## Overview

KubeDB uses Prometheus [exporter](https://prometheus.io/docs/instrumenting/exporters/#databases) images to export Prometheus metrics for respective databases. Following diagram shows the logical flow of database monitoring with KubeDB.

<p align="center">
  <img alt="Database Monitoring Flow"  src="/docs/images/concepts/monitoring/database-monitoring-overview.svg">
</p>

When a user creates a database crd with `spec.monitor` section configured, KubeDB operator provisions the respective database and injects an exporter image as sidecar to the database pod. It also creates a dedicated stats service with name `{database-crd-name}-stats` for monitoring. Prometheus server can scrape metrics using this stats service.

## Configure Monitoring

In order to enable monitoring for a database, you have to configure `spec.monitor` section. KubeDB provides following options to configure `spec.monitor` section:

|                Field                |    Type    |                                                                                     Uses                                                                                      |
| ----------------------------------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `spec.monitor.agent`                | `Required` | Type of the monitoring agent that will be used to monitor this database. It can be `prometheus.io/builtin` or `prometheus.io/coreos-operator`.                              |
| `spec.monitor.prometheus.namespace` | `Optional` | Namespace where the Prometheus server is running or will be deployed. For agent type `prometheus.io/coreos-operator`, `ServiceMonitor` crd will be created in this namespace. |
| `spec.monitor.prometheus.labels`    | `Optional` | Labels for `ServiceMonitor`  crd.                                                                                                                                             |
| `spec.monitor.prometheus.port`      | `Optional` | Port number where the exporter side car will serve metrics.                                                                                                                   |
| `spec.monitor.prometheus.interval`  | `Optional` | Interval at which metrics should be scraped.                                                                                                                                  |
| `spec.monitor.args`                 | `Optional` | Arguments to pass to the exporter sidecar.                                                                                                                                    |
| `spec.monitor.env`                  | `Optional` | List of environment variables to set in the exporter sidecar container.                                                                                                       |
| `spec.monitor.resources`            | `Optional` | Resources required by exporter sidecar container.                                                                                                                             |
| `spec.monitor.securityContext`      | `Optional` | Security options the exporter should run with.                                                                                                                                |

## Sample Configuration

A sample YAML for Redis crd with `spec.monitor` section configured to enable monitoring with CoreOS [prometheus-operator](https://github.com/coreos/prometheus-operator) is shown below.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  name: sample-redis
  namespace: databases
spec:
  version: "4.0-v1"
  terminationPolicy: WipeOut
  configSource: # configure Redis to use password for authentication
    configMap:
      name: redis-config
  storageType: Durable
  storage:
    storageClassName: default
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 5Gi
  monitor:
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: monitoring
      labels:
        k8s-app: prometheus
    args:
    - --redis.password=$(REDIS_PASSWORD)
    env:
    - name: REDIS_PASSWORD
      valueFrom:
        secretKeyRef:
          name: _name_of_secret_with_redis_password
          key: password # key with the password
    resources:
      requests:
        memory: 512Mi
        cpu: 200m
      limits:
        memory: 512Mi
        cpu: 250m
    securityContext:
      runAsUser: 2000
      allowPrivilegeEscalation: false
```

Assume that above Redis is configured to use basic authentication. So, exporter image also need to provide password to collect metrics. We have provided it through `spec.monitor.args` field.

Here, we have specified that we are going to monitor this server using CoreOS prometheus-operator through `spec.monitor.agent: prometheus.io/coreos-operator`. KubeDB will create a `ServiceMonitor` crd in `monitoring` namespace and this `ServiceMonitor` will have `k8s-app: prometheus` label.

## Next Steps

- Learn how to monitor Elasticsearch database with KubeDB using [builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md) and using [CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Learn how to monitor PostgreSQL database with KubeDB using [builtin-Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md) and using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Learn how to monitor MySQL database with KubeDB using [builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md) and using [CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Learn how to monitor MongoDB database with KubeDB using [builtin-Prometheus](/docs/guides/mongodb/monitoring/using-builtin-prometheus.md) and using [CoreOS Prometheus Operator](/docs/guides/mongodb/monitoring/using-coreos-prometheus-operator.md).
- Learn how to monitor Redis server with KubeDB using [builtin-Prometheus](/docs/guides/redis/monitoring/using-builtin-prometheus.md) and using [CoreOS Prometheus Operator](/docs/guides/redis/monitoring/using-coreos-prometheus-operator.md).
- Learn how to monitor Memcached server with KubeDB using [builtin-Prometheus](/docs/guides/memcached/monitoring/using-builtin-prometheus.md) and using [CoreOS Prometheus Operator](/docs/guides/memcached/monitoring/using-coreos-prometheus-operator.md).
