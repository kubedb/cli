---
title: Initialize MySQL using Script
menu:
  docs_0.8.0-beta.2:
    identifier: my-using-script-initialization
    name: Using Script
    parent: my-initialization-mysql
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Initialize MySQL using Script

This tutorial will show you how to use KubeDB to initialize a MySQL database with \*.sql, \*.sh and/or \*.sql.gz script.
In this tutorial we will use .sql script stored in GitHub repository [kubedb/mysql-init-scripts](https://github.com/kubedb/mysql-init-scripts).

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. This tutorial will also use a phpMyAdmin to connect and test MySQL database, once it is running. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/mysql/demo-0.yaml
namespace "demo" created

$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/mysql/quickstart/demo-1.yaml
deployment "myadmin" created
service "myadmin" created

$ kubectl get pods -n demo --watch
NAME                      READY     STATUS              RESTARTS   AGE
myadmin-c4db4df95-4tgkx   0/1       ContainerCreating   0          27s
myadmin-c4db4df95-4tgkx   1/1       Running             0          1m

$ kubectl get service -n demo
NAME                TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
myadmin             LoadBalancer   10.101.247.127   <pending>     80:32673/TCP   50s

$ minikube ip
192.168.99.100
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{myadmin-svc-nodeport}_.
You can also get this URl by running the following command:

```console
$ minikube service myadmin -n demo --url
http://192.168.99.100:32673
```

According to the above example, this URL will be [http://192.168.99.100:32673](http://192.168.99.100:32673). The login informations to phpMyAdmin _(host, username and password)_ will be retrieved later in this tutorial.

## Create a MySQL database with Init-Script

Below is the `MySQL` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-init-script
  namespace: demo
spec:
  version: 8.0
  doNotPause: true
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    scriptSource:
      gitRepo:
        repository: "https://github.com/kubedb/mysql-init-scripts.git"
        directory: .

```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/mysql/Initialization/demo-1.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/mysql/Initialization/demo-1.yaml"
mysql "mysql-init-script" created
```

Here,

- `spec.version` is the version of MySQL database. In this tutorial, a MySQL 8.0 database is going to be created.
- `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.
- `spec.init.scriptSource` specifies a script source used to initialize the database before database server starts. The scripts will be executed alphabatically. In this tutorial, a sample .js script from the git repository `https://github.com/kubedb/mysql-init-scripts.git` is used to create a test database. You can use other [volume sources](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes) instead of `gitrepo`.  The \*.sql, \*sql.gz and/or \*.sh sripts that are stored inside the root folder will be executed alphabatically. The scripts inside child folders will be skipped.

KubeDB operator watches for `MySQL` objects using Kubernetes api. When a `MySQL` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching MySQL object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. No MySQL specific RBAC roles are required for [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb describe my -n demo mysql-init-script
Name:		mysql-init-script
Namespace:	demo
StartTimestamp:	Fri, 09 Feb 2018 17:18:14 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:
  Name:			mysql-init-script
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Fri, 09 Feb 2018 17:18:15 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		mysql-init-script
  Type:		ClusterIP
  IP:		10.101.136.66
  Port:		db	3306/TCP

Database Secret:
  Name:	mysql-init-script-auth
  Type:	Opaque
  Data
  ====
  password:	16 bytes
  user:		4 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From             Type       Reason       Message
  ---------   --------   -----     ----             --------   ------       -------
  9m          9m         1         MySQL operator   Normal     Successful   Successfully patched StatefulSet
  9m          9m         1         MySQL operator   Normal     Successful   Successfully patched MySQL
  9m          9m         1         MySQL operator   Normal     Successful   Successfully created StatefulSet
  9m          9m         1         MySQL operator   Normal     Successful   Successfully created MySQL
  9m          9m         1         MySQL operator   Normal     Successful   Successfully created Service



$ kubectl get statefulset -n demo
NAME                DESIRED   CURRENT   AGE
mysql-init-script   1         1         10m


$ kubectl get pvc -n demo
NAME                       STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-mysql-init-script-0   Bound     pvc-ec9fd7b2-0d8a-11e8-9091-08002751ae8c   50Mi       RWO            standard       12m


$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                           STORAGECLASS   REASON    AGE
pvc-ec9fd7b2-0d8a-11e8-9091-08002751ae8c   50Mi       RWO            Delete           Bound     demo/data-mysql-init-script-0   standard                 12m


$ kubectl get service -n demo
NAME                TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
kubedb              ClusterIP   None            <none>        <none>     13m
mysql-init-script   ClusterIP   10.101.136.66   <none>        3306/TCP   13m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified MySQL object:

```yaml
$ kubedb get my -n demo mysql-init-script -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-09T11:18:14Z
  finalizers:
  - kubedb.com
  generation: 0
  name: mysql-init-script
  namespace: demo
  resourceVersion: "28624"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mysqls/mysql-init-script
  uid: ebbcc002-0d8a-11e8-9091-08002751ae8c
spec:
  databaseSecret:
    secretName: mysql-init-script-auth
  doNotPause: true
  init:
    scriptSource:
      gitRepo:
        directory: .
        repository: https://github.com/kubedb/mysql-init-scripts.git
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: 8
status:
  creationTime: 2018-02-09T11:18:14Z
  phase: Running
```

Please note that KubeDB operator has created a new Secret called `mysql-init-script-auth` *(format: {mysql-object-name}-auth)* for storing the password for MySQL superuser. This secret contains a `user` key which contains the *username* for MySQL superuser and a `password` key which contains the *password* for MySQL superuser.
If you want to use an existing secret please specify that when creating the MySQL object using `spec.databaseSecret.secretName`. While creating this secret manually, make sure the secret contains these two keys containing data `user` and `password` and also make sure of using `root` as value of `user`.

Now, you can connect to this database from the phpMyAdmin dashboard using the database pod IP and and `mysql` user password.

```console
$ kubectl get pods mysql-init-script-0 -n demo -o yaml | grep IP
  hostIP: 192.168.99.100
  podIP: 172.17.0.5

$ kubectl get secrets -n demo mysql-init-script-auth -o jsonpath='{.data.\user}' | base64 -d
root

$ kubectl get secrets -n demo mysql-init-script-auth -o jsonpath='{.data.\password}' | base64 -d
h1sPb6ZTHQmKC1ng
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{myadmin-svc-nodeport}_. To log into the phpMyAdmin, use host __`172.17.0.5`__ , username __`root`__ and password __`h1sPb6ZTHQmKC1ng`__.

As you can see here, the initial script has successfully created a table named `kubedb_table` in `mysql` database and inserted three rows of data into that table successfully.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete my,drmn,snap -n demo --all --force

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Initialize [MySQL with Snapshot](/docs/guides/mysql/initialization/using-snapshot.md).
- [Snapshot and Restore](/docs/guides/mysql/snapshot/backup-and-restore.md) process of MySQL databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mysql/snapshot/scheduled-backup.md) of MySQL databases using KubeDB.
- Initialize [MySQL with Script](/docs/guides/mysql/initialization/using-script.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
