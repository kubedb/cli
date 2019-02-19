---
title: Initialize MySQL from Snapshot
menu:
  docs_0.9.0:
    identifier: my-using-snapshot-initialization
    name: From Snapshot
    parent: my-initialization-mysql
    weight: 15
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Initialize MySQL with Snapshot

This tutorial will show you how to use KubeDB to initialize a MySQL database with an existing snapshot.

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- This tutorial assumes that you have created a namespace `demo` and a snapshot `snapshot-infant`. Follow the steps [here](/docs/guides/mysql/snapshot/backup-and-restore.md) to create a database and take [instant snapshot](/docs/guides/mysql/snapshot/backup-and-restore.md#instant-backups), if you have not done so already. If you have changed the name of either namespace or snapshot object, please modify the YAMLs used in this tutorial accordingly.

> Note: The yaml files that are used in this tutorial are stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create MySQL with Init-Snapshot

Below is the `MySQL` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-init-snapshot
  namespace: demo
spec:
  version: "8.0-v2"
  databaseSecret:
    secretName: mysql-infant-auth
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  init:
    snapshotSource:
      name: snap-mysql-infant
      namespace: demo
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/Initialization/demo-2.yaml
mysql.kubedb.com/mysql-init-snapshot created
```

Here,

- `spec.init.snapshotSource.name` refers to a Snapshot object for a MySQL database in the same namespaces as this new `mysql-init-snapshot` MySQL object.

Now, wait several seconds. KubeDB operator will create a new `StatefulSet`. Then KubeDB operator launches a Kubernetes Job to initialize the new database using the data from `snap-mysql-infant` Snapshot.

```console
$ kubedb get my -n demo
NAME                  VERSION   STATUS         AGE
mysql-infant          8.0-v2    Running        8m
mysql-init-snapshot   8.0-v2    Initializing   1m

$ kubedb get my -n demo
NAME                  VERSION   STATUS    AGE
mysql-infant          8.0-v2    Running   20m
mysql-init-snapshot   8.0-v2    Running   13m

$ kubedb describe my -n demo mysql-init-snapshot
Name:               mysql-init-snapshot
Namespace:          demo
CreationTimestamp:  Thu, 27 Sep 2018 17:54:16 +0600
Labels:             <none>
Annotations:        kubedb.com/initialized=
Replicas:           1  total
Status:             Running
  StorageType:      Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:
  Name:               mysql-init-snapshot
  CreationTimestamp:  Thu, 27 Sep 2018 17:54:17 +0600
  Labels:               kubedb.com/kind=MySQL
                        kubedb.com/name=mysql-init-snapshot
  Annotations:        <none>
  Replicas:           824642013116 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         mysql-init-snapshot
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-init-snapshot
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.104.217.79
  Port:         db  3306/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.5:3306

Database Secret:
  Name:         mysql-infant-auth
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-infant
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  password:  16 bytes
  user:      4 bytes

No Snapshots.

Events:
  Type    Reason                Age   From            Message
  ----    ------                ----  ----            -------
  Normal  Successful            13m   MySQL operator  Successfully created Service
  Normal  Successful            12m   MySQL operator  Successfully created MySQL
  Normal  Successful            12m   MySQL operator  Successfully created StatefulSet
  Normal  Initializing          12m   MySQL operator  Initializing from Snapshot: "snap-mysql-infant"
  Normal  Successful            12m   MySQL operator  Successfully patched StatefulSet
  Normal  Successful            12m   MySQL operator  Successfully patched MySQL
  Normal  SuccessfulInitialize  6m    Job Controller  Successfully completed initialization
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mysql/mysql-infant mysql/mysql-init-snapshot -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mysql/mysql-infant mysql/mysql-init-snapshot

kubectl patch -n demo drmn/mysql-infant drmn/mysql-init-snapshot -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mysql-infant drmn/mysql-init-snapshot

kubectl delete ns demo
```

## Next Steps

- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Detail concepts of [MySQLVersion object](/docs/concepts/catalog/mysql.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
