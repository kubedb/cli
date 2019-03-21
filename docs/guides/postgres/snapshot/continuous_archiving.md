---
title: Continuous Archiving of PostgreSQL
menu:
  docs_0.11.0:
    identifier: pg-continuous-archiving-snapshot
    name: WAL Archiving
    parent: pg-snapshot-postgres
    weight: 20
menu_name: docs_0.11.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Continuous Archiving with WAL-G

KubeDB also supports continuous archiving of PostgreSQL using [WAL-G ](https://github.com/wal-g/wal-g). Users can use any one of `s3`, `gcs`, `azure`, or `swift` as cloud storage destination. 

**What is this Continuous Archiving**

PostgreSQL maintains a write ahead log (WAL) in the `pg_xlog/` subdirectory of the cluster's data directory.  The existence of the log makes it possible to use a third strategy for backing up databases and if recovery is needed, restore from the backed-up WAL files to bring the system back to last known state.

**Continuous Archiving Setup**

Following additional parameters are set in `postgresql.conf` for *primary* server

```console
archive_command = 'wal-g wal-push %p'
archive_timeout = 60
```

Here, these commands are used to push files to the cloud.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

> Note: YAML files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Next Steps

- Learn about archiving to [Amazon S3](/docs/guides/postgres/snapshot/archiving_to_s3.md).
- Learn about archiving to [Google Cloud Storage](/docs/guides/postgres/snapshot/archiving_to_gcs.md).
- Learn about archiving to [Azure Storage](/docs/guides/postgres/snapshot/archiving_to_azure.md).
- Learn about archiving to [OpenStack Object Storage (Swift)](/docs/guides/postgres/snapshot/archiving_to_swift.md).
