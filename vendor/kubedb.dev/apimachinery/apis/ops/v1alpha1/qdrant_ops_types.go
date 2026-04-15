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
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeQdrantOpsRequest     = "qdops"
	ResourceKindQdrantOpsRequest     = "QdrantOpsRequest"
	ResourceSingularQdrantOpsRequest = "qdrantopsrequest"
	ResourcePluralQdrantOpsRequest   = "qdrantopsrequests"
)

// QdrantDBOpsRequest defines a Qdrant DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=qdrantopsrequests,singular=qdrantopsrequest,shortName=qdops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type QdrantOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              QdrantOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus     `json:"status,omitempty"`
}

// QdrantOpsRequestSpec is the spec for QdrantOpsRequest
type QdrantOpsRequestSpec struct {
	// Specifies the Qdrant reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type QdrantOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading qdrant
	UpdateVersion *QdrantUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *QdrantHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *QdrantVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *QdrantVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of qdrant
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *QdrantTLSSpec `json:"tls,omitempty"`
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
type QdrantOpsRequestType string

// QdrantUpdateVersionSpec contains the update version information of a qdrant cluster
type QdrantUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// QdrantReplicaReadinessCriteria is the criteria for checking readiness of a Qdrant pod
// after updating, horizontal scaling etc.
type QdrantReplicaReadinessCriteria struct{}

// QdrantHorizontalScalingSpec contains the horizontal scaling information of a Qdrant cluster
type QdrantHorizontalScalingSpec struct {
	// Number of node
	Node *int32 `json:"node,omitempty"`
}

// QdrantVerticalScalingSpec contains the vertical scaling information of a Qdrant cluster
type QdrantVerticalScalingSpec struct {
	// Resource spec for nodes
	Node *PodResources `json:"node,omitempty"`
}

// QdrantVolumeExpansionSpec is the spec for Qdrant volume expansion
type QdrantVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for nodes
	Node *resource.Quantity `json:"node,omitempty"`
}

type QdrantTLSSpec struct {
	// +optional
	api.QdrantTLSConfig `json:",inline,omitempty"`

	// RotateCertificates tells operator to initiate certificate rotation
	// +optional
	RotateCertificates bool `json:"rotateCertificates,omitempty"`

	// Remove tells operator to remove TLS configuration
	// +optional
	Remove bool `json:"remove,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QdrantOpsRequestList is a list of QdrantOpsRequests
type QdrantOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of QdrantOpsRequest CRD objects
	Items []QdrantOpsRequest `json:"items,omitempty"`
}
