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
	ResourceCodeRabbitMQOpsRequest     = "rmops"
	ResourceKindRabbitMQOpsRequest     = "RabbitMQOpsRequest"
	ResourceSingularRabbitMQOpsRequest = "rabbitmqopsrequest"
	ResourcePluralRabbitMQOpsRequest   = "rabbitmqopsrequests"
)

// RabbitMQDBOpsRequest defines a RabbitMQ DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rabbitmqopsrequests,singular=rabbitmqopsrequest,shortName=rmops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type RabbitMQOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RabbitMQOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus       `json:"status,omitempty"`
}

// RabbitMQOpsRequestSpec is the spec for RabbitMQOpsRequest
type RabbitMQOpsRequestSpec struct {
	// Specifies the RabbitMQ reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type RabbitMQOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading rabbitmq
	UpdateVersion *RabbitMQUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *RabbitMQHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *RabbitMQVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *RabbitMQVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of rabbitmq
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
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
type RabbitMQOpsRequestType string

// RabbitMQUpdateVersionSpec contains the update version information of a rabbitmq cluster
type RabbitMQUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// RabbitMQReplicaReadinessCriteria is the criteria for checking readiness of a RabbitMQ pod
// after updating, horizontal scaling etc.
type RabbitMQReplicaReadinessCriteria struct{}

// RabbitMQHorizontalScalingSpec contains the horizontal scaling information of a RabbitMQ cluster
type RabbitMQHorizontalScalingSpec struct {
	// Number of node
	Node *int32 `json:"node,omitempty"`
}

// RabbitMQVerticalScalingSpec contains the vertical scaling information of a RabbitMQ cluster
type RabbitMQVerticalScalingSpec struct {
	// Resource spec for nodes
	Node *PodResources `json:"node,omitempty"`
}

// RabbitMQVolumeExpansionSpec is the spec for RabbitMQ volume expansion
type RabbitMQVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for nodes
	Node *resource.Quantity `json:"node,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RabbitMQOpsRequestList is a list of RabbitMQOpsRequests
type RabbitMQOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of RabbitMQOpsRequest CRD objects
	Items []RabbitMQOpsRequest `json:"items,omitempty"`
}
