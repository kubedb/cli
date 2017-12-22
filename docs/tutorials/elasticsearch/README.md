---
title: Elasticsearch
menu:
  docs_0.8.0:
    identifier: tutorials-elasticsearch-readme
    name: Overview
    parent: tutorials-elasticsearch
    weight: 10
menu_name: docs_0.8.0
section_menu_id: tutorials
url: /docs/0.8.0/tutorials/elasticsearch/
aliases:
  - /docs/0.8.0/tutorials/elasticsearch/README/
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
  version: 5.6.4
  replicas: 1
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
$ kubedb create -f ./docs/examples/elasticsearch/demo-1.yaml
validating "./docs/examples/elasticsearch/demo-1.yaml"
elasticsearch "e1" created
```

Here,
 - `spec.version` is the version of Elasticsearch database. In this tutorial, an Elasticsearch 5.6.4 database is going to be created.
 - `spec.replicas` is the number of nodes in the Elasticsearch cluster. Here, we are creating a single node Elasticsearch cluster.
 - `spec.doNotPause` tells KubeDB operator that if this CRD object is deleted, it should be automatically reverted. This should be set to `true` for production databases to avoid accidental deletion.
 - `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

KubeDB operator watches for `Elasticsearch` objects using Kubernetes api. When a `Elasticsearch` object is created, KubeDB operator will create a new StatefulSet and two ClusterIP Service with the matching name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. If [RBAC is enabled](/docs/tutorials/rbac.md), a ClusterRole, ServiceAccount and ClusterRoleBinding with the matching tpr name will be created and used as the service account name for the corresponding StatefulSet.


```console
$ kubedb describe es e1 -n demo
Name:               e1
Namespace:          demo
CreationTimestamp:  Thu, 14 Dec 2017 10:04:23 +0600
Status:             Running
Replicas:           1 total
Volume:
  StorageClass: standard
  Capacity:     50Mi
  Access Modes: RWO

StatefulSet:
  Name:                 e1
  Replicas:             1 current / 1 desired
  CreationTimestamp:    Thu, 14 Dec 2017 10:04:29 +0600
  Pods Status:          1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		e1
  Type:		ClusterIP
  IP:		10.108.4.122
  Port:		http	9200/TCP

Service:
  Name:		e1-master
  Type:		ClusterIP
  IP:		10.103.97.44
  Port:		transport	9300/TCP

Database Secret:
  Name:	e1-auth
  Type:	Opaque
  Data
  ====
  sg_roles.yml:             312 bytes
  sg_roles_mapping.yml:     73 bytes
  ADMIN_PASSWORD:           8 bytes
  READALL_PASSWORD:         8 bytes
  sg_action_groups.yml:     430 bytes
  sg_config.yml:            240 bytes
  sg_internal_users.yml:    156 bytes

Certificate Secret:
  Name:	e1-cert
  Type:	Opaque
  Data
  ====
  ca.pem:           1139 bytes
  client-key.pem:   1675 bytes
  client.pem:       1151 bytes
  keystore.jks:     3050 bytes
  sgadmin.jks:      3011 bytes
  truststore.jks:   864 bytes

Topology:
  Type                 Pod       StartTime                       Phase
  ----                 ---       ---------                       -----
  master|client|data   e1-0      2017-12-14 10:04:30 +0600 +06   Running

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason               Message
  ---------   --------   -----     ----                     --------   ------               -------
  16s         16s        1         Elasticsearch operator   Normal     SuccessfulCreate     Successfully created Elasticsearch
  46s         46s        1         Elasticsearch operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  59s         59s        1         Elasticsearch operator   Normal     SuccessfulValidate   Successfully validate Elasticsearch
  59s         59s        1         Elasticsearch operator   Normal     Creating             Creating Kubernetes objects
```

```console
$ kubectl get pvc -n demo
NAME        STATUS    VOLUME                                     CAPACITY   ACCESSMODES   STORAGECLASS   AGE
data-e1-0   Bound     pvc-35683016-dfec-11e7-9e33-08002726ce5b   50Mi       RWO           standard       12m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS     CLAIM            STORAGECLASS   REASON    AGE
pvc-35683016-dfec-11e7-9e33-08002726ce5b   50Mi       RWO           Delete          Bound      demo/data-e1-0   standard                 12m

$ kubectl get service -n demo
NAME        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
e1          10.99.174.203    <none>        9200/TCP   13m
e1-master   10.103.121.146   <none>        9300/TCP   13m
kubedb      None             <none>                   13
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified CRD object:

```yaml
$ kubedb get es -n demo e1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: e1
  namespace: demo
