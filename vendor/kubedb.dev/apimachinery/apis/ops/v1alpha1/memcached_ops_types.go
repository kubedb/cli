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
	ResourceCodeMemcachedOpsRequest     = "mcops"
	ResourceKindMemcachedOpsRequest     = "MemcachedOpsRequest"
	ResourceSingularMemcachedOpsRequest = "memcachedopsrequest"
	ResourcePluralMemcachedOpsRequest   = "memcachedopsrequests"
)

// MemcachedOpsRequest defines a Memcached DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=memcachedopsrequests,singular=memcachedopsrequest,shortName=mcops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MemcachedOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MemcachedOpsRequestSpec   `json:"spec,omitempty"`
	Status            MemcachedOpsRequestStatus `json:"status,omitempty"`
}

// MemcachedOpsRequestSpec is the spec for MemcachedOpsRequest
type MemcachedOpsRequestSpec struct {
	// Specifies the Memcached reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type"`
	// Specifies information necessary for upgrading Memcached
	Upgrade *MemcachedUpgradeSpec `json:"upgrade,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *MemcachedHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *MemcachedVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *MemcachedVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of Memcached
	Configuration *MemcachedCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
}

// MemcachedReplicaReadinessCriteria is the criteria for checking readiness of a Memcached pod
// after updating, horizontal scaling etc.
type MemcachedReplicaReadinessCriteria struct{}

type MemcachedUpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                             `json:"targetVersion,omitempty"`
	ReadinessCriteria *MemcachedReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

// HorizontalScaling is the spec for Memcached horizontal scaling
type MemcachedHorizontalScalingSpec struct{}

// MemcachedVerticalScalingSpec is the spec for Memcached vertical scaling
type MemcachedVerticalScalingSpec struct {
	ReadinessCriteria *MemcachedReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

// MemcachedVolumeExpansionSpec is the spec for Memcached volume expansion
type MemcachedVolumeExpansionSpec struct{}

type MemcachedCustomConfigurationSpec struct{}

type MemcachedCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty"`
	Data      map[string]string          `json:"data,omitempty"`
	Remove    bool                       `json:"remove,omitempty"`
}

// MemcachedOpsRequestStatus is the status for MemcachedOpsRequest
type MemcachedOpsRequestStatus struct {
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

// MemcachedOpsRequestList is a list of MemcachedOpsRequests
type MemcachedOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MemcachedOpsRequest CRD objects
	Items []MemcachedOpsRequest `json:"items,omitempty"`
}
