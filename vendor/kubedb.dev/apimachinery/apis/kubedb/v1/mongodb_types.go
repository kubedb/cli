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
	ResourceCodeMongoDB     = "mg"
	ResourceKindMongoDB     = "MongoDB"
	ResourceSingularMongoDB = "mongodb"
	ResourcePluralMongoDB   = "mongodbs"
)

// MongoDB defines a MongoDB database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mongodbs,singular=mongodb,shortName=mg,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MongoDB struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MongoDBSpec   `json:"spec,omitempty"`
	Status            MongoDBStatus `json:"status,omitempty"`
}

type MongoDBSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of MongoDB to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a MongoDB database.
	Replicas *int32 `json:"replicas,omitempty"`

	// MongoDB replica set
	ReplicaSet *MongoDBReplicaSet `json:"replicaSet,omitempty"`

	// MongoDB sharding topology.
	ShardTopology *MongoDBShardingTopology `json:"shardTopology,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// EphemeralStorage spec to specify the configuration of ephemeral storage type.
	EphemeralStorage *core.EmptyDirVolumeSource `json:"ephemeralStorage,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// ClusterAuthMode for replicaset or sharding. (default will be x509 if sslmode is not `disabled`.)
	// See available ClusterAuthMode: https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-clusterauthmode
	ClusterAuthMode ClusterAuthMode `json:"clusterAuthMode,omitempty"`

	// SSLMode for both standalone and clusters. (default, disabled.)
	// See more options: https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-sslmode
	SSLMode SSLMode `json:"sslMode,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e mongod.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// Secret for KeyFileSecret. Contains keyfile `key.txt` if spec.clusterAuthMode == keyFile || sendKeyFile
	KeyFileSecret *core.LocalObjectReference `json:"keyFileSecret,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// StorageEngine can be wiredTiger (default) or inMemory
	// See available StorageEngine: https://docs.mongodb.com/manual/core/storage-engines/
	StorageEngine StorageEngine `json:"storageEngine,omitempty"`

	// AllowedSchemas defines the types of database schemas that MAY refer to
	// a database instance and the trusted namespaces where those schema resources MAY be
	// present.
	//
	// +kubebuilder:default={namespaces:{from: Same}}
	// +optional
	AllowedSchemas *AllowedConsumers `json:"allowedSchemas,omitempty"`

	// Mongo Arbiter component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/replica-set-arbiter/
	// +optional
	// +nullable
	Arbiter *MongoArbiterNode `json:"arbiter,omitempty"`

	// Hidden component of mongodb which is invisible to client applications
	// More info: https://www.mongodb.com/docs/manual/core/replica-set-hidden-member/
	// +optional
	// +nullable
	Hidden *MongoHiddenNode `json:"hidden,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Archiver controls database backup using Archiver CR
	// +optional
	Archiver *Archiver `json:"archiver,omitempty"`
}

// +kubebuilder:validation:Enum=server;client;metrics-exporter
type MongoDBCertificateAlias string

const (
	MongoDBServerCert          MongoDBCertificateAlias = "server"
	MongoDBClientCert          MongoDBCertificateAlias = "client"
	MongoDBMetricsExporterCert MongoDBCertificateAlias = "metrics-exporter"
)

// ClusterAuthMode represents the clusterAuthMode of mongodb clusters ( replicaset or sharding)
// ref: https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-clusterauthmode
// +kubebuilder:validation:Enum=keyFile;sendKeyFile;sendX509;x509
type ClusterAuthMode string

const (
	// ClusterAuthModeKeyFile represents `keyFile` mongodb clusterAuthMode. In this mode, Use a keyfile for authentication. Accept only keyfiles.
	ClusterAuthModeKeyFile ClusterAuthMode = "keyFile"

	// ClusterAuthModeSendKeyFile represents `sendKeyFile` mongodb clusterAuthMode.
	// This mode is for rolling upgrade purposes. Send a keyfile for authentication but can accept both keyfiles
	// and x.509 certificates.
	ClusterAuthModeSendKeyFile ClusterAuthMode = "sendKeyFile"

	// ClusterAuthModeSendX509 represents `sendx509` mongodb clusterAuthMode. This mode is usually for rolling upgrade purposes.
	// Send the x.509 certificate for authentication but can accept both keyfiles and x.509 certificates.
	ClusterAuthModeSendX509 ClusterAuthMode = "sendX509"

	// ClusterAuthModeX509 represents `x509` mongodb clusterAuthMode. This is the recommended clusterAuthMode.
	// Send the x.509 certificate for authentication and accept only x.509 certificates.
	ClusterAuthModeX509 ClusterAuthMode = "x509"
)

// SSLMode represents available sslmodes of mongodb.
// ref: https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-sslmode
// +kubebuilder:validation:Enum=disabled;allowSSL;preferSSL;requireSSL
type SSLMode string

const (
	// SSLModeDisabled represents `disabled` sslMode. It ensures that the server does not use TLS/SSL.
	SSLModeDisabled SSLMode = "disabled"

	// SSLModeAllowSSL represents `allowSSL` sslMode. It ensures that the connections between servers do not use TLS/SSL. For incoming connections,
	// the server accepts both TLS/SSL and non-TLS/non-SSL.
	SSLModeAllowSSL SSLMode = "allowSSL"

	// SSLModePreferSSL represents `preferSSL` sslMode. It ensures that the connections between servers use TLS/SSL. For incoming connections,
	// the server accepts both TLS/SSL and non-TLS/non-SSL.
	SSLModePreferSSL SSLMode = "preferSSL"

	// SSLModeRequireSSL represents `requiteSSL` sslmode. It ensures that the server uses and accepts only TLS/SSL encrypted connections.
	SSLModeRequireSSL SSLMode = "requireSSL"
)

// StorageEngine represents storage engine of mongodb clusters.
// ref: https://docs.mongodb.com/manual/core/storage-engines/
type StorageEngine string

const (
	// StorageEngineWiredTiger represents `wiredTiger` storage engine of mongodb.
	StorageEngineWiredTiger StorageEngine = "wiredTiger"

	// StorageEngineInMemory represents `inMemory` storage engine of mongodb.
	StorageEngineInMemory StorageEngine = "inMemory"
)

type MongoDBReplicaSet struct {
	// Name of replicaset
	Name string `json:"name"`
}

type MongoDBShardingTopology struct {
	// Shard component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-shards/
	Shard MongoDBShardNode `json:"shard"`

	// Config Server (metadata) component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-config-servers/
	ConfigServer MongoDBConfigNode `json:"configServer"`

	// Mongos (router) component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-query-router/
	Mongos MongoDBMongosNode `json:"mongos"`
}

type MongoDBShardNode struct {
	// Shards represents number of shards for shard type of node
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-shards/
	Shards int32 `json:"shards"`

	// MongoDB sharding node configs
	MongoDBNode `json:",inline"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// EphemeralStorage spec to specify the configuration of ephemeral storage type.
	EphemeralStorage *core.EmptyDirVolumeSource `json:"ephemeralStorage,omitempty"`
}

