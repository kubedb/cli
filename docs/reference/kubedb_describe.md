---
title: Kubedb Describe
menu:
  docs_0.8.0:
    identifier: kubedb-describe
    name: Kubedb Describe
    parent: reference
menu_name: docs_0.8.0
section_menu_id: reference
---
## kubedb describe

Show details of a specific resource or group of resources

### Synopsis

Show details of a specific resource or group of resources. This command joins many API calls together to form a detailed description of a given resource or group of resources.Valid resource types include:  * all  * elasticsearches  * postgreses  * mysqls  * mongodbs  * redises  * memcacheds  * snapshots  * dormantdatabases

Use "kubedb api-resources" for a complete list of supported resources.

```
kubedb describe (-f FILENAME | TYPE [NAME_PREFIX | -l label] | TYPE/NAME)
```

### Examples

```
  # Describe a elasticsearch
  kubedb describe elasticsearches elasticsearch-demo
  
  # Describe a postgres
  kubedb describe pg/postgres-demo
  
  # Describe all dormantdatabases
  kubedb describe drmn
```

### Options

```
      --all-namespaces          If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
  -f, --filename strings        Filename, directory, or URL to files containing the resource to describe
  -h, --help                    help for describe
      --include-uninitialized   If true, the kubectl command applies to uninitialized objects. If explicitly set to false, this flag overrides other flags that make the kubectl commands apply to uninitialized objects, e.g., "--all". Objects with empty metadata.initializers are regarded as initialized.
  -R, --recursive               Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
  -l, --selector string         Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)
      --show-events             If true, display events related to the described object. (default true)
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --analytics                        Send analytical events to Google Analytics (default true)
      --as string                        Username to impersonate for the operation
      --as-group stringArray             Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --cache-dir string                 Default HTTP cache directory (default "/home/tamal/.kube/http-cache")
      --certificate-authority string     Path to a cert file for the certificate authority
      --client-certificate string        Path to a client certificate file for TLS
      --client-key string                Path to a client key file for TLS
      --cluster string                   The name of the kubeconfig cluster to use
      --context string                   The name of the kubeconfig context to use
      --insecure-skip-tls-verify         If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string                Path to the kubeconfig file to use for CLI requests.
      --log-backtrace-at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log-dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --match-server-version             Require server version to match client version
  -n, --namespace string                 If present, the namespace scope for this CLI request
      --request-timeout string           The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                    The address and port of the Kubernetes API server
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
      --token string                     Bearer token for authentication to the API server
      --user string                      The name of the kubeconfig user to use
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [kubedb](/docs/reference/kubedb.md)	 - Command line interface for KubeDB


