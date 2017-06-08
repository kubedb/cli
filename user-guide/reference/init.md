### kubedb init

```bash
$ kubedb init --help

Create or upgrade unified operator for kubedb databases.

Examples:
  # Create operator with version canary.
  kubedb init --version=0.1.0
  
  # Upgrade operator to use another version.
  kubedb init --version=0.1.0 --upgrade

Options:
  -n, --namespace='default': Namespace name. Operator will be deployed in this namespace.
      --upgrade=false: If present, Upgrade operator to use provided version
      --version='0.1.0': Operator version

Usage:
  kubedb init [flags] [options]

Use "kubedb init options" for a list of global command-line options (applies to all commands).
```

We can provide operator version using `--version` flag.

Also we can say in which namespace we want to deploy this operator providing `--namespace` flag.

And this same command can be used to upgrade operator to another version. We need to pass `--upgrade` flag for that.