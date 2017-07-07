> New to KubeDB? Please start [here](/docs/tutorial.md).

# List Databases

`kubedb get` command allows users to list or find any KubeDB object. To list all Postgres objects in `default` namespace, run the following command:

```sh
$ kubedb get postgres

NAME            STATUS    AGE
postgres-demo   Running   5h
postgres-dev    Running   4h
postgres-prod   Running   30m
postgres-qa     Running   2h
```

To get YAML of an object, use `--output=yaml` flag.

```yaml
$ kubedb get postgres postgres-demo --output=yaml

apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: postgres-demo
  namespace: default
spec:
  databaseSecret:
    secretName: postgres-demo-admin-auth
  version: "9.5"
status:
  creationTime: 2017-06-05T04:10:06Z
  phase: Running
```

To get JSON of an object, use `--output=json` flag.

```sh
$ kubedb get postgres postgres-demo --output=json
```

To list all KubeDB objects, use following command:

```sh
$ kubedb get all -o wide

NAME                    VERSION   STATUS    AGE
es/elasticsearch-demo   2.3.1     Running   17m

NAME               VERSION   STATUS    AGE
pg/postgres-demo   9.5       Running   3h
pg/postgres-dev    9.5       Running   3h
pg/postgres-prod   9.5       Running   3h
pg/postgres-qa     9.5       Running   3h

NAME                                 DATABASE                BUCKET             STATUS      AGE
snap/postgres-demo-20170605-073557   pg/postgres-demo        gs:bucket-name     Succeeded   9m
snap/snapshot-20170505-1147          pg/postgres-demo        gs:bucket-name     Succeeded   1h
snap/snapshot-xyz                    es/elasticsearch-demo   local:/directory   Succeeded   5m
```

Flag `--output=wide` is used to print additional information.

You can print labels with objects. The following command will list all Snapshots with their corresponding labels.

```sh
$ kubedb get snap --show-labels

NAME                            DATABASE                STATUS      AGE       LABELS
postgres-demo-20170605-073557   pg/postgres-demo        Succeeded   11m       kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
snapshot-20170505-1147          pg/postgres-demo        Succeeded   1h        kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
snapshot-xyz                    es/elasticsearch-demo   Succeeded   6m        kubedb.com/kind=Elastic,kubedb.com/name=elasticsearch-demo
```

You can also filter list using `--selector` flag.

```sh
$ kubedb get snap --selector='kubedb.com/kind=Postgres' --show-labels

NAME                            DATABASE           STATUS      AGE       LABELS
postgres-demo-20170605-073557   pg/postgres-demo   Succeeded   14m       kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
snapshot-20170505-1147          pg/postgres-demo   Succeeded   2h        kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
```

To print only object name, run the following command:
```sh
$ kubedb get all -o name

elastic/elasticsearch-demo
postgres/postgres-demo
postgres/postgres-dev
postgres/postgres-prod
postgres/postgres-qa
snapshot/postgres-demo-20170605-073557
snapshot/snapshot-20170505-1147
snapshot/snapshot-xyz
```

To learn about various options of `get` command, please visit [here](/docs/reference/kubedb_get.md).
