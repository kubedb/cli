/*
Copyright The KubeDB Authors.

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

package v1alpha1

import (
	core "k8s.io/api/core/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

type InitSpec struct {
	ScriptSource *ScriptSourceSpec `json:"scriptSource,omitempty" protobuf:"bytes,1,opt,name=scriptSource"`
	// Deprecated
	SnapshotSource *SnapshotSourceSpec    `json:"snapshotSource,omitempty" protobuf:"bytes,2,opt,name=snapshotSource"`
	PostgresWAL    *PostgresWALSourceSpec `json:"postgresWAL,omitempty" protobuf:"bytes,3,opt,name=postgresWAL"`
	// Name of stash restoreSession in same namespace of kubedb object.
	// ref: https://github.com/stashed/stash/blob/09af5d319bb5be889186965afb04045781d6f926/apis/stash/v1beta1/restore_session_types.go#L22
	StashRestoreSession *core.LocalObjectReference `json:"stashRestoreSession,omitempty" protobuf:"bytes,4,opt,name=stashRestoreSession"`
}

type ScriptSourceSpec struct {
	ScriptPath        string `json:"scriptPath,omitempty" protobuf:"bytes,1,opt,name=scriptPath"`
	core.VolumeSource `json:",inline,omitempty" protobuf:"bytes,2,opt,name=volumeSource"`
}

type SnapshotSourceSpec struct {
	Namespace string `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	Name      string `json:"name" protobuf:"bytes,2,opt,name=name"`
	// Arguments to the restore job
	Args []string `json:"args,omitempty" protobuf:"bytes,3,rep,name=args"`
}

type BackupScheduleSpec struct {
	CronExpression string `json:"cronExpression,omitempty" protobuf:"bytes,1,opt,name=cronExpression"`

	// Snapshot Spec
	store.Backend `json:",inline" protobuf:"bytes,2,opt,name=backend"`

	// StorageType can be durable or ephemeral.
	// If not given, database storage type will be used.
	// +optional
	StorageType *StorageType `json:"storageType,omitempty" protobuf:"bytes,3,opt,name=storageType,casttype=StorageType"`

	// PodTemplate is an optional configuration for pods used to take database snapshots
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,4,opt,name=podTemplate"`

	// PodVolumeClaimSpec is used to specify temporary storage for backup/restore Job.
	// If not given, database's PvcSpec will be used.
	// If storageType is durable, then a PVC will be created using this PVCSpec.
	// If storageType is ephemeral, then an empty directory will be created of size PvcSpec.Resources.Requests[core.ResourceStorage].
	// +optional
	PodVolumeClaimSpec *core.PersistentVolumeClaimSpec `json:"podVolumeClaimSpec,omitempty" protobuf:"bytes,5,opt,name=podVolumeClaimSpec"`
}

// LeaderElectionConfig contains essential attributes of leader election.
// ref: https://github.com/kubernetes/client-go/blob/6134db91200ea474868bc6775e62cc294a74c6c6/tools/leaderelection/leaderelection.go#L105-L114
type LeaderElectionConfig struct {
	// LeaseDuration is the duration in second that non-leader candidates will
	// wait to force acquire leadership. This is measured against time of
	// last observed ack. Default 15
	LeaseDurationSeconds int32 `json:"leaseDurationSeconds" protobuf:"varint,1,opt,name=leaseDurationSeconds"`
	// RenewDeadline is the duration in second that the acting master will retry
	// refreshing leadership before giving up. Normally, LeaseDuration * 2 / 3.
	// Default 10
	RenewDeadlineSeconds int32 `json:"renewDeadlineSeconds" protobuf:"varint,2,opt,name=renewDeadlineSeconds"`
	// RetryPeriod is the duration in second the LeaderElector clients should wait
	// between tries of actions. Normally, LeaseDuration / 3.
	// Default 2
	RetryPeriodSeconds int32 `json:"retryPeriodSeconds" protobuf:"varint,3,opt,name=retryPeriodSeconds"`
}

type DatabasePhase string

const (
	// used for Databases that are currently running
	DatabasePhaseRunning DatabasePhase = "Running"
	// used for Databases that are currently creating
	DatabasePhaseCreating DatabasePhase = "Creating"
	// used for Databases that are currently initializing
	DatabasePhaseInitializing DatabasePhase = "Initializing"
	// used for Databases that are Failed
	DatabasePhaseFailed DatabasePhase = "Failed"
)

type StorageType string

const (
	// default storage type and requires spec.storage to be configured
	StorageTypeDurable StorageType = "Durable"
	// Uses emptyDir as storage
	StorageTypeEphemeral StorageType = "Ephemeral"
)

type TerminationPolicy string

const (
	// Pauses database into a DormantDatabase
	TerminationPolicyPause TerminationPolicy = "Pause"
	// Deletes database pods, service, pvcs but leave the snapshot data intact. This will not create a DormantDatabase.
	TerminationPolicyDelete TerminationPolicy = "Delete"
	// Deletes database pods, service, pvcs and snapshot data. This will not create a DormantDatabase.
	TerminationPolicyWipeOut TerminationPolicy = "WipeOut"
	// Rejects attempt to delete database using ValidationWebhook. This replaces spec.doNotPause = true
	TerminationPolicyDoNotTerminate TerminationPolicy = "DoNotTerminate"
)
