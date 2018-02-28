---
title: Initialize Elasticsearch from Snapshot
menu:
  docs_0.8.0-beta.2:
    identifier: es-snapshot-source-initialization
    name: Using Snapshot
    parent: es-initialization-elasticsearch
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---

> Don't know how backup works?  Check [tutorial](/docs/guides/elasticsearch/snapshot/instant_backup.md) on Instant Backup.

# Initialize Elasticsearch with Snapshot

KubeDB supports Elasticsearch database initialization.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: Yaml files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

This tutorial will show you how to use KubeDB to initialize a Elasticsearch database with an existing Snapshot.

So, we need a Snapshot object in `Succeeded` phase to perform this initialization .

Follow these steps to prepare this tutorial

- Create Elasticsearch object `infant-elasticsearch`, if not exists.

    ```console
    $ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/elasticsearch/quickstart/infant-elasticsearch.yaml
    validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/elasticsearch/quickstart/infant-elasticsearch.yaml"
    elasticsearch "infant-elasticsearch" created
    ```

    ```console
    $ kubedb get es -n demo infant-elasticsearch
    NAME                   STATUS    AGE
    infant-elasticsearch   Running   57s
    ```

- Populate database with some data. Follow [this](https://github.com/kubedb/cli/blob/master/docs/guides/elasticsearch/snapshot/instant_backup.md#populate-database).
- Create storage Secret.<br>In this tutorial, we need a storage Secret for backup process

    ```console
    $ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
    $ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
    $ kubectl create secret -n demo generic gcs-secret \
        --from-file=./GOOGLE_PROJECT_ID \
        --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
    secret "gcs-secret" created
    ```

- Take an instant backup, if not available. Follow [this](https://github.com/kubedb/cli/blob/master/docs/guides/elasticsearch/snapshot/instant_backup.md#instant-backup).

```console
$ kubedb get snap -n demo --selector="kubedb.com/kind=Elasticsearch,kubedb.com/name=infant-elasticsearch"
NAME               DATABASE                  STATUS      AGE
instant-snapshot   es/infant-elasticsearch   Succeeded   39s
```

## Initialize with Snapshot source

Specify the Snapshot `name` and `namespace` in the `spec.init.snapshotSource` field of your new Elasticsearch object.

See the example Elasticsearch object below

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: recovered-es
  namespace: demo
spec:
  version: 5.6
  databaseSecret:
    secretName: infant-elasticsearch-auth
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    snapshotSource:
      name: instant-snapshot
      namespace: demo
```

Here,

- `spec.init.snapshotSource` specifies Snapshot object information to be used in this initialization process.
	- `snapshotSource.name` refers to a Snapshot object `name`.
	- `snapshotSource.namespace` refers to a Snapshot object `namespace`.

Snapshot `instant-snapshot` in `demo` namespace belongs to Elasticsearch `infant-elasticsearch`:

```console
$ kubedb get snap -n demo instant-snapshot
NAME               DATABASE                  STATUS      AGE
instant-snapshot   es/infant-elasticsearch   Succeeded   2m
```

> Note: Elasticsearch `recovered-es` must have same `admin` user password as Elasticsearch `infant-elasticsearch`.

[//]: # (Describe authentication part. This should match with existing one)

Now, create the Elasticsearch object.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/elasticsearch/initialization/recovered-es.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/elasticsearch/initialization/recovered-es.yaml"
elasticsearch "recovered-es" created
```

When Elasticsearch database is ready, KubeDB operator launches a Kubernetes Job to initialize this database using the data from Snapshot `instant-snapshot`.

As a final step of initialization, KubeDB Job controller adds `kubedb.com/initialized` annotation in initialized Elasticsearch object.
This prevents further invocation of initialization process.

```console
$ kubedb describe es -n demo recovered-es -S=false -W=false
Name:			    recovered-es
Namespace:		    demo
CreationTimestamp:  Wed, 14 Feb 2018 17:31:14 +0600
Status:			    Running
Annotations:		kubedb.com/initialized
Init:
  snapshotSource:
    namespace:	demo
    name:	    instant-snapshot
Volume:
  StorageClass:	standard
  Capacity:	    50Mi
  Access Modes:	RWO
StatefulSet:	recovered-es
Service:	    recovered-es, recovered-es-master
Secrets:	    recovered-es-cert, infant-elasticsearch-auth

Topology:
  Type                 Pod              StartTime                       Phase
  ----                 ---              ---------                       -----
  master|client|data   recovered-es-0   2018-02-14 17:31:23 +0600 +06   Running

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason               Message
  ---------   --------   -----     ----                     --------   ------               -------
  22m         22m        1         Job Controller           Normal     SuccessfulSnapshot   Successfully completed initialization
  23m         23m        1         Elasticsearch operator   Normal     Successful           Successfully patched Elasticsearch
  24m         24m        1         Elasticsearch operator   Normal     Successful           Successfully patched StatefulSet
  24m         24m        1         Elasticsearch operator   Normal     Initializing         Initializing from Snapshot: "instant-snapshot"
  24m         24m        1         Elasticsearch operator   Normal     Successful           Successfully created Elasticsearch
  25m         25m        1         Elasticsearch operator   Normal     Successful           Successfully created StatefulSet
  25m         25m        1         Elasticsearch operator   Normal     Successful           Successfully created Service
  25m         25m        1         Elasticsearch operator   Normal     Successful           Successfully created Service
```

In this tutorial, we will expose ClusterIP Service `recovered-es` to connect database from local.

```console
$ kubectl expose svc -n demo recovered-es --name=recovered-es-exposed --port=9200 --protocol=TCP --type=NodePort
service "recovered-es-exposed" exposed
```

Now lets check data in Elasticsearch `recovered-es` using Service `recovered-es-exposed`

Connection information:

- address: Use Service URL `$ minikube service recovered-es-exposed -n demo --https --url`

Run following command to get `admin` user password

```console
$ kubectl get secrets -n demo infant-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
3pk6qxxo⏎
```

```console
export es_service=$(minikube service recovered-es-exposed -n demo --https --url)
export es_admin_pass=$(kubectl get secrets -n demo infant-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d)
curl -XGET --user "admin:$es_admin_pass" "$es_service/test/snapshot/1?pretty" --insecure
```

```json
{
  "_index" : "test",
  "_type" : "snapshot",
  "_id" : "1",
  "_version" : 1,
  "found" : true,
  "_source" : {
    "title" : "Snapshot",
    "text" : "Testing instand backup",
    "date" : "2018/02/13"
  }
}
```

Elasticsearch `recovered-es` is successfully initialized with Snapshot `instant-snapshot`

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete es,drmn,snap -n demo --all --force

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Learn how to [schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Learn how to configure [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