type MongoDBConfigNode struct {
	// MongoDB config server node configs
	MongoDBNode `json:",inline"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// EphemeralStorage spec to specify the configuration of ephemeral storage type.
	EphemeralStorage *core.EmptyDirVolumeSource `json:"ephemeralStorage,omitempty"`
}

type MongoDBMongosNode struct {
	// MongoDB mongos node configs
	MongoDBNode `json:",inline"`
}

type MongoArbiterNode struct {
	// ConfigSecret is an optional field to provide custom configuration file for database (i.e mongod.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`
}

type MongoHiddenNode struct {
	// ConfigSecret is an optional field to provide custom configuration file for database (i.e mongod.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// Replicas represents number of replicas of this specific node.
	// If current node has replicaset enabled, then replicas is the amount of replicaset nodes.
	Replicas int32 `json:"replicas"`

	// Storage to specify how storage shall be used.
	Storage core.PersistentVolumeClaimSpec `json:"storage"`
}

type MongoDBNode struct {
	// Replicas represents number of replicas of this specific node.
	// If current node has replicaset enabled, then replicas is the amount of replicaset nodes.
	Replicas int32 `json:"replicas"`

	// Prefix is the name prefix of this node.
	Prefix string `json:"prefix,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e mongod.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`
}

type MongoDBStatus struct {
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

type MongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MongoDB TPR objects
	Items []MongoDB `json:"items,omitempty"`
}
