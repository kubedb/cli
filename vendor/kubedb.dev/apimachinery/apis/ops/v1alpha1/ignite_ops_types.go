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
	ResourceCodeIgniteOpsRequest     = "igops"
	ResourceKindIgniteOpsRequest     = "IgniteOpsRequest"
	ResourceSingularIgniteOpsRequest = "igniteopsrequest"
	ResourcePluralIgniteOpsRequest   = "igniteopsrequests"
)

// IgniteDBOpsRequest defines a Ignite DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=igniteopsrequests,singular=igniteopsrequest,shortName=igops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type IgniteOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              IgniteOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus     `json:"status,omitempty"`
}

// IgniteOpsRequestSpec is the spec for IgniteOpsRequest
type IgniteOpsRequestSpec struct {
	// Specifies the Ignite reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type IgniteOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading ignite
	UpdateVersion *IgniteUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *IgniteHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *IgniteVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *IgniteVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of ignite
	Configuration *IgniteCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;RotateAuth
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, RotateAuth)
type IgniteOpsRequestType string

// IgniteUpdateVersionSpec contains the update version information of a ignite cluster
type IgniteUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// IgniteHorizontalScalingSpec contains the horizontal scaling information of a Ignite cluster
type IgniteHorizontalScalingSpec struct {
	// Number of node
	Ignite *int32 `json:"node,omitempty"`
}

// IgniteVerticalScalingSpec contains the vertical scaling information of a Ignite cluster
type IgniteVerticalScalingSpec struct {
	// Resource spec for nodes
	Ignite *PodResources `json:"node,omitempty"`
}

// IgniteVolumeExpansionSpec is the spec for Ignite volume expansion
type IgniteVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for nodes
	Ignite *resource.Quantity `json:"node,omitempty"`
}

// IgniteCustomConfigurationSpec is the spec for Reconfiguring the ignite Settings
type IgniteCustomConfigurationSpec struct {
	// ConfigSecret is an optional field to provide custom configuration file for database.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
	// ApplyConfig is an optional field to provide ignite configuration.
	// Provided configuration will be applied to config files stored in ConfigSecret.
	// If the ConfigSecret is missing, the operator will create a new k8s secret by the
	// following naming convention: {db-name}-user-config .
	// Expected input format:
	//	applyConfig:
	//		node-configuration.xml: |
	//			key=value
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
	// If set to "true", the user provided configuration will be removed.
	// The ignite cluster will start will default configuration that is generated by the operator.
	// +optional
	RemoveCustomConfig bool `json:"removeCustomConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IgniteOpsRequestList is a list of IgniteOpsRequests
type IgniteOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of IgniteOpsRequest CRD objects
	Items []IgniteOpsRequest `json:"items,omitempty"`
}
