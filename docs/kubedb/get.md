# kubedb get

## Example

##### Get help
```bash
$ kubedb --help

kubedb controls k8sdb ThirdPartyResource objects.

Find more information at https://github.com/k8sdb/kubedb.

Basic Commands (Intermediate):
  get         Display one or many resources

Other Commands:
  help        Help about any command

Use "kubedb <command> --help" for more information about a given command.
```


##### Help for get command

```bash
$ kubedb get --help

Display one or many resources.

Valid resource types include:

  * all
  * elastic
  * postgres
  * databasesnapshot
  * deleteddatabase

Examples:
  # List all elastic in ps output format.
  kubedb get elastics

  # List all elastic in ps output format with more information (such as version).
  kubedb get elastics -o wide

  # List a single postgres with specified NAME in ps output format.
  kubedb get postgres database

  # List a single databasesnapshot in JSON output format.
  kubedb get -o json databasesnapshot snapshot-xyz

  # List all postgreses and elastics together in ps output format.
  kubedb get postgreses,elastics

  # List one or more resources by their type and names.
  kubedb get elastic/es-db postgres/pg-db

Options:
      --all-namespaces=false: If present, list the requested object(s) across all namespaces. Namespace in current
context is ignored even if specified with --namespace.
  -o, --output='': Output format. One of: json|yaml|wide|name.
  -a, --show-all=false: When printing, show all resources (default hide terminated pods.)
      --show-kind=false: If present, list the resource type for the requested object(s).
      --show-labels=false: When printing, show all labels as the last column (default hide labels column)

Usage:
  kubedb get [options]

Use "kubedb get options" for a list of global command-line options (applies to all commands).
```


##### Get Elastic
```bash
$ kubedb get elastic

NAME                      STATUS    AGE
es/elasticsearch-demo     Running   5h
es/elasticsearch-demo-1   Running   4h
```

##### Get All
```bash
$ kubedb get all

NAME                      STATUS    AGE
es/elasticsearch-demo     Running   5h
es/elasticsearch-demo-1   Running   4h

NAME               STATUS    AGE
pg/postgres-demo   Running   1h

NAME               STATUS      AGE
dbs/snapshot-xyz   Succeeded   27m

NAME                     STATUS    AGE
ddb/e2e-elastic-v4xgwz   Deleted   9m
```

##### Get Postgres with labels
```bash
$ kubedb get postgres --show-labels

NAME            STATUS    AGE       LABELS
postgres-demo   Running   1h        k8sdb.com/type=postgres
```

##### Get Elastic with wide
```bash
$ kubedb get elastic -o wide

NAME                   STATUS    VERSION   AGE
elasticsearch-demo     Running   canary    6h
elasticsearch-demo-1   Running   canary    5h
```

##### Get YAML
```bash
$ kubedb get pg postgres-demo -o yaml

apiVersion: k8sdb.com/v1beta1
kind: Postgres
metadata:
  annotations:
    postgres.k8sdb.com/version: canary-db
  creationTimestamp: 2017-05-05T07:04:06Z
  labels:
    k8sdb.com/type: postgres
  name: postgres-demo
  namespace: default
spec:
  databaseSecret:
    secretName: postgres-demo-admin-auth
  replicas: 1
  serviceAccountName: governing-postgres
  version: canary-db
status:
  creationTime: 2017-05-05T07:04:06Z
  phase: Running
```
