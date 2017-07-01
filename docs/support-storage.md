> New to KubeDB? Please start [here](/docs/tutorial.md).

### Add Storage

**T**o add PersistentVolume support, we need to add following StorageSpec in `spec`

```yaml
spec:
  storage:
    class: "gp2"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: "10Gi"
```

Here we must have to add following storage information in `spec.storage`:

* `class:` StorageClass (`kubectl get storageclasses`)
* `resources:` ResourceRequirements for PersistentVolumeClaimSpec

**A**s we have used storage information in our database yaml, StatefulSet will be created with PersistentVolumeClaim.
