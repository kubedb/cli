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
	ofstv1 "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeDocumentDBOpsRequest     = "dcops"
	ResourceKindDocumentDBOpsRequest     = "DocumentDBOpsRequest"
	ResourceSingularDocumentDBOpsRequest = "documentdbopsrequest"
	ResourcePluralDocumentDBOpsRequest   = "documentdbopsrequests"
)

// DocumentDBOpsRequest defines a DocumentDB DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=documentdbopsrequests,singular=documentdbopsrequest,shortName=dcops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type DocumentDBOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DocumentDBOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus         `json:"status,omitempty"`
}

type DocumentDBTLSSpec struct {
	TLSSpec `json:",inline,omitempty"`
}

// DocumentDBOpsRequestSpec is the spec for DocumentDBOpsRequest
type DocumentDBOpsRequestSpec struct {
	// Specifies the DocumentDB reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type DocumentDBOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading DocumentDB
	UpdateVersion *DocumentDBUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *DocumentDBHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *DocumentDBVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *DocumentDBVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of DocumentDB
	Configuration *DocumentDBCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *DocumentDBTLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Try to reconnect standby's with primary
	ReconnectStandby *DocumentDBReconnectStandby `json:"reconnectStandby,omitempty"`
	// Forcefully do a failover to the given candidate
	ForceFailOver *DocumentDBForceFailOver `json:"forceFailOver,omitempty"`
	// Set given key pairs to raft storage
	SetRaftKeyPair *DocumentDBSetRaftKeyPair `json:"setRaftKeyPair,omitempty"`
	// Specifies information necessary for migrating storageClass or data
	Migration *DocumentDBMigrationSpec `json:"migration,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;RotateAuth;ReconnectStandby;ForceFailOver;SetRaftKeyPair;StorageMigration
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, RotateAuth, ReconnectStandby, ForceFailOver, SetRaftKeyPair, StorageMigration)
type DocumentDBOpsRequestType string

type DocumentDBUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// +kubebuilder:validation:Enum=Synchronous;Asynchronous
type DocumentDBStreamingMode string

const (
	SynchronousDocumentDBStreamingMode  DocumentDBStreamingMode = "Synchronous"
	AsynchronousDocumentDBStreamingMode DocumentDBStreamingMode = "Asynchronous"
)

// +kubebuilder:validation:Enum=Hot;Warm
type DocumentDBStandbyMode string

const (
	HotDocumentDBStandbyMode  DocumentDBStandbyMode = "Hot"
	WarmDocumentDBStandbyMode DocumentDBStandbyMode = "Warm"
)

type DocumentDBPrimaryCandidate string

// HorizontalScaling is the spec for DocumentDB horizontal scaling
type DocumentDBHorizontalScalingSpec struct {
	Replicas *int32 `json:"replicas,omitempty"`
	// Standby mode
	// +kubebuilder:default="Hot"
	StandbyMode *DocumentDBStandbyMode `json:"standbyMode,omitempty"`

	// Streaming mode
	// +kubebuilder:default="Asynchronous"
	StreamingMode *DocumentDBStreamingMode `json:"streamingMode,omitempty"`

	// +optional
	ReadReplicas []DocumentDBReadReplicaHzScalingSpec `json:"readReplicas,omitempty"`
}

type DocumentDBReadReplicaHzScalingSpec struct {
	// Name specifies the name of the read replica
	Name string `json:"name"`
	// Number of instances to deploy for a DocumentDB database.
	Replicas *int32 `json:"replicas,omitempty"`
	// Compute Resources required by the sidecar container.
	// +optional
	Resources core.ResourceRequirements `json:"resources,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// +optional
	// +mapType=atomic
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty"`
	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`
	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
	// PodPlacementPolicy is the reference of the podPlacementPolicy
	// +kubebuilder:default={name:"default"}
	// +optional
	PodPlacementPolicy *core.LocalObjectReference `json:"podPlacementPolicy,omitempty"`
	// ServiceTemplate is an optional configuration for services used to expose database
	// +optional
	ServiceTemplate *ofstv1.ServiceTemplateSpec `json:"serviceTemplate,omitempty"`
	// +optional
	Remove bool `json:"remove,omitempty"`
}

// DocumentDBVerticalScalingSpec is the spec for DocumentDB vertical scaling
type DocumentDBVerticalScalingSpec struct {
	DocumentDB   *PodResources                    `json:"documentdb,omitempty"`
	Exporter     *ContainerResources              `json:"exporter,omitempty"`
	Coordinator  *ContainerResources              `json:"coordinator,omitempty"`
	Arbiter      *PodResources                    `json:"arbiter,omitempty"`
	ReadReplicas []DocumentDBReadReplicaResources `json:"readReplicas,omitempty"`
}

type DocumentDBReadReplicaResources struct {
	DocumentDB *PodResources `json:"documentdb,omitempty"`
	Name       string        `json:"name,omitempty"`
}

type DocumentDBMigrationSpec struct {
	StorageClassName   *string                            `json:"storageClassName"`
	OldPVReclaimPolicy core.PersistentVolumeReclaimPolicy `json:"oldPVReclaimPolicy,omitempty"`
}

// DocumentDBVolumeExpansionSpec is the spec for DocumentDB volume expansion
type DocumentDBVolumeExpansionSpec struct {
	// volume specification for DocumentDB
	DocumentDB *resource.Quantity  `json:"documentdb,omitempty"`
	Arbiter    *resource.Quantity  `json:"arbiter,omitempty"`
	Mode       VolumeExpansionMode `json:"mode"`
}

type DocumentDBCustomConfigurationSpec struct {
	Tuning              *DocumentDBTuningConfig `json:"tuning,omitempty"`
	ReconfigurationSpec `json:",inline,omitempty"`
}

type DocumentDBCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty"`
	Data      map[string]string          `json:"data,omitempty"`
	Remove    bool                       `json:"remove,omitempty"`
}

