## kubedb init

Create or upgrade KubeDB operator

### Synopsis


Install or upgrade KubeDB operator.

```
kubedb init [flags]
```

### Examples

```
  # Install latest released operator.
  kubedb init
  
  # Upgrade operator to use another version.
  kubedb init --version=0.7.1 --upgrade
```

### Options

```
  -h, --help                        help for init
      --operator-namespace string   Name of namespace where operator will be deployed. (default "kube-system")
      --rbac                        If true, uses RBAC with operator and database objects
      --upgrade                     If present, Upgrade operator to use provided version
      --version string              Operator version (default "0.7.1")
```

### Options inherited from parent commands

```
      --analytics             Send analytical events to Google Analytics (default true)
      --kube-context string   name of the kubeconfig context to use
```

### SEE ALSO
* [kubedb](kubedb.md)	 - Command line interface for KubeDB


