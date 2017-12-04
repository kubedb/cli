---
title: CLI | KubeDB
menu:
  docs_0.7.1:
    identifier: tutorials-cli
    name: CLI
    parent: tutorials
    weight: 60
menu_name: docs_0.7.1
section_menu_id: tutorials
---

> New to KubeDB? Please start [here](/docs/tutorials/README.md).

# Manage KubeDB objects using CLIs

## KubeDB CLI
KubeDB comes with its own cli. It is called `kubedb` cli. `kubedb` can be used to deploy KubeDB operator in a cluster and manage all KubeDB tprs. `kubedb` cli also performs various validations to improve ux. To install KubeDB cli on your workstation, follow the steps [here](/docs/install.md).

### How to Create objects

`kubedb create` creates a database tpr in `default` namespace by default. Following command will create a Postgres CRD as specified in `postgres.yaml`.

```console
$ kubedb create -f ./docs/examples/postgres/postgres.yaml

postgres "postgres-demo" created
```

You can provide namespace as a flag `--namespace`. Provided namespace should match with namespace specified in input file.

```console
$ kubedb create -f postgres.yaml --namespace=kube-system

postgres "postgres-demo" created
```

`kubedb create` command also considers `stdin` as input.

```console
cat postgres.yaml | kubedb create -f -
```

To learn about various options of `create` command, please visit [here](/docs/reference/kubedb_create.md).

###  How to List Objects

`kubedb get` command allows users to list or find any KubeDB object. To list all Postgres objects in `default` namespace, run the following command:

```console
$ kubedb get postgres

NAME            STATUS    AGE
postgres-demo   Running   5h
postgres-dev    Running   4h
postgres-prod   Running   30m
postgres-qa     Running   2h
```

To get YAML of an object, use `--output=yaml` flag.

```yaml
$ kubedb get postgres postgres-demo --output=yaml

apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: postgres-demo
  namespace: default
spec:
  databaseSecret:
    secretName: postgres-demo-admin-auth
  version: "9.5"
status:
  creationTime: 2017-06-05T04:10:06Z
  phase: Running
```

To get JSON of an object, use `--output=json` flag.

```console
$ kubedb get postgres postgres-demo --output=json
```

To list all KubeDB objects, use following command:

```console
$ kubedb get all -o wide

NAME                    VERSION   STATUS    AGE
es/elasticsearch-demo   2.3.1     Running   17m

NAME               VERSION   STATUS    AGE
pg/postgres-demo   9.5       Running   3h
pg/postgres-dev    9.5       Running   3h
pg/postgres-prod   9.5       Running   3h
pg/postgres-qa     9.5       Running   3h

NAME                                 DATABASE                BUCKET             STATUS      AGE
snap/postgres-demo-20170605-073557   pg/postgres-demo        gs:bucket-name     Succeeded   9m
snap/snapshot-20170505-1147          pg/postgres-demo        gs:bucket-name     Succeeded   1h
snap/snapshot-xyz                    es/elasticsearch-demo   local:/directory   Succeeded   5m
```

Flag `--output=wide` is used to print additional information.

List command supports short names for each object types. You can use it like `kubedb get <short-name>`. Below are the short name for KubeDB objects:
 - Postgres: `pg`
 - Elasticsearch: `es`
 - Snapshot: `snap`
 - DormantDatabase: `drmn`

You can print labels with objects. The following command will list all Snapshots with their corresponding labels.

```console
$ kubedb get snap --show-labels

NAME                            DATABASE                STATUS      AGE       LABELS
postgres-demo-20170605-073557   pg/postgres-demo        Succeeded   11m       kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
snapshot-20170505-1147          pg/postgres-demo        Succeeded   1h        kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
snapshot-xyz                    es/elasticsearch-demo   Succeeded   6m        kubedb.com/kind=Elasticsearch,kubedb.com/name=elasticsearch-demo
```

You can also filter list using `--selector` flag.

