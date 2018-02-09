> New to KubeDB Postgres?  Quick start [here](/docs/guides/postgres/quickstart.md).

# KubeDB Snapshot

KubeDB operator maintains another Custom Resource Definition (CRD) for database backups called Snapshot. Snapshot object is used to take backup or restore from a backup.

### Before You Begin

In this tutorial, we will take instant backup of a PostgreSQL database using KubeDB Snapshot object.

So, lets create a Postgres object first following [this tutorial](/docs/guides/postgres/initialization/script_source.md#script-source).

```console
$ kubedb get pg -n demo script-postgres
NAME              STATUS    AGE
script-postgres   Running   24m
```

We will take backup of this PostgreSQL database `script-postgres`.

## Instant Backup

Snapshot provides a declarative configuration for backup behavior in a Kubernetes native way.

Below is the Snapshot object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Snapshot
metadata:
  name: instant-snapshot
  namespace: demo
  labels:
    kubedb.com/kind: Postgres
spec:
  databaseName: script-postgres
  storageSecretName: gcs-secret
  gcs:
    bucket: kubedb
```

Here,

 - `metadata.labels` should include the type of database.
 - `spec.databaseName` indicates the Postgres object name, `p1`, whose snapshot is taken.
 - `spec.storageSecretName` points to the Secret containing the credentials for snapshot storage destination.
 - `spec.gcs.bucket` points to the bucket name used to store the snapshot data.

In this case, `kubedb.com/kind: Postgres` tells KubeDB operator that this Snapshot belongs to a Postgres object.
Only Postgres controller will handle this Snapshot object.

> Note: Snapshot and Secret objects must be in the same namespace as Postgres, `p1`, in our case.


##### Snapshot Storage Secret

Storage Secret should contain credentials that will be used to access storage destination.
In this tutorial, snapshot data will be stored in a Google Cloud Storage (GCS) bucket.

For that a storage Secret is needed with following 2 keys:

| Key                               | Description                                                |
|-----------------------------------|------------------------------------------------------------|
| `GOOGLE_PROJECT_ID`               | `Required`. Google Cloud project ID                        |
| `GOOGLE_SERVICE_ACCOUNT_JSON_KEY` | `Required`. Google Cloud service account json key          |

```console
$ echo -n '<your-project-id>' > GOOGLE_PROJECT_ID
$ mv downloaded-sa-json.key > GOOGLE_SERVICE_ACCOUNT_JSON_KEY
$ kubectl create secret -n demo generic gcs-secret \
    --from-file=./GOOGLE_PROJECT_ID \
    --from-file=./GOOGLE_SERVICE_ACCOUNT_JSON_KEY
secret "gcs-secret" created
```

```yaml
$ kubectl get secret -n demo gcs-secret -o yaml
apiVersion: v1
data:
  GOOGLE_PROJECT_ID: PHlvdXItcHJvamVjdC1pZD4=
  GOOGLE_SERVICE_ACCOUNT_JSON_KEY: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3V...9tIgp9Cg==
kind: Secret
metadata:
  creationTimestamp: 2018-02-05T06:10:50Z
  name: gcs-secret
  namespace: demo
  resourceVersion: "3869"
  selfLink: /api/v1/namespaces/demo/secrets/gcs-secret
  uid: 5055ce8e-0a3b-11e8-b4de-42010a8000be
type: Opaque
```

##### Snapshot Storage Backend

KubeDB supports various cloud providers (_S3_, _GCS_, _Azure_, _OpenStack_ _Swift_ and/or locally mounted volumes) as snapshot storage backend.
In this tutorial, _GCS_ backend is used.

To configure this backend, following parameters are available:

| Parameter                | Description                                                                     |
|--------------------------|---------------------------------------------------------------------------------|
| `spec.gcs.bucket`        | `Required`. Name of bucket                                                      |
| `spec.gcs.prefix`        | `Optional`. Path prefix into bucket where snapshot data will be stored          |

> An open source project [osm](https://github.com/appscode/osm) is used to store snapshot data into cloud.

To lean how to configure other storage destinations for snapshot data, please visit [here](/docs/concepts/snapshot.md).

Now, create the Snapshot object.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/postgres/snapshot/instant-snapshot.yaml
validating "https://raw.githubusercontent.com/kubedb/cli/postgres-docs/docs/examples/postgres/snapshot/instant-snapshot.yaml"
snapshot "instant-snapshot" created
```

Lets see Snapshot list of Postgres `script-postgres`.

```console
$ kubedb get snap -n demo --selector="kubedb.com/kind=Postgres,kubedb.com/name=script-postgres"
NAME               DATABASE             STATUS    AGE
instant-snapshot   pg/script-postgres   Running   42s
```

KubeDB operator watches for Snapshot objects using Kubernetes API. When a Snapshot object is created, it will launch a Job that runs the `pg_dumpall` command and
uploads the output **sql** file to cloud storage using [osm](https://github.com/appscode/osm).

Snapshot data is stored in a folder called `{bucket}/{prefix}/kubedb/{namespace}/{Postgres name}/{Snapshot name}/`.

Once the snapshot Job is completed, you can see the output of the `pg_dumpall` command stored in the GCS bucket.

<p align="center">
  <kbd>
    <img alt="snapshot-console"  src="/docs/images/postgres/instant-snapshot.png">
  </kbd>
</p>

From the above image, you can see that the snapshot data file `dumpfile.sql` is stored in your bucket.

If you open this `dumpfile.sql` file, you will see the query to create `dashboard` TABLE.

```console
--
-- Name: dashboard; Type: TABLE; Schema: data; Owner: postgres
--

CREATE TABLE dashboard (
    id bigint NOT NULL,
    version integer NOT NULL,
    slug character varying(255) NOT NULL,
    title character varying(255) NOT NULL,
    data text NOT NULL,
    org_id bigint NOT NULL,
    created timestamp without time zone NOT NULL,
    updated timestamp without time zone NOT NULL,
    updated_by integer,
    created_by integer
);


ALTER TABLE dashboard OWNER TO postgres;
```


Lets see the Snapshot list for Postgres `script-postgres` by running `kubedb describe` command.

```console
$ kubedb describe pg -n demo script-postgres -S=false -W=false
Name:           script-postgres
Namespace:      demo
StartTimestamp: Thu, 08 Feb 2018 15:55:11 +0600
Status:         Running
Init:
  scriptSource:
    Type:       GitRepo (a volume that is pulled from git when the pod is created)
    Repository: https://github.com/kubedb/postgres-init-scripts.git
    Directory:  .
Volume:
  StorageClass: standard
  Capacity:     50Mi
  Access Modes: RWO
StatefulSet:    script-postgres
Service:        script-postgres, script-postgres-primary
Secrets:        script-postgres-auth

Topology:
  Type      Pod                 StartTime                       Phase
  ----      ---                 ---------                       -----
  primary   script-postgres-0   2018-02-08 15:55:29 +0600 +06   Running

Snapshots:
  Name               Bucket      StartTime                         CompletionTime                    Phase
  ----               ------      ---------                         --------------                    -----
  instant-snapshot   gs:kubedb   Thu, 08 Feb 2018 16:30:29 +0600   Thu, 08 Feb 2018 16:31:54 +0600   Succeeded

Events:
  FirstSeen   LastSeen   Count     From                  Type       Reason               Message
  ---------   --------   -----     ----                  --------   ------               -------
  11m         11m        1         Job Controller        Normal     SuccessfulSnapshot   Successfully completed snapshot
  12m         12m        1         Snapshot Controller   Normal     Starting             Backup running
  48m         48m        1         Postgres operator     Normal     Successful           Successfully patched StatefulSet
  48m         48m        1         Postgres operator     Normal     Successful           Successfully patched Postgres
  48m         48m        1         Postgres operator     Normal     Successful           Successfully created StatefulSet
  48m         48m        1         Postgres operator     Normal     Successful           Successfully created Postgres
  48m         48m        1         Postgres operator     Normal     Successful           Successfully created Service
  48m         48m        1         Postgres operator     Normal     Successful           Successfully created Service
```


## Cleanup Snapshot

If you want to delete snapshot data from storage, you can delete Snapshot object.

```console
$ kubectl delete snap -n demo instant-snapshot
snapshot "instant-snapshot" deleted
```

Once Snapshot object is deleted, you can't revert this process and snapshot data from storage will be deleted permanently.

## Next Steps
- Setup [Continuous Archiving](/docs/guides/postgres/snapshot/continuous_archiving.md) in PostgreSQL using `wal-g`
- Learn how to [schedule backup](/docs/guides/postgres/snapshot/scheduled_backup.md)  of PostgreSQL database.
- Learn about initializing [PostgreSQL from KubeDB Snapshot](/docs/guides/postgres/initialization/snapshot_source.md).
- Want to setup PostgreSQL cluster? Check how to [configure Highly Available PostgreSQL Cluster](/docs/guides/postgres/clustering/ha_cluster.md)
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
