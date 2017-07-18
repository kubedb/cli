# Using PostgreSQL
This tutorial will show you how to use KubeDB to run a PostgreSQL database. 

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). 

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. This tutorial will also use a PGAdmin to connect and test PostgreSQL database, once it is running. Run the following command to prepare your cluster for this tutorial:

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
KubeDB implements a `Postgres` TPR to define the specification of a PostgreSQL database. Below is the `Postgres` object created in this tutorial.

```yaml
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

$ kubedb create -f ./docs/examples/tutorial/postgres/demo-1.yaml 
validating "./docs/examples/tutorial/postgres/demo-1.yaml"
postgres "p1" created
```

Here,
 - `spec.version` is the version of PostgreSQL database. In this tutorial, a PostgreSQL 9.5 database is going to be created.

 - `spec.doNotPause` tells KubeDB operator that if this tpr is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.

 - `spec.storage` specifies the StorageClass of PVC dynamically allocated to store data for this database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests. If no storage spec is given, an `emptyDir` is used.

 - `spec.init.scriptSource` specifies a bash script used to initialize the database after it is created. In this tutorial, `run.sh` script from the git repository `https://github.com/k8sdb/postgres-init-scripts.git` is used to create a `dashboard` table in `data` schema.

KubeDB operator watches for `Postgres` objects using Kubernetes api. When a `Postgres` object is created, KubeDB operator will create a new StatefulSet and a ClusterIP Service with the matching tpr name. KubeDB operator will also create a governing service for StatefulSets with the name `kubedb`, if one is not already present. If [RBAC is enabled](/docs/rbac.md), a ClusterRole, ServiceAccount and ClusterRoleBinding with the matching tpr name will be created and used as the service account name for the corresponding StatefulSet.

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


$ kubectl get statefulset -n demo
NAME      DESIRED   CURRENT   AGE
p1        1         1         1m

$ kubectl get pvc -n demo
NAME        STATUS    VOLUME                                     CAPACITY   ACCESSMODES   STORAGECLASS   AGE
data-p1-0   Bound     pvc-e90b87d4-6b5a-11e7-b9ca-080027f73ab7   50Mi       RWO           standard       1m

$ kubectl get pv -n demo
NAME                                       CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS    CLAIM            STORAGECLASS   REASON    AGE
pvc-e90b87d4-6b5a-11e7-b9ca-080027f73ab7   50Mi       RWO           Delete          Bound     demo/data-p1-0   standard                 1m

$ kubectl get service -n demo
NAME      CLUSTER-IP   EXTERNAL-IP   PORT(S)        AGE
kubedb    None         <none>                       3m
p1        10.0.0.143   <none>        5432/TCP       3m
pgadmin   10.0.0.120   <pending>     80:30576/TCP   6m
```


KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified tpr:

```yaml
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


Please note that KubeDB operator has created a new Secret called `p1-admin-auth` (format: {tpr-name}-admin-auth) for storing the password for `postgres` superuser. This secret contains a `.admin` key with a ini formatted key-value pairs. If you want to use an existing secret please specify that when creating the tpr using `spec.databaseSecret.secretName`.

Now, you can connect to this database from the PGAdmin dasboard using the database pod IP and `postgres` user password. 

```sh
$ kubectl get pods p1-0 -n demo -o yaml | grep IP
  hostIP: 192.168.99.100
  podIP: 172.17.0.6

$ kubectl get secrets -n demo p1-admin-auth -o jsonpath={'.data.\.admin'} | base64 -d
POSTGRES_PASSWORD=R9keKKRTqSJUPtNC
```

![Using p1 from PGAdmin4](/docs/images/tutorial/postgres/p1-pgadmin.gif)


