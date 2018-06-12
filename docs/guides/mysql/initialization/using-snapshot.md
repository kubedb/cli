---
title: Initialize MySQL from Snapshot
menu:
  docs_0.8.0:
    identifier: my-using-snapshot-initialization
    name: From Snapshot
    parent: my-initialization-mysql
    weight: 15
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Initialize MySQL with Snapshot

This tutorial will show you how to use KubeDB to initialize a MySQL database with an existing snapshot.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

This tutorial assumes that you have created a namespace `demo` and a snapshot `snapshot-infant`. Follow the steps [here](/docs/guides/mysql/snapshot/backup-and-restore.md) to create a database and take [instant snapshot](/docs/guides/mysql/snapshot/backup-and-restore.md#instant-backups), if you have not done so already. If you have changed the name of either namespace or snapshot object, please modify the YAMLs used in this tutorial accordingly.

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create MySQL with Init-Snapshot

Below is the `MySQL` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-init-snapshot
  namespace: demo
spec:
  version: "8.0"
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    snapshotSource:
      name: snap-mysql-infant
      namespace: demo
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/mysql/Initialization/demo-2.yaml
mysql "mysql-init-snapshot" created
```

Here,

- `spec.init.snapshotSource.name` refers to a Snapshot object for a MySQL database in the same namespaces as this new `mysql-init-snapshot` MySQL object.

Now, wait several seconds. KubeDB operator will create a new `StatefulSet`. Then KubeDB operator launches a Kubernetes Job to initialize the new database using the data from `snap-mysql-infant` Snapshot.

```console
$ kubedb get my -n demo
NAME                  STATUS         AGE
mysql-infant          Running        15m
mysql-init-snapshot   Initializing   11s

$ kubedb get my -n demo
NAME                  STATUS    AGE
mysql-infant          Running   17m
mysql-init-snapshot   Running   2m

$ kubedb describe my -n demo mysql-init-snapshot
Name:		mysql-init-snapshot
Namespace:	demo
StartTimestamp:	Mon, 12 Feb 2018 11:09:12 +0600
Status:		Running
Annotations:	kubedb.com/initialized=
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			mysql-init-snapshot
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Mon, 12 Feb 2018 11:09:13 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		mysql-init-snapshot
  Type:		ClusterIP
  IP:		10.105.136.174
  Port:		db	3306/TCP

Database Secret:
  Name:	mysql-init-snapshot-auth
  Type:	Opaque
  Data
  ====
  password:	16 bytes
  user:		4 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From             Type       Reason               Message
  ---------   --------   -----     ----             --------   ------               -------
  19s         19s        1         Job Controller   Normal     SuccessfulSnapshot   Successfully completed initialization
  2m          2m         1         MySQL operator   Normal     Successful           Successfully patched StatefulSet
  2m          2m         1         MySQL operator   Normal     Successful           Successfully patched MySQL
  2m          2m         1         MySQL operator   Normal     Successful           Successfully created StatefulSet
  2m          2m         1         MySQL operator   Normal     Initializing         Initializing from Snapshot: "snap-mysql-infant"
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo mysql/mysql-init-snapshot -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo mysql/mysql-init-snapshot

$ kubectl patch -n demo drmn/mysql-init-snapshot -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/mysql-init-snapshot

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
