> New to KubeDB? Please start [here](/docs/guides/README.md).

# Database Scheduled Snapshots

This tutorial will show you how to use KubeDB to take scheduled snapshot of a MongoDB database.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.1/docs/examples/mongodb/demo-0.yaml
namespace "demo" created

$ kubectl get ns
NAME          STATUS    AGE
default       Active    1h
demo          Active    1m
kube-public   Active    1h
kube-system   Active    1h
```

Note that the yaml files that are used in this tutorial, stored in [docs/examples](https://github.com/kubedb/cli/tree/master/docs/examples) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Scheduled Backups

KubeDB supports taking periodic backups for a database using a [cron expression](https://github.com/robfig/cron/blob/v2/doc.go#L26).  KubeDB operator will launch a Job periodically that runs the `mongodump` command and uploads the output bson file to various cloud providers S3, GCS, Azure, OpenStack Swift and/or locally mounted volumes using [osm](https://github.com/appscode/osm).

In this tutorial, snapshots will be stored in a Google Cloud Storage (GCS) bucket. To do so, a secret is needed that has the following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret generic mg-snap-secret -n demo \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "mg-snap-secret" created
```

```yaml
$ kubectl get secret mg-snap-secret -n demo -o yaml
apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2018-02-02T10:02:09Z
  name: mg-snap-secret
  namespace: demo
  resourceVersion: "48679"
  selfLink: /api/v1/namespaces/demo/secrets/mg-snap-secret
  uid: 220a7c60-0800-11e8-946f-080027c05a6e
type: Opaque
```

To learn how to configure other storage destinations for Snapshots, please visit [here](/docs/concepts/snapshot.md).  Now, create the `MongoDB` object with scheduled snapshot.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-scheduled
  namespace: demo
spec:
  version: 3.4
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
  init:
    scriptSource:
      gitRepo:
        repository: "https://github.com/kubedb/mongodb-init-scripts.git"
        directory: .
  backupSchedule:
    cronExpression: "@every 1m"
    storageSecretName: mg-snap-secret
    gcs:
      bucket: restic
```

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.1/docs/examples/mongodb/snapshot/demo-4.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/0.8.0-beta.1/docs/examples/mongodb/snapshot/demo-4.yaml"
mongodb "mgo-scheduled" created
```

It is also possible to add  backup scheduler to an existing `MongoDB`. You just have to edit the `MongoDB` CRD and add below spec:

```yaml
$ kubedb edit mg {db-name} -n demo
spec:
  backupSchedule:
    cronExpression: '@every 1m'
    gcs:
      bucket: restic
    storageSecretName: mg-snap-secret
```

Once the `spec.backupSchedule` is added, KubeDB operator will create a new Snapshot object on each tick of the cron expression. This triggers KubeDB operator to create a Job as it would for any regular instant backup process. You can see the snapshots as they are created using `kubedb get snap` command.

```console
$ kubedb get snap -n demo
NAME                            DATABASE           STATUS      AGE
mgo-scheduled-20180202-104632   mg/mgo-scheduled   Succeeded   4m
mgo-scheduled-20180202-104737   mg/mgo-scheduled   Succeeded   3m
mgo-scheduled-20180202-104837   mg/mgo-scheduled   Succeeded   2m
mgo-scheduled-20180202-104937   mg/mgo-scheduled   Running     10s
```

you should see the output of the `mongodump` command for each snapshot stored in the GCS bucket.

![snapshot-console](/docs/images/mongodb/mgo-scheduled.png)

From the above image, you can see that the snapshot output is stored in a folder called `{bucket}/kubedb/{namespace}/{mongodb-object}/{snapshot}/`.

## Remove Scheduler

To remove scheduler, edit the MongoDB object  to remove `spec.backupSchedule` section.

```yaml
$ kubedb edit mg mgo-scheduled -n demo
apiVersion: kubedb.com/v1alpha1
kind: MongoDB
metadata:
  name: mgo-scheduled
  namespace: demo
  ...
spec:
# backupSchedule:
#   cronExpression: '@every 1m'
#   gcs:
#     bucket: restic
#   storageSecretName: mg-snap-secret
  databaseSecret:
    secretName: mgo-scheduled-auth
  init:
    scriptSource:
      gitRepo:
        directory: .
        repository: https://github.com/kubedb/mongodb-init-scripts.git
  storage:
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 50Mi
    storageClassName: standard
  version: 3.4
status:
  creationTime: 2018-02-02T10:46:18Z
  phase: Running
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubedb delete mg,drmn,snap -n demo --all --force

$ kubectl delete ns demo
namespace "demo" deleted
```

## Next Steps

- See the list of supported storage providers for snapshots [here](/docs/concepts/snapshot.md).
- Initialize [MongoDB with Script](/docs/guides/mongodb/initialization/using-script.md).
- Initialize [MongoDB with Snapshot](/docs/guides/mongodb/initialization/using-snapshot.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box CoreOS Prometheus Operator](/docs/guides/mongodb/monitoring/using-coreos-prometheus-operator.md).
- Monitor your MongoDB database with KubeDB using [out-of-the-box builtin-Prometheus](/docs/guides/mongodb/monitoring/using-builtin-prometheus.md).
- Use [private Docker registry](/docs/guides/mongodb/private-registry/using-private-registry.md) to deploy MongoDB with KubeDB.
- Detail concepts of [MongoDB object](/docs/concepts/databases/mongodb.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
