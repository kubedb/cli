> New to KubeDB? Please start [here](/docs/tutorials/README.md).

# Running Memcached
This tutorial will show you how to use KubeDB to run a Memcached database.

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). 

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f ./docs/examples/memcached/demo-0.yaml
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    1h
demo          Active    21s
kube-public   Active    1h
kube-system   Active    1h
```

## Create a Memcached database
KubeDB implements a `Memcached` CRD to define the specification of a Memcached database. Below is the `Memcached` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  name: mc1
  namespace: demo
spec:
  replicas: 3
  version: 1.5.3-alpine
  doNotPause: true
  resources:
    requests:
      memory: 50Mi

$ kubedb create -f ./docs/examples/memcached/demo-1.yaml
validating "./docs/examples/memcached/demo-1.yaml"
memcached "mc1" created
```

Here,
 - `spec.replicas` is an optional field that specifies the number of desired Instances/Replicas of Memcached database. It defaults to 1.

 - `spec.version` is the version of Memcached database. In this
   tutorial, 3 instances of Memcached 1.5.3-alpine database is going to be created.

 - `spec.doNotPause` tells KubeDB operator that if this CRD is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.

 - `spec.resource` is an optional field that specifies how much CPU and memory (RAM) each Container needs. To learn details about Managing Compute Resources for Containers, please visit [here](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/).

KubeDB operator watches for `Memcached` objects using Kubernetes api. When a `Memcached` object is created, KubeDB operator will create a new `Deployment` and a ClusterIP Service with the matching crd name. If [RBAC is enabled](/docs/tutorials/rbac.md), a ClusterRole, ServiceAccount and ClusterRoleBinding with the matching crd name will be created and used as the service account name for the corresponding Deployment.

```console
$ kubedb describe mc -n demo mc1
Name:		mc1
Namespace:	demo
StartTimestamp:	Thu, 07 Dec 2017 16:08:10 +0600
Status:		Running

Deployment:		
  Name:			mc1
  Replicas:		3 current / 3 desired
  CreationTimestamp:	Thu, 07 Dec 2017 16:08:14 +0600
  Pods Status:		3 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:	
  Name:		mc1
  Type:		ClusterIP
  IP:		10.106.213.85
  Port:		db	11211/TCP

Events:
  FirstSeen   LastSeen   Count     From                 Type       Reason               Message
  ---------   --------   -----     ----                 --------   ------               -------
  2m          2m         1         Memcached operator   Normal     SuccessfulCreate     Successfully created Deployment
  2m          2m         1         Memcached operator   Normal     SuccessfulCreate     Successfully created Memcached
  3m          3m         1         Memcached operator   Normal     SuccessfulValidate   Successfully validate Memcached
  3m          3m         1         Memcached operator   Normal     Creating             Creating Kubernetes objects


$ kubectl get deployment -n demo
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
mc1       3         3         3            3           4m


$ kubectl get service -n demo
NAME      TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)     AGE
mc1       ClusterIP   10.106.213.85   <none>        11211/TCP   12m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified crd:

```yaml
$ kubedb get mc -n demo mc1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-07T10:08:10Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  name: mc1
  namespace: demo
  resourceVersion: "8816"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/memcacheds/mc1
  uid: 876225d5-db36-11e7-b052-0800275191b8
spec:
  doNotPause: true
  replicas: 3
  resources:
    requests:
      memory: 50Mi
  version: 1.5.3-alpine
status:
  creationTime: 2017-12-07T10:08:10Z
  phase: Running
```

Now, you can connect to this Memcached cluster from inside the cluster..

```console
$ kubectl get pods -n demo
NAME                   READY     STATUS    RESTARTS   AGE
mc1-66bc4d558d-6xkcz   1/1       Running   0          49m
mc1-66bc4d558d-h48dw   1/1       Running   0          49m
mc1-66bc4d558d-ns9zd   1/1       Running   0          49m

$ kubectl get pods mc1-66bc4d558d-6xkcz -n demo -o yaml | grep IP
  hostIP: 192.168.99.100
  podIP: 172.17.0.6
  
# Exec into kubedb operator pod
$ kubectl exec -it $(kubectl get pods --all-namespaces -l app=kubedb -o jsonpath='{.items[0].metadata.name}') -n kube-system sh

~ $ ps aux
PID   USER     TIME   COMMAND
    1 nobody     0:01 /operator run --address=:8080 --rbac=false --v=3
   27 nobody     0:00 sh
   33 nobody     0:00 ps aux

# Connect Memcached cluster through telnet
~ $ telnet 172.17.0.6 11211

# Save a data
> set my_key 0 2592000 1
2          
STORED

# Meaning:
# 0       => no flags
# 2592000 => TTL (Time-To-Live) in [s]
# 1       => size in byte
# 2       => value

