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
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceCodeElasticsearchAutoscaler     = "esscaler"
	ResourceKindElasticsearchAutoscaler     = "ElasticsearchAutoscaler"
	ResourceSingularElasticsearchAutoscaler = "elasticsearchautoscaler"
	ResourcePluralElasticsearchAutoscaler   = "elasticsearchautoscalers"
)

// ElasticsearchAutoscaler is the configuration for a horizontal pod
// autoscaler, which automatically manages the replica count of any resource
// implementing the scale subresource based on the metrics specified.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=elasticsearchautoscalers,singular=elasticsearchautoscaler,shortName=esscaler,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
type ElasticsearchAutoscaler struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// spec is the specification for the behaviour of the autoscaler.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status.
	// +optional
	Spec ElasticsearchAutoscalerSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// status is the current information about the autoscaler.
	// +optional
	Status ElasticsearchAutoscalerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ElasticsearchAutoscalerSpec is the specification of the behavior of the autoscaler.
type ElasticsearchAutoscalerSpec struct {
	DatabaseRef *core.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`

	Compute *ElasticsearchComputeAutoscalerSpec `json:"compute,omitempty" protobuf:"bytes,2,opt,name=compute"`
	Storage *ElasticsearchStorageAutoscalerSpec `json:"storage,omitempty" protobuf:"bytes,3,opt,name=storage"`
}

type ElasticsearchComputeAutoscalerSpec struct {
	Node             *ComputeAutoscalerSpec                      `json:"node,omitempty" protobuf:"bytes,1,opt,name=node"`
	Topology         *ElasticsearchComputeTopologyAutoscalerSpec `json:"topology,omitempty" protobuf:"bytes,2,opt,name=topology"`
	DisableScaleDown bool                                        `json:"disableScaleDown,omitempty" protobuf:"varint,3,opt,name=disableScaleDown"`
}

type ElasticsearchComputeTopologyAutoscalerSpec struct {
	Master *ComputeAutoscalerSpec `json:"master,omitempty" protobuf:"bytes,1,opt,name=master"`
	Data   *ComputeAutoscalerSpec `json:"data,omitempty" protobuf:"bytes,2,opt,name=data"`
	Ingest *ComputeAutoscalerSpec `json:"ingest,omitempty" protobuf:"bytes,3,opt,name=ingest"`
}

type ElasticsearchStorageAutoscalerSpec struct {
	Node     *StorageAutoscalerSpec                      `json:"node,omitempty" protobuf:"bytes,1,opt,name=node"`
	Topology *ElasticsearchStorageTopologyAutoscalerSpec `json:"topology,omitempty" protobuf:"bytes,2,opt,name=topology"`
}

type ElasticsearchStorageTopologyAutoscalerSpec struct {
	Master *StorageAutoscalerSpec `json:"master,omitempty" protobuf:"bytes,1,opt,name=master"`
	Data   *StorageAutoscalerSpec `json:"data,omitempty" protobuf:"bytes,2,opt,name=data"`
	Ingest *StorageAutoscalerSpec `json:"ingest,omitempty" protobuf:"bytes,3,opt,name=ingest"`
}

// ElasticsearchAutoscalerStatus describes the runtime state of the autoscaler.
type ElasticsearchAutoscalerStatus struct {
	// observedGeneration is the most recent generation observed by this autoscaler.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// Conditions is the set of conditions required for this autoscaler to scale its target,
	// and indicates whether or not those conditions are met.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []kmapi.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
}

// ElasticsearchAutoscalerConditionType are the valid conditions of
// a ElasticsearchAutoscaler.
type ElasticsearchAutoscalerConditionType string

var (
	// ConfigDeprecated indicates that this VPA configuration is deprecated
	// and will stop being supported soon.
	ElasticsearchAutoscalerConfigDeprecated ElasticsearchAutoscalerConditionType = "ConfigDeprecated"
	// ConfigUnsupported indicates that this VPA configuration is unsupported
	// and recommendations will not be provided for it.
	ElasticsearchAutoscalerConfigUnsupported ElasticsearchAutoscalerConditionType = "ConfigUnsupported"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ElasticsearchAutoscalerList is a list of ElasticsearchAutoscaler objects.
type ElasticsearchAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	// items is the list of elasticsearch database autoscaler objects.
	Items []ElasticsearchAutoscaler `json:"items" protobuf:"bytes,2,rep,name=items"`
}
