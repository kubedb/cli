---
title: Elasticsearch
menu:
  docs_0.12.0:
    identifier: elasticsearch-db
    name: Elasticsearch
    parent: databases
    weight: 10
menu_name: docs_0.12.0
section_menu_id: concepts
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Elasticsearch

## What is Elasticsearch

`Elasticsearch` is a Kubernetes `Custom Resource Definitions` (CRD). It provides declarative configuration for [Elasticsearch](https://www.elastic.co/products/elasticsearch) in a Kubernetes native way. You only need to describe the desired database configuration in an Elasticsearch object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## Elasticsearch Spec

As with all other Kubernetes objects, an Elasticsearch needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example Elasticsearch object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: e1
  namespace: demo
spec:
  version: "6.3-v1"
  replicas: 2
  authPlugin: "SearchGuard"
  enableSSL: true
  certificateSecret:
    secretName: e1-certs
  databaseSecret:
    secretName: e1-auth
  storageType: "Durable"
  storage:
    storageClassName: standard
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  init:
    scriptSource:
      configMap:
        name: es-init-script
  backupSchedule:
    cronExpression: "@every 2m"
    storageSecretName: gcs-secret
    gcs:
      bucket: kubedb-qa
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
      name: es-custom-config
  podTemplate:
    annotations:
      passMe: ToDatabasePod
    controller:
      annotations:
        passMe: ToStatefulSet
    spec:
      schedulerName: my-scheduler
      nodeSelector:
        disktype: ssd
      imagePullSecrets:
      - name: myregistrykey
      env:
      - name: ES_JAVA_OPTS
        value: "-Xms128m -Xmx128m"
      resources:
        requests:
          memory: "64Mi"
          cpu: "250m"
        limits:
          memory: "128Mi"
          cpu: "500m"
  serviceTemplate:
    annotations:
      passMe: ToService
    spec:
      type: NodePort
      ports:
      - name:  http
        port:  9200
        targetPort: http
  updateStrategy:
    type: "RollingUpdate"
  terminationPolicy: "DoNotTerminate"
```

### spec.version

`spec.version` is a required field that specifies the name of the [ElasticsearchVersion](/docs/concepts/catalog/elasticsearch.md) crd where the docker images are specified. Currently, when you install KubeDB, it creates the following `ElasticsearchVersion` crd,

- `5.6`, `5.6-v1`, `5.6.4`, `5.6.4-v1`
- `6.2`, `6.2-v1`, `6.2.4`, `6.2.4-v1`, `6.3`, `6.3-v1`, `6.3.0`, `6.3.0-v1`

### spec.topology

`spec.topology` is an optional field that provides a way to configure different types of nodes for Elasticsearch cluster. This field enables you to specify how many nodes you want to act as master, data and client node. You can also specify how much storage and resources to use for each types of nodes independently.

You can specify the following things in `spec.topology` field,

- `spec.topology.master`
    - `.replicas` is an optional field to specify how many pods we want as `master` node. If not set, this defaults to 1.
    - `.prefix` is an optional field to be used as the prefix of StatefulSet name.
    - `.storage` is an optional field that specifies how much storage to use for `master` node.
    - `.resources` is an optional field that specifies how much compute resources to request for `master` node.
- `spec.topology.data`
    - `.replicas` is an optional field to specify how many pods we want as `data` node. If not set, this defaults to 1.
    - `.prefix` is an optional field to be used as the prefix of StatefulSet name.
    - `.storage` is an optional field that specifies how much storage to use for `data` node.
    - `.resources` is an optional field that specifies how much compute resources to request for `data` node.
- `spec.topology.client`
    - `.replicas` is an optional field to specify how many pods we want as `client` node. If not set, this defaults to 1.
    - `.prefix` is an optional field to be used as the prefix of StatefulSet name.
    - `.storage` is an optional field that specifies how much storage to use for `client` node.
    - `.resources` is an optional field that specifies how much compute resources to request for `client` node.

> Note: Any two of them can't have the same prefix.

A sample configuration for `spec.topology` field is shown below,

```yaml
topology:
  master:
    prefix: master
    replicas: 1
    storage:
      storageClassName: "standard"
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
  data:
    prefix: data
    replicas: 3
    storage:
      storageClassName: "standard"
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
    resources:
      requests:
        memory: "512Mi"
        cpu: "250m"
      limits:
        memory: "1Gi"
        cpu: "500m"
  client:
    prefix: client
    replicas: 2
    storage:
      storageClassName: "standard"
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
```

If you specify `spec.topology` field then you are not allowed to specify following fields in Elasticsearch crd.

- `spec.replicas`
- `spec.storage`
- `spec.podTemplate.spec.resources`

> If you do not specify `spec.topology` field, all nodes of your Elasticsearch cluster will work as `master`, `data` and `client` node simultaneously.

### spec.replicas

`spec.replicas` is an optional field that can be used if `spec.topology` is not specified. This field specifies the number of pods in the Elasticsearch cluster. The default value of this field is 1.

### spec.authPlugin

`spec.authPlugin` is an optional field that specifies which plugin to use for authentication. Currently, this field accepts `None` or `SearchGuard`. By default, KubeDB uses [Search Guard](https://github.com/floragunncom/search-guard) for authentication. If you specify `None` in this field, KubeDB will disable Search Guard plugin and your database will not be protected anymore.

### spec.enableSSL

`spec.enableSSL` is an optional field that specifies whether to enable SSL for [Search Guard](https://github.com/floragunncom/search-guard). The default value of this field is `false`.

### spec.certificateSecret

`spec.certificateSecret` is an optional field that points a Secret used to hold the following information for the certificate.

- `root.pem:` The root CA in `pem` format
- `root.jks:` The root CA in `jks` format
- `node.jks:` The node certificate used for the transport layer
- `client.jks:` The client certificate used for http layer
- `sgadmin.jks:` The admin certificate used to change the Search Guard configuration.
- `key_pass:` The key password used to encrypt certificates.

If not set, KubeDB operator creates a new Secret `{Elasticsearch name}-cert` with generated certificates. If you want to use an existing secret, please specify that when creating Elasticsearch using `spec.certificateSecret.secretName`.

### spec.databaseSecret

`spec.databaseSecret` is an optional field that points to a Secret used to hold credentials and [search guard](https://github.com/floragunncom/search-guard) configuration.

Following keys are used to hold credentials

- `ADMIN_USERNAME:` Username for superuser.
- `ADMIN_PASSWORD:` Password for superuser.
- `READALL_USERNAME`  Username for `readall` user.
- `READALL_PASSWORD:` Password for `readall` user.

Following keys are used for search-guard configuration

- `sg_config.yml:` Configure authenticators and authorization backends
- `sg_internal_users.yml:` user and hashed passwords (hash with hasher.sh)
- `sg_roles_mapping.yml:` map backend roles, hosts and users to roles
- `sg_action_groups.yml:` define permission groups
- `sg_roles.yml:` define the roles and the associated permissions

If not set, KubeDB operator creates a new Secret `{Elasticsearch name}-auth` with generated credentials and default search-guard configuration. If you want to use an existing secret, please specify that when creating Elasticsearch using `spec.databaseSecret.secretName`.

### spec.storageType

`spec.storageType` is an optional field that specifies the type of storage to use for database. It can be either `Durable` or `Ephemeral`. The default value of this field is `Durable`. If `Ephemeral` is used then KubeDB will create Elasticsearch database using [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) volume. In this case, you don't have to specify `spec.storage` field.

### spec.storage

If you don't set `spec.storageType:`  to `Ephemeral` and if you don't specify `spec.topology` field then `spec.storage` field is required. This field specifies the StorageClass of PVCs dynamically allocated to store data for the database. This storage spec will be passed to the StatefulSet created by KubeDB operator to run database pods. You can specify any StorageClass available in your cluster with appropriate resource requests.

 - `spec.storage.storageClassName` is the name of the StorageClass used to provision PVCs. PVCs don’t necessarily have to request a class. A PVC with its storageClassName set equal to "" is always interpreted to be requesting a PV with no class, so it can only be bound to PVs with no class (no annotation or one set equal to ""). A PVC with no storageClassName is not quite the same and is treated differently by the cluster depending on whether the DefaultStorageClass admission plugin is turned on.
 - `spec.storage.accessModes` uses the same conventions as Kubernetes PVCs when requesting storage with specific access modes.
 - `spec.storage.resources` can be used to request specific quantities of storage. This follows the same resource model used by PVCs.

To learn how to configure `spec.storage`, please visit the links below:
 - https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims

### spec.init

`spec.init` is an optional section that can be used to initialize a newly created Elasticsearch cluster from prior snapshots. To initialize from prior snapshots, set the `spec.init.snapshotSource` section when creating an Elasticsearch object. In this case, SnapshotSource must have the following information:

 - `name:` Name of the Snapshot
 - `namespace:` Namespace of the Snapshot

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: elasticsearch-db
spec:
  version: "5.6-v1"
  databaseSecret:
    secretName: old-elasticsearch-auth
  init:
    snapshotSource:
      name: "snapshot-xyz"
      namespace: demo
```

In the above example, Elasticsearch cluster will be initialized from Snapshot `snapshot-xyz` in `demo` namespace. Here, KubeDB operator will launch a Job to initialize Elasticsearch, once StatefulSet pods are running. For details tutorial on how to initialize Elasticsearch from snapshot, please visit [here](/docs/guides/elasticsearch/initialization/snapshot_source.md).

### spec.backupSchedule

KubeDB supports taking periodic snapshots for Elasticsearch database. This is an optional section in `.spec`. When `spec.backupSchedule` section is added, KubeDB operator immediately takes a backup to validate this information. After that, at each tick kubeDB operator creates a [Snapshot](/docs/concepts/snapshot.md) object. This triggers operator to create a Job to take backup.

You have to specify following fields to take periodic backup of your Elasticsearch database:

 - `spec.backupSchedule.cronExpression` is a required [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26). This specifies the schedule for backup operations.
 - `spec.backupSchedule.{storage}` is a required field that is used as the destination for storing snapshot data. KubeDB supports cloud storage providers like S3, GCS, Azure, and OpenStack Swift. It also supports any locally mounted Kubernetes volumes, like NFS, Ceph, etc. Only one backend can be used at a time. To learn how to configure this, please visit [here](/docs/concepts/snapshot.md).

You can also specify a template for pod of backup job through `spec.backupSchedule.podTemplate`. KubeDB will use the information you have provided in `podTemplate` to create the backup job. KubeDB accept following fields to set in `spec.backupSchedule.podTemplate`:

- metadata
  - annotations (pod's annotation)
- controller
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

Elasticsearch managed by KubeDB can be monitored with builtin-Prometheus and CoreOS-Prometheus operator out-of-the-box. To learn more,

- [Monitor Elasticsearch with builtin Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md)
- [Monitor Elasticsearch with CoreOS Prometheus operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md)

### spec.configSource

`spec.configSource` is an optional field that allows users to provide custom configuration for Elasticsearch. This field accepts a [`VolumeSource`](https://github.com/kubernetes/api/blob/release-1.11/core/v1/types.go#L47). So you can use any kubernetes supported volume source such as `configMap`, `secret`, `azureDisk` etc. To learn more about how to use a custom configuration file see [here](/docs/guides/elasticsearch/custom-config/overview.md).

### spec.podTemplate

KubeDB allows providing a template for database pod through `spec.podTemplate`. KubeDB operator will pass the information provided in `spec.podTemplate` to the StatefulSet created for Elasticsearch database.

KubeDB accept following fields to set in `spec.podTemplate:`

- metadata
  - annotations (pod's annotation)
- controller
  - annotations (statefulset's annotation)
- spec:
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

Uses of some fields of `spec.podTemplate` are described below,

#### spec.podTemplate.spec.env

`spec.podTemplate.spec.env` is an optional field that specifies the environment variables to pass to the Elasticsearch docker image. To know about supported environment variables, please visit [here](https://github.com/pires/docker-elasticsearch#environment-variables).

A list of the supported environment variables, their permission to use in KubeDB and their default value is given below.

|      Environment variables      | Permission to use |                                           Default value                                            |
| ------------------------------- | :---------------: | -------------------------------------------------------------------------------------------------- |
| CLUSTER_NAME                    |     `allowed`     | `metadata.name`                                                                                    |
| NODE_NAME                       |   `not allowed`   | Pod name                                                                                           |
| NODE_MASTER                     |   `not allowed`   | KubeDB sets it based on `Elasticsearch` crd sepcification                                           |
| NODE_DATA                       |   `not allowed`   | KubeDB sets it based on `Elasticsearch` crd sepcification                                           |
| NETWORK_HOST                    |     `allowed`     | `_site_`                                                                                           |
| HTTP_ENABLE                     |     `allowed`     | If `spec.topology` is not specified then `true`. Otherwise, `false` for Master node and Data node. |
| HTTP_CORS_ENABLE                |     `allowed`     | `true`                                                                                             |
| HTTP_CORS_ALLOW_ORIGIN          |     `allowed`     | `*`                                                                                                |
| NUMBER_OF_MASTERS               |     `allowed`     | `(replicas/2)+1`                                                                                   |
| MAX_LOCAL_STORAGE_NODES         |     `allowed`     | `1`                                                                                                |
| ES_JAVA_OPTS                    |     `allowed`     | `-Xms128m -Xmx128m`                                                                                |
| ES_PLUGINS_INSTALL              |     `allowed`     | Not set                                                                                            |
| SHARD_ALLOCATION_AWARENESS      |     `allowed`     | `""`                                                                                               |
| SHARD_ALLOCATION_AWARENESS_ATTR |     `allowed`     | `""`                                                                                               |
| MEMORY_LOCK                     |     `allowed`     | `true`                                                                                             |
| REPO_LOCATIONS                  |     `allowed`     | `""`                                                                                               |
| PROCESSORS                      |     `allowed`     | `1`                                                                                                |

Note that, KubeDB does not allow `NODE_NAME`, `NODE_MASTER`, and `NODE_DATA` environment variables to set in `spec.podTemplate.spec.env`. KubeDB operator set them based on Elasticsearch crd specification.

If you try to set any these forbidden environment variable in Elasticsearch crd, KubeDB operator will reject the request with following error,

```ini
Error from server (Forbidden): error when creating "./elasticsearch.yaml": admission webhook "elasticsearch.validators.kubedb.com" denied the request: environment variable NODE_NAME is forbidden to use in Elasticsearch spec
```

Also, note that KubeDB does not allow to update the environment variables as updating them does not have any effect once the database is created.  If you try to update environment variables, KubeDB operator will reject the request with following error,

```ini
Error from server (BadRequest): error when applying patch:
...
for: "./elasticsearch.yaml": admission webhook "elasticsearch.validators.kubedb.com" denied the request: precondition failed for:
...
At least one of the following was changed:
    apiVersion
    kind
    name
    namespace
    spec.version
    spec.topology.*.prefix
    spec.topology.*.storage
    spec.enableSSL
    spec.certificateSecret
    spec.databaseSecret
    spec.storageType
    spec.storage
    spec.nodeSelector
    spec.init
    spec.env
```

#### spec.podTemplate.spec.imagePullSecrets

`spec.podTemplate.spec.imagePullSecrets` is an optional field that points to secrets to be used for pulling docker image when you are using a private docker registry. For more details on how to use private docker registry, please visit [here](/docs/guides/elasticsearch/private-registry/using-private-registry.md).

#### spec.podTemplate.spec.nodeSelector

`spec.podTemplate.spec.nodeSelector` is an optional field that specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). To learn more, see [here](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) .

#### spec.podTemplate.spec.resources

`spec.podTemplate.spec.resources` is an optional field. If you didn't specify `spec.topology` field then this can be used to request compute resources required by the database pods. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).

### spec.serviceTemplate

You can also provide a template for the services created by KubeDB operator for Elasticsearch database through `spec.serviceTemplate`. This will allow you to set the type and other properties of the services.

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

You can specify [update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#update-strategies) of StatefulSet created by KubeDB for Elasticsearch database thorough `spec.updateStrategy` field. The default value of this field is `RollingUpdate`. In future, we will use this field to determine how automatic migration from old KubeDB version to new one should behave.

### spec.terminationPolicy

`terminationPolicy` gives flexibility whether to `nullify`(reject) the delete operation of `Elasticsearch` crd or which resources KubeDB should keep or delete when you delete `Elasticsearch` crd. KubeDB provides following four termination policies:

- DoNotTerminate
- Pause (`Default`)
- Delete
- WipeOut

When, `terminationPolicy` is `DoNotTerminate`, KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to provide safety from accidental deletion of database. If admission webhook is enabled, KubeDB prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`.

Following table show what KubeDB does when you delete Elasticsearch crd for different termination policies,

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

- Learn how to use KubeDB to run an Elasticsearch database [here](/docs/guides/elasticsearch/README.md).
- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
