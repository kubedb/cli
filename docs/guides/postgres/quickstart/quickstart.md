---
title: PostgreSQL Quickstart
menu:
  docs_0.9.0:
    identifier: pg-quickstart-quickstart
    name: Overview
    parent: pg-quickstart-postgres
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Running PostgreSQL

This tutorial will show you how to use KubeDB to run a PostgreSQL database.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/postgres/lifecycle.png">
</p>

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace/demo created
```

> Note: YAML files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

>We have designed this tutorial to demonstrate a production setup of KubeDB managed PostgreSQL. If you just want to try out KubeDB, you can bypass some of the safety features following the tips [here](/docs/guides/postgres/quickstart/quickstart.md#tips-for-testing).

## Install pgAdmin

This tutorial will also use a pgAdmin to connect and test PostgreSQL database, once it is running.

Run the following command to install pgAdmin,

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/postgres/quickstart/pgadmin.yaml
deployment.apps/pgadmin created
service/pgadmin created

$ kubectl get pods -n demo --watch
NAME                      READY     STATUS              RESTARTS   AGE
pgadmin-5b4b96779-lfpfh   0/1       ContainerCreating   0          1m
pgadmin-5b4b96779-lfpfh   1/1       Running   0         2m
^C⏎
```

Now, you can open pgAdmin on your browser using following address `http://<cluster ip>:<NodePort of pgadmin service>`.

If you are using minikube then open pgAdmin in your browser by running `minikube service pgadmin -n demo`. Or you can get the URL of Service `pgadmin` by running following command

```console
$ minikube service pgadmin -n demo --url
http://192.168.99.100:31983
```

To log into the pgAdmin, use username __`admin`__ and password __`admin`__.

## Find Available StorageClass

We will have to provide `StorageClass` in Postgres crd specification. Check available `StorageClass` in your cluster using following command,

```console
$ kubectl get storageclass
NAME                 PROVISIONER                AGE
standard (default)   k8s.io/minikube-hostpath   5h
```

Here, we have `standard` StorageClass in our cluster.

## Find Available PostgresVersion

When you have installed KubeDB, it has created `PostgresVersion` crd for all supported PostgreSQL versions. Let's check available PostgresVersions by,

```console
$ kubectl get postgresversions
NAME       VERSION   DB_IMAGE                   DEPRECATED   AGE
10.2       10.2      kubedb/postgres:10.2       true         1h
10.2-v1    10.2      kubedb/postgres:10.2-v2    true         1h
10.2-v2    10.2      kubedb/postgres:10.2-v3                 1h
10.6       10.6      kubedb/postgres:10.6                    1h
11.1       11.1      kubedb/postgres:11.1                    1h
9.6        9.6       kubedb/postgres:9.6        true         1h
9.6-v1     9.6       kubedb/postgres:9.6-v2     true         1h
9.6-v2     9.6       kubedb/postgres:9.6-v3                  2h
9.6.7      9.6.7     kubedb/postgres:9.6.7      true         1h
9.6.7-v1   9.6.7     kubedb/postgres:9.6.7-v2   true         1h
9.6.7-v2   9.6.7     kubedb/postgres:9.6.7-v3                1h
```

Notice the `DEPRECATED` column. Here, `true` means that this PostgresVersion is deprecated for current KubeDB version. KubeDB will not work for deprecated PostgresVersion.

