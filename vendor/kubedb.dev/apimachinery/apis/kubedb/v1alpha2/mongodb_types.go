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
	ResourceCodeMongoDB     = "mg"
	ResourceKindMongoDB     = "MongoDB"
	ResourceSingularMongoDB = "mongodb"
	ResourcePluralMongoDB   = "mongodbs"

	MongoDBInitInstallContainerName   = "copy-config"
	MongoDBInitBootstrapContainerName = "bootstrap"
	MongoDBConfigDirectoryName        = "configdir"
	MongoDBContainerName              = ResourceSingularMongoDB
)

// MongoDB defines a MongoDB database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mongodbs,singular=mongodb,shortName=mg,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MongoDB struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MongoDBSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            MongoDBStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type MongoDBSpec struct {
	// Version of MongoDB to be deployed.
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`

	// Number of instances to deploy for a MongoDB database.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// MongoDB replica set
	ReplicaSet *MongoDBReplicaSet `json:"replicaSet,omitempty" protobuf:"bytes,3,opt,name=replicaSet"`

	// MongoDB sharding topology.
	ShardTopology *MongoDBShardingTopology `json:"shardTopology,omitempty" protobuf:"bytes,4,opt,name=shardTopology"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty" protobuf:"bytes,5,opt,name=storageType,casttype=StorageType"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty" protobuf:"bytes,6,opt,name=storage"`

	// Database authentication secret
	AuthSecret *core.LocalObjectReference `json:"authSecret,omitempty" protobuf:"bytes,7,opt,name=authSecret"`

	// ClusterAuthMode for replicaset or sharding. (default will be x509 if sslmode is not `disabled`.)
	// See available ClusterAuthMode: https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-clusterauthmode
	ClusterAuthMode ClusterAuthMode `json:"clusterAuthMode,omitempty" protobuf:"bytes,8,opt,name=clusterAuthMode,casttype=ClusterAuthMode"`

	// SSLMode for both standalone and clusters. (default, disabled.)
	// See more options: https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-sslmode
	SSLMode SSLMode `json:"sslMode,omitempty" protobuf:"bytes,9,opt,name=sslMode,casttype=SSLMode"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty" protobuf:"bytes,10,opt,name=init"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,11,opt,name=monitor"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e mongod.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty" protobuf:"bytes,12,opt,name=configSecret"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,13,opt,name=podTemplate"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty" protobuf:"bytes,14,rep,name=serviceTemplates"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty" protobuf:"bytes,15,opt,name=tls"`

	// Secret for KeyFileSecret. Contains keyfile `key.txt` if spec.clusterAuthMode == keyFile || sendKeyFile
	KeyFileSecret *core.LocalObjectReference `json:"keyFileSecret,omitempty" protobuf:"bytes,16,opt,name=keyFileSecret"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty" protobuf:"varint,17,opt,name=halted"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty" protobuf:"bytes,18,opt,name=terminationPolicy,casttype=TerminationPolicy"`

	// StorageEngine can be wiredTiger (default) or inMemory
	// See available StorageEngine: https://docs.mongodb.com/manual/core/storage-engines/
	StorageEngine StorageEngine `json:"storageEngine,omitempty" protobuf:"bytes,19,opt,name=storageEngine,casttype=StorageEngine"`
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
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
}

type MongoDBShardingTopology struct {
	// Shard component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-shards/
	Shard MongoDBShardNode `json:"shard" protobuf:"bytes,1,opt,name=shard"`

	// Config Server (metadata) component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-config-servers/
	ConfigServer MongoDBConfigNode `json:"configServer" protobuf:"bytes,2,opt,name=configServer"`

	// Mongos (router) component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-query-router/
	Mongos MongoDBMongosNode `json:"mongos" protobuf:"bytes,3,opt,name=mongos"`
}

type MongoDBShardNode struct {
	// Shards represents number of shards for shard type of node
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-shards/
	Shards int32 `json:"shards" protobuf:"varint,1,opt,name=shards"`

	// MongoDB sharding node configs
	MongoDBNode `json:",inline" protobuf:"bytes,2,opt,name=mongoDBNode"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty" protobuf:"bytes,3,opt,name=storage"`
}

type MongoDBConfigNode struct {
	// MongoDB config server node configs
	MongoDBNode `json:",inline" protobuf:"bytes,1,opt,name=mongoDBNode"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty" protobuf:"bytes,2,opt,name=storage"`
}

type MongoDBMongosNode struct {
	// MongoDB mongos node configs
	MongoDBNode `json:",inline" protobuf:"bytes,1,opt,name=mongoDBNode"`
}

type MongoDBNode struct {
	// Replicas represents number of replicas of this specific node.
	// If current node has replicaset enabled, then replicas is the amount of replicaset nodes.
	Replicas int32 `json:"replicas" protobuf:"varint,1,opt,name=replicas"`

	// Prefix is the name prefix of this node.
	Prefix string `json:"prefix,omitempty" protobuf:"bytes,2,opt,name=prefix"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e mongod.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty" protobuf:"bytes,3,opt,name=configSecret"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,4,opt,name=podTemplate"`
}

type MongoDBStatus struct {
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

type MongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of MongoDB TPR objects
	Items []MongoDB `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
