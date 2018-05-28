---
title: CLI | KubeDB
menu:
  docs_0.8.0-rc.0:
    identifier: rd-cli-cli
    name: Quickstart
    parent: rd-cli-redis
    weight: 10
menu_name: docs_0.8.0-rc.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Manage KubeDB objects using CLIs

## KubeDB CLI

KubeDB comes with its own cli. It is called `kubedb` cli. `kubedb` can be used to manage any KubeDB object. `kubedb` cli also performs various validations to improve ux. To install KubeDB cli on your workstation, follow the steps [here](/docs/setup/install.md).

### How to Create objects

`kubedb create` creates a database CRD object in `default` namespace by default. Following command will create a Redis object as specified in `redis.yaml`.

```console
$ kubedb create -f redis-demo.yaml
redis "redis-demo" created
```

You can provide namespace as a flag `--namespace`. Provided namespace should match with namespace specified in input file.

```console
$ kubedb create -f redis-demo.yaml --namespace=kube-system
redis "redis-demo" created
```

`kubedb create` command also considers `stdin` as input.

```console
cat redis-demo.yaml | kubedb create -f -
```

To learn about various options of `create` command, please visit [here](/docs/reference/kubedb_create.md).

### How to List Objects

`kubedb get` command allows users to list or find any KubeDB object. To list all Redis objects in `default` namespace, run the following command:

```console
$ kubedb get redis
NAME          STATUS    AGE
redis-demo    Running   5h
redis-dev     Running   4h
redis-prod    Running   30m
redis-qa      Running   2h
```

To get YAML of an object, use `--output=yaml` flag.

```yaml
$ kubedb get redis redis-demo --output=yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  clusterName: ""
  creationTimestamp: 2018-03-01T11:08:10Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb: cli-demo
  name: redis-demo
  namespace: default
  resourceVersion: "16072"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/redises/redis-demo
  uid: d42b66be-1d40-11e8-8599-0800272b52b5
spec:
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: "4"
status:
  creationTime: 2018-03-01T11:08:11Z
  phase: Running
```

To get JSON of an object, use `--output=json` flag.

```console
$ kubedb get redis redis-demo --output=json
```

To list all KubeDB objects, use following command:

```console
$ kubedb get all -o wide
NAME                VERSION     STATUS   AGE
rd/redis-demo       4          Running  3h
rd/redis-dev        4          Running  3h
rd/redis-prod       4          Running  3h
rd/redis-qa         4          Running  3h
```

Flag `--output=wide` is used to print additional information.

List command supports short names for each object types. You can use it like `kubedb get <short-name>`. Below are the short name for KubeDB objects:

- Redis: `rd`
- DormantDatabase: `drmn`

You can print labels with objects. The following command will list all Redis with their corresponding labels.

```console
$ kubedb get rd --show-labels
NAME         STATUS    AGE       LABELS
redis-demo   Running   1m        kubedb=cli-demo
```

To print only object name, run the following command:

```console
$ kubedb get all -o name
redis/redis-demo
redis/redis-dev
redis/redis-prod
redis/redis-qa
```

To learn about various options of `get` command, please visit [here](/docs/reference/kubedb_get.md).

### How to Describe Objects

`kubedb describe` command allows users to describe any KubeDB object. The following command will describe Redis database `redis-demo` with relevant information.

```console
$ kubedb describe rd redis-demo
Name:		redis-demo
Namespace:	default
StartTimestamp:	Thu, 01 Mar 2018 17:08:10 +0600
Labels:		kubedb=cli-demo
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			redis-demo
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Thu, 01 Mar 2018 17:08:12 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		redis-demo
  Type:		ClusterIP
  IP:		10.101.218.235
  Port:		db	6379/TCP

Events:
  FirstSeen   LastSeen   Count     From             Type       Reason       Message
  ---------   --------   -----     ----             --------   ------       -------
  7m          7m         1         Redis operator   Normal     Successful   Successfully created Redis
  7m          7m         1         Redis operator   Normal     Successful   Successfully created StatefulSet
  7m          7m         1         Redis operator   Normal     Successful   Successfully created Service
```

`kubedb describe` command provides following basic information about a Redis database.

- StatefulSet
- Storage (Persistent Volume)
- Service
- Monitoring system (If available)

To hide details about StatefulSet & Service, use flag `--show-workload=false`
To hide events on KubeDB object, use flag `--show-events=false`

To describe all Redis objects in `default` namespace, use following command

```console
$ kubedb describe rd
```

To describe all Redis objects from every namespace, provide `--all-namespaces` flag.

```console
$ kubedb describe rd --all-namespaces
```

To describe all KubeDB objects from every namespace, use the following command:

```console
$ kubedb describe all --all-namespaces
```

You can also describe KubeDB objects with matching labels. The following command will describe all Redis objects with specified labels from every namespace.

```console
$ kubedb describe rd --all-namespaces --selector='group=dev'
```

To learn about various options of `describe` command, please visit [here](/docs/reference/kubedb_describe.md).

### How to Edit Objects

`kubedb edit` command allows users to directly edit any KubeDB object. It will open the editor defined by _KUBEDB_EDITOR_, or _EDITOR_ environment variables, or fall back to `nano`.

Lets edit an existing running Redis object to setup [Monitoring](/docs/guides/redis/monitoring/using-builtin-prometheus.md). The following command will open Redis `redis-demo` in editor.

```console
$ kubedb edit rd redis-demo
#spec:
#  monitor:
#    agent: prometheus.io/builtin

redis "redis-demo" edited
```

#### Edit Restrictions

Various fields of a KubeDB object can't be edited using `edit` command. The following fields are restricted from updates for all KubeDB objects:

- apiVersion
- kind
- metadata.name
- metadata.namespace
- status

If StatefulSets exists for a Redis database, following fields can't be modified as well.

- spec.version
- spec.storage
- spec.nodeSelector

For DormantDatabase, `spec.origin` can't be edited using `kubedb edit`

To learn about various options of `edit` command, please visit [here](/docs/reference/kubedb_edit.md).

### How to Delete Objects

`kubedb delete` command will delete an object in `default` namespace by default unless namespace is provided. The following command will delete a Redis `redis-dev` in default namespace

```console
$ kubedb delete redis redis-dev
redis "redis-dev" deleted
```

You can also use YAML files to delete objects. The following command will delete a redis using the type and name specified in `redis.yaml`.

```console
$ kubedb delete -f redis-demo.yaml
redis "redis-dev" deleted
```

`kubedb delete` command also takes input from `stdin`.

```console
cat redis-demo.yaml | kubedb delete -f -
```

To delete database with matching labels, use `--selector` flag. The following command will delete redis with label `redis.kubedb.com/name=redis-demo`.

```console
$ kubedb delete redis -l redis.kubedb.com/name=redis-demo
```

To delete all Redis without following further steps, add flag `--force`

```console
$ kubedb delete redis -n kube-system --all --force
redis "redis-demo" deleted
```

To learn about various options of `delete` command, please visit [here](/docs/reference/kubedb_delete.md).

## Using Kubectl

You can use Kubectl with KubeDB objects like any other CRDs. Below are some common examples of using Kubectl with KubeDB objects.

```console
# List objects
$ kubectl get redis
$ kubectl get redis.kubedb.com

# Delete objects
$ kubectl delete redis <name>
```

## Next Steps

- Learn how to use KubeDB to run a Redis database [here](/docs/guides/redis/README.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
