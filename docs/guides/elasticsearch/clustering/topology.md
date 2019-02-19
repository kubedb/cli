---
title: Elasticsearch Cluster Topology
menu:
  docs_0.9.0:
    identifier: es-topology-clustering
    name: Topology
    parent: es-clustering-elasticsearch
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Elasticsearch Topology

KubeDB Elasticsearch supports multi-node database cluster.

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

## Create multi-node Elasticsearch

Elasticsearch can be created with multiple nodes. If you want to create an Elasticsearch cluster with three nodes, you need to set `spec.replicas` to `3`. In this case, all of these three nodes will act as *master*, *data* and *client*.

Check following Elasticsearch object

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: multi-node-es
  namespace: demo
spec:
  version: "6.3-v1"
  replicas: 3
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

Here,

- `spec.replicas` is the number of nodes in the Elasticsearch cluster. Here, we are creating a three node Elasticsearch cluster.

> Note: If `spec.topology` is set, you won't able to `spec.replicas`. KubeDB will reject the create request for Elasticsearch crd from validating webhook.

Create example above with following command

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/clustering/multi-node-es.yaml
elasticsearch.kubedb.com/multi-node-es created
```

Let's describe Elasticsearch object `multi-node-es` while Running

```console
$ kubedb describe es -n demo multi-node-es
Name:               multi-node-es
Namespace:          demo
CreationTimestamp:  Fri, 05 Oct 2018 12:51:43 +0600
Labels:             <none>
Annotations:        kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"kubedb.com/v1alpha1","kind":"Elasticsearch","metadata":{"annotations":{},"name":"multi-node-es","namespace":"demo"},"spec":{"replicas":3...
Status:             Running
Replicas:           3  total
  StorageType:      Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:          
  Name:               multi-node-es
  CreationTimestamp:  Fri, 05 Oct 2018 12:51:45 +0600
  Labels:               kubedb.com/kind=Elasticsearch
                        kubedb.com/name=multi-node-es
                        node.role.client=set
                        node.role.data=set
                        node.role.master=set
  Annotations:        <none>
  Replicas:           824639906344 desired | 3 total
  Pods Status:        3 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         multi-node-es
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=multi-node-es
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.101.96.68
  Port:         http  9200/TCP
  TargetPort:   http/TCP
  Endpoints:    192.168.1.18:9200,192.168.1.19:9200,192.168.1.20:9200

Service:        
  Name:         multi-node-es-master
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=multi-node-es
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.103.195.123
  Port:         transport  9300/TCP
  TargetPort:   transport/TCP
  Endpoints:    192.168.1.18:9300,192.168.1.19:9300,192.168.1.20:9300

Database Secret:
  Name:         multi-node-es-auth
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=multi-node-es
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  ADMIN_USERNAME:         5 bytes
  READALL_USERNAME:       7 bytes
  sg_action_groups.yml:   430 bytes
  sg_roles.yml:           312 bytes
  ADMIN_PASSWORD:         8 bytes
  READALL_PASSWORD:       8 bytes
  sg_config.yml:          242 bytes
  sg_internal_users.yml:  156 bytes
  sg_roles_mapping.yml:   73 bytes

Certificate Secret:
  Name:         multi-node-es-cert
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=multi-node-es
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  key_pass:     6 bytes
  node.jks:     3009 bytes
  root.jks:     864 bytes
  sgadmin.jks:  3010 bytes

Topology:
  Type                Pod              StartTime                      Phase
  ----                ---              ---------                      -----
  master|client|data  multi-node-es-0  2018-10-05 12:51:46 +0600 +06  Running
  client|data|master  multi-node-es-1  2018-10-05 12:52:03 +0600 +06  Running
  master|client|data  multi-node-es-2  2018-10-05 12:52:33 +0600 +06  Running

No Snapshots.

Events:
  Type    Reason      Age   From                    Message
  ----    ------      ----  ----                    -------
  Normal  Successful  4m    Elasticsearch operator  Successfully created Service
  Normal  Successful  4m    Elasticsearch operator  Successfully created Service
  Normal  Successful  2m    Elasticsearch operator  Successfully created StatefulSet
  Normal  Successful  2m    Elasticsearch operator  Successfully created Elasticsearch
  Normal  Successful  2m    Elasticsearch operator  Successfully patched StatefulSet
  Normal  Successful  1m    Elasticsearch operator  Successfully patched Elasticsearch
  Normal  Successful  1m    Elasticsearch operator  Successfully patched StatefulSet
  Normal  Successful  1m    Elasticsearch operator  Successfully patched Elasticsearch
```

Here, we can see in `Topology` section that all three Pods are acting as *master*, *data* and *client*.

## Create Elasticsearch with dedicated node

If you want to use separate node for *master*, *data* and *client* role, you need to configure `spec.topology` field of Elasticsearch crd.

In this tutorial, we will create following Elasticsearch with topology

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: topology-es
  namespace: demo
spec:
  version: "6.3-v1"
  storageType: Durable
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
      replicas: 2
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
            storage: 1Gi
```

Here,

- `spec.topology` point to the number of pods we want as dedicated `master`, `client` and `data` nodes and also specify prefix, storage, resources for the pods.

Let's create this Elasticsearch object

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/clustering/topology-es.yaml
elasticsearch.kubedb.com/topology-es created
```

When this object is created, Elasticsearch database has started with 5 pods under 3 different StatefulSets.

```console
$ kubectl get statefulset -n demo --show-labels --selector="kubedb.com/kind=Elasticsearch,kubedb.com/name=topology-es"
NAME                 DESIRED   CURRENT   AGE       LABELS
client-topology-es   2         2         1m        kubedb.com/kind=Elasticsearch,kubedb.com/name=topology-es,node.role.client=set
data-topology-es     2         2         48s       kubedb.com/kind=Elasticsearch,kubedb.com/name=topology-es,node.role.data=set
master-topology-es   1         1         1m        kubedb.com/kind=Elasticsearch,kubedb.com/name=topology-es,node.role.master=set
```

Three StatefulSets are created for *client*, *data* and *master* node respectively.

- client-topology-es

    ```yaml
    spec:
      topology:
        client:
          prefix: client
          replicas: 2
          storage:
            storageClassName: "standard"
            accessModes:
            - ReadWriteOnce
            resources:
              requests:
                storage: 1Gi
    ```

    This configuration creates a StatefulSet named `client-topology-es` for client node

  - `spec.replicas` is set to `2`. Two dedicated nodes is created as client.
  - Label `node.role.client: set` is added in Pods
  - Each Pod will receive a single PersistentVolume with a StorageClass of **standard** and **1Gi** of provisioned storage.

- data-topology-es

    ```yaml
    spec:
      topology:
        data:
          prefix: data
          replicas: 2
          storage:
            storageClassName: "standard"
            accessModes:
            - ReadWriteOnce
            resources:
              requests:
                storage: 1Gi
    ```

  This configuration creates a StatefulSet named `data-topology-es` for data node.

  - `spec.replicas` is set to `2`. Two dedicated nodes is created for data.
  - Label `node.role.data: set` is added in Pods
  - Each Pod will receive a single PersistentVolume with a StorageClass of **standard** and **1 Gib** of provisioned storage. 

- master-topology-es

    ```yaml
    spec:
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
    ```

    This configuration creates a StatefulSet named `data-topology-es` for master node

  - `spec.replicas` is set to `1`. One dedicated node is created as master.
  - Label `node.role.master: set` is added in Pods
  - Each Pod will receive a single PersistentVolume with a StorageClass of **standard** and **1Gi** of provisioned storage.

> Note: StatefulSet name format: `{topology-prefix}-{elasticsearch-name}`

Let's describe this Elasticsearch

```console
$ kubedb describe es -n demo topology-es
Name:               topology-es
Namespace:          demo
CreationTimestamp:  Fri, 05 Oct 2018 14:40:24 +0600
Labels:             <none>
Annotations:        kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"kubedb.com/v1alpha1","kind":"Elasticsearch","metadata":{"annotations":{},"name":"topology-es","namespace":"demo"},"spec":{"storageType":...
Status:             Creating
  StorageType:      Durable
No volumes.

StatefulSet:          
  Name:               client-topology-es
  CreationTimestamp:  Fri, 05 Oct 2018 14:40:26 +0600
  Labels:               kubedb.com/kind=Elasticsearch
                        kubedb.com/name=topology-es
                        node.role.client=set
  Annotations:        <none>
  Replicas:           824640231228 desired | 2 total
  Pods Status:        2 Running / 0 Waiting / 0 Succeeded / 0 Failed

StatefulSet:          
  Name:               data-topology-es
  CreationTimestamp:  Fri, 05 Oct 2018 14:41:21 +0600
  Labels:               kubedb.com/kind=Elasticsearch
                        kubedb.com/name=topology-es
                        node.role.data=set
  Annotations:        <none>
  Replicas:           824640232652 desired | 2 total
  Pods Status:        2 Running / 0 Waiting / 0 Succeeded / 0 Failed

StatefulSet:          
  Name:               master-topology-es
  CreationTimestamp:  Fri, 05 Oct 2018 14:41:00 +0600
  Labels:               kubedb.com/kind=Elasticsearch
                        kubedb.com/name=topology-es
                        node.role.master=set
  Annotations:        <none>
  Replicas:           824641372604 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         topology-es
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=topology-es
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.98.164.122
  Port:         http  9200/TCP
  TargetPort:   http/TCP
  Endpoints:    192.168.1.26:9200,192.168.1.27:9200

Service:        
  Name:         topology-es-master
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=topology-es
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.111.168.40
  Port:         transport  9300/TCP
  TargetPort:   transport/TCP
  Endpoints:    192.168.1.28:9300

Database Secret:
  Name:         topology-es-auth
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=topology-es
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  ADMIN_USERNAME:         5 bytes
  READALL_PASSWORD:       8 bytes
  sg_internal_users.yml:  156 bytes
  ADMIN_PASSWORD:         8 bytes
  READALL_USERNAME:       7 bytes
  sg_action_groups.yml:   430 bytes
  sg_config.yml:          242 bytes
  sg_roles.yml:           312 bytes
  sg_roles_mapping.yml:   73 bytes

Certificate Secret:
  Name:         topology-es-cert
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=topology-es
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  key_pass:     6 bytes
  node.jks:     3008 bytes
  root.jks:     864 bytes
  sgadmin.jks:  3010 bytes

Topology:
  Type    Pod                   StartTime                      Phase
  ----    ---                   ---------                      -----
  client  client-topology-es-0  2018-10-05 14:40:27 +0600 +06  Running
  client  client-topology-es-1  2018-10-05 14:40:44 +0600 +06  Running
  data    data-topology-es-0    2018-10-05 14:41:22 +0600 +06  Running
  data    data-topology-es-1    2018-10-05 14:41:49 +0600 +06  Running
  master  master-topology-es-0  2018-10-05 14:41:01 +0600 +06  Running

No Snapshots.

Events:
  Type    Reason      Age   From                    Message
  ----    ------      ----  ----                    -------
  Normal  Successful  2m    Elasticsearch operator  Successfully created Service
  Normal  Successful  2m    Elasticsearch operator  Successfully created Service
  Normal  Successful  1m    Elasticsearch operator  Successfully created StatefulSet
  Normal  Successful  1m    Elasticsearch operator  Successfully created StatefulSet
  Normal  Successful  18s   Elasticsearch operator  Successfully created StatefulSet
```

Here, we can see from `Topology` section that 2 pods working as *client*, 2 pods working as *data* and 1 pod working as *master*.

Two services are also created for this Elasticsearch object.

- Service *`quick-elasticsearch`* targets all Pods which are acting as *client* node
- Service *`quick-elasticsearch-master`* targets all Pods which are acting as *master* node

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/multi-node-es es/topology-es -p '{"spec":{"terminationPolicy": "WipeOut"}}' --type="merge"
$ kubectl delete -n demo es/multi-node-es es/topology-es

$ kubectl delete ns demo
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
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
