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
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeQdrant     = "qd"
	ResourceKindQdrant     = "Qdrant"
	ResourceSingularQdrant = "qdrant"
	ResourcePluralQdrant   = "qdrants"
)

// +kubebuilder:validation:Enum=Standalone;Distributed
type QdrantMode string

const (
	QdrantStandalone  QdrantMode = "Standalone"
	QdrantDistributed QdrantMode = "Distributed"
)

// +kubebuilder:validation:Enum=server;client
type QdrantCertificateAlias string

const (
	QdrantServerCert QdrantCertificateAlias = "server"
	QdrantClientCert QdrantCertificateAlias = "client"
)

// Qdrant is the Schema for the Qdrant API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=qdrants,singular=qdrant,shortName=qd,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Qdrant struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrantSpec   `json:"spec,omitempty"`
	Status QdrantStatus `json:"status,omitempty"`
}

// QdrantSpec defines the desired state of Qdrant.
type QdrantSpec struct {
	// Version of Qdrant to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Qdrant database.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Qdrant cluster mode
	// +optional
	Mode QdrantMode `json:"mode,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// disable security. It disables authentication security of user.
	// If unset, default is false
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// +optional
	Configuration *ConfigurationSpec `json:"configuration,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// TLS contains tls configurations for client and server.
	TLS *QdrantTLSConfig `json:"tls,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`
}

type QdrantTLSConfig struct {
	kmapi.TLSConfig `json:",inline"`
	// +optional
	P2P *bool `json:"p2p"`
	// +optional
	Client *bool `json:"client"`
}

// QdrantStatus defines the observed state of Qdrant.
type QdrantStatus struct {
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QdrantList contains a list of Qdrant.
type QdrantList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []Qdrant `json:"items"`
}
