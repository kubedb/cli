---
title: Memcached Quickstart
menu:
  docs_0.8.0-beta.2:
    identifier: mc-quickstart-quickstart
    name: Overview
    parent: mc-quickstart-memcached
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Memcached QuickStart

This tutorial will show you how to use KubeDB to run a Memcached database.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/memcached/memcached-lifecycle.png" width="600" height="373">
</p>

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/memcached/demo-0.yaml
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    45m
demo          Active    10s
kube-public   Active    45m
kube-system   Active    45m
```

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Create a Memcached database

KubeDB implements a `Memcached` CRD to define the specification of a Memcached database. Below is the `Memcached` object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  name: memcd-quickstart
  namespace: demo
spec:
  replicas: 3
  version: 1.5.4
  doNotPause: true
  resources:
    requests:
      memory: 64Mi
      cpu: 250m
    limits:
      memory: 128Mi
      cpu: 500m
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/memcached/quickstart/demo-1.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/memcached/quickstart/demo-1.yaml"
memcached "memcd-quickstart" created
```

Here,

- `spec.replicas` is an optional field that specifies the number of desired Instances/Replicas of Memcached database. It defaults to 1.
- `spec.version` is the version of Memcached database. In this tutorial, a Memcached 1.5.4 database is going to be created.
- `spec.doNotPause` tells KubeDB operator that if this object is deleted, it should be automatically reverted. This should be set to true for production databases to avoid accidental deletion.
- `spec.resource` is an optional field that specifies how much CPU and memory (RAM) each Container needs. To learn details about Managing Compute Resources for Containers, please visit [here](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/).

