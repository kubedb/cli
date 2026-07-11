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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

// +k8s:openapi-gen=false
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Migration struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Migration
	// +required
	Spec MigrationSpec `json:"spec"`

	// status defines the observed state of Migration
	// +optional
	Status MigrationStatus `json:"status,omitzero"`
}

// MigrationSpec defines the desired state of Migration
type MigrationSpec struct {
	// Source defines the source database configuration
	Source *Source `json:"source" protobuf:"bytes,1,opt,name=source"`

	// Target defines the target database configuration
	Target *Target `json:"target" protobuf:"bytes,2,opt,name=target"`

	// JobDefaults specifies default settings for migration jobs
	JobDefaults *JobDefaults `json:"jobDefaults,omitempty"`

	// JobTemplate specifies runtime configurations for the backup/restore Job
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
}

// JobDefaults defines default settings for migration jobs
type JobDefaults struct {
	// ImagePullPolicy specifies the image pull policy for the migrator Job
	// +kubebuilder:default=IfNotPresent
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// BackoffLimit specifies the number of retries before marking the job as failed
	// +kubebuilder:default=6
	// +optional
	BackoffLimit *int32 `json:"backoffLimit,omitempty"`

	// TTLSecondsAfterFinished specifies the TTL for completed jobs
	// +optional
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty"`

	// ActiveDeadlineSeconds specifies the duration in seconds relative to the startTime
	// that the job may be active before the system tries to terminate it
	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty"`
}

// MigrationStatus defines the observed state of Migration.
type MigrationStatus struct {
	// Phase represents the current phase of migration
	// +optional
	// +kubebuilder:default:=Pending
	Phase MigrationPhase `json:"phase,omitempty"`

	// Progress contains the current progress of migration
	// +optional
	Progress *Progress `json:"progress,omitempty"`

	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// MigrationPhase represents the current phase of migration
type MigrationPhase string

const (
	// MigrationPhasePending indicates the migration is pending
	MigrationPhasePending MigrationPhase = "Pending"
	// MigrationPhaseRunning indicates the migration is in progress
	MigrationPhaseRunning MigrationPhase = "Running"
	// MigrationPhaseSucceeded indicates the migration completed successfully
	MigrationPhaseSucceeded MigrationPhase = "Succeeded"
	// MigrationPhaseFailed indicates the migration failed
	MigrationPhaseFailed MigrationPhase = "Failed"
)

// Progress contains the current progress of migration
type Progress struct {
	// DBType indicates the type of database
	// +optional
	DBType string `json:"dbType,omitempty"`

	// Info contains the additional information about the current progress
	// +optional
	Info map[string]string `json:"info,omitempty"`
}

// MigrationList contains a list of Migration

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Migration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Migration{}, &MigrationList{})
}
