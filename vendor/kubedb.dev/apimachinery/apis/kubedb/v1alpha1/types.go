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

package v1alpha1

import (
	core "k8s.io/api/core/v1"
)

type InitSpec struct {
	// Initializer object in same namespace of a KubeDB object.
	// ref: https://github.com/stashed/stash/blob/09af5d319bb5be889186965afb04045781d6f926/apis/stash/v1beta1/restore_session_types.go#L22
	Initializer *core.TypedLocalObjectReference `json:"initializer,omitempty" protobuf:"bytes,1,opt,name=initializer"`
	Script      *ScriptSourceSpec               `json:"script,omitempty" protobuf:"bytes,2,opt,name=script"`
	PostgresWAL *PostgresWALSourceSpec          `json:"postgresWAL,omitempty" protobuf:"bytes,3,opt,name=postgresWAL"`
}

type ScriptSourceSpec struct {
	ScriptPath        string `json:"scriptPath,omitempty" protobuf:"bytes,1,opt,name=scriptPath"`
	core.VolumeSource `json:",inline,omitempty" protobuf:"bytes,2,opt,name=volumeSource"`
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

// +kubebuilder:validation:Enum=Running;Creating;Initializing;Paused;Halted;Failed
type DatabasePhase string

const (
	// used for Databases that are currently running
	DatabasePhaseRunning DatabasePhase = "Running"
	// used for Databases that are currently creating
	DatabasePhaseCreating DatabasePhase = "Creating"
	// used for Databases that are currently initializing
	DatabasePhaseInitializing DatabasePhase = "Initializing"
	// used for Databases that are paused
	DatabasePhasePaused DatabasePhase = "Paused"
	// used for Databases that are halted
	DatabasePhaseHalted DatabasePhase = "Halted"
	// used for Databases that are failed
	DatabasePhaseFailed DatabasePhase = "Failed"
)

// +kubebuilder:validation:Enum=Durable;Ephemeral
type StorageType string

const (
	// default storage type and requires spec.storage to be configured
	StorageTypeDurable StorageType = "Durable"
	// Uses emptyDir as storage
	StorageTypeEphemeral StorageType = "Ephemeral"
)

// +kubebuilder:validation:Enum=Halt;Delete;WipeOut;DoNotTerminate
type TerminationPolicy string

const (
	// Pauses database into a DormantDatabase
	// Deprecated: Use spec.halted = true
	TerminationPolicyPause TerminationPolicy = "Pause"
	// Deletes database pods, service but leave the PVCs and stash backup data intact.
	TerminationPolicyHalt TerminationPolicy = "Halt"
	// Deletes database pods, service, pvcs but leave the stash backup data intact.
	TerminationPolicyDelete TerminationPolicy = "Delete"
	// Deletes database pods, service, pvcs and stash backup data.
	TerminationPolicyWipeOut TerminationPolicy = "WipeOut"
	// Rejects attempt to delete database using ValidationWebhook. This replaces spec.doNotPause = true
	TerminationPolicyDoNotTerminate TerminationPolicy = "DoNotTerminate"
)
