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

This tutorial will show you how to upgrade KubeDB from previous version to 0.10.0.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB 0.9.0 cli on your workstation and KubeDB operator in your cluster following the steps [here](https://kubedb.com/docs/0.9.0/setup/install/).

## Previous operator and sample database

In this tutorial we are using helm to install kubedb 0.9.0 release. But, user can install kubedb operator from script too.

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update

# Step 1: Install kubedb operator chart
$ helm install appscode/kubedb --name kubedb-operator --version 0.9.0 \
  --namespace kube-system

# Step 2: wait until crds are registered
$ kubectl get crds -l app=kubedb -w
NAME                               AGE
dormantdatabases.kubedb.com        6s
elasticsearches.kubedb.com         12s
elasticsearchversions.kubedb.com   8s
etcds.kubedb.com                   8s
etcdversions.kubedb.com            8s
memcacheds.kubedb.com              6s
memcachedversions.kubedb.com       6s
mongodbs.kubedb.com                7s
mongodbversions.kubedb.com         6s
mysqls.kubedb.com                  7s
mysqlversions.kubedb.com           7s
postgreses.kubedb.com              8s
postgresversions.kubedb.com        7s
redises.kubedb.com                 6s
redisversions.kubedb.com           6s
snapshots.kubedb.com               6s

# Step 3(a): Install KubeDB catalog of database versions
$ helm install appscode/kubedb-catalog --name kubedb-catalog --version 0.9.0 \
  --namespace kube-system

# Step 3(b): Or, if previously installed, upgrade KubeDB catalog of database versions
$ helm upgrade kubedb-catalog appscode/kubedb-catalog --version 0.9.0 \
  --namespace kube-system

$ helm ls
NAME           	REVISION	UPDATED                 	STATUS  	CHART               	APP VERSION	NAMESPACE  
kubedb-catalog 	1       	Fri Feb  8 11:21:34 2019	DEPLOYED	kubedb-catalog-0.9.0	0.9.0      	kube-system
kubedb-operator	1       	Fri Feb  8 11:18:46 2019	DEPLOYED	kubedb-0.9.0        	0.9.0      	kube-system
```

Also a sample Postgres database (with scheduled backups) compatible with kubedb-0.9.0 to examine successful upgrade. Read the guide [here](https://kubedb.com/docs/0.9.0/guides/postgres/snapshot/scheduled_backup/#create-postgres-with-backupschedule) to learn about scheduled backup in details.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: scheduled-pg
  namespace: demo
spec:
  version: "9.6-v2"
  replicas: 3
  standbyMode: Hot
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  backupSchedule:
    cronExpression: "@every 1m"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb-qa
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
$ kubectl get snap -n demo
NAME                           DATABASENAME   STATUS      AGE
scheduled-pg-20190208-053512   scheduled-pg   Succeeded   5m
scheduled-pg-20190208-053612   scheduled-pg   Succeeded   4m
scheduled-pg-20190208-053712   scheduled-pg   Succeeded   3m
scheduled-pg-20190208-053812   scheduled-pg   Succeeded   2m
scheduled-pg-20190208-053912   scheduled-pg   Succeeded   1m
scheduled-pg-20190208-054012   scheduled-pg   Succeeded   29s
```

Node status:

```console
$ kubectl get pods -n demo --show-labels
NAME             READY   STATUS    RESTARTS   AGE   LABELS
scheduled-pg-0   1/1     Running   0          16m   controller-revision-hash=scheduled-pg-75f67456c9,kubedb.com/kind=Postgres,kubedb.com/name=scheduled-pg,kubedb.com/role=primary,statefulset.kubernetes.io/pod-name=scheduled-pg-0
scheduled-pg-1   1/1     Running   0          15m   controller-revision-hash=scheduled-pg-75f67456c9,kubedb.com/kind=Postgres,kubedb.com/name=scheduled-pg,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=scheduled-pg-1
scheduled-pg-2   1/1     Running   0          15m   controller-revision-hash=scheduled-pg-75f67456c9,kubedb.com/kind=Postgres,kubedb.com/name=scheduled-pg,kubedb.com/role=replica,statefulset.kubernetes.io/pod-name=scheduled-pg-2
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
  xfd0mqa3S2Ir0tTP
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
  59 |       10 | postgres | scheduled-pg-1   | 172.17.0.10 |                 |       38916 | 2019-02-08 05:24:47.575398+00 |              | streaming | 0/4000300     | 0/4000300      | 0/4000300      | 0/4000300       |             0 | async
  63 |       10 | postgres | scheduled-pg-2   | 172.17.0.11 |                 |       59954 | 2019-02-08 05:24:51.028714+00 |              | streaming | 0/4000300     | 0/4000300      | 0/4000300      | 0/4000300       |             0 | async
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
$ helm upgrade --install kubedb-operator appscode/kubedb --version 0.10.0
$ helm install appscode/kubedb-catalog --name kubedb-catalog --version 0.10.0 --namespace default

$ helm ls
NAME           	REVISION	UPDATED                 	STATUS  	CHART               	APP VERSION	NAMESPACE  
kubedb-catalog 	2       	Fri Feb  8 12:12:45 2019	DEPLOYED	kubedb-catalog-0.10.0	0.10.0      	kube-system
kubedb-operator	2       	Fri Feb  8 12:11:57 2019	DEPLOYED	kubedb-0.10.0        	0.10.0      	kube-system
```

For Bash script installation, uninstall first, then install again with 0.10.0 script. See [0.9.0 installation guide](https://kubedb.com/docs/0.9.0/setup/install/).

## Stale CRD objects

At this state, the operator is skipping this `scheduled-pg` Postgres. Because, Postgres version `9.6-v1` is deprecated in `kubedb 0.10.0`. You can see the skipped event message in postgres database event. 

```console
$ kubedb describe pg -n demo scheduled-pg
....

Events:
  Type     Reason              Age   From             Message
  ----     ------              ----  ----             -------
  ...
  Normal   Starting            22m   KubeDB operator  Backup running
  Normal   SuccessfulSnapshot  22m   KubeDB operator  Successfully completed snapshot
  Normal   Successful          21m   KubeDB operator  Successfully patched StatefulSet
  Normal   Successful          21m   KubeDB operator  Successfully patched Postgres
  Warning  Invalid             4m    KubeDB operator  postgres demo/scheduled-pg is using deprecated version 9.6-v1. Skipped processing
```

The scheduled snapshot also in paused state.

```concole
$ kubectl get snap -n demo
NAME                           DATABASENAME   STATUS      AGE
...
scheduled-pg-20190208-060615   scheduled-pg   Succeeded   23m
scheduled-pg-20190208-060715   scheduled-pg   Succeeded   22m
scheduled-pg-20190208-060815   scheduled-pg   Succeeded   21m
scheduled-pg-20190208-060915   scheduled-pg   Succeeded   20m
scheduled-pg-20190208-061015   scheduled-pg   Succeeded   19m
scheduled-pg-20190208-061115   scheduled-pg   Succeeded   18m
```

Some other fields in CRD also got deprecated and some are added. The good thing is, kubedb operator will handle those changes in it's mutating webhook. (So, always try to run kubedb with webhooks enabled). But, user has to update the db-version on his own.

## Upgrade CRD objects

Note that, Once the DB version is updated, kubedb-operator will update the statefulsets too. This update strategy can be modified by `spec.updateStrategy`. Read [here](https://kubedb.com/docs/0.9.0/concepts/databases/postgres/#spec-updatestrategy) for details about updateStrategy.

Now, Before updating CRD, find Available PostgresVersion.

```console
$ kubectl get postgresversions
NAME       VERSION   DB_IMAGE                   DEPRECATED   AGE
10.2       10.2      kubedb/postgres:10.2       true         57m
10.2-v1    10.2      kubedb/postgres:10.2-v2    true         57m
10.2-v2    10.2      kubedb/postgres:10.2-v3                 6m
10.6       10.6      kubedb/postgres:10.6                    6m
11.1       11.1      kubedb/postgres:11.1                    6m
9.6        9.6       kubedb/postgres:9.6        true         57m
9.6-v1     9.6       kubedb/postgres:9.6-v2     true         57m
9.6-v2     9.6       kubedb/postgres:9.6-v3                  6m
9.6.7      9.6.7     kubedb/postgres:9.6.7      true         57m
9.6.7-v1   9.6.7     kubedb/postgres:9.6.7-v2   true         57m
9.6.7-v2   9.6.7     kubedb/postgres:9.6.7-v3                6m
```

Notice the `DEPRECATED` column. Here, `true` means that this PostgresVersion is deprecated for current KubeDB version. KubeDB will not work for deprecated PostgresVersion. To know more about what is `PostgresVersion` crd and why there is `10.2` and `10.2-v2` variation, please visit [here](/docs/concepts/catalog/postgres.md).

Now, Update the CRD and set `Spec.version` to `9.6-v2`.

```console
kubectl edit pg -n demo scheduled-pg
```

See the changed object (defaulted).

```yaml
$ kubectl get pg -n demo scheduled-pg -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  creationTimestamp: "2019-02-08T05:24:08Z"
  finalizers:
  - kubedb.com
  generation: 5
  name: scheduled-pg
  namespace: demo
  resourceVersion: "18594"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/postgreses/scheduled-pg
  uid: c2a335ed-2b61-11e9-884f-080027973433
spec:
  backupSchedule:
    cronExpression: '@every 1m'
    gcs:
      bucket: kubedb-qa
    podTemplate:
      controller: {}
      metadata: {}
      spec:
        resources: {}
    storageSecretName: gcs-secret
  databaseSecret:
    secretName: scheduled-pg-auth
  leaderElection:
    leaseDurationSeconds: 15
    renewDeadlineSeconds: 10
    retryPeriodSeconds: 2
  podTemplate:
    controller: {}
    metadata: {}
    spec:
      resources: {}
  replicas: 3
  serviceTemplate:
    metadata: {}
    spec: {}
  standbyMode: Hot
  storage:
    accessModes:
    - ReadWriteOnce
    dataSource: null
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  storageType: Durable
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
  version: 9.6-v2
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
...
scheduled-pg-20190208-060715   scheduled-pg   Succeeded   32m
scheduled-pg-20190208-060815   scheduled-pg   Succeeded   31m
scheduled-pg-20190208-060915   scheduled-pg   Succeeded   30m
scheduled-pg-20190208-061015   scheduled-pg   Succeeded   29m
scheduled-pg-20190208-061115   scheduled-pg   Succeeded   28m
scheduled-pg-20190208-063804   scheduled-pg   Succeeded   1m
scheduled-pg-20190208-063904   scheduled-pg   Succeeded   47s
```

## Varify Data persistence

### Connect to master-node

Get master node

```console
$ kubectl get pods -n demo --selector="kubedb.com/role=primary"
NAME             READY   STATUS    RESTARTS   AGE
scheduled-pg-0   1/1     Running   0          3m12s
```

Exec into postgres master

```console
$ kubectl exec -it -n demo scheduled-pg-0 bash

bash-4.3# psql -h localhost -U postgres

```

Postgres replication state

```console
postgres=# SELECT * FROM pg_stat_replication;
 pid | usesysid | usename  | application_name | client_addr | client_hostname | client_port |         backend_start         | backend_xmin |   state   | sent_location | write_location | flush_location | replay_location | sync_priority | sync_state 
-----+----------+----------+------------------+-------------+-----------------+-------------+-------------------------------+--------------+-----------+---------------+----------------+----------------+-----------------+---------------+------------
  27 |       10 | postgres | scheduled-pg-1   | 172.17.0.10 |                 |       40084 | 2019-02-08 06:39:16.160849+00 |              | streaming | 0/C000108     | 0/C000108      | 0/C000108      | 0/C000108       |             0 | async
  42 |       10 | postgres | scheduled-pg-2   | 172.17.0.11 |                 |       33624 | 2019-02-08 06:41:12.778783+00 |              | streaming | 0/C000108     | 0/C000108      | 0/C000108      | 0/C000108       |             0 | async
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
