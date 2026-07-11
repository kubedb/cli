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
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindBranchWork     = "BranchWork"
	ResourceSingularBranchWork = "branchwork"
	ResourcePluralBranchWorks  = "branchworks"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=branchworks,singular=branchwork,shortName=bw,categories={kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="TargetCluster",type="string",JSONPath=".spec.targetCluster"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.branch.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
//
// BranchWork is a hub-only cross-cluster delivery object, similar to OCM ManifestWork. A source-cluster
// Branch operator creates it in its own hub namespace to drive a branch into a different cluster; the
// addon manager translates it into a ManifestWork in the target cluster's namespace and syncs status
// back. Same-cluster branches never use it.
type BranchWork struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of BranchWork
	// +required
	Spec BranchWorkSpec `json:"spec"`

	// status defines the observed state of BranchWork
	// +optional
	Status BranchWorkStatus `json:"status,omitzero"`
}

// BranchWorkSpec carries the target cluster and the manifests to create there.
type BranchWorkSpec struct {
	// TargetCluster is the cluster the manifests are delivered to.
	TargetCluster string `json:"targetCluster"`

	// Manifests are the resources to create in the target cluster (the Branch CR, the target Database,
	// the auth and config secrets, the VolumeSnapshot(s), and anything else the target needs).
	// +optional
	Manifests []runtime.RawExtension `json:"manifests,omitempty"`

	// RefreshGeneration is bumped by the source each time it re-ships a fresh snapshot; the creator
	// keys its swap off a change here.
	// +optional
	RefreshGeneration int64 `json:"refreshGeneration,omitempty"`

	// DeleteOption mirrors ManifestWork deletion semantics for the delivered resources.
	// +optional
	DeleteOption BranchDeletionPolicy `json:"deleteOption,omitempty"`
}

// BranchWorkStatus mirrors the delivered Branch status plus per-manifest apply/feedback conditions,
// modeled on OCM ManifestWork.
type BranchWorkStatus struct {
	// ObservedRefreshGeneration is the spec.refreshGeneration the target has acted on.
	// +optional
	ObservedRefreshGeneration int64 `json:"observedRefreshGeneration,omitempty"`

	// Branch mirrors the target-side Branch status so the source can copy it into its own Branch.
	// +optional
	Branch *BranchStatus `json:"branch,omitempty"`

	// Manifests carries the per-manifest apply/feedback conditions from the target (ManifestWork-style).
	// +optional
	Manifests []BranchWorkManifestCondition `json:"manifests,omitempty"`

	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// BranchWorkManifestCondition reports the apply/feedback status of a single delivered manifest,
// modeled on OCM ManifestWork's ManifestCondition.
type BranchWorkManifestCondition struct {
	// Ordinal is the index of the manifest in spec.manifests.
	Ordinal int32 `json:"ordinal"`

	// ResourceMeta identifies the delivered resource.
	// +optional
	ResourceMeta BranchWorkResourceMeta `json:"resourceMeta,omitempty"`

	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// BranchWorkResourceMeta identifies a delivered resource.
type BranchWorkResourceMeta struct {
	// +optional
	Group string `json:"group,omitempty"`
	// +optional
	Version string `json:"version,omitempty"`
	// +optional
	Kind string `json:"kind,omitempty"`
	// +optional
	Resource string `json:"resource,omitempty"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
}

// BranchWorkList contains a list of BranchWork

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BranchWorkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []BranchWork `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BranchWork{}, &BranchWorkList{})
}
