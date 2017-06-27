## kubedb describe

Show details of a specific resource or group of resources

### Synopsis


Show details of a specific resource or group of resources. This command joins many API calls together to form a detailed description of a given resource or group of resources.Valid resource types include: 

  * all  
  * elastic  
  * postgres  
  * snapshot  
  * dormantdatabase

```
kubedb describe (TYPE [NAME_PREFIX] | TYPE/NAME) [flags]
```

### Examples

```
  # Describe a elastic
  kubedb describe elastics elasticsearch-demo
  
  # Describe a postgres
  kubedb describe pg/postgres-demo
  
  # Describe all dormantdatabases
  kubedb describe drmn
```

### Options

```
      --all-namespaces     If present, describe the requested object(s) across all namespaces. Namespace specified with --namespace will be ignored.
  -h, --help               help for describe
  -n, --namespace string   Describe object(s) from this namespace. (default "default")
  -l, --selector string    Selector (label query) to filter on, supports '=', '==', and '!='.
      --show-events        If true, display events related to the described object. (default true)
```

### Options inherited from parent commands

```
      --analytics             Send events to Google Analytics (default true)
      --kube-context string   name of the kubeconfig context to use
```

### SEE ALSO
* [kubedb](kubedb.md)	 - Command line interface for KubeDB


