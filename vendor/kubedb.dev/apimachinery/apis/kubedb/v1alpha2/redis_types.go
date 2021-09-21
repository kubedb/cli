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
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Redis struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              RedisSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            RedisStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type RedisSpec struct {
	// Version of Redis to be deployed.
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`

	// Number of instances to deploy for a Redis database.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// Default is "Standalone". If set to "Cluster", ClusterSpec is required and redis servers will
	// start in cluster mode
	Mode RedisMode `json:"mode,omitempty" protobuf:"bytes,3,opt,name=mode,casttype=RedisMode"`

	SentinelRef *RedisSentinelRef `json:"sentinelRef,omitempty" protobuf:"bytes,4,opt,name=sentinelRef"`

	// Redis cluster configuration for running redis servers in cluster mode. Required if Mode is set to "Cluster"
	Cluster *RedisClusterSpec `json:"cluster,omitempty" protobuf:"bytes,5,opt,name=cluster"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty" protobuf:"bytes,6,opt,name=storageType,casttype=StorageType"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty" protobuf:"bytes,7,opt,name=storage"`

	// Database authentication secret
	AuthSecret *core.LocalObjectReference `json:"authSecret,omitempty" protobuf:"bytes,8,opt,name=authSecret"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty" protobuf:"bytes,9,opt,name=init"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,10,opt,name=monitor"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e redis.conf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty" protobuf:"bytes,11,opt,name=configSecret"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,12,opt,name=podTemplate"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty" protobuf:"bytes,13,rep,name=serviceTemplates"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty" protobuf:"bytes,14,opt,name=tls"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty" protobuf:"varint,15,opt,name=halted"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty" protobuf:"bytes,16,opt,name=terminationPolicy,casttype=TerminationPolicy"`

	// Coordinator defines attributes of the coordinator container
	// +optional
	Coordinator CoordinatorSpec `json:"coordinator,omitempty" protobuf:"bytes,17,opt,name=coordinator"`
}

// +kubebuilder:validation:Enum=server;client;metrics-exporter
type RedisCertificateAlias string

const (
	RedisServerCert          RedisCertificateAlias = "server"
	RedisClientCert          RedisCertificateAlias = "client"
	RedisMetricsExporterCert RedisCertificateAlias = "metrics-exporter"
)

type RedisClusterSpec struct {
	// Number of master nodes. It must be >= 3. If not specified, defaults to 3.
	Master *int32 `json:"master,omitempty" protobuf:"varint,1,opt,name=master"`

	// Number of replica(s) per master node. If not specified, defaults to 1.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`
}

type RedisSentinelRef struct {
	// Name of the refereed sentinel
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Namespace where refereed sentinel has been deployed
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
}

type RedisStatus struct {
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

type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Redis TPR objects
	Items []Redis `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
