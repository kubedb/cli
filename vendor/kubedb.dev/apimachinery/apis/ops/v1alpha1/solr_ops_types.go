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
	ResourceCodeSolrOpsRequest     = "slops"
	ResourceKindSolrOpsRequest     = "SolrOpsRequest"
	ResourceSingularSolrOpsRequest = "solropsrequest"
	ResourcePluralSolrOpsRequest   = "solropsrequests"
)

// SolrDBOpsRequest defines a Solr DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=solropsrequests,singular=solropsrequest,shortName=slops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SolrOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SolrOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus   `json:"status,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;VerticalScaling;VolumeExpansion;Reconfigure;Restart
// ENUM(UpdateVersion, VerticalScaling, VolumeExpansion, Reconfigure, Restart)
type SolrOpsRequestType string

// DruidOpsRequestSpec is the spec for DruidOpsRequest
type SolrOpsRequestSpec struct {
	// Specifies the Druid reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type SolrOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading Solr
	UpdateVersion *SolrUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *SolrVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *SolrVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Specifies information necessary for custom configuration of solr
	Configuration *SolrCustomConfigurationSpec `json:"configuration,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

type SolrVerticalScalingSpec struct {
	// Resource spec for combined nodes
	Node *PodResources `json:"node,omitempty"`
	// Resource spec for data nodes
	Data *PodResources `json:"data,omitempty"`
	// Resource spec for overseer nodes
	Overseer *PodResources `json:"overseer,omitempty"`
	// Resource spec for overseer nodes
	Coordinator *PodResources `json:"coordinator,omitempty"`
}

type SolrVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for combined nodes
	Node *resource.Quantity `json:"node,omitempty"`
	// volume specification for data nodes
	Data *resource.Quantity `json:"data,omitempty"`
	// volume specification for overseer nodes
	Overseer *resource.Quantity `json:"overseer,omitempty"`
	// volume specification for overseer nodes
	Coordinator *resource.Quantity `json:"coordinator,omitempty"`
}

type SolrUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// SolrCustomConfigurationSpec is the spec for Reconfiguring the solr Settings
type SolrCustomConfigurationSpec struct {
	// ConfigSecret is an optional field to provide custom configuration file for database.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
	// ApplyConfig is an optional field to provide solr configuration.
	// Provided configuration will be applied to config files stored in ConfigSecret.
	// If the ConfigSecret is missing, the operator will create a new k8s secret by the
	// following naming convention: {db-name}-user-config .
	// Expected input format:
	//	applyConfig:
	//		solr.xml: |
	//			key=value
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
	// If set to "true", the user provided configuration will be removed.
	// The solr cluster will start will default configuration that is generated by the operator.
	// +optional
	RemoveCustomConfig bool `json:"removeCustomConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SolrOpsRequestList is a list of DruidOpsRequests
type SolrOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of SolrOpsRequest CRD objects
	Items []SolrOpsRequest `json:"items,omitempty"`
}
