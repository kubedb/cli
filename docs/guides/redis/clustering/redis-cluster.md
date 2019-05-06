---
title: Redis Cluster Guide
menu:
  docs_0.12.0:
    identifier: rd-cluster
    name: Clustering Guide
    parent: rd-clustering-redis
    weight: 15
menu_name: docs_0.12.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# KubeDB - Redis Cluster

This tutorial will show you how to use KubeDB to provision a Redis cluster.

## Before You Begin

Before proceeding:

- Read [redis clustering concept](/docs/guides/redis/clustering/overview.md) to learn about Redis clustering.

- You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace/demo created
  ```

> Note: The yaml files used in this tutorial are stored in [docs/examples/redis](https://github.com/kubedb/cli/tree/master/docs/examples/redis) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Deploy Redis Cluster

To deploy a Redis Cluster, specify `spec.mode` and `spec.cluster` fields in `Redis` CRD.

The following is an example `Redis` object which creates a Redis cluster with three master nodes each of which has one replica node.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  name: redis-cluster
  namespace: demo
spec:
  version: 4.0-v2
  mode: Cluster
  cluster:
    master: 3
    replicas: 1
  storageType: Durable
  storage:
    resources:
      requests:
        storage: 1Gi
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.12.0/docs/examples/redis/clustering/demo-1.yaml
redis.kubedb.com/redis-cluster created
```

Here,

- `spec.mode` specifies the mode for Redis. Here we have used `Cluster` to tell the operator that we want to deploy Redis in cluster mode.
- `spec.cluster` represents the cluster configuration.
  - `master` denotes the number of master nodes.
  - `replicas` denotes the number of replica nodes per master.
- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. So, each members will have a pod of this storage configuration. You can specify any StorageClass available in your cluster with appropriate resource requests.

