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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeMilvusAutoscaler     = "mvscaler"
	ResourceKindMilvusAutoscaler     = "MilvusAutoscaler"
	ResourceSingularMilvusAutoscaler = "milvusautoscaler"
	ResourcePluralMilvusAutoscaler   = "milvusautoscalers"
)

// MilvusAutoscaler is the configuration for a milvus database
// autoscaler, which automatically manages pod resources based on historical and
// real time resource utilization.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=milvusautoscalers,singular=milvusautoscaler,shortName=mvscaler,categories={autoscaler,kubedb,appscode}
// +kubebuilder:subresource:status
type MilvusAutoscaler struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the behavior of the autoscaler.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status.
	Spec MilvusAutoscalerSpec `json:"spec"`

	// Current information about the autoscaler.
	// +optional
	Status AutoscalerStatus `json:"status,omitempty"`
}

// MilvusAutoscalerSpec is the specification of the behavior of the autoscaler.
type MilvusAutoscalerSpec struct {
	DatabaseRef *core.LocalObjectReference `json:"databaseRef"`

	// This field will be used to control the behaviour of ops-manager
	OpsRequestOptions *OpsRequestOptions `json:"opsRequestOptions,omitempty"`

	Compute *MilvusComputeAutoscalerSpec `json:"compute,omitempty"`
	Storage *MilvusStorageAutoscalerSpec `json:"storage,omitempty"`
}

type MilvusComputeAutoscalerSpec struct {
	// +optional
	NodeTopology *NodeTopology `json:"nodeTopology,omitempty"`

	// Standalone mode
	Node *ComputeAutoscalerSpec `json:"node,omitempty"`

	// Distributed mode
	Proxy         *ComputeAutoscalerSpec `json:"proxy,omitempty"`
	MixCoord      *ComputeAutoscalerSpec `json:"mixcoord,omitempty"`
	DataNode      *ComputeAutoscalerSpec `json:"datanode,omitempty"`
	QueryNode     *ComputeAutoscalerSpec `json:"querynode,omitempty"`
	StreamingNode *ComputeAutoscalerSpec `json:"streamingnode,omitempty"`
}

type MilvusStorageAutoscalerSpec struct {
	// Standalone mode
	Node *StorageAutoscalerSpec `json:"node,omitempty"`

	// Distributed mode (only StreamingNode has persistent storage)
	StreamingNode *StorageAutoscalerSpec `json:"streamingnode,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// MilvusAutoscalerList is a list of MilvusAutoscaler objects.
type MilvusAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata"`

	// items is the list of milvus database autoscaler objects.
	Items []MilvusAutoscaler `json:"items"`
}
