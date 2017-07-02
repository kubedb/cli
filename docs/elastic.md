### Create Elastic

**L**ets create a simple elasticsearch database using following yaml.

```yaml
apiVersion: kubedb.com/v1alpha1
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






# Create Database

we can create a database supported by **kubedb** using this CLI.

Lets create a Postgres database.

### kubedb create

`kubedb create` command will create an object in `default` namespace by default unless namespace is specified by input.

Following command will create a Postgres TPR as specified in `postgres.yaml`.

```bash
$ kubedb create -f postgres.yaml

postgres "postgres-demo" created
```

We can provide namespace as a flag `--namespace`.

```bash
$ kubedb create -f postgres.yaml --namespace=kube-system

postgres "postgres-demo" created
```

> Provided namespace should match with namespace specified in input file.

If input file do not specify namespace, object will be created in `default` namespace if not provided.


`kubedb create` command also considers `stdin` as input.

```bash
cat postgres.yaml | kubedb create -f -
```

### Add Storage

**T**o add PersistentVolume support, we need to add following StorageSpec in `spec`

```yaml
spec:
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


##### Click [here](../reference/create.md) to get command details.

### Initialize Database

We now support initialization from two sources.

1. ScriptSource
2. SnapshotSource

We can use one of them to initialize out database.

#### ScriptSource

**W**hen providing ScriptSource to initialize,
a script is run while starting up database.

ScriptSource must have following information:
1. `scriptPath:` ScriptPath (The script you want to run)
2. [VolumeSource](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes) (Where your script and other data will be stored)

##### Example to use GitRepo

```yaml
spec:
  init:
    scriptSource:
      scriptPath: "kubernetes-gitRepo/run.sh"
      gitRepo:
        repository: "https://github.com/appscode/kubernetes-gitRepo.git"
```
When database is starting up, script `run.sh` will be executed.

> **Note:** all path used in script should be relative

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
