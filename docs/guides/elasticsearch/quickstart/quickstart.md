---
title: Elasticsearch Quickstart
menu:
  docs_0.8.0:
    identifier: es-quickstart-quickstart
    name: Overview
    parent: es-quickstart-elasticsearch
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Elasticsearch QuickStart

This tutorial will show you how to use KubeDB to run a Elasticsearch database.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/elasticsearch/lifecycle.png" width="600" height="660">
</p>

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

## Create a Elasticsearch database

KubeDB implements a Elasticsearch CRD to define the specification of a Elasticsearch database.

Below is the Elasticsearch object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: quick-elasticsearch
  namespace: demo
spec:
  version: "5.6"
  doNotPause: true
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
```

Here,

- `spec.version` is the version of Elasticsearch database. In this tutorial, a Elasticsearch 5.6 database is created.
- `spec.doNotPause` prevents user from deleting this object if admission webhook is enabled.
- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet
 created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests.
 If no storage spec is given, an `emptyDir` is used.

Create example above with following command

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/elasticsearch/quickstart/quick-elasticsearch.yaml
elasticsearch "quick-elasticsearch" created
```

KubeDB operator watches for Elasticsearch objects using Kubernetes api. When an Elasticsearch object is created, KubeDB operator creates a new StatefulSet and two ClusterIP Service with the matching name.
KubeDB operator will also create a governing service for StatefulSet with the name `kubedb`, if one is not already present.

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created.

```console
$ kubedb get es -n demo quick-elasticsearch -o wide
NAME                  VERSION   STATUS    AGE
quick-elasticsearch   5.6       Running   33m
```

Lets describe Elasticsearch object `quick-elasticsearch`

```console
$ kubedb describe es -n demo quick-elasticsearch
Name:               quick-elasticsearch
Namespace:          demo
CreationTimestamp:  Mon, 19 Feb 2018 16:10:45 +0600
Status:             Running
Volume:
  StorageClass: standard
  Capacity:     50Mi
  Access Modes: RWO

StatefulSet:
  Name:                 quick-elasticsearch
  Replicas:             1 current / 1 desired
  CreationTimestamp:    Mon, 19 Feb 2018 16:10:55 +0600
  Pods Status:          1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		quick-elasticsearch
  Type:		ClusterIP
  IP:		10.11.255.214
  Port:		http	9200/TCP

Service:
  Name:		quick-elasticsearch-master
  Type:		ClusterIP
  IP:		10.11.255.170
  Port:		transport	9300/TCP

Certificate Secret:
  Name:	quick-elasticsearch-cert
  Type:	Opaque
  Data
  ====
  key_pass:     6 bytes
  node.jks:     3013 bytes
  root.jks:     864 bytes
  sgadmin.jks:  3009 bytes

Database Secret:
  Name:	quick-elasticsearch-auth
  Type:	Opaque
  Data
  ====
  ADMIN_PASSWORD:           8 bytes
  READALL_PASSWORD:         8 bytes
  sg_action_groups.yml:     430 bytes
  sg_config.yml:            240 bytes
  sg_internal_users.yml:    156 bytes
  sg_roles.yml:             312 bytes
  sg_roles_mapping.yml:     73 bytes

Topology:
  Type                 Pod                     StartTime                       Phase
  ----                 ---                     ---------                       -----
  master|client|data   quick-elasticsearch-0   2018-02-19 16:11:02 +0600 +06   Running

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                     Type       Reason       Message
  ---------   --------   -----     ----                     --------   ------       -------
  10s         10s        1         Elasticsearch operator   Normal     Successful   Successfully patched Elasticsearch
  40s         40s        1         Elasticsearch operator   Normal     Successful   Successfully patched StatefulSet
  48s         48s        1         Elasticsearch operator   Normal     Successful   Successfully created Elasticsearch
  1m          1m         1         Elasticsearch operator   Normal     Successful   Successfully created StatefulSet
  1m          1m         1         Elasticsearch operator   Normal     Successful   Successfully created Service
  1m          1m         1         Elasticsearch operator   Normal     Successful   Successfully created Service
```

```console
$ kubectl get service -n demo --selector=kubedb.com/kind=Elasticsearch,kubedb.com/name=quick-elasticsearch
NAME                         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
quick-elasticsearch          ClusterIP   10.97.171.52     <none>        9200/TCP   22m
quick-elasticsearch-master   ClusterIP   10.105.209.152   <none>        9300/TCP   22m
```

Two services for each Elasticsearch object.

- Service *`quick-elasticsearch`* targets all Pods which are acting as *client* node
- Service *`quick-elasticsearch-master`* targets all Pods which are acting as *master* node

KubeDB supports Elasticsearch clustering where Pod can be any of these three role: *master*, *data* or *client*.

If you see `Topology` section in `kubedb describe` result, you will know role(s) of each Pod.

```console
Topology:
  Type                 Pod                     StartTime                       Phase
  ----                 ---                     ---------                       -----
  data|master|client   quick-elasticsearch-0   2018-02-12 11:56:47 +0600 +06   Running
```