## Taking Snapshots
Now, you can easily take a snapshot of this database by creating a `Snapshot` tpr. When a `Snapshot` tpr is created, KubeDB operator will launch a Job that runs the `pg_dump` command and uploads the output sql file to various cloud providers S3, GCS, Azure, OpenStack Swift and locally mounted volumes using [osm](https://github.com/appscode/osm).

In this tutorial, snapshots will be stored in a Google Cloud Storage (GCS) bucket. To do so, a secret is needed that has the following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```sh
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic pg-snap-secret -n demo \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "pg-snap-secret" created
```

```yaml
$ kubectl get secret pg-snap-secret -o yaml

apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2017-07-17T18:06:51Z
  name: pg-snap-secret
  namespace: demo
  resourceVersion: "5461"
  selfLink: /api/v1/namespaces/demo/secrets/pg-snap-secret
  uid: a6983b00-5c02-11e7-bb52-08002711f4aa
type: Opaque
```


To lean how to configure other storage destinations for Snapshots, please visit [here](/docs/snapshot.md). Now, create the Snapshot tpr.

```
$ kubedb create -f ./docs/examples/tutorial/postgres/demo-2.yaml 
validating "./docs/examples/tutorial/postgres/demo-2.yaml"
snapshot "p1-xyz" created

$ kubedb get snap -n demo
NAME      DATABASE   STATUS    AGE
p1-xyz    pg/p1      Running   22s
```

```yaml
$ kubedb get snap -n demo p1-xyz -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  creationTimestamp: 2017-07-18T02:18:00Z
  labels:
    kubedb.com/kind: Postgres
    kubedb.com/name: p1
  name: p1-xyz
  namespace: demo
  resourceVersion: "2973"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/snapshots/p1-xyz
  uid: 5269701f-6b5f-11e7-b9ca-080027f73ab7
spec:
  databaseName: p1
  gcs:
    bucket: restic
  resources: {}
  storageSecretName: snap-secret
status:
  completionTime: 2017-07-18T02:19:11Z
  phase: Succeeded
  startTime: 2017-07-18T02:18:00Z
```

Here,

- `metadata.labels` should include the type of database `kubedb.com/kind: Postgres` whose snapshot will be taken.

- `spec.databaseName` points to the databse whose snapshot is taken.

- `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.

- `spec.gcs.bucket` points to the bucket name used to store the snapshot data.


You can also run the `kubedb describe` command to see the recent snapshots taken for a database.

```sh
$ kubedb describe pg -n demo p1
Name:		p1
Namespace:	demo
StartTimestamp:	Mon, 17 Jul 2017 18:46:24 -0700
Status:		Running
Volume:
  StorageClass:	standard
  Capacity:	50Mi
  Access Modes:	RWO

Service:	
  Name:		p1
  Type:		ClusterIP
  IP:		10.0.0.143
  Port:		db	5432/TCP

Database Secret:
  Name:	p1-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	35 bytes

Snapshots:
  Name     Bucket      StartTime                         CompletionTime                    Phase
  ----     ------      ---------                         --------------                    -----
  p1-xyz   gs:restic   Mon, 17 Jul 2017 19:18:00 -0700   Mon, 17 Jul 2017 19:19:11 -0700   Succeeded

Events:
  FirstSeen   LastSeen   Count     From                  Type       Reason               Message
  ---------   --------   -----     ----                  --------   ------               -------
  1m          1m         1         Snapshot Controller   Normal     SuccessfulSnapshot   Successfully completed snapshot
  2m          2m         1         Snapshot Controller   Normal     Starting             Backup running
  33m         33m        1         Postgres operator     Normal     SuccessfulValidate   Successfully validate Postgres
  33m         33m        1         Postgres operator     Normal     SuccessfulCreate     Successfully created StatefulSet
  33m         33m        1         Postgres operator     Normal     SuccessfulCreate     Successfully created Postgres
  34m         34m        1         Postgres operator     Normal     Creating             Creating Kubernetes objects
  34m         34m        1         Postgres operator     Normal     SuccessfulValidate   Successfully validate Postgres
```

Once the snapshot Job is complete, you should see the output of the `pg_dump` command stored in the GCS buckeet.

![snapshot-console](/docs/images/tutorial/postgres/p1-xyz-snapshot.png)



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
