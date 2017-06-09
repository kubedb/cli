# Take Backup

We need to create a Snapshot object to initiate backup process. 

Here is a template of Snapshot object

```yaml
apiVersion: kubedb.com/v1beta1
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

**Google Cloud**

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

**Amazon S3**

* `provider: s3`
* `config:`
    ```yaml
    access_key_id: "access_key_id"
    region: "region"
    secret_key: "secret_key"
    ```

**Azure**

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

NAME           STATUS      BUCKET                AGE
snapshot-xyz   Succeeded   bucket-for-snapshot   5m
```
