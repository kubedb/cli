---
title: Monitor Memcached using Builtin Prometheus Discovery
menu:
  docs_0.9.0:
    identifier: mc-using-builtin-prometheus-monitoring
    name: Builtin Prometheus Discovery
    parent: mc-monitoring-memcached
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Monitoring Memcached with builtin Prometheus

This tutorial will show you how to monitor Memcached server using builtin [Prometheus](https://github.com/prometheus/prometheus) scrapper.

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Install KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- If you are not familiar with how to configure Prometheus to scrape metrics from various Kubernetes resources, please read the tutorial from [here](https://github.com/appscode/third-party-tools/tree/master/monitoring/prometheus/builtin).

- To learn how Prometheus monitoring works with KubeDB in general, please visit [here](/docs/concepts/database-monitoring/overview.md).

- To keep Prometheus resources isolated, we are going to use a separate namespace called `monitoring` to deploy respective monitoring resources. We are going to deploy database in `demo` namespace.

  ```console
  $ kubectl create ns monitoring
  namespace/monitoring created

  $ kubectl create ns demo
  namespace/demo created
  ```

> Note: YAML files used in this tutorial are stored in [docs/examples/memcached](https://github.com/kubedb/cli/tree/master/docs/examples/memcached) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Deploy Memcached server with Monitoring Enabled

At first, let's deploy a Memcached server with monitoring enabled. Below is the Memcached object that we are going to create.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  name: builtin-prom-memcd
  namespace: demo
spec:
  replicas: 1
  version: "1.5.4-v1"
  terminationPolicy: WipeOut
  podTemplate:
    spec:
      resources:
        limits:
          cpu: 500m
          memory: 128Mi
        requests:
          cpu: 250m
          memory: 64Mi
  monitor:
    agent: prometheus.io/builtin
```

Here,

- `spec.monitor.agent: prometheus.io/builtin` specifies that we are going to monitor this server using builtin Prometheus scrapper.

Let's create the Memcached crd we have shown above.

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/memcached/monitoring/builtin-prom-memcd.yaml
memcached.kubedb.com/builtin-prom-memcd created
```

Now, wait for the database to go into `Running` state.

```console
$ kubectl get mc -n demo builtin-prom-memcd
NAME                 VERSION    STATUS    AGE
builtin-prom-memcd   1.5.4-v1   Running   1m
```

KubeDB will create a separate stats service with name `{Memcached crd name}-stats` for monitoring purpose.

```console
$ kubectl get svc -n demo --selector="kubedb.com/name=builtin-prom-memcd"
NAME                       TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)     AGE
builtin-prom-memcd         ClusterIP   10.105.40.31    <none>        11211/TCP   2m6s
builtin-prom-memcd-stats   ClusterIP   10.110.89.251   <none>        56790/TCP   94s
```

Here, `builtin-prom-memcd-stats` service has been created for monitoring purpose. Let's describe the service.

```console
$ kubectl describe svc -n demo builtin-prom-memcd-stats
Name:              builtin-prom-memcd-stats
Namespace:         demo
Labels:            kubedb.com/kind=Memcached
                   kubedb.com/name=builtin-prom-memcd
Annotations:       monitoring.appscode.com/agent: prometheus.io/builtin
                   prometheus.io/path: /metrics
                   prometheus.io/port: 56790
                   prometheus.io/scrape: true
Selector:          kubedb.com/kind=Memcached,kubedb.com/name=builtin-prom-memcd
Type:              ClusterIP
IP:                10.110.89.251
Port:              prom-http  56790/TCP
TargetPort:        prom-http/TCP
Endpoints:         172.17.0.5:56790,172.17.0.7:56790,172.17.0.8:56790
Session Affinity:  None
Events:            <none>
```

You can see that the service contains following annotations.

```console
prometheus.io/path: /metrics
prometheus.io/port: 56790
prometheus.io/scrape: true
```

The Prometheus server will discover the service endpoint using these specifications and will scrape metrics from the exporter.

## Configure Prometheus Server

Now, we have to configure a Prometheus scrapping job to scrape the metrics using this service. We are going to configure scrapping job similar to this [kubernetes-service-endpoints](https://github.com/appscode/third-party-tools/tree/master/monitoring/prometheus/builtin#kubernetes-service-endpoints) job that scrapes metrics from endpoints of a service.

Let's configure a Prometheus scrapping job to collect metrics from this service.

```yaml
- job_name: 'kubedb-databases'
  kubernetes_sd_configs:
  - role: endpoints
  # by default Prometheus server select all kubernetes services as possible target.
  # relabel_config is used to filter only desired endpoints
  relabel_configs:
  # keep only those services that has "prometheus.io/scrape","prometheus.io/path" and "prometheus.io/port" anootations
  - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape, __meta_kubernetes_service_annotation_prometheus_io_port]
    separator: ;
    regex: true;(.*)
    action: keep
  # currently KubeDB supported databases uses only "http" scheme to export metrics. so, drop any service that uses "https" scheme.
  - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
    action: drop
    regex: https
  # only keep the stats services created by KubeDB for monitoring purpose which has "-stats" suffix
  - source_labels: [__meta_kubernetes_service_name]
    separator: ;
    regex: (.*-stats)
    action: keep
  # service created by KubeDB will have "kubedb.com/kind" and "kubedb.com/name" annotations. keep only those services that have these annotations.
  - source_labels: [__meta_kubernetes_service_label_kubedb_com_kind]
    separator: ;
    regex: (.*)
    action: keep
  # read the metric path from "prometheus.io/path: <path>" annotation
  - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
    action: replace
    target_label: __metrics_path__
    regex: (.+)
  # read the port from "prometheus.io/port: <port>" annotation and update scrapping address accordingly
  - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
    action: replace
    target_label: __address__
    regex: ([^:]+)(?::\d+)?;(\d+)
    replacement: $1:$2
  - action: labelmap
    regex: __meta_kubernetes_service_label_(.+)
  # add service namespace as label to the scrapped metrics
  - source_labels: [__meta_kubernetes_namespace]
    action: replace
    target_label: kubernetes_namespace
  # add service name as label to the scrapped metrics
  - source_labels: [__meta_kubernetes_service_name]
    action: replace
    target_label: kubernetes_name
```

### Configure Existing Prometheus Server

If you already have a Prometheus server running, you have to add above scrapping job in the `ConfigMap` used to configure the Prometheus server. Then, you have to restart it for the updated configuration to take effect.

>If you don't use a persistent volume for Prometheus storage, you will lose your previously scrapped data on restart.

### Deploy New Prometheus Server

If you don't have any existing Prometheus server running, you have to deploy one. In this section, we are going to deploy a Prometheus server in `monitoring` namespace to collect metrics using this stats service.

**Create ConfigMap:**

At first, create a ConfigMap with the scrapping configuration. Bellow, the YAML of ConfigMap that we are going to create in this tutorial.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
  labels:
    app: prometheus-demo
  namespace: monitoring
data:
  prometheus.yml: |-
    global:
      scrape_interval: 5s
      evaluation_interval: 5s
    scrape_configs:
    - job_name: 'kubedb-databases'
      honor_labels: true
      scheme: http
      kubernetes_sd_configs:
      - role: endpoints
      # by default Prometheus server select all kubernetes services as possible target.
      # relabel_config is used to filter only desired endpoints
      relabel_configs:
      # keep only those services that has "prometheus.io/scrape","prometheus.io/path" and "prometheus.io/port" anootations
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape, __meta_kubernetes_service_annotation_prometheus_io_port]
        separator: ;
        regex: true;(.*)
        action: keep
      # currently KubeDB supported databases uses only "http" scheme to export metrics. so, drop any service that uses "https" scheme.
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
        action: drop
        regex: https
      # only keep the stats services created by KubeDB for monitoring purpose which has "-stats" suffix
      - source_labels: [__meta_kubernetes_service_name]
        separator: ;
        regex: (.*-stats)
        action: keep
      # service created by KubeDB will have "kubedb.com/kind" and "kubedb.com/name" annotations. keep only those services that have these annotations.
      - source_labels: [__meta_kubernetes_service_label_kubedb_com_kind]
        separator: ;
        regex: (.*)
        action: keep
      # read the metric path from "prometheus.io/path: <path>" annotation
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      # read the port from "prometheus.io/port: <port>" annotation and update scrapping address accordingly
      - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
      # add service namespace as label to the scrapped metrics
      - source_labels: [__meta_kubernetes_namespace]
        separator: ;
        regex: (.*)
        target_label: namespace
        replacement: $1
        action: replace
      # add service name as a label to the scrapped metrics
      - source_labels: [__meta_kubernetes_service_name]
        separator: ;
        regex: (.*)
        target_label: service
        replacement: $1
        action: replace
      # add stats service's labels to the scrapped metrics
      - action: labelmap
        regex: __meta_kubernetes_service_label_(.+)
```

Let's create above `ConfigMap`,

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/monitoring/builtin-prometheus/prom-config.yaml
configmap/prometheus-config created
```

**Create RBAC:**

If you are using an RBAC enabled cluster, you have to give necessary RBAC permissions for Prometheus. Let's create necessary RBAC stuffs for Prometheus,

```console
$ kubectl apply -f https://raw.githubusercontent.com/appscode/third-party-tools/master/monitoring/prometheus/builtin/artifacts/rbac.yaml
clusterrole.rbac.authorization.k8s.io/prometheus created
serviceaccount/prometheus created
clusterrolebinding.rbac.authorization.k8s.io/prometheus created
```

>YAML for the RBAC resources created above can be found [here](https://github.com/appscode/third-party-tools/blob/master/monitoring/prometheus/builtin/artifacts/rbac.yaml).

**Deploy Prometheus:**

Now, we are ready to deploy Prometheus server. We are going to use following [deployment](https://github.com/appscode/third-party-tools/blob/master/monitoring/prometheus/builtin/artifacts/deployment.yaml) to deploy Prometheus server.

Let's deploy the Prometheus server.

```console
$ kubectl apply -f https://raw.githubusercontent.com/appscode/third-party-tools/master/monitoring/prometheus/builtin/artifacts/deployment.yaml
deployment.apps/prometheus created
```

### Verify Monitoring Metrics

Prometheus server is listening to port `9090`. We are going to use [port forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) to access Prometheus dashboard.

At first, let's check if the Prometheus pod is in `Running` state.

```console
$ kubectl get pod -n monitoring -l=app=prometheus
NAME                          READY   STATUS    RESTARTS   AGE
prometheus-8568c86d86-95zhn   1/1     Running   0          77s
```

Now, run following command on a separate terminal to forward 9090 port of `prometheus-8568c86d86-95zhn` pod,

```console
$ kubectl port-forward -n monitoring prometheus-8568c86d86-95zhn 9090
Forwarding from 127.0.0.1:9090 -> 9090
Forwarding from [::1]:9090 -> 9090
```

Now, we can access the dashboard at `localhost:9090`. Open [http://localhost:9090](http://localhost:9090) in your browser. You should see the endpoints of `builtin-prom-memcd-stats` service as targets.

<p align="center">
  <img alt="Prometheus Target" height="100%" src="/docs/images/memcached/monitoring/mc-builtin-prom-target.png" style="padding:10px">
</p>

Check the labels marked with red rectangle. These labels confirm that the metrics are coming from `Memcached` server `builtin-prom-memcd` through stats service `builtin-prom-memcd-stats`.

Now, you can view the collected metrics and create a graph from homepage of this Prometheus dashboard. You can also use this Prometheus server as data source for [Grafana](https://grafana.com/) and create beautiful dashboard with collected metrics.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run following commands

```console
$ kubectl delete -n demo mc/builtin-prom-memcd

$ kubectl delete -n monitoring deployment.apps/prometheus

$ kubectl delete -n monitoring clusterrole.rbac.authorization.k8s.io/prometheus
$ kubectl delete -n monitoring serviceaccount/prometheus
$ kubectl delete -n monitoring clusterrolebinding.rbac.authorization.k8s.io/prometheus

$ kubectl delete ns demo
$ kubectl delete ns monitoring
```

## Next Steps

- Monitor your Memcached server with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/memcached/monitoring/using-coreos-prometheus-operator.md).
- Use [private Docker registry](/docs/guides/memcached/private-registry/using-private-registry.md) to deploy Memcached with KubeDB.
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
