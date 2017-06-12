## kubedb delete

Delete resources by filenames, stdin, resources and names, or by resources and label selector

### Synopsis


Delete resources by filenames, stdin, resources and names, or by resources and label selector. JSON and YAML formats are accepted. 

Note that the delete command does NOT do resource version checks

```
kubedb delete ([-f FILENAME] | TYPE [(NAME | -l label)]) [flags]
```

### Examples

```
  # Delete a elastic using the type and name specified in elastic.json.
  kubedb delete -f ./elastic.json
  
  # Delete a postgres based on the type and name in the JSON passed into stdin.
  cat postgres.json | kubedb delete -f -
  
  # Delete elastic with label elastic.kubedb.com/name=elasticsearch-demo.
  kubedb delete elastic -l elastic.kubedb.com/name=elasticsearch-demo
```

### Options

```
  -f, --filename stringSlice   Filename to use to delete the resource
  -h, --help                   help for delete
  -n, --namespace string       Delete object(s) from this namespace. (default "default")
  -o, --output string          Output mode. Use "-o name" for shorter output (resource/name).
  -R, --recursive              Process the directory used in -f, --filename recursively.
  -l, --selector string        Selector (label query) to filter on.
```

### Options inherited from parent commands

```
      --analytics             Send events to Google Analytics (default true)
      --kube-context string   name of the kubeconfig context to use
```

### SEE ALSO
* [kubedb](kubedb.md)	 - Controls kubedb objects


