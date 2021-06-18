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
	ResourceCodeMariaDBOpsRequest     = "mariaops"
	ResourceKindMariaDBOpsRequest     = "MariaDBOpsRequest"
	ResourceSingularMariaDBOpsRequest = "mariadbopsrequest"
	ResourcePluralMariaDBOpsRequest   = "mariadbopsrequests"
)

// MariaDBOpsRequest defines a MariaDB DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mariadbopsrequests,singular=mariadbopsrequest,shortName=mariaops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MariaDBOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MariaDBOpsRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            MariaDBOpsRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MariaDBOpsRequestSpec is the spec for MariaDBOpsRequest
type MariaDBOpsRequestSpec struct {
	// Specifies the MariaDB reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=OpsRequestType"`
	// Specifies information necessary for upgrading MariaDB
	Upgrade *MariaDBUpgradeSpec `json:"upgrade,omitempty" protobuf:"bytes,3,opt,name=upgrade"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *MariaDBHorizontalScalingSpec `json:"horizontalScaling,omitempty" protobuf:"bytes,4,opt,name=horizontalScaling"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *MariaDBVerticalScalingSpec `json:"verticalScaling,omitempty" protobuf:"bytes,5,opt,name=verticalScaling"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *MariaDBVolumeExpansionSpec `json:"volumeExpansion,omitempty" protobuf:"bytes,6,opt,name=volumeExpansion"`
	// Specifies information necessary for custom configuration of MariaDB
	Configuration *MariaDBCustomConfigurationSpec `json:"configuration,omitempty" protobuf:"bytes,7,opt,name=configuration"`
	// Specifies information necessary for configuring TLS
	TLS *MariaDBTLSSpec `json:"tls,omitempty" protobuf:"bytes,8,opt,name=tls"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty" protobuf:"bytes,9,opt,name=restart"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty" protobuf:"bytes,10,opt,name=timeout"`
}

type MariaDBUpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty" protobuf:"bytes,1,opt,name=targetVersion"`
}

type MariaDBHorizontalScalingSpec struct {
	// Number of nodes/members of the group
	Member *int32 `json:"member,omitempty" protobuf:"varint,1,opt,name=member"`
	// specifies the weight of the current member/Node
	MemberWeight int32 `json:"memberWeight,omitempty" protobuf:"varint,2,opt,name=memberWeight"`
}

type MariaDBVerticalScalingSpec struct {
	MariaDB  *core.ResourceRequirements `json:"mariadb,omitempty" protobuf:"bytes,1,opt,name=mariadb"`
	Exporter *core.ResourceRequirements `json:"exporter,omitempty" protobuf:"bytes,2,opt,name=exporter"`
}

// MariaDBVolumeExpansionSpec is the spec for MariaDB volume expansion
type MariaDBVolumeExpansionSpec struct {
	MariaDB *resource.Quantity `json:"mariadb,omitempty" protobuf:"bytes,1,opt,name=mariadb"`
}

type MariaDBCustomConfigurationSpec struct {
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty" protobuf:"bytes,1,opt,name=configSecret"`
	InlineConfig       string                     `json:"inlineConfig,omitempty" protobuf:"bytes,2,opt,name=inlineConfig"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty" protobuf:"varint,3,opt,name=removeCustomConfig"`
}

type MariaDBCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty" protobuf:"bytes,1,opt,name=configMap"`
	Data      map[string]string          `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`
	Remove    bool                       `json:"remove,omitempty" protobuf:"varint,3,opt,name=remove"`
}

type MariaDBTLSSpec struct {
	TLSSpec `json:",inline,omitempty" protobuf:"bytes,1,opt,name=tLSSpec"`

	// Indicates that the database server need to be encrypted connections(ssl)
	// +optional
	RequireSSL *bool `json:"requireSSL,omitempty" protobuf:"varint,2,opt,name=requireSSL"`
}

// MariaDBOpsRequestStatus is the status for MariaDBOpsRequest
type MariaDBOpsRequestStatus struct {
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

// MariaDBOpsRequestList is a list of MariaDBOpsRequests
type MariaDBOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of MariaDBOpsRequest CRD objects
	Items []MariaDBOpsRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
