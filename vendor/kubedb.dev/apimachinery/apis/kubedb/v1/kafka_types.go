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

package v1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeKafka     = "kf"
	ResourceKindKafka     = "Kafka"
	ResourceSingularKafka = "kafka"
	ResourcePluralKafka   = "kafkas"
)

// Kafka is the Schema for the kafka API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:path=kafkas,singular=kafka,shortName=kf,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Kafka struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaSpec   `json:"spec,omitempty"`
	Status KafkaStatus `json:"status,omitempty"`
}

// KafkaSpec defines the desired state of Kafka
type KafkaSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of Kafka to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Kafka database.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Kafka topology for node specification
	// +optional
	Topology *KafkaClusterTopology `json:"topology,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// To enable ssl for http layer
	EnableSSL bool `json:"enableSSL,omitempty"`

	// disable security. It disables authentication security of user.
	// If unset, default is false
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e config.properties).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// Keystore encryption secret
	// +optional
	KeystoreCredSecret *SecretReference `json:"keystoreCredSecret,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// TLS contains tls configurations
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// CruiseControl is used to re-balance Kafka cluster
	// +optional
	CruiseControl *KafkaCruiseControl `json:"cruiseControl,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`
}

// KafkaClusterTopology defines kafka topology node specs for controller node and broker node
// dedicated controller nodes contains metadata for brokers and broker nodes contains data
// both nodes must be configured in topology mode
type KafkaClusterTopology struct {
	Controller *KafkaNode `json:"controller,omitempty"`
	Broker     *KafkaNode `json:"broker,omitempty"`
}

type KafkaNode struct {
	// Replicas represents number of replica for this specific type of node
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// suffix to append with node name
	// +optional
	Suffix string `json:"suffix,omitempty"`

	// Storage to specify how storage shall be used.
	// +optional
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`
}

// KafkaStatus defines the observed state of Kafka
type KafkaStatus struct {
	// Specifies the current phase of the database
	// +optional
	Phase DatabasePhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

type KafkaCruiseControl struct {
	// Configuration for cruise-control
	// +optional
	ConfigSecret *SecretReference `json:"configSecret,omitempty"`

	// Replicas represents number of replica for this specific type of node
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Suffix to append with node name
	// +optional
	Suffix string `json:"suffix,omitempty"`

	// Compute Resources required by the sidecar container.
	// +optional
	Resources core.ResourceRequirements `json:"resources,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	BrokerCapacity *KafkaBrokerCapacity `json:"brokerCapacity,omitempty"`
}

type KafkaBrokerCapacity struct {
	InBoundNetwork  string `json:"inBoundNetwork,omitempty"`
	OutBoundNetwork string `json:"outBoundNetwork,omitempty"`
}

// +kubebuilder:validation:Enum=controller;broker;combined
type KafkaNodeRoleType string

const (
	KafkaNodeRoleController KafkaNodeRoleType = "controller"
	KafkaNodeRoleBroker     KafkaNodeRoleType = "broker"
	KafkaNodeRoleCombined   KafkaNodeRoleType = "combined"
)

// +kubebuilder:validation:Enum=BROKER;CONTROLLER;INTERNAL;CC
type KafkaListenerType string

const (
	KafkaListenerBroker     KafkaListenerType = "BROKER"
	KafkaListenerController KafkaListenerType = "CONTROLLER"
	KafkaListenerLocal      KafkaListenerType = "LOCAL"
	KafkaListenerCC         KafkaListenerType = "CC"
)

// +kubebuilder:validation:Enum=ca;transport;http;client;server
type KafkaCertificateAlias string

const (
	KafkaCACert        KafkaCertificateAlias = "ca"
	KafkaTransportCert KafkaCertificateAlias = "transport"
	KafkaHTTPCert      KafkaCertificateAlias = "http"
	KafkaClientCert    KafkaCertificateAlias = "client"
	KafkaServerCert    KafkaCertificateAlias = "server"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KafkaList contains a list of Kafka
type KafkaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kafka `json:"items"`
}
