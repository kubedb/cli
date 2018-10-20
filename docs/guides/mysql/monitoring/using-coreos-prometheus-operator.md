---
title: Monitor MySQL using Coreos Prometheus Operator
menu:
  docs_0.9.0-rc.0:
    identifier: my-using-coreos-prometheus-operator-monitoring
    name: Coreos Prometheus Operator
    parent: my-monitoring-mysql
    weight: 10
menu_name: docs_0.9.0-rc.0
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

> Note: The yaml files that are used in this tutorial are stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Deploy CoreOS-Prometheus Operator

Run the following command to deploy CoreOS-Prometheus operator.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/monitoring/coreos-operator/demo-0.yaml
namespace/demo created
clusterrole.rbac.authorization.k8s.io/prometheus-operator created
serviceaccount/prometheus-operator created
clusterrolebinding.rbac.authorization.k8s.io/prometheus-operator created
deployment.extensions/prometheus-operator created
```

Wait for running the Deployment’s Pods.

```console
$ kubectl get pods -n demo --watch
NAME                                   READY     STATUS              RESTARTS   AGE
prometheus-operator-857455484c-dg4qg   0/1       ContainerCreating   0          34s
prometheus-operator-857455484c-dg4qg   1/1       Running             0         45s
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
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/monitoring/coreos-operator/demo-1.yaml
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
default               1         3m
prometheus            1         1m
prometheus-operator   1         3m
```

### Prometheus Dashboard

Now to open prometheus dashboard on Browser:

```console
$ kubectl get svc -n demo
NAME                  TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
prometheus            LoadBalancer   10.100.197.251   <pending>     9090:30900/TCP   1m
prometheus-operated   ClusterIP      None             <none>        9090/TCP         1m

$ minikube ip
192.168.99.100

$ minikube service prometheus -n demo --url
http://192.168.99.100:30900
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{prometheus-svc-nodeport}_ to visit Prometheus Dashboard. According to the above example, this URL will be [http://192.168.99.100:30900](http://192.168.99.100:30900).

If you are not using minikube, browse prometheus dashboard using following address `http://{Node's ExternalIP}:{NodePort of prometheus-service}`.

## Create a MySQL database

KubeDB implements a `MySQL` CRD to define the specification of a MySQL database. Below is the `MySQL` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-mon-coreos
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
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
```

The `MySQL` CRD object contains `monitor` field in it's `spec`.  It is also possible to add CoreOS-Prometheus monitor to an existing `MySQL` database by adding the below part in it's `spec` field.

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

Run the following command to deploy the above `MySQL` CRD object.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/mysql/monitoring/coreos-operator/demo-1.yaml
mysql.kubedb.com/mysql-mon-coreos created
```

Here,

- `spec.monitor` specifies that CoreOS Prometheus operator is used to monitor this database instance. A ServiceMonitor should be created in the `demo` namespace with label `app=kubedb`. The exporter endpoint should be scrapped every 10 seconds.

KubeDB will create a separate stats service with name `<mysql-crd-name>-stats` for monitoring purpose. KubeDB operator will configure this monitoring service once the MySQL is successfully running.

```console
$ kubedb get my -n demo
NAME               VERSION   STATUS     AGE
mysql-mon-coreos   8.0-v1    Creating   22s

$ kubedb describe my -n demo mysql-mon-coreos
Name:               mysql-mon-coreos
Namespace:          demo
CreationTimestamp:  Thu, 27 Sep 2018 16:29:36 +0600
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
  Name:               mysql-mon-coreos
  CreationTimestamp:  Thu, 27 Sep 2018 16:29:39 +0600
  Labels:               kubedb.com/kind=MySQL
                        kubedb.com/name=mysql-mon-coreos
  Annotations:        <none>
  Replicas:           824640215820 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         mysql-mon-coreos
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-mon-coreos
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.97.243.29
  Port:         db  3306/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.7:3306

Service:        
  Name:         mysql-mon-coreos-stats
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-mon-coreos
  Annotations:    monitoring.appscode.com/agent=prometheus.io/coreos-operator
  Type:         ClusterIP
  IP:           10.109.38.68
  Port:         prom-http  56790/TCP
  TargetPort:   prom-http/TCP
  Endpoints:    172.17.0.7:56790

Database Secret:
  Name:         mysql-mon-coreos-auth
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-mon-coreos
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  password:  16 bytes
  user:      4 bytes

Monitoring System:
  Agent:  prometheus.io/coreos-operator
  Prometheus:
    Port:       56790
    Namespace:  demo
    Labels:     app=kubedb
    Interval:   10s

No Snapshots.

Events:
  Type    Reason      Age   From            Message
  ----    ------      ----  ----            -------
  Normal  Successful  1m    MySQL operator  Successfully created Service
  Normal  Successful  1m    MySQL operator  Successfully created StatefulSet
  Normal  Successful  1m    MySQL operator  Successfully created MySQL
  Normal  Successful  56s   MySQL operator  Successfully created stats service
  Normal  Successful  52s   MySQL operator  Successfully patched StatefulSet
  Normal  Successful  52s   MySQL operator  Successfully patched MySQL
  Normal  Successful  51s   MySQL operator  Successfully patched StatefulSet
  Normal  Successful  51s   MySQL operator  Successfully patched MySQL
```

Since `spec.monitoring` was configured, a ServiceMonitor object is created accordingly. You can verify it running the following commands:

```yaml
$ kubectl get servicemonitor -n demo
NAME                           AGE
kubedb-demo-mysql-mon-coreos   1m

$ kubectl get servicemonitor -n demo kubedb-demo-mysql-mon-coreos -o yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  creationTimestamp: 2018-09-27T10:30:20Z
  generation: 1
  labels:
    app: kubedb
    monitoring.appscode.com/service: mysql-mon-coreos-stats.demo
  name: kubedb-demo-mysql-mon-coreos
  namespace: demo
  resourceVersion: "6257"
  selfLink: /apis/monitoring.coreos.com/v1/namespaces/demo/servicemonitors/kubedb-demo-mysql-mon-coreos
  uid: 55a3ae53-c240-11e8-b2cc-080027d9f35e
spec:
  endpoints:
  - interval: 10s
    path: /kubedb.com/v1alpha1/namespaces/demo/mysqls/mysql-mon-coreos/metrics
    port: prom-http
    targetPort: 0
  namespaceSelector:
    matchNames:
    - demo
  selector:
    matchLabels:
      kubedb.com/kind: MySQL
      kubedb.com/name: mysql-mon-coreos
```

Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.
![prometheus-coreos](/docs/images/mysql/mysql-coreos.png)

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mysql/mysql-mon-coreos -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mysql/mysql-mon-coreos

kubectl patch -n demo drmn/mysql-mon-coreos -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mysql-mon-coreos

kubectl delete -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/monitoring/coreos-operator/demo-1.yaml
kubectl delete -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/monitoring/coreos-operator/demo-0.yaml

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
