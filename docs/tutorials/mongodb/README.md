> New to KubeDB? Please start [here](/docs/tutorials/README.md).

# Running MongoDB
This tutorial will show you how to use KubeDB to run a MongoDB database.

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). 

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. This tutorial will also use a phpMyAdmin to connect and test MongoDB database, once it is running. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f ./docs/examples/mongodb/demo-0.yaml
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    1h
demo          Active    1m
kube-public   Active    1h
kube-system   Active    1h
```

## Create a MongoDB database
KubeDB implements a `MongoDB` CRD to define the specification of a MongoDB database. Below is the `MongoDB` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo1
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
  init:
    scriptSource:
      gitRepo:
        repository: "https://github.com/kubedb/mongodb-init-scripts.git"
        directory: .


$ kubedb create -f ./docs/examples/mongodb/demo-1.yaml
validating "./docs/examples/mongodb/demo-1.yaml"
mongodb "mg1" created
```

Here,
 - `spec.version` is the version of MongoDB database. In this tutorial, a MongoDB 3.4 database is going to be created.

 - `spec.doNotPause` tells KubeDB operator that if this object is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.

 - `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

 - `spec.init.scriptSource` specifies a sql script source used to initialize the database after it is created. The sql scripts will be executed alphabatically. In this tutorial, a sample sql script from the git repository `https://github.com/kubedb/mongodb-init-scripts.git` is used to create a test database.

KubeDB operator watches for `MongoDB` objects using Kubernetes api. When a `MongoDB` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching MongoDB object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. If [RBAC is enabled](/docs/tutorials/rbac.md), a ClusterRole, ServiceAccount and ClusterRoleBinding with the matching object name will be created and used as the service account name for the corresponding StatefulSet.

```console
$ kubedb describe mg -n demo mgo1
Name:		mgo1
Namespace:	demo
StartTimestamp:	Mon, 11 Dec 2017 12:38:54 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:		
  Name:			mgo1
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Mon, 11 Dec 2017 12:38:58 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:	
  Name:		mgo1
  Type:		ClusterIP
  IP:		10.105.200.75
  Port:		db	27017/TCP

Database Secret:
  Name:	mgo1-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	16 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From               Type       Reason               Message
  ---------   --------   -----     ----               --------   ------               -------
  1m          1m         1         MongoDB operator   Normal     SuccessfulValidate   Successfully validate MongoDB
  1m          1m         1         MongoDB operator   Normal     SuccessfulValidate   Successfully validate MongoDB
  1m          1m         1         MongoDB operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  1m          1m         1         MongoDB operator   Normal     SuccessfulCreate     Successfully created MongoDB
  5m          5m         1         MongoDB operator   Normal     SuccessfulValidate   Successfully validate MongoDB
  5m          5m         1         MongoDB operator   Normal     Creating             Creating Kubernetes objects


$ kubectl get statefulset -n demo
NAME      DESIRED   CURRENT   AGE
mgo1      1         1         6m


$ kubectl get pvc -n demo
NAME          STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-mgo1-0   Bound     pvc-f7746b65-de3d-11e7-879f-0800279fc284   50Mi       RWO            standard       6m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM              STORAGECLASS   REASON    AGE
pvc-f7746b65-de3d-11e7-879f-0800279fc284   50Mi       RWO            Delete           Bound     demo/data-mgo1-0   standard                 6m

$ kubectl get service -n demo
NAME      TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)     AGE
kubedb    ClusterIP   None            <none>        <none>      7m
mgo1      ClusterIP   10.105.200.75   <none>        27017/TCP   7m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified MongoDB object:

```yaml
$ kubedb get mg -n demo mgo1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-11T06:38:54Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  name: mgo1
  namespace: demo
  resourceVersion: "717"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mongodbs/mgo1
  uid: f52ab165-de3d-11e7-879f-0800279fc284
spec:
  doNotPause: true
  databaseSecret:
    secretName: mgo1-admin-auth
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
  version: 3.4
status:
  creationTime: 2017-12-11T06:38:54Z
  phase: Running
```

Please note that KubeDB operator has created a new Secret called `mgo1-admin-auth` (format: {mongodb-object-name}-admin-auth) for storing the password for `mongodb` superuser. This secret contains a `.admin` key which contains the password for `mongodb` superuser. If you want to use an existing secret please specify that when creating the MongoDB object using `spec.databaseSecret.secretName`.

Now, you can connect to this database through [mongo-shell](https://docs.mongodb.com/v3.4/mongo/). In this tutorial, we are connecting to the MongoDB server from inside of pod. 
```console
$ kubectl get secrets -n demo mgo1-admin-auth -o jsonpath='{.data.\.admin}' | base64 -d
aaqCftpLsaGDLVIo

