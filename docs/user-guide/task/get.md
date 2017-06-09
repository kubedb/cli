# Get TPR objects

Following command will list all Postgres objects in `default` namespace.

```bash
$ kubedb get postgres

NAME            STATUS    AGE
postgres-demo   Running   5h
postgres-dev    Running   4h
postgres-prod   Running   30m
postgres-qa     Running   2h
```

To get YAML of an object, we can provide `--output=yaml` flag

```bash
$ kubedb get postgres postgres-demo --output=yaml

apiVersion: kubedb.com/v1beta1
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

To get JSON of an object, we can provide `--output=json` flag

```bash
$ kubedb get postgres postgres-demo --output=json
```

To list all objects of all supported TPR, we can use following command

```bash
$ kubedb get all -o wide

NAME                    STATUS    VERSION   AGE
es/elasticsearch-demo   Running   2.3.1     17m

NAME               STATUS    VERSION   AGE
pg/postgres-demo   Running   9.5       3h
pg/postgres-dev    Running   9.5       3h
pg/postgres-prod   Running   9.5       3h
pg/postgres-qa     Running   9.5       3h

NAME                                 STATUS      BUCKET          AGE
snap/postgres-demo-20170605-073557   Succeeded   bucket-name     9m
snap/snapshot-20170505-1147          Succeeded   bucket-name     1h
snap/snapshot-xyz                    Succeeded   bucket-name     5m
```

Flag `--output=wide` is used to print additional information.

We can print labels with objects

Following command will list all Snapshots with their corresponding labels.

```bash
$ kubedb get snap --show-labels

NAME                            STATUS      AGE       LABELS
postgres-demo-20170605-073557   Succeeded   11m       kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
snapshot-20170505-1147          Succeeded   1h        kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
snapshot-xyz                    Succeeded   6m        kubedb.com/kind=Elastic,kubedb.com/name=elasticsearch-demo
```

We can also filter list using `--selector` flag.

```bash
$ kubedb get snap --selector='kubedb.com/kind=Postgres' --show-labels

NAME                            STATUS      AGE       LABELS
postgres-demo-20170605-073557   Succeeded   14m       kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
snapshot-20170505-1147          Succeeded   2h        kubedb.com/kind=Postgres,kubedb.com/name=postgres-demo
```

To print only object name, we can use this command
```bash
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

##### Click [here](../reference/get.md) to get command details.
