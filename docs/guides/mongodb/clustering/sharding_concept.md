---
title: MongoDB Sharding Concept
menu:
  docs_0.11.0:
    identifier: mg-clustering-sharding-concept
    name: ReplicaSet Concept
    parent: mg-clustering-mongodb
    weight: 20
menu_name: docs_0.11.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# MongoDB Sharding

Sharding is a method for distributing data across multiple machines. MongoDB uses sharding to support deployments with very large data sets and high throughput operations. This section introduces sharding in MongoDB as well as the components and architecture of sharding.


## Sharding Members

A MongoDB sharded cluster consists of the following components:

- shard: Each shard contains a subset of the sharded data. As of MongoDB 3.6, shards must be deployed as a replica set.
- mongos: The mongos acts as a query router, providing an interface between client applications and the sharded cluster.
- config servers: Config servers store metadata and configuration settings for the cluster. As of MongoDB 3.4, config servers must be deployed as a replica set (CSRS).

## Production Configuration

In a production cluster, ensure that data is redundant and that your systems are highly available. Consider the following for a production sharded cluster deployment:

- Deploy Config Servers as a 3 member replica set
- Deploy each Shard as a 3 member replica set
- Deploy one or more mongos routers

### Primary

The primary is the only member in the replica set that receives write operations. MongoDB applies write operations on the primary and then records the operations on the primary’s oplog. Secondary members replicate this log and apply the operations to their data sets.

In the following three-member replica set, the primary accepts all write operations. Then the secondaries replicate the oplog to apply to their data sets.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/mongodb/primary-nodes.png" width="500" height="408">
</p>

All members of the replica set can accept read operations. However, by default, an application directs its read operations to the primary member.

The replica set can have at most one primary. If the current primary becomes unavailable, an election determines the new primary. See [Replica Set Elections](https://docs.mongodb.com/manual/core/replica-set-elections/) for more details.

### Secondaries

A secondary maintains a copy of the primary’s data set. To replicate data, a secondary applies operations from the primary’s oplog to its own data set in an asynchronous process. A replica set can have one or more secondaries.

The following three-member replica set has two secondary members. The secondaries replicate the primary’s oplog and apply the operations to their data sets.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/mongodb/secondary-nodes.png" width="530" height="200">
</p>

Although clients cannot write data to secondaries, clients can read data from secondary members. See [Read Preference](https://docs.mongodb.com/manual/core/read-preference/) for more information on how clients direct read operations to replica sets.

A secondary can become a primary. If the current primary becomes unavailable, the replica set holds an election to choose which of the secondaries becomes the new primary.

### Arbiter

An arbiter does not have a copy of data set and cannot become a primary. Replica sets may have arbiters to add a vote in elections for primary. Arbiters always have exactly 1 election vote, and thus allow replica sets to have an uneven number of voting members without the overhead of an additional member that replicates data.

Changed in version 3.6: Starting in MongoDB 3.6, arbiters have priority 0. When you upgrade a replica set to MongoDB 3.6, if the existing configuration has an arbiter with priority 1, MongoDB 3.6 reconfigures the arbiter to have priority 0.

> IMPORTANT: Do not run an arbiter on systems that also host the primary or the secondary members of the replica set. [[reference]](https://docs.mongodb.com/manual/core/replica-set-members/#arbiter).

## Asynchronous Replication

Secondaries apply operations from the primary asynchronously. By applying operations after the primary, sets can continue to function despite the failure of one or more members.

## Automatic Failover

When a primary does not communicate with the other members of the set for more than the configured electionTimeoutMillis period (10 seconds by default), an eligible secondary calls for an election to nominate itself as the new primary. The cluster attempts to complete the election of a new primary and resume normal operations.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/mongodb/mongodb-election.png" width="500" height="378">
</p>

The replica set cannot process write operations until the election completes successfully. The replica set can continue to serve read queries if such queries are configured to run on secondaries while the primary is offline.

The median time before a cluster elects a new primary should not typically exceed 12 seconds, assuming default [replica configuration settings](https://docs.mongodb.com/manual/reference/replica-configuration/#rsconf.settings). This includes time required to mark the primary as unavailable and call and complete an election. You can tune this time period by modifying the [settings.electionTimeoutMillis](https://docs.mongodb.com/manual/reference/replica-configuration/#rsconf.settings.electionTimeoutMillis) replication configuration option. Factors such as network latency may extend the time required for replica set elections to complete, which in turn affects the amount of time your cluster may operate without a primary. These factors are dependent on your particular cluster architecture.

Lowering the electionTimeoutMillis replication configuration option from the default 10000 (10 seconds) can result in faster detection of primary failure. However, the cluster may call elections more frequently due to factors such as temporary network latency even if the primary is otherwise healthy. This can result in increased [rollbacks](https://docs.mongodb.com/manual/core/replica-set-rollbacks/#replica-set-rollback) for [w : 1](https://docs.mongodb.com/manual/reference/write-concern/#wc-w) write operations.

Your application connection logic should include tolerance for automatic failovers and the subsequent elections.

## Read Operations

By default, clients read from the primary; however, clients can specify a read preference to send read operations to secondaries. Asynchronous replication to secondaries means that reads from secondaries may return data that does not reflect the state of the data on the primary. For information on reading from replica sets, see [Read Preference](https://docs.mongodb.com/manual/core/read-preference/).

[Multi-document transactions](https://docs.mongodb.com/manual/core/transactions/) that contain read operations must use read preference primary.

All operations in a given transaction must route to the same member.

## Transactions

Starting in MongoDB 4.0, multi-document transactions are available for replica sets.

[Multi-document transactions](https://docs.mongodb.com/manual/core/transactions/) that contain read operations must use read preference primary.

All operations in a given transaction must route to the same member.

## Change Streams

Starting in MongoDB 3.6, [change streams](https://docs.mongodb.com/manual/changeStreams/) are available for replica sets and sharded clusters. Change streams allow applications to access real-time data changes without the complexity and risk of tailing the oplog. Applications can use change streams to subscribe to all data changes on a collection or collections.

## Next Steps

- [Deploy MongoDB ReplicaSet](/docs/guides/mongodb/clustering/replicaset.md) using KubeDB.
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Detail concepts of [MongoDBVersion object](/docs/concepts/catalog/mongodb.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).

NB: The images in this page are taken from [MongoDB website](https://docs.mongodb.com/manual/replication/).
