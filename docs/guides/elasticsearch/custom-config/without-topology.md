---
title: Using Custom Configuration in Elasticsearch without Topology
menu:
  docs_0.9.0:
    identifier: es-custom-config-without-topology
    name: Without Topology
    parent: es-custom-config
    weight: 30
menu_name: docs_0.9.0
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
namespace/demo created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: YAML files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

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
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/custom-config/es-custom.yaml
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
  version: "6.2.4-v1"
  replicas: 2
  configSource:
    configMap:
      name: es-custom-config
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

Now, wait for few minutes. KubeDB will create necessary secrets, services, and statefulsets.

Check resources created in `demo` namespace by KubeDB,

```console
$ kubectl get all -n demo -l=kubedb.com/name=custom-elasticsearch
NAME                         READY     STATUS    RESTARTS   AGE
pod/custom-elasticsearch-0   1/1       Running   0          6m
pod/custom-elasticsearch-1   1/1       Running   0          5m

NAME                                  TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
service/custom-elasticsearch          ClusterIP   10.110.0.164   <none>        9200/TCP   6m
service/custom-elasticsearch-master   ClusterIP   10.97.53.188   <none>        9300/TCP   6m

NAME                                    DESIRED   CURRENT   AGE
statefulset.apps/custom-elasticsearch   2         2         6m
```

Check secrets created by KubeDB,

```console
$ kubectl get secret -n demo -l=kubedb.com/name=custom-elasticsearch
NAME                        TYPE      DATA      AGE
custom-elasticsearch-auth   Opaque    9         6m
custom-elasticsearch-cert   Opaque    4         6m
```

Once everything is created, Elasticsearch will go to `Running` state. Check that Elasticsearch is in running state.

```console
$ kubectl get es -n demo custom-elasticsearch
NAME                   VERSION    STATUS    AGE
custom-elasticsearch   6.2.4-v1   Running   7m
```

## Verify Configuration

Now, we will connect with the Elasticsearch cluster we have created. We will query for nodes settings to verify that the cluster is using the custom configuration we have provided.

At first, forward `9200` port of `custom-elasticsearch-0` pod. Run following command on a separate terminal,

```console
$ kubectl port-forward -n demo custom-elasticsearch-0 9200
Forwarding from 127.0.0.1:9200 -> 9200
Forwarding from [::1]:9200 -> 9200
```

Now, we can connect to the database at `localhost:9200`. Let's find out necessary connection information first.

**Connection information:**

- Address: `localhost:9200`
- Username: Run following command to get *username*

  ```console
  $ kubectl get secrets -n demo custom-elasticsearch-auth -o jsonpath='{.data.\ADMIN_USERNAME}' | base64 -d
    admin
  ```

- Password: Run following command to get *password*

  ```console
  $ kubectl get secrets -n demo custom-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
    fqc6rkha
  ```

Now, we will query for settings of all nodes in an Elasticsearch cluster,

```console
$ curl --user "admin:fqc6rkha" "localhost:9200/_nodes/_all/settings"
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
$ kubectl patch -n demo es/custom-elasticsearch -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo es/custom-elasticsearch

$ kubectl delete  -n demo configmap/es-custom-config

$ kubectl delete ns demo
```

To uninstall KubeDB follow this [guide](/docs/setup/uninstall.md).
