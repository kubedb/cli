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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeSnapshot     = "snap"
	ResourceKindSnapshot     = "Snapshot"
	ResourceSingularSnapshot = "snapshot"
	ResourcePluralSnapshot   = "snapshots"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=snapshots,singular=snapshot,shortName=snap,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="DatabaseName",type="string",JSONPath=".spec.databaseName"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Snapshot struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              SnapshotSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            SnapshotStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type SnapshotSpec struct {
	// Database name
	DatabaseName string `json:"databaseName" protobuf:"bytes,1,opt,name=databaseName"`

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

type SnapshotPhase string

const (
	// used for Snapshots that are currently running
	SnapshotPhaseRunning SnapshotPhase = "Running"
	// used for Snapshots that are Succeeded
	SnapshotPhaseSucceeded SnapshotPhase = "Succeeded"
	// used for Snapshots that are Failed
	SnapshotPhaseFailed SnapshotPhase = "Failed"
)

type SnapshotStatus struct {
	StartTime      *metav1.Time  `json:"startTime,omitempty" protobuf:"bytes,1,opt,name=startTime"`
	CompletionTime *metav1.Time  `json:"completionTime,omitempty" protobuf:"bytes,2,opt,name=completionTime"`
	Phase          SnapshotPhase `json:"phase,omitempty" protobuf:"bytes,3,opt,name=phase,casttype=SnapshotPhase"`
	Reason         string        `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,5,opt,name=observedGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Snapshot CRD objects
	Items []Snapshot `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
