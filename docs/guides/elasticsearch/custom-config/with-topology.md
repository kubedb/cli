---
title: Using Custom Configuration in Elasticsearch with Topology
menu:
  docs_0.8.0:
    identifier: es-custom-config-with-topology
    name: With Topology
    parent: es-custom-config
    weight: 20
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Custom Configuration in Elasticsearch with Topology

This tutorial will show you how to use custom configuration in an Elasticsearch cluster in KubeDB specifying `spec.topology` field.

If you don't know how KubeDB handles custom configuration for an Elasticsearch cluster, please visit [here](/docs/guides/elasticsearch/custom-config/overview.md).

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: Yaml files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

## Use Custom Configuration

At first, let's create four configuration files namely `master-config.yml`, `client-config.yml`, `data-config.yml` and `common-config.yalm`.

Content of `master-config.yml`,

```yaml
node:
  name:  es-node-master
path:
  data: ["/data/elasticsearch/master-datadir"]
```

Content of `client-config.yml`,

```yaml
node:
  name:  es-node-client
path:
  data: ["/data/elasticsearch/client-datadir"]
```

Content of `data-config.yml`,

```yaml
node:
  name:  es-node-data
path:
  data: ["/data/elasticsearch/data-datadir"]
```

Content of `common-config.yml`,

```yaml
path:
  logs: /data/elasticsearch/common-logdir
```

Now, let's create a configMap with these configuration files,

```console
 $ kubectl create configmap -n demo es-custom-config \
                        --from-file=./common-config.yml \
                        --from-file=./master-config.yml \
                        --from-file=./data-config.yml \
                        --from-file=./client-config.yml
configmap/es-custom-config created
```

Check that the configMap has these configuration files,

```console
$ kubectl get configmap -n demo es-custom-config -o yaml
apiVersion: v1
data:
  client-config.yml: |-
    node:
      name:  es-node-client
    path:
      data: ["/data/elasticsearch/client-datadir"]
  common-config.yml: |
    path:
      logs: /data/elasticsearch/common-logdir
  data-config.yml: |-
    node:
      name:  es-node-data
    path:
      data: ["/data/elasticsearch/data-datadir"]
  master-config.yml: |-
    node:
      name:  es-node-master
    path:
      data: ["/data/elasticsearch/master-datadir"]
kind: ConfigMap
metadata:
  ...
  name: es-custom-config
  namespace: demo
  ...
```

Now, create an Elasticsearch crd with topology specified,

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/elasticsearch/custom-config/es-custom-with-topology.yaml
elasticsearch.kubedb.com/custom-elasticsearch created
```

Bellow is the YAML for the Elasticsearch crd we just created.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: custom-elasticsearch
  namespace: demo
spec:
  version: "6.2.4"
  doNotPause: true
  configSource:
    configMap:
      name: es-custom-config
  topology:
    master:
      prefix: master
      replicas: 1
      storage:
        storageClassName: "standard"
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
    data:
      prefix: data
      replicas: 1
      storage:
        storageClassName: "standard"
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
    client:
      prefix: client
      replicas: 2
      storage:
        storageClassName: "standard"
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: 50Mi
```

Now, wait for few minutes. KubeDB will create necessary secrets, services, and statefulsets.

Check resources created in `demo` namespace by KubeDB,

```console
$ kubectl get all -n demo
NAME                                READY     STATUS    RESTARTS   AGE
pod/client-custom-elasticsearch-0   1/1       Running   0          9m
pod/client-custom-elasticsearch-1   1/1       Running   0          9m
pod/data-custom-elasticsearch-0     1/1       Running   0          7m
pod/master-custom-elasticsearch-0   1/1       Running   0          8m

NAME                                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/custom-elasticsearch          ClusterIP   10.101.163.110   <none>        9200/TCP   10m
service/custom-elasticsearch-master   ClusterIP   10.103.59.244    <none>        9300/TCP   10m
service/kubedb                        ClusterIP   None             <none>        <none>     10m

NAME                                           DESIRED   CURRENT   AGE
statefulset.apps/client-custom-elasticsearch   2         2         9m
statefulset.apps/data-custom-elasticsearch     1         1         7m
statefulset.apps/master-custom-elasticsearch   1         1         8m
```

Check secrets created by KubeDB,

```console
$ kubectl get secret -n demo
NAME                        TYPE                                  DATA      AGE
custom-elasticsearch-auth   Opaque                                7         10m
custom-elasticsearch-cert   Opaque                                4         10m
default-token-qnf27         kubernetes.io/service-account-token   3         14m
```

Once everything is created, Elasticsearch will go to `Running` state. Check that Elasticsearch is in running state.

```console
$ kubectl get es -n demo custom-elasticsearch
NAME                   VERSION   STATUS    AGE
custom-elasticsearch   6.2.4     Running   14m
```

