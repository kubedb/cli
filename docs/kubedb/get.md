# kubedb get

## Example

##### Get Help
```bash
$ kubedb get --help

Usage:
  kubedb get [flags]

Flags:
      --all-namespaces        If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
      --kube-context string   name of the kubeconfig context to use
  -o, --output string         Output format. One of: json|yaml|wide|name.
  -a, --show-all              When printing, show all resources (default hide terminated pods.)
      --show-kind             If present, list the resource type for the requested object(s).
      --show-labels           When printing, show all labels as the last column (default hide labels column)
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

NAME                   STATUS    VERSION      AGE
elasticsearch-demo     Running   2.3.1-v2.3   6h
elasticsearch-demo-1   Running   2.3.1-v2.3   5h
```

##### Get YAML
```bash
$ kubedb get elastic elasticsearch-demo -o yaml

apiVersion: k8sdb.com/v1beta1
kind: Elastic
metadata:
  annotations:
    elastic.k8sdb.com/version: 2.3.1-v2.3
  creationTimestamp: 2017-04-27T04:46:52Z
  labels:
    k8sdb.com/type: elastic
  name: elasticsearch-demo
  namespace: default
  resourceVersion: "2056"
  selfLink: /apis/k8sdb.com/v1beta1/namespaces/default/elastics/elasticsearch-demo
  uid: 88819884-2b04-11e7-b948-080027d28b41
spec:
  replicas: 1
  serviceAccountName: governing-elasticsearch
  version: 2.3.1-v2.3
status:
  DatabaseStatus: Running
  creationTime: 2017-04-27T04:46:52Z
```