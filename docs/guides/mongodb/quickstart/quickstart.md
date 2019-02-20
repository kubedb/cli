---
title: MongoDB Quickstart
menu:
  docs_0.9.0:
    identifier: mg-quickstart-quickstart
    name: Overview
    parent: mg-quickstart-mongodb
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# MongoDB QuickStart

This tutorial will show you how to use KubeDB to run a MongoDB database.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/mongodb/mgo-lifecycle.png">
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
  ```

> Note: The yaml files used in this tutorial are stored in [docs/examples/mongodb](https://github.com/kubedb/cli/tree/master/docs/examples/mongodb) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Find Available MongoDBVersion

When you have installed KubeDB, it has created `MongoDBVersion` crd for all supported MongoDB versions. Check 0

```console
$ kubectl get mongodbversions
NAME     VERSION   DB_IMAGE              DEPRECATED   AGE
3.4      3.4       kubedb/mongo:3.4      true         19h
3.4-v1   3.4       kubedb/mongo:3.4-v1   true         19h
3.4-v2   3.4       kubedb/mongo:3.4-v2                19h
3.6      3.6       kubedb/mongo:3.6      true         19h
3.6-v1   3.6       kubedb/mongo:3.6-v1   true         19h
3.6-v2   3.6       kubedb/mongo:3.6-v2                19h
4.0      4.0.5     kubedb/mongo:4.0                   19h
4.0.5    4.0.5     kubedb/mongo:4.0.5                 19h
4.1.7    4.1.7     kubedb/mongo:4.1.7                 19h
```

## Create a MongoDB database

KubeDB implements a `MongoDB` CRD to define the specification of a MongoDB database. Below is the `MongoDB` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-quickstart
  namespace: demo
spec:
  version: "3.4-v2"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  terminationPolicy: DoNotTerminate
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mongodb/quickstart/demo-1.yaml
mongodb.kubedb.com/mgo-quickstart created
```

Here,

- `spec.version` is name of the MongoDBVersion crd where the docker images are specified. In this tutorial, a MongoDB 3.4-v2 database is created.
- `spec.storageType` specifies the type of storage that will be used for MongoDB database. It can be `Durable` or `Ephemeral`. Default value of this field is `Durable`. If `Ephemeral` is used then KubeDB will create MongoDB database using `EmptyDir` volume. In this case, you don't have to specify `spec.storage` field. This is useful for testing purpose.
- `spec.storage` specifies PVC spec that will be dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests.
- `spec.terminationPolicy` gives flexibility whether to `nullify`(reject) the delete operation of `MongoDB` crd or which resources KubeDB should keep or delete when you delete `MongoDB` crd. If admission webhook is enabled, It prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`. Learn details of all `TerminationPolicy` [here]

> Note: spec.storage section is used to create PVC for database pod. It will create PVC with storage size specified instorage.resources.requests field. Don't specify limits here. PVC does not get resized automatically.

KubeDB operator watches for `MongoDB` objects using Kubernetes api. When a `MongoDB` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching MongoDB object name. KubeDB operator will also create a governing service for StatefulSets with the name `<mongodb-name>-gvr`. No MongoDB specific RBAC permission is required in [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb describe mg -n demo mgo-quickstart
Name:               mgo-quickstart
Namespace:          demo
CreationTimestamp:  Wed, 06 Feb 2019 11:25:26 +0600
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
  Name:               mgo-quickstart
  CreationTimestamp:  Wed, 06 Feb 2019 11:25:27 +0600
  Labels:               kubedb.com/kind=MongoDB
                        kubedb.com/name=mgo-quickstart
  Annotations:        <none>
  Replicas:           824637515568 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         mgo-quickstart
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-quickstart
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.111.237.184
  Port:         db  27017/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.8:27017

Service:
  Name:         mgo-quickstart-gvr
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-quickstart
  Annotations:    service.alpha.kubernetes.io/tolerate-unready-endpoints=true
  Type:         ClusterIP
  IP:           None
  Port:         db  27017/TCP
  TargetPort:   27017/TCP
  Endpoints:    172.17.0.8:27017

Database Secret:
  Name:         mgo-quickstart-auth
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-quickstart
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  password:  16 bytes
  username:  4 bytes

No Snapshots.

Events:
  Type    Reason      Age   From              Message
  ----    ------      ----  ----              -------
  Normal  Successful  8m    MongoDB operator  Successfully created Service
  Normal  Successful  8m    MongoDB operator  Successfully created StatefulSet
  Normal  Successful  8m    MongoDB operator  Successfully created MongoDB
  Normal  Successful  8m    MongoDB operator  Successfully created appbinding
  Normal  Successful  8m    MongoDB operator  Successfully patched StatefulSet
  Normal  Successful  8m    MongoDB operator  Successfully patched MongoDB


$ kubectl get statefulset -n demo
NAME             DESIRED   CURRENT   AGE
mgo-quickstart   1         1         20m

$ kubectl get pvc -n demo
NAME                       STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
datadir-mgo-quickstart-0   Bound    pvc-9c738632-29cf-11e9-aebf-080027875192   1Gi        RWO            standard       16m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                           STORAGECLASS   REASON   AGE
pvc-9c738632-29cf-11e9-aebf-080027875192   1Gi        RWO            Delete           Bound    demo/datadir-mgo-quickstart-0   standard                17m

$ kubectl get service -n demo
NAME                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
mgo-quickstart       ClusterIP   10.111.237.184   <none>        27017/TCP   18m
mgo-quickstart-gvr   ClusterIP   None             <none>        27017/TCP   18m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified MongoDB object:

```yaml
$ kubedb get mg -n demo mgo-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  creationTimestamp: "2019-02-06T05:25:26Z"
  finalizers:
  - kubedb.com
  generation: 2
  name: mgo-quickstart
  namespace: demo
  resourceVersion: "67614"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mongodbs/mgo-quickstart
  uid: 9c57c2c0-29cf-11e9-aebf-080027875192
