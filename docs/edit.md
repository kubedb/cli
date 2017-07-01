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
Various fields of a KubeDb object can't be edited using `edit` command. The following fields are restricted from updates for all KubeDB objects:
* apiVersion
* kind
* .metadata.name
* .metadata.namespace
* status


If StatefulSet exists for a database, following fields can't be modified as well.

Postgres:
* .spec.version
* .spec.storage
* .spec.databaseSecret
* .spec.nodeSelector
* .spec.init

Elastic:
* .spec.version
* .spec.storage
* .spec.nodeSelector
* .spec.init

For DormantDatabase, .spec.origin can't be edited using `kbuedb edit`

To learn about various options of `edit` command, please visit [here](/docs/reference/kubedb_edit.md).
