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
	ResourceCodeKafkaOpsRequest     = "kfops"
	ResourceKindKafkaOpsRequest     = "KafkaOpsRequest"
	ResourceSingularKafkaOpsRequest = "kafkaopsrequest"
	ResourcePluralKafkaOpsRequest   = "kafkaopsrequests"
)

// kafkaDBOpsRequest defines a Kafka DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=kafkaopsrequests,singular=kafkaopsrequest,shortName=kfops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type KafkaOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KafkaOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus    `json:"status,omitempty"`
}

// KafkaOpsRequestSpec is the spec for KafkaOpsRequest
type KafkaOpsRequestSpec struct {
	// Specifies the Kafka reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type KafkaOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading Kafka
	UpdateVersion *KafkaUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *KafkaHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *KafkaVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *KafkaVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of Kafka
	Configuration *KafkaCustomConfigurationSpec `json:"configuration,omitempty"`
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

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS)
type KafkaOpsRequestType string

// KafkaReplicaReadinessCriteria is the criteria for checking readiness of a Kafka pod
// after updating, horizontal scaling etc.
type KafkaReplicaReadinessCriteria struct{}

// KafkaUpdateVersionSpec contains the update version information of a kafka cluster
type KafkaUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// KafkaHorizontalScalingSpec contains the horizontal scaling information of a Kafka cluster
type KafkaHorizontalScalingSpec struct {
	// Number of combined (i.e. broker, controller) node
	Node *int32 `json:"node,omitempty"`
	// Node topology specification
	Topology *KafkaHorizontalScalingTopologySpec `json:"topology,omitempty"`
}

// KafkaHorizontalScalingTopologySpec contains the horizontal scaling information in cluster topology mode
type KafkaHorizontalScalingTopologySpec struct {
	// Number of broker nodes
	Broker *int32 `json:"broker,omitempty"`
	// Number of controller nodes
	Controller *int32 `json:"controller,omitempty"`
}

// KafkaVerticalScalingSpec contains the vertical scaling information of a Kafka cluster
type KafkaVerticalScalingSpec struct {
	// Resource spec for combined nodes
	Node *PodResources `json:"node,omitempty"`
	// Resource spec for broker
	Broker *PodResources `json:"broker,omitempty"`
	// Resource spec for controller
	Controller *PodResources `json:"controller,omitempty"`
}

// KafkaVolumeExpansionSpec is the spec for Kafka volume expansion
type KafkaVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for combined nodes
	Node *resource.Quantity `json:"node,omitempty"`
	// volume specification for broker
	Broker *resource.Quantity `json:"broker,omitempty"`
	// volume specification for controller
	Controller *resource.Quantity `json:"controller,omitempty"`
}

// KafkaCustomConfigurationSpec is the spec for Reconfiguring the Kafka Settings
type KafkaCustomConfigurationSpec struct {
	// ConfigSecret is an optional field to provide custom configuration file for database.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
	// ApplyConfig is an optional field to provide Kafka configuration.
	// Provided configuration will be applied to config files stored in ConfigSecret.
	// If the ConfigSecret is missing, the operator will create a new k8s secret by the
	// following naming convention: {db-name}-user-config .
	// Expected input format:
	//	applyConfig:
	//		file-name.properties: |
	//			key=value
	//		server.properties: |
	//			log.retention.ms=10000
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
	// If set to "true", the user provided configuration will be removed.
	// The Kafka cluster will start will default configuration that is generated by the operator.
	// +optional
	RemoveCustomConfig bool `json:"removeCustomConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KafkaOpsRequestList is a list of KafkaOpsRequests
type KafkaOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of KafkaOpsRequest CRD objects
	Items []KafkaOpsRequest `json:"items,omitempty"`
}