spec:
  databaseSecret:
    secretName: mgo-quickstart-auth
  podTemplate:
    controller: {}
    metadata: {}
    spec:
      livenessProbe:
        exec:
          command:
          - mongo
          - --eval
          - db.adminCommand('ping')
        failureThreshold: 3
        periodSeconds: 10
        successThreshold: 1
        timeoutSeconds: 5
      readinessProbe:
        exec:
          command:
          - mongo
          - --eval
          - db.adminCommand('ping')
        failureThreshold: 3
        periodSeconds: 10
        successThreshold: 1
        timeoutSeconds: 1
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
  terminationPolicy: DoNotTerminate
  updateStrategy:
    type: RollingUpdate
  version: 3.4-v2
status:
  observedGeneration: 2$4213139756412538772
  phase: Running
```

Please note that KubeDB operator has created a new Secret called `mgo-quickstart-auth` *(format: {mongodb-object-name}-auth)* for storing the password for `mongodb` superuser. This secret contains a `username` key which contains the *username* for MongoDB superuser and a `password` key which contains the *password* for MongoDB superuser.

If you want to use custom or existing secret please specify that when creating the MongoDB object using `spec.databaseSecret.secretName`. While creating this secret manually, make sure the secret contains these two keys containing data `user` and `password`. For more details, please see [here](/docs/concepts/databases/mongodb.md#specdatabasesecret).

Now, you can connect to this database through [mongo-shell](https://docs.mongodb.com/v3.4/mongo/). In this tutorial, we are connecting to the MongoDB server from inside the pod.

```console
$ kubectl get secrets -n demo mgo-quickstart-auth -o jsonpath='{.data.\username}' | base64 -d
root

$ kubectl get secrets -n demo mgo-quickstart-auth -o jsonpath='{.data.\password}' | base64 -d
t1QJppJw_B13UIES

$ kubectl exec -it mgo-quickstart-0 -n demo sh

> mongo admin

> db.auth("root","t1QJppJw_B13UIES")
1

> show dbs
admin  0.000GB
local  0.000GB
mydb   0.000GB

> show users
{
    "_id" : "admin.root",
    "user" : "root",
    "db" : "admin",
    "roles" : [
        {
            "role" : "root",
            "db" : "admin"
        }
    ]
}

> use newdb
switched to db newdb

> db.movie.insert({"name":"batman"});
WriteResult({ "nInserted" : 1 })

> db.movie.find().pretty()
{ "_id" : ObjectId("5a2e435d7ec14e7bda785f16"), "name" : "batman" }

> exit
bye
```

## DoNotTerminate Property

When, `terminationPolicy` is `DoNotTerminate`, KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `DoNotTerminate` feature. If admission webhook is enabled, It prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`. You can see this below:

```console
$ kubedb delete mg mgo-quickstart -n demo
Error from server (BadRequest): admission webhook "mongodb.validators.kubedb.com" denied the request: mongodb "mgo-quickstart" can't be paused. To delete, change spec.terminationPolicy
```

Now, run `kubedb edit mg mgo-quickstart -n demo` to set `spec.terminationPolicy` to `Pause` (which creates `domantdatabase` when mongodb is deleted and keeps PVC, snapshots, Secrets intact) or remove this field (which default to `Pause`). Then you will be able to delete/pause the database.

