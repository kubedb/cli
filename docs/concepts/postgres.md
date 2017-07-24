> New to KubeDB? Please start [here](/docs/tutorials/README.md).

# Postgres

## What is Postgres
A `Postgres` is a Kubernetes `Third Party Object` (TPR). It provides declarative configuration for [PostgreSQL](https://www.postgresql.org/) in a Kubernetes native way. You only need to describe the desired database configuration in a Postgres object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## Postgres Spec
As with all other Kubernetes objects, a Postgres needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example Postgres object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: p1
  namespace: demo
spec:
  version: 9.5
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  databaseSecret:
    secretName: p1-admin-auth
  nodeSelector:
    disktype: ssd
  init:
    scriptSource:
      scriptPath: "postgres-init-scripts/run.sh"
      gitRepo:
        repository: "https://github.com/k8sdb/postgres-init-scripts.git"
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: snap-secret
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
    agent: coreos-prometheus-operator
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
`spec.version` is a required field specifying the version of PostgreSQL database. Currently the supported value is `9.5`.


### spec.storage
`spec.storage` is an optional field that specifies the StorageClass of PVCs dynamically allocated to store data for the database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

 - `spec.storage.storageClassName` is the name of the StorageClass used to provision PVCs. PVCs don’t necessarily have to request a class. A PVC with its storageClassName set equal to "" is always interpreted to be requesting a PV with no class, so it can only be bound to PVs with no class (no annotation or one set equal to ""). A PVC with no storageClassName is not quite the same and is treated differently by the cluster depending on whether the DefaultStorageClass admission plugin is turned on.

 - `spec.storage.accessModes` uses the same conventions as Kubernetes PVCs when requesting storage with specific access modes.

 - `spec.storage.resources` can be used to request specific quantities of storage. This follows the same resource model used by PVCs.

To learn how to configure `spec.storage`, please visit the links below:
 - https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims


### spec.databaseSecret
`spec.databaseSecret` is an optional field that points to a Secret used to hold credentials for `postgres` super user. If not set, KubeDB operator creates a new Secret `{tpr-name}-admin-auth` for storing the password for `postgres` superuser for each Postgres object. If you want to use an existing secret please specify that when creating the tpr using `spec.databaseSecret.secretName`.

This secret contains a `.admin` key with a ini formatted key-value pairs. Example:
```ini
POSTGRES_PASSWORD=vPlT2PzewCaC3XZP
```


### spec.nodeSelector
`spec.nodeSelector` is an optional field that specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). To learn more, see [here](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) .


### spec.init
`spec.init` is an optional section that can be used to initialize a newly created Postgres database. PostgreSQL databases can be initialized in one of two ways:

#### Initialize via Script
To initialize a PostgreSQL database using a script (shell script, db migrator, etc.), set the `spec.init.scriptSource` section when creating a Postgres object. ScriptSource must have following information:

 - `scriptPath:` The script you want to run. Note that all path used in script should be relative.
 - [VolumeSource](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes): Where your script is loaded from.

Below is an example showing how a shell script from a git repository can be used to initialize a PostgreSQL database.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: postgres-db
spec:
  version: 9.5
  init:
    scriptSource:
      scriptPath: "postgres-init-scripts/run.sh"
      gitRepo:
        repository: "https://github.com/k8sdb/postgres-init-scripts.git"
```

In the above example, KubeDB operator will launch a Job to execute `run.sh` script once StatefulSet pods are running.


#### Initialize from Snapshots
To initialize from prior snapshots, set the `spec.init.snapshotSource` section when creating a Postgres object. In this case, SnapshotSource must have following information:

 - `name:` Name of the Snapshot

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: postgres-db
spec:
  version: 9.5
  init:
    snapshotSource:
      name: "snapshot-xyz"
```

In the above example, PostgreSQL database will be initialized from Snapshot `snapshot-xyz` in `default` namespace. Here, KubeDB operator will launch a Job to initialize PostgreSQL once StatefulSet pods are running.

### spec.backupSchedule
KubeDB supports taking periodic snapshots for Postgres database. This is an optional section in `.spec`. When `spec.backupSchedule` section is added, KubeDB operator immediately takes a backup to validate this information. After that, at each tick kubeDB operator creates a [Snapshot](/docs/concepts/snapshot.md) object. This triggers operator to create a Job to take backup. If used, set the various sub-fields accordingly.

 - `spec.backupSchedule.cronExpression` is a required [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). This specifies the schedule for backup operations.

 - `spec.backupSchedule.{storage}` is a required field that is used as the destination for storing snapshot data. KubeDB supports cloud storage providers like S3, GCS, Azure and OpenStack Swift. It also supports any locally mounted Kubernetes volumes, like NFS, Ceph , etc. Only one backend can be used at a time. To learn how to configure this, please visit [here](/docs/concepts/snapshots.md).

 - `spec.backupSchedule.resources` is an optional field that can request compute resources required by Jobs used to take snapshot or initialize databases from snapshot.  To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).


### spec.doNotPause
`spec.doNotPause` is an optional field that tells KubeDB operator that if this tpr is deleted, whether it should be reverted automatically. This should be set to `true` for production databases to avoid accidental deletion. If not set or set to false, deleting a Postgres object put the database into a dormant state. THe StatefulSet for a DormantDatabase is deleted but the underlying PVCs are left intact. This allows user to resume the database later.


### spec.monitor
To learn how to monitor Postgres databases, please visit [here](/docs/concepts/monitoring.md).


### spec.resources
`spec.resources` is an optional field. This can be used to request compute resources required by the database pods. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).


## Next Steps
- Learn how to use KubeDB to run a PostgreSQL database [here](/docs/tutorials/postgres.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Thinking about monitoring your database? KubeDB works [out-of-the-box with Prometheus](/docs/tutorials/monitoring.md).
- Learn how to use KubeDB in a [RBAC](/docs/tutorials/rbac.md) enabled cluster.
- Wondering what features are coming next? Please visit [here](/ROADMAP.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/CONTRIBUTING.md).
