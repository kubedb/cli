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

const (
	ResourceKindMigrator     = "Migrator"
	ResourceSingularMigrator = "migrator"
	ResourcePluralMigrator   = "migrators"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=migrators,singular=migrator,shortName=mgtr,categories={kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="DBType",type="string",JSONPath=".status.progress.dbType"
// +kubebuilder:printcolumn:name="Stage",type="string",JSONPath=".status.progress.info.Stage"
// +kubebuilder:printcolumn:name="Lag",type="string",JSONPath=".status.progress.info.Lag"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Migrator struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Migrator
	// +required
	Spec MigratorSpec `json:"spec"`

	// status defines the observed state of Migrator
	// +optional
	Status MigratorStatus `json:"status,omitzero"`
}

// MigratorSpec defines the desired state of Migrator
type MigratorSpec struct {
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

// MigratorStatus defines the observed state of Migrator.
type MigratorStatus struct {
	// Phase represents the current phase of migration
	// +optional
	// +kubebuilder:default:=Pending
	Phase MigratorPhase `json:"phase,omitempty"`

	// Progress contains the current progress of migration
	// +optional
	Progress *Progress `json:"progress,omitempty"`

	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// MigratorPhase represents the current phase of migration
type MigratorPhase string

const (
	// MigratorPhasePending indicates the migration is pending
	MigratorPhasePending MigratorPhase = "Pending"
	// MigratorPhaseRunning indicates the migration is in progress
	MigratorPhaseRunning MigratorPhase = "Running"
	// MigratorPhaseSucceeded indicates the migration completed successfully
	MigratorPhaseSucceeded MigratorPhase = "Succeeded"
	// MigratorPhaseFailed indicates the migration failed
	MigratorPhaseFailed MigratorPhase = "Failed"
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

// MigratorList contains a list of Migrator

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MigratorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Migrator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Migrator{}, &MigratorList{})
}
