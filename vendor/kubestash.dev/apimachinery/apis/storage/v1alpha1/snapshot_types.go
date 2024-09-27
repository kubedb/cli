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
	"kubestash.dev/apimachinery/apis"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

type BackupType string

const (
	ResourceKindSnapshot     = "Snapshot"
	ResourceSingularSnapshot = "snapshot"
	ResourcePluralSnapshot   = "snapshots"

	BackupTypeFull        BackupType = "FullBackup"
	BackupTypeIncremental BackupType = "IncrementalBackup"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=snapshots,singular=snapshot,categories={kubestash,appscode}
// +kubebuilder:printcolumn:name="Repository",type="string",JSONPath=".spec.repository"
// +kubebuilder:printcolumn:name="Session",type="string",JSONPath=".spec.session"
// +kubebuilder:printcolumn:name="Snapshot-Time",type="string",JSONPath=".status.snapshotTime"
// +kubebuilder:printcolumn:name="Deletion-Policy",type="string",JSONPath=".spec.deletionPolicy"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Snapshot represents the state of a backup run to a particular Repository.
// Multiple components of the same target may be backed up in the same Snapshot.
// This is a namespaced CRD. It should be in the same namespace as the respective Repository.
// KubeStash operator is responsible for creating Snapshot CR.
// Snapshot is not supposed to be created/edited by the end user.
type Snapshot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnapshotSpec   `json:"spec,omitempty"`
	Status SnapshotStatus `json:"status,omitempty"`
}

// SnapshotSpec specifies the information regarding the application that is being backed up,
// the Repository where the backed up data is being stored, and the session which is
// responsible for this snapshot etc.
type SnapshotSpec struct {
	// SnapshotID represents a "Universally Unique Lexicographically Sortable Identifier" (ULID) for the Snapshot.
	// For more details about ULID, please see: https://github.com/oklog/ulid
	// +optional
	SnapshotID string `json:"snapshotID,omitempty"`

	// Type specifies whether this snapshot represents a full or incremental backup
	Type BackupType `json:"type,omitempty"`

	// Repository specifies the name of the Repository where this Snapshot is being stored.
	Repository string `json:"repository,omitempty"`

	// Session specifies the name of the session which is responsible for this Snapshot
	Session string `json:"session,omitempty"`

	// BackupSession represents the name of the respective BackupSession which is responsible for this Snapshot.
	// +optional
	BackupSession string `json:"backupSession,omitempty"`

	// Version denotes the respective data organization structure inside the Repository
	Version string `json:"version,omitempty"`

	// AppRef specifies the reference of the application that has been backed up in this Snapshot.
	AppRef kmapi.TypedObjectReference `json:"appRef,omitempty"`

	// DeletionPolicy specifies what to do when you delete a Snapshot CR.
	// The valid values are:
	// - "Delete": This will delete just the Snapshot CR from the cluster but keep the backed up data in the remote backend. This is the default behavior.
	// - "WipeOut": This will delete the Snapshot CR as well as the backed up data from the backend.
	// +kubebuilder:default=Delete
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// Paused specifies whether the Snapshot is paused or not. If the Snapshot is paused,
	// KubeStash will not process any further event for the Snapshot.
	// +optional
	Paused bool `json:"paused,omitempty"`
}

// SnapshotStatus defines the observed state of Snapshot
type SnapshotStatus struct {
	// Phase represents the backup state of this Snapshot
	// +optional
	Phase SnapshotPhase `json:"phase,omitempty"`

	// VerificationStatus specifies whether this Snapshot has been verified or not
	// +optional
	VerificationStatus VerificationStatus `json:"verificationStatus,omitempty"`

	// SnapshotTime represents the timestamp when this Snapshot was taken.
	// +optional
	SnapshotTime *metav1.Time `json:"snapshotTime,omitempty"`

	// LastUpdateTime specifies the timestamp when this Snapshot was last updated.
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`

	// Size represents the size of the Snapshot
	// +optional
	Size string `json:"size,omitempty"`

	// Integrity represents whether the Snapshot data has been corrupted or not
	// +optional
	Integrity *bool `json:"integrity,omitempty"`

	// Conditions represents list of conditions regarding this Snapshot
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`

	// TotalComponents represents the number of total components for this Snapshot
	// +optional
	TotalComponents int32 `json:"totalComponents,omitempty"`

	// Components represents the backup information of the individual components of this Snapshot
	// +optional
	// +mapType=granular
	Components map[string]Component `json:"components,omitempty"`
}

// SnapshotPhase represent the overall progress of this Snapshot
// +kubebuilder:validation:Enum=Pending;Running;Succeeded;Failed
type SnapshotPhase string

const (
	SnapshotPending   SnapshotPhase = "Pending"
	SnapshotRunning   SnapshotPhase = "Running"
	SnapshotSucceeded SnapshotPhase = "Succeeded"
	SnapshotFailed    SnapshotPhase = "Failed"
)

// VerificationStatus represents whether the Snapshot has been verified or not.
// +kubebuilder:validation:Enum=Verified;NotVerified;VerificationFailed
type VerificationStatus string

