---
title: Postgres
menu:
  docs_0.9.0-beta.0:
    identifier: postgres-db
    name: Postgres
    parent: databases
    weight: 30
menu_name: docs_0.9.0-beta.0
section_menu_id: concepts
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Postgres

## What is Postgres

`Postgres` is a Kubernetes `Custom Resource Definitions` (CRD). It provides declarative configuration for [PostgreSQL](https://www.postgresql.org/) in a Kubernetes native way. You only need to describe the desired database configuration in a Postgres object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## Postgres Spec

As with all other Kubernetes objects, a Postgres needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section.

Below is an example Postgres object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: p1
  namespace: demo
spec:
  version: "9.6-v1"
  replicas: 2
  standbyMode: Hot
  streamingMode: asynchronous
  archiver:
    storage:
      storageSecretName: s3-secret
      s3:
        bucket: kubedb
  databaseSecret:
    secretName: p1-auth
  storageType: "Durable"
  storage:
    storageClassName: standard
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    scriptSource:
      configMap:
        name: pg-init-script
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb
      prefix: demo
  monitor:
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
  configSource:
    configMap:
      name: pg-custom-config
  podTemplate:
    annotation:
      passMe: ToDatabasePod
    controller:
      annotation:
        passMe: ToStatefulSet
    spec:
      schedulerName: my-scheduler
      nodeSelector:
        disktype: ssd
      imagePullSecrets:
      - name: myregistrykey
      env:
      - name: POSTGRES_DB
        value: pgdb
      resources:
        requests:
          memory: "64Mi"
          cpu: "250m"
        limits:
          memory: "128Mi"
          cpu: "500m"
  serviceTemplate:
    annotation:
      passMe: ToService
    spec:
      type: NodePort
      ports:
      - name:  http
        port:  5432
        targetPort: http
  updateStrategy:
    type: RollingUpdate
  terminationPolicy: "DoNotTerminate"
```

### spec.version

`spec.version` is a required field that specifies the name of the [PostgresVersion](/docs/concepts/catalog/postgres.md) crd where the docker images are specified. Currently, when you install KubeDB, it creates the following `PostgresVersion` crd,

 - `9.6.7-v1`, `9.6.7`, `9.6-v1`, `9.6`
 - `10.2-v1`, `10.2`

### spec.replicas

`spec.replicas` specifies the total number of primary and standby nodes in Postgres database cluster configuration. One pod is selected as Primary and others are acted as standby replicas. To know more about how to setup a HA PostgreSQL cluster in KubeDB, please visit [here](/docs/guides/postgres/clustering/ha_cluster.md).

### spec.standbyMode

`spec.standby` is an optional field that specifies the standby mode (_Warm / Hot_) to use for standby replicas. In **hot standby** mode, standby replicas can accept connection and run read-only queries. In **warm standby** mode, standby replicas can't accept connection and only used for replication purpose.

### spec.streamingMode

`spec.streamingMode` is an optional field that specifies the streaming mode (_synchronous / asynchronous_) of the standby replicas. KubeDB currently supports only **asynchronous** streaming mode.

### spec.archiver

`spec.archiver` is an optional field which specifies storage information that will be used by `wal-g`.

 - `spec.archiver.storage.storageSecretName` points to the Secret containing the credentials for cloud storage destination.
 - `spec.archiver.storage.s3.bucket` points to the bucket name used to store continuous archiving data.

Continuous archiving data will be stored in a folder called `{bucket}/{prefix}/kubedb/{namespace}/{postgres-name}/archive/`.

To know more about how to configure Postgres to archive WAL data continuously in AWS S3 bucket, please visit [here](/docs/guides/postgres/snapshot/continuous_archiving.md).

### spec.databaseSecret

`spec.databaseSecret` is an optional field that points to a Secret used to hold credentials for `postgres` database. If not set, KubeDB operator creates a new Secret with name `{postgres-name}-auth` that hold *username* and *password* for `postgres` database.

If you want to use an existing or custom secret, please specify that when creating the Postgres object using `spec.databaseSecret.secretName`. This Secret should contain superuser *username* as `POSTGRES_USER` key and superuser *password* as `POSTGRES_PASSWORD` key.

Example:

```console
$ kubectl create secret generic p1-auth -n demo \
--from-literal=POSTGRES_USER=not@user \
--from-literal=POSTGRES_PASSWORD=not@secret
secret "p1-auth" created
```

```console
$ kubectl get secret -n demo p1-auth -o yaml
apiVersion: v1
data:
  POSTGRES_PASSWORD: bm90QHNlY3JldA==
  POSTGRES_USER: bm90QHVzZXI=
kind: Secret
metadata:
  creationTimestamp: 2018-09-03T11:25:39Z
  name: p1-auth
  namespace: demo
  resourceVersion: "1677"
  selfLink: /api/v1/namespaces/demo/secrets/p1-auth
  uid: 15b3e8a1-af6c-11e8-996d-0800270d7bae
type: Opaque
```

### spec.storageType

`spec.storageType` is an optional field that specifies the type of storage to use for database. It can be either `Durable` or `Ephemeral`. The default value of this field is `Durable`. If `Ephemeral` is used then KubeDB will create Postgres database using [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) volume. In this case, you don't have to specify `spec.storage` field.

### spec.storage

If you don't set `spec.storageType:`  to `Ephemeral` then `spec.storage` field is required. This field specifies the StorageClass of PVCs dynamically allocated to store data for the database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests.

 - `spec.storage.storageClassName` is the name of the StorageClass used to provision PVCs. PVCs don’t necessarily have to request a class. A PVC with its storageClassName set equal to "" is always interpreted to be requesting a PV with no class, so it can only be bound to PVs with no class (no annotation or one set equal to ""). A PVC with no storageClassName is not quite the same and is treated differently by the cluster depending on whether the DefaultStorageClass admission plugin is turned on.
 - `spec.storage.accessModes` uses the same conventions as Kubernetes PVCs when requesting storage with specific access modes.
 - `spec.storage.resources` can be used to request specific quantities of storage. This follows the same resource model used by PVCs.

To learn how to configure `spec.storage`, please visit the links below:
 - https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims

### spec.init

`spec.init` is an optional section that can be used to initialize a newly created Postgres database. PostgreSQL databases can be initialized from these three ways:

  1. Initialize from Script
  2. Initialize from Snapshot
  3. Initialize from WAL archive
  
#### Initialize via Script

To initialize a PostgreSQL database using a script (shell script, db migrator, etc.), set the `spec.init.scriptSource` section when creating a Postgres object. `scriptSource` must have the following information:

- [VolumeSource](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes): Where your script is loaded from.

Below is an example showing how a script from a configMap can be used to initialize a PostgreSQL database.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: postgres-db
spec:
  version: "9.6-v1"
  init:
    scriptSource:
      configMap:
        name: pg-init-script
```

In the above example, Postgres will execute provided script once the database is running. For more details tutorial on how to initialize from script, please visit [here](/docs/guides/postgres/initialization/script_source.md).

#### Initialize from Snapshots

To initialize from prior Snapshot, set the `spec.init.snapshotSource` section when creating a Postgres object. In this case, `snapshotSource` must have the following information:

- `name:` Name of the Snapshot
- `namespace:` Namespace of the Snapshot

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: postgres-db
spec:
  version: "9.6-v1"
  databaseSecret:
    secretName: postgres-old-auth
  init:
    snapshotSource:
      name: "snapshot-xyz"
      namespace: "demo"
```

In the above example, PostgreSQL database will be initialized from Snapshot `snapshot-xyz` of `demo` namespace. Here, KubeDB operator will launch a Job to initialize PostgreSQL once StatefulSet pods are running.

When initializing from Snapshot, superuser credentials must have to match with the previous one. For example, let's say, Snapshot `snapshot-xyz` is for Postgres `postgres-old`. In this case, new Postgres `postgres-db` should use the same credentials for superuser of `postgres-old`. Otherwise, the restoration process will fail.

For more details tutorial on how to initialize from snapshot, please visit [here](/docs/guides/postgres/initialization/snapshot_source.md).

#### Initialize from WAL archive

To initialize from WAL archive, set the `spec.init.postgresWAL` section when creating a Postgres object.

Below is an example showing how to initialize a PostgreSQL database from WAL archive.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: postgres-db
spec:
  version: "9.6-v1"
  databaseSecret:
    secretName: postgres-old
  init:
    postgresWAL:
      storageSecretName: s3-secret
      s3:
        endpoint: 's3.amazonaws.com'
        bucket: kubedb
        prefix: 'kubedb/demo/old-pg/archive'
```

In the above example, PostgreSQL database will be initialized from WAL archive.

When initializing from WAL archive, superuser credentials must have to match with the previous one. For example, let's say, we want to initialize this database from `postgres-old` WAL archive. In this case, superuser credentials of new Postgres should be the same as `postgres-old`. Otherwise, the restoration process will be failed.

For more details tutorial on how to initialize from wal archive, please visit [here](/docs/guides/postgres/initialization/wal_source.md).

### spec.backupSchedule

KubeDB supports taking periodic snapshots for Postgres database. This is an optional section in `.spec`. When `spec.backupSchedule` section is added, KubeDB operator immediately takes a backup to validate this information. After that, at each tick kubeDB operator creates a [Snapshot](/docs/concepts/snapshot.md) object. This triggers operator to create a Job to take backup.

You have to specify following fields to take periodic backup of your Postgres database:

 - `spec.backupSchedule.cronExpression` is a required [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). This specifies the schedule for backup operations.
 - `spec.backupSchedule.{storage}` is a required field that is used as the destination for storing snapshot data. KubeDB supports cloud storage providers like S3, GCS, Azure, and OpenStack Swift. It also supports any locally mounted Kubernetes volumes, like NFS, Ceph, etc. Only one backend can be used at a time. To learn how to configure this, please visit [here](/docs/concepts/snapshot.md).

You can also specify a template for pod of backup job through `spec.backupSchedule.podTemplate`. KubeDB will use the information you have provided in `podTemplate` to create the backup job. KubeDB accept following fields to set in `spec.backupSchedule.podTemplate`:

- metadata
  - annotations (pod's annotation)
- controller
  - annotations (job's annotation)
- spec:
  - args
  - env
  - resources
  - imagePullSecrets
  - initContainers
  - nodeSelector
  - affinity
  - schedulerName
  - tolerations
  - priorityClassName
  - priority
  - securityContext
  - livenessProbe
  - readinessProbe
  - lifecycle

### spec.monitor

PostgreSQL managed by KubeDB can be monitored with builtin-Prometheus and CoreOS-Prometheus operator out-of-the-box. To learn more,

- [Monitor PostgreSQL with builtin Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md)
- [Monitor PostgreSQL with CoreOS Prometheus operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md)

### spec.configSource

`spec.configSource` is an optional field that allows users to provide custom configuration for PostgreSQL. This field accepts a [`VolumeSource`](https://github.com/kubernetes/api/blob/release-1.11/core/v1/types.go#L47). You can use any kubernetes supported volume source such as `configMap`, `secret`, `azureDisk` etc. To learn more about how to use a custom configuration file see [here](/docs/guides/postgres/custom-config/using-custom-config.md).

### spec.podTemplate

KubeDB allows providing a template for database pod through `spec.podTemplate`. KubeDB operator will pass the information provided in `spec.podTemplate` to the StatefulSet created for Postgres database.

KubeDB accept following fields to set in `spec.podTemplate:`

- metadata
  - annotations (pod's annotation)
- controller
  - annotations (statefulset's annotation)
- spec:
  - env
  - resources
  - initContainers
  - imagePullSecrets
  - nodeSelector
  - affinity
  - schedulerName
  - tolerations
  - priorityClassName
  - priority
  - securityContext
  - livenessProbe
  - readinessProbe
  - lifecycle

Uses of some field of `spec.podTemplate` is described below,

#### spec.podTemplate.spec.env

`spec.podTemplate.spec.env` is an optional field that specifies the environment variables to pass to the Postgres docker image. To know about supported environment variables, please visit [here](https://hub.docker.com/_/postgres/).

Note that, KubeDB does not allow `POSTGRES_USER` and `POSTGRES_PASSWORD` environment variable to set in `spec.podTemplate.spec.env`. If you want to set the superuser *username* and *password*, please use `spec.databaseSecret` instead described earlier.

If you try to set `POSTGRES_USER` or `POSTGRES_PASSWORD` environment variable in Postgres crd, KubeDB operator will reject the request with following error,

```ini
Error from server (Forbidden): error when creating "./postgres.yaml": admission webhook "postgres.validators.kubedb.com" denied the request: environment variable POSTGRES_PASSWORD is forbidden to use in Postgres spec
```

Also, note that KubeDB does not allow to update the environment variables as updating them does not have any effect once the database is created.  If you try to update environment variables, KubeDB operator will reject the request with following error,

```ini
Error from server (BadRequest): error when applying patch:
...
for: "./postgres.yaml": admission webhook "postgres.validators.kubedb.com" denied the request: precondition failed for:
...
At least one of the following was changed:
    apiVersion
    kind
    name
    namespace
    spec.standby
    spec.streaming
    spec.archiver
    spec.databaseSecret
    spec.storageType
    spec.storage
    spec.podTemplate.spec.nodeSelector
    spec.init
    spec.podTemplate.spec.env
```

#### spec.podTemplate.spec.imagePullSecrets

`spec.podTemplate.spec.imagePullSecrets` is an optional field that points to secrets to be used for pulling docker image if you are using a private docker registry. For more details on how to use private docker registry, please visit [here](/docs/guides/postgres/private-registry/using-private-registry.md).

#### spec.podTemplate.spec.nodeSelector

`spec.podTemplate.spec.nodeSelector` is an optional field that specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). To learn more, see [here](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) .

#### spec.podTemplate.spec.resources

`spec.podTemplate.spec.resources` is an optional field. This can be used to request compute resources required by the database pods. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).

### spec.serviceTemplate

You can also provide a template for the services created by KubeDB operator for Postgres database through `spec.serviceTemplate`. This will allow you to set the type and other properties of the services.

KubeDB allows following fields to set in `spec.serviceTemplate`:

- annotations
- type
- ports
- clusterIP
- externalIPs
- loadBalancerIP
- loadBalancerSourceRanges
- externalTrafficPolicy
- healthCheckNodePort

### spec.updateStrategy

You can specify [update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#update-strategies) of StatefulSet created by KubeDB for Postgres database thorough `spec.updateStrategy` field. The default value of this field is `RollingUpdate`. In future, we will use this field to determine how automatic migration from old KubeDB version to new one should behave.

### spec.terminationPolicy

`terminationPolicy` gives flexibility whether to `nullify`(reject) the delete operation of `Postgres` crd or which resources KubeDB should keep or delete when you delete `Postgres` crd. KubeDB provides following four termination policies:

- DoNotTerminate
- Pause (`Default`)
- Delete
- WipeOut

When, `terminationPolicy` is `DoNotTerminate`, KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to provide safety from accidental deletion of database. If admission webhook is enabled, KubeDB prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`.

Following table show what KubeDB does when you delete Postgres crd for different termination policies,

|                Behaviour                 | DoNotTerminate |  Pause   |  Delete  | WipeOut  |
| ---------------------------------------- | :------------: | :------: | :------: | :------: |
| 1. Nullify Delete operation              |    &#10003;    | &#10007; | &#10007; | &#10007; |
| 2. Create Dormant Database               |    &#10007;    | &#10003; | &#10007; | &#10007; |
| 3. Delete StatefulSet                    |    &#10007;    | &#10003; | &#10003; | &#10003; |
| 4. Delete Services                       |    &#10007;    | &#10003; | &#10003; | &#10003; |
| 5. Delete PVCs                           |    &#10007;    | &#10007; | &#10003; | &#10003; |
| 6. Delete Secrets                        |    &#10007;    | &#10007; | &#10007; | &#10003; |
| 7. Delete Snapshots                      |    &#10007;    | &#10007; | &#10007; | &#10003; |
| 8. Delete Snapshot data from bucket      |    &#10007;    | &#10007; | &#10007; | &#10003; |
| 9. Delete archieved WAL data from bucket |    &#10007;    | &#10007; | &#10007; | &#10003; |

If you don't specify `spec.terminationPolicy` KubeDB uses `Pause` termination policy by default.

## Next Steps

- Learn how to use KubeDB to run a PostgreSQL database [here](/docs/guides/postgres/README.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
