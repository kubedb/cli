---
title: KubeDB Overview
description: KubeDB Overview
menu:
  docs_0.8.0-beta.2:
    identifier: overview-concepts
    name: Overview
    parent: what-is-kubedb
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: concepts
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
