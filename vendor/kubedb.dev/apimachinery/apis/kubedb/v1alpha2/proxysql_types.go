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
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeProxySQL     = "prx"
	ResourceKindProxySQL     = "ProxySQL"
	ResourceSingularProxySQL = "proxysql"
	ResourcePluralProxySQL   = "proxysqls"
)

// +kubebuilder:validation:Enum=Galera;GroupReplication
type LoadBalanceMode string

const (
	LoadBalanceModeGalera           LoadBalanceMode = "Galera"
	LoadBalanceModeGroupReplication LoadBalanceMode = "GroupReplication"
)

// ProxySQL defines a percona variation of Mysql database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=proxysqls,singular=proxysql,shortName=prx,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ProxySQL struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProxySQLSpec   `json:"spec,omitempty"`
	Status            ProxySQLStatus `json:"status,omitempty"`
}

type MySQLUser struct {
	Username string `json:"username"`

	// +optional
	Active *int `json:"active,omitempty"`

	// +optional
	UseSSL int `json:"use_ssl,omitempty"`

	// +optional
	DefaultHostgroup int `json:"default_hostgroup,omitempty"`

	// +optional
	DefaultSchema string `json:"default_schema,omitempty"`

	// +optional
	SchemaLocked int `json:"schema_locked,omitempty"`

	// +optional
	TransactionPersistent *int `json:"transaction_persistent,omitempty"`

	// +optional
	FastForward int `json:"fast_forward,omitempty"`

	// +optional
	Backend *int `json:"backend,omitempty"`

	// +optional
	Frontend *int `json:"frontend,omitempty"`

	// +optional
	MaxConnections *int32 `json:"max_connections,omitempty"`

	// +optional
	Attributes string `json:"attributes,omitempty"`

	// +optional
	Comment string `json:"comment,omitempty"`
}

type ProxySQLConfiguration struct {
	// +optional
	MySQLUsers []MySQLUser `json:"mysqlUsers,omitempty"`

	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	MySQLQueryRules []*runtime.RawExtension `json:"mysqlQueryRules,omitempty"`

	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	MySQLVariables *runtime.RawExtension `json:"mysqlVariables,omitempty"`

	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	AdminVariables *runtime.RawExtension `json:"adminVariables,omitempty"`
}

type ProxySQLSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// +optional
	// SyncUsers is a boolean type and when enabled, operator fetches all users created in the backend server to the
	// ProxySQL server . Password changes are also synced in proxysql when it is enabled.
	SyncUsers bool `json:"syncUsers,omitempty"`

	// +optional
	// InitConfiguration contains information with which the proxysql will bootstrap (only 4 tables are configurable)
	InitConfiguration *ProxySQLConfiguration `json:"initConfig,omitempty"`

	// Version of ProxySQL to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for ProxySQL. Currently we support only replicas = 1.
	// TODO: If replicas > 1, proxysql will be clustered
	Replicas *int32 `json:"replicas,omitempty"`

	// Mode specifies the type of MySQL/Percona-XtraDB/MariaDB cluster for which proxysql
	// will be configured. It must be either "Galera" or "GroupReplication"
	Mode *LoadBalanceMode `json:"mode,omitempty"`

	// Backend refers to the AppBinding of the backend MySQL/MariaDB/Percona-XtraDB server
	Backend *core.LocalObjectReference `json:"backend,omitempty"`

	// ProxySQL secret containing username and password for root user and proxysql user
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// Monitor is used monitor proxysql instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for proxysql (i.e custom-proxysql.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose proxysql
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`
}

// +kubebuilder:validation:Enum=server;archiver;metrics-exporter
type ProxySQLCertificateAlias string

const (
	ProxySQLServerCert          ProxySQLCertificateAlias = "server"
	ProxySQLClientCert          ProxySQLCertificateAlias = "client"
	ProxySQLMetricsExporterCert ProxySQLCertificateAlias = "metrics-exporter"
)

type ProxySQLStatus struct {
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

type ProxySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ProxySQL TPR objects
	Items []ProxySQL `json:"items,omitempty"`
}
