# Create Database

we can create a database supported by **kubedb** using this CLI.

Lets create a Postgres database.

### kubedb create

`kubedb create` command will create an object in `default` namespace by default unless namespace is specified by input.

Following command will create a Postgres TPR as specified in `postgres.yaml`.

```bash
$ kubedb create -f postgres.yaml

postgres "postgres-demo" created
```

We can provide namespace as a flag `--namespace`.

```bash
$ kubedb create -f postgres.yaml --namespace=kube-system

postgres "postgres-demo" created
```

> Provided namespace should match with namespace specified in input file.

If input file do not specify namespace, object will be created in `default` namespace if not provided.


`kubedb create` command also considers `stdin` as input.

```bash
cat postgres.yaml | kubedb create -f -
```

##### Click [here](../reference/create.md) to get command details.
