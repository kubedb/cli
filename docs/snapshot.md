> New to KubeDB? Please start [here](/docs/tutorial.md).

# Backup Database

We need to create a Snapshot object to initiate backup process. 

Here is a template of Snapshot object

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: "snapshot-xyz"
  labels:
    kubedb.com/kind: <database TPR kind: Postgres|Elastic>
spec:
  databaseName: "database-demo"
  bucketName: "bucket-for-snapshot"
  storageSecret:
    secretName: "secret-for-bucket"
```

This will create a Snapshot object in `default` namespace.

**L**et me describe this YAML in details.

While taking backup of any database, we must provide three information.

1. Database name (`databaseName:` in spec)
2. Database kind (`kubedb.com/kind:` in labels)

```yaml
metadata:
  labels:
    kubedb.com/kind: <Postgres|Elastic>
spec:
  databaseName: "database-demo"
```

3. Storage information
    * Bucket name
    * Secret to access bucket

```yaml
spec:
  bucketName: "bucket-for-snapshot"
  storageSecret:
    secretName: "secret-for-bucket"
```

Storage secret example:

```yaml
apiVersion: v1
data:
  config: anNvbjogfAog-------dF9pZDogdGlnZXJ3b3Jrcy1rdWJlCg==
  provider: Z29vZ2xl
kind: Secret
metadata:
  name: secret-for-bucket
type: Opaque
```

**T**his storage secret must have two key:
1. Provider (`provider:`)
2. Config (`config: `)

Example:

## Google Cloud Storage (GCS)

* `provider: google`
* `config:`
    ```yaml
    json: |
        {
          "type": "service_account",
          "project_id": "project_id",
          "private_key_id": "private_key_id",
          "private_key": "private_key",
          "client_email": "client_email",
          "client_id": "client_id",
          "auth_uri": "auth_uri",
          "token_uri": "token_uri",
          "auth_provider_x509_cert_url": "auth_provider_x509_cert_url",
          "client_x509_cert_url": "client_x509_cert_url"
        }
    project_id: "project_id"
    ```

## Amazon S3

* `provider: s3`
* `config:`
    ```yaml
    access_key_id: "access_key_id"
    region: "region"
    secret_key: "secret_key"
    ```

## Microsoft Azure Storage
KubeDB can store database snapshots in Microsoft Azure Storage. To configure this, the following secret keys are needed:

| Key                     | Description                                                |
|-------------------------|------------------------------------------------------------|
| `provider`       | `Required`. Password used to encrypt snapshots by `snapshot` |
| `config`    | `Required`. Azure Storage account name                     |
| `AZURE_ACCOUNT_KEY`     | `Required`. Azure Storage account key                      |



* `provider: azure`
* `config:`
    ```yaml
    account: "account_id"
    key: "key_value"
    ```

Before starting backup process, controller will validate storage secret by creating an empty file
in specified bucket using this secret.

**L**ets create a Snapshot object using `snapshot.yaml`.

```bash
$ kubedb create -f snapshot.yaml

snapshot "snapshot-xyz" created
```

We can see its status.

```bash
$ kubedb get snap snapshot-xyz -o wide

NAME           DATABASE              BUCKET    STATUS      AGE
snapshot-xyz   es/elasticsearch-db   snapshot    Succeeded   24m
```









### Local
`Local` backend refers to a local path inside `stash` sidecar container. Any Kubernetes supported [persistent volume](https://kubernetes.io/docs/concepts/storage/volumes/) can be used here. Some examples are: `emptyDir` for testing, NFS, Ceph, GlusterFS, etc. To configure this backend, following secret keys are needed:

| Key               | Description                                                |
|-------------------|------------------------------------------------------------|
| `RESTIC_PASSWORD` | `Required`. Password used to encrypt snapshots by `snapshot` |

```sh
$ echo -n 'changeit' > RESTIC_PASSWORD
$ kubectl create secret generic local-secret --from-file=./RESTIC_PASSWORD
secret "local-secret" created
```

```yaml
$ kubectl get secret local-secret -o yaml

apiVersion: v1
data:
  RESTIC_PASSWORD: Y2hhbmdlaXQ=
kind: Secret
metadata:
  creationTimestamp: 2017-06-28T12:06:19Z
  name: stash-local
  namespace: default
  resourceVersion: "1440"
  selfLink: /api/v1/namespaces/default/secrets/stash-local
  uid: 31a47380-5bfa-11e7-bb52-08002711f4aa
