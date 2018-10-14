---
title: MySQL
menu:
  docs_0.9.0-rc.0:
    identifier: mysql-db
    name: MySQL
    parent: databases
    weight: 25
menu_name: docs_0.9.0-rc.0
section_menu_id: concepts
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# MySQL

## What is MySQL

`MySQL` is a Kubernetes `Custom Resource Definitions` (CRD). It provides declarative configuration for [MySQL](https://www.mysql.com/) in a Kubernetes native way. You only need to describe the desired database configuration in a MySQL object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## MySQL Spec

As with all other Kubernetes objects, a MySQL needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example MySQL object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: m1
  namespace: demo
spec:
  version: "8.0-v1"
  databaseSecret:
    secretName: m1-auth
  storageType: "Durable"
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    scriptSource:
      configMap:
        name: mg-init-script
  backupSchedule:
    cronExpression: "@every 6h"
    storageSecretName: ms-snap-secret
    gcs:
      bucket: kubedb
      prefix: demo
  monitor:
    agent: prometheus.io/coreos-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
  configSource:
    configMap:
      name: my-custom-config
  podTemplate:
    annotation:
      passMe: ToDatabasePod
    controller:
      annotation:
        passMe: ToStatefulSet
    spec:
      schedulerName: my-scheduler
      nodeSelector:
        disktype: ssd
      imagePullSecrets:
      - name: myregistrykey
      args:
      - --character-set-server=utf8mb4
      env:
      - name: MYSQL_DATABASE
        value: myDB
      resources:
        requests:
          memory: "64Mi"
          cpu: "250m"
        limits:
          memory: "128Mi"
          cpu: "500m"
  serviceTemplate:
    annotation:
      passMe: ToService
    spec:
      type: NodePort
      ports:
      - name:  http
        port:  9200
        targetPort: http
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
```

### spec.version

`spec.version` is a required field specifying the name of the [MySQLVersion](/docs/concepts/catalog/mysql.md) crd where the docker images are specified. Currently, when you install KubeDB, it creates the following `MySQLVersion` crd,

- `8.0-v1`, `8.0`, `8-v1`, `8`
- `5.7-v1`, `5.7`, `5-v1`, `5`

### spec.databaseSecret

`spec.databaseSecret` is an optional field that points to a Secret used to hold credentials for `mysql` root user. If not set, KubeDB operator creates a new Secret `{mysql-object-name}-auth` for storing the password for `mysql` root user for each MySQL object. If you want to use an existing secret please specify that when creating the MySQL object using `spec.databaseSecret.secretName`.

This secret contains a `user` key and a `password` key which contains the `username` and `password` respectively for `mysql` root user. Here, the value of `user` key is fixed to be `root`.

Example:

```console
$ kubectl create secret generic m1-auth -n demo \
--from-literal=user=root \
--from-literal=password=6q8u_2jMOW-OOZXk
secret "m1-auth" created
```

```yaml
apiVersion: v1
data:
  password: NnE4dV8yak1PVy1PT1pYaw==
  user: cm9vdA==
kind: Secret
metadata:
  ...
  name: m1-auth
  namespace: demo
  ...
type: Opaque
```

### spec.storageType

`spec.storageType` is an optional field that specifies the type of storage to use for database. It can be either `Durable` or `Ephemeral`. The default value of this field is `Durable`. If `Ephemeral` is used then KubeDB will create MySQL database using [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) volume. In this case, you don't have to specify `spec.storage` field.

### spec.storage

Since 0.9.0-rc.0, If you set `spec.storageType:` to `Durable`, then  `spec.storage` is a required field that specifies the StorageClass of PVCs dynamically allocated to store data for the database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests.

- `spec.storage.storageClassName` is the name of the StorageClass used to provision PVCs. PVCs don’t necessarily have to request a class. A PVC with its storageClassName set equal to "" is always interpreted to be requesting a PV with no class, so it can only be bound to PVs with no class (no annotation or one set equal to ""). A PVC with no storageClassName is not quite the same and is treated differently by the cluster depending on whether the DefaultStorageClass admission plugin is turned on.
- `spec.storage.accessModes` uses the same conventions as Kubernetes PVCs when requesting storage with specific access modes.
- `spec.storage.resources` can be used to request specific quantities of storage. This follows the same resource model used by PVCs.

To learn how to configure `spec.storage`, please visit the links below:

- https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims

### spec.init

`spec.init` is an optional section that can be used to initialize a newly created MySQL database. MySQL databases can be initialized in one of two ways:

  1. Initialize from Script
  2. Initialize from Snapshot

#### Initialize via Script

To initialize a MySQL database using a script (shell script, sql script etc.), set the `spec.init.scriptSource` section when creating a MySQL object. It will execute files alphabetically with extensions `.sh` , `.sql`  and `.sql.gz` that are found in the repository. The scripts inside child folders will be skipped. ScriptSource must have following information:

- [VolumeSource](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes): Where your script is loaded from.

Below is an example showing how a script from a configMap can be used to initialize a MySQL database.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: m1
spec:
  version: 8.0-v1
  init:
    scriptSource:
      configMap:
        name: mysql-init-script
```

In the above example, KubeDB operator will launch a Job to execute all js script of `mysql-init-script` in alphabetical order once StatefulSet pods are running. For more details tutorial on how to initialize from script, please visit [here](/docs/guides/mysql/initialization/script-source.md).

#### Initialize from Snapshots

To initialize from prior snapshots, set the `spec.init.snapshotSource` section when creating a MySQL object. In this case, SnapshotSource must have following information:

- `name:` Name of the Snapshot
- `namespace:` Namespace of the Snapshot

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: m1
spec:
  version: 8.0-v1
  init:
    snapshotSource:
      name: "snapshot-xyz"
      namespace: "demo"
```

In the above example, MySQL database will be initialized from Snapshot `snapshot-xyz` in `demo` namespace. Here, KubeDB operator will launch a Job to initialize MySQL once StatefulSet pods are running.

When initializing from Snapshot, root user credentials must have to match with the previous one. For example, let's say, Snapshot `snapshot-xyz` is for MySQL `mysql-old`. In this case, new MySQL `mysql-db` should use the same credentials for root user of `mysql-old`. Otherwise, the restoration process will fail.

For more details tutorial on how to initialize from snapshot, please visit [here](/docs/guides/mysql/initialization/snapshot-source.md).

### spec.backupSchedule

KubeDB supports taking periodic snapshots for MySQL database. This is an optional section in `.spec`. When `spec.backupSchedule` section is added, KubeDB operator immediately takes a backup to validate this information. After that, at each tick kubeDB operator creates a [Snapshot](/docs/concepts/snapshot.md) object. This triggers operator to create a Job to take backup. If used, set the various sub-fields accordingly.

- `spec.backupSchedule.cronExpression` is a required [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). This specifies the schedule for backup operations.
- `spec.backupSchedule.{storage}` is a required field that is used as the destination for storing snapshot data. KubeDB supports cloud storage providers like S3, GCS, Azure and OpenStack Swift. It also supports any locally mounted Kubernetes volumes, like NFS, Ceph, etc. Only one backend can be used at a time. To learn how to configure this, please visit [here](/docs/concepts/snapshot.md).

You can also specify a template for pod of backup job through `spec.backupSchedule.podTemplate`. KubeDB will use the information you have provided in `podTemplate` to create the backup job. KubeDB accept following fields to set in `spec.backupSchedule.podTemplate`:

- metadata:
  - annotations (pod's annotation)
- controller:
  - annotations (job's annotation)
- spec:
  - args
  - env
  - resources
  - imagePullSecrets
  - initContainers
  - nodeSelector
  - affinity
  - schedulerName
  - tolerations
  - priorityClassName
  - priority
  - securityContext
  - livenessProbe
  - readinessProbe
  - lifecycle

### spec.monitor

MySQL managed by KubeDB can be monitored with builtin-Prometheus and CoreOS-Prometheus operator out-of-the-box. To learn more,

- [Monitor MySQL with builtin Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md)
- [Monitor MySQL with CoreOS Prometheus operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md)

### spec.configSource

`spec.configSource` is an optional field that allows users to provide custom configuration for MySQL. This field accepts a [`VolumeSource`](https://github.com/kubernetes/api/blob/release-1.11/core/v1/types.go#L47). So you can use any kubernetes supported volume source such as `configMap`, `secret`, `azureDisk` etc. To learn more about how to use a custom configuration file see [here](/docs/guides/mysql/configuration/using-custom-config.md).

### spec.podTemplate

KubeDB allows providing a template for database pod through `spec.podTemplate`. KubeDB operator will pass the information provided in `spec.podTemplate` to the StatefulSet created for MySQL database.

KubeDB accept following fields to set in `spec.podTemplate:`

- metadata:
  - annotations (pod's annotation)
- controller:
  - annotations (statefulset's annotation)
- spec:
  - args
  - env
  - resources
  - initContainers
  - imagePullSecrets
  - nodeSelector
  - affinity
  - schedulerName
  - tolerations
  - priorityClassName
  - priority
  - securityContext
  - livenessProbe
  - readinessProbe
  - lifecycle

Uses of some field of `spec.podTemplate` is described below,

#### spec.podTemplate.spec.args

`spec.podTemplate.spec.args` is an optional field. This can be used to provide additional arguments to database installation. To learn about available args of `mysqld`, visit [here](https://dev.mysql.com/doc/refman/8.0/en/server-options.html).

#### spec.podTemplate.spec.env

`spec.podTemplate.spec.env` is an optional field that specifies the environment variables to pass to the MySQL docker image. To know about supported environment variables, please visit [here](https://hub.docker.com/_/mysql/).

Note that, Kubedb does not allow `MYSQL_ROOT_PASSWORD`, `MYSQL_ALLOW_EMPTY_PASSWORD`, `MYSQL_RANDOM_ROOT_PASSWORD`, and `MYSQL_ONETIME_PASSWORD` environment variables to set in `spec.env`. If you want to set the root password, please use `spec.databaseSecret` instead described earlier.

If you try to set any of the forbidden environment variables i.e. `MYSQL_ROOT_PASSWORD` in MySQL crd, Kubed operator will reject the request with following error,

```ini
Error from server (Forbidden): error when creating "./mysql.yaml": admission webhook "mysql.validators.kubedb.com" denied the request: environment variable MYSQL_ROOT_PASSWORD is forbidden to use in MySQL spec
```

Also note that Kubedb does not allow to update the environment variables as updating them does not have any effect once the database is created.  If you try to update environment variables, Kubedb operator will reject the request with following error,

```ini
Error from server (BadRequest): error when applying patch:
...
for: "./mysql.yaml": admission webhook "mysql.validators.kubedb.com" denied the request: precondition failed for:
...At least one of the following was changed:
    apiVersion
    kind
    name
    namespace
    spec.databaseSecret
    spec.init
    spec.storageType
    spec.storage
    spec.podTemplate.spec.nodeSelector
    spec.podTemplate.spec.env
```

#### spec.podTemplate.spec.imagePullSecrets

`KubeDB` provides the flexibility of deploying MySQL database from a private Docker registry. `spec.podTemplate.spec.imagePullSecrets` is an optional field that points to secrets to be used for pulling docker image if you are using a private docker registry. To learn how to deploy MySQL from a private registry, please visit [here](/docs/guides/mysql/private-registry/using-private-registry.md).

#### spec.podTemplate.spec.nodeSelector

`spec.podTemplate.spec.nodeSelector` is an optional field that specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). To learn more, see [here](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) .

#### spec.podTemplate.spec.resources

`spec.podTemplate.spec.resources` is an optional field. This can be used to request compute resources required by the database pods. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).

### spec.serviceTemplate

You can also provide a template for the services created by KubeDB operator for MySQL database through `spec.serviceTemplate`. This will allow you to set the type and other properties of the services.

KubeDB allows following fields to set in `spec.serviceTemplate`:

- metadata:
  - annotations
- spec:
  - type
  - ports
  - clusterIP
  - externalIPs
  - loadBalancerIP
  - loadBalancerSourceRanges
  - externalTrafficPolicy
  - healthCheckNodePort

### spec.updateStrategy

You can specify [update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#update-strategies) of StatefulSet created by KubeDB for MySQL database thorough `spec.updateStrategy` field. The default value of this field is `RollingUpdate`. In future, we will use this field to determine how automatic migration from old KubeDB version to new one should behave.

### spec.terminationPolicy

`terminationPolicy` gives flexibility whether to `nullify`(reject) the delete operation of `MySQL` crd or which resources KubeDB should keep or delete when you delete `MySQL` crd. KubeDB provides following four termination policies:

- DoNotTerminate
- Pause (`Default`)
- Delete
- WipeOut

When, `terminationPolicy` is `DoNotTerminate`, KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `DoNotTerminate` feature. If admission webhook is enabled, `DoNotTerminate` prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`.

Following table show what KubeDB does when you delete MySQL crd for different termination policies,

|              Behaviour              | DoNotTerminate |  Pause   |  Delete  | WipeOut  |
| ----------------------------------- | :------------: | :------: | :------: | :------: |
| 1. Block Delete operation           |    &#10003;    | &#10007; | &#10007; | &#10007; |
| 2. Create Dormant Database          |    &#10007;    | &#10003; | &#10007; | &#10007; |
| 3. Delete StatefulSet               |    &#10007;    | &#10003; | &#10003; | &#10003; |
| 4. Delete Services                  |    &#10007;    | &#10003; | &#10003; | &#10003; |
| 5. Delete PVCs                      |    &#10007;    | &#10007; | &#10003; | &#10003; |
| 6. Delete Secrets                   |    &#10007;    | &#10007; | &#10007; | &#10003; |
| 7. Delete Snapshots                 |    &#10007;    | &#10007; | &#10007; | &#10003; |
| 8. Delete Snapshot data from bucket |    &#10007;    | &#10007; | &#10007; | &#10003; |

If you don't specify `spec.terminationPolicy` KubeDB uses `Pause` termination policy by default.

## Next Steps

- Learn how to use KubeDB to run a MySQL database [here](/docs/guides/mysql/README.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
