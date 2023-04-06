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

//go:generate go-enum --mustparse --names --values
package v1alpha1

import (
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ElasticsearchOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus            `json:"status,omitempty"`
}

// ElasticsearchOpsRequestSpec is the spec for ElasticsearchOpsRequest
type ElasticsearchOpsRequestSpec struct {
	// Specifies the Elasticsearch reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type ElasticsearchOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading Elasticsearch
	UpdateVersion *ElasticsearchUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for upgrading Elasticsearch
	// Deprecated: use UpdateVersion
	Upgrade *ElasticsearchUpdateVersionSpec `json:"upgrade,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *ElasticsearchHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *ElasticsearchVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *ElasticsearchVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of Elasticsearch
	Configuration *ElasticsearchCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=Upgrade;UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS
// ENUM(Upgrade, UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS)
type ElasticsearchOpsRequestType string

// ElasticsearchReplicaReadinessCriteria is the criteria for checking readiness of an Elasticsearch database
type ElasticsearchReplicaReadinessCriteria struct{}

type ElasticsearchUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// ElasticsearchHorizontalScalingSpec contains the horizontal scaling information of an Elasticsearch cluster
type ElasticsearchHorizontalScalingSpec struct {
	// Number of combined (i.e. master, data, ingest) node
	Node *int32 `json:"node,omitempty"`
	// Node topology specification
	Topology *ElasticsearchHorizontalScalingTopologySpec `json:"topology,omitempty"`
}

// ElasticsearchHorizontalScalingTopologySpec contains the horizontal scaling information in cluster topology mode
type ElasticsearchHorizontalScalingTopologySpec struct {
	// Number of master nodes
	Master *int32 `json:"master,omitempty"`
	// Number of ingest nodes
	Ingest *int32 `json:"ingest,omitempty"`
	// Number of data nodes
	Data         *int32 `json:"data,omitempty"`
	DataContent  *int32 `json:"dataContent,omitempty"`
	DataHot      *int32 `json:"dataHot,omitempty"`
	DataWarm     *int32 `json:"dataWarm,omitempty"`
	DataCold     *int32 `json:"dataCold,omitempty"`
	DataFrozen   *int32 `json:"dataFrozen,omitempty"`
	ML           *int32 `json:"ml,omitempty"`
	Transform    *int32 `json:"transform,omitempty"`
	Coordinating *int32 `json:"coordinating,omitempty"`
}

// ElasticsearchVerticalScalingSpec is the spec for Elasticsearch vertical scaling
type ElasticsearchVerticalScalingSpec struct {
	// Resource spec for combined nodes
	Node *core.ResourceRequirements `json:"node,omitempty"`
	// Resource spec for exporter sidecar
	Exporter *core.ResourceRequirements `json:"exporter,omitempty"`
	// Specifies the resource spec for cluster in topology mode
	Topology *ElasticsearchVerticalScalingTopologySpec `json:"topology,omitempty"`
}

// ElasticsearchVerticalScalingTopologySpec is the resource spec in the cluster topology mode
type ElasticsearchVerticalScalingTopologySpec struct {
	Master       *core.ResourceRequirements `json:"master,omitempty"`
	Ingest       *core.ResourceRequirements `json:"ingest,omitempty"`
	Data         *core.ResourceRequirements `json:"data,omitempty"`
	DataContent  *core.ResourceRequirements `json:"dataContent,omitempty"`
	DataHot      *core.ResourceRequirements `json:"dataHot,omitempty"`
	DataWarm     *core.ResourceRequirements `json:"dataWarm,omitempty"`
	DataCold     *core.ResourceRequirements `json:"dataCold,omitempty"`
	DataFrozen   *core.ResourceRequirements `json:"dataFrozen,omitempty"`
	ML           *core.ResourceRequirements `json:"ml,omitempty"`
	Transform    *core.ResourceRequirements `json:"transform,omitempty"`
	Coordinating *core.ResourceRequirements `json:"coordinating,omitempty"`
}

// ElasticsearchVolumeExpansionSpec is the spec for Elasticsearch volume expansion
type ElasticsearchVolumeExpansionSpec struct {
	// +kubebuilder:default="Online"
	Mode *VolumeExpansionMode `json:"mode,omitempty"`
	// volume specification for combined nodes
	Node *resource.Quantity `json:"node,omitempty"`
	// volume specification for nodes in cluster topology
	Topology *ElasticsearchVolumeExpansionTopologySpec `json:"topology,omitempty"`
}

// ElasticsearchVolumeExpansionTopologySpec is the spec for Elasticsearch volume expansion in topology mode
type ElasticsearchVolumeExpansionTopologySpec struct {
	// volume specification for master nodes
	Master *resource.Quantity `json:"master,omitempty"`
	// volume specification for ingest nodes
	Ingest *resource.Quantity `json:"ingest,omitempty"`
	// volume specification for data nodes
	Data         *resource.Quantity `json:"data,omitempty"`
	DataContent  *resource.Quantity `json:"dataContent,omitempty"`
	DataHot      *resource.Quantity `json:"dataHot,omitempty"`
	DataWarm     *resource.Quantity `json:"dataWarm,omitempty"`
	DataCold     *resource.Quantity `json:"dataCold,omitempty"`
	DataFrozen   *resource.Quantity `json:"dataFrozen,omitempty"`
	ML           *resource.Quantity `json:"ml,omitempty"`
	Transform    *resource.Quantity `json:"transform,omitempty"`
	Coordinating *resource.Quantity `json:"coordinating,omitempty"`
}

// ElasticsearchCustomConfigurationSpec is the spec for Reconfiguring the Elasticsearch Settings
type ElasticsearchCustomConfigurationSpec struct {
	// ConfigSecret is an optional field to provide custom configuration file for database.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
	// SecureConfigSecret is an optional field to provide secure settings for database.
	//	- Ref: https://www.elastic.co/guide/en/elasticsearch/reference/7.14/secure-settings.html
	// +optional
	SecureConfigSecret *core.LocalObjectReference `json:"secureConfigSecret,omitempty"`
	// ApplyConfig is an optional field to provide Elasticsearch configuration.
	// Provided configuration will be applied to config files stored in ConfigSecret.
	// If the ConfigSecret is missing, the operator will create a new k8s secret by the
	// following naming convention: {db-name}-user-config .
	// Expected input format:
	//	applyConfig:
	//		file-name.yml: |
	//			key: value
	//		elasticsearch.yml: |
	//			thread_pool:
	//				write:
	//					size: 30
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
	// If set to "true", the user provided configuration will be removed.
	// The Elasticsearch cluster will start will default configuration that is generated by the operator.
	// +optional
	RemoveCustomConfig bool `json:"removeCustomConfig,omitempty"`
	// If set to "true", the user provided secure settings will be removed.
	// The elasticsearch.keystore will start will default password (i.e. "").
	// +optional
	RemoveSecureCustomConfig bool `json:"removeSecureCustomConfig,omitempty"`
}

type ElasticsearchCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty"`
	Data      map[string]string          `json:"data,omitempty"`
	Remove    bool                       `json:"remove,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElasticsearchOpsRequestList is a list of ElasticsearchOpsRequests
type ElasticsearchOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ElasticsearchOpsRequest CRD objects
	Items []ElasticsearchOpsRequest `json:"items,omitempty"`
}
