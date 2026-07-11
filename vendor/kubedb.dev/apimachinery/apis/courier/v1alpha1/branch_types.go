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
)

const (
	ResourceKindBranch     = "Branch"
	ResourceSingularBranch = "branch"
	ResourcePluralBranches = "branches"

	// BranchCleanupFinalizer is set on a Branch when the operator adopts it and cleared only
	// after ordered teardown finishes.
	BranchCleanupFinalizer = "courier.kubedb.com/branch-cleanup"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=branches,singular=branch,shortName=br,categories={kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Mode",type="string",JSONPath=".status.mode"
// +kubebuilder:printcolumn:name="Target",type="string",JSONPath=".status.targetRef.name"
// +kubebuilder:printcolumn:name="Freshness",type="string",JSONPath=".status.freshness"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Branch struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Branch
	// +required
	Spec BranchSpec `json:"spec"`

	// status defines the observed state of Branch
	// +optional
	Status BranchStatus `json:"status,omitzero"`
}

// BranchSpec defines the desired state of Branch. One Branch CR is one branch, and it doubles as the
// session object.
type BranchSpec struct {
	// Source is the KubeDB Database whose storage is cloned. Branch has no external source.
	Source BranchSource `json:"source"`

	// Target describes only what differs from the source: the target cluster, namespace, name,
	// StorageClass, and cpu/memory. Everything else is copied from the source Database.
	Target BranchTarget `json:"target"`

	// ResetRootPassword resets the branch's root password after provisioning, so the source's
	// password does not unlock the branch.
	// +optional
	ResetRootPassword bool `json:"resetRootPassword,omitempty"`

	// DataMassageImage is an optional user-provided container run as the LAST step to massage or
	// anonymize the branch data.
	// +optional
	DataMassageImage string `json:"dataMassageImage,omitempty"`

	// Schedule optionally refreshes the branch on a cron cadence. Omit for a one-shot branch.
	// +optional
	Schedule *BranchSchedule `json:"schedule,omitempty"`

	// HistoryLimit bounds status.history (default: last 3 successful, last 2 failed).
	// +optional
	HistoryLimit *BranchHistoryLimit `json:"historyLimit,omitempty"`

	// VolumeSnapshotClassName is the VolumeSnapshotClass used wherever courier creates a
	// VolumeSnapshot — snapshotting the source PVCs and, for cross-namespace/cross-cluster
	// branches, the importing snapshot in the target. It must match the CSI driver backing the
	// volumes. When empty, courier auto-resolves the default class for the driver.
	// +optional
	VolumeSnapshotClassName string `json:"volumeSnapshotClassName,omitempty"`

	// DeletionPolicy decides the target's fate on Branch deletion.
	// +kubebuilder:default=Delete
	// +optional
	DeletionPolicy BranchDeletionPolicy `json:"deletionPolicy,omitempty"`
}

// BranchSource points at a KubeDB Database.
type BranchSource struct {
	// DatabaseRef refers to the source KubeDB Database (kind and name).
	DatabaseRef corev1.TypedLocalObjectReference `json:"databaseRef"`

	// Namespace of the source Database. Defaults to the Branch's namespace when empty.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// BranchTarget describes the target Database. spec.target.cluster equal to the source's cluster is a
// same-cluster branch; a different cluster is a cross-cluster branch. Omit cluster for a same-cluster
// (Local) branch.
type BranchTarget struct {
	// ClusterName is the target cluster name. Empty (or equal to the source's own cluster) means a
	// same-cluster (Local) branch; a different ClusterName selects cross-cluster (OCM).
	// +optional
	ClusterName string `json:"clusterName,omitempty"`

	// Namespace of the target Database.
	Namespace string `json:"namespace"`

	// Name of the target Database.
	Name string `json:"name"`

	// StorageClassName is the StorageClass in the TARGET cluster.
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`

	// Resources are the cpu/memory requests and limits in the TARGET cluster.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// IssuerRef references a cert-manager Issuer or ClusterIssuer in the TARGET cluster. TLS secrets are
	// namespace and cluster scoped, so a branch cannot reuse the source's; when the source Database has
	// TLS enabled the operator points the branch's TLS at this issuer and KubeDB mints fresh
	// certificates for the branch. Required for a TLS-enabled source, ignored otherwise.
	// +optional
	IssuerRef *corev1.TypedLocalObjectReference `json:"issuerRef,omitempty"`
}

// BranchSchedule holds the refresh cadence.
type BranchSchedule struct {
	// Cron is the refresh schedule in standard cron syntax.
	Cron string `json:"cron"`
}

// BranchHistoryLimit bounds status.history.
type BranchHistoryLimit struct {
	// Success is the number of successful runs to retain (default 3).
	// +kubebuilder:default=3
	// +optional
	Success *int32 `json:"success,omitempty"`

	// Failed is the number of failed runs to retain (default 2).
	// +kubebuilder:default=2
	// +optional
	Failed *int32 `json:"failed,omitempty"`
}

// BranchDeletionPolicy decides the target Database's fate on Branch deletion.
// +kubebuilder:validation:Enum=Delete;Orphan
type BranchDeletionPolicy string

const (
	// BranchDeletionPolicyDelete tears the branch down (default).
	BranchDeletionPolicyDelete BranchDeletionPolicy = "Delete"
	// BranchDeletionPolicyOrphan keeps the target as a standalone KubeDB Database.
	BranchDeletionPolicyOrphan BranchDeletionPolicy = "Orphan"
)

// BranchStatus defines the observed state of Branch.
type BranchStatus struct {
	// Phase is the current phase of the branch.
	// +optional
	Phase BranchPhase `json:"phase,omitempty"`

	// Mode is how this operator is participating in the branch.
	// +optional
	Mode BranchMode `json:"mode,omitempty"`

	// TargetRef references the branched Database.
	// +optional
	TargetRef *corev1.ObjectReference `json:"targetRef,omitempty"`

	// Snapshot references the source snapshot the current branch was cloned from.
	// +optional
	Snapshot *BranchSnapshotRef `json:"snapshot,omitempty"`

	// LastRefreshAt is the time of the last successful refresh.
	// +optional
	LastRefreshAt *metav1.Time `json:"lastRefreshAt,omitempty"`

	// Freshness is the human-readable age of the branch data since the last refresh (e.g. "3m", "1h2m").
	// +optional
	Freshness metav1.Duration `json:"freshness,omitempty"`

	// History is the bounded refresh history (bounded by spec.historyLimit).
	// +optional
	History []BranchRun `json:"history,omitempty"`

	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// BranchPhase is the lifecycle phase of a Branch.
// +kubebuilder:validation:Enum=Pending;Snapshotting;Cloning;Provisioning;Massaging;Ready;Refreshing;Deleting;Failed
type BranchPhase string

const (
	BranchPhasePending      BranchPhase = "Pending"
	BranchPhaseSnapshotting BranchPhase = "Snapshotting"
	BranchPhaseCloning      BranchPhase = "Cloning"
	BranchPhaseProvisioning BranchPhase = "Provisioning"
	BranchPhaseMassaging    BranchPhase = "Massaging"
	BranchPhaseReady        BranchPhase = "Ready"
	BranchPhaseRefreshing   BranchPhase = "Refreshing"
	BranchPhaseDeleting     BranchPhase = "Deleting"
	BranchPhaseFailed       BranchPhase = "Failed"
)

// BranchMode is how the branch operator participates in a branch.
// +kubebuilder:validation:Enum=Local;Initiator;Creator
type BranchMode string

const (
	// BranchModeLocal is a same-cluster branch (the operator runs the whole flow).
	BranchModeLocal BranchMode = "Local"
	// BranchModeInitiator is the source cluster of a cross-cluster branch.
	BranchModeInitiator BranchMode = "Initiator"
	// BranchModeCreator is the target cluster of a cross-cluster branch.
	BranchModeCreator BranchMode = "Creator"
)

// BranchSnapshotType is the kind of snapshot backing a branch.
// +kubebuilder:validation:Enum=VolumeGroupSnapshot;VolumeSnapshot
type BranchSnapshotType string

const (
	BranchSnapshotTypeVolumeGroupSnapshot BranchSnapshotType = "VolumeGroupSnapshot"
	BranchSnapshotTypeVolumeSnapshot      BranchSnapshotType = "VolumeSnapshot"
)

// BranchSnapshotRef references the source snapshot.
type BranchSnapshotRef struct {
	// Type is the snapshot kind (VolumeGroupSnapshot preferred, VolumeSnapshot fallback).
	Type BranchSnapshotType `json:"type,omitempty"`
	// Ref is the name of the snapshot object.
	Ref string `json:"ref,omitempty"`
}

// BranchRunResult is the outcome of a refresh run.
// +kubebuilder:validation:Enum=Succeeded;Failed
type BranchRunResult string

const (
	BranchRunSucceeded BranchRunResult = "Succeeded"
	BranchRunFailed    BranchRunResult = "Failed"
)

// BranchRun is one entry in the refresh history.
type BranchRun struct {
	// At is when the run finished.
	At metav1.Time `json:"at,omitempty"`
	// Result is the outcome of the run.
	Result BranchRunResult `json:"result,omitempty"`
	// Message is an optional human-readable detail (for a failed run).
	// +optional
	Message string `json:"message,omitempty"`
}

// BranchList contains a list of Branch

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BranchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Branch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Branch{}, &BranchList{})
}
