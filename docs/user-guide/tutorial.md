# Tutorial

### Initialize unified operator

**F**irst of all, we need an unified [operator](https://github.com/k8sdb/operator) to handle supported TPR.

We can deploy this operator using `kubedb` CLI.

```bash
$ kubedb init

Successfully created operator deployment.
Successfully created operator service.
```

This will deploy a controller that can handle our TPR objects.

For more, see [init](operation/init.md) operation.

When operator is ready, we can create database.

### Create database

**W**e are supporting two kinds of databases. Lets see how to create them.

* Create [Postgres](database/postgres/create.md) database
* Create [Elasticsearch](database/elastic/create.md) database

### Describe database

**W**e can describe our database using `kubedb` CLI. 

```bash
$ kubedb describe <pg|es> database-demo
```
The `kubedb describe` command provides following basic information of TPR object.
* StatefulSet
* Storage (Persistent Volume) (If available)
* Service
* Secret (If available)
* Snapshots (If any)
* Monitoring system (If available)

For more, see [describe](task/describe.md) operation.

### Take instant backup

**W**e can take backup of a running database using Snapshot TPR object.

This Snapshot object will contain all necessary information to take backup and upload backup data to cloud.

For details explanation, see [here](database/backup.md) .


### Lock database

**T**o prevent a database from deleting, we can set `doNotPause: true`.

```bash
$ kubedb edit <pg|es> database-demo

# Add following in spec
#  doNotPause: true
```

To undo this, set `doNotPause: false`.

### Pause database

**P**ausing a database means operator will only delete _StatefulSet_ & _Service_. Database can be resumed later.

Lets pause a database by deleting it.

```bash
$ kubedb delete <pg|es> database-demo
```

For deleting this Postgres object, operator will create a DormantDatabase object `database-demo`.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  labels:
    kubedb.com/kind: <Postgres|Elastic>
  name: database-demo
spec:
  origin:
    metadata:
      name: database-demo
      namespace: default
    spec:
      -----
      -----
status:
  creationTime: 2017-06-06T08:11:10Z
  pausingTime: 2017-06-06T08:11:50Z
  phase: Paused
```

This DormantDatabase object will contain original database information in `spec.origin`.

**W**hen this DormantDatabase object will be created, operator will delete following workloads:

* StatefulSet
* Service

And following kubernetes objects will be intact:

* Secret
* GoverningService
* PersistentVolumeClaim

We can get back our database again by resuming it.

### Resume database

**T**o resume our database `database-demo` from DormantDatabase object,
we need to edit this DormantDatabase object to set `resume: true`.

```bash
$ kubedb edit drmn database-demo

# Add following in spec
#  resume: true

dormantdatabase "database-demo" edited
```

Operator will detect this modification and will create database TPR object `database-demo`.

**T**his new database TPR object will use same _Secret_ and StatefulSet will be created with existing PersistentVolumeClaim.

### Wipeout database

**T**o wipeout a database, we need to set `wipeOut: true` in DormantDatabase object

```bash
$ kubedb edit drmn database-demo

# Add following in spec
#  wipeOut: true

dormantdatabase "database-demo" edited
```

This will delete following kubernetes objects:

* Secret
* PersistentVolumeClaim

And also delete all Snapshot objects.

Once database is Wipedout, there is no way to resume it.
