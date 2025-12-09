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
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeNeo4j     = "neo"
	ResourceKindNeo4j     = "Neo4j"
	ResourceSingularNeo4j = "neo4j"
	ResourcePluralNeo4j   = "neo4js"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=neo4js,singular=neo4j,shortName=neo,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Neo4j struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Neo4jSpec   `json:"spec,omitempty"`
	Status Neo4jStatus `json:"status,omitempty"`
}

// Neo4jSpec defines the desired state of Neo4j.
type Neo4jSpec struct {
	// Version of Neo4j to deploy. Must correspond to a supported Neo4jVersion in the catalog.
	Version string `json:"version"`

	// Number of Neo4j instances (pods) to run. If omitted, the operator uses its default.
	Replicas *int32 `json:"replicas,omitempty"`

	// StorageType selects the data storage mode: "Durable" (default, uses PVCs) or "Ephemeral" (emptyDir).
	StorageType StorageType `json:"storageType,omitempty"`

	// PVC template used when StorageType is "Durable". Ignored when StorageType is "Ephemeral".
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// DisableSecurity disables authentication when set to true. Defaults to false. Not recommended for production.
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// AuthSecret references a Secret containing database credentials (for example, username/password).
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// ConfigSecret references a Secret that provides a custom configuration file .
	// When set, this configuration takes precedence over the operator defaults.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate customizes the pods running Neo4j (resources, environment variables, probes, affinity, etc.).
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates customizes the Services that expose Neo4j endpoints (for example, primary, replicas, admin).
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// DeletionPolicy controls behavior when this resource is deleted (for example, Delete, WipeOut, DoNotTerminate).
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// Indicates that the Neo4j Protocols that are required to be disabled on bootstrap.
	// +optional
	DisabledProtocols []Neo4jProtocol `json:"disabledProtocols,omitempty"`

	// HealthChecker configures health checks performed by the operator.
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`
}

// Neo4jStatus defines the observed state of Neo4j.
type Neo4jStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
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

// +kubebuilder:validation:Enum=tcp-bolt;neo4j;http;https;tcp-backup;graphite;prometheus;jmx;tcp-boltrouting;tcp-raft;tcp-tx
type Neo4jProtocol string

const (
	Neo4jProtocolTCPBolt        Neo4jProtocol = "tcp-bolt"
	Neo4jProtocolNeo4j          Neo4jProtocol = "neo4j" // alias for bolt
	Neo4jProtocolHTTP           Neo4jProtocol = "http"
	Neo4jProtocolHTTPS          Neo4jProtocol = "https"
	Neo4jProtocolTCPBackup      Neo4jProtocol = "tcp-backup"
	Neo4jProtocolGraphite       Neo4jProtocol = "graphite"
	Neo4jProtocolPrometheus     Neo4jProtocol = "prometheus"
	Neo4jProtocolJMX            Neo4jProtocol = "jmx"
	Neo4jProtocolTCPBoltRouting Neo4jProtocol = "tcp-boltrouting"
	Neo4jProtocolTCPRaft        Neo4jProtocol = "tcp-raft"
	Neo4jProtocolTCPTx          Neo4jProtocol = "tcp-tx"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Neo4jList contains a list of Neo4j.
type Neo4jList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Neo4j `json:"items"`
}
