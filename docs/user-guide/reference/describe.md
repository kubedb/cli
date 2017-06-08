### kubedb describe

```bash
$ kubedb describe --help

Show details of a specific resource or group of resources. This command joins many API calls together to form a detailed
description of a given resource or group of resources.Valid resource types include:

  * all
  * elastic
  * postgres
  * snapshot
  * dormantdatabase

Examples:
  # Describe a elastic
  kubedb describe elastics elasticsearch-demo

  # Describe a postgres
  kubedb describe pg/postgres-demo

  # Describe all dormantdatabases
  kubedb describe drmn

Options:
      --all-namespaces=false: If present, describe the requested object(s) across all namespaces. Namespace specified with --namespace will be ignored.
  -n, --namespace='default': Describe object(s) from this namespace.
  -l, --selector='': Selector (label query) to filter on, supports '=', '==', and '!='.
      --show-events=true: If true, display events related to the described object.

Usage:
  kubedb describe (TYPE [NAME_PREFIX] | TYPE/NAME) [flags] [options]

Use "kubedb describe options" for a list of global command-line options (applies to all commands).
```

`kubedb describe` command can describe multiple objects at a time.

This command will describe all objects of a single resource type if we do not provide object name.

To describe only a single object, we need to provide object name after resource type.

We can use `--all-namespaces` flag to describe all objects from every namespaces.

By default, `kubedb describe` shows events of an object. We need to provide flag `--show-events=false` to hide it.
