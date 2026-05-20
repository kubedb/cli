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
	ResourceCodeOracleOpsRequest     = "oracleops"
	ResourceKindOracleOpsRequest     = "OracleOpsRequest"
	ResourceSingularOracleOpsRequest = "oracleopsrequest"
	ResourcePluralOracleOpsRequest   = "oracleopsrequests"
)

// OracleDBOpsRequest defines a Oracle DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=oracleopsrequests,singular=oracleopsrequest,shortName=oraops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type OracleOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              OracleOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus     `json:"status,omitempty"`
}

// OracleOpsRequestSpec is the spec for OracleOpsRequest
type OracleOpsRequestSpec struct {
	// Specifies the Oracle reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type OracleOpsRequestType `json:"type"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *OracleVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	// Specifies information necessary for volume expansion
	VolumeExpansion *OracleVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of oracle
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Specifies information necessary for migrating storage
	Migration *OracleMigrationSpec `json:"migration,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

type OracleMigrationSpec struct {
	StorageClassName   *string                            `json:"storageClassName"`
	OldPVReclaimPolicy core.PersistentVolumeReclaimPolicy `json:"oldPVReclaimPolicy,omitempty"`
}

// +kubebuilder:validation:Enum=Restart;Reconfigure;StorageMigration;VerticalScaling;VolumeExpansion;RotateAuth
// ENUM(Restart, Reconfigure, StorageMigration, VerticalScaling, VolumeExpansion, RotateAuth)
type OracleOpsRequestType string

// OracleVerticalScalingSpec contains the vertical scaling information of an Oracle cluster
type OracleVerticalScalingSpec struct {
	// Resource spec for nodes
	Node *PodResources `json:"node,omitempty"`
}

// OracleVolumeExpansionSpec is the spec for Oracle volume expansion
type OracleVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for Oracle database nodes
	Node *resource.Quantity `json:"node,omitempty"`
	// volume specification for DataGuard observer
	Observer *resource.Quantity `json:"observer,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OracleOpsRequestList is a list of OracleOpsRequests
type OracleOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of OracleOpsRequest CRD objects
	Items []OracleOpsRequest `json:"items,omitempty"`
}
