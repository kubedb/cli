---
title: Memcached
menu:
  docs_0.9.0-beta.0:
    identifier: memcached-db
    name: Memcached
    parent: databases
    weight: 15
menu_name: docs_0.9.0-beta.0
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
  version: 1.5.3-alpine
  nodeSelector:
    disktype: ssd
  doNotPause: true
  configSource:
    configMap:
      name: mc-custom-config
  env:
    - name:  MY_ENV
      value: "env-example"
  monitor:
    agent: coreos-prometheus-operator
    prometheus:
      namespace: demo
      labels:
        app: kubedb
      interval: 10s
  resources:
    requests:
      memory: "64Mi"
      cpu: "250m"
    limits:
      memory: "128Mi"
      cpu: "500m"
```

### spec.replicas

`spec.replicas` is an optional field that specifies the number of desired Instances/Replicas of Memcached database. If you do not specify .spec.replicas, then it defaults to 1.

### spec.version

`spec.version` is a required field specifying the version of Memcached database. Here the database version is [`1.5.3-alpine`](https://hub.docker.com/r/library/memcached/tags/).

### spec.nodeSelector

`spec.nodeSelector` is an optional field that specifies a map of key-value pairs. For the pod to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels (it can have additional labels as well). To learn more, see [here](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) .

### spec.doNotPause

`spec.doNotPause` is an optional field that tells KubeDB operator that if this Memcached object is deleted, whether it should be reverted automatically. This should be set to `true` for production databases to avoid accidental deletion. If not set or set to false, deleting a Memcached object put the database into a dormant state.

### spec.imagePullSecret

`KubeDB` provides the flexibility of deploying Memcached database from a private Docker registry. To learn how to deploym Memcached from a private registry, please visit [here](/docs/guides/memcached/private-registry/using-private-registry.md).

### spec.configSource

`spec.configSource` is an optional field that allows users to provide custom configuration for Memcached. This field accepts a [`VolumeSource`](https://github.com/kubernetes/api/blob/release-1.11/core/v1/types.go#L47). So you can use any kubernetes supported volume source such as `configMap`, `secret`, `azureDisk` etc. To learn more about how to use a custom configuration file see [here](/docs/guides/memcached/custom-config/using-custom-config.md).

### spec.env

`spec.env` is an optional field that specifies the environment variables to pass to the Memcached docker image.

Note that, Kubedb does not allow to update the environment variables. If you try to update environment variables, Kubedb operator will reject the request with following error,

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
	spec.version
	spec.nodeSelector
	spec.env
```

### spec.monitor

Memcached managed by KubeDB can be monitored with builtin-Prometheus and CoreOS-Prometheus operator out-of-the-box. To learn more,

- [Monitor Memcached with builtin Prometheus](/docs/guides/memcached/monitoring/using-builtin-prometheus.md)
- [Monitor Memcached with CoreOS Prometheus operator](/docs/guides/memcached/monitoring/using-coreos-prometheus-operator.md)

### spec.resources

`spec.resources` is an optional field. This can be used to request compute resources required by the database pods. To learn more, visit [here](http://kubernetes.io/docs/user-guide/compute-resources/).

## Next Steps

- Learn how to use KubeDB to run a Memcached database [here](/docs/guides/memcached/README.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
