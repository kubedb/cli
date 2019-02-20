---
title: Memcached
menu:
  docs_0.9.0:
    identifier: memcached-db
    name: Memcached
    parent: databases
    weight: 15
menu_name: docs_0.9.0
section_menu_id: concepts
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Memcached

## What is Memcached

`Memcached` is a Kubernetes `Custom Resource Definitions` (CRD). It provides declarative configuration for [Memcached](https://memcached.org/) in a Kubernetes native way. You only need to describe the desired database configuration in a Memcached object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## Memcached Spec

As with all other Kubernetes objects, a Memcached needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example of a Memcached object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Memcached
metadata:
  name: mc1
  namespace: demo
spec:
  replicas: 3
  version: 1.5.3-v1
  monitor:
    agent: coreos-prometheus-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
  configSource:
    configMap:
      name: mc-custom-config
  podTemplate:
    annotations:
      passMe: ToDatabasePod
    controller:
      annotations:
        passMe: ToDeployment
    spec:
      schedulerName: my-scheduler
      nodeSelector:
        disktype: ssd
      imagePullSecrets:
      - name: myregistrykey
      args:
      - "-u memcache"
      env:
      - name: TEST_ENV
        value: "value"
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
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
```

### spec.replicas

`spec.replicas` is an optional field that specifies the number of desired Instances/Replicas of Memcached server. If you do not specify .spec.replicas, then it defaults to 1.

### spec.version

`spec.version` is a required field specifying the name of the [MemcachedVersion](/docs/concepts/catalog/memcached.md) crd where the docker images are specified. Currently, when you install KubeDB, it creates the following `MemcachedVersion` crd,

- `1.5.4`, `1.5.4-v1`, `1.5`, `1.5-v1`

### spec.monitor

Memcached managed by KubeDB can be monitored with builtin-Prometheus and CoreOS-Prometheus operator out-of-the-box. To learn more,

- [Monitor Memcached with builtin Prometheus](/docs/guides/memcached/monitoring/using-builtin-prometheus.md)
- [Monitor Memcached with CoreOS Prometheus operator](/docs/guides/memcached/monitoring/using-coreos-prometheus-operator.md)

### spec.configSource

`spec.configSource` is an optional field that allows users to provide custom configuration for Memcached. This field accepts a [`VolumeSource`](https://github.com/kubernetes/api/blob/release-1.11/core/v1/types.go#L47). So you can use any kubernetes supported volume source such as `configMap`, `secret`, `azureDisk` etc. To learn more about how to use a custom configuration file see [here](/docs/guides/memcached/custom-config/using-custom-config.md).

### spec.podTemplate

KubeDB allows providing a template for database pod through `spec.podTemplate`. KubeDB operator will pass the information provided in `spec.podTemplate` to the Deployment created for Memcached server.

KubeDB accept following fields to set in `spec.podTemplate:`

- metadata
  - annotations (pod's annotation)
- controller
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

`spec.podTemplate.spec.args` is an optional field. This can be used to provide additional arguments to database installation.

#### spec.podTemplate.spec.env

`spec.env` is an optional field that specifies the environment variables to pass to the Memcached docker image.

Note that, KubeDB does not allow to update the environment variables. If you try to update environment variables, KubeDB operator will reject the request with following error,

```ini
Error from server (BadRequest): error when applying patch:
...
for: "./mc.yaml": admission webhook "memcached.validators.kubedb.com" denied the request: precondition failed for:
...
At least one of the following was changed:
	apiVersion
	kind
	name
	namespace
	spec.podTemplate.spec.nodeSelector
    spec.podTemplate.spec.env
```

#### spec.podTemplate.spec.imagePullSecrets

`KubeDB` provides the flexibility of deploying Memcached server from a private Docker registry. To learn how to deploym Memcached from a private registry, please visit [here](/docs/guides/memcached/private-registry/using-private-registry.md).

#### spec.podTemplate.spec.nodeSelector

`spec.nodeSelector` is an optional field that specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). To learn more, see [here](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) .

#### spec.podTemplate.spec.resources

`spec.resources` is an optional field. This can be used to request compute resources required by the database pods. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).

### spec.serviceTemplate

You can also provide a template for the services created by KubeDB operator for Memcached server through `spec.serviceTemplate`. This will allow you to set the type and other properties of the services.

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

You can specify [update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#update-strategies) of Deployment created by KubeDB for Memcached server thorough `spec.updateStrategy` field. The default value of this field is `RollingUpdate`. In future, we will use this field to determine how automatic migration from old KubeDB version to new one should behave.

### spec.terminationPolicy

`terminationPolicy` gives flexibility whether to `nullify`(reject) the delete operation of `Memcached` crd or which resources KubeDB should keep or delete when you delete `Memcached` crd. KubeDB provides following four termination policies:

- DoNotTerminate
- Pause (`Default`)
- Delete
- WipeOut

When, `terminationPolicy` is `DoNotTerminate`, KubeDB takes advantage of `ValidationWebhook` feature in Kubernetes 1.9.0 or later clusters to implement `DoNotTerminate` feature. If admission webhook is enabled, `DoNotTerminate` prevents users from deleting the database as long as the `spec.terminationPolicy` is set to `DoNotTerminate`.

Following table show what KubeDB does when you delete Memcached crd for different termination policies,

|              Behaviour              | DoNotTerminate |  Pause   |  Delete  | WipeOut  |
| ----------------------------------- | :------------: | :------: | :------: | :------: |
| 1. Block Delete operation           |    &#10003;    | &#10007; | &#10007; | &#10007; |
| 2. Create Dormant Database          |    &#10007;    | &#10003; | &#10007; | &#10007; |
| 3. Delete StatefulSet               |    &#10007;    | &#10003; | &#10003; | &#10003; |
| 4. Delete Services                  |    &#10007;    | &#10003; | &#10003; | &#10003; |

If you don't specify `spec.terminationPolicy` KubeDB uses `Pause` termination policy by default.

## Next Steps

- Learn how to use KubeDB to run a Memcached server [here](/docs/guides/memcached/README.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
