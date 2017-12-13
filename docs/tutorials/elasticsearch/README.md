---
title: Elasticsearch
menu:
  docs_0.7.1:
    identifier: tutorials-elasticsearch-readme
    name: Overview
    parent: tutorials-elasticsearch
    weight: 10
menu_name: docs_0.7.1
section_menu_id: tutorials
url: /docs/0.7.1/tutorials/elasticsearch/
aliases:
  - /docs/0.7.1/tutorials/elasticsearch/README/
---

> New to KubeDB? Please start [here](/docs/tutorials/README.md).

# Running Elasticsearch
This tutorial will show you how to use KubeDB to run an Elasticsearch database.

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f ./docs/examples/elasticsearch/demo-0.yaml
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    3m
demo          Active    5s
kube-public   Active    3m
kube-system   Active    3m
```

## Create an Elasticsearch database
KubeDB implements a `Elasticsearch` CRD to define the specification of an Elasticsearch database. Below is the `Elasticsearch` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: e1
  namespace: demo
spec:
  version: 2.3.1
  replicas: 1
  doNotPause: true
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi

$ kubedb create -f ./docs/examples/elasticsearch/demo-1.yaml
validating "./docs/examples/elasticsearch/demo-1.yaml"
elasticsearch "e1" created
```

Here,
 - `spec.version` is the version of Elasticsearch database. In this tutorial, an Elasticsearch 2.3.1 cluster is going to be created.

 - `spec.replicas` is the number of pods in the Elasticsearch cluster. In this tutorial, a single node Elasticsearch cluster is going to be created.

 - `spec.doNotPause` tells KubeDB operator that if this object is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.

 - `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

KubeDB operator watches for `Elasticsearch` objects using Kubernetes api. When a `Elasticsearch` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching Elasticsearch object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. If [RBAC is enabled](/docs/tutorials/rbac.md), a ClusterRole, ServiceAccount and ClusterRoleBinding with the matching object name will be created and used as the service account name for the corresponding StatefulSet.

```console
$ kubedb describe es e1 -n demo
Name:			e1
Namespace:		demo
CreationTimestamp:	Tue, 18 Jul 2017 14:35:41 -0700
Status:			Running
Replicas:		1  total
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

Service:	
  Name:		e1
  Type:		ClusterIP
  IP:		10.0.0.238
  Port:		db	9200/TCP
  Port:		cluster	9300/TCP

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason               Message
  ---------   --------   -----     ----                     --------   ------               -------
  6m          6m         1         Elasticsearch operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  6m          6m         1         Elasticsearch operator   Normal     SuccessfulCreate     Successfully created Elasticsearch
  8m          8m         1         Elasticsearch operator   Normal     SuccessfulValidate   Successfully validate Elasticsearch
  8m          8m         1         Elasticsearch operator   Normal     Creating             Creating Kubernetes objects


$ kubectl get statefulset -n demo
NAME      DESIRED   CURRENT   AGE
e1        1         1         8m

$ kubectl get pvc -n demo
NAME        STATUS    VOLUME                                     CAPACITY   ACCESSMODES   STORAGECLASS   AGE
data-e1-0   Bound     pvc-0d32d0e8-6c01-11e7-b566-080027691dbf   50Mi RWO           standard       8m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS    CLAIM            STORAGECLASS   REASON    AGE
pvc-0d32d0e8-6c01-11e7-b566-080027691dbf   50Mi RWO           Delete          Bound     demo/data-e1-0   standard                 8m

$ kubectl get service -n demo
NAME      CLUSTER-IP   EXTERNAL-IP   PORT(S)             AGE
e1        10.0.0.238   <none>        9200/TCP,9300/TCP   9m
kubedb    None         <none>                            9m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified Elasticsearch object:

```yaml
$ kubedb get es -n demo e1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  creationTimestamp: 2017-07-18T21:35:41Z
  name: e1
  namespace: demo
  resourceVersion: "608"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/elasticsearchs/e1
  uid: 0c174082-6c01-11e7-b566-080027691dbf
spec:
  doNotPause: true
  replicas: 1
  resources: {}
  storage:
    accessModes:
    - ReadWriteOnce
    storageClassName: standard
    resources:
      requests:
        storage: 50Mi
  version: 2.3.1
status:
  creationTime: 2017-07-18T21:35:41Z
  phase: Running
```


Please note that KubeDB operator has created a new Secret called `e1-admin-auth` (format: {elasticsearch-object-name}-admin-auth) for storing the password for `postgres` superuser. This secret contains a `.admin` key with a ini formatted key-value pairs. If you want to use an existing secret please specify that when creating the Elasticsearch object using `spec.databaseSecret.secretName`.

Now, you can connect to this Elasticsearch cluster from inside the cluster.

