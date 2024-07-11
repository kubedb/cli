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

package v1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
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
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Postgres struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PostgresSpec   `json:"spec,omitempty"`
	Status            PostgresStatus `json:"status,omitempty"`
}
type PostgreSQLMode string

const (
	PostgreSQLModeStandAlone    PostgreSQLMode = "Standalone"
	PostgreSQLModeRemoteReplica PostgreSQLMode = "RemoteReplica"
	PostgreSQLModeCluster       PostgreSQLMode = "Cluster"
)

type PostgresSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of Postgres to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Postgres database.
	Replicas *int32 `json:"replicas,omitempty"`

	// Standby mode
	StandbyMode *PostgresStandbyMode `json:"standbyMode,omitempty"`

	// Streaming mode
	StreamingMode *PostgresStreamingMode `json:"streamingMode,omitempty"`

	// + optional
	Mode *PostgreSQLMode `json:"mode,omitempty"`
	// RemoteReplica implies that the instance will be a MySQL Read Only Replica,
	// and it will take reference of  appbinding of the source
	// +optional
	RemoteReplica *RemoteReplicaSpec `json:"remoteReplica,omitempty"`

	// Leader election configuration
	// +optional
	LeaderElection *PostgreLeaderElectionConfig `json:"leaderElection,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// ClientAuthMode for sidecar or sharding. (default will be md5. [md5;scram;cert])
	ClientAuthMode PostgresClientAuthMode `json:"clientAuthMode,omitempty"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	SSLMode PostgresSSLMode `json:"sslMode,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e postgresql.conf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// EnforceFsGroup Is Used when the storageClass's CSI Driver doesn't support FsGroup properties properly.
	// If It's true then The Init Container will run as RootUser and
	// the init-container will set user's permission for the mounted pvc volume with which coordinator and postgres containers are going to run.
	// In postgres it is /var/pv
	// +optional
	EnforceFsGroup bool `json:"enforceFsGroup,omitempty"`

	// AllowedSchemas defines the types of database schemas that MAY refer to
	// a database instance and the trusted namespaces where those schema resources MAY be
	// present.
	//
	// +kubebuilder:default={namespaces:{from: Same}}
	// +optional
	AllowedSchemas *AllowedConsumers `json:"allowedSchemas,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Archiver controls database backup using Archiver CR
	// +optional
	Archiver *Archiver `json:"archiver,omitempty"`

	// Arbiter controls spec for arbiter pods
	// +optional
	Arbiter *ArbiterSpec `json:"arbiter,omitempty"`

	// +optional
	Replication *PostgresReplication `json:"replication,omitempty"`
}

type WALLimitPolicy string

const (
	WALKeepSize     WALLimitPolicy = "WALKeepSize"
	ReplicationSlot WALLimitPolicy = "ReplicationSlot"
	WALKeepSegment  WALLimitPolicy = "WALKeepSegment"
)

type PostgresReplication struct {
	WALLimitPolicy WALLimitPolicy `json:"walLimitPolicy"`

	// +optional
	WalKeepSizeInMegaBytes *int32 `json:"walKeepSize,omitempty"`
	// +optional
	WalKeepSegment *int32 `json:"walKeepSegment,omitempty"`
	// +optional
	MaxSlotWALKeepSizeInMegaBytes *int32 `json:"maxSlotWALKeepSize,omitempty"`
}

type ArbiterSpec struct {
	// Compute Resources required by the sidecar container.
	// +optional
	Resources core.ResourceRequirements `json:"resources,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	// +mapType=atomic
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty"`
}

// PostgreLeaderElectionConfig contains essential attributes of leader election.
type PostgreLeaderElectionConfig struct {
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
	Phase DatabasePhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
	// +optional
	AuthSecret *Age `json:"authSecret,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PostgresList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Postgres CRD objects
	Items []Postgres `json:"items,omitempty"`
}

type RecoveryTarget struct {
	// TargetTime specifies the time stamp up to which recovery will proceed.
	TargetTime string `json:"targetTime,omitempty"`
	// TargetTimeline specifies recovering into a particular timeline.
	// The default is to recover along the same timeline that was current when the base backup was taken.
	TargetTimeline string `json:"targetTimeline,omitempty"`
	// TargetXID specifies the transaction ID up to which recovery will proceed.
	TargetXID string `json:"targetXID,omitempty"`
	// TargetInclusive specifies whether to include ongoing transaction in given target point.
	TargetInclusive *bool `json:"targetInclusive,omitempty"`
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
