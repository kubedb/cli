---
title: Kubedb Delete
menu:
  docs_0.9.0-rc.1:
    identifier: kubedb-delete
    name: Kubedb Delete
    parent: reference
menu_name: docs_0.9.0-rc.1
section_menu_id: reference
---
## kubedb delete

Delete resources by filenames, stdin, resources and names, or by resources and label selector

### Synopsis

Delete resources by filenames, stdin, resources and names, or by resources and label selector. JSON and YAML formats are accepted. 

Note that the delete command does NOT do resource version checks

```
kubedb delete ([-f FILENAME] | TYPE [(NAME | -l label | --all)])
```

### Examples

```
  # Delete a elasticsearch using the type and name specified in elastic.json.
  kubedb delete -f ./elastic.json
  
  # Delete a postgres based on the type and name in the JSON passed into stdin.
  cat postgres.json | kubedb delete -f -
  
  # Delete elasticsearch with label elasticsearch.kubedb.com/name=elasticsearch-demo.
  kubedb delete elasticsearch -l elasticsearch.kubedb.com/name=elasticsearch-demo
  
  # Delete all mysql objects
  kubedb delete mysql --all
```

### Options

```
      --all                     Delete all resources, including uninitialized ones, in the namespace of the specified resource types.
      --cascade                 If true, cascade the deletion of the resources managed by this resource (e.g. Pods created by a ReplicationController).  Default true. (default true)
      --field-selector string   Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. --field-selector key1=value1,key2=value2). The server only supports a limited number of field queries per type.
  -f, --filename strings        containing the resource to delete.
      --force                   Only used when grace-period=0. If true, immediately remove resources from API and bypass graceful deletion. Note that immediate deletion of some resources may result in inconsistency or data loss and requires confirmation.
      --grace-period int        Period of time in seconds given to the resource to terminate gracefully. Ignored if negative. Set to 1 for immediate shutdown. Can only be set to 0 when --force is true (force deletion). (default -1)
  -h, --help                    help for delete
      --ignore-not-found        Treat "resource not found" as a successful delete. Defaults to "true" when --all is specified.
      --include-uninitialized   If true, the kubectl command applies to uninitialized objects. If explicitly set to false, this flag overrides other flags that make the kubectl commands apply to uninitialized objects, e.g., "--all". Objects with empty metadata.initializers are regarded as initialized.
      --now                     If true, resources are signaled for immediate shutdown (same as --grace-period=1).
  -o, --output string           Output mode. Use "-o name" for shorter output (resource/name).
  -R, --recursive               Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
  -l, --selector string         Selector (label query) to filter on, not including uninitialized ones.
      --timeout duration        The length of time to wait before giving up on a delete, zero means determine a timeout from the size of the object
      --wait                    If true, wait for resources to be gone before returning. This waits for finalizers. (default true)
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
      --stderrthreshold severity         logs at or above this threshold go to stderr
      --token string                     Bearer token for authentication to the API server
      --user string                      The name of the kubeconfig user to use
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [kubedb](/docs/reference/kubedb.md)	 - Command line interface for KubeDB


