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
	kubedbApiV1Alpha2 "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeCassandraOpsRequest     = "casops"
	ResourceKindCassandraOpsRequest     = "CassandraOpsRequest"
	ResourceSingularCassandraOpsRequest = "cassandraopsrequest"
	ResourcePluralCassandraOpsRequest   = "cassandraopsrequests"
)

// CassandraDBOpsRequest defines a Cassandra DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=cassandraopsrequests,singular=cassandraopsrequest,shortName=casops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type CassandraOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CassandraOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus        `json:"status,omitempty"`
}

// CassandraOpsRequestSpec is the spec for CassandraOpsRequest
type CassandraOpsRequestSpec struct {
	// Specifies information necessary for custom configuration of Cassandra
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies the Cassandra reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type CassandraOpsRequestType `json:"type"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *CassandraHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for upgrading cassandra
	UpdateVersion *CassandraUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *CassandraVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *CassandraVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Keystore encryption secret
	// +optional
	KeystoreCredSecret *kubedbApiV1Alpha2.SecretReference `json:"keystoreCredSecret,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;VerticalScaling;Restart;VolumeExpansion;HorizontalScaling;Reconfigure;ReconfigureTLS;RotateAuth
// ENUM(UpdateVersion, VerticalScaling, Restart, VolumeExpansion, HorizontalScaling, Reconfigure, ReconfigureTLS, RotateAuth)
type CassandraOpsRequestType string

// CassandraHorizontalScalingSpec contains the horizontal scaling information of a Cassandra cluster
type CassandraHorizontalScalingSpec struct {
	// Number of node
	Node *int32 `json:"node,omitempty"`
}

// CassandraUpdateVersionSpec contains the update version information of a cassandra cluster
type CassandraUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// CassandraVerticalScalingSpec contains the vertical scaling information of a Cassandra cluster
type CassandraVerticalScalingSpec struct {
	// Resource spec for nodes
	Node *PodResources `json:"node,omitempty"`
}

// CassandraVolumeExpansionSpec is the spec for Cassandra volume expansion
type CassandraVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for nodes
	Node *resource.Quantity `json:"node,omitempty"`
}

// CassandraCustomConfigurationSpec is the spec for Reconfiguring the cassandra Settings
type CassandraCustomConfigurationSpec struct {
	// ConfigSecret is an optional field to provide custom configuration file for database.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
	// ApplyConfig is an optional field to provide cassandra configuration.
	// Provided configuration will be applied to config files stored in ConfigSecret.
	// If the ConfigSecret is missing, the operator will create a new k8s secret by the
	// following naming convention: {db-name}-user-config .
	// Expected input format:
	//	applyConfig:
	//		cassandra.conf: |
	//			key=value
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
	// If set to "true", the user provided configuration will be removed.
	// The cassandra cluster will start will default configuration that is generated by the operator.
	// +optional
	RemoveCustomConfig bool `json:"removeCustomConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CassandraOpsRequestList is a list of CassandraOpsRequests
type CassandraOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of CassandraOpsRequest CRD objects
	Items []CassandraOpsRequest `json:"items,omitempty"`
}
