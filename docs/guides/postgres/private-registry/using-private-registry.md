---
title: Run PostgreSQL using Private Registry
menu:
  docs_0.8.0:
    identifier: pg-using-private-registry-private-registry
    name: Quickstart
    parent: pg-private-registry-postgres
    weight: 10
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using private Docker registry

KubeDB operator supports using private Docker registry. This tutorial will show you how to use KubeDB to run PostgreSQL database using private Docker images.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: Yaml files used in this tutorial are stored in [docs/examples/postgres](https://github.com/kubedb/cli/tree/master/docs/examples/postgres) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

You will also need a docker private [registry](https://docs.docker.com/registry/) or [private repository](https://docs.docker.com/docker-hub/repos/#private-repositories).
In this tutorial we will use private repository of [docker hub](https://hub.docker.com/).

You have to push the required images from KubeDB's [Docker hub account](https://hub.docker.com/r/kubedb/) into your private registry.

For Postgres, push the following images to your private registry.

- [kubedb/operator](https://hub.docker.com/r/kubedb/operator)
- [kubedb/postgres](https://hub.docker.com/r/kubedb/postgres)
- [kubedb/postgres-tools](https://hub.docker.com/r/kubedb/postgres-tools)

```console
$ export DOCKER_REGISTRY=<your-registry>

$ docker pull kubedb/operator:0.9.0-beta.1 ; docker tag kubedb/operator:0.9.0-beta.1 $DOCKER_REGISTRY/operator:0.9.0-beta.1 ; docker push $DOCKER_REGISTRY/operator:0.9.0-beta.1
$ docker pull kubedb/postgres:9.6 ; docker tag kubedb/postgres:9.6 $DOCKER_REGISTRY/postgres:9.6 ; docker push $DOCKER_REGISTRY/postgres:9.6
$ docker pull kubedb/postgres-tools:9.6 ; docker tag kubedb/postgres-tools:9.6 $DOCKER_REGISTRY/postgres-tools:9.6 ; docker push $DOCKER_REGISTRY/postgres-tools:9.6
```

## Create ImagePullSecret

ImagePullSecrets is a type of a Kubernetes Secret whose sole purpose is to pull private images from a Docker registry.
It allows you to specify the url of the docker registry, credentials for logging in and the image name of your private docker image.

Run the following command, substituting the appropriate uppercase values to create an image pull secret for your private Docker registry:

```console
$ kubectl create secret docker-registry myregistrykey \
  --docker-server=DOCKER_REGISTRY_SERVER \
  --docker-username=DOCKER_USER \
  --docker-email=DOCKER_EMAIL \
  --docker-password=DOCKER_PASSWORD

secret "myregistrykey" created.
```

If you wish to follow other ways to pull private images see [official docs](https://kubernetes.io/docs/concepts/containers/images/) of kubernetes.

> Note; If you are using `kubectl` 1.9.0, update to 1.9.1 or later to avoid this [issue](https://github.com/kubernetes/kubernetes/issues/57427).

## Install KubeDB operator

When installing KubeDB operator, set the flags `--docker-registry` and `--image-pull-secret` to appropriate value.
Follow the steps to [install KubeDB operator](/docs/setup/install.md) properly in cluster so that to points to the DOCKER_REGISTRY you wish to pull images from.

## Deploy PostgreSQL database from Private Registry

While deploying PostgreSQL from private repository, you have to add `myregistrykey` secret in Postgres `spec.imagePullSecrets`.

Below is the Postgres object we will create in this tutorial

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: pvt-reg-postgres
  namespace: demo
spec:
  version: "9.6"
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  imagePullSecrets:
    - name: myregistrykey
```

Now run the command to create this Postgres object:

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.1/docs/examples/postgres/private-registry/pvt-reg-postgres.yaml
postgres "pvt-reg-postgres" created
```

To check if the images pulled successfully from the repository, see if the PostgreSQL is in Running state:

```console
$ kubectl get pods -n demo --selector="kubedb.com/name=pvt-reg-postgres" --watch
NAME                 READY     STATUS    RESTARTS   AGE
pvt-reg-postgres-0   1/1       Running   0          41s
^C⏎
```

## Snapshot

We don't need to add `imagePullSecret` for Snapshot objects. Just create Snapshot object and KubeDB operator will reuse the `ImagePullSecret` from Postgres object.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo pg/pvt-reg-postgres -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo pg/pvt-reg-postgres

$ kubectl patch -n demo drmn/pvt-reg-postgres -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/pvt-reg-postgres

$ kubectl delete ns demo
```

If you would like to uninstall KubeDB operator, please follow the steps [here](/docs/setup/uninstall.md).

## Next Steps

- Learn about [taking instant backup](/docs/guides/postgres/snapshot/instant_backup.md) of PostgreSQL database using KubeDB Snapshot.
- Learn how to [schedule backup](/docs/guides/postgres/snapshot/scheduled_backup.md)  of PostgreSQL database.
- Learn about initializing [PostgreSQL with Script](/docs/guides/postgres/initialization/script_source.md).
- Learn about initializing [PostgreSQL from KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Monitor your PostgreSQL database with KubeDB using [built-in Prometheus](/docs/guides/postgres/monitoring/using-builtin-prometheus.md).
- Monitor your PostgreSQL database with KubeDB using [CoreOS Prometheus Operator](/docs/guides/postgres/monitoring/using-coreos-prometheus-operator.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
