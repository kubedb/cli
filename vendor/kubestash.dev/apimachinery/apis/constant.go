/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apis

import "time"

const (
	KubeStashKey              = "kubestash.com"
	KubeStashApp              = "kubestash.com/app"
	KubeStashCleanupFinalizer = "kubestash.com/cleanup"
	KubeDBGroupName           = "kubedb.com"
	ElasticsearchGroupName    = "elasticsearch.kubedb.com"
)

const (
	KindStatefulSet           = "StatefulSet"
	KindDaemonSet             = "DaemonSet"
	KindDeployment            = "Deployment"
	KindClusterRole           = "ClusterRole"
	KindRole                  = "Role"
	KindPersistentVolumeClaim = "PersistentVolumeClaim"
	KindReplicaSet            = "ReplicaSet"
	KindReplicationController = "ReplicationController"
	KindJob                   = "Job"
	KindVolumeSnapshot        = "VolumeSnapshot"
	KindNamespace             = "Namespace"
	KindEmpty                 = ""
)

const (
	PrefixTrigger         = "trigger"
	PrefixInit            = "init"
	PrefixUpload          = "upload"
	PrefixCleanup         = "cleanup"
	PrefixRetentionPolicy = "retentionpolicy"
	PrefixPopulate        = "populate"
	PrefixPrime           = "prime"
)

const (
	KubeStashBackupComponent      = "kubestash-backup"
	KubeStashRestoreComponent     = "kubestash-restore"
	KubeStashInitializerComponent = "kubestash-initializer"
	KubeStashUploaderComponent    = "kubestash-uploader"
	KubeStashCleanerComponent     = "kubestash-cleaner"
	KubeStashHookComponent        = "kubestash-hook"
	KubeStashPopulatorComponent   = "kubestash-populator"
)

// Keys for offshoot labels
const (
	KubeStashInvokerName      = "kubestash.com/invoker-name"
	KubeStashInvokerNamespace = "kubestash.com/invoker-namespace"
	KubeStashInvokerKind      = "kubestash.com/invoker-kind"
	KubeStashSessionName      = "kubestash.com/session-name"
)

// Keys for snapshots labels
const (
	KubeStashRepoName        = "kubestash.com/repo-name"
	KubeStashAppRefKind      = "kubestash.com/app-ref-kind"
	KubeStashAppRefNamespace = "kubestash.com/app-ref-namespace"
	KubeStashAppRefName      = "kubestash.com/app-ref-name"
)

// Keys for structure logging
const (
	KeyTargetKind      = "target_kind"
	KeyTargetName      = "target_name"
	KeyTargetNamespace = "target_namespace"
	KeyReason          = "reason"
	KeyName            = "name"
)

// Keys for BackupBlueprint
const (
	VariablesKey       = "variables.kubestash.com"
	BackupBlueprintKey = "blueprint.kubestash.com"

	KeyBlueprintName      = BackupBlueprintKey + "/name"
	KeyBlueprintNamespace = BackupBlueprintKey + "/namespace"
	KeyBlueprintSessions  = BackupBlueprintKey + "/sessions"
)

// RBAC related
const (
	KubeStashBackupJobClusterRole          = "kubestash-backup-job"
	KubeStashRestoreJobClusterRole         = "kubestash-restore-job"
	KubeStashCronJobClusterRole            = "kubestash-cron-job"
	KubeStashBackendJobClusterRole         = "kubestash-backend-job"
	KubeStashStorageInitializerClusterRole = "kubestash-storage-initializer-job"
	KubeStashPopulatorJobClusterRole       = "kubestash-populator-job"
)

// Reconciliation related
const (
	RequeueTimeInterval = 10 * time.Second
	Requeue             = true
	DoNotRequeue        = false
)

// Local Network Volume Accessor related
const (
	KubeStashNetVolAccessor = "kubestash-netvol-accessor"
	TempDirVolumeName       = "kubestash-tmp-volume"
	TempDirMountPath        = "/kubestash-tmp"
	OperatorContainer       = "operator"
	KubeStashContainer      = "kubestash"
)

// Volume populator related constants
const (
	PopulatorKey                = "populator.kubestash.com"
	KeyPopulatedFrom            = PopulatorKey + "/populated-from"
	KeyAppName                  = PopulatorKey + "/app-name"
	KubeStashPopulatorContainer = "kubestash-populator"
)

const (
	ComponentPod            = "pod"
	ComponentDump           = "dump"
	ComponentWal            = "wal"
	ComponentManifest       = "manifest"
	ComponentVolumeSnapshot = "volumesnapshot"
	ComponentDashboard      = "dashboard"
)

const (
	EnvComponentName = "COMPONENT_NAME"
	KeyPodOrdinal    = "POD_ORDINAL"
	KeyPVCName       = "PVC_NAME"
	KeyDBVersion     = "DB_VERSION"
	KeyInterimVolume = "INTERIM_VOLUME"

	InterimVolumeName = "kubestash-interim-volume"
	OwnerKey          = ".metadata.controller"
	SnapshotVersionV1 = "v1"
	DirRepository     = "repository"
)
