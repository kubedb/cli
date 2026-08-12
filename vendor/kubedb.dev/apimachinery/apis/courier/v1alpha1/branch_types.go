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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
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
// +kubebuilder:printcolumn:name="Target",type="string",JSONPath=".spec.target.name"
// +kubebuilder:printcolumn:name="Freshness",type="date",JSONPath=".status.lastSuccessfulRefreshTime"
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

// PostActions is a container run as a Job against the branched Database after it is Ready,
// typically to massage or anonymize the cloned data. The Job receives the connection
// environment for the branch (host, port, and credentials from its auth secret).
type PostActions struct {
	// Image is the container image to run.
	Image string `json:"image,omitempty"`

	// JobDefaults specifies default settings for the Job (pull policy, backoff limit, TTL,
	// active deadline).
	// +optional
	JobDefaults *JobDefaults `json:"jobDefaults,omitempty"`

	// JobTemplate specifies runtime configurations for the Job, so the user can customize
	// scheduling, resources, security context, volumes, etc.
	// +optional
	JobTemplate ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
}

// BranchSpec defines the desired state of Branch. One Branch CR is one branch, and it doubles as the
// session object.
type BranchSpec struct {
	// Source is the KubeDB Database whose storage is cloned. Branch has no external source.
	Source BranchSource `json:"source"`

	// Target describes only what differs from the source: the target cluster, namespace, name,
	// StorageClass, and cpu/memory. Everything else is copied from the source Database.
	Target BranchTarget `json:"target"`

	// ResetRootPassword gives the branch its own root credential instead of the source's, so the
	// source password does not unlock the branch.
	// +optional
	ResetRootPassword bool `json:"resetRootPassword,omitempty"`

	// PostActions are user-provided containers run as the LAST step of a branch, after the target
	// Database is Ready. Each action runs as a Job
	// against the branch — typically to massage or anonymize the cloned data. They run in order,
	// and the branch only becomes Ready once every action has succeeded.
	// +optional
	PostActions []PostActions `json:"postActions,omitempty"`

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

	// Resources lists the objects this branch owns, for audit and cleanup visibility.
	// +optional
	Resources *BranchOwnedResources `json:"resources,omitempty"`

	// --- snapshot provenance (current generation) ---

	// Snapshot describes the source snapshot set backing the current branch generation.
	// +optional
	Snapshot *BranchSnapshotStatus `json:"snapshot,omitempty"`

	// --- refresh / freshness ---

	// RefreshGeneration is the current refresh generation; it bumps on each scheduled refresh.
	// +optional
	RefreshGeneration int64 `json:"refreshGeneration,omitempty"`

	// LastRefreshTime is the time of the last refresh ATTEMPT, regardless of outcome.
	// +optional
	LastRefreshTime *metav1.Time `json:"lastRefreshTime,omitempty"`

	// LastSuccessfulRefreshTime is the time of the last SUCCESSFUL refresh. It is how fresh the
	// branch data is: consumers compute the data age as now - lastSuccessfulRefreshTime. A failed
	// tick updates LastRefreshTime but not this field.
	// +optional
	LastSuccessfulRefreshTime *metav1.Time `json:"lastSuccessfulRefreshTime,omitempty"`

	// NextRefreshTime is the next scheduled refresh, computed from spec.schedule.cron. Nil for
	// a one-shot branch.
	// +optional
	NextRefreshTime *metav1.Time `json:"nextRefreshTime,omitempty"`

	// History is the bounded refresh history (bounded by spec.historyLimit).
	// +optional
	History []BranchRun `json:"history,omitempty"`

	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// BranchOwnedResources references the objects created and owned by a branch.
type BranchOwnedResources struct {
	// ClonedPVCs are the target PVCs cloned from the source snapshots, ordered by ordinal.
	// +optional
	ClonedPVCs []string `json:"clonedPVCs,omitempty"`

	// AuthSecret is the branch's auth Secret (the credential matching the cloned data).
	// +optional
	AuthSecret string `json:"authSecret,omitempty"`

	// ConfigSecret is the branch's config Secret, when the engine uses one.
	// +optional
	ConfigSecret string `json:"configSecret,omitempty"`

	// PostActionJob is the Job name of the current generation's first
	// spec.postActions entry, set only when spec.postActions is used.
	// +optional
	PostActionJob string `json:"postActionJob,omitempty"`
}

// BranchPhase is the lifecycle phase of a Branch.
// +kubebuilder:validation:Enum=Pending;Snapshotting;Cloning;Provisioning;ActionsRunning;Ready;Refreshing;Deleting;Failed
type BranchPhase string

const (
	BranchPhasePending        BranchPhase = "Pending"
	BranchPhaseSnapshotting   BranchPhase = "Snapshotting"
	BranchPhaseCloning        BranchPhase = "Cloning"
	BranchPhaseProvisioning   BranchPhase = "Provisioning"
	BranchPhaseActionsRunning BranchPhase = "ActionsRunning"
	BranchPhaseReady          BranchPhase = "Ready"
	BranchPhaseRefreshing     BranchPhase = "Refreshing"
	BranchPhaseDeleting       BranchPhase = "Deleting"
	BranchPhaseFailed         BranchPhase = "Failed"
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

// BranchSnapshotStatus describes the source snapshot set backing the current branch generation.
type BranchSnapshotStatus struct {
	// Strategy is how the source was snapshotted: VolumeGroupSnapshot (group-consistent) or
	// VolumeSnapshot (per-PVC fallback, used when the driver has no VolumeGroupSnapshotClass).
	// +optional
	Strategy BranchSnapshotType `json:"strategy,omitempty"`

	// Generation is the refresh generation these snapshots belong to.
	// +optional
	Generation int64 `json:"generation,omitempty"`

	// GroupRef is the VolumeGroupSnapshot object name, set only when Strategy is
	// VolumeGroupSnapshot.
	// +optional
	GroupRef string `json:"groupRef,omitempty"`

	// Ready is true when every member snapshot is readyToUse.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Members is one entry per source data PVC, ordered by ordinal and aligned to the cloned
	// target PVCs.
	// +optional
	Members []BranchSnapshotMember `json:"members,omitempty"`
}

// BranchSnapshotMember is one source VolumeSnapshot backing a single data PVC of the branch.
type BranchSnapshotMember struct {
	// Name is the VolumeSnapshot object name.
	Name string `json:"name"`

	// SourcePVC is the source PVC this snapshot was taken from (its ordinal maps to the cloned
	// target PVC).
	// +optional
	SourcePVC string `json:"sourcePVC,omitempty"`

	// ReadyToUse mirrors the VolumeSnapshot's readyToUse status.
	// +optional
	ReadyToUse bool `json:"readyToUse,omitempty"`

	// RestoreSize is the snapshot's restore size, when reported by the CSI driver.
	// +optional
	RestoreSize *resource.Quantity `json:"restoreSize,omitempty"`

	// CreationTime is when the snapshot was taken, when reported by the CSI driver.
	// +optional
	CreationTime *metav1.Time `json:"creationTime,omitempty"`
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
