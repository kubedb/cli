> New to KubeDB? Please start [here](/docs/tutorial.md).

# Postgreses

## What is Postgres
A `Postgres` is a Kubernetes `Third Party Object` (TPR). It provides declarative configuration for [PostgreSQL](https://www.postgresql.org/) in a Kubernetes native way. You only need to describe the desired database configuration in a Postgres object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## Postgres Spec
As with all other Kubernetes objects, a Postgres needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example Postgres object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: postgres-db
spec:
  version: 9.5
```

```sh
$ kubedb create -f ./docs/examples/postgres/postgres.yaml

postgres "postgres-db" created
```

Once the Postgres object is created, KubeDB operator will detect it and create the following Kubernetes objects in the same namespace:
* StatefulSet (name: **postgres-db**-pg)
* Service (name: **postgres-db**)
* GoverningService (If not available) (name: **kubedb**)
* Secret (name: **postgres-db**-admin-auth)

Since secret name is not provided during creating Postgres object, a secret will be created with random password.

```yaml
$ kubectl get secret postgres-db-admin-auth -o yaml
apiVersion: v1
data:
  .admin: UE9TVEdSRVNfUEFTU1dPUkQ9dlBsVDJQemV3Q2FDM1haUAo=
kind: Secret
metadata:
  labels:
    kubedb.com/kind: Postgres
  name: postgres-db-admin-auth
  namespace: default
type: Opaque
```

The `.admin` contains a `ini` formatted key/value pairs. 

```ini
POSTGRES_PASSWORD=vPlT2PzewCaC3XZP
```
> **Note:** default username is **`postgres`**

To confirm the new PostgreSQL database is ready, run the following command:

```sh
$ kubedb get postgres postgres-db -o wide

NAME          VERSION   STATUS    AGE
postgres-db   9.5       Running   34m
```

This database does not have any PersistentVolume behind StatefulSet pods.


### Using PersistentVolume
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

```sh
$ kubectl get pvc --selector='kubedb.com/kind=Postgres,kubedb.com/name=postgres-db'

NAME                 STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
data-postgres-db-0   Bound     pvc-a1a95954-4a75-11e7-8b69-12f236046fba   10Gi       RWO           2m
```


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