spec:
  certificateSecret:
    secretName: e1-cert
  databaseSecret:
    secretName: e1-auth
  doNotPause: true
  replicas: 1
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: 5.6.4
status:
  creationTime: 2017-12-14T04:04:24Z
  phase: Running
```

Please note that KubeDB operator has created a new Secret called `e1-auth` (format: {crd-name}-auth) for storing the password for `admin` user and `e1-cert` (format: {crd-name}-cert) for storing certificates. If you want to use an existing secret please specify that when creating the CRD using `spec.databaseSecret.secretName`.

Lets edit Service `e1` to set `type: NodePort`.

```console
$ kubectl edit svc e1 -n demo
spec:
  type: NodePort
```

This will provide us URL for `e1` service of our Elasticsearch database

```console
$ minikube service -n demo e1 --https --url
https://192.168.99.100:30653

$ kubectl get secrets -n demo e1-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
wcun4tcq⏎
```

Now, lets connect to this Elasticsearch cluster.

```console
$ curl --user admin:wcun4tcq https://192.168.99.100:30653 --insecure
```

```json
{
  "name" : "e1-0",
  "cluster_name" : "e1",
  "cluster_uuid" : "11TBmi74ThmaSjXUStAA5w",
  "version" : {
    "number" : "5.6.4",
    "build_hash" : "8bbedf5",
    "build_date" : "2017-10-31T18:55:38.105Z",
    "build_snapshot" : false,
    "lucene_version" : "6.6.1"
  },
  "tagline" : "You Know, for Search"
}
```


![Connect to ES](/docs/images/elasticsearch/connect-es.gif)

### Elasticsearch Topology
We can create Elasticsearch database with dedicated pods for `master`, `client` and `data` nodes. Below is the Elasticsearch object created with topology configuration.

```
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: e2
  namespace: demo
spec:
  version: 5.6.4
  topology:
    master:
      replicas: 1
      prefix: master
    data:
      replicas: 2
      prefix: data
    client:
      replicas: 1
      prefix: client
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
```
Here,
- `spec.topology` point to the number of pods we want as dedicated `master`, `client` and `data` nodes and also specify prefix for their StatefulSet name

Now Elasticsearch database has started with 4 pods under 3 different StatefulSets.

```console
$ kubedb describe es -n demo e2
# Only showing additional information
StatefulSet:
  Name:                 client-e2
  Replicas:             1 current / 1 desired
  CreationTimestamp:    Thu, 14 Dec 2017 11:18:25 +0600
  Pods Status:          1 Running / 0 Waiting / 0 Succeeded / 0 Failed

StatefulSet:
  Name:                 data-e2
  Replicas:             2 current / 2 desired
  CreationTimestamp:    Thu, 14 Dec 2017 11:19:29 +0600
  Pods Status:          2 Running / 0 Waiting / 0 Succeeded / 0 Failed

StatefulSet:
  Name:                 master-e2
  Replicas:             1 current / 1 desired
  CreationTimestamp:    Thu, 14 Dec 2017 11:19:07 +0600
  Pods Status:          1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Topology:
  Type     Pod           StartTime                       Phase
  ----     ---           ---------                       -----
  client   client-e2-0   2017-12-14 11:18:32 +0600 +06   Running
  data     data-e2-0     2017-12-14 11:19:36 +0600 +06   Running
  data     data-e2-1     2017-12-14 11:19:55 +0600 +06   Running
  master   master-e2-0   2017-12-14 11:19:14 +0600 +06   Running
```


## Database Snapshots

### Instant Backups
Now, you can easily take a snapshot of this database by creating a `Snapshot` CRD object. When a `Snapshot` object is created, KubeDB operator will launch a Job that runs [elasticdump](https://github.com/taskrabbit/elasticsearch-dump) command and uploads snapshot data to various cloud providers  _S3_, _GCS_, _Azure_, _OpenStack_ _Swift_ and/or locally mounted volumes using [osm](https://github.com/appscode/osm).

In this tutorial, snapshots will be stored in a Google Cloud Storage (GCS) bucket. To do so, a secret is needed that has the following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic snap-secret -n demo \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "snap-secret" created
```

```yaml
$ kubectl get secret snap-secret -n demo -o yaml

apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  name: snap-secret
  namespace: demo
type: Opaque
```


