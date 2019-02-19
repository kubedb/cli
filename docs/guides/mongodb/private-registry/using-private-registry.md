---
title: Run MongoDB using Private Registry
menu:
  docs_0.9.0:
    identifier: mg-using-private-registry-private-registry
    name: Quickstart
    parent: mg-private-registry-mongodb
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using private Docker registry

KubeDB operator supports using private Docker registry. This tutorial will show you how to use KubeDB to run MongoDB database using private Docker images.

## Before You Begin

- Read [concept of MongoDB Version Catalog](/docs/concepts/catalog/mongodb.md) to learn detail concepts of `MongoDBVersion` object.

- you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

- You will also need a docker private [registry](https://docs.docker.com/registry/) or [private repository](https://docs.docker.com/docker-hub/repos/#private-repositories).  In this tutorial we will use private repository of [docker hub](https://hub.docker.com/).

- You have to push the required images from KubeDB's [Docker hub account](https://hub.docker.com/r/kubedb/) into your private registry. For mongodb, push `DB_IMAGE`, `TOOLS_IMAGE`, `EXPORTER_IMAGE` of following MongoDBVersions, where `deprecated` is not true, to your private registry.

  ```console
  $ kubectl get mongodbversions -n kube-system  -o=custom-columns=NAME:.metadata.name,VERSION:.spec.version,DB_IMAGE:.spec.db.image,TOOLS_IMAGE:.spec.tools.image,EXPORTER_IMAGE:.spec.exporter.image,DEPRECATED:.spec.deprecated
  NAME     VERSION   DB_IMAGE              TOOLS_IMAGE                 EXPORTER_IMAGE                   DEPRECATED
  3.4      3.4       kubedb/mongo:3.4      kubedb/mongo-tools:3.4      kubedb/operator:0.8.0            true
  3.4-v1   3.4       kubedb/mongo:3.4-v1   kubedb/mongo-tools:3.4-v2   kubedb/mongodb_exporter:v1.0.0   true
  3.4-v2   3.4       kubedb/mongo:3.4-v2   kubedb/mongo-tools:3.4-v2   kubedb/mongodb_exporter:v1.0.0   <none>
  3.6      3.6       kubedb/mongo:3.6      kubedb/mongo-tools:3.6      kubedb/operator:0.8.0            true
  3.6-v1   3.6       kubedb/mongo:3.6-v1   kubedb/mongo-tools:3.6-v2   kubedb/mongodb_exporter:v1.0.0   true
  3.6-v2   3.6       kubedb/mongo:3.6-v2   kubedb/mongo-tools:3.6-v2   kubedb/mongodb_exporter:v1.0.0   <none>
  4.0      4.0.5     kubedb/mongo:4.0      kubedb/mongo-tools:4.0      kubedb/mongodb_exporter:v1.0.0   <none>
  4.0.5    4.0.5     kubedb/mongo:4.0.5    kubedb/mongo-tools:4.0.5    kubedb/mongodb_exporter:v1.0.0   <none>
  4.1.7    4.1.7     kubedb/mongo:4.1.7    kubedb/mongo-tools:4.1.7    kubedb/mongodb_exporter:v1.0.0   <none>
  ```

  Docker hub repositories:

  - [kubedb/operator](https://hub.docker.com/r/kubedb/operator)
  - [kubedb/mongo](https://hub.docker.com/r/kubedb/mongo)
  - [kubedb/mongo-tools](https://hub.docker.com/r/kubedb/mongo-tools)
  - [kubedb/mongodb_exporter](https://hub.docker.com/r/kubedb/mongodb_exporter)

- Update KubeDB catalog for private Docker registry. Ex:

  ```yaml
  apiVersion: catalog.kubedb.com/v1alpha1
  kind: MongoDBVersion
  metadata:
    name: "3.4-v2"
    labels:
      app: kubedb
  spec:
    version: "3.4"
    db:
      image: "PRIVATE_DOCKER_REGISTRY/mongo:3.4-v2"
    exporter:
      image: "PRIVATE_DOCKER_REGISTRY/mongodb_exporter:v1.0.0"
    tools:
      image: "PRIVATE_DOCKER_REGISTRY/mongo-tools:3.4-v2"
  
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

## Create Demo namespace

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create ns demo
namespace/demo created
```

## Deploy MongoDB database from Private Registry

While deploying `MongoDB` from private repository, you have to add `myregistrykey` secret in `MongoDB` `spec.imagePullSecrets`.
Below is the MongoDB CRD object we will create.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-pvt-reg
  namespace: demo
spec:
  version: "3.4-v2"
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

Now run the command to deploy this `MongoDB` object:

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/mongodb/private-registry/demo-2.yaml
mongodb.kubedb.com/mgo-pvt-reg created
```

To check if the images pulled successfully from the repository, see if the `MongoDB` is in running state:

```console
$ kubectl get pods -n demo -w
NAME            READY     STATUS              RESTARTS   AGE
mgo-pvt-reg-0   0/1       Pending             0          0s
mgo-pvt-reg-0   0/1       Pending             0          0s
mgo-pvt-reg-0   0/1       ContainerCreating   0          0s
mgo-pvt-reg-0   1/1       Running             0          5m


$ kubedb get mg -n demo
NAME          VERSION   STATUS    AGE
mgo-pvt-reg   3.4-v2    Running   38s
```

## Snapshot

You can specify `imagePullSecret` for Snapshot objects in `spec.podTemplate.spec.imagePullSecrets` field of Snapshot object. If you are using scheduled backup, you can also provide `imagePullSecret` in `backupSchedule.podTemplate.spec.imagePullSecrets` field of MongoDB crd. KubeDB also reuses `imagePullSecret` for Snapshot object from `spec.podTemplate.spec.imagePullSecrets` field of MongoDB crd.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
kubectl patch -n demo mg/mgo-pvt-reg -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
kubectl delete -n demo mg/mgo-pvt-reg

kubectl patch -n demo drmn/mgo-pvt-reg -p '{"spec":{"wipeOut":true}}' --type="merge"
kubectl delete -n demo drmn/mgo-pvt-reg

kubectl delete ns demo
```

## Next Steps

- [Snapshot and Restore](/docs/guides/mongodb/snapshot/backup-and-restore.md) process of MongoDB databases using KubeDB.
- Take [Scheduled Snapshot](/docs/guides/mongodb/snapshot/scheduled-backup.md) of MongoDB databases using KubeDB.
- Initialize [MongoDB with Script](/docs/guides/mongodb/initialization/using-script.md).
- Initialize [MongoDB with Snapshot](/docs/guides/mongodb/initialization/using-snapshot.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mongodb/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mongodb/monitoring/using-builtin-prometheus.md).
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Detail concepts of [MongoDBVersion object](/docs/concepts/catalog/mongodb.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
