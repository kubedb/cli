---
title: Weclome | KubeDB
description: Welcome to CLI
menu:
  docs_0.8.0-beta.0:
    identifier: readme-cli
    name: Readme
    parent: welcome
    weight: -1
menu_name: docs_0.8.0-beta.0
section_menu_id: welcome
url: /docs/0.8.0-beta.0/welcome/
aliases:
  - /docs/0.8.0-beta.0/
  - /docs/0.8.0-beta.0/README/
---

# KubeDB

Running production quality databases in Kubernetes can be tricky. KubeDB is a framework for writing operators for any database that support the following operational requirements:

 - Create a database declaratively using CRD.
 - Take one-off backups or period backups to various cloud stores, eg,, S3, GCS, etc.
 - Restore from backup or clone any database.
 - Native integration with Prometheus for monitoring via [CoreOS Prometheus Operator](https://github.com/coreos/prometheus-operator).
 - Apply deletion lock to avoid accidental deletion of database.
 - Keep track of deleted databases, cleanup prior snapshots with a single command.
 - Use cli to manage databases like kubectl for Kubernetes.

Currently KubeDB includes support for following datastores:
 - Postgres
 - Elasticsearch
 - MySQL
 - MongoDB
 - Redis
 - Memcached

From here you can learn all about KubeDB's architecture and how to deploy and use KubeDB.

- [Concepts](/docs/concepts/). Concepts explain some significant aspect of KubeDB. This is where you can learn about what KubeDB does and how it does it.

- [Setup](/docs/setup/). Setup contains instructions for installing the KubeDB in various cloud providers.

- [Guides](/docs/guides/). Guides show you how to perform tasks with KubeDB.

- [Reference](/docs/reference/). Detailed exhaustive lists of command-line options, configuration options, API definitions, and procedures.

We're always looking for help improving our documentation, so please don't hesitate to [file an issue](https://github.com/kubedb/project/issues/new) if you see some problem. Or better yet, submit your own [contributions](/docs/CONTRIBUTING.md) to help make our docs better.

---

**KubeDB binaries collects anonymous usage statistics to help us learn how the software is being used and how we can improve it. To disable stats collection, run the operator with the flag** `--analytics=false`.

---
