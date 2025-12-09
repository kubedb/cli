/*
Copyright 2025.

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
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeMilvus     = "mv"
	ResourceKindMilvus     = "Milvus"
	ResourceSingularMilvus = "milvus"
	ResourcePluralMilvus   = "milvuses"
)

// +kubebuilder:validation:Enum=Standalone;Distributed
type MilvusMode string

// Package v1alpha2 contains API Schema definitions for the  v1alpha2 API group.
// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=milvuses,singular=milvus,shortName=mv,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

type Milvus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MilvusSpec   `json:"spec,omitempty"`
	Status MilvusStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen=true
// MilvusSpec defines the desired state of Milvus
type MilvusSpec struct {
	// Version of Milvus to be deployed
	Version string `json:"version"`

	// Meta contains configuration for etcd meta storage
	MetaStorage *MetaStorageSpec `json:"metaStorage,omitempty"`

	// ObjectStorage contains specification for druid to connect to the object storage
	ObjectStorage *ObjectStorageSpec `json:"objectStorage"`

	// Milvus cluster topology
	// +optional
	Topology *MilvusTopology `json:"topology,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

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

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e config.properties).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`
}

// +k8s:deepcopy-gen=true
type MilvusTopology struct {
	// If set to -
	// "Standalone", Standalone is required, and Milvus will start a Standalone Mode
	// "Distributed", DistributedSpec is required, and Milvus will start a Distributed Mode
	Mode *MilvusMode `json:"mode,omitempty"`
}

// +k8s:deepcopy-gen=true
// Meta Storage defines the configuration for etcd meta storage
type MetaStorageSpec struct {
	// ExternallyManaged indicates whether etcd is managed outside this operator.
	// If true, only endpoints are used. If false, an EtcdCluster CR is created.
	// +optional
	ExternallyManaged bool `json:"externallyManaged,omitempty"`

	// Endpoints are the client endpoints of etcd (e.g., ["http://etcd-svc:2379"]).
	// Required when ExternallyManaged=true.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Endpoints []string `json:"endpoints,omitempty"`

	// Size is the expected size of the cluster.
	// Required when ExternallyManaged=false. Ignored otherwise.
	// +kubebuilder:validation:Minimum=1
	// +optional
	Size int `json:"size,omitempty"`

	// StorageType can be durable (default) or ephemeral
	// +optional
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	// +optional
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
}

// +k8s:deepcopy-gen=true
// ObjectStorageStorageSpec defines the configuration for MinIO or S3 object storage
type ObjectStorageSpec struct {
	// ConfigSecret should contain the necessary data to connect to external MinIO
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
}

// +k8s:deepcopy-gen=true
// MilvusStatus defines the observed state of Milvus
type MilvusStatus struct {
	// Specifies the current phase of the database
	// +optional
	Phase DatabasePhase `json:"phase,omitempty"`

	// ObservedGeneration is the most recent generation observed for this resource
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// MilvusList contains a list of Milvus
type MilvusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Milvus `json:"items"`
}
