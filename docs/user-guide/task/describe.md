# Describe database

Lets describe a Postgres database `postgres-demo`.

### kubedb describe

Following command will describe Postgres database object `postgres-demo` with relevant information.

```bash
$ kubedb describe pg postgres-demo

Name:		postgres-demo
Namespace:	default
StartTimestamp:	Mon, 05 Jun 2017 10:10:06 +0600
Status:		Running
No volumes.

StatefulSet:
  Name:			postgres-demo-pg
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Mon, 05 Jun 2017 10:10:14 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		postgres-demo
  Type:		ClusterIP
  IP:		10.0.0.36
  Port:		port	5432/TCP

Database Secret:
  Name:	postgres-demo-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	35 bytes

Snapshots:
  Name                     Bucket          StartTime                         CompletionTime                    Phase
  ----                     ------          ---------                         --------------                    -----
  postgres-demo-20170605-073557   database-test   Mon, 05 Jun 2017 13:35:57 +0600   Mon, 05 Jun 2017 13:36:10 +0600   Succeeded
  snapshot-20170505-1147          database-test   Mon, 05 Jun 2017 11:48:06 +0600   Mon, 05 Jun 2017 12:01:39 +0600   Succeeded

Events:
  FirstSeen   LastSeen   Count     From                  Type       Reason               Message
  ---------   --------   -----     ----                  --------   ------               -------
  3m          3m         1         Snapshot Controller   Normal     Starting             Backup running
  21m         21m        1         Postgres operator     Normal     SuccessfulCreate     Successfully created StatefulSet
  21m         21m        1         Postgres operator     Normal     SuccessfulCreate     Successfully created Postgres
  29m         29m        1         Postgres operator     Normal     SuccessfulValidate   Successfully validate Postgres
  29m         29m        1         Postgres operator     Normal     Creating             Creating Kubernetes objects
```

The `kubedb describe` command provides following basic information of TPR object.
* StatefulSet
* Storage (Persistent Volume)
* Service
* Secret (If available)
* Snapshots (If any)
* Monitoring system (If available)

This command also shows events unless `--show-events=false`

To describe all Postgres objects in `default` namespace, use following command
```bash
$ kubedb describe pg
```

If we want to describe all Postgres objects from every namespace, we need to provide flag `--all-namespaces`.

```bash
$ kubedb describe pg --all-namespaces
```

And to describe all objects of all resource type from every namespace, use following

```bash
$ kubedb describe all --all-namespaces
```

This `all` represents all supported TPR.

We can also describe some selected objects

```bash
$ kubedb describe pg,es --all-namespaces --selector='group=dev'
```

This command will describe all Elastic & Postgres object with specified labels from every namespace.

##### Click [here](../reference/describe.md) to get command details.
