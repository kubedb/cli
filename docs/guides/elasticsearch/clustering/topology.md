---
title: Elasticsearch Cluster Topology
menu:
  docs_0.8.0-beta.2:
    identifier: es-topology-clustering
    name: Topology
    parent: es-clustering-elasticsearch
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Elasticsearch Topology

KubeDB Elasticsearch supports multi-node database cluster.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

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

## Create multi-node Elasticsearch

Elasticsearch can be created with multiple nodes. If you want to create Elasticsearch cluster with three nodes, you need to set `spec.replicas` to `3`.
In this case, all of these three nodes will act as *master*, *data* and *client*.

Check following Elasticsearch object

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: multi-node-es
  namespace: demo
spec:
  version: 5.6
  replicas: 3
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
```

Here,

- `spec.replicas` is the number of nodes in the Elasticsearch cluster. Here, we are creating a three node Elasticsearch cluster.

> Note: If `spec.topology` is set, `spec.replicas` has no effect.

Create example above with following command

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/elasticsearch/clustering/multi-node-es.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/elasticsearch/clustering/multi-node-es.yaml"
elasticsearch "multi-node-es" created
```

Lets describe Elasticsearch object `multi-node-es` while Running

```console
$ kubedb describe es -S=false -W=false -n demo multi-node-es
Name:			multi-node-es
Namespace:		demo
CreationTimestamp:	Tue, 20 Feb 2018 14:36:03 +0600
Status:			Running
Replicas:		3  total
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO
StatefulSet:	multi-node-es
Service:	multi-node-es, multi-node-es-master
Secrets:	multi-node-es-auth, multi-node-es-cert

Topology:
  Type                 Pod               StartTime                       Phase
  ----                 ---               ---------                       -----
  master|client|data   multi-node-es-0   2018-02-20 14:36:13 +0600 +06   Running
  master|client|data   multi-node-es-1   2018-02-20 14:36:24 +0600 +06   Running
  master|client|data   multi-node-es-2   2018-02-20 14:36:45 +0600 +06   Running

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason       Message
  ---------   --------   -----     ----                     --------   ------       -------
  12m         12m        1         Elasticsearch operator   Normal     Successful   Successfully patched Elasticsearch
  13m         13m        1         Elasticsearch operator   Normal     Successful   Successfully patched StatefulSet
  13m         13m        1         Elasticsearch operator   Normal     Successful   Successfully created Elasticsearch
  13m         13m        1         Elasticsearch operator   Normal     Successful   Successfully created StatefulSet
  14m         14m        1         Elasticsearch operator   Normal     Successful   Successfully created Service
  14m         14m        1         Elasticsearch operator   Normal     Successful   Successfully created Service
```

Here, we can see in Topology section that all three Pods are acting as *master*, *data* and *client*.

## Create Elasticsearch with dedicated node

If you want to use separate node for *master*, *data* and *client* role, you need to configure `spec.topology`

In this tutorial, we will create following Elasticsearch with topology

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: topology-es
  namespace: demo
spec:
  version: 5.6
  topology:
    master:
      prefix: master
      replicas: 1
    data:
      prefix: data
      replicas: 2
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

Here,

- `spec.topology` point to the number of pods we want as dedicated `master`, `client` and `data` nodes and also specify prefix for their StatefulSet name

