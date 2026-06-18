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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeHanaDBOpsRequest     = "hdbops"
	ResourceKindHanaDBOpsRequest     = "HanaDBOpsRequest"
	ResourceSingularHanaDBOpsRequest = "hanadbopsrequest"
	ResourcePluralHanaDBOpsRequest   = "hanadbopsrequests"
)

// HanaDBOpsRequest defines a HanaDB DBA operation.
// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=hanadbopsrequests,singular=hanadbopsrequest,shortName=hdbops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type HanaDBOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HanaDBOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus     `json:"status,omitempty"`
}

type HanaDBTLSSpec struct {
	TLSSpec `json:",inline,omitempty"`
}

type HanaDBOpsRequestSpec struct {
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	Type        HanaDBOpsRequestType      `json:"type"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *HanaDBVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *HanaDBVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of HanaDB
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *HanaDBTLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart   *RestartSpec          `json:"restart,omitempty"`
	Migration *StorageMigrationSpec `json:"migration,omitempty"`
	Timeout   *metav1.Duration      `json:"timeout,omitempty"`
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// HanaDBVerticalScalingSpec is the spec for HanaDB vertical scaling.
type HanaDBVerticalScalingSpec struct {
	HanaDB      *PodResources       `json:"hanadb,omitempty"`
	Coordinator *ContainerResources `json:"coordinator,omitempty"`
	Exporter    *ContainerResources `json:"exporter,omitempty"`
}

// HanaDBVolumeExpansionSpec is the spec for HanaDB volume expansion.
type HanaDBVolumeExpansionSpec struct {
	HanaDB *resource.Quantity  `json:"hanadb,omitempty"`
	Mode   VolumeExpansionMode `json:"mode"`
}

// +kubebuilder:validation:Enum=VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;RotateAuth;StorageMigration
// ENUM(VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, RotateAuth, StorageMigration)
type HanaDBOpsRequestType string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HanaDBOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HanaDBOpsRequest `json:"items,omitempty"`
}
