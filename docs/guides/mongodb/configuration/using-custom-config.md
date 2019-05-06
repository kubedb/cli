---
title: Run MongoDB with Custom Configuration
menu:
  docs_0.12.0:
    identifier: mg-custom-config-quickstart
    name: Quickstart
    parent: mg-configuration
    weight: 10
menu_name: docs_0.12.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Custom Configuration File

KubeDB supports providing custom configuration for MongoDB. This tutorial will show you how to use KubeDB to run a MongoDB database with custom configuration.

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace/demo created
  ```

> Note: The yaml files used in this tutorial are stored in [docs/examples/mongodb](https://github.com/kubedb/cli/tree/master/docs/examples/mongodb) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Overview

MongoDB allows to configure database via configuration file. The default configuration file for MongoDB deployed by `KubeDB` can be found in `/data/configdb/mongod.conf`. When MongoDB starts, it will look for custom configuration file in `/configdb-readonly/mongod.conf`. If configuration file exist, this custom configuration will overwrite the existing default one.

> To learn available configuration option of MongoDB see [Configuration File Options](https://docs.mongodb.com/manual/reference/configuration-options/).

At first, you have to create a config file named `mongod.conf`. Then you have to put this file into a [volume](https://kubernetes.io/docs/concepts/storage/volumes/). You have to specify this volume  in `spec.configSource` section while creating MongoDB crd. KubeDB will mount this volume into `/configdb-readonly/` directory of the database pod.

In this tutorial, we will configure [net.maxIncomingConnections](https://docs.mongodb.com/manual/reference/configuration-options/#net.maxIncomingConnections) (default value: 65536) via a custom config file. We will use configMap as volume source.

## Custom Configuration

At first, create `mongod.conf` file containing required configuration settings.

```ini
$ cat mongod.conf
net:
   maxIncomingConnections: 10000
```

Here, `maxIncomingConnections` is set to `10000`, whereas the default value is 65536.

Now, create a configMap with this configuration file.

```console
$ kubectl create configmap -n demo mg-custom-config --from-file=./mongod.conf
configmap/mg-custom-config created
```

Verify the config map has the configuration file.

```yaml
$  kubectl get configmap -n demo mg-custom-config -o yaml
apiVersion: v1
data:
  mongod.conf: |+
    net:
       maxIncomingConnections: 10000
kind: ConfigMap
metadata:
  creationTimestamp: "2019-02-06T10:03:45Z"
  name: mg-custom-config
  namespace: demo
  resourceVersion: "91905"
  selfLink: /api/v1/namespaces/demo/configmaps/mg-custom-config
  uid: 7da0467c-29f6-11e9-aebf-080027875192
```

Now, create MongoDB crd specifying `spec.configSource` field.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-custom-config
  namespace: demo
spec:
  version: "3.4-v3"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  configSource:
      configMap:
        name: mg-custom-config
```

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.12.0/docs/examples/mongodb/configuration/demo-1.yaml
mongodb.kubedb.com/mgo-custom-config created
```

Now, wait a few minutes. KubeDB operator will create necessary PVC, statefulset, services, secret etc. If everything goes well, we will see that a pod with the name `mgo-custom-config-0` has been created.

Check that the statefulset's pod is running

```console
$ kubectl get pod -n demo mgo-custom-config-0
NAME                  READY     STATUS    RESTARTS   AGE
mgo-custom-config-0   1/1       Running   0          1m
```

Now, we will check if the database has started with the custom configuration we have provided.

Now, you can connect to this database through [mongo-shell](https://docs.mongodb.com/v3.4/mongo/). In this tutorial, we are connecting to the MongoDB server from inside the pod.

```console
$ kubectl get secrets -n demo mgo-custom-config-auth -o jsonpath='{.data.\username}' | base64 -d
root

$ kubectl get secrets -n demo mgo-custom-config-auth -o jsonpath='{.data.\password}' | base64 -d
ErialNojWParBFoP

$ kubectl exec -it mgo-custom-config-0 -n demo sh

> mongo admin

> db.auth("root","ErialNojWParBFoP")
1

> db._adminCommand( {getCmdLineOpts: 1})
{
	"argv" : [
		"mongod",
		"--dbpath=/data/db",
		"--auth",
		"--bind_ip=0.0.0.0",
		"--port=27017",
		"--config=/data/configdb/mongod.conf"
	],
	"parsed" : {
		"config" : "/data/configdb/mongod.conf",
		"net" : {
			"bindIp" : "0.0.0.0",
			"maxIncomingConnections" : 10000,
			"port" : 27017
		},
		"security" : {
			"authorization" : "enabled"
		},
		"storage" : {
			"dbPath" : "/data/db"
		}
	},
	"ok" : 1
}

> exit
bye
```

As we can see from the configuration of running mongodb, the value of `maxIncomingConnections` has been set to 10000 successfully.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mg/mgo-custom-config -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mg/mgo-custom-config

kubectl patch -n demo drmn/mgo-custom-config -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mgo-custom-config

kubectl delete -n demo configmap mg-custom-config

kubectl delete ns demo
```

## Next Steps

- [Snapshot and Restore](/docs/guides/mongodb/snapshot/backup-and-restore.md) process of MongoDB databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mongodb/snapshot/scheduled-backup.md) of MongoDB databases using KubeDB.
- Initialize [MongoDB with Script](/docs/guides/mongodb/initialization/using-script.md).
- Initialize [MongoDB with Snapshot](/docs/guides/mongodb/initialization/using-snapshot.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mongodb/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mongodb/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mongodb/private-registry/using-private-registry.md) to deploy MongoDB with KubeDB.
- Use [kubedb cli](/docs/guides/mongodb/cli/cli.md) to manage databases like kubectl for Kubernetes.
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Detail concepts of [MongoDBVersion object](/docs/concepts/catalog/mongodb.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
