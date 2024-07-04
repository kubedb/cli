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
)

const (
	ResourceKindNodeTopology = "NodeTopology"
	ResourceNodeTopology     = "nodetopology"
	ResourceNodeTopologies   = "nodetopologies"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=nodetopologies,singular=nodetopology,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type NodeTopology struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NodeTopologySpec `json:"spec,omitempty"`
}

type NodeTopologySpec struct {
	// +optional
	Description string `json:"description"`
	// +optional
	NodeSelectionPolicy NodeSelectionPolicy `json:"nodeSelectionPolicy,omitempty"`
	TopologyKey         string              `json:"topologyKey"`
	NodeGroups          []NodeGroup         `json:"nodeGroups,omitempty"`

	// Requirements are layered with GetLabels and applied to every node.
	// +kubebuilder:validation:XValidation:message="requirements with operator 'In' must have a value defined",rule="self.all(x, x.operator == 'In' ? x.values.size() != 0 : true)"
	// +kubebuilder:validation:XValidation:message="requirements operator 'Gt' or 'Lt' must have a single positive integer value",rule="self.all(x, (x.operator == 'Gt' || x.operator == 'Lt') ? (x.values.size() == 1 && int(x.values[0]) >= 0) : true)"
	// +kubebuilder:validation:MaxItems:=30
	// +optional
	Requirements []core.NodeSelectorRequirement `json:"requirements,omitempty" hash:"ignore"`
}

type NodeGroup struct {
	TopologyValue string `json:"topologyValue"`
	// Allocatable represents the total resources of a node.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#capacity
	// Deprecated: Use resources instead.
	// +optional
	Allocatable core.ResourceList `json:"allocatable,omitempty"`
	// Resources represents the requested and limited resources of a machine type.
	// +optional
	Resources core.ResourceRequirements `json:"resources,omitempty"`
	// Cost is the cost of the running an ondeamd machine for a month
	// +optional
	Cost *ResourceCost `json:"cost,omitempty"`
}

type ResourceCost struct {
	Unit  string `json:"unit"`
	Price string `json:"price"`
}

type NodeSelectionPolicy string

const (
	NodeSelectionPolicyLabelSelector NodeSelectionPolicy = "LabelSelector"
	NodeSelectionPolicyTaint         NodeSelectionPolicy = "Taint"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type NodeTopologyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeTopology `json:"items,omitempty"`
}