To lean how to configure other storage destinations for Snapshots, please visit [here](/docs/concepts/snapshot.md). Now, create the Snapshot object.

```console
$ kubedb create -f ./docs/examples/elasticsearch/demo-2.yaml
validating "./docs/examples/elasticsearch/demo-2.yaml"
snapshot "e1-xyz" created

$ kubedb get snap -n demo
NAME      DATABASE   STATUS      AGE
e1-xyz    es/e1      Succeeded   2m
```

```yaml
$ kubedb get snap -n demo e1-xyz -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  labels:
    kubedb.com/kind: Elasticsearch
    kubedb.com/name: e1
  name: e1-xyz
  namespace: demo
spec:
  databaseName: e1
  gcs:
    bucket: kubedb
  storageSecretName: snap-secret
status:
  completionTime: 2017-12-14T05:43:40Z
  phase: Succeeded
  startTime: 2017-12-14T05:41:46Z
```

Here,

- `metadata.labels` should include the type of database `kubedb.com/kind: Elasticsearch`.
- `spec.databaseName` points to the database whose snapshot is taken.
- `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.
- `spec.gcs.bucket` points to the bucket name used to store the snapshot data.


You can also run the `kubedb describe` command to see the recent snapshots taken for a database.

```console
$ kubedb describe es -n demo e1 -S=false -W=false
Name:               e1
Namespace:          demo
CreationTimestamp:  Thu, 14 Dec 2017 11:30:50 +0600
Status:             Running

Topology:
  Type                 Pod       StartTime                       Phase
  ----                 ---       ---------                       -----
  master|client|data   e1-0      2017-12-14 11:31:10 +0600 +06   Running

Snapshots:
  Name     Bucket      StartTime                         CompletionTime                    Phase
  ----     ------      ---------                         --------------                    -----
  e1-xyz   gs:kubedb   Thu, 14 Dec 2017 11:41:46 +0600   Thu, 14 Dec 2017 11:43:40 +0600   Succeeded

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason               Message
  ---------   --------   -----     ----                     --------   ------               -------
  2m          2m         1         Snapshot Controller      Normal     SuccessfulSnapshot   Successfully completed snapshot
  4m          4m         1         Snapshot Controller      Normal     Starting             Backup running
  14m         14m        1         Elasticsearch operator   Normal     SuccessfulCreate     Successfully created Elasticsearch
  15m         15m        1         Elasticsearch operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  15m         15m        1         Elasticsearch operator   Normal     Creating             Creating Kubernetes objects
  15m         15m        1         Elasticsearch operator   Normal     SuccessfulValidate   Successfully validate Elasticsearch
```

Once the snapshot Job is complete, you should see the output of the [elasticdump](https://github.com/taskrabbit/elasticsearch-dump) process stored in the GCS bucket.

![snapshot-console](/docs/images/elasticsearch/e1-xyz-snapshot.png)

From the above image, you can see that the snapshot output is stored in a folder called `{bucket}/kubedb/{namespace}/{CRD object}/{snapshot}/`.


### Scheduled Backups
KubeDB supports taking periodic backups for a database using a [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). To take periodic backups, edit the Elasticsearch object to add `spec.backupSchedule` section.

```yaml
$ kubedb edit es e1 -n demo
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: snap-secret
    gcs:
      bucket: kubedb
```

Once the `spec.backupSchedule` is added, KubeDB operator will create a new Snapshot object on each tick of the cron expression. This triggers KubeDB operator to create a Job as it would for any regular instant backup process. You can see the snapshots as they are created using `kubedb get snap` command.

```console
$ kubedb get snap -n demo
NAME                 DATABASE   STATUS      AGE
e1-20171214-055100   es/e1      Running     31s
e1-xyz               es/e1      Succeeded   9m
```

### Restore from Snapshot
You can create a new database from a previously taken Snapshot. Specify the Snapshot name in the `spec.init.snapshotSource` field of a new Elasticsearch object. See the example `recovered` object below:

```yaml
# See full YAML file here: /docs/examples/elasticsearch/demo-4.yaml
  init:
    snapshotSource:
      namespace: demo
      name: e1-xyz
