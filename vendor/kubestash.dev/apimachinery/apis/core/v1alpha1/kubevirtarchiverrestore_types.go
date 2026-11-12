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
	v1 "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindKubeVirtArchiverRestore     = "KubeVirtArchiverRestore"
	ResourceSingularKubeVirtArchiverRestore = "kubevirtarchiverrestore"
	ResourcePluralKubeVirtArchiverRestore   = "kubevirtarchiverrestores"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubevirtarchiverrestores,singular=kubevirtarchiverrestore,shortName=kvarestore,categories={archiver,kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Checkpoint",type="string",JSONPath=".status.checkpoint"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// KubeVirtArchiverRestore restores one VirtualMachine from the complete CBT
// recovery points in a KubeVirt Repository. It is intentionally a KubeVirt-
// specific entrypoint: generic RestoreSession continues to restore an
// explicitly supplied source without understanding checkpoint chains.
type KubeVirtArchiverRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeVirtArchiverRestoreSpec   `json:"spec,omitempty"`
	Status KubeVirtArchiverRestoreStatus `json:"status,omitempty"`
}

type KubeVirtArchiverRestoreSpec struct {
	// Source identifies the Repository that contains the archived VM. The
	// Repository is imported with its Snapshot metadata when its BackupStorage
	// is initialized in another cluster, so it is also the cross-cluster restore
	// identity. Its namespace defaults to the restore's own.
	Source KubeVirtArchiveSource `json:"source"`

	// Target is the name of a new VirtualMachine, created in the restore's
	// namespace. It may keep the source's name when that namespace differs. The
	// controller fails rather than overwriting an existing target.
	Target v1.ObjectReference `json:"target"`

	// RecoveryPoint selects a complete checkpoint. An omitted recoveryPoint
	// selects the newest complete point in the Repository.
	// +optional
	RecoveryPoint *KubeVirtRecoveryPoint `json:"recoveryPoint,omitempty"`
}

type KubeVirtArchiveSource struct {
	// RepositoryRef identifies the Repository that contains the archived VM.
	RepositoryRef v1.ObjectReference `json:"repositoryRef"`
}

type KubeVirtRecoveryPoint struct {
	// Checkpoint selects one exact complete checkpoint.
	// +optional
	Checkpoint string `json:"checkpoint,omitempty"`

	// TargetTime is an RFC3339 timestamp. The selection is inclusive.
	// +optional
	TargetTime *metav1.Time `json:"targetTime,omitempty"`
}

type KubeVirtArchiverRestoreStatus struct {
	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase is Pending, Running, Succeeded, or Failed.
	// +optional
	Phase RestorePhase `json:"phase,omitempty"`

	// Conditions reports validation, checkpoint selection, and both restore stages.
	// +optional
	Conditions []v1.Condition `json:"conditions,omitempty"`

	// Snapshot is the resident Snapshot frozen for this restore.
	// +optional
	Snapshot *v1.ObjectReference `json:"snapshot,omitempty"`

	// Repository is the Repository selected for this restore.
	// +optional
	Repository *v1.ObjectReference `json:"repository,omitempty"`

	// Checkpoint and RecoveryTime identify the immutable recovery point selected.
	// +optional
	Checkpoint string `json:"checkpoint,omitempty"`
	// +optional
	RecoveryTime *metav1.Time `json:"recoveryTime,omitempty"`

	// VolumeRestoreSession and ManifestRestoreSession are the owned standard
	// RestoreSessions used for the two stages.
	// +optional
	VolumeRestoreSession *v1.ObjectReference `json:"volumeRestoreSession,omitempty"`
	// +optional
	ManifestRestoreSession *v1.ObjectReference `json:"manifestRestoreSession,omitempty"`
}

const (
	TypeKubeVirtArchiverRestoreValidated         = "Validated"
	TypeKubeVirtArchiverRecoveryPointSelected    = "RecoveryPointSelected"
	TypeKubeVirtArchiverVolumeRestoreSucceeded   = "VolumeRestoreSucceeded"
	TypeKubeVirtArchiverManifestRestoreSucceeded = "ManifestRestoreSucceeded"
)

// +kubebuilder:object:root=true

type KubeVirtArchiverRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeVirtArchiverRestore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeVirtArchiverRestore{}, &KubeVirtArchiverRestoreList{})
}
