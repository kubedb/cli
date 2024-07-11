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
	"k8s.io/apimachinery/pkg/util/intstr"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeElasticsearch     = "es"
	ResourceKindElasticsearch     = "Elasticsearch"
	ResourceSingularElasticsearch = "elasticsearch"
	ResourcePluralElasticsearch   = "elasticsearches"
)

// Elasticsearch defines a Elasticsearch database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=elasticsearches,singular=elasticsearch,shortName=es,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Elasticsearch struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ElasticsearchSpec   `json:"spec,omitempty"`
	Status            ElasticsearchStatus `json:"status,omitempty"`
}

type ElasticsearchSpec struct {
	// AutoOps contains configuration of automatic ops-request-recommendation generation
	// +optional
	AutoOps AutoOpsSpec `json:"autoOps,omitempty"`

	// Version of Elasticsearch to be deployed.
	Version string `json:"version"`

	// Number of instances to deploy for a Elasticsearch database.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Elasticsearch topology for node specification
	// +optional
	Topology *ElasticsearchClusterTopology `json:"topology,omitempty"`

	// To enable ssl for http layer
	EnableSSL bool `json:"enableSSL,omitempty"`

	// disable security of authPlugin (ie, xpack or searchguard). It disables authentication security of user.
	// If unset, default is false
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSecret is an optional field to provide custom configuration file for database.
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`

	// SecureConfigSecret is an optional field to provide secure settings for database.
	//	- Ref: https://www.elastic.co/guide/en/elasticsearch/reference/7.14/secure-settings.html
	// Secure settings are store at "ES_CONFIG_DIR/elasticsearch.keystore" file (contents are encoded with password),
	// once the keystore created.
	// Expects a k8s secret name with data format:
	//	data:
	//		key: value
	//		password: KEYSTORE_PASSWORD
	//		s3.client.default.access_key: ACCESS_KEY
	//		s3.client.default.secret_key: SECRET_KEY
	// +optional
	SecureConfigSecret *core.LocalObjectReference `json:"secureConfigSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// An eviction is allowed if at most "maxUnavailable" pods selected by
	// "selector" are unavailable after the eviction, i.e. even in absence of
	// the evicted pod. For example, one can prevent all voluntary evictions
	// by specifying 0. This is a mutually exclusive setting with "minAvailable".
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`

	// TLS contains tls configurations
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// InternalUsers contains internal user configurations.
	// Expected Input format:
	// internalUsers:
	//   <username1>:
	//		...
	//   <username2>:
	//		...
	// +optional
	InternalUsers map[string]ElasticsearchUserSpec `json:"internalUsers,omitempty"`

	// RolesMapping contains roles mapping configurations.
	// Expected Input format:
	// rolesMapping:
	//   <role1>:
	//		...
	//   <role2>:
	//		...
	// +optional
	RolesMapping map[string]ElasticsearchRoleMapSpec `json:"rolesMapping,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// KernelSettings contains the additional kernel settings.
	// +optional
	KernelSettings *KernelSettings `json:"kernelSettings,omitempty"`

	// HeapSizePercentage specifies both the initial heap allocation (xms) percentage and the maximum heap allocation (xmx) percentage.
	// Elasticsearch bootstrap fails, if -Xms and -Xmx are not equal.
	// Error: initial heap size [X] not equal to maximum heap size [Y]; this can cause resize pauses.
	// It will be applied to all nodes. If the node level `heapSizePercentage` is specified,  this global value will be overwritten.
	// It defaults to 50% of memory limit.
	// +optional
	// +kubebuilder:default=50
	HeapSizePercentage *int32 `json:"heapSizePercentage,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 1}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`
}

type ElasticsearchClusterTopology struct {
	Master       ElasticsearchNode  `json:"master"`
	Ingest       ElasticsearchNode  `json:"ingest"`
	Data         *ElasticsearchNode `json:"data,omitempty"`
	DataContent  *ElasticsearchNode `json:"dataContent,omitempty"`
	DataHot      *ElasticsearchNode `json:"dataHot,omitempty"`
	DataWarm     *ElasticsearchNode `json:"dataWarm,omitempty"`
	DataCold     *ElasticsearchNode `json:"dataCold,omitempty"`
	DataFrozen   *ElasticsearchNode `json:"dataFrozen,omitempty"`
	ML           *ElasticsearchNode `json:"ml,omitempty"`
	Transform    *ElasticsearchNode `json:"transform,omitempty"`
	Coordinating *ElasticsearchNode `json:"coordinating,omitempty"`
}

