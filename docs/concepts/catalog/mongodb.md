---
title: MongoDBVersion
menu:
  docs_0.9.0:
    identifier: mongodb-version
    name: MongoDBVersion
    parent: catalog
    weight: 30
menu_name: docs_0.9.0
section_menu_id: concepts
---

# MongoDBVersion

## What is MongoDBVersion

`MongoDBVersion` is a Kubernetes `Custom Resource Definitions` (CRD). It provides a declarative configuration to specify the docker images to be used for [MongoDB](https://www.mongodb.com/) database deployed with KubeDB in Kubernetes native way.

When you install KubeDB, `MongoDBVersion` crd will be created automatically for every supported MongoDB versions. You have to specify the name of `MongoDBVersion` crd in `spec.version` field of [MongoDB](/docs/concepts/databases/mongodb.md) crd. Then, KubeDB will use the docker images specified in the `MongoDBVersion` crd to create your expected database.

Using a separate crd for specifying respective docker images allows us to modify the images independent of KubeDB operator. This will also allow the users to use a custom image for the database.

## MongoDBVersion Specification

As with all other Kubernetes objects, a MongoDBVersion needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section.

```yaml
apiVersion: catalog.kubedb.com/v1alpha1
kind: MongoDBVersion
metadata:
  name: "3.6-v2"
  labels:
    app: kubedb
spec:
  version: "3.6"
  deprecated: false
  db:
    image: "${KUBEDB_DOCKER_REGISTRY}/mongo:3.6-v1"
  exporter:
    image: "${KUBEDB_DOCKER_REGISTRY}/mongodb_exporter:v1.0.0"
  tools:
    image: "${KUBEDB_DOCKER_REGISTRY}/mongo-tools:3.6-v1"
```

### metadata.name

`metadata.name` is a required field that specify the name of the `MongoDBVersion` crd. You have to specify this name in `spec.version` field of [MongoDB](/docs/concepts/databases/mongodb.md) crd.

We follow this convention for naming MongoDBVersion crd:

- Name format: `{Original PosgreSQL image verion}-{modification tag}`

We modify original MongoDB docker image to support MongoDB clustering and re-tag the image with v1, v2 etc. modification tag. An image with higher modification tag will have more feature than the images with lower modification tag. Hence, it is recommended to use MongoDBVersion crd with highest modification tag to enjoy the latest features.

### spec.version

`spec.version` is a required field that specifies the original version of MongoDB database that has been used to build the docker image specified in `spec.db.image` field.

### spec.deprecated

`spec.deprecated` is an optional field that specifies whether the docker images specified here is supported by the current KubeDB operator. For example, we have modified `kubedb/mongo:3.6` docker image to support MongoDB clustering and re-tagged as `kubedb/mongo:3.6-v1`. So, we have marked `kubedb/mongo:3.6` as deprecated for KubeDB `0.9.0-rc.0`.

The default value of this field is `false`. If `spec.depcrecated` is set to `true`, KubeDB operator will skip processing this CRD object and will add a event to the CRD object specifying that the DB version is deprecated.

### spec.db.image

`spec.db.image` is a required field that specifies the docker image which will be used to create Statefulset by KubeDB operator to create expected MongoDB database.

### spec.exporter.image

`spec.exporter.image` is required field that specifies the image which will be used to export Prometheus metrics.

### spec.tools.image

`spec.tools.image` is a required field that specifies the image which will be used to take backup and initialize database from a snapshot.

## Next Steps

- Learn about MongoDB crd [here](/docs/concepts/databases/mongodb.md).
- Deploy your first MongoDB database with KubeDB by following the guide [here](/docs/guides/mongodb/quickstart/quickstart.md).
