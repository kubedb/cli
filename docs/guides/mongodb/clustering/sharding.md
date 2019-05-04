---
title: MongoDB Sharding Guide
menu:
  docs_0.11.0:
    identifier: mg-clustering-sharding
    name: Sharding Guide
    parent: mg-clustering-mongodb
    weight: 25
menu_name: docs_0.11.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# MongoDB Sharding

This tutorial will show you how to use KubeDB to run a sharded MongoDB cluster.

## Before You Begin

Before proceeding:

- Read [mongodb sharding concept](/docs/guides/mongodb/clustering/sharding_concept.md) to learn about MongoDB Sharding clustering.

- You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace/demo created
  ```

> Note: The yaml files used in this tutorial are stored in [docs/examples/mongodb](https://github.com/kubedb/cli/tree/master/docs/examples/mongodb) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Deploy Sharded MongoDB Cluster

To deploy a MongoDB Sharding, user have to specify `spec.replicaSet` option in `Mongodb` CRD.

The following is an example of a `Mongodb` object which creates MongoDB Sharding of three members.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mongo-sh
  namespace: demo
spec:
  version: 3.6-v3
  shardTopology:
    configServer:
      replicas: 3
      storage:
        resources:
          requests:
            storage: 1Gi
        storageClassName: standard
    mongos:
      replicas: 2
      strategy:
        type: RollingUpdate
    shard:
      replicas: 3
      shards: 3
      storage:
        resources:
          requests:
            storage: 1Gi
        storageClassName: standard
```

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.11.0//docs/examples/mongodb/clustering/mongo-sharding.yaml
mongodb.kubedb.com/mongo-sh created
```

Here,

- `spec.shardTopology` represents the topology configuration for sharding.
  - `shard` represents configuration for Shard component of mongodb.
    - `shards` represents number of shards for a mongodb deployment. Each shard is deployed as a [replicaset](/docs/guides/mongodb/clustering/replication_concept.md).
    - `replicas` represents number of replicas of each shard replicaset.
    - `prefix` represents the prefix of each shard node.
    - `configSource` is an optional field to provide custom configuration file for shards (i.e mongod.cnf). If specified, this file will be used as configuration file otherwise a default configuration file will be used.
    - `podTemplate` is an optional configuration for pods.
    - `storage` to specify pvc spec for each node of sharding. You can specify any StorageClass available in your cluster with appropriate resource requests.
  - `configServer` represents configuration for ConfigServer component of mongodb.
    - `replicas` represents number of replicas for configServer replicaset. Here, configServer is deployed as a replicaset of mongodb.
    - `prefix` represents the prefix of configServer nodes.
    - `configSource` is an optional field to provide custom configuration file for configSource (i.e mongod.cnf). If specified, this file will be used as configuration file otherwise a default configuration file will be used.
    - `podTemplate` is an optional configuration for pods.
    - `storage` to specify pvc spec for each node of configServer. You can specify any StorageClass available in your cluster with appropriate resource requests.
  - `mongos` represents configuration for Mongos component of mongodb. `Mongos` instances run as stateless components (deployment).
    - `replicas` represents number of replicas of `Mongos` instance. Here, Mongos is not deployed as replicaset.
    - `prefix` represents the prefix of mongos nodes.
    - `configSource` is an optional field to provide custom configuration file for mongos (i.e mongod.cnf). If specified, this file will be used as configuration file otherwise a default configuration file will be used.
    - `podTemplate` is an optional configuration for pods.
    - `strategy` is the deployment strategy to use to replace existing pods with new ones. This is optional. If not provided, kubernetes will use default deploymentStrategy, ie. `RollingUpdate`. See more about [Deployment Strategy](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy).
- `spec.certificateSecret` (optional) is a secret name that contains keyfile(a random string) against `key.txt` key. Each mongod replica set instances in the topology uses the contents of the keyfile as the shared password for authenticating other members in the replicaset. Only mongod instances with the correct keyfile can join the replica set. _User can provide the `CertificateSecret` by creating a secret with key `key.txt`. See [here](https://docs.mongodb.com/manual/tutorial/enforce-keyfile-access-control-in-existing-replica-set/#create-a-keyfile) to create the string for `CertificateSecret`._ If `CertificateSecret` is not given, KubeDB operator will generate a `CertificateSecret` itself.

KubeDB operator watches for `MongoDB` objects using Kubernetes api. When a `MongoDB` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching MongoDB object name. KubeDB operator will also create governing services for StatefulSets with the name `<mongodb-name>-<node-type>-gvr`. No MongoDB specific RBAC permission is required in [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

MongoDB `mongo-sh` state,

```console
$ kubectl get mg -n demo
NAME       VERSION   STATUS    AGE
mongo-sh   3.6-v3    Running   9m41s
```

`Sharding` and `ConfigServer` nodes are deployed as statefulset.

```console
$ kubectl get statefulset -n demo
NAME                 READY   AGE
mongo-sh-configsvr   3/3     11m
mongo-sh-shard0      3/3     10m
mongo-sh-shard1      3/3     8m59s
mongo-sh-shard2      3/3     7m45s
```

`Mongos` nodes are deployed as deployment.

```console
$ kubectl get deployments -n demo
NAME              READY   UP-TO-DATE   AVAILABLE   AGE
mongo-sh-mongos   2/2     2            2           8m41s
```

All PVCs and PVs for MongoDB `mongo-sh`,

```console
$ kubectl get pvc -n demo
NAME                           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
datadir-mongo-sh-configsvr-0   Bound    pvc-1db4185e-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       16m
datadir-mongo-sh-configsvr-1   Bound    pvc-330cc6ee-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       16m
datadir-mongo-sh-configsvr-2   Bound    pvc-3db2d3f5-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       15m
datadir-mongo-sh-shard0-0      Bound    pvc-49b7cc3b-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       15m
datadir-mongo-sh-shard0-1      Bound    pvc-5b781770-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       15m
datadir-mongo-sh-shard0-2      Bound    pvc-6ba3263e-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       14m
datadir-mongo-sh-shard1-0      Bound    pvc-75feb227-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       14m
datadir-mongo-sh-shard1-1      Bound    pvc-89bb7bb3-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       13m
datadir-mongo-sh-shard1-2      Bound    pvc-98c96ae4-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       13m
datadir-mongo-sh-shard2-0      Bound    pvc-a1eebcd2-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       13m
datadir-mongo-sh-shard2-1      Bound    pvc-b231fb18-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       12m
datadir-mongo-sh-shard2-2      Bound    pvc-c5bb265f-6a5f-11e9-a871-080027a851ba   1Gi        RWO            standard       12m


