> New to KubeDB? Please start [here](/docs/tutorial.md).

# Elastics

## What is Elastic
A `Elastic` is a Kubernetes `Third Party Object` (TPR). It provides declarative configuration for [Elasticsearch](https://www.elastic.co/products/elasticsearch) in a Kubernetes native way. You only need to describe the desired database configuration in a Elastic object, and the KubeDB operator will create Kubernetes objects in the desired state for you.

## Elastic Spec
As with all other Kubernetes objects, a Elastic needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example Elastic object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elastic
metadata:
  name: elasticsearch-db
spec:
  version: 2.3.1
  replicas: 1
```

```sh
$ kubedb create -f  ./docs/examples/elastic/elastic-with-storage.yaml

elastic "elasticsearch-db" created
```

Once the Elastic object is created, KubeDB operator will detect it and create the following Kubernetes objects in the same namespace:
* StatefulSet (name: **elasticsearch-db**-es)
* Service (name: **elasticsearch-db**)
* GoverningService (If not available) (name: **kubedb**)

To confirm the new Elasticsearch is ready, run the following command:
```sh
$ kubedb get elastic elasticsearch-db -o wide

NAME               VERSION   STATUS    AGE
elasticsearch-db   2.3.1     Running   37m
```

This database does not have any PersistentVolume behind StatefulSet pods.


### Using PersistentVolume
To use PersistentVolume, add the `spec.storage` section when creating Elastic object.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elastic
metadata:
  name: elasticsearch-db
spec:
  version: 2.3.1
  replicas: 1
  storage:
    class: "gp2"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: "10Gi"
```

Here we must have to add following storage information in `spec.storage`:

* `class:` StorageClass (`kubectl get storageclasses`)
* `resources:` ResourceRequirements for PersistentVolumeClaimSpec

As `spec.storage` fields are set, StatefulSet will be created with dynamically provisioned PersistentVolumeClaim. Following command will list PVCs for this database.

```bash
$ kubectl get pvc --selector='kubedb.com/kind=Elastic,kubedb.com/name=elasticsearch-db'

NAME                         STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
data-elasticsearch-db-pg-0   Bound     pvc-a1a95954-4a75-11e7-8b69-12f236046fba   10Gi       RWO           2m
```


### Database Initialization
Elasticsearch databases can be created from a previously takes Snapshot. To initialize from prior snapshot, set the `spec.init.snapshotSource` section when creating an Elastic object.

In this case, SnapshotSource must have following information:
1. `namespace:` Namespace of Snapshot object
2. `name:` Name of the Snapshot

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elastic
metadata:
  name: elasticsearch-db
spec:
  version: 2.3.1
  replicas: 1
  init:
    snapshotSource:
      name: "snapshot-xyz"
```

In the above example, Elasticsearch database will be initialized from Snapshot `snapshot-xyz` in `default` namespace. Here,  KubeDB operator will launch a Job to initialize Elasticsearch once StatefulSet pods are running.
