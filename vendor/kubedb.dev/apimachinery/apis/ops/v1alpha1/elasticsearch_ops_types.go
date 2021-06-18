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
	"k8s.io/apimachinery/pkg/api/resource"
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
	DatabaseRef core.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=OpsRequestType"`
	// Specifies information necessary for upgrading Elasticsearch
	Upgrade *ElasticsearchUpgradeSpec `json:"upgrade,omitempty" protobuf:"bytes,3,opt,name=upgrade"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *ElasticsearchHorizontalScalingSpec `json:"horizontalScaling,omitempty" protobuf:"bytes,4,opt,name=horizontalScaling"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *ElasticsearchVerticalScalingSpec `json:"verticalScaling,omitempty" protobuf:"bytes,5,opt,name=verticalScaling"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *ElasticsearchVolumeExpansionSpec `json:"volumeExpansion,omitempty" protobuf:"bytes,6,opt,name=volumeExpansion"`
	// Specifies information necessary for custom configuration of Elasticsearch
	Configuration *ElasticsearchCustomConfigurationSpec `json:"configuration,omitempty" protobuf:"bytes,7,opt,name=configuration"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty" protobuf:"bytes,8,opt,name=tls"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty" protobuf:"bytes,9,opt,name=restart"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty" protobuf:"bytes,10,opt,name=timeout"`
}

type ElasticsearchUpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty" protobuf:"bytes,1,opt,name=targetVersion"`
}

// ElasticsearchHorizontalScalingSpec contains the horizontal scaling information of an Elasticsearch cluster
type ElasticsearchHorizontalScalingSpec struct {
	// Number of combined (i.e. master, data, ingest) node
	Node *int32 `json:"node,omitempty" protobuf:"varint,1,opt,name=node"`
	// Node topology specification
	Topology *ElasticsearchHorizontalScalingTopologySpec `json:"topology,omitempty" protobuf:"bytes,2,opt,name=topology"`
}

// ElasticsearchHorizontalScalingTopologySpec contains the horizontal scaling information in cluster topology mode
type ElasticsearchHorizontalScalingTopologySpec struct {
	// Number of master nodes
	Master *int32 `json:"master,omitempty" protobuf:"varint,1,opt,name=master"`
	// Number of data nodes
	Data *int32 `json:"data,omitempty" protobuf:"varint,2,opt,name=data"`
	// Number of ingest nodes
	Ingest *int32 `json:"ingest,omitempty" protobuf:"varint,3,opt,name=ingest"`
}

// ElasticsearchVerticalScalingSpec is the spec for Elasticsearch vertical scaling
type ElasticsearchVerticalScalingSpec struct {
	// Resource spec for combined nodes
	Node *core.ResourceRequirements `json:"node,omitempty" protobuf:"bytes,1,opt,name=node"`
	// Resource spec for exporter sidecar
	Exporter *core.ResourceRequirements `json:"exporter,omitempty" protobuf:"bytes,2,opt,name=exporter"`
	// Specifies the resource spec for cluster in topology mode
	Topology *ElasticsearchVerticalScalingTopologySpec `json:"topology,omitempty" protobuf:"bytes,3,opt,name=topology"`
}

// ElasticsearchVerticalScalingTopologySpec is the resource spec in the cluster topology mode
type ElasticsearchVerticalScalingTopologySpec struct {
	// Resource spec for master nodes
	Master *core.ResourceRequirements `json:"master,omitempty" protobuf:"bytes,1,opt,name=master"`
	// Resource spec for data nodes
	Data *core.ResourceRequirements `json:"data,omitempty" protobuf:"bytes,2,opt,name=data"`
	// Resource spec for ingest nodes
	Ingest *core.ResourceRequirements `json:"ingest,omitempty" protobuf:"bytes,3,opt,name=ingest"`
}

// ElasticsearchVolumeExpansionSpec is the spec for Elasticsearch volume expansion
type ElasticsearchVolumeExpansionSpec struct {
	// volume specification for combined nodes
	Node *resource.Quantity `json:"node,omitempty" protobuf:"bytes,1,opt,name=node"`
	// volume specification for nodes in cluster topology
	Topology *ElasticsearchVolumeExpansionTopologySpec `json:"topology,omitempty" protobuf:"bytes,2,opt,name=topology"`
}

// ElasticsearchVolumeExpansionTopologySpec is the spec for Elasticsearch volume expansion in topology mode
type ElasticsearchVolumeExpansionTopologySpec struct {
	// volume specification for master nodes
	Master *resource.Quantity `json:"master,omitempty" protobuf:"bytes,1,opt,name=master"`
	// volume specification for data nodes
	Data *resource.Quantity `json:"data,omitempty" protobuf:"bytes,2,opt,name=data"`
	// volume specification for ingest nodes
	Ingest *resource.Quantity `json:"ingest,omitempty" protobuf:"bytes,3,opt,name=ingest"`
}

type ElasticsearchCustomConfigurationSpec struct {
}

type ElasticsearchCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty" protobuf:"bytes,1,opt,name=configMap"`
	Data      map[string]string          `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`
	Remove    bool                       `json:"remove,omitempty" protobuf:"varint,3,opt,name=remove"`
}

// ElasticsearchOpsRequestStatus is the status for ElasticsearchOpsRequest
type ElasticsearchOpsRequestStatus struct {
	// Specifies the current phase of the ops request
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
