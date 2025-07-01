/*
Copyright 2023.

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
	ofst "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeDruid     = "dr"
	ResourceKindDruid     = "Druid"
	ResourceSingularDruid = "druid"
	ResourcePluralDruid   = "druids"
)

// Druid is the Schema for the druids API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=druids,singular=druid,shortName=dr,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Druid struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DruidSpec   `json:"spec,omitempty"`
	Status DruidStatus `json:"status,omitempty"`
}

// DruidSpec defines the desired state of Druid
type DruidSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of Druid to be deployed.
	Version string `json:"version"`

	// Druid topology for node specification
	// +optional
	Topology *DruidClusterTopology `json:"topology,omitempty"`

	// disable security. It disables authentication security of user.
	// If unset, default is false
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e. config.properties).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// To enable ssl for http layer
	EnableSSL bool `json:"enableSSL,omitempty"`

	// Keystore encryption secret
	// +optional
	KeystoreCredSecret *SecretReference `json:"keystoreCredSecret,omitempty"`

	// TLS contains tls configurations
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// MetadataStorage contains information for Druid to connect to external dependency metadata storage
	// +optional
	MetadataStorage *MetadataStorage `json:"metadataStorage,omitempty"`

	// DeepStorage contains specification for druid to connect to the deep storage
	DeepStorage *DeepStorageSpec `json:"deepStorage"`

	// ZooKeeper contains information for Druid to connect to external dependency metadata storage
	// +optional
	ZookeeperRef *ZookeeperRef `json:"zookeeperRef,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 30, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`
}

type DruidClusterTopology struct {
	Coordinators *DruidNode `json:"coordinators,omitempty"`
	// +optional
	Overlords *DruidNode `json:"overlords,omitempty"`

	MiddleManagers *DruidDataNode `json:"middleManagers,omitempty"`

	Historicals *DruidDataNode `json:"historicals,omitempty"`

	Brokers *DruidNode `json:"brokers,omitempty"`
	// +optional
	Routers *DruidNode `json:"routers,omitempty"`
}

type DruidNode struct {
	// Replicas represents number of replicas for the specific type of node
	// +kubebuilder:default=1
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Suffix to append with node name
	// +optional
	Suffix string `json:"suffix,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`
}

type DruidDataNode struct {
	// DruidDataNode has all the characteristics of DruidNode
	DruidNode `json:",inline"`

	// StorageType specifies if the storage
	// of this node is durable (default) or ephemeral.
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// EphemeralStorage spec to specify the configuration of ephemeral storage type.
	EphemeralStorage *core.EmptyDirVolumeSource `json:"ephemeralStorage,omitempty"`
}

type MetadataStorage struct {
	// Name and namespace of the appbinding of metadata storage
	// +optional
	*kmapi.ObjectReference `json:",omitempty"`

	// If not KubeDB managed, then specify type of the metadata storage
	// +optional
	Type DruidMetadataStorageType `json:"type,omitempty"`

	// If Druid has the permission to create new tables
	// +optional
	CreateTables *bool `json:"createTables,omitempty"`

	// +optional
	LinkedDB string `json:"linkedDB,omitempty"`

	// +optional
	ExternallyManaged bool `json:"externallyManaged,omitempty"`

	// Version of the MySQL/PG used
	// +optional
	Version *string `json:"version,omitempty"`
}

type DeepStorageSpec struct {
	// Specifies the storage type to be used by druid
	// Possible values: s3, google, azure, hdfs
	Type DruidDeepStorageType `json:"type"`

	// deepStorage.configSecret should contain the necessary data
	// to connect to the deep storage
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
}

type ZookeeperRef struct {
	// Name and namespace of appbinding of zookeeper
	// +optional
	*kmapi.ObjectReference `json:",omitempty"`

	// Base ZooKeeperSpec path
	// +optional
	PathsBase string `json:"pathsBase,omitempty"`

	// +optional
	ExternallyManaged bool `json:"externallyManaged,omitempty"`

	// Version of the ZK used
	// +optional
	Version *string `json:"version,omitempty"`
}

// DruidStatus defines the observed state of Druid
type DruidStatus struct {
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

// DruidList contains a list of Druid
type DruidList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Druid `json:"items"`
}

// +kubebuilder:validation:Enum=coordinators;overlords;brokers;routers;middleManagers;historicals
type DruidNodeRoleType string

const (
	DruidNodeRoleCoordinators   DruidNodeRoleType = "coordinators"
	DruidNodeRoleOverlords      DruidNodeRoleType = "overlords"
	DruidNodeRoleBrokers        DruidNodeRoleType = "brokers"
	DruidNodeRoleRouters        DruidNodeRoleType = "routers"
	DruidNodeRoleMiddleManagers DruidNodeRoleType = "middleManagers"
	DruidNodeRoleHistoricals    DruidNodeRoleType = "historicals"
)

// +kubebuilder:validation:Enum=MySQL;PostgreSQL
type DruidMetadataStorageType string

const (
	DruidMetadataStorageMySQL      DruidMetadataStorageType = "MySQL"
	DruidMetadataStoragePostgreSQL DruidMetadataStorageType = "PostgreSQL"
)

// +kubebuilder:validation:Enum=s3;google;azure;hdfs
type DruidDeepStorageType string

const (
	DruidDeepStorageS3     DruidDeepStorageType = "s3"
	DruidDeepStorageGoogle DruidDeepStorageType = "google"
	DruidDeepStorageAzure  DruidDeepStorageType = "azure"
	DruidDeepStorageHDFS   DruidDeepStorageType = "hdfs"
)

// +kubebuilder:validation:Enum=server;client
type DruidCertificateAlias string

const (
	DruidServerCert DruidCertificateAlias = "server"
	DruidClientCert DruidCertificateAlias = "client"
)
