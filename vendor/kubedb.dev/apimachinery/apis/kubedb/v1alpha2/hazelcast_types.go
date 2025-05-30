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

package v1alpha2

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeHazelcast     = "hz"
	ResourceKindHazelcast     = "Hazelcast"
	ResourceSingularHazelcast = "hazelcast"
	ResourcePluralHazelcast   = "hazelcasts"
)

// Hazelcast is the Schema for the hazelcasts API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=hazelcasts,singular=hazelcast,shortName=hz,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Hazelcast struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HazelcastSpec   `json:"spec,omitempty"`
	Status HazelcastStatus `json:"status,omitempty"`
}
type HazelcastSpec struct {
	// Version of Hazelcast to be deployed
	Version string `json:"version"`

	// lisence for hazelcast community edition
	LicenseSecret *core.SecretReference `json:"licenseSecret,omitempty"`

	// Number of instances to deploy for a Hazelcast database
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// StorageType van be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Custom java environment variables to be integrated
	// +optional
	JavaOpts []string `json:"javaOpts,omitempty"`

	// Disable security. It disables authentication security of users.
	// If unset, default is false
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 20, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// To enable ssl for http layer
	EnableSSL bool `json:"enableSSL,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// Keystore encryption secret
	// +optional
	KeystoreSecret *core.LocalObjectReference `json:"keystoreSecret,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`
}

// HazelcastStatus defines the observed state of Hazelcast.
type HazelcastStatus struct {
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

// +kubebuilder:validation:Enum=ca;transport;http;client;server
type HazelcastCertificateAlias string

const (
	HazelcastCACert        HazelcastCertificateAlias = "ca"
	HazelcastTransportCert HazelcastCertificateAlias = "transport"
	HazelcastHTTPCert      HazelcastCertificateAlias = "http"
	HazelcastClientCert    HazelcastCertificateAlias = "client"
	HazelcastServerCert    HazelcastCertificateAlias = "server"
)

// HazelcastList contains a list of Hazelcast.

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HazelcastList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Hazelcast `json:"items"`
}
