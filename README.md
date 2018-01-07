[![Go Report Card](https://goreportcard.com/badge/github.com/kubedb/cli)](https://goreportcard.com/report/github.com/kubedb/cli)

[Website](https://kubedb.com) • [Slack](http://slack.kubernetes.io) • [Twitter](https://twitter.com/KubeDB)

# KubeDB
Running production quality database in Kubernetes can be tricky to say the least. In the early days of Kubernetes, replication controllers were used to run a single pod for a database. With the introduction of StatefulSet, it became easy to run a docker container for any database. But what about monitoring, taking periodic backups, restoring from backups or cloning from an existing database? KubeDB is a framework for writing operators for any database that support the following operational requirements:

 - Create a database declaratively using CRD.
 - Take one-off backups or periodic backups to various cloud stores, eg, S3, GCS, Azure etc.
 - Restore from backup or clone any database.
 - Native integration with Prometheus for monitoring via [CoreOS Prometheus Operator](https://github.com/coreos/prometheus-operator).
 - Apply deletion lock to avoid accidental deletion of database.
 - Keep track of deleted databases, cleanup prior snapshots with a single command.
 - Use cli to manage databases like kubectl for Kubernetes.

KubeDB is developed at [AppsCode](https://twitter.com/AppsCodeHQ) to run their SAAS platform on Kubernetes. Currently KubeDB includes support for following datastores:
 - Postgres
 - Elasticsearch
 - MySQL
 - MongoDB
 - Redis
 - Memcached

## Supported Versions
Please pick a version of KubeDB that matches your Kubernetes installation.

| KubeDB Version                                                                     | Docs                                                        | Kubernetes Version |
|------------------------------------------------------------------------------------|-------------------------------------------------------------|--------------------|
| [0.8.0-beta.0](https://github.com/kubedb/cli/releases/tag/0.8.0-beta.0) (uses CRD) | [User Guide](https://kubedb.com/docs/0.8.0-beta.0/)         | 1.7.x +            |
| [0.6.0](https://github.com/kubedb/cli/releases/tag/0.6.0) (uses TPR)               | [User Guide](https://github.com/kubedb/cli/tree/0.6.0/docs) | 1.5.x - 1.7.x      |

## Installation
To install KubeDB, please follow the guide [here](https://kubedb.com/docs/latest/setup/install/).

## Using KubeDB
Want to learn how to use KubeDB? Please start [here](https://kubedb.com/docs/latest/guides/).

## Contribution guidelines
Want to help improve KubeDB? Please start [here](https://kubedb.com/docs/latest/welcome/contributing/).

---

**KubeDB binaries collects anonymous usage statistics to help us learn how the software is being used and how we can improve it. To disable stats collection, run the operator with the flag** `--analytics=false`.

---

## Support
We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [Kubernetes Slack team](https://kubernetes.slack.com/messages/C8149MREV/) channel `#kubedb`. To sign up, use our [Slack inviter](http://slack.kubernetes.io/).

To receive product annoucements, please join our [mailing list](https://groups.google.com/forum/#!forum/kubedb) or follow us on [Twitter](https://twitter.com/KubeDB). Our mailing list is also used to share design docs shared via Google docs.

If you have found a bug with KubeDB or want to request for new features, please [file an issue](https://github.com/kubedb/project/issues/new).
