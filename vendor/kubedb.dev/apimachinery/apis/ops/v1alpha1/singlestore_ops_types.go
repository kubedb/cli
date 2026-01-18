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
	ResourceCodeSinglestoreOpsRequest     = "sdbops"
	ResourceKindSinglestoreOpsRequest     = "SinglestoreOpsRequest"
	ResourceSingularSinglestoreOpsRequest = "singlestoreopsrequest"
	ResourcePluralSinglestoreOpsRequest   = "singlestoreopsrequests"
)

// SinglestoreOpsRequest defines a Singlestore DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=singlestoreopsrequests,singular=singlestoreopsrequest,shortName=sdbops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SinglestoreOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SinglestoreOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus          `json:"status,omitempty"`
}

// SinglestoreOpsRequestSpec is the spec for SinglestoreOpsRequest
type SinglestoreOpsRequestSpec struct {
	// Specifies the Singlestore reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type SinglestoreOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading SingleStore Version
	UpdateVersion *SinglestoreUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *SinglestoreHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *SinglestoreVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *SinglestoreVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Specifies information necessary for custom configuration of Singlestore
	Configuration *SinglestoreCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
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
type SinglestoreOpsRequestType string

// SinglestoreUpdateVersionSpec contains the update version information of a kafka cluster
type SinglestoreUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// SinglestoreHorizontalScalingSpec contains the horizontal scaling information of a Singlestore cluster
type SinglestoreHorizontalScalingSpec struct {
	// number of Aggregator node
	Aggregator *int32 `json:"aggregator,omitempty"`
	// number of Leaf node
	Leaf *int32 `json:"leaf,omitempty"`
}

// SinglestoreVerticalScalingSpec contains the vertical scaling information of a Singlestore cluster
type SinglestoreVerticalScalingSpec struct {
	// Resource spec for standalone node
	Node *PodResources `json:"node,omitempty"`
	// Resource spec for Aggregator
	Aggregator *PodResources `json:"aggregator,omitempty"`
	// Resource spec for Leaf
	Leaf *PodResources `json:"leaf,omitempty"`
	// Resource spec for Coordinator container
	Coordinator *ContainerResources `json:"coordinator,omitempty"`
}

// SinglestoreVolumeExpansionSpec is the spec for Singlestore volume expansion
type SinglestoreVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// Volume specification for standalone
	Node *resource.Quantity `json:"node,omitempty"`
	// Volume specification for Aggregator
	Aggregator *resource.Quantity `json:"aggregator,omitempty"`
	// Volume specification for Leaf
	Leaf *resource.Quantity `json:"leaf,omitempty"`
}

// SinglestoreCustomConfigurationSpec is the spec for Singlestore reconfiguration
type SinglestoreCustomConfigurationSpec struct {
	// Custom Configuration specification for standalone
	Node *ReconfigurationSpec `json:"node,omitempty"`
	// Custom Configuration specification for Aggregator
	Aggregator *ReconfigurationSpec `json:"aggregator,omitempty"`
	// Custom Configuration specification for Leaf
	Leaf *ReconfigurationSpec `json:"leaf,omitempty"`
}

type SinglestoreCustomConfiguration struct {
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty"`
	ApplyConfig        map[string]string          `json:"applyConfig,omitempty"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SinglestoreOpsRequestList is a list of SinglestoreOpsRequests
type SinglestoreOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of SinglestoreOpsRequest CRD objects
	Items []SinglestoreOpsRequest `json:"items,omitempty"`
}
