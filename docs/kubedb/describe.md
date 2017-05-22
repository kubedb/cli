# kubedb describe

## Example

##### Help for describe command

```bash
$ kubedb describe --help

Show details of a specific resource or group of resources. This command joins many API calls together to form a detailed
description of a given resource or group of resources.Valid resource types include:

  * all
  * elastic
  * postgres
  * snapshot
  * dormantdatabase

Examples:
  # Describe a elastic
  kubedb describe elastics elasticsearch-demo

  # Describe a postgres
  kubedb describe pg/postgres-demo

  # Describe all dormantdatabases
  kubedb describe drmn

Options:
      --all-namespaces=false: If present, list the requested object(s) across all namespaces.
      --show-events=true: If true, display events related to the described object.

Usage:
  kubedb describe (TYPE [NAME_PREFIX] | TYPE/NAME) [flags] [options]

Use "kubedb describe options" for a list of global command-line options (applies to all commands).
```

##### Describe

```bash
Name:		postgres-demo
Namespace:	default
StartTimestamp:	Thu, 11 May 2017 15:10:50 +0600
Labels::	k8sdb.com/type=postgres
Status:		Running
Replicas:	1  total
Annotations:	postgres.k8sdb.com/version=canary-db
No volumes.

StatefulSet:
  Name:			k8sdb-postgres-demo
  Replicas:		1 current / 1 desired
  CreationTimestamp:	Thu, 11 May 2017 15:10:50 +0600
  Pods Status:		1 Running / 0 Waiting / 0 Succeeded / 0 Failed

Service:
  Name:		postgres-demo
  Type:		ClusterIP
  IP:		10.0.241.155
  Port:		port	5432/TCP

Database Secret:
  Name:	postgres-demo-admin-auth
  Type:	Opaque
  Data
  ====
  .admin:	35 bytes

No Snapshots.

Events:
  FirstSeen   LastSeen   Count     From                  Type       Reason               Message
  ---------   --------   -----     ----                  --------   ------               -------
  37m         37m        1         Postgres Controller   Normal     SuccessfulCreate     Successfully created Postgres
  38m         38m        1         Postgres Controller   Normal     SuccessfulValidate   Successfully validate Postgres
  38m         38m        1         Postgres Controller   Normal     Creating             Creating Kubernetes objects
```
