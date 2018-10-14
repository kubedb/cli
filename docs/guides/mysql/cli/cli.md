---
title: CLI | KubeDB
menu:
  docs_0.9.0-beta.0:
    identifier: my-cli-cli
    name: Quickstart
    parent: my-cli-mysql
    weight: 10
menu_name: docs_0.9.0-beta.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Manage KubeDB objects using CLIs

## KubeDB CLI

KubeDB comes with its own cli. It is called `kubedb` cli. `kubedb` can be used to manage any KubeDB object. `kubedb` cli also performs various validations to improve ux. To install KubeDB cli on your workstation, follow the steps [here](/docs/setup/install.md).

### How to Create objects

`kubedb create` creates a database CRD object in `default` namespace by default. Following command will create a MySQL object as specified in `mysql.yaml`.

```console
$ kubedb create -f mysql-demo.yaml
mysql "mysql-demo" created
```

You can provide namespace as a flag `--namespace`. Provided namespace should match with namespace specified in input file.

```console
$ kubedb create -f mysql-demo.yaml --namespace=kube-system
mysql "mysql-demo" created
```

`kubedb create` command also considers `stdin` as input.

```console
cat mysql-demo.yaml | kubedb create -f -
```

To learn about various options of `create` command, please visit [here](/docs/reference/kubedb_create.md).

### How to List Objects

`kubedb get` command allows users to list or find any KubeDB object. To list all MySQL objects in `default` namespace, run the following command:

```console
$ kubedb get mysql
NAME          STATUS    AGE
mysql-demo    Running   5h
mysql-dev     Running   4h
mysql-prod    Running   30m
mysql-qa      Running   2h
```

To get YAML of an object, use `--output=yaml` flag.

```yaml
$ kubedb get mysql mysql-demo --output=yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  clusterName: ""
  creationTimestamp: 2018-03-01T07:02:10Z
  finalizers:
  - kubedb.com
  generation: 0
  name: mysql-demo
  namespace: default
  resourceVersion: "6910"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/mysqls/mysql-demo
  uid: 76379db5-1d1e-11e8-8599-0800272b52b5
spec:
  databaseSecret:
    secretName: mysql-demo-auth
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: "8.0"
status:
  creationTime: 2018-03-01T07:02:10Z
  phase: Running
```

To get JSON of an object, use `--output=json` flag.

```console
$ kubedb get mysql mysql-demo --output=json
```

To list all KubeDB objects, use following command:

```console
$ kubedb get all -o wide
NAME                VERSION     STATUS   AGE
my/mysql-demo       8.0         Running  3h
my/mysql-dev        8.0         Running  3h
my/mysql-prod       8.0         Running  3h
my/mysql-qa         8.0         Running  3h

NAME                                DATABASE              BUCKET              STATUS      AGE
snap/mysql-demo-20170605-073557     my/mysql-demo         gs:bucket-name      Succeeded   9m
snap/snapshot-20171212-114700       my/mysql-demo         gs:bucket-name      Succeeded   1h
```

Flag `--output=wide` is used to print additional information.

List command supports short names for each object types. You can use it like `kubedb get <short-name>`. Below are the short name for KubeDB objects:

- MySQL: `my`
- Snapshot: `snap`
- DormantDatabase: `drmn`

You can print labels with objects. The following command will list all Snapshots with their corresponding labels.

```console
$ kubedb get snap --show-labels
NAME                          DATABASE              STATUS      AGE       LABELS
mysql-demo-20170605-073557    my/mysql-demo         Succeeded   11m       kubedb.com/kind=MySQL,kubedb.com/name=mysql-demo
snapshot-20171212-114700      my/mysql-demo         Succeeded   1h        kubedb.com/kind=MySQL,kubedb.com/name=mysql-demo
```

You can also filter list using `--selector` flag.

```console
$ kubedb get snap --selector='kubedb.com/kind=MySQL' --show-labels
NAME                          DATABASE         STATUS      AGE       LABELS
mysql-demo-20171212-073557    my/mysql-demo    Succeeded   14m       kubedb.com/kind=MySQL,kubedb.com/name=mysql-demo
snapshot-20171212-114700      my/mysql-demo    Succeeded   2h        kubedb.com/kind=MySQL,kubedb.com/name=mysql-demo
```

To print only object name, run the following command:

```console
$ kubedb get all -o name
mysql/mysql-demo
mysql/mysql-dev
mysql/mysql-prod
mysql/mysql-qa
snapshot/mysql-demo-20170605-073557
snapshot/snapshot-20170505-114700
```

To learn about various options of `get` command, please visit [here](/docs/reference/kubedb_get.md).

### How to Describe Objects

