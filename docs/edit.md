> New to KubeDB? Please start [here](/docs/tutorial.md).

# Edit Database

`kubedb edit` command allows users to directly edit any KubeDB object. It will open the editor defined by _KUBEDB_EDITOR_, or _EDITOR_ environment variables, or fall back to `nano`.

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
Various fields of a KubeDb object can't be edited using `edit` command. The following fields are restricted from updates for all KubeDB objects:
* _apiVersion_
* _kind_
* _metadata.name_
* _metadata.namespace_
* _status_


If StatefulSet exists for a database, following fields can't be modified as well.

Postgres:
* _spec.version_
* _spec.storage_
* _spec.databaseSecret_
* _spec.nodeSelector_
* _spec.init_

Elastic:
* _spec.version_
* _spec.storage_
* _spec.nodeSelector_
* _spec.init_

For DormantDatabase, _spec.origin_ can't be edited using `kbuedb edit`

To learn about various options of `edit` command, please visit [here](/docs/reference/kubedb_edit.md).
