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
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ProxySQLSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            ProxySQLStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type ProxySQLSpec struct {
	// Version of ProxySQL to be deployed.
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`

	// Number of instances to deploy for ProxySQL. Currently we support only replicas = 1.
	// TODO: If replicas > 1, proxysql will be clustered
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// Mode specifies the type of MySQL/Percona-XtraDB/MariaDB cluster for which proxysql
	// will be configured. It must be either "Galera" or "GroupReplication"
	Mode *LoadBalanceMode `json:"mode,omitempty" protobuf:"bytes,3,opt,name=mode,casttype=LoadBalanceMode"`

	// Backend specifies the information about backend MySQL/Percona-XtraDB/MariaDB servers
	Backend *ProxySQLBackendSpec `json:"backend,omitempty" protobuf:"bytes,4,opt,name=backend"`

	// ProxySQL secret containing username and password for root user and proxysql user
	AuthSecret *core.LocalObjectReference `json:"authSecret,omitempty" protobuf:"bytes,5,opt,name=authSecret"`

	// Monitor is used monitor proxysql instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,6,opt,name=monitor"`

	// ConfigSecret is an optional field to provide custom configuration file for proxysql (i.e custom-proxysql.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty" protobuf:"bytes,7,opt,name=configSecret"`

	// PodTemplate is an optional configuration for pods used to expose proxysql
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,8,opt,name=podTemplate"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty" protobuf:"bytes,9,rep,name=serviceTemplates"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty" protobuf:"bytes,10,opt,name=tls"`
}

// +kubebuilder:validation:Enum=server;archiver;metrics-exporter
type ProxySQLCertificateAlias string

const (
	ProxySQLServerCert          ProxySQLCertificateAlias = "server"
	ProxySQLArchiverCert        ProxySQLCertificateAlias = "archiver"
	ProxySQLMetricsExporterCert ProxySQLCertificateAlias = "metrics-exporter"
)

type ProxySQLBackendSpec struct {
	// Ref lets one to locate the typed referenced object
	// (in our case, it is the MySQL/Percona-XtraDB/ProxySQL object)
	// inside the same namespace.
	Ref *core.TypedLocalObjectReference `json:"ref,omitempty" protobuf:"bytes,7,opt,name=ref"`

	// Number of backend servers.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,8,opt,name=replicas"`
}

type ProxySQLStatus struct {
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

type ProxySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of ProxySQL TPR objects
	Items []ProxySQL `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
