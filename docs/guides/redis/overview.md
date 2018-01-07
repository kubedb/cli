---
title: Redis
menu:
  docs_0.8.0-beta.0:
    identifier: guides-redis-readme
    name: Overview
    parent: guides-redis
    weight: 10
menu_name: docs_0.8.0-beta.0
section_menu_id: guides
aliases:
  - /docs/0.8.0-beta.0/guides/redis/
---

> New to KubeDB? Please start [here](/docs/guides/README.md).

# Running Redis
This tutorial will show you how to use KubeDB to run an Redis database.

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f ./docs/examples/redis/demo-0.yaml
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    23h
demo          Active    35s
kube-public   Active    23h
kube-system   Active    23h
```

## Create an Redis database
KubeDB implements a `Redis` CRD to define the specification of an Redis database. Below is the `Redis` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  name: r1
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


$ kubedb create -f ./docs/examples/redis/demo-1.yaml
validating "./docs/examples/redis/demo-1.yaml"
redis "r1" created
```

Here,
 - `spec.version` is the version of Redis database. In this tutorial, an Redis 4 database is going to be created.

 - `spec.doNotPause` tells KubeDB operator that if this object is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.

 - `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

KubeDB operator watches for `Redis` objects using Kubernetes api. When a `Redis` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching Redis object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. If [RBAC is enabled](/docs/guides/rbac.md), a ClusterRole, ServiceAccount and ClusterRoleBinding with the matching object name will be created and used as the service account name for the corresponding StatefulSet.

```console
$ kubedb describe rd r1 -n demo
Name:		r1
Namespace:	demo
StartTimestamp:	Tue, 12 Dec 2017 12:02:05 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			r1
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Tue, 12 Dec 2017 12:02:15 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		r1
  Type:		ClusterIP
  IP:		10.102.1.255
  Port:		db	6379/TCP

Events:
  FirstSeen   LastSeen   Count     From             Type       Reason               Message
  ---------   --------   -----     ----             --------   ------               -------
  2m          2m         1         Redis operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  2m          2m         1         Redis operator   Normal     SuccessfulCreate     Successfully created Redis
  2m          2m         1         Redis operator   Normal     SuccessfulValidate   Successfully validate Redis
  2m          2m         1         Redis operator   Normal     Creating             Creating Kubernetes objects


$ kubectl get statefulset -n demo
NAME      DESIRED   CURRENT   AGE
r1        1         1         3m


$ kubectl get pvc -n demo
NAME        STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-r1-0   Bound     pvc-0117d2e5-df02-11e7-9e8f-0800279fc284   50Mi       RWO            standard       3m


$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM            STORAGECLASS   REASON    AGE
pvc-0117d2e5-df02-11e7-9e8f-0800279fc284   50Mi       RWO            Delete           Bound     demo/data-r1-0   standard                 3m


$ kubectl get service -n demo
NAME      TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
kubedb    ClusterIP   None           <none>        <none>     4m
r1        ClusterIP   10.102.1.255   <none>        6379/TCP   4m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified Redis object:

```yaml
$ kubedb get rd -n demo r1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-12T06:02:05Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  name: r1
  namespace: demo
  resourceVersion: "24131"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/redises/r1
  uid: fb09548e-df01-11e7-9e8f-0800279fc284
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
  creationTime: 2017-12-12T06:02:05Z
  phase: Running
```

Now, you can connect to this database through [redis-cli](https://redis.io/topics/rediscli). In this tutorial, we are connecting to the Redis server from inside of pod.
```console
$ kubectl exec -it r1-0 -n demo sh

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

Since the Redis object created in this tutorial has `spec.doNotPause` set to true, if you delete the Redis object, KubeDB operator will recreate the object and essentially nullify the delete operation. You can see this below:

```console
$ kubedb delete rd r1 -n demo
error: Redis "r1" can't be paused. To continue delete, unset spec.doNotPause and retry.
```

Now, run `kubedb edit rd r1 -n demo` to set `spec.doNotPause` to false or remove this field (which default to false). Then if you delete the Redis object, KubeDB operator will delete the StatefulSet and its pods, but leaves the PVCs unchanged. In KubeDB parlance, we say that `r1` Redis database has entered into dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase object.

```yaml
$ kubedb delete rd -n demo r1
redis "r1" deleted

$ kubedb get drmn -n demo r1
NAME      STATUS    AGE
r1        Pausing   9s

$ kubedb get drmn -n demo r1
NAME      STATUS    AGE
r1        Paused    50s

$ kubedb get drmn -n demo r1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-12T06:20:46Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: Redis
  name: r1
  namespace: demo
  resourceVersion: "25424"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/r1
  uid: 96f5512b-df04-11e7-9e8f-0800279fc284
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: r1
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
  creationTime: 2017-12-12T06:20:46Z
  pausingTime: 2017-12-12T06:21:36Z
  phase: Paused
```

Here,
 - `spec.origin` is the spec of the original spec of the original Redis object.

 - `status.phase` points to the current database state `Paused`.


## Resume Dormant Database

To resume the database from the dormant state, set `spec.resume` to `true` in the DormantDatabase object.

```yaml
$ kubedb edit drmn -n demo r1
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-12T06:20:46Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: Redis
  name: r1
  namespace: demo
  resourceVersion: "25424"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/r1
  uid: 96f5512b-df04-11e7-9e8f-0800279fc284
spec:
  resume: true
  origin:
    metadata:
      creationTimestamp: null
      name: r1
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
  creationTime: 2017-12-12T06:20:46Z
  pausingTime: 2017-12-12T06:21:36Z
  phase: Paused
```
KubeDB operator will notice that `spec.resume` is set to true. KubeDB operator will delete the DormantDatabase object and create a new Redis object using the original spec. This will in turn start a new StatefulSet which will mount the originally created PVCs. Thus the original database is resumed.


## Wipeout Dormant Database
You can also wipe out a DormantDatabase by setting `spec.wipeOut` to true. KubeDB operator will delete the PVCs once the `spec.wipeOut` is set to `true`. There is no way to resume a wiped out database. So, be sure before you wipe out a database.

```yaml
$ kubedb edit drmn -n demo r1
# set spec.wipeOut: true

$ kubedb get drmn -n demo r1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-12T06:25:29Z
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: Redis
  name: r1
  namespace: demo
  resourceVersion: "25827"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/r1
  uid: 3f71dac4-df05-11e7-9e8f-0800279fc284
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: r1
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
  creationTime: 2017-12-12T06:25:29Z
  pausingTime: 2017-12-12T06:26:39Z
  phase: WipedOut
  wipeOutTime: 2017-12-12T06:26:39Z


$ kubedb get drmn -n demo
NAME      STATUS     AGE
r1        WipedOut   1m
```


## Delete Dormant Database
You still have a record that there used to be an Redis database `r1` in the form of a DormantDatabase database `r1`. Since you have already wiped out the database, you can delete the DormantDatabase object.

```console
$ kubedb delete drmn r1 -n demo
dormantdatabase "r1" deleted
```

## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall KubeDB operator, please follow the steps [here](/docs/setup/uninstall.md).


## Next Steps
- Learn about the details of Redis object [here](/docs/concepts/databases/redis.md).
- Thinking about monitoring your database? KubeDB works [out-of-the-box with Prometheus](/docs/guides/monitoring.md).
- Learn how to use KubeDB in a [RBAC](/docs/guides/rbac.md) enabled cluster.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