Lets create this Elasticsearch object

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/elasticsearch/clustering/topology-es.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/elasticsearch/clustering/topology-es.yaml"
elasticsearch "topology-es" created
```

When this object is created, Elasticsearch database has started with 5 pods under 3 different StatefulSets.

```console
$ kubectl get statefulset -n demo --show-labels --selector="kubedb.com/kind=Elasticsearch,kubedb.com/name=topology-es"
NAME                 DESIRED   CURRENT   AGE       LABELS
client-topology-es   2         2         6m        kubedb.com/kind=Elasticsearch,kubedb.com/name=topology-es,node.role.client=set
data-topology-es     2         2         2m        kubedb.com/kind=Elasticsearch,kubedb.com/name=topology-es,node.role.data=set
master-topology-es   1         1         2m        kubedb.com/kind=Elasticsearch,kubedb.com/name=topology-es,node.role.master=set
```

Three StatefulSets are created

- client-topology-es

    ```yaml
    spec:
      topology:
        client:
          prefix: client
          replicas: 2
    ```

    This configuration creates a StatefulSet named `client-topology-es` for client node

    - `spec.replicas` is set to `2`. Two dedicated nodes is created as client.
    - Label `node.role.client: set` is added in Pods

- data-topology-es

    ```yaml
    spec:
      topology:
        data:
          prefix: data
          replicas: 2
    ```

    This configuration creates a StatefulSet named `data-topology-es` for data node

    - `spec.replicas` is set to `2`. Two dedicated nodes is created for data.

- master-topology-es

    ```yaml
    spec:
      topology:
        master:
          prefix: master
          replicas: 1
    ```

    This configuration creates a StatefulSet named `data-topology-es` for master node

    - `spec.replicas` is set to `1`. One dedicated node is created as master.
    - Label `node.role.master: set` is added in Pods


> Note: StatefulSet name format: `{topology-prefix}-{elasticsearch-name}`

Lets describe this Elasticsearch

```console
$ kubedb describe es -S=false -W=false -n demo topology-es
Name:			topology-es
Namespace:		demo
CreationTimestamp:	Tue, 20 Feb 2018 16:34:43 +0600
Status:			Running
Replicas:		0  total
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO
StatefulSet:	client-topology-es, data-topology-es, master-topology-es
Service:	topology-es, topology-es-master
Secrets:	topology-es-auth, topology-es-cert

Topology:
  Type     Pod                    StartTime                       Phase
  ----     ---                    ---------                       -----
  client   client-topology-es-0   2018-02-20 16:34:50 +0600 +06   Running
  client   client-topology-es-1   2018-02-20 16:38:23 +0600 +06   Running
  data     data-topology-es-0     2018-02-20 16:39:12 +0600 +06   Running
  data     data-topology-es-1     2018-02-20 16:39:40 +0600 +06   Running
  master   master-topology-es-0   2018-02-20 16:38:44 +0600 +06   Running

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason       Message
  ---------   --------   -----     ----                     --------   ------       -------
  23m         23m        1         Elasticsearch operator   Normal     Successful   Successfully patched Elasticsearch
  24m         24m        1         Elasticsearch operator   Normal     Successful   Successfully patched StatefulSet
  24m         24m        1         Elasticsearch operator   Normal     Successful   Successfully patched StatefulSet
  24m         24m        1         Elasticsearch operator   Normal     Successful   Successfully patched StatefulSet
  24m         24m        1         Elasticsearch operator   Normal     Successful   Successfully created Elasticsearch
  25m         25m        1         Elasticsearch operator   Normal     Successful   Successfully created StatefulSet
  26m         26m        1         Elasticsearch operator   Normal     Successful   Successfully created StatefulSet
  26m         26m        1         Elasticsearch operator   Normal     Successful   Successfully created StatefulSet
  30m         30m        1         Elasticsearch operator   Normal     Successful   Successfully created Service
  30m         30m        1         Elasticsearch operator   Normal     Successful   Successfully created Service
```

We can see in Topology section that

Two Pods are dedicated as client

```
Topology:
  Type     Pod                    StartTime                       Phase
  ----     ---                    ---------                       -----
  client   client-topology-es-0   2018-02-20 16:34:50 +0600 +06   Running
  client   client-topology-es-1   2018-02-20 16:38:23 +0600 +06   Running
```

Two Pods for data node

```
Topology:
  Type     Pod                    StartTime                       Phase
  ----     ---                    ---------                       -----
  data     data-topology-es-0     2018-02-20 16:39:12 +0600 +06   Running
  data     data-topology-es-1     2018-02-20 16:39:40 +0600 +06   Running
```

And one Pod as master node

```
Topology:
  Type     Pod                    StartTime                       Phase
  ----     ---                    ---------                       -----
  master   master-topology-es-0   2018-02-20 16:38:44 +0600 +06   Running
```


Two services are also created for this Elasticsearch object.

 - Service *`quick-elasticsearch`* targets all Pods which are acting as *client* node
 - Service *`quick-elasticsearch-master`* targets all Pods which are acting as *master* node


## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete es,drmn,snap -n demo --all --force

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Learn about [taking instant backup](/docs/guides/elasticsearch/snapshot/instant_backup.md) of Elasticsearch database using KubeDB.
- Learn how to [schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Learn about initializing [Elasticsearch with Snapshot](/docs/guides/elasticsearch/initialization/snapshot_source.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Detail concepts of [Elasticsearch object](/docs/concepts/databases/elasticsearch.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
