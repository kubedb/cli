/*
Copyright 2023.

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
	ResourceCodeFerretDB     = "fr"
	ResourceKindFerretDB     = "FerretDB"
	ResourceSingularFerretDB = "ferretdb"
	ResourcePluralFerretDB   = "ferretdbs"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ferretdbs,singular=ferretdb,shortName=fr,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Namespace",type="string",JSONPath=".metadata.namespace"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type FerretDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FerretDBSpec   `json:"spec,omitempty"`
	Status FerretDBStatus `json:"status,omitempty"`
}

type FerretDBSpec struct {
	// Version of FerretDB to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a FerretDB database.
	Replicas *int32 `json:"replicas,omitempty"`

	// Database authentication secret.
	// Use this only when backend is internally managed.
	// For externally managed backend, we will get the authSecret from AppBinding
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// See more options: https://docs.ferretdb.io/security/tls-connections/
	// +optional
	SSLMode SSLMode `json:"sslMode,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// StorageType can be durable (default) or ephemeral for KubeDB Backend
	// +optional
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used for KubeDB Backend.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy TerminationPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Monitor is used monitor database instance and KubeDB Backend
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	Backend *FerretDBBackend `json:"backend"`
}

type FerretDBStatus struct {
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

type FerretDBBackend struct {
	// PostgresRef refers to the AppBinding of the backend Postgres server
	// +optional
	PostgresRef *kmapi.ObjectReference `json:"postgresRef,omitempty"`
	// Which versions pg will be used as backend of ferretdb. default 13.13 when backend internally managed
	// +optional
	Version *string `json:"version,omitempty"`
	// A DB inside backend specifically made for ferretdb
	// +optional
	LinkedDB          string `json:"linkedDB,omitempty"`
	ExternallyManaged bool   `json:"externallyManaged"`
}

// +kubebuilder:validation:Enum=server;client
type FerretDBCertificateAlias string

const (
	FerretDBServerCert FerretDBCertificateAlias = "server"
	FerretDBClientCert FerretDBCertificateAlias = "client"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FerretDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FerretDB `json:"items"`
}
