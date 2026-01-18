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
	ResourceCodePgpoolOpsRequest     = "ppops"
	ResourceKindPgpoolOpsRequest     = "PgpoolOpsRequest"
	ResourceSingularPgpoolOpsRequest = "pgpoolopsrequest"
	ResourcePluralPgpoolOpsRequest   = "pgpoolopsrequests"
)

// PgpoolDBOpsRequest defines a Pgpool DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pgpoolopsrequests,singular=pgpoolopsrequest,shortName=ppops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PgpoolOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PgpoolOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus     `json:"status,omitempty"`
}

// PgpoolOpsRequestSpec is the spec for PgpoolOpsRequest
type PgpoolOpsRequestSpec struct {
	// Specifies the Pgpool reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type PgpoolOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading pgpool
	UpdateVersion *PgpoolUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *PgpoolHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *PgpoolVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for custom configuration of Pgpool
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *PgpoolTLSSpec `json:"tls,omitempty"`
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

type PgpoolTLSSpec struct {
	TLSSpec `json:",inline,omitempty"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	// +optional
	SSLMode v1alpha2.PgpoolSSLMode `json:"sslMode,omitempty"`

	// ClientAuthMode for both standalone and clusters. (default will be md5. [md5;scram;cert])
	// +optional
	ClientAuthMode v1alpha2.PgpoolClientAuthMode `json:"clientAuthMode,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;VerticalScaling;Reconfigure;Restart;HorizontalScaling;ReconfigureTLS;RotateAuth
// ENUM(UpdateVersion, Restart, Reconfigure, VerticalScaling, HorizontalScaling, ReconfigureTLS, RotateAuth)
type PgpoolOpsRequestType string

// PgpoolUpdateVersionSpec contains the update version information of a pgpool cluster
type PgpoolUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// PgpoolHorizontalScalingSpec contains the horizontal scaling information of a Pgpool cluster
type PgpoolHorizontalScalingSpec struct {
	// Number of node
	Node *int32 `json:"node,omitempty"`
}

// PgpoolVerticalScalingSpec contains the vertical scaling information of a Pgpool cluster
type PgpoolVerticalScalingSpec struct {
	// Resource spec for nodes
	Node *PodResources `json:"node,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgpoolOpsRequestList is a list of PgpoolOpsRequests
type PgpoolOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PgpoolOpsRequest CRD objects
	Items []PgpoolOpsRequest `json:"items,omitempty"`
}
