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
	ResourceCodeEtcd     = "etc"
	ResourceKindEtcd     = "Etcd"
	ResourceSingularEtcd = "etcd"
	ResourcePluralEtcd   = "etcds"
)

// Etcd defines a Etcd database.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// +kubebuilder:object:root=true
// +kubebuilder:skipversion
type Etcd struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              EtcdSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            EtcdStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type EtcdSpec struct {
	// Version of Etcd to be deployed.
	Version types.StrYo `json:"version" protobuf:"bytes,1,opt,name=version,casttype=gomodules.xyz/encoding/json/types.StrYo"`

	// Number of instances to deploy for a Etcd database.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty" protobuf:"bytes,3,opt,name=storageType,casttype=StorageType"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty" protobuf:"bytes,4,opt,name=storage"`

	// Database authentication secret
	DatabaseSecret *core.SecretVolumeSource `json:"databaseSecret,omitempty" protobuf:"bytes,5,opt,name=databaseSecret"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty" protobuf:"bytes,6,opt,name=init"`

	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty" protobuf:"bytes,7,opt,name=backupSchedule"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,8,opt,name=monitor"`

	// etcd cluster TLS configuration
	TLS *TLSPolicy `json:"tls,omitempty" protobuf:"bytes,9,opt,name=tls"`

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

type TLSPolicy struct {
	Member         *MemberSecret `json:"member,omitempty" protobuf:"bytes,1,opt,name=member"`
	OperatorSecret string        `json:"operatorSecret,omitempty" protobuf:"bytes,2,opt,name=operatorSecret"`
}

type MemberSecret struct {
	// PeerSecret is the secret containing TLS certs used by each etcd member pod
	// for the communication between etcd peers.
	PeerSecret string `json:"peerSecret,omitempty" protobuf:"bytes,1,opt,name=peerSecret"`
	// ServerSecret is the secret containing TLS certs used by each etcd member pod
	// for the communication between etcd server and its clients.
	ServerSecret string `json:"serverSecret,omitempty" protobuf:"bytes,2,opt,name=serverSecret"`
}

type EtcdStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=DatabasePhase"`
	Reason string        `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty" protobuf:"bytes,3,opt,name=observedGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EtcdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Etcd TPR objects
	Items []Etcd `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
