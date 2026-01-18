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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodePgBouncerOpsRequest     = "pbops"
	ResourceKindPgBouncerOpsRequest     = "PgBouncerOpsRequest"
	ResourceSingularPgBouncerOpsRequest = "pgbounceropsrequest"
	ResourcePluralPgBouncerOpsRequest   = "pgbounceropsrequests"
)

// PgBouncerOpsRequest defines a PgBouncer DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pgbounceropsrequests,singular=pgbounceropsrequest,shortName=pbops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PgBouncerOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PgBouncerOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus        `json:"status,omitempty"`
}

// PgBouncerOpsRequestSpec is the spec for PgBouncerOpsRequest
type PgBouncerOpsRequestSpec struct {
	// Specifies the PgBouncer reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type PgBouncerOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading PgBouncer
	UpdateVersion *PgBouncerUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *PgBouncerHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *PgBouncerVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for custom configuration of PgBouncer
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *PgBouncerTLSSpec `json:"tls,omitempty"`
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

// +kubebuilder:validation:Enum=HorizontalScaling;VerticalScaling;UpdateVersion;Reconfigure;RotateAuth;Restart;ReconfigureTLS
// ENUM(HorizontalScaling, VerticalScaling, UpdateVersion, Reconfigure, RotateAuth, Restart, ReconfigureTLS)
type PgBouncerOpsRequestType string

type PgBouncerUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// HorizontalScaling is the spec for PgBouncer horizontal scaling
type PgBouncerHorizontalScalingSpec struct {
	Replicas *int32 `json:"replicas,omitempty"`
}

// PgBouncerVerticalScalingSpec is the spec for PgBouncer vertical scaling
type PgBouncerVerticalScalingSpec struct {
	PgBouncer *PodResources       `json:"pgbouncer,omitempty"`
	Exporter  *ContainerResources `json:"exporter,omitempty"`
}

type PgBouncerTLSSpec struct {
	TLSSpec `json:",inline,omitempty"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	// +optional
	SSLMode dbapi.PgBouncerSSLMode `json:"sslMode,omitempty"`

	// ClientAuthMode for sidecar or sharding. (default will be md5. [md5;scram])
	// +optional
	ClientAuthMode dbapi.PgBouncerClientAuthMode `json:"clientAuthMode,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgBouncerOpsRequestList is a list of PgBouncerOpsRequests
type PgBouncerOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PgBouncerOpsRequest CRD objects
	Items []PgBouncerOpsRequest `json:"items,omitempty"`
}
