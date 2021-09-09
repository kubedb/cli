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
	"gomodules.xyz/encoding/json/types"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1alpha1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodePerconaXtraDB     = "px"
	ResourceKindPerconaXtraDB     = "PerconaXtraDB"
	ResourceSingularPerconaXtraDB = "perconaxtradb"
	ResourcePluralPerconaXtraDB   = "perconaxtradbs"
)

// PerconaXtraDB defines a percona variation of Mysql database.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// +kubebuilder:object:root=true
// +kubebuilder:skipversion
type PerconaXtraDB struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              PerconaXtraDBSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            PerconaXtraDBStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type PerconaXtraDBSpec struct {
	// Version of PerconaXtraDB to be deployed.
	Version types.StrYo `json:"version" protobuf:"bytes,1,opt,name=version,casttype=gomodules.xyz/encoding/json/types.StrYo"`

	// Number of instances to deploy for PerconaXtraDB
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// PXC is the cluster specification for PerconaXtraDB Cluster
	PXC *PXCSpec `json:"pxc,omitempty" protobuf:"bytes,3,opt,name=pxc"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty" protobuf:"bytes,4,opt,name=storageType,casttype=StorageType"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty" protobuf:"bytes,5,opt,name=storage"`

	// Database authentication secret
	DatabaseSecret *core.SecretVolumeSource `json:"databaseSecret,omitempty" protobuf:"bytes,6,opt,name=databaseSecret"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty" protobuf:"bytes,7,opt,name=init"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,8,opt,name=monitor"`

	// ConfigSource is an optional field to provide custom configuration file for database (i.e custom-mysql.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty" protobuf:"bytes,9,opt,name=configSource"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,10,opt,name=podTemplate"`

	// ServiceTemplate is an optional configuration for service used to expose database
	// +optional
	ServiceTemplate ofst.ServiceTemplateSpec `json:"serviceTemplate,omitempty" protobuf:"bytes,11,opt,name=serviceTemplate"`

	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	UpdateStrategy apps.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty" protobuf:"bytes,12,opt,name=updateStrategy"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty" protobuf:"bytes,13,opt,name=terminationPolicy,casttype=TerminationPolicy"`
}

type PXCSpec struct {
	// Name of the cluster and should be identical on all nodes.
	ClusterName string `json:"clusterName,omitempty" protobuf:"bytes,1,opt,name=clusterName"`

	// Proxysql configuration
	Proxysql ProxysqlSpec `json:"proxysql,omitempty" protobuf:"bytes,2,opt,name=proxysql"`
}

type ProxysqlSpec struct {
	// Number of Proxysql nodes. Currently we support only replicas = 1.
	// TODO: If replicas > 1, proxysql will be clustered
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`

	// PodTemplate is an optional configuration for pods used to expose proxysql
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,2,opt,name=podTemplate"`
}

type PerconaXtraDBStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=DatabasePhase"`
	Reason string        `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty" protobuf:"bytes,3,opt,name=observedGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PerconaXtraDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of PerconaXtraDB TPR objects
	Items []PerconaXtraDB `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
