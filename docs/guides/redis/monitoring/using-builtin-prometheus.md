---
title: Monitoring Redis Using Builtin Prometheus Discovery
menu:
  docs_0.9.0-beta.0:
    identifier: rd-using-builtin-prometheus-monitoring
    name: Builtin Prometheus Discovery
    parent: rd-monitoring-redis
    weight: 10
menu_name: docs_0.9.0-beta.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Prometheus with KubeDB

This tutorial will show you how to monitor KubeDB databases using [Prometheus](https://prometheus.io/).

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

  ```console
  $ kubectl create ns demo
  namespace "demo" created
  
  $ kubectl get ns
  NAME          STATUS    AGE
  demo          Active    10s
  ```

> Note: The yaml files used in this tutorial are stored in [docs/examples/redis](https://github.com/kubedb/cli/tree/master/docs/examples/redis) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

## Monitor with builtin Prometheus

User can define `spec.monitor` either while creating the CRD object, Or can update the spec of existing CRD object to add the `spec.monitor` part. Below is the `Redis` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  name: redis-mon-prometheus
  namespace: demo
spec:
  version: "4.0-v1"
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

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.1/docs/examples/redis/monitoring/builtin-prometheus/demo-1.yaml
redis.kubedb.com/redis-mon-prometheus created
```

Here,

- `spec.monitor` specifies that built-in [Prometheus](https://github.com/prometheus/prometheus) is used to monitor this database instance. KubeDB operator will configure the service of this database in a way that the Prometheus server will automatically find out the service endpoint aka `Redis Exporter` and will receive metrics from exporter.

KubeDB will create a separate stats service with name `<redis-crd-name>-stats` for monitoring purpose. KubeDB operator will configure this monitoring service once the Redis is successfully running.

```console
$ kubedb get rd -n demo
NAME                   VERSION   STATUS    AGE
redis-mon-prometheus   4.0-v1    Running   2m

$ kubedb describe rd -n demo redis-mon-prometheus
Name:               redis-mon-prometheus
Namespace:          demo
CreationTimestamp:  Mon, 01 Oct 2018 12:34:20 +0600
Labels:             <none>
Annotations:        <none>
Replicas:           1  total
Status:             Running
  StorageType:      Durable
Volume:
  StorageClass:  standard
  Capacity:      50Mi
  Access Modes:  RWO

StatefulSet:
  Name:               redis-mon-prometheus
  CreationTimestamp:  Mon, 01 Oct 2018 12:34:22 +0600
  Labels:               kubedb.com/kind=Redis
                        kubedb.com/name=redis-mon-prometheus
  Annotations:        <none>
  Replicas:           824641421356 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         redis-mon-prometheus
  Labels:         kubedb.com/kind=Redis
                  kubedb.com/name=redis-mon-prometheus
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.98.125.255
  Port:         db  6379/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.4:6379

Service:
  Name:         redis-mon-prometheus-stats
  Labels:         kubedb.com/kind=Redis
                  kubedb.com/name=redis-mon-prometheus
  Annotations:    monitoring.appscode.com/agent=prometheus.io/builtin
                  prometheus.io/path=/kubedb.com/v1alpha1/namespaces/demo/redises/redis-mon-prometheus/metrics
                  prometheus.io/port=56790
                  prometheus.io/scrape=true
  Type:         ClusterIP
  IP:           10.104.85.239
  Port:         prom-http  56790/TCP
  TargetPort:   prom-http/TCP
  Endpoints:    172.17.0.4:56790

Monitoring System:
  Agent:  prometheus.io/builtin
  Prometheus:
    Port:  56790

No Snapshots.

Events:
  Type    Reason      Age   From            Message
  ----    ------      ----  ----            -------
  Normal  Successful  2m    Redis operator  Successfully created Service
  Normal  Successful  1m    Redis operator  Successfully created StatefulSet
  Normal  Successful  1m    Redis operator  Successfully created Redis
  Normal  Successful  1m    Redis operator  Successfully created stats service
  Normal  Successful  1m    Redis operator  Successfully patched StatefulSet
  Normal  Successful  1m    Redis operator  Successfully patched Redis
```

Since `spec.monitoring` was configured, the database monitoring service is configured accordingly. You can verify it running the following commands:

```console
$ kubectl get services -n demo
NAME                         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)     AGE
kubedb                       ClusterIP   None            <none>        <none>      2m
redis-mon-prometheus         ClusterIP   10.98.125.255   <none>        6379/TCP    2m
redis-mon-prometheus-stats   ClusterIP   10.104.85.239   <none>        56790/TCP   1m
```

```yaml
$ kubectl get services redis-mon-prometheus-stats -n demo -o yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    monitoring.appscode.com/agent: prometheus.io/builtin
    prometheus.io/path: /kubedb.com/v1alpha1/namespaces/demo/redises/redis-mon-prometheus/metrics
    prometheus.io/port: "56790"
    prometheus.io/scrape: "true"
  creationTimestamp: 2018-10-01T06:35:03Z
  labels:
    kubedb.com/kind: Redis
    kubedb.com/name: redis-mon-prometheus
  name: redis-mon-prometheus-stats
  namespace: demo
  ownerReferences:
  - apiVersion: kubedb.com/v1alpha1
    blockOwnerDeletion: false
    kind: Redis
    name: redis-mon-prometheus
    uid: 076b13b3-c544-11e8-9ba7-0800274bef12
  resourceVersion: "10495"
  selfLink: /api/v1/namespaces/demo/services/redis-mon-prometheus-stats
  uid: 211f339b-c544-11e8-9ba7-0800274bef12
spec:
  clusterIP: 10.104.85.239
  ports:
  - name: prom-http
    port: 56790
    protocol: TCP
    targetPort: prom-http
  selector:
    kubedb.com/kind: Redis
    kubedb.com/name: redis-mon-prometheus
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

Create above ConfigMap

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.1/docs/examples/monitoring/builtin-prometheus/demo-1.yaml
configmap/prometheus-server-conf created
```

Now, the below yaml is used to deploy Prometheus in kubernetes:

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
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.1/docs/examples/monitoring/builtin-prometheus/demo-2.yaml
clusterrole.rbac.authorization.k8s.io/prometheus-server created
serviceaccount/prometheus-server created
clusterrolebinding.rbac.authorization.k8s.io/prometheus-server created
deployment.apps/prometheus-server created
service/prometheus-service created

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
NAME                         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
kubedb                       ClusterIP   None             <none>        <none>           5m
prometheus-service           NodePort    10.108.252.226   <none>        9090:30901/TCP   37s
redis-mon-prometheus         ClusterIP   10.98.125.255    <none>        6379/TCP         5m
redis-mon-prometheus-stats   ClusterIP   10.104.85.239    <none>        56790/TCP        5m

$ minikube ip
192.168.99.100

$ minikube service prometheus-service -n demo --url
http://192.168.99.100:30901
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{prometheus-svc-nodeport}_ to visit Prometheus Dashboard. According to the above example, this URL will be [http://192.168.99.100:30901](http://192.168.99.100:30901).

If you are not using minikube, browse prometheus dashboard using following address `http://{Node's ExternalIP}:{NodePort of prometheus-service}`.

Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.

![prometheus-builtin](/docs/images/redis/redis-builtin.png)

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo rd/redis-mon-prometheus -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo rd/redis-mon-prometheus

kubectl patch -n demo drmn/redis-mon-prometheus -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/redis-mon-prometheus

kubectl delete -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.1/docs/examples/monitoring/builtin-prometheus/demo-2.yaml

kubectl delete ns demo
```

## Next Steps

- Monitor your Redis database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/redis/monitoring/using-coreos-prometheus-operator.md).
- Detail concepts of [Redis object](/docs/concepts/databases/redis.md).
- Use [private Docker registry](/docs/guides/redis/private-registry/using-private-registry.md) to deploy Redis with KubeDB.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
