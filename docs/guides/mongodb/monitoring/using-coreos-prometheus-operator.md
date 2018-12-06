---
title: Monitor MongoDB using Coreos Prometheus Operator
menu:
  docs_0.9.0-rc.2:
    identifier: mg-using-coreos-prometheus-operator-monitoring
    name: Coreos Prometheus Operator
    parent: mg-monitoring-mongodb
    weight: 15
menu_name: docs_0.9.0-rc.2
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Prometheus (CoreOS operator) with KubeDB

This tutorial will show you how to monitor KubeDB databases using Prometheus via [CoreOS Prometheus Operator](https://github.com/coreos/prometheus-operator).

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

> Note: The yaml files used in this tutorial are stored in [docs/examples/mongodb](https://github.com/kubedb/cli/tree/master/docs/examples/mongodb) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Deploy CoreOS-Prometheus Operator

Run the following command to deploy CoreOS-Prometheus operator.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.2/docs/examples/monitoring/coreos-operator/demo-0.yaml
namespace/demo created
clusterrole.rbac.authorization.k8s.io/prometheus-operator created
serviceaccount/prometheus-operator created
clusterrolebinding.rbac.authorization.k8s.io/prometheus-operator created
deployment.extensions/prometheus-operator created
```

Wait for running the Deployment’s Pods.

```console
$ kubectl get pods -n demo
NAME                                   READY     STATUS    RESTARTS   AGE
prometheus-operator-857455484c-skbnp   1/1       Running   0          21s
```

This CoreOS-Prometheus operator will create some supported Custom Resource Definition (CRD).

```console
$ kubectl get crd
NAME                                          CREATED AT
...
alertmanagers.monitoring.coreos.com           2018-09-24T12:42:22Z
prometheuses.monitoring.coreos.com            2018-09-24T12:42:22Z
servicemonitors.monitoring.coreos.com         2018-09-24T12:42:22Z
...
```

Once the Prometheus operator CRDs are registered, run the following command to create a Prometheus.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.2/docs/examples/monitoring/coreos-operator/demo-1.yaml
clusterrole.rbac.authorization.k8s.io/prometheus created
serviceaccount/prometheus created
clusterrolebinding.rbac.authorization.k8s.io/prometheus created
prometheus.monitoring.coreos.com/prometheus created
service/prometheus created

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
NAME                  TYPE           CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
prometheus            LoadBalancer   10.97.56.77   <pending>     9090:30900/TCP   34s
prometheus-operated   ClusterIP      None          <none>        9090/TCP         33s

$ minikube ip
192.168.99.100

$ minikube service prometheus -n demo --url
http://192.168.99.100:30900
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{prometheus-svc-nodeport}_ to visit Prometheus Dashboard. According to the above example, this URL will be [http://192.168.99.100:30900](http://192.168.99.100:30900).

If you are not using minikube, browse prometheus dashboard using following address `http://{Node's ExternalIP}:{NodePort of prometheus-service}`.

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

## Monitor Mongodb with CoreOS Prometheus

KubeDB implements a `MongoDB` CRD to define the specification of a MongoDB database. Below is the `MongoDB` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-mon-coreos
  namespace: demo
spec:
  version: "3.4-v1"
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

The `MongoDB` CRD object contains `monitor` field in it's `spec`.  It is also possible to add CoreOS-Prometheus monitor to an existing `MongoDB` database by adding the below part in it's `spec` field.

Here, `spec.monitor.prometheus.labels` is the `serviceMonitorSelector` that we found earlier.

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

| Keys                                | Value  | Description                                                                                                  |
|-------------------------------------|--------|--------------------------------------------------------------------------------------------------------------|
| `spec.monitor.agent`                | string | `Required`. Indicates the monitoring agent used. Only valid value currently is `coreos-prometheus-operator`  |
| `spec.monitor.prometheus.namespace` | string | `Required`. Indicates namespace where service monitors are created. This must be the same namespace of the Prometheus instance. |
| `spec.monitor.prometheus.labels`    | map    | `Required`. Indicates labels applied to service monitor.                                                     |
| `spec.monitor.prometheus.interval`  | string | `Optional`. Indicates the scrape interval for database exporter endpoint (eg, '10s')                         |
| `spec.monitor.prometheus.port`      | int    | `Optional`. Indicates the port for database exporter endpoint (default is `56790`)                           |

__Known Limitations:__ If the database password is updated, exporter must be restarted to use the new credentials. This issue is tracked [here](https://github.com/kubedb/project/issues/53).

Run the following command to deploy the above `MongoDB` CRD object.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.2/docs/examples/mongodb/monitoring/coreos-operator/demo-1.yaml
mongodb.kubedb.com/mgo-mon-coreos created
```

Here,

- `spec.monitor` specifies that CoreOS Prometheus operator is used to monitor this database instance. A ServiceMonitor should be created in the `demo` namespace with label `app=kubedb`. The exporter endpoint should be scrapped every 10 seconds.

KubeDB will create a separate stats service with name `<mongodb-crd-name>-stats` for monitoring purpose. KubeDB operator will configure this monitoring service once the MongoDB is successfully running.

```console
$ kubedb get mg -n demo
NAME             VERSION   STATUS    AGE
mgo-mon-coreos   3.4-v1    Running   1m

$ kubedb describe mg -n demo mgo-mon-coreos
Name:               mgo-mon-coreos
Namespace:          demo
CreationTimestamp:  Tue, 25 Sep 2018 11:56:23 +0600
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
  Name:               mgo-mon-coreos
  CreationTimestamp:  Tue, 25 Sep 2018 11:56:26 +0600
  Labels:               kubedb.com/kind=MongoDB
                        kubedb.com/name=mgo-mon-coreos
  Annotations:        <none>
  Replicas:           824636593232 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         mgo-mon-coreos
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-mon-coreos
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.110.192.170
  Port:         db  27017/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.7:27017

Service:
  Name:         mgo-mon-coreos-gvr
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-mon-coreos
  Annotations:    service.alpha.kubernetes.io/tolerate-unready-endpoints=true
  Type:         ClusterIP
  IP:           None
  Port:         db  27017/TCP
  TargetPort:   27017/TCP
  Endpoints:    172.17.0.7:27017

Service:
  Name:         mgo-mon-coreos-stats
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-mon-coreos
  Annotations:    monitoring.appscode.com/agent=prometheus.io/coreos-operator
  Type:         ClusterIP
  IP:           10.103.111.174
  Port:         prom-http  56790/TCP
  TargetPort:   prom-http/TCP
  Endpoints:    172.17.0.7:56790

Database Secret:
  Name:         mgo-mon-coreos-auth
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-mon-coreos
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  user:      4 bytes
  password:  16 bytes

Monitoring System:
  Agent:  prometheus.io/coreos-operator
  Prometheus:
    Port:       56790
    Namespace:  demo
    Labels:     app=kubedb
    Interval:   10s

No Snapshots.

Events:
  Type    Reason      Age   From              Message
  ----    ------      ----  ----              -------
  Normal  Successful  1m    MongoDB operator  Successfully created Service
  Normal  Successful  26s   MongoDB operator  Successfully created StatefulSet
  Normal  Successful  26s   MongoDB operator  Successfully created MongoDB
  Normal  Successful  24s   MongoDB operator  Successfully created stats service
  Normal  Successful  22s   MongoDB operator  Successfully patched StatefulSet
  Normal  Successful  22s   MongoDB operator  Successfully patched MongoDB
  Normal  Successful  21s   MongoDB operator  Successfully patched StatefulSet
  Normal  Successful  21s   MongoDB operator  Successfully patched MongoDB
```

Since `spec.monitoring` was configured, a ServiceMonitor object is created accordingly. You can verify it running the following commands:

```yaml
$ kubectl get servicemonitor -n demo
NAME                         AGE
kubedb-demo-mgo-mon-coreos   1m

$ kubectl get servicemonitor -n demo kubedb-demo-mgo-mon-coreos -o yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  creationTimestamp: 2018-09-25T05:57:17Z
  generation: 1
  labels:
    app: kubedb
    monitoring.appscode.com/service: mgo-mon-coreos-stats.demo
  name: kubedb-demo-mgo-mon-coreos
  namespace: demo
  resourceVersion: "9093"
  selfLink: /apis/monitoring.coreos.com/v1/namespaces/demo/servicemonitors/kubedb-demo-mgo-mon-coreos
  uid: dbec02e9-c087-11e8-b4a9-0800272618ed
spec:
  endpoints:
  - interval: 10s
    path: /metrics
    port: prom-http
    targetPort: 0
  namespaceSelector:
    matchNames:
    - demo
  selector:
    matchLabels:
      kubedb.com/kind: MongoDB
      kubedb.com/name: mgo-mon-coreos
```

Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mg/mgo-mon-coreos -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mg/mgo-mon-coreos

kubectl patch -n demo drmn/mgo-mon-coreos -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mgo-mon-coreos

kubectl delete -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.2/docs/examples/monitoring/coreos-operator/demo-1.yaml
kubectl delete -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.2/docs/examples/monitoring/coreos-operator/demo-0.yaml

kubectl delete ns demo
```

## Next Steps

- Monitor your MongoDB database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mongodb/monitoring/using-builtin-prometheus.md).
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Detail concepts of [MongoDBVersion object](/docs/concepts/catalog/mongodb.md).
- [Snapshot and Restore](/docs/guides/mongodb/snapshot/backup-and-restore.md) process of MongoDB databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mongodb/snapshot/scheduled-backup.md) of MongoDB databases using KubeDB.
- Initialize [MongoDB with Script](/docs/guides/mongodb/initialization/using-script.md).
- Initialize [MongoDB with Snapshot](/docs/guides/mongodb/initialization/using-snapshot.md).
- Use [private Docker registry](/docs/guides/mongodb/private-registry/using-private-registry.md) to deploy MongoDB with KubeDB.
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
