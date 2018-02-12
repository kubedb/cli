> New to KubeDB Postgres?  Quick start [here](/docs/guides/postgres/quickstart.md).

# Private Docker Registry

KubeDB supports Postgres docker images from non-public docker registry. A *docker-registry* type Secret is used to provide necessary information to KubeDB.

This tutorial will show you how to create this *docker-registry* type Secret and add that Secret in Postgres object.
If you wish to follow other ways to pull private images see [official docs](https://kubernetes.io/docs/concepts/containers/images/) of kubernetes.

## Before You Begin

At first, You will need a docker [registry](https://docs.docker.com/registry/) of your own or a [private repository](https://docs.docker.com/docker-hub/repos/#private-repositories) in docker hub.

In this tutorial, we will use private repository of [docker hub](https://hub.docker.com/).

Push necessary KubeDB images for Postgres into your repository.

For Postgres, these three images are needed to be pushed into your private repository for running KubeDB operator smoothly.

 - [kubedb/operator](https://hub.docker.com/r/kubedb/operator)
 - [kubedb/mongo](https://hub.docker.com/r/kubedb/postgres)
 - [kubedb/mongo-tools](https://hub.docker.com/r/kubedb/postgres-tools)


```console
$ export DOCKER_REGISTRY=<your-registry>

# Pull and Push kubedb/operator
$ docker pull kubedb/operator:0.8.0-beta.0-4
$ docker tag kubedb/operator:0.8.0-beta.0-4 $DOCKER_REGISTRY/operator:0.8.0-beta.0-4
$ docker push $DOCKER_REGISTRY/operator:0.8.0-beta.0-4

# Pull and Push kubedb/postgres
$ docker pull kubedb/postgres:9.6
$ docker tag kubedb/postgres:9.6 $DOCKER_REGISTRY/postgres:9.6
$ docker push $DOCKER_REGISTRY/postgres:9.6

# Pull and Push kubedb/postgres-tools
$ docker pull kubedb/postgres-tools:9.6
$ docker tag kubedb/postgres-tools:9.6 $DOCKER_REGISTRY/postgres-tools:9.6
$ docker push $DOCKER_REGISTRY/postgres-tools:9.6
```

KubeDB needs to be installed by providing `--docker-registry` and `--image-pull-secret` value.

Follow the steps to [install `KubeDB operator`](/docs/setup/install.md) properly in cluster so that it points to the DOCKER_REGISTRY you wish to pull images from.

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created
```

## Create *docker-registry* type Secret

Kubernetes Secret with `type: kubernetes.io/dockercfg` is used to keep docker registry credentials. This Secret is used in Pod as `spec.imagePullSecrets`.

It allows you to specify the server address of the docker registry and credentials to access it.

Now create a Secret `private-registry` in Namespace `demo`

```console
$ kubectl create secret docker-registry private-registry \
    --docker-server=<server location for Docker registry> \
    --docker-username=<username for Docker registry authentication> \
    --docker-password=<password for Docker registry authentication> \
    --docker-email=<email for Docker registry>

secret "private-registry" created
```

```yaml
$ kubectl get secret -n demo private-registry -o yaml
apiVersion: v1
data:
  .dockercfg: "PHlvdSBkb2NrZXIgY29uZmlnPgo="
kind: Secret
metadata:
  creationTimestamp: 2018-02-09T08:11:43Z
  name: private-registry
  namespace: demo
  resourceVersion: "26220"
  selfLink: /api/v1/namespaces/demo/secrets/private-registry
  uid: dd1d6b45-0d70-11e8-9632-080027966966
type: kubernetes.io/dockercfg
```

## Create Postgres with Private Registry

While deploying Postgres from private repository, you have to set `spec.imagePullSecrets`

Below is the Postgres object created in this tutorial

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: pvt-reg-postgres
  namespace: demo
spec:
  version: 9.6
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  imagePullSecrets:
    - name: private-registry
```

Now run the command to create this Postgres object:

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/postgres/private-registry/pvt-reg-postgres.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.1/docs/examples/postgres/private-registry/pvt-reg-postgres.yaml"
postgres "pvt-reg-postgres" created
```

To check if the images pulled successfully from the repository, see if the Postgres is in Running state:

```console
$ kubectl get pods -n demo --selector="kubedb.com/name=pvt-reg-postgres" --watch
NAME                 READY     STATUS    RESTARTS   AGE
pvt-reg-postgres-0   1/1       Running   0          41s
^C⏎
```

## Snapshot

You don't need to add `imagePullSecret` in Snapshot objects. KubeDB operator will re-use the `spec.imagePullSecrets` from Postgres object.

## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete pg,drmn,snap -n demo --all --force
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
