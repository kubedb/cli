---
title: Continuous Archiving to Swift
menu:
  docs_0.11.0:
    identifier: pg-continuous-archiving-snapshot
    name: Archiving to Swift
    parent: pg-snapshot-postgres
    weight: 40
menu_name: docs_0.11.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Continuous Archiving  to Swift

**WAL-G** is used to handle continuous archiving mechanism. Please refer to [continuous archiving in kubeDB](/docs/guides/postgres/snapshot/continuous_archiving.md) to know more about it.

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

For archiving, we need storage Secret, and storage backend information. Below is a Postgres object created with Continuous Archiving support to backup WAL files to Swift Storage.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: wal-postgres
  namespace: demo
spec:
  version: "11.1-v1"
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
      storageSecretName: swift-secret
      gcs:
        bucket: kubedb
```

Here,

- `spec.archiver.storage` specifies storage information that will be used by `WAL-G`
  - `storage.storageSecretName` points to the Secret containing the credentials for cloud storage destination.
  - `storage.swift` points to Swift storage configuration.
  - `storage.swift.container` points to the bucket name used to store continuous archiving data.

**Archiver Storage Secret**

Storage Secret should contain credentials that will be used to access storage destination.

Storage Secret for **WAL-G** is needed with the full set of v1, v2 or v3 authentication keys from the following list:

| Key                      | Description                        |
| ------------------------ | ---------------------------------- |
| `ST_AUTH`                | For keystone v1 authentication     |
| `ST_USER`                | For keystone v1 authentication     |
| `ST_KEY`                 | For keystone v1 authentication     |
| `OS_AUTH_URL`            | For keystone v2 authentication     |
| `OS_REGION_NAME`         | For keystone v2 authentication     |
| `OS_USERNAME`            | For keystone v2 authentication     |
| `OS_PASSWORD`            | For keystone v2 authentication     |
| `OS_TENANT_ID`           | For keystone v2 authentication     |
| `OS_TENANT_NAME`         | For keystone v2 authentication     |
| `OS_AUTH_URL`            | For keystone v3 authentication     |
| `OS_REGION_NAME`         | For keystone v3 authentication     |
| `OS_USERNAME`            | For keystone v3 authentication     |
| `OS_PASSWORD`            | For keystone v3 authentication     |
| `OS_USER_DOMAIN_NAME`    | For keystone v3 authentication     |
| `OS_PROJECT_NAME`        | For keystone v3 authentication     |
| `OS_PROJECT_DOMAIN_NAME` | For keystone v3 authentication     |
| `OS_STORAGE_URL`         | For authentication based on tokens |
| `OS_AUTH_TOKEN`          | For authentication based on tokens |

```console
$ echo -n '<your-swift-storage-account-name>' > AZURE_ACCOUNT_NAME
$ echo -n '<your-swift-storage-account-key>' > AZURE_ACCOUNT_KEY
$ kubectl create secret generic swift-secret \
    --from-file=./AZURE_ACCOUNT_NAME \
    --from-file=./AZURE_ACCOUNT_KEY
secret "swift-secret" created
```

```yaml
$ kubectl get secret swift-secret -o yaml
apiVersion: v1
data:
  OS_AUTH_URL: PHlvdXItYXV0aC11cmw+
  OS_PASSWORD: PHlvdXItcGFzc3dvcmQ+
  OS_REGION_NAME: PHlvdXItcmVnaW9uPg==
  OS_TENANT_ID: PHlvdXItdGVuYW50LWlkPg==
  OS_TENANT_NAME: PHlvdXItdGVuYW50LW5hbWU+
  OS_USERNAME: PHlvdXItdXNlcm5hbWU+
kind: Secret
metadata:
  creationTimestamp: 2017-07-03T19:17:39Z
  name: swift-secret
  namespace: default
  resourceVersion: "36381"
  selfLink: /api/v1/namespaces/default/secrets/swift-secret
  uid: 47b4bcab-6024-11e7-879a-080027726d6b
type: Opaqu
```

**Archiver Storage Backend**

To configure GCS backend, following parameters are available:

| Parameter              | Description                                                  |
| ---------------------- | ------------------------------------------------------------ |
| `spec.azure.container` | `Required`. Name of Storage container                        |
| `spec.azure.prefix`    | `Optional`. Path prefix into bucket where snapshot will be stored |

Now create this Postgres object with Continuous Archiving support.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.11.0/docs/examples/postgres/snapshot/wal-postgres-swift.yaml
postgres.kubedb.com/wal-postgres created
```

When database is ready, **WAL-G** takes a base backup and uploads it to the cloud storage defined by storage backend.

Archived data is stored in a folder called `{container}/{prefix}/kubedb/{namespace}/{postgres-name}/archive/`.

You can see continuous archiving data stored in swift container.

<p align="center">
  <kbd>
    <img alt="continuous-archiving"  src="/docs/images/postgres/wal-postgres-swift.png">
  </kbd>
</p>



From the above image, you can see that the archived data is stored in a folder `kubedb/kubedb/demo/wal-postgres/archive`.

## Termination Policy

If termination policy of this `wal-postgres` is set to `WipeOut` or, If `Spec.WipeOut` of dormant database is set to `true`, then the data in cloud backend will be deleted.

The data will be intact in other scenarios.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo pg/wal-postgres -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo pg/wal-postgres

kubectl delete -n demo secret/swift-secret
kubectl delete ns demo
```

## Next Steps

- Learn about initializing [PostgreSQL from WAL](/docs/guides/postgres/initialization/script_source.md) files stored in cloud.

