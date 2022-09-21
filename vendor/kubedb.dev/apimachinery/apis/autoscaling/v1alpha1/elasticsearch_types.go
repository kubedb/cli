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
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeElasticsearchAutoscaler     = "esscaler"
	ResourceKindElasticsearchAutoscaler     = "ElasticsearchAutoscaler"
	ResourceSingularElasticsearchAutoscaler = "elasticsearchautoscaler"
	ResourcePluralElasticsearchAutoscaler   = "elasticsearchautoscalers"
)

// ElasticsearchAutoscaler is the configuration for a horizontal pod
// autoscaler, which automatically manages the replica count of any resource
// implementing the scale subresource based on the metrics specified.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=elasticsearchautoscalers,singular=elasticsearchautoscaler,shortName=esscaler,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
type ElasticsearchAutoscaler struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification for the behaviour of the autoscaler.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status.
	// +optional
	Spec ElasticsearchAutoscalerSpec `json:"spec,omitempty"`

	// status is the current information about the autoscaler.
	// +optional
	Status AutoscalerStatus `json:"status,omitempty"`
}

// ElasticsearchAutoscalerSpec is the specification of the behavior of the autoscaler.
type ElasticsearchAutoscalerSpec struct {
	DatabaseRef *core.LocalObjectReference `json:"databaseRef"`

	// This field will be used to control the behaviour of ops-manager
	OpsRequestOptions *ElasticsearchOpsRequestOptions `json:"opsRequestOptions,omitempty"`

	Compute *ElasticsearchComputeAutoscalerSpec `json:"compute,omitempty"`
	Storage *ElasticsearchStorageAutoscalerSpec `json:"storage,omitempty"`
}

type ElasticsearchComputeAutoscalerSpec struct {
	Node     *ComputeAutoscalerSpec                      `json:"node,omitempty"`
	Topology *ElasticsearchComputeTopologyAutoscalerSpec `json:"topology,omitempty"`
}

type ElasticsearchComputeTopologyAutoscalerSpec struct {
	Master *ComputeAutoscalerSpec `json:"master,omitempty"`
	Data   *ComputeAutoscalerSpec `json:"data,omitempty"`
	Ingest *ComputeAutoscalerSpec `json:"ingest,omitempty"`
}

type ElasticsearchStorageAutoscalerSpec struct {
	Node     *StorageAutoscalerSpec                      `json:"node,omitempty"`
	Topology *ElasticsearchStorageTopologyAutoscalerSpec `json:"topology,omitempty"`
}

type ElasticsearchStorageTopologyAutoscalerSpec struct {
	Master *StorageAutoscalerSpec `json:"master,omitempty"`
	Data   *StorageAutoscalerSpec `json:"data,omitempty"`
	Ingest *StorageAutoscalerSpec `json:"ingest,omitempty"`
}

type ElasticsearchOpsRequestOptions struct {
	// Specifies the Readiness Criteria
	ReadinessCriteria *opsapi.ElasticsearchReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`

	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply opsapi.ApplyOption `json:"apply,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ElasticsearchAutoscalerList is a list of ElasticsearchAutoscaler objects.
type ElasticsearchAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata"`

	// items is the list of elasticsearch database autoscaler objects.
	Items []ElasticsearchAutoscaler `json:"items"`
}
