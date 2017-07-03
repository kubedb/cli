> New to KubeDB? Please start [here](/docs/tutorial.md).

# Elastics

## What is Elastic
A `Elastic` is a Kubernetes `Third Party Object` (TPR). It provides declarative configuration for [Elasticsearch](https://www.elastic.co/products/elasticsearch) in a Kubernetes native way. You only need to describe the desired database configuration in a Elastic object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## Elastic Spec
As with all other Kubernetes objects, a Elastic needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example Elastic object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elastic
metadata:
  name: elasticsearch-db
spec:
  version: 2.3.1
  replicas: 1
  storage:
    class: "gp2"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: "10Gi"
```

Save this yaml as `elasticsearch-db.yaml` and create Elastic object.

```sh
kubedb create -f  ./docs/examples/elastic/elastic-with-storage.yaml

elastic "elasticsearch-db" created
```

**O**ur deployed unified operator will detect this object and will create workloads.

For this object, following kubernetes objects will be created in same namespace:
* StatefulSet (name: **elasticsearch-db**-es)
* Service (name: **elasticsearch-db**)
* GoverningService (If not available) (name: **kubedb**)


**N**ow lets see whether our database is ready or not.

```bash
$ kubedb get elastic elasticsearch-db -o wide

NAME               VERSION   STATUS    AGE
elasticsearch-db   2.3.1     Running   37m
```

This database do not have any PersistentVolume behind StatefulSet.

#### Add storage support


**W**e can create a Elastic database that will use PersistentVolumeClaim in StatefulSet.


**T**o add PersistentVolume support, we need to add following StorageSpec in `spec`

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elastic
metadata:
  name: elasticsearch-db
spec:
  version: 2.3.1
  replicas: 1
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

**A**s we have used storage information in our database yaml, StatefulSet will be created with PersistentVolumeClaim.

Following command will list `pvc` for this database.

```bash
$ kubectl get pvc --selector='kubedb.com/kind=Elastic,kubedb.com/name=elasticsearch-db'

NAME                         STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
data-elasticsearch-db-pg-0   Bound     pvc-a1a95954-4a75-11e7-8b69-12f236046fba   10Gi       RWO           2m
```




### Initialize Database

We now support initialization from two sources.

2. SnapshotSource

We can use one of them to initialize out database.

#### SnapshotSource

**D**atabase can also be initialized with Snapshot data.

In this case, SnapshotSource must have following information:
1. `namespace:` Namespace of Snapshot object
2. `name:` Name of the Snapshot

If SnapshotSource is provided to initialize database,
a job will do that initialization when database is running.

##### Example

```yaml
spec:
  init:
    snapshotSource:
      name: "snapshot-xyz"
```

Database will be initialized from backup data of Snapshot `snapshot-xyz` in `default` namespace.



















#### Schedule Backup

**W**e can also schedule automatic backup by providing BackupSchedule information in `spec.backupSchedule`.

How to add information in Elastic `spec` to schedule automatic backup? See [here](../schedule-backup.md).


#### Monitor Database

**W**e can also monitor our elasticsearch database.
To enable monitoring, we need to set MonitorSpec in Elastic `spec`.

How to set monitoring? See [here](../monitor-database.md).
