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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeNeo4jOpsRequest     = "neoops"
	ResourceKindNeo4jOpsRequest     = "Neo4jOpsRequest"
	ResourceSingularNeo4jOpsRequest = "neo4jopsrequest"
	ResourcePluralNeo4jOpsRequest   = "neo4jopsrequests"
)

// Neo4jDBOpsRequest defines a Neo4j DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=neo4jopsrequests,singular=neo4jopsrequest,shortName=neoops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Neo4jOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Neo4jOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus    `json:"status,omitempty"`
}

// Neo4jOpsRequestSpec is the spec for Neo4jOpsRequest
type Neo4jOpsRequestSpec struct {
	// Specifies the Neo4j reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type Neo4jOpsRequestType `json:"type"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *Neo4jTLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for custom configuration of Neo4j
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *Neo4jHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *Neo4jVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *Neo4jVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for upgrading Neo4j
	UpdateVersion *Neo4jUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// Neo4jVolumeExpansionSpec is the spec for Neo4j volume expansion
type Neo4jVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for servers
	Server *resource.Quantity `json:"server,omitempty"`
}

// Neo4jUpdateVersionSpec contains the update version information of a Neo4j cluster
type Neo4jUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// ReallocateStrategy defines how reallocation should be performed
type ReallocateStrategy string

const (
	StrategyIncremental ReallocateStrategy = "incremental" // safe, batch by batch
	StrategyFull        ReallocateStrategy = "full"        // all at once
	StrategyNone        ReallocateStrategy = "none"        // no reallocation
)

// ReallocateConfig defines reallocation behaviour
type ReallocateConfig struct {
	// +kubebuilder:validation:Enum=incremental;full;none
	// +kubebuilder:default=incremental
	Strategy ReallocateStrategy `json:"strategy,omitempty"`
	// only used when Strategy == incremental
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:default=1
	// +optional
	BatchSize *int32 `json:"batchSize,omitempty"`
}

// Neo4jHorizontalScalingSpec contains the horizontal scaling information of a Neo4j cluster
type Neo4jHorizontalScalingSpec struct {
	// Number of server
	// +kubebuilder:validation:Minimum=1
	Server *int32 `json:"server,omitempty"`
	// how to handle reallocation after scaling
	// +optional
	// +kubebuilder:default={strategy: "incremental", batchSize: 1}
	Reallocate *ReallocateConfig `json:"reallocate,omitempty"`
}

type Neo4jVerticalScalingSpec struct {
	// Resource spec for neo4j servers
	Server *PodResources `json:"server,omitempty"`
}

type Neo4jTLSSpec struct {
	// Neo4jTLSSpec contains updated tls configurations for client and server.
	// +optional
	dbapi.Neo4jTLSConfig `json:",inline,omitempty"`

	// RotateCertificates tells operator to initiate certificate rotation
	// +optional
	RotateCertificates bool `json:"rotateCertificates,omitempty"`

	// Remove tells operator to remove TLS configuration
	// +optional
	Remove bool `json:"remove,omitempty"`
}

// +kubebuilder:validation:Enum=Restart;ReconfigureTLS;RotateAuth;Reconfigure;VerticalScaling;HorizontalScaling;VolumeExpansion;UpdateVersion
// ENUM(Restart,ReconfigureTLS,RotateAuth,Reconfigure,HorizontalScaling,VerticalScaling,VolumeExpansion,UpdateVersion)
type Neo4jOpsRequestType string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Neo4jOpsRequestList is a list of Neo4jOpsRequests
type Neo4jOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Neo4jOpsRequest CRD objects
	Items []Neo4jOpsRequest `json:"items,omitempty"`
}
