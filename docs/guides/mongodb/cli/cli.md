---
title: CLI | KubeDB
menu:
  docs_0.9.0:
    identifier: mg-cli-cli
    name: Quickstart
    parent: mg-cli-mongodb
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Manage KubeDB objects using CLIs

## KubeDB CLI

KubeDB comes with its own cli. It is called `kubedb` cli. `kubedb` can be used to manage any KubeDB object. `kubedb` cli also performs various validations to improve ux. To install KubeDB cli on your workstation, follow the steps [here](/docs/setup/install.md).

### How to Create objects

`kubedb create` creates a database CRD object in `default` namespace by default. Following command will create a MongoDB object as specified in `mongodb.yaml`.

```console
$ kubedb create -f mongodb-demo.yaml
mongodb.kubedb.com/mongodb-demo created
```

You can provide namespace as a flag `--namespace`. Provided namespace should match with namespace specified in input file.

```console
$ kubedb create -f mongodb-demo.yaml --namespace=kube-system
mongodb.kubedb.com/mongodb-demo
```

`kubedb create` command also considers `stdin` as input.

```console
cat mongodb-demo.yaml | kubedb create -f -
```

To learn about various options of `create` command, please visit [here](/docs/reference/kubedb_create.md).

### How to List Objects

`kubedb get` command allows users to list or find any KubeDB object. To list all MongoDB objects in `default` namespace, run the following command:

```console
$ kubedb get mongodb
NAME           VERSION   STATUS    AGE
mongodb-demo   3.4-v2    Running   13m
mongodb-dev    3.4-v2    Running   11m
mongodb-prod   3.4-v2    Running   11m
mongodb-qa     3.4-v2    Running   10m
```

To get YAML of an object, use `--output=yaml` flag.

```yaml
$ kubedb get mongodb mongodb-demo --output=yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  creationTimestamp: "2019-02-06T10:31:04Z"
  finalizers:
  - kubedb.com
  generation: 2
  name: mongodb-demo
  namespace: default
  resourceVersion: "94703"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/default/mongodbs/mongodb-demo
  uid: 4eaaba0e-29fa-11e9-aebf-080027875192
spec:
  databaseSecret:
    secretName: mongodb-demo-auth
  podTemplate:
    controller: {}
    metadata: {}
    spec:
      livenessProbe:
        exec:
          command:
          - mongo
          - --eval
          - db.adminCommand('ping')
        failureThreshold: 3
        periodSeconds: 10
        successThreshold: 1
        timeoutSeconds: 5
      readinessProbe:
        exec:
          command:
          - mongo
          - --eval
          - db.adminCommand('ping')
        failureThreshold: 3
        periodSeconds: 10
        successThreshold: 1
        timeoutSeconds: 1
      resources: {}
  replicas: 1
  serviceTemplate:
    metadata: {}
    spec: {}
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
  version: 3.4-v2
status:
  observedGeneration: 2$4213139756412538772
  phase: Running
```

To get JSON of an object, use `--output=json` flag.

```console
kubedb get mongodb mongodb-demo --output=json
```

To list all KubeDB objects, use following command:

```console
$ kubedb get kubedb -o wide
NAME                VERSION     STATUS  AGE
mg/mongodb-demo     3.4         Running 3h
mg/mongodb-dev      3.4         Running 3h
mg/mongodb-prod     3.4         Running 3h
mg/mongodb-qa       3.4         Running 3h

NAME                                DATABASE                BUCKET              STATUS      AGE
snap/mongodb-demo-20170605-073557   mg/mongodb-demo         gs:bucket-name      Succeeded   9m
snap/snapshot-20171212-114700       mg/mongodb-demo         gs:bucket-name      Succeeded   1h
```

Flag `--output=wide` is used to print additional information.

List command supports short names for each object types. You can use it like `kubedb get <short-name>`. Below are the short name for KubeDB objects:

- MongoDB: `mg`
- Snapshot: `snap`
- DormantDatabase: `drmn`

You can print labels with objects. The following command will list all Snapshots with their corresponding labels.

```console
$ kubedb get snap --show-labels
NAME                            DATABASE                STATUS      AGE       LABELS
mongodb-demo-20170605-073557    mg/mongodb-demo         Succeeded   11m       kubedb.com/kind=MongoDB,kubedb.com/name=mongodb-demo
snapshot-20171212-114700        mg/mongodb-demo         Succeeded   1h        kubedb.com/kind=MongoDB,kubedb.com/name=mongodb-demo
```

You can also filter list using `--selector` flag.

```console
$ kubedb get snap --selector='kubedb.com/kind=MongoDB' --show-labels
NAME                            DATABASE           STATUS      AGE       LABELS
mongodb-demo-20171212-073557    mg/mongodb-demo    Succeeded   14m       kubedb.com/kind=MongoDB,kubedb.com/name=mongodb-demo
snapshot-20171212-114700        mg/mongodb-demo    Succeeded   2h        kubedb.com/kind=MongoDB,kubedb.com/name=mongodb-demo
```

To print only object name, run the following command:

```console
$ kubedb get all -o name
mongodb/mongodb-demo
mongodb/mongodb-dev
mongodb/mongodb-prod
mongodb/mongodb-qa
snapshot/mongodb-demo-20170605-073557
snapshot/snapshot-20170505-114700
```

To learn about various options of `get` command, please visit [here](/docs/reference/kubedb_get.md).

### How to Describe Objects

`kubedb describe` command allows users to describe any KubeDB object. The following command will describe MongoDB database `mongodb-demo` with relevant information.