```console
$ kubedb get snap --selector='kubedb.com/kind=Postgres' --show-labels

NAME                            DATABASE           STATUS      AGE       LABELS
postgres-demo-20170605-073557   pg/postgres-demo   Succeeded   14m       kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
snapshot-20170505-1147          pg/postgres-demo   Succeeded   2h        kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
```

To print only object name, run the following command:
```console
$ kubedb get all -o name

elastic/elasticsearch-demo
postgres/postgres-demo
postgres/postgres-dev
postgres/postgres-prod
postgres/postgres-qa
snapshot/postgres-demo-20170605-073557
snapshot/snapshot-20170505-1147
snapshot/snapshot-xyz
```

To learn about various options of `get` command, please visit [here](/docs/reference/kubedb_get.md).

### How to Describe Objects

`kubedb describe` command allows users to describe any KubeDB object. The following command will describe Postgres database `postgres-demo` with relevant information.

```console
$ kubedb describe pg postgres-demo

Name:		postgres-demo
Namespace:	default
StartTimestamp:	Mon, 05 Jun 2017 10:10:06 +0600
Status:		Running
No volumes.

StatefulSet:
  Name:			postgres-demo
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Mon, 05 Jun 2017 10:10:14 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		postgres-demo
  Type:		ClusterIP
  IP:		10.0.0.36
  Port:		port	5432/TCP

Database Secret:
  Name:	postgres-demo-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	35 bytes

Snapshots:
  Name                     Bucket          StartTime                         CompletionTime                    Phase
  ----                     ------          ---------                         --------------                    -----
  postgres-demo-20170605-073557   database-test   Mon, 05 Jun 2017 13:35:57 +0600   Mon, 05 Jun 2017 13:36:10 +0600   Succeeded
  snapshot-20170505-1147          database-test   Mon, 05 Jun 2017 11:48:06 +0600   Mon, 05 Jun 2017 12:01:39 +0600   Succeeded

Events:
  FirstSeen   LastSeen   Count     From                  Type       Reason               Message
  ---------   --------   -----     ----                  --------   ------               -------
  3m          3m         1         Snapshot Controller   Normal     Starting             Backup running
  21m         21m        1         Postgres operator     Normal     SuccessfulCreate     Successfully created StatefulSet
  21m         21m        1         Postgres operator     Normal     SuccessfulCreate     Successfully created Postgres
  29m         29m        1         Postgres operator     Normal     SuccessfulValidate   Successfully validate Postgres
  29m         29m        1         Postgres operator     Normal     Creating             Creating Kubernetes objects
```

`kubedb describe` command provides following basic information about a database.

* StatefulSet
* Storage (Persistent Volume)
* Service
* Secret (If available)
* Snapshots (If any)
* Monitoring system (If available)

This command also shows events unless `--show-events=false`

To describe all Postgres objects in `default` namespace, use following command
```console
$ kubedb describe pg
```

To describe all Postgres objects from every namespace, provide `--all-namespaces` flag.
```console
$ kubedb describe pg --all-namespaces
```

To describe all KubeDB objects from every namespace, use the following command:
```console
$ kubedb describe all --all-namespaces
```

You can also describe KubeDb objects with matching labels. The following command will describe all Elasticsearch & Postgres objects with specified labels from every namespace.

```bash
$ kubedb describe pg,es --all-namespaces --selector='group=dev'
```

To learn about various options of `describe` command, please visit [here](/docs/reference/kubedb_describe.md).


### How to Edit Objects

`kubedb edit` command allows users to directly edit any KubeDB object. It will open the editor defined by _KUBEDB_EDITOR_, or _EDITOR_ environment variables, or fall back to `nano`.

Lets edit an existing running Postgres database to setup [Scheduled Backup](/docs/backup.md). The following command will open Postgres `postgres-demo` in editor.

```bash
$ kubedb edit pg postgres-demo

# Add following under Spec to configure periodic backups
#  backupSchedule:
#    cronExpression: "@every 6h"
#    storageSecretName: "secret-name"
#   gcs:
#      bucket: "bucket-name"

postgres "postgres-demo" edited
```

