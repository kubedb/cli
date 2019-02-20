---
title: Run Elasticsearch using Private Registry
menu:
  docs_0.9.0:
    identifier: es-using-private-registry-private-registry
    name: Quickstart
    parent: es-private-registry-elasticsearch
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using private Docker registry

KubeDB operator supports using private Docker registry. This tutorial will show you how to run KubeDB managed Elasticsearch database using private Docker images.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace/demo created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: YAML files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Prepare Private Docker Registry

You will also need a docker private [registry](https://docs.docker.com/registry/) or [private repository](https://docs.docker.com/docker-hub/repos/#private-repositories). In this tutorial we will use private repository of [docker hub](https://hub.docker.com/).

You have to push the required images from KubeDB's [Docker hub account](https://hub.docker.com/r/kubedb/) into your private registry.

For Elasticsearch, push the following images to your private registry.

- [kubedb/operator](https://hub.docker.com/r/kubedb/operator)
- [kubedb/elasticsearch](https://hub.docker.com/r/kubedb/elasticsearch)
- [kubedb/elasticsearch-tools](https://hub.docker.com/r/kubedb/elasticsearch-tools)
- [kubedb/elasticsearch_exporter](https://hub.docker.com/r/kubedb/elasticsearch_exporter)

```console
$ export DOCKER_REGISTRY=<your-registry>

$ docker pull kubedb/operator:0.9.0 ; docker tag kubedb/operator:0.9.0 $DOCKER_REGISTRY/operator:0.9.0 ; docker push $DOCKER_REGISTRY/operator:0.9.0
$ docker pull kubedb/elasticsearch:6.3-v1 ; docker tag kubedb/elasticsearch:6.3-v1 $DOCKER_REGISTRY/elasticsearch:6.3-v1 ; docker push $DOCKER_REGISTRY/elasticsearch:6.3-v1
$ docker pull kubedb/elasticsearch-tools:6.3-v1 ; docker tag kubedb/elasticsearch-tools:6.3-v1 $DOCKER_REGISTRY/elasticsearch-tools:6.3-v1 ; docker push $DOCKER_REGISTRY/elasticsearch-tools:6.3-v1
$ docker pull kubedb/elasticsearch_exporter:1.0.2 ; docker tag kubedb/elasticsearch_exporter:1.0.2 $DOCKER_REGISTRY/elasticsearch_exporter:1.0.2 ; docker push $DOCKER_REGISTRY/elasticsearch_exporter:1.0.2
```

## Create ImagePullSecret

ImagePullSecrets is a type of a Kubernetes Secret whose sole purpose is to pull private images from a Docker registry. It allows you to specify the url of the docker registry, credentials for logging in and the image name of your private docker image.

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

## Create ElasticsearchVersion CRD

KubeDB uses images specified in ElasticsearchVersion crd for database, backup and exporting prometheus metrics. You have to create a ElasticsearchVersion crd specifying images from your private registry. Then, you have to point this ElasticsearchVersion crd in `spec.version` field of Elasticsearch object. For more details about ElasticsearchVersion crd, please visit [here](/docs/concepts/catalog/elasticsearch.md).

Here, is an example of ElasticsearchVersion crd. Replace `<YOUR_PRIVATE_REGISTRY>` with your private registry.

```yaml
apiVersion: catalog.kubedb.com/v1alpha1
kind: ElasticsearchVersion
metadata:
  name: "pvt-6.3"
  labels:
    app: kubedb
spec:
  version: "6.3"
  db:
    image: "<YOUR_PRIVATE_REGISTRY>/elasticsearch:6.3-v1"
  exporter:
    image: "<YOUR_PRIVATE_REGISTRY>/elasticsearch_exporter:1.0.2"
  tools:
    image: "<YOUR_PRIVATE_REGISTRY>/elasticsearch-tools:6.3-v1"
```

Now, create the ElasticsearchVersion crd,

```console
$ kubectl apply -f pvt-elasticsearchversion.yaml
elasticsearchversion.kubedb.com/pvt-6.3 created
```

## Install KubeDB operator

When installing KubeDB operator, set the flags `--docker-registry` and `--image-pull-secret` to appropriate value. Follow the guide for customizing installer to see how to pass those flags from [here](/docs/setup/install.md#customizing-installer).

## Deploy Elasticsearch database from Private Registry

While deploying Elasticsearch from private repository, you have to add `myregistrykey` secret in Elasticsearch `spec.podTemplate.spec.imagePullSecrets`.

Below is the YAML for Elasticsearch crd that will be created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: pvt-reg-elasticsearch
  namespace: demo
spec:
  version: "pvt-6.3"
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

Now run the command to deploy this Elasticsearch object:

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/private-registry/private-registry.yaml
elasticsearch.kubedb.com/pvt-reg-elasticsearch created
```

To check if the images pulled successfully from the repository, see if the Elasticsearch is in running state:

```console
$ kubectl get es -n demo pvt-reg-elasticsearch -o wide
NAME                    VERSION   STATUS       AGE
pvt-reg-elasticsearch   pvt-6.3   Running      33m
```

## Snapshot

You can specify `imagePullSecret` for Snapshot objects in `spec.podTemplate.spec.imagePullSecrets` field of Snapshot object. If you are using scheduled backup, you can also provide `imagePullSecret` in `backupSchedule.podTemplate.spec.imagePullSecrets` field of Elasticsearch crd. KubeDB also reuses `imagePullSecret` for Snapshot object from `spec.podTemplate.spec.imagePullSecrets` field of Elasticsearch crd.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/pvt-reg-elasticsearch -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo es/pvt-reg-elasticsearch

$ kubectl delete ns demo
```

## Next Steps

- Learn about [taking instant backup](/docs/guides/elasticsearch/snapshot/instant_backup.md) of Elasticsearch database using KubeDB.
- Learn how to [schedule backup](/docs/guides/elasticsearch/snapshot/scheduled_backup.md)  of Elasticsearch database.
- Learn about initializing [Elasticsearch with Snapshot](/docs/guides/elasticsearch/initialization/snapshot_source.md).
- Learn how to configure [Elasticsearch Topology](/docs/guides/elasticsearch/clustering/topology.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` builtin-Prometheus](/docs/guides/elasticsearch/monitoring/using-builtin-prometheus.md).
- Monitor your Elasticsearch database with KubeDB using [`out-of-the-box` CoreOS Prometheus Operator](/docs/guides/elasticsearch/monitoring/using-coreos-prometheus-operator.md).
- Detail concepts of [Elasticsearch object](/docs/concepts/databases/elasticsearch.md).
- Detail concepts of [Snapshot object](/docs/concepts/snapshot.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