KubeDB operator watches for `Redis` objects using Kubernetes API. When a `Redis` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching Redis object name. KubeDB operator will also create a governing service for StatefulSets with the name `<redis-name>-gvr`. No Redis specific RBAC permission is required in [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb describe rd -n demo redis-cluster
Name:               redis-cluster
Namespace:          demo
CreationTimestamp:  Tue, 19 Feb 2019 19:28:59 +0600
Labels:             <none>
Annotations:        <none>
Replicas:           1  total
Status:             Running
  StorageType:      Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:
  Name:               redis-cluster-shard0
  CreationTimestamp:  Tue, 19 Feb 2019 19:28:59 +0600
  Labels:               kubedb.com/kind=Redis
                        kubedb.com/name=redis-cluster
  Annotations:        <none>
  Replicas:           824640878220 desired | 2 total
  Pods Status:        2 Running / 0 Waiting / 0 Succeeded / 0 Failed

StatefulSet:
  Name:               redis-cluster-shard1
  CreationTimestamp:  Tue, 19 Feb 2019 19:29:07 +0600
  Labels:               kubedb.com/kind=Redis
                        kubedb.com/name=redis-cluster
  Annotations:        <none>
  Replicas:           824640879052 desired | 2 total
  Pods Status:        2 Running / 0 Waiting / 0 Succeeded / 0 Failed

StatefulSet:
  Name:               redis-cluster-shard2
  CreationTimestamp:  Tue, 19 Feb 2019 19:29:13 +0600
  Labels:               kubedb.com/kind=Redis
                        kubedb.com/name=redis-cluster
  Annotations:        <none>
  Replicas:           824640879900 desired | 2 total
  Pods Status:        2 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         redis-cluster
  Labels:         kubedb.com/kind=Redis
                  kubedb.com/name=redis-cluster
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.100.246.86
  Port:         db  6379/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.10:6379,172.17.0.11:6379,172.17.0.12:6379 + 3 more...

No Snapshots.

Events:
  Type    Reason      Age   From            Message
  ----    ------      ----  ----            -------
  Normal  Successful  1m    Redis operator  Successfully created ConfigMap
  Normal  Successful  1m    Redis operator  Successfully created Service
  Normal  Successful  1m    Redis operator  Successfully created StatefulSet
  Normal  Successful  1m    Redis operator  Successfully created StatefulSet
  Normal  Successful  1m    Redis operator  Successfully created StatefulSet
  Normal  Successful  24s   Redis operator  Successfully created Redis
  Normal  Successful  24s   Redis operator  Successfully created appbinding
  Normal  Successful  24s   Redis operator  Successfully patched StatefulSet
  Normal  Successful  24s   Redis operator  Successfully patched StatefulSet
  Normal  Successful  24s   Redis operator  Successfully patched StatefulSet
  Normal  Successful  20s   Redis operator  Successfully patched Redis

$ kubectl get statefulset -n demo
NAME                   READY   AGE
redis-cluster-shard0   2/2     107s
redis-cluster-shard1   2/2     99s
redis-cluster-shard2   2/2     93s

$ kubectl get pvc -n demo
NAME                          STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-redis-cluster-shard0-0   Bound    pvc-50b82133-344a-11e9-b1be-0800275426d2   1Gi        RWO            standard       2m2s
data-redis-cluster-shard0-1   Bound    pvc-51ee270b-344a-11e9-b1be-0800275426d2   1Gi        RWO            standard       2m
data-redis-cluster-shard1-0   Bound    pvc-550d1008-344a-11e9-b1be-0800275426d2   1Gi        RWO            standard       114s
data-redis-cluster-shard1-1   Bound    pvc-564b493a-344a-11e9-b1be-0800275426d2   1Gi        RWO            standard       112s
data-redis-cluster-shard2-0   Bound    pvc-58c40c52-344a-11e9-b1be-0800275426d2   1Gi        RWO            standard       108s
data-redis-cluster-shard2-1   Bound    pvc-5c761601-344a-11e9-b1be-0800275426d2   1Gi        RWO            standard       102s

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                              STORAGECLASS   REASON   AGE
pvc-50b82133-344a-11e9-b1be-0800275426d2   1Gi        RWO            Delete           Bound    demo/data-redis-cluster-shard0-0   standard                2m21s
pvc-51ee270b-344a-11e9-b1be-0800275426d2   1Gi        RWO            Delete           Bound    demo/data-redis-cluster-shard0-1   standard                2m15s
pvc-550d1008-344a-11e9-b1be-0800275426d2   1Gi        RWO            Delete           Bound    demo/data-redis-cluster-shard1-0   standard                2m13s
pvc-564b493a-344a-11e9-b1be-0800275426d2   1Gi        RWO            Delete           Bound    demo/data-redis-cluster-shard1-1   standard                2m8s
pvc-58c40c52-344a-11e9-b1be-0800275426d2   1Gi        RWO            Delete           Bound    demo/data-redis-cluster-shard2-0   standard                2m3s
pvc-5c761601-344a-11e9-b1be-0800275426d2   1Gi        RWO            Delete           Bound    demo/data-redis-cluster-shard2-1   standard                2m

$ kubectl get service -n demo
NAME            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
kubedb          ClusterIP   None            <none>        <none>     2m39s
redis-cluster   ClusterIP   10.100.246.86   <none>        6379/TCP   2m39s
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified `Redis` object:

```yaml
$ kubedb get rd -n demo redis-cluster -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  creationTimestamp: "2019-02-19T13:28:59Z"
  finalizers:
  - kubedb.com
  generation: 2
  name: redis-cluster
  namespace: demo
  resourceVersion: "569405"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/redises/redis-cluster
  uid: 509b50ba-344a-11e9-b1be-0800275426d2
spec:
  cluster:
    master: 3
    replicas: 1
  configSource:
    configMap:
      defaultMode: 511
      name: redis-cluster
  mode: Cluster
  podTemplate:
    controller: {}
    metadata: {}
    spec:
      resources: {}
  replicas: 1
  serviceTemplate:
    metadata: {}
    spec: {}
  storage:
    accessModes:
    - ReadWriteOnce
    dataSource: null
    resources:
      requests:
        storage: 1Gi
    storageClassName: standard
  storageType: Durable
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
  version: 4.0-v2
status:
  observedGeneration: 2$4213139756412538772
  phase: Running
```

## Check Cluster Scenario

The operator creates a cluster according to the newly created `Redis` object. This cluster has 3 masters and one replica per master. And every node in the cluster is responsible for a subset of the total **16384** hash slots.

```console
# first list the redis pods list
$ kubectl get pods --all-namespaces -o jsonpath='{range.items[*]}{.metadata.name} ---------- {.status.podIP}:6379{"\\n"}{end}' | grep redis
redis-cluster-shard0-0 ---------- 172.17.0.4:6379
redis-cluster-shard0-1 ---------- 172.17.0.8:6379
redis-cluster-shard1-0 ---------- 172.17.0.10:6379
redis-cluster-shard1-1 ---------- 172.17.0.11:6379
redis-cluster-shard2-0 ---------- 172.17.0.12:6379
redis-cluster-shard2-1 ---------- 172.17.0.13:6379

# enter into any pod's container named redis
$ kubectl exec -it redis-cluster-shard0-0 -n demo -c redis -- sh
/data #

# now inside this container, see which ones are the masters
# which ones are the replicas
/data # redis-cli -c cluster nodes
e49d748b815db355f29e670ba62b669627384273 172.17.0.10:6379@16379 master - 0 1550585018746 2 connected 5461-10922
0cfa026b921bd07c36a95734edf5ccd73cda5d50 172.17.0.12:6379@16379 master - 0 1550585019750 3 connected 10923-16383
e425958b07698cc69e62e0f2c94f5951155cbe71 172.17.0.11:6379@16379 replica e49d748b815db355f29e670ba62b669627384273 0 1550585018540 2 connected
37ae24f6f2442cefc34fc5d3c678b2aff8f13d26 172.17.0.4:6379@16379 myself,master - 0 1550585019000 1 connected 0-5460
4136ea1a767fd26d76ad7f8066d7eab994850048 172.17.0.8:6379@16379 replica 37ae24f6f2442cefc34fc5d3c678b2aff8f13d26 0 1550585018540 1 connected
981f26c1e2d16f56109ca74ee79aaa5cd5e62a79 172.17.0.13:6379@16379 replica 0cfa026b921bd07c36a95734edf5ccd73cda5d50 0 1550585019000 3 connected
```

- redis-cluster-shard0-0
  - `ip` 172.17.0.4:6379
  - `role` master
  - `node-id` 37ae24f6f2442cefc34fc5d3c678b2aff8f13d26
  - `slot` 0-5460
- redis-cluster-shard0-1
  - `ip` 172.17.0.8:6379
  - `role` replica
  - `node-id` 4136ea1a767fd26d76ad7f8066d7eab994850048
- redis-cluster-shard1-0
  - `ip` 172.17.0.10:6379
  - `role` master
  - `node-id` e49d748b815db355f29e670ba62b669627384273
  - `slot` 5461-10922
- redis-cluster-shard1-1
  - `ip` 172.17.0.11:6379
  - `role` replica
  - `node-id` e425958b07698cc69e62e0f2c94f5951155cbe71
- redis-cluster-shard2-0
  - `ip` 172.17.0.12:6379
  - `role` master
  - `node-id` 0cfa026b921bd07c36a95734edf5ccd73cda5d50
  - `slot` 5461-10922
- redis-cluster-shard2-1
  - `ip` 172.17.0.13:6379
  - `role` replica
  - `node-id` 981f26c1e2d16f56109ca74ee79aaa5cd5e62a79

Every replica node will serve for the same hash slot as its master.

## Data Availability

Now, you can connect to this database through [redis-cli](https://redis.io/topics/rediscli). In this tutorial, we will insert data, and we will see whether we can get the data from any other node (any master or replica) or not.

> Read the comment written for the following commands. They contain the instructions and explanations of the commands.

```console
# here the hash slot for key 'hello' is 866 which is in 1st node
# named 'redis-cluster-shard0-0' (0-5460)
$ kubectl exec -it redis-cluster-shard0-0 -n demo -c redis -- redis-cli -c cluster keyslot hello
(integer) 866

# connect to any node
kubectl exec -it redis-cluster-shard0-0 -n demo -c redis -- sh
/data #

# now ensure that you are connected to the 1st pod
/data # redis-cli -c -h 172.17.0.4
172.17.0.4:6379>

# set 'world' as value for the key 'hello'
172.17.0.4:6379> set hello world
OK
172.17.0.4:6379> exit

# switch the connection to the replica of the current master and get the data
/data # redis-cli -c -h 172.17.0.8
172.17.0.8:6379> get hello
-> Redirected to slot [866] located at 172.17.0.4:6379
"world"
172.17.0.4:6379> exit

# switch the connection to any other node
# get the data
/data # redis-cli -c -h 172.17.0.11
172.17.0.11:6379> get hello
-> Redirected to slot [866] located at 172.17.0.4:6379
"world"
172.17.0.4:6379> exit
```

## Automatic Failover

To test automatic failover, we will force a master node to restart. Since the master node (`pod`) becomes unavailable, the rest of the members will elect a replica (one of its replica in case of more than one replica under this master) of this master node as the new master. When the old master comes back, it will join the cluster as the new replica of the new master.

> Read the comment written for the following commands. They contain the instructions and explanations of the commands.

```console
# connect to any node and get the master nodes info
$ kubectl exec -it redis-cluster-shard0-0 -n demo -c redis -- sh
/data # redis-cli -c cluster nodes | grep master
e49d748b815db355f29e670ba62b669627384273 172.17.0.10:6379@16379 master - 0 1550589457000 2 connected 5461-10922
0cfa026b921bd07c36a95734edf5ccd73cda5d50 172.17.0.12:6379@16379 master - 0 1550589457739 3 connected 10923-16383
37ae24f6f2442cefc34fc5d3c678b2aff8f13d26 172.17.0.4:6379@16379 myself,master - 0 1550589456000 1 connected 0-5460

# let's crash node 172.17.0.4 with the `DEBUG SEGFAULT` command
/data # redis-cli -h 172.17.0.4 debug segfault
Error: Server closed the connection

# now again connect to a node and get the master nodes info
$ kubectl exec -it redis-cluster-shard0-0 -n demo -c redis -- sh
/data # redis-cli -c cluster nodes | grep master
e49d748b815db355f29e670ba62b669627384273 172.17.0.10:6379@16379 master - 0 1550589881100 2 connected 5461-10922
4136ea1a767fd26d76ad7f8066d7eab994850048 172.17.0.8:6379@16379 master - 0 1550589880000 4 connected 0-5460
0cfa026b921bd07c36a95734edf5ccd73cda5d50 172.17.0.12:6379@16379 master - 0 1550589881000 3 connected 10923-16383

/data # redis-cli -c cluster nodes
e425958b07698cc69e62e0f2c94f5951155cbe71 172.17.0.11:6379@16379 replica e49d748b815db355f29e670ba62b669627384273 0 1550590186590 2 connected
e49d748b815db355f29e670ba62b669627384273 172.17.0.10:6379@16379 master - 0 1550590186990 2 connected 5461-10922
981f26c1e2d16f56109ca74ee79aaa5cd5e62a79 172.17.0.13:6379@16379 replica 0cfa026b921bd07c36a95734edf5ccd73cda5d50 0 1550590186000 3 connected
37ae24f6f2442cefc34fc5d3c678b2aff8f13d26 172.17.0.4:6379@16379 myself,replica 4136ea1a767fd26d76ad7f8066d7eab994850048 0 1550590185000 1 connected
4136ea1a767fd26d76ad7f8066d7eab994850048 172.17.0.8:6379@16379 master - 0 1550590184585 4 connected 0-5460
0cfa026b921bd07c36a95734edf5ccd73cda5d50 172.17.0.12:6379@16379 master - 0 1550590186000 3 connected 10923-16383

/data # exit
```

Notice that 172.17.0.8 is the new master and  172.17.0.4 is the replica of  172.17.0.8.

## Cleaning up

Clean what you created in this tutorial.

```yaml
$ kubedb edit rd -n demo redis-cluster -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  name: redis-cluster
  namespace: demo
spec:
  ...
  terminationPolicy: WipeOut
  ...
status:
  ...
  phase: Running

$ kubedb delete rd redis-cluster -n demo
redis.kubedb.com "redis-cluster" deleted
```

## Next Steps

- Monitor your Redis database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/redis/monitoring/using-coreos-prometheus-operator.md).
- Monitor your Redis database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/redis/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/redis/private-registry/using-private-registry.md) to deploy Redis with KubeDB.
- Detail concepts of [Redis object](/docs/concepts/databases/redis.md).
- Detail concepts of [RedisVersion object](/docs/concepts/catalog/redis.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
