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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindRepository     = "Repository"
	ResourceSingularRepository = "repository"
	ResourcePluralRepository   = "repositories"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=repositories,singular=repository,shortName=repo,categories={kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Integrity",type="boolean",JSONPath=".status.integrity"
// +kubebuilder:printcolumn:name="Snapshot-Count",type="integer",JSONPath=".status.snapshotCount"
// +kubebuilder:printcolumn:name="Size",type="string",JSONPath=".status.size"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Last-Successful-Backup",type="date",format="date-time",JSONPath=".status.lastBackupTime"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Repository specifies the information about the targeted application that has been backed up
// and the BackupStorage where the backed up data is being stored. It also holds a list of recent
// Snapshots that have been taken in this Repository.
// Repository is a namespaced object. It must be in the same namespace as the targeted application.
type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepositorySpec   `json:"spec,omitempty"`
	Status RepositoryStatus `json:"status,omitempty"`
}

// RepositorySpec specifies the application reference and the BackupStorage reference.It also specifies
// what should be the behavior when a Repository CR is deleted from the cluster.
type RepositorySpec struct {
	// AppRef refers to the application that is being backed up in this Repository.
	AppRef kmapi.TypedObjectReference `json:"appRef,omitempty"`

	// StorageRef refers to the BackupStorage CR which contain the backend information where the backed
	// up data will be stored. The BackupStorage could be in a different namespace. However, the Repository
	// namespace must be allowed to use the BackupStorage.
	StorageRef kmapi.ObjectReference `json:"storageRef,omitempty"`

	// Path represents the directory inside the BackupStorage where this Repository is storing its data
	// This path is relative to the path of BackupStorage.
	Path string `json:"path,omitempty"`

	// DeletionPolicy specifies what to do when you delete a Repository CR.
	// The valid values are:
	// "Delete": This will delete the respective Snapshot CRs from the cluster but keep the backed up data in the remote backend. This is the default behavior.
	// "WipeOut": This will delete the respective Snapshot CRs as well as the backed up data from the backend.
	// +kubebuilder:default=Delete
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// EncryptionSecret refers to the Secret containing the encryption key which will be used to encode/decode the backed up data.
	// You can refer to a Secret of a different namespace.
	// If you don't provide the namespace field, KubeStash will look for the Secret in the same namespace as the BackupConfiguration / BackupBatch.
	EncryptionSecret *kmapi.ObjectReference `json:"encryptionSecret,omitempty"`

	// Paused specifies whether the Repository is paused or not. If the Repository is paused,
	// KubeStash will not process any further event for the Repository.
	// +optional
	Paused bool `json:"paused,omitempty"`
}

// RepositoryStatus defines the observed state of Repository
type RepositoryStatus struct {
	// Phase represents the current state of the Repository.
	// +optional
	Phase RepositoryPhase `json:"phase,omitempty"`

	// LastBackupTime specifies the timestamp when the last successful backup has been taken
	// +optional
	LastBackupTime *metav1.Time `json:"lastBackupTime,omitempty"`

	// Integrity specifies whether the backed up data of this Repository has been corrupted or not
	// +optional
	Integrity *bool `json:"integrity,omitempty"`

	// SnapshotCount specifies the number of current Snapshots stored in this Repository
	// +optional
	SnapshotCount *int32 `json:"snapshotCount,omitempty"`

	// Size specifies the amount of backed up data stored in the Repository
	// +optional
	Size string `json:"size,omitempty"`

	// RecentSnapshots holds a list of recent Snapshot information that has been taken in this Repository
	// +optional
	RecentSnapshots []SnapshotInfo `json:"recentSnapshots,omitempty"`

	// Conditions represents list of conditions regarding this Repository
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`

	// ComponentPaths represents list of component paths in this Repository
	// +optional
	ComponentPaths []string `json:"componentPaths,omitempty"`
}

// RepositoryPhase specifies the current state of the Repository
// +kubebuilder:validation:Enum=NotReady;Ready
type RepositoryPhase string

const (
	RepositoryNotReady RepositoryPhase = "NotReady"
	RepositoryReady    RepositoryPhase = "Ready"
)

// SnapshotInfo specifies some basic information about the Snapshots stored in this Repository
type SnapshotInfo struct {
	// Name represents the name of the Snapshot
	Name string `json:"name,omitempty"`

	// Phase represents the phase of the Snapshot
	// +optional
	Phase SnapshotPhase `json:"phase,omitempty"`

	// Session represents the name of the session that is responsible for this Snapshot
	Session string `json:"session,omitempty"`

	// Size represents the size of the Snapshot
	// +optional
	Size string `json:"size,omitempty"`

	// SnapshotTime represents the time when this Snapshot was taken
	// +optional
	SnapshotTime *metav1.Time `json:"snapshotTime,omitempty"`
}

const (
	TypeRepositoryInitialized               = "RepositoryInitialized"
	ReasonRepositoryInitializationSucceeded = "RepositoryInitializationSucceeded"
	ReasonRepositoryInitializationFailed    = "RepositoryInitializationFailed"
)

//+kubebuilder:object:root=true

// RepositoryList contains a list of Repository
type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Repository{}, &RepositoryList{})
}
