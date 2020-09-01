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

package v1alpha1

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodePostgres     = "pg"
	ResourceKindPostgres     = "Postgres"
	ResourceSingularPostgres = "postgres"
	ResourcePluralPostgres   = "postgreses"
)

// Postgres defines a Postgres database.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=postgreses,singular=postgres,shortName=pg,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Postgres struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              PostgresSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            PostgresStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type PostgresSpec struct {
	// Version of Postgres to be deployed.
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`

	// Number of instances to deploy for a Postgres database.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// Standby mode
	StandbyMode *PostgresStandbyMode `json:"standbyMode,omitempty" protobuf:"bytes,3,opt,name=standbyMode,casttype=PostgresStandbyMode"`

	// Streaming mode
	StreamingMode *PostgresStreamingMode `json:"streamingMode,omitempty" protobuf:"bytes,4,opt,name=streamingMode,casttype=PostgresStreamingMode"`

	// Archive for wal files
	Archiver *PostgresArchiverSpec `json:"archiver,omitempty" protobuf:"bytes,5,opt,name=archiver"`

	// Leader election configuration
	// +optional
	LeaderElection *LeaderElectionConfig `json:"leaderElection,omitempty" protobuf:"bytes,6,opt,name=leaderElection"`

	// Database authentication secret
	DatabaseSecret *core.SecretVolumeSource `json:"databaseSecret,omitempty" protobuf:"bytes,7,opt,name=databaseSecret"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty" protobuf:"bytes,8,opt,name=storageType,casttype=StorageType"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty" protobuf:"bytes,9,opt,name=storage"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty" protobuf:"bytes,10,opt,name=init"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty" protobuf:"bytes,12,opt,name=monitor"`

	// ConfigSource is an optional field to provide custom configuration file for database (i.e postgresql.conf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty" protobuf:"bytes,13,opt,name=configSource"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty" protobuf:"bytes,14,opt,name=podTemplate"`

	// ServiceTemplate is an optional configuration for service used to expose database
	// +optional
	ServiceTemplate ofst.ServiceTemplateSpec `json:"serviceTemplate,omitempty" protobuf:"bytes,15,opt,name=serviceTemplate"`

	// ReplicaServiceTemplate is an optional configuration for service used to expose postgres replicas
	// +optional
	ReplicaServiceTemplate ofst.ServiceTemplateSpec `json:"replicaServiceTemplate,omitempty" protobuf:"bytes,16,opt,name=replicaServiceTemplate"`

	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	UpdateStrategy apps.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty" protobuf:"bytes,17,opt,name=updateStrategy"`

	// TLS contains tls configurations for client and server.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty" protobuf:"bytes,18,opt,name=tls"`

	// Indicates that the database is paused and controller will not sync any changes made to this spec.
	// +optional
	Paused bool `json:"paused,omitempty" protobuf:"varint,19,opt,name=paused"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty" protobuf:"varint,20,opt,name=halted"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty" protobuf:"bytes,21,opt,name=terminationPolicy,casttype=TerminationPolicy"`
}

// +kubebuilder:validation:Enum=server;archiver;metrics-exporter
type PostgresCertificateAlias string

const (
	PostgresServerCert          PostgresCertificateAlias = "server"
	PostgresArchiverCert        PostgresCertificateAlias = "archiver"
	PostgresMetricsExporterCert PostgresCertificateAlias = "metrics-exporter"
)

type PostgresArchiverSpec struct {
	Storage *store.Backend `json:"storage,omitempty" protobuf:"bytes,1,opt,name=storage"`
	// wal_keep_segments
}

type PostgresStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=DatabasePhase"`
	Reason string        `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PostgresList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Postgres CRD objects
	Items []Postgres `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}

type PostgresWALSourceSpec struct {
	BackupName    string          `json:"backupName,omitempty" protobuf:"bytes,1,opt,name=backupName"`
	PITR          *RecoveryTarget `json:"pitr,omitempty" protobuf:"bytes,2,opt,name=pitr"`
	store.Backend `json:",inline,omitempty" protobuf:"bytes,3,opt,name=backend"`
}

type RecoveryTarget struct {
	// TargetTime specifies the time stamp up to which recovery will proceed.
	TargetTime string `json:"targetTime,omitempty" protobuf:"bytes,1,opt,name=targetTime"`
	// TargetTimeline specifies recovering into a particular timeline.
	// The default is to recover along the same timeline that was current when the base backup was taken.
	TargetTimeline string `json:"targetTimeline,omitempty" protobuf:"bytes,2,opt,name=targetTimeline"`
	// TargetXID specifies the transaction ID up to which recovery will proceed.
	TargetXID string `json:"targetXID,omitempty" protobuf:"bytes,3,opt,name=targetXID"`
	// TargetInclusive specifies whether to include ongoing transaction in given target point.
	TargetInclusive *bool `json:"targetInclusive,omitempty" protobuf:"varint,4,opt,name=targetInclusive"`
}

// +kubebuilder:validation:Enum=Hot;Warm
type PostgresStandbyMode string

const (
	HotPostgresStandbyMode  PostgresStandbyMode = "Hot"
	WarmPostgresStandbyMode PostgresStandbyMode = "Warm"

	// Deprecated
	DeprecatedHotStandby PostgresStandbyMode = "hot"
	// Deprecated
	DeprecatedWarmStandby PostgresStandbyMode = "warm"
)

// +kubebuilder:validation:Enum=Synchronous;Asynchronous
type PostgresStreamingMode string

const (
	SynchronousPostgresStreamingMode  PostgresStreamingMode = "Synchronous"
	AsynchronousPostgresStreamingMode PostgresStreamingMode = "Asynchronous"

	// Deprecated
	DeprecatedAsynchronousStreaming PostgresStreamingMode = "asynchronous"
)
