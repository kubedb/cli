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
    class: "standard"
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
`spec.version` is the version of PostgreSQL database. In this tutorial, a PostgreSQL 9.5 database is going to be created.

### spec.storage
`spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.


Class
A claim can request a particular class by specifying the name of a StorageClass using the attribute storageClassName. Only PVs of the requested class, ones with the same storageClassName as the PVC, can be bound to the PVC.
PVCs don’t necessarily have to request a class. A PVC with its storageClassName set equal to "" is always interpreted to be requesting a PV with no class, so it can only be bound to PVs with no class (no annotation or one set equal to ""). A PVC with no storageClassName is not quite the same and is treated differently by the cluster depending on whether the DefaultStorageClass admission plugin is turned on.

Access Modes
Claims use the same conventions as volumes when requesting storage with specific access modes.
Resources
Claims, like pods, can request specific quantities of a resource. In this case, the request is for storage. The same resource model applies to both volumes and claims.




https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims

https://kubernetes.io/docs/tutorials/stateful-application/basic-stateful-set/#writing-to-stable-storage






To use PersistentVolume, add the `spec.storage` section when creating Postgres object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: postgres-db
spec:
  version: 9.5
  storage:
    class: "gp2"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: "10Gi"
```

Here we must have to add following storage information in `spec.storage`:

* `class:` StorageClass (`kubectl get storageclasses`)
* `resources:` ResourceRequirements for PersistentVolumeClaimSpec

As `spec.storage` fields are set, StatefulSet will be created with dynamically provisioned PersistentVolumeClaim. Following command will list PVCs for this database.

```console
$ kubectl get pvc --selector='kubedb.com/kind=Postgres,kubedb.com/name=postgres-db'

NAME                 STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
data-postgres-db-0   Bound     pvc-a1a95954-4a75-11e7-8b69-12f236046fba   10Gi       RWO           2m
```





### spec.databaseSecret


Please note that KubeDB operator has created a new Secret called `p1-admin-auth` (format: {tpr-name}-admin-auth) for storing the password for `postgres` superuser. This secret contains a `.admin` key with a ini formatted key-value pairs. If you want to use an existing secret please specify that when creating the tpr using `spec.databaseSecret.secretName`.

Now, you can connect to this database from the PGAdmin dasboard using the database pod IP and `postgres` user password. 

```console
$ kubectl get pods p1-0 -n demo -o yaml | grep IP
  hostIP: 192.168.99.100
  podIP: 172.17.0.6

$ kubectl get secrets -n demo p1-admin-auth -o jsonpath='{.data.\.admin}' | base64 -d
POSTGRES_PASSWORD=R9keKKRTqSJUPtNC
```

The `.admin` contains a `ini` formatted key/value pairs. 

```ini
POSTGRES_PASSWORD=vPlT2PzewCaC3XZP
```





### spec.nodeSelector

nodeSelector is the simplest form of constraint. nodeSelector is a field of PodSpec. It specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). The most common usage is one key-value pair


https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector






### spec.init
`spec.resources` refers to compute resources required by the `stash` sidecar container. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).




### Database Initialization
PostgreSQL databases can be initialized in one of two ways:

#### Initialize via Script
To initialize a PostgreSQL database using a script (shell script, db migrator, etc.), set the `spec.init.scriptSource` section when creating a Postgres object. ScriptSource must have following information:
1. `scriptPath:` ScriptPath (The script you want to run)
2. [VolumeSource](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes) (Where your script and other data is stored)

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

> **Note:** all path used in script should be relative

#### Initialize from Snapshots
To initialize from prior snapshot, set the `spec.init.snapshotSource` section when creating a Postgres object.

In this case, SnapshotSource must have following information:
1. `namespace:` Namespace of Snapshot object
2. `name:` Name of the Snapshot

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

In the above example, PostgreSQL database will be initialized from Snapshot `snapshot-xyz` in `default` namespace. Here,  KubeDB operator will launch a Job to initialize PostgreSQL once StatefulSet pods are running.














### spec.backupSchedule
`spec.resources` refers to compute resources required by the `stash` sidecar container. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).

KubeDB supports taking periodic backups for a database using a [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). To take periodic backups, edit the Postgres tpr to add `spec.backupSchedule` section.


## Schedule Backups
Scheduled backups are supported for all types of databases. To schedule backups, add the following `BackupScheduleSpec` in `spec` of a database tpr.
All snapshot storage backends are supported for scheduled backup.

```yaml
spec:
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: "secret-for-bucket"
    s3:
      endpoint: 's3.amazonaws.com'
      bucket: kubedb-qa
```

`spec.backupSchedule.schedule` is a [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26) that indicates how often backups are taken.

When `spec.backupSchedule` section is added, KubeDB operator immediately takes a backup to validate this information. After that, at each tick kubeDB operator creates a Snapshot object. This triggers operator to create a Job to take backup.








### spec.doNotPause
`spec.doNotPause` tells KubeDB operator that if this tpr is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.


### spec.monitor
To learn how to monitor Postgres databases, please visit [here](/docs/concepts/monitoring.md).


### spec.resources
`spec.resources` refers to compute resources required by the `stash` sidecar container. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).


## Next Steps
- Learn how to use KubeDB to run a PostgreSQL database [here](/docs/tutorials/postgres.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Thinking about monitoring your database? KubeDB works [out-of-the-box with Prometheus](/docs/tutorials/monitoring.md).
- Learn how to use KubeDB in a [RBAC](/docs/tutorials/rbac.md) enabled cluster.
- Wondering what features are coming next? Please visit [here](/ROADMAP.md). 
- Want to hack on KubeDB? Check our [contribution guidelines](/CONTRIBUTING.md).
