---
title: Monitoring
menu:
  docs_0.8.0-beta.0:
    identifier: guides-monitoring
    name: Monitoring
    parent: guides
    weight: 80
menu_name: docs_0.8.0-beta.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/guides/README.md).

# Using Prometheus with KubeDB
This tutorial will show you how to monitor KubeDB databases using Prometheus via [CoreOS Prometheus Operator](https://github.com/coreos/prometheus-operator).

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f ./docs/examples/monitoring/demo-0.yaml 
namespace "demo" created
deployment "prometheus-operator" created

$ kubectl get pods -n demo --watch
NAME                                  READY     STATUS    RESTARTS   AGE
prometheus-operator-449376836-4pkzn   1/1       Running   0          15s
^C⏎                                                                                                                                                             

$ kubectl get crd 
NAME                                    DESCRIPTION                           VERSION(S)
alertmanager.monitoring.coreos.com      Managed Alertmanager cluster          v1alpha1
prometheus.monitoring.coreos.com        Managed Prometheus server             v1alpha1
service-monitor.monitoring.coreos.com   Prometheus monitoring for a service   v1alpha1
```

Once the Prometheus operator CRDs are registered, run the following command to create a Prometheus.

```console
 $ kubectl create -f ./docs/examples/monitoring/demo-1.yaml
prometheus "prometheus" created
service "prometheus" created

$ kubectl get svc -n demo
NAME                  CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
prometheus            10.0.0.247   <pending>     9090:30900/TCP   14s
prometheus-operated   None         <none>        9090/TCP         14s

$ minikube ip
192.168.99.100
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{prometheus-svc-nodeport}_ to visit Prometheus Dashboard. According to the above example, this URL will be [http://192.168.99.100:30900](http://192.168.99.100:30900).

## Create a PostgreSQL database
KubeDB implements a `Postgres` CRD to define the specification of a PostgreSQL database. Below is the `Postgres` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: pmon
  namespace: demo
spec:
  version: 9.5
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  monitor:
    agent: coreos-prometheus-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s

$ kubedb create -f ./docs/examples/monitoring/demo-2.yaml 
validating "./docs/examples/monitoring/demo-2.yaml"
postgres "pmon" created
```

Here,
 - `spec.version` is the version of PostgreSQL database. In this tutorial, a PostgreSQL 9.5 database is going to be created.

 - `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

 - `spec.monitor` specifies that CoreOS Prometheus operator is used to monitor this database instance. A ServiceMonitor should be created in the `demo` namespace with label `app=kubedb`. The exporter endpoint should be scrapped every 10 seconds.

KubeDB operator watches for `Postgres` objects using Kubernetes api. When a `Postgres` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching tpr name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. If [RBAC is enabled](/docs/guides/rbac.md), a ClusterRole, ServiceAccount and ClusterRoleBinding with the matching tpr name will be created and used as the service account name for the corresponding StatefulSet.

```console
$ kubedb get pg -n demo
NAME      STATUS     AGE
pmon      Creating   1m

$ kubedb get pg -n demo
NAME      STATUS    AGE
pmon      Running   1m

$ kubedb describe pg -n demo pmon
Name:		pmon
Namespace:	demo
StartTimestamp:	Mon, 17 Jul 2017 23:46:03 -0700
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

Service:	
  Name:		pmon
  Type:		ClusterIP
  IP:		10.0.0.216
  Port:		db	5432/TCP
  Port:		http	56790/TCP

Database Secret:
  Name:	pmon-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	35 bytes

Monitoring System:
  Agent:	coreos-prometheus-operator
  Prometheus:
    Namespace:	demo
    Labels:	app=kubedb
    Interval:	10s

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                Type       Reason               Message
  ---------   --------   -----     ----                --------   ------               -------
  13s         13s        1         Postgres operator   Normal     SuccessfulValidate   Successfully validate Postgres
  15s         15s        1         Postgres operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  15s         15s        1         Postgres operator   Normal     SuccessfulCreate     Successfully created Postgres
  15s         15s        1         Postgres operator   Normal     SuccessfulCreate     Successfully added monitoring system.
  1m          1m         1         Postgres operator   Normal     SuccessfulValidate   Successfully validate Postgres
  1m          1m         1         Postgres operator   Normal     Creating             Creating Kubernetes objects
```


Since `spec.monitoring` was configured, s ServiceMonitor object is created accordingly. You can verify it running the following commands:

```yaml
$ kubectl get servicemonitor -n demo
NAME               KIND
kubedb-demo-pmon   ServiceMonitor.v1alpha1.monitoring.coreos.com

$ kubectl get servicemonitor -n demo kubedb-demo-pmon -o yaml
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  creationTimestamp: 2017-07-18T06:47:44Z
  labels:
    app: kubedb
  name: kubedb-demo-pmon
  namespace: demo
  resourceVersion: "644"
  selfLink: /apis/monitoring.coreos.com/v1alpha1/namespaces/demo/servicemonitors/kubedb-demo-pmon
  uid: 00fc6ae8-6b85-11e7-aad2-080027036663
spec:
  endpoints:
  - interval: 10s
    path: /kubedb.com/v1alpha1/namespaces/demo/postgreses/pmon/metrics
    port: http
    targetPort: 0
  namespaceSelector:
    matchNames:
    - demo
  selector:
    matchLabels:
      kubedb.com/kind: Postgres
      kubedb.com/name: pmon
```


Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.

![Prometheus Dashboard](/docs/images/monitoring/prometheus.gif)


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall KubeDB operator, please follow the steps [here](/docs/setup/uninstall.md).


## Next Steps
- Learn about the details of monitoring support [here](/docs/concepts/monitoring.md).
- Learn how to use KubeDB to run a PostgreSQL database [here](/docs/guides/postgres/overview.md).
- Learn how to use KubeDB to run an Elasticsearch database [here](/docs/guides/elasticsearch/overview.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md). 
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
