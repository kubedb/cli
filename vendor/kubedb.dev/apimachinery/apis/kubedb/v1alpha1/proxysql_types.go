package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeProxySQL     = "psql"
	ResourceKindProxySQL     = "ProxySQL"
	ResourceSingularProxySQL = "proxysql"
	ResourcePluralProxySQL   = "proxysqls"
)

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
// +kubebuilder:resource:path=proxysqls,singular=proxysql,categories={datastore,kubedb,appscode,all}
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

type ProxySQLSpec struct {
	// Version of ProxySQL to be deployed.
	Version types.StrYo `json:"version"`

	// Number of instances to deploy for ProxySQL. Currently we support only replicas = 1.
	// TODO: If replicas > 1, proxysql will be clustered
	Replicas *int32 `json:"replicas,omitempty"`

	// Mode specifies the type of MySQL/Percona-XtraDB/MariaDB cluster for which proxysql
	// will be configured. It must be either "Galera" or "GroupReplication"
	Mode *LoadBalanceMode `json:"mode,omitempty"`

	// Backend specifies the information about backend MySQL/Percona-XtraDB/MariaDB servers
	Backend *ProxySQLBackendSpec `json:"backend,omitempty"`

	// ProxySQL secret containing username and password for root user and proxysql user
	ProxySQLSecret *core.SecretVolumeSource `json:"proxysqlSecret,omitempty"`

	// Monitor is used monitor proxysql instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSource is an optional field to provide custom configuration file for proxysql (i.e custom-proxysql.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose proxysql
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplate is an optional configuration for service used to expose proxysql
	// +optional
	ServiceTemplate ofst.ServiceTemplateSpec `json:"serviceTemplate,omitempty"`

	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	UpdateStrategy apps.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty"`
}

type ProxySQLBackendSpec struct {
	// Ref lets one to locate the typed referenced object
	// (in our case, it is the MySQL/Percona-XtraDB/MariaDB object)
	// inside the same namespace.
	Ref *core.TypedLocalObjectReference `json:"ref,omitempty" protobuf:"bytes,7,opt,name=ref"`

	// Number of backend servers.
	Replicas *int32 `json:"replicas,omitempty"`
}

type ProxySQLStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty"`
	Reason string        `json:"reason,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ProxySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ProxySQL TPR objects
	Items []ProxySQL `json:"items,omitempty"`
}
