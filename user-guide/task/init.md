# Unified operator

First of all, we need an unified [operator](https://github.com/k8sdb/operator) to handle supported TPR.

We can deploy this operator using `kubedb` CLI.

### kubedb init

The `kubedb init` command will start an unified operator for kubedb databases. This command can also be used to upgrade version of operator.

Following command will create a deployment with image `kubedb/operator:0.1.0` and a service in `default` namespace

```bash
$ kubedb init --namespace='default' --version='0.1.0'

Successfully created operator deployment.
Successfully created operator service.
```

Any existing operator can also be upgraded using this command.

```bash
$ kubedb init --version='0.2.0' --upgrade

Successfully upgraded operator deployment.
```

##### Click [here](../reference/init.md) to get command details.