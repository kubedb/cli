---
title: Continuous Archiving to MinIO
menu:
  docs_0.12.0:
    identifier: pg-continuous-archiving-minio
    name: WAL Archiving to MinIO
    parent: pg-snapshot-postgres
    weight: 50
menu_name: docs_0.12.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Continuous Archiving to MinIO

[MinIO](https://docs.min.io/) is an open source object storage server compatible with Amazon S3 APIs. **WAL-G** is used to continuously archive PostgreSQL WAL files to MinIO. Please refer to [continuous archiving in KubeDB](/docs/guides/postgres/snapshot/continuous_archiving.md) to learn more.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

> Note: YAML files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace/demo created
```

## Create PostgreSQL with Continuous Archiving

For archiving, a running instance of MinIO server is required. If you don't have one, see [here](https://github.com/appscode/third-party-tools/tree/master/storage/minio) to get a MinIO server configured, and running in your cluster. We will use this storage backend's information to execute connect, archive, and restore operations.
Below is a Postgres object created with Continuous Archiving support that backs up WAL files to a MinIO server.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: wal-postgres-minio
  namespace: demo
spec:
  version: "9.6.7-v4"
  replicas: 2
  updateStrategy:
    type: RollingUpdate
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
        endpoint: https://minio-service.demo.svc:443/
```

Here,

- `spec.archiver.storage` specifies storage information that will be used by `WAL-G`
  - `storage.storageSecretName` points to the Secret containing the credentials for cloud storage destination.
  - `storage.s3` points to s3 api based storage configuration.
  - `storage.s3.bucket` points to the bucket name used to store continuous archiving data.
  - `storage.s3.endpoint` points to the storage location where the bucket can be found.

**Archiver Storage Secret**

Storage Secret should contain credentials that will be used to access storage destination.

Storage Secret for **WAL-G** is needed with the following 2 keys:

| Key                     | Description                               |
| ----------------------- | ----------------------------------------- |
| `AWS_ACCESS_KEY_ID`     | `Required`. AWS / MinIO access key ID     |
| `AWS_SECRET_ACCESS_KEY` | `Required`. AWS / MinIO secret access key |
| `CA_CERT_DATA`          | `Optional`. AWS / MinIO certificate data  |

for MinIO server secured with custom CA,
necessary certificates have to provided in the storage secret as `CA_CERT_DATA` to establish secure connection.

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

To create secret using custom CA:

```console
$ echo -n '<your-aws-access-key-id-here>' > AWS_ACCESS_KEY_ID
$ echo -n '<your-aws-secret-access-key-here>' > AWS_SECRET_ACCESS_KEY
$ kubectl create secret -n demo generic s3-secret \
    --from-file=./AWS_ACCESS_KEY_ID \
    --from-file=./AWS_SECRET_ACCESS_KEY
    --from-file=./CA_CERT_DATA
secret "s3-secret" created
```

```yaml
$ kubectl get secret -n demo s3-secret -o yaml
apiVersion: v1
data:
  AWS_ACCESS_KEY_ID: M0o4Q0FJVjVWU1VWRkUzTTlKWlM=
  AWS_SECRET_ACCESS_KEY: eHlFZ3I0aDVIRXJxYmdReG5RV055Y1RJRTVvK2x3MnZPQ1BPM2Z1Zw==
  CA_CERT_DATA: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN1RENDQWFDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFOTVFzd0NRWURWUVFERXdKallUQWUKRncweE9UQTBNVGd3TmpJNU5EWmFGdzB5T1RBME1UVXdOakk1TkRaYU1BMHhDekFKQmdOVkJBTVRBbU5oTUlJQgpJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBdngyK0xPTUZrN1dBci9pb3BQSS9uNU1qClJDTm9PZ3FEQk1TNktETnl1YXI5NDhxTm9jYlFUMGMzT09vYmsyWDJHckZYK0ZoUCtXQ3l1U2c1NWYzUmhkSUsKb2FDNWpNR1MzNDUzR3JOK2EwMS9HRkRGTnRYdTI3UUJZd2laWFU1eHBPLzRtRUE2L2MwbkVheGxocjZGSlpYcQpPMGJpaTVaWjJMYmhVOTAvNnBFZDJUQ1FBTGhQZEhuaFlnQUkzNExTU1JsWHh3ckJIMDhWTzRURHJxWW9icDh5CmZYbVNqVDFZRFZZL1M3OEV3dGtDcVY3cG9scTRFdGJoeTJGdWdONzEwektxay8zckFFSm40aitnQW5KT05UNUoKSmhHcTU1cXQ2aWFHNGFEUythVWhFQXhEWE9XWWRkQVl6YVBybnE2U0xBNnVlNEhDd0hyOFhXZytjRVJ1TndJRApBUUFCb3lNd0lUQU9CZ05WSFE4QkFmOEVCQU1DQXFRd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBTkJna3Foa2lHCjl3MEJBUXNGQUFPQ0FRRUFKblY4MFYvUHlFS3BvTGtXeXRvUmlncklsNENDdmwvb2pqMjNVQWZPMmcwcGxsc0MKVWV3VnpONkY0Mmh2dHV0RG9CVDBpdlN5TWdRdTc0UExCSkE1UGlYM1ZNVkEvaGxqbVlpQmdOcVVYejN4UmNjeAoyU1lISnkwalc4R0FhekRhN2RzdTE0ZzMzT3V4WUFFdE1raW5ub3o5OFR0NENHWjdwMWduZi9Id2NPWHdhREIyCllQaThtZWZHZzJSMDdZb1NLbVQreTNTUytpUUsyTHJzZnltRlFOY1oySWViZkN3K0dTZG9qdm1uR28wTDY4U1oKNFN2TGl1MUVYRUZGQThhODBCelMzTUFhZy9uS3JhZytzRllBVXNCYTFOY2JnMWd2VVlMb1RqNlRmbEJaYWNWOQpndTdyV1lkeThGd2VGUVJ6c1pScnJrclJyaE9Pb3Z5cFkxR3lKdz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
kind: Secret
metadata:
  creationTimestamp: 2019-04-24T12:08:47Z
  name: s3-secret
  namespace: storage
  resourceVersion: "27413"
  selfLink: /api/v1/namespaces/storage/secrets/s3-secret
  uid: b6a58c8e-6689-11e9-9a61-0800275cdf2b
type: Opaque
```

**Archiver Storage Backend**

To configure s3 backend, following parameters are available:

| Parameter                           | Description                                                        |
| ----------------------------------- | ------------------------------------------------------------------ |
| `spec.archiver.storage.s3.endpoint` | `Required`. For MinIO, use the URL to your server                  |
| `spec.archiver.storage.s3.bucket`   | `Required`. Name of Bucket                                         |
| `spec.archiver.storage.s3.prefix`   | `Optional`. Path prefix into bucket where WAL files will be stored |

Now create this Postgres object with continuous archiving support.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.12.0/docs/examples/postgres/snapshot/wal-postgres-minio.yaml
postgres.kubedb.com/wal-postgres created
```

When database is ready, **WAL-G** takes a base backup and uploads it to the cloud storage defined by storage backend.

Archived data is stored in a folder called `{bucket}/{prefix}/kubedb/{namespace}/{postgres-name}/archive/`.

## Termination Policy

If termination policy of this `wal-postgres` is set to `WipeOut` or, If `Spec.WipeOut` of dormant database is set to `true`, then the data in cloud backend will be deleted.

The data will be intact in other scenarios.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo pg/wal-postgres -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo pg/wal-postgres

kubectl delete ns demo
```

## Next Steps

- Learn about initializing [PostgreSQL from WAL](/docs/guides/postgres/initialization/script_source.md) files stored in backup servers.
