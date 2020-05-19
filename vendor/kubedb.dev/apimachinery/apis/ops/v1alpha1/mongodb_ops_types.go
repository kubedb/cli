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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceCodeMongoDBOpsRequest     = "mgops"
	ResourceKindMongoDBOpsRequest     = "MongoDBOpsRequest"
	ResourceSingularMongoDBOpsRequest = "mongodbopsrequest"
	ResourcePluralMongoDBOpsRequest   = "mongodbopsrequests"
)

// MongoDBOpsRequest defines a MongoDB DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mongodbopsrequests,singular=mongodbopsrequest,shortName=mgops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MongoDBOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MongoDBOpsRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            MongoDBOpsRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MongoDBOpsRequestSpec is the spec for MongoDBOpsRequest
type MongoDBOpsRequestSpec struct {
	// Specifies the MongoDB reference
	DatabaseRef v1.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`
	// Specifies the modification request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=OpsRequestType"`
	// Specifies information necessary for upgrading mongodb
	Upgrade *UpgradeSpec `json:"upgrade,omitempty" protobuf:"bytes,3,opt,name=upgrade"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *MongoDBHorizontalScalingSpec `json:"horizontalScaling,omitempty" protobuf:"bytes,4,opt,name=horizontalScaling"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *MongoDBVerticalScalingSpec `json:"verticalScaling,omitempty" protobuf:"bytes,5,opt,name=verticalScaling"`
}

// MongoDBShardNode is the spec for mongodb Shard
type MongoDBShardNode struct {
	Shards   int32 `json:"shards,omitempty" protobuf:"bytes,1,opt,name=shards"`
	Replicas int32 `json:"replicas,omitempty" protobuf:"bytes,2,opt,name=replicas"`
}

// ConfigNode is the spec for mongodb ConfigServer
type ConfigNode struct {
	Replicas int32 `json:"replicas,omitempty" protobuf:"bytes,1,opt,name=replicas"`
}

// MongosNode is the spec for mongodb Mongos
type MongosNode struct {
	Replicas int32 `json:"replicas,omitempty" protobuf:"bytes,1,opt,name=replicas"`
}

// HorizontalScaling is the spec for mongodb horizontal scaling
type MongoDBHorizontalScalingSpec struct {
	Shard        *MongoDBShardNode `json:"shard,omitempty" protobuf:"bytes,1,opt,name=shard"`
	ConfigServer *ConfigNode       `json:"configServer,omitempty" protobuf:"bytes,2,opt,name=configServer"`
	Mongos       *MongosNode       `json:"mongos,omitempty" protobuf:"bytes,3,opt,name=mongos"`
	Replicas     *int32            `json:"replicas,omitempty" protobuf:"bytes,4,opt,name=replicas"`
}

// MongoDBVerticalScalingSpec is the spec for mongodb vertical scaling
type MongoDBVerticalScalingSpec struct {
	Containers   []ContainerResources `json:"containers,omitempty" protobuf:"bytes,1,rep,name=containers"`
	Mongos       []ContainerResources `json:"mongos,omitempty" protobuf:"bytes,2,rep,name=mongos"`
	ConfigServer []ContainerResources `json:"configServer,omitempty" protobuf:"bytes,3,rep,name=configServer"`
	Shard        []ContainerResources `json:"shard,omitempty" protobuf:"bytes,4,rep,name=shard"`
}

// MongoDBOpsRequestStatus is the status for MongoDBOpsRequest
type MongoDBOpsRequestStatus struct {
	Phase OpsRequestPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=OpsRequestPhase"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,2,opt,name=observedGeneration"`
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty" protobuf:"bytes,3,rep,name=conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBOpsRequestList is a list of MongoDBOpsRequests
type MongoDBOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of MongoDBOpsRequest CRD objects
	Items []MongoDBOpsRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
