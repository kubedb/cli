---
title: Monitor Elasticsearch using Builtin Prometheus Discovery
menu:
  docs_0.8.0-beta.2:
    identifier: es-using-builtin-prometheus-monitoring
    name: Builtin Prometheus Discovery
    parent: es-monitoring-elasticsearch
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Prometheus with KubeDB

This tutorial will show you how to monitor Elasticsearch database using [Prometheus](https://prometheus.io/).

## Before You begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: Yaml files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

This tutorial assumes that you are familiar with Elasticsearch concept.

## Monitor with builtin Prometheus

Below is the Elasticsearch object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: builtin-prom-es
  namespace: demo
spec:
  version: "5.6"
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  monitor:
    agent: prometheus.io/builtin
```

Here,

- `spec.monitor` specifies that built-in [prometheus](https://github.com/prometheus/prometheus) is used to monitor this database instance.

Run following command to create example above.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/elasticsearch/monitoring/builtin-prom-es.yaml
elasticsearch "builtin-prom-es" created
```

KubeDB operator will configure its service once the Elasticsearch is successfully running.

```console
$ kubedb get es -n demo builtin-prom-es
NAME              STATUS    AGE
builtin-prom-es   Running   5m
```

You can verify it running the following commands:

```console
$ kubectl get svc -n demo --selector="kubedb.com/name=builtin-prom-es"
NAME                     TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)              AGE
builtin-prom-es          ClusterIP   10.101.27.83     <none>        9200/TCP,56790/TCP   41s
builtin-prom-es-master   ClusterIP   10.111.170.101   <none>        9300/TCP             41s
```

Lets describe Service `builtin-prom-es`

```console
$ kubectl describe svc -n demo builtin-prom-es
Name:              builtin-prom-es
Namespace:         demo
Labels:            kubedb.com/kind=Elasticsearch
                   kubedb.com/name=builtin-prom-es
Annotations:       monitoring.appscode.com/agent=prometheus.io/builtin
                   prometheus.io/path=/kubedb.com/v1alpha1/namespaces/demo/elasticsearches/builtin-prom-es/metrics
                   prometheus.io/port=56790
                   prometheus.io/scrape=true
Selector:          kubedb.com/kind=Elasticsearch,kubedb.com/name=builtin-prom-es,node.role.client=set
Type:              ClusterIP
IP:                10.101.27.83
Port:              http  9200/TCP
TargetPort:        %!d(string=http)/TCP
Endpoints:         172.17.0.8:9200
Port:              prom-http  56790/TCP
TargetPort:        %!d(string=prom-http)/TCP
Endpoints:         172.17.0.8:56790
Session Affinity:  None
```

You can see that the service contains following annotations.

```console
prometheus.io/path=/kubedb.com/v1alpha1/namespaces/demo/elasticsearches/builtin-prom-es/metrics
prometheus.io/port=56790
prometheus.io/scrape=true
```

The prometheus server will discover the service endpoint aka `Elasticsearch Exporter` using these specifications and will scrap metrics from exporter.

## Deploy and configure Prometheus server

The prometheus server is needed to configure so that it can discover endpoints of services. If a Prometheus server is already running in cluster
and if it is configured in a way that it can discover service endpoints, no extra configuration will be needed.

If there is no existing Prometheus server running, rest of this tutorial will create a Prometheus server with appropriate configuration.

The configuration file of Prometheus server will be provided by ConfigMap. Create following ConfigMap with Prometheus configuration.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-server-conf
  labels:
    name: prometheus-server-conf
  namespace: demo
data:
  prometheus.yml: |-
    global:
      scrape_interval: 5s
      evaluation_interval: 5s
    scrape_configs:
    - job_name: 'kubernetes-service-endpoints'

      kubernetes_sd_configs:
      - role: endpoints

      relabel_configs:
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
        action: replace
        target_label: __scheme__
        regex: (https?)
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
      - action: labelmap
        regex: __meta_kubernetes_service_label_(.+)
      - source_labels: [__meta_kubernetes_namespace]
        action: replace
        target_label: kubernetes_namespace
      - source_labels: [__meta_kubernetes_service_name]
        action: replace
        target_label: kubernetes_name
```

Create above ConfigMap

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/monitoring/builtin-prometheus/demo-1.yaml
configmap "prometheus-server-conf" created
```

Now, the below YAML is used to deploy Prometheus in kubernetes :

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus-server
  namespace: demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus-server
  template:
    metadata:
      labels:
        app: prometheus-server
    spec:
      containers:
        - name: prometheus
          image: prom/prometheus:v2.1.0
          args:
            - "--config.file=/etc/prometheus/prometheus.yml"
            - "--storage.tsdb.path=/prometheus/"
          ports:
            - containerPort: 9090
          volumeMounts:
            - name: prometheus-config-volume
              mountPath: /etc/prometheus/
            - name: prometheus-storage-volume
              mountPath: /prometheus/
      volumes:
        - name: prometheus-config-volume
          configMap:
            defaultMode: 420
            name: prometheus-server-conf
        - name: prometheus-storage-volume
          emptyDir: {}
```

#### In RBAC enabled cluster

If RBAC *is* enabled, run the following command to prepare your cluster for this tutorial

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/monitoring/builtin-prometheus/rbac/demo-2.yaml
clusterrole "prometheus-server" created
serviceaccount "prometheus-server" created
clusterrolebinding "prometheus-server" created
deployment "prometheus-server" created
service "prometheus-service" created
```

Watch the Deployment’s Pods.

```console
$ kubectl get pods -n demo --selector=app=prometheus-server --watch
NAME                                 READY     STATUS              RESTARTS   AGE
prometheus-server-6b8476d6c5-kx78z   0/1       ContainerCreating   0          1m
prometheus-server-6b8476d6c5-kx78z   1/1       Running   0         1m
```

And also verify RBAC stuffs

```console
$ kubectl get clusterrole prometheus-server -n demo
NAME                AGE
prometheus-server   1m
```

```console
$ kubectl get clusterrolebinding prometheus-server -n demo
NAME                AGE
prometheus-server   2m
```

#### In RBAC \*not\* enabled cluster

If RBAC *is not* enabled, Run the following command to deploy prometheus in kubernetes

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/monitoring/builtin-prometheus/demo-2.yaml
deployment "prometheus-server" created
service "prometheus-service" created
```

Watch the Deployment’s Pods.

```console
$ kubectl get pods -n demo --selector=app=prometheus-server --watch
NAME                                 READY     STATUS              RESTARTS   AGE
prometheus-server-6b8476d6c5-kx78z   0/1       ContainerCreating   0          1m
prometheus-server-6b8476d6c5-kx78z   1/1       Running   0         1m
```

### Prometheus Dashboard

Now open prometheus dashboard on browser by running `minikube service prometheus-service -n demo`.

Or you can get the URL of `prometheus-service` Service by running following command

```console
$ minikube service prometheus-service -n demo --url
http://192.168.99.100:30901
```

Now, if you go to the Prometheus Dashboard, you will see this database endpoint in target list.

<p align="center">
  <kbd>
    <img alt="builtin-prom-elasticsearch"  src="/docs/images/elasticsearch/builtin-prom-es.png">
  </kbd>
</p>

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run following commands

```console
$ kubectl patch -n demo es/builtin-prom-es -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo es/builtin-prom-es

$ kubectl patch -n demo drmn/builtin-prom-es -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/builtin-prom-es

# In rbac enabled cluster,
# $ kubectl delete clusterrole prometheus-server
# $ kubectl delete clusterrolebindings  prometheus-server
# $ kubectl delete serviceaccounts -n demo  prometheus-server

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Learn about [taking instant backup](/docs/guides/elasticsearch/snapshot/instant_backup.md) of Elasticsearch database using KubeDB.
- Learn how to [schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Learn about initializing [Elasticsearch with Snapshot](/docs/guides/elasticsearch/initialization/snapshot_source.md).
- Learn how to configure [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
