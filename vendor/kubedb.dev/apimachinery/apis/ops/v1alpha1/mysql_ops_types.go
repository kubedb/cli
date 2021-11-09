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
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceCodeMySQLOpsRequest     = "myops"
	ResourceKindMySQLOpsRequest     = "MySQLOpsRequest"
	ResourceSingularMySQLOpsRequest = "mysqlopsrequest"
	ResourcePluralMySQLOpsRequest   = "mysqlopsrequests"
)

// MySQLOpsRequest defines a MySQL DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mysqlopsrequests,singular=mysqlopsrequest,shortName=myops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MySQLOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MySQLOpsRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            MySQLOpsRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MySQLOpsRequestSpec is the spec for MySQLOpsRequest
type MySQLOpsRequestSpec struct {
	// Specifies the MySQL reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=OpsRequestType"`
	// Specifies information necessary for upgrading MySQL
	Upgrade *MySQLUpgradeSpec `json:"upgrade,omitempty" protobuf:"bytes,3,opt,name=upgrade"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *MySQLHorizontalScalingSpec `json:"horizontalScaling,omitempty" protobuf:"bytes,4,opt,name=horizontalScaling"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *MySQLVerticalScalingSpec `json:"verticalScaling,omitempty" protobuf:"bytes,5,opt,name=verticalScaling"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *MySQLVolumeExpansionSpec `json:"volumeExpansion,omitempty" protobuf:"bytes,6,opt,name=volumeExpansion"`
	// Specifies information necessary for custom configuration of MySQL
	Configuration *MySQLCustomConfigurationSpec `json:"configuration,omitempty" protobuf:"bytes,7,opt,name=configuration"`
	// Specifies information necessary for configuring TLS
	TLS *MySQLTLSSpec `json:"tls,omitempty" protobuf:"bytes,8,opt,name=tls"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty" protobuf:"bytes,9,opt,name=restart"`
}

// MySQLReplicaReadinessCriteria is the criteria for checking readiness of a MySQL pod
// after updating, horizontal scaling etc.
type MySQLReplicaReadinessCriteria struct {
}

type MySQLUpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                         `json:"targetVersion,omitempty" protobuf:"bytes,1,opt,name=targetVersion"`
	ReadinessCriteria *MySQLReplicaReadinessCriteria `json:"readinessCriteria,omitempty" protobuf:"bytes,2,opt,name=readinessCriteria"`
}

type MySQLHorizontalScalingSpec struct {
	// Number of nodes/members of the group
	Member *int32 `json:"member,omitempty" protobuf:"varint,1,opt,name=member"`
}

type MySQLVerticalScalingSpec struct {
	MySQL       *core.ResourceRequirements `json:"mysql,omitempty" protobuf:"bytes,1,opt,name=mysql"`
	Exporter    *core.ResourceRequirements `json:"exporter,omitempty" protobuf:"bytes,2,opt,name=exporter"`
	Coordinator *core.ResourceRequirements `json:"coordinator,omitempty" protobuf:"bytes,3,opt,name=coordinator"`
}

// MySQLVolumeExpansionSpec is the spec for MySQL volume expansion
type MySQLVolumeExpansionSpec struct {
	MySQL *resource.Quantity `json:"mysql,omitempty" protobuf:"bytes,1,opt,name=mysql"`
}

type MySQLCustomConfigurationSpec struct {
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty" protobuf:"bytes,2,opt,name=configSecret"`
	InlineConfig       string                     `json:"inlineConfig,omitempty" protobuf:"bytes,3,opt,name=inlineConfig"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty" protobuf:"varint,4,opt,name=removeCustomConfig"`
}

type MySQLTLSSpec struct {
	TLSSpec `json:",inline,omitempty" protobuf:"bytes,1,opt,name=tLSSpec"`

	// Indicates that the database server need to be encrypted connections(ssl)
	// +optional
	RequireSSL *bool `json:"requireSSL,omitempty" protobuf:"varint,2,opt,name=requireSSL"`
}

// MySQLOpsRequestStatus is the status for MySQLOpsRequest
type MySQLOpsRequestStatus struct {
	Phase OpsRequestPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=OpsRequestPhase"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,2,opt,name=observedGeneration"`
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty" protobuf:"bytes,3,rep,name=conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLOpsRequestList is a list of MySQLOpsRequests
type MySQLOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of MySQLOpsRequest CRD objects
	Items []MySQLOpsRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
