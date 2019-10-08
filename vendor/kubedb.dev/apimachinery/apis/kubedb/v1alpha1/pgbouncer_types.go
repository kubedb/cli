package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	v1 "kmodules.xyz/offshoot-api/api/v1"
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
	Version types.StrYo `json:"version"`
	// Number of instances to deploy for a PgBouncer instance.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// ServiceTemplate is an optional configuration for service used to expose database
	// +optional
	ServiceTemplate v1.ServiceTemplateSpec `json:"serviceTemplate,omitempty"`
	// PodTemplate is an optional configuration for pods
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`
	// Databases to proxy by connection pooling
	// +optional
	Databases []Databases `json:"databases,omitempty"`
	// ConnectionPoolConfig defines Connection pool configuration
	// +optional
	ConnectionPool *ConnectionPoolConfig `json:"connectionPool,omitempty"`
	// UserListSecretRef is a secret with a list of PgBouncer user and passwords
	// +optional
	UserListSecretRef *core.LocalObjectReference `json:"userListSecretRef,omitempty"`
	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`
}

type Databases struct {
	//Alias to uniquely identify a target database running inside a specific Postgres instance
	Alias string `json:"alias"`
	// DatabaseRef specifies the database appbinding reference in any namespace
	DatabaseRef appcat.AppReference `json:"databaseRef"`
	// DatabaseName is the name of the target database inside a Postgres instance
	DatabaseName string `json:"databaseName"`
	//UserName is used to bind a single user to a specific database connection
	// +optional
	UserName string `json:"username,omitempty"`
	//Password is to authenticate the user specified in Username field
	// +optional
	Password string `json:"password,omitempty"`
}

type ConnectionPoolConfig struct {
	//Port is the port number on which PgBouncer listens to clients. Default: 5432.
	// +optional
	Port *int32 `json:"port,omitempty"`
	//PoolMode is the pooling mechanism type. Default: session.
	// +optional
	PoolMode string `json:"poolMode,omitempty"`
	//MaxClientConnections is the maximum number of allowed client connections. Default: 100.
	// +optional
	MaxClientConnections *int `json:"maxClientConnections,omitempty"`
	//DefaultPoolSize specifies how many server connections to allow per user/database pair. Default: 20.
	// +optional
	DefaultPoolSize *int `json:"defaultPoolSize,omitempty"`
	//MinPoolSize is used to add more server connections to pool if below this number. Default: 0 (disabled).
	// +optional
	MinPoolSize *int `json:"minPoolSize,omitempty"`
	//ReservePoolSize specifies how many additional connections to allow to a pool. 0 disables. Default: 0 (disabled)
	// +optional
	ReservePoolSize *int `json:"reservePoolSize,omitempty"`
	//ReservePoolTimeoutSeconds is the number of seconds in which if a client has not been serviced,
	//pgbouncer enables use of additional connections from reserve pool. 0 disables. Default: 5.0
	// +optional
	ReservePoolTimeoutSeconds *int `json:"reservePoolTimeoutSeconds,omitempty"`
	//MaxDBConnections is the maximum number of connections allowed per-database. Default: unlimited.
	// +optional
	MaxDBConnections *int `json:"maxDBConnections,omitempty"`
	//MaxUserConnections is the maximum number of users allowed per-database. Default: unlimited.
	// +optional
	MaxUserConnections *int `json:"maxUserConnections,omitempty"`
	//StatsPeriodSeconds sets how often the averages shown in various SHOW commands are updated
	//and how often aggregated statistics are written to the log
	// +optional
	StatsPeriodSeconds *int `json:"statsPeriodSeconds,omitempty"`
	//AdminUsers specifies an array of users who can act as PgBouncer administrators
	// +optional
	AdminUsers []string `json:"adminUsers,omitempty"`
	//AuthType specifies how to authenticate users. Default: md5 (md5+plain text)
	// +optional
	AuthType string `json:"authType,omitempty"`
	//AuthUser looks up any user not specified in auth_file from pg_shadow. Default: not set.
	// +optional
	AuthUser string `json:"authUser,omitempty"`
	//IgnoreStartupParameters specifies comma-separated startup parameters that
	//pgbouncer knows are handled by admin and it can ignore them
	// +optional
	IgnoreStartupParameters string `json:"ignoreStartupParameters,omitempty"`
}

type UserList struct {
	//SecretName points to a secret that holds a file containing list of users
	SecretName string `json:"secretName"`
	//SecretNamespace specifies the namespace of specified secret.
	//By default, uses the same namespace as pgbouncer if left empty, not default namespace.
	// +optional
	SecretNamespace string `json:"secretNamespace,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PgBouncer CRD objects
	Items []PgBouncer `json:"items,omitempty"`
}

type PgBouncerStatus struct {
	//Phase specifies the current state of PgBouncer server
	Phase DatabasePhase `json:"phase,omitempty"`
	//Reason is used to explain phases of interest of the server.
	Reason string `json:"reason,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
}
