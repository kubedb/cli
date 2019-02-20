---
title: Setup HA Postgres Cluster
menu:
  docs_0.9.0:
    identifier: pg-ha-cluster-clustering
    name: HA Setup
    parent: pg-clustering-postgres
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Configuring Highly Available PostgreSQL Cluster

In PostgreSQL, multiple servers can work together to serve high availability and load balancing. These servers will be either in *Master* or *Standby* mode.

In *master* mode, server that can modify data. In *standby* mode, the server continuously applies WAL received from the master server. The standby server can read WAL from a WAL archive (see restore_command) or directly from the master over a TCP connection (streaming replication).

Standby servers can be either *warm standby* or *hot standby* server.

## Warm Standby

A standby server that cannot be connected to until it is promoted to a *master* server is called a *warm standby* server.
*Standby* servers are by default *warm standby* unless we make them *hot standby*.

The following is an example of a `Postgres` object which creates PostgreSQL cluster of three servers.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: warm-postgres
  namespace: demo
spec:
  version: "9.6-v2"
  replicas: 3
  standbyMode: Warm
  storageType: Ephemeral
```

In this examples:

- This `Postgres` object creates three PostgreSQL servers, indicated by the **`replicas`** field.
- One server will be *primary* and two others will be *warm standby* servers, as instructed by **`spec.standbyMode`**

## Hot Standby

A standby server that can accept connections and serves read-only queries is called a *hot standby* server.

The following `Postgres` object will create PostgreSQL cluster with *hot standby* servers.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: hot-postgres
  namespace: demo
spec:
  version: "9.6-v2"
  replicas: 3
  standbyMode: Hot
  storageType: Ephemeral
```

In this examples:

- This `Postgres` object creates three PostgreSQL servers, indicated by the **`replicas`** field.
- One server will be *primary* and two others will be *hot standby* servers, as instructed by **`spec.standbyMode`**

## High Availability

Database servers can work together to allow a second server to take over quickly if the *primary* server fails. This is called high availability. When *primary* server is unavailable, *standby* servers go through a leader election process to take control as *primary* server. PostgreSQL database with high availability feature can either have *warm standby* or *hot standby* servers.

To enable high availability, you need to create PostgreSQL with multiple server. Set `spec.replicas` to more than one in Postgres.

[//]: # (For more information on failover process, [read here])

## Load Balancing

*Master* server along with *standby* server(s) can serve the same data. This is called load balancing. In our setup, we only support read-only *standby* server.
To enable load balancing, you need to setup *hot standby* PostgreSQL cluster.

Read about [hot standby](#hot-standby) and its setup in Postgres.

## Replication

There are many approaches available to scale PostgreSQL beyond running on a single server.

Now KubeDB supports only following one:

- **Streaming Replication** provides *asynchronous* replication to one or more *standby* servers.

These *standby* servers can also be *hot standby* server. This is the fastest type of replication available as
WAL data is sent immediately rather than waiting for a whole segment to be produced and shipped.

    KubeDB PostgreSQL support [Streaming Replication](/docs/guides/postgres/clustering/streaming_replication.md)

## Next Steps

- Learn how to setup [Streaming Replication](/docs/guides/postgres/clustering/streaming_replication.md)