`kubedb describe` command allows users to describe any KubeDB object. The following command will describe MySQL database `mysql-demo` with relevant information.

```console
$ kubedb describe my mysql-demo
Name:		mysql-demo
Namespace:	default
StartTimestamp:	Thu, 01 Mar 2018 15:03:52 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			mysql-demo
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Thu, 01 Mar 2018 13:02:12 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		mysql-demo
  Type:		ClusterIP
  IP:		10.97.55.246
  Port:		db	3306/TCP

Database Secret:
  Name:	mysql-demo-auth
  Type:	Opaque
  Data
  ====
  password:	16 bytes
  user:		4 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From             Type       Reason       Message
  ---------   --------   -----     ----             --------   ------       -------
  2m          2m         1         MySQL operator   Normal     Successful   Successfully patched StatefulSet
  2m          2m         1         MySQL operator   Normal     Successful   Successfully patched MySQL
  2m          2m         1         MySQL operator   Normal     Successful   Successfully patched StatefulSet
  2m          2m         1         MySQL operator   Normal     Successful   Successfully patched MySQL
```

`kubedb describe` command provides following basic information about a MySQL database.

- StatefulSet
- Storage (Persistent Volume)
- Service
- Secret (If available)
- Snapshots (If any)
- Monitoring system (If available)

To hide details about StatefulSet & Service, use flag `--show-workload=false`
To hide details about Secret, use flag `--show-secret=false`
To hide events on KubeDB object, use flag `--show-events=false`

To describe all MySQL objects in `default` namespace, use following command

```console
$ kubedb describe my
```

To describe all MySQL objects from every namespace, provide `--all-namespaces` flag.

```console
$ kubedb describe my --all-namespaces
```

To describe all KubeDB objects from every namespace, use the following command:

```console
$ kubedb describe all --all-namespaces
```

You can also describe KubeDB objects with matching labels. The following command will describe all MySQL objects with specified labels from every namespace.

```console
$ kubedb describe my --all-namespaces --selector='group=dev'
```

To learn about various options of `describe` command, please visit [here](/docs/reference/kubedb_describe.md).

### How to Edit Objects

`kubedb edit` command allows users to directly edit any KubeDB object. It will open the editor defined by _KUBEDB_EDITOR_, or _EDITOR_ environment variables, or fall back to `nano`.

Lets edit an existing running MySQL object to setup [Scheduled Backup](/docs/guides/mysql/snapshot/scheduled-backup.md). The following command will open MySQL `mysql-demo` in editor.

```console
$ kubedb edit my mysql-demo

# Add following under Spec to configure periodic backups
# backupSchedule:
#   cronExpression: '@every 1m'
#   storageSecretName: my-snap-secret
#   gcs:
#     bucket: bucket-name

mysql "mysql-demo" edited
```

#### Edit Restrictions

Various fields of a KubeDB object can't be edited using `edit` command. The following fields are restricted from updates for all KubeDB objects:

- apiVersion
- kind
- metadata.name
- metadata.namespace
- status

If StatefulSets exists for a MySQL database, following fields can't be modified as well.

- spec.version
- spec.databaseSecret
- spec.storage
- spec.nodeSelector
- spec.init

For DormantDatabase, `spec.origin` can't be edited using `kubedb edit`

To learn about various options of `edit` command, please visit [here](/docs/reference/kubedb_edit.md).

### How to Delete Objects

`kubedb delete` command will delete an object in `default` namespace by default unless namespace is provided. The following command will delete a MySQL `mysql-dev` in default namespace

```console
$ kubedb delete mysql mysql-dev
mysql "mysql-dev" deleted
```

You can also use YAML files to delete objects. The following command will delete a mysql using the type and name specified in `mysql.yaml`.

```console
$ kubedb delete -f mysql-demo.yaml
mysql "mysql-dev" deleted
```

`kubedb delete` command also takes input from `stdin`.

```console
cat mysql-demo.yaml | kubedb delete -f -
```

To delete database with matching labels, use `--selector` flag. The following command will delete mysql with label `mysql.kubedb.com/name=mysql-demo`.

```console
$ kubedb delete mysql -l mysql.kubedb.com/name=mysql-demo
```

To learn about various options of `delete` command, please visit [here](/docs/reference/kubedb_delete.md).

## Using Kubectl

You can use Kubectl with KubeDB objects like any other CRDs. Below are some common examples of using Kubectl with KubeDB objects.

```console
# List objects
$ kubectl get mysql
$ kubectl get mysql.kubedb.com

# Delete objects
$ kubectl delete mysql <name>
```

## Next Steps

- Learn how to use KubeDB to run a MySQL database [here](/docs/guides/mysql/README.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
