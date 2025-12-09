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
	ResourceCodeZooKeeperOpsRequest     = "zkops"
	ResourceKindZooKeeperOpsRequest     = "ZooKeeperOpsRequest"
	ResourceSingularZooKeeperOpsRequest = "zookeeperopsrequest"
	ResourcePluralZooKeeperOpsRequest   = "zookeeperopsrequests"
)

// ZooKeeperDBOpsRequest defines a ZooKeeper DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=zookeeperopsrequests,singular=zookeeperopsrequest,shortName=zkops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ZooKeeperOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ZooKeeperOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus        `json:"status,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;RotateAuth
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, RotateAuth)
type ZooKeeperOpsRequestType string

// ZooKeeperOpsRequestSpec is the spec for ZooKeeperOpsRequest
type ZooKeeperOpsRequestSpec struct {
	// Specifies the ZooKeeper reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type ZooKeeperOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading zookeeper
	UpdateVersion *ZooKeeperUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *ZooKeeperHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *ZooKeeperVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *ZooKeeperVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of zookeeper
	Configuration *ZooKeeperCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
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

// ZooKeeperpdateVersionSpec contains the update version information of a zookeeper cluster
type ZooKeeperUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// ZooKeeperReplicaReadinessCriteria is the criteria for checking readiness of a ZooKeeper pod
// after updating, horizontal scaling etc.
type ZooKeeperReplicaReadinessCriteria struct{}

// ZooKeeperHorizontalScalingSpec contains the horizontal scaling information of a ZooKeeper cluster
type ZooKeeperHorizontalScalingSpec struct {
	// Number of replicas
	Replicas *int32 `json:"replicas,omitempty"`
}

// ZooKeeperVerticalScalingSpec contains the vertical scaling information of a ZooKeeper cluster
type ZooKeeperVerticalScalingSpec struct {
	// Resource spec for nodes
	Node *PodResources `json:"node,omitempty"`
}

// ZooKeeperVolumeExpansionSpec is the spec for ZooKeeper volume expansion
type ZooKeeperVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for nodes
	Node *resource.Quantity `json:"node,omitempty"`
}

// ZooKeeperCustomConfigurationSpec is the spec for Reconfiguring the zookeeper Settings
type ZooKeeperCustomConfigurationSpec struct {
	// ConfigSecret is an optional field to provide custom configuration file for database.
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
	// ApplyConfig is an optional field to provide zookeeper configuration.
	// Provided configuration will be applied to config files stored in ConfigSecret.
	// If the ConfigSecret is missing, the operator will create a new k8s secret by the
	// following naming convention: {db-name}-user-config .
	// Expected input format:
	//	applyConfig:
	//		zookeeper.conf: |
	//			key=value
	// +optional
	ApplyConfig map[string]string `json:"applyConfig,omitempty"`
	// If set to "true", the user provided configuration will be removed.
	// The zookeeper cluster will start will default configuration that is generated by the operator.
	// +optional
	RemoveCustomConfig bool `json:"removeCustomConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZooKeeperOpsRequestList is a list of ZooKeeperOpsRequests
type ZooKeeperOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ZooKeeperOpsRequest CRD objects
	Items []ZooKeeperOpsRequest `json:"items,omitempty"`
}
