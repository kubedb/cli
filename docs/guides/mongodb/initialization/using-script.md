---
title: Initialize MongoDB using Script
menu:
  docs_0.8.0:
    identifier: mg-using-script-initialization
    name: Using Script
    parent: mg-initialization-mongodb
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Initialize MongoDB using Script

This tutorial will show you how to use KubeDB to initialize a MongoDB database with .js and/or .sh script.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    45m
demo          Active    10s
kube-public   Active    45m
kube-system   Active    45m
```

In this tutorial we will use .js script stored in GitHub repository [kubedb/mongodb-init-scripts](https://github.com/kubedb/mongodb-init-scripts).

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create a MongoDB database with Init-Script

Below is the `MongoDB` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-init-script
  namespace: demo
spec:
  version: "3.4"
  doNotPause: true
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    scriptSource:
      gitRepo:
        repository: "https://github.com/kubedb/mongodb-init-scripts.git"
        directory: .

```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/mongodb/Initialization/demo-1.yaml
mongodb "mgo-init-script" created
```

Here,

- `spec.version` is the version of MongoDB database. In this tutorial, a MongoDB 3.4 database is going to be created.
- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. Since release 0.8.0, a storage spec is required for MongoDB.
- `spec.init.scriptSource` specifies a script source used to initialize the database before database server starts. The scripts will be executed alphabatically. In this tutorial, a sample .js script from the git repository `https://github.com/kubedb/mongodb-init-scripts.git` is used to create a test database. You can use other [volume sources](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes) instead of `gitrepo`.  The \*.js and/or \*.sh sripts that are stored inside the root folder will be executed alphabatically. The scripts inside child folders will be skipped.

KubeDB operator watches for `MongoDB` objects using Kubernetes api. When a `MongoDB` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching MongoDB object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. No MongoDB specific RBAC roles are required for [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb describe mg -n demo mgo-init-script
Name:		mgo-init-script
Namespace:	demo
StartTimestamp:	Tue, 06 Feb 2018 09:56:07 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			mgo-init-script
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Tue, 06 Feb 2018 09:56:12 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		mgo-init-script
  Type:		ClusterIP
  IP:		10.106.175.209
  Port:		db	27017/TCP

Database Secret:
  Name:	mgo-init-script-auth
  Type:	Opaque
  Data
  ====
  password:	16 bytes
  user:		4 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From               Type       Reason       Message
  ---------   --------   -----     ----               --------   ------       -------
  6s          6s         1         MongoDB operator   Normal     Successful   Successfully patched StatefulSet
  6s          6s         1         MongoDB operator   Normal     Successful   Successfully patched MongoDB
  9s          9s         1         MongoDB operator   Normal     Successful   Successfully created StatefulSet
  9s          9s         1         MongoDB operator   Normal     Successful   Successfully created MongoDB
  18s         18s        1         MongoDB operator   Normal     Successful   Successfully created Service



$ kubectl get statefulset -n demo
NAME              DESIRED   CURRENT   AGE
mgo-init-script   1         1         46s


$ kubectl get pvc -n demo
NAME                     STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-mgo-init-script-0   Bound     pvc-ac84fbb9-0af1-11e8-a107-080027869227   50Mi       RWO            standard       1m


$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                         STORAGECLASS   REASON    AGE
pvc-ac84fbb9-0af1-11e8-a107-080027869227   50Mi       RWO            Delete           Bound     demo/data-mgo-init-script-0   standard


$ kubectl get service -n demo
NAME              TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
kubedb            ClusterIP   None             <none>        <none>      2m
mgo-init-script   ClusterIP   10.106.175.209   <none>        27017/TCP   2m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified MongoDB object:

```yaml
$ kubedb get mg -n demo mgo-init-script -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-06T03:56:07Z
  finalizers:
  - kubedb.com
  generation: 0
  name: mgo-init-script
  namespace: demo
  resourceVersion: "4827"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mongodbs/mgo-init-script
  uid: a9348cad-0af1-11e8-a107-080027869227
spec:
  databaseSecret:
    secretName: mgo-init-script-auth
  doNotPause: true
  init:
    scriptSource:
      gitRepo:
        directory: .
        repository: https://github.com/kubedb/mongodb-init-scripts.git
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: "3.4"
status:
  creationTime: 2018-02-06T03:56:12Z
  phase: Running
```

Please note that KubeDB operator has created a new Secret called `mgo-init-script-auth` *(format: {mongodb-object-name}-auth)* for storing the password for MongoDB superuser. This secret contains a `user` key which contains the *username* for MongoDB superuser and a `password` key which contains the *password* for MongoDB superuser.
If you want to use an existing secret please specify that when creating the MongoDB object using `spec.databaseSecret.secretName`. While creating this secret manually, make sure the secret contains these two keys containing data `user` and `password`.

```console
$ kubectl get secrets -n demo mgo-init-script-auth -o json
apiVersion: v1
data:
  password: STJ1YnNiU3BUNzFOZUhXSA==
  user: cm9vdA==
kind: Secret
metadata:
  creationTimestamp: 2018-02-06T03:56:12Z
  labels:
    kubedb.com/kind: MongoDB
    kubedb.com/name: mgo-init-script
  name: mgo-init-script-auth
  namespace: demo
  resourceVersion: "4789"
  selfLink: /api/v1/namespaces/demo/secrets/mgo-init-script-auth
  uid: ac33c72d-0af1-11e8-a107-080027869227
type: Opaque
```

Now, you can connect to this database through [mongo-shell](https://docs.mongodb.com/v3.4/mongo/). In this tutorial, we are connecting to the MongoDB server from inside the pod.

```console
$ kubectl get secrets -n demo mgo-init-script-auth -o jsonpath='{.data.\user}' | base64 -d
root

$ kubectl get secrets -n demo mgo-init-script-auth -o jsonpath='{.data.\password}' | base64 -d
I2ubsbSpT71NeHWH

$ kubectl exec -it mgo-init-script-0 -n demo sh

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

> db.auth("root","I2ubsbSpT71NeHWH")
1

> show dbs
admin  0.000GB
local  0.000GB
mydb   0.000GB

> use mydb
switched to db mydb

> db.movie.find()
{ "_id" : ObjectId("5a72b2b1e1a0770e3bdb56f1"), "name" : "batman" }

> exit
bye
```

As you can see here, the initial script has successfully created a database named `mydb` and inserted data into that database successfully.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo mg/mgo-init-script -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo mg/mgo-init-script

$ kubectl patch -n demo drmn/mgo-init-script -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/mgo-init-script

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Initialize [MongoDB with Snapshot](/docs/guides/mongodb/initialization/using-snapshot.md).
- [Snapshot and Restore](/docs/guides/mongodb/snapshot/backup-and-restore.md) process of MongoDB databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mongodb/snapshot/scheduled-backup.md) of MongoDB databases using KubeDB.
- Initialize [MongoDB with Script](/docs/guides/mongodb/initialization/using-script.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mongodb/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mongodb/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mongodb/private-registry/using-private-registry.md) to deploy MongoDB with KubeDB.
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