Learn details of all `TerminationPolicy` [here](docs/concepts/databases/mongodb.md#specterminationpolicy)

## Pause Database

When [TerminationPolicy](/docs/concepts/databases/mongodb.md#specterminationpolicy) is set to `Pause`, it will pause the MongoDB database instead of deleting it. Here, If you delete the MongoDB object, KubeDB operator will delete the StatefulSet and its pods but leaves the PVCs unchanged. In KubeDB parlance, we say that `mgo-quickstart` MongoDB database has entered into the dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase object.

```console
$ kubedb delete mg mgo-quickstart -n demo
mongodb.kubedb.com "mgo-quickstart" deleted

$ kubedb get drmn -n demo mgo-quickstart
NAME             STATUS    AGE
mgo-quickstart   Pausing   39s

$ kubedb get drmn -n demo mgo-quickstart
NAME             STATUS    AGE
mgo-quickstart   Paused    21s
```

```yaml
$ kubedb get drmn -n demo mgo-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: "2019-02-06T06:22:25Z"
  finalizers:
  - kubedb.com
  generation: 1
  labels:
    kubedb.com/kind: MongoDB
  name: mgo-quickstart
  namespace: demo
  resourceVersion: "73033"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mgo-quickstart
  uid: 91eb4a12-29d7-11e9-aebf-080027875192
spec:
  origin:
    metadata:
      creationTimestamp: "2019-02-06T06:06:54Z"
      name: mgo-quickstart
      namespace: demo
    spec:
      mongodb:
        databaseSecret:
          secretName: mgo-quickstart-auth
        podTemplate:
          controller: {}
          metadata: {}
          spec:
            livenessProbe:
              exec:
                command:
                - mongo
                - --eval
                - db.adminCommand('ping')
              failureThreshold: 3
              periodSeconds: 10
              successThreshold: 1
              timeoutSeconds: 5
            readinessProbe:
              exec:
                command:
                - mongo
                - --eval
                - db.adminCommand('ping')
              failureThreshold: 3
              periodSeconds: 10
              successThreshold: 1
              timeoutSeconds: 1
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
        version: 3.4-v2
status:
  observedGeneration: 1$16440556888999634490
  pausingTime: "2019-02-06T06:22:37Z"
  phase: Paused
```

Here,

- `spec.origin` is the spec of the original spec of the original MongoDB object.
- `status.phase` points to the current database state `Paused`.

## Resume Dormant Database

To resume the database from the dormant state, create same `MongoDB` object with same Spec.

In this tutorial, the dormant database can be resumed by creating original MongoDB object.

The below command will resume the DormantDatabase `mgo-quickstart`.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mongodb/quickstart/demo-1.yaml
mongodb.kubedb.com/mgo-quickstart created
```

Now, if you exec into the database, you can see that the datas are intact.

## WipeOut DormantDatabase

You can wipe out a DormantDatabase while deleting the object by setting `spec.wipeOut` to true. KubeDB operator will delete any relevant resources of this `MongoDB` database (i.e, PVCs, Secrets, Snapshots). It will also delete snapshot data stored in the Cloud Storage buckets.

```yaml
$ kubedb edit drmn -n demo mgo-quickstart
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  name: mgo-quickstart
  namespace: demo
  ...
spec:
  wipeOut: true
  ...
status:
  phase: Paused
  ...
```

If `spec.wipeOut` is not set to true while deleting the `dormantdatabase` object, then only this object will be deleted and `kubedb-operator` won't delete related Secrets, PVCs, and Snapshots. So, users still can access the stored data in the cloud storage buckets as well as PVCs.

## Delete DormantDatabase

As it is already discussed above, `DormantDatabase` can be deleted with or without wiping out the resources. To delete the `dormantdatabase`,

```console
$ kubectl delete drmn mgo-quickstart -n demo
dormantdatabase "mgo-quickstart" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mg/mgo-quickstart -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mg/mgo-quickstart

kubectl patch -n demo drmn/mgo-quickstart -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mgo-quickstart

kubectl delete ns demo
```

## Tips for Testing

If you are just testing some basic functionalities, you might want to avoid additional hassles due to some safety features that are great for production environment. You can follow these tips to avoid them.

1. **Use `storageType: Ephemeral`**. Databases are precious. You might not want to lose your data in your production environment if database pod fail. So, we recommend to use `spec.storageType: Durable` and provide storage spec in `spec.storage` section. For testing purpose, you can just use `spec.storageType: Ephemeral`. KubeDB will use [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) for storage. You will not require to provide `spec.storage` section.
2. **Use `terminationPolicy: WipeOut`**. It is nice to be able to resume database from previous one. So, we create `DormantDatabase` and preserve all your `PVCs`, `Secrets`, `Snapshots` etc. If you don't want to resume database, you can just use `spec.terminationPolicy: WipeOut`. It will not create `DormantDatabase` and it will delete everything created by KubeDB for a particular MongoDB crd when you delete the crd. For more details about termination policy, please visit [here](/docs/concepts/databases/mongodb.md#specterminationpolicy).

## Next Steps

- [Snapshot and Restore](/docs/guides/mongodb/snapshot/backup-and-restore.md) process of MongoDB databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mongodb/snapshot/scheduled-backup.md) of MongoDB databases using KubeDB.
- Initialize [MongoDB with Script](/docs/guides/mongodb/initialization/using-script.md).
- Initialize [MongoDB with Snapshot](/docs/guides/mongodb/initialization/using-snapshot.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mongodb/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mongodb/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mongodb/private-registry/using-private-registry.md) to deploy MongoDB with KubeDB.
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Detail concepts of [MongoDBVersion object](/docs/concepts/catalog/mongodb.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