type DocumentDBReconnectStandby struct {
	// ReadyTimeOut is the time to wait for standby`s to become ready
	// +optional
	ReadyTimeOut *metav1.Duration `json:"readyTimeOut,omitempty"`
}

type DocumentDBForceFailOver struct {
	Candidates []DocumentDBPrimaryCandidate `json:"candidates,omitempty"`
}

type DocumentDBSetRaftKeyPair struct {
	KeyPair map[string]string `json:"keyPair,omitempty"`
}

// DocumentDBTuningConfig defines configuration for DocumentDB performance tuning
type DocumentDBTuningConfig struct {
	// Profile defines a predefined tuning profile for different workload types.
	// +optional
	Profile *DocumentDBProfile `json:"profile,omitempty"`

	// MaxConnections defines the maximum number of concurrent connections.
	// +optional
	MaxConnections *int32 `json:"maxConnections,omitempty"`

	// StorageType defines the type of storage for tuning purposes.
	// +optional
	StorageType *DocumentDBStorageType `json:"storageType,omitempty"`

	// DisableAutoTune disables automatic tuning entirely.
	// +optional
	DisableAutoTune bool `json:"disableAutoTune,omitempty"`
}

// DocumentDBProfile defines predefined tuning profiles
// +kubebuilder:validation:Enum=web;oltp;dw;mixed;desktop
type DocumentDBProfile string

const (
	DocumentDBTuningProfileWeb     DocumentDBProfile = "web"
	DocumentDBTuningProfileOLTP    DocumentDBProfile = "oltp"
	DocumentDBTuningProfileDW      DocumentDBProfile = "dw"
	DocumentDBTuningProfileMixed   DocumentDBProfile = "mixed"
	DocumentDBTuningProfileDesktop DocumentDBProfile = "desktop"
)

// DocumentDBStorageType defines storage types for tuning purposes
// +kubebuilder:validation:Enum=ssd;hdd;san
type DocumentDBStorageType string

const (
	DocumentDBStorageTypeSSD DocumentDBStorageType = "ssd"
	DocumentDBStorageTypeHDD DocumentDBStorageType = "hdd"
	DocumentDBStorageTypeSAN DocumentDBStorageType = "san"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DocumentDBOpsRequestList is a list of DocumentDBOpsRequests
type DocumentDBOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of DocumentDBOpsRequest CRD objects
	Items []DocumentDBOpsRequest `json:"items,omitempty"`
}
