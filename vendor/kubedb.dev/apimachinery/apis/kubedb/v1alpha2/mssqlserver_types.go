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
	ResourceCodeMSSQLServer     = "ms"
	ResourceKindMSSQLServer     = "MSSQLServer"
	ResourceSingularMSSQLServer = "mssqlserver"
	ResourcePluralMSSQLServer   = "mssqlservers"
)

// +kubebuilder:validation:Enum=AvailabilityGroup;RemoteReplica
type MSSQLServerMode string

const (
	MSSQLServerModeAvailabilityGroup MSSQLServerMode = "AvailabilityGroup"
	MSSQLServerModeRemoteReplica     MSSQLServerMode = "RemoteReplica"
)

// +kubebuilder:validation:Enum=server;client;endpoint
type MSSQLServerCertificateAlias string

const (
	MSSQLServerServerCert   MSSQLServerCertificateAlias = "server"
	MSSQLServerClientCert   MSSQLServerCertificateAlias = "client"
	MSSQLServerEndpointCert MSSQLServerCertificateAlias = "endpoint"
)

// +kubebuilder:validation:Enum=Passive;ReadOnly;All
type SecondaryAccessMode string

const (
	// Passive  = secondary is passive, no connections allowed
	SecondaryAccessModePassive SecondaryAccessMode = "Passive"
	// ReadOnly = secondary allows read-intent only
	SecondaryAccessModeReadOnly SecondaryAccessMode = "ReadOnly"
	// All = secondary allows any connections
	SecondaryAccessModeAll SecondaryAccessMode = "All"
)

// MSSQLServer defines a MSSQLServer database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mssqlservers,singular=mssqlserver,shortName=ms,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MSSQLServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MSSQLServerSpec   `json:"spec,omitempty"`
	Status MSSQLServerStatus `json:"status,omitempty"`
}

// MSSQLServerSpec defines the desired state of MSSQLServer
type MSSQLServerSpec struct {
	// Version of MSSQLServer to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a MSSQLServer database. In case of MSSQLServer Availability Group.
	Replicas *int32 `json:"replicas,omitempty"`

	// MSSQLServer cluster topology
	// +optional
	Topology *MSSQLServerTopology `json:"topology,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// ConfigSecret is an optional field to provide a custom configuration file for the database (i.e., mssql.conf).
	// If specified, this file will be used as a configuration file, otherwise a default configuration file will be used.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// Init is used to initialize a database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// TLS contains tls configurations for client and server.
	TLS *MSSQLServerTLSConfig `json:"tls,omitempty"`

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

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// Archiver controls database backup using Archiver CR
	// +optional
	Archiver *Archiver `json:"archiver,omitempty"`

	// Arbiter controls spec for arbiter pods
	// +optional
	Arbiter *ArbiterSpec `json:"arbiter,omitempty"`
}

type MSSQLServerTLSConfig struct {
	kmapi.TLSConfig `json:",inline"`

	// +optional
	ClientTLS *bool `json:"clientTLS"`
}

type MSSQLServerTopology struct {
	// If set to -
	// "AvailabilityGroup", MSSQLAvailabilityGroupSpec is required and MSSQLServer servers will start an Availability Group
	Mode *MSSQLServerMode `json:"mode,omitempty"`

	// AvailabilityGroup info for MSSQLServer
	// +optional
	AvailabilityGroup *MSSQLServerAvailabilityGroupSpec `json:"availabilityGroup,omitempty"`
}

// MSSQLServerAvailabilityGroupSpec defines the availability group spec for MSSQLServer
type MSSQLServerAvailabilityGroupSpec struct {
	// AvailabilityDatabases is an array of databases to be included in the availability group
	// +optional
	Databases []string `json:"databases,omitempty"`

	// Leader election configuration
	// +optional
	LeaderElection *MSSQLServerLeaderElectionConfig `json:"leaderElection,omitempty"`

	// SecondaryAccessMode controls which connections are allowed to secondary replicas.
	// https://learn.microsoft.com/en-us/sql/t-sql/statements/create-availability-group-transact-sql?view=sql-server-ver16#secondary_role---
	// +optional
	// +kubebuilder:default=Passive
	SecondaryAccessMode SecondaryAccessMode `json:"secondaryAccessMode,omitempty"`
}

// MSSQLServerStatus defines the observed state of MSSQLServer
type MSSQLServerStatus struct {
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

// MSSQLServerLeaderElectionConfig contains essential attributes of leader election.
type MSSQLServerLeaderElectionConfig struct {
	// Period between Node.Tick invocations
	// +kubebuilder:default="100ms"
	// +optional
	Period metav1.Duration `json:"period,omitempty"`

	// ElectionTick is the number of Node.Tick invocations that must pass between
	//	elections. That is, if a follower does not receive any message from the
	//  leader of current term before ElectionTick has elapsed, it will become
	//	candidate and start an election. ElectionTick must be greater than
	//  HeartbeatTick. We suggest ElectionTick = 10 * HeartbeatTick to avoid
	//  unnecessary leader switching. default value is 10.
	// +default=10
	// +kubebuilder:default=10
	// +optional
	ElectionTick int32 `json:"electionTick,omitempty"`

	// HeartbeatTick is the number of Node.Tick invocations that must pass between
	// heartbeats. That is, a leader sends heartbeat messages to maintain its
	// leadership every HeartbeatTick ticks. default value is 1.
	// +default=1
	// +kubebuilder:default=1
	// +optional
	HeartbeatTick int32 `json:"heartbeatTick,omitempty"`

	// TransferLeadershipInterval retry interval for transfer leadership
	// to the healthiest node
	// +kubebuilder:default="1s"
	// +optional
	TransferLeadershipInterval *metav1.Duration `json:"transferLeadershipInterval,omitempty"`

	// TransferLeadershipTimeout retry timeout for transfer leadership
	// to the healthiest node
	// +kubebuilder:default="60s"
	// +optional
	TransferLeadershipTimeout *metav1.Duration `json:"transferLeadershipTimeout,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MSSQLServerList contains a list of MSSQLServer
type MSSQLServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MSSQLServer `json:"items"`
}
