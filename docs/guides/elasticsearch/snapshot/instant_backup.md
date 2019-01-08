---
title: Instant Backup of Elasticsearch
menu:
  docs_0.9.0:
    identifier: es-instant-backup-snapshot
    name: Instant Backup
    parent: es-snapshot-elasticsearch
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# KubeDB Snapshot

KubeDB operator maintains another Custom Resource Definition (CRD) for database backups called Snapshot. Snapshot object is used to take backup or restore from a backup.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: YAML files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Prepare Database

We need an Elasticsearch object in `Running` phase to perform backup operation. If you do not already have an Elasticsearch instance running, create one first.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/quickstart/infant-elasticsearch.yaml
elasticsearch "infant-elasticsearch" created
```

Below the YAML for the Elasticsearch crd we have created above.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: infant-elasticsearch
  namespace: demo
spec:
  version: "6.3-v1"
  storageType: Ephemeral
```

Here, we have used `spec.storageType: Ephemeral`. So, we don't need to specify storage section. KubeDB will use [emptyDir]((https://kubernetes.io/docs/concepts/storage/volumes/#emptydir)) volume for this database.

Verify that the Elasticsearch is running,

```console
$ kubedb get es -n demo infant-elasticsearch
NAME                   STATUS    AGE
infant-elasticsearch   Running   11m
```

### Populate database

Let's insert some data so that we can verify that the snapshot contains those data. Check how to connect with the database from [here](/docs/guides/elasticsearch/quickstart/quickstart.md#connect-with-elasticsearch-database).

```console
$ curl -XPUT --user "admin:fqvzdvz3" "localhost:9200/test/snapshot/1?pretty" -H 'Content-Type: application/json' -d'
{
    "title": "Snapshot",
    "text":  "Testing instand backup",
    "date":  "2018/02/13"
}
'
```

```console
$ curl -XGET --user "admin:fqvzdvz3" "localhost:9200/test/snapshot/1?pretty"
```

```json
{
  "_index" : "test",
  "_type" : "snapshot",
  "_id" : "1",
  "_version" : 1,
  "found" : true,
  "_source" : {
    "title" : "Snapshot",
    "text" : "Testing instand backup",
    "date" : "2018/02/13"
  }
}
```

Now, we are ready to take backup of this database `infant-elasticsearch`.

## Instant backup

Snapshot provides a declarative configuration for backup behavior in a Kubernetes native way.

Below is the Snapshot object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: instant-snapshot
  namespace: demo
  labels:
    kubedb.com/kind: Elasticsearch
spec:
  databaseName: infant-elasticsearch
  storageSecretName: gcs-secret
  gcs:
    bucket: kubedb
```

Here,

- `metadata.labels` should include the type of database.
- `spec.databaseName` indicates the Elasticsearch object name, `infant-elasticsearch`, whose snapshot is taken.
- `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.
- `spec.gcs.bucket` points to the bucket name used to store the snapshot data.

In this case, `kubedb.com/kind: Elasticsearch` tells KubeDB operator that this Snapshot belongs to an Elasticsearch object. Only Elasticsearch controller will handle this Snapshot object.

> Note: Snapshot and Secret objects must be in the same namespace as Elasticsearch, `infant-elasticsearch`.

#### Snapshot Storage Secret

Storage Secret should contain credentials that will be used to access storage destination.
In this tutorial, snapshot data will be stored in a Google Cloud Storage (GCS) bucket.

For that a storage Secret is needed with following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret -n demo generic gcs-secret \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "gcs-secret" created
```

```yaml
$ kubectl get secret -n demo gcs-secret -o yaml
apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2018-02-13T06:35:36Z
  name: gcs-secret
  namespace: demo
  resourceVersion: "4308"
  selfLink: /api/v1/namespaces/demo/secrets/gcs-secret
  uid: 19a77054-1088-11e8-9e42-0800271bdbb6
