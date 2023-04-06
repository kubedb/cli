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
	ResourceCodePerconaXtraDBOpsRequest     = "pxcops"
	ResourceKindPerconaXtraDBOpsRequest     = "PerconaXtraDBOpsRequest"
	ResourceSingularPerconaXtraDBOpsRequest = "perconaxtradbopsrequest"
	ResourcePluralPerconaXtraDBOpsRequest   = "perconaxtradbopsrequests"
)

// PerconaXtraDBOpsRequest defines a PerconaXtraDB (percona variation for MySQL database) DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=perconaxtradbopsrequests,singular=perconaxtradbopsrequest,shortName=pxcops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PerconaXtraDBOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PerconaXtraDBOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus            `json:"status,omitempty"`
}

// PerconaXtraDBOpsRequestSpec is the spec for PerconaXtraDBOpsRequest
type PerconaXtraDBOpsRequestSpec struct {
	// Specifies the PerconaXtraDB reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type PerconaXtraDBOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading PerconaXtraDB
	UpdateVersion *PerconaXtraDBUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for upgrading PerconaXtraDB
	// Deprecated: use UpdateVersion
	Upgrade *PerconaXtraDBUpdateVersionSpec `json:"upgrade,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *PerconaXtraDBHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *PerconaXtraDBVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *PerconaXtraDBVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of PerconaXtraDB
	Configuration *PerconaXtraDBCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *PerconaXtraDBTLSSpec `json:"tls,omitempty"`
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
type PerconaXtraDBOpsRequestType string

// PerconaXtraDBReplicaReadinessCriteria is the criteria for checking readiness of an PerconaXtraDB database
type PerconaXtraDBReplicaReadinessCriteria struct{}

type PerconaXtraDBUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

type PerconaXtraDBHorizontalScalingSpec struct {
	// Number of nodes/members of the group
	Member *int32 `json:"member,omitempty"`
	// specifies the weight of the current member/Node
	MemberWeight int32 `json:"memberWeight,omitempty"`
}

type PerconaXtraDBVerticalScalingSpec struct {
	PerconaXtraDB *core.ResourceRequirements `json:"perconaxtradb,omitempty"`
	Exporter      *core.ResourceRequirements `json:"exporter,omitempty"`
	Coordinator   *core.ResourceRequirements `json:"coordinator,omitempty"`
}

// PerconaXtraDBVolumeExpansionSpec is the spec for PerconaXtraDB volume expansion
type PerconaXtraDBVolumeExpansionSpec struct {
	PerconaXtraDB *resource.Quantity   `json:"perconaxtradb,omitempty"`
	Mode          *VolumeExpansionMode `json:"mode,omitempty"`
}

type PerconaXtraDBCustomConfigurationSpec struct {
	// ConfigSecret is an optional field to provide custom configuration file for database.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
	// Deprecated
	InlineConfig string `json:"inlineConfig,omitempty"`
	// If set to "true", the user provided configuration will be removed.
	// PerconaXtraDB will start will default configuration that is generated by the operator.
	// +optional
	RemoveCustomConfig bool `json:"removeCustomConfig,omitempty"`
	// ApplyConfig is an optional field to provide PerconaXtraDB configuration.
	// Provided configuration will be applied to config files stored in ConfigSecret.
	// If the ConfigSecret is missing, the operator will create a new k8s secret by the
	// following naming convention: {db-name}-user-config .
	// Expected input format:
	//	applyConfig:
	//		file-name.cnf: |
	//			[mysqld]
	//			key1: value1
	//			key2: value2
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
}

// PerconaXtraDBTLSSpec specifies information necessary for configuring TLS
type PerconaXtraDBTLSSpec struct {
	TLSSpec `json:",inline,omitempty"`

	// Indicates that the database server need to be encrypted connections(ssl)
	// +optional
	RequireSSL *bool `json:"requireSSL,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PerconaXtraDBOpsRequestList is a list of PerconaXtraDBOpsRequests
type PerconaXtraDBOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PerconaXtraDBOpsRequest CRD objects
	Items []PerconaXtraDBOpsRequest `json:"items,omitempty"`
}
