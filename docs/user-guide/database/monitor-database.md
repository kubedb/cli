### Monitor Database

We now support only Promethues to minitor database.

**T**o set monitoring by Promethues, we need to set following in `spec.monitor.prometheus`

* `namespace:` Namespace (ServiceMonitor will be created in this namespace)
* `labels:` Operator will detect ServiceMonitor using this labels.
* `interval:` Interval at which metrics should be scraped

```yaml
spec:
  monitor:
    prometheus:
      namespace: default
      labels:
        app: kubedb-exporter
      interval: 10s
```

**A**s we have used monitoring information in our database yaml,
Prometheus will start collecting matrices fot this database from exporter.
