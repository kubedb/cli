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
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodePgBouncer     = "pb"
	ResourceKindPgBouncer     = "PgBouncer"
	ResourceSingularPgBouncer = "pgbouncer"
	ResourcePluralPgBouncer   = "pgbouncers"
)

// PgBouncer defines a PgBouncer Server.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pgbouncers,singular=pgbouncer,shortName=pb,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PgBouncer struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PgBouncerSpec   `json:"spec,omitempty"`
	Status            PgBouncerStatus `json:"status,omitempty"`
}

type PgBouncerSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of PgBouncer to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a PgBouncer instance.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// PodTemplate is an optional configuration for pods.
	// +optional
	PodTemplate ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// Database to proxy by connection pooling.
	Database Database `json:"database"`

	// ConnectionPoolConfig defines Connection pool configuration.
	// +optional
	ConnectionPool *ConnectionPoolConfig `json:"connectionPool,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e mongod.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// Monitor is used monitor database instance.
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	SSLMode PgBouncerSSLMode `json:"sslMode,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Indicates that the database is halted and all offshoot Kubernetes resources are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`
}

// +kubebuilder:validation:Enum=server;archiver;metrics-exporter
type PgBouncerCertificateAlias string

const (
	PgBouncerServerCert          PgBouncerCertificateAlias = "server"
	PgBouncerClientCert          PgBouncerCertificateAlias = "client"
	PgBouncerMetricsExporterCert PgBouncerCertificateAlias = "metrics-exporter"
)

type Database struct {
	// SyncUsers is a boolean type and when enabled, operator fetches users of backend server from externally managed
	// secrets to the PgBouncer server. Secrets updation or deletion are also synced in pgBouncer when it is enabled.
	// +optional
	SyncUsers bool `json:"syncUsers,omitempty"`

	// DatabaseRef specifies the database appbinding reference in any namespace.
	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// DatabaseName is the name of the target database inside a Postgres instance.
	DatabaseName string `json:"databaseName"`
}

type ConnectionPoolConfig struct {
	// Port is the port number on which PgBouncer listens to clients. Default: 5432.
	// +optional
	Port *int32 `json:"port,omitempty"`
	// PoolMode is the pooling mechanism type. Default: session.
	// +optional
	PoolMode string `json:"poolMode,omitempty"`
	// MaxClientConnections is the maximum number of allowed client connections. Default: 100.
	// +optional
	MaxClientConnections *int64 `json:"maxClientConnections,omitempty"`
	// DefaultPoolSize specifies how many server connections to allow per user/database pair. Default: 20.
	// +optional
	DefaultPoolSize *int64 `json:"defaultPoolSize,omitempty"`
	// MinPoolSize is used to add more server connections to pool if below this number. Default: 0 (disabled).
	// +optional
	MinPoolSize *int64 `json:"minPoolSize,omitempty"`
	// ReservePoolSize specifies how many additional connections to allow to a pool. 0 disables. Default: 0 (disabled).
	// +optional
	ReservePoolSize *int64 `json:"reservePoolSize,omitempty"`
	// ReservePoolTimeoutSeconds is the number of seconds in which if a client has not been serviced,
	// pgbouncer enables use of additional connections from reserve pool. 0 disables. Default: 5.0.
	// +optional
	ReservePoolTimeoutSeconds *int64 `json:"reservePoolTimeoutSeconds,omitempty"`
	// MaxDBConnections is the maximum number of connections allowed per-database. Default: 0 (unlimited).
	// +optional
	MaxDBConnections *int64 `json:"maxDBConnections,omitempty"`
	// MaxUserConnections is the maximum number of users allowed per-database. Default: 0 (unlimited).
	// +optional
	MaxUserConnections *int64 `json:"maxUserConnections,omitempty"`
	// StatsPeriodSeconds sets how often the averages shown in various SHOW commands are updated
	// and how often aggregated statistics are written to the log. Default: 60
	// +optional
	StatsPeriodSeconds *int64 `json:"statsPeriodSeconds,omitempty"`
	// AuthType specifies how to authenticate users. Default: md5 (md5+plain text).
	// +optional
	AuthType PgBouncerClientAuthMode `json:"authType,omitempty"`
	// IgnoreStartupParameters specifies comma-separated startup parameters that
	// pgbouncer knows are handled by admin and it can ignore them. Default: empty
	// +optional
	IgnoreStartupParameters string `json:"ignoreStartupParameters,omitempty"`
	// AdminUsers specifies an array of users who can act as PgBouncer administrators.
	// +optional
	// AdminUsers []string `json:"adminUsers,omitempty"`
	// AuthUser looks up any user not specified in auth_file from pg_shadow. Default: not set.
	// +optional
	// AuthUser string `json:"authUser,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PgBouncer CRD objects.
	Items []PgBouncer `json:"items,omitempty"`
}

