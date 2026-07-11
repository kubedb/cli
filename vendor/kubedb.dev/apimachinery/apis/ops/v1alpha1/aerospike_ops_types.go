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
	ResourceCodeAerospikeOpsRequest     = "asops"
	ResourceKindAerospikeOpsRequest     = "AerospikeOpsRequest"
	ResourceSingularAerospikeOpsRequest = "aerospikeopsrequest"
	ResourcePluralAerospikeOpsRequest   = "aerospikeopsrequests"
)

// AerospikeOpsRequest defines an Aerospike DBA operation.
// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=aerospikeopsrequests,singular=aerospikeopsrequest,shortName=asops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type AerospikeOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AerospikeOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus        `json:"status,omitempty"`
}

type AerospikeOpsRequestSpec struct {
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	Type        AerospikeOpsRequestType   `json:"type"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *AerospikeVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec     `json:"restart,omitempty"`
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// AerospikeVerticalScalingSpec is the spec for Aerospike vertical scaling.
type AerospikeVerticalScalingSpec struct {
	Aerospike *PodResources       `json:"aerospike,omitempty"`
	Exporter  *ContainerResources `json:"exporter,omitempty"`
	// Mode selects how the vertical scaling is actuated. Defaults to Restart.
	// +optional
	// +kubebuilder:default=Restart
	Mode VerticalScalingMode `json:"mode,omitempty"`
}

// +kubebuilder:validation:Enum=VerticalScaling;Restart
// ENUM(VerticalScaling, Restart)
type AerospikeOpsRequestType string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AerospikeOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AerospikeOpsRequest `json:"items,omitempty"`
}
