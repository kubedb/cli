---
title: Run MySQL using Private Registry
menu:
  docs_0.9.0:
    identifier: my-using-private-registry-private-registry
    name: Quickstart
    parent: my-private-registry-mysql
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Deploy MySQL from private Docker registry

KubeDB operator supports using private Docker registry. This tutorial will show you how to use KubeDB to run MySQL database using private Docker images.

## Before You Begin

- Read [concept of MySQL Version Catalog](/docs/concepts/catalog/mysql.md) to learn detail concepts of `MySQLVersion` object.

- You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- You will also need a docker private [registry](https://docs.docker.com/registry/) or [private repository](https://docs.docker.com/docker-hub/repos/#private-repositories).  In this tutorial we will use private repository of [docker hub](https://hub.docker.com/).

- You have to push the required images from KubeDB's [Docker hub account](https://hub.docker.com/r/kubedb/) into your private registry. For mysql, push `DB_IMAGE`, `TOOLS_IMAGE`, `EXPORTER_IMAGE` of following MySQLVersions, where `deprecated` is not true, to your private registry.

  ```console
  $ kubectl get mysqlversions -n kube-system  -o=custom-columns=NAME:.metadata.name,VERSION:.spec.version,DB_IMAGE:.spec.db.image,TOOLS_IMAGE:.spec.tools.image,EXPORTER_IMAGE:.spec.exporter.image,DEPRECATED:.spec.deprecated
  NAME     VERSION   DB_IMAGE              TOOLS_IMAGE                 EXPORTER_IMAGE                   DEPRECATED
  5        5         kubedb/mysql:5        kubedb/mysql-tools:5        kubedb/operator:0.8.0            true
  5-v1     5         kubedb/mysql:5-v1     kubedb/mysql-tools:5-v2     kubedb/mysqld-exporter:v0.11.0   true
  5.7      5.7       kubedb/mysql:5.7      kubedb/mysql-tools:5.7      kubedb/operator:0.8.0            true
  5.7-v1   5.7       kubedb/mysql:5.7-v1   kubedb/mysql-tools:5.7-v2   kubedb/mysqld-exporter:v0.11.0   <none>
  8        8         kubedb/mysql:8        kubedb/mysql-tools:8        kubedb/operator:0.8.0            true
  8-v1     8         kubedb/mysql:8-v1     kubedb/mysql-tools:8-v2     kubedb/mysqld-exporter:v0.11.0   true
  8.0      8.0       kubedb/mysql:8.0      kubedb/mysql-tools:8.0      kubedb/operator:0.8.0            true
  8.0-v1   8.0       kubedb/mysql:8.0-v1   kubedb/mysql-tools:8.0-v2   kubedb/mysqld-exporter:v0.11.0   <none>
  8.0-v2   8.0       kubedb/mysql:8.0-v2   kubedb/mysql-tools:8.0-v3   kubedb/mysqld-exporter:v0.11.0   <none>
  8.0.14   8.0.14    kubedb/mysql:8.0.14   kubedb/mysql-tools:8.0.14   kubedb/mysqld-exporter:v0.11.0   <none>
  8.0.3    8.0.3     kubedb/mysql:8.0.3    kubedb/mysql-tools:8.0.3    kubedb/mysqld-exporter:v0.11.0   <none>
  ```

  Docker hub repositories:

  - [kubedb/operator](https://hub.docker.com/r/kubedb/operator)
  - [kubedb/mysql](https://hub.docker.com/r/kubedb/mysql)
  - [kubedb/mysql-tools](https://hub.docker.com/r/kubedb/mysql-tools)
  - [kubedb/mysqld-exporter](https://hub.docker.com/r/kubedb/mysqld-exporter)

- Update KubeDB catalog for private Docker registry. Ex:

  ```yaml
  apiVersion: catalog.kubedb.com/v1alpha1
  kind: MySQLVersion
  metadata:
    name: "8.0-v2"
    labels:
      app: kubedb
  spec:
    version: "8.0"
    db:
      image: "PRIVATE_DOCKER_REGISTRY/mysql:8.0-v2"
    exporter:
      image: "PRIVATE_DOCKER_REGISTRY/mysqld-exporter:v0.11.0"
    tools:
      image: "PRIVATE_DOCKER_REGISTRY/mysql-tools:8.0-v3"
  
  ```

- To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

  ```console
  $ kubectl create ns demo
  namespace/demo created
   ```

## Create ImagePullSecret

ImagePullSecrets is a type of a Kubernete Secret whose sole purpose is to pull private images from a Docker registry. It allows you to specify the url of the docker registry, credentials for logging in and the image name of your private docker image.

Run the following command, substituting the appropriate uppercase values to create an image pull secret for your private Docker registry:

```console
$ kubectl create secret docker-registry -n demo myregistrykey \
  --docker-server=DOCKER_REGISTRY_SERVER \
  --docker-username=DOCKER_USER \
  --docker-email=DOCKER_EMAIL \
  --docker-password=DOCKER_PASSWORD
secret/myregistrykey created
```

If you wish to follow other ways to pull private images see [official docs](https://kubernetes.io/docs/concepts/containers/images/) of kubernetes.

NB: If you are using `kubectl` 1.9.0, update to 1.9.1 or later to avoid this [issue](https://github.com/kubernetes/kubernetes/issues/57427).

## Install KubeDB operator

When installing KubeDB operator, set the flags `--docker-registry` and `--image-pull-secret` to appropriate value. Follow the steps to [install KubeDB operator](/docs/setup/install.md) properly in cluster so that to points to the DOCKER_REGISTRY you wish to pull images from.

## Deploy MySQL database from Private Registry

While deploying `MySQL` from private repository, you have to add `myregistrykey` secret in `MySQL` `spec.imagePullSecrets`.
Below is the MySQL CRD object we will create.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MySQL
metadata:
  name: mysql-pvt-reg
  namespace: demo
spec:
  version: "8.0-v2"
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
  podTemplate:
    spec:
      imagePullSecrets:
      - name: myregistrykey
```

Now run the command to deploy this `MySQL` object:

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mysql/private-registry/demo-2.yaml
mysql.kubedb.com/mysql-pvt-reg created
```

To check if the images pulled successfully from the repository, see if the `MySQL` is in running state:

```console
$ kubectl get pods -n demo
NAME              READY     STATUS    RESTARTS   AGE
mysql-pvt-reg-0   1/1       Running   0          56s
```

## Snapshot

You can specify `imagePullSecret` for Snapshot objects in `spec.podTemplate.spec.imagePullSecrets` field of Snapshot object. If you are using scheduled backup, you can also provide `imagePullSecret` in `backupSchedule.podTemplate.spec.imagePullSecrets` field of MySQL crd. KubeDB also reuses `imagePullSecret` for Snapshot object from `spec.podTemplate.spec.imagePullSecrets` field of MySQL crd.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mysql/mysql-pvt-reg -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mysql/mysql-pvt-reg

kubectl patch -n demo drmn/mysql-pvt-reg -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mysql-pvt-reg

kubectl delete ns demo
```

## Next Steps

- [Snapshot and Restore](/docs/guides/mysql/snapshot/backup-and-restore.md) process of MySQL databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mysql/snapshot/scheduled-backup.md) of MySQL databases using KubeDB.
- Initialize [MySQL with Script](/docs/guides/mysql/initialization/using-script.md).
- Initialize [MySQL with Snapshot](/docs/guides/mysql/initialization/using-snapshot.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mysql/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MySQL database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mysql/monitoring/using-builtin-prometheus.md).
- Detail concepts of [MySQL object](/docs/concepts/databases/mysql.md).
- Detail concepts of [MySQLVersion object](/docs/concepts/catalog/mysql.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
