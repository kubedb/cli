---
title: CLI | KubeDB
menu:
  docs_0.8.0-beta.2:
    identifier: mc-cli-cli
    name: Quickstart
    parent: mc-cli-memcached
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Manage KubeDB objects using CLIs

## KubeDB CLI

KubeDB comes with its own cli. It is called `kubedb` cli. `kubedb` can be used to manage any KubeDB object. `kubedb` cli also performs various validations to improve ux. To install KubeDB cli on your workstation, follow the steps [here](/docs/setup/install.md).

### How to Create objects

`kubedb create` creates a database CRD object in `default` namespace by default. Following command will create a Memcached object as specified in `memcached.yaml`.

```console
$ kubedb create -f memcached-demo.yaml
validating "memcached-demo.yaml"
memcached "memcached-demo" created
```

You can provide namespace as a flag `--namespace`. Provided namespace should match with namespace specified in input file.

```console
$ kubedb create -f memcached-demo.yaml --namespace=kube-system
validating "memcached-demo.yaml"
memcached "memcached-demo" created
```

`kubedb create` command also considers `stdin` as input.

```console
cat memcached-demo.yaml | kubedb create -f -
```

To learn about various options of `create` command, please visit [here](/docs/reference/kubedb_create.md).

### How to List Objects

`kubedb get` command allows users to list or find any KubeDB object. To list all Memcached objects in `default` namespace, run the following command:

```console
$ kubedb get memcached
NAME              STATUS    AGE
memcached-demo    Running   5h
memcached-dev     Running   4h
memcached-prod    Running   30m
memcached-qa      Running   2h
```

To get YAML of an object, use `--output=yaml` flag.

```yaml
$ kubedb get memcached memcached-demo --output=yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  clusterName: ""
  creationTimestamp: 2018-03-01T11:30:18Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb: cli-demo
  name: memcached-demo
  namespace: default
  resourceVersion: "17038"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/memcacheds/memcached-demo
  uid: eb84ab44-1d43-11e8-8599-0800272b52b5
spec:
  replicas: 3
  resources:
    limits:
      cpu: 500m
      memory: 128Mi
    requests:
      cpu: 250m
      memory: 64Mi
  version: 1.5.4
status:
  creationTime: 2018-03-01T11:30:18Z
  phase: Running
```

To get JSON of an object, use `--output=json` flag.

```console
$ kubedb get memcached memcached-demo --output=json
```

To list all KubeDB objects, use following command:

```console
$ kubedb get all -o wide
NAME                    VERSION     STATUS   AGE
mc/memcached-demo       1.5.4       Running  3h
mc/memcached-dev        1.5.4       Running  3h
mc/memcached-prod       1.5.4       Running  3h
mc/memcached-qa         1.5.4       Running  3h
```

Flag `--output=wide` is used to print additional information.

List command supports short names for each object types. You can use it like `kubedb get <short-name>`. Below are the short name for KubeDB objects:

- Memcached: `mc`
- DormantDatabase: `drmn`

You can print labels with objects. The following command will list all Memcached with their corresponding labels.

```console
$ kubedb get mc --show-labels
NAME             STATUS    AGE       LABELS
memcached-demo   Running   2m        kubedb=cli-demo
```

To print only object name, run the following command:

```console
$ kubedb get all -o name
memcached/memcached-demo
memcached/memcached-dev
memcached/memcached-prod
memcached/memcached-qa
```

To learn about various options of `get` command, please visit [here](/docs/reference/kubedb_get.md).

### How to Describe Objects

`kubedb describe` command allows users to describe any KubeDB object. The following command will describe Memcached database `memcached-demo` with relevant information.