type ElasticsearchNode struct {
	// Replicas represents number of replica for this specific type of node
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// +optional
	Suffix string `json:"suffix,omitempty"`
	// HeapSizePercentage specifies both the initial heap allocation (-Xms) percentage and the maximum heap allocation (-Xmx) percentage.
	// Node level values have higher precedence than global values.
	// +optional
	HeapSizePercentage *int32 `json:"heapSizePercentage,omitempty"`
	// Storage to specify how storage shall be used.
	// +optional
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`
	// An eviction is allowed if at most "maxUnavailable" pods selected by
	// "selector" are unavailable after the eviction, i.e. even in absence of
	// the evicted pod. For example, one can prevent all voluntary evictions
	// by specifying 0. This is a mutually exclusive setting with "minAvailable".
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
}

// +kubebuilder:validation:Enum=ca;transport;http;admin;client;archiver;metrics-exporter
type ElasticsearchCertificateAlias string

const (
	ElasticsearchCACert              ElasticsearchCertificateAlias = "ca"
	ElasticsearchTransportCert       ElasticsearchCertificateAlias = "transport"
	ElasticsearchHTTPCert            ElasticsearchCertificateAlias = "http"
	ElasticsearchAdminCert           ElasticsearchCertificateAlias = "admin"
	ElasticsearchClientCert          ElasticsearchCertificateAlias = "client"
	ElasticsearchArchiverCert        ElasticsearchCertificateAlias = "archiver"
	ElasticsearchMetricsExporterCert ElasticsearchCertificateAlias = "metrics-exporter"
)

type ElasticsearchInternalUser string

const (
	ElasticsearchInternalUserElastic              ElasticsearchInternalUser = "elastic"
	ElasticsearchInternalUserAdmin                ElasticsearchInternalUser = "admin"
	ElasticsearchInternalUserKibanaserver         ElasticsearchInternalUser = "kibanaserver"
	ElasticsearchInternalUserKibanaSystem         ElasticsearchInternalUser = "kibana_system"
	ElasticsearchInternalUserLogstashSystem       ElasticsearchInternalUser = "logstash_system"
	ElasticsearchInternalUserBeatsSystem          ElasticsearchInternalUser = "beats_system"
	ElasticsearchInternalUserApmSystem            ElasticsearchInternalUser = "apm_system"
	ElasticsearchInternalUserRemoteMonitoringUser ElasticsearchInternalUser = "remote_monitoring_user"
	ElasticsearchInternalUserKibanaro             ElasticsearchInternalUser = "kibanaro"
	ElasticsearchInternalUserLogstash             ElasticsearchInternalUser = "logstash"
	ElasticsearchInternalUserReadall              ElasticsearchInternalUser = "readall"
	ElasticsearchInternalUserSnapshotrestore      ElasticsearchInternalUser = "snapshotrestore"
	ElasticsearchInternalUserMetricsExporter      ElasticsearchInternalUser = "metrics_exporter"
)

// ElasticsearchUserSpec specifies the security plugin internal user structure.
// Both 'json' and 'yaml' tags are used in structure metadata.
// The `json` tags (camel case) are used while taking input from users.
// The `yaml` tags (snake case) are used by the operator to generate internal_users.yml file.
// For Elastic-Stack built-in users, there is no yaml files, instead the operator is responsible for
// creating/syncing the users. For the fields that are only used by operator,
// the metadata yaml tag is kept empty ("-") so that they do not interrupt in other distributions YAML generation.
type ElasticsearchUserSpec struct {
	// Specifies the hash of the password.
	// +optional
	Hash string `json:"-" yaml:"hash,omitempty"`

	// Specifies The full name of the user
	// Only applicable for xpack authplugin
	FullName string `json:"full_name,omitempty" yaml:"-"`

	// Specifies Arbitrary metadata that you want to associate with the user
	// Only applicable for xpack authplugin
	Metadata map[string]string `json:"metadata,omitempty" yaml:"-"`

	// Specifies the email of the user.
	// Only applicable for xpack authplugin
	Email string `json:"email,omitempty" yaml:"-"`

	// A set of roles the user has. The roles determine the user’s access permissions.
	// To create a user without any roles, specify an empty list: []
	// Only applicable for xpack authplugin
	Roles []string `json:"roles,omitempty" yaml:"-"`

	// Specifies the k8s secret name that holds the user credentials.
	// Default to "<resource-name>-<username>-cred".
	// +optional
	SecretName string `json:"secretName,omitempty" yaml:"-"`

	// Specifies the reserved status.
	// Resources that have this set to true can’t be changed using the REST API or Kibana.
	// Default to "false".
	// +optional
	Reserved bool `json:"reserved,omitempty" yaml:"reserved,omitempty"`

	// Specifies the hidden status.
	// Resources that have this set to true are not returned by the REST API
	// and not visible in Kibana.
	// Default to "false".
	// +optional
	Hidden bool `json:"hidden,omitempty" yaml:"hidden,omitempty"`

	// Specifies a list of backend roles assigned to this user.
	// Backend roles can come from the internal user database,
	// LDAP groups, JSON web token claims or SAML assertions.
	// +optional
	BackendRoles []string `json:"backendRoles,omitempty" yaml:"backend_roles,omitempty"`

	// Specifies a list of searchguard security plugin roles assigned to this user.
	// +optional
	SearchGuardRoles []string `json:"searchGuardRoles,omitempty" yaml:"search_guard_roles,omitempty"`

	// Specifies a list of opendistro security plugin roles assigned to this user.
	// +optional
	OpendistroSecurityRoles []string `json:"opendistroSecurityRoles,omitempty" yaml:"opendistro_security_roles,omitempty"`

	// Specifies one or more custom attributes,
	// which can be used in index names and DLS queries.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty" yaml:"attributes,omitempty"`

	// Specifies the description of the user
	// +optional
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// Specifies the role mapping structure.
// Both 'json' and 'yaml' tags are used in structure metadata.
// The `json` tags (camel case) are used while taking input from users.
// The `yaml` tags (snake case) are used by the operator to generate roles_mapping.yml file.
type ElasticsearchRoleMapSpec struct {
	// Specifies the reserved status.
	// Resources that have this set to true can’t be changed using the REST API or Kibana.
	// Default to "false".
	// +optional
	Reserved bool `json:"reserved,omitempty" yaml:"reserved,omitempty"`

	// Specifies the hidden status.
	// Resources that have this set to true are not returned by the REST API
	// and not visible in Kibana.
	// Default to "false".
	// +optional
	Hidden bool `json:"hidden,omitempty" yaml:"hidden,omitempty"`

	// Specifies a list of backend roles assigned to this role.
	// Backend roles can come from the internal user database,
	// LDAP groups, JSON web token claims or SAML assertions.
	// +optional
	BackendRoles []string `json:"backendRoles,omitempty" yaml:"backend_roles,omitempty"`

	// Specifies a list of hosts assigned to this role.
	// +optional
	Hosts []string `json:"hosts,omitempty" yaml:"hosts,omitempty"`

	// Specifies a list of users assigned to this role.
	// +optional
	Users []string `json:"users,omitempty" yaml:"users,omitempty"`

	// Specifies a list of backend roles (migrated from ES-version6) assigned to this role.
	AndBackendRoles []string `json:"andBackendRoles,omitempty" yaml:"and_backend_roles,omitempty"`
}

type ElasticsearchStatus struct {
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

type ElasticsearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Elasticsearch CRD objects
	Items []Elasticsearch `json:"items,omitempty"`
}

type ElasticsearchNodeRoleType string

const (
	ElasticsearchNodeRoleTypeCombined            ElasticsearchNodeRoleType = "combined"
	ElasticsearchNodeRoleTypeMaster              ElasticsearchNodeRoleType = "master"
	ElasticsearchNodeRoleTypeData                ElasticsearchNodeRoleType = "data"
	ElasticsearchNodeRoleTypeDataContent         ElasticsearchNodeRoleType = "data-content"
	ElasticsearchNodeRoleTypeDataHot             ElasticsearchNodeRoleType = "data-hot"
	ElasticsearchNodeRoleTypeDataWarm            ElasticsearchNodeRoleType = "data-warm"
	ElasticsearchNodeRoleTypeDataCold            ElasticsearchNodeRoleType = "data-cold"
	ElasticsearchNodeRoleTypeDataFrozen          ElasticsearchNodeRoleType = "data-frozen"
	ElasticsearchNodeRoleTypeIngest              ElasticsearchNodeRoleType = "ingest"
	ElasticsearchNodeRoleTypeML                  ElasticsearchNodeRoleType = "ml"
	ElasticsearchNodeRoleTypeRemoteClusterClient ElasticsearchNodeRoleType = "remote-cluster-client"
	ElasticsearchNodeRoleTypeTransform           ElasticsearchNodeRoleType = "transform"
	ElasticsearchNodeRoleTypeCoordinating        ElasticsearchNodeRoleType = "coordinating"
)
