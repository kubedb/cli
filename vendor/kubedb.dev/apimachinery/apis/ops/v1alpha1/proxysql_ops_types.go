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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindProxySQLOpsRequest     = "ProxySQLOpsRequest"
	ResourceSingularProxySQLOpsRequest = "proxysqlopsrequest"
	ResourcePluralProxySQLOpsRequest   = "proxysqlopsrequests"
)

// ProxySQLOpsRequest defines a ProxySQL load-balancer DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=proxysqlopsrequests,singular=proxysqlopsrequest,shortName=proxyops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ProxySQLOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ProxySQLOpsRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            ProxySQLOpsRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ProxySQLOpsRequestSpec is the spec for ProxySQLOpsRequest
type ProxySQLOpsRequestSpec struct {
	// Specifies the ProxySQL reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=OpsRequestType"`
	// Specifies information necessary for upgrading ProxySQL
	Upgrade *ProxySQLUpgradeSpec `json:"upgrade,omitempty" protobuf:"bytes,3,opt,name=upgrade"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *ProxySQLHorizontalScalingSpec `json:"horizontalScaling,omitempty" protobuf:"bytes,4,opt,name=horizontalScaling"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *ProxySQLVerticalScalingSpec `json:"verticalScaling,omitempty" protobuf:"bytes,5,opt,name=verticalScaling"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *ProxySQLVolumeExpansionSpec `json:"volumeExpansion,omitempty" protobuf:"bytes,6,opt,name=volumeExpansion"`
	// Specifies information necessary for custom configuration of ProxySQL
	Configuration *ProxySQLCustomConfigurationSpec `json:"configuration,omitempty" protobuf:"bytes,7,opt,name=configuration"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty" protobuf:"bytes,8,opt,name=tls"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty" protobuf:"bytes,9,opt,name=restart"`
}

// ProxySQLReplicaReadinessCriteria is the criteria for checking readiness of a ProxySQL pod
// after updating, horizontal scaling etc.
type ProxySQLReplicaReadinessCriteria struct {
}

type ProxySQLUpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                            `json:"targetVersion,omitempty" protobuf:"bytes,1,opt,name=targetVersion"`
	ReadinessCriteria *ProxySQLReplicaReadinessCriteria `json:"readinessCriteria,omitempty" protobuf:"bytes,2,opt,name=readinessCriteria"`
}

// HorizontalScaling is the spec for ProxySQL horizontal scaling
type ProxySQLHorizontalScalingSpec struct {
}

// ProxySQLVerticalScalingSpec is the spec for ProxySQL vertical scaling
type ProxySQLVerticalScalingSpec struct {
	ReadinessCriteria *ProxySQLReplicaReadinessCriteria `json:"readinessCriteria,omitempty" protobuf:"bytes,1,opt,name=readinessCriteria"`
}

// ProxySQLVolumeExpansionSpec is the spec for ProxySQL volume expansion
type ProxySQLVolumeExpansionSpec struct {
}

type ProxySQLCustomConfigurationSpec struct {
}

type ProxySQLCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty" protobuf:"bytes,1,opt,name=configMap"`
	Data      map[string]string          `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`
	Remove    bool                       `json:"remove,omitempty" protobuf:"varint,3,opt,name=remove"`
}

// ProxySQLOpsRequestStatus is the status for ProxySQLOpsRequest
type ProxySQLOpsRequestStatus struct {
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

// ProxySQLOpsRequestList is a list of ProxySQLOpsRequests
type ProxySQLOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of ProxySQLOpsRequest CRD objects
	Items []ProxySQLOpsRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
