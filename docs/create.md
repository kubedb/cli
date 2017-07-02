> New to KubeDB? Please start [here](/docs/tutorial.md).

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
