---
title: Upgrade Manual
menu:
  docs_0.9.0:
    identifier: pg-upgrade-manual
    name: Overview
    parent: pg-postgres-guides
    weight: 15
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# KubeDB Upgrade Manual

This tutorial will show you how to upgrade KubeDB from previous version to 0.9.0.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB 0.8.0 cli on your workstation and KubeDB operator in your cluster following the steps [here](https://kubedb.com/docs/0.8.0/setup/install/).

## Previous operator and sample database

In this tutorial we are using helm to install kubedb 0.8.0 release. But, user can install kubedb operator from script too.
Follow [Install instructions](https://github.com/kubedb/project/issues/262) to install kubedb-operator 0.8.0.

```console
$ helm ls
NAME               REVISION    UPDATED                     STATUS      CHART           APP VERSION    NAMESPACE
kubedb-operator    1           Wed Dec 19 15:42:37 2018    DEPLOYED    kubedb-0.8.0    0.8.0          default  
```

Also a sample Postgres database (with scheduled backups) compatible with kubedb-0.8.0 to examine successful upgrade. Read the guide [here](https://kubedb.com/docs/0.8.0/guides/postgres/snapshot/scheduled_backup/#create-postgres-with-backupschedule) to learn about scheduled backup in details.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: scheduled-pg
  namespace: demo
spec:
  version: "9.6"
  replicas: 3
  standbyMode: hot
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  backupSchedule:
    cronExpression: "@every 1m"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb-dev
```

Now create secret and deploy database.

```console
$ kubectl create ns demo
namespace/demo created

$ kubectl create secret -n demo generic gcs-secret \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret/gcs-secret created

$ kubectl create -f scheduled-pg.yaml
postgres.kubedb.com/scheduled-pg created
```

See running scheduled snapshots,

```console
NAME                           DATABASE          STATUS      AGE
scheduled-pg-20181219-094347   pg/scheduled-pg   Succeeded   5m
scheduled-pg-20181219-094447   pg/scheduled-pg   Succeeded   4m
scheduled-pg-20181219-094547   pg/scheduled-pg   Succeeded   3m
scheduled-pg-20181219-094647   pg/scheduled-pg   Succeeded   2m
scheduled-pg-20181219-094747   pg/scheduled-pg   Succeeded   1m
scheduled-pg-20181219-094847   pg/scheduled-pg   Succeeded   29s
```

Node status:

```console
$ kubectl get pods -n demo --show-labels
NAME             READY   STATUS    RESTARTS   AGE     LABELS
scheduled-pg-0   1/1     Running   0          6m12s   controller-revision-hash=scheduled-pg-598d87f567,kubedb.com/kind=Postgres,kubedb.com/name=scheduled-pg,kubedb.com/role=primary,statefulset.kubernetes.io/pod-name=scheduled-pg-0
scheduled-pg-1   1/1     Running   0          5m48s   controller-revision-hash=scheduled-pg-598d87f567,kubedb.com/kind=Postgres,kubedb.com/name=scheduled-pg,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=scheduled-pg-1
scheduled-pg-2   1/1     Running   0          5m47s   controller-revision-hash=scheduled-pg-598d87f567,kubedb.com/kind=Postgres,kubedb.com/name=scheduled-pg,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=scheduled-pg-2
```

### Connect with PostgreSQL database

Now, you can connect to this database using `scheduled-pg.demo` service and *password* created in `scheduled-pg-auth` secret.

**Connection information:**

- Host name/address: you can use any of these
  - Service: `scheduled-pg.demo`
  - Pod IP: (`$ kubectl get pods scheduled-pg-0 -n demo -o yaml | grep podIP`)

  But, In this tutorial we will exec into each pod to insert data and see data availability in replica nodes. So, `localhost` as host is fine.
- Port: `5432`
- Maintenance database: `postgres`
- Username: `postgres`
- Password: Run the following command to get *password*,

  ```console
  $ kubectl get secrets -n demo scheduled-pg-auth -o jsonpath='{.data.\POSTGRES_PASSWORD}' | base64 -d
  9HQ1iTll9IzM7AkL
  ```

### Connect to master-node through cli

```console
$ kubectl exec -it -n demo scheduled-pg-0 bash

bash-4.3# psql -h localhost -U postgres

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
```

Postgres replication state

```console
postgres=# SELECT * FROM pg_stat_replication;
 pid | usesysid | usename  | application_name | client_addr | client_hostname | client_port |         backend_start         | backend_xmin |   state   | sent_location | write_location | flush_location | replay_location | sync_priority | sync_state 
-----+----------+----------+------------------+-------------+-----------------+-------------+-------------------------------+--------------+-----------+---------------+----------------+----------------+-----------------+---------------+------------
  58 |       10 | postgres | scheduled-pg-2   | 172.17.0.11 |                 |       41610 | 2018-12-28 12:42:33.979832+00 |              | streaming | 0/4000060     | 0/4000060      | 0/4000060      | 0/4000060       |             0 | async
  61 |       10 | postgres | scheduled-pg-1   | 172.17.0.10 |                 |       58026 | 2018-12-28 12:42:34.495538+00 |              | streaming | 0/4000060     | 0/4000060      | 0/4000060      | 0/4000060       |             0 | async
(2 rows)
```

Insert data

```console
postgres=# CREATE DATABASE testdb;
CREATE DATABASE

postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 testdb    | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
(4 rows)

postgres=# \c testdb
You are now connected to database "testdb" as user "postgres".

testdb=# \d
No relations found.

testdb=# CREATE TABLE COMPANY(
   ID INT PRIMARY KEY     NOT NULL,
   NAME           TEXT    NOT NULL,
   AGE            INT     NOT NULL,
   ADDRESS        CHAR(50),
   SALARY         REAL,
   JOIN_DATE	  DATE
);
CREATE TABLE

testdb=# \d
          List of relations
 Schema |  Name   | Type  |  Owner   
--------+---------+-------+----------
 public | company | table | postgres
(1 row)

testdb=# INSERT INTO COMPANY (ID,NAME,AGE,ADDRESS,SALARY,JOIN_DATE) VALUES (1, 'Paul', 32, 'California', 20000.00,'2001-07-13'), (2, 'Allen', 25, 'Texas', 20000.00, '2007-12-13'),(3, 'Teddy', 23, 'Norway', 20000.00, '2007-12-13' ), (4, 'Mark', 25, 'Rich-Mond ', 65000.00, '2007-12-13' ), (5, 'David', 27, 'Texas', 85000.00, '2007-12-13');
INSERT 0 5

testdb=# SELECT * FROM company;
 id | name  | age |                      address                       | salary | join_date  
----+-------+-----+----------------------------------------------------+--------+------------
  1 | Paul  |  32 | California                                         |  20000 | 2001-07-13
  2 | Allen |  25 | Texas                                              |  20000 | 2007-12-13
  3 | Teddy |  23 | Norway                                             |  20000 | 2007-12-13
  4 | Mark  |  25 | Rich-Mond                                          |  65000 | 2007-12-13
  5 | David |  27 | Texas                                              |  85000 | 2007-12-13
(5 rows)
```

Exit from postgres cli

```console
testdb=# \q
bash-4.3# exit
exit
```

### Connect to replica-node-1 through cli

Connect to `scheduled-pg-1` and see availability of data

```console
$ kubectl exec -it -n demo scheduled-pg-1 bash

bash-4.3# psql -h localhost -U postgres
psql (9.6.7)
Type "help" for help.

postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 testdb    | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
(4 rows)

postgres=# \c testdb
You are now connected to database "testdb" as user "postgres".

testdb=# \d
          List of relations
 Schema |  Name   | Type  |  Owner   
--------+---------+-------+----------
 public | company | table | postgres
(1 row)

testdb=# SELECT * FROM company;
 id | name  | age |                      address                       | salary | join_date  
----+-------+-----+----------------------------------------------------+--------+------------
  1 | Paul  |  32 | California                                         |  20000 | 2001-07-13
  2 | Allen |  25 | Texas                                              |  20000 | 2007-12-13
  3 | Teddy |  23 | Norway                                             |  20000 | 2007-12-13
  4 | Mark  |  25 | Rich-Mond                                          |  65000 | 2007-12-13
  5 | David |  27 | Texas                                              |  85000 | 2007-12-13
(5 rows)
```

### Connect to replica-node-2 through cli

Connect to `scheduled-pg-2` and see availability of data

```console
$ kubectl exec -it -n demo scheduled-pg-2 bash
bash-4.3# psql -h localhost -U postgres
psql (9.6.7)
Type "help" for help.

postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 testdb    | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
(4 rows)

postgres=# \c testdb
You are now connected to database "testdb" as user "postgres".

testdb=# \d
          List of relations
 Schema |  Name   | Type  |  Owner   
--------+---------+-------+----------
 public | company | table | postgres
(1 row)

testdb=# SELECT * FROM company;
 id | name  | age |                      address                       | salary | join_date  
----+-------+-----+----------------------------------------------------+--------+------------
  1 | Paul  |  32 | California                                         |  20000 | 2001-07-13
  2 | Allen |  25 | Texas                                              |  20000 | 2007-12-13
  3 | Teddy |  23 | Norway                                             |  20000 | 2007-12-13
  4 | Mark  |  25 | Rich-Mond                                          |  65000 | 2007-12-13
  5 | David |  27 | Texas                                              |  85000 | 2007-12-13
(5 rows)

testdb=# \q
bash-4.3# exit
exit
```

## Upgrade kubedb-operator

For helm, `upgrade` command works fine.

```console
$ helm upgrade --install kubedb-operator appscode/kubedb --version 0.9.0
$ helm install appscode/kubedb-catalog --name kubedb-catalog --version 0.9.0 --namespace default

$ helm ls
NAME               REVISION    UPDATED                     STATUS      CHART                   APP VERSION    NAMESPACE
kubedb-catalog     1           Wed Dec 19 15:50:40 2018    DEPLOYED    kubedb-catalog-0.9.0    0.9.0          default  
kubedb-operator    2           Wed Dec 19 15:49:55 2018    DEPLOYED    kubedb-0.9.0            0.9.0          default  
```

For Bash script installation, uninstall first, then install again with 0.9.0 script. See [0.9.0 installation guide](https://kubedb.com/docs/0.9.0/setup/install/).

## Stale CRD objects

At this state, the operator is skipping this `scheduled-pg` Postgres. Because, Postgres version `9.6` is deprecated in `kubedb 0.9.0`.
The scheduled snapshot also in paused state.

```concole
$ kubectl get snap -n demo
NAME                           DATABASENAME   STATUS      AGE
scheduled-pg-20181219-094347   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094447   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094547   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094647   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094747   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094847   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094947   scheduled-pg   Succeeded   1h
```

Some [other fields](https://github.com/kubedb/apimachinery/blob/6319e29148b40f1ac9a7ea312394754e83feba8e/apis/kubedb/v1alpha1/postgres_types.go#L90-L127) in CRD also got deprecated and some are added. The good thing is, kubedb operator will handle those changes in it's mutating webhook. (So, always try to run kubedb with webhooks enabled). But, user has to update the db-version on his own.

## Upgrade CRD objects

Note that, Once the DB version is updated, kubedb-operator will update the statefulsets too. This update strategy can be modified by `spec.updateStrategy`. Read [here](https://kubedb.com/docs/0.9.0/concepts/databases/postgres/#spec-updatestrategy) for details about updateStrategy.

Now, Before updating CRD, find Available PostgresVersion.

```console
$ kubectl get postgresversions
NAME       VERSION   DB_IMAGE                   DEPRECATED   AGE
10.2       10.2      kubedb/postgres:10.2       true         1h
10.2-v1    10.2      kubedb/postgres:10.2-v2                 1h
9.6        9.6       kubedb/postgres:9.6        true         1h
9.6-v1     9.6       kubedb/postgres:9.6-v2                  1h
9.6.7      9.6.7     kubedb/postgres:9.6.7      true         1h
9.6.7-v1   9.6.7     kubedb/postgres:9.6.7-v2                1h
```

Notice the `DEPRECATED` column. Here, `true` means that this PostgresVersion is deprecated for current KubeDB version. KubeDB will not work for deprecated PostgresVersion. To know more about what is `PostgresVersion` crd and why there is `10.2` and `10.2-v1` variation, please visit [here](/docs/concepts/catalog/postgres.md).

Now, Update the CRD and set `Spec.version` to `9.6-v1`.

```console
kubectl edit pg -n demo scheduled-pg
```

See the changed object (defaulted).

```yaml
$ kubectl get pg -n demo scheduled-pg -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  creationTimestamp: "2018-12-28T12:42:21Z"
  finalizers:
  - kubedb.com
  generation: 5
  name: scheduled-pg
  namespace: demo
  resourceVersion: "6080"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/postgreses/scheduled-pg
  uid: 04d98e2c-0a9e-11e9-a1ef-0800270037ae
spec:
  backupSchedule:
    cronExpression: '@every 5m'
    gcs:
      bucket: kubedb-dev
    podTemplate:
      controller: {}
      metadata: {}
      spec:
        resources: {}
    storageSecretName: gcs-secret
  databaseSecret:
    secretName: scheduled-pg-auth
  podTemplate:
    controller: {}
    metadata: {}
    spec:
      resources: {}
  replicas: 3
  serviceTemplate:
    metadata: {}
    spec: {}
  standbyMode: hot
  storage:
    accessModes:
    - ReadWriteOnce
    dataSource: null
    resources:
      requests:
        storage: 1Gi
    storageClassName: standard
  storageType: Durable
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
  version: 9.6-v1
status:
  observedGeneration: 5$4214054550087021099
  phase: Running
```

Watch for pod updates in `RollingUpdate` UpdateStrategy.

```console
$ kubectl get po -n demo -w
NAMESPACE     NAME     READY   STATUS    RESTARTS   AGE
...
demo   scheduled-pg-2   1/1   Terminating   0     107m
demo   scheduled-pg-2   0/1   Terminating   0     107m
demo   scheduled-pg-2   0/1   Terminating   0     107m
demo   scheduled-pg-2   0/1   Terminating   0     107m
demo   scheduled-pg-2   0/1   Pending   0     0s
demo   scheduled-pg-2   0/1   Pending   0     0s
demo   scheduled-pg-2   0/1   ContainerCreating   0     0s
demo   scheduled-pg-2   0/1   ContainerCreating   0     19s
demo   scheduled-pg-2   1/1   Running   0     19s
demo   scheduled-pg-1   1/1   Terminating   0     108m
demo   scheduled-pg-1   0/1   Terminating   0     108m
demo   scheduled-pg-1   0/1   Terminating   0     108m
demo   scheduled-pg-1   0/1   Terminating   0     108m
demo   scheduled-pg-1   0/1   Pending   0     0s
demo   scheduled-pg-1   0/1   Pending   0     0s
demo   scheduled-pg-1   0/1   ContainerCreating   0     0s
demo   scheduled-pg-1   0/1   ContainerCreating   0     1s
demo   scheduled-pg-1   1/1   Running   0     2s
demo   scheduled-pg-0   1/1   Terminating   0     108m
demo   scheduled-pg-0   0/1   Terminating   0     109m
demo   scheduled-pg-0   0/1   Terminating   0     109m
demo   scheduled-pg-0   0/1   Terminating   0     109m
demo   scheduled-pg-0   0/1   Pending   0     0s
demo   scheduled-pg-0   0/1   Pending   0     0s
demo   scheduled-pg-0   0/1   ContainerCreating   0     0s
demo   scheduled-pg-0   0/1   ContainerCreating   0     1s
demo   scheduled-pg-0   1/1   Running   0     2s
```

Watch for statefulset states.

```console
$ kubectl get statefulset -n demo -w
NAME           READY   AGE
scheduled-pg   3/3     105m
scheduled-pg   3/3   107m
scheduled-pg   3/3   107m
scheduled-pg   2/3   107m
scheduled-pg   2/3   107m
scheduled-pg   3/3   108m
scheduled-pg   2/3   108m
scheduled-pg   2/3   108m
scheduled-pg   3/3   108m
scheduled-pg   2/3   109m
scheduled-pg   2/3   109m
scheduled-pg   3/3   109m
```

The scheduled snapshot is in working state now.

```console
$ kubectl get snap -n demo
NAME                           DATABASENAME   STATUS      AGE
scheduled-pg-20181219-094347   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094447   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094547   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094647   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094747   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094847   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-094947   scheduled-pg   Succeeded   1h
scheduled-pg-20181219-111319   scheduled-pg   Succeeded   2m
scheduled-pg-20181219-111419   scheduled-pg   Succeeded   1m
scheduled-pg-20181219-111519   scheduled-pg   Running     4s
```

## Varify Data persistence

### Connect to master-node

```console
$ kubectl exec -it -n demo scheduled-pg-0 bash

bash-4.3# psql -h localhost -U postgres

```

Postgres replication state

```console
postgres=# SELECT * FROM pg_stat_replication;
 pid | usesysid | usename  | application_name | client_addr | client_hostname | client_port |         backend_start         | backend_xmin |   state   | sent_location | write_location | flush_location | replay_location | sync_priority | sync_state 
-----+----------+----------+------------------+-------------+-----------------+-------------+-------------------------------+--------------+-----------+---------------+----------------+----------------+-----------------+---------------+------------
  24 |       10 | postgres | scheduled-pg-1   | 172.17.0.10 |                 |       35058 | 2018-12-28 12:55:35.811518+00 |              | streaming | 0/8000060     | 0/8000060      | 0/8000060      | 0/8000060       |             0 | async
  36 |       10 | postgres | scheduled-pg-2   | 172.17.0.8  |                 |       33610 | 2018-12-28 12:57:31.013411+00 |              | streaming | 0/8000060     | 0/8000060      | 0/8000060      | 0/8000060       |             0 | async
(2 rows)
```

Data availability

```console
postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 testdb    | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
(4 rows)

postgres=# \c testdb
You are now connected to database "testdb" as user "postgres".

testdb=# \d
          List of relations
 Schema |  Name   | Type  |  Owner   
--------+---------+-------+----------
 public | company | table | postgres
(1 row)

testdb=# SELECT * FROM company;
 id | name  | age |                      address                       | salary | join_date  
----+-------+-----+----------------------------------------------------+--------+------------
  1 | Paul  |  32 | California                                         |  20000 | 2001-07-13
  2 | Allen |  25 | Texas                                              |  20000 | 2007-12-13
  3 | Teddy |  23 | Norway                                             |  20000 | 2007-12-13
  4 | Mark  |  25 | Rich-Mond                                          |  65000 | 2007-12-13
  5 | David |  27 | Texas                                              |  85000 | 2007-12-13
(5 rows)


testdb=# \q
bash-4.3# exit
exit
```

### Connect to replica-node-1

Connect to `scheduled-pg-1` and see availability of data

```console
$ kubectl exec -it -n demo scheduled-pg-1 bash

bash-4.3# psql -h localhost -U postgres

postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 testdb    | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
(4 rows)

postgres=# \c testdb
You are now connected to database "testdb" as user "postgres".

testdb=# \d
          List of relations
 Schema |  Name   | Type  |  Owner   
--------+---------+-------+----------
 public | company | table | postgres
(1 row)

testdb=# SELECT * FROM company;
 id | name  | age |                      address                       | salary | join_date  
----+-------+-----+----------------------------------------------------+--------+------------
  1 | Paul  |  32 | California                                         |  20000 | 2001-07-13
  2 | Allen |  25 | Texas                                              |  20000 | 2007-12-13
  3 | Teddy |  23 | Norway                                             |  20000 | 2007-12-13
  4 | Mark  |  25 | Rich-Mond                                          |  65000 | 2007-12-13
  5 | David |  27 | Texas                                              |  85000 | 2007-12-13
(5 rows)
```

### Connect to replica-node-2

Connect to `scheduled-pg-2` and see availability of data

```console
$ kubectl exec -it -n demo scheduled-pg-2 bash

bash-4.3# psql -h localhost -U postgres


postgres=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 testdb    | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
(4 rows)

postgres=# \c testdb
You are now connected to database "testdb" as user "postgres".

testdb=# \d
          List of relations
 Schema |  Name   | Type  |  Owner   
--------+---------+-------+----------
 public | company | table | postgres
(1 row)

testdb=# SELECT * FROM company;
 id | name  | age |                      address                       | salary | join_date  
----+-------+-----+----------------------------------------------------+--------+------------
  1 | Paul  |  32 | California                                         |  20000 | 2001-07-13
  2 | Allen |  25 | Texas                                              |  20000 | 2007-12-13
  3 | Teddy |  23 | Norway                                             |  20000 | 2007-12-13
  4 | Mark  |  25 | Rich-Mond                                          |  65000 | 2007-12-13
  5 | David |  27 | Texas                                              |  85000 | 2007-12-13
(5 rows)

testdb=# \q
bash-4.3# exit
exit
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo pg/scheduled-pg -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo pg/scheduled-pg

kubectl delete ns demo
```

## Next Steps

- Learn about [custom PostgresVersions](/docs/guides/postgres/custom-versions/setup.md).
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Detail concepts of [Postgres object](/docs/concepts/databases/postgres.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
