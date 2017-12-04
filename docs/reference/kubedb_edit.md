---
title: Kubedb Edit
menu:
  docs_0.7.1:
    identifier: kubedb-edit
    name: Kubedb Edit
    parent: reference
menu_name: docs_0.7.1
section_menu_id: reference
---
## kubedb edit

Edit a resource on the server

### Synopsis

Edit a resource from the default editor. 

The edit command allows you to directly edit any API resource you can retrieve via the command line tools. It will open the editor defined by your KUBEDB _EDITOR, or EDITOR environment variables, or fall back to 'nano'

```
kubedb edit (RESOURCE/NAME) [flags]
```

### Examples

```
  # Edit the elastic named 'elasticsearch-demo':
  kubedb edit es/elasticsearch-demo
  
  # Use an alternative editor
  KUBEDB_EDITOR="nano" kubedb edit es/elasticsearch-demo
```

### Options

```
  -h, --help               help for edit
  -n, --namespace string   Edit object(s) in this namespace. (default "default")
  -o, --output string      Output format. One of: yaml|json. (default "yaml")
```

### Options inherited from parent commands

```
      --analytics             Send analytical events to Google Analytics (default true)
      --kube-context string   name of the kubeconfig context to use
```

### SEE ALSO

* [kubedb](/docs/reference/kubedb.md)	 - Command line interface for KubeDB


