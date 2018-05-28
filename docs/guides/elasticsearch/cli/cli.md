---
title: CLI | KubeDB
menu:
  docs_0.8.0-rc.0:
    identifier: es-cli-cli
    name: Quickstart
    parent: es-cli-elasticsearch
    weight: 10
menu_name: docs_0.8.0-rc.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Manage KubeDB objects using CLIs

## KubeDB CLI

KubeDB comes with its own cli. It is called `kubedb` cli. `kubedb` can be used to manage any KubeDB object. `kubedb` cli also performs various validations to improve ux. To install KubeDB cli on your workstation, follow the steps [here](/docs/setup/install.md).

### How to Create objects

`kubedb create` creates a database CRD object in `default` namespace by default. Following command will create a Elasticsearch object as specified in `elasticsearch.yaml`.

```console
$ kubedb create -f elasticsearch-demo.yaml
elasticsearch "elasticsearch-demo" created
```

You can provide namespace as a flag `--namespace`. Provided namespace should match with namespace specified in input file.

```console
$ kubedb create -f elasticsearch-demo.yaml --namespace=kube-system
elasticsearch "elasticsearch-demo" created
```

`kubedb create` command also considers `stdin` as input.

```console
cat elasticsearch-demo.yaml | kubedb create -f -
```

To learn about various options of `create` command, please visit [here](/docs/reference/kubedb_create.md).

### How to List Objects

`kubedb get` command allows users to list or find any KubeDB object. To list all Elasticsearch objects in `default` namespace, run the following command:

```console
$ kubedb get elasticsearch
NAME                 STATUS    AGE
elasticsearch-demo   Running   5h
elasticsearch-dev    Running   4h
elasticsearch-prod   Running   30m
elasticsearch-qa     Running   2h
```

To get YAML of an object, use `--output=yaml` flag.

```yaml
$ kubedb get elasticsearch elasticsearch-demo --output=yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  clusterName: ""
  creationTimestamp: 2018-03-02T05:46:59Z
  finalizers:
  - kubedb.com
  generation: 0
  name: elasticsearch-demo
  namespace: default
  resourceVersion: "23697"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/elasticsearches/elasticsearch-demo
  uid: 1fe8f422-1ddd-11e8-b65f-0800272b52b5
spec:
  certificateSecret:
    secretName: elasticsearch-demo-cert
  databaseSecret:
    secretName: elasticsearch-demo-auth
  doNotPause: true
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: "5.6"
status:
  creationTime: 2018-03-02T05:46:59Z
  phase: Running
```

To get JSON of an object, use `--output=json` flag.

```console
$ kubedb get elasticsearch elasticsearch-demo --output=json
```

To list all KubeDB objects, use following command:

```console
$ kubedb get all -o wide
NAME                     VERSION     STATUS  AGE
es/elasticsearch-demo    5.6         Running 3h
es/elasticsearch-dev     5.6         Running 3h
es/elasticsearch-prod    5.6         Running 3h
es/elasticsearch-qa      5.6         Running 3h

NAME                                     DATABASE                     BUCKET              STATUS      AGE
snap/elasticsearch-demo-20170605-073557  es/elasticsearch-demo        gs:bucket-name      Succeeded   9m
snap/snapshot-20171212-114700            es/elasticsearch-demo        gs:bucket-name      Succeeded   1h
```

Flag `--output=wide` is used to print additional information.

List command supports short names for each object types. You can use it like `kubedb get <short-name>`. Below are the short name for KubeDB objects:

- Elasticsearch: `es`
- Snapshot: `snap`
- DormantDatabase: `drmn`

You can print labels with objects. The following command will list all Snapshots with their corresponding labels.

```console
$ kubedb get snap --show-labels
NAME                                 DATABASE                     STATUS      AGE       LABELS
elasticsearch-demo-20170605-073557   es/elasticsearch-demo        Succeeded   11m       kubedb.com/kind=Elasticsearch,kubedb.com/name=elasticsearch-demo
snapshot-20171212-114700             es/elasticsearch-demo        Succeeded   1h        kubedb.com/kind=Elasticsearch,kubedb.com/name=elasticsearch-demo
```

You can also filter list using `--selector` flag.

```console
$ kubedb get snap --selector='kubedb.com/kind=Elasticsearch' --show-labels
NAME                                 DATABASE                STATUS      AGE       LABELS
elasticsearch-demo-20171212-073557   es/elasticsearch-demo   Succeeded   14m       kubedb.com/kind=Elasticsearch,kubedb.com/name=elasticsearch-demo
snapshot-20171212-114700             es/elasticsearch-demo   Succeeded   2h        kubedb.com/kind=Elasticsearch,kubedb.com/name=elasticsearch-demo
```

To print only object name, run the following command:

```console
$ kubedb get all -o name
elasticsearch/elasticsearch-demo
elasticsearch/elasticsearch-dev
elasticsearch/elasticsearch-prod
elasticsearch/elasticsearch-qa
snapshot/elasticsearch-demo-20170605-073557
snapshot/snapshot-20170505-114700
snapshot/snapshot-xyz
```

To learn about various options of `get` command, please visit [here](/docs/reference/kubedb_get.md).

### How to Describe Objects

`kubedb describe` command allows users to describe any KubeDB object. The following command will describe Elasticsearch database `elasticsearch-demo` with relevant information.

