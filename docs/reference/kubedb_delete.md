---
title: Kubedb Delete
menu:
  docs_0.8.0:
    identifier: kubedb-delete
    name: Kubedb Delete
    parent: reference
menu_name: docs_0.8.0
section_menu_id: reference
---
## kubedb delete

Delete resources by filenames, stdin, resources and names, or by resources and label selector

### Synopsis

Delete resources by filenames, stdin, resources and names, or by resources and label selector. JSON and YAML formats are accepted. 

Note that the delete command does NOT do resource version checks

```
kubedb delete ([-f FILENAME] | TYPE [(NAME | -l label | --all)]) [flags]
```

### Examples

```
  # Delete a elasticsearch using the type and name specified in elastic.json.
  kubedb delete -f ./elastic.json
  
  # Delete a postgres based on the type and name in the JSON passed into stdin.
  cat postgres.json | kubedb delete -f -
  
  # Delete elasticsearch with label elasticsearch.kubedb.com/name=elasticsearch-demo.
  kubedb delete elasticsearch -l elasticsearch.kubedb.com/name=elasticsearch-demo
  
  # Force delete a mysql object
  kubedb delete mysql ms-demo --force
  
  # Delete all mysql objects
  kubedb delete mysql --all
```

### Options

```
      --all                    Delete all resources, including uninitialized ones, in the namespace of the specified resource types.
  -f, --filename stringSlice   Filename to use to delete the resource
      --force                  Immediate deletion of some resources may result in inconsistency or data loss.
  -h, --help                   help for delete
  -n, --namespace string       Delete object(s) from this namespace. (default "default")
  -o, --output string          Output mode. Use "-o name" for shorter output (resource/name).
  -R, --recursive              Process the directory used in -f, --filename recursively.
  -l, --selector string        Selector (label query) to filter on.
```

### Options inherited from parent commands

```
      --analytics             Send analytical events to Google Analytics (default true)
      --kube-context string   name of the kubeconfig context to use
```

### SEE ALSO

* [kubedb](/docs/reference/kubedb.md)	 - Command line interface for KubeDB


