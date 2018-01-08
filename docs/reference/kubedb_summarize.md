---
title: Kubedb Summarize
menu:
  docs_0.8.0-beta.0:
    identifier: kubedb-summarize
    name: Kubedb Summarize
    parent: reference
menu_name: docs_0.8.0-beta.0
section_menu_id: reference
---
## kubedb summarize

Export summary report

### Synopsis

Export summary report

```
kubedb summarize [flags]
```

### Options

```
  -h, --help                        help for summarize
      --index string                Export summary report for this only.
  -n, --namespace string            Export summary report of the requested object from this namespace. (default "default")
      --operator-namespace string   Name of namespace where operator is running (default "kube-system")
      --output string               Directory used to store summary report
```

### Options inherited from parent commands

```
      --analytics             Send analytical events to Google Analytics (default true)
      --kube-context string   name of the kubeconfig context to use
```

### SEE ALSO

* [kubedb](/docs/reference/kubedb.md)	 - Command line interface for KubeDB


