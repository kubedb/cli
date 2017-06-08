### Create Postgres

**L**ets create a simple postgres database using following yaml.

```yaml
apiVersion: kubedb.com/v1beta1
kind: Postgres
metadata:
  name: postgres-db
spec:
  version: 9.5
```

Save this yaml as `postgres-db.yaml` and create Postgres object.

```bash
$ cat postgres-db.yaml | kubedb create -f -

postgres "postgres-db" created
```

**O**ur deployed unified operator will detect this object and will create workloads.

For this object, following kubernetes objects will be created in same namespace:
* StatefulSet (name: **postgres-db**-pg)
* Service (name: **postgres-db**)
* GoverningService (If not available) (name: **kubedb**)
* Secret (name: **postgres-db**-admin-auth)

**A**s secret name is not provided in yaml, a secret will be created with random password.

```bash
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

This secret contains following `ini` data under `.admin` key

```ini
POSTGRES_PASSWORD=vPlT2PzewCaC3XZP
```
> **Note:** default username is **`postgres`**

**N**ow lets see whether our database is ready or not.

```bash
$ kubedb get postgres postgres-db -o wide

NAME          STATUS    VERSION   AGE
postgres-db   Running   9.5       6m
```

This database do not have any PersistentVolume behind StatefulSet.

#### Add storage support

**W**e can create a Postgres database that will use PersistentVolumeClaim in StatefulSet.

How to add storage information in Postgres `spec`? See [here](../support-storage.md).

Following command will list `pvc` for this database.

```bash
$ kubectl get pvc --selector='kubedb.com/kind=Postgres,kubedb.com/name=postgres-db'

NAME                    STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
data-postgres-db-pg-0   Bound     pvc-a1a95954-4a75-11e7-8b69-12f236046fba   10Gi       RWO           2m
```

#### Initialize Database

**W**hen we are creating a new Postgres, we can also initialize this database with existing data.

How to initialize database? See [here](../initialize-database.md).


#### Schedule Backup

**W**e can also schedule automatic backup by providing BackupSchedule information in `spec.backupSchedule`.

How to add information in Postgres `spec` to schedule automatic backup? See [here](../schedule-backup.md).


#### Monitor Database

**W**e can also monitor our postgres database.
To enable monitoring, we need to set MonitorSpec in Postgres `spec`.

How to set monitoring? See [here](../monitor-database.md).

