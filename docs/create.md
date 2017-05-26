# kubedb create

## Example

##### Help for create command

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

##### Create from file
```bash
$ kubedb create -f ./elastic.json

elastic "elasticsearch-demo" created
```

##### Create from stdin
```bash
$ cat ./elastic.json | kubedb create -f -

elastic "elasticsearch-demo" created
```

##### Create from folder
```bash
$ kubedb create -f resources -R

es "elasticsearch-demo" created
pg "postgres-demo" created
```
