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
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeAerospike     = "ar"
	ResourceKindAerospike     = "Aerospike"
	ResourceSingularAerospike = "aerospike"
	ResourcePluralAerospike   = "aerospikes"
)

// +kubebuilder:validation:Enum=Standalone;Cluster
type AerospikeMode string

const (
	AerospikeModeStandalone AerospikeMode = "Standalone"
	AerospikeModeCluster    AerospikeMode = "Cluster"
)

// Aerospike is the Schema for the aerospikes API
// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=aerospikes,singular=aerospike,shortName=ar,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Aerospike struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            AerospikeSpec   `json:"spec,omitempty"`
	Status          AerospikeStatus `json:"status,omitempty"`
}

// AerospikeSpec defines the desired state of Aerospike
type AerospikeSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of Aerospike to be deployed.
	// +optional
	Version string `json:"version"`

	// Number of instances to deploy for a Aerospike instance.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Default is "Standalone".
	Mode AerospikeMode `json:"mode,omitempty"`

	Cluster *AerospikeClusterSpec `json:"cluster,omitempty"`

	// Aerospike secret containing username and password for aerospike pcp user
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// +optional
	Configuration *ConfigurationSpec `json:"configuration,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose Aerospike
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose Aerospike
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Monitor is used to monitor Aerospike instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	// +optional
	SSLMode AerospikeSSLMode `json:"sslMode,omitempty"`

	// ClientAuthMode for sidecar or sharding. (default will be md5. [md5;scram;cert])
	// +optional
	ClientAuthMode AerospikeClientAuthMode `json:"clientAuthMode,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`
}

type AerospikeClusterSpec struct {
	Replicas          *int32 `json:"replicas,omitempty"`
	ReplicationFactor *int32 `json:"replicationFactor,omitempty"`
}

// AerospikeStatus defines the observed state of Aerospike
type AerospikeStatus struct {
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AerospikeList contains a list of Aerospike
type AerospikeList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []Aerospike `json:"items"`
}

// +kubebuilder:validation:Enum=server;client;metrics-exporter
type AerospikeCertificateAlias string

const (
	AerospikeServerCert          AerospikeCertificateAlias = "server"
	AerospikeClientCert          AerospikeCertificateAlias = "client"
	AerospikeMetricsExporterCert AerospikeCertificateAlias = "metrics-exporter"
)

// ref: https://www.postgresql.org/docs/13/libpq-ssl.html
// +kubebuilder:validation:Enum=disable;allow;prefer;require;verify-ca;verify-full
type AerospikeSSLMode string

const (
	// AerospikeSSLModeDisable represents `disable` sslMode. It ensures that the server does not use TLS/SSL.
	AerospikeSSLModeDisable AerospikeSSLMode = "disable"

	// AerospikeSSLModeAllow represents `allow` sslMode. 	I don't care about security,
	// but I will pay the overhead of encryption if the server insists on it.
	AerospikeSSLModeAllow AerospikeSSLMode = "allow"

	// AerospikeSSLModePrefer represents `preferSSL` sslMode.
	// I don't care about encryption, but I wish to pay the overhead of encryption if the server supports it.
	AerospikeSSLModePrefer AerospikeSSLMode = "prefer"

	// AerospikeSSLModeRequire represents `requiteSSL` sslmode. I want my data to be encrypted, and I accept the overhead.
	// I trust that the network will make sure I always connect to the server I want.
	AerospikeSSLModeRequire AerospikeSSLMode = "require"

	// AerospikeSSLModeVerifyCA represents `verify-ca` sslmode. I want my data encrypted, and I accept the overhead.
	// I want to be sure that I connect to a server that I trust.
	AerospikeSSLModeVerifyCA AerospikeSSLMode = "verify-ca"

	// AerospikeSSLModeVerifyFull represents `verify-full` sslmode. I want my data encrypted, and I accept the overhead.
	// I want to be sure that I connect to a server I trust, and that it's the one I specify.
	AerospikeSSLModeVerifyFull AerospikeSSLMode = "verify-full"
)

// AerospikeClientAuthMode represents the ClientAuthMode of Aerospike clusters ( replicaset )
// ref: https://www.postgresql.org/docs/12/auth-methods.html
// +kubebuilder:validation:Enum=md5;scram;cert
type AerospikeClientAuthMode string

const (
	// AerospikeClientAuthModeMD5 uses a custom less secure challenge-response mechanism.
	// It prevents password sniffing and avoids storing passwords on the server in plain text but provides no protection
	// if an attacker manages to steal the password hash from the server.
	// Also, the MD5 hash algorithm is nowadays no longer considered secure against determined attacks
	AerospikeClientAuthModeMD5 AerospikeClientAuthMode = "md5"

	// AerospikeClientAuthModeScram performs SCRAM-SHA-256 authentication, as described in RFC 7677.
	// It is a challenge-response scheme that prevents password sniffing on untrusted connections
	// and supports storing passwords on the server in a cryptographically hashed form that is thought to be secure.
	// This is the most secure of the currently provided methods, but it is not supported by older client libraries.
	AerospikeClientAuthModeScram AerospikeClientAuthMode = "scram"

	// AerospikeClientAuthModeCert represents `cert clientcert=1` auth mode where client need to provide cert and private key for authentication.
	// When server is config with this auth method. Client can't connect with Aerospike server with password. They need
	// to Send the client cert and client key certificate for authentication.
	AerospikeClientAuthModeCert AerospikeClientAuthMode = "cert"
)

var _ Accessor = &Aerospike{}

func (p *Aerospike) GetObjectMeta() meta.ObjectMeta {
	return p.ObjectMeta
}

func (p *Aerospike) GetConditions() []kmapi.Condition {
	return p.Status.Conditions
}

func (p *Aerospike) SetCondition(cond kmapi.Condition) {
	p.Status.Conditions = setCondition(p.Status.Conditions, cond)
}

func (p *Aerospike) RemoveCondition(typ string) {
	p.Status.Conditions = removeCondition(p.Status.Conditions, typ)
}
