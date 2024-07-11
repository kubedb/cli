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
	ResourceCodeRedis     = "rd"
	ResourceKindRedis     = "Redis"
	ResourceSingularRedis = "redis"
	ResourcePluralRedis   = "redises"
)

// +kubebuilder:validation:Enum=Standalone;Cluster;Sentinel
type RedisMode string

const (
	RedisModeStandalone RedisMode = "Standalone"
	RedisModeCluster    RedisMode = "Cluster"
	RedisModeSentinel   RedisMode = "Sentinel"
)

// Redis defines a Redis database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=redises,singular=redis,shortName=rd,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Redis struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RedisSpec   `json:"spec,omitempty"`
	Status            RedisStatus `json:"status,omitempty"`
}

type RedisSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of Redis to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Redis database.
	Replicas *int32 `json:"replicas,omitempty"`

	// Default is "Standalone". If set to "Cluster", ClusterSpec is required and redis servers will
	// start in cluster mode
	Mode RedisMode `json:"mode,omitempty"`

	SentinelRef *RedisSentinelRef `json:"sentinelRef,omitempty"`

	// Redis cluster configuration for running redis servers in cluster mode. Required if Mode is set to "Cluster"
	Cluster *RedisClusterSpec `json:"cluster,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// If disable Auth true then don't create any auth secret
	// +optional
	DisableAuth bool `json:"disableAuth,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e redis.conf).
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
}

// +kubebuilder:validation:Enum=server;client;metrics-exporter
type RedisCertificateAlias string

const (
	RedisServerCert          RedisCertificateAlias = "server"
	RedisClientCert          RedisCertificateAlias = "client"
	RedisMetricsExporterCert RedisCertificateAlias = "metrics-exporter"
)

type RedisClusterSpec struct {
	// Number of shards. It must be >= 3. If not specified, defaults to 3.
	Shards *int32 `json:"shards,omitempty"`

	// Number of replica(s) per shard. If not specified, defaults to 2.
	Replicas *int32 `json:"replicas,omitempty"`
}

type RedisSentinelRef struct {
	// Name of the refereed sentinel
	Name string `json:"name,omitempty"`

	// Namespace where refereed sentinel has been deployed
	Namespace string `json:"namespace,omitempty"`
}

type RedisStatus struct {
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

type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Redis TPR objects
	Items []Redis `json:"items,omitempty"`
}
