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
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeRedisOpsRequest     = "rdops"
	ResourceKindRedisOpsRequest     = "RedisOpsRequest"
	ResourceSingularRedisOpsRequest = "redisopsrequest"
	ResourcePluralRedisOpsRequest   = "redisopsrequests"
)

// RedisOpsRequest defines a Redis DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=redisopsrequests,singular=redisopsrequest,shortName=rdops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type RedisOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RedisOpsRequestSpec   `json:"spec,omitempty"`
	Status            RedisOpsRequestStatus `json:"status,omitempty"`
}

// RedisOpsRequestSpec is the spec for RedisOpsRequest
type RedisOpsRequestSpec struct {
	// Specifies the Redis reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type"`
	// Specifies information necessary for upgrading Redis
	Upgrade *RedisUpgradeSpec `json:"upgrade,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *RedisHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *RedisVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *RedisVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of Redis
	Configuration *RedisCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
}

// RedisReplicaReadinessCriteria is the criteria for checking readiness of a Redis pod
// after updating, horizontal scaling etc.
type RedisReplicaReadinessCriteria struct{}

type RedisUpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                         `json:"targetVersion,omitempty"`
	ReadinessCriteria *RedisReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

type RedisHorizontalScalingSpec struct {
	// Number of Masters in the cluster
	Master *int32 `json:"master,omitempty"`
	// specifies the number of replica for the master
	Replicas *int32 `json:"replicas,omitempty"`
}

// RedisVerticalScalingSpec is the spec for Redis vertical scaling
type RedisVerticalScalingSpec struct {
	Redis       *core.ResourceRequirements `json:"redis,omitempty"`
	Exporter    *core.ResourceRequirements `json:"exporter,omitempty"`
	Coordinator *core.ResourceRequirements `json:"coordinator,omitempty"`
}

// RedisVolumeExpansionSpec is the spec for Redis volume expansion
type RedisVolumeExpansionSpec struct {
	Redis *resource.Quantity `json:"redis,omitempty"`
}

type RedisCustomConfigurationSpec struct {
	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate        ofst.PodTemplateSpec       `json:"podTemplate,omitempty"`
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty"`
	InlineConfig       string                     `json:"inlineConfig,omitempty"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty"`
}

// RedisOpsRequestStatus is the status for RedisOpsRequest
type RedisOpsRequestStatus struct {
	Phase OpsRequestPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisOpsRequestList is a list of RedisOpsRequests
type RedisOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of RedisOpsRequest CRD objects
	Items []RedisOpsRequest `json:"items,omitempty"`
}
