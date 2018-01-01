---
title: Kubedb Init
menu:
  docs_0.8.0:
    identifier: kubedb-init
    name: Kubedb Init
    parent: reference
menu_name: docs_0.8.0
section_menu_id: reference
---
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
  kubedb init --version=0.8.0 --upgrade
```

### Options

```
      --address string              Address to listen on for web interface and telemetry. (default ":8080")
      --docker-registry string      User provided docker repository (default "kubedb")
      --governing-service string    Governing service for database statefulset (default "kubedb")
  -h, --help                        help for init
      --operator-namespace string   Name of namespace where operator will be deployed. (default "kube-system")
      --rbac                        If true, uses RBAC with operator and database objects
      --upgrade                     If present, Upgrade operator to use provided version
      --version string              Operator version (default "0.8.0")
```

### Options inherited from parent commands

```
      --analytics             Send analytical events to Google Analytics (default true)
      --kube-context string   name of the kubeconfig context to use
```

### SEE ALSO

* [kubedb](/docs/reference/kubedb.md)	 - Command line interface for KubeDB