```console
$ kubedb describe es elasticsearch-demo
Name:			elasticsearch-demo
Namespace:		default
CreationTimestamp:	Fri, 02 Mar 2018 11:46:59 +0600
Status:			Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			elasticsearch-demo
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Fri, 02 Mar 2018 11:47:00 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		elasticsearch-demo
  Type:		ClusterIP
  IP:		10.110.91.198
  Port:		http	9200/TCP

Service:
  Name:		elasticsearch-demo-master
  Type:		ClusterIP
  IP:		10.100.126.132
  Port:		transport	9300/TCP

Database Secret:
  Name:	elasticsearch-demo-auth
  Type:	Opaque
  Data
  ====
  sg_internal_users.yml:	156 bytes
  sg_roles.yml:			312 bytes
  sg_roles_mapping.yml:		73 bytes
  ADMIN_PASSWORD:		8 bytes
  READALL_PASSWORD:		8 bytes
  sg_action_groups.yml:		430 bytes
  sg_config.yml:		240 bytes

Certificate Secret:
  Name:	elasticsearch-demo-cert
  Type:	Opaque
  Data
  ====
  root.jks:	864 bytes
  sgadmin.jks:	3010 bytes
  key_pass:	6 bytes
  node.jks:	3014 bytes

Topology:
  Type                 Pod                    StartTime                       Phase
  ----                 ---                    ---------                       -----
  client|data|master   elasticsearch-demo-0   2018-03-02 11:47:00 +0600 +06   Running

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason       Message
  ---------   --------   -----     ----                     --------   ------       -------
  3m          3m         1         Elasticsearch operator   Normal     Successful   Successfully patched Elasticsearch
  3m          3m         1         Elasticsearch operator   Normal     Successful   Successfully patched StatefulSet
  4m          4m         1         Elasticsearch operator   Normal     Successful   Successfully created StatefulSet
```

`kubedb describe` command provides following basic information about a database.

- StatefulSet
- Storage (Persistent Volume)
- Service
- Secret (If available)
- Topology (If available)
- Snapshots (If any)
- Monitoring system (If available)

To hide details about StatefulSet & Service, use flag `--show-workload=false`
To hide details about Secret, use flag `--show-secret=false`
To hide events on KubeDB object, use flag `--show-events=false`

To describe all Elasticsearch objects in `default` namespace, use following command

```console
$ kubedb describe es
```

To describe all Elasticsearch objects from every namespace, provide `--all-namespaces` flag.

```console
$ kubedb describe es --all-namespaces
```

To describe all KubeDB objects from every namespace, use the following command:

```console
$ kubedb describe all --all-namespaces
```

You can also describe KubeDb objects with matching labels. The following command will describe all Elasticsearch objects with specified labels from every namespace.

```console
$ kubedb describe es --all-namespaces --selector='group=dev'
```

To learn about various options of `describe` command, please visit [here](/docs/reference/kubedb_describe.md).

### How to Edit Objects

`kubedb edit` command allows users to directly edit any KubeDB object. It will open the editor defined by _KUBEDB_EDITOR_, or _EDITOR_ environment variables, or fall back to `nano`.

Lets edit an existing running Elasticsearch object to setup [Scheduled Backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md). The following command will open Elasticsearch `elasticsearch-demo` in editor.

```console
$ kubedb edit es elasticsearch-demo

# Add following under Spec to configure periodic backups
#  backupSchedule:
#    cronExpression: "@every 6h"
#    storageSecretName: "secret-name"
#    gcs:
#      bucket: "bucket-name"

elasticsearch "elasticsearch-demo" edited
```

#### Edit restrictions

Various fields of a KubeDb object can't be edited using `edit` command. The following fields are restricted from updates for all KubeDB objects:

- apiVersion
- kind
- metadata.name
- metadata.namespace
- status

If StatefulSets or Deployments exists for a database, following fields can't be modified as well.

Elasticsearch:

- spec.version
- spec.storage
- spec.nodeSelector
- spec.init

For DormantDatabase, _spec.origin_ can't be edited using `kubedb edit`

To learn about various options of `edit` command, please visit [here](/docs/reference/kubedb_edit.md).

### How to Delete Objects

`kubedb delete` command will delete an object in `default` namespace by default unless namespace is provided. The following command will delete a Elasticsearch `elasticsearch-dev` in default namespace

```console
$ kubedb delete elasticsearch elasticsearch-dev
elasticsearch "elasticsearch-dev" deleted
```

You can also use YAML files to delete objects. The following command will delete a elasticsearch using the type and name specified in `elasticsearch.yaml`.

```console
$ kubedb delete -f elasticsearch.yaml
elasticsearch "elasticsearch-dev" deleted
```

`kubedb delete` command also takes input from `stdin`.

```console
cat elasticsearch.yaml | kubedb delete -f -
```

To delete database with matching labels, use `--selector` flag. The following command will delete elasticsearch with label `elasticsearch.kubedb.com/name=elasticsearch-demo`.

```console
$ kubedb delete elasticsearch -l elasticsearch.kubedb.com/name=elasticsearch-demo
```

To learn about various options of `delete` command, please visit [here](/docs/reference/kubedb_delete.md).

## Using Kubectl

You can use Kubectl with KubeDB objects like any other CRDs. Below are some common examples of using Kubectl with KubeDB objects.

```console
# List objects
$ kubectl get elasticsearch
$ kubectl get elasticsearch.kubedb.com

# Delete objects
$ kubectl delete elasticsearch <name>
```

## Next Steps

- Learn how to use KubeDB to run a Elasticsearch database [here](/docs/guides/elasticsearch/README.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
