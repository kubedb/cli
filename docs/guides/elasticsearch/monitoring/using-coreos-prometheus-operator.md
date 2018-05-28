---
title: Monitoring Elasticsearch using Coreos Prometheus Operator
menu:
  docs_0.8.0-rc.0:
    identifier: es-using-coreos-prometheus-operator-monitoring
    name: Coreos Prometheus Operator
    parent: es-monitoring-elasticsearch
    weight: 10
menu_name: docs_0.8.0-rc.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Prometheus (CoreOS operator) with KubeDB

This tutorial will show you how to monitor Elasticsearch database using Prometheus via [CoreOS Prometheus Operator](https://github.com/coreos/prometheus-operator).

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

## Deploy CoreOS-Prometheus Operator

#### In RBAC enabled cluster

If RBAC *is* enabled, run the following command to prepare your cluster for this tutorial

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-rc.0/docs/examples/monitoring/coreos-operator/rbac/demo-0.yaml
clusterrole "prometheus-operator" created
serviceaccount "prometheus-operator" created
clusterrolebinding "prometheus-operator" created
deployment "prometheus-operator" created
```

Watch the Deployment’s Pods.

```console
$ kubectl get pods -n demo --selector=operator=prometheus --watch
NAME                                   READY     STATUS    RESTARTS   AGE
prometheus-operator-79cb9dcd4b-24khh   1/1       Running   0          46s
```

This CoreOS-Prometheus operator will create some supported Custom Resource Definition (CRD).

```console
$ kubectl get crd
NAME                                    AGE
alertmanagers.monitoring.coreos.com     3m
prometheuses.monitoring.coreos.com      3m
servicemonitors.monitoring.coreos.com   3m
```

Once the Prometheus CRDs are registered, run the following command to create a Prometheus.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-rc.0/docs/examples/monitoring/coreos-operator/rbac/demo-1.yaml
clusterrole "prometheus" created
serviceaccount "prometheus" created
clusterrolebinding "prometheus" created
prometheus "prometheus" created
service "prometheus" created
```

Verify RBAC stuffs

```console
$ kubectl get clusterroles
NAME                  AGE
prometheus            1m
prometheus-operator   5m
```

```console
$ kubectl get clusterrolebindings
NAME                  AGE
prometheus            1m
prometheus-operator   5m
```

#### In RBAC *not* enabled cluster

If RBAC *is not* enabled, Run the following command to prepare your cluster for this tutorial:

```console
$ https://raw.githubusercontent.com/kubedb/cli/0.8.0-rc.0/docs/examples/monitoring/coreos-operator/demo-0.yaml
deployment "prometheus-operator" created
```

Watch the Deployment’s Pods.

```console
$ kubectl get pods -n demo --selector=operator=prometheus --watch
NAME                                   READY     STATUS    RESTARTS   AGE
prometheus-operator-79cb9dcd4b-24khh   1/1       Running   0          46s
```

This CoreOS-Prometheus operator will create some supported Custom Resource Definition (CRD).

```console
$ kubectl get crd
NAME                                    AGE
alertmanagers.monitoring.coreos.com     3m
prometheuses.monitoring.coreos.com      3m
servicemonitors.monitoring.coreos.com   3m
```

Once the Prometheus operator CRDs are registered, run the following command to create a Prometheus.

```console
$ https://raw.githubusercontent.com/kubedb/cli/0.8.0-rc.0/docs/examples/monitoring/coreos-operator/demo-1.yaml
prometheus "prometheus" created
service "prometheus" created
```

### Prometheus Dashboard

Now open prometheus dashboard on browser by running `minikube service prometheus -n demo`.

Or you can get the URL of `prometheus` Service by running following command

```console
$ minikube service prometheus -n demo --url
http://192.168.99.100:30900
```

Now, if you go to the Prometheus Dashboard, you will see that target list is now empty.

## Monitor Elasticsearch with CoreOS Prometheus

Below is the Elasticsearch object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: coreos-prom-es
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
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
```

Here,

- `monitor.agent` indicates the monitoring agent. Currently only valid value currently is `coreos-prometheus-operator`
- `monitor.prometheus` specifies the information for monitoring by prometheus
  - `prometheus.namespace` specifies the namespace where ServiceMonitor is created.
  - `prometheus.labels` specifies the labels applied to ServiceMonitor.
  - `prometheus.port` indicates the port for Elasticsearch exporter endpoint (default is `56790`)
  - `prometheus.interval` indicates the scraping interval (eg, '10s')

Now create this Elasticsearch object with monitoring spec

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-rc.0/docs/examples/elasticsearch/monitoring/coreos-prom-es.yaml
elasticsearch "coreos-prom-es" created
```

KubeDB operator will create a ServiceMonitor object once the Elasticsearch is successfully running.

```console
$ kubedb get es -n demo builtin-prom-es
NAME              STATUS    AGE
builtin-prom-es   Running   5m
```

You can verify it running the following commands

```console
$ kubectl get servicemonitor -n demo --selector="app=kubedb"
NAME                         AGE
kubedb-demo-coreos-prom-es   1m
```

Now, if you go the Prometheus Dashboard, you will see this database endpoint in target list.

<p align="center">
  <kbd>
    <img alt="prometheus-builtin"  src="/docs/images/elasticsearch/coreos-prom-es.png">
  </kbd>
</p>

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run following commands

```console
$ kubectl patch -n demo es/coreos-prom-es -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo es/coreos-prom-es

$ kubectl patch -n demo drmn/coreos-prom-es -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/coreos-prom-es

# In rbac enabled cluster,
# $ kubectl delete clusterrolebindings prometheus-operator  prometheus
# $ kubectl delete clusterrole prometheus-operator prometheus

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Learn about [taking instant backup](/docs/guides/elasticsearch/snapshot/instant_backup.md) of Elasticsearch database using KubeDB.
- Learn how to [schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Learn about initializing [Elasticsearch with Snapshot](/docs/guides/elasticsearch/initialization/snapshot_source.md).
- Learn how to configure [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md).
- Detail concepts of [Elasticsearch object](/docs/concepts/databases/elasticsearch.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