# View the data
> get my_key
VALUE my_key 0 1
2
END
```


## Pause Database

Since the Memcached crd created in this crd has `spec.doNotPause` set to true, if you delete the crd, KubeDB operator will recreate the crd and essentially nullify the delete operation. You can see this below:

```console
$ kubedb delete mc mc1 -n demo
error: Memcached "mc1" can't be paused. To continue delete, unset spec.doNotPause and retry.
```

Now, run `kubedb edit mc mc1 -n demo` to set `spec.doNotPause` to false or remove this field (which default to false). Then if you delete the Memcached crd, KubeDB operator will delete the Deployment and its pods. In KubeDB parlance, we say that `mc1` Memcached database has entered into dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase crd.

```yaml
$ kubedb delete mc -n demo mc1
memcached "mc1" deleted

$ kubedb get drmn -n demo mc1
NAME      STATUS    AGE
mc1        Pausing   20s

$ kubedb get drmn -n demo mc1
NAME      STATUS    AGE
mc1       Paused    53s

$ kubedb get drmn -n demo mc1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-07T11:19:36Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: Memcached
  name: mc1
  namespace: demo
  resourceVersion: "13639"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mc1
  uid: 825ce6f1-db40-11e7-b052-0800275191b8
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: mc1
      namespace: demo
    spec:
      memcached:
        replicas: 3
        resources:
          requests:
            memory: 50Mi
        version: 1.5.3-alpine
status:
  creationTime: 2017-12-07T11:19:36Z
  pausingTime: 2017-12-07T11:20:07Z
  phase: Paused
```

Here,
 - `spec.origin` is the spec of the original spec of the original Memcached crd.

 - `status.phase` points to the current database state `Paused`.


## Resume Dormant Database

To resume the database from the dormant state, set `spec.resume` to `true` in the DormantDatabase crd.

```yaml
$ kubedb edit drmn -n demo mc1
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-07T11:22:39Z
  deletionGracePeriodSeconds: null
  deletionTimestamp: null
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: Memcached
  name: mc1
  namespace: demo
  resourceVersion: "13875"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mc1
  uid: eede27dd-db40-11e7-b052-0800275191b8
spec:
  resume: true
  origin:
    metadata:
      creationTimestamp: null
      name: mc1
      namespace: demo
    spec:
      memcached:
        replicas: 3
        resources:
          requests:
            memory: 50Mi
        version: 1.5.3-alpine
status:
  creationTime: 2017-12-07T11:22:39Z
  phase: Pausing
```

KubeDB operator will notice that `spec.resume` is set to true. KubeDB operator will delete the DormantDatabase crd and create a new Memcached crd using the original spec. This will in turn start a new Deployment.This way the memcached database is resumed.


## Delete Dormant Database
To delete a dormant database, it needs to be wiped out first. You can also wipe out a DormantDatabase by setting `spec.wipeOut` to true. There is no way to resume a wiped out database. So, be sure before you wipe out a database.
```yaml
$ kubedb edit drmn -n demo mc1
# set spec.wipeOut: true

$ kubedb get drmn -n demo mc1 -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2017-12-07T11:32:21Z
  generation: 0
  initializers: null
  labels:
    kubedb.com/kind: Memcached
  name: mc1
  namespace: demo
  resourceVersion: "14790"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/mc1
  uid: 49d89efe-db42-11e7-b052-0800275191b8
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: mc1
      namespace: demo
    spec:
      memcached:
        replicas: 3
        resources:
          requests:
            memory: 50Mi
        version: 1.5.3-alpine
  wipeOut: true
status:
  creationTime: 2017-12-07T11:32:21Z
  pausingTime: 2017-12-07T11:33:11Z
  phase: WipedOut
  wipeOutTime: 2017-12-07T11:34:08Z

$ kubedb get drmn -n demo
NAME      STATUS     AGE
mc1       WipedOut   2m

```

 You still have a record that there used to be an Memcached database `mc1` in the form of a DormantDatabase database `mc1`. Since you have already wiped out the database, you can delete the DormantDatabase crd.

```console
$ kubedb delete drmn mc1 -n demo
dormantdatabase "mc1" deleted

$ kubedb get drmn -n demo
No resources found.
```

## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall KubeDB operator, please follow the steps [here](/docs/uninstall.md).


## Next Steps
- Learn about the details of Memcached crd [here](/docs/concepts/memcached.md).
- Thinking about monitoring your database? KubeDB works [out-of-the-box with Prometheus](/docs/tutorials/monitoring.md).
- Learn how to use KubeDB in a [RBAC](/docs/tutorials/rbac.md) enabled cluster.
- Wondering what features are coming next? Please visit [here](/ROADMAP.md). 
- Want to hack on KubeDB? Check our [contribution guidelines](/CONTRIBUTING.md).
