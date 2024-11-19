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
	ResourceCodeMySQL     = "my"
	ResourceKindMySQL     = "MySQL"
	ResourceSingularMySQL = "mysql"
	ResourcePluralMySQL   = "mysqls"
)

// +kubebuilder:validation:Enum=GroupReplication;InnoDBCluster;RemoteReplica;SemiSync
type MySQLMode string

const (
	MySQLModeGroupReplication MySQLMode = "GroupReplication"
	MySQLModeInnoDBCluster    MySQLMode = "InnoDBCluster"
	MySQLModeRemoteReplica    MySQLMode = "RemoteReplica"
	MySQLModeSemiSync         MySQLMode = "SemiSync"
)

// +kubebuilder:validation:Enum=Single-Primary;Multi-Primary
type MySQLGroupMode string

const (
	MySQLGroupModeSinglePrimary MySQLGroupMode = "Single-Primary"
	MySQLGroupModeMultiPrimary  MySQLGroupMode = "Multi-Primary"
)

// Mysql defines a Mysql database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mysqls,singular=mysql,shortName=my,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MySQL struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MySQLSpec   `json:"spec,omitempty"`
	Status            MySQLStatus `json:"status,omitempty"`
}

type MySQLSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of MySQL to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a MySQL database. In case of MySQL group
	// replication, max allowed value is 9 (default 3).
	// (see ref: https://dev.mysql.com/doc/refman/5.7/en/group-replication-frequently-asked-questions.html)
	Replicas *int32 `json:"replicas,omitempty"`

	// MySQL cluster topology
	Topology *MySQLTopology `json:"topology,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database (i.e custom-mysql.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// Indicates that the database server need to be encrypted connections(ssl)
	// +optional
	RequireSSL bool `json:"requireSSL,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// Indicated whether to use DNS or IP address to address pods in a db cluster.
	// If IP address is used, HostNetwork will be used. Defaults to DNS.
	// +kubebuilder:default=DNS
	// +optional
	// +default="DNS"
	UseAddressType AddressType `json:"useAddressType,omitempty"`

	// AllowedSchemas defines the types of database schemas that may refer to
	// a database instance and the trusted namespaces where those schema resources may be
	// present.
	//
	// +kubebuilder:default={namespaces:{from: Same}}
	// +optional
	AllowedSchemas *AllowedConsumers `json:"allowedSchemas,omitempty"`

	// AllowedReadReplicas defines the types of read replicas that MAY be attached to a
	// MySQL instance and the trusted namespaces where those Read Replica resources MAY be
	// present.
	//
	// Support: Core
	// +kubebuilder:default={namespaces:{from: Same}}
	// +optional
	AllowedReadReplicas *AllowedConsumers `json:"allowedReadReplicas,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Archiver controls database backup using Archiver CR
	// +optional
	Archiver *Archiver `json:"archiver,omitempty"`
}

// +kubebuilder:validation:Enum=server;client;metrics-exporter
type MySQLCertificateAlias string

const (
	MySQLServerCert          MySQLCertificateAlias = "server"
	MySQLClientCert          MySQLCertificateAlias = "client"
	MySQLMetricsExporterCert MySQLCertificateAlias = "metrics-exporter"
	MySQLRouterCert          MySQLCertificateAlias = "router"
)

type MySQLTopology struct {
	// If set to -
	// "GroupReplication", GroupSpec is required and MySQL servers will start  a replication group
	Mode *MySQLMode `json:"mode,omitempty"`

	// Group replication info for MySQL
	// +optional
	Group *MySQLGroupSpec `json:"group,omitempty"`

	// InnoDBCluster replication info for MySQL InnodbCluster
	// +optional
	InnoDBCluster *MySQLInnoDBClusterSpec `json:"innoDBCluster,omitempty"`

	// RemoteReplica implies that the instance will be a MySQL Read Only Replica
	// and it will take reference of  appbinding of the source
	// +optional
	RemoteReplica *RemoteReplicaSpec `json:"remoteReplica,omitempty"`
	// +optional
	SemiSync *SemiSyncSpec `json:"semiSync,omitempty"`
}

// +kubebuilder:validation:Enum= Clone;PseudoTransaction

type ErrantTransactionRecoveryPolicy string

const (
	ErrantTransactionRecoveryPolicyClone             ErrantTransactionRecoveryPolicy = "Clone"
	ErrantTransactionRecoveryPolicyPseudoTransaction ErrantTransactionRecoveryPolicy = "PseudoTransaction"
)

type SemiSyncSpec struct {
	// count of slave to wait for before commit
	// +kubebuilder:default=1
	//+kubebuilder:validation:Minimum=1
	SourceWaitForReplicaCount int `json:"sourceWaitForReplicaCount,omitempty"`
	// +kubebuilder:default="24h"
	SourceTimeout metav1.Duration `json:"sourceTimeout,omitempty"`
	// recovery method if the slave has any errant transaction
	// +kubebuilder:default=PseudoTransaction
	ErrantTransactionRecoveryPolicy *ErrantTransactionRecoveryPolicy `json:"errantTransactionRecoveryPolicy"`
}

type MySQLGroupSpec struct {
	// TODO: "Multi-Primary" needs to be implemented
	// Group Replication can be deployed in either "Single-Primary" or "Multi-Primary" mode
	// +kubebuilder:default=Single-Primary
	Mode *MySQLGroupMode `json:"mode,omitempty"`

	// Group name is a version 4 UUID
	// ref: https://dev.mysql.com/doc/refman/5.7/en/group-replication-options.html#sysvar_group_replication_group_name
	Name string `json:"name,omitempty"`
}

type MySQLInnoDBClusterSpec struct {
	// +kubebuilder:default=Single-Primary
	// +optional
	Mode *MySQLGroupMode `json:"mode,omitempty"`

	Router MySQLRouterSpec `json:"router,omitempty"`
}

type MySQLRouterSpec struct {
	// +optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum:=1
	Replicas *int32 `json:"replicas,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose MySQL router
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`
}

type MySQLStatus struct {
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

type MySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MySQL TPR objects
	Items []MySQL `json:"items,omitempty"`
}
