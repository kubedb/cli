---
title: CLI | KubeDB
menu:
  docs_0.9.0-rc.2:
    identifier: mc-cli-cli
    name: Quickstart
    parent: mc-cli-memcached
    weight: 10
menu_name: docs_0.9.0-rc.2
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
memcached.kubedb.com/memcached-demo created
```

You can provide namespace as a flag `--namespace`. Provided namespace should match with namespace specified in input file.

```console
$ kubedb create -f memcached-demo.yaml --namespace=kube-system
memcached.kubedb.com/memcached-demo created
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
NAME             VERSION    STATUS    AGE
memcached-demo   1.5.4-v1   Running   40s
memcached-dev    1.5.4-v1   Running   40s
memcached-prod   1.5.4-v1   Running   40s
memcached-qa     1.5.4-v1   Running   40s
```

To get YAML of an object, use `--output=yaml` flag.

```yaml
$ kubedb get memcached memcached-demo --output=yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  creationTimestamp: 2018-10-04T05:58:57Z
  finalizers:
  - kubedb.com
  generation: 1
  labels:
    kubedb: cli-demo
  name: memcached-demo
  namespace: default
  resourceVersion: "6883"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/memcacheds/memcached-demo
  uid: 953df4d1-c79a-11e8-bb11-0800272ad446
spec:
  podTemplate:
    controller: {}
    metadata: {}
    spec:
      resources:
        limits:
          cpu: 500m
          memory: 128Mi
        requests:
          cpu: 250m
          memory: 64Mi
  replicas: 3
  serviceTemplate:
    metadata: {}
    spec: {}
  strategy:
    type: RollingUpdate
  terminationPolicy: Pause
  version: 1.5.4-v1
status:
  observedGeneration: 1$7916315637361465932
  phase: Running
```

To get JSON of an object, use `--output=json` flag.

```console
kubedb get memcached memcached-demo --output=json
```

To list all KubeDB objects, use following command:

```console
$ kubedb get all -o wide
NAME                    VERSION     STATUS   AGE
mc/memcached-demo       1.5.4-v1    Running  3h
mc/memcached-dev        1.5.4-v1    Running  3h
mc/memcached-prod       1.5.4-v1    Running  3h
mc/memcached-qa         1.5.4-v1    Running  3h
```

Flag `--output=wide` is used to print additional information.

List command supports short names for each object types. You can use it like `kubedb get <short-name>`. Below are the short name for KubeDB objects:

- Memcached: `mc`
- DormantDatabase: `drmn`

You can print labels with objects. The following command will list all Memcached with their corresponding labels.

```console
$ kubedb get mc --show-labels
NAME             VERSION    STATUS    AGE       LABELS
memcached-demo   1.5.4-v1   Running   2m        kubedb=cli-demo
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
Name:               memcached-demo
Namespace:          default
CreationTimestamp:  Thu, 04 Oct 2018 11:58:57 +0600
Labels:             kubedb=cli-demo
Annotations:        <none>
Replicas:           3  total
Status:             Running

Deployment:
  Name:               memcached-demo
  CreationTimestamp:  Thu, 04 Oct 2018 11:58:59 +0600
  Labels:               kubedb=cli-demo
                        kubedb.com/kind=Memcached
                        kubedb.com/name=memcached-demo
  Annotations:          deployment.kubernetes.io/revision=1
  Replicas:           3 desired | 3 updated | 3 total | 3 available | 0 unavailable
  Pods Status:        3 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         memcached-demo
  Labels:         kubedb=cli-demo
                  kubedb.com/kind=Memcached
                  kubedb.com/name=memcached-demo
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.102.208.191
  Port:         db  11211/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.4:11211,172.17.0.5:11211,172.17.0.6:11211

No Snapshots.

Events:
  Type    Reason      Age   From                Message
  ----    ------      ----  ----                -------
  Normal  Successful  2m    Memcached operator  Successfully created Service
  Normal  Successful  2m    Memcached operator  Successfully created StatefulSet
  Normal  Successful  2m    Memcached operator  Successfully created Memcached
  Normal  Successful  2m    Memcached operator  Successfully patched StatefulSet
  Normal  Successful  2m    Memcached operator  Successfully patched Memcached
```

`kubedb describe` command provides following basic information about a Memcached database.

- Deployment
- Service
- Monitoring system (If available)

To hide events on KubeDB object, use flag `--show-events=false`

To describe all Memcached objects in `default` namespace, use following command

```console
kubedb describe mc
```

To describe all Memcached objects from every namespace, provide `--all-namespaces` flag.

```console
kubedb describe mc --all-namespaces
```

To describe all KubeDB objects from every namespace, use the following command:

```console
kubedb describe all --all-namespaces
```

You can also describe KubeDB objects with matching labels. The following command will describe all Memcached objects with specified labels from every namespace.

```console
kubedb describe mc --all-namespaces --selector='group=dev'
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

- spec.nodeSelector
- spec.podTemplate.spec.nodeSelector
- spec.podTemplate.spec.env

For DormantDatabase, `spec.origin` can't be edited using `kubedb edit`

To learn about various options of `edit` command, please visit [here](/docs/reference/kubedb_edit.md).

### How to Delete Objects

`kubedb delete` command will delete an object in `default` namespace by default unless namespace is provided. The following command will delete a Memcached `memcached-dev` in default namespace

```console
$ kubedb delete memcached memcached-dev
memcached.kubedb.com "memcached-dev" deleted
```

You can also use YAML files to delete objects. The following command will delete a memcached using the type and name specified in `memcached.yaml`.

```console
$ kubedb delete -f memcached-demo.yaml
memcached.kubedb.com "memcached-dev" deleted
```

`kubedb delete` command also takes input from `stdin`.

```console
cat memcached-demo.yaml | kubedb delete -f -
```

To delete database with matching labels, use `--selector` flag. The following command will delete memcached with label `memcached.kubedb.com/name=memcached-demo`.

```console
kubedb delete memcached -l memcached.kubedb.com/name=memcached-demo
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
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
