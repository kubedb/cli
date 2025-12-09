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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodePostgresOpsRequest     = "pgops"
	ResourceKindPostgresOpsRequest     = "PostgresOpsRequest"
	ResourceSingularPostgresOpsRequest = "postgresopsrequest"
	ResourcePluralPostgresOpsRequest   = "postgresopsrequests"
)

// PostgresOpsRequest defines a PostgreSQL DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=postgresopsrequests,singular=postgresopsrequest,shortName=pgops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PostgresOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PostgresOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus       `json:"status,omitempty"`
}
type PostgresTLSSpec struct {
	TLSSpec `json:",inline,omitempty"`

	// SSLMode for both standalone and clusters. [disable;allow;prefer;require;verify-ca;verify-full]
	// +optional
	SSLMode dbapi.PostgresSSLMode `json:"sslMode,omitempty"`

	// ClientAuthMode for sidecar or sharding. (default will be md5. [md5;scram;cert])
	// +optional
	ClientAuthMode dbapi.PostgresClientAuthMode `json:"clientAuthMode,omitempty"`
}

// PostgresOpsRequestSpec is the spec for PostgresOpsRequest
type PostgresOpsRequestSpec struct {
	// Specifies the Postgres reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type PostgresOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading Postgres
	UpdateVersion *PostgresUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *PostgresHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *PostgresVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *PostgresVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of Postgres
	Configuration *PostgresCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *PostgresTLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Try to reconnect standby's with primary
	ReconnectStandby *PostgresReconnectStandby `json:"reconnectStandby,omitempty"`
	// Forcefully do a failover to the given candidate
	ForceFailOver *PostgresForceFailOver `json:"forceFailOver,omitempty"`
	// Set given key pairs to raft storage
	SetRaftKeyPair *PostgresSetRaftKeyPair `json:"setRaftKeyPair,omitempty"`
	// Specifies information necessary for migrating storageClass or data
	Migration *PostgresMigrationSpec `json:"migration,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// +kubebuilder:validation:Enum=Upgrade;UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;RotateAuth;ReconnectStandby;ForceFailOver;SetRaftKeyPair;StorageMigration
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, RotateAuth, ReconnectStandby, ForceFailOver, SetRaftKeyPair, StorageMigration)
type PostgresOpsRequestType string

type PostgresUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// +kubebuilder:validation:Enum=Synchronous;Asynchronous
type PostgresStreamingMode string

const (
	SynchronousPostgresStreamingMode  PostgresStreamingMode = "Synchronous"
	AsynchronousPostgresStreamingMode PostgresStreamingMode = "Asynchronous"
)

// +kubebuilder:validation:Enum=Hot;Warm
type PostgresStandbyMode string

const (
	HotPostgresStandbyMode  PostgresStandbyMode = "Hot"
	WarmPostgresStandbyMode PostgresStandbyMode = "Warm"
)

type PostgresPrimaryCandidate string

// HorizontalScaling is the spec for Postgres horizontal scaling
type PostgresHorizontalScalingSpec struct {
	Replicas *int32 `json:"replicas,omitempty"`
	// Standby mode
	// +kubebuilder:default="Hot"
	StandbyMode *PostgresStandbyMode `json:"standbyMode,omitempty"`

	// Streaming mode
	// +kubebuilder:default="Asynchronous"
	StreamingMode *PostgresStreamingMode `json:"streamingMode,omitempty"`
}

// PostgresVerticalScalingSpec is the spec for Postgres vertical scaling
type PostgresVerticalScalingSpec struct {
	Postgres    *PodResources       `json:"postgres,omitempty"`
	Exporter    *ContainerResources `json:"exporter,omitempty"`
	Coordinator *ContainerResources `json:"coordinator,omitempty"`
	Arbiter     *PodResources       `json:"arbiter,omitempty"`
}

type PostgresMigrationSpec struct {
	StorageClassName   *string                            `json:"storageClassName"`
	OldPVReclaimPolicy core.PersistentVolumeReclaimPolicy `json:"oldPVReclaimPolicy,omitempty"`
}

// PostgresVolumeExpansionSpec is the spec for Postgres volume expansion
type PostgresVolumeExpansionSpec struct {
	// volume specification for Postgres
	Postgres *resource.Quantity  `json:"postgres,omitempty"`
	Arbiter  *resource.Quantity  `json:"arbiter,omitempty"`
	Mode     VolumeExpansionMode `json:"mode"`
}

type PostgresCustomConfigurationSpec struct {
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty"`
	ApplyConfig        map[string]string          `json:"applyConfig,omitempty"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty"`
}

type PostgresCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty"`
	Data      map[string]string          `json:"data,omitempty"`
	Remove    bool                       `json:"remove,omitempty"`
}

type PostgresReconnectStandby struct {
	// ReadyTimeOut is the time to wait for standby`s to become ready
	// +optional
	ReadyTimeOut *metav1.Duration `json:"readyTimeOut,omitempty"`
}

type PostgresForceFailOver struct {
	Candidates []PostgresPrimaryCandidate `json:"candidates,omitempty"`
}

type PostgresSetRaftKeyPair struct {
	KeyPair map[string]string `json:"keyPair,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresOpsRequestList is a list of PostgresOpsRequests
type PostgresOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PostgresOpsRequest CRD objects
	Items []PostgresOpsRequest `json:"items,omitempty"`
}
