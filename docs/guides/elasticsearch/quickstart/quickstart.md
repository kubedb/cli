---
title: Elasticsearch Quickstart
menu:
  docs_0.9.0:
    identifier: es-quickstart-quickstart
    name: Overview
    parent: es-quickstart-elasticsearch
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Elasticsearch QuickStart

This tutorial will show you how to use KubeDB to run an Elasticsearch database.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/elasticsearch/lifecycle.png">
</p>

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

>We have designed this tutorial to demonstrate a production setup of KubeDB managed Elasticsearch. If you just want to try out KubeDB, you can bypass some of the safety features following the tips [here](/docs/guides/elasticsearch/quickstart/quickstart.md#tips-for-testing).

## Find Available StorageClass

We will have to provide `StorageClass` in Elasticsearch crd specification. Check available `StorageClass` in your cluster using the following command,

```console
$ kubectl get storageclass
NAME       PROVISIONER        AGE
standard   external/pharmer   1m

```

Here, we have `standard` StorageClass in our cluster.

## Find Available ElasticsearchVersion

When you have installed KubeDB, it has created `ElasticsearchVersion` crd for all supported Elasticsearch versions. Let's check available ElasticsearchVersions by,

```console
$ kubectl get elasticsearchversions
NAME       VERSION   DB_IMAGE                        DEPRECATED   AGE
5.6        5.6       kubedb/elasticsearch:5.6        true         21m
5.6-v1     5.6       kubedb/elasticsearch:5.6-v1                  21m
5.6.4      5.6.4     kubedb/elasticsearch:5.6.4      true         21m
5.6.4-v1   5.6.4     kubedb/elasticsearch:5.6.4-v1                21m
6.2        6.2       kubedb/elasticsearch:6.2        true         21m
6.2-v1     6.2       kubedb/elasticsearch:6.2-v1                  21m
6.2.4      6.2.4     kubedb/elasticsearch:6.2.4      true         21m
6.2.4-v1   6.2.4     kubedb/elasticsearch:6.2.4-v1                21m
6.3        6.3       kubedb/elasticsearch:6.3        true         21m
6.3-v1     6.3       kubedb/elasticsearch:6.3-v1                  21m
6.3.0      6.3.0     kubedb/elasticsearch:6.3.0      true         21m
6.3.0-v1   6.3.0     kubedb/elasticsearch:6.3.0-v1                20m
```

Notice the `DEPRECATED` column. Here, `true` means that this ElasticsearchVersion is deprecated for current KubeDB version. KubeDB will not work for deprecated ElasticsearchVersion.

In this tutorial, we will use `6.3-v1` ElasticsearchVersion crd to create Elasticsearch database. To know more about what is `ElasticsearchVersion` crd and why there is `6.3` and `6.3-v1` variation, please visit [here](/docs/concepts/catalog/elasticsearch.md). You can also see supported ElasticsearchVersion in KubeDB 0.9.0 from [here](/docs/guides/elasticsearch/README.md#supported-elasticsearchversion-crd).

## Create an Elasticsearch database

KubeDB implements an Elasticsearch CRD to define the specification of an Elasticsearch database.

Below is the Elasticsearch object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: quick-elasticsearch
  namespace: demo
spec:
  version: "6.3-v1"
  storageType: Durable
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  terminationPolicy: DoNotTerminate
```

Here,

- `spec.version` is name of the ElasticsearchVersion crd. In this tutorial, an Elasticsearch 6.3 database is created.
- `spec.storageType` specifies the type of storage that will be used for Elasticsearch database. It can be `Durable` or `Ephemeral`. The default value of this field is `Durable`. If `Ephemeral` is used then KubeDB will create Elasticsearch database using `EmptyDir` volume. In this case, you don't have to specify `spec.storage` field. This is useful for testing purpose.
- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If you don't specify `spec.storageType: Ephemeral`, then this field is required.
- `spec.terminationPolicy` specifies what KubeDB should do when user try to delete Elasticsearch crd. Termination policy `DoNotTerminate` prevents a user from deleting this object if admission webhook is enabled.

>Note: `spec.storage` section is used to create PVC for database pod. It will create PVC with storage size specified in`storage.resources.requests` field. Don't specify `limits` here. PVC does not get resized automatically.

Let's create Elasticsearch crd that is shown above with following command

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/quickstart/quick-elasticsearch.yaml
elasticsearch.kubedb.com/quick-elasticsearch created
```

KubeDB operator watches for Elasticsearch objects using Kubernetes api. When an Elasticsearch object is created, KubeDB operator creates a new StatefulSet and two ClusterIP Service with the matching name.

KubeDB operator will also create a governing service for StatefulSet with the name `kubedb`, if one is not already present.

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created.

```console
$ kubectl get es -n demo quick-elasticsearch
NAME                  VERSION   STATUS    AGE
quick-elasticsearch   6.3-v1    Running   3m
```

Let's describe Elasticsearch object `quick-elasticsearch`

```console
$ kubedb describe es -n demo quick-elasticsearch
Name:               quick-elasticsearch
Namespace:          demo
CreationTimestamp:  Fri, 28 Sep 2018 11:33:29 +0600
Labels:             <none>
Annotations:        kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"kubedb.com/v1alpha1","kind":"Elasticsearch","metadata":{"annotations":{},"name":"quick-elasticsearch","namespace":"demo"},"spec":{"doNot...
Status:             Running
Replicas:           1  total
  StorageType:      Durable
Volume:
  StorageClass:  standard
  Capacity:      1Gi
  Access Modes:  RWO

StatefulSet:
  Name:               quick-elasticsearch
  CreationTimestamp:  Fri, 28 Sep 2018 11:33:36 +0600
  Labels:               kubedb.com/kind=Elasticsearch
                        kubedb.com/name=quick-elasticsearch
                        node.role.client=set
                        node.role.data=set
                        node.role.master=set
  Annotations:        <none>
  Replicas:           824640716856 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         quick-elasticsearch
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=quick-elasticsearch
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.100.103.159
  Port:         http  9200/TCP
  TargetPort:   http/TCP
  Endpoints:    192.168.1.5:9200

Service:
  Name:         quick-elasticsearch-master
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=quick-elasticsearch
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.99.42.193
  Port:         transport  9300/TCP
  TargetPort:   transport/TCP
  Endpoints:    192.168.1.5:9300

Certificate Secret:
  Name:         quick-elasticsearch-cert
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=quick-elasticsearch
  Annotations:  <none>

Type:  Opaque

Data
====
  key_pass:     6 bytes
  node.jks:     3014 bytes
  root.jks:     864 bytes
  sgadmin.jks:  3009 bytes

Database Secret:
  Name:         quick-elasticsearch-auth
  Labels:         kubedb.com/kind=Elasticsearch
                  kubedb.com/name=quick-elasticsearch
  Annotations:  <none>

Type:  Opaque

Data
====
  READALL_PASSWORD:       8 bytes
  READALL_USERNAME:       7 bytes
  sg_action_groups.yml:   430 bytes
  sg_config.yml:          242 bytes
  sg_internal_users.yml:  156 bytes
  sg_roles_mapping.yml:   73 bytes
  ADMIN_PASSWORD:         8 bytes
  ADMIN_USERNAME:         5 bytes
  sg_roles.yml:           312 bytes

Topology:
  Type                Pod                    StartTime                      Phase
  ----                ---                    ---------                      -----
  master|client|data  quick-elasticsearch-0  2018-09-28 11:33:42 +0600 +06  Running

No Snapshots.

Events:
  Type    Reason      Age   From                    Message
  ----    ------      ----  ----                    -------
  Normal  Successful  3m    Elasticsearch operator  Successfully created Service
  Normal  Successful  3m    Elasticsearch operator  Successfully created Service
  Normal  Successful  2m    Elasticsearch operator  Successfully created StatefulSet
  Normal  Successful  2m    Elasticsearch operator  Successfully created Elasticsearch
  Normal  Successful  2m    Elasticsearch operator  Successfully patched StatefulSet
  Normal  Successful  1m    Elasticsearch operator  Successfully patched Elasticsearch
  Normal  Successful  1m    Elasticsearch operator  Successfully patched StatefulSet
  Normal  Successful  1m    Elasticsearch operator  Successfully patched Elasticsearch
```

```console
$ kubectl get service -n demo --selector=kubedb.com/kind=Elasticsearch,kubedb.com/name=quick-elasticsearch
NAME                         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
quick-elasticsearch          ClusterIP   10.100.103.159   <none>        9200/TCP   5m
quick-elasticsearch-master   ClusterIP   10.99.42.193     <none>        9300/TCP   5m
```

Two services for each Elasticsearch object.

- Service *`quick-elasticsearch`* targets all Pods which are acting as *client* node
- Service *`quick-elasticsearch-master`* targets all Pods which are acting as *master* node

KubeDB supports Elasticsearch clustering where pods can be any of these three roles: *master*, *data* or *client*.

If you see `Topology` section in `kubedb describe` result, you will know role(s) of each Pod.

```console
Topology:
  Type                Pod                    StartTime                      Phase
  ----                ---                    ---------                      -----
  master|client|data  quick-elasticsearch-0  2018-09-28 11:33:42 +0600 +06  Running
```

Here, we have created an Elasticsearch database with a single node. This single node is acting as *master*, *data* and *client*.

To learn about how to configure an Elasticsearch cluster, please visit [here](/docs/guides/elasticsearch/clustering/topology.md).

Please note that KubeDB operator has created two new Secrets for Elasticsearch object.

1. `quick-elasticsearch-auth` for storing the passwords and [search-guard](https://github.com/floragunncom/search-guard) configuration.
2. `quick-elasticsearch-cert` for storing certificates used for SSL connection.

#### Secret for authentication & configuration

Auth secret is used to authenticate user for Elasticsearch database and configure Search Guard plugin.

```console
$ kubectl get secret -n demo quick-elasticsearch-auth -o yaml
```

```yaml
apiVersion: v1
data:
  ADMIN_PASSWORD: Y2JjaXdjZmg=
  ADMIN_USERNAME: YWRtaW4=
  READALL_PASSWORD: YW02b21zY2g=
  READALL_USERNAME: cmVhZGFsbA==
  sg_action_groups.yml: ClVOTElNSVRFRDoKICAtICIqIgoKUkVBRDoKICAtICJpbmRpY2VzOmRhdGEvcmVhZCoiCiAgLSAiaW5kaWNlczphZG1pbi9tYXBwaW5ncy9maWVsZHMvZ2V0KiIKCkNMVVNURVJfQ09NUE9TSVRFX09QU19STzoKICAtICJpbmRpY2VzOmRhdGEvcmVhZC9tZ2V0IgogIC0gImluZGljZXM6ZGF0YS9yZWFkL21zZWFyY2giCiAgLSAiaW5kaWNlczpkYXRhL3JlYWQvbXR2IgogIC0gImluZGljZXM6ZGF0YS9yZWFkL2Nvb3JkaW5hdGUtbXNlYXJjaCoiCiAgLSAiaW5kaWNlczphZG1pbi9hbGlhc2VzL2V4aXN0cyoiCiAgLSAiaW5kaWNlczphZG1pbi9hbGlhc2VzL2dldCoiCgpDTFVTVEVSX0tVQkVEQl9TTkFQU0hPVDoKICAtICJpbmRpY2VzOmRhdGEvcmVhZC9zY3JvbGwqIgoKSU5ESUNFU19LVUJFREJfU05BUFNIT1Q6CiAgLSAiaW5kaWNlczphZG1pbi9nZXQiCg==
  sg_config.yml: CnNlYXJjaGd1YXJkOgogIGR5bmFtaWM6CiAgICBhdXRoYzoKICAgICAgYmFzaWNfaW50ZXJuYWxfYXV0aF9kb21haW46CiAgICAgICAgZW5hYmxlZDogdHJ1ZQogICAgICAgIG9yZGVyOiA0CiAgICAgICAgaHR0cF9hdXRoZW50aWNhdG9yOgogICAgICAgICAgdHlwZTogYmFzaWMKICAgICAgICAgIGNoYWxsZW5nZTogdHJ1ZQogICAgICAgIGF1dGhlbnRpY2F0aW9uX2JhY2tlbmQ6CiAgICAgICAgICB0eXBlOiBpbnRlcm5hbAo=
  sg_internal_users.yml: CmFkbWluOgogIGhhc2g6ICQyYSQxMCRaQ0ROZVdyLjFiNGhJUVFCcno0TmpPaW9OTG9YVjZLRDJ4UFNEMTZ6di5IMHZFRUQvV0J3dQoKcmVhZGFsbDoKICBoYXNoOiAkMmEkMTAkSmpzUkkvVDBhb2dRb3hDcDlQZXV6dWd6Umw5UUZIMzg5aFJZUmQ0eUI5dU9lVFVGRlpiTzIK
  sg_roles.yml: CnNnX2FsbF9hY2Nlc3M6CiAgY2x1c3RlcjoKICAgIC0gVU5MSU1JVEVECiAgaW5kaWNlczoKICAgICcqJzoKICAgICAgJyonOgogICAgICAgIC0gVU5MSU1JVEVECiAgdGVuYW50czoKICAgIGFkbV90ZW5hbnQ6IFJXCiAgICB0ZXN0X3RlbmFudF9ybzogUlcKCnNnX3JlYWRhbGw6CiAgY2x1c3RlcjoKICAgIC0gQ0xVU1RFUl9DT01QT1NJVEVfT1BTX1JPCiAgICAtIENMVVNURVJfS1VCRURCX1NOQVBTSE9UCiAgaW5kaWNlczoKICAgICcqJzoKICAgICAgJyonOgogICAgICAgIC0gUkVBRAogICAgICAgIC0gSU5ESUNFU19LVUJFREJfU05BUFNIT1QK
  sg_roles_mapping.yml: CnNnX2FsbF9hY2Nlc3M6CiAgdXNlcnM6CiAgICAtIGFkbWluCgpzZ19yZWFkYWxsOgogIHVzZXJzOgogICAgLSByZWFkYWxsCg==
kind: Secret
metadata:
  creationTimestamp: 2018-09-28T05:33:36Z
  labels:
    kubedb.com/kind: Elasticsearch
    kubedb.com/name: quick-elasticsearch
  name: quick-elasticsearch-auth
  namespace: demo
  ...
type: Opaque
```

> Note: Auth Secret name format: `{elasticsearch-name}-auth`

This Secret contains:

- `ADMIN_USERNAME` *username* for superuser used in search-guard configuration as an internal user.
- `ADMIN_PASSWORD` *password* for the superuser.
- `READALL_USERNAME` *username* for `readall` user with read-only permission only.
- `READALL_PASSWORD` *password* for the `readall` user.
- Followings are used as search-guard configuration
  - `sg_action_groups.yml`
  - `sg_config.yml`
  - `sg_internal_users.yml`
  - `sg_roles.yml`
  - `sg_roles_mapping.yml`

To know more about search-guard configuration, please visit [here](/docs/guides/elasticsearch/search-guard/configuration.md).

#### Secret for certificates

Certificate secret contains SSL certificates that are used to secure communication with Elasticsearch database.

```console
$ kubectl get secret -n demo quick-elasticsearch-cert -o yaml
```

```yaml
apiVersion: v1
data:
  key_pass: ZWR0aGd3
  node.jks: <base64 encoded node certificate in jks format>
  root.jks: <base64 encoded root CA in jks format>
  sgadmin.jks: <base64 encoded admin certificate used to change the Search Guard configuration>
kind: Secret
metadata:
  creationTimestamp: 2018-09-28T05:33:35Z
  labels:
    kubedb.com/kind: Elasticsearch
    kubedb.com/name: quick-elasticsearch
  name: quick-elasticsearch-cert
  namespace: demo
  ...
type: Opaque
```

> Note: Cert Secret name format: `{elasticsearch-name}-cert`

To know more about how to create TLS secure Elasticsearch database with KubeDB, please visit [here](/docs/guides/elasticsearch/search-guard/use-tls.md).

## Connect with Elasticsearch Database

We will use [port forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) to connect with our Elasticsearch database. Then we will use `curl` to send `http` request to check cluster health to verify that our Elasticsearch database is working well.

Let's forward `9200` port of our database pod. Run following command on a separate terminal,

```console
$ kubectl port-forward -n demo quick-elasticsearch-0 9200
Forwarding from 127.0.0.1:9200 -> 9200
Forwarding from [::1]:9200 -> 9200
```

Now, we can connect to the database at `localhost:9200`. Let's find out necessary connection information first.

**Connection information:**

- Address: `localhost:9200`
- Username: Run following command to get *username*

  ```console
  $ kubectl get secrets -n demo quick-elasticsearch-auth -o jsonpath='{.data.\ADMIN_USERNAME}' | base64 -d
  admin
  ```

- Password: Run following command to get *password*

  ```console
  $ kubectl get secrets -n demo quick-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
  cbciwcfh
  ```

Now let's check health of our Elasticsearch database.

```console
curl --user "admin:cbciwcfh" "localhost:9200/_cluster/health?pretty"
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

Requst format: `curl --user "$USERNAME:$PASSWORD" "$ADDRESS/_cluster/health?pretty"`

From the health information above, we can see that our Elasticsearch cluster's status is `green`. That means everything is going well.

## Pause Elasticsearch

KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `DoNotTerminate` termination policy. If admission webhook is enabled, it prevents user from deleting the database as long as the `spec.terminationPolicy` is set `DoNotTerminate`.

In this tutorial, Elasticsearch `quick-elasticsearch` is created with `spec.terminationPolicy: DoNotTerminate`. So if you try to delete this Elasticsearch object, admission webhook will nullify the delete operation.

```console
$ kubectl delete es -n demo quick-elasticsearch
Error from server (BadRequest): admission webhook "elasticsearch.validators.kubedb.com" denied the request: elasticsearch "quick-elasticsearch" can't be paused. To delete, change spec.terminationPolicy
```

To pause the database, we have to set `spec.terminationPolicy:` to `Pause` by updating it,

```console
$ kubectl edit es -n demo quick-elasticsearch
spec:
  terminationPolicy: Pause
```

Now, if you delete the Elasticsearch object, KubeDB operator will create a matching DormantDatabase object. KubeDB operator watches for DormantDatabase objects and it will take necessary steps when a DormantDatabase object is created.

KubeDB operator will delete the StatefulSet and its Pods, but leaves the Secret, PVCs unchanged.

```console
$ kubectl delete es -n demo quick-elasticsearch
elasticsearch.kubedb.com "quick-elasticsearch" deleted
```

Check DormantDatabase entry

```console
$ kubectl get drmn -n demo quick-elasticsearch
NAME                  STATUS    AGE
quick-elasticsearch   Paused    29s
```

In KubeDB parlance, we say that Elasticsearch `quick-elasticsearch`  has entered into dormant state.

Let's see, what we have in this DormantDatabase object

```console
$ kubectl get drmn -n demo quick-elasticsearch -o yaml
```

```yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: 2018-09-28T08:56:15Z
  finalizers:
  - kubedb.com
  generation: 1
  labels:
    kubedb.com/kind: Elasticsearch
  name: quick-elasticsearch
  namespace: demo
  resourceVersion: "23969"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/quick-elasticsearch
  uid: 5b1a99dd-c2fc-11e8-aac4-8a5cc86ecf00
spec:
  origin:
    metadata:
      annotations:
        kubectl.kubernetes.io/last-applied-configuration: |
          {"apiVersion":"kubedb.com/v1alpha1","kind":"Elasticsearch","metadata":{"annotations":{},"name":"quick-elasticsearch","namespace":"demo"},"spec":{"terminationPolicy":true,"storage":{"accessModes":["ReadWriteOnce"],"resources":{"requests":{"storage":"1Gi"}},"storageClassName":"standard"},"storageType":"Durable","version":"6.3-v1"}}
      creationTimestamp: 2018-09-28T05:33:29Z
      name: quick-elasticsearch
      namespace: demo
    spec:
      elasticsearch:
        certificateSecret:
          secretName: quick-elasticsearch-cert
        databaseSecret:
          secretName: quick-elasticsearch-auth
        podTemplate:
          controller: {}
          metadata: {}
          spec:
            resources: {}
        replicas: 1
        serviceTemplate:
          metadata: {}
          spec: {}
        storage:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 1Gi
          storageClassName: standard
        storageType: Durable
        terminationPolicy: Pause
        updateStrategy:
          type: RollingUpdate
        version: 6.3-v1
status:
  observedGeneration: 1$10263513872796756591
  pausingTime: 2018-09-28T08:56:24Z
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
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/quickstart/quick-elasticsearch.yaml
elasticsearch.kubedb.com/quick-elasticsearch created
```

## WipeOut DormantDatabase

You can wipe out a DormantDatabase while deleting the object by setting `spec.wipeOut` to `true`. KubeDB operator will delete any relevant resources of this `Elasticsearch` database (i.e, PVCs, Secrets, Snapshots). It will also delete snapshot data stored in the Cloud Storage buckets.

```yaml
$ kubectl edit drmn -n demo quick-elasticsearch
spec:
  wipeOut: true
```

You can also set `wipeOut: true` by patching the DormantDatabase,

```console
$ kubectl patch -n demo drmn/quick-elasticsearch -p '{"spec":{"wipeOut":true}}' --type="merge"
```

If `spec.wipeOut` is not set to `true` while deleting the `dormantdatabase` object, then only this object will be deleted and KubeDB operator won't delete related Secrets, PVCs and Snapshots. So, user still can access the stored data in the cloud storage buckets as well as PVCs.

## Delete DormantDatabase

As it is already discussed above, `DormantDatabase` can be deleted with or without wiping out the resources. To delete the `dormantdatabase`,

```console
$ kubectl delete drmn -n demo quick-elasticsearch
dormantdatabase.kubedb.com "quick-elasticsearch" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/quick-elasticsearch -p '{"spec":{"terminationPolicy": "WipeOut"}}' --type="merge"
$ kubectl delete -n demo es/quick-elasticsearch

$ kubectl delete ns demo
```
## Tips for Testing

If you are just testing some basic functionalities, you might want to avoid additional hassles due to some safety features that are great for production environment. You can follow these tips to avoid them.

1. **Use `storageType: Ephemeral`**. Databases are precious. You might not want to lose your data in your production environment if database pod fail. So, we recommend to use `spec.storageType: Durable` and provide storage spec in `spec.storage` section. For testing purpose, you can just use `spec.storageType: Ephemeral`. KubeDB will use [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) for storage. You will not require to provide `spec.storage` section.
2. **Use `terminationPolicy: WipeOut`**. It is nice to be able to resume database from previous one. So, we create `DormantDatabase` and preserve all your `PVCs`, `Secrets`, `Snapshots` etc. If you don't want to resume database, you can just use `spec.terminationPolicy: WipeOut`. It will not create `DormantDatabase` and it will delete everything created by KubeDB for a particular Elasticsearch crd when you delete the crd. For more details about termination policy, please visit [here](/docs/concepts/databases/elasticsearch.md#specterminationpolicy).

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
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