```console
$ kubedb describe mc memcached-demo
Name:		memcached-demo
Namespace:	default
StartTimestamp:	Thu, 01 Mar 2018 17:30:18 +0600
Labels:		kubedb=cli-demo
Replicas:	3  total
Status:		Running

Deployment:		
  Name:			memcached-demo
  Replicas:		3 current / 3 desired
  CreationTimestamp:	Thu, 01 Mar 2018 17:30:20 +0600
  Pods Status:		3 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:	
  Name:		memcached-demo
  Type:		ClusterIP
  IP:		10.100.158.166
  Port:		db	11211/TCP

Events:
  FirstSeen   LastSeen   Count     From                 Type       Reason       Message
  ---------   --------   -----     ----                 --------   ------       -------
  2m          2m         1         Memcached operator   Normal     Successful   Successfully created Memcached
  2m          2m         1         Memcached operator   Normal     Successful   Successfully created Deployment
  2m          2m         1         Memcached operator   Normal     Successful   Successfully created Service
```

`kubedb describe` command provides following basic information about a Memcached database.

- Deployment
- Service
- Monitoring system (If available)

To hide details about Deployment & Service, use flag `--show-workload=false`
To hide events on KubeDB object, use flag `--show-events=false`

To describe all Memcached objects in `default` namespace, use following command

```console
$ kubedb describe mc
```

To describe all Memcached objects from every namespace, provide `--all-namespaces` flag.

```console
$ kubedb describe mc --all-namespaces
```

To describe all KubeDB objects from every namespace, use the following command:

```console
$ kubedb describe all --all-namespaces
```

You can also describe KubeDB objects with matching labels. The following command will describe all Memcached objects with specified labels from every namespace.

```console
$ kubedb describe mc --all-namespaces --selector='group=dev'
```

To learn about various options of `describe` command, please visit [here](/docs/reference/kubedb_describe.md).

### How to Edit Objects

`kubedb edit` command allows users to directly edit any KubeDB object. It will open the editor defined by _KUBEDB_EDITOR_, or _EDITOR_ environment variables, or fall back to `nano`.

Lets edit an existing running Memcached object to setup [Monitoring](/docs/guides/memcached/monitoring/using-builtin-prometheus.md). The following command will open Memcached `memcached-demo` in editor.

```console
$ kubedb edit mc memcached-demo

#spec:
#  monitor:
#    agent: prometheus.io/builtin

memcached "memcached-demo" edited
```

#### Edit Restrictions

Various fields of a KubeDB object can't be edited using `edit` command. The following fields are restricted from updates for all KubeDB objects:

- apiVersion
- kind
- metadata.name
- metadata.namespace
- status

If Deployment exists for a Memcached database, following fields can't be modified as well.

- spec.version
- spec.nodeSelector

For DormantDatabase, `spec.origin` can't be edited using `kubedb edit`

To learn about various options of `edit` command, please visit [here](/docs/reference/kubedb_edit.md).

### How to Delete Objects

`kubedb delete` command will delete an object in `default` namespace by default unless namespace is provided. The following command will delete a Memcached `memcached-dev` in default namespace

```console
$ kubedb delete memcached memcached-dev
memcached "memcached-dev" deleted
```

You can also use YAML files to delete objects. The following command will delete a memcached using the type and name specified in `memcached.yaml`.

```console
$ kubedb delete -f memcached-demo.yaml
memcached "memcached-dev" deleted
```

`kubedb delete` command also takes input from `stdin`.

```console
cat memcached-demo.yaml | kubedb delete -f -
```

To delete database with matching labels, use `--selector` flag. The following command will delete memcached with label `memcached.kubedb.com/name=memcached-demo`.

```console
$ kubedb delete memcached -l memcached.kubedb.com/name=memcached-demo
```

To delete all Memcached without following further steps, add flag `--force`

```console
$ kubedb delete memcached -n kube-system --all --force
memcached "memcached-demo" deleted
```

To learn about various options of `delete` command, please visit [here](/docs/reference/kubedb_delete.md).

## Using Kubectl

You can use Kubectl with KubeDB objects like any other CRDs. Below are some common examples of using Kubectl with KubeDB objects.

```console
# List objects
$ kubectl get memcached
$ kubectl get memcached.kubedb.com

# Delete objects
$ kubectl delete memcached <name>
```

## Next Steps

- Learn how to use KubeDB to run a Memcached database [here](/docs/guides/memcached/README.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
