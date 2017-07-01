> New to KubeDB? Please start [here](/docs/tutorial.md).

# kubedb edit

`edit` command allows users to directly edit any KubeDB object. It will open the editor defined by _KUBEDB_EDITOR_, or _EDITOR_ environment variables, or fall back to `nano`.

Lets edit an existing running Postgres database to setup [Scheduled Backup](/docs/backup.md). The following command will open Postgres `postgres-demo` in editor.

```bash
$ kubedb edit pg postgres-demo

# Add following under Spec to configure periodic backups
#  backupSchedule:
#    cronExpression: "@every 6h"
#    bucketName: "bucket-name"
#    storageSecret:
#      secretName: "secret-name"

postgres "postgres-demo" edited
```

## Edit restrictions
`kubedb edit` will allow us to edit only supported fields.

Following fields are restricted to be modified for all supported TPR objects using `kbuedb edit`

* _apiVersion_
* _kind_
* **metadata**._name_
* **metadata**._namespace_
* _status_


If StatefulSet exists for a database, following fields can't be modified as well.

Postgres:
* **spec**._version_
* **spec**._storage_
* **spec**._databaseSecret_
* **spec**._nodeSelector_
* **spec**._init_

Elastic:
* **spec**._version_
* **spec**._storage_
* **spec**._nodeSelector_
* **spec**._init_

For DormantDatabase, **spec**._origin_ can't be edited using `kbuedb edit`

To learn about various options of `get` command, please visit [here](/docs/reference/kubedb_edit.md).
