---
title: Upgrading MySQL
menu:
  docs_0.10.0:
    identifier: my-upgrade-manul
    name: Manual
    parent: my-upgrading-mysql
    weight: 10
menu_name: docs_0.10.0
section_menu_id: guides
---

# Upgrading MySQL

If you want to upgrade your existing MySQL databases, you can follow the procedure below.
>**NOTE**: Upgrading may cause some downtime for your database(s) (normally not more than a few minutes).

## Check your current version(s)

To get an overview of the MySQL versions currently running in your cluster, run the following command with the KubeDB CLI:

```console
$ kubedb get mysql --all-namespaces
NAMESPACE   NAME            VERSION   STATUS         AGE
ns1         database-ns1    5.7-v1    Running        59d
ns2         database-ns2    5.7-v1    Running        38d
```

## Make a backup!

Before you start the upgrade, it's highly recommended that you make a backup. In case something goes wrong, you can [recreate a database from the snapshot](https://kubedb.com/docs/0.10.0/guides/mysql/initialization/using-snapshot/). You have two options here:

- [Instant Backup](https://kubedb.com/docs/0.10.0/guides/mysql/snapshot/backup-and-restore/)
- [Scheduled Backup](https://kubedb.com/docs/0.10.0/guides/mysql/snapshot/scheduled-backup/)

## Check MySQL upgrade paths

You should check the supported MySQL [upgrade paths](https://dev.mysql.com/doc/refman/8.0/en/upgrade-paths.html) to find out to which version you can upgrade. For example, upgrade from 5.7.9+ to 8.0 is supported, but **5.6 to 8.0 is not**.

### PHP Compatibility Note for MySQL 8.0+

At the moment of writing, PHP [doesn't support](https://secure.php.net/manual/en/mysqli.requirements.php) the new `caching_sha2_password` authentication method on MySQL 8+ (tested on PHP 7.2 and 7.3). You will get errors like:

```
mysqli_real_connect(): The server requested authentication method unknown to the client [caching_sha2_password]

mysqli_real_connect(): (HY000/2054): The server requested authentication method unknown to the client
```

There are ways to work around this and fall back to `mysql_native_password`, but for now we strongly recommend not to upgrade to MySQL 8.0 if you're using PHP. This way, you will use the safer `caching_sha2_password` method as soon as PHP adds support.

### Upgrade from 8.0.3

>**NOTE**: Upgrading from version 8.0.3 (RC) to 8.0.14 (GA) is not supported. Following this procedure might be painful, and could raise unexpected errors. Proceed at your own risk!

8.0.3 (Release Candidate) has been offered by KubeDB. After you try to upgrade to 8.0.14+, when the pod starts, it will go into CrashLoopBackOff. You'll get the following errors:

```console
2019-03-13T21:14:44.908546Z 0 [System] [MY-010116] [Server] /usr/sbin/mysqld (mysqld 8.0.14) starting as process 1
2019-03-13T21:14:46.033145Z 1 [ERROR] [MY-013326] [Server] Upgrading the server from server version '0' is not supported.
2019-03-13T21:14:46.033156Z 1 [ERROR] [MY-010334] [Server] Failed to initialize DD Storage Engine
2019-03-13T21:14:46.033401Z 0 [ERROR] [MY-010020] [Server] Data Dictionary initialization failed.
2019-03-13T21:14:46.033585Z 0 [ERROR] [MY-010119] [Server] Aborting
2019-03-13T21:14:47.656766Z 0 [System] [MY-010910] [Server] /usr/sbin/mysqld: Shutdown complete (mysqld 8.0.14)  MySQL Community Server - GPL.
```

The reason for this is that 8.0.3 is a [non-GA release](https://mysqlserverteam.com/upgrading-to-mysql-8-0-here-is-what-you-need-to-know/): *"As of MySQL 8.0.11, the server version is written to the data dictionary tablespace." ... "The immediate consequence for MySQL 8.0.11 is that it refuses to upgrade from MySQL 8.0 DMRs (8.0.1, 8.0.1, 8.0.2, 8.0.3, 8.0.4) to MySQL 8.0.11 because this path is not safe."*

The best option you have in this case is to:

1. Export only the tables you need (**not** the built-in MySQL tables) with [mysqldump](https://dev.mysql.com/doc/refman/8.0/en/mysqldump-sql-format.html). Export the users you need as well (note: exporting the root user will overwrite its password on your new installation!).
2. Create a completely new MySQL database with the [regular KubeDB procedure](https://kubedb.com/docs/0.10.0/guides/mysql/initialization/using-script/).
3. Restore the tables & users you exported. You might run into issues like `ERROR 3723 (HY000) at line 821: The table 'events' may not be created in the reserved tablespace 'mysql'`. You will have to fix these manually in the dumpfile, as this is a non-supported upgrade path.
4. Verify if everything is working properly.
5. Delete the old 8.0.3 database.

## Perform the upgrade

The upgrade consists of two steps:

1. Upgrade container version in Kubernetes/KubeDB
2. Upgrade MySQL data with the `mysql_upgrade` command

### 1. Upgrade container version in Kubernetes/KubeDB
First, we'll edit the kubedb instance. Open the instance editor:

```console
kubedb edit mysql/DB_NAME -n YOUR_NAMESPACE
```

Replace `version: "X.X.X"` with the version you need, for example 8.0.14:

```yaml
  version: "8.0.14"
```

After you've edited the YAML spec, save the file with the `:wq` command. The deployment of the new version will now start:

```console
mysql.kubedb.com/database edited
```

Check if the version of the database has changed:

```console
$ kubedb get mysql/DB_NAME -n YOUR_NAMESPACE
NAMESPACE   NAME            VERSION   STATUS         AGE
ns1         database-ns1    8.0.14   Running        59d
```

To double-check, open the Kubernetes dashboard or check with `kubectl`:

```console
$ kubectl get pod -n YOUR_NAMESPACE POD_NAME
NAME            READY    STATUS    RESTARTS   AGE
database-ns1    1/1      Running   0          3m41s
```

### 2. Upgrade MySQL data with the `mysql_upgrade` command

>**NOTE**: This step is only necessary if you are currently running or upgrading to **MySQL 8.0.15 or lower**. Starting from 8.0.16, this step is performed automatically by the MySQL server on startup. For more details, consult the [MySQL documentation](https://dev.mysql.com/doc/refman/8.0/en/mysql-upgrade.html).

In order to perform the upgrade, we need to run the `mysql_upgrade` command within the database pod.

```console
kubectl get pod -n YOUR_NAMESPACE
NAME                                                        READY     STATUS      RESTARTS   AGE
database-0                                                  1/1       Running     0          1h
```

1. Get the MySQL root password: `kubectl get secret DB_NAME-auth -n YOUR_NAMESPACE -o 'go-template={{index .data "password"}}' | base64 --decode`
2. Log into the pod with the pod name from above: `kubectl exec -it POD_NAME -n YOUR_NAMESPACE -- /bin/bash`
3. Run `mysql_upgrade -u root -p`, press Enter, and enter the root password from step 1.
4. Done!
