---
title: Kubedb Get
menu:
  docs_0.10.0:
    identifier: kubedb-get
    name: Kubedb Get
    parent: reference
menu_name: docs_0.10.0
section_menu_id: reference
---
## kubedb get

Display one or many resources

### Synopsis

Display one or many resources.

Use "kubedb api-resources" for a complete list of supported resources.

```
kubedb get [(-o|--output=)json|yaml|wide|custom-columns=...|custom-columns-file=...|go-template=...|go-template-file=...|jsonpath=...|jsonpath-file=...] (TYPE[.VERSION][.GROUP] [NAME | -l label] | TYPE[.VERSION][.GROUP]/NAME ...) [flags]
```

### Examples

```
  # List all elasticsearch in ps output format.
  kubedb get es
  
  # List all elasticsearch in ps output format with more information (such as version).
  kubedb get elasticsearches -o wide
  
  # List a single postgres with specified NAME in ps output format.
  kubedb get postgres database
  
  # List a single snapshot in JSON output format.
  kubedb get -o json snapshot snapshot-xyz
  
  # List all postgreses and elastics together in ps output format.
  kubedb get postgreses,elastics
  
  # List one or more resources by their type and names.
  kubedb get es/es-db postgres/pg-db
  
  Valid resource types include:
  * all
  * etcds
  * elasticsearches
  * postgreses
  * mysqls
  * mongodbs
  * redises
  * memcacheds
  * snapshots
  * dormantdatabases
```

### Options

```
      --all-namespaces                If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
      --allow-missing-template-keys   If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats. (default true)
      --chunk-size int                Return large lists in chunks rather than all at once. Pass 0 to disable. This flag is beta and may change in the future. (default 500)
      --export                        If true, use 'export' for the resources.  Exported resources are stripped of cluster-specific information.
      --field-selector string         Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. --field-selector key1=value1,key2=value2). The server only supports a limited number of field queries per type.
  -f, --filename strings              Filename, directory, or URL to files identifying the resource to get from a server.
  -h, --help                          help for get
      --ignore-not-found              If the requested object does not exist the command will return exit code 0.
      --include-uninitialized         If true, the kubectl command applies to uninitialized objects. If explicitly set to false, this flag overrides other flags that make the kubectl commands apply to uninitialized objects, e.g., "--all". Objects with empty metadata.initializers are regarded as initialized.
  -L, --label-columns strings         Accepts a comma separated list of labels that are going to be presented as columns. Names are case-sensitive. You can also use multiple flag options like -L label1 -L label2...
      --no-headers                    When using the default or custom-column output format, don't print headers (default print headers).
  -o, --output string                 Output format. One of: json|yaml|wide|name|custom-columns=...|custom-columns-file=...|go-template=...|go-template-file=...|jsonpath=...|jsonpath-file=... See custom columns [http://kubernetes.io/docs/user-guide/kubectl-overview/#custom-columns], golang template [http://golang.org/pkg/text/template/#pkg-overview] and jsonpath template [http://kubernetes.io/docs/user-guide/jsonpath].
      --raw string                    Raw URI to request from the server.  Uses the transport specified by the kubeconfig file.
  -R, --recursive                     Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
  -l, --selector string               Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)
      --server-print                  If true, have the server return the appropriate table output. Supports extension APIs and CRDs. (default true)
      --show-kind                     If present, list the resource type for the requested object(s).
      --show-labels                   When printing, show all labels as the last column (default hide labels column)
      --sort-by string                If non-empty, sort list types using this field specification.  The field specification is expressed as a JSONPath expression (e.g. '{.metadata.name}'). The field in the API resource specified by this JSONPath expression must be an integer or a string.
      --template string               Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].
  -w, --watch                         After listing/getting the requested object, watch for changes. Uninitialized objects are excluded if no object name is provided.
      --watch-only                    Watch for changes to the requested object(s), without listing/getting first.
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --as string                        Username to impersonate for the operation
      --as-group stringArray             Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --cache-dir string                 Default HTTP cache directory (default "/home/tamal/.kube/http-cache")
      --certificate-authority string     Path to a cert file for the certificate authority
      --client-certificate string        Path to a client certificate file for TLS
      --client-key string                Path to a client key file for TLS
      --cluster string                   The name of the kubeconfig cluster to use
      --context string                   The name of the kubeconfig context to use
      --enable-analytics                 Send analytical events to Google Analytics (default true)
      --insecure-skip-tls-verify         If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string                Path to the kubeconfig file to use for CLI requests.
      --log-backtrace-at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log-dir string                   If non-empty, write log files in this directory
      --log-flush-frequency duration     Maximum number of seconds between log flushes (default 5s)
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


