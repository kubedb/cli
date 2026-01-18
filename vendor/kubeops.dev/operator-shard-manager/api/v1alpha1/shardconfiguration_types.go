/*
Copyright AppsCode Inc. and Contributors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindShardConfiguration = "ShardConfiguration"
	ResourceShardConfiguration     = "shardconfiguration"
	ResourceShardConfigurations    = "shardconfigurations"
)

// ShardConfigurationSpec defines the desired state of ShardConfiguration.
type ShardConfigurationSpec struct {
	// +kubebuilder:validation:MinItems=1
	Controllers []kmapi.TypedObjectReference `json:"controllers,omitempty"`
	// +kubebuilder:validation:MinItems=1
	Resources []ResourceInfo `json:"resources,omitempty"`
}

type ResourceInfo struct {
	kmapi.TypeReference          `json:",inline"`
	ShardKey                     *string `json:"shardKey,omitempty"`
	UseCooperativeShardMigration bool    `json:"useCooperativeShardMigration,omitempty"`
}

type ControllerAllocation struct {
	kmapi.TypedObjectReference `json:",inline"`
	Pods                       []string `json:"pods,omitempty"`
}

// ShardConfigurationStatus defines the observed state of ShardConfiguration.
type ShardConfigurationStatus struct {
	// Specifies the current phase of the App
	// +optional
	Phase Phase `json:"phase,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=12
	Conditions []kmapi.Condition `json:"conditions,omitempty"`

	Controllers []ControllerAllocation `json:"controllers,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Current;Failed
type Phase string

const (
	PhasePending Phase = "Pending"
	PhaseCurrent Phase = "Current"
	PhaseFailed  Phase = "Failed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// ShardConfiguration is the Schema for the shardconfigurations API.
type ShardConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ShardConfigurationSpec   `json:"spec,omitempty"`
	Status ShardConfigurationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// ShardConfigurationList contains a list of ShardConfiguration.
type ShardConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ShardConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ShardConfiguration{}, &ShardConfigurationList{})
}
