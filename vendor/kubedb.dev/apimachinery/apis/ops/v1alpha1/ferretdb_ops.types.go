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
	"kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeFerretDBOpsRequest     = "frops"
	ResourceKindFerretDBOpsRequest     = "FerretDBOpsRequest"
	ResourceSingularFerretDBOpsRequest = "ferretdbopsrequest"
	ResourcePluralFerretDBOpsRequest   = "ferretdbopsrequests"
)

// FerretDBDBOpsRequest defines a FerretDB DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=ferretdbopsrequests,singular=ferretdbopsrequest,shortName=frops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type FerretDBOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FerretDBOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus       `json:"status,omitempty"`
}

// FerretDBOpsRequestSpec is the spec for FerretDBOpsRequest
type FerretDBOpsRequestSpec struct {
	// Specifies the FerretDB reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type FerretDBOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading ferretdb
	UpdateVersion *FerretDBUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *FerretDBHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *FerretDBVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *FerretDBTLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

type FerretDBTLSSpec struct {
	TLSSpec `json:",inline,omitempty"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	// +optional
	SSLMode v1alpha2.SSLMode `json:"sslMode,omitempty"`

	// ClientAuthMode for both standalone and clusters. (default will be md5. [md5;scram;cert])
	// +optional
	ClientAuthMode v1alpha2.ClusterAuthMode `json:"clientAuthMode,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;VerticalScaling;Restart;HorizontalScaling;ReconfigureTLS;RotateAuth
// ENUM(UpdateVersion, Restart, VerticalScaling, HorizontalScaling, ReconfigureTLS, RotateAuth)
type FerretDBOpsRequestType string

// FerretDBUpdateVersionSpec contains the update version information of a ferretdb cluster
type FerretDBUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// FerretDBHorizontalScalingSpec contains the horizontal scaling information of a FerretDB cluster
type FerretDBHorizontalScalingSpec struct {
	Primary   *FerretDBHorizontalScalingReplicas `json:"primary,omitempty"`
	Secondary *FerretDBHorizontalScalingReplicas `json:"secondary,omitempty"`
}

type FerretDBHorizontalScalingReplicas struct {
	Replicas *int32 `json:"replicas,omitempty"`
}

// FerretDBVerticalScalingSpec contains the vertical scaling information of a FerretDB cluster
type FerretDBVerticalScalingSpec struct {
	Primary   *PodResources `json:"primary,omitempty"`
	Secondary *PodResources `json:"secondary,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FerretDBOpsRequestList is a list of FerretDBOpsRequests
type FerretDBOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of FerretDBOpsRequest CRD objects
	Items []FerretDBOpsRequest `json:"items,omitempty"`
}
