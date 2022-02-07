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
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
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
// +kubebuilder:resource:path=pgbouncers,singular=pgbouncer,shortName=pb,categories={proxy,kubedb,appscode,all}
// +kubebuilder:subresource:status
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
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`
	// Databases to proxy by connection pooling.
	// +optional
	Databases []Databases `json:"databases,omitempty"`
	// ConnectionPoolConfig defines Connection pool configuration.
	// +optional
	ConnectionPool *ConnectionPoolConfig `json:"connectionPool,omitempty"`
	// UserListSecretRef is a secret with a list of PgBouncer user and passwords.
	// +optional
	UserListSecretRef *core.LocalObjectReference `json:"userListSecretRef,omitempty"`
	// Monitor is used monitor database instance.
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`
	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`
}

// +kubebuilder:validation:Enum=server;archiver;metrics-exporter
type PgBouncerCertificateAlias string

const (
	PgBouncerServerCert          PgBouncerCertificateAlias = "server"
	PgBouncerClientCert          PgBouncerCertificateAlias = "client"
	PgBouncerMetricsExporterCert PgBouncerCertificateAlias = "metrics-exporter"
)

type Databases struct {
	// Alias to uniquely identify a target database running inside a specific Postgres instance.
	Alias string `json:"alias"`
	// DatabaseRef specifies the database appbinding reference in any namespace.
	DatabaseRef appcat.AppReference `json:"databaseRef"`
	// DatabaseName is the name of the target database inside a Postgres instance.
	DatabaseName string `json:"databaseName"`
	// AuthSecretRef points to a secret that contains the credentials
	// (username and password) of an existing user of this database.
	// It is used to bind a single user to this specific database connection.
	// +optional
	AuthSecretRef *core.LocalObjectReference `json:"authSecretRef,omitempty"`
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
	// MaxDBConnections is the maximum number of connections allowed per-database. Default: unlimited.
	// +optional
	MaxDBConnections *int64 `json:"maxDBConnections,omitempty"`
	// MaxUserConnections is the maximum number of users allowed per-database. Default: unlimited.
	// +optional
	MaxUserConnections *int64 `json:"maxUserConnections,omitempty"`
	// StatsPeriodSeconds sets how often the averages shown in various SHOW commands are updated
	// and how often aggregated statistics are written to the log.
	// +optional
	StatsPeriodSeconds *int64 `json:"statsPeriodSeconds,omitempty"`
	// AdminUsers specifies an array of users who can act as PgBouncer administrators.
	// +optional
	AdminUsers []string `json:"adminUsers,omitempty"`
	// AuthType specifies how to authenticate users. Default: md5 (md5+plain text).
	// +optional
	AuthType string `json:"authType,omitempty"`
	// AuthUser looks up any user not specified in auth_file from pg_shadow. Default: not set.
	// +optional
	AuthUser string `json:"authUser,omitempty"`
	// IgnoreStartupParameters specifies comma-separated startup parameters that
	// pgbouncer knows are handled by admin and it can ignore them.
	// +optional
	IgnoreStartupParameters string `json:"ignoreStartupParameters,omitempty"`
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
