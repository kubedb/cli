---
title: Run Redis using Private Registry
menu:
  docs_0.8.0-beta.2:
    identifier: rd-using-private-registry-private-registry
    name: Quickstart
    parent: rd-private-registry-redis
    weight: 10
menu_name: docs_0.8.0-beta.2
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using private Docker registry

KubeDB operator supports using private Docker registry. This tutorial will show you how to use KubeDB to run Redis database using private Docker images.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

You will also need a docker private [registry](https://docs.docker.com/registry/) or [private repository](https://docs.docker.com/docker-hub/repos/#private-repositories).  In this tutorial we will use private repository of [docker hub](https://hub.docker.com/).

You have to push the required images from KubeDB's [Docker hub account](https://hub.docker.com/r/kubedb/) into your private registry. For redis, push the following images to your private registry.

- [kubedb/operator](https://hub.docker.com/r/kubedb/operator)
- [kubedb/redis](https://hub.docker.com/r/kubedb/redis)

```console
$ export DOCKER_REGISTRY=<your-registry>

$ docker pull kubedb/operator:0.8.0-beta.2 ; docker tag kubedb/operator:0.8.0-beta.2 $DOCKER_REGISTRY/operator:0.8.0-beta.2 ; docker push $DOCKER_REGISTRY/operator:0.8.0-beta.2
$ docker pull kubedb/redis:4 ; docker tag kubedb/redis:4 $DOCKER_REGISTRY/redis:4 ; docker push $DOCKER_REGISTRY/redis:4
```

## Create ImagePullSecret

ImagePullSecrets is a type of a Kubernete Secret whose sole purpose is to pull private images from a Docker registry. It allows you to specify the url of the docker registry, credentials for logging in and the image name of your private docker image.

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

NB: If you are using `kubectl` 1.9.0, update to 1.9.1 or later to avoid this [issue](https://github.com/kubernetes/kubernetes/issues/57427).

## Install KubeDB operator

When installing KubeDB operator, set the flags `--docker-registry` and `--image-pull-secret` to appropriate value. Follow the steps to [install KubeDB operator](/docs/setup/install.md) properly in cluster so that to points to the DOCKER_REGISTRY you wish to pull images from.

## Create Demo namespace

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    45m
demo          Active    10s
kube-public   Active    45m
kube-system   Active    45m
```

## Deploy Redis database from Private Registry

While deploying `Redis` from private repository, you have to add `myregistrykey` secret in `Redis` `spec.imagePullSecrets`.
Below is the Redis CRD object we will create.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Redis
metadata:
  name: redis-pvt-reg
  namespace: demo
spec:
  version: "4"
  doNotPause: true
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

Now run the command to deploy this `Redis` object:

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.2/docs/examples/redis/private-registry/demo-2.yaml
redis "redis-pvt-reg" created
```

To check if the images pulled successfully from the repository, see if the `Redis` is in running state:

```console
$ kubectl get pods -n demo -w
NAME              READY     STATUS              RESTARTS   AGE
redis-pvt-reg-0   0/1       Pending             0          0s
redis-pvt-reg-0   0/1       Pending             0          0s
redis-pvt-reg-0   0/1       ContainerCreating   0          0s
redis-pvt-reg-0   1/1       Running             0          2m


$ kubedb get rd -n demo
NAME            STATUS    AGE
redis-pvt-reg   Running   15s
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo rd/redis-pvt-reg -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo rd/redis-pvt-reg

$ kubectl patch -n demo drmn/redis-pvt-reg -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/redis-pvt-reg

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- Monitor your Redis database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/redis/monitoring/using-coreos-prometheus-operator.md).
- Monitor your Redis database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/redis/monitoring/using-builtin-prometheus.md).
- Detail concepts of [Redis object](/docs/concepts/databases/redis.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
