---
title: Monitor Memcached using Coreos Prometheus Operator
menu:
  docs_0.9.0-rc.0:
    identifier: mc-using-coreos-prometheus-operator-monitoring
    name: Coreos Prometheus Operator
    parent: mc-monitoring-memcached
    weight: 15
menu_name: docs_0.9.0-rc.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Prometheus (CoreOS operator) with KubeDB

This tutorial will show you how to monitor KubeDB databases using Prometheus via [CoreOS Prometheus Operator](https://github.com/coreos/prometheus-operator).

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

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
$ kubectl get pods -n demo
NAME                                   READY     STATUS    RESTARTS   AGE
prometheus-operator-857455484c-q4qlr   1/1       Running   0          21s
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

## Monitor Memcached with CoreOS Prometheus

KubeDB implements a `Memcached` CRD to define the specification of a Memcached database. Below is the `Memcached` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  name: memcd-mon-coreos
  namespace: demo
spec:
  replicas: 3
  version: "1.5.4-v1"
  podTemplate:
    spec:
      resources:
        limits:
          cpu: 500m
          memory: 128Mi
        requests:
          cpu: 250m
          memory: 64Mi
  monitor:
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
```

The `Memcached` CRD object contains `monitor` field in it's `spec`.  It is also possible to add CoreOS-Prometheus monitor to an existing `Memcached` database by adding the below part in it's `spec` field.

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

|  Keys  |  Value |  Description |
|--------|--------|--------------|
| `spec.monitor.agent` | string | `Required`. Indicates the monitoring agent used. Only valid value currently is `coreos-prometheus-operator` |
| `spec.monitor.prometheus.namespace` | string | `Required`. Indicates namespace where service monitors are created. This must be the same namespace of the Prometheus instance. |
| `spec.monitor.prometheus.labels` | map | `Required`. Indicates labels applied to service monitor.                                                    |
| `spec.monitor.prometheus.interval` | string | `Optional`. Indicates the scrape interval for database exporter endpoint (eg, '10s')                        |
| `spec.monitor.prometheus.port` | int |`Optional`. Indicates the port for database exporter endpoint (default is `56790`)|

__Known Limitations:__ If the database password is updated, exporter must be restarted to use the new credentials. This issue is tracked [here](https://github.com/kubedb/project/issues/53).

Run the following command to deploy the above `Memcached` CRD object.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/memcached/monitoring/coreos-operator/demo-1.yaml
memcached.kubedb.com/memcd-mon-coreos created
```

Here,

- `spec.monitor` specifies that CoreOS Prometheus operator is used to monitor this database instance. A ServiceMonitor should be created in the `demo` namespace with label `app=kubedb`. The exporter endpoint should be scrapped every 10 seconds.

KubeDB will create a separate stats service with name `<memcached-crd-name>-stats` for monitoring purpose. KubeDB operator will configure this monitoring service once the Memcached is successfully running.

```console
$ kubedb get mc -n demo
NAME               VERSION    STATUS    AGE
memcd-mon-coreos   1.5.4-v1   Running   1m

$ kubedb describe mc -n demo memcd-mon-coreos
Name:               memcd-mon-coreos
Namespace:          demo
CreationTimestamp:  Wed, 03 Oct 2018 16:46:01 +0600
Labels:             <none>
Annotations:        <none>
Replicas:           3  total
Status:             Running

Deployment:
  Name:               memcd-mon-coreos
  CreationTimestamp:  Wed, 03 Oct 2018 16:46:03 +0600
  Labels:               kubedb.com/kind=Memcached
                        kubedb.com/name=memcd-mon-coreos
  Annotations:          deployment.kubernetes.io/revision=1
  Replicas:           3 desired | 3 updated | 3 total | 3 available | 0 unavailable
  Pods Status:        3 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         memcd-mon-coreos
  Labels:         kubedb.com/kind=Memcached
                  kubedb.com/name=memcd-mon-coreos
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.101.172.16
  Port:         db  11211/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.6:11211,172.17.0.7:11211,172.17.0.8:11211

Service:
  Name:         memcd-mon-coreos-stats
  Labels:         kubedb.com/kind=Memcached
                  kubedb.com/name=memcd-mon-coreos
  Annotations:    monitoring.appscode.com/agent=prometheus.io/coreos-operator
  Type:         ClusterIP
  IP:           10.99.230.104
  Port:         prom-http  56790/TCP
  TargetPort:   prom-http/TCP
  Endpoints:    172.17.0.6:56790,172.17.0.7:56790,172.17.0.8:56790

Monitoring System:
  Agent:  prometheus.io/coreos-operator
  Prometheus:
    Port:       56790
    Namespace:  demo
    Labels:     app=kubedb
    Interval:   10s

No Snapshots.

Events:
  Type    Reason      Age   From                Message
  ----    ------      ----  ----                -------
  Normal  Successful  2m    Memcached operator  Successfully created Service
  Normal  Successful  1m    Memcached operator  Successfully created StatefulSet
  Normal  Successful  1m    Memcached operator  Successfully created Memcached
  Normal  Successful  1m    Memcached operator  Successfully created stats service
  Normal  Successful  1m    Memcached operator  Successfully patched StatefulSet
  Normal  Successful  1m    Memcached operator  Successfully patched Memcached
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
  creationTimestamp: 2018-10-03T10:46:48Z
  generation: 1
  labels:
    app: kubedb
    monitoring.appscode.com/service: memcd-mon-coreos-stats.demo
  name: kubedb-demo-memcd-mon-coreos
  namespace: demo
  resourceVersion: "28890"
  selfLink: /apis/monitoring.coreos.com/v1/namespaces/demo/servicemonitors/kubedb-demo-memcd-mon-coreos
  uid: a0b48fd2-c6f9-11e8-8ebc-0800275bbbee
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
      kubedb.com/kind: Memcached
      kubedb.com/name: memcd-mon-coreos
```

Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.
![prometheus-coreos](/docs/images/memcached/memcached-coreos.png)

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mc/memcd-mon-coreos -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mc/memcd-mon-coreos

kubectl patch -n demo drmn/memcd-mon-coreos -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/memcd-mon-coreos

kubectl delete -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/monitoring/coreos-operator/demo-1.yaml
kubectl delete -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/docs/examples/monitoring/coreos-operator/demo-0.yaml

kubectl delete ns demo
```

## Next Steps

- Monitor your Memcached database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/memcached/monitoring/using-builtin-prometheus.md).
- Detail concepts of [Memcached object](/docs/concepts/databases/memcached.md).
- Use [private Docker registry](/docs/guides/memcached/private-registry/using-private-registry.md) to deploy Memcached with KubeDB.
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
