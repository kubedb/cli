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
	ResourceCodeMySQL     = "my"
	ResourceKindMySQL     = "MySQL"
	ResourceSingularMySQL = "mysql"
	ResourcePluralMySQL   = "mysqls"
)

type MySQLClusterMode string

const (
	MySQLClusterModeGroup MySQLClusterMode = "GroupReplication"
)

type MySQLGroupMode string

const (
	MySQLGroupModeSinglePrimary MySQLGroupMode = "Single-Primary"
	MySQLGroupModeMultiPrimary  MySQLGroupMode = "Multi-Primary"
)

// Mysql defines a Mysql database.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// +kubebuilder:object:root=true
// +kubebuilder:skipversion
type MySQL struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MySQLSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            MySQLStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type MySQLSpec struct {
	// Version of MySQL to be deployed.
	Version types.StrYo `json:"version" protobuf:"bytes,1,opt,name=version,casttype=gomodules.xyz/encoding/json/types.StrYo"`

	// Number of instances to deploy for a MySQL database. In case of MySQL group
	// replication, max allowed value is 9 (default 3).
	// (see ref: https://dev.mysql.com/doc/refman/5.7/en/group-replication-frequently-asked-questions.html)
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// MySQL cluster topology
	Topology *MySQLClusterTopology `json:"topology,omitempty" protobuf:"bytes,3,opt,name=topology"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty" protobuf:"bytes,4,opt,name=storageType,casttype=StorageType"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty" protobuf:"bytes,5,opt,name=storage"`

	// Database authentication secret
	DatabaseSecret *core.SecretVolumeSource `json:"databaseSecret,omitempty" protobuf:"bytes,6,opt,name=databaseSecret"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty" protobuf:"bytes,7,opt,name=init"`

	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty" protobuf:"bytes,8,opt,name=backupSchedule"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,9,opt,name=monitor"`

	// ConfigSource is an optional field to provide custom configuration file for database (i.e custom-mysql.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty" protobuf:"bytes,10,opt,name=configSource"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,11,opt,name=podTemplate"`

	// ServiceTemplate is an optional configuration for service used to expose database
	// +optional
	ServiceTemplate ofst.ServiceTemplateSpec `json:"serviceTemplate,omitempty" protobuf:"bytes,12,opt,name=serviceTemplate"`

	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	UpdateStrategy apps.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty" protobuf:"bytes,13,opt,name=updateStrategy"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty" protobuf:"bytes,14,opt,name=terminationPolicy,casttype=TerminationPolicy"`
}

type MySQLClusterTopology struct {
	// If set to -
	// "GroupReplication", GroupSpec is required and MySQL servers will start  a replication group
	Mode *MySQLClusterMode `json:"mode,omitempty" protobuf:"bytes,1,opt,name=mode,casttype=MySQLClusterMode"`

	// Group replication info for MySQL
	Group *MySQLGroupSpec `json:"group,omitempty" protobuf:"bytes,2,opt,name=group"`
}

type MySQLGroupSpec struct {
	// TODO: "Multi-Primary" needs to be implemented
	// Group Replication can be deployed in either "Single-Primary" or "Multi-Primary" mode
	Mode *MySQLGroupMode `json:"mode,omitempty" protobuf:"bytes,1,opt,name=mode,casttype=MySQLGroupMode"`

	// Group name is a version 4 UUID
	// ref: https://dev.mysql.com/doc/refman/5.7/en/group-replication-options.html#sysvar_group_replication_group_name
	Name string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`

	// On a replication master and each replication slave, the --server-id
	// option must be specified to establish a unique replication ID in the
	// range from 1 to 2^32 − 1. “Unique”, means that each ID must be different
	// from every other ID in use by any other replication master or slave.
	// ref: https://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_server_id
	//
	// So, BaseServerID is needed to calculate a unique server_id for each member.
	BaseServerID *uint `json:"baseServerID,omitempty" protobuf:"varint,3,opt,name=baseServerID"`
}

type MySQLStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=DatabasePhase"`
	Reason string        `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty" protobuf:"bytes,3,opt,name=observedGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of MySQL TPR objects
	Items []MySQL `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
