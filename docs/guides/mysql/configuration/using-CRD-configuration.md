---
title: Run MySQL with Custom Configuration
menu:
  docs_0.9.0:
    identifier: my-crd-configuration
    name: Using CRD Config
    parent: my-custom-config
    weight: 15
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Run MySQL with Custom Configuration

KubeDB supports providing custom configuration for MySQL via [PodTemplate](/docs/concepts/databases/mysql.md#specpodtemplate). This tutorial will show you how to use KubeDB to run a MySQL database with custom configuration using PodTemplate.

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

  ```console
  $ kubectl create ns demo
  namespace/demo created
  ```

> Note: YAML files used in this tutorial are stored in [docs/examples/mysql](https://github.com/kubedb/cli/tree/master/docs/examples/mysql) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Overview

KubeDB allows providing a template for database pod through `spec.podTemplate`. KubeDB operator will pass the information provided in `spec.podTemplate` to the StatefulSet created for MySQL database.

KubeDB accept following fields to set in `spec.podTemplate:`

- metadata:
  - annotations (pod's annotation)
- controller:
  - annotations (statefulset's annotation)
- spec:
  - env
  - resources
  - initContainers
  - imagePullSecrets
  - nodeSelector
  - affinity
  - schedulerName
  - tolerations
  - priorityClassName
  - priority
  - securityContext

Read about the fields in details in [PodTemplate concept](/docs/concepts/databases/mysql.md#specpodtemplate),

## CRD Configuration

Below is the YAML for the MySQL created in this example. Here, [`spec.podTemplate.spec.env`](/docs/concepts/databases/mysql.md#specpodtemplatespecenv) specifies environment variables and [`spec.podTemplate.spec.args`](/docs/concepts/databases/mysql.md#specpodtemplatespecargs) provides extra arguments for [MySQL Docker Image](https://hub.docker.com/_/mysql/). 

In this tutorial, an initial database `myDB` will be created by providing `env` `MYSQL_DATABASE` while the server character set will be set to `utf8mb4` by adding extra `args`. Note that, `character-set-server` in `MySQL 5.7` is `latin1`.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-misc-config
  namespace: demo
spec:
  version: "5.7-v1"
  storageType: "Durable"
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  podTemplate:
    spec:
      env:
      - name: MYSQL_DATABASE
        value: myDB
      args:
      - --character-set-server=utf8mb4
      resources:
        requests:
          memory: "1Gi"
          cpu: "250m"
  terminationPolicy: Pause
  updateStrategy:
    type: RollingUpdate
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/configuration/mysql-misc-config.yaml
mysql.kubedb.com/mysql-misc-config created
```

Now, wait a few minutes. KubeDB operator will create necessary PVC, statefulset, services, secret etc. If everything goes well, we will see that a pod with the name `mysql-misc-config-0` has been created.

Check that the statefulset's pod is running

```console
$ kubectl get pod -n demo
NAME                  READY     STATUS    RESTARTS   AGE
mysql-misc-config-0   1/1       Running   0          1m
```

Check the pod's log to see if the database is ready

```console
$ kubectl logs -f -n demo mysql-misc-config-0
Initializing database
.....
Database initialized
Initializing certificates
...
Certificates initialized
MySQL init process in progress...
....
MySQL init process done. Ready for start up.
....
2018-10-02T09:34:33.694994Z 0 [Note] mysqld: ready for connections.
Version: '5.7.23'  socket: '/var/run/mysqld/mysqld.sock'  port: 3306  MySQL Community Server (GPL)
....
```

Once we see `[Note] /usr/sbin/mysqld: ready for connections.` in the log, the database is ready.

Now, we will check if the database has started with the custom configuration we have provided.

First, deploy [phpMyAdmin](https://hub.docker.com/r/phpmyadmin/phpmyadmin/) to connect with the MySQL database we have just created.

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/quickstart/demo-1.yaml
deployment.extensions/myadmin created
service/myadmin created
```

Then, open your browser and go to the following URL: _http://{cluster-ip}:{myadmin-svc-nodeport}_. For minikube you can get this URL by running the following command:

```console
$ minikube service myadmin -n demo --url
http://192.168.99.100:30039
```

Now, let's connect to the database from the phpMyAdmin dashboard using the database pod IP and MySQL user password.

```console
$ kubectl get pods mysql-misc-config-0 -n demo -o yaml | grep IP
  hostIP: 10.0.2.15
  podIP: 172.17.0.6

$ kubectl get secrets -n demo mysql-misc-config-auth -o jsonpath='{.data.\user}' | base64 -d
root

$ kubectl get secrets -n demo mysql-misc-config-auth -o jsonpath='{.data.\password}' | base64 -d
MLO5_fPVKcqPiEu9
```

Once, you have connected to the database with phpMyAdmin go to **SQL** tab and run sql to see all databases `SHOW DATABASES;` and to see charcter-set configuration `SHOW VARIABLES LIKE 'char%';`. You will see a database called `myDB` is created and also all the character-set is set to `utf8mb4`.

![mysql_all_databases](/docs/images/mysql/mysql-all-databases.png)

![mysql_charset](/docs/images/mysql/mysql-charset.png)

## Snapshot Configuration

`Snapshot` also has the scope to be configured through `spec.podTemplate`. In this tutorial, an extra argument is passed to snapshot crd so that the backup job uses `--default-character-set=utf8mb4` while taking backup.

Below is the Snapshot CRD that is deployed in this tutorial. Create a secret `my-snap-secret` from [here](/docs/guides/mysql/snapshot/backup-and-restore.md#instant-backups) for snapshot. 

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: snap-mysql-config
  namespace: demo
  labels:
    kubedb.com/kind: MySQL
spec:
  databaseName: mysql-misc-config
  storageSecretName: my-snap-secret
  gcs:
    bucket: kubedb-qa
  podTemplate:
    spec:
      args:
      - --all-databases
      - --default-character-set=utf8mb4
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/configuration/snapshot-misc-conf.yaml 
snapshot.kubedb.com/snap-mysql-config created


$ kubedb get snap -n demo
NAME                DATABASENAME        STATUS      AGE
snap-mysql-config   mysql-misc-config   Succeeded   1m
```

## Scheduled Backups

To configure BackupScheduler, add the require changes in PodTemplate just like snapshot object.

```yaml
$ kubedb edit my mysql-misc-config -n demo
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-misc-config
  namespace: demo
  ...
spec:
  backupSchedule:
    cronExpression: '@every 1m'
    storageSecretName: my-snap-secret
    gcs:
      bucket: kubedb-qa
    podTemplate:
      controller: {}
      metadata: {}
      spec:
        args:
        - --all-databases
        - --default-character-set=utf8mb4
        resources: {}
  ...
status:
  observedGeneration: 3$4212299729528774793
  phase: Running
```

```console
$ kubedb get snap -n demo
NAME                                DATABASENAME        STATUS      AGE
mysql-misc-config-20181002-105247   mysql-misc-config   Succeeded   3m
mysql-misc-config-20181002-105349   mysql-misc-config   Succeeded   2m
mysql-misc-config-20181002-105449   mysql-misc-config   Succeeded   1m
mysql-misc-config-20181002-105549   mysql-misc-config   Succeeded   43s
snap-mysql-config                   mysql-misc-config   Succeeded   12m
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo my/mysql-misc-config -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo my/mysql-misc-config

kubectl patch -n demo drmn/mysql-misc-config -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mysql-misc-config

kubectl delete -n demo configmap my-custom-config
kubectl delete deployment -n demo myadmin
kubectl delete service -n demo myadmin

kubectl delete ns demo
```

If you would like to uninstall KubeDB operator, please follow the steps [here](/docs/setup/uninstall.md).

## Next Steps

- [Quickstart MySQL](/docs/guides/mysql/quickstart/quickstart.md) with KubeDB Operator.
- [Snapshot and Restore](/docs/guides/mysql/snapshot/backup-and-restore.md) process of MySQL databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mysql/snapshot/scheduled-backup.md) of MySQL databases using KubeDB.
- Initialize [MySQL with Script](/docs/guides/mysql/initialization/using-script.md).
- Initialize [MySQL with Snapshot](/docs/guides/mysql/initialization/using-snapshot.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mysql/private-registry/using-private-registry.md) to deploy MySQL with KubeDB.
- Use [kubedb cli](/docs/guides/mysql/cli/cli.md) to manage databases like kubectl for Kubernetes.
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
