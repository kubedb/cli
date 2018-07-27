---
title: Monitor Memcached using Builtin Prometheus Discovery
menu:
  docs_0.8.0:
    identifier: mc-using-builtin-prometheus-monitoring
    name: Builtin Prometheus Discovery
    parent: mc-monitoring-memcached
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Prometheus with KubeDB

This tutorial will show you how to monitor KubeDB databases using [Prometheus](https://prometheus.io/).

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    45m
demo          Active    10s
kube-public   Active    45m
kube-system   Active    45m
```

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create a Memcached database

KubeDB implements a `Memcached` CRD to define the specification of a Memcached database. Below is the `Memcached` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  name: memcd-mon-prometheus
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
    agent: prometheus.io/builtin
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/memcached/monitoring/builtin-prometheus/demo-1.yaml
memcached "memcd-mon-prometheus" created
```

Here,

- `spec.version` is the version of Memcached database. In this tutorial, a Memcached 1.5.4 database is going to be created.
- `spec.resource` is an optional field that specifies how much CPU and memory (RAM) each Container needs. To learn details about Managing Compute Resources for Containers, please visit [here](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/).
- `spec.monitor` specifies that built-in [Prometheus](https://github.com/prometheus/prometheus) is used to monitor this database instance. KubeDB operator will configure the service of this database in a way that the Prometheus server will automatically find out the service endpoint aka `Memcached Exporter` and will receive metrics from exporter.

KubeDB operator watches for `Memcached` objects using Kubernetes api. When a `Memcached` object is created, KubeDB operator will create a new Deployment and a ClusterIP Service with the matching crd name.

```console
$ kubedb get mc -n demo
NAME                   STATUS    AGE
memcd-mon-prometheus   Running   50s


$ kubedb describe mc -n demo memcd-mon-prometheus
Name:		memcd-mon-prometheus
Namespace:	demo
StartTimestamp:	Tue, 13 Feb 2018 12:24:34 +0600
Status:		Running

Deployment:
  Name:			memcd-mon-prometheus
  Replicas:		3 current / 3 desired
  CreationTimestamp:	Tue, 13 Feb 2018 12:24:35 +0600
  Pods Status:		3 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		memcd-mon-prometheus
  Type:		ClusterIP
  IP:		10.99.92.24
  Port:		db		11211/TCP
  Port:		prom-http	56790/TCP

Monitoring System:
  Agent:	prometheus.io/builtin
  Prometheus:
    Namespace:
    Interval:

Events:
  FirstSeen   LastSeen   Count     From                 Type       Reason       Message
  ---------   --------   -----     ----                 --------   ------       -------
  45s         45s        1         Memcached operator   Normal     Successful   Successfully patched Deployment
  45s         45s        1         Memcached operator   Normal     Successful   Successfully patched Memcached
  46s         46s        1         Memcached operator   Normal     Successful   Successfully created Deployment
  46s         46s        1         Memcached operator   Normal     Successful   Successfully created Memcached
  1m          1m         1         Memcached operator   Normal     Successful   Successfully created Service
```

Since `spec.monitoring` was configured, the database service object is configured accordingly. You can verify it running the following commands:

```console
$ kubectl get services -n demo
NAME                   TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)               AGE
memcd-mon-prometheus   ClusterIP   10.99.92.24   <none>        11211/TCP,56790/TCP   1m
```

```yaml
$ kubectl get services memcd-mon-prometheus -n demo -o yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    monitoring.appscode.com/agent: prometheus.io/builtin
    prometheus.io/path: /kubedb.com/v1alpha1/namespaces/demo/memcacheds/memcd-mon-prometheus/metrics
    prometheus.io/port: "56790"
    prometheus.io/scrape: "true"
  creationTimestamp: 2018-02-13T06:24:34Z
  labels:
    kubedb.com/kind: Memcached
    kubedb.com/name: memcd-mon-prometheus
  name: memcd-mon-prometheus
  namespace: demo
  resourceVersion: "5191"
  selfLink: /api/v1/namespaces/demo/services/memcd-mon-prometheus
  uid: 8efa423f-1086-11e8-801e-080027e82bd4
spec:
  clusterIP: 10.99.92.24
  ports:
  - name: db
    port: 11211
    protocol: TCP
    targetPort: db
  - name: prom-http
    port: 56790
    protocol: TCP
    targetPort: prom-http
  selector:
    kubedb.com/kind: Memcached
    kubedb.com/name: memcd-mon-prometheus
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}
```

We can see that the service contains these specific annotations. The Prometheus server will discover the exporter using these specifications.

```yaml
prometheus.io/path: ...
prometheus.io/port: ...
prometheus.io/scrape: ...
```

## Deploy and configure Prometheus Server

The Prometheus server is needed to configure so that it can discover endpoints of services. If a Prometheus server is already running in cluster and if it is configured in a way that it can discover service endpoints, no extra configuration will be needed. If there is no existing Prometheus server running, rest of this tutorial will create a Prometheus server with appropriate configuration.

The configuration file to `Prometheus-Server` will be provided by `ConfigMap`. The below config map will be created:

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

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/monitoring/builtin-prometheus/demo-1.yaml
configmap "prometheus-server-conf" created
```

Now, the below yaml is used to deploy Prometheus in kubernetes :

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

Run the following command to deploy prometheus-server

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/monitoring/builtin-prometheus/demo-2.yaml
clusterrole "prometheus-server" created
serviceaccount "prometheus-server" created
clusterrolebinding "prometheus-server" created
deployment "prometheus-server" created
service "prometheus-service" created

# Verify RBAC stuffs
$ kubectl get clusterroles
NAME                AGE
prometheus-server   57s

$ kubectl get clusterrolebindings
NAME                AGE
prometheus-server   1m

$ kubectl get serviceaccounts -n demo
NAME                SECRETS   AGE
default             1         48m
prometheus-server   1         1m
```

### Prometheus Dashboard

Now to open prometheus dashboard on Browser:

```console
$ kubectl get svc -n demo
NAME                   TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)               AGE
memcd-mon-prometheus   ClusterIP      10.99.92.24    <none>        11211/TCP,56790/TCP   4m
prometheus-service     LoadBalancer   10.99.73.172   <pending>     9090:30901/TCP        1m

$ minikube ip
192.168.99.100

$ minikube service prometheus-service -n demo --url
http://192.168.99.100:30901
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{prometheus-svc-nodeport}_ to visit Prometheus Dashboard. According to the above example, this URL will be [http://192.168.99.100:30901](http://192.168.99.100:30901).

Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.
![prometheus-builtin](/docs/images/memcached/memcached-builtin.png)

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo mc/memcd-mon-prometheus -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo mc/memcd-mon-prometheus

$ kubectl patch -n demo drmn/memcd-mon-prometheus -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/memcd-mon-prometheus

$ kubectl delete clusterrole prometheus-server
$ kubectl delete clusterrolebindings  prometheus-server
$ kubectl delete serviceaccounts -n demo  prometheus-server

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Monitor your Memcached database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/memcached/monitoring/using-coreos-prometheus-operator.md).
- Detail concepts of [Memcached object](/docs/concepts/databases/memcached.md).
- Use [private Docker registry](/docs/guides/memcached/private-registry/using-private-registry.md) to deploy Memcached with KubeDB.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