$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                               STORAGECLASS   REASON   AGE
pvc-1db4185e-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-configsvr-0   standard                17m
pvc-330cc6ee-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-configsvr-1   standard                16m
pvc-3db2d3f5-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-configsvr-2   standard                16m
pvc-49b7cc3b-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-shard0-0      standard                16m
pvc-5b781770-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-shard0-1      standard                15m
pvc-6ba3263e-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-shard0-2      standard                15m
pvc-75feb227-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-shard1-0      standard                14m
pvc-89bb7bb3-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-shard1-1      standard                14m
pvc-98c96ae4-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-shard1-2      standard                13m
pvc-a1eebcd2-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-shard2-0      standard                13m
pvc-b231fb18-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-shard2-1      standard                13m
pvc-c5bb265f-6a5f-11e9-a871-080027a851ba   1Gi        RWO            Delete           Bound    demo/datadir-mongo-sh-shard2-2      standard                12m
```

Services created for MongoDB `mongo-sh`

```console
$ kubectl get svc -n demo
NAME                     TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
mongo-sh                 ClusterIP   10.108.188.201   <none>        27017/TCP   18m
mongo-sh-configsvr-gvr   ClusterIP   None             <none>        27017/TCP   18m
mongo-sh-shard0-gvr      ClusterIP   None             <none>        27017/TCP   18m
mongo-sh-shard1-gvr      ClusterIP   None             <none>        27017/TCP   18m
mongo-sh-shard2-gvr      ClusterIP   None             <none>        27017/TCP   18m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. It has also defaulted some field of crd object. Run the following command to see the modified MongoDB object:

```yaml
$ kubedb get mg -n demo mongo-sh -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  creationTimestamp: "2019-04-29T09:13:56Z"
  finalizers:
  - kubedb.com
  generation: 3
  name: mongo-sh
  namespace: demo
  resourceVersion: "25825"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mongodbs/mongo-sh
  uid: 1d83622c-6a5f-11e9-a871-080027a851ba
spec:
  certificateSecret:
    secretName: mongo-sh-keyfile
  databaseSecret:
    secretName: mongo-sh-auth
  serviceTemplate:
    metadata: {}
    spec: {}
  shardTopology:
    configServer:
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
          securityContext:
            fsGroup: 999
            runAsNonRoot: true
            runAsUser: 999
      replicas: 3
      storage:
        dataSource: null
        resources:
          requests:
            storage: 1Gi
        storageClassName: standard
    mongos:
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
          securityContext:
            fsGroup: 999
            runAsNonRoot: true
            runAsUser: 999
      replicas: 2
      strategy:
        type: RollingUpdate
    shard:
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
          securityContext:
            fsGroup: 999
            runAsNonRoot: true
            runAsUser: 999
      replicas: 3
      shards: 3
      storage:
        dataSource: null
        resources:
          requests:
            storage: 1Gi
        storageClassName: standard
  storageType: Durable
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
  version: 3.6-v3
status:
  observedGeneration: 3$4212299729528774793
  phase: Running
```

Please note that KubeDB operator has created a new Secret called `mongo-sh-auth` _(format: {mongodb-object-name}-auth)_ for storing the password for `mongodb` superuser. This secret contains a `username` key which contains the _username_ for MongoDB superuser and a `password` key which contains the _password_ for MongoDB superuser.

If you want to use custom or existing secret please specify that when creating the MongoDB object using `spec.databaseSecret.secretName`. While creating this secret manually, make sure the secret contains these two keys containing data `username` and `password`. For more details, please see [here](/docs/concepts/databases/mongodb.md#specdatabasesecret).

## Connection Information

- Hostname/address: you can use any of these
  - Service: `mongo-sh.demo`
  - Pod IP: (`$ kubectl get po -n demo -l mongodb.kubedb.com/node.mongos=mongo-sh-mongos -o yaml | grep podIP`)
- Port: `27017`
- Username: Run following command to get _username_,

  ```console
  $ kubectl get secrets -n demo mongo-sh-auth -o jsonpath='{.data.\username}' | base64 -d
  root
  ```

- Password: Run the following command to get _password_,

  ```console
  $ kubectl get secrets -n demo mongo-sh-auth -o jsonpath='{.data.\password}' | base64 -d
  7QiqLcuSCmZ8PU5a
  ```

Now, you can connect to this database through [mongo-shell](https://docs.mongodb.com/v3.6/mongo/).

## Sharded Data

In this tutorial, we will insert sharded and unsharded document, and we will see if the data actually sharded across cluster or not.

```console
$ kubectl get po -n demo -l mongodb.kubedb.com/node.mongos=mongo-sh-mongos
NAME                               READY   STATUS    RESTARTS   AGE
mongo-sh-mongos-69b557f9f5-2kz68   1/1     Running   0          49m
mongo-sh-mongos-69b557f9f5-5hvh2   1/1     Running   0          49m

$ kubectl exec -it mongo-sh-mongos-69b557f9f5-2kz68 -n demo bash

mongodb@mongo-sh-mongos-69b557f9f5-2kz68:/$ mongo admin -u root -p 7QiqLcuSCmZ8PU5a
MongoDB shell version v3.6.12
connecting to: mongodb://127.0.0.1:27017/admin?gssapiServiceName=mongodb
Implicit session: session { "id" : UUID("8b7abf57-09e4-4e30-b4a0-a37ebf065e8f") }
MongoDB server version: 3.6.12
Welcome to the MongoDB shell.
For interactive help, type "help".
For more comprehensive documentation, see
	http://docs.mongodb.org/
Questions? Try the support group
	http://groups.google.com/group/mongodb-user
2019-04-29T10:09:17.311+0000 I STORAGE  [main] In File::open(), ::open for '/home/mongodb/.mongorc.js' failed with No such file or directory
mongos> isMaster;
2019-04-29T10:11:16.128+0000 E QUERY    [thread1] ReferenceError: isMaster is not defined :
@(shell):1:1
mongos> isMaster();
2019-04-29T10:11:20.020+0000 E QUERY    [thread1] ReferenceError: isMaster is not defined :
@(shell):1:1

mongos>
```

To detect if the MongoDB instance that your client is connected to is mongos, use the isMaster command. When a client connects to a mongos, isMaster returns a document with a `msg` field that holds the string `isdbgrid`.

```console
mongos> rs.isMaster()
{
	"ismaster" : true,
	"msg" : "isdbgrid",
	"maxBsonObjectSize" : 16777216,
	"maxMessageSizeBytes" : 48000000,
	"maxWriteBatchSize" : 100000,
	"localTime" : ISODate("2019-04-29T10:12:00.145Z"),
	"logicalSessionTimeoutMinutes" : 30,
	"maxWireVersion" : 6,
	"minWireVersion" : 0,
	"ok" : 1,
	"operationTime" : Timestamp(1556532710, 1),
	"$clusterTime" : {
		"clusterTime" : Timestamp(1556532710, 1),
		"signature" : {
			"hash" : BinData(0,"6W7pmWBdVSzY0x+BxQj74d0WhXg="),
			"keyId" : NumberLong("6685242219722440730")
		}
	}
}
```

`mongo-sh` Shard status,

```console
mongos> sh.status()
--- Sharding Status ---
  sharding version: {
  	"_id" : 1,
  	"minCompatibleVersion" : 5,
  	"currentVersion" : 6,
  	"clusterId" : ObjectId("5cc6c061f439d076e04d737b")
  }
  shards:
        {  "_id" : "shard0",  "host" : "shard0/mongo-sh-shard0-0.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-1.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-2.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard1",  "host" : "shard1/mongo-sh-shard1-0.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-1.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-2.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard2",  "host" : "shard2/mongo-sh-shard2-0.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-1.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-2.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
  active mongoses:
        "3.6.12" : 2
  autosplit:
        Currently enabled: yes
  balancer:
        Currently enabled:  yes
        Currently running:  no
        Failed balancer rounds in last 5 attempts:  0
        Migration Results for the last 24 hours:
                No recent migrations
  databases:
        {  "_id" : "config",  "primary" : "config",  "partitioned" : true }
                config.system.sessions
                        shard key: { "_id" : 1 }
                        unique: false
                        balancing: true
                        chunks:
                                shard0	1
                        { "_id" : { "$minKey" : 1 } } -->> { "_id" : { "$maxKey" : 1 } } on : shard0 Timestamp(1, 0)
```

Shard collection `test.testcoll` and insert document. See [`sh.shardCollection(namespace, key, unique, options)`](https://docs.mongodb.com/manual/reference/method/sh.shardCollection/#sh.shardCollection) for details about `shardCollection` command.

```console
mongos> sh.enableSharding("test");
{
	"ok" : 1,
	"operationTime" : Timestamp(1556535000, 8),
	"$clusterTime" : {
		"clusterTime" : Timestamp(1556535000, 8),
		"signature" : {
			"hash" : BinData(0,"84KefOzN8tKmsPfr6IrnBUxF9NM="),
			"keyId" : NumberLong("6685242219722440730")
		}
	}
}

mongos> sh.shardCollection("test.testcoll", {"myfield": 1});
{
	"collectionsharded" : "test.testcoll",
	"collectionUUID" : UUID("68ff9452-40bb-41a2-b35a-405132f90cd3"),
	"ok" : 1,
	"operationTime" : Timestamp(1556535010, 8),
	"$clusterTime" : {
		"clusterTime" : Timestamp(1556535010, 8),
		"signature" : {
			"hash" : BinData(0,"IgVzMa8qE4UBzjc2gOZJX5kZ3T4="),
			"keyId" : NumberLong("6685242219722440730")
		}
	}
}

mongos> use test;
switched to db test

mongos> db.testcoll.insert({"myfield": "a", "otherfield": "b"});
WriteResult({ "nInserted" : 1 })

mongos> db.testcoll.insert({"myfield": "c", "otherfield": "d", "kube" : "db" });
WriteResult({ "nInserted" : 1 })

mongos> db.testcoll.find();
{ "_id" : ObjectId("5cc6d6f656a9ddd30be2c12a"), "myfield" : "a", "otherfield" : "b" }
{ "_id" : ObjectId("5cc6d71e56a9ddd30be2c12b"), "myfield" : "c", "otherfield" : "d", "kube" : "db" }

```

Run [`sh.status()`](https://docs.mongodb.com/manual/reference/method/sh.status/) to see whether the `test` database has sharding enabled, and the primary shard for the `test` database.

The Sharded Collection section `sh.status.databases.<collection>` provides information on the sharding details for sharded collection(s) (E.g. `test.testcoll`). For each sharded collection, the section displays the shard key, the number of chunks per shard(s), the distribution of documents across chunks, and the tag information, if any, for shard key range(s).

```console
mongos> sh.status();
--- Sharding Status ---
  sharding version: {
  	"_id" : 1,
  	"minCompatibleVersion" : 5,
  	"currentVersion" : 6,
  	"clusterId" : ObjectId("5cc6c061f439d076e04d737b")
  }
  shards:
        {  "_id" : "shard0",  "host" : "shard0/mongo-sh-shard0-0.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-1.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-2.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard1",  "host" : "shard1/mongo-sh-shard1-0.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-1.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-2.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard2",  "host" : "shard2/mongo-sh-shard2-0.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-1.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-2.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
  active mongoses:
        "3.6.12" : 2
  autosplit:
        Currently enabled: yes
  balancer:
        Currently enabled:  yes
        Currently running:  no
        Failed balancer rounds in last 5 attempts:  0
        Migration Results for the last 24 hours:
                No recent migrations
  databases:
        {  "_id" : "config",  "primary" : "config",  "partitioned" : true }
                config.system.sessions
                        shard key: { "_id" : 1 }
                        unique: false
                        balancing: true
                        chunks:
                                shard0	1
                        { "_id" : { "$minKey" : 1 } } -->> { "_id" : { "$maxKey" : 1 } } on : shard0 Timestamp(1, 0)
        {  "_id" : "test",  "primary" : "shard1",  "partitioned" : true }
                test.testcoll
                        shard key: { "myfield" : 1 }
                        unique: false
                        balancing: true
                        chunks:
                                shard1	1
                        { "myfield" : { "$minKey" : 1 } } -->> { "myfield" : { "$maxKey" : 1 } } on : shard1 Timestamp(1, 0)
```

Now create another database where partiotioned is not applied and see how the data is stored.

```
mongos> use demo
switched to db demo

mongos> db.testcoll2.insert({"myfield": "ccc", "otherfield": "d", "kube" : "db" });
WriteResult({ "nInserted" : 1 })

mongos> db.testcoll2.insert({"myfield": "aaa", "otherfield": "d", "kube" : "db" });
WriteResult({ "nInserted" : 1 })


mongos> db.testcoll2.find()
{ "_id" : ObjectId("5cc6dc831b6d9b3cddc947ec"), "myfield" : "ccc", "otherfield" : "d", "kube" : "db" }
{ "_id" : ObjectId("5cc6dce71b6d9b3cddc947ed"), "myfield" : "aaa", "otherfield" : "d", "kube" : "db" }
```

Now, eventually `sh.status()`

```
mongos> sh.status()
--- Sharding Status ---
  sharding version: {
  	"_id" : 1,
  	"minCompatibleVersion" : 5,
  	"currentVersion" : 6,
  	"clusterId" : ObjectId("5cc6c061f439d076e04d737b")
  }
  shards:
        {  "_id" : "shard0",  "host" : "shard0/mongo-sh-shard0-0.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-1.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-2.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard1",  "host" : "shard1/mongo-sh-shard1-0.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-1.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-2.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard2",  "host" : "shard2/mongo-sh-shard2-0.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-1.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-2.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
  active mongoses:
        "3.6.12" : 2
  autosplit:
        Currently enabled: yes
  balancer:
        Currently enabled:  yes
        Currently running:  no
        Failed balancer rounds in last 5 attempts:  0
        Migration Results for the last 24 hours:
                No recent migrations
  databases:
        {  "_id" : "config",  "primary" : "config",  "partitioned" : true }
                config.system.sessions
                        shard key: { "_id" : 1 }
                        unique: false
                        balancing: true
                        chunks:
                                shard0	1
                        { "_id" : { "$minKey" : 1 } } -->> { "_id" : { "$maxKey" : 1 } } on : shard0 Timestamp(1, 0)
        {  "_id" : "demo",  "primary" : "shard2",  "partitioned" : false }
        {  "_id" : "test",  "primary" : "shard1",  "partitioned" : true }
                test.testcoll
                        shard key: { "myfield" : 1 }
                        unique: false
                        balancing: true
                        chunks:
                                shard1	1
                        { "myfield" : { "$minKey" : 1 } } -->> { "myfield" : { "$maxKey" : 1 } } on : shard1 Timestamp(1, 0)
```

Here, `demo` database is not partitioned and all collections under `demo` database are stored in it's primary shard, which is `shard2`.

## Update number of ShardTopology Instances

User can increase or decrease the number of router/mongos `spec.shardTopology.mongos.replicas`.

At this moment, decreasing the number of shards and replicasets is not handled from KubeDB end. But, User can increase the number of instances if needed.

Here a table of allowed actions are given for mongodb `ShardTopology`,

|                                 Instances                                  | Increase | Decrease |
| :------------------------------------------------------------------------: | :------: | :------: |
|        # of replicas of Mongos `spec.shardTopology.mongos.replicas`        | &#10003; | &#10003; |
|               # of Shards `spec.shardTopology.shard.shards`                | &#10003; | &#10007; |
|     # of replicaset of each Shard `spec.shardTopology.shard.replicas`      | &#10003; | &#10007; |
| # of replicaset of ConfigServer `spec.shardTopology.configServer.replicas` | &#10003; | &#10007; |

Now edit MongoDB `mongo-sh` to increase `spec.shardTopology.shard.shards` to 4 and increase `spec.shardTopology.mongos` to 3.

```console
$ kubectl edit mg -n demo mongo-sh
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mongo-sh
  namespace: demo
  ...
spec:
  shardTopology:
    mongos:
      replicas: 3 # set 3
      ...
    shard:
      shards: 4 # set 4
      ...
  ...
```

Watch for pod changes,

```console
$ kubectl get po --all-namespaces -w
NAMESPACE     NAME                                    READY   STATUS    RESTARTS   AGE
demo          mongo-sh-configsvr-0                    1/1     Running   0          8m12s
demo          mongo-sh-configsvr-1                    1/1     Running   0          7m41s
demo          mongo-sh-configsvr-2                    1/1     Running   0          7m17s
demo          mongo-sh-mongos-69b557f9f5-2qvb7        1/1     Running   0          3m44s
demo          mongo-sh-mongos-69b557f9f5-6z2s4        1/1     Running   0          3m44s
demo          mongo-sh-shard0-0                       1/1     Running   0          7m4s
demo          mongo-sh-shard0-1                       1/1     Running   0          6m37s
demo          mongo-sh-shard0-2                       1/1     Running   0          6m20s
demo          mongo-sh-shard1-0                       1/1     Running   0          5m55s
demo          mongo-sh-shard1-1                       1/1     Running   0          5m29s
demo          mongo-sh-shard1-2                       1/1     Running   0          5m14s
demo          mongo-sh-shard2-0                       1/1     Running   0          4m59s
demo          mongo-sh-shard2-1                       1/1     Running   0          4m32s
demo          mongo-sh-shard2-2                       1/1     Running   0          4m11s
kube-system   coredns-fb8b8dccf-nzb5q                 1/1     Running   1          165m
kube-system   coredns-fb8b8dccf-tqldv                 1/1     Running   1          165m
kube-system   etcd-minikube                           1/1     Running   0          164m
kube-system   kube-addon-manager-minikube             1/1     Running   0          164m
kube-system   kube-apiserver-minikube                 1/1     Running   0          163m
kube-system   kube-controller-manager-minikube        1/1     Running   0          164m
kube-system   kube-proxy-qznbv                        1/1     Running   0          165m
kube-system   kube-scheduler-minikube                 1/1     Running   0          163m
kube-system   kubedb-operator-c8cd5c69-c5nkc          1/1     Running   0          44m
kube-system   kubernetes-dashboard-79dd6bfc48-l88rk   1/1     Running   4          165m
kube-system   storage-provisioner                     1/1     Running   0          164m
demo          mongo-sh-shard3-0                       0/1     Pending   0          0s
demo          mongo-sh-shard3-0                       0/1     Pending   0          0s
demo          mongo-sh-shard3-0                       0/1     Pending   0          9s
demo          mongo-sh-shard3-0                       0/1     Init:0/2   0          9s
demo          mongo-sh-shard3-0                       0/1     Init:1/2   0          10s
demo          mongo-sh-shard3-0                       0/1     Init:1/2   0          11s
demo          mongo-sh-shard3-0                       0/1     PodInitializing   0          30s
demo          mongo-sh-shard3-0                       0/1     Running           0          31s
demo          mongo-sh-shard3-0                       1/1     Running           0          41s
demo          mongo-sh-shard3-1                       0/1     Pending           0          0s
demo          mongo-sh-shard3-1                       0/1     Pending           0          0s
demo          mongo-sh-shard3-1                       0/1     Pending           0          5s
demo          mongo-sh-shard3-1                       0/1     Init:0/2          0          5s
demo          mongo-sh-shard3-1                       0/1     Init:1/2          0          6s
demo          mongo-sh-shard3-1                       0/1     Init:1/2          0          7s
demo          mongo-sh-shard3-1                       0/1     PodInitializing   0          15s
demo          mongo-sh-shard3-1                       0/1     Running           0          16s
demo          mongo-sh-shard3-1                       1/1     Running           0          23s
demo          mongo-sh-shard3-2                       0/1     Pending           0          0s
demo          mongo-sh-shard3-2                       0/1     Pending           0          0s
demo          mongo-sh-shard3-2                       0/1     Pending           0          2s
demo          mongo-sh-shard3-2                       0/1     Init:0/2          0          2s
demo          mongo-sh-shard3-2                       0/1     Init:1/2          0          4s
demo          mongo-sh-shard3-2                       0/1     Init:1/2          0          5s
demo          mongo-sh-shard3-2                       0/1     PodInitializing   0          14s
demo          mongo-sh-shard3-2                       0/1     Running           0          15s
demo          mongo-sh-shard3-2                       1/1     Running           0          22s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Pending           0          0s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Pending           0          0s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Init:0/2          0          0s
demo          mongo-sh-mongos-598658d8f9-6j544        0/1     Pending           0          0s
demo          mongo-sh-mongos-598658d8f9-6j544        0/1     Pending           0          0s
demo          mongo-sh-mongos-598658d8f9-6j544        0/1     Init:0/2          0          0s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Init:0/2          0          1s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Init:1/2          0          3s
demo          mongo-sh-mongos-598658d8f9-6j544        0/1     Init:1/2          0          3s
demo          mongo-sh-mongos-598658d8f9-6j544        0/1     Init:1/2          0          4s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Init:1/2          0          4s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     PodInitializing   0          7s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Running           0          8s
demo          mongo-sh-mongos-598658d8f9-6j544        0/1     PodInitializing   0          8s
demo          mongo-sh-mongos-598658d8f9-6j544        0/1     Running           0          9s
demo          mongo-sh-mongos-598658d8f9-6j544        1/1     Running           0          12s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Terminating       0          12s
demo          mongo-sh-mongos-598658d8f9-tmwvb        0/1     Pending           0          0s
demo          mongo-sh-mongos-598658d8f9-tmwvb        0/1     Pending           0          1s
demo          mongo-sh-mongos-598658d8f9-tmwvb        0/1     Init:0/2          0          1s
demo          mongo-sh-mongos-598658d8f9-tmwvb        0/1     Init:1/2          0          3s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Terminating       0          15s
demo          mongo-sh-mongos-598658d8f9-tmwvb        0/1     Init:1/2          0          4s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Terminating       0          19s
demo          mongo-sh-mongos-69b557f9f5-w77cs        0/1     Terminating       0          19s
demo          mongo-sh-mongos-598658d8f9-tmwvb        0/1     PodInitializing   0          7s
demo          mongo-sh-mongos-598658d8f9-tmwvb        0/1     Running           0          8s
demo          mongo-sh-mongos-598658d8f9-tmwvb        1/1     Running           0          10s
demo          mongo-sh-mongos-69b557f9f5-6z2s4        1/1     Terminating       0          6m6s
demo          mongo-sh-mongos-598658d8f9-464xr        0/1     Pending           0          0s
demo          mongo-sh-mongos-598658d8f9-464xr        0/1     Pending           0          0s
demo          mongo-sh-mongos-598658d8f9-464xr        0/1     Init:0/2          0          0s
demo          mongo-sh-mongos-69b557f9f5-6z2s4        0/1     Terminating       0          6m7s
demo          mongo-sh-mongos-598658d8f9-464xr        0/1     Init:1/2          0          2s
demo          mongo-sh-mongos-598658d8f9-464xr        0/1     Init:1/2          0          3s
demo          mongo-sh-mongos-598658d8f9-464xr        0/1     PodInitializing   0          6s
demo          mongo-sh-mongos-69b557f9f5-6z2s4        0/1     Terminating       0          6m13s
demo          mongo-sh-mongos-69b557f9f5-6z2s4        0/1     Terminating       0          6m13s
demo          mongo-sh-mongos-598658d8f9-464xr        0/1     Running           0          7s
demo          mongo-sh-mongos-598658d8f9-464xr        1/1     Running           0          14s
demo          mongo-sh-mongos-69b557f9f5-2qvb7        1/1     Terminating       0          6m21s
demo          mongo-sh-mongos-69b557f9f5-2qvb7        0/1     Terminating       0          6m22s
demo          mongo-sh-mongos-69b557f9f5-2qvb7        0/1     Terminating       0          6m23s
demo          mongo-sh-mongos-69b557f9f5-2qvb7        0/1     Terminating       0          6m23s
```

You can see that an extra statefulset `mongo-sh-shard3` is created as 4th shard and one extra mongos instance also came up.

Notice that, all new mongos instances came up replacing old instances because of some changes in `shard` config. This update strategy follows `spec.shardTopology.mongos.strategy`, which is optional. If not provided, kubernetes will use default deploymentStrategy, ie. `RollingUpdate`. See more about [Deployment Strategy](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy).

```console
$ kubectl get deploy -n demo -w
NAME              READY   UP-TO-DATE   AVAILABLE   AGE
mongo-sh-mongos   2/2     2            2           12m
mongo-sh-mongos   2/3     2            2           13m
mongo-sh-mongos   2/3     2            2           13m
mongo-sh-mongos   2/3     2            2           13m
mongo-sh-mongos   2/3     3            2           13m
mongo-sh-mongos   3/3     3            3           13m
```

Now check `sh.status()` in `mongos`,

```console
$ kubectl get po -n demo -l mongodb.kubedb.com/node.mongos=mongo-sh-mongos
NAME                               READY   STATUS    RESTARTS   AGE
mongo-sh-mongos-598658d8f9-6j544   1/1     Running   0          17m
mongo-sh-mongos-598658d8f9-s8gn4   1/1     Running   0          9m54s
mongo-sh-mongos-598658d8f9-tmwvb   1/1     Running   0          16m


$ kubectl exec -it mongo-sh-mongos-598658d8f9-6j544 -n demo bash

mongodb@mongo-sh-mongos-598658d8f9-6j544:/$ mongo admin -u root -p 7QiqLcuSCmZ8PU5a

mongos> sh.status()
--- Sharding Status ---
  sharding version: {
  	"_id" : 1,
  	"minCompatibleVersion" : 5,
  	"currentVersion" : 6,
  	"clusterId" : ObjectId("5cc7ed91d06f28b1b3c64c66")
  }
  shards:
        {  "_id" : "shard0",  "host" : "shard0/mongo-sh-shard0-0.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-1.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-2.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard1",  "host" : "shard1/mongo-sh-shard1-0.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-1.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-2.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard2",  "host" : "shard2/mongo-sh-shard2-0.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-1.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-2.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard3",  "host" : "shard3/mongo-sh-shard3-0.mongo-sh-shard3-gvr.demo.svc.cluster.local:27017,mongo-sh-shard3-1.mongo-sh-shard3-gvr.demo.svc.cluster.local:27017,mongo-sh-shard3-2.mongo-sh-shard3-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
  active mongoses:
        "3.6.12" : 3
  autosplit:
        Currently enabled: yes
  balancer:
        Currently enabled:  yes
        Currently running:  no
        Failed balancer rounds in last 5 attempts:  0
        Migration Results for the last 24 hours:
                No recent migrations
  databases:
        {  "_id" : "config",  "primary" : "config",  "partitioned" : true }
                config.system.sessions
                        shard key: { "_id" : 1 }
                        unique: false
                        balancing: true
                        chunks:
                                shard0	1
                        { "_id" : { "$minKey" : 1 } } -->> { "_id" : { "$maxKey" : 1 } } on : shard0 Timestamp(1, 0)
        {  "_id" : "demo",  "primary" : "shard2",  "partitioned" : false }
        {  "_id" : "test",  "primary" : "shard1",  "partitioned" : true }
                test.testcoll
                        shard key: { "myfield" : 1 }
                        unique: false
                        balancing: true
                        chunks:
                                shard1	1
                        { "myfield" : { "$minKey" : 1 } } -->> { "myfield" : { "$maxKey" : 1 } } on : shard1 Timestamp(1, 0)
```

## Pause Database

When, `terminationPolicy` is `DoNotTerminate`, KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `DoNotTerminate` feature. If admission webhook is enabled, It prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`.

Since the MongoDB object created in this tutorial has `spec.terminationPolicy` set to `Pause` (default), if you delete the MongoDB object, KubeDB operator will create a dormant database while deleting the StatefulSet and its pods but leaves the PVCs unchanged.

```console
$ kubedb delete mg mongo-sh -n demo
mongodb.kubedb.com "mongo-sh" deleted

$ kubedb get drmn -n demo mongo-sh
NAME       STATUS    AGE
mongo-sh   Pausing   13s

$ kubedb get drmn -n demo mongo-sh
NAME       STATUS   AGE
mongo-sh   Paused   52s
```

```yaml
$ kubedb get drmn -n demo mongo-sh -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: "2019-04-29T11:24:24Z"
  finalizers:
  - kubedb.com
  generation: 1
  labels:
    kubedb.com/kind: MongoDB
  name: mongo-sh
  namespace: demo
  resourceVersion: "35082"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mongo-sh
  uid: 579c2c2d-6a71-11e9-a871-080027a851ba
spec:
  origin:
    metadata:
      creationTimestamp: "2019-04-29T09:13:56Z"
      name: mongo-sh
      namespace: demo
    spec:
      mongodb:
        certificateSecret:
          secretName: mongo-sh-keyfile
        databaseSecret:
          secretName: mongo-sh-auth
        serviceTemplate:
          metadata: {}
          spec: {}
        shardTopology:
          configServer:
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
                securityContext:
                  fsGroup: 999
                  runAsNonRoot: true
                  runAsUser: 999
            replicas: 3
            storage:
              dataSource: null
              resources:
                requests:
                  storage: 1Gi
              storageClassName: standard
          mongos:
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
                securityContext:
                  fsGroup: 999
                  runAsNonRoot: true
                  runAsUser: 999
            replicas: 2
            strategy:
              type: RollingUpdate
          shard:
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
                securityContext:
                  fsGroup: 999
                  runAsNonRoot: true
                  runAsUser: 999
            replicas: 3
            shards: 3
            storage:
              dataSource: null
              resources:
                requests:
                  storage: 1Gi
              storageClassName: standard
        storageType: Durable
        terminationPolicy: Pause
        updateStrategy:
          type: RollingUpdate
        version: 3.6-v3
status:
  observedGeneration: 1$16440556888999634490
  pausingTime: "2019-04-29T11:24:41Z"
  phase: Paused
```

Here,

- `spec.origin` is the spec of the original spec of the original MongoDB object.
- `status.phase` points to the current database state `Paused`.

## Resume Dormant Database

To resume the database from the dormant state, create same `MongoDB` object with same Spec.

In this tutorial, the dormant database can be resumed by creating original MongoDB object.

The below command will resume the DormantDatabase `mongo-sh`.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.11.0/docs/examples/mongodb/clustering/mongo-sh.yaml
mongodb.kubedb.com/mongo-sh created
```

```console
$ kubectl get mg -n demo
NAME       VERSION   STATUS    AGE
mongo-sh   3.6-v3    Running   6m27s
```

Now, If you again exec into `pod` and look for previous data, you will see that, all the data persists.

```console
$ kubectl get po -n demo -l mongodb.kubedb.com/node.mongos=mongo-sh-mongos
NAME                               READY   STATUS    RESTARTS   AGE
mongo-sh-mongos-69b557f9f5-62j76   1/1     Running   0          3m52s
mongo-sh-mongos-69b557f9f5-tdn69   1/1     Running   0          3m52s


$ kubectl exec -it mongo-sh-mongos-69b557f9f5-62j76 -n demo bash

mongodb@mongo-sh-mongos-69b557f9f5-62j76:/$ mongo admin -u root -p 7QiqLcuSCmZ8PU5a

mongos> use test;
switched to db test

mongos> db.testcoll.find();
{ "_id" : ObjectId("5cc6d6f656a9ddd30be2c12a"), "myfield" : "a", "otherfield" : "b" }
{ "_id" : ObjectId("5cc6d71e56a9ddd30be2c12b"), "myfield" : "c", "otherfield" : "d", "kube" : "db" }

mongos> sh.status()
--- Sharding Status ---
  sharding version: {
  	"_id" : 1,
  	"minCompatibleVersion" : 5,
  	"currentVersion" : 6,
  	"clusterId" : ObjectId("5cc6c061f439d076e04d737b")
  }
  shards:
        {  "_id" : "shard0",  "host" : "shard0/mongo-sh-shard0-0.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-1.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017,mongo-sh-shard0-2.mongo-sh-shard0-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard1",  "host" : "shard1/mongo-sh-shard1-0.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-1.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017,mongo-sh-shard1-2.mongo-sh-shard1-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
        {  "_id" : "shard2",  "host" : "shard2/mongo-sh-shard2-0.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-1.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017,mongo-sh-shard2-2.mongo-sh-shard2-gvr.demo.svc.cluster.local:27017",  "state" : 1 }
  active mongoses:
        "3.6.12" : 2
  autosplit:
        Currently enabled: yes
  balancer:
        Currently enabled:  yes
        Currently running:  no
        Failed balancer rounds in last 5 attempts:  5
        Last reported error:  Could not find host matching read preference { mode: "primary" } for set shard2
        Time of Reported error:  Mon Apr 29 2019 11:30:33 GMT+0000 (UTC)
        Migration Results for the last 24 hours:
                No recent migrations
  databases:
        {  "_id" : "config",  "primary" : "config",  "partitioned" : true }
                config.system.sessions
                        shard key: { "_id" : 1 }
                        unique: false
                        balancing: true
                        chunks:
                                shard0	1
                        { "_id" : { "$minKey" : 1 } } -->> { "_id" : { "$maxKey" : 1 } } on : shard0 Timestamp(1, 0)
        {  "_id" : "demo",  "primary" : "shard2",  "partitioned" : false }
        {  "_id" : "test",  "primary" : "shard1",  "partitioned" : true }
                test.testcoll
                        shard key: { "myfield" : 1 }
                        unique: false
                        balancing: true
                        chunks:
                                shard1	1
                        { "myfield" : { "$minKey" : 1 } } -->> { "myfield" : { "$maxKey" : 1 } } on : shard1 Timestamp(1, 0)

```

## WipeOut DormantDatabase

You can wipe out a DormantDatabase while deleting the object by setting `spec.wipeOut` to true. KubeDB operator will delete any relevant resources of this `MongoDB` database (i.e, PVCs, Secrets, Snapshots). It will also delete snapshot data stored in the Cloud Storage buckets.

```console
$ kubedb delete mg mongo-sh -n demo
mongodb.kubedb.com "mongo-sh" deleted
```

```yaml
$ kubedb edit drmn -n demo mongo-sh
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  name: mongo-sh
  namespace: demo
  ...
spec:
  wipeOut: true
  ...
status:
  phase: Paused
  ...
```

If `spec.wipeOut` is not set to true while deleting the `dormantdatabase` object, then only this object will be deleted and `kubedb-operator` won't delete related Secrets, PVCs, and Snapshots. So, users still can access the stored data in the cloud storage buckets as well as PVCs.

## Delete DormantDatabase

As it is already discussed above, `DormantDatabase` can be deleted with or without wiping out the resources. To delete the `dormantdatabase`,

```console
$ kubectl delete drmn mongo-sh -n demo
dormantdatabase.kubedb.com "mongo-sh" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mg/mongo-sh -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mg/mongo-sh

kubectl patch -n demo drmn/mongo-sh -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mongo-sh

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
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Detail concepts of [MongoDBVersion object](/docs/concepts/catalog/mongodb.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
