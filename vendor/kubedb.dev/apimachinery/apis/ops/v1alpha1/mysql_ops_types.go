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
	ResourceCodeMySQLOpsRequest     = "myops"
	ResourceKindMySQLOpsRequest     = "MySQLOpsRequest"
	ResourceSingularMySQLOpsRequest = "mysqlopsrequest"
	ResourcePluralMySQLOpsRequest   = "mysqlopsrequests"
)

// MySQLOpsRequest defines a MySQL DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mysqlopsrequests,singular=mysqlopsrequest,shortName=myops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MySQLOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MySQLOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus    `json:"status,omitempty"`
}

// MySQLOpsRequestSpec is the spec for MySQLOpsRequest
type MySQLOpsRequestSpec struct {
	// Specifies the MySQL reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type MySQLOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading MySQL
	UpdateVersion *MySQLUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for upgrading MySQL
	// Deprecated: use UpdateVersion
	Upgrade *MySQLUpdateVersionSpec `json:"upgrade,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *MySQLHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *MySQLVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *MySQLVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of MySQL
	Configuration *MySQLCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *MySQLTLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=Upgrade;UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS
// ENUM(Upgrade, UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS)
type MySQLOpsRequestType string

// MySQLReplicaReadinessCriteria is the criteria for checking readiness of a MySQL pod
// after updating, horizontal scaling etc.
type MySQLReplicaReadinessCriteria struct{}

type MySQLUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                         `json:"targetVersion,omitempty"`
	ReadinessCriteria *MySQLReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

type MySQLHorizontalScalingSpec struct {
	// Number of nodes/members of the group
	Member *int32 `json:"member,omitempty"`
}

type MySQLVerticalScalingSpec struct {
	MySQL       *core.ResourceRequirements `json:"mysql,omitempty"`
	Exporter    *core.ResourceRequirements `json:"exporter,omitempty"`
	Coordinator *core.ResourceRequirements `json:"coordinator,omitempty"`
}

// MySQLVolumeExpansionSpec is the spec for MySQL volume expansion
type MySQLVolumeExpansionSpec struct {
	MySQL *resource.Quantity `json:"mysql,omitempty"`
	// +kubebuilder:default="Online"
	Mode *VolumeExpansionMode `json:"mode,omitempty"`
}

type MySQLCustomConfigurationSpec struct {
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty"`
	InlineConfig       string                     `json:"inlineConfig,omitempty"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty"`
}

type MySQLTLSSpec struct {
	TLSSpec `json:",inline,omitempty"`

	// Indicates that the database server need to be encrypted connections(ssl)
	// +optional
	RequireSSL *bool `json:"requireSSL,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLOpsRequestList is a list of MySQLOpsRequests
type MySQLOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MySQLOpsRequest CRD objects
	Items []MySQLOpsRequest `json:"items,omitempty"`
}
