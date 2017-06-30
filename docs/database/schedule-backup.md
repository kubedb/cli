### Schedule Backup

**T**o schedule backup, we need to add following BackupScheduleSpec in `spec`

```yaml
spec:
  backupSchedule:
    cronExpression: "@every 6h"
    bucketName: "bucket-for-snapshot"
    storageSecret:
      secretName: "secret-for-bucket"
```

> **Note:** storage can also be used here

When database TPR object is running,
operator immediately takes a backup to validate this information.

And after successful backup, operator will set a scheduler to take backup `every 6h`.

See backup process in [details](backup.md).

