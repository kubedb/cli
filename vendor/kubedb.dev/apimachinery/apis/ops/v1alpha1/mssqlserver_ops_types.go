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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeMSSQLServerOpsRequest     = "msops"
	ResourceKindMSSQLServerOpsRequest     = "MSSQLServerOpsRequest"
	ResourceSingularMSSQLServerOpsRequest = "mssqlserveropsrequest"
	ResourcePluralMSSQLServerOpsRequest   = "mssqlserveropsrequests"
)

// MSSQLServerDBOpsRequest defines a MSSQLServer DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mssqlserveropsrequests,singular=mssqlserveropsrequest,shortName=msops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MSSQLServerOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MSSQLServerOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus          `json:"status,omitempty"`
}

// MSSQLServerOpsRequestSpec is the spec for MSSQLServerOpsRequest
type MSSQLServerOpsRequestSpec struct {
	// Specifies the MSSQLServer reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type MSSQLServerOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading MSSQL
	UpdateVersion *MSSQLServerUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *MSSQLServerHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *MSSQLServerVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *MSSQLServerVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of MSSQLServer
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *MSSQLServerTLSSpec `json:"tls,omitempty"`
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

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;RotateAuth
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, RotateAuth)
type MSSQLServerOpsRequestType string

// MSSQLServerReplicaReadinessCriteria is the criteria for checking readiness of a MSSQLServer pod
// after updating, horizontal scaling etc.
type MSSQLServerReplicaReadinessCriteria struct{}

// MSSQLServerUpdateVersionSpec contains the update version information of a MSSQLServer cluster
type MSSQLServerUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// MSSQLServerHorizontalScalingSpec contains the horizontal scaling information of a MSSQLServer cluster
type MSSQLServerHorizontalScalingSpec struct {
	// Number of Replicas of MSSQLServer Availability Group
	Replicas *int32 `json:"replicas,omitempty"`
}

// MSSQLServerVerticalScalingSpec contains the vertical scaling information of a MSSQLServer cluster
type MSSQLServerVerticalScalingSpec struct {
	MSSQLServer *PodResources       `json:"mssqlserver,omitempty"`
	Exporter    *ContainerResources `json:"exporter,omitempty"`
	Coordinator *ContainerResources `json:"coordinator,omitempty"`
}

// MSSQLServerVolumeExpansionSpec is the spec for MSSQLServer volume expansion
type MSSQLServerVolumeExpansionSpec struct {
	// volume specification for MSSQLServer
	MSSQLServer *resource.Quantity  `json:"mssqlserver,omitempty"`
	Mode        VolumeExpansionMode `json:"mode"`
}

type MSSQLServerTLSSpec struct {
	// SQLServerTLSSpec contains updated tls configurations for client and server.
	// +optional
	dbapi.MSSQLServerTLSConfig `json:",inline,omitempty"`

	// RotateCertificates tells operator to initiate certificate rotation
	// +optional
	RotateCertificates bool `json:"rotateCertificates,omitempty"`

	// Remove tells operator to remove TLS configuration
	// +optional
	Remove bool `json:"remove,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MSSQLServerOpsRequestList is a list of MSSQLServerOpsRequests
type MSSQLServerOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MSSQLServerOpsRequest CRD objects
	Items []MSSQLServerOpsRequest `json:"items,omitempty"`
}
