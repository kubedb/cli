> New to KubeDB? Please start [here](/docs/tutorials/README.md).

# Elasticsearch

## What is Elasticsearch
A `Elasticsearch` is a Kubernetes `Third Party Object` (TPR). It provides declarative configuration for [Elasticsearch](https://www.elastic.co/products/elasticsearch) in a Kubernetes native way. You only need to describe the desired database configuration in a Elasticsearch object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## Elasticsearch Spec
As with all other Kubernetes objects, a Elasticsearch needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example Elasticsearch object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: e1
  namespace: demo
spec:
  version: 2.3.1
  replicas: 1
  storage:
    class: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  nodeSelector:
    disktype: ssd
  init:
    snapshotSource:
      name: "snapshot-xyz"
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
  monitor:
    agent: coreos-prometheus-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
  doNotPause: true
  resources:
    requests:
      memory: "64Mi"
      cpu: "250m"
    limits:
      memory: "128Mi"
      cpu: "500m"
```

### spec.version
`spec.version` is a required field specifying the version of Elasticsearch cluster. Currently the supported value is `2.3.1`.


## spec.replicas
`spec.replicas` specifies the number of pods in the Elasticsearch cluster. If not set, this defaults to 1.


### spec.storage
`spec.storage` is an optional field that specifies the StorageClass of PVCs dynamically allocated to store data for the database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

 - `spec.storage.class` is the name of the StorageClass used to provision PVCs. PVCs don’t necessarily have to request a class. A PVC with its storageClassName set equal to "" is always interpreted to be requesting a PV with no class, so it can only be bound to PVs with no class (no annotation or one set equal to ""). A PVC with no storageClassName is not quite the same and is treated differently by the cluster depending on whether the DefaultStorageClass admission plugin is turned on.

 - `spec.storage.accessModes` uses the same conventions as Kubernetes PVCs when requesting storage with specific access modes.

 - `spec.storage.resources` can be used to request specific quantities of storage. This follows the same resource model used by PVCs.

To learn how to configure `spec.storage`, please visit the links below:
 - https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims


### spec.nodeSelector
`spec.nodeSelector` is an optional field that specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). To learn more, see [here](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) .


### spec.init
`spec.init` is an optional section that can be used to initialize a newly created Elasticsearch cluster from prior snapshots. To initialize from prior snapshots, set the `spec.init.snapshotSource` section when creating a Elasticsearch object. In this case, SnapshotSource must have following information:

 - `name:` Name of the Snapshot

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: elasticsearch-db
spec:
  version: 2.3.1
  replicas: 1
  init:
    snapshotSource:
      name: "snapshot-xyz"
```

In the above example, Elasticsearch cluster will be initialized from Snapshot `snapshot-xyz` in `default` namespace. Here, KubeDB operator will launch a Job to initialize Elasticsearch, once StatefulSet pods are running.

### spec.backupSchedule
KubeDB supports taking periodic snapshots for Elasticsearch database. This is an optional section in `.spec`. When `spec.backupSchedule` section is added, KubeDB operator immediately takes a backup to validate this information. After that, at each tick kubeDB operator creates a [Snapshot](/docs/concepts/snapshot.md) object. This triggers operator to create a Job to take backup. If used, set the various sub-fields accordingly.

 - `spec.backupSchedule.cronExpression` is a required [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). This specifies the schedule for backup operations.

 - `spec.backupSchedule.{storage}` is a required field that is used as the destination for storing snapshot data. KubeDB supports cloud storage providers like S3, GCS, Azure and OpenStack Swift. It also supports any locally mounted Kubernetes volumes, like NFS, Ceph , etc. Only one backend can be used at a time. To learn how to configure this, please visit [here](/docs/concepts/snapshots.md).

 - `spec.backupSchedule.resources` is an optional field that can request compute resources required by Jobs used to take snapshot or initialize databases from snapshot.  To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).


### spec.doNotPause
`spec.doNotPause` is an optional field that tells KubeDB operator that if this tpr is deleted, whether it should be reverted automatically. This should be set to `true` for production databases to avoid accidental deletion. If not set or set to false, deleting a Elasticsearch object put the database into a dormant state. THe StatefulSet for a DormantDatabase is deleted but the underlying PVCs are left intact. This allows user to resume the database later.


### spec.monitor
To learn how to monitor Elasticsearch databases, please visit [here](/docs/concepts/monitoring.md).


### spec.resources
`spec.resources` is an optional field. This can be used to request compute resources required by the database pods. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).


## Next Steps
- Learn how to use KubeDB to run an Elasticsearch database [here](/docs/tutorials/elasticsearch.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Thinking about monitoring your database? KubeDB works [out-of-the-box with Prometheus](/docs/tutorials/monitoring.md).
- Learn how to use KubeDB in a [RBAC](/docs/tutorials/rbac.md) enabled cluster.
- Wondering what features are coming next? Please visit [here](/ROADMAP.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/CONTRIBUTING.md).