In this tutorial, we will use `10.2-v2` PostgresVersion crd to create PostgreSQL database. To know more about what is `PostgresVersion` crd and why there is `10.2` and `10.2-v2` variation, please visit [here](/docs/concepts/catalog/postgres.md). You can also see supported PostgresVersion [here](/docs/guides/postgres/README.md#supported-postgresversion-crd).

## Create a PostgreSQL database

KubeDB implements a Postgres CRD to define the specification of a PostgreSQL database.

Below is the Postgres object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: quick-postgres
  namespace: demo
spec:
  version: "10.2-v2"
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

- `spec.version` is name of the PostgresVersion crd where the docker images are specified. In this tutorial, a PostgreSQL 10.2 database is created.
- `spec.storageType` specifies the type of storage that will be used for Postgres database. It can be `Durable` or `Ephemeral`. Default value of this field is `Durable`. If `Ephemeral` is used then KubeDB will create Postgres database using `EmptyDir` volume. In this case, you don't have to specify `spec.storage` field. This is useful for testing purpose.
- `spec.storage` specifies the size and StorageClass of PVC that will be dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If you don't specify `spec.storageType: Ephemeral`, then this field is required.
- `spec.terminationPolicy` specifies what KubeDB should do when user try to delete Postgres crd. Termination policy `DoNotTerminate` prevents a user from deleting this object if admission webhook is enabled.

>Note: `spec.storage` section is used to create PVC for database pod. It will create PVC with storage size specified in`storage.resources.requests` field. Don't specify `limits` here. PVC does not get resized automatically.

Let's create Postgres crd,

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/doc-upd-mrf/docs/examples/postgres/quickstart/quick-postgres.yaml
postgres.kubedb.com/quick-postgres created
```

KubeDB operator watches for Postgres objects using Kubernetes api. When a Postgres object is created, KubeDB operator will create a new StatefulSet and two ClusterIP Service with the matching name. KubeDB operator will also create a governing service for StatefulSet with the name `kubedb`, if one is not already present.

If you are using RBAC enabled cluster, PostgreSQL specific RBAC permission is required. For details, please visit [here](/docs/guides/postgres/quickstart/rbac.md).

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created.

```console
$  kubectl get pg -n demo quick-postgres -o wide
NAME             VERSION   STATUS     AGE
quick-postgres   10.2-v2   Creating   13s
```

Let's describe Postgres object `quick-postgres`

```console
$ kubedb describe pg -n demo quick-postgres
Name:               quick-postgres
Namespace:          demo
CreationTimestamp:  Thu, 07 Feb 2019 17:03:11 +0600
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
  Name:               quick-postgres
  CreationTimestamp:  Thu, 07 Feb 2019 17:03:11 +0600
  Labels:               kubedb.com/kind=Postgres
                        kubedb.com/name=quick-postgres
  Annotations:        <none>
  Replicas:           824641589664 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         quick-postgres
  Labels:         kubedb.com/kind=Postgres
                  kubedb.com/name=quick-postgres
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.100.86.27
  Port:         api  5432/TCP
  TargetPort:   api/TCP
  Endpoints:    172.17.0.8:5432

Service:        
  Name:         quick-postgres-replicas
  Labels:         kubedb.com/kind=Postgres
                  kubedb.com/name=quick-postgres
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.103.133.93
  Port:         api  5432/TCP
  TargetPort:   api/TCP
  Endpoints:    172.17.0.8:5432

Database Secret:
  Name:         quick-postgres-auth
  Labels:         kubedb.com/kind=Postgres
                  kubedb.com/name=quick-postgres
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  POSTGRES_PASSWORD:  16 bytes
  POSTGRES_USER:      8 bytes

Topology:
  Type     Pod               StartTime                      Phase
  ----     ---               ---------                      -----
  primary  quick-postgres-0  2019-02-07 17:03:12 +0600 +06  Running

No Snapshots.

Events:
  Type    Reason      Age   From             Message
  ----    ------      ----  ----             -------
  Normal  Successful  51s   KubeDB operator  Successfully created Service
  Normal  Successful  51s   KubeDB operator  Successfully created Service
  Normal  Successful  25s   KubeDB operator  Successfully created StatefulSet
  Normal  Successful  25s   KubeDB operator  Successfully created Postgres
  Normal  Successful  25s   KubeDB operator  Successfully created appbinding
  Normal  Successful  25s   KubeDB operator  Successfully patched StatefulSet
  Normal  Successful  25s   KubeDB operator  Successfully patched Postgres
```

KubeDB has created two services for the Postgres object.

```console
$ kubectl get service -n demo --selector=kubedb.com/kind=Postgres,kubedb.com/name=quick-postgres
NAME                      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
quick-postgres            ClusterIP   10.108.152.107   <none>        5432/TCP   3m
quick-postgres-replicas   ClusterIP   10.105.175.166   <none>        5432/TCP   3m
```

Here,

- Service *`quick-postgres`* targets only one Pod which is acting as *primary* server
- Service *`quick-postgres-replicas`* targets all Pods created by StatefulSet

KubeDB supports PostgreSQL clustering where Pod can be either *primary* or *standby*. To learn how to configure highly available PostgreSQL cluster, click [here](/docs/guides/postgres/clustering/ha_cluster.md).

Here, we have created a PostgreSQL database with single node, *primary* only.

## Connect with PostgreSQL database

KubeDB operator has created a new Secret called `quick-postgres-auth` for storing the *username* and *password* for `postgres` database.

```yaml
 $ kubectl get secret -n demo quick-postgres-auth -o yaml
apiVersion: v1
data:
  POSTGRES_PASSWORD: REQ4aTU2VUJJY3M2M1BWTw==
  POSTGRES_USER: cG9zdGdyZXM=
kind: Secret
metadata:
  creationTimestamp: 2018-09-03T11:25:39Z
  labels:
    kubedb.com/kind: Postgres
    kubedb.com/name: quick-postgres
  name: quick-postgres-auth
  namespace: demo
  resourceVersion: "1677"
  selfLink: /api/v1/namespaces/demo/secrets/quick-postgres-auth
  uid: 15b3e8a1-af6c-11e8-996d-0800270d7bae
type: Opaque
```

This secret contains superuser name for `postgres` database as `POSTGRES_USER` key and
password as `POSTGRES_PASSWORD` key. By default, superuser name is `postgres` and password is randomly generated.

If you want to use custom password, please create the secret manually and specify that when creating the Postgres object using `spec.databaseSecret.secretName`. For more details see [here](/docs/concepts/databases/postgres.md#specdatabasesecret).

> Note: Auth Secret name format: `{postgres-name}-auth`

Now, you can connect to this database from the pgAdmin dashboard using `quick-postgres.demo` service and *username* and *password* created in `quick-postgres-auth` secret.

**Connection information:**

- Host name/address: you can use any of these
  - Service: `quick-postgres.demo`
  - Pod IP: (`$ kubectl get pods quick-postgres-0 -n demo -o yaml | grep podIP`)
- Port: `5432`
- Maintenance database: `postgres`

- Username: Run following command to get *username*,

  ```console
  $ kubectl get secrets -n demo quick-postgres-auth -o jsonpath='{.data.\POSTGRES_USER}' | base64 -d
  postgres
  ```

- Password: Run the following command to get *password*,

  ```console
  $ kubectl get secrets -n demo quick-postgres-auth -o jsonpath='{.data.\POSTGRES_PASSWORD}' | base64 -d
  DD8i56UBIcs63PVO
  ```

Now, go to pgAdmin dashboard and connect to the database using the connection information as shown below,

<p align="center">
  <kbd>
    <img alt="quick-postgres"  src="/docs/images/postgres/quick-postgres.gif">
  </kbd>
</p>

## Pause Database

KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `DoNotTerminate` termination policy. If admission webhook is enabled, it prevents user from deleting the database as long as the `spec.terminationPolicy` is set `DoNotTerminate`.

In this tutorial, Postgres `quick-postgres` is created with `spec.terminationPolicy: DoNotTerminate`. So if you try to delete this Postgres object, admission webhook will nullify the delete operation.

```console
$  kubedb delete pg -n demo quick-postgres
Error from server (BadRequest): admission webhook "postgres.validators.kubedb.com" denied the request: postgres "quick-postgres" can't be paused. To delete, change spec.terminationPolicy
```

To pause the database, we have to set `spec.terminationPolicy:` to `Pause` by updating it,

```console
$ kubectl edit pg -n demo quick-postgres
spec:
  terminationPolicy: Pause
```

Now, if you delete the Postgres object, KubeDB operator will create a matching DormantDatabase object. This DormantDatabase object can be used to resume the database. KubeDB operator will delete the StatefulSet and its Pods but leaves the Secret, PVCs unchanged.

Let's delete the Postgres object,

```console
$ kubectl delete pg -n demo quick-postgres
postgres.kubedb.com "quick-postgres" deleted
```

Check DormantDatabase has been created successfully,

```console
$ kubectl get drmn -n demo quick-postgres
NAME             STATUS    AGE
quick-postgres   Paused    5m
```

In KubeDB parlance, we say that Postgres `quick-postgres`  has entered into the dormant state.

Let's see, what we have in this DormantDatabase object

```yaml
$ kubectl get drmn -n demo quick-postgres -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: "2019-02-07T11:05:44Z"
  finalizers:
  - kubedb.com
  generation: 1
  labels:
    kubedb.com/kind: Postgres
  name: quick-postgres
  namespace: demo
  resourceVersion: "39020"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/quick-postgres
  uid: 50a9c42e-2ac8-11e9-9d44-080027154f61
spec:
  origin:
    metadata:
      creationTimestamp: "2019-02-07T11:03:11Z"
      name: quick-postgres
      namespace: demo
    spec:
      postgres:
        databaseSecret:
          secretName: quick-postgres-auth
        leaderElection:
          leaseDurationSeconds: 15
          renewDeadlineSeconds: 10
          retryPeriodSeconds: 2
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
          dataSource: null
          resources:
            requests:
              storage: 1Gi
          storageClassName: standard
        storageType: Durable
        terminationPolicy: Pause
        updateStrategy:
          type: RollingUpdate
        version: 10.2-v2
status:
  observedGeneration: 1$8378748355133368567
  pausingTime: "2019-02-07T11:05:56Z"
  phase: Paused
```

Here,

- `spec.origin` contains original Postgres object.
- `status.phase` points to the current database state `Paused`.

## Resume DormantDatabase

To resume the database from the dormant state, create same Postgres object with same Spec.

In this tutorial, the DormantDatabase `quick-postgres` can be resumed by creating original Postgres object.

Let's create the original Postgres object,

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/doc-upd-mrf/docs/examples/postgres/quickstart/quick-postgres.yaml
postgres.kubedb.com/quick-postgres created
```

This will resume the previous database. All data that was inserted in previous database will be available again.

When the database is resumed, respective DormantDatabase object will be removed. Verify that the DormantDatabase object has been removed,

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/doc-upd-mrf/docs/examples/postgres/quickstart/quick-postgres.yaml
postgres "quick-postgres" created
```

## WipeOut DormantDatabase

You can wipe out a DormantDatabase while deleting the object by setting `spec.wipeOut` to true. KubeDB operator will delete any relevant resources of this `PostgresSQL` database (i.e, PVCs, Secrets, Snapshots). It will also delete snapshot data stored in the Cloud Storage buckets.

```yaml
$ kubectl edit drmn -n demo quick-postgres
spec:
  wipeOut: true
```

If `spec.wipeOut` is not set to true while deleting the `dormantdatabase` object, then only this object will be deleted and `kubedb-operator` won't delete related Secrets, PVCs, and Snapshots. So, users still can access the stored data in the cloud storage buckets as well as PVCs.

## Delete DormantDatabase

As it is already discussed above, `DormantDatabase` can be deleted with or without wiping out the resources. To delete the `dormantdatabase`,

```console
$ kubectl delete drmn -n demo quick-postgres
dormantdatabase.kubedb.com "quick-postgres" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo pg/quick-postgres -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo pg/quick-postgres

kubectl delete ns demo
```

## Tips for Testing

If you are just testing some basic functionalities, you might want to avoid additional hassles due to some safety features that are great for production environment. You can follow these tips to avoid them.

1. **Use `storageType: Ephemeral`**. Databases are precious. You might not want to lose your data in your production environment if database pod fail. So, we recommend to use `spec.storageType: Durable` and provide storage spec in `spec.storage` section. For testing purpose, you can just use `spec.storageType: Ephemeral`. KubeDB will use [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) for storage. You will not require to provide `spec.storage` section.
2. **Use `terminationPolicy: WipeOut`**. It is nice to be able to resume database from previous one. So, we create `DormantDatabase` and preserve all your `PVCs`, `Secrets`, `Snapshots` etc. If you don't want to resume database, you can just use `spec.terminationPolicy: WipeOut`. It will not create `DormantDatabase` and it will delete everything created by KubeDB for a particular Postgres crd when you delete the crd. For more details about termination policy, please visit [here](/docs/concepts/databases/postgres.md#specterminationpolicy).

## Next Steps

- Learn about [taking instant backup](/docs/guides/postgres/snapshot/instant_backup.md) of PostgreSQL database using KubeDB Snapshot.
- Learn how to [schedule backup](/docs/guides/postgres/snapshot/scheduled_backup.md)  of PostgreSQL database.
- Learn about initializing [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Learn about initializing [PostgreSQL from KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- Learn about [custom PostgresVersions](/docs/guides/postgres/custom-versions/setup.md).
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Detail concepts of [Postgres object](/docs/concepts/databases/postgres.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Use [private Docker registry](/docs/guides/postgres/private-registry/using-private-registry.md) to deploy PostgreSQL with KubeDB.
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
