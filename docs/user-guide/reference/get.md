### kubedb get

```bash
$ kubedb get --help

Display one or many resources.

Valid resource types include:

  * all
  * elastic
  * postgres
  * snapshot
  * dormantdatabase

Examples:
  # List all elastic in ps output format.
  kubedb get elastics

  # List all elastic in ps output format with more information (such as version).
  kubedb get elastics -o wide

  # List a single postgres with specified NAME in ps output format.
  kubedb get postgres database

  # List a single snapshot in JSON output format.
  kubedb get -o json snapshot snapshot-xyz

  # List all postgreses and elastics together in ps output format.
  kubedb get postgreses,elastics

  # List one or more resources by their type and names.
  kubedb get elastic/es-db postgres/pg-db

Options:
      --all-namespaces=false: If present, list the requested object(s) across all namespaces. Namespace specified with
--namespace will be ignored.
  -n, --namespace='default': List the requested object(s) from this namespace.
  -o, --output='': Output format. One of: json|yaml|wide|name.
  -l, --selector='': Selector (label query) to filter on, supports '=', '==', and '!='.
  -a, --show-all=false: When printing, show all resources (default hide terminated pods.)
      --show-kind=false: If present, list the resource type for the requested object(s).
      --show-labels=false: When printing, show all labels as the last column (default hide labels column)

Usage:
  kubedb get [flags] [options]

Use "kubedb get options" for a list of global command-line options (applies to all commands).
```

To list requested object(s) across all namespaces, use `--all-namespaces=true`

Flag `--selector` can be used to filter against labels.

For output format, we can pass `--output` flag.
* `--output=yaml` to print requested objects in YAML format
* `--output=json` to print requested objects in JSON format
* `--output=wide` to print requested objects with additional information
* `--output=name` to print requested objects' name only

We can also print objects by specifying names after resource type.

See examples in `kubedb get --help`
