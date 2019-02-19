---
title: Run MySQL with Custom Configuration
menu:
  docs_0.9.0:
    identifier: my-custom-config-file
    name: Using Config File
    parent: my-custom-config
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Custom Configuration File

KubeDB supports providing custom configuration for MySQL. This tutorial will show you how to use KubeDB to run a MySQL database with custom configuration.

## Before You Begin

- At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

- Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

  ```console
  $ kubectl create ns demo
  namespace/demo created
  
  $ kubectl get ns demo
  NAME    STATUS  AGE
  demo    Active  5s
  ```

> Note: YAML files used in this tutorial are stored in [docs/examples/mysql](https://github.com/kubedb/cli/tree/master/docs/examples/mysql) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Overview

MySQL allows to configure database via configuration file. The default configuration for MySQL can be found in `/etc/mysql/my.cnf` file. When MySQL starts, it will look for custom configuration file in `/etc/mysql/conf.d` directory. If configuration file exist, MySQL instance will use combined startup setting from both `/etc/mysql/my.cnf` and `*.cnf` files in `/etc/mysql/conf.d` directory. This custom configuration will overwrite the existing default one. To know more about configuring MySQL see [here](https://dev.mysql.com/doc/refman/8.0/en/server-configuration.html).

At first, you have to create a config file with `.cnf` extension with your desired configuration. Then you have to put this file into a [volume](https://kubernetes.io/docs/concepts/storage/volumes/). You have to specify this volume  in `spec.configSource` section while creating MySQL crd. KubeDB will mount this volume into `/etc/mysql/conf.d` directory of the database pod.

In this tutorial, we will configure [max_connections](https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_max_connections) and [read_buffer_size](https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_read_buffer_size) via a custom config file. We will use configMap as volume source.

## Custom Configuration

At first, let's create `my-config.cnf` file setting `max_connections` and `read_buffer_size` parameters.

```console
cat <<EOF > my-config.cnf
[mysqld]
max_connections = 200
read_buffer_size = 1048576
EOF

$ cat my-config.cnf
[mysqld]
max_connections = 200
read_buffer_size = 1048576
```

Here, `read_buffer_size` is set to 1MB in bytes.

Now, create a configMap with this configuration file.

```console
 $ kubectl create configmap -n demo my-custom-config --from-file=./my-config.cnf
configmap/my-custom-config created
```

Verify the config map has the configuration file.

```yaml
$ kubectl get configmap -n demo my-custom-config -o yaml
apiVersion: v1
data:
  my-config.cnf: |
    [mysqld]
    max_connections = 200
    read_buffer_size = 1048576
kind: ConfigMap
metadata:
  name: my-custom-config
  namespace: demo
  ...
```

Now, create MySQL crd specifying `spec.configSource` field.

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/configuration/mysql-custom.yaml
mysql.kubedb.com/custom-mysql created
```

Below is the YAML for the MySQL crd we just created.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: custom-mysql
  namespace: demo
spec:
  version: "8.0-v2"
  configSource:
    configMap:
      name: my-custom-config
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

Now, wait a few minutes. KubeDB operator will create necessary PVC, statefulset, services, secret etc. If everything goes well, we will see that a pod with the name `custom-mysql-0` has been created.

Check that the statefulset's pod is running

```console
$ kubectl get pod -n demo
NAME             READY     STATUS    RESTARTS   AGE
custom-mysql-0   1/1       Running   0          44s
```

Check the pod's log to see if the database is ready

```console
$ kubectl logs -f -n demo custom-mysql-0
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
2018-07-10T06:12:46.957611Z 0 [Note] /usr/sbin/mysqld: ready for connections. Version: '8.0.3-rc-log'  socket: '/var/run/mysqld/mysqld.sock'  port: 3306  MySQL Community Server (GPL)
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
$ kubectl get pods custom-mysql-0 -n demo -o yaml | grep IP
  hostIP: 10.0.2.15
  podIP: 172.17.0.6

$ kubectl get secrets -n demo custom-mysql-auth -o jsonpath='{.data.\user}' | base64 -d
root

$ kubectl get secrets -n demo custom-mysql-auth -o jsonpath='{.data.\password}' | base64 -d
MLO5_fPVKcqPiEu9
```

Once, you have connected to the database with phpMyAdmin go to **Variables** tab and search for `max_connections` and `read_buffer_size`. Here are some screenshot showing those configured variables.
![max_connections](/docs/images/mysql/max_connection.png)

![read_buffer_size](/docs/images/mysql/read_buffer_size.png)

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo my/custom-mysql -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo my/custom-mysql

kubectl patch -n demo drmn/custom-mysql -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/custom-mysql

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
