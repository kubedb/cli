---
title: Redis Quickstart
menu:
  docs_0.9.0:
    identifier: rd-quickstart-quickstart
    name: Overview
    parent: rd-quickstart-redis
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Redis QuickStart

This tutorial will show you how to use KubeDB to run a Redis server.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/redis/redis-lifecycle.png">
</p>

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/) is required to run KubeDB. Check the available StorageClass in cluster.

  ```console
  $ kubectl get storageclasses
  NAME                 PROVISIONER                AGE
  standard (default)   k8s.io/minikube-hostpath   4h
  ```

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace/demo created

  $ kubectl get ns
  NAME          STATUS    AGE
  demo          Active    10s
  ```

> Note: The yaml files used in this tutorial are stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Find Available RedisVersion

When you have installed KubeDB, it has created `RedisVersion` crd for all supported Redis versions. Check 0

```console
$ kubectl get redisversions
NAME       VERSION   DB_IMAGE                DEPRECATED   AGE
4          4         kubedb/redis:4          true         30m
4-v1       4         kubedb/redis:4-v1                    30m
4.0        4.0       kubedb/redis:4.0        true         30m
4.0-v1     4.0       kubedb/redis:4.0-v1                  30m
4.0.6      4.0.6     kubedb/redis:4.0.6-v1   true         30m
4.0.6-v1   4.0.6     kubedb/redis:4.0.6-v1                30m
```

## Create a Redis server

KubeDB implements a `Redis` CRD to define the specification of a Redis server. Below is the `Redis` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  name: redis-quickstart
  namespace: demo
spec:
  version: "4.0-v1"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/redis/quickstart/demo-1.yaml
redis.kubedb.com/redis-quickstart created
```

Here,

- `spec.version` is name of the RedisVersion crd where the docker images are specified. In this tutorial, a Redis 4.0-v1 database is created.
- `spec.storageType` specifies the type of storage that will be used for Redis server. It can be `Durable` or `Ephemeral`. Default value of this field is `Durable`. If `Ephemeral` is used then KubeDB will create Redis server using `EmptyDir` volume. In this case, you don't have to specify `spec.storage` field. This is useful for testing purpose.
- `spec.storage` specifies PVC spec that will be dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests.
- `spec.terminationPolicy` gives flexibility whether to `nullify`(reject) the delete operation of `Redis` crd or which resources KubeDB should keep or delete when you delete `Redis` crd. If admission webhook is enabled, It prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`. Learn details of all `TerminationPolicy` [here](/docs/concepts/databases/redis.md#specterminationpolicy)

> Note: spec.storage section is used to create PVC for database pod. It will create PVC with storage size specified instorage.resources.requests field. Don't specify limits here. PVC does not get resized automatically.

KubeDB operator watches for `Redis` objects using Kubernetes api. When a `Redis` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching Redis object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. No Redis specific RBAC permission is required in [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb get rd -n demo
NAME               VERSION   STATUS    AGE
redis-quickstart   4.0-v1    Running   1m

$ kubedb describe rd -n demo redis-quickstart
Name:               redis-quickstart
Namespace:          demo
CreationTimestamp:  Mon, 01 Oct 2018 12:01:23 +0600
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
  Name:               redis-quickstart
  CreationTimestamp:  Mon, 01 Oct 2018 12:01:25 +0600
  Labels:               kubedb.com/kind=Redis
                        kubedb.com/name=redis-quickstart
  Annotations:        <none>
  Replicas:           824641951004 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         redis-quickstart
  Labels:         kubedb.com/kind=Redis
                  kubedb.com/name=redis-quickstart
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.108.149.205
  Port:         db  6379/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.4:6379

No Snapshots.

Events:
  Type    Reason      Age   From            Message
  ----    ------      ----  ----            -------
  Normal  Successful  1m    Redis operator  Successfully created Service
  Normal  Successful  53s   Redis operator  Successfully created StatefulSet
  Normal  Successful  53s   Redis operator  Successfully created Redis
  Normal  Successful  52s   Redis operator  Successfully patched StatefulSet
  Normal  Successful  52s   Redis operator  Successfully patched Redis

$ kubectl get statefulset -n demo
NAME               DESIRED   CURRENT   AGE
redis-quickstart   1         1         1m

$ kubectl get pvc -n demo
NAME                      STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-redis-quickstart-0   Bound     pvc-6e457226-c53f-11e8-9ba7-0800274bef12   1Gi        RWO            standard       2m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                          STORAGECLASS   REASON    AGE
pvc-6e457226-c53f-11e8-9ba7-0800274bef12   1Gi        RWO            Delete           Bound     demo/data-redis-quickstart-0   standard                 2m

$ kubectl get service -n demo
NAME               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
kubedb             ClusterIP   None             <none>        <none>     2m
redis-quickstart   ClusterIP   10.108.149.205   <none>        6379/TCP   2m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified Redis object:

```yaml
$ kubedb get rd -n demo redis-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  creationTimestamp: 2018-10-01T06:01:23Z
  finalizers:
  - kubedb.com
  generation: 1
  name: redis-quickstart
  namespace: demo
  resourceVersion: "7841"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/redises/redis-quickstart
  uid: 6cc214c9-c53f-11e8-9ba7-0800274bef12
spec:
  mode: Standalone
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
    resources:
      requests:
        storage: 1Gi
    storageClassName: standard
  storageType: Durable
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
  version: 4.0-v1
status:
  observedGeneration: 1$4210395375389091791
  phase: Running
```

Now, you can connect to this database through [redis-cli](https://redis.io/topics/rediscli). In this tutorial, we are connecting to the Redis server from inside of pod.

```console
$ kubectl exec -it redis-quickstart-0 -n demo sh

> redis-cli

127.0.0.1:6379> ping
PONG

