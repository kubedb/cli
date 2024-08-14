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

//go:generate go-enum --mustparse --names --values
package v1alpha1

import (
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeDruidOpsRequest     = "drops"
	ResourceKindDruidOpsRequest     = "DruidOpsRequest"
	ResourceSingularDruidOpsRequest = "druidopsrequest"
	ResourcePluralDruidOpsRequest   = "druidopsrequests"
)

// DruidDBOpsRequest defines a Druid DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=druidopsrequests,singular=druidopsrequest,shortName=drops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type DruidOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DruidOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus    `json:"status,omitempty"`
}

// DruidOpsRequestSpec is the spec for DruidOpsRequest
type DruidOpsRequestSpec struct {
	// Specifies the Druid reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type DruidOpsRequestType `json:"type"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *DruidVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *DruidVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=VerticalScaling;VolumeExpansion;Restart
// ENUM(VerticalScaling, VolumeExpansion, Restart)
type DruidOpsRequestType string

// DruidVerticalScalingSpec contains the vertical scaling information of a Druid cluster
type DruidVerticalScalingSpec struct {
	Coordinators   *PodResources `json:"coordinators,omitempty"`
	Overlords      *PodResources `json:"overlords,omitempty"`
	MiddleManagers *PodResources `json:"middleManagers,omitempty"`
	Historicals    *PodResources `json:"historicals,omitempty"`
	Brokers        *PodResources `json:"brokers,omitempty"`
	Routers        *PodResources `json:"routers,omitempty"`
}

// DruidVolumeExpansionSpec is the spec for Druid volume expansion
type DruidVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for middleManagers nodes
	MiddleManagers *resource.Quantity `json:"middleManagers,omitempty"`
	// volume specification for historicals nodes
	Historicals *resource.Quantity `json:"historicals,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DruidOpsRequestList is a list of DruidOpsRequests
type DruidOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of DruidOpsRequest CRD objects
	Items []DruidOpsRequest `json:"items,omitempty"`
}