#### Edit restrictions
Various fields of a KubeDb object can't be edited using `edit` command. The following fields are restricted from updates for all KubeDB objects:
* _apiVersion_
* _kind_
* _metadata.name_
* _metadata.namespace_
* _status_


If StatefulSet exists for a database, following fields can't be modified as well.

Postgres:
* _spec.version_
* _spec.storage_
* _spec.databaseSecret_
* _spec.nodeSelector_
* _spec.init_

Elasticsearch:
* _spec.version_
* _spec.storage_
* _spec.nodeSelector_
* _spec.init_

For DormantDatabase, _spec.origin_ can't be edited using `kubedb edit`

To learn about various options of `edit` command, please visit [here](/docs/reference/kubedb_edit.md).

### How to Summarize Databases
`kubedb summarize` command can be used to generate a JSON formatted summary report for any supported database. The summary contains various stats on database tables and/or indices like, number of rows, the maximum id. This report is intended to be used as a tool to quickly verify whether backup/restore process has worked properly or not. To learn about various options of `summarize` command, please visit [here](/docs/reference/kubedb_summarize.md).

```console
$ kubedb summarize pg p1 -n demo
E0719 08:32:47.285561   16159 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:36226: bind: cannot assign requested address
E0719 08:32:47.791193   16159 portforward.go:317] error copying from local connection to remote stream: read tcp4 127.0.0.1:36226->127.0.0.1:52904: read: connection reset by peer
Summary report for "postgreses/p1" has been stored in 'report-20170719-153247.json'
```

`kubed compare` command compares two summary reports for the same type of database. By default it dumps a git diff-like output on terminal. To learn about various options of `compare` command, please visit [here](/docs/reference/kubedb_compare.md).


```yaml
$ kubedb compare report-20170719-152824.json report-20170719-153247.json
Comparison result has been stored in 'result-20170719-153401.txt'.

 {
   "apiVersion": "kubedb.com/v1alpha1",
   "kind": "Postgres",
   "metadata": {
     "creationTimestamp": "2017-07-19T15:11:45Z",
     "name": "p1",
     "namespace": "demo"
   },
   "status": {
-    "completionTime": "2017-07-19T15:28:24Z",
+    "completionTime": "2017-07-19T15:32:47Z",
-    "startTime": "2017-07-19T15:28:24Z"
+    "startTime": "2017-07-19T15:32:47Z"
   },
   "summary": {
     "postgres": {
       "postgres": {
...
```

### How to Delete Objects

`kubedb delete` command will delete an object in `default` namespace by default unless namespace is provided. The following command will delete a Postgres `postgres-dev` in default namespace

```console
$ kubedb delete postgres postgres-dev

postgres "postgres-dev" deleted
```

You can also use YAML files to delete objects. The following command will delete a postgres using the type and name specified in `postgres.yaml`.

```console
$ kubedb delete -f postgres.yaml

postgres "postgres-dev" deleted
```

`kubedb delete` command also takes input from `stdin`.

```console
cat postgres.yaml | kubedb delete -f -
```

To delete database with matching labels, use `--selector` flag. The following command will delete postgres with label `postgres.kubedb.com/name=postgres-demo`.
```console
$ kubedb delete postgres -l postgres.kubedb.com/name=postgres-demo
```

To learn about various options of `delete` command, please visit [here](/docs/reference/kubedb_delete.md).

## Using Kubectl
Kubectl has limited support for CRDs in general. You can use Kubectl with KubeDB objects like any other CRDs. Below are some common examples of using Kubectl with KubeDB objects.

```console
# List objects
$ kubectl get postgres
$ kubectl get postgres.kubedb.com

# Delete objects
$ kubectl delete postgres <name>
```


## Next Steps
- Learn how to use KubeDB to run a PostgreSQL database [here](/docs/tutorials/postgres/README.md).
- Learn how to use KubeDB to run an Elasticsearch database [here](/docs/tutorials/elasticsearch/README.md).
- Wondering what features are coming next? Please visit [here](/ROADMAP.md). 
- Want to hack on KubeDB? Check our [contribution guidelines](/CONTRIBUTING.md).
