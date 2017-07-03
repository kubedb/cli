> New to KubeDB? Please start [here](/docs/tutorial.md).

# Create Database

`kubedb create` creates a database tpr in `default` namespace by default. Following command will create a Postgres TPR as specified in `postgres.yaml`.

```sh
$ kubedb create -f ./docs/examples/postgres/postgres.yaml

postgres "postgres-demo" created
```

You can provide namespace as a flag `--namespace`. Provided namespace should match with namespace specified in input file.

```sh
$ kubedb create -f postgres.yaml --namespace=kube-system

postgres "postgres-demo" created
```

`kubedb create` command also considers `stdin` as input.

```sh
cat postgres.yaml | kubedb create -f -
```

To learn about various options of `create` command, please visit [here](/docs/reference/kubedb_create.md).