const (
	SnapshotVerified           VerificationStatus = "Verified"
	SnapshotNotVerified        VerificationStatus = "NotVerified"
	SnapshotVerificationFailed VerificationStatus = "VerificationFailed"
)

// Component represents the backup information of individual components
type Component struct {
	// Path specifies the path inside the Repository where the backed up data for this component has been stored.
	// This path is relative to Repository path.
	Path string `json:"path,omitempty"`

	// Phase represents the backup phase of the component
	// +optional
	Phase ComponentPhase `json:"phase,omitempty"`

	// Size represents the size of the restic repository for this component
	// +optional
	Size string `json:"size,omitempty"`

	// Duration specifies the total time taken to complete the backup process for this component
	// +optional
	Duration string `json:"duration,omitempty"`

	// Integrity represents the result of the restic repository integrity check for this component
	// +optional
	Integrity *bool `json:"integrity,omitempty"`

	// Error specifies the reason in case of backup failure for the component
	// +optional
	Error string `json:"error,omitempty"`

	// Driver specifies the name of the tool that has been used to upload the underlying backed up data
	Driver apis.Driver `json:"driver,omitempty"`

	// ResticStats specifies the "Restic" driver specific information
	// +optional
	ResticStats []ResticStats `json:"resticStats,omitempty"`

	// WalGStats specifies the "WalG" driver specific information
	// +optional
	WalGStats *WalGStats `json:"walGStats,omitempty"`

	// VolumeSnapshotterStats specifies the "VolumeSnapshotter" driver specific information
	// +optional
	VolumeSnapshotterStats []VolumeSnapshotterStats `json:"volumeSnapshotterStats,omitempty"`
	// WalSegments specifies a list of wall segment for individual component
	WalSegments []WalSegment `json:"walSegments,omitempty"`
}

// ComponentPhase represents the backup phase of the individual component.
// +kubebuilder:validation:Enum=Pending;Running;Succeeded;Failed
type ComponentPhase string

const (
	ComponentPhasePending   ComponentPhase = "Pending"
	ComponentPhaseRunning   ComponentPhase = "Running"
	ComponentPhaseSucceeded ComponentPhase = "Succeeded"
	ComponentPhaseFailed    ComponentPhase = "Failed"
)

// ResticStats specifies the "Restic" driver specific information
type ResticStats struct {
	// Id represents the restic snapshot id
	Id string `json:"id,omitempty"`

	// Uploaded specifies the amount of data that has been uploaded in the restic snapshot.
	// +optional
	Uploaded string `json:"uploaded,omitempty"`

	// HostPath represents the backup path for which restic snapshot is taken.
	// +optional
	HostPath string `json:"hostPath,omitempty"`

	// Size represents the restic snapshot size
	// +optional
	Size string `json:"size,omitempty"`
}

// VolumeSnapshotterStats specifies the "VolumeSnapshotter" driver specific information
type VolumeSnapshotterStats struct {

	// PVCName represents the backup PVC name for which volumeSnapshot is created.
	// +optional
	PVCName string `json:"pvcName,omitempty"`

	// HostPath represents the corresponding path of PVC for which volumeSnapshot is created.
	// +optional
	HostPath string `json:"hostPath,omitempty"`

	// VolumeSnapshotName represents the name of created volumeSnapshot.
	// +optional
	VolumeSnapshotName string `json:"volumeSnapshotName,omitempty"`

	// VolumeSnapshotTime indicates the timestamp at which the volumeSnapshot was created.
	VolumeSnapshotTime *metav1.Time `json:"volumeSnapshotTime,omitempty"`
}

// WalGStats specifies the information specific to the "WalG" driver.
type WalGStats struct {
	// Id represents the WalG snapshot ID.
	Id string `json:"id,omitempty"`

	// Databases represents the list of target backup databases.
	// +optional
	Databases []string `json:"databases,omitempty"`

	// StartTime represents the WalG backup start time.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// StopTime represents the WalG backup stop time.
	// +optional
	StopTime *metav1.Time `json:"stopTime,omitempty"`
}

// WalSegment specifies the "WalG" driver specific information
type WalSegment struct {
	Start *metav1.Time `json:"start,omitempty"`
	End   *metav1.Time `json:"end,omitempty"`
}

const (
	TypeSnapshotMetadataUploaded               = "SnapshotMetadataUploaded"
	ReasonFailedToUploadSnapshotMetadata       = "FailedToUploadSnapshotMetadata"
	ReasonSuccessfullyUploadedSnapshotMetadata = "SuccessfullyUploadedSnapshotMetadata"

	TypeRecentSnapshotListUpdated               = "RecentSnapshotListUpdated"
	ReasonFailedToUpdateRecentSnapshotList      = "FailedToUpdateRecentSnapshotList"
	ReasonSuccessfullyUpdatedRecentSnapshotList = "SuccessfullyUpdatedRecentSnapshotList"

	TypeBackupIncomplete                           = "BackupIncomplete"
	ReasonBackupExecutorTerminatedBeforeCompletion = "BackupExecutorTerminatedBeforeCompletion"
)

//+kubebuilder:object:root=true

// SnapshotList contains a list of Snapshot
type SnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Snapshot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Snapshot{}, &SnapshotList{})
}
