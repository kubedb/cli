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
	ResourceCodeMariaDBOpsRequest     = "mariaops"
	ResourceKindMariaDBOpsRequest     = "MariaDBOpsRequest"
	ResourceSingularMariaDBOpsRequest = "mariadbopsrequest"
	ResourcePluralMariaDBOpsRequest   = "mariadbopsrequests"
)

// MariaDBOpsRequest defines a MariaDB DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mariadbopsrequests,singular=mariadbopsrequest,shortName=mariaops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MariaDBOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MariaDBOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus      `json:"status,omitempty"`
}

// MariaDBOpsRequestSpec is the spec for MariaDBOpsRequest
type MariaDBOpsRequestSpec struct {
	// Specifies the MariaDB reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type MariaDBOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading MariaDB
	UpdateVersion *MariaDBUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *MariaDBHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *MariaDBVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *MariaDBVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of MariaDB
	Configuration *MariaDBCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *MariaDBTLSSpec `json:"tls,omitempty"`
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

// +kubebuilder:validation:Enum=Upgrade;UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;RotateAuth
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, RotateAuth)
type MariaDBOpsRequestType string

// MariaDBReplicaReadinessCriteria is the criteria for checking readiness of an MariaDB database
type MariaDBReplicaReadinessCriteria struct{}

type MariaDBUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

type MariaDBHorizontalScalingSpec struct {
	// Specifies whether horizontal scaling is applied to the MaxScale Server.
	// When set to true, it enables horizontal scaling for the MaxScale Server.
	MaxScale bool `json:"maxscale,omitempty"`
	// Number of nodes/members of the group
	Member *int32 `json:"member,omitempty"`
	// specifies the weight of the current member/PodResources
	MemberWeight int32 `json:"memberWeight,omitempty"`
}

type MariaDBVerticalScalingSpec struct {
	MariaDB     *PodResources       `json:"mariadb,omitempty"`
	MaxScale    *PodResources       `json:"maxscale,omitempty"`
	Exporter    *ContainerResources `json:"exporter,omitempty"`
	Coordinator *ContainerResources `json:"coordinator,omitempty"`
}

// MariaDBVolumeExpansionSpec is the spec for MariaDB volume expansion
type MariaDBVolumeExpansionSpec struct {
	MariaDB  *resource.Quantity  `json:"mariadb,omitempty"`
	MaxScale *resource.Quantity  `json:"maxscale,omitempty"`
	Mode     VolumeExpansionMode `json:"mode"`
}

type MariaDBCustomConfigurationSpec struct {
	// ConfigSecret is an optional field to provide custom configuration file for database.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
	// If set to "true", the user provided configuration will be removed.
	// MariaDB will start will default configuration that is generated by the operator.
	// +optional
	RemoveCustomConfig bool `json:"removeCustomConfig,omitempty"`
	// ApplyConfig is an optional field to provide MariaDB configuration.
	// Provided configuration will be applied to config files stored in ConfigSecret.
	// If the ConfigSecret is missing, the operator will create a new k8s secret by the
	// following naming convention: {db-name}-user-config .
	// Expected input format:
	//	applyConfig:
	//		file-name.cnf: |
	//			[mysqld]
	//			key1: value1
	//			key2: value2
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
}

type MariaDBCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty"`
	Data      map[string]string          `json:"data,omitempty"`
	Remove    bool                       `json:"remove,omitempty"`
}

type MariaDBTLSSpec struct {
	TLSSpec `json:",inline,omitempty"`

	// Indicates that the database server need to be encrypted connections(ssl)
	// +optional
	RequireSSL *bool `json:"requireSSL,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MariaDBOpsRequestList is a list of MariaDBOpsRequests
type MariaDBOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MariaDBOpsRequest CRD objects
	Items []MariaDBOpsRequest `json:"items,omitempty"`
}
