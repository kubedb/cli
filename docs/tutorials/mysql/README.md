
> New to KubeDB? Please start [here](/docs/tutorials/README.md).

# Running MySQL
This tutorial will show you how to use KubeDB to run a MySQL database.

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). 

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. This tutorial will also use a phpMyAdmin to connect and test MySQL database, once it is running. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f ./docs/examples/mysql/demo-0.yaml
namespace "demo" created
deployment "myadmin" created
service "myadmin" created


$ kubectl get pods -n demo --watch
NAME                      READY     STATUS              RESTARTS   AGE
myadmin-fccf65985-ppgbh   0/1       ContainerCreating   0          9s
myadmin-fccf65985-ppgbh   1/1       Running   			0         40s


$ kubectl get service -n demo
NAME      CLUSTER-IP   EXTERNAL-IP   PORT(S)        AGE
pgadmin   10.0.0.92    <pending>     80:31188/TCP   1m


$ minikube ip
192.168.99.100
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{myadmin-svc-nodeport}_. 
You can Also get this URl by running the following command: 
```console
$ minikube service myadmin -n demo --url
http://192.168.99.100:31833
```
According to the above example, this URL will be [http://192.168.99.100:31833](http://192.168.99.100:31833). The logging informations to phpMyAdmin _(host, username and password)_ will be retrieved later in this tutorial.

## Create a MySQL database
KubeDB implements a `MySQL` CRD to define the specification of a MySQL database. Below is the `MySQL` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: m1
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
        repository: "https://github.com/the-redback/mysql-init-script.git"
        directory: .


$ kubedb create -f ./docs/examples/mysql/demo-1.yaml
validating "./docs/examples/mysql/demo-1.yaml"
mysql "m1" created
```

Here,
 - `spec.version` is the version of MySQL database. In this tutorial, a MySQL 8.0 database is going to be created.

 - `spec.doNotPause` tells KubeDB operator that if this tpr is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.

 - `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

 - `spec.init.scriptSource` specifies a sql script source used to initialize the database after it is created. The sql scripts will be executed alphabatically. In this tutorial, a sample sql script from the git repository `https://github.com/kubedb/mysql-init-scripts.git` is used to create a `dashboard` table in _test_ database.

KubeDB operator watches for `MySQL` objects using Kubernetes api. When a `MySQL` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching tpr name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. If [RBAC is enabled](/docs/tutorials/rbac.md), a ClusterRole, ServiceAccount and ClusterRoleBinding with the matching tpr name will be created and used as the service account name for the corresponding StatefulSet.

```console
$ kubedb describe ms -n demo m1
Name:		m1
Namespace:	demo
StartTimestamp:	Fri, 08 Dec 2017 12:26:30 +0600
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

StatefulSet:		
  Name:			m1
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Fri, 08 Dec 2017 12:26:34 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:	
  Name:		m1
  Type:		ClusterIP
  IP:		10.97.157.102
  Port:		db	3306/TCP

Database Secret:
  Name:	m1-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	16 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From             Type       Reason               Message
  ---------   --------   -----     ----             --------   ------               -------
  8m          8m         1         mysql operator   Normal     SuccessfulValidate   Successfully validate MySQL
  8m          8m         1         mysql operator   Normal     SuccessfulValidate   Successfully validate MySQL
  8m          8m         1         mysql operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  8m          8m         1         mysql operator   Normal     SuccessfulCreate     Successfully created MySQL
  8m          8m         1         mysql operator   Normal     SuccessfulValidate   Successfully validate MySQL
  8m          8m         1         mysql operator   Normal     Creating             Creating Kubernetes objects


$ kubectl get statefulset -n demo
NAME      DESIRED   CURRENT   AGE
m1        1         1         9m


$ kubectl get pvc -n demo
NAME        STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-m1-0   Bound     pvc-bd074404-dbe0-11e7-983b-08002790ff42   50Mi       RWO            standard       9m


$ kubectl get pv -n demo
NAME        STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-m1-0   Bound     pvc-bd074404-dbe0-11e7-983b-08002790ff42   50Mi       RWO            standard       9m


$ kubectl get service -n demo
NAME      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
kubedb    ClusterIP      None             <none>        <none>         28m
m1        ClusterIP      10.97.157.102    <none>        3306/TCP       10m
myadmin   LoadBalancer   10.104.146.121   <pending>     80:31833/TCP   44m
```


KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified tpr:

```yaml
$ kubedb get ms -n demo m1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-08T06:26:30Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  name: m1
  namespace: demo
  resourceVersion: "2727"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/mysqls/m1
  uid: ba353d07-dbe0-11e7-983b-08002790ff42
spec:
  databaseSecret:
    secretName: m1-admin-auth
  doNotPause: true
  init:
    scriptSource:
      gitRepo:
        directory: .
        repository: https://github.com/the-redback/mysql-init-script.git
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: 8
status:
  creationTime: 2017-12-08T06:26:30Z
  phase: Running
```

Please note that KubeDB operator has created a new Secret called `m1-admin-auth` (format: {tpr-name}-admin-auth) for storing the password for `mysql` superuser. This secret contains a `.admin` key with a ini formatted key-value pairs. If you want to use an existing secret please specify that when creating the tpr using `spec.databaseSecret.secretName`.

Now, you can connect to this database from the phpMyAdmin dashboard using the database pod IP and `mysql` user password. 
```console
$ kubectl get pods m1-0 -n demo -o yaml | grep IP
  hostIP: 192.168.99.100
  podIP: 172.17.0.6

$ kubectl get secrets -n demo m1-admin-auth -o jsonpath='{.data.\.admin}' | base64 -d
nMdePjKp4vP90AQF
```
Now, open your browser and go to the following URL: _http://{minikube-ip}:{pgadmin-svc-nodeport}_. To log into the phpMyAdmin, use host __`172.17.0.6`__ , username __`root`__ and password __`nMdePjKp4vP90AQF`__.


---Rest of the doc