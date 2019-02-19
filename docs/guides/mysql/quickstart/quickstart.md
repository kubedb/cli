---
title: MySQL Quickstart
menu:
  docs_0.9.0:
    identifier: my-quickstart-quickstart
    name: Overview
    parent: my-quickstart-mysql
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# MySQL QuickStart

This tutorial will show you how to use KubeDB to run a MySQL database.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/mysql/mysql-lifecycle.png">
</p>

> Note: The yaml files used in this tutorial are stored in [docs/examples/mysql](https://github.com/kubedb/cli/tree/master/docs/examples/mysql) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/) is required to run KubeDB. Check the available StorageClass in cluster.

  ```console
  $ kubectl get storageclasses
  NAME                 PROVISIONER                AGE
  standard (default)   k8s.io/minikube-hostpath   4h
  ```

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. This tutorial will also use a [phpMyAdmin](https://hub.docker.com/r/phpmyadmin/phpmyadmin/) deployment to connect and test MySQL database, once it is running. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace/demo created
  
  $ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/quickstart/demo-1.yaml
  deployment.extensions/myadmin created
  service/myadmin created
  
  $ kubectl get pods -n demo --watch
  NAME                      READY     STATUS              RESTARTS   AGE
  myadmin-c4db4df95-8lk74   0/1       ContainerCreating   0          27s
  myadmin-c4db4df95-8lk74   1/1       Running             0          1m
  
  $ kubectl get service -n demo
  NAME      TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
  myadmin   LoadBalancer   10.105.73.16   <pending>     80:30158/TCP   23m
  
  $ minikube ip
  192.168.99.100
  ```

  Now, open your browser and go to the following URL: _http://{minikube-ip}:{myadmin-svc-nodeport}_.
  You can also get this URl by running the following command:

  ```console
  $ minikube service myadmin -n demo --url
  http://192.168.99.100:30158
  ```

According to the above example, this URL will be [http://192.168.99.100:30158](http://192.168.99.100:30158). The login informations to phpMyAdmin _(host, username and password)_ will be retrieved later in this tutorial.

## Find Available MySQLVersion

When you have installed KubeDB, it has created `MySQLVersion` crd for all supported MySQL versions. Check 0

```console
$ kubectl get mysqlversions
NAME     VERSION   DB_IMAGE                  DEPRECATED   AGE
5        5         kubedb/mysql:5        true         29s
5-v1     5         kubedb/mysql:5-v1     true         29s
5.7      5.7       kubedb/mysql:5.7      true         29s
5.7-v1   5.7       kubedb/mysql:5.7-v1                29s
8        8         kubedb/mysql:8        true         29s
8-v1     8         kubedb/mysql:8-v1     true         29s
8.0      8.0       kubedb/mysql:8.0      true         29s
8.0-v1   8.0       kubedb/mysql:8.0-v1                29s
8.0-v2   8.0       kubedb/mysql:8.0-v2                29s
8.0.14   8.0.14    kubedb/mysql:8.0.14                29s
8.0.3    8.0.3     kubedb/mysql:8.0.3                 29s
```

## Create a MySQL database

KubeDB implements a `MySQL` CRD to define the specification of a MySQL database. Below is the `MySQL` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-quickstart
  namespace: demo
spec:
  version: "8.0-v2"
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

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/quickstart/demo-2.yaml
mysql.kubedb.com/mysql-quickstart created
```

Here,

- `spec.version` is the name of the MySQLVersion CRD where the docker images are specified. In this tutorial, a MySQL 8.0-v2 database is going to be created.
- `spec.storageType` specifies the type of storage that will be used for MySQL database. It can be `Durable` or `Ephemeral`. Default value of this field is `Durable`. If `Ephemeral` is used then KubeDB will create MySQL database using `EmptyDir` volume. In this case, you don't have to specify `spec.storage` field. This is useful for testing purpose.
- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests.
- `spec.terminationPolicy` gives flexibility whether to `nullify`(reject) the delete operation of `MySQL` crd or which resources KubeDB should keep or delete when you delete `MySQL` crd. If admission webhook is enabled, It prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`. Learn details of all `TerminationPolicy` [here](docs/concepts/databases/mysql.md#specterminationpolicy)

> Note: spec.storage section is used to create PVC for database pod. It will create PVC with storage size specified instorage.resources.requests field. Don't specify limits here. PVC does not get resized automatically.

KubeDB operator watches for `MySQL` objects using Kubernetes api. When a `MySQL` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching MySQL object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. No MySQL specific RBAC permission is required in [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb describe my -n demo mysql-quickstart
Name:               mysql-quickstart
Namespace:          demo
CreationTimestamp:  Wed, 06 Feb 2019 17:17:55 +0600
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
  Name:               mysql-quickstart
  CreationTimestamp:  Wed, 06 Feb 2019 17:17:55 +0600
  Labels:               kubedb.com/kind=MySQL
                        kubedb.com/name=mysql-quickstart
  Annotations:        <none>
  Replicas:           824641282668 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:        
  Name:         mysql-quickstart
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-quickstart
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.99.24.103
  Port:         db  3306/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.8:3306

Database Secret:
  Name:         mysql-quickstart-auth
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-quickstart
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
  Normal  Successful  53s   KubeDB operator  Successfully created StatefulSet
  Normal  Successful  53s   KubeDB operator  Successfully created MySQL
  Normal  Successful  53s   KubeDB operator  Successfully created appbinding
  Normal  Successful  53s   KubeDB operator  Successfully patched StatefulSet
  Normal  Successful  53s   KubeDB operator  Successfully patched MySQL

$ kubectl get statefulset -n demo
NAME               READY   AGE
mysql-quickstart   1/1     2m22s

$ kubectl get pvc -n demo
NAME                      STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-mysql-quickstart-0   Bound     pvc-652e02c7-0d7f-11e8-9091-08002751ae8c   1Gi        RWO            standard       10m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                          STORAGECLASS   REASON    AGE
pvc-652e02c7-0d7f-11e8-9091-08002751ae8c   1Gi        RWO            Delete           Bound     demo/data-mysql-quickstart-0   standard                 11m

$ kubectl get service -n demo
NAME               TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
kubedb             ClusterIP      None            <none>        <none>         11m
myadmin            LoadBalancer   10.105.73.16    <pending>     80:30158/TCP   41m
mysql-quickstart   ClusterIP      10.104.50.139   <none>        3306/TCP       11m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified MySQL object:

```yaml
$ kubedb get my -n demo mysql-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  creationTimestamp: "2019-02-06T11:17:55Z"
  finalizers:
  - kubedb.com
  generation: 2
  name: mysql-quickstart
  namespace: demo
  resourceVersion: "2158"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mysqls/mysql-quickstart
  uid: d9b37fc3-2a00-11e9-a088-080027ab5700
spec:
  databaseSecret:
    secretName: mysql-quickstart-auth
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
  terminationPolicy: DoNotTerminate
  updateStrategy:
    type: RollingUpdate
  version: 8.0-v2
status:
  observedGeneration: 2$4213139756412538772
  phase: Running
```

## Connect with MySQL database

KubeDB operator has created a new Secret called `mysql-quickstart-auth` *(format: {mysql-object-name}-auth)* for storing the password for `mysql` superuser. This secret contains a `username` key which contains the *username* for MySQL superuser and a `password` key which contains the *password* for MySQL superuser.

If you want to use an existing secret please specify that when creating the MySQL object using `spec.databaseSecret.secretName`. While creating this secret manually, make sure the secret contains these two keys containing data `username` and `password` and also make sure of using `root` as value of `username`. For more details see [here](/docs/concepts/databases/mysql.md#specdatabasesecret).

Now, you can connect to this database from the phpMyAdmin dashboard using the database pod IP and and `mysql` user password.

```console
$ kubectl get pods mysql-quickstart-0 -n demo -o yaml | grep podIP
  podIP: 172.17.0.6

$ kubectl get secrets -n demo mysql-quickstart-auth -o jsonpath='{.data.\username}' | base64 -d
root

$ kubectl get secrets -n demo mysql-quickstart-auth -o jsonpath='{.data.\password}' | base64 -d
l0yKjI1E7IMohsGR
```

---
Note: In MySQL:8.0-v1 (ie, 8.0.14), connection to phpMyAdmin may give error as it is using `caching_sha2_password` and `sha256_password` authentication plugins over `mysql_native_password`. If the error happens do the following for work around. But, It's not recommended to change authentication plugins. See [here](https://stackoverflow.com/questions/49948350/phpmyadmin-on-mysql-8-0) for alternative solutions.

```console
kubectl exec -it -n demo mysql-quickstart-0 -- mysql -u root --password=l0yKjI1E7IMohsGR -e "ALTER USER root IDENTIFIED WITH mysql_native_password BY 'l0yKjI1E7IMohsGR';"
```
---

Now, open your browser and go to the following URL: _http://{minikube-ip}:{myadmin-svc-nodeport}_. To log into the phpMyAdmin, use host __`mysql-quickstart.demo`__ or __`172.17.0.6`__ , username __`root`__ and password __`pefjWeXoAQ9PaRZv`__.

## DoNotTerminate Property

When, `terminationPolicy` is `DoNotTerminate`, KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `DoNotTerminate` feature. If admission webhook is enabled, It prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`. You can see this below:

```console
$ kubedb delete my mysql-quickstart -n demo
Error from server (BadRequest): admission webhook "mysql.validators.kubedb.com" denied the request: mysql "mysql-quickstart" can't be paused. To delete, change spec.terminationPolicy
```

Now, run `kubedb edit my mysql-quickstart -n demo` to set `spec.terminationPolicy` to `Pause` (which creates `domantdatabase` when mysql is deleted and keeps PVC, snapshots, Secrets intact) or remove this field (which default to `Pause`). Then you will be able to delete/pause the database.

Learn details of all `TerminationPolicy` [here](docs/concepts/databases/mysql.md#specterminationpolicy)

## Pause Database

When [TerminationPolicy](/docs/concepts/databases/mysql.md#specterminationpolicy) is set to `Pause`, it will pause the MySQL database instead of deleting it. Here, If you delete the MySQL object, KubeDB operator will delete the StatefulSet and its pods but leaves the PVCs unchanged. In KubeDB parlance, we say that `mgo-quickstart` MySQL database has entered into the dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase object.

```console
$ kubedb delete my mysql-quickstart -n demo
mysql.kubedb.com "mysql-quickstart" deleted

$ kubedb get drmn -n demo mysql-quickstart
NAME               STATUS    AGE
mysql-quickstart   Pausing   14s

$ kubedb get drmn -n demo mysql-quickstart
NAME               STATUS    AGE
mysql-quickstart   Paused    39s
```

```yaml
$ kubedb get drmn -n demo mysql-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  creationTimestamp: "2019-02-07T09:54:20Z"
  finalizers:
  - kubedb.com
  generation: 1
  labels:
    kubedb.com/kind: MySQL
  name: mysql-quickstart
  namespace: demo
  resourceVersion: "32852"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mysql-quickstart
  uid: 575bf14f-2abe-11e9-9d44-080027154f61
spec:
  origin:
    metadata:
      creationTimestamp: "2019-02-07T09:46:04Z"
      name: mysql-quickstart
      namespace: demo
    spec:
      mysql:
        databaseSecret:
          secretName: mysql-quickstart-auth
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
        version: 8.0-v2
status:
  observedGeneration: 1$5984877185736766566
  pausingTime: "2019-02-07T09:54:24Z"
  phase: Paused
```

Here,

- `spec.origin` is the spec of the original spec of the original MySQL object.
- `status.phase` points to the current database state `Paused`.

## Resume Dormant Database

To resume the database from the dormant state, create same `MySQL` object with same Spec.

In this tutorial, the dormant database can be resumed by creating original `MySQL` object.

The below command will resume the DormantDatabase `mysql-quickstart` that was created before.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/quickstart/demo-2.yaml
mysql.kubedb.com/mysql-quickstart created
```

Now, if you exec into the database, you can see that the datas are intact.

## WipeOut DormantDatabase

You can wipe out a DormantDatabase while deleting the objet by setting `spec.wipeOut` to true. KubeDB operator will delete any relevant resources of this `MySQL` database (i.e, PVCs, Secrets, Snapshots). It will also delete snapshot data stored in the Cloud Storage buckets.

```yaml
$ kubedb delete my mysql-quickstart -n demo
mysql.kubedb.com "mysql-quickstart" deleted

$ kubedb edit drmn -n demo mysql-quickstart
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  name: mysql-quickstart
  namespace: demo
  ...
spec:
  wipeOut: true
  ...
status:
  phase: Paused
  ...
```

If `spec.wipeOut` is not set to true while deleting the `dormantdatabase` object, then only this object will be deleted and `kubedb-operator` won't delete related Secrets, PVCs and Snapshots. So, user still can access the stored data in the cloud storage buckets as well as PVCs.

## Delete DormantDatabase

As it is already discussed above, `DormantDatabase` can be deleted with or without wiping out the resources. To delete the `dormantdatabase`,

```console
$ kubedb delete drmn mysql-quickstart -n demo
dormantdatabase.kubedb.com "mysql-quickstart" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mysql/mysql-quickstart -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mysql/mysql-quickstart

kubectl patch -n demo drmn/mysql-quickstart -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mysql-quickstart

kubectl delete ns demo
```

## Tips for Testing

If you are just testing some basic functionalities, you might want to avoid additional hassles due to some safety features that are great for production environment. You can follow these tips to avoid them.

1. **Use `storageType: Ephemeral`**. Databases are precious. You might not want to lose your data in your production environment if database pod fail. So, we recommend to use `spec.storageType: Durable` and provide storage spec in `spec.storage` section. For testing purpose, you can just use `spec.storageType: Ephemeral`. KubeDB will use [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) for storage. You will not require to provide `spec.storage` section.
2. **Use `terminationPolicy: WipeOut`**. It is nice to be able to resume database from previous one. So, we create `DormantDatabase` and preserve all your `PVCs`, `Secrets`, `Snapshots` etc. If you don't want to resume database, you can just use `spec.terminationPolicy: WipeOut`. It will not create `DormantDatabase` and it will delete everything created by KubeDB for a particular MySQL crd when you delete the crd. For more details about termination policy, please visit [here](/docs/concepts/databases/mysql.md#specterminationpolicy).

## Next Steps

- [Snapshot and Restore](/docs/guides/mysql/snapshot/backup-and-restore.md) process of MySQL databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mysql/snapshot/scheduled-backup.md) of MySQL databases using KubeDB.
- Initialize [MySQL with Script](/docs/guides/mysql/initialization/using-script.md).
- Initialize [MySQL with Snapshot](/docs/guides/mysql/initialization/using-snapshot.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Detail concepts of [MySQLVersion object](/docs/concepts/catalog/mysql.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
