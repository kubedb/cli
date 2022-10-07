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
	ResourceCodePostgresAutoscaler     = "pgscaler"
	ResourceKindPostgresAutoscaler     = "PostgresAutoscaler"
	ResourceSingularPostgresAutoscaler = "postgresautoscaler"
	ResourcePluralPostgresAutoscaler   = "postgresautoscalers"
)

// PostgresAutoscaler holds the configuration for autoscaling a Postgres database.
// which automatically manages pod resources based on historical and
// real time resource utilization.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=postgresautoscalers,singular=postgresautoscaler,shortName=pgscaler,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
type PostgresAutoscaler struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification for the behaviour of the autoscaler.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status.
	// +optional
	Spec PostgresAutoscalerSpec `json:"spec,omitempty"`

	// status is the current information about the autoscaler.
	// +optional
	Status AutoscalerStatus `json:"status,omitempty"`
}

// PostgresAutoscalerSpec describes the desired functionality of the PostgresAutoscaler.
type PostgresAutoscalerSpec struct {
	DatabaseRef *core.LocalObjectReference `json:"databaseRef"`

	// This field will be used to control the behaviour of ops-manager
	OpsRequestOptions *PostgresOpsRequestOptions `json:"opsRequestOptions,omitempty"`

	Compute *PostgresComputeAutoscalerSpec `json:"compute,omitempty"`
	Storage *PostgresStorageAutoscalerSpec `json:"storage,omitempty"`
}

type PostgresComputeAutoscalerSpec struct {
	Postgres *ComputeAutoscalerSpec `json:"postgres,omitempty"`
}

type PostgresStorageAutoscalerSpec struct {
	Postgres *StorageAutoscalerSpec `json:"postgres,omitempty"`
}

type PostgresOpsRequestOptions struct {
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply opsapi.ApplyOption `json:"apply,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresAutoscalerList is a list of horizontal pod autoscaler objects.
type PostgresAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is the list of horizontal pod autoscaler objects.
	Items []PostgresAutoscaler `json:"items"`
}
