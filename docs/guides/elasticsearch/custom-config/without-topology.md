---
title: Using Custom Configuration in Elasticsearch without Topology
menu:
  docs_0.8.0:
    identifier: es-custom-config-without-topology
    name: Without Topology
    parent: es-custom-config
    weight: 30
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Custom Configuration in Elasticsearch without Topology

This tutorial will show you how to use custom configuration in an Elasticsearch cluster in KubeDB without specifying `spec.topology` field.

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

At first, let's create two configuration files namely `master-config.yml`, `data-config.yml` and `common-config.yalm` respectively.

Content of `master-config.yml`,

```yaml
node:
  name:  es-node-master
path:
  data: ["/data/elasticsearch/master-datadir"]
```

Content of `data-config.yml`,

```yaml
node:
  name:  es-node-data
path:
  data: ["/data/elasticsearch/data-datadir"]
http:
  compression: false
```

Content of `common-config.yml`,

```yaml
path:
  logs: /data/elasticsearch/common-logdir
```

This time we have added an additional field `http.compression: false` in `data-config.yml` file. By default, this field is set to `true`.

Now, let's create a configMap with these configuration files,

```console
 $ kubectl create configmap  -n demo es-custom-config \
                        --from-file=./common-config.yml \
                        --from-file=./master-config.yml \
                        --from-file=./data-config.yml
configmap/es-custom-config created
```

Check that the configMap has these configuration files,

```console
$ kubectl get configmap -n demo es-custom-config -o yaml
apiVersion: v1
data:
  common-config.yml: |
    path:
      logs: /data/elasticsearch/common-logdir
  data-config.yml: |-
    node:
      name:  es-node-data
    path:
      data: ["/data/elasticsearch/data-datadir"]
    http:
      compression: false
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

Now, create an Elasticsearch crd without topology,

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/elasticsearch/custom-config/es-custom.yaml
elasticsearch.kubedb.com/custom-elasticsearch created
```

Below is the YAML for the Elasticsearch crd we just created.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: custom-elasticsearch
  namespace: demo
spec:
  version: "6.2.4"
  replicas: 2
  doNotPause: true
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

Check resources created in `demo` namespace by KubeDB,

```console
$ kubectl get all -n demo
NAME                         READY     STATUS    RESTARTS   AGE
pod/custom-elasticsearch-0   1/1       Running   0          1m
pod/custom-elasticsearch-1   1/1       Running   0          48s

NAME                                  TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/custom-elasticsearch          ClusterIP   10.97.175.61    <none>        9200/TCP   1m
service/custom-elasticsearch-master   ClusterIP   10.100.153.28   <none>        9300/TCP   1m
service/kubedb                        ClusterIP   None            <none>        <none>     1m

NAME                                    DESIRED   CURRENT   AGE
statefulset.apps/custom-elasticsearch   2         2         1m
```

Check secrets created by KubeDB,

```console
$ kubectl get secret -n demo
NAME                        TYPE                                  DATA      AGE
custom-elasticsearch-auth   Opaque                                7         2m
custom-elasticsearch-cert   Opaque                                4         2m
default-token-58ddc         kubernetes.io/service-account-token   3         2m
```

Once everything is created, Elasticsearch will go to `Running` state. Check that Elasticsearch is in running state.

```console
$ kubectl get es -n demo custom-elasticsearch
NAME                   VERSION   STATUS    AGE
custom-elasticsearch   6.2.4     Running   2m
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
NAME                TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
custom-es-exposed   NodePort   10.96.7.103   <none>        9200:32413/TCP   19s
```

To connect with the Elasticsearch cluster we have to use the NodePort of the service along with the cluster's IP address.

For minikube, we can get the url by,

```console
$ minikube service custom-es-exposed -n demo --url
http://192.168.99.100:32413
```

Run the following command to get `admin` user password

```console
$ kubectl get secret -n demo custom-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
ieeu57oj
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
        "total": 2,
        "successful": 2,
        "failed": 0
    },
    "cluster_name": "custom-elasticsearch",
    "nodes": {
        "qhH6AJ9JTU6SlYHKF3kwOQ": {
            "name": "es-node-master",
            ...
            "roles": [
                "master",
                "data",
                "ingest"
            ],
            "settings": {
                ...
                "node": {
                    "name": "es-node-master",
                    "data": "true",
                    "ingest": "true",
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
                "http": {
                    ...
                    "compression": "false",
                    ...
                },
            }
        },
        "mNgVB7i6RkOPtRmXzjqPFA": {
            "name": "es-node-master",
            ...
            "roles": [
                "master",
                "data",
                "ingest"
            ],
            "settings": {
               ...
                "node": {
                    "name": "es-node-master",
                    "data": "true",
                    "ingest": "true",
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
                "http": {
                    ...
                    "compression": "false",
                    ...
                },
               ...
            }
        }
    }
}
```

We have total two nodes in our Elasticsearch cluster. Here, we have an array of these node's settings information. Here, `"roles"` field represents if the node is working as either a master, ingest/client or data node. We can see from the response that all the nodes are working as master, ingest/client and data node simultaneously.

From the response above, we can see that `"node.name"` and `"path.data"` keys are set to the value we have specified in `master-config.yml` file. Note that, we had also specified these keys in `data-config.yml` file but they were overridden by the value in `master-config.yml` file. This happened because config values in `master-config.yml` file has higher precedence than the `data-config.yml` file. However, note that `"http.compress"` field has been applied from `data-config.yml` file as `master-config.yml` file does not have this field.

Also note that, the `"path.logs"` field of each node is set to the value we have specified in `common-config.yml` file.

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