KubeDB operator watches for `Memcached` objects using Kubernetes api. When a `Memcached` object is created, KubeDB operator will create a new Deployment and a ClusterIP Service with the matching Memcached object name. No Memcached specific RBAC permission is required in [RBAC enabled clusters](/docs/setup/install.md#using-yaml).

```console
$ kubedb get mc -n demo
NAME               STATUS    AGE
memcd-quickstart   Running   1m

$ kubedb describe mc -n demo memcd-quickstart
Name:		memcd-quickstart
Namespace:	demo
StartTimestamp:	Tue, 13 Feb 2018 10:53:47 +0600
Status:		Running

Deployment:
  Name:			memcd-quickstart
  Replicas:		3 current / 3 desired
  CreationTimestamp:	Tue, 13 Feb 2018 10:53:48 +0600
  Pods Status:		3 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		memcd-quickstart
  Type:		ClusterIP
  IP:		10.103.23.178
  Port:		db	11211/TCP

Events:
  FirstSeen   LastSeen   Count     From                 Type       Reason       Message
  ---------   --------   -----     ----                 --------   ------       -------
  1m          1m         1         Memcached operator   Normal     Successful   Successfully created Deployment
  1m          1m         1         Memcached operator   Normal     Successful   Successfully created Memcached
  2m          2m         1         Memcached operator   Normal     Successful   Successfully created Service

$ kubectl get deployment -n demo
NAME               DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
memcd-quickstart   3         3         3            3           2m

$ kubectl get service -n demo
NAME               TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)     AGE
memcd-quickstart   ClusterIP   10.103.23.178   <none>        11211/TCP   3m
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created. Run the following command to see the modified Memcached object:

```yaml
$ kubedb get mc -n demo memcd-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-13T04:53:47Z
  finalizers:
  - kubedb.com
  generation: 0
  name: memcd-quickstart
  namespace: demo
  resourceVersion: "1321"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/memcacheds/memcd-quickstart
  uid: e01b4c39-1079-11e8-801e-080027e82bd4
spec:
  doNotPause: true
  replicas: 3
  resources:
    limits:
      cpu: 500m
      memory: 128Mi
    requests:
      cpu: 250m
      memory: 64Mi
  version: 1.5.4
status:
  creationTime: 2018-02-13T04:53:47Z
  phase: Running

```

Now, you can connect to this Memcached cluster using `telnet`.
Here, we will connect to Memcached database from local-machine through port-forwarding.

```console
$ kubectl get pods -n demo
NAME                                READY     STATUS    RESTARTS   AGE
memcd-quickstart-667cd68854-gs69q   1/1       Running   0          4m
memcd-quickstart-667cd68854-hpkbb   1/1       Running   0          4m
memcd-quickstart-667cd68854-jlmwh   1/1       Running   0          4m

// We will connect to `memcd-quickstart-667cd68854-gs69q` pod from local-machine using port-frowarding.
$ kubectl port-forward -n demo memcd-quickstart-667cd68854-gs69q 11211
Forwarding from 127.0.0.1:11211 -> 11211

# Connect Memcached cluster from localmachine through telnet.
~ $ telnet 127.0.0.1 11211
Trying 127.0.0.1...
Connected to 127.0.0.1.

# Save data Command:
set my_key 0 2592000 1
2
# Output:
STORED

# Meaning:
# 0       => no flags
# 2592000 => TTL (Time-To-Live) in [s]
# 1       => size in byte
# 2       => value

# View data command
get my_key
# Output
VALUE my_key 0 1
2
END

# Exit
quit
```

## Pause Database

KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `doNotPause` feature. If admission webhook is enabled, It prevents user from deleting the database as long as the `spec.doNotPause` is set to true. Since the Memcached object created in this tutorial has `spec.doNotPause` set to true, if you delete the Memcached object, KubeDB operator will nullify the delete operation. You can see this below:

```console
$ kubedb delete mc memcd-quickstart -n demo
error: Memcached "memcd-quickstart" can't be paused. To continue delete, unset spec.doNotPause and retry.
```

Now, run `kubedb edit mc memcd-quickstart -n demo` to set `spec.doNotPause` to false or remove this field (which default to false). Then if you delete the Memcached object, KubeDB operator will delete the Deployment and its pods. In KubeDB parlance, we say that `memcd-quickstart` Memcached database has entered into dormant state. This is represented by KubeDB operator by creating a matching DormantDatabase object.

```console
$ kubedb delete mc memcd-quickstart -n demo
memcached "memcd-quickstart" deleted

$ kubedb get drmn -n demo memcd-quickstart
NAME               STATUS    AGE
memcd-quickstart   Pausing   6s

$ kubedb get drmn -n demo memcd-quickstart
NAME               STATUS    AGE
memcd-quickstart   Paused    9s
```

```yaml
$ kubedb get drmn -n demo memcd-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-13T05:34:05Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: Memcached
  name: memcd-quickstart
  namespace: demo
  resourceVersion: "2854"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/memcd-quickstart
  uid: 814d99f6-107f-11e8-801e-080027e82bd4
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: memcd-quickstart
      namespace: demo
    spec:
      memcached:
        replicas: 3
        resources:
          limits:
            cpu: 500m
            memory: 128Mi
          requests:
            cpu: 250m
            memory: 64Mi
        version: 1.5.4
status:
  creationTime: 2018-02-13T05:34:05Z
  pausingTime: 2018-02-13T05:34:14Z
  phase: Paused
```

Here,

- `spec.origin` is the spec of the original spec of the original Memcached object.
- `status.phase` points to the current database state `Paused`.

## Resume Dormant Database

To resume the database from the dormant state, set `spec.resume` to `true` in the DormantDatabase object.

```yaml
$ kubedb edit drmn -n demo memcd-quickstart
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  name: memcd-quickstart
  namespace: demo
  ...
spec:
  resume: true
  ...
status:
  phase: Paused
  ...
```

KubeDB operator will notice that `spec.resume` is set to true. KubeDB operator will delete the DormantDatabase object and create a new Memcached object using the original spec. It will start fresh as there is no persistent volume for Memcached.

Please note that the dormant database can also be resumed by creating same `Memcached` database by using same Specs. In this tutorial, the dormant database can be resumed by creating `Memcached` database using demo-1.yaml file. The below command resumes the dormant database `memcd-quickstart` that was created before.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/memcached/quickstart/demo-1.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/memcached/quickstart/demo-1.yaml"
memcached "memcd-quickstart" created
```

## Wipeout Dormant Database

You can also wipe out a DormantDatabase by setting `spec.wipeOut` to true. There is no way to resume a wiped out database. So, be sure before you wipe out a database.

Create dormant database again and set `spec.wipeOut` to true:

```yaml
$ kubedb delete mc memcd-quickstart -n demo
memcached "memcd-quickstart" deleted

$ kubedb edit drmn -n demo memcd-quickstart
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  name: memcd-quickstart
  namespace: demo
  ...
spec:
  wipeOut: true
  ...
status:
  phase: Paused
  ...

$ kubedb get drmn -n demo memcd-quickstart -o yaml
apiVersion: kubedb.com/v1alpha1
kind: DormantDatabase
metadata:
  clusterName: ""
  creationTimestamp: 2018-02-13T05:59:44Z
  finalizers:
  - kubedb.com
  generation: 0
  labels:
    kubedb.com/kind: Memcached
  name: memcd-quickstart
  namespace: demo
  resourceVersion: "4093"
  selfLink: /apis/kubedb.com/v1alpha1/namespaces/demo/dormantdatabases/memcd-quickstart
  uid: 16e9a846-1083-11e8-801e-080027e82bd4
spec:
  origin:
    metadata:
      creationTimestamp: null
      name: memcd-quickstart
      namespace: demo
    spec:
      memcached:
        replicas: 3
        resources:
          limits:
            cpu: 500m
            memory: 128Mi
          requests:
            cpu: 250m
            memory: 64Mi
        version: 1.5.4
  wipeOut: true
status:
  creationTime: 2018-02-13T05:59:44Z
  pausingTime: 2018-02-13T05:59:57Z
  phase: WipedOut
  wipeOutTime: 2018-02-13T06:04:41Z

$ kubedb get drmn -n demo
NAME               STATUS     AGE
memcd-quickstart   WipedOut   5m
```

## Delete Dormant Database

You still have a record that there used to be a Memcached database `memcd-quickstart` in the form of a DormantDatabase database `memcd-quickstart`. Since you have already wiped out the database, you can delete the DormantDatabase object.

```console
$ kubedb delete drmn memcd-quickstart -n demo
dormantdatabase "memcd-quickstart" deleted
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete mc memcd-quickstart -n demo --force
$ kubedb delete drmn memcd-quickstart -n demo --force

# or
# $ kubedb delete mc,drmn -n demo --all --force

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Monitor your Memcached database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/memcached/monitoring/using-coreos-prometheus-operator.md).
- Monitor your Memcached database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/memcached/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/memcached/private-registry/using-private-registry.md) to deploy Memcached with KubeDB.
- Detail concepts of [Memcached object](/docs/concepts/databases/memcached.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
