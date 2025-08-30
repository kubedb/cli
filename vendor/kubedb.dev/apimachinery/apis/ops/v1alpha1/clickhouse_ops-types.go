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
	olddbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeClickHouseOpsRequest     = "chops"
	ResourceKindClickHouseOpsRequest     = "ClickHouseOpsRequest"
	ResourceSingularClickHouseOpsRequest = "clickhouseopsrequest"
	ResourcePluralClickHouseOpsRequest   = "clickhouseopsrequests"
)

// ClickHouseDBOpsRequest defines a ClickHouse DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clickhouseopsrequests,singular=clickhouseopsrequest,shortName=chops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ClickHouseOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClickHouseOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus         `json:"status,omitempty"`
}

// ClickHouseOpsRequestSpec is the spec for ClickHouseOpsRequest
type ClickHouseOpsRequestSpec struct {
	// Specifies the ClickHouse reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type ClickHouseOpsRequestType `json:"type"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *ClickHouseHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for upgrading clickhouse
	UpdateVersion *ClickHouseUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *ClickHouseVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *ClickHouseVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of clickhouse
	Configuration *ClickHouseCustomConfigurationSpec `json:"configuration,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *ClickHouseTLSSpec `json:"tls,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=Restart;VerticalScaling;HorizontalScaling;UpdateVersion;VolumeExpansion;Reconfigure;ReconfigureTLS;RotateAuth
// ENUM(Restart, VerticalScaling, HorizontalScaling, UpdateVersion, VolumeExpansion, Reconfigure, ReconfigureTLS, RotateAuth)
type ClickHouseOpsRequestType string

// ClickHouseUpdateVersionSpec contains the update version information of a clickhouse cluster
type ClickHouseUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// ClickHouseVolumeExpansionSpec is the spec for ClickHouse volume expansion
type ClickHouseVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for nodes
	Node *resource.Quantity `json:"node,omitempty"`
}

// ClickHouseCustomConfigurationSpec is the spec for Reconfiguring the clickhouse Settings
type ClickHouseCustomConfigurationSpec struct {
	// ConfigSecret is an optional field to provide custom configuration file for database.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
	// ApplyConfig is an optional field to provide clickhouse configuration.
	// Provided configuration will be applied to config files stored in ConfigSecret.
	// If the ConfigSecret is missing, the operator will create a new k8s secret by the
	// following naming convention: {db-name}-user-config .
	// Expected input format:
	//	applyConfig:
	//		file-name.yml: |
	//			key: value
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
	// If set to "true", the user provided configuration will be removed.
	// The clickhouse cluster will start will default configuration that is generated by the operator.
	// +optional
	RemoveCustomConfig bool `json:"removeCustomConfig,omitempty"`
}

// ClickHouseTLSSpec contains necessary information for configuring TLS
type ClickHouseTLSSpec struct {
	// SSLVerificationMode specifies how server certificates should be verified
	SSLVerificationMode olddbapi.SSLVerificationMode `json:"sslVerificationMode,omitempty"`
	// TLSSpec holds the TLS certificate and key configuration.
	TLSSpec `json:",inline,omitempty"`
}

// ClickHouseVerticalScalingSpec contains the vertical scaling information of a clickhouse cluster
type ClickHouseVerticalScalingSpec struct {
	Node *PodResources `json:"node,omitempty"`
}

// ClickHouseHorizontalScalingSpec contains the horizontal scaling information of a clickhouse cluster
type ClickHouseHorizontalScalingSpec struct {
	// Number of node
	Replicas *int32 `json:"replicas,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClickHouseOpsRequestList is a list of ClickHouseOpsRequests
type ClickHouseOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ClickHouseOpsRequest CRD objects
	Items []ClickHouseOpsRequest `json:"items,omitempty"`
}
