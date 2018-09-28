---
title: Kubedb Api-Resources
menu:
  docs_0.9.0-beta.0:
    identifier: kubedb-api-resources
    name: Kubedb Api-Resources
    parent: reference
menu_name: docs_0.9.0-beta.0
section_menu_id: reference
---
## kubedb api-resources

Print the supported API resources on the server

### Synopsis

Print the supported API resources on the server

```
kubedb api-resources [flags]
```

### Examples

```
  # Print the supported API Resources
  kubedb api-resources
  # Print the supported API Resources with more information
  kubedb api-resources -o wide
  # Print the supported namespaced resources
  kubedb api-resources --namespaced=true
  # Print the supported non-namespaced resources
  kubedb api-resources --namespaced=false
  # Print the supported API Resources with specific APIGroup
  kubedb api-resources --api-group=extensions
```

### Options

```
      --cached          Use the cached list of resources if available.
  -h, --help            help for api-resources
      --namespaced      If false, non-namespaced resources will be returned, otherwise returning namespaced resources by default. (default true)
      --no-headers      When using the default or custom-column output format, don't print headers (default print headers).
  -o, --output string   Output format. One of: wide|name.
      --verbs strings   Limit to resources that support the specified verbs.
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