## Verify Configuration

Now, we will connect with the Elasticsearch cluster we have created. We will query for nodes settings to verify that the cluster is using the custom configuration we have provided.

At first, expose Service `custom-elasticsearch`,

```console
$ kubectl expose svc -n demo custom-elasticsearch --name=custom-es-exposed --port=9200 --protocol=TCP --type=NodePort
service/custom-es-exposed exposed
```

Verify the Service exposed successfully,

```console
$ kubectl get svc -n demo custom-es-exposed
NAME                TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
custom-es-exposed   NodePort   10.99.144.141   <none>        9200:32604/TCP   1m
```

To connect with the Elasticsearch cluster we have to use the NodePort of the service along with the cluster's IP address.

For minikube, we can get the url by,

```console
$ minikube service custom-es-exposed -n demo --url
http://192.168.99.100:32604
```

Run the following command to get `admin` user password

```console
$ kubectl get secret -n demo custom-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
ssmhw6nq⏎
```

Let's export this `url` and `password` for later use,

```ini
$ export es_service=$(minikube service custom-es-exposed -n demo --url)
$ export es_admin_pass=$(kubectl get secrets -n demo custom-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d)
```

Now, we will query for settings of all nodes in an Elasticsearch cluster,

```console
$ curl --user "admin:$es_admin_pass" "$es_service/_nodes/_all/settings"
```

This will return a large JSON with nodes settings information. Here is the prettified JSON response,

```json
{
    "_nodes": {
        "total": 4,
        "successful": 4,
        "failed": 0
    },
    "cluster_name": "custom-elasticsearch",
    "nodes": {
        "fA7g2r7rTV--FZzusuctww": {
            "name": "es-node-client",
            ...
            "roles": [
                "ingest"
            ],
            "settings": {
                ...
                "node": {
                    "name": "es-node-client",
                    "data": "false",
                    "ingest": "true",
                    "master": "false"
                },
                "path": {
                    "data": [
                        "/data/elasticsearch/client-datadir"
                    ],
                    "logs": "/data/elasticsearch/common-logdir",
                    "home": "/elasticsearch"
                },
                ...
            }
        },
        "_8HsT6oZTAGf9Gmz0kInsA": {
            "name": "es-node-client",
            "roles": [
                "ingest"
            ],
            "settings": {
                ...
                "node": {
                    "name": "es-node-client",
                    "data": "false",
                    "ingest": "true",
                    "master": "false"
                },
                "path": {
                    "data": [
                        "/data/elasticsearch/client-datadir"
                    ],
                    "logs": "/data/elasticsearch/common-logdir",
                    "home": "/elasticsearch"
                },
                ...
            }
        },
        "pT1cxPVNQU-UBkjcj6JSzw": {
            "name": "es-node-master",
            ...
            "roles": [
                "master"
            ],
            "settings": {
                ...
                "node": {
                    "name": "es-node-master",
                    "data": "false",
                    "ingest": "false",
                    "master": "true"
                },
                "path": {
                    "data": [
                        "/data/elasticsearch/master-datadir"
                    ],
                    "logs": "/data/elasticsearch/common-logdir",
                    "home": "/elasticsearch"
                },
                ...
            }
        },
        "tBecrUhUTlO9x5kXlPAR5A": {
            "name": "es-node-data",
            ...
            "roles": [
                "data"
            ],
            "settings": {
                ...
                "node": {
                    "name": "es-node-data",
                    "data": "true",
                    "ingest": "false",
                    "master": "false"
                },
                "path": {
                    "data": [
                        "/data/elasticsearch/data-datadir"
                    ],
                    "logs": "/data/elasticsearch/common-logdir",
                    "home": "/elasticsearch"
                },
                ...
            }
        }
    }
}
```

We have total four (1 master + 2 client + 1 data) nodes in our Elasticsearch cluster. Here, we have an array of these node's settings information. Here, `"roles"` field represents if the node is working as either a master, ingest/client or data node.

From the response above, we can see that `"node.name"` and `"path.data"` fields are set according to node rules to the value we have specified in configuration files.

Note that, the `"path.logs"` field of each node is set to the value we have specified in `common-config.yml` file.

## Cleanup

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/custom-elasticsearch -p '{"spec":{"doNotPause":false}}' --type="merge"

$ kubectl delete -n demo es/custom-elasticsearch

$ kubectl patch -n demo drmn/custom-elasticsearch -p '{"spec":{"wipeOut":true}}' --type="merge"

$ kubectl delete -n demo drmn/custom-elasticsearch

$ kubectl delete  -n demo configmap/es-custom-config

$ kubectl delete ns demo
```

To uninstall KubeDB follow this [guide](/docs/setup/uninstall.md).
