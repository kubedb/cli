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
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeWeaviate     = "wv"
	ResourceKindWeaviate     = "Weaviate"
	ResourceSingularWeaviate = "weaviate"
	ResourcePluralWeaviate   = "weaviates"
)

// Weaviate is the Schema for the Weaviate API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=weaviates,singular=weaviate,shortName=wv,categories={datastore,vectordb,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Weaviate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WeaviateSpec   `json:"spec,omitempty"`
	Status WeaviateStatus `json:"status,omitempty"`
}

// WeaviateSpec defines the desired state of Weaviate.
type WeaviateSpec struct {
	// Version of Weaviate to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Weaviate database.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Replication configuration for the Weaviate cluster.
	// This controls the data replication factor per collection.
	// +optional
	Replication *ReplicationConfig `json:"replication,omitempty"`

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

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// Configuration is an optional field to provide custom configuration file for database (i.e conf.yaml).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	// You can provide custom configurations using Secret or ApplyConfig.
	// +optional
	Configuration *WeaviateConfiguration `json:"configuration,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *WeaviateTLSConfig `json:"tls,omitempty"`

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

	// Monitor is used to monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`
}

// WeaviateStatus defines the observed state of Weaviate.
type WeaviateStatus struct {
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

// WeaviateList contains a list of Weaviate.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type WeaviateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Weaviate `json:"items"`
}

// ReplicationConfig defines replication settings for Weaviate.
type ReplicationConfig struct {
	// Factor is the number of replicas for each data object.
	// Set to 1 for no replication (default), 2-3 for production HA.
	// +optional
	// +kubebuilder:minimum=1
	// +kubebuilder:maximum=5
	Factor int32 `json:"factor,omitempty"`
}

type WeaviateTLSConfig struct {
	kmapi.TLSConfig `json:",inline"`

	// ClientAuth controls whether the REST HTTPS listener requires clients to present a valid certificate.
	// If unset, client certificate authentication is enabled for backward compatibility.
	// +optional
	ClientAuth *bool `json:"clientAuth,omitempty"`
}

// +kubebuilder:validation:Enum=server;client
type WeaviateCertificateAlias string

const (
	WeaviateServerCert WeaviateCertificateAlias = "server"
	WeaviateClientCert WeaviateCertificateAlias = "client"
)

type WeaviateConfiguration struct {
	ConfigurationSpec `json:",inline,omitempty"`

	// BackupConfigSecret is an optional field to provide environment variables
	// from a Kubernetes Secret for backup or other purposes.
	// These env vars will be injected into the database container.
	// +optional
	BackupConfigSecret *core.LocalObjectReference `json:"backupConfigSecret,omitempty"`
}

var _ Accessor = &Weaviate{}

func (w *Weaviate) GetObjectMeta() metav1.ObjectMeta {
	return w.ObjectMeta
}

func (w *Weaviate) GetConditions() []kmapi.Condition {
	return w.Status.Conditions
}

func (w *Weaviate) SetCondition(cond kmapi.Condition) {
	w.Status.Conditions = setCondition(w.Status.Conditions, cond)
}

func (w *Weaviate) RemoveCondition(typ string) {
	w.Status.Conditions = removeCondition(w.Status.Conditions, typ)
}
