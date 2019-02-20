---
title: MySQLVersion
menu:
  docs_0.9.0:
    identifier: mysql-version
    name: MySQLVersion
    parent: catalog
    weight: 30
menu_name: docs_0.9.0
section_menu_id: concepts
---

# MySQLVersion

## What is MySQLVersion

`MySQLVersion` is a Kubernetes `Custom Resource Definitions` (CRD). It provides a declarative configuration to specify the docker images to be used for [MySQL](https://www.mysql.org/) database deployed with KubeDB in Kubernetes native way.

When you install KubeDB, `MySQLVersion` crd will be created automatically for every supported MySQL versions. You have to specify the name of `MySQLVersion` crd in `spec.version` field of [MySQL](/docs/concepts/databases/mysql.md) crd. Then, KubeDB will use the docker images specified in the `MySQLVersion` crd to create your expected database.

Using a separate crd for specifying respective docker images allows us to modify the images independent of KubeDB operator. This will also allow the users to use a custom image for the database.

## MySQLVersion Specification

As with all other Kubernetes objects, a MySQLVersion needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section.

```yaml
apiVersion: catalog.kubedb.com/v1alpha1
kind: MySQLVersion
metadata:
  name: "8.0-v2"
  labels:
    app: kubedb
spec:
  version: "8.0"
  db:
    image: "${KUBEDB_DOCKER_REGISTRY}/mysql:8.0-v2"
  exporter:
    image: "${KUBEDB_DOCKER_REGISTRY}/mysqld-exporter:v0.11.0"
  tools:
    image: "${KUBEDB_DOCKER_REGISTRY}/mysql-tools:8.0-v2"
```

### metadata.name

`metadata.name` is a required field that specify the name of the `MySQLVersion` crd. You have to specify this name in `spec.version` field of [MySQL](/docs/concepts/databases/mysql.md) crd.

We follow this convention for naming MySQLVersion crd:

- Name format: `{Original MySQL image version}-{modification tag}`

We modify original MySQL docker image to support additional features. An image with higher modification tag will have more feature than the images with lower modification tag. Hence, it is recommended to use MySQLVersion crd with highest modification tag to take advantage of the latest features.

### spec.version

`spec.version` is a required field that specifies the original version of MySQL database that has been used to build the docker image specified in `spec.db.image` field.

### spec.deprecated

`spec.deprecated` is an optional field that specifies whether the docker images specified here is supported by the current KubeDB operator. For example, we have modified `kubedb/mysql:8.0` docker image to support custom configuration and re-tagged as `kubedb/mysql:8.0-v2`. Now, KubeDB `0.9.0-rc.0` supports providing custom configuration which required `kubedb/mysql:8.0-v2` docker image. So, we have marked `kubedb/mysql:8.0` as deprecated for KubeDB `0.9.0-rc.0`.

The default value of this field is `false`. If `spec.depcrecated` is set `true`, KubeDB operator will not create the database and other respective resources for this version.

### spec.db.image

`spec.db.image` is a required field that specifies the docker image which will be used to create Statefulset by KubeDB operator to create expected MySQL database.

### spec.exporter.image

`spec.exporter.image` is required field that specifies the image which will be used to export Prometheus metrics.

### spec.tools.image

`spec.tools.image` is a required field that specifies the image which will be used to take backup and initialize database from snapshot.

## Next Steps

- Learn about MySQL crd [here](/docs/concepts/databases/mysql.md).
- Deploy your first MySQL database with KubeDB by following the guide [here](/docs/guides/mysql/quickstart/quickstart.md).