$ kubectl exec -it mgo1-0 -n demo sh

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


## Database Snapshots

### Instant Backups
Now, you can easily take a snapshot of this database by creating a `Snapshot` object. When a `Snapshot` object is created, KubeDB operator will launch a Job that runs the `mongodump` command and uploads the output bson file to various cloud providers S3, GCS, Azure, OpenStack Swift and/or locally mounted volumes using [osm](https://github.com/appscode/osm).

In this tutorial, snapshots will be stored in a Google Cloud Storage (GCS) bucket. To do so, a secret is needed that has the following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic mg-snap-secret -n demo \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "mg-snap-secret" created
```

```yaml
$ kubectl get secret mg-snap-secret -n demo -o yaml
apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2017-12-11T08:40:34Z
  name: mg-snap-secret
  namespace: demo
  resourceVersion: "4140"
  selfLink: /api/v1/namespaces/demo/secrets/mg-snap-secret
  uid: f44826c4-de4e-11e7-879f-0800279fc284
type: Opaque
```

To lean how to configure other storage destinations for Snapshots, please visit [here](/docs/concepts/snapshot.md). Now, create the Snapshot object.

```console
$ kubedb create -f ./docs/examples/mongodb/demo-2.yaml
validating "./docs/examples/mongodb/demo-2.yaml"
snapshot "mgo-xyz" created

$ kubedb get snap -n demo
NAME      DATABASE   STATUS    AGE
mgo-xyz   mg/mgo1    Running   20s
```

```yaml
$ kubedb get snap -n demo mgo-xyz -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-11T08:41:54Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: MongoDB
    kubedb.com/name: mgo1
    snapshots.kubedb.com/status: Running
  name: mgo-xyz
  namespace: demo
  resourceVersion: "4243"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/snapshots/mgo-xyz
  uid: 2421bd0a-de4f-11e7-879f-0800279fc284
spec:
  databaseName: mgo1
  gcs:
    bucket: restic
  storageSecretName: mg-snap-secret
status:
  phase: Running
  startTime: 2017-12-11T08:41:54Z
```

Here,

- `metadata.labels` should include the type of database `kubedb.com/kind: MongoDB` whose snapshot will be taken.

- `spec.databaseName` points to the database whose snapshot is taken.

- `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.

- `spec.gcs.bucket` points to the bucket name used to store the snapshot data.


You can also run the `kubedb describe` command to see the recent snapshots taken for a database.

```console
$ kubedb describe mg -n demo mgo1
Name:		mgo1
Namespace:	demo
StartTimestamp:	Mon, 11 Dec 2017 12:38:54 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:		
  Name:			mgo1
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Mon, 11 Dec 2017 12:38:58 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:	
  Name:		mgo1
  Type:		ClusterIP
  IP:		10.105.200.75
  Port:		db	27017/TCP

Database Secret:
  Name:	mgo1-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	16 bytes

Snapshots:
  Name      Bucket      StartTime                         CompletionTime                    Phase
  ----      ------      ---------                         --------------                    -----
  mgo-xyz   gs:restic   Mon, 11 Dec 2017 14:41:54 +0600   Mon, 11 Dec 2017 14:48:58 +0600   Succeeded

Events:
  FirstSeen   LastSeen   Count     From                  Type       Reason               Message
  ---------   --------   -----     ----                  --------   ------               -------
  2m          2m         1         Snapshot Controller   Normal     SuccessfulSnapshot   Successfully completed snapshot
  9m          9m         1         Snapshot Controller   Normal     Starting             Backup running
```

Once the snapshot Job is complete, you should see the output of the `mongodump` command stored in the GCS bucket.

![snapshot-console](/docs/images/mongodb/mgo-xyz-snapshot.png)

From the above image, you can see that the snapshot output is stored in a folder called `{bucket}/kubedb/{namespace}/{mongodb-object}/{snapshot}/`.


### Scheduled Backups
KubeDB supports taking periodic backups for a database using a [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). To take periodic backups, edit the MongoDB  to add `spec.backupSchedule` section.

```yaml
$ kubedb edit mg mgo1 -n demo
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-11T06:38:54Z
  generation: 0
  initializers: null
  name: mgo1
  namespace: demo
  resourceVersion: "5378"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mongodbs/mgo1
  uid: f52ab165-de3d-11e7-879f-0800279fc284
spec:
  doNotPause: true
  databaseSecret:
    secretName: mgo1-admin-auth
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
  version: 3.4
  backupSchedule:
    cronExpression: '@every 1m'
    gcs:
      bucket: restic
    storageSecretName: mg-snap-secret
status:
  creationTime: 2017-12-11T06:38:54Z
  phase: Running
```

Once the `spec.backupSchedule` is added, KubeDB operator will create a new Snapshot object on each tick of the cron expression. This triggers KubeDB operator to create a Job as it would for any regular instant backup process. You can see the snapshots as they are created using `kubedb get snap` command.
```console
$ kubedb get snap -n demo
NAME                   DATABASE   STATUS      AGE
mgo-xyz                mg/mgo1    Succeeded   21m
mgo1-20171211-090039   mg/mgo1    Succeeded   2m
mgo1-20171211-090159   mg/mgo1    Succeeded   1m
mgo1-20171211-090259   mg/mgo1    Running     3s
```

### Restore from Snapshot
You can create a new database from a previously taken Snapshot. Specify the Snapshot name in the `spec.init.snapshotSource` field of a new MongoDB object. See the example `recovered` object below:

```yaml
$ cat ./docs/examples/mongodb/demo-4.yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: recovered
  namespace: demo
spec:
  version: 3.4
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    snapshotSource:
      name: mgo-xyz


$ kubedb create -f ./docs/examples/mongodb/demo-4.yaml
validating "./docs/examples/mongodb/demo-4.yaml"
mongodb "recovered" created
```

Here,
 - `spec.init.snapshotSource.name` refers to a Snapshot object for a MongoDB database in the same namespaces as this new `recovered` MongoDB object.

Now, wait several seconds. KubeDB operator will create a new StatefulSet. Then KubeDB operator launches a Kubernetes Job to initialize the new database using the data from `mgo-xyz` Snapshot.

```console
$ kubedb get mg -n demo
NAME        STATUS    AGE
mgo1        Running   2h
recovered   Running   1m


$ kubedb describe mg -n demo recovered
Name:		recovered
Namespace:	demo
StartTimestamp:	Mon, 11 Dec 2017 15:04:34 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:		
  Name:			recovered
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Mon, 11 Dec 2017 15:04:38 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:	
  Name:		recovered
  Type:		ClusterIP
  IP:		10.103.115.39
  Port:		db	27017/TCP

Database Secret:
  Name:	recovered-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	16 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From               Type       Reason                 Message
  ---------   --------   -----     ----               --------   ------                 -------
  1m          1m         1         MongoDB operator   Normal     SuccessfulValidate     Successfully validate MongoDB
  1m          1m         1         MongoDB operator   Normal     SuccessfulValidate     Successfully validate MongoDB
  1m          1m         1         MongoDB operator   Normal     SuccessfulInitialize   Successfully completed initialization
  1m          1m         1         MongoDB operator   Normal     SuccessfulCreate       Successfully created MongoDB
  1m          1m         1         MongoDB operator   Normal     SuccessfulCreate       Successfully created StatefulSet
  1m          1m         1         MongoDB operator   Normal     Initializing           Initializing from Snapshot: "mgo-xyz"
  2m          2m         1         MongoDB operator   Normal     SuccessfulValidate     Successfully validate MongoDB
  2m          2m         1         MongoDB operator   Normal     Creating               Creating Kubernetes objects
```

## Pause Database

Since the MongoDB object created in this tutorial has `spec.doNotPause` set to true, if you delete the MongoDB object, KubeDB operator will recreate the object and essentially nullify the delete operation. You can see this below:

```console
$ kubedb delete mg mgo1 -n demo
error: MongoDB "mgo1" can't be paused. To continue delete, unset spec.doNotPause and retry.
```

Now, run `kubedb edit mg mgo1 -n demo` to set `spec.doNotPause` to false or remove this field (which default to false). Then if you delete the MongoDB object, KubeDB operator will delete the StatefulSet and its pods, but leaves the PVCs unchanged. In KubeDB parlance, we say that `mgo1` MongoDB database has entered into dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase object.

```yaml
$ kubedb delete mg mgo1 -n demo
mongodb "mgo1" deleted

$ kubedb get drmn -n demo mgo1
NAME      STATUS    AGE
mgo1      Pausing   39s

$ kubedb get drmn -n demo mgo1
NAME      STATUS    AGE
mgo1      Paused    1m


$ kubedb get drmn -n demo mgo1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  annotations:
    mongodbs.kubedb.com/init: '{"scriptSource":{"gitRepo":{"repository":"https://github.com/kubedb/mongodb-init-scripts.git","directory":"."}}}'
  clusterName: ""
  creationTimestamp: 2017-12-11T09:17:01Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: MongoDB
  name: mgo1
  namespace: demo
  resourceVersion: "7029"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mgo1
  uid: 0ba4a0df-de54-11e7-879f-0800279fc284
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: mgo1
      namespace: demo
    spec:
      mongodb:
        databaseSecret:
          secretName: mgo1-admin-auth
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
  creationTime: 2017-12-11T09:17:01Z
  pausingTime: 2017-12-11T09:18:11Z
  phase: Paused
```

Here,
 - `spec.origin` is the spec of the original spec of the original MongoDB object.

 - `status.phase` points to the current database state `Paused`.


## Resume Dormant Database

To resume the database from the dormant state, set `spec.resume` to `true` in the DormantDatabase object.

```yaml
$ kubedb edit drmn -n demo mgo1
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  annotations:
    mongodbs.kubedb.com/init: '{"scriptSource":{"gitRepo":{"repository":"https://github.com/kubedb/mongodb-init-scripts.git","directory":"."}}}'
  clusterName: ""
  creationTimestamp: 2017-12-11T09:17:01Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: MongoDB
  name: mgo1
  namespace: demo
  resourceVersion: "7029"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mgo1
  uid: 0ba4a0df-de54-11e7-879f-0800279fc284
spec:
  resume: true
  origin:
    metadata:
      creationTimestamp: null
      name: mgo1
      namespace: demo
    spec:
      mongodb:
        databaseSecret:
          secretName: mgo1-admin-auth
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
  creationTime: 2017-12-11T09:17:01Z
  pausingTime: 2017-12-11T09:18:11Z
  phase: Paused

```

KubeDB operator will notice that `spec.resume` is set to true. KubeDB operator will delete the DormantDatabase object and create a new MongoDB object using the original spec. This will in turn start a new StatefulSet which will mount the originally created PVCs. Thus the original database is resumed.

## Wipeout Dormant Database
You can also wipe out a DormantDatabase by setting `spec.wipeOut` to true. KubeDB operator will delete the PVCs, delete any relevant Snapshot objects for this database and also delete snapshot data stored in the Cloud Storage buckets. There is no way to resume a wiped out database. So, be sure before you wipe out a database.

```yaml
$ kubedb edit drmn -n demo mgo1
# set spec.wipeOut: true

$ kubedb get drmn -n demo mgo1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-11T09:22:11Z
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: MongoDB
  name: mgo1
  namespace: demo
  resourceVersion: "7497"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mgo1
  uid: c4b9cb9c-de54-11e7-879f-0800279fc284
spec:
  origin:
    metadata:
      annotations:
        mongodbs.kubedb.com/init: '{"scriptSource":{"gitRepo":{"repository":"https://github.com/kubedb/mongodb-init-scripts.git","directory":"."}}}'
      creationTimestamp: null
      name: mgo1
      namespace: demo
    spec:
      mongodb:
        databaseSecret:
          secretName: mgo1-admin-auth
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
  creationTime: 2017-12-11T09:22:11Z
  pausingTime: 2017-12-11T09:23:11Z
  phase: WipedOut
  wipeOutTime: 2017-12-11T09:24:09Z
  

$ kubedb get drmn -n demo
NAME      STATUS     AGE
mgo1      WipedOut   3m
```


## Delete Dormant Database
You still have a record that there used to be a MongoDB database `mgo1` in the form of a DormantDatabase database `mgo1`. Since you have already wiped out the database, you can delete the DormantDatabase object. 

```console
$ kubedb delete drmn mgo1 -n demo
dormantdatabase "mgo1" deleted
```

## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
namespace "demo" deleted
```

If you would like to uninstall KubeDB operator, please follow the steps [here](/docs/uninstall.md).


## Next Steps
- Learn about the details of MongoDB object [here](/docs/concepts/mongodb.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Thinking about monitoring your database? KubeDB works [out-of-the-box with Prometheus](/docs/tutorials/monitoring.md).
- Learn how to use KubeDB in a [RBAC](/docs/tutorials/rbac.md) enabled cluster.
- Wondering what features are coming next? Please visit [here](/ROADMAP.md). 
- Want to hack on KubeDB? Check our [contribution guidelines](/CONTRIBUTING.md).
