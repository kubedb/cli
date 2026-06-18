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
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeDocumentDB     = "dc"
	ResourceKindDocumentDB     = "DocumentDB"
	ResourceSingularDocumentDB = "documentdb"
	ResourcePluralDocumentDB   = "documentdbs"
)

// +kubebuilder:validation:Enum=Hot;Warm
type DocDBStandbyMode string

const (
	HotDocDBStandbyMode  DocDBStandbyMode = "Hot"
	WarmDocDBStandbyMode DocDBStandbyMode = "Warm"
)

// +kubebuilder:validation:Enum=Synchronous;Asynchronous
type DocDBStreamingMode string

const (
	SynchronousDocDBStreamingMode  DocDBStreamingMode = "Synchronous"
	AsynchronousDocDBStreamingMode DocDBStreamingMode = "Asynchronous"
)

// DocDBClientAuthMode represents the ClientAuthMode of DocumentDB clusters ( replicaset )
// +kubebuilder:validation:Enum=scram;cert
type DocDBClientAuthMode string

const (

	// ClientAuthModeScram performs SCRAM-SHA-256 authentication, as described in RFC 7677.
	// It is a challenge-response scheme that prevents password sniffing on untrusted connections
	// and supports storing passwords on the server in a cryptographically hashed form that is thought to be secure.
	// This is the most secure of the currently provided methods, but it is not supported by older client libraries.
	DocDBClientAuthModeScram DocDBClientAuthMode = "scram"

	// ClientAuthModeCert represents `cert clientcert=1` auth mode where client need to provide cert and private key for authentication.
	// When server is config with this auth method. Client can't connect with postgreSQL server with password. They need
	// to Send the client cert and client key certificate for authentication.
	DocDBClientAuthModeCert DocDBClientAuthMode = "cert"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=documentdbs,singular=documentdb,shortName=docdb,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Namespace",type="string",JSONPath=".metadata.namespace"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type DocumentDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DocumentDBSpec   `json:"spec,omitempty"`
	Status DocumentDBStatus `json:"status,omitempty"`
}

type DocumentDBSpec struct {
	// Version of DocumentDB to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a documentdb database.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Standby mode
	StandbyMode *DocDBStandbyMode `json:"standbyMode,omitempty"`

	// Streaming mode
	StreamingMode *DocDBStreamingMode `json:"streamingMode,omitempty"`

	// ClientAuthMode for sidecar or sharding. (default will be md5. [md5;scram;cert])
	ClientAuthMode DocDBClientAuthMode `json:"clientAuthMode,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Leader election configuration
	// +optional
	LeaderElection *DocumentDBLeaderElectionConfig `json:"leaderElection,omitempty"`

	// +optional
	Replication *DocumentDBReplication `json:"replication,omitempty"`

	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// AdminAuthSecret specifies the admin auth secret for "default_user"
	// +optional
	AdminAuthSecret *SecretReference `json:"adminAuthSecret,omitempty"`

	// Configuration is an optional field to provide custom configuration and performance tuning for the database.
	// +optional
	Configuration *DocumentDBConfiguration `json:"configuration,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

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

	// Init is used to initialize the database from a script or git repo.
	// +optional
	Init *InitSpec `json:"init,omitempty"`
}

type DocumentDBConfiguration struct {
	ConfigurationSpec `json:",inline,omitempty"`
	// Tuning defines performance tuning options for the database.
	// +optional
	Tuning *DocumentDBTuningConfig `json:"tuning,omitempty"`
}

// DocumentDBTuningConfig defines configuration for DocumentDB (PostgreSQL) performance tuning
type DocumentDBTuningConfig struct {
	// Profile defines a predefined tuning profile for different workload types.
	// If specified, other tuning parameters will be calculated based on this profile.
	// +optional
	Profile *DocumentDBProfile `json:"profile,omitempty"`

	// MaxConnections defines the maximum number of concurrent connections.
	// If not specified, it will be calculated based on available memory and tuning profile.
	// +optional
	MaxConnections *int32 `json:"maxConnections,omitempty"`

	// StorageType defines the type of storage for tuning purposes.
	// If not specified, it will be inferred from StorageClass or default to HDD.
	// +optional
	StorageType *DocumentDBStorageType `json:"storageType,omitempty"`

	// DisableAutoTune disables automatic tuning entirely.
	// If set to true, no tuning will be applied.
	// +optional
	DisableAutoTune bool `json:"disableAutoTune,omitempty"`
}

