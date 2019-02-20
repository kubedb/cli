---
title: Continuous Archiving of PostgreSQL
menu:
  docs_0.9.0:
    identifier: pg-continuous-archiving-snapshot
    name: WAL Archiving
    parent: pg-snapshot-postgres
    weight: 20
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Continuous Archiving with wal-g

KubeDB PostgreSQL also supports continuous archiving using [wal-g ](https://github.com/wal-g/wal-g). Now **wal-g** supports _S3_ and _GCP_ as cloud storage.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace/demo created
```

> Note: YAML files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create PostgreSQL with Continuous Archiving

Below is the Postgres object created with Continuous Archiving support.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: wal-postgres
  namespace: demo
spec:
  version: "9.6-v2"
  replicas: 2
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  archiver:
    storage:
      storageSecretName: s3-secret
      s3:
        bucket: kubedb
```

Here,

- `spec.archiver.storage` specifies storage information that will be used by `wal-g`
  - `storage.storageSecretName` points to the Secret containing the credentials for cloud storage destination.
  - `storage.s3` points to s3 storage configuration.
  - `storage.s3.bucket` points to the bucket name used to store continuous archiving data.
  - `storage.gcs` points to GCS storage configuration.
  - `storage.gcs.bucket` points to the bucket name used to store continuous archiving data.

User can use either s3 or gcs. In this tutorial, s3 is used for wal-g archiving. `gcs` is similar to this tutorial. Follow [this link](/docs/concepts/snapshot/#google-cloud-storage-gcs) to know how to create secret for `gcs` storage. 

**What is this Continuous Archiving**

PostgreSQL maintains a write ahead log (WAL) in the `pg_xlog/` subdirectory of the cluster's data directory.  The existence of the log makes it possible to use a third strategy for backing up databases and if recovery is needed, restore from the backed-up WAL files to bring the system to a current state.

**Continuous Archiving Setup**

KubeDB PostgreSQL supports [wal-g](https://github.com/wal-g/wal-g) for this continuous archiving.

Following additional parameters are set in `postgresql.conf` for *primary* server

```console
archive_command = 'wal-g wal-push %p'
archive_timeout = 60
```

Here, these commands are used to push and pull WAL files respectively from cloud.

**wal-g** is used to handle this continuous archiving mechanism. For this we need storage Secret and need to provide storage backend information.

**Archiver Storage Secret**

Storage Secret should contain credentials that will be used to access storage destination.

Storage Secret for **wal-g** is needed with following 2 keys:

|           Key           |                Description                |
| ----------------------- | ----------------------------------------- |
| `AWS_ACCESS_KEY_ID`     | `Required`. AWS / Minio access key ID     |
| `AWS_SECRET_ACCESS_KEY` | `Required`. AWS / Minio secret access key |

```console
$ echo -n '<your-aws-access-key-id-here>' > AWS_ACCESS_KEY_ID
$ echo -n '<your-aws-secret-access-key-here>' > AWS_SECRET_ACCESS_KEY
$ kubectl create secret -n demo generic s3-secret \
    --from-file=./AWS_ACCESS_KEY_ID \
    --from-file=./AWS_SECRET_ACCESS_KEY
secret "s3-secret" created
```

```yaml
$ kubectl get secret -n demo s3-secret -o yaml
apiVersion: v1
data:
  AWS_ACCESS_KEY_ID: PHlvdXItYXdzLWFjY2Vzcy1rZXktaWQtaGVyZT4=
  AWS_SECRET_ACCESS_KEY: PHlvdXItYXdzLXNlY3JldC1hY2Nlc3Mta2V5LWhlcmU+
kind: Secret
metadata:
  creationTimestamp: 2018-02-06T09:12:37Z
  name: s3-secret
  namespace: demo
  resourceVersion: "59225"
  selfLink: /api/v1/namespaces/demo/secrets/s3-secret
  uid: dfbe6b06-0b1d-11e8-9fb9-42010a800064
type: Opaque
```

**Archiver Storage Backend**

**wal-g** supports both _S3_ and __GCS__ cloud providers.

To configure s3 backend, following parameters are available:

|     Parameter      |                           Description                            |
| ------------------ | ---------------------------------------------------------------- |
| `spec.s3.endpoint` | `Required`. For S3, use `s3.amazonaws.com`                       |
| `spec.s3.bucket`   | `Required`. Name of Bucket                                       |
| `spec.s3.prefix`   | `Optional`. Path prefix into bucket where snapshot will be store |

Now create this Postgres object with Continuous Archiving support.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/postgres/snapshot/wal-postgres.yaml
postgres.kubedb.com/wal-postgres created
```

When database is ready, **wal-g** takes a base backup and uploads it to cloud storage defined by storage backend.

Archived data is stored in a folder called `{bucket}/{prefix}/kubedb/{namespace}/{postgres-name}/archive/`.

you can see continuous archiving data stored in S3 bucket.

<p align="center">
  <kbd>
    <img alt="continuous-archiving"  src="/docs/images/postgres/wal-postgres.png">
  </kbd>
</p>

From the above image, you can see that the archived data is stored in a folder `kubedb/kubedb/demo/wal-postgres/archive`.

## Termination Policy

If termination policy of this `wal-postgres` is set to `WipeOut` or, If `Spec.WipeOut` of dormant database is set to `true`, then the data in cloud backend will be deleted.

Other than that, the data will be intact in other scenarios.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo pg/wal-postgres -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo pg/wal-postgres

kubectl delete -n demo secret/s3-secret
kubectl delete ns demo
```

## Next Steps

- Learn about initializing [PostgreSQL from WAL](/docs/guides/postgres/initialization/script_source.md) files stored in cloud.
