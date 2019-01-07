---
title: Point-in-Time Recovery (PITR) | Postgres
menu:
  docs_0.9.0:
    identifier: pg-pitr
    name: Point-in-Time Recovery
    parent: pg-initialization-postgres
    weight: 25
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Point-in-Time Recovery (PITR) from WAL Source

KubeDB supports Point-in-Time Recovery (PITR) from WAL archive. You can recover your PostgreSQL database to an identical state of a particular time point. This tutorial will show you how to perform Point-in-Time Recovery of a PostgreSQL database with KubeDB.

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- If you are not familiar with Point-in-Time Recovery please read the doc from [here](https://www.postgresql.org/docs/current/continuous-archiving.html).

- You also need to be familiar with WAL Archiving of PostgreSQL database with KubeDB. If you are not familiar with it, please read the guide from [here](/docs/guides/postgres/snapshot/continuous_archiving.md).

To keep things isolated, we are going to use a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created
```

> Note: YAML files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Overview

You can specify `pitr` field while initializing PostgreSQL database from WAL source. This `pitr` filed let you configure the target state in time that you want to achieve after the recovery.

KubeDB allows to configure the following fields in `init.postgresWAL.pitr`:

|       Field       |                         Default Value                          |                                           Uses                                            |
| ----------------- | -------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| `targetTime`      | Up to end of WAL logs                                          | `targetTime` specifies the timestamp up to which recovery will proceed.                   |
| `targetTimeline`  | Same timeline that was current when the base backup was taken. | `targetTimeline` specifies the timeline that you want to recover.                         |
| `targetXID`       | `nil`                                                          | `targetXID` specifies the transaction ID up to which recovery will proceed.               |
| `targetInclusive` | `true`                                                         | `targetInclusive` specifies whether to include ongoing transaction in given target point. |

In this tutorial, we are going to create a sample database and configure continuous WAL archiving in `S3` bucket. Then, we are going to insert some data in this database. Finally, we are going to recover the database from WAL archive up to a particular time point.

## Prepare WAL Archive

Let's deploy a sample database and configure continuous WAL archiving to `S3` bucket.

**Create Storage Secret:**

At first, create a secret for `S3` bucket.

```console
$ echo -n '<your-aws-access-key-id-here>' > AWS_ACCESS_KEY_ID
$ echo -n '<your-aws-secret-access-key-here>' > AWS_SECRET_ACCESS_KEY
$ kubectl create secret -n demo generic s3-secret \
    --from-file=./AWS_ACCESS_KEY_ID \
    --from-file=./AWS_SECRET_ACCESS_KEY
secret "s3-secret" created
```

**Deploy Sample Database:**

Below is the YAML for sample database we are going to deploy,

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: pg-original
  namespace: demo
spec:
  version: "10.2-v1"
  replicas: 1
  terminationPolicy: Delete
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  archiver:
    storage:
      storageSecretName: s3-secret
      s3:
        bucket: kubedb
```

We have configured the above database to continuously backup WAL logs into `kubedb` bucket.

Let's create the database we have shown above,

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/postgres/initialization/pitr/pg-original.yaml
postgres.kubedb.com/pg-original created
```

Now, wait for the database to go into `Running` state.

```console
$ kubectl get pg -n demo pg-original
NAME          VERSION   STATUS    AGE
pg-original   10.2-v1    Running   1m
```

**Insert Sample Data:**

Now, `exec` into the database pod and insert some sample data. Here, we are going to create a table named `pitrDemo`. In order to track the insertion time of a data, we are going to use a separate column named `created_at`. Whenever, a new data is inserted, `created_at` column will be automatically set to the insertion time.

```console
$ kubectl exec -it -n demo pg-original-0 sh
# login as "postgres" superuser.
/ # psql -U postgres
psql (10.2)
Type "help" for help.

# list available databases
postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
(3 rows)

# connect to "postgres" database
postgres=# \c postgres
You are now connected to database "postgres" as user "postgres".

# create a table named "pitrDemo"
postgres=# create table pitrDemo( id serial PRIMARY KEY, message varchar(255) NOT NULL, created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP );
CREATE TABLE

# insert sample data 1
postgres=# insert into pitrDemo(message)  values('row 1 created');
INSERT 0 1

# now wait for 10 minutes then insert sample data 2
postgres=# insert into pitrDemo(message)  values('row 2 created');
INSERT 0 1

# wait for another 10 minutes then insert sample data 3
postgres=# insert into pitrDemo(message)  values('row 3 created');
INSERT 0 1

# let's view the inserted data. check "created_at" column.
# we are going to use this timestamp for point-in-time recovery.
postgres=# select * from pitrDemo;
 id |    message    |          created_at           
----+---------------+-------------------------------
  1 | row 1 created | 2019-01-07 08:28:58.752729+00
  2 | row 2 created | 2019-01-07 08:39:36.77533+00
  3 | row 3 created | 2019-01-07 08:50:12.081161+00
(3 rows)

# quit from the database
postgres=# \q

# exit from the pod
/ # exit
```

> Note that we have used `created_at` column for tracking data insertion time easily. You don't have to use this column if you know when a particular data was inserted.

## Point-in-Time Recovery

Now, we are going to recover PostgreSQL database to an identical state of a specific time point from the archived WAL of `pg-original` database.

At first, let's delete the `pg-original` database so that it does not keep holding the WAL archive.

```console
$ kubectl delete -n demo pg/pg-original
postgres.kubedb.com "pg-original" deleted
```

**Recover up to first sample data:**

We had inserted first sample data at `2019-01-07 08:28:58.752729+00`. So, if we recover with `targetTime: "2019-01-07 08:30:00+00"`, recovered database will contain only first sample data.

Below is the YAML we are going to create to recover up to `"2019-01-07 08:30:00+00"` timestamp.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: pitr-1
  namespace: demo
spec:
  version: "10.2-v1"
  replicas: 1
  terminationPolicy: Delete
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  init:
    postgresWAL:
      storageSecretName: s3-secret
      s3:
        bucket: kubedb
        prefix: 'kubedb/demo/pg-original/archive'
      pitr:
        targetTime: "2019-01-07 08:30:00+00"
```

Let's create the above Postgres crd,

```console
$ kubectl apply -f ./docs/examples/postgres/initialization/pitr/pitr-1.yaml
postgres.kubedb.com/pitr-1 created
```

Now, wait for the Postgres crd `pitr-1` to go into `Running` state.

```console
$ kubectl get pg -n demo pitr-1
NAME     VERSION   STATUS    AGE
pitr-1   10.2-v1    Running   2m
```

Once, the database is running, `exec` into the database pod and check if the recovered database contains only first sample data.

```console
$ kubectl exec -it -n demo pitr-1-0 sh
/ # psql -U postgres
psql (10.2)
Type "help" for help.

postgres=# select * from pitrDemo;
 id |    message    |          created_at           
----+---------------+-------------------------------
  1 | row 1 created | 2019-01-07 08:28:58.752729+00
(1 row)

postgres=#
```

So, we can see that the PostgreSQL database has been recovered to a state where only one sample data was inserted.

**Recover up to second sample data:**

Now, let's recover up to second sample data. We had inserted second sample data at `2019-01-07 08:39:36.77533+00`. So, if we recover with `targetTime: "2019-01-07 08:40:00+00"`, recovered database will contain only the first two sample data.

Below is the YAML we are going to create to recover up to `"2019-01-07 08:40:00+00"` timestamp.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: pitr-2
  namespace: demo
spec:
  version: "10.2-v1"
  replicas: 1
  terminationPolicy: Delete
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  init:
    postgresWAL:
      storageSecretName: s3-secret
      s3:
        bucket: kubedb
        prefix: 'kubedb/demo/pg-original/archive'
      pitr:
        targetTime: "2019-01-07 08:40:00+00"
```

Let's create the above Postgres crd,

```console
$ kubectl apply -f ./docs/examples/postgres/initialization/pitr/pitr-2.yaml
postgres.kubedb.com/pitr-2 created
```

Now, wait for the Postgres crd `pitr-2` to go into `Running` state.

```console
$ kubectl get pg -n demo pitr-2
NAME     VERSION   STATUS    AGE
pitr-2   10.2-v1    Running   2m
```

Once, the database is running, `exec` into the database pod and check if the recovered database contains only the first two sample data.

```console
$ kubectl exec -it -n demo pitr-2-0 sh
/ # psql -U postgres
psql (10.2)
Type "help" for help.

postgres=# select * from pitrDemo;
 id |    message    |          created_at           
----+---------------+-------------------------------
  1 | row 1 created | 2019-01-07 08:28:58.752729+00
  2 | row 2 created | 2019-01-07 08:39:36.77533+00
(2 rows)
```

So, we can see that the PostgreSQL database has been recovered to a state where only two sample data were inserted.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo pg/pitr-1 -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo pg/pitr-1

$ kubectl patch -n demo pg/pitr-2 -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo pg/pitr-2

$ kubectl delete -n demo secret s3-secret
$ kubectl delete ns demo
```
