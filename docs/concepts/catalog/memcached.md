---
title: MemcachedVersion
menu:
  docs_0.11.0:
    identifier: memcached-version
    name: MemcachedVersion
    parent: catalog
    weight: 30
menu_name: docs_0.11.0
section_menu_id: concepts
---

# MemcachedVersion

## What is MemcachedVersion

`MemcachedVersion` is a Kubernetes `Custom Resource Definitions` (CRD). It provides a declarative configuration to specify the docker images to be used for [Memcached](https://memcached.org) database deployed with KubeDB in Kubernetes native way.

When you install KubeDB, `MemcachedVersion` crd will be created automatically for every supported Memcached versions. You have to specify the name of `MemcachedVersion` crd in `spec.version` field of [Memcached](/docs/concepts/databases/memcached.md) crd. Then, KubeDB will use the docker images specified in the `MemcachedVersion` crd to create your expected database.

Using a separate crd for specifying respective docker images allows us to modify the images independent of KubeDB operator. This will also allow the users to use a custom image for the database.

## MemcachedVersion Specification

As with all other Kubernetes objects, a MemcachedVersion needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section.

```yaml
apiVersion: catalog.kubedb.com/v1alpha1
kind: MemcachedVersion
metadata:
  name: "1.5.4-v1"
  labels:
    app: kubedb
spec:
  version: "1.5.4"
  db:
    image: "${KUBEDB_DOCKER_REGISTRY}/memcached:1.5.4-v1"
  exporter:
    image: "${KUBEDB_DOCKER_REGISTRY}/memcached-exporter:v0.4.1"
```

### metadata.name

`metadata.name` is a required field that specify the name of the `MemcachedVersion` crd. You have to specify this name in `spec.version` field of [Memcached](/docs/concepts/databases/memcached.md) crd.

We follow this convention for naming MemcachedVersion crd:

- Name format: `{Original Memcached image version}-{modification tag}`

We modify original Memcached docker image to support additional features. An image with higher modification tag will have more feature than the images with lower modification tag. Hence, it is recommended to use MemcachedVersion crd with highest modification tag to take advantage of the latest features.

### spec.version

`spec.version` is a required field that specifies the original version of Memcached server that has been used to build the docker image specified in `spec.db.image` field.

### spec.deprecated

`spec.deprecated` is an optional field that specifies whether the docker images specified here is supported by the current KubeDB operator. For example, we have modified `kubedb/memcached:1.5.4` docker image to support custom configuration and re-tagged as `kubedb/memcached:1.5.4-v1`. Now, KubeDB `0.9.0-rc.0` supports providing custom configuration which required `kubedb/memcached:1.5.4-v1` docker image. So, we have marked `kubedb/memcached:1.5.4` as deprecated for KubeDB `0.9.0-rc.0`.

The default value of this field is `false`. If `spec.depcrecated` is set `true`, KubeDB operator will not create the database and other respective resources for this version.

### spec.db.image

`spec.db.image` is a required field that specifies the docker image which will be used to create Statefulset by KubeDB operator to create expected Memcached server.

### spec.exporter.image

`spec.exporter.image` is required field that specifies the image which will be used to export Prometheus metrics.

## Next Steps

- Learn about Memcached crd [here](/docs/concepts/databases/memcached.md).
- Deploy your first Memcached server with KubeDB by following the guide [here](/docs/guides/memcached/quickstart/quickstart.md).
