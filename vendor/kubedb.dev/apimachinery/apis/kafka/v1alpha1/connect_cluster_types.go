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

package v1alpha1

import (
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeConnectCluster     = "kcc"
	ResourceKindConnectCluster     = "ConnectCluster"
	ResourceSingularConnectCluster = "connectcluster"
	ResourcePluralConnectCluster   = "connectclusters"
)

// ConnectCluster defines a framework for connecting Kafka with external systems

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=kcc,scope=Namespaced
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ConnectCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectClusterSpec   `json:"spec,omitempty"`
	Status ConnectClusterStatus `json:"status,omitempty"`
}

// ConnectClusterSpec defines the desired state of ConnectCluster
type ConnectClusterSpec struct {
	// Version of ConnectCluster to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Kafka database.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Kafka app-binding reference
	// KafkaRef is a required field, where ConnectCluster will store its metadata
	KafkaRef *kmapi.ObjectReference `json:"kafkaRef"`

	// disable security. It disables authentication security of user.
	// If unset, default is false
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// kafka connect cluster authentication secret
	// +optional
	AuthSecret *dbapi.SecretReference `json:"authSecret,omitempty"`

	// To enable https
	EnableSSL bool `json:"enableSSL,omitempty"`

	// Keystore encryption secret
	// +optional
	KeystoreCredSecret *dbapi.SecretReference `json:"keystoreCredSecret,omitempty"`

	// TLS contains tls configurations
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// List of connector-plugins
	ConnectorPlugins []string `json:"connectorPlugins,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for kafka connect cluster (i.e distributed.properties).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []dbapi.NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy dbapi.DeletionPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 20, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`
}

// ConnectClusterStatus defines the observed state of ConnectCluster
type ConnectClusterStatus struct {
	// Specifies the current phase of the database
	// +optional
	Phase ConnectClusterPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:validation:Enum=Provisioning;Ready;NotReady;Critical;Unknown
type ConnectClusterPhase string

const (
	ConnectClusterPhaseProvisioning ConnectClusterPhase = "Provisioning"
	ConnectClusterPhaseReady        ConnectClusterPhase = "Ready"
	ConnectClusterPhaseNotReady     ConnectClusterPhase = "NotReady"
	ConnectClusterPhaseCritical     ConnectClusterPhase = "Critical"
	ConnectClusterPhaseUnknown      ConnectClusterPhase = "Unknown"
)

// +kubebuilder:validation:Enum=ca;transport;http;client;server
type ConnectClusterCertificateAlias string

const (
	ConnectClusterCACert        ConnectClusterCertificateAlias = "ca"
	ConnectClusterTransportCert ConnectClusterCertificateAlias = "transport"
	ConnectClusterHTTPCert      ConnectClusterCertificateAlias = "http"
	ConnectClusterClientCert    ConnectClusterCertificateAlias = "client"
	ConnectClusterServerCert    ConnectClusterCertificateAlias = "server"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConnectClusterList contains a list of ConnectCluster
type ConnectClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConnectCluster `json:"items"`
}

// +kubebuilder:validation:Enum=statndalone;distributed
type ConnectClusterNodeRoleType string

const (
	ConnectClusterNodeRoleStandalone  ConnectClusterNodeRoleType = "standalone"
	ConnectClusterNodeRoleDistributed ConnectClusterNodeRoleType = "distributed"
)
