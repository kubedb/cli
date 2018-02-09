> New to KubeDB Postgres?  Quick start [here](/docs/guides/postgres/quickstart.md).

# Using Prometheus (CoreOS operator) with KubeDB

This tutorial will show you how to monitor PostgreSQL database using [Prometheus](https://prometheus.io/).

## Monitor with builtin Prometheus

Below is the Postgres object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: builtin-prom-postgres
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
    agent: prometheus.io/builtin
```

Here,

 - `spec.monitor` specifies that built-in [prometheus](https://github.com/prometheus/prometheus) is used to monitor this database instance.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/postgres/monitoring/builtin-prom-postgres.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/postgres/monitoring/builtin-prom-postgres.yaml"
postgres "builtin-prom-postgres" created
```

KubeDB operator will configure the service of this database in a way that the Prometheus server will automatically find out the service endpoint aka `Postgres Exporter` and
will receive metrics from exporter.

You can verify it running the following commands:

```yaml
$ kubectl get svc -n demo --selector="kubedb.com/name=builtin-prom-postgres"
NAME                            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
builtin-prom-postgres           ClusterIP   10.107.124.174   <none>        56790/TCP   16s
builtin-prom-postgres-primary   ClusterIP   10.96.250.3      <none>        5432/TCP    16s
```

Lets describe Service `builtin-prom-postgres`

```console
$ kubectl describe svc -n demo builtin-prom-postgres
Name:              builtin-prom-postgres
Namespace:         demo
Labels:            kubedb.com/kind=Postgres
                   kubedb.com/name=builtin-prom-postgres
Annotations:       monitoring.appscode.com/agent=prometheus.io/builtin
                   prometheus.io/path=/kubedb.com/v1alpha1/namespaces/demo/postgreses/builtin-prom-postgres/metrics
                   prometheus.io/port=56790
                   prometheus.io/scrape=true
Selector:          kubedb.com/kind=Postgres,kubedb.com/name=builtin-prom-postgres
Type:              ClusterIP
IP:                10.107.124.174
Port:              prom-http  56790/TCP
TargetPort:        %!d(string=prom-http)/TCP
Endpoints:         172.17.0.8:56790
Session Affinity:  None
```

You can see that the service contains following annotations. The prometheus server will discover the exporter using these specifications.

```console
prometheus.io/path=/kubedb.com/v1alpha1/namespaces/demo/postgreses/builtin-prom-postgres/metrics
prometheus.io/port=56790
prometheus.io/scrape=true
```

## Deploy and configure Prometheus Server

The prometheus server is needed to configure so that it can discover endpoints of services. If a Prometheus server is already running in cluster
and if it is configured in a way that it can discover service endpoints, no extra configuration will be needed.

If there is no existing Prometheus server running, rest of this tutorial will create a Prometheus server with appropriate configuration.

The configuration file of Prometheus server will be provided by ConfigMap. The below config map will be created
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
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/monitoring/builtin-prometheus/demo-1.yaml
configmap "prometheus-server-conf" created
```

Now, the below YAML is used to deploy Prometheus in kubernetes :

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

#### In RBAC enabled cluster
If RBAC *is* enabled, Run the following command to deploy prometheus in kubernetes

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/monitoring/builtin-prometheus/rbac/demo-2.yaml
clusterrole "prometheus-server" created
serviceaccount "prometheus-server" created
clusterrolebinding "prometheus-server" created
deployment "prometheus-server" created
service "prometheus-service" created
```

Verify RBAC stuffs

```console
$ kubectl get clusterrole,clusterrolebinding,sa prometheus-server -n demo
NAME                             AGE
clusterroles/prometheus-server   1m

NAME                                    AGE
clusterrolebindings/prometheus-server   1m

NAME                   SECRETS   AGE
sa/prometheus-server   1         1m
```


#### In RBAC \*not\* enabled cluster

If RBAC *is not* enabled, Run the following command to deploy prometheus in kubernetes

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/monitoring/builtin-prometheus/demo-2.yaml
deployment "prometheus-server" created
service "prometheus-service" created
```

#### Prometheus Dashboard

Now open prometheus dashboard on browser by running `minikube service prometheus-service -n demo`.

Or you can get the URL of `prometheus-service` Service by running following command

```console
$ minikube service prometheus-service -n demo --url
http://192.168.99.100:30901
```

Now, if you go the Prometheus Dashboard, you should see that this database endpoint as one of the targets.

<p align="center">
  <kbd>
    <img alt="builtin-prom-postgres"  src="/docs/images/postgres/builtin-prom-postgres.png">
  </kbd>
</p>

## Next Steps
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using_coreos_prometheus_operator.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
