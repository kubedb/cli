# kubedb delete

## Example

##### Help for delete command

```bash
$ kubedb delete --help

Delete resources by filenames, stdin, resources and names, or by resources and label selector. JSON and YAML formats are
accepted.

Note that the delete command does NOT do resource version checks

Examples:
  # Delete a elastic using the type and name specified in elastic.json.
  kubedb delete -f ./elastic.json

  # Delete a postgres based on the type and name in the JSON passed into stdin.
  cat postgres.json | kubedb delete -f -

  # Delete elastic with label elastic.kubedb.com/name=elasticsearch-demo.
  kubedb delete elastic -l elastic.kubedb.com/name=elasticsearch-demo

Options:
  -f, --filename=[]: Filename to use to delete the resource
  -n, --namespace='default': Delete object(s) from this namespace.
  -o, --output='': Output mode. Use "-o name" for shorter output (resource/name).
  -R, --recursive=false: Process the directory used in -f, --filename recursively.
  -l, --selector='': Selector (label query) to filter on.

Usage:
  kubedb delete ([-f FILENAME] | TYPE [(NAME | -l label)]) [flags] [options]

Use "kubedb delete options" for a list of global command-line options (applies to all commands).
```

##### Delete
```bash
$ kubedb delete pg/postgres-demo

postgres "postgres-demo" deleted
```

##### Delete from file
```bash
$ kubedb delete -f ./elastic.json

elastic "elasticsearch-demo" deleted
```

##### Delete from stdin
```bash
$ cat ./elastic.json | kubedb delete -f -

elastic "elasticsearch-demo" deleted
```
