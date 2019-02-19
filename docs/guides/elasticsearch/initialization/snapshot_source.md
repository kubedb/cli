---
title: Initialize Elasticsearch from Snapshot
menu:
  docs_0.9.0:
    identifier: es-snapshot-source-initialization
    name: Using Snapshot
    parent: es-initialization-elasticsearch
    weight: 10
menu_name: docs_0.9.0
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
namespace/demo created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: YAML files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Prepare Snapshot

This tutorial will show you how to use KubeDB to initialize an Elasticsearch database with an existing Snapshot. So, we need a Snapshot to perform this initialization. If you don't have a Snapshot already, create one by following the tutorial [here](/docs/guides/elasticsearch/snapshot/instant_backup.md).

If you have changed either namespace or snapshot object name, please modify the YAMLs used in this tutorial accordingly.

## Initialize with Snapshot source

You have to specify the Snapshot `name` and `namespace` in the `spec.init.snapshotSource` field of your new Elasticsearch object.

Below is the YAML for Elasticsearch object that will be created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: recovered-es
  namespace: demo
spec:
  version: "6.3-v1"
  databaseSecret:
    secretName: infant-elasticsearch-auth
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
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
$ kubectl get snap -n demo instant-snapshot
NAME               DATABASENAME           STATUS      AGE
instant-snapshot   infant-elasticsearch   Succeeded   51m
```

> Note: Elasticsearch `recovered-es` must have same superuser credentials as Elasticsearch `infant-elasticsearch`.

[//]: # (Describe authentication part. This should match with existing one)

Now, create the Elasticsearch object.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/initialization/recovered-es.yaml
elasticsearch.kubedb.com/recovered-es created
```

When Elasticsearch database is ready, KubeDB operator launches a Kubernetes Job to initialize this database using the data from Snapshot `instant-snapshot`.

As a final step of initialization, KubeDB Job controller adds `kubedb.com/initialized` annotation in initialized Elasticsearch object. This prevents further invocation of initialization process.

```console
$ kubedb describe es -n demo recovered-es
Name:               recovered-es
Namespace:          demo
CreationTimestamp:  Mon, 08 Oct 2018 12:37:19 +0600
Labels:             <none>
Annotations:        kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"kubedb.com/v1alpha1","kind":"Elasticsearch","metadata":{"annotations":{},"name":"recovered-es","namespace":"demo"},"spec":{"databaseSecr...
                    kubedb.com/initialized
Status:             Running
Replicas:           1  total
Init:
  snapshotSource:
    namespace:  demo
    name:       instant-snapshot
  StorageType:  Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:          
  Name:               recovered-es
  CreationTimestamp:  Mon, 08 Oct 2018 12:37:21 +0600
  Labels:               kubedb.com/kind=Elasticsearch
                        kubedb.com/name=recovered-es
                        node.role.client=set
                        node.role.data=set
                        node.role.master=set
  Annotations:        <none>
  Replicas:           824638233976 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         recovered-es
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=recovered-es
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.104.209.94
  Port:         http  9200/TCP
  TargetPort:   http/TCP
  Endpoints:    192.168.1.14:9200

Service:        
  Name:         recovered-es-master
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=recovered-es
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.110.233.136
  Port:         transport  9300/TCP
  TargetPort:   transport/TCP
  Endpoints:    192.168.1.14:9300

Database Secret:
  Name:         infant-elasticsearch-auth
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=infant-elasticsearch
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  ADMIN_PASSWORD:         8 bytes
  ADMIN_USERNAME:         5 bytes
  sg_action_groups.yml:   430 bytes
  sg_internal_users.yml:  156 bytes
  sg_roles.yml:           312 bytes
  sg_roles_mapping.yml:   73 bytes
  READALL_PASSWORD:       8 bytes
  READALL_USERNAME:       7 bytes
  sg_config.yml:          242 bytes

Certificate Secret:
  Name:         recovered-es-cert
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=recovered-es
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  sgadmin.jks:  3011 bytes
  key_pass:     6 bytes
  node.jks:     3008 bytes
  root.jks:     864 bytes

Topology:
  Type                Pod             StartTime                      Phase
  ----                ---             ---------                      -----
  master|client|data  recovered-es-0  2018-10-08 12:37:22 +0600 +06  Running

No Snapshots.

Events:
  Type    Reason                Age   From                    Message
  ----    ------                ----  ----                    -------
  Normal  Successful            35m   Elasticsearch operator  Successfully created Service
  Normal  Successful            35m   Elasticsearch operator  Successfully created Service
  Normal  Successful            35m   Elasticsearch operator  Successfully created StatefulSet
  Normal  Successful            34m   Elasticsearch operator  Successfully created Elasticsearch
  Normal  Initializing          34m   Elasticsearch operator  Initializing from Snapshot: "instant-snapshot"
  Normal  Successful            34m   Elasticsearch operator  Successfully patched StatefulSet
  Normal  Successful            34m   Elasticsearch operator  Successfully patched Elasticsearch
  Normal  SuccessfulInitialize  33m   Job Controller          Successfully completed initialization
```

## Verify initialization

Let's connect to our Elasticsearch `recovered-es` to verify that the database has been successfully initialized.

At first, forward `9200` port of `recovered-es` pod. Run following command on a separate terminal,

```console
$ kubectl port-forward -n demo recovered-es-0 9200
Forwarding from 127.0.0.1:9200 -> 9200
Forwarding from [::1]:9200 -> 9200
```

Now, we can connect to the database at `localhost:9200`. Let's find out necessary connection information first.

**Connection information:**

- Address: `localhost:9200`
- Username: Run following command to get *username*

  ```console
  $ kubectl get secrets -n demo infant-elasticsearch-auth -o jsonpath='{.data.\ADMIN_USERNAME}' | base64 -d
  admin
  ```

- Password: Run following command to get *password*

  ```console
  $ kubectl get secrets -n demo infant-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
  cfgn547j
  ```

We had created an index `test` before taking snapshot of `infant-elasticsearch` database. Let's check this index is present in newly initialized database `recovered-es`.

```console
$ curl -XGET --user "admin:cfgn547j" "localhost:9200/test/snapshot/1?pretty"
```

```json
{
  "_index" : "test",
  "_type" : "snapshot",
  "_id" : "1",
  "_version" : 33,
  "found" : true,
  "_source" : {
    "title" : "Snapshot",
    "text" : "Testing instand backup",
    "date" : "2018/02/13"
  }
}
```

We can see from above output that `test` index is present in `recovered-es` database. That's means our database has been initialized from snapshot successfully.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/infant-elasticsearch es/recovered-es -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo es/infant-elasticsearch es/recovered-es

$ kubectl delete ns demo
```

## Next Steps

- Learn how to [schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Learn how to configure [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Use [private Docker registry](/docs/guides/elasticsearch/private-registry/using-private-registry.md) to deploy Elasticsearch with KubeDB.
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
