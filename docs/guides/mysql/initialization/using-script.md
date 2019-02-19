---
title: Initialize MySQL using Script
menu:
  docs_0.9.0:
    identifier: my-using-script-initialization
    name: Using Script
    parent: my-initialization-mysql
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Initialize MySQL using Script

This tutorial will show you how to use KubeDB to initialize a MySQL database with \*.sql, \*.sh and/or \*.sql.gz script.
In this tutorial we will use .sql script stored in GitHub repository [kubedb/mysql-init-scripts](https://github.com/kubedb/mysql-init-scripts).

> Note: The yaml files that are used in this tutorial are stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. This tutorial will also use a phpMyAdmin to connect and test MySQL database, once it is running. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace/demo created
  
  $ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/quickstart/demo-1.yaml
  deployment.extensions/myadmin created
  service/myadmin created
  
  $ kubectl get pods -n demo
  NAME                       READY     STATUS    RESTARTS   AGE
  myadmin-584d845666-rtkzg   1/1       Running   0          9m
  
  $ kubectl get service -n demo
  NAME      TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
  myadmin   LoadBalancer   10.108.49.82   <pending>     80:30192/TCP   9m
  
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
  
## Prepare Initialization Scripts

MySQL supports initialization with `.sh`, `.sql` and `.sql.gz` files. In this tutorial, we will use `init.sql` script from [mysql-init-scripts](https://github.com/kubedb/mysql-init-scripts) git repository to create a TABLE `kubedb_table` in `mysql` database.

As [gitRepo](https://kubernetes.io/docs/concepts/storage/volumes/#gitrepo) volume has been deprecated, we will use a ConfigMap as script source. You can use any Kubernetes supported [volume](https://kubernetes.io/docs/concepts/storage/volumes) as script source.

At first, we will create a ConfigMap from `init.sql` file. Then, we will provide this ConfigMap as script source in `init.scriptSource` of MySQL crd spec.

Let's create a ConfigMap with initialization script,

```console
$ kubectl create configmap -n demo my-init-script \
--from-literal=init.sql="$(curl -fsSL https://raw.githubusercontent.com/kubedb/mysql-init-scripts/master/init.sql)"
configmap/my-init-script created
```

## Create a MySQL database with Init-Script

Below is the `MySQL` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-init-script
  namespace: demo
spec:
  version: "8.0-v2"
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  init:
    scriptSource:
      configMap:
        name: my-init-script
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/Initialization/demo-1.yaml
mysql.kubedb.com/mysql-init-script created
```

Here,

- `spec.init.scriptSource` specifies a script source used to initialize the database before database server starts. The scripts will be executed alphabatically. In this tutorial, a sample .sql script from the git repository `https://github.com/kubedb/mysql-init-scripts.git` is used to create a test database. You can use other [volume sources](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes) instead of `ConfigMap`.  The \*.sql, \*sql.gz and/or \*.sh sripts that are stored inside the root folder will be executed alphabatically. The scripts inside child folders will be skipped.

KubeDB operator watches for `MySQL` objects using Kubernetes api. When a `MySQL` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching MySQL object name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. No MySQL specific RBAC roles are required for [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb describe my -n demo mysql-init-script
Name:               mysql-init-script
Namespace:          demo
CreationTimestamp:  Thu, 27 Sep 2018 17:06:37 +0600
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
  Name:               mysql-init-script
  CreationTimestamp:  Thu, 27 Sep 2018 17:06:39 +0600
  Labels:               kubedb.com/kind=MySQL
                        kubedb.com/name=mysql-init-script
  Annotations:        <none>
  Replicas:           824637787500 desired | 1 total
  Pods Status:        1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:         mysql-init-script
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-init-script
  Annotations:  <none>
  Type:         ClusterIP
  IP:           10.102.60.242
  Port:         db  3306/TCP
  TargetPort:   db/TCP
  Endpoints:    172.17.0.6:3306

Database Secret:
  Name:         mysql-init-script-auth
  Labels:         kubedb.com/kind=MySQL
                  kubedb.com/name=mysql-init-script
  Annotations:  <none>
  
Type:  Opaque
  
Data
====
  password:  16 bytes
  user:      4 bytes

No Snapshots.

Events:
  Type    Reason      Age   From            Message
  ----    ------      ----  ----            -------
  Normal  Successful  1m    MySQL operator  Successfully created Service
  Normal  Successful  41s   MySQL operator  Successfully created StatefulSet
  Normal  Successful  41s   MySQL operator  Successfully created MySQL
  Normal  Successful  40s   MySQL operator  Successfully patched StatefulSet
  Normal  Successful  40s   MySQL operator  Successfully patched MySQL
  Normal  Successful  37s   MySQL operator  Successfully patched StatefulSet
  Normal  Successful  37s   MySQL operator  Successfully patched MySQL

$ kubectl get statefulset -n demo
NAME                DESIRED   CURRENT   AGE
mysql-init-script   1         1         1m

$ kubectl get pvc -n demo
NAME                       STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-mysql-init-script-0   Bound     pvc-68e49ec6-c245-11e8-b2cc-080027d9f35e   1Gi        RWO            standard       1m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                           STORAGECLASS   REASON    AGE
pvc-68e49ec6-c245-11e8-b2cc-080027d9f35e   1Gi        RWO            Delete           Bound     demo/data-mysql-init-script-0   standard                 1m

$ kubectl get service -n demo
NAME                TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
kubedb              ClusterIP      None            <none>        <none>         2m
myadmin             LoadBalancer   10.108.49.82    <pending>     80:30192/TCP   22m
mysql-init-script   ClusterIP      10.102.60.242   <none>        3306/TCP       2m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified MySQL object:

```yaml
$ kubedb get my -n demo mysql-init-script -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  creationTimestamp: 2018-09-27T11:06:37Z
  finalizers:
  - kubedb.com
  generation: 2
  name: mysql-init-script
  namespace: demo
  resourceVersion: "9070"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mysqls/mysql-init-script
  uid: 677921ee-c245-11e8-b2cc-080027d9f35e
spec:
  databaseSecret:
    secretName: mysql-init-script-auth
  init:
    scriptSource:
      configMap:
        name: my-init-script
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
  version: 8.0-v2
status:
  observedGeneration: 2$4213139756412538772
  phase: Running
```

KubeDB operator has created a new Secret called `mysql-init-script-auth` *(format: {mysql-object-name}-auth)* for storing the password for MySQL superuser. This secret contains a `username` key which contains the *username* for MySQL superuser and a `password` key which contains the *password* for MySQL superuser.
If you want to use an existing secret please specify that when creating the MySQL object using `spec.databaseSecret.secretName`. While creating this secret manually, make sure the secret contains these two keys containing data `username` and `password` and also make sure of using `root` as value of `username`.

Now, you can connect to this database from the phpMyAdmin dashboard using the database pod IP and and `mysql` user password.

```console
$ kubectl get pods mysql-init-script-0 -n demo -o yaml | grep IP
  hostIP: 10.0.2.15
  podIP: 172.17.0.6

$ kubectl get secrets -n demo mysql-init-script-auth -o jsonpath='{.data.\user}' | base64 -d
root

$ kubectl get secrets -n demo mysql-init-script-auth -o jsonpath='{.data.\password}' | base64 -d
1Pc7bwSygrv1MX1Q
```

---
Note: In MySQL:8.0-v2 (ie, 8.0.14), connection to phpMyAdmin may give error as it is using `caching_sha2_password` and `sha256_password` authentication plugins over `mysql_native_password`. If the error happens do the following for work around. But, It's not recommended to change authentication plugins. See [here](https://stackoverflow.com/questions/49948350/phpmyadmin-on-mysql-8-0) for alternative solutions.

```console
kubectl exec -it -n demo mysql-quickstart-0 -- mysql -u root --password=1Pc7bwSygrv1MX1Q -e "ALTER USER root IDENTIFIED WITH mysql_native_password BY '1Pc7bwSygrv1MX1Q';"
```

---

Now, open your browser and go to the following URL: _http://{minikube-ip}:{myadmin-svc-nodeport}_. To log into the phpMyAdmin, use host __`172.17.0.6`__ , username __`root`__ and password __`1Pc7bwSygrv1MX1Q`__.

As you can see here, the initial script has successfully created a table named `kubedb_table` in `mysql` database and inserted three rows of data into that table successfully.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mysql/mysql-init-script -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mysql/mysql-init-script

kubectl patch -n demo drmn/mysql-init-script -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mysql-init-script

kubectl delete ns demo
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
- Detail concepts of [MySQLVersion object](/docs/concepts/catalog/mysql.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