type: Opaque
```

Now, you can create a Snapshot tpr using this secret. Following parameters are availble for `Local` backend.

| Parameter      | Description                                                                                 |
|----------------|---------------------------------------------------------------------------------------------|
| `local.path`   | `Required`. Path where this volume will be mounted in the sidecar container. Example: /repo |
| `local.volume` | `Required`. Any Kubernetes volume                                                           |

```sh
$ kubectl create -f ./docs/examples/snapshot/local/local-snapshot.yaml
snapshot "local-snapshot" created
```

```yaml
$ kubectl get snapshot local-snapshot -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  creationTimestamp: 2017-06-28T12:14:48Z
  name: local-snapshot
  namespace: default
  resourceVersion: "2000"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/snapshots/local-snapshot
  uid: 617e3487-5bfb-11e7-bb52-08002711f4aa
spec:
  selector:
    matchLabels:
      app: local-snapshot
  fileGroups:
  - path: /source/data
    retentionPolicy:
      keepLast: 5
      prune: true
  backend:
    local:
      path: /repo
      volume:
        emptyDir: {}
        name: repo
    repositorySecretName: local-secret
  schedule: '@every 1m'
  volumeMounts:
  - mountPath: /source/data
    name: source-data
```


### AWS S3
Stash supports AWS S3 service or [Minio](https://minio.io/) servers as backend. To configure this backend, following secret keys are needed:

| Key                     | Description                                                |
|-------------------------|------------------------------------------------------------|
| `AWS_ACCESS_KEY_ID`     | `Required`. AWS / Minio access key ID                      |
| `AWS_SECRET_ACCESS_KEY` | `Required`. AWS / Minio secret access key                  |

```sh
$ echo -n '<your-aws-access-key-id-here>' > AWS_ACCESS_KEY_ID
$ echo -n '<your-aws-secret-access-key-here>' > AWS_SECRET_ACCESS_KEY
$ kubectl create secret generic s3-secret \
    --from-file=./AWS_ACCESS_KEY_ID \
    --from-file=./AWS_SECRET_ACCESS_KEY
secret "s3-secret" created
```

```yaml
$ kubectl get secret s3-secret -o yaml

apiVersion: v1
data:
  AWS_ACCESS_KEY_ID: PHlvdXItYXdzLWFjY2Vzcy1rZXktaWQtaGVyZT4=
  AWS_SECRET_ACCESS_KEY: PHlvdXItYXdzLXNlY3JldC1hY2Nlc3Mta2V5LWhlcmU+
kind: Secret
metadata:
  creationTimestamp: 2017-06-28T12:22:33Z
  name: s3-secret
  namespace: default
  resourceVersion: "2511"
  selfLink: /api/v1/namespaces/default/secrets/s3-secret
  uid: 766d78bf-5bfc-11e7-bb52-08002711f4aa
type: Opaque
```

Now, you can create a Snapshot tpr using this secret. Following parameters are availble for `S3` backend.

| Parameter     | Description                                                                     |
|---------------|---------------------------------------------------------------------------------|
| `s3.endpoint` | `Required`. For S3, use `s3.amazonaws.com`. If your bucket is in a different location, S3 server (s3.amazonaws.com) will redirect snapshot to the correct endpoint. For an S3-compatible server that is not Amazon (like Minio), or is only available via HTTP, you can specify the endpoint like this: `http://server:port`. |
| `s3.bucket`   | `Required`. Name of Bucket                                                      |

```sh
$ kubectl create -f ./docs/examples/snapshot/s3/s3-snapshot.yaml
snapshot "s3-snapshot" created
```

```yaml
$ kubectl get snapshot s3-snapshot -o yaml

apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  creationTimestamp: 2017-06-28T12:58:10Z
  name: s3-snapshot
  namespace: default
  resourceVersion: "4889"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/snapshots/s3-snapshot
  uid: 7036ba69-5c01-11e7-bb52-08002711f4aa
spec:
  selector:
    matchLabels:
      app: s3-snapshot
  fileGroups:
  - path: /source/data
    retentionPolicy:
      keepLast: 5
      prune: true
  backend:
    s3:
      endpoint: 's3.amazonaws.com'
      bucket: stash-qa
      prefix: demo
    repositorySecretName: s3-secret
  schedule: '@every 1m'
  volumeMounts:
  - mountPath: /source/data
    name: source-data
```


### Google Cloud Storage (GCS)
Stash supports Google Cloud Storage(GCS) as backend. To configure this backend, following secret keys are needed:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```sh
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic gcs-secret \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "gcs-secret" created
```

```yaml
$ kubectl get secret gcs-secret -o yaml

apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2017-06-28T13:06:51Z
  name: gcs-secret
  namespace: default
  resourceVersion: "5461"
  selfLink: /api/v1/namespaces/default/secrets/gcs-secret
  uid: a6983b00-5c02-11e7-bb52-08002711f4aa
type: Opaque
```

