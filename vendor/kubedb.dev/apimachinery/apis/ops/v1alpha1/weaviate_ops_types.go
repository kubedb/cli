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
	ResourceCodeWeaviateOpsRequest     = "wvops"
	ResourceKindWeaviateOpsRequest     = "WeaviateOpsRequest"
	ResourceSingularWeaviateOpsRequest = "weaviateopsrequest"
	ResourcePluralWeaviateOpsRequest   = "weaviateopsrequests"
)

// WeaviateDBOpsRequest defines a Weaviate DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=weaviateopsrequests,singular=weaviateopsrequest,shortName=wvops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type WeaviateOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WeaviateOpsRequestSpec `json:"spec,omitempty"`
	Status OpsRequestStatus       `json:"status,omitempty"`
}

// WeaviateOpsRequestSpec is the spec for WeaviateOpsRequest
type WeaviateOpsRequestSpec struct {
	// Specifies the Weaviate reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type WeaviateOpsRequestType `json:"type"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *WeaviateHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *WeaviateVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *WeaviateVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Specifies information necessary for custom configuration of weaviate
	Configuration *WeaviateReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for migrating storageClass or data
	Migration *WeaviateMigrationSpec `json:"migration,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;RotateAuth;StorageMigration
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, RotateAuth, StorageMigration)
type WeaviateOpsRequestType string

// WeaviateUpdateVersionSpec contains the update version information of a Weaviate cluster
type WeaviateUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// WeaviateReplicaReadinessCriteria is the criteria for checking readiness of a Weaviate pod
// after updating, horizontal scaling etc.
type WeaviateReplicaReadinessCriteria struct{}

// WeaviateHorizontalScalingSpec contains the horizontal scaling information of a Weaviate cluster
type WeaviateHorizontalScalingSpec struct {
	// Number of node
	Node *int32 `json:"node,omitempty"`
}

// WeaviateVerticalScalingSpec contains the vertical scaling information of a Weaviate cluster
type WeaviateVerticalScalingSpec struct {
	// Resource spec for nodes
	Node *PodResources `json:"node,omitempty"`
}

// WeaviateVolumeExpansionSpec is the spec for Weaviate volume expansion
type WeaviateVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for nodes
	Node *resource.Quantity `json:"node,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WeaviateOpsRequestList is a list of WeaviateOpsRequests
type WeaviateOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of WeaviateOpsRequest CRD objects
	Items []WeaviateOpsRequest `json:"items,omitempty"`
}

// WeaviateReconfigurationSpec is the Weaviate-specific reconfiguration spec.
// It embeds the generic ReconfigurationSpec and adds Weaviate-specific fields.
type WeaviateReconfigurationSpec struct {
	ReconfigurationSpec `json:",inline,omitempty"`
	// BackupConfigSecret is an optional field to provide environment variables
	// from a Kubernetes Secret for the database container.
	// +optional
	BackupConfigSecret *core.LocalObjectReference `json:"backupConfigSecret,omitempty"`
}

type WeaviateMigrationSpec struct {
	StorageClassName   *string                            `json:"storageClassName"`
	OldPVReclaimPolicy core.PersistentVolumeReclaimPolicy `json:"oldPVReclaimPolicy,omitempty"`
}
