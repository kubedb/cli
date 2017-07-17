# Using PostgreSQL
This tutorial will show you how to use KubeDB to run a PostgreSQL database. 

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). 

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/install.md).

TO keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. This tutorial will also use a PGAdmin to connect and test PostgreSQL database, once it is running. Run the following command to prepare your cluster for this tutorial:

```sh
$ kubectl create -f ./docs/examples/tutorial/postgres/demo-0.yaml 
namespace "demo" created
deployment "pgadmin" created
service "pgadmin" created

$ kubectl get pods -n demo --watch
NAME                      READY     STATUS              RESTARTS   AGE
pgadmin-538449054-s046r   0/1       ContainerCreating   0          13s
pgadmin-538449054-s046r   1/1       Running   0          1m
^C⏎                                                                                                                                                             

$ kubectl get service -n demo
NAME      CLUSTER-IP   EXTERNAL-IP   PORT(S)        AGE
pgadmin   10.0.0.92    <pending>     80:31188/TCP   1m

$ minikube ip
192.168.99.100
```

Now, open your browser and go to the following URL: _http://{minikube-ip}:{pgadmin-svc-nodeport}_. According to the above example, this URL will be [http://192.168.99.100:31188](http://192.168.99.100:31188).

## Create a PostgreSQL database

```yaml
$ kubedb create -f ./docs/examples/tutorial/postgres/demo-1.yaml 
validating "./docs/examples/tutorial/postgres/demo-1.yaml"
postgres "p1" created

apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: p1
  namespace: demo
spec:
  version: 9.5
  doNotPause: true
  storage:
    class: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi      
  init:
    scriptSource:
      scriptPath: "postgres-init-scripts/run.sh"
      gitRepo:
        repository: "https://github.com/k8sdb/postgres-init-scripts.git"
```

```yaml
$ kubedb get pg -n demo p1
NAME      STATUS     AGE
p1        Creating   21s

$ kubedb get pg -n demo p1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  creationTimestamp: 2017-07-17T22:31:34Z
  name: p1
  namespace: demo
  resourceVersion: "2677"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/postgreses/p1
  uid: b02ccec1-6b3f-11e7-bdc0-080027aa4456
spec:
  databaseSecret:
    secretName: p1-admin-auth
  doNotPause: true
  init:
    scriptSource:
      gitRepo:
        repository: https://github.com/k8sdb/postgres-init-scripts.git
      scriptPath: postgres-init-scripts/run.sh
  resources: {}
  storage:
    accessModes:
    - ReadWriteOnce
    class: standard
    resources:
      requests:
        storage: 50Mi
  version: "9.5"
status:
  creationTime: 2017-07-17T22:31:34Z
  phase: Running
```

```sh
$ kubedb describe pg -n demo p1
Name:		p1
Namespace:	demo
StartTimestamp:	Mon, 17 Jul 2017 15:31:34 -0700
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

Service:	
  Name:		p1
  Type:		ClusterIP
  IP:		10.0.0.161
  Port:		db	5432/TCP

Database Secret:
  Name:	p1-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	35 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                Type       Reason               Message
  ---------   --------   -----     ----                --------   ------               -------
  2m          2m         1         Postgres operator   Normal     SuccessfulValidate   Successfully validate Postgres
  2m          2m         1         Postgres operator   Normal     SuccessfulCreate     Successfully created Postgres
  2m          2m         1         Postgres operator   Normal     SuccessfulCreate     Successfully created StatefulSet
  3m          3m         1         Postgres operator   Normal     SuccessfulValidate   Successfully validate Postgres
  3m          3m         1         Postgres operator   Normal     Creating             Creating Kubernetes objects
```








In this tutorial, we are going to backup the `/source/data` folder of a `busybox` pod into a local backend. First deploy the following `busybox` Deployment in your cluster. Here we are using a git repository as source volume for demonstration purpose.

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: stash-demo
  name: stash-demo
  namespace: default
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: stash-demo
      name: busybox
    spec:
      containers:
      - command:
        - sleep
        - "3600"
        image: busybox
        imagePullPolicy: IfNotPresent
        name: busybox
        volumeMounts:
        - mountPath: /source/data
          name: source-data
      restartPolicy: Always
      volumes:
      - gitRepo:
          repository: https://github.com/appscode/stash-data.git
        name: source-data
```




# Tutorial

### Initialize unified operator

**F**irst of all, we need an unified [operator](https://github.com/k8sdb/operator) to handle supported TPR.

We can deploy this operator using `kubedb` CLI.

```bash
$ kubedb init

Successfully created operator deployment.
Successfully created operator service.
```

This will deploy a controller that can handle our TPR objects.

For more, see [init](operation/init.md) operation.

When operator is ready, we can create database.

### Create database

**W**e are supporting two kinds of databases. Lets see how to create them.

* Create [Postgres](database/postgres/create.md) database
* Create [Elasticsearch](database/elastic/create.md) database

### Describe database

**W**e can describe our database using `kubedb` CLI. 

```bash
$ kubedb describe <pg|es> database-demo
```
The `kubedb describe` command provides following basic information of TPR object.
* StatefulSet
* Storage (Persistent Volume) (If available)
* Service
* Secret (If available)
* Snapshots (If any)
* Monitoring system (If available)

For more, see [describe](task/describe.md) operation.

### Take instant backup

**W**e can take backup of a running database using Snapshot TPR object.

This Snapshot object will contain all necessary information to take backup and upload backup data to cloud.

For details explanation, see [here](database/backup.md) .


### Lock database

**T**o prevent a database from deleting, we can set `doNotPause: true`.

```bash
$ kubedb edit <pg|es> database-demo

# Add following in spec
#  doNotPause: true
```

To undo this, set `doNotPause: false`.

### Pause database

**P**ausing a database means operator will only delete _StatefulSet_ & _Service_. Database can be resumed later.

Lets pause a database by deleting it.

```bash
$ kubedb delete <pg|es> database-demo
```

For deleting this Postgres object, operator will create a DormantDatabase object `database-demo`.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  labels:
    kubedb.com/kind: <Postgres|Elastic>
  name: database-demo
spec:
  origin:
    metadata:
      name: database-demo
      namespace: default
    spec:
      -----
      -----
status:
  creationTime: 2017-06-06T08:11:10Z
  pausingTime: 2017-06-06T08:11:50Z
  phase: Paused
```

This DormantDatabase object will contain original database information in `spec.origin`.

**W**hen this DormantDatabase object will be created, operator will delete following workloads:

* StatefulSet
* Service

And following kubernetes objects will be intact:

* Secret
* GoverningService
* PersistentVolumeClaim

We can get back our database again by resuming it.

### Resume database

**T**o resume our database `database-demo` from DormantDatabase object,
we need to edit this DormantDatabase object to set `resume: true`.

```bash
$ kubedb edit drmn database-demo

# Add following in spec
#  resume: true

dormantdatabase "database-demo" edited
```

Operator will detect this modification and will create database TPR object `database-demo`.

**T**his new database TPR object will use same _Secret_ and StatefulSet will be created with existing PersistentVolumeClaim.

### Wipeout database

**T**o wipeout a database, we need to set `wipeOut: true` in DormantDatabase object

```bash
$ kubedb edit drmn database-demo

# Add following in spec
#  wipeOut: true

dormantdatabase "database-demo" edited
```

This will delete following kubernetes objects:

* Secret
* PersistentVolumeClaim

And also delete all Snapshot objects.

Once database is Wipedout, there is no way to resume it.