Now, you can create a Snapshot tpr using this secret. Following parameters are availble for `gcs` backend.

| Parameter      | Description                                                                     |
|----------------|---------------------------------------------------------------------------------|
| `gcs.location` | `Required`. Name of Google Cloud region.                                        |
| `gcs.bucket`   | `Required`. Name of Bucket                                                      |

```sh
$ kubectl create -f ./docs/examples/snapshot/gcs/gcs-snapshot.yaml
snapshot "gcs-snapshot" created
```

```yaml
$ kubectl get snapshot gcs-snapshot -o yaml

apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  creationTimestamp: 2017-06-28T13:11:43Z
  name: gcs-snapshot
  namespace: default
  resourceVersion: "5781"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/snapshots/gcs-snapshot
  uid: 54b1bad3-5c03-11e7-bb52-08002711f4aa
spec:
  selector:
    matchLabels:
      app: gcs-snapshot
  fileGroups:
  - path: /source/data
    retentionPolicy:
      keepLast: 5
      prune: true
  backend:
    gcs:
      location: /repo
      bucket: stash-qa
      prefix: demo
    repositorySecretName: gcs-secret
  schedule: '@every 1m'
  volumeMounts:
  - mountPath: /source/data
    name: source-data
```


### Microsoft Azure Storage
Stash supports Microsoft Azure Storage as backend. To configure this backend, following secret keys are needed:

| Key                     | Description                                                |
|-------------------------|------------------------------------------------------------|
| `AZURE_ACCOUNT_NAME`    | `Required`. Azure Storage account name                     |
| `AZURE_ACCOUNT_KEY`     | `Required`. Azure Storage account key                      |

```sh
$ echo -n '<your-azure-storage-account-name>' > AZURE_ACCOUNT_NAME
$ echo -n '<your-azure-storage-account-key>' > AZURE_ACCOUNT_KEY
$ kubectl create secret generic azure-secret \
    --from-file=./AZURE_ACCOUNT_NAME \
    --from-file=./AZURE_ACCOUNT_KEY
secret "azure-secret" created
```

```yaml
$ kubectl get secret azure-secret -o yaml

apiVersion: v1
data:
  AZURE_ACCOUNT_KEY: PHlvdXItYXp1cmUtc3RvcmFnZS1hY2NvdW50LWtleT4=
  AZURE_ACCOUNT_NAME: PHlvdXItYXp1cmUtc3RvcmFnZS1hY2NvdW50LW5hbWU+
kind: Secret
metadata:
  creationTimestamp: 2017-06-28T13:27:16Z
  name: azure-secret
  namespace: default
  resourceVersion: "6809"
  selfLink: /api/v1/namespaces/default/secrets/azure-secret
  uid: 80f658d1-5c05-11e7-bb52-08002711f4aa
type: Opaque
```

Now, you can create a Snapshot tpr using this secret. Following parameters are availble for `Azure` backend.

| Parameter     | Description                                                                     |
|---------------|---------------------------------------------------------------------------------|
| `azure.container` | `Required`. Name of Storage container                                       |

```sh
$ kubectl create -f ./docs/examples/snapshot/azure/azure-snapshot.yaml
snapshot "azure-snapshot" created
```

```yaml
$ kubectl get snapshot azure-snapshot -o yaml

apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  creationTimestamp: 2017-06-28T13:31:14Z
  name: azure-snapshot
  namespace: default
  resourceVersion: "7070"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/snapshots/azure-snapshot
  uid: 0e8eb89b-5c06-11e7-bb52-08002711f4aa
spec:
  selector:
    matchLabels:
      app: azure-snapshot
  fileGroups:
  - path: /source/data
    retentionPolicy:
      keepLast: 5
      prune: true
  backend:
    azure:
      container: stashqa
      prefix: demo
    repositorySecretName: azure-secret
  schedule: '@every 1m'
  volumeMounts:
  - mountPath: /source/data
    name: source-data
```














## Schedule Backup

**T**o schedule backup, we need to add following BackupScheduleSpec in `spec`

```yaml
spec:
  backupSchedule:
    cronExpression: "@every 6h"
    bucketName: "bucket-for-snapshot"
    storageSecret:
      secretName: "secret-for-bucket"
```

> **Note:** storage can also be used here

When database TPR object is running,
operator immediately takes a backup to validate this information.

And after successful backup, operator will set a scheduler to take backup `every 6h`.

See backup process in [details](backup.md).
