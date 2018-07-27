---
title: Monitor Memcached using Coreos Prometheus Operator
menu:
  docs_0.8.0:
    identifier: mc-using-coreos-prometheus-operator-monitoring
    name: Coreos Prometheus Operator
    parent: mc-monitoring-memcached
    weight: 15
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Prometheus (CoreOS operator) with KubeDB

This tutorial will show you how to monitor KubeDB databases using Prometheus via [CoreOS Prometheus Operator](https://github.com/coreos/prometheus-operator).

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Deploy CoreOS-Prometheus Operator

Run the following command to deploy CoreOS-Prometheus operator.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/monitoring/coreos-operator/demo-0.yaml
namespace "demo" created
clusterrole "prometheus-operator" created
serviceaccount "prometheus-operator" created
clusterrolebinding "prometheus-operator" created
deployment "prometheus-operator" created

$ kubectl get pods -n demo --watch
NAME                                   READY     STATUS    RESTARTS   AGE
prometheus-operator-79cb9dcd4b-2njgq   1/1       Running   0          2m

$ kubectl get crd
NAME                                    AGE
alertmanagers.monitoring.coreos.com     11m
prometheuses.monitoring.coreos.com      11m
servicemonitors.monitoring.coreos.com   11m
```

Once the Prometheus operator CRDs are registered, run the following command to create a Prometheus.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/monitoring/coreos-operator/demo-1.yaml
clusterrole "prometheus" created
serviceaccount "prometheus" created
clusterrolebinding "prometheus" created
prometheus "prometheus" created
service "prometheus" created

# Verify RBAC stuffs
$ kubectl get clusterroles
NAME                  AGE
prometheus            48s
prometheus-operator   1m

$ kubectl get clusterrolebindings
NAME                  AGE
prometheus            7s
prometheus-operator   25s

$ kubectl get serviceaccounts -n demo
NAME                  SECRETS   AGE
default               1         5m
prometheus            1         4m
prometheus-operator   1         5m
```

### Prometheus Dashboard

Now to open prometheus dashboard on Browser:

```console
$ kubectl get svc -n demo
NAME                  TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
prometheus            LoadBalancer   10.104.95.211   <pending>     9090:30900/TCP   5s
prometheus-operated   ClusterIP      None            <none>        9090/TCP         5s

$ minikube ip
192.168.99.100

$ minikube service prometheus -n demo --url
http://192.168.99.100:30900
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{prometheus-svc-nodeport}_ to visit Prometheus Dashboard. According to the above example, this URL will be [http://192.168.99.100:30900](http://192.168.99.100:30900).

## Create a Memcached database

KubeDB implements a `Memcached` CRD to define the specification of a Memcached database. Below is the `Memcached` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  name: memcd-mon-coreos
  namespace: demo
spec:
  replicas: 3
  version: "1.5.4"
  doNotPause: true
  resources:
    requests:
      memory: 64Mi
      cpu: 250m
    limits:
      memory: 128Mi
      cpu: 500m
  monitor:
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
```

The `Memcached` CRD object contains `monitor` field in it's `spec`.  It is also possible to add CoreOS-Prometheus monitor to an existing `Memcached` database by adding the below part in it's `spec` field.

```yaml
spec:
  monitor:
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
```

|  Keys  |  Value |  Description |
|--|--|--|
| `spec.monitor.agent` | string | `Required`. Indicates the monitoring agent used. Only valid value currently is `coreos-prometheus-operator` |
| `spec.monitor.prometheus.namespace` | string | `Required`. Indicates namespace where service monitors are created. This must be the same namespace of the Prometheus instance. |
| `spec.monitor.prometheus.labels` | map | `Required`. Indicates labels applied to service monitor.                                                    |
| `spec.monitor.prometheus.interval` | string | `Optional`. Indicates the scrape interval for database exporter endpoint (eg, '10s')                        |
| `spec.monitor.prometheus.port` | int |`Optional`. Indicates the port for database exporter endpoint (default is `56790`)|

__Known Limitations:__ If the database password is updated, exporter must be restarted to use the new credentials. This issue is tracked [here](https://github.com/kubedb/project/issues/53).

Run the following command to deploy the above `Memcached` CRD object.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/memcached/monitoring/coreos-operator/demo-1.yaml
memcached "memcd-mon-coreos" created
```

Here,

- `spec.version` is the version of Memcached database. In this tutorial, a Memcached 1.5.4 database is going to be created.
- `spec.resource` is an optional field that specifies how much CPU and memory (RAM) each Container needs. To learn details about Managing Compute Resources for Containers, please visit [here](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/).
- `spec.monitor` specifies that CoreOS Prometheus operator is used to monitor this database instance. A ServiceMonitor should be created in the `demo` namespace with label `app=kubedb`. The exporter endpoint should be scrapped every 10 seconds.

KubeDB operator watches for `Memcached` objects using Kubernetes api. When a `Memcached` object is created, KubeDB operator will create a new Deployment and a ClusterIP Service with the matching crd name.

```console
$ kubedb get mc -n demo
NAME               STATUS    AGE
memcd-mon-coreos   Running   1m

$ kubedb describe mc -n demo memcd-mon-coreos
Name:		memcd-mon-coreos
Namespace:	demo
StartTimestamp:	Tue, 13 Feb 2018 12:17:21 +0600
Status:		Running

Deployment:
  Name:			memcd-mon-coreos
  Replicas:		3 current / 3 desired
  CreationTimestamp:	Tue, 13 Feb 2018 12:17:22 +0600
  Pods Status:		3 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		memcd-mon-coreos
  Type:		ClusterIP
  IP:		10.104.166.35
  Port:		db		11211/TCP
  Port:		prom-http	56790/TCP

Monitoring System:
  Agent:	prometheus.io/coreos-operator
  Prometheus:
    Namespace:	demo
    Labels:	app=kubedb
    Interval:	10s

Events:
  FirstSeen   LastSeen   Count     From                 Type       Reason       Message
  ---------   --------   -----     ----                 --------   ------       -------
  1m          1m         1         Memcached operator   Normal     Successful   Successfully patched Deployment
  1m          1m         1         Memcached operator   Normal     Successful   Successfully patched Memcached
  1m          1m         1         Memcached operator   Normal     Successful   Successfully created Deployment
  1m          1m         1         Memcached operator   Normal     Successful   Successfully created Memcached
  1m          1m         1         Memcached operator   Normal     Successful   Successfully created Service
```

Since `spec.monitoring` was configured, a ServiceMonitor object is created accordingly. You can verify it running the following commands:

```yaml
$ kubectl get servicemonitor -n demo
NAME                           AGE
kubedb-demo-memcd-mon-coreos   1m

$ kubectl get servicemonitor -n demo kubedb-demo-memcd-mon-coreos -o yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-13T06:17:36Z
  labels:
    app: kubedb
    monitoring.appscode.com/service: memcd-mon-coreos.demo
  name: kubedb-demo-memcd-mon-coreos
  namespace: demo
  resourceVersion: "4743"
  selfLink: /apis/monitoring.coreos.com/v1/namespaces/demo/servicemonitors/kubedb-demo-memcd-mon-coreos
  uid: 961a507b-1085-11e8-801e-080027e82bd4
spec:
  endpoints:
  - interval: 10s
    path: /kubedb.com/v1alpha1/namespaces/demo/memcacheds/memcd-mon-coreos/metrics
    port: prom-http
    targetPort: 0
  namespaceSelector:
    matchNames:
    - demo
  selector:
    matchLabels:
      kubedb.com/kind: Memcached
      kubedb.com/name: memcd-mon-coreos
```

Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.
![prometheus-coreos](/docs/images/memcached/memcached-coreos.png)

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo mc/memcd-mon-coreos -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo mc/memcd-mon-coreos

$ kubectl patch -n demo drmn/memcd-mon-coreos -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/memcd-mon-coreos

$ kubectl delete clusterrolebindings prometheus-operator  prometheus
$ kubectl delete clusterrole prometheus-operator prometheus

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Monitor your Memcached database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/memcached/monitoring/using-builtin-prometheus.md).
- Detail concepts of [Memcached object](/docs/concepts/databases/memcached.md).
- Use [private Docker registry](/docs/guides/memcached/private-registry/using-private-registry.md) to deploy Memcached with KubeDB.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
