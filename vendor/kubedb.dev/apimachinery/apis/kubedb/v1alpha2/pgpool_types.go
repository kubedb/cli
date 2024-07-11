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
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodePgpool     = "pp"
	ResourceKindPgpool     = "Pgpool"
	ResourceSingularPgpool = "pgpool"
	ResourcePluralPgpool   = "pgpools"
)

// Pgpool is the Schema for the pgpools API
// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=pgpools,singular=pgpool,shortName=pp,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Pgpool struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            PgpoolSpec   `json:"spec,omitempty"`
	Status          PgpoolStatus `json:"status,omitempty"`
}

// PgpoolSpec defines the desired state of Pgpool
type PgpoolSpec struct {
	// SyncUsers is a boolean type and when enabled, operator fetches all users created in the backend server to the
	// Pgpool server . Password changes are also synced in pgpool when it is enabled.
	// +optional
	SyncUsers bool `json:"syncUsers,omitempty"`

	// Version of Pgpool to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Pgpool instance.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// PostgresRef refers to the AppBinding of the backend PostgreSQL server
	PostgresRef *kmapi.ObjectReference `json:"postgresRef"`

	// Pgpool secret containing username and password for pgpool pcp user
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// ConfigSecret is a configuration secret which will be created with default and InitConfiguration
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose Pgpool
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// InitConfiguration contains information with which the Pgpool will bootstrap
	// +optional
	InitConfiguration *PgpoolConfiguration `json:"initConfig,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose Pgpool
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Monitor is used to monitor Pgpool instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy TerminationPolicy `json:"deletionPolicy,omitempty"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	SSLMode PgpoolSSLMode `json:"sslMode,omitempty"`

	// ClientAuthMode for sidecar or sharding. (default will be md5. [md5;scram;cert])
	// +kubebuilder:default=md5
	ClientAuthMode PgpoolClientAuthMode `json:"clientAuthMode,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`
}

// PgpoolStatus defines the observed state of Pgpool
type PgpoolStatus struct {
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

type PgpoolConfiguration struct {
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	PgpoolConfig *runtime.RawExtension `json:"pgpoolConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgpoolList contains a list of Pgpool
type PgpoolList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []Pgpool `json:"items"`
}

// +kubebuilder:validation:Enum=server;client;metrics-exporter
type PgpoolCertificateAlias string

const (
	PgpoolServerCert          PgpoolCertificateAlias = "server"
	PgpoolClientCert          PgpoolCertificateAlias = "client"
	PgpoolMetricsExporterCert PgpoolCertificateAlias = "metrics-exporter"
)

// ref: https://www.postgresql.org/docs/13/libpq-ssl.html
// +kubebuilder:validation:Enum=disable;allow;prefer;require;verify-ca;verify-full
type PgpoolSSLMode string

const (
	// PgpoolSSLModeDisable represents `disable` sslMode. It ensures that the server does not use TLS/SSL.
	PgpoolSSLModeDisable PgpoolSSLMode = "disable"

	// PgpoolSSLModeAllow represents `allow` sslMode. 	I don't care about security,
	// but I will pay the overhead of encryption if the server insists on it.
	PgpoolSSLModeAllow PgpoolSSLMode = "allow"

	// PgpoolSSLModePrefer represents `preferSSL` sslMode.
	// I don't care about encryption, but I wish to pay the overhead of encryption if the server supports it.
	PgpoolSSLModePrefer PgpoolSSLMode = "prefer"

	// PgpoolSSLModeRequire represents `requiteSSL` sslmode. I want my data to be encrypted, and I accept the overhead.
	// I trust that the network will make sure I always connect to the server I want.
	PgpoolSSLModeRequire PgpoolSSLMode = "require"

	// PgpoolSSLModeVerifyCA represents `verify-ca` sslmode. I want my data encrypted, and I accept the overhead.
	// I want to be sure that I connect to a server that I trust.
	PgpoolSSLModeVerifyCA PgpoolSSLMode = "verify-ca"

	// PgpoolSSLModeVerifyFull represents `verify-full` sslmode. I want my data encrypted, and I accept the overhead.
	// I want to be sure that I connect to a server I trust, and that it's the one I specify.
	PgpoolSSLModeVerifyFull PgpoolSSLMode = "verify-full"
)

// PgpoolClientAuthMode represents the ClientAuthMode of Pgpool clusters ( replicaset )
// ref: https://www.postgresql.org/docs/12/auth-methods.html
// +kubebuilder:validation:Enum=md5;scram;cert
type PgpoolClientAuthMode string

const (
	// PgpoolClientAuthModeMD5 uses a custom less secure challenge-response mechanism.
	// It prevents password sniffing and avoids storing passwords on the server in plain text but provides no protection
	// if an attacker manages to steal the password hash from the server.
	// Also, the MD5 hash algorithm is nowadays no longer considered secure against determined attacks
	PgpoolClientAuthModeMD5 PgpoolClientAuthMode = "md5"

	// PgpoolClientAuthModeScram performs SCRAM-SHA-256 authentication, as described in RFC 7677.
	// It is a challenge-response scheme that prevents password sniffing on untrusted connections
	// and supports storing passwords on the server in a cryptographically hashed form that is thought to be secure.
	// This is the most secure of the currently provided methods, but it is not supported by older client libraries.
	PgpoolClientAuthModeScram PgpoolClientAuthMode = "scram"

	// PgpoolClientAuthModeCert represents `cert clientcert=1` auth mode where client need to provide cert and private key for authentication.
	// When server is config with this auth method. Client can't connect with pgpool server with password. They need
	// to Send the client cert and client key certificate for authentication.
	PgpoolClientAuthModeCert PgpoolClientAuthMode = "cert"
)
