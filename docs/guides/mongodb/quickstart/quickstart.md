---
title: MongoDB Quickstart
menu:
  docs_0.8.0-beta.2:
    identifier: mg-quickstart-quickstart
    name: Overview
    parent: mg-quickstart-mongodb
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# MongoDB QuickStart

This tutorial will show you how to use KubeDB to run a MongoDB database.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/mongodb/mgo-lifecycle.png" width="600" height="373">
</p>

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/mongodb/demo-0.yaml
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    45m
demo          Active    10s
kube-public   Active    45m
kube-system   Active    45m
```

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create a MongoDB database

KubeDB implements a `MongoDB` CRD to define the specification of a MongoDB database. Below is the `MongoDB` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-quickstart
  namespace: demo
spec:
  version: 3.4
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
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/mongodb/quickstart/demo-1.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/mongodb/quickstart/demo-1.yaml"
mongodb "mgo-quickstart" created
```

Here,

- `spec.version` is the version of MongoDB database. In this tutorial, a MongoDB 3.4 database is going to be created.
- `spec.doNotPause` tells KubeDB operator that if this object is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.
- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

KubeDB operator watches for `MongoDB` objects using Kubernetes api. When a `MongoDB` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching MongoDB object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. No MongoDB specific RBAC permission is required in [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb describe mg -n demo mgo-quickstart
Name:		mgo-quickstart
Namespace:	demo
StartTimestamp:	Fri, 02 Feb 2018 15:11:58 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			mgo-quickstart
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Fri, 02 Feb 2018 15:11:24 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		mgo-quickstart
  Type:		ClusterIP
  IP:		10.103.114.139
  Port:		db	27017/TCP

Database Secret:
  Name:	mgo-quickstart-auth
  Type:	Opaque
  Data
  ====
  password:	16 bytes
  user:		4 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From               Type       Reason       Message
  ---------   --------   -----     ----               --------   ------       -------
  2m          2m         1         MongoDB operator   Normal     Successful   Successfully patched StatefulSet
  2m          2m         1         MongoDB operator   Normal     Successful   Successfully patched MongoDB
  2m          2m         1         MongoDB operator   Normal     Successful   Successfully patched StatefulSet
  2m          2m         1         MongoDB operator   Normal     Successful   Successfully patched MongoDB


$ kubectl get statefulset -n demo
NAME             DESIRED   CURRENT   AGE
mgo-quickstart   1         1         4m

$ kubectl get pvc -n demo
NAME                    STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-mgo-quickstart-0   Bound     pvc-16158aae-07fa-11e8-946f-080027c05a6e   50Mi       RWO            standard       2m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                        STORAGECLASS   REASON    AGE
pvc-16158aae-07fa-11e8-946f-080027c05a6e   50Mi       RWO            Delete           Bound     demo/data-mgo-quickstart-0   standard                 3m

$ kubectl get service -n demo
NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
kubedb           ClusterIP   None             <none>        <none>      3m
mgo-quickstart   ClusterIP   10.107.133.189   <none>        27017/TCP   3m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified MongoDB object:

```yaml
$ kubedb get mg -n demo mgo-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-02T09:18:39Z
  finalizers:
  - kubedb.com
  generation: 0
  name: mgo-quickstart
  namespace: demo
  resourceVersion: "46856"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mongodbs/mgo-quickstart
  uid: 0de4d2a2-07fa-11e8-946f-080027c05a6e
spec:
  databaseSecret:
    secretName: mgo-quickstart-auth
  doNotPause: true
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: 3.4
status:
  creationTime: 2018-02-02T09:18:50Z
  phase: Running
```

Please note that KubeDB operator has created a new Secret called `mgo-quickstart-auth` *(format: {mongodb-object-name}-auth)* for storing the password for `mongodb` superuser. This secret contains a `user` key which contains the *username* for MongoDB superuser and a `password` key which contains the *password* for MongoDB superuser.
If you want to use an existing secret please specify that when creating the MongoDB object using `spec.databaseSecret.secretName`. While creating this secret manually, make sure the secret contains these two keys containing data `user` and `password`.

Now, you can connect to this database through [mongo-shell](https://docs.mongodb.com/v3.4/mongo/). In this tutorial, we are connecting to the MongoDB server from inside the pod.

```console
$ kubectl get secrets -n demo mgo-quickstart-auth -o jsonpath='{.data.\user}' | base64 -d
root

$ kubectl get secrets -n demo mgo-quickstart-auth -o jsonpath='{.data.\password}' | base64 -d
aaqCftpLsaGDLVIo

$ kubectl exec -it mgo-quickstart-0 -n demo sh

> mongo admin
MongoDB shell version v3.4.10
connecting to: mongodb://127.0.0.1:27017/admin
MongoDB server version: 3.4.10
Welcome to the MongoDB shell.
For interactive help, type "help".
For more comprehensive documentation, see
	http://docs.mongodb.org/
Questions? Try the support group
	http://groups.google.com/group/mongodb-user

> db.auth("root","aaqCftpLsaGDLVIo")
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

## Pause Database

KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `doNotPause` feature. If admission webhook is enabled, It prevents user from deleting the database as long as the `spec.doNotPause` is set to true. Since the MongoDB object created in this tutorial has `spec.doNotPause` set to true, if you delete the MongoDB object, KubeDB operator will nullify the delete operation. You can see this below:

```console
$ kubedb delete mg mgo-quickstart -n demo
error: MongoDB "mgo-quickstart" can't be paused. To continue delete, unset spec.doNotPause and retry.
```

Now, run `kubedb edit mg mgo-quickstart -n demo` to set `spec.doNotPause` to false or remove this field (which default to false). Then if you delete the MongoDB object, KubeDB operator will delete the StatefulSet and its pods, but leaves the PVCs unchanged. In KubeDB parlance, we say that `mgo-quickstart` MongoDB database has entered into dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase object.

```console
$ kubedb delete mg mgo-quickstart -n demo
mongodb "mgo-quickstart" deleted

$ kubedb get drmn -n demo mgo-quickstart
NAME             STATUS    AGE
mgo-quickstart   Pausing   39s

$ kubedb get drmn -n demo mgo-quickstart
NAME             STATUS    AGE
mgo-quickstart   Paused    1m
```

```yaml
$ kubedb get drmn -n demo mgo-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-02T09:24:49Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: MongoDB
  name: mgo-quickstart
  namespace: demo
  resourceVersion: "47107"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mgo-quickstart
  uid: eadf575b-07fa-11e8-946f-080027c05a6e
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: mgo-quickstart
      namespace: demo
    spec:
      mongodb:
        databaseSecret:
          secretName: mgo-quickstart-auth
        resources: {}
        storage:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 50Mi
          storageClassName: standard
        version: "3.4"
status:
  creationTime: 2018-02-02T09:24:50Z
  pausingTime: 2018-02-02T09:25:11Z
  phase: Paused
```

Here,

- `spec.origin` is the spec of the original spec of the original MongoDB object.
- `status.phase` points to the current database state `Paused`.

## Resume Dormant Database

To resume the database from the dormant state, set `spec.resume` to `true` in the DormantDatabase object.

```yaml
$ kubedb edit drmn -n demo mgo-quickstart
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  name: mgo-quickstart
  namespace: demo
  ...
spec:
  resume: true
  ...
status:
  phase: Paused
  ...
```

KubeDB operator will notice that `spec.resume` is set to true. KubeDB operator will delete the DormantDatabase object and create a new MongoDB object using the original spec. This will in turn start a new StatefulSet which will mount the originally created PVCs. Thus the original database is resumed.

Please note that the dormant database can also be resumed by creating same `MongoDB` database by using same Specs. In this tutorial, the dormant database can be resumed by creating `MongoDB` database using demo-1.yaml file. The below command resumes the dormant database `mgo-quickstart` that was created before.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/mongodb/quickstart/demo-1.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/mongodb/quickstart/demo-1.yaml"
mongodb "mgo-quickstart" created
```

## Wipeout Dormant Database

You can also wipe out a DormantDatabase by setting `spec.wipeOut` to true. KubeDB operator will delete the PVCs, delete any relevant Snapshot objects for this database and also delete snapshot data stored in the Cloud Storage buckets. There is no way to resume a wiped out database. So, be sure before you wipe out a database.

```yaml
$ kubedb edit drmn -n demo mgo-quickstart
# set spec.wipeOut: true

$ kubedb get drmn -n demo mgo-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-02T09:30:06Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: MongoDB
  name: mgo-quickstart
  namespace: demo
  resourceVersion: "47440"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mgo-quickstart
  uid: a75e6314-07fb-11e8-946f-080027c05a6e
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: mgo-quickstart
      namespace: demo
    spec:
      mongodb:
        databaseSecret:
          secretName: mgo-quickstart-auth
        resources: {}
        storage:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 50Mi
          storageClassName: standard
        version: "3.4"
  wipeOut: true
status:
  creationTime: 2018-02-02T09:30:07Z
  pausingTime: 2018-02-02T09:30:21Z
  phase: WipedOut
  wipeOutTime: 2018-02-02T09:30:51Z

$ kubedb get drmn -n demo
NAME             STATUS     AGE
mgo-quickstart   WipedOut   1m
```

## Delete Dormant Database

You still have a record that there used to be a MongoDB database `mgo-quickstart` in the form of a DormantDatabase database `mgo-quickstart`. Since you have already wiped out the database, you can delete the DormantDatabase object.

```console
$ kubedb delete drmn mgo-quickstart -n demo
dormantdatabase "mgo-quickstart" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete mg mgo-quickstart -n demo --force
$ kubedb delete drmn mgo-quickstart -n demo --force

# or
# $ kubedb delete mg,drmn,snap -n demo --all --force

$ kubectl delete ns demo
namespace "demo" deleted
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
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
