# kubedb init

## Example

##### Help for init command

```bash
$ kubedb init --help

Create or upgrade unified operator for k8sdb databases.

Examples:
  # Create operator with version canary.
  kubedb init --version=canary
  
  # Upgrade operator to use another version.
  kubedb init --version=canary --upgrade

Options:
  -n, --namespace='default': Namespace name. Operator will be deployed in this namespace.
      --upgrade=false: If present, Upgrade operator to use provided version
      --version='': Operator version

Usage:
  kubedb init [flags] [options]

Use "kubedb init options" for a list of global command-line options (applies to all commands).
```

##### Create
```bash
$ kubedb init --version=canary

Successfully created operator deployment.
```

##### Upgrade
```bash
$ kubedb init --version=canary --upgrade

Successfully upgraded operator deployment.
```