```console
$ kubectl get pods e1-0 -n demo -o yaml | grep IP
  hostIP: 192.168.99.100
  podIP: 172.17.0.5

# Exec into kubedb operator pod
$ kubectl exec -it $(kubectl get pods --all-namespaces -l app=kubedb -o jsonpath='{.items[0].metadata.name}') -n kube-system sh

~ $ ps aux
PID   USER     TIME   COMMAND
    1 nobody     0:00 /operator run --address=:8080 --rbac=false --v=3
   18 nobody     0:00 sh
   26 nobody     0:00 ps aux
~ $ wget -qO- http://172.17.0.5:9200
{
  "name" : "e1-0.demo",
  "cluster_name" : "e1",
  "version" : {
    "number" : "2.3.1",
    "build_hash" : "bd980929010aef404e7cb0843e61d0665269fc39",
    "build_timestamp" : "2016-04-04T12:25:05Z",
    "build_snapshot" : false,
    "lucene_version" : "5.5.0"
  },
  "tagline" : "You Know, for Search"
}
```

![Using e1 from esAdmin4](/docs/images/elasticsearch/e1.gif)


## Database Snapshots

### Instant Backups
Now, you can easily take a snapshot of this database by creating a `Snapshot` object. When a `Snapshot` object is created, KubeDB operator will launch a Job that runs [elasticdump](https://github.com/taskrabbit/elasticsearch-dump) command and uploads snapshot data to various cloud providers S3, GCS, Azure, OpenStack Swift and/or locally mounted volumes using [osm](https://github.com/appscode/osm).

In this tutorial, snapshots will be stored in a Google Cloud Storage (GCS) bucket. To do so, a secret is needed that has the following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic es-snap-secret -n demo \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "es-snap-secret" created
```

```yaml
$ kubectl get secret es-snap-secret -o yaml

apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2017-07-17T18:06:51Z
  name: es-snap-secret
  namespace: demo
  resourceVersion: "5461"
  selfLink: /api/v1/namespaces/demo/secrets/es-snap-secret
  uid: a6983b00-5c02-11e7-bb52-08002711f4aa
type: Opaque
```


To lean how to configure other storage destinations for Snapshots, please visit [here](/docs/concepts/snapshot.md). Now, create the Snapshot object.

```
$ kubedb create -f ./docs/examples/elasticsearch/demo-2.yaml
validating "./docs/examples/elasticsearch/demo-2.yaml"
snapshot "e1-xyz" created

$ kubedb get snap -n demo
NAME      DATABASE   STATUS    AGE
e1-xyz    es/e1      Running   22s
```

```yaml
$ kubedb get snap -n demo e1-xyz -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  creationTimestamp: 2017-07-18T22:21:40Z
  labels:
    kubedb.com/kind: Elasticsearch
    kubedb.com/name: e1
  name: e1-xyz
  namespace: demo
  resourceVersion: "3713"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/snapshots/e1-xyz
  uid: 78d99dfe-6c07-11e7-b566-080027691dbf
spec:
  databaseName: e1
  gcs:
    bucket: restic
  resources: {}
  storageSecretName: snap-secret
status:
  completionTime: 2017-07-18T22:23:53Z
  phase: Succeeded
  startTime: 2017-07-18T22:21:40Z
```

Here,

- `metadata.labels` should include the type of database `kubedb.com/kind: Elasticsearch` whose snapshot will be taken.

- `spec.databaseName` points to the database whose snapshot is taken.

- `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.

- `spec.gcs.bucket` points to the bucket name used to store the snapshot data.


You can also run the `kubedb describe` command to see the recent snapshots taken for a database.

```console
$ kubedb describe es -n demo e1
Name:			e1
Namespace:		demo
CreationTimestamp:	Tue, 18 Jul 2017 14:35:41 -0700
Status:			Running
Replicas:		1  total
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

Service:	
  Name:		e1
  Type:		ClusterIP
  IP:		10.0.0.238
  Port:		db	9200/TCP
  Port:		cluster	9300/TCP

Snapshots:
  Name     Bucket      StartTime                         CompletionTime                    Phase
  ----     ------      ---------                         --------------                    -----
  e1-xyz   gs:restic   Tue, 18 Jul 2017 15:21:40 -0700   Tue, 18 Jul 2017 15:23:53 -0700   Succeeded

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason               Message
  ---------   --------   -----     ----                     --------   ------               -------
  4m          4m         1         Snapshot Controller      Normal     SuccessfulSnapshot   Successfully completed snapshot
  6m          6m         1         Snapshot Controller      Normal     Starting             Backup running
  50m         50m        1         Elasticsearch operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  50m         50m        1         Elasticsearch operator   Normal     SuccessfulCreate     Successfully created Elasticsearch
  52m         52m        1         Elasticsearch operator   Normal     SuccessfulValidate   Successfully validate Elasticsearch
  52m         52m        1         Elasticsearch operator   Normal     Creating             Creating Kubernetes objects
```

Once the snapshot Job is complete, you should see the output of the [elasticdump](https://github.com/taskrabbit/elasticsearch-dump) process stored in the GCS bucket.

![snapshot-console](/docs/images/elasticsearch/e1-xyz-snapshot.png)

From the above image, you can see that the snapshot output is stored in a folder called `{bucket}/kubedb/{namespace}/{elasticsearch-object}/{snapshot}/`.


### Scheduled Backups
KubeDB supports taking periodic backups for a database using a [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). To take periodic backups, edit the Elasticsearch object to add `spec.backupSchedule` section.

```yaml
$ kubedb edit es e1 -n demo

apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: e1
  namespace: demo
spec:
  version: 2.3.1
  replicas: 1
  doNotPause: true
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  backupSchedule:
    cronExpression: "@every 1m"
    storageSecretName: snap-secret
    gcs:
      bucket: restic
```

Once the `spec.backupSchedule` is added, KubeDB operator will create a new Snapshot object on each tick of the cron expression. This triggers KubeDB operator to create a Job as it would for any regular instant backup process. You can see the snapshots as they are created using `kubedb get snap` command.
```console
$ kubedb get snap -n demo
NAME                 DATABASE   STATUS      AGE
e1-20170718-223046   es/e1      Succeeded   8m
e1-20170718-223206   es/e1      Running     7m
e1-xyz               es/e1      Succeeded   18m
```

### Restore from Snapshot
You can create a new database from a previously taken Snapshot. Specify the Snapshot name in the `spec.init.snapshotSource` field of a new Elasticsearch object. See the example `recovered` object below:

```yaml
$ cat ./docs/examples/elasticsearch/demo-4.yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: recovered
  namespace: demo
spec:
  version: 2.3.1
  doNotPause: true
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    snapshotSource:
      name: e1-xyz

$ kubedb create -f ./docs/examples/elasticsearch/demo-4.yaml
validating "./docs/examples/elasticsearch/demo-4.yaml"
elasticsearch "recovered" created
```

Here,
 - `spec.init.snapshotSource.name` refers to a Snapshot object for a Elasticsearch database in the same namespaces as this new `recovered` Elasticsearch object.

Now, wait several seconds. KubeDB operator will create a new StatefulSet. Then KubeDB operator launches a Kubernetes Job to initialize the new database using the data from `e1-xyz` Snapshot.

```console
$ kubedb get es -n demo
NAME        STATUS    AGE
e1          Running   1h
recovered   Running   49s


$ kubedb describe es -n demo recovered
Name:			recovered
Namespace:		demo
CreationTimestamp:	Tue, 18 Jul 2017 15:41:45 -0700
Status:			Running
Replicas:		0  total
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

Service:	
  Name:		recovered
  Type:		ClusterIP
  IP:		10.0.0.65
  Port:		db	9200/TCP
  Port:		cluster	9300/TCP

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason                 Message
  ---------   --------   -----     ----                     --------   ------                 -------
  1m          1m         1         Elasticsearch operator   Normal     SuccessfulInitialize   Successfully completed initialization
  1m          1m         1         Elasticsearch operator   Normal     SuccessfulCreate       Successfully created Elasticsearch
  1m          1m         1         Elasticsearch operator   Normal     SuccessfulValidate     Successfully validate Elasticsearch
  1m          1m         1         Elasticsearch operator   Normal     Creating               Creating Kubernetes objects
  1m          1m         1         Elasticsearch operator   Normal     Initializing           Initializing from Snapshot: "e1-xyz"
```


## Pause Database

Since the Elasticsearch object created in this tutorial has `spec.doNotPause` set to true, if you delete the Elasticsearch object, KubeDB operator will recreate the object and essentially nullify the delete operation. You can see this below:

```console
$ kubedb delete es e1 -n demo
error: Elasticsearch "e1" can't be paused. To continue delete, unset spec.doNotPause and retry.
```

Now, run `kubedb edit es e1 -n demo` to set `spec.doNotPause` to false or remove this field (which default to false). Then if you delete the Elasticsearch object, KubeDB operator will delete the StatefulSet and its pods, but leaves the PVCs unchanged. In KubeDB parlance, we say that `e1` Elasticsearch database has entered into dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase object.

```yaml
$ kubedb delete es -n demo e1
elasticsearch "e1" deleted

$ kubedb get drmn -n demo e1
NAME      STATUS    AGE
e1        Pausing   20s

$ kubedb get drmn -n demo e1
NAME      STATUS    AGE
e1        Paused    3m

$ kubedb get drmn -n demo e1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: 2017-07-18T22:47:51Z
  labels:
    kubedb.com/kind: Elasticsearch
  name: e1
  namespace: demo
  resourceVersion: "6216"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/e1
  uid: 21464b6c-6c0b-11e7-b566-080027691dbf
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: e1
      namespace: demo
    spec:
      elasticsearch:
        backupSchedule:
          cronExpression: '@every 1m'
          gcs:
            bucket: restic
          resources: {}
          storageSecretName: snap-secret
        replicas: 1
        resources: {}
        storage:
          accessModes:
          - ReadWriteOnce
          storageClassName: standard
          resources:
            requests:
              storage: 50Mi
        version: 2.3.1
status:
  creationTime: 2017-07-18T22:47:51Z
  pausingTime: 2017-07-18T22:48:01Z
  phase: Paused
```

Here,
 - `spec.origin` is the spec of the original spec of the original Elasticsearch object.

 - `status.phase` points to the current database state `Paused`.


## Resume Dormant Database

To resume the database from the dormant state, set `spec.resume` to `true` in the DormantDatabase object.

```yaml
$ kubedb edit drmn -n demo e1

apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: 2017-07-18T22:47:51Z
  labels:
    kubedb.com/kind: Elasticsearch
  name: e1
  namespace: demo
  resourceVersion: "6216"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/e1
  uid: 21464b6c-6c0b-11e7-b566-080027691dbf
spec:
  resume: true
  origin:
    metadata:
      creationTimestamp: null
      name: e1
      namespace: demo
    spec:
      elasticsearch:
        backupSchedule:
          cronExpression: '@every 1m'
          gcs:
            bucket: restic
          resources: {}
          storageSecretName: snap-secret
        replicas: 1
        resources: {}
        storage:
          accessModes:
          - ReadWriteOnce
          storageClassName: standard
          resources:
            requests:
              storage: 50Mi
        version: 2.3.1
status:
  creationTime: 2017-07-18T22:47:51Z
  pausingTime: 2017-07-18T22:48:01Z
  phase: Paused
```

KubeDB operator will notice that `spec.resume` is set to true. KubeDB operator will delete the DormantDatabase object and create a new Elasticsearch object using the original spec. This will in turn start a new StatefulSet which will mount the originally created PVCs. Thus the original database is resumed.

## Wipeout Dormant Database
You can also wipe out a DormantDatabase by setting `spec.wipeOut` to true. KubeDB operator will delete the PVCs, delete any relevant Snapshot objects for this database and also delete snapshot data stored in the Cloud Storage buckets. There is no way to resume a wiped out database. So, be sure before you wipe out a database.

```yaml
$ kubedb edit drmn -n demo e1
# set spec.wipeOut: true

$ kubedb get drmn -n demo e1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: 2017-07-18T22:51:42Z
  labels:
    kubedb.com/kind: Elasticsearch
  name: e1
  namespace: demo
  resourceVersion: "6653"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/e1
  uid: aacfbbec-6c0b-11e7-b566-080027691dbf
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: e1
      namespace: demo
    spec:
      elasticsearch:
        backupSchedule:
          cronExpression: '@every 1m'
          gcs:
            bucket: restic
          resources: {}
          storageSecretName: snap-secret
        replicas: 1
        resources: {}
        storage:
          accessModes:
          - ReadWriteOnce
          storageClassName: standard
          resources:
            requests:
              storage: 50Mi
        version: 2.3.1
  wipeOut: true
status:
  creationTime: 2017-07-18T22:51:42Z
  pausingTime: 2017-07-18T22:51:52Z
  phase: WipedOut
  wipeOutTime: 2017-07-18T22:52:37Z

$ kubedb get drmn -n demo
NAME      STATUS     AGE
e1        WipedOut   1m
```


## Delete Dormant Database
You still have a record that there used to be an Elasticsearch database `e1` in the form of a DormantDatabase database `e1`. Since you have already wiped out the database, you can delete the DormantDatabase object.

```console
$ kubedb delete drmn e1 -n demo
dormantdatabase "e1" deleted
```

## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall KubeDB operator, please follow the steps [here](/docs/uninstall.md).


## Next Steps
- Learn about the details of Elasticsearch object [here](/docs/concepts/elasticsearch.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Thinking about monitoring your database? KubeDB works [out-of-the-box with Prometheus](/docs/tutorials/monitoring.md).
- Learn how to use KubeDB in a [RBAC](/docs/tutorials/rbac.md) enabled cluster.
- Wondering what features are coming next? Please visit [here](/ROADMAP.md). 
- Want to hack on KubeDB? Check our [contribution guidelines](/CONTRIBUTING.md).
