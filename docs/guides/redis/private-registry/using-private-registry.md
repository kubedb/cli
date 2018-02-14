
> New to KubeDB? Please start [here](/docs/guides/README.md).

# Deploy Redis from Private Docker Registry
`KubeDB` can be installed in a way that it only uses images from a specific docker-registry (may be private images) by providing the flag `--docker-registry=<your-registry>`.

This tutorial will show you how to use KubeDB to run Redis database using Private Docker images. In this tutorial we will create a `ImagePullSecret` and add that secret in `Redis` CRD object specs. If you wish to follow other ways to pull private images see [official docs](https://kubernetes.io/docs/concepts/containers/images/) of kubernetes.

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

You will also need a docker private [registry](https://docs.docker.com/registry/) or [private repository](https://docs.docker.com/docker-hub/repos/#private-repositories).  In this tutorial we will use private repository of [docker hub](https://hub.docker.com/).

Push necessary `KubeDB` images in your repository. For redis, total two images need to be pushed in private registry or repository for running `KubeDB` operator smoothly.

 - [kubedb/operator](https://hub.docker.com/r/kubedb/operator)
 - [kubedb/redis](https://hub.docker.com/r/kubedb/redis)


```console
$ export DOCKER_REGISTRY=<your-registry>

$ docker pull kubedb/operator:0.8.0-beta.0-4 ; docker tag kubedb/operator:0.8.0-beta.0-4 $DOCKER_REGISTRY/operator:0.8.0-beta.0-4 ; docker push $DOCKER_REGISTRY/operator:0.8.0-beta.0-4
$ docker pull kubedb/redis:4 ; docker tag kubedb/redis:4 $DOCKER_REGISTRY/redis:4 ; docker push $DOCKER_REGISTRY/redis:4
```

KubeDB needs to be installed by providing `--docker-registry` and `--image-pull-secret` value. Follow the steps to [install `KubeDB-Operator`](/docs/setup/install.md) properly in cluster so that to points to the DOCKER_REGISTRY you wish to pull images from.

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.1/docs/examples/redis/demo-0.yaml
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    45m
demo          Active    10s
kube-public   Active    45m
kube-system   Active    45m
```

## Create ImagePullSecret
ImagePullSecrets is a type of a Kubernete Secret whose sole purpose is to pull private images from a Docker registry. It allows you to specify the Url of the docker registry, credentials for logging in and the image name of your private docker image.

### Log in to Docker
Before creating `ImagePullSecret`, log in to your private registry manually. This will create a ~/.docker directory and a ~/.docker/config.json file. See here for details of [docker login](https://docs.docker.com/engine/reference/commandline/login/) options.

```console
# docker login <your-registry-server>
$ docker login
Username: <docker-hub-username>
Password:
Login Succeeded
```

When prompted, enter your Docker username and password.

The login process creates or updates a config.json file that holds an authorization token.
View the config.json file. The output contains a section similar to this:

```console
$ cat ~/.docker/config.json
{
    "auths": {
        "https://index.docker.io/v1/": {
            "auth": "c3R...zE2"
        }
    }
}
```

### Create a Secret that holds your authorization token
We will create a secret named `myregistrykey` in `demo` namespace so that kubernetes can pull `redis` and `redis-tools` images from private repository.

```console
$ cat ~/.docker/config.json | base64
<base-64-encoded-json>

$ gedit image-pull-secret.yaml
```


Now paste the below yaml in the gedit editor and replace `<base-64-encoded-json-here>` with base64 encoded `.docker/config.json`.

```yaml
apiVersion: v1
kind: Secret
metadata:
 name: myregistrykey
 namespace: demo
data:
 .dockerconfigjson: <base-64-encoded-json-here>
type: kubernetes.io/dockerconfigjson
```
Now save the yaml file and run the below command to create secret:

```console
$ kubectl create -f image-pull-secret.yaml
secret "myregistrykey" created
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
  version: 4
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
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.1/docs/examples/redis/private-registry/demo-2.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.1/docs/examples/redis/private-registry/demo-2.yaml"
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
$ kubedb delete rd,drmn -n demo --all --force

$ kubectl delete ns demo
namespace "demo" deleted
```


## Next Steps
- Monitor your Redis database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/redis/monitoring/using-coreos-prometheus-operator.md).
- Monitor your Redis database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/redis/monitoring/using-builtin-prometheus.md).
- Use [Private Docker Registry](/docs/guides/redis/private-registry/using-private-registry.md) to deploy Redis with KubeDB.
- Detail concepts of [Redis object](/docs/concepts/databases/redis.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