type: Opaque
```

#### Snapshot storage backend

KubeDB supports various cloud providers (_S3_, _GCS_, _Azure_, _OpenStack_ _Swift_ and/or locally mounted volumes) as snapshot storage backend. In this tutorial, _GCS_ backend is used.

To configure this backend, following parameters are available:

| Parameter                | Description                                                                     |
|--------------------------|---------------------------------------------------------------------------------|
| `spec.gcs.bucket`        | `Required`. Name of bucket                                                      |
| `spec.gcs.prefix`        | `Optional`. Path prefix into bucket where snapshot data will be stored          |

> An open source project [osm](https://github.com/appscode/osm) is used to store snapshot data into cloud.

To learn how to configure other storage destinations for snapshot data, please visit [here](/docs/concepts/snapshot.md).

Now, create the Snapshot object.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/snapshot/instant-snapshot.yaml
snapshot.kubedb.com/instant-snapshot created
```

Let's see Snapshot list of Elasticsearch `infant-elasticsearch`.

```console
$ kubectl get snap -n demo --selector=kubedb.com/kind=Elasticsearch,kubedb.com/name=infant-elasticsearch
NAME               DATABASENAME           STATUS      AGE
instant-snapshot   infant-elasticsearch   Succeeded   47s
```

KubeDB operator watches for Snapshot objects using Kubernetes API. When a Snapshot object is created, it will launch a Job that runs the [elasticdump](https://github.com/taskrabbit/elasticsearch-dump) command and uploads the output files to cloud storage using [osm](https://github.com/appscode/osm).

Snapshot data is stored in a folder called `{bucket}/{prefix}/kubedb/{namespace}/{elasticsearch}/{snapshot}/`.

Once the snapshot Job is completed, you can see the output of the `elasticdump` command stored in the GCS bucket.

<p align="center">
  <kbd>
    <img alt="snapshot-console"  src="/docs/images/elasticsearch/instant-backup.png">
  </kbd>
</p>

From the above image, you can see that the snapshot data files for index `test` are stored in your bucket.

If you open this `test.data.json` file, you will see the data you have created previously.

```json
{
   "_index":"test",
   "_type":"snapshot",
   "_id":"1",
   "_score":1,
   "_source":{
      "title":"Snapshot",
      "text":"Testing instand backup",
      "date":"2018/02/13"
   }
}
```

Let's see the Snapshot list for Elasticsearch `infant-elasticsearch` by running `kubedb describe` command.

```console
$ kubedb describe es -n demo infant-elasticsearch
Name:               infant-elasticsearch
Namespace:          demo
CreationTimestamp:  Fri, 05 Oct 2018 16:45:56 +0600
Labels:             <none>
Annotations:        kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"kubedb.com/v1alpha1","kind":"Elasticsearch","metadata":{"annotations":{},"name":"infant-elasticsearch","namespace":"demo"},"spec":{"repl...
Status:             Running
Replicas:           1  total
  StorageType:      Ephemeral
Volume:
  Capacity:  0

StatefulSet:          
  Name:               infant-elasticsearch
  CreationTimestamp:  Fri, 05 Oct 2018 16:45:58 +0600
  Labels:               kubedb.com/kind=Elasticsearch
                        kubedb.com/name=infant-elasticsearch
                        node.role.client=set
                        node.role.data=set
                        node.role.master=set
  Annotations:        <none>
  Replicas:           824639991608 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

...
Topology:
  Type                Pod                     StartTime                      Phase
  ----                ---                     ---------                      -----
  master|client|data  infant-elasticsearch-0  2018-10-05 16:45:58 +0600 +06  Running

Snapshots:
  Name              Bucket     StartTime                        CompletionTime                   Phase
  ----              ------     ---------                        --------------                   -----
  instant-snapshot  gs:kubedb  Fri, 05 Oct 2018 17:27:55 +0600  Fri, 05 Oct 2018 17:28:10 +0600  Succeeded

Events:
  Type    Reason              Age   From                    Message
  ----    ------              ----  ----                    -------
  Normal  Successful          44m   Elasticsearch operator  Successfully created Service
  Normal  Successful          44m   Elasticsearch operator  Successfully created Service
  Normal  Successful          44m   Elasticsearch operator  Successfully created StatefulSet
  Normal  Successful          44m   Elasticsearch operator  Successfully created Elasticsearch
  Normal  Successful          44m   Elasticsearch operator  Successfully patched StatefulSet
  Normal  Successful          43m   Elasticsearch operator  Successfully patched Elasticsearch
  Normal  Successful          43m   Elasticsearch operator  Successfully patched StatefulSet
  Normal  Successful          43m   Elasticsearch operator  Successfully patched Elasticsearch
  Normal  Starting            2m    Job Controller          Backup running
  Normal  SuccessfulSnapshot  2m    Job Controller          Successfully completed snapshot
```

From the above output, we can see in `Snapshots:` section that we have one successful snapshot.

## Delete Snapshot

If you want to delete snapshot data from storage, you can delete Snapshot object.

```console
$ kubectl delete snap -n demo instant-snapshot
snapshot "instant-snapshot" deleted
```

Once Snapshot object is deleted, you can't revert this process and snapshot data from storage will be deleted permanently.

## Customizing Snapshot

You can customize pod template spec and volume claim spec for the backup and restore jobs. For details options read [this doc](/docs/concepts/snapshot.md).

Some common customization sample is shown below.

**Specify PVC Template:**

Backup and recovery job needs a temporary storage to hold `dump` files before it can be uploaded to cloud backend or inserted into database. By default, KubeDB reads storage specification from `spec.storage` section of database crd and creates PVC with similar specification for backup or recovery job. However, if you want to specify custom PVC template, you can do it through `spec.podVolumeClaimSpec` field of Snapshot crd. This is particularly helpful when you want to use different `storageclass` for backup or recovery job than the database.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: instant-snapshot
  namespace: demo
  labels:
    kubedb.com/kind: Elasticsearch
spec:
  databaseName: infant-elasticsearch
  storageSecretName: gcs-secret
  gcs:
    bucket: kubedb-dev
  podVolumeClaimSpec:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi # make sure size is larger or equal than your database size
```

**Specify Resources for Backup/Recovery Job:**

You can specify resources for backup or recovery job through `spec.podTemplate.spec.resources` field.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: instant-snapshot
  namespace: demo
  labels:
    kubedb.com/kind: Elasticsearch
spec:
  databaseName: infant-elasticsearch
  storageSecretName: gcs-secret
  gcs:
    bucket: kubedb-dev
  podTemplate:
    spec:
      resources:
        requests:
          memory: "64Mi"
          cpu: "250m"
        limits:
          memory: "128Mi"
          cpu: "500m"
```

**Provide Annotation for Backup/Recovery Job:**

If you need to add some annotations to backup or recovery job, you can specify this in `spec.podTemplate.controller.annotations`. You can also specify annotation for the pod created by backup or recovery job through `spec.podTemplate.annotations` field.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: instant-snapshot
  namespace: demo
  labels:
    kubedb.com/kind: Elasticsearch
spec:
  databaseName: infant-elasticsearch
  storageSecretName: gcs-secret
  gcs:
    bucket: kubedb-dev
  podTemplate:
    annotations:
      passMe: ToBackupJobPod
    controller:
      annotations:
        passMe: ToBackupJob
```

**Pass Arguments to Backup/Recovery Job:**

KubeDB also allows to pass extra arguments for backup or recovery job. You can provide these arguments through `spec.podTemplate.spec.args` field of Snapshot crd.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: instant-snapshot
  namespace: demo
  labels:
    kubedb.com/kind: Elasticsearch
spec:
  databaseName: infant-elasticsearch
  storageSecretName: gcs-secret
  gcs:
    bucket: kubedb-dev
  podTemplate:
    spec:
      args:
      - --extra-args-to-backup-command
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/infant-elasticsearch -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo es/infant-elasticsearch

$ kubectl delete ns demo
```

## Next Steps

- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Learn how to [schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Learn about initializing [Elasticsearch with Snapshot](/docs/guides/elasticsearch/initialization/snapshot_source.md).
- Learn how to configure [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
