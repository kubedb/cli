> New to KubeDB? Please start [here](/docs/tutorial.md).

# Delete Database

`kubedb delete` command will delete an object in `default` namespace by default unless namespace is provided. The following command will delete a Postgres `postgres-dev` in default namespace

```sh
$ kubedb delete postgres postgres-dev

postgres "postgres-dev" deleted
```

You can also use YAML files to delete objects. The following command will delete a postgres using the type and name specified in `postgres.yaml`.

```sh
$ kubedb delete -f postgres.yaml

postgres "postgres-dev" deleted
```

`kubedb delete` command also takes input from `stdin`.

```sh
cat postgres.yaml | kubedb delete -f -
```

To delete database with matching labels, use `--selector` flag. The following command will delete postgres with label `postgres.kubedb.com/name=postgres-demo`.
```sh
$ kubedb delete postgres -l postgres.kubedb.com/name=postgres-demo
```

To learn about various options of `delete` command, please visit [here](/docs/reference/kubedb_delete.md).
