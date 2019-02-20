---
title: MongoDB ReplicaSet Guide
menu:
  docs_0.9.0:
    identifier: mg-clustering-replicaset
    name: ReplicaSet Guide
    parent: mg-clustering-mongodb
    weight: 15
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# KubeDB - MongoDB ReplicaSet

This tutorial will show you how to use KubeDB to run a MongoDB ReplicaSet.

## Before You Begin

Before proceeding:

- Read [mongodb replication concept](/docs/guides/mongodb/clustering/replication_concept.md) to learn about MongoDB Replica Set clustering.

- You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace/demo created
  ```

> Note: The yaml files used in this tutorial are stored in [docs/examples/mongodb](https://github.com/kubedb/cli/tree/master/docs/examples/mongodb) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Deploy MongoDB ReplicaSet

To deploy a MongoDB ReplicaSet, user have to specify `spec.replicaSet` option in `Mongodb` CRD.

The following is an example of a `Mongodb` object which creates MongoDB ReplicaSet of three members.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-replicaset
  namespace: demo
spec:
  version: "3.6-v2"
  replicas: 3
  replicaSet:
    name: rs0
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mongodb/clustering/demo-1.yaml
mongodb.kubedb.com/mgo-replicaset created
```

Here,

- `spec.replicaSet` represents the configuration for replicaset.
  - `name` denotes the name of mongodb replicaset.
  - `KeyFileSecret` (optional) is a secret name that contains keyfile (a random string)against `key.txt` key. Each mongod instances in the replica set uses the contents of the keyfile as the shared password for authenticating other members in the deployment. Only mongod instances with the correct keyfile can join the replica set. _User can provide the `KeyFileSecret` by creating a secret with key `key.txt`. See [here](https://docs.mongodb.com/manual/tutorial/enforce-keyfile-access-control-in-existing-replica-set/#create-a-keyfile) to create the string for `KeyFileSecret`._ If `KeyFileSecret` is not given, KubeDB operator will generate a `KeyFileSecret` itself.
- `spec.replicas` denotes the number of members in `rs0` mongodb replicaset.
- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. So, each members will have a pod of this storage configuration. You can specify any StorageClass available in your cluster with appropriate resource requests.

KubeDB operator watches for `MongoDB` objects using Kubernetes api. When a `MongoDB` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching MongoDB object name. KubeDB operator will also create a governing service for StatefulSets with the name `<mongodb-name>-gvr`. No MongoDB specific RBAC permission is required in [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb describe mg -n demo mgo-replicaset
Name:               mgo-replicaset
Namespace:          demo
CreationTimestamp:  Wed, 06 Feb 2019 16:08:15 +0600
Labels:             <none>
Annotations:        <none>
Replicas:           3  total
Status:             Running
  StorageType:      Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:
  Name:               mgo-replicaset
  CreationTimestamp:  Wed, 06 Feb 2019 16:08:15 +0600
  Labels:               kubedb.com/kind=MongoDB
                        kubedb.com/name=mgo-replicaset
  Annotations:        <none>
  Replicas:           824637261744 desired | 3 total
  Pods Status:        3 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         mgo-replicaset
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-replicaset
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.97.174.220
  Port:         db  27017/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.7:27017,172.17.0.8:27017,172.17.0.9:27017

Service:
  Name:         mgo-replicaset-gvr
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-replicaset
  Annotations:    service.alpha.kubernetes.io/tolerate-unready-endpoints=true
  Type:         ClusterIP
  IP:           None
  Port:         db  27017/TCP
  TargetPort:   27017/TCP
  Endpoints:    172.17.0.7:27017,172.17.0.8:27017,172.17.0.9:27017

Database Secret:
  Name:         mgo-replicaset-auth
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mgo-replicaset
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  password:  16 bytes
  username:  4 bytes

No Snapshots.

Events:
  Type    Reason      Age   From             Message
  ----    ------      ----  ----             -------
  Normal  Successful  1m    KubeDB operator  Successfully created Service
  Normal  Successful  26s   KubeDB operator  Successfully created StatefulSet
  Normal  Successful  26s   KubeDB operator  Successfully created MongoDB
  Normal  Successful  26s   KubeDB operator  Successfully created appbinding
  Normal  Successful  26s   KubeDB operator  Successfully patched StatefulSet
  Normal  Successful  26s   KubeDB operator  Successfully patched MongoDB


$ kubectl get statefulset -n demo
NAME             READY   AGE
mgo-replicaset   3/3     105s

$ kubectl get pvc -n demo
NAME                       STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
datadir-mgo-replicaset-0   Bound     pvc-597784c9-c093-11e8-b4a9-0800272618ed   1Gi        RWO            standard       1h
datadir-mgo-replicaset-1   Bound     pvc-8ca7a9d9-c093-11e8-b4a9-0800272618ed   1Gi        RWO            standard       1h
datadir-mgo-replicaset-2   Bound     pvc-b7d8a624-c093-11e8-b4a9-0800272618ed   1Gi        RWO            standard       1h

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                           STORAGECLASS   REASON    AGE
pvc-597784c9-c093-11e8-b4a9-0800272618ed   1Gi        RWO            Delete           Bound     demo/datadir-mgo-replicaset-0   standard                 1h
pvc-8ca7a9d9-c093-11e8-b4a9-0800272618ed   1Gi        RWO            Delete           Bound     demo/datadir-mgo-replicaset-1   standard                 1h
pvc-b7d8a624-c093-11e8-b4a9-0800272618ed   1Gi        RWO            Delete           Bound     demo/datadir-mgo-replicaset-2   standard                 1h

$ kubectl get service -n demo
NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)     AGE
mgo-replicaset       ClusterIP   10.97.174.220   <none>        27017/TCP   119s
mgo-replicaset-gvr   ClusterIP   None            <none>        27017/TCP   119s
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified MongoDB object:

```yaml
$ kubedb get mg -n demo mgo-replicaset -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  creationTimestamp: "2019-02-06T10:08:15Z"
  finalizers:
  - kubedb.com
  generation: 3
  name: mgo-replicaset
  namespace: demo
  resourceVersion: "92643"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mongodbs/mgo-replicaset
  uid: 1e8cd706-29f7-11e9-aebf-080027875192
spec:
  databaseSecret:
    secretName: mgo-replicaset-auth
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
  replicaSet:
    keyFile:
      secretName: mgo-replicaset-keyfile
    name: rs0
  replicas: 3
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
  version: 3.6-v2
status:
  observedGeneration: 3$4212299729528774793
  phase: Running
```

Please note that KubeDB operator has created a new Secret called `mgo-replicaset-auth` *(format: {mongodb-object-name}-auth)* for storing the password for `mongodb` superuser. This secret contains a `username` key which contains the *username* for MongoDB superuser and a `password` key which contains the *password* for MongoDB superuser.

If you want to use custom or existing secret please specify that when creating the MongoDB object using `spec.databaseSecret.secretName`. While creating this secret manually, make sure the secret contains these two keys containing data `username` and `password`. For more details, please see [here](/docs/concepts/databases/mongodb.md#specdatabasesecret).

## Redundancy and Data Availability

Now, you can connect to this database through [mongo-shell](https://docs.mongodb.com/v3.4/mongo/). In this tutorial, we will insert document on primary member, and we will see if the data becomes available on secondary members.

At first, insert data inside primary member `rs0:PRIMARY`.

```console
$ kubectl get secrets -n demo mgo-replicaset-auth -o jsonpath='{.data.\username}' | base64 -d
root

$ kubectl get secrets -n demo mgo-replicaset-auth -o jsonpath='{.data.\password}' | base64 -d
5O4R2ze2bWXcWsdP

$ kubectl exec -it mgo-replicaset-0 -n demo bash

mongodb@mgo-replicaset-0:/$ mongo admin -u root -p 5O4R2ze2bWXcWsdP
MongoDB shell version v3.6.6
connecting to: mongodb://127.0.0.1:27017/admin
MongoDB server version: 3.6.6
Welcome to the MongoDB shell.

rs0:PRIMARY> > rs.isMaster().primary
mgo-replicaset-0.mgo-replicaset-gvr.demo.svc.cluster.local:27017

rs0:PRIMARY> show dbs
admin   0.000GB
config  0.000GB
local   0.000GB

rs0:PRIMARY> show users
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


rs0:PRIMARY> use newdb
switched to db newdb

rs0:PRIMARY> db.movie.insert({"name":"batman"});
WriteResult({ "nInserted" : 1 })

rs0:PRIMARY> db.movie.find().pretty()
{ "_id" : ObjectId("5b5efeea9d097ca0600694a3"), "name" : "batman" }

rs0:PRIMARY> exit
bye
```

Now, check the redundancy and data availability in secondary members.
We will exec in `mgo-replicaset-1`(which is secondary member right now) to check the data availability.

```console
$ kubectl exec -it mgo-replicaset-1 -n demo bash
mongodb@mgo-replicaset-1:/$ mongo admin -u root -p 5O4R2ze2bWXcWsdP
MongoDB shell version v3.6.6
connecting to: mongodb://127.0.0.1:27017/admin
MongoDB server version: 3.6.6
Welcome to the MongoDB shell.

rs0:SECONDARY> rs.slaveOk()
rs0:SECONDARY> > show dbs
admin   0.000GB
config  0.000GB
local   0.000GB
newdb   0.000GB

rs0:SECONDARY> show users
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

rs0:SECONDARY> use newdb
switched to db newdb

rs0:SECONDARY> db.movie.find().pretty()
{ "_id" : ObjectId("5b5efeea9d097ca0600694a3"), "name" : "batman" }

rs0:SECONDARY> exit
bye

```

## Automatic Failover

To test automatic failover, we will force the primary member to restart. As the primary member (`pod`) becomes unavailable, the rest of the members will elect a primary member by election.

```console
$ kubectl get pods -n demo
NAME               READY     STATUS    RESTARTS   AGE
mgo-replicaset-0   1/1       Running   0          1h
mgo-replicaset-1   1/1       Running   0          1h
mgo-replicaset-2   1/1       Running   0          1h

$ kubectl delete pod -n demo mgo-replicaset-0
pod "mgo-replicaset-0" deleted

$ kubectl get pods -n demo
NAME               READY     STATUS        RESTARTS   AGE
mgo-replicaset-0   1/1       Terminating   0          1h
mgo-replicaset-1   1/1       Running       0          1h
mgo-replicaset-2   1/1       Running       0          1h

```

Now verify the automatic failover, Let's exec in `mgo-replicaset-1` pod,

```console
$ kubectl exec -it mgo-replicaset-1 -n demo bash
mongodb@mgo-replicaset-1:/$ mongo admin -u root -p 5O4R2ze2bWXcWsdP
MongoDB shell version v3.6.6
connecting to: mongodb://127.0.0.1:27017/admin
MongoDB server version: 3.6.6
Welcome to the MongoDB shell.

rs0:SECONDARY> rs.isMaster().primary
mgo-replicaset-2.mgo-replicaset-gvr.demo.svc.cluster.local:27017

# Also verify, data persistency
rs0:SECONDARY> rs.slaveOk()
rs0:SECONDARY> > show dbs
admin   0.000GB
config  0.000GB
local   0.000GB
newdb   0.000GB

rs0:SECONDARY> show users
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

rs0:SECONDARY> use newdb
switched to db newdb

rs0:SECONDARY> db.movie.find().pretty()
{ "_id" : ObjectId("5b5efeea9d097ca0600694a3"), "name" : "batman" }
```

## Pause Database

When, `terminationPolicy` is `DoNotTerminate`, KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `DoNotTerminate` feature. If admission webhook is enabled, It prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`.

Since the MongoDB object created in this tutorial has `spec.terminationPolicy` set to `Pause` (default), if you delete the MongoDB object, KubeDB operator will create a dormant database while deleting the StatefulSet and its pods but leaves the PVCs unchanged.

```console
$ kubedb delete mg mgo-replicaset -n demo
mongodb.kubedb.com "mgo-replicaset" deleted

$ kubedb get drmn -n demo mgo-replicaset
NAME             STATUS    AGE
mgo-replicaset   Pausing   25s

$ kubedb get drmn -n demo mgo-replicaset
NAME             STATUS    AGE
mgo-replicaset   Paused    1m
```

```yaml
$ kubedb get drmn -n demo mgo-replicaset -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: "2019-02-06T10:20:21Z"
  finalizers:
  - kubedb.com
  generation: 1
  labels:
    kubedb.com/kind: MongoDB
  name: mgo-replicaset
  namespace: demo
  resourceVersion: "93530"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mgo-replicaset
  uid: cf0829ce-29f8-11e9-aebf-080027875192
spec:
  origin:
    metadata:
      creationTimestamp: "2019-02-06T10:08:15Z"
      name: mgo-replicaset
      namespace: demo
    spec:
      mongodb:
        databaseSecret:
          secretName: mgo-replicaset-auth
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
        replicaSet:
          keyFile:
            secretName: mgo-replicaset-keyfile
          name: rs0
        replicas: 3
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
        version: 3.6-v2
status:
  observedGeneration: 1$16440556888999634490
  pausingTime: "2019-02-06T10:20:27Z"
  phase: Paused
```

Here,

- `spec.origin` is the spec of the original spec of the original MongoDB object.
- `status.phase` points to the current database state `Paused`.

## Resume Dormant Database

To resume the database from the dormant state, create same `MongoDB` object with same Spec.

In this tutorial, the dormant database can be resumed by creating original MongoDB object.

The below command will resume the DormantDatabase `mgo-replicaset`.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mongodb/clustering/demo-1.yaml
mongodb.kubedb.com/mgo-replicaset created
```

Now, If you again exec into `pod` and look for previous data, you will see that, all the data persists.

## WipeOut DormantDatabase

You can wipe out a DormantDatabase while deleting the object by setting `spec.wipeOut` to true. KubeDB operator will delete any relevant resources of this `MongoDB` database (i.e, PVCs, Secrets, Snapshots). It will also delete snapshot data stored in the Cloud Storage buckets.

```yaml
$ kubedb edit drmn -n demo mgo-replicaset
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  name: mgo-replicaset
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
$ kubectl delete drmn mgo-replicaset -n demo
dormantdatabase "mgo-replicaset" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mg/mgo-replicaset -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mg/mgo-replicaset

kubectl patch -n demo drmn/mgo-replicaset -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mgo-replicaset

kubectl delete ns demo
```

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
