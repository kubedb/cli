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
	ResourceCodeElasticsearchOpsRequest     = "esops"
	ResourceKindElasticsearchOpsRequest     = "ElasticsearchOpsRequest"
	ResourceSingularElasticsearchOpsRequest = "elasticsearchopsrequest"
	ResourcePluralElasticsearchOpsRequest   = "elasticsearchopsrequests"
)

// ElasticsearchOpsRequest defines a Elasticsearch DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=elasticsearchopsrequests,singular=elasticsearchopsrequest,shortName=esops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ElasticsearchOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ElasticsearchOpsRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            ElasticsearchOpsRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ElasticsearchOpsRequestSpec is the spec for ElasticsearchOpsRequest
type ElasticsearchOpsRequestSpec struct {
	// Specifies the Elasticsearch reference
	DatabaseRef v1.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`
	// Specifies the modification request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=OpsRequestType"`
	// Specifies the field information that needed to be upgraded
	Upgrade *UpgradeSpec `json:"upgrade,omitempty" protobuf:"bytes,3,opt,name=upgrade"`
	// HorizontalScaling specifies the horizontal scaling.
	HorizontalScaling *ElasticsearchHorizontalScalingSpec `json:"horizontalScaling,omitempty" protobuf:"bytes,4,opt,name=horizontalScaling"`
}

// ElasticsearchHorizontalScalingSpec contains the horizontal scaling information of an Elasticsearch cluster
type ElasticsearchHorizontalScalingSpec struct {
	// Number of master nodes
	Master *int32 `json:"master,omitempty" protobuf:"bytes,1,opt,name=master"`
	// Number of data nodes
	Data *int32 `json:"data,omitempty" protobuf:"bytes,2,opt,name=data"`
	// Number of client nodes
	Client *int32 `json:"client,omitempty" protobuf:"bytes,3,opt,name=client"`
}

// ElasticsearchOpsRequestStatus is the status for ElasticsearchOpsRequest
type ElasticsearchOpsRequestStatus struct {
	// Specifies the current phase of the modification request
	// +optional
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

// ElasticsearchOpsRequestList is a list of ElasticsearchOpsRequests
type ElasticsearchOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of ElasticsearchOpsRequest CRD objects
	Items []ElasticsearchOpsRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
