---
title: Monitor PostgreSQL using Coreos Prometheus Operator
menu:
  docs_0.8.0-beta.2:
    identifier: pg-using-coreos-prometheus-operator-monitoring
    name: Coreos Prometheus Operator
    parent: pg-monitoring-postgres
    weight: 15
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Prometheus (CoreOS operator) with KubeDB

This tutorial will show you how to monitor PostgreSQL using Prometheus via [CoreOS Prometheus Operator](https://github.com/coreos/prometheus-operator).

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: Yaml files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

This tutorial assumes that you are familiar with PostgreSQL concept.

## Deploy CoreOS-Prometheus Operator

### In RBAC enabled cluster

If RBAC is enabled, Run the following command to prepare your cluster for this tutorial:

```console
 $ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/monitoring/coreos-operator/rbac/demo-0.yaml
clusterrole "prometheus-operator" created
serviceaccount "prometheus-operator" created
clusterrolebinding "prometheus-operator" created
deployment "prometheus-operator" created
```

Watch the Deployment’s Pods.

```console
$ kubectl get pods -n demo --watch
NAME                                   READY     STATUS    RESTARTS   AGE
prometheus-operator-79cb9dcd4b-2njgq   1/1       Running   0          2m
```

This CoreOS-Prometheus operator will create some supported Custom Resource Definition (CRD).

```console
$ kubectl get crd
NAME                                    AGE
alertmanagers.monitoring.coreos.com     11m
prometheuses.monitoring.coreos.com      11m
servicemonitors.monitoring.coreos.com   11m
```

Once the Prometheus operator CRDs are registered, run the following command to create a Prometheus.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/monitoring/coreos-operator/rbac/demo-1.yaml
clusterrole "prometheus" created
serviceaccount "prometheus" created
clusterrolebinding "prometheus" created
prometheus "prometheus" created
service "prometheus" created
```

Verify RBAC stuffs

```consolw
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

### In RBAC *not* enabled cluster

If RBAC is not enabled, Run the following command to prepare your cluster for this tutorial:

```console
$ https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/monitoring/coreos-operator/demo-0.yaml
namespace "demo" created
deployment "prometheus-operator" created
```

Watch the Deployment’s Pods.

```console
$ kubectl get pods -n demo --watch
NAME                                   READY     STATUS              RESTARTS   AGE
prometheus-operator-5dcd844486-nprmk   0/1       ContainerCreating   0          27s
prometheus-operator-5dcd844486-nprmk   1/1       Running   0         46s
```

This CoreOS-Prometheus operator will create some supported Custom Resource Definition (CRD).

```console
$ kubectl get crd
NAME                                    AGE
alertmanagers.monitoring.coreos.com     45s
prometheuses.monitoring.coreos.com      44s
servicemonitors.monitoring.coreos.com   44s
```

Once the Prometheus operator CRDs are registered, run the following command to create a Prometheus.

```console
$ https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/monitoring/coreos-operator/demo-1.yaml
prometheus "prometheus" created
service "prometheus" created

```

#### Prometheus Dashboard

Now open prometheus dashboard on browser by running `minikube service prometheus-service -n demo`.

Or you can get the URL of `prometheus` Service by running following command

```console
$ minikube service prometheus -n demo --url
http://192.168.99.100:30900
```

## Monitor PostgreSQL with CoreOS Prometheus

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: coreos-prom-postgres
  namespace: demo
spec:
  version: 9.6
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
      - `prometheus.port` indicates the port for PostgreSQL exporter endpoint (default is `56790`)
      - `prometheus.interval` indicates the scraping interval (eg, '10s')


Now create PostgreSQL with monitoring spec

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/postgres/monitoring/coreos-prom-postgres.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/postgres/monitoring/coreos-prom-postgres.yaml"
postgres "coreos-prom-postgres" created
```

KubeDB operator will create a ServiceMonitor object once the PostgreSQL is successfully running.

```yaml
$ kubectl get servicemonitor -n demo
NAME                               AGE
kubedb-demo-coreos-prom-postgres   1s
```

Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.

<p align="center">
  <kbd>
    <img alt="prometheus-builtin"  src="/docs/images/postgres/coreos-prom-postgres.png">
  </kbd>
</p>

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete pg,drmn,snap -n demo --all --force
$ kubectl delete ns demo
```

## Next Steps
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
