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
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeRedisSentinelOpsRequest     = "rdsops"
	ResourceKindRedisSentinelOpsRequest     = "RedisSentinelOpsRequest"
	ResourceSingularRedisSentinelOpsRequest = "redissentinelopsrequest"
	ResourcePluralRedisSentinelOpsRequest   = "redissentinelopsrequests"
)

// RedisSentinelOpsRequest defines a RedisSentinel DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=redissentinelopsrequests,singular=redissentinelopsrequest,shortName=rdsops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type RedisSentinelOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RedisSentinelOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus            `json:"status,omitempty"`
}

// RedisSentinelOpsRequestSpec is the spec for RedisSentinelOpsRequest
type RedisSentinelOpsRequestSpec struct {
	// Specifies the RedisSentinel reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type RedisSentinelOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading RedisSentinel
	UpdateVersion *RedisSentinelUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *RedisSentinelHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *RedisSentinelVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for custom configuration of RedisSentinel
	Configuration *RedisSentinelCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;Restart;Reconfigure;ReconfigureTLS
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, Restart, Reconfigure, ReconfigureTLS)
type RedisSentinelOpsRequestType string

// RedisSentinelReplicaReadinessCriteria is the criteria for checking readiness of a RedisSentinel pod
// after updating, horizontal scaling etc.
type RedisSentinelReplicaReadinessCriteria struct{}

type RedisSentinelUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                                 `json:"targetVersion,omitempty"`
	ReadinessCriteria *RedisSentinelReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

type RedisSentinelHorizontalScalingSpec struct {
	// specifies the number of replica for the master
	Replicas *int32 `json:"replicas,omitempty"`
}

// RedisSentinelVerticalScalingSpec is the spec for RedisSentinel vertical scaling
type RedisSentinelVerticalScalingSpec struct {
	RedisSentinel *core.ResourceRequirements `json:"redissentinel,omitempty"`
	Exporter      *core.ResourceRequirements `json:"exporter,omitempty"`
	Coordinator   *core.ResourceRequirements `json:"coordinator,omitempty"`
}

// RedisSentinelVolumeExpansionSpec is the spec for RedisSentinel volume expansion
type RedisSentinelVolumeExpansionSpec struct {
	// +kubebuilder:default="Online"
	Mode          *VolumeExpansionMode `json:"mode,omitempty"`
	RedisSentinel *resource.Quantity   `json:"redissentinel,omitempty"`
}

type RedisSentinelCustomConfigurationSpec struct {
	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate        ofst.PodTemplateSpec       `json:"podTemplate,omitempty"`
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty"`
	InlineConfig       string                     `json:"inlineConfig,omitempty"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisSentinelOpsRequestList is a list of RedisSentinelOpsRequests
type RedisSentinelOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of RedisSentinelOpsRequest CRD objects
	Items []RedisSentinelOpsRequest `json:"items,omitempty"`
}