// DocumentDBProfile defines predefined tuning profiles
// +kubebuilder:validation:Enum=web;oltp;dw;mixed;desktop
type DocumentDBProfile string

const (
	DocumentDBTuningProfileWeb     DocumentDBProfile = "web"
	DocumentDBTuningProfileOLTP    DocumentDBProfile = "oltp"
	DocumentDBTuningProfileDW      DocumentDBProfile = "dw"
	DocumentDBTuningProfileMixed   DocumentDBProfile = "mixed"
	DocumentDBTuningProfileDesktop DocumentDBProfile = "desktop"
)

// DocumentDBStorageType defines storage types for tuning purposes
// +kubebuilder:validation:Enum=ssd;hdd;san
type DocumentDBStorageType string

const (
	DocumentDBStorageTypeSSD DocumentDBStorageType = "ssd"
	DocumentDBStorageTypeHDD DocumentDBStorageType = "hdd"
	DocumentDBStorageTypeSAN DocumentDBStorageType = "san"
)

type DocumentDBStatus struct {
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

type DocumentDBReplication struct {
	// WALimitPolicy defines which WAL retention policy to use.
	WALLimitPolicy WALLimitPolicy `json:"walLimitPolicy"`

	// +optional
	WalKeepSizeInMegaBytes *int32 `json:"walKeepSize,omitempty"`
	// +optional
	WalKeepSegment *int32 `json:"walKeepSegment,omitempty"`
	// +optional
	MaxSlotWALKeepSizeInMegaBytes *int32 `json:"maxSlotWALKeepSize,omitempty"`

	// ForceFailoverAcceptingDataLossAfter is the maximum time to wait before running a force failover process
	// This is helpful for a scenario where the old primary is not available and it has the most updated wal lsn
	// Doing force failover may or may not end up loosing data depending on any wrtie transaction
	// in the range lagged lsn between the new primary and the old primary
	// +optional
	ForceFailoverAcceptingDataLossAfter *metav1.Duration `json:"forceFailoverAcceptingDataLossAfter,omitempty"`
}

type DocumentDBLeaderElectionConfig struct {
	// LeaseDuration is the duration in second that non-leader candidates will
	// wait to force acquire leadership. This is measured against time of
	// last observed ack. Default 15
	// Deprecated
	LeaseDurationSeconds int32 `json:"leaseDurationSeconds,omitempty"`
	// RenewDeadline is the duration in second that the acting master will retry
	// refreshing leadership before giving up. Normally, LeaseDuration * 2 / 3.
	// Default 10
	// Deprecated
	RenewDeadlineSeconds int32 `json:"renewDeadlineSeconds,omitempty"`
	// RetryPeriod is the duration in second the LeaderElector clients should wait
	// between tries of actions. Normally, LeaseDuration / 3.
	// Default 2
	// Deprecated
	RetryPeriodSeconds int32 `json:"retryPeriodSeconds,omitempty"`

	// MaximumLagBeforeFailover is used as maximum lag tolerance for the cluster.
	// when ever a replica is lagging more than MaximumLagBeforeFailover
	// this node need to sync manually with the primary node. default value is 32MB
	// +default=33554432
	// +kubebuilder:default=33554432
	// +optional
	MaximumLagBeforeFailover uint64 `json:"maximumLagBeforeFailover,omitempty"`

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

type DocumentDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DocumentDB `json:"items"`
}

func (m *DocumentDB) GetObjectMeta() metav1.ObjectMeta {
	return m.ObjectMeta
}

func (m *DocumentDB) GetConditions() []kmapi.Condition {
	return m.Status.Conditions
}

func (m *DocumentDB) SetCondition(cond kmapi.Condition) {
	m.Status.Conditions = setCondition(m.Status.Conditions, cond)
}

func (m *DocumentDB) RemoveCondition(typ string) {
	m.Status.Conditions = removeCondition(m.Status.Conditions, typ)
}
