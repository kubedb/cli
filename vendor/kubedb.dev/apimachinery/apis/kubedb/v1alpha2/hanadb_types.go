/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed u`nder the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeHanaDB     = "hdb"
	ResourceKindHanaDB     = "HanaDB"
	ResourceSingularHanaDB = "hanadb"
	ResourcePluralHanaDB   = "hanadbs"
)

// +kubebuilder:validation:Enum=Standalone;SystemReplication
type HanaDBMode string

const (
	HanaDBModeStandalone        HanaDBMode = "Standalone"
	HanaDBModeSystemReplication HanaDBMode = "SystemReplication"
)

// HanaDB is the Schema for the hanadbs API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=hanadbs,singular=hanadb,shortName=hdb,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type HanaDB struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of HanaDB
	// +required
	Spec HanaDBSpec `json:"spec"`

	// status defines the observed state of HanaDB
	// +optional
	Status HanaDBStatus `json:"status,omitempty,omitzero"`
}

// HanaDBSpec defines the desired state of HanaDB
type HanaDBSpec struct {
	// Version of HanaDB to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a HanaDB database
	Replicas *int32 `json:"replicas,omitempty"`

	// Topology defines the deployment mode (e.g., standalone or system replication).
	// +optional
	Topology *HanaDBTopology `json:"topology,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// Configuration holds the custom config for hanadb
	// +optional
	Configuration *ConfigurationSpec `json:"configuration,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

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
}

// HanaDBTopology defines the deployment mode for HanaDB
type HanaDBTopology struct {
	// Mode specifies the deployment mode.
	// +optional
	Mode *HanaDBMode `json:"mode,omitempty"`
}

// HanaDBStatus defines the observed state of HanaDB.
type HanaDBStatus struct {
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
// HanaDBList contains a list of HanaDB
type HanaDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HanaDB `json:"items"`
}
