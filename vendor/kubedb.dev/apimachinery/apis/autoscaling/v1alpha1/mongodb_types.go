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
	ResourceCodeMongoDBAutoscaler     = "mgautoscaler"
	ResourceKindMongoDBAutoscaler     = "MongoDBAutoscaler"
	ResourceSingularMongoDBAutoscaler = "mongodbautoscaler"
	ResourcePluralMongoDBAutoscaler   = "mongodbautoscalers"
)

// MongoDBAutoscaler is the configuration for a mongodb database
// autoscaler, which automatically manages pod resources based on historical and
// real time resource utilization.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mongodbautoscalers,singular=mongodbautoscaler,shortName=mgautoscaler,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
type MongoDBAutoscaler struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the behavior of the autoscaler.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status.
	Spec MongoDBAutoscalerSpec `json:"spec" protobuf:"bytes,2,name=spec"`

	// Current information about the autoscaler.
	// +optional
	Status MongoDBAutoscalerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MongoDBAutoscalerSpec is the specification of the behavior of the autoscaler.
type MongoDBAutoscalerSpec struct {
	DatabaseRef *core.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`

	Compute *MongoDBComputeAutoscalerSpec `json:"compute,omitempty" protobuf:"bytes,2,opt,name=compute"`
	Storage *MongoDBStorageAutoscalerSpec `json:"storage,omitempty" protobuf:"bytes,3,opt,name=storage"`
}

type MongoDBComputeAutoscalerSpec struct {
	Standalone       *ComputeAutoscalerSpec `json:"standalone,omitempty" protobuf:"bytes,1,opt,name=standalone"`
	ReplicaSet       *ComputeAutoscalerSpec `json:"replicaSet,omitempty" protobuf:"bytes,2,opt,name=replicaSet"`
	ConfigServer     *ComputeAutoscalerSpec `json:"configServer,omitempty" protobuf:"bytes,3,opt,name=configServer"`
	Shard            *ComputeAutoscalerSpec `json:"shard,omitempty" protobuf:"bytes,4,opt,name=shard"`
	Mongos           *ComputeAutoscalerSpec `json:"mongos,omitempty" protobuf:"bytes,5,opt,name=mongos"`
	DisableScaleDown bool                   `json:"disableScaleDown,omitempty" protobuf:"varint,6,opt,name=disableScaleDown"`
}

type MongoDBStorageAutoscalerSpec struct {
	Standalone   *StorageAutoscalerSpec `json:"standalone,omitempty" protobuf:"bytes,1,opt,name=standalone"`
	ReplicaSet   *StorageAutoscalerSpec `json:"replicaSet,omitempty" protobuf:"bytes,2,opt,name=replicaSet"`
	ConfigServer *StorageAutoscalerSpec `json:"configServer,omitempty" protobuf:"bytes,3,opt,name=configServer"`
	Shard        *StorageAutoscalerSpec `json:"shard,omitempty" protobuf:"bytes,4,opt,name=shard"`
}

// MongoDBAutoscalerStatus describes the runtime state of the autoscaler.
type MongoDBAutoscalerStatus struct {
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

// MongoDBAutoscalerConditionType are the valid conditions of
// a MongoDBAutoscaler.
type MongoDBAutoscalerConditionType string

var (
	// ConfigDeprecated indicates that this VPA configuration is deprecated
	// and will stop being supported soon.
	MongoDBAutoscalerConfigDeprecated MongoDBAutoscalerConditionType = "ConfigDeprecated"
	// ConfigUnsupported indicates that this VPA configuration is unsupported
	// and recommendations will not be provided for it.
	MongoDBAutoscalerConfigUnsupported MongoDBAutoscalerConditionType = "ConfigUnsupported"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// MongoDBAutoscalerList is a list of MongoDBAutoscaler objects.
type MongoDBAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	// items is the list of mongodb database autoscaler objects.
	Items []MongoDBAutoscaler `json:"items" protobuf:"bytes,2,rep,name=items"`
}
