---
title: Monitoring Single Elasticsearch Cluster in KubeDB
menu:
  docs_0.8.0:
    identifier: es-x-pack-monitoring-single-cluster
    name: Monitoring Single Cluster
    parent: es-x-pack
    weight: 20
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Kibana with KubeDB Elasticsearch

This tutorial will show you how to use X-Pack monitoring in an Elasticsearch cluster in KubeDB.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

As KubeDB uses [Search Guard](https://search-guard.com/) plugin for authentication and authorization, you have to know how to configure Search Guard for both Elasticsearch cluster and Kibana. If you don't know, please visit [here](https://docs.search-guard.com/latest/main-concepts).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: Yaml files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

## Overview

At first, we will create some necessary Search Guard configuration and roles to give an user permission to monitor an Elasticsearch cluster from Kibana. We will create a secret with this configuration files. Then we will provide this secret in `spec.databaseSecret` field of Elasticsearch crd so that our Elasticsearch cluster start with this configuration. We will also configure Elasticsearch cluster with a [custom configuration](/docs/guides/elasticsearch/custom-config/overview.md) file.

Then, we will deploy Kibana with Search Guard plugin installed. We will configure Kibana to connect with our Elasticsearch cluster and view monitoring data from it.

For this tutorial, we will use Elasticsearch 6.3.0 with Search Guard plugin 23.0 and Kibana 6.3.0 with Search Guard plugin 14 installed.

## Deploy Elasticsearch Cluster

Let's create necessary Search Guard configuration files. Here, we will create two users `admin` and `monitor`. User `admin` will have all permissions on the cluster and user `monitor` will have some limited permission to view only monitoring data. Here, are the contents of Search Guard configuration files,

**sg_action_groups.yml:**

```yaml
###### UNLIMITED ######
UNLIMITED:
  readonly: true
  permissions:
    - "*"

###### CLUSTER LEVEL #####
CLUSTER_MONITOR:
  readonly: true
  permissions:
    - "cluster:monitor/*"

CLUSTER_COMPOSITE_OPS_RO:
  readonly: true
  permissions:
    - "indices:data/read/mget"
    - "indices:data/read/msearch"
    - "indices:data/read/mtv"
    - "indices:data/read/coordinate-msearch*"
    - "indices:admin/aliases/exists*"
    - "indices:admin/aliases/get*"
    - "indices:data/read/scroll"

CLUSTER_COMPOSITE_OPS:
  readonly: true
  permissions:
    - "indices:data/write/bulk"
    - "indices:admin/aliases*"
    - "indices:data/write/reindex"
    - CLUSTER_COMPOSITE_OPS_RO

###### INDEX LEVEL ######
INDICES_ALL:
  readonly: true
  permissions:
    - "indices:*"

READ:
  readonly: true
  permissions:
    - "indices:data/read*"
    - "indices:admin/mappings/fields/get*"
    - "indices:admin/mappings/get*"
```

**sg_roles.yaml:**

```yaml
### Admin
sg_all_access:
  readonly: true
  cluster:
    - UNLIMITED
  indices:
    '*':
      '*':
        - UNLIMITED
  tenants:
    admin_tenant: RW

### X-Pack Monitoring
sg_xp_monitoring:
  cluster:
    - cluster:admin/xpack/monitoring/*
    - cluster:admin/ingest/pipeline/put       
    - cluster:admin/ingest/pipeline/get
    - indices:admin/template/get
    - indices:admin/template/put
    - CLUSTER_MONITOR
    - CLUSTER_COMPOSITE_OPS
  indices:
    '?monitor*':
      '*':
        - INDICES_ALL
    '?marvel*':
      '*':
        - INDICES_ALL
    '?kibana*':
      '*':
        - READ
    '*':
      '*':
        - indices:data/read/field_caps
```

**sg_internal_users.yml:**

```yaml
#password is: admin@secret
admin:
  readonly: true
  hash: $2y$12$skma87wuFFtxtGWegeAiIeTtUH1nnOfIRZzwwhBlzXjg0DdM4gLeG
  roles:
    - admin

#password is: monitor@secret
monitor:
  readonly: true
  hash: $2y$12$JDTXih3AqV/1MDRYQ.KIY.u68CkzCIq.xiiqwtRJx3cjN0YmFavTe
  roles:
    - monitor
```

Here, we have used `admin@secret` password for `admin` user and  `monitor@secret` password for `monitor` user.

**sg_roles_mapping.yml:**

```yaml
sg_all_access:
  readonly: true
  backendroles:
    - admin

sg_xp_monitoring:
  readonly: true
  backendroles:
    - monitor
```

**sg_config.yml:**

```yaml
searchguard:
  dynamic:
    authc:
      kibana_auth_domain:
        enabled: true
        order: 0
        http_authenticator:
          type: basic
          challenge: false
        authentication_backend:
          type: internal
      basic_internal_auth_domain: 
        http_enabled: true
        transport_enabled: true
        order: 1
        http_authenticator:
          type: basic
          challenge: true
        authentication_backend:
          type: internal
```

Now, create a secret with these Search Guard configuration files.

```console
 $ kubectl create secret generic -n demo  es-auth \
                        --from-file=./sg_action_groups.yml \
                        --from-file=./sg_config.yml \
                        --from-file=./sg_internal_users.yml \
                        --from-file=./sg_roles_mapping.yml \
                        --from-file=./sg_roles.yml
secret/es-auth created
```

Verify the secret has desired configuration files,

```yaml
$ kubectl get secret -n demo es-auth -o yaml
apiVersion: v1
data:
  sg_action_groups.yml: <base64 encoded content>
  sg_config.yml: <base64 encoded content>
  sg_internal_users.yml: <base64 encoded content>
  sg_roles.yml: <base64 encoded content>
  sg_roles_mapping.yml: <base64 encoded content>
kind: Secret
metadata:
  ...
  name: es-auth
  namespace: demo
  ...
type: Opaque
```

As we are using Search Guard plugin for authentication, we need to ensure that `x-pack` security is not enabled. We will ensure that by providing `xpack.security.enabled: false` in `common-config.yml` file and we will use this file to configure our Elasticsearch cluster. As Search Guard does not support `local` exporter, we will use `http` exporter and set `host` filed to `http://127.0.0.1:9200` to store monitoring data in same cluster.

 Let's create `common-config.yml` with following configuration,

```yaml
xpack.security.enabled: false
xpack.monitoring.enabled: true
xpack.monitoring.collection.enabled: true
xpack.monitoring.exporters:
  my-http-exporter:
    type: http
    host: ["http://127.0.0.1:9200"]
    auth:
      username: monitor
      password: monitor@secret
```

Create a ConfigMap using this file,

```console
$ kubectl create configmap -n demo es-custom-config \
                        --from-file=./common-config.yaml
configmap/es-custom-config created
```

Verify that the ConfigMap has desired configuration,

```yaml
$  kubectl get configmap -n demo es-custom-config -o yaml
apiVersion: v1
data:
  common-config.yaml: |-
    xpack.security.enabled: false
    xpack.monitoring.enabled: true
    xpack.monitoring.collection.enabled: true
    xpack.monitoring.exporters:
      my-http-exporter:
        type: http
        host: ["http://127.0.0.1:9200"]
        auth:
          username: monitor
          password: monitor@secret
kind: ConfigMap
metadata:
  ...
  name: es-custom-config
  namespace: demo
  ...
```

Now, create Elasticsearch crd specifying  `spec.databaseSecret` and `spec.configSource` field.

```console
$ kubectl apply -f kubectl apply -f ./es-mon-demo.yaml 
elasticsearch.kubedb.com/es-mon-demo created
```

Below is the YAML for the Elasticsearch crd we just created.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: es-mon-demo
  namespace: demo
spec:
  version: "6.3.0"
  doNotPause: false
  replicas: 1
  databaseSecret:
    secretName: es-auth
  configSource:
    configMap:
      name: es-custom-config
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
```

Now, wait for few minutes. KubeDB will create necessary secrets, services, and statefulsets.

Check resources created in demo namespace by KubeDB,

```console
$  kubectl get all -n demo
NAME                READY     STATUS    RESTARTS   AGE
pod/es-mon-demo-0   1/1       Running   0          54s

NAME                         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/es-mon-demo          ClusterIP   10.105.183.97   <none>        9200/TCP   1m
service/es-mon-demo-master   ClusterIP   10.102.108.64   <none>        9300/TCP   59s
service/kubedb               ClusterIP   None            <none>        <none>     13m

NAME                           DESIRED   CURRENT   AGE
statefulset.apps/es-mon-demo   1         1         54s
```

Once everything is created, Elasticsearch will go to Running state. Check that Elasticsearch is in running state.

```console
$ kubectl get es -n demo es-mon-demo
NAME          VERSION   STATUS    AGE
es-mon-demo   6.3.0     Running   1m
```

Now, check elasticsearch log to see if the cluster is ready to accept requests,

```console
$ kubectl logs -n demo es-mon-demo-0 -f
...
Starting runit...
...
Elasticsearch Version: 6.3.0
Search Guard Version: 6.3.0-23.0
Connected as CN=sgadmin,O=Elasticsearch Operator
Contacting elasticsearch cluster 'elasticsearch' and wait for YELLOW clusterstate ...
Clustername: es-mon-demo
Clusterstate: GREEN
Number of nodes: 1
Number of data nodes: 1
...
Done with success
...
```

Once you see `Done with success` success line in the log, the cluster is ready to accept requests. Now, it is time to connect with Kibana.

## Deploy Kibana

In order to view monitoring data from Kibana, we need to configure `kibana.yml` with appropriate configuration.

KubeDB has created a service with name `es-mon-demo` in `demo` namespace for the Elasticsearch cluster. We will use this service in `elasticsearch.url` field. Kibana will use this service to connect with the Elasticsearch cluster.

Let's, configure `kibana.yml` as below,

```yaml
xpack.security.enabled: false
xpack.monitoring.enabled: true
xpack.monitoring.kibana.collection.enabled: true
xpack.monitoring.ui.enabled: true

server.host: 0.0.0.0

elasticsearch.url: "http://es-mon-demo.demo.svc"
elasticsearch.username: "monitor"
elasticsearch.password: "monitor@secret"

searchguard.auth.type: "basicauth"
searchguard.cookie.secure: false

```

Notice the `elasticsearch.username` and `elasticsearch.password` field. Kibana will connect to Elasticsearch cluster with this credentials. They must match with the credentials we have provided in `sg_internal_users.yml` file for `monitor` user while creating the cluster.

Now, create a ConfigMap with `kibana.yml` file. We will mount this ConfigMap in Kibana deployment so that Kibana starts with this configuration.

```conlose
$ kubectl create configmap -n demo kibana-config \
                        --from-file=./kibana.yml
configmap/kibana-config created
```

Finally, deploy Kibana deployment,
```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/elasticsearch/x-pack/kibana-deployment.yaml
deployment.apps/kibana created
```

Below is the YAML for the Kibana deployment we just created.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kibana
  namespace: demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kibana
  template:
    metadata:
      labels:
        app: kibana
    spec:
      containers:
      - name: kibana
        image: kubedbci/kibana:6.3.0
        volumeMounts:
        - name:  kibana-config
          mountPath: /usr/share/kibana/config
      volumes:
      - name:  kibana-config
        configMap:
          name: kibana-config
```

Now, wait for few minutes. Let the Kibana pod  go in`Running` state. Check pod is in `Running` using this command,

```console
 $ kubectl get pods -n demo -l app=kibana
NAME                      READY     STATUS    RESTARTS   AGE
kibana-84b8cbcf7c-mg699   1/1       Running   0          3m
```

Now, watch the Kibana pod's log to see if Kibana is ready to access,

```console
$ kubectl logs -n demo kibana-84b8cbcf7c-mg699 -f
...
{"type":"log","@timestamp":"2018-08-27T09:50:47Z","tags":["listening","info"],"pid":1,"message":"Server running at http://0.0.0.0:5601"}
```

Once you see `"message":"Server running at http://0.0.0.0:5601"` in the log, Kibana is ready. Now it is time to access Kibana UI.

In order to access Kibana UI from outside of the cluster, we need to create a service with type `NodePort`.
```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/elasticsearch/x-pack/kibana-service.yaml
service/kibana created
```

Below is the YAML for the service we just created.

```yaml
kind: Service
apiVersion: v1
metadata:
  name:  kibana
  namespace: demo
spec:
  selector:
    app:  kibana
  type:  NodePort
  ports:
  - name:  http
    protocol: TCP
    port:  5601
    targetPort: 5601
```

Now, we can access the Kibana UI using `Cluster's IP address` and the `NodePort` of the service.

```console
$ kubectl get service -n demo kibana 
NAME      TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
kibana    NodePort   10.102.33.55   <none>        5601:32473/TCP   1h
```

Here, `32473` is the desired `NodePort`. Now you can open `http://{cluster-ip-address:32473}` in your browser to access Kibana UI.

For minikube, we can get the url by,

```console
$ minikube service -n demo kibana --url
http://192.168.99.100:32473
```

When you will open the address, you will be greeted with Search Guard login UI. Login with following credentials: `username: monitor` and `password: monitor@secret`.

After login, go to `Monitoring` tab in Kibana UI. You will see Kibana has connected with the Elasticsearch cluster and showing monitoring data. Some screenshots of monitoring `es-mon-demo` cluster is given below.

![Kibana Monitoring Home](/docs/images/elasticsearch/x-pack/monitoring-home.png)

![Kibana Monitoring Node](/docs/images/elasticsearch/x-pack/monitoring-node.png)

![Kibana Monitoring Overview](/docs/images/elasticsearch/x-pack/monitoring-overview.png)

## Cleanup

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/es-mon-demo -p '{"spec":{"doNotPause":false}}' --type="merge"

$ kubectl delete -n demo es/es-mon-demo

$ kubectl patch -n demo drmn/es-mon-demo -p '{"spec":{"wipeOut":true}}' --type="merge"

$ kubectl delete -n demo drmn/es-mon-demo

$ kubectl delete  -n demo configmap/es-custom-config

$ kubectl delete -n demo secret/es-auth

$ kubectl delete -n demo configmap/kibana-config

$ kubectl delete -n demo deployment/kibana

$ kubectl delete -n demo service/kibana

$ kubectl delete ns demo
```

To uninstall KubeDB follow this [guide](/docs/setup/uninstall.md).