#save data
127.0.0.1:6379> SET mykey "Hello"
OK

# view data
127.0.0.1:6379> GET mykey
"Hello"

127.0.0.1:6379> exit
```

## DoNotTerminate Property

When, `terminationPolicy` is `DoNotTerminate`, KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `DoNotTerminate` feature. If admission webhook is enabled, It prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`. You can see this below:

```console
$ kubedb delete rd redis-quickstart -n demo
Error from server (BadRequest): admission webhook "redis.validators.kubedb.com" denied the request: redis "redis-quickstart" can't be paused. To delete, change spec.terminationPolicy
```

Now, run `kubedb edit rd redis-quickstart -n demo` to set `spec.terminationPolicy` to `Pause` (which creates `domantdatabase` when redis is deleted and keeps PVCs intact) or remove this field (which default to `Pause`). Then you will be able to delete/pause the database. 

Learn details of all `TerminationPolicy` [here](/docs/concepts/databases/redis.md#specterminationpolicy)

## Pause Database

When [TerminationPolicy](/docs/concepts/databases/redis.md#specterminationpolicy) is set to `Pause`, it will pause the Redis server instead of deleting it. Here, If you delete the Redis object, KubeDB operator will delete the StatefulSet and its pods but leaves the PVCs unchanged. In KubeDB parlance, we say that `redis-quickstart` Redis server has entered into the dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase object.

```console
$ kubedb delete rd redis-quickstart -n demo
redis.kubedb.com "redis-quickstart" deleted

$ kubedb get drmn -n demo redis-quickstart
NAME               STATUS    AGE
redis-quickstart   Pausing   17s

$ kubedb get drmn -n demo redis-quickstart
NAME               STATUS    AGE
redis-quickstart   Paused    1m
```

```yaml
$ kubedb get drmn -n demo redis-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: 2018-10-01T06:09:58Z
  finalizers:
  - kubedb.com
  generation: 1
  labels:
    kubedb.com/kind: Redis
  name: redis-quickstart
  namespace: demo
  resourceVersion: "8445"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/redis-quickstart
  uid: 9fb52903-c540-11e8-9ba7-0800274bef12
spec:
  origin:
    metadata:
      creationTimestamp: 2018-10-01T06:01:23Z
      name: redis-quickstart
      namespace: demo
    spec:
      redis:
        mode: Standalone
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
          resources:
            requests:
              storage: 1Gi
          storageClassName: standard
        storageType: Durable
        terminationPolicy: Pause
        updateStrategy:
          type: RollingUpdate
        version: 4.0-v1
status:
  observedGeneration: 1$4235806204804343739
  pausingTime: 2018-10-01T06:10:17Z
  phase: Paused
```

Here,

- `spec.origin` is the spec of the original spec of the original Redis object.
- `status.phase` points to the current database state `Paused`.

## Resume Dormant Database

To resume the database from the dormant state, create same `Redis` object with same Spec.

In this tutorial, the dormant database can be resumed by creating original `Redis` object.

The below command will resume the DormantDatabase `redis-quickstart`.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/redis/quickstart/demo-1.yaml
redis.kubedb.com/redis-quickstart created
```

Now, if you exec into the database, you can see that the datas are intact.

## WipeOut DormantDatabase

You can wipe out a DormantDatabase while deleting the objet by setting `spec.wipeOut` to true. KubeDB operator will delete any relevant resources of this `Redis` database (i.e, PVCs, Secrets).

```console
$ kubedb edit drmn -n demo redis-quickstart
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  name: redis-quickstart
  namespace: demo
  ...
spec:
  wipeOut: true
  ...
status:
  phase: Paused
  ...
```

If `spec.wipeOut` is not set to true while deleting the `dormantdatabase` object, then only this object will be deleted and `kubedb-operator` won't delete related Secrets and PVCs. So, users can still access the stored data in the cloud storage buckets (if there is any) as well as PVCs.

## Delete DormantDatabase

As it is already discussed above, `DormantDatabase` can be deleted with or without wiping out the resources. To delete the `dormantdatabase`,

```console
$ kubedb delete drmn redis-quickstart -n demo
dormantdatabase.kubedb.com "redis-quickstart" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo rd/redis-quickstart -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo rd/redis-quickstart

kubectl patch -n demo drmn/redis-quickstart -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/redis-quickstart

kubectl delete ns demo
```

## Tips for Testing

If you are just testing some basic functionalities, you might want to avoid additional hassles due to some safety features that are great for production environment. You can follow these tips to avoid them.

1. **Use `storageType: Ephemeral`**. Databases are precious. You might not want to lose your data in your production environment if database pod fail. So, we recommend to use `spec.storageType: Durable` and provide storage spec in `spec.storage` section. For testing purpose, you can just use `spec.storageType: Ephemeral`. KubeDB will use [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) for storage. You will not require to provide `spec.storage` section.
2. **Use `terminationPolicy: WipeOut`**. It is nice to be able to resume database from previous one. So, we create `DormantDatabase` and preserve all your `PVCs`, `Secrets`, `Snapshots` etc. If you don't want to resume database, you can just use `spec.terminationPolicy: WipeOut`. It will not create `DormantDatabase` and it will delete everything created by KubeDB for a particular Redis crd when you delete the crd. For more details about termination policy, please visit [here](/docs/concepts/databases/redis.md#specterminationpolicy).

## Next Steps

- Monitor your Redis server with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/redis/monitoring/using-coreos-prometheus-operator.md).
- Monitor your Redis server with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/redis/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/redis/private-registry/using-private-registry.md) to deploy Redis with KubeDB.
- Detail concepts of [Redis object](/docs/concepts/databases/redis.md).
- Detail concepts of [RedisVersion object](/docs/concepts/catalog/redis.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
