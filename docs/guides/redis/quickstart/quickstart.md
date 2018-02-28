---
title: Redis Quickstart
menu:
  docs_0.8.0-beta.2:
    identifier: rd-quickstart-quickstart
    name: Overview
    parent: rd-quickstart-redis
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Redis QuickStart

This tutorial will show you how to use KubeDB to run a Redis database.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/redis/redis-lifecycle.png" width="600" height="373">
</p>

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/redis/demo-0.yaml
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    45m
demo          Active    10s
kube-public   Active    45m
kube-system   Active    45m
```

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create a Redis database

KubeDB implements a `Redis` CRD to define the specification of a Redis database. Below is the `Redis` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  name: redis-quickstart
  namespace: demo
spec:
  version: 4
  doNotPause: true
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/redis/quickstart/demo-1.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/redis/quickstart/demo-1.yaml"
redis "redis-quickstart" created
```

Here,

- `spec.version` is the version of Redis database. In this tutorial, a Redis 4 database is going to be created.

- `spec.doNotPause` tells KubeDB operator that if this object is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.

- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

KubeDB operator watches for `Redis` objects using Kubernetes api. When a `Redis` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching Redis object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. No Redis specific RBAC permission is required in [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb get rd -n demo
NAME               STATUS    AGE
redis-quickstart   Running   1m


$ kubedb describe rd -n demo redis-quickstart
Name:		redis-quickstart
Namespace:	demo
StartTimestamp:	Mon, 12 Feb 2018 16:41:39 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			redis-quickstart
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Mon, 12 Feb 2018 16:41:41 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		redis-quickstart
  Type:		ClusterIP
  IP:		10.101.253.6
  Port:		db	6379/TCP

Events:
  FirstSeen   LastSeen   Count     From             Type       Reason       Message
  ---------   --------   -----     ----             --------   ------       -------
  1m          1m         1         Redis operator   Normal     Successful   Successfully created StatefulSet
  1m          1m         1         Redis operator   Normal     Successful   Successfully created Redis
  1m          1m         1         Redis operator   Normal     Successful   Successfully created Service



$ kubectl get statefulset -n demo
NAME               DESIRED   CURRENT   AGE
redis-quickstart   1         1         2m


$ kubectl get pvc -n demo
NAME                      STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-redis-quickstart-0   Bound     pvc-4fbc09fb-0fe1-11e8-a2d6-08002751ae8c   50Mi       RWO            standard       2m


$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                          STORAGECLASS   REASON    AGE
pvc-4fbc09fb-0fe1-11e8-a2d6-08002751ae8c   50Mi       RWO            Delete           Bound     demo/data-redis-quickstart-0   standard                 3m


$ kubectl get service -n demo
NAME               TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
kubedb             ClusterIP   None           <none>        <none>     3m
redis-quickstart   ClusterIP   10.101.253.6   <none>        6379/TCP   3m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified Redis object:

```yaml
$ kubedb get rd -n demo redis-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-12T10:41:39Z
  finalizers:
  - kubedb.com
  generation: 0
  name: redis-quickstart
  namespace: demo
  resourceVersion: "46523"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/redises/redis-quickstart
  uid: 4ecf9d7c-0fe1-11e8-a2d6-08002751ae8c
spec:
  doNotPause: true
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: 4
status:
  creationTime: 2018-02-12T10:41:40Z
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

## Pause Database

KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `doNotPause` feature. If admission webhook is enabled, It prevents user from deleting the database as long as the `spec.doNotPause` is set to true. Since the Redis object created in this tutorial has `spec.doNotPause` set to true, if you delete the Redis object, KubeDB operator will nullify the delete operation. You can see this below:

```console
$ kubedb delete rd redis-quickstart -n demo
error: Redis "redis-quickstart" can't be paused. To continue delete, unset spec.doNotPause and retry.
```

Now, run `kubedb edit rd redis-quickstart -n demo` to set `spec.doNotPause` to false or remove this field (which default to false). Then if you delete the Redis object, KubeDB operator will delete the StatefulSet and its pods, but leaves the PVCs unchanged. In KubeDB parlance, we say that `redis-quickstart` Redis database has entered into dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase object.

```console
$ kubedb delete rd redis-quickstart -n demo
redis "redis-quickstart" deleted


$ kubedb get drmn -n demo redis-quickstart
NAME               STATUS    AGE
redis-quickstart   Pausing   6s


$ kubedb get drmn -n demo redis-quickstart
NAME               STATUS    AGE
redis-quickstart   Paused    10s
```

```yaml
$ kubedb get drmn -n demo redis-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-12T10:48:07Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: Redis
  name: redis-quickstart
  namespace: demo
  resourceVersion: "46767"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/redis-quickstart
  uid: 360017bb-0fe2-11e8-a2d6-08002751ae8c
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: redis-quickstart
      namespace: demo
    spec:
      redis:
        resources: {}
        storage:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 50Mi
          storageClassName: standard
        version: "4"
status:
  creationTime: 2018-02-12T10:48:08Z
  pausingTime: 2018-02-12T10:48:15Z
  phase: Paused
```

Here,

- `spec.origin` is the spec of the original spec of the original Redis object.

- `status.phase` points to the current database state `Paused`.

## Resume Dormant Database

To resume the database from the dormant state, set `spec.resume` to `true` in the DormantDatabase object.

```yaml
$ kubedb edit drmn -n demo redis-quickstart
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  name: redis-quickstart
  namespace: demo
  ...
spec:
  resume: true
  ...
status:
  phase: Paused
  ...
```

KubeDB operator will notice that `spec.resume` is set to true. KubeDB operator will delete the DormantDatabase object and create a new Redis object using the original spec. This will in turn start a new StatefulSet which will mount the originally created PVCs. Thus the original database is resumed.

Please note that the dormant database can also be resumed by creating same `Redis` database by using same Specs. In this tutorial, the dormant database can be resumed by creating `Redis` database using demo-1.yaml file. The below command resumes the dormant database `redis-quickstart` that was created before.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/redis/quickstart/demo-1.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/redis/quickstart/demo-1.yaml"
redis "redis-quickstart" created
```

## Wipeout Dormant Database

You can also wipe out a DormantDatabase by setting `spec.wipeOut` to true. KubeDB operator will delete the PVCs that is used by `redis-quickstart`. There is no way to resume a wiped out database. So, be sure before you wipe out a database.

Create dormant database again and set `spec.wipeOut` to true:

```yaml
$ kubedb delete rd redis-quickstart -n demo
redis "redis-quickstart" deleted

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


$ kubedb get drmn -n demo redis-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-12T10:56:24Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: Redis
  name: redis-quickstart
  namespace: demo
  resourceVersion: "47168"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/redis-quickstart
  uid: 5e085768-0fe3-11e8-a2d6-08002751ae8c
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: redis-quickstart
      namespace: demo
    spec:
      redis:
        resources: {}
        storage:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 50Mi
          storageClassName: standard
        version: "4"
  wipeOut: true
status:
  creationTime: 2018-02-12T10:56:24Z
  pausingTime: 2018-02-12T10:56:35Z
  phase: WipedOut
  wipeOutTime: 2018-02-12T10:57:24Z


$ kubedb get drmn -n demo
NAME               STATUS     AGE
redis-quickstart   WipedOut   2m
```

## Delete Dormant Database

You still have a record that there used to be a Redis database `redis-quickstart` in the form of a DormantDatabase database `redis-quickstart`. Since you have already wiped out the database, you can delete the DormantDatabase object.

```console
$ kubedb delete drmn redis-quickstart -n demo
dormantdatabase "redis-quickstart" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete rd redis-quickstart -n demo --force
$ kubedb delete drmn redis-quickstart -n demo --force

# or
# $ kubedb delete rd,drmn -n demo --all --force

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Monitor your Redis database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/redis/monitoring/using-coreos-prometheus-operator.md).
- Monitor your Redis database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/redis/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/redis/private-registry/using-private-registry.md) to deploy Redis with KubeDB.
- Detail concepts of [Redis object](/docs/concepts/databases/redis.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