type PgBouncerStatus struct {
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

// +kubebuilder:validation:Enum=disable;allow;prefer;require;verify-ca;verify-full
type PgBouncerSSLMode string

const (
	// PgBouncerSSLModeDisable represents `disable` sslMode. It ensures that the server does not use TLS/SSL.
	PgBouncerSSLModeDisable PgBouncerSSLMode = "disable"

	// PgBouncerSSLModeAllow represents `allow` sslMode. 	I don't care about security,
	// but I will pay the overhead of encryption if the server insists on it.
	PgBouncerSSLModeAllow PgBouncerSSLMode = "allow"

	// PgBouncerSSLModePrefer represents `preferSSL` sslMode.
	// I don't care about encryption, but I wish to pay the overhead of encryption if the server supports it.
	PgBouncerSSLModePrefer PgBouncerSSLMode = "prefer"

	// PgBouncerSSLModeRequire represents `requiteSSL` sslmode. I want my data to be encrypted, and I accept the overhead.
	// I trust that the network will make sure I always connect to the server I want.
	PgBouncerSSLModeRequire PgBouncerSSLMode = "require"

	// PgBouncerSSLModeVerifyCA represents `verify-ca` sslmode. I want my data encrypted, and I accept the overhead.
	// I want to be sure that I connect to a server that I trust.
	PgBouncerSSLModeVerifyCA PgBouncerSSLMode = "verify-ca"

	// PgBouncerSSLModeVerifyFull represents `verify-full` sslmode. I want my data encrypted, and I accept the overhead.
	// I want to be sure that I connect to a server I trust, and that it's the one I specify.
	PgBouncerSSLModeVerifyFull PgBouncerSSLMode = "verify-full"
)

// PgBouncerClientAuthMode represents the ClientAuthMode of PgBouncer clusters ( replicaset )
// We are allowing md5, scram-sha-256, cert as ClientAuthMode
// +kubebuilder:validation:Enum=md5;scram-sha-256;cert;
type PgBouncerClientAuthMode string

const (
	// ClientAuthModeMD5 uses a custom less secure challenge-response mechanism.
	// It prevents password sniffing and avoids storing passwords on the server in plain text but provides no protection
	// if an attacker manages to steal the password hash from the server.
	// Also, the MD5 hash algorithm is nowadays no longer considered secure against determined attacks
	PgBouncerClientAuthModeMD5 PgBouncerClientAuthMode = "md5"

	// ClientAuthModeScram performs SCRAM-SHA-256 authentication, as described in RFC 7677.
	// It is a challenge-response scheme that prevents password sniffing on untrusted connections
	// and supports storing passwords on the server in a cryptographically hashed form that is thought to be secure.
	// This is the most secure of the currently provided methods, but it is not supported by older client libraries.
	PgBouncerClientAuthModeScram PgBouncerClientAuthMode = "scram-sha-256"

	// ClientAuthModeCert represents `cert clientcert=1` auth mode where client need to provide cert and private key for authentication.
	// When server is config with this auth method. Client can't connect with pgbouncer server with password. They need
	// to Send the client cert and client key certificate for authentication.
	PgBouncerClientAuthModeCert PgBouncerClientAuthMode = "cert"
)
