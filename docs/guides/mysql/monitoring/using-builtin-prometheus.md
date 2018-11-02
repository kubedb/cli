---
title: Monitor MySQL using Builtin Prometheus Discovery
menu:
  docs_0.9.0-rc.0:
    identifier: my-using-builtin-prometheus-monitoring
    name: Builtin Prometheus Discovery
    parent: my-monitoring-mysql
    weight: 10
menu_name: docs_0.9.0-rc.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Prometheus with KubeDB

This tutorial will show you how to monitor KubeDB databases using [Prometheus](https://prometheus.io/).

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace "demo" created

  $ kubectl get ns
  NAME          STATUS    AGE
  demo          Active    10s
  ```

> Note: The yaml files that are used in this tutorial are stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Monitor with builtin Prometheus

User can define `spec.monitor` either while creating the CRD object, Or can update the spec of existing CRD object to add the `spec.monitor` part. Below is the `MySQL` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-mon-prometheus
  namespace: demo
spec:
  version: "8.0-v1"
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
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/mysql/monitoring/builtin-prometheus/demo-1.yaml
mysql.kubedb.com/mysql-mon-prometheus created
```

Here,

- `spec.monitor` specifies that built-in [Prometheus](https://github.com/prometheus/prometheus) is used to monitor this database instance. KubeDB operator will configure the service of this database in a way that the Prometheus server will automatically find out the service endpoint aka `MySQL Exporter` and will receive metrics from exporter.

KubeDB will create a separate stats service with name `<mysql-crd-name>-stats` for monitoring purpose. KubeDB operator will configure this monitoring service once the MySQL is successfully running.

```console
$ kubedb get my -n demo
NAME                   VERSION   STATUS    AGE
mysql-mon-prometheus   8.0-v1    Running   13m

$ kubedb describe my -n demo mysql-mon-prometheus
Name:               mysql-mon-prometheus
Namespace:          demo
CreationTimestamp:  Thu, 27 Sep 2018 16:02:43 +0600
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
  Name:               mysql-mon-prometheus
  CreationTimestamp:  Thu, 27 Sep 2018 16:02:45 +0600
  Labels:               kubedb.com/kind=MySQL
                        kubedb.com/name=mysql-mon-prometheus
  Annotations:        <none>
  Replicas:           824639361756 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         mysql-mon-prometheus
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-mon-prometheus
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.105.118.238
  Port:         db  3306/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.5:3306

Service:        
  Name:         mysql-mon-prometheus-stats
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-mon-prometheus
  Annotations:    monitoring.appscode.com/agent=prometheus.io/builtin
                  prometheus.io/path=/metrics
                  prometheus.io/port=56790
                  prometheus.io/scrape=true
  Type:         ClusterIP
  IP:           10.110.18.171
  Port:         prom-http  56790/TCP
  TargetPort:   prom-http/TCP
  Endpoints:    172.17.0.5:56790

Database Secret:
  Name:         mysql-mon-prometheus-auth
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-mon-prometheus
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  password:  16 bytes
  user:      4 bytes

Monitoring System:
  Agent:  prometheus.io/builtin
  Prometheus:
    Port:  56790

No Snapshots.

Events:
  Type    Reason      Age   From            Message
  ----    ------      ----  ----            -------
  Normal  Successful  13m   MySQL operator  Successfully created Service
  Normal  Successful  13m   MySQL operator  Successfully created StatefulSet
  Normal  Successful  13m   MySQL operator  Successfully created MySQL
  Normal  Successful  13m   MySQL operator  Successfully created stats service
  Normal  Successful  12m   MySQL operator  Successfully patched StatefulSet
  Normal  Successful  12m   MySQL operator  Successfully patched MySQL
  Normal  Successful  12m   MySQL operator  Successfully patched StatefulSet
  Normal  Successful  12m   MySQL operator  Successfully patched MySQL
```

Since `spec.monitoring` was configured, the database monitoring service object is configured accordingly. You can verify it running the following commands:

```console
$ kubectl get services -n demo
NAME                         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
kubedb                       ClusterIP   None             <none>        <none>      31m
mysql-mon-prometheus         ClusterIP   10.105.118.238   <none>        3306/TCP    14m
mysql-mon-prometheus-stats   ClusterIP   10.110.18.171    <none>        56790/TCP   13m
```

```yaml
$ kubectl get services mysql-mon-prometheus-stats -n demo -o yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    monitoring.appscode.com/agent: prometheus.io/builtin
    prometheus.io/path: /metrics
    prometheus.io/port: "56790"
    prometheus.io/scrape: "true"
  creationTimestamp: 2018-09-27T10:03:36Z
  labels:
    kubedb.com/kind: MySQL
    kubedb.com/name: mysql-mon-prometheus
  name: mysql-mon-prometheus-stats
  namespace: demo
  ownerReferences:
  - apiVersion: kubedb.com/v1alpha1
    blockOwnerDeletion: false
    kind: MySQL
    name: mysql-mon-prometheus
    uid: 7a30757a-c23c-11e8-b2cc-080027d9f35e
  resourceVersion: "4015"
  selfLink: /api/v1/namespaces/demo/services/mysql-mon-prometheus-stats
  uid: 99cdcd40-c23c-11e8-b2cc-080027d9f35e
spec:
  clusterIP: 10.110.18.171
  ports:
  - name: prom-http
    port: 56790
    protocol: TCP
    targetPort: prom-http
  selector:
    kubedb.com/kind: MySQL
    kubedb.com/name: mysql-mon-prometheus
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
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/monitoring/builtin-prometheus/demo-1.yaml
configmap/prometheus-server-conf created
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
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/monitoring/builtin-prometheus/demo-2.yaml
clusterrole.rbac.authorization.k8s.io/prometheus-server created
serviceaccount/prometheus-server created
clusterrolebinding.rbac.authorization.k8s.io/prometheus-server created
deployment.apps/prometheus-server created
service/prometheus-service created

# Verify RBAC stuffs
$ kubectl get clusterroles prometheus-server
NAME                AGE
prometheus-server   28s

$ kubectl get clusterrolebindings prometheus-server
NAME                AGE
prometheus-server   59s

$ kubectl get serviceaccounts -n demo
NAME                SECRETS   AGE
default             1         52m
prometheus-server   1         1m
```

### Prometheus Dashboard

Now to open prometheus dashboard on Browser:

```console
$ kubectl get svc -n demo
NAME                         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
kubedb                       ClusterIP   None             <none>        <none>           34m
mysql-mon-prometheus         ClusterIP   10.105.118.238   <none>        3306/TCP         17m
mysql-mon-prometheus-stats   ClusterIP   10.110.18.171    <none>        56790/TCP        16m
prometheus-service           NodePort    10.100.155.55    <none>        9090:30901/TCP   1m

$ minikube ip
192.168.99.100

$ minikube service prometheus-service -n demo --url
http://192.168.99.100:30901
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{prometheus-svc-nodeport}_ to visit Prometheus Dashboard. According to the above example, this URL will be [http://192.168.99.100:30901](http://192.168.99.100:30901).

If you are not using minikube, browse prometheus dashboard using following address `http://{Node's ExternalIP}:{NodePort of prometheus-service}`.

Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.
![prometheus-builtin](/docs/images/mysql/mysql-builtin.png)

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mysql/mysql-mon-prometheus -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mysql/mysql-mon-prometheus

kubectl patch -n demo drmn/mysql-mon-prometheus -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mysql-mon-prometheus

kubectl delete -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/monitoring/builtin-prometheus/demo-2.yaml

kubectl delete ns demo
```

## Next Steps

- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Detail concepts of [MySQLVersion object](/docs/concepts/catalog/mysql.md).
- [Snapshot and Restore](/docs/guides/mysql/snapshot/backup-and-restore.md) process of MySQL databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mysql/snapshot/scheduled-backup.md) of MySQL databases using KubeDB.
- Initialize [MySQL with Script](/docs/guides/mysql/initialization/using-script.md).
- Initialize [MySQL with Snapshot](/docs/guides/mysql/initialization/using-snapshot.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
