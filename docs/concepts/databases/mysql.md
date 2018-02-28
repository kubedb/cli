---
title: MySQL
menu:
  docs_0.8.0-beta.2:
    identifier: mysql-db
    name: MySQL
    parent: databases
    weight: 25
menu_name: docs_0.8.0-beta.2
section_menu_id: concepts
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# MySQL

## What is MySQL

`MySQL` is a Kubernetes `Custom Resource Definitions` (CRD). It provides declarative configuration for [MySQL](https://www.mysql.com/) in a Kubernetes native way. You only need to describe the desired database configuration in a MySQL object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## MySQL Spec

As with all other Kubernetes objects, a MySQL needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example MySQL object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: m1
  namespace: demo
spec:
  version: 8.0
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  databaseSecret:
    secretName: m1-auth
  nodeSelector:
    disktype: ssd
  init:
    scriptSource:
      gitRepo:
        repository: "https://github.com/kubedb/mysql-init-scripts.git"
        directory: .
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: ms-snap-secret
    gcs:
      bucket: restic
      prefix: demo
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
  doNotPause: true
  monitor:
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
  resources:
    requests:
      memory: "64Mi"
      cpu: "250m"
    limits:
      memory: "128Mi"
      cpu: "500m"
```

### spec.version

`spec.version` is a required field specifying the version of MySQL database. Official [mysql docker images](https://hub.docker.com/r/library/mysql/tags/) will be used for the specific version.

### spec.storage

`spec.storage` is an optional field that specifies the StorageClass of PVCs dynamically allocated to store data for the database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

- `spec.storage.storageClassName` is the name of the StorageClass used to provision PVCs. PVCs don’t necessarily have to request a class. A PVC with its storageClassName set equal to "" is always interpreted to be requesting a PV with no class, so it can only be bound to PVs with no class (no annotation or one set equal to ""). A PVC with no storageClassName is not quite the same and is treated differently by the cluster depending on whether the DefaultStorageClass admission plugin is turned on.

- `spec.storage.accessModes` uses the same conventions as Kubernetes PVCs when requesting storage with specific access modes.

- `spec.storage.resources` can be used to request specific quantities of storage. This follows the same resource model used by PVCs.

To learn how to configure `spec.storage`, please visit the links below:

- https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims

### spec.databaseSecret

`spec.databaseSecret` is an optional field that points to a Secret used to hold credentials for `mysql` super user. If not set, KubeDB operator creates a new Secret `{mysql-object-name}-auth` for storing the password for `mysql` superuser for each MySQL object. If you want to use an existing secret please specify that when creating the MySQL object using `spec.databaseSecret.secretName`.

This secret contains a `user` key and a `password` key which contains the `username` and `password` respectively for `mysql` root user. Here the value of `user` key is fixed to be `root`. Example:

```ini
apiVersion: v1
data:
  password: NnE4dV8yak1PVy1PT1pYaw==
  user: cm9vdA==
kind: Secret
metadata:
  ...
  name: m1-auth
  namespace: demo
  ...
type: Opaque
```

### spec.nodeSelector

`spec.nodeSelector` is an optional field that specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). To learn more, see [here](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) .

### spec.init

`spec.init` is an optional section that can be used to initialize a newly created MySQL database. MySQL databases can be initialized in one of two ways:

#### Initialize via Script

To initialize a MySQL database using a script (shell script, sql script etc.), set the `spec.init.scriptSource` section when creating a MySQL object. It will execute files alphabetically with extensions `.sh` , `.sql`  and `.sql.gz` that are found in the repository. ScriptSource must have following information:

- [VolumeSource](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes): Where your script is loaded from.

Below is an example showing how a shell script from a git repository can be used to initialize a MySQL database.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: m1
spec:
  version: 8.0
  init:
    scriptSource:
      gitRepo:
        repository: "https://github.com/kubedb/mysql-init-scripts.git"
        directory: .
```

In the above example, KubeDB operator will launch a Job to execute all sql script of `mysql-init-script` repo in alphabetical  order once StatefulSet pods are running.

#### Initialize from Snapshots

To initialize from prior snapshots, set the `spec.init.snapshotSource` section when creating a MySQL object. In this case, SnapshotSource must have following information:

- `name:` Name of the Snapshot

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: m1
spec:
  version: 8.0
  init:
    snapshotSource:
      name: "snapshot-xyz"
```

In the above example, MySQL database will be initialized from Snapshot `snapshot-xyz` in `default` namespace. Here, KubeDB operator will launch a Job to initialize MySQL once StatefulSet pods are running.

### spec.backupSchedule

KubeDB supports taking periodic snapshots for MySQL database. This is an optional section in `.spec`. When `spec.backupSchedule` section is added, KubeDB operator immediately takes a backup to validate this information. After that, at each tick kubeDB operator creates a [Snapshot](/docs/concepts/snapshot.md) object. This triggers operator to create a Job to take backup. If used, set the various sub-fields accordingly.

- `spec.backupSchedule.cronExpression` is a required [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). This specifies the schedule for backup operations.

- `spec.backupSchedule.{storage}` is a required field that is used as the destination for storing snapshot data. KubeDB supports cloud storage providers like S3, GCS, Azure and OpenStack Swift. It also supports any locally mounted Kubernetes volumes, like NFS, Ceph , etc. Only one backend can be used at a time. To learn how to configure this, please visit [here](/docs/concepts/snapshot.md).

- `spec.backupSchedule.resources` is an optional field that can request compute resources required by Jobs used to take snapshot or initialize databases from snapshot.  To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).

### spec.doNotPause

`spec.doNotPause` is an optional field that tells KubeDB operator that if this MySQL object is deleted, whether it should be reverted automatically. This should be set to `true` for production databases to avoid accidental deletion. If not set or set to false, deleting a MySQL object put the database into a dormant state. THe StatefulSet for a DormantDatabase is deleted but the underlying PVCs are left intact. This allows user to resume the database later.

### spec.imagePullSecret

`KubeDB` provides the flexibility of deploying MySQL database from a private Docker registry. To learn how to deploym MySQL from a private registry, please visit [here](/docs/guides/mysql/private-registry/using-private-registry.md).

### spec.monitor

MySQL managed by KubeDB can be monitored with builtin-Prometheus and CoreOS-Prometheus operator out-of-the-box. To learn more,

- [Monitor MySQL with builtin Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md)
- [Monitor MySQL with CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md)

### spec.resources

`spec.resources` is an optional field. This can be used to request compute resources required by the database pods. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).

## Next Steps

- Learn how to use KubeDB to run a MySQL database [here](/docs/guides/mysql/README.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