Here, we create a Elasticsearch database with single node. This single node will act as *master*, *data* and *client*.

To learn how to configure Elasticsearch cluster, click [here](/docs/guides/elasticsearch/clustering/topology.md).

Please note that KubeDB operator has created two new Secrets for Elasticsearch object.

1. `quick-elasticsearch-auth` for storing the passwords and [search-guard](https://github.com/floragunncom/search-guard) configuration.
2. `quick-elasticsearch-cert` for storing certificates used for SSL connection.

##### Secret for authentication & configuration

```console
$ kubectl get secret -n demo quick-elasticsearch-auth -o yaml
```

```yaml
apiVersion: v1
data:
  ADMIN_PASSWORD: MmJpaGQ1NGc=
  READALL_PASSWORD: YTJkcDZjamc=
  sg_action_groups.yml: ClVOTElNSVRFRDoKICAtICIqIgoKUk...AiaW5kaWNlczphZG1pbi9nZXQiCg==
  sg_config.yml: CnNlYXJjaGd1YXJkOgogIGR5bmFtaW...AgICAgICAgICB0eXBlOiBpbnRlcm4K
  sg_internal_users.yml: CmFkbWluOgogIGhhc2g6ICQyYSQxMC...JEdkxVSzZrUE1xT1hJVTZYbnN0OWEK
  sg_roles.yml: CnNnX2FsbF9hY2Nlc3M6CiAgY2x1c3...5ESUNFU19LVUJFREJfU05BUFNIT1QK
  sg_roles_mapping.yml: CnNnX2FsbF9hY2Nlc3M6CiAgdXNlcn...VzZXJzOgogICAgLSByZWFkYWxsCg==
kind: Secret
metadata:
  creationTimestamp: 2018-02-14T08:24:05Z
  labels:
    kubedb.com/kind: Elasticsearch
    kubedb.com/name: quick-elasticsearch
  name: quick-elasticsearch-auth
  namespace: demo
  resourceVersion: "4376"
  selfLink: /api/v1/namespaces/demo/secrets/quick-elasticsearch-auth
  uid: 6baf55bd-1160-11e8-a344-08002716e6a0
type: Opaque
```

> Note: Auth Secret name format: `{elasticsearch-name}-auth`

This Secret contains:

- `ADMIN_PASSWORD` password for `admin` user used in search-guard configuration as internal user.
- `READALL_PASSWORD` password for `readall` user with read-only permission only.
- Followings are used as search-guard configuration
  - `sg_action_groups.yml`
  - `sg_config.yml`
  - `sg_internal_users.yml`
  - `sg_roles.yml`
  - `sg_roles_mapping.yml`

See details about [search-guard configuration](/docs/guides/elasticsearch/search-guard/configuration.md)

##### Secret for certificates

```console
$ kubectl get secret -n demo quick-elasticsearch-cert -o yaml
```

```yaml
apiVersion: v1
data:
  key_pass: b2xxeHN1
  node.jks: /u3+7QAAAAIAAAABAAAA...A0+i8Kj9XQUo1V/Qg==
  root.jks: /u3+7QAAAAIAAAABAAAA...tBkRsCa+uTUYjiatf7j
  sgadmin.jks: /u3+7QAAAAIAAAABAAAA...e5h2S9Y3e429E/9P1qw
kind: Secret
metadata:
  creationTimestamp: 2018-02-19T10:10:53Z
  labels:
    kubedb.com/kind: Elasticsearch
    kubedb.com/name: quick-elasticsearch
  name: quick-elasticsearch-cert
  namespace: demo
  resourceVersion: "1778"
  selfLink: /api/v1/namespaces/demo/secrets/quick-elasticsearch-cert
  uid: 2b5abae5-155d-11e8-a001-42010a8000d5
type: Opaque
```

> Note: Cert Secret name format: `{elasticsearch-name}-cert`

This Secret contains SSL certificates. See details about [SSL certificates](/docs/guides/elasticsearch/search-guard/search_guard.md) needed for Elasticsearch.

#### Connect Elasticsearch

In this tutorial, we will expose ClusterIP Service `quick-elasticsearch` to connect database from local.

```console
$ kubectl expose svc -n demo quick-elasticsearch --name=quick-es-exposed --port=9200 --protocol=TCP --type=NodePort
service "quick-es-exposed" exposed
```

```console
$ kubectl get svc -n demo quick-es-exposed
NAME               TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
quick-es-exposed   NodePort   10.110.247.179   <none>        9200:32403/TCP   1m
```

Following will provide URL for `quick-es-exposed` Service to access Elasticsearch database.

```console
$ minikube service quick-es-exposed -n demo --url
http://192.168.99.100:32403
```

Now, you can connect to this database using curl.

Connection information:

- address: Use Service URL `$ minikube service quick-es-exposed -n demo --url`

Run following command to get `admin` user password

```console
$ kubectl get secrets -n demo quick-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
ud7cagcu⏎
```

Check health of the elasticsearch database

```console
export es_service=$(minikube service quick-es-exposed -n demo --url)
export es_admin_pass=$(kubectl get secrets -n demo quick-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d)
curl --user "admin:$es_admin_pass" "$es_service/_cluster/health?pretty"
```

```json
{
  "cluster_name" : "quick-elasticsearch",
  "status" : "green",
  "timed_out" : false,
  "number_of_nodes" : 1,
  "number_of_data_nodes" : 1,
  "active_primary_shards" : 1,
  "active_shards" : 1,
  "relocating_shards" : 0,
  "initializing_shards" : 0,
  "unassigned_shards" : 0,
  "delayed_unassigned_shards" : 0,
  "number_of_pending_tasks" : 0,
  "number_of_in_flight_fetch" : 0,
  "task_max_waiting_in_queue_millis" : 0,
  "active_shards_percent_as_number" : 100.0
}
```

## Pause Elasticsearch

KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `doNotPause` feature. If admission webhook is enabled,
It prevents user from deleting the database as long as the `spec.doNotPause` is set `true`.

In this tutorial, Elasticsearch `quick-elasticsearch` is created with `spec.doNotPause: true`. So, if you delete this Elasticsearch object,
admission webhook will nullify the delete operation.

```console
$ kubedb delete es -n demo quick-elasticsearch
error: Elasticsearch "quick-elasticsearch " can't be paused. To continue delete, unset spec.doNotPause and retry.
```

To continue with this tutorial, unset `spec.doNotPause` by updating Elasticsearch object

```console
$ kubedb edit es -n demo quick-elasticsearch
spec:
  doNotPause: false
```

Now, if you delete the Elasticsearch object, KubeDB operator will create a matching DormantDatabase object.
KubeDB operator watches for DormantDatabase objects and it will take necessary steps when a DormantDatabase object is created.

KubeDB operator will delete the StatefulSet and its Pods, but leaves the Secret, PVCs unchanged.

```console
$ kubedb delete es -n demo quick-elasticsearch
elasticsearch "quick-elasticsearch" deleted
```

Check DormantDatabase entry

```console
$ kubedb get drmn -n demo quick-elasticsearch
NAME                  STATUS    AGE
quick-elasticsearch   Paused    1m
```

In KubeDB parlance, we say that Elasticsearch `quick-elasticsearch`  has entered into dormant state.

Lets see, what we have in this DormantDatabase object

```yaml
$ kubedb get drmn -n demo quick-elasticsearch -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-13T05:41:49Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: Elasticsearch
  name: quick-elasticsearch
  namespace: demo
  resourceVersion: "2072"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/quick-elasticsearch
  uid: 9624f877-1080-11e8-9e42-0800271bdbb6
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: quick-elasticsearch
      namespace: demo
    spec:
      elasticsearch:
        certificateSecret:
          secretName: quick-elasticsearch-cert
        databaseSecret:
          secretName: quick-elasticsearch-auth
        resources: {}
        storage:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 50Mi
          storageClassName: standard
        version: "5.6"
status:
  creationTime: 2018-02-13T05:41:49Z
  pausingTime: 2018-02-13T05:42:12Z
  phase: Paused
```

Here,

- `spec.origin` contains original Elasticsearch object.
- `status.phase` points to the current database state `Paused`.

## Resume DormantDatabase

To resume the database from the dormant state, create same Elasticsearch object with same Spec.

In this tutorial, the DormantDatabase `quick-elasticsearch` can be resumed by creating original Elasticsearch object.

The below command will resume the DormantDatabase `quick-elasticsearch`

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0/docs/examples/elasticsearch/quickstart/quick-elasticsearch.yaml
elasticsearch "quick-elasticsearch" created
```

## WipeOut DormantDatabase

You can wipe out a DormantDatabase while deleting the objet by setting `spec.wipeOut` to true. KubeDB operator will delete any relevant resources of this `Elasticsearch` database (i.e, PVCs, Secrets, Snapshots). It will also delete snapshot data stored in the Cloud Storage buckets.

```yaml
$ kubedb edit drmn -n demo quick-elasticsearch
spec:
  wipeOut: true
```

If `spec.wipeOut` is not set to true while deleting the `dormantdatabase` object, then only this object will be deleted and `kubedb-operator` won't delete related Secrets, PVCs and Snapshots. So, user still can access the stored data in the cloud storage buckets as well as PVCs.

## Delete DormantDatabase

As it is already discussed above, `DormantDatabase` can be deleted with or without wiping out the resources. To delete the `dormantdatabase`,

```console
$ kubedb delete drmn -n demo quick-elasticsearch
dormantdatabase "quick-elasticsearch" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/quick-elasticsearch -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo es/quick-elasticsearch

$ kubectl patch -n demo drmn/quick-elasticsearch -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/quick-elasticsearch

$ kubectl delete ns demo
```

## Next Steps

- Learn about [taking instant backup](/docs/guides/elasticsearch/snapshot/instant_backup.md) of Elasticsearch database using KubeDB.
- Learn how to [schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Learn about initializing [Elasticsearch with Snapshot](/docs/guides/elasticsearch/initialization/snapshot_source.md).
- Learn how to configure [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Detail concepts of [Elasticsearch object](/docs/concepts/databases/elasticsearch.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
