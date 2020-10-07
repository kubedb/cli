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
	ResourceCodePgBouncerOpsRequest     = "pbops"
	ResourceKindPgBouncerOpsRequest     = "PgBouncerOpsRequest"
	ResourceSingularPgBouncerOpsRequest = "pgbounceropsrequest"
	ResourcePluralPgBouncerOpsRequest   = "pgbounceropsrequests"
)

// PgBouncerOpsRequest defines a PgBouncer DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pgbounceropsrequests,singular=pgbounceropsrequest,shortName=pbops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PgBouncerOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              PgBouncerOpsRequestSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            PgBouncerOpsRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// PgBouncerOpsRequestSpec is the spec for PgBouncerOpsRequest
type PgBouncerOpsRequestSpec struct {
	// Specifies the PgBouncer reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef" protobuf:"bytes,1,opt,name=databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=OpsRequestType"`
	// Specifies information necessary for upgrading PgBouncer
	Upgrade *PgBouncerUpgradeSpec `json:"upgrade,omitempty" protobuf:"bytes,3,opt,name=upgrade"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *PgBouncerHorizontalScalingSpec `json:"horizontalScaling,omitempty" protobuf:"bytes,4,opt,name=horizontalScaling"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *PgBouncerVerticalScalingSpec `json:"verticalScaling,omitempty" protobuf:"bytes,5,opt,name=verticalScaling"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *PgBouncerVolumeExpansionSpec `json:"volumeExpansion,omitempty" protobuf:"bytes,6,opt,name=volumeExpansion"`
	// Specifies information necessary for custom configuration of PgBouncer
	Configuration *PgBouncerCustomConfigurationSpec `json:"configuration,omitempty" protobuf:"bytes,7,opt,name=configuration"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty" protobuf:"bytes,8,opt,name=tls"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty" protobuf:"bytes,9,opt,name=restart"`
}

// PgBouncerReplicaReadinessCriteria is the criteria for checking readiness of a PgBouncer pod
// after updating, horizontal scaling etc.
type PgBouncerReplicaReadinessCriteria struct {
}

type PgBouncerUpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                             `json:"targetVersion,omitempty" protobuf:"bytes,1,opt,name=targetVersion"`
	ReadinessCriteria *PgBouncerReplicaReadinessCriteria `json:"readinessCriteria,omitempty" protobuf:"bytes,2,opt,name=readinessCriteria"`
}

// HorizontalScaling is the spec for PgBouncer horizontal scaling
type PgBouncerHorizontalScalingSpec struct {
}

// PgBouncerVerticalScalingSpec is the spec for PgBouncer vertical scaling
type PgBouncerVerticalScalingSpec struct {
	ReadinessCriteria *PgBouncerReplicaReadinessCriteria `json:"readinessCriteria,omitempty" protobuf:"bytes,1,opt,name=readinessCriteria"`
}

// PgBouncerVolumeExpansionSpec is the spec for PgBouncer volume expansion
type PgBouncerVolumeExpansionSpec struct {
}

type PgBouncerCustomConfigurationSpec struct {
}

type PgBouncerCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty" protobuf:"bytes,1,opt,name=configMap"`
	Data      map[string]string          `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`
	Remove    bool                       `json:"remove,omitempty" protobuf:"varint,3,opt,name=remove"`
}

// PgBouncerOpsRequestStatus is the status for PgBouncerOpsRequest
type PgBouncerOpsRequestStatus struct {
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

// PgBouncerOpsRequestList is a list of PgBouncerOpsRequests
type PgBouncerOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of PgBouncerOpsRequest CRD objects
	Items []PgBouncerOpsRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
