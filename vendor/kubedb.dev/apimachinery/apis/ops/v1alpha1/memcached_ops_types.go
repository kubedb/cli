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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Spec              MemcachedOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus        `json:"status,omitempty"`
}

// MemcachedOpsRequestSpec is the spec for MemcachedOpsRequest
type MemcachedOpsRequestSpec struct {
	// Specifies the Memcached reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type MemcachedOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading Memcached
	UpdateVersion *MemcachedUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for upgrading Memcached
	// Deprecated: use UpdateVersion
	Upgrade *MemcachedUpdateVersionSpec `json:"upgrade,omitempty"`
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
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=Upgrade;UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS
// ENUM(Upgrade, UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS)
type MemcachedOpsRequestType string

// MemcachedReplicaReadinessCriteria is the criteria for checking readiness of a Memcached pod
// after updating, horizontal scaling etc.
type MemcachedReplicaReadinessCriteria struct{}

type MemcachedUpdateVersionSpec struct {
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MemcachedOpsRequestList is a list of MemcachedOpsRequests
type MemcachedOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MemcachedOpsRequest CRD objects
	Items []MemcachedOpsRequest `json:"items,omitempty"`
}
