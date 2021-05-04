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
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodePostgres     = "pg"
	ResourceKindPostgres     = "Postgres"
	ResourceSingularPostgres = "postgres"
	ResourcePluralPostgres   = "postgreses"
)

// Postgres defines a Postgres database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=postgreses,singular=postgres,shortName=pg,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Postgres struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              PostgresSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            PostgresStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type PostgresSpec struct {
	// Version of Postgres to be deployed.
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`

	// Number of instances to deploy for a Postgres database.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// Standby mode
	StandbyMode *PostgresStandbyMode `json:"standbyMode,omitempty" protobuf:"bytes,3,opt,name=standbyMode,casttype=PostgresStandbyMode"`

	// Streaming mode
	StreamingMode *PostgresStreamingMode `json:"streamingMode,omitempty" protobuf:"bytes,4,opt,name=streamingMode,casttype=PostgresStreamingMode"`

	// Leader election configuration
	// +optional
	LeaderElection *PostgreLeaderElectionConfig `json:"leaderElection,omitempty" protobuf:"bytes,5,opt,name=leaderElection"`

	// Database authentication secret
	AuthSecret *core.LocalObjectReference `json:"authSecret,omitempty" protobuf:"bytes,6,opt,name=authSecret"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty" protobuf:"bytes,7,opt,name=storageType,casttype=StorageType"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty" protobuf:"bytes,8,opt,name=storage"`

	// ClientAuthMode for sidecar or sharding. (default will be md5. [md5;scram;cert])
	ClientAuthMode PostgresClientAuthMode `json:"clientAuthMode,omitempty" protobuf:"bytes,9,opt,name=clientAuthMode,casttype=PostgresClientAuthMode"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	SSLMode PostgresSSLMode `json:"sslMode,omitempty" protobuf:"bytes,10,opt,name=sslMode,casttype=PostgresSSLMode"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty" protobuf:"bytes,11,opt,name=init"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,12,opt,name=monitor"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e postgresql.conf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty" protobuf:"bytes,13,opt,name=configSecret"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,14,opt,name=podTemplate"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty" protobuf:"bytes,15,rep,name=serviceTemplates"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty" protobuf:"bytes,16,opt,name=tls"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty" protobuf:"varint,17,opt,name=halted"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty" protobuf:"bytes,18,opt,name=terminationPolicy,casttype=TerminationPolicy"`
}

// PostgreLeaderElectionConfig contains essential attributes of leader election.
type PostgreLeaderElectionConfig struct {
	// LeaseDuration is the duration in second that non-leader candidates will
	// wait to force acquire leadership. This is measured against time of
	// last observed ack. Default 15
	// Deprecated
	LeaseDurationSeconds int32 `json:"leaseDurationSeconds,omitempty" protobuf:"varint,1,opt,name=leaseDurationSeconds"`
	// RenewDeadline is the duration in second that the acting master will retry
	// refreshing leadership before giving up. Normally, LeaseDuration * 2 / 3.
	// Default 10
	// Deprecated
	RenewDeadlineSeconds int32 `json:"renewDeadlineSeconds,omitempty" protobuf:"varint,2,opt,name=renewDeadlineSeconds"`
	// RetryPeriod is the duration in second the LeaderElector clients should wait
	// between tries of actions. Normally, LeaseDuration / 3.
	// Default 2
	// Deprecated
	RetryPeriodSeconds int32 `json:"retryPeriodSeconds,omitempty" protobuf:"varint,3,opt,name=retryPeriodSeconds"`

	// MaximumLagBeforeFailover is used as maximum lag tolerance for the cluster.
	// when ever a replica is lagging more than MaximumLagBeforeFailover
	// this node need to sync manually with the primary node. default value is 32MB
	// +default=33554432
	// +kubebuilder:default:=33554432
	// +optional
	MaximumLagBeforeFailover uint64 `json:"maximumLagBeforeFailover,omitempty" protobuf:"varint,4,opt,name=maximumLagBeforeFailover"`

	// Period between Node.Tick invocations
	// +kubebuilder:default:="100ms"
	// +optional
	Period metav1.Duration `json:"period,omitempty" protobuf:"bytes,5,opt,name=period"`

	// ElectionTick is the number of Node.Tick invocations that must pass between
	//	elections. That is, if a follower does not receive any message from the
	//  leader of current term before ElectionTick has elapsed, it will become
	//	candidate and start an election. ElectionTick must be greater than
	//  HeartbeatTick. We suggest ElectionTick = 10 * HeartbeatTick to avoid
	//  unnecessary leader switching. default value is 10.
	// +default=10
	// +kubebuilder:default:=10
	// +optional
	ElectionTick int32 `json:"electionTick,omitempty" protobuf:"varint,6,opt,name=electionTick"`

	// HeartbeatTick is the number of Node.Tick invocations that must pass between
	// heartbeats. That is, a leader sends heartbeat messages to maintain its
	// leadership every HeartbeatTick ticks. default value is 1.
	// +default=1
	// +kubebuilder:default:=1
	// +optional
	HeartbeatTick int32 `json:"heartbeatTick,omitempty" protobuf:"varint,7,opt,name=heartbeatTick"`
}

