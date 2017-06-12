### Create Elastic

**L**ets create a simple elasticsearch database using following yaml.

```yaml
apiVersion: kubedb.com/v1beta1
kind: Elastic
metadata:
  name: elasticsearch-db
spec:
  version: 2.3.1
  replicas: 1
```

Save this yaml as `elasticsearch-db.yaml` and create Elastic object.

```bash
$ cat elasticsearch-db.yaml | kubedb create -f -

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

How to add storage information in Elastic `spec`? See [here](../support-storage.md).

Following command will list `pvc` for this database.

```bash
$ kubectl get pvc --selector='kubedb.com/kind=Elastic,kubedb.com/name=elasticsearch-db'

NAME                         STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
data-elasticsearch-db-pg-0   Bound     pvc-a1a95954-4a75-11e7-8b69-12f236046fba   10Gi       RWO           2m
```

#### Initialize Database

**W**hen we are creating a new Elastic, we can also initialize this database with existing data.

> **Note:** Elastic database supports only SnapshotSource to initialize.

How to initialize database using SnapshotSource? See [here](../initialize-database.md#snapshotsource).


#### Schedule Backup

**W**e can also schedule automatic backup by providing BackupSchedule information in `spec.backupSchedule`.

How to add information in Elastic `spec` to schedule automatic backup? See [here](../schedule-backup.md).


#### Monitor Database

**W**e can also monitor our elasticsearch database.
To enable monitoring, we need to set MonitorSpec in Elastic `spec`.

How to set monitoring? See [here](../monitor-database.md).
