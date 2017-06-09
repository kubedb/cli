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
