---
title: PostgresVersion
menu:
  docs_0.12.0:
    identifier: postgres-version
    name: PostgresVersion
    parent: catalog
    weight: 30
menu_name: docs_0.12.0
section_menu_id: concepts
---

# PostgresVersion

## What is PostgresVersion

`PostgresVersion` is a Kubernetes `Custom Resource Definitions` (CRD). It provides a declarative configuration to specify the docker images to be used for [PostgreSQL](https://www.postgresql.org/) database deployed with KubeDB in Kubernetes native way.

When you install KubeDB, `PostgresVersion` crd will be created automatically for every supported PostgreSQL versions. You have to specify the name of `PostgresVersion` crd in `spec.version` field of [Postgres](/docs/concepts/databases/postgres.md) crd. Then, KubeDB will use the docker images specified in the `PostgresVersion` crd to create your expected database.

Using a separate crd for specifying respective docker images, and pod security policy names allow us to modify the images, and policies independent of KubeDB operator. This will also allow the users to use a custom image for the database. For more details about how to use custom image with Postgres in KubeDB, please visit [here](/docs/guides/postgres/custom-versions/setup.md).

## PostgresVersion Specification

As with all other Kubernetes objects, a PostgresVersion needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section.

```yaml
apiVersion: catalog.kubedb.com/v1alpha1
kind: PostgresVersion
metadata:
  name: "10.2-v2"
  labels:
    app: kubedb
spec:
  version: "10.2"
  deprecated: false
  db:
    image: "kubedb/postgres:10.2-v2"
  exporter:
    image: "kubedb/postgres_exporter:v0.4.6"
  tools:
    image: "kubedb/postgres-tools:10.2-v2"
  podSecurityPolicies:
    databasePolicyName: "postgres-db"
    snapshotterPolicyName: "postgres-snapshot"
```

### metadata.name

`metadata.name` is a required field that specifies the name of the `PostgresVersion` crd. You have to specify this name in `spec.version` field of [Postgres](/docs/concepts/databases/postgres.md) crd.

We follow this convention for naming PostgresVersion crd:
- Name format: `{Original PostgreSQL image version}-{modification tag}`

We modify original PostgreSQL docker image to support additional features like WAL archiving, clustering etc. and re-tag the image with v1, v2 etc. modification tag. An image with higher modification tag will have more feature than the images with lower modification tag. Hence, it is recommended to use PostgresVersion crd with highest modification tag to take advantage of the latest features.

### spec.version

`spec.version` is a required field that specifies the original version of PostgreSQL database that has been used to build the docker image specified in `spec.db.image` field.

### spec.deprecated

`spec.deprecated` is an optional field that specifies whether the docker images specified here is supported by the current KubeDB operator. For example, we have modified `kubedb/postgres:10.2` docker image to support custom configuration and re-tagged as `kubedb/postgres:10.2-v2`. Now, KubeDB `0.9.0-rc.0` supports providing custom configuration which required `kubedb/postgres:10.2-v2` docker image. So, we have marked `kubedb/postgres:10.2` as deprecated in KubeDB `0.9.0-rc.0`.

The default value of this field is `false`. If `spec.depcrecated` is set `true`, KubeDB operator will not create the database and other respective resources for this version.

### spec.db.image

`spec.db.image` is a required field that specifies the docker image which will be used to create Statfulset by KubeDB operator to create expected PostgreSQL database.

### spec.exporter.image

`spec.exporter.image` is a required field that specifies the image which will be used to export Prometheus metrics.

### spec.tools.image

`spec.tools.image` is a required field that specifies the image which will be used to take backup and initialize database from snapshot.

### spec.podSecurityPolicies.databasePolicyName

`spec.podSecurityPolicies.databasePolicyName` is a required field that specifies the name of the pod security policy required to get the database server pod(s) running.

### spec.podSecurityPolicies.snapshotterPolicyName

`spec.podSecurityPolicies.snapshotterPolicyName` is a required field that specifies the name of the pod security policy required to get the snapshotter pod(s) running. To use user-defined policies, names of the policies have to be set in `spec.podSecurityPolicies` and in the list of allowed policy names in KubeDB operator like below:

```console
helm upgrade kubedb-operator appscode/kubedb --namespace kube-system \
  --set additionalPodSecurityPolicies[0]=custom-db-policy \
  --set additionalPodSecurityPolicies[1]=custom-snapshotter-policy
```

## Next Steps

- Learn about Postgres crd [here](/docs/concepts/databases/postgres.md).
- Deploy your first PostgreSQL database with KubeDB by following the guide [here](/docs/guides/postgres/quickstart/quickstart.md).