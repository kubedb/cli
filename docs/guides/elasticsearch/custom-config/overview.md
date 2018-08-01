---
title: Run Elasticsearch with Custom Configuration
menu:
  docs_0.8.0:
    identifier: es-custom-config-overview
    name: Overview
    parent: es-custom-config
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Custom Configuration File

KubeDB supports providing custom configuration files for Elasticsearch cluster. This tutorial will give you an overview of how custom configuration files for Elasticsearch cluster are managed in KubeDB.

## Overview

Elasticsearch use configuration from `elasticsearch.yaml` file. In KubeDB this file is located in `/elasticsearch/config` directory of elasticsearch pods. To know more about configuring the Elasticsearch cluster see [here](https://www.elastic.co/guide/en/elasticsearch/reference/current/settings.html).

KubeDB allows the users to provide 4 different configuration files. These are,

1. master-config.yaml
2. client-config.yaml
3. data-config.yaml
4. common-config.yaml

**master-config.yaml:** The users can provide master node specific configurations in this file. These configurations will be applied to all master nodes.

**client-config.yaml:** The users can provide client/ingest node specific configurations in this file. These configurations will be applied to all client nodes.

**data-config.yaml:** The users can provide data node specific configurations in this file. These configurations will be applied to all data nodes.

**common-config.yaml:** The user can provide a common configuration file. These configurations will be applied to all nodes in the Elasticsearch cluster.

As KubeDB provides two different way to create Elasticsearch cluster, this custom configurations will be applied according to respective modes. These are,

1. With Topology
2. Without Topology

**With Topology:**
When the users create Elasticsearch cluster using KubeDB specifying `spec.topology` field, KubeDB creates Elasticsearch nodes according to the specification. In this case, the custom configurations will be applied according to node rule. i.e. configurations of `master-config.yaml` file will be applied to only master nodes and so on. However, configurations of `common-config.yaml` file will be applied to all nodes.

**Without Topology:**
When the users create Elasticsearch cluster using KubeDB without specifying `spec.topology` field, KubeDB creates all nodes as combined node. All nodes act as master, client and data node. In this case, configurations from all configuration files will be applied to all nodes. However, there is a precedence of the configuration files. If more than one configuration file has same configuration field, this configuration will be applied from the highest precedented configuration file. The precedence order is,

1. master-config.yaml
2. client-config.yaml
3. data-config.yaml
4. common-config.yaml

`master-config.yaml` has the highest precedence and `common-config.yaml` has the lowest precedence.

## How to Provide Configuration Files

At first, you have to create configuration files with named specified earlier with your desired configurations. Then you have to put these files into a [volume](https://kubernetes.io/docs/concepts/storage/volumes/). You have to specify this volume in `spec.configSource` section while creating Elasticsearch crd. KubeDB will mount this volume into `/elasticsearch/custom-config` directory of the elasticsearch pod. Configurations from these files will be appended to `elasticsearch.yaml` file according to cluster mode described earlier. Finally, Elasticsearch binary will use this configuration file.

## Next Steps

- [Learn how to use custom configuration specifying topology](/docs/guides/elasticsearch/custom-config/with-topology.md).
- [Learn how to use custom configuration without specifying topology](/docs/guides/elasticsearch/custom-config/without-topology.md).