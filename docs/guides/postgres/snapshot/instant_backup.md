---
title: Instant Backup of PostgreSQL
menu:
  docs_0.9.0:
    identifier: pg-instant-backup-snapshot
    name: Instant Backup
    parent: pg-snapshot-postgres
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Database Snapshot

KubeDB operator maintains another Custom Resource Definition (CRD) for database backups called Snapshot. Snapshot object is used to take backup or restore from a backup. For more details about Snapshot please visit [here](/docs/concepts/snapshot.md).

This tutorial will show how to take instant backup of PostgreSQL database deployed with KubeDB.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: YAML files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Prepare Database

We need an Postgres database running to perform backup operation. If you don't have a Postgres instance running, create one and initialize it by following the tutorial [here](/docs/guides/postgres/initialization/script_source.md).

## Instant Backup

KubeDB operator watches for Snapshot objects using Kubernetes API. When a Snapshot object is created, it will launch a Job that runs the `pg_dumpall` command and uploads the output **sql** file to cloud storage using [osm](https://github.com/appscode/osm).

Below the Snapshot object that will be created in this tutorial,

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: instant-snapshot
  namespace: demo
  labels:
    kubedb.com/kind: Postgres
spec:
  databaseName: script-postgres
  storageSecretName: gcs-secret
  gcs:
    bucket: kubedb
```

Here,

- `metadata.labels` should include the type of database.
- `spec.databaseName` indicates the Postgres object name, `script-postgres`, whose snapshot is taken.
- `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.
- `spec.gcs.bucket` points to the bucket name used to store the snapshot data.

In this case, `kubedb.com/kind: Postgres` tells KubeDB operator that this Snapshot belongs to a Postgres object. Only PostgreSQL controller will handle this Snapshot object.

> Note: Snapshot and Secret objects must be in the same namespace as Postgres, `script-postgres`, in our case.

### Snapshot Storage Secret

Storage Secret should contain credentials that will be used to access storage destination. In this tutorial, snapshot data will be stored in a Google Cloud Storage (GCS) bucket.

For GCS bucket, a storage Secret require to have following 2 keys:

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
$  kubectl get secret -n demo gcs-secret -o yaml
apiVersion: v1
data:
  GOOGLE_PROJECT_ID: <base64 encoded project id>
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: <base64 encoded service-account-json-key>
kind: Secret
metadata:
  creationTimestamp: 2018-09-04T06:08:01Z
  name: gcs-secret
  namespace: demo
  resourceVersion: "11716"
  selfLink: /api/v1/namespaces/demo/secrets/gcs-secret
  uid: e0aef5a7-b008-11e8-9990-0800279292a5
type: Opaque

```

### Snapshot Storage Backend

KubeDB supports various cloud providers (_S3_, _GCS_, _Azure_, _OpenStack_ _Swift_ and/or locally mounted volumes) as snapshot storage backend. In this tutorial, _GCS_ backend is used.

To configure this backend, following parameters are available:

| Parameter                | Description                                                                     |
|--------------------------|---------------------------------------------------------------------------------|
| `spec.gcs.bucket`        | `Required`. Name of bucket                                                      |
| `spec.gcs.prefix`        | `Optional`. Path prefix into bucket where snapshot data will be stored          |

To learn how to configure other storage destinations for snapshot data, please visit [here](/docs/concepts/snapshot.md).

Now, let's create a Snapshot object.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/postgres/snapshot/instant-snapshot.yaml
snapshot.kubedb.com/instant-snapshot created
```

Verify that the Snapshot has been successfully created,

```console
$ kubectl get snap -n demo --selector="kubedb.com/kind=Postgres,kubedb.com/name=script-postgres"
NAME               DATABASENAME      STATUS    AGE
instant-snapshot   script-postgres   Running   58s
```

Notice that the `STATUS` field is showing `Running`. It means the backup is running.

Snapshot data is stored in the backend in following directory `{bucket}/{prefix}/kubedb/{namespace}/{PostgreSQL name}/{Snapshot name}/`.

Once the snapshot Job is completed, you can see the output of the `pg_dumpall` command stored in the GCS bucket.

Verify that the backup has been completed successfully using following command,

```console
$ kubectl get snap -n demo --selector="kubedb.com/kind=Postgres,kubedb.com/name=script-postgres"
NAME               DATABASENAME      STATUS      AGE
instant-snapshot   script-postgres   Succeeded   36s
```

Here, `STATUS` `Succeeded` means the backup has been completed successfully. Now, navigate to the bucket to see the backed up file.

<p align="center">
  <kbd>
    <img alt="snapshot-console"  src="/docs/images/postgres/instant-snapshot.png">
  </kbd>
</p>

From the above image, you can see that the snapshot data file `dumpfile.sql` is stored in your bucket.

If you open this `dumpfile.sql` file, you will see the query to create `dashboard` TABLE.

```console

--
-- Name: dashboard; Type: TABLE; Schema: data; Owner: postgres
--

CREATE TABLE dashboard (
    id bigint NOT NULL,
    version integer NOT NULL,
    slug character varying(255) NOT NULL,
    title character varying(255) NOT NULL,
    data text NOT NULL,
    org_id bigint NOT NULL,
    created timestamp without time zone NOT NULL,
    updated timestamp without time zone NOT NULL,
    updated_by integer,
    created_by integer
);


ALTER TABLE dashboard OWNER TO postgres;
```

You can see the Snapshot list for Postgres `script-postgres` by running `kubedb describe` command.

```console
$ kubedb describe pg -n demo script-postgres
Name:               script-postgres
Namespace:          demo
CreationTimestamp:  Tue, 04 Sep 2018 11:55:22 +0600
Labels:             <none>
Annotations:        kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"kubedb.com/v1alpha1","kind":"Postgres","metadata":{"annotations":{},"name":"script-postgres","namespace":"demo"},"spec":{"init":{"script...
Replicas:           1  total
Status:             Running
Init:
  scriptSource:
Volume:
    Type:       ConfigMap (a volume populated by a ConfigMap)
    Name:       pg-init-script
    Optional:   false
  StorageType:  Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:          
  Name:               script-postgres
  CreationTimestamp:  Tue, 04 Sep 2018 11:55:25 +0600
  Labels:               kubedb.com/kind=Postgres
                        kubedb.com/name=script-postgres
  Annotations:        <none>
  Replicas:           824640513680 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         script-postgres
  Labels:         kubedb.com/kind=Postgres
                  kubedb.com/name=script-postgres
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.99.38.101
  Port:         api  5432/TCP
  TargetPort:   api/TCP
  Endpoints:    172.17.0.6:5432

Service:        
  Name:         script-postgres-replicas
  Labels:         kubedb.com/kind=Postgres
                  kubedb.com/name=script-postgres
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.106.251.230
  Port:         api  5432/TCP
  TargetPort:   api/TCP
  Endpoints:    172.17.0.6:5432

Database Secret:
  Name:         script-postgres-auth
  Labels:         kubedb.com/kind=Postgres
                  kubedb.com/name=script-postgres
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  POSTGRES_PASSWORD:  16 bytes
  POSTGRES_USER:      8 bytes

Topology:
  Type     Pod                StartTime                      Phase
  ----     ---                ---------                      -----
  primary  script-postgres-0  2018-09-04 11:55:32 +0600 +06  Running

Snapshots:
  Name              Bucket     StartTime                        CompletionTime                   Phase
  ----              ------     ---------                        --------------                   -----
  instant-snapshot  gs:kubedb  Tue, 04 Sep 2018 12:10:54 +0600  Tue, 04 Sep 2018 12:11:45 +0600  Succeeded

Events:
  Type    Reason              Age   From               Message
  ----    ------              ----  ----               -------
  Normal  Successful          33m   Postgres operator  Successfully created Service
  Normal  Successful          33m   Postgres operator  Successfully created Service
  Normal  Successful          31m   Postgres operator  Successfully created StatefulSet
  Normal  Successful          31m   Postgres operator  Successfully created Postgres
  Normal  Successful          31m   Postgres operator  Successfully patched StatefulSet
  Normal  Successful          31m   Postgres operator  Successfully patched Postgres
  Normal  Starting            17m   Job Controller     Backup running
  Normal  SuccessfulSnapshot  16m   Job Controller     Successfully completed snapshot
```

## Cleanup Snapshot

If you want to delete snapshot data from storage, you can delete Snapshot object.

```console
$ kubectl delete snap -n demo instant-snapshot
snapshot "instant-snapshot" deleted
```

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
    kubedb.com/kind: Postgres
spec:
  databaseName: script-postgres
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
    kubedb.com/kind: Postgres
spec:
  databaseName: script-postgres
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
    kubedb.com/kind: Postgres
spec:
  databaseName: script-postgres
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
    kubedb.com/kind: Postgres
spec:
  databaseName: script-postgres
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
$ kubectl patch -n demo pg/script-postgres -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo pg/script-postgres

$ kubectl delete -n demo configmap/pg-init-script
$ kubectl delete -n demo secret/gcs-secret

$ kubectl delete ns demo
```

## Next Steps

- Setup [Continuous Archiving](/docs/guides/postgres/snapshot/continuous_archiving.md) in PostgreSQL using `wal-g`
- Learn how to [schedule backup](/docs/guides/postgres/snapshot/scheduled_backup.md)  of PostgreSQL database.
- Learn about initializing [PostgreSQL from KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
