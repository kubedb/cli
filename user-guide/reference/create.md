### kubedb create

```bash
$ kubedb create --help

Create a resource by filename or stdin. 

JSON and YAML formats are accepted.

Examples:
  # Create a elastic using the data in elastic.json.
  kubedb create -f ./elastic.json
  
  # Create a elastic based on the JSON passed into stdin.
  cat elastic.json | kubedb create -f -

Options:
  -f, --filename=[]: Filename to use to create the resource
  -n, --namespace='default': Create object(s) in this namespace.
  -R, --recursive=false: Process the directory used in -f, --filename recursively.

Usage:
  kubedb create [flags] [options]

Use "kubedb create options" for a list of global command-line options (applies to all commands).
```

We can say in which namespace we want to create this object providing `--namespace` flag.

We can also provide `--recursive` flag to process the directory used in `-f`, `--filename` recursively.

`kubedb create` command also supports input from _stdin_ with `-f -`
