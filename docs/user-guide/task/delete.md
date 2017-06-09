# Delete TPR object

we can delete supported TPR objects using this CLI.

Lets delete a Postgres database.

### kubedb delete

`kubedb delete` command will delete an object in `default` namespace by default unless namespace is provided.

Following command will delete a Postgres `postgres-dev` in default namespace

```bash
$ kubedb delete postgres postgres-dev

postgres "postgres-dev" deleted
```

We can use `postgres.yaml` file to delete objects.

```bash
$ kubedb delete -f postgres.yaml

postgres "postgres-dev" deleted
```

This will delete a postgres using the type and name specified in postgres.yaml

`kubedb delete` command also takes input from `stdin`.

```bash
cat postgres.yaml | kubedb delete -f -
```

Also we can also filter using `--selector` flag.

```bash
$ kubedb delete postgres -l postgres.kubedb.com/name=postgres-demo
```

This will delete postgres with label postgres.kubedb.com/name=postgres-demo.

##### Click [here](../reference/delete.md) to get command details.