```console
$ kubedb describe mg mongodb-demo
Name:               mongodb-demo
Namespace:          default
CreationTimestamp:  Wed, 06 Feb 2019 16:31:04 +0600
Labels:             <none>
Annotations:        <none>
Replicas:           1  total
Status:             Running
  StorageType:      Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:
  Name:               mongodb-demo
  CreationTimestamp:  Wed, 06 Feb 2019 16:31:05 +0600
  Labels:               kubedb.com/kind=MongoDB
                        kubedb.com/name=mongodb-demo
  Annotations:        <none>
  Replicas:           824639727120 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         mongodb-demo
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mongodb-demo
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.96.245.200
  Port:         db  27017/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.8:27017

Service:
  Name:         mongodb-demo-gvr
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mongodb-demo
  Annotations:    service.alpha.kubernetes.io/tolerate-unready-endpoints=true
  Type:         ClusterIP
  IP:           None
  Port:         db  27017/TCP
  TargetPort:   27017/TCP
  Endpoints:    172.17.0.8:27017

Database Secret:
  Name:         mongodb-demo-auth
  Labels:         kubedb.com/kind=MongoDB
                  kubedb.com/name=mongodb-demo
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  password:  16 bytes
  username:  4 bytes

No Snapshots.

Events:
  Type    Reason      Age   From             Message
  ----    ------      ----  ----             -------
  Normal  Successful  2m    KubeDB operator  Successfully created Service
  Normal  Successful  2m    KubeDB operator  Successfully created StatefulSet
  Normal  Successful  2m    KubeDB operator  Successfully created MongoDB
  Normal  Successful  2m    KubeDB operator  Successfully created appbinding
  Normal  Successful  2m    KubeDB operator  Successfully patched StatefulSet
  Normal  Successful  2m    KubeDB operator  Successfully patched MongoDB
```

`kubedb describe` command provides following basic information about a MongoDB database.

- StatefulSet
- Storage (Persistent Volume)
- Service
- Secret (If available)
- Snapshots (If any)
- Monitoring system (If available)

To hide events on KubeDB object, use flag `--show-events=false`

To describe all MongoDB objects in `default` namespace, use following command

```console
kubedb describe mg
```

To describe all MongoDB objects from every namespace, provide `--all-namespaces` flag.

```console
kubedb describe mg --all-namespaces
```

To describe all KubeDB objects from every namespace, use the following command:

```console
kubedb describe all --all-namespaces
```

You can also describe KubeDb objects with matching labels. The following command will describe all MongoDB objects with specified labels from every namespace.

```console
kubedb describe mg --all-namespaces --selector='group=dev'
```

To learn about various options of `describe` command, please visit [here](/docs/reference/kubedb_describe.md).

### How to Edit Objects

`kubedb edit` command allows users to directly edit any KubeDB object. It will open the editor defined by _KUBEDB_EDITOR_, or _EDITOR_ environment variables, or fall back to `vim`.

Lets edit an existing running MongoDB object to setup [Scheduled Backup](/docs/guides/mongodb/snapshot/scheduled-backup.md). The following command will open MongoDB `mongodb-demo` in editor.

```console
$ kubedb edit mg mongodb-demo

# Add following under Spec to configure periodic backups
#  backupSchedule:
#    cronExpression: '@every 1m'
#    storageSecretName: mg-snap-secret
#    gcs:
#      bucket: bucket-name

mongodb "mongodb-demo" edited
```

#### Edit Restrictions

Various fields of a KubeDb object can't be edited using `edit` command. The following fields are restricted from updates for all KubeDB objects:

- apiVersion
- kind
- metadata.name
- metadata.namespace

If StatefulSets exists for a MongoDB database, following fields can't be modified as well.

- spec.ReplicaSet
- spec.databaseSecret
- spec.init
- spec.storageType
- spec.storage
- spec.podTemplate.spec.nodeSelector

For DormantDatabase, `spec.origin` can't be edited using `kubedb edit`

To learn about various options of `edit` command, please visit [here](/docs/reference/kubedb_edit.md).

### How to Delete Objects

`kubedb delete` command will delete an object in `default` namespace by default unless namespace is provided. The following command will delete a MongoDB `mongodb-dev` in default namespace

```console
$ kubedb delete mongodb mongodb-dev
mongodb.kubedb.com "mongodb-dev" deleted
```

You can also use YAML files to delete objects. The following command will delete a mongodb using the type and name specified in `mongodb.yaml`.

```console
$ kubedb delete -f mongodb-demo.yaml
mongodb.kubedb.com "mongodb-dev" deleted
```

`kubedb delete` command also takes input from `stdin`.

```console
cat mongodb-demo.yaml | kubedb delete -f -
```

To delete database with matching labels, use `--selector` flag. The following command will delete mongodb with label `mongodb.kubedb.com/name=mongodb-demo`.

```console
kubedb delete mongodb -l mongodb.kubedb.com/name=mongodb-demo
```

To learn about various options of `delete` command, please visit [here](/docs/reference/kubedb_delete.md).

## Using Kubectl

You can use Kubectl with KubeDB objects like any other CRDs. Below are some common examples of using Kubectl with KubeDB objects.

```console
# Create objects
$ kubectl create -f

# List objects
$ kubectl get mongodb
$ kubectl get mongodb.kubedb.com

# Delete objects
$ kubectl delete mongodb <name>
```

## Next Steps

- Learn how to use KubeDB to run a MongoDB database [here](/docs/guides/mongodb/README.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
