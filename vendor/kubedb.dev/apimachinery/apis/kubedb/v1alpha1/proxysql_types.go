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

package v1alpha1

import (
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
// +kubebuilder:resource:path=proxysqls,singular=proxysql,categories={datastore,kubedb,appscode,all}
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
	ProxySQLSecret *core.SecretVolumeSource `json:"proxysqlSecret,omitempty" protobuf:"bytes,5,opt,name=proxysqlSecret"`

	// Monitor is used monitor proxysql instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,6,opt,name=monitor"`

	// ConfigSource is an optional field to provide custom configuration file for proxysql (i.e custom-proxysql.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty" protobuf:"bytes,7,opt,name=configSource"`

	// PodTemplate is an optional configuration for pods used to expose proxysql
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,8,opt,name=podTemplate"`

	// ServiceTemplate is an optional configuration for service used to expose proxysql
	// +optional
	ServiceTemplate ofst.ServiceTemplateSpec `json:"serviceTemplate,omitempty" protobuf:"bytes,9,opt,name=serviceTemplate"`

	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	UpdateStrategy apps.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty" protobuf:"bytes,10,opt,name=updateStrategy"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *TLSConfig `json:"tls,omitempty" protobuf:"bytes,11,opt,name=tls"`

	// Indicates that the database is paused and controller will not sync any changes made to this spec.
	// +optional
	Paused bool `json:"paused,omitempty" protobuf:"varint,12,opt,name=paused"`
}

type ProxySQLBackendSpec struct {
	// Ref lets one to locate the typed referenced object
	// (in our case, it is the MySQL/Percona-XtraDB/MariaDB object)
	// inside the same namespace.
	Ref *core.TypedLocalObjectReference `json:"ref,omitempty" protobuf:"bytes,7,opt,name=ref"`

	// Number of backend servers.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,8,opt,name=replicas"`
}

type ProxySQLStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=DatabasePhase"`
	Reason string        `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ProxySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of ProxySQL TPR objects
	Items []ProxySQL `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
