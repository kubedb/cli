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
	ResourceCodeDB2OpsRequest     = "db2ops"
	ResourceKindDB2OpsRequest     = "DB2OpsRequest"
	ResourceSingularDB2OpsRequest = "db2opsrequest"
	ResourcePluralDB2OpsRequest   = "db2opsrequests"
)

// DB2OpsRequest defines a DB2 DBA operation.
// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=db2opsrequests,singular=db2opsrequest,shortName=db2ops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type DB2OpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DB2OpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus  `json:"status,omitempty"`
}

type DB2OpsRequestSpec struct {
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	Type        DB2OpsRequestType         `json:"type"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *DB2VerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec     `json:"restart,omitempty"`
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// DB2VerticalScalingSpec is the spec for DB2 vertical scaling.
type DB2VerticalScalingSpec struct {
	DB2      *PodResources       `json:"db2,omitempty"`
	Exporter *ContainerResources `json:"exporter,omitempty"`
	// Mode selects how the vertical scaling is actuated. Defaults to Restart.
	// +optional
	// +kubebuilder:default=Restart
	Mode VerticalScalingMode `json:"mode,omitempty"`
}

// +kubebuilder:validation:Enum=VerticalScaling;Restart
// ENUM(VerticalScaling, Restart)
type DB2OpsRequestType string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DB2OpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DB2OpsRequest `json:"items,omitempty"`
}
