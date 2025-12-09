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
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeOracle     = "ora"
	ResourceKindOracle     = "Oracle"
	ResourceSingularOracle = "oracle"
	ResourcePluralOracle   = "oracles"
)

// +kubebuilder:validation:Enum=Standalone;DataGuard
// OracleMode defines supported deployment modes
type OracleMode string

const (
	OracleModeStandalone OracleMode = "Standalone"
	OracleModeDataGuard  OracleMode = "DataGuard"
	// OracleModeRAC        OracleMode = "RAC"
	// OracleModeSharding   OracleMode = "Sharding"
)

// +kubebuilder:validation:Enum=server;metrics-exporter;client
type OracleCertificateAlias string

const (
	OracleServerCert          OracleCertificateAlias = "server"
	OracleClientCert          OracleCertificateAlias = "client"
	OracleMetricsExporterCert OracleCertificateAlias = "metrics-exporter"
)

// ProtectionMode defines the data protection mode for a resource.
// +kubebuilder:validation:Enum=MaximumAvailability;MaximumPerformance;MaximumProtection
type ProtectionMode string

const (
	// ProtectionModeMaximumAvailability provides high availability with possible trade-offs in performance.
	ProtectionModeMaximumAvailability ProtectionMode = "MaximumAvailability"

	// ProtectionModeMaximumPerformance optimizes for speed with reduced redundancy.
	ProtectionModeMaximumPerformance ProtectionMode = "MaximumPerformance"

	// ProtectionModeMaximumProtection ensures maximum data durability at the cost of performance.
	ProtectionModeMaximumProtection ProtectionMode = "MaximumProtection"
)

// SyncMode defines the synchronization mode for data replication.
// +kubebuilder:validation:Enum=SYNC;ASYNC
type SyncMode string

const (
	// SyncModeSync indicates synchronous replication (strong consistency).
	SyncModeSync SyncMode = "SYNC"

	// SyncModeAsync indicates asynchronous replication (eventual consistency).
	SyncModeAsync SyncMode = "ASYNC"
)

// StandbyType defines the type of standby configuration.
// +kubebuilder:validation:Enum=PHYSICAL;LOGICAL
type StandbyType string

const (
	// StandbyTypePhysical indicates a physical standby (block-level replication).
	StandbyTypePhysical StandbyType = "PHYSICAL"

	// StandbyTypeLogical indicates a logical standby (SQL-level replication).
	StandbyTypeLogical StandbyType = "LOGICAL"
)

// OracleListenerProtocol defines the protocol used for Oracle database listeners.
// +kubebuilder:validation:Enum=TCP;TCPS
type OracleListenerProtocol string

const (
	// OracleListenerProtocolTCP indicates standard TCP protocol (unencrypted)
	OracleListenerProtocolTCP OracleListenerProtocol = "TCP"

	// OracleListenerProtocolTCPS indicates TCP with SSL/TLS encryption
	OracleListenerProtocolTCPS OracleListenerProtocol = "TCPS"
)

// OracleSpec defines the desired state of Oracle.
type OracleSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of Oracle to be deployed.
	// +optional
	Version string `json:"version"`

	// Deployment mode of Oracle
	Mode OracleMode `json:"mode,omitempty"`

	// future versions standard;express;free
	// +kubebuilder:validation:Enum=enterprise
	Edition string `json:"edition,omitempty"`

	// Number of instances (for RAC or DataGuard primary+standby)
	Replicas *int32 `json:"replicas,omitempty"`

	// Core storage type: durable or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e config.properties).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

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
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// DataGuard configuration, required if Mode is "DataGuard".
	// +optional
	DataGuard *DataGuardSpec `json:"dataGuard,omitempty"`

	// Listener is for Oracle Net Listener
	// +optional
	Listener *ListenerSpec `json:"listener,omitempty"`
	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// TLS configuration for secure client connections
	// +optional
	TCPSConfig *OracleTCPSConfig `json:"tcpsConfig,omitempty"`
}

type OracleTCPSConfig struct {
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`
	// Listener TCPS configuration
	// +optional
	TCPSListener *ListenerSpec `json:"tcpsListener,omitempty"`
}

// ListenerSpec defines a TNS listener (TCP or TCPS)
type ListenerSpec struct {
	// Listener name
	// +optional
	Name string `json:"name,omitempty"`
	// Port number
	// +kubebuilder:validation:Minimum=1025
	Port *int32 `json:"port,omitempty"`
	// Database Service
	Service *string `json:"service,omitempty"`
	// Protocol (TCP, TCPS)
	Protocol OracleListenerProtocol `json:"protocol,omitempty"`
}

type DataGuardSpec struct {
	// Protection mode: MaximumAvailability, MaximumPerformance, MaximumProtection
	ProtectionMode ProtectionMode `json:"protectionMode,omitempty"`
	// SyncMode specifies the synchronization mode (e.g., "SYNC", "ASYNC").
	SyncMode SyncMode `json:"syncMode,omitempty"`
	// StandbyType specifies the type of standby (e.g., "PHYSICAL", "LOGICAL").
	// +optional
	StandbyType StandbyType `json:"standbyType,omitempty"`
	// Oracle Failover Assistant
	FastStartFailover *FastStartFailover `json:"fastStartFailover,omitempty"`
	// ApplyLag allowed
	ApplyLagThreshold *int32 `json:"applyLagThreshold,omitempty"`
	// Dataguard observer spec
	Observer *ObserverSpec `json:"observer,omitempty"`
	// Transport Lag
	TransportLagThreshold *int32 `json:"transportLagThreshold,omitempty"`
}

type FastStartFailover struct {
	// FastStartFailoverThreshold configuration property defines the number of seconds the master observer attempts
	// to reconnect to the primary database before initiating a fast-start failover.
	FastStartFailoverThreshold *int32 `json:"fastStartFailoverThreshold,omitempty"`
}

type ObserverSpec struct {
	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`
	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
}

// OracleStatus defines the observed state of Oracle.
type OracleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
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

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=oracles,singular=oracle,shortName=ora,categories={datastore,oracle,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Mode",type="string",JSONPath=".spec.mode"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Oracle is the Schema for the oracles API.
type Oracle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OracleSpec   `json:"spec,omitempty"`
	Status OracleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OracleList contains a list of Oracle.
type OracleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Oracle `json:"items"`
}
