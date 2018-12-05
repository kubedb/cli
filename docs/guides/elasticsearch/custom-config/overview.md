---
title: Run Elasticsearch with Custom Configuration
menu:
  docs_0.9.0:
    identifier: es-custom-config-overview
    name: Overview
    parent: es-custom-config
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Custom Configuration File

KubeDB supports providing custom configuration for an Elasticsearch cluster. This tutorial will give you an overview of how custom configuration files for an Elasticsearch cluster are managed in KubeDB.

## Overview

Elasticsearch is configured using `elasticsearch.yml` file. In KubeDB, this file is located in `/elasticsearch/config` directory of elasticsearch pods. To know more about configuring the Elasticsearch cluster see [here](https://www.elastic.co/guide/en/elasticsearch/reference/current/settings.html).

KubeDB allows users to provide 4 different configuration files. These are,

1. master-config.yml
2. client-config.yml
3. data-config.yml
4. common-config.yml

**master-config.yml:** Users can provide master node specific configuration in this file. This configuration will be applied to all master nodes.

**client-config.yml:** Users can provide client/ingest node specific configuration in this file. This configuration will be applied to all client nodes.

**data-config.yml:** Users can provide data node specific configuration in this file. This configuration will be applied to all data nodes.

**common-config.yml:** The user can provide a common configuration file. This configuration will be applied to all nodes in the Elasticsearch cluster.

As KubeDB provides two different way to create an Elasticsearch cluster, this custom configuration will be applied according to respective modes. These are:

1. With Topology
2. Without Topology

**With Topology:**
When users create an Elasticsearch cluster using KubeDB specifying `spec.topology` field, KubeDB creates Elasticsearch nodes according to the specification. In this case, the custom configuration will be applied according to node rule. i.e. configuration of `master-config.yml` file will be applied to only master nodes and so on. However, configuration of `common-config.yml` file will be applied to all nodes.

**Without Topology:**
When users create an Elasticsearch cluster using KubeDB without specifying `spec.topology` field, KubeDB creates all nodes as combined node, i.e., every node acts as master, client and data node. In this case, configuration from all configuration files will be merged and applied to all nodes. If the same configuration key appears in more than one configuration file, the value from the highest precedence configuration file will be used in `elasticsearch.yml`. The precedence order of configuration files from lowest to highest is:

1. common-config.yml
2. data-config.yml
3. client-config.yml
4. master-config.yml

`common-config.yml` has the lowest precedence and `master-config.yml` has the highest precedence.

## How to Provide Configuration Files

At first, you have to create configuration files with name specified earlier with your desired configuration. Then you have to put these files into a [volume](https://kubernetes.io/docs/concepts/storage/volumes/). You have to specify this volume in `spec.configSource` section while creating Elasticsearch crd. KubeDB will mount this volume into `/elasticsearch/custom-config` directory of the elasticsearch pod. Configurations from these files will be merged to `elasticsearch.yml` file according to cluster mode described earlier. Finally, Elasticsearch server will use this configuration file.

## Next Steps

- Learn how to use custom configuration specifying topology from [here](/docs/guides/elasticsearch/custom-config/with-topology.md).
- Learn how to use custom configuration without specifying topology from [here](/docs/guides/elasticsearch/custom-config/without-topology.md).