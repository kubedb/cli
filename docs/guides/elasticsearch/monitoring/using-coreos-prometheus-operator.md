---
title: Monitoring Elasticsearch using Coreos Prometheus Operator
menu:
  docs_0.9.0-rc.2:
    identifier: es-using-coreos-prometheus-operator-monitoring
    name: Coreos Prometheus Operator
    parent: es-monitoring-elasticsearch
    weight: 10
menu_name: docs_0.9.0-rc.2
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

> Note: Yaml files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

This tutorial assumes that you are familiar with Elasticsearch concept.

## Deploy CoreOS-Prometheus Operator

Run the following command to deploy CoreOS-Prometheus operator.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.2/docs/examples/monitoring/coreos-operator/demo-0.yaml
namespace/demo configured
clusterrole.rbac.authorization.k8s.io/prometheus-operator created
serviceaccount/prometheus-operator created
clusterrolebinding.rbac.authorization.k8s.io/prometheus-operator created
deployment.extensions/prometheus-operator created
```

Wait for running the Deployment’s Pods.

```console
$ kubectl get pods -n demo --selector=operator=prometheus
NAME                                   READY     STATUS    RESTARTS   AGE
prometheus-operator-857455484c-mbzsp   1/1       Running   0          57s
```

This CoreOS-Prometheus operator will create some supported Custom Resource Definition (CRD).

```console
$ kubectl get crd
NAME                                          CREATED AT
...
alertmanagers.monitoring.coreos.com           2018-10-08T12:53:46Z
prometheuses.monitoring.coreos.com            2018-10-08T12:53:46Z
servicemonitors.monitoring.coreos.com         2018-10-08T12:53:47Z
...
```

Once the Prometheus CRDs are registered, run the following command to create a Prometheus.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.2/docs/examples/monitoring/coreos-operator/demo-1.yaml
clusterrole.rbac.authorization.k8s.io/prometheus created
serviceaccount/prometheus created
clusterrolebinding.rbac.authorization.k8s.io/prometheus created
prometheus.monitoring.coreos.com/prometheus created
service/prometheus created
```

Verify RBAC stuffs

```console
$ kubectl get clusterroles
NAME                             AGE
...
prometheus                       28s
prometheus-operator              10m
...
```

```console
$ kubectl get clusterrolebindings
NAME                  AGE
...
prometheus            2m
prometheus-operator   11m
...
```

### Prometheus Dashboard

Now open prometheus dashboard on browser by running `minikube service prometheus -n demo`.

Or you can get the URL of `prometheus` Service by running following command

```console
$ minikube service prometheus -n demo --url
http://192.168.99.100:30900
```

If you are not using minikube, browse prometheus dashboard using following address `http://{Node's ExternalIP}:{NodePort of prometheus-service}`.

Now, if you go to the Prometheus Dashboard, you will see that target list is now empty.

## Find out required label for ServiceMonitor

First, check created objects of `Prometheus` kind.

```console
$ kubectl get prometheus --all-namespaces
NAMESPACE   NAME         AGE
demo        prometheus   20m
```

Now if we see the full spec of `prometheus` of `Prometheus` kind, we will see a field called `serviceMonitorSelector`. The value of `matchLabels` under `serviceMonitorSelector` part, is the required label for `KubeDB` monitoring spec `monitor.prometheus.labels`.

```yaml
 $ kubectl get prometheus -n demo prometheus -o yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  creationTimestamp: 2018-11-15T10:40:57Z
  generation: 1
  name: prometheus
  namespace: demo
  resourceVersion: "1661"
  selfLink: /apis/monitoring.coreos.com/v1/namespaces/demo/prometheuses/prometheus
  uid: ef59e6e6-e8c2-11e8-8e44-08002771fd7b
spec:
  resources:
    requests:
      memory: 400Mi
  serviceAccountName: prometheus
  serviceMonitorSelector:
    matchLabels:
      app: kubedb
  version: v1.7.0
```

In this tutorial, the required label is `app: kubedb`.

## Monitor Elasticsearch with CoreOS Prometheus

Below is the Elasticsearch object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: coreos-prom-es
  namespace: demo
spec:
  version: "6.3-v1"
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
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.2/docs/examples/elasticsearch/monitoring/coreos-prom-es.yaml
elasticsearch.kubedb.com/coreos-prom-es created
```

KubeDB operator will create a ServiceMonitor object once the Elasticsearch is successfully running.

```console
$ kubectl get es -n demo coreos-prom-es
NAME             VERSION   STATUS    AGE
coreos-prom-es   6.3-v1    Running   1m
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
$ kubectl patch -n demo es/coreos-prom-es -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo es/coreos-prom-es

$ kubectl delete -n demo deployment/prometheus-operator
$ kubectl delete -n demo service/prometheus
$ kubectl delete -n demo service/prometheus-operated
$ kubectl delete -n demo statefulset.apps/prometheus-prometheus

$ kubectl delete clusterrolebindings prometheus-operator  prometheus
$ kubectl delete clusterrole prometheus-operator prometheus

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
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