```

```console
$ kubedb create -f ./docs/examples/elasticsearch/demo-4.yaml
validating "./docs/examples/elasticsearch/demo-4.yaml"
elasticsearch "recovered" created
```

Here,
- `spec.init.snapshotSource` specifies Snapshot object information to be used in restoration process.
	- `snapshotSource.name` refers to a Snapshot object name
	- `snapshotSource.namespace` refers to a Snapshot object namespace

Now, wait several seconds. KubeDB operator will create a new StatefulSet. Then KubeDB operator launches a Kubernetes Job to initialize the new database using the data from `e1-xyz` Snapshot.

```console
$ kubedb get es -n demo
NAME        STATUS    AGE
e1          Running   1h
recovered   Running   49s

$ kubedb describe es -n demo recovered -S=false -W=false
Name:               recovered
Namespace:          demo
CreationTimestamp:  Thu, 14 Dec 2017 12:08:57 +0600
Status:             Running
Replicas:           1  total
Init:
  snapshotSource:
    namespace:  demo
    name:       e1-xyz
StatefulSet:    recovered
Service:        recovered, recovered-master
Secrets:        recovered-auth, recovered-cert

Topology:
  Type                 Pod             StartTime                       Phase
  ----                 ---             ---------                       -----
  data|master|client   recovered-0   2017-12-14 12:09:15 +0600 +06   Running

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason                 Message
  ---------   --------   -----     ----                     --------   ------                 -------
  1m          1m         1         Elasticsearch operator   Normal     SuccessfulCreate       Successfully created Elasticsearch
  1m          1m         1         Elasticsearch operator   Normal     SuccessfulInitialize   Successfully completed initialization
  3m          3m         1         Elasticsearch operator   Normal     Initializing           Initializing from Snapshot: "e1-xyz"
  3m          3m         1         Elasticsearch operator   Normal     SuccessfulCreate       Successfully created StatefulSet
  4m          4m         1         Elasticsearch operator   Normal     SuccessfulValidate     Successfully validate Elasticsearch
  4m          4m         1         Elasticsearch operator   Normal     Creating               Creating Kubernetes objects
```


## Pause Database

Since the Elasticsearch `e1` has `spec.doNotPause` set to true, if you delete the object, KubeDB operator will recreate original Elasticsearch object and essentially nullify the delete operation. You can see this below:

```console
$ kubedb delete es e1 -n demo
error: Elasticsearch "e1" can't be paused. To continue delete, unset spec.doNotPause and retry.
```

Now, run `kubedb edit es e1 -n demo` to set `spec.doNotPause` to false or remove this field (which default to false). Then if you delete the Elasticsearch object, KubeDB operator will delete the StatefulSet and its pods, but leaves the PVCs unchanged. In KubeDB parlance, we say that **e1** Elasticsearch database has entered into dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase CRD object.

```console
$ kubedb delete es -n demo e1
elasticsearch "e1" deleted

$ kubedb get drmn -n demo e1
NAME    STATUS  AGE
e1      Paused  3m
```

```yaml
$ kubedb get drmn -n demo e1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  labels:
    kubedb.com/kind: Elasticsearch
  name: e1
  namespace: demo
spec:
  origin:
    metadata:
      name: e1
      namespace: demo
    spec:
      elasticsearch:
        certificateSecret:
          secretName: e1-cert
        databaseSecret:
          secretName: e1-auth
        replicas: 1
        storage:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 50Mi
          storageClassName: standard
        version: 5.6.4
status:
  creationTime: 2017-12-14T06:18:53Z
  pausingTime: 2017-12-14T06:19:07Z
  phase: Paused
```

Here,
 - `spec.origin` is the spec of the original spec of the original Elasticsearch object.
 - `status.phase` points to the current database state `Paused`.


## Resume Dormant Database

To resume the database from the dormant state, set `spec.resume` to `true` in the DormantDatabase object.

```yaml
$ kubedb edit drmn -n demo e1
spec:
  resume: true
```

KubeDB operator will notice that `spec.resume` is set to true. KubeDB operator will delete the DormantDatabase object and create a new Elasticsearch using the original spec. This will in turn start a new StatefulSet which will mount the originally created PVCs. Thus the original database is resumed.

## Wipeout Dormant Database
You can also wipe out a DormantDatabase by setting `spec.wipeOut` to true. KubeDB operator will delete the PVCs, delete any relevant Snapshot for this database and also delete snapshot data stored in the Cloud Storage buckets. There is no way to resume a wiped out database. So, be sure before you wipe out a database.

```yaml
$ kubedb edit drmn -n demo e1
spec:
  wipeOut: true
```

When database is completely wiped out, you can see status `WipedOut`

```console
$ kubedb get drmn -n demo e1
NAME      STATUS     AGE
e1        WipedOut   4m
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
