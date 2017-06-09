# Edit TPR object

The edit command allows us to directly edit any TPR object we can retrieve via this CLI.
It will open the editor defined by _KUBEDB_EDITOR_, or _EDITOR_ environment variables, or fall back to `nano`.


Lets edit our existing running Postgres database to set Scheduled Backup.

### kubedb edit

Following command will open Postgres `postgres-demo` in editor.

```bash
$ kubedb edit pg postgres-demo

# Add following in Spec to schedule backup
#  backupSchedule:
#    cronExpression: "@every 6h"
#    bucketName: "bucket-name"
#    storageSecret:
#      secretName: "secret-name"

postgres "postgres-demo" edited
```

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

##### Click [here](../reference/edit.md) to get command details.