// +kubebuilder:validation:Enum=server;archiver;metrics-exporter
type PostgresCertificateAlias string

const (
	PostgresServerCert          PostgresCertificateAlias = "server"
	PostgresClientCert          PostgresCertificateAlias = "client"
	PostgresArchiverCert        PostgresCertificateAlias = "archiver"
	PostgresMetricsExporterCert PostgresCertificateAlias = "metrics-exporter"
)

type PostgresStatus struct {
	// Specifies the current phase of the database
	// +optional
	Phase DatabasePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=DatabasePhase"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,2,opt,name=observedGeneration"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty" protobuf:"bytes,3,rep,name=conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PostgresList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Postgres CRD objects
	Items []Postgres `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}

type RecoveryTarget struct {
	// TargetTime specifies the time stamp up to which recovery will proceed.
	TargetTime string `json:"targetTime,omitempty" protobuf:"bytes,1,opt,name=targetTime"`
	// TargetTimeline specifies recovering into a particular timeline.
	// The default is to recover along the same timeline that was current when the base backup was taken.
	TargetTimeline string `json:"targetTimeline,omitempty" protobuf:"bytes,2,opt,name=targetTimeline"`
	// TargetXID specifies the transaction ID up to which recovery will proceed.
	TargetXID string `json:"targetXID,omitempty" protobuf:"bytes,3,opt,name=targetXID"`
	// TargetInclusive specifies whether to include ongoing transaction in given target point.
	TargetInclusive *bool `json:"targetInclusive,omitempty" protobuf:"varint,4,opt,name=targetInclusive"`
}

// +kubebuilder:validation:Enum=Hot;Warm
type PostgresStandbyMode string

const (
	HotPostgresStandbyMode  PostgresStandbyMode = "Hot"
	WarmPostgresStandbyMode PostgresStandbyMode = "Warm"
)

// +kubebuilder:validation:Enum=Synchronous;Asynchronous
type PostgresStreamingMode string

const (
	SynchronousPostgresStreamingMode  PostgresStreamingMode = "Synchronous"
	AsynchronousPostgresStreamingMode PostgresStreamingMode = "Asynchronous"
)

// ref: https://www.postgresql.org/docs/13/libpq-ssl.html
// +kubebuilder:validation:Enum=disable;allow;prefer;require;verify-ca;verify-full
type PostgresSSLMode string

const (
	// PostgresSSLModeDisable represents `disable` sslMode. It ensures that the server does not use TLS/SSL.
	PostgresSSLModeDisable PostgresSSLMode = "disable"

	// PostgresSSLModeAllow represents `allow` sslMode. 	I don't care about security,
	// but I will pay the overhead of encryption if the server insists on it.
	PostgresSSLModeAllow PostgresSSLMode = "allow"

	// PostgresSSLModePrefer represents `preferSSL` sslMode.
	// I don't care about encryption, but I wish to pay the overhead of encryption if the server supports it.
	PostgresSSLModePrefer PostgresSSLMode = "prefer"

	// PostgresSSLModeRequire represents `requiteSSL` sslmode. I want my data to be encrypted, and I accept the overhead.
	// I trust that the network will make sure I always connect to the server I want.
	PostgresSSLModeRequire PostgresSSLMode = "require"

	// PostgresSSLModeVerifyCA represents `verify-ca` sslmode. I want my data encrypted, and I accept the overhead.
	// I want to be sure that I connect to a server that I trust.
	PostgresSSLModeVerifyCA PostgresSSLMode = "verify-ca"

	// PostgresSSLModeVerifyFull represents `verify-full` sslmode. I want my data encrypted, and I accept the overhead.
	// I want to be sure that I connect to a server I trust, and that it's the one I specify.
	PostgresSSLModeVerifyFull PostgresSSLMode = "verify-full"
)

// PostgresClientAuthMode represents the ClientAuthMode of PostgreSQL clusters ( replicaset )
// ref: https://www.postgresql.org/docs/12/auth-methods.html
// +kubebuilder:validation:Enum=md5;scram;cert
type PostgresClientAuthMode string

const (
	// ClientAuthModeMD5 uses a custom less secure challenge-response mechanism.
	// It prevents password sniffing and avoids storing passwords on the server in plain text but provides no protection
	// if an attacker manages to steal the password hash from the server.
	// Also, the MD5 hash algorithm is nowadays no longer considered secure against determined attacks
	ClientAuthModeMD5 PostgresClientAuthMode = "md5"

	// ClientAuthModeScram performs SCRAM-SHA-256 authentication, as described in RFC 7677.
	// It is a challenge-response scheme that prevents password sniffing on untrusted connections
	// and supports storing passwords on the server in a cryptographically hashed form that is thought to be secure.
	// This is the most secure of the currently provided methods, but it is not supported by older client libraries.
	ClientAuthModeScram PostgresClientAuthMode = "scram"

	// ClientAuthModeCert represents `cert clientcert=1` auth mode where client need to provide cert and private key for authentication.
	// When server is config with this auth method. Client can't connect with postgreSQL server with password. They need
	// to Send the client cert and client key certificate for authentication.
	ClientAuthModeCert PostgresClientAuthMode = "cert"
)
