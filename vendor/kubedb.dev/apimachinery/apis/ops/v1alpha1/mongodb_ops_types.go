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
	ResourceCodeMongoDBOpsRequest     = "mgops"
	ResourceKindMongoDBOpsRequest     = "MongoDBOpsRequest"
	ResourceSingularMongoDBOpsRequest = "mongodbopsrequest"
	ResourcePluralMongoDBOpsRequest   = "mongodbopsrequests"
)

// MongoDBOpsRequest defines a MongoDB DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mongodbopsrequests,singular=mongodbopsrequest,shortName=mgops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MongoDBOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MongoDBOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus      `json:"status,omitempty"`
}

// MongoDBOpsRequestSpec is the spec for MongoDBOpsRequest
type MongoDBOpsRequestSpec struct {
	// Specifies the MongoDB reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type MongoDBOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading MongoDB
	UpdateVersion *MongoDBUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *MongoDBHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *MongoDBVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *MongoDBVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of MongoDB
	Configuration *MongoDBCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Specifies information necessary for reprovisioning database
	Reprovision *Reprovision `json:"reprovision,omitempty"`
	// Specifies information necessary for setting up Archiver for database
	Archiver *ArchiverOptions `json:"archiver,omitempty"`
	// Horizons specifies the information for setting up replicaset horizons.
	Horizons *Horizons `json:"horizons,omitempty"`

	// Specifies the Readiness Criteria
	ReadinessCriteria *MongoDBReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// +kubebuilder:validation:Enum=Upgrade;UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;Reprovision;RotateAuth;Horizons
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, Reprovision, RotateAuth, Horizons)
type MongoDBOpsRequestType string

// MongoDBReplicaReadinessCriteria is the criteria for checking readiness of a MongoDB pod
// after restarting the pod
type MongoDBReplicaReadinessCriteria struct {
	// +kubebuilder:validation:Minimum:=0
	OplogMaxLagSeconds int32 `json:"oplogMaxLagSeconds,omitempty"`
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:Maximum:=100
	ObjectsCountDiffPercentage int32 `json:"objectsCountDiffPercentage,omitempty"`
}

type MongoDBUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// MongoDBShardNode is the spec for mongodb Shard
type MongoDBShardNode struct {
	Shards   int32 `json:"shards,omitempty"`
	Replicas int32 `json:"replicas,omitempty"`
}

// ConfigNode is the spec for mongodb ConfigServer
type ConfigNode struct {
	Replicas int32 `json:"replicas,omitempty"`
}

// MongosNode is the spec for mongodb Mongos
type MongosNode struct {
	Replicas int32 `json:"replicas,omitempty"`
}

type HiddenNode struct {
	Replicas int32 `json:"replicas,omitempty"`
}

// HorizontalScaling is the spec for mongodb horizontal scaling
type MongoDBHorizontalScalingSpec struct {
	Shard        *MongoDBShardNode `json:"shard,omitempty"`
	ConfigServer *ConfigNode       `json:"configServer,omitempty"`
	Mongos       *MongosNode       `json:"mongos,omitempty"`
	Hidden       *HiddenNode       `json:"hidden,omitempty"`
	Replicas     *int32            `json:"replicas,omitempty"`
}

// MongoDBVerticalScalingSpec is the spec for mongodb vertical scaling
type MongoDBVerticalScalingSpec struct {
	Standalone   *PodResources       `json:"standalone,omitempty"`
	ReplicaSet   *PodResources       `json:"replicaSet,omitempty"`
	Mongos       *PodResources       `json:"mongos,omitempty"`
	ConfigServer *PodResources       `json:"configServer,omitempty"`
	Shard        *PodResources       `json:"shard,omitempty"`
	Arbiter      *PodResources       `json:"arbiter,omitempty"`
	Hidden       *PodResources       `json:"hidden,omitempty"`
	Exporter     *ContainerResources `json:"exporter,omitempty"`
	Coordinator  *ContainerResources `json:"coordinator,omitempty"`
}

// MongoDBVolumeExpansionSpec is the spec for mongodb volume expansion
type MongoDBVolumeExpansionSpec struct {
	Mode         VolumeExpansionMode `json:"mode"`
	Standalone   *resource.Quantity  `json:"standalone,omitempty"`
	ReplicaSet   *resource.Quantity  `json:"replicaSet,omitempty"`
	ConfigServer *resource.Quantity  `json:"configServer,omitempty"`
	Shard        *resource.Quantity  `json:"shard,omitempty"`
	Hidden       *resource.Quantity  `json:"hidden,omitempty"`
}

type MongoDBCustomConfigurationSpec struct {
	Standalone   *MongoDBCustomConfiguration `json:"standalone,omitempty"`
	ReplicaSet   *MongoDBCustomConfiguration `json:"replicaSet,omitempty"`
	Mongos       *MongoDBCustomConfiguration `json:"mongos,omitempty"`
	ConfigServer *MongoDBCustomConfiguration `json:"configServer,omitempty"`
	Shard        *MongoDBCustomConfiguration `json:"shard,omitempty"`
	Arbiter      *MongoDBCustomConfiguration `json:"arbiter,omitempty"`
	Hidden       *MongoDBCustomConfiguration `json:"hidden,omitempty"`
}

type MongoDBCustomConfiguration struct {
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty"`
	ApplyConfig        map[string]string          `json:"applyConfig,omitempty"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty"`
}

type Horizons struct {
	DNS string `json:"dns"`
	// Pods contain the host:port for all the replicas. Its length will be same as db.spec.replicas
	// +optional
	Pods []string `json:"pods,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBOpsRequestList is a list of MongoDBOpsRequests
type MongoDBOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MongoDBOpsRequest CRD objects
	Items []MongoDBOpsRequest `json:"items,omitempty"`
}
