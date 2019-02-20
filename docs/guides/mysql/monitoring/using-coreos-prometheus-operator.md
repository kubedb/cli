---
title: Monitor MySQL using Coreos Prometheus Operator
menu:
  docs_0.9.0:
    identifier: my-using-coreos-prometheus-operator-monitoring
    name: Coreos Prometheus Operator
    parent: my-monitoring-mysql
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Monitoring MySQL Using CoreOS Prometheus Operator

CoreOS [prometheus-operator](https://github.com/coreos/prometheus-operator) provides simple and Kubernetes native way to deploy and configure Prometheus server. This tutorial will show you how to use CoreOS Prometheus operator to monitor MySQL database deployed with KubeDB.

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- To learn how Prometheus monitoring works with KubeDB in general, please visit [here](/docs/concepts/database-monitoring/overview.md).

- To keep Prometheus resources isolated, we are going to use a separate namespace called `monitoring` to deploy respective monitoring resources. We are going to deploy database in `demo` namespace.

  ```console
  $ kubectl create ns monitoring
  namespace/monitoring created

  $ kubectl create ns demo
  namespace/demo created
  ```

- We need a CoreOS [prometheus-operator](https://github.com/coreos/prometheus-operator) instance running. If you don't already have a running instance, deploy one following the docs from [here](https://github.com/appscode/third-party-tools/blob/master/monitoring/prometheus/coreos-operator/README.md).

- If you already don't have a Prometheus server running, deploy one following tutorial from [here](https://github.com/appscode/third-party-tools/blob/master/monitoring/prometheus/coreos-operator/README.md#deploy-prometheus-server).

> Note: YAML files used in this tutorial are stored in [docs/examples/mysql](https://github.com/kubedb/cli/tree/master/docs/examples/mysql) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Find out required labels for ServiceMonitor

We need to know the labels used to select `ServiceMonitor` by a `Prometheus` crd. We are going to provide these labels in `spec.monitor.prometheus.labels` field of MySQL crd so that KubeDB creates `ServiceMonitor` object accordingly.

At first, let's find out the available Prometheus server in our cluster.

```console
$ kubectl get prometheus --all-namespaces
NAMESPACE    NAME         AGE
monitoring   prometheus   18m
```

> If you don't have any Prometheus server running in your cluster, deploy one following the guide specified in **Before You Begin** section.

Now, let's view the YAML of the available Prometheus server `prometheus` in `monitoring` namespace.

```yaml
$ kubectl get prometheus -n monitoring prometheus -o yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"monitoring.coreos.com/v1","kind":"Prometheus","metadata":{"annotations":{},"labels":{"prometheus":"prometheus"},"name":"prometheus","namespace":"monitoring"},"spec":{"replicas":1,"resources":{"requests":{"memory":"400Mi"}},"serviceAccountName":"prometheus","serviceMonitorSelector":{"matchLabels":{"k8s-app":"prometheus"}}}}
  creationTimestamp: 2019-01-03T13:41:51Z
  generation: 1
  labels:
    prometheus: prometheus
  name: prometheus
  namespace: monitoring
  resourceVersion: "44402"
  selfLink: /apis/monitoring.coreos.com/v1/namespaces/monitoring/prometheuses/prometheus
  uid: 5324ad98-0f5d-11e9-b230-080027f306f3
spec:
  replicas: 1
  resources:
    requests:
      memory: 400Mi
  serviceAccountName: prometheus
  serviceMonitorSelector:
    matchLabels:
      k8s-app: prometheus
```

Notice the `spec.serviceMonitorSelector` section. Here, `k8s-app: prometheus` label is used to select `ServiceMonitor` crd. So, we are going to use this label in `spec.monitor.prometheus.labels` field of MySQL crd.

## Deploy MySQL with Monitoring Enabled

At first, let's deploy an MySQL database with monitoring enabled. Below is the MySQL object that we are going to create.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: coreos-prom-mysql
  namespace: demo
spec:
  version: "8.0-v2"
  terminationPolicy: WipeOut
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  monitor:
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: monitoring
      labels:
        k8s-app: prometheus
      interval: 10s
```

Here,

- `monitor.agent:  prometheus.io/coreos-operator` indicates that we are going to monitor this server using CoreOS prometheus operator.
- `monitor.prometheus.namespace: monitoring` specifies that KubeDB should create `ServiceMonitor` in `monitoring` namespace.

- `monitor.prometheus.labels` specifies that KubeDB should create `ServiceMonitor` with these labels.

- `monitor.prometheus.interval` indicates that the Prometheus server should scrape metrics from this database with 10 seconds interval.

Let's create the MySQL object that we have shown above,

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/monitoring/coreos-prom-mysql.yaml
mysql.kubedb.com/coreos-prom-mysql created
```

Now, wait for the database to go into `Running` state.

```console
$ kubectl get mysql -n demo coreos-prom-mysql
NAME                VERSION   STATUS    AGE
coreos-prom-mysql   8.0-v2    Running   17s
```

KubeDB will create a separate stats service with name `{MySQL crd name}-stats` for monitoring purpose.

```console
$ kubectl get svc -n demo --selector="kubedb.com/name=coreos-prom-mysql"
NAME                      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
coreos-prom-mysql         ClusterIP   10.103.216.168   <none>        3306/TCP    40s
coreos-prom-mysql-stats   ClusterIP   10.111.190.240   <none>        56790/TCP   38s
```

Here, `coreos-prom-mysql-stats` service has been created for monitoring purpose.

Let's describe this stats service.

```yaml
$ kubectl describe svc -n demo coreos-prom-mysql-stats
Name:              coreos-prom-mysql-stats
Namespace:         demo
Labels:            kubedb.com/kind=MySQL
                   kubedb.com/name=coreos-prom-mysql
Annotations:       monitoring.appscode.com/agent: prometheus.io/coreos-operator
Selector:          kubedb.com/kind=MySQL,kubedb.com/name=coreos-prom-mysql
Type:              ClusterIP
IP:                10.111.190.240
Port:              prom-http  56790/TCP
TargetPort:        prom-http/TCP
Endpoints:         172.17.0.7:56790
Session Affinity:  None
Events:            <none>
```

Notice the `Labels` and `Port` fields. `ServiceMonitor` will use these information to target its endpoints.

KubeDB will also create a `ServiceMonitor` crd in `monitoring` namespace that select the endpoints of `coreos-prom-mysql-stats` service. Verify that the `ServiceMonitor` crd has been created.

```console
$ kubectl get servicemonitor -n monitoring
NAME                            AGE
kubedb-demo-coreos-prom-mysql   1m
```

Let's verify that the `ServiceMonitor` has the label that we had specified in `spec.monitor` section of MySQL crd.

```yaml
$ kubectl get servicemonitor -n monitoring kubedb-demo-coreos-prom-mysql -o yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  creationTimestamp: 2019-01-03T15:36:03Z
  generation: 1
  labels:
    k8s-app: prometheus
    monitoring.appscode.com/service: coreos-prom-mysql-stats.demo
  name: kubedb-demo-coreos-prom-mysql
  namespace: monitoring
  resourceVersion: "53056"
  selfLink: /apis/monitoring.coreos.com/v1/namespaces/monitoring/servicemonitors/kubedb-demo-coreos-prom-mysql
  uid: 474d6032-0f6d-11e9-b230-080027f306f3
spec:
  endpoints:
  - honorLabels: true
    interval: 10s
    path: /metrics
    port: prom-http
  namespaceSelector:
    matchNames:
    - demo
  selector:
    matchLabels:
      kubedb.com/kind: MySQL
      kubedb.com/name: coreos-prom-mysql
```

Notice that the `ServiceMonitor` has label `k8s-app: prometheus` that we had specified in MySQL crd.

Also notice that the `ServiceMonitor` has selector which match the labels we have seen in the `coreos-prom-mysql-stats` service. It also, target the `prom-http` port that we have seen in the stats service.

## Verify Monitoring Metrics

At first, let's find out the respective Prometheus pod for `prometheus` Prometheus server.

```console
$ kubectl get pod -n monitoring -l=app=prometheus
NAME                      READY   STATUS    RESTARTS   AGE
prometheus-prometheus-0   3/3     Running   1          63m
```

Prometheus server is listening to port `9090` of `prometheus-prometheus-0` pod. We are going to use [port forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) to access Prometheus dashboard.

Run following command on a separate terminal to forward the port 9090 of `prometheus-prometheus-0` pod,

```console
$ kubectl port-forward -n monitoring prometheus-prometheus-0 9090
Forwarding from 127.0.0.1:9090 -> 9090
Forwarding from [::1]:9090 -> 9090
```

Now, we can access the dashboard at `localhost:9090`. Open [http://localhost:9090](http://localhost:9090) in your browser. You should see `prom-http` endpoint of `coreos-prom-mysql-stats` service as one of the targets.

<p align="center">
  <img alt="Prometheus Target" src="/docs/images/mysql/monitoring/mysql-coreos-prom-target.png" style="padding:10px">
</p>

Check the `endpoint` and `service` labels marked by red rectangle. It verifies that the target is our expected database. Now, you can view the collected metrics and create a graph from homepage of this Prometheus dashboard. You can also use this Prometheus server as data source for [Grafana](https://grafana.com/) and create beautiful dashboard with collected metrics.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run following commands

```console
# cleanup database
kubectl delete -n demo my/coreos-prom-mysql

# cleanup prometheus resources
kubectl delete -n monitoring prometheus prometheus
kubectl delete -n monitoring clusterrolebinding prometheus
kubectl delete -n monitoring clusterrole prometheus
kubectl delete -n monitoring serviceaccount prometheus
kubectl delete -n monitoring service prometheus-operated

# cleanup prometheus operator resources
kubectl delete -n monitoring deployment prometheus-operator
kubectl delete -n dmeo serviceaccount prometheus-operator
kubectl delete clusterrolebinding prometheus-operator
kubectl delete clusterrole prometheus-operator

# delete namespace
kubectl delete ns monitoring
kubectl delete ns demo
```

## Next Steps

- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Detail concepts of [MySQLVersion object](/docs/concepts/catalog/mysql.md).
- [Snapshot and Restore](/docs/guides/mysql/snapshot/backup-and-restore.md) process of MySQL databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mysql/snapshot/scheduled-backup.md) of MySQL databases using KubeDB.
- Initialize [MySQL with Script](/docs/guides/mysql/initialization/using-script.md).
- Initialize [MySQL with Snapshot](/docs/guides/mysql/initialization/using-snapshot.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
