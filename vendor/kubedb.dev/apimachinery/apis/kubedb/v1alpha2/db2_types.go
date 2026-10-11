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

package v1alpha2

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeDB2     = "db2"
	ResourceKindDB2     = "DB2"
	ResourceSingularDB2 = "db2"
	ResourcePluralDB2   = "db2s"
)

// DB2 is the Schema for the db2s API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=db2s,singular=db2,shortName=db2,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".apiVersion"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type DB2 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DB2Spec   `json:"spec,omitempty"`
	Status DB2Status `json:"status,omitempty"`
}

// DB2Spec defines the desired state of DB2.
type DB2Spec struct {
	// Version of DB2 to be deployed.
	Version string `json:"version,omitempty"`

	// Number of instances to deploy for a DB2 database.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Init is used to initialize the database from a script or git repo.
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// HADR configures High Availability Disaster Recovery replication.
	// Only meaningful when Replicas > 1; setting it with Replicas == 1 is rejected.
	// +optional
	HADR *DB2HADRSpec `json:"hadr,omitempty"`
}

// DB2HADRSpec configures Db2 HADR. Pod ordinal 0 is the first primary, ordinal 1
// the principal standby and ordinals 2+ auxiliary standbys, which Db2 forces to
// SUPERASYNC.
//
// HADR is per-database in Db2, not per-instance: every database is configured,
// seeded, started and taken over on its own. Only the databases listed in
// Databases are protected. A database created by hand on the primary is not
// replicated and would not survive a failover; the operator reports such
// databases in the HADRUnreplicatedDatabases condition rather than enrolling
// them behind the user's back.
type DB2HADRSpec struct {
	// Databases are the databases to replicate. All of them fail over together,
	// so every one is always primary on the same pod.
	// +optional
	Databases []DB2HADRDatabase `json:"databases,omitempty"`

	// DatabaseName is the single database to replicate.
	//
	// Deprecated: use Databases. When Databases is empty this is treated as a
	// one-entry list; when Databases is set it is ignored.
	// +optional
	DatabaseName string `json:"databaseName,omitempty"`

	// SyncMode applies to the principal standby. Auxiliary standbys are always
	// SUPERASYNC regardless of this value.
	// +kubebuilder:validation:Enum=SYNC;NEARSYNC;ASYNC;SUPERASYNC
	// +kubebuilder:default=NEARSYNC
	// +optional
	SyncMode DB2HADRSyncMode `json:"syncMode,omitempty"`

	// TimeoutSeconds is HADR_TIMEOUT: how long a member waits without hearing
	// from its peer before it treats the connection as lost. A node that dies
	// silently is only noticed this long after its last heartbeat.
	// +kubebuilder:default=60
	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`

	// PeerWindowSeconds is HADR_PEER_WINDOW. For SYNC/NEARSYNC it must be at
	// least timeoutSeconds + 60. Db2 counts the peer window from the last
	// heartbeat, but notices a silent failure only timeoutSeconds later: with a
	// shorter window the standby goes straight to REMOTE_CATCHUP_PENDING, and a
	// lossless takeover is no longer possible. The extra 60 seconds cover the
	// primary lease running out and the operator promoting the standby.
	// +kubebuilder:default=120
	// +optional
	PeerWindowSeconds int32 `json:"peerWindowSeconds,omitempty"`

	// Seed controls how a standby is initialized from the primary.
	// +optional
	Seed DB2HADRSeedSpec `json:"seed,omitempty"`
}

// DB2HADRDatabase is one HADR-replicated database.
type DB2HADRDatabase struct {
	// Name is the Db2 database name: a letter followed by up to 7 letters or
	// digits. It is stored upper-case, as Db2 reports it.
	Name string `json:"name"`

	// Port is this database's HADR_LOCAL_SVC / HADR_REMOTE_SVC. Db2 needs a
	// distinct port for every HADR database in an instance. Defaulted to the
	// lowest free port from 55000, and immutable once set: moving a running
	// database to another port would disconnect its standbys.
	// +optional
	Port int32 `json:"port,omitempty"`
}

type DB2HADRSyncMode string

const (
	DB2HADRSyncModeSync       DB2HADRSyncMode = "SYNC"
	DB2HADRSyncModeNearSync   DB2HADRSyncMode = "NEARSYNC"
	DB2HADRSyncModeAsync      DB2HADRSyncMode = "ASYNC"
	DB2HADRSyncModeSuperAsync DB2HADRSyncMode = "SUPERASYNC"
)

type DB2HADRSeedSpec struct {
	// Method is the seed transport. Only Stream is implemented; SharedVolume and
	// ObjectStorage are deliberately absent from the enum until they are built,
	// because accepting a method the operator cannot act on is worse than not
	// offering it.
	// +kubebuilder:validation:Enum=Stream
	// +kubebuilder:default=Stream
	// +optional
	Method DB2HADRSeedMethod `json:"method,omitempty"`
}

type DB2HADRSeedMethod string

const (
	DB2HADRSeedMethodStream DB2HADRSeedMethod = "Stream"
)

// DB2Status defines the observed state of DB2.
type DB2Status struct {
	// Specifies the current phase of the database
	// +optional
	Phase DatabasePhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DB2List contains a list of DB2
type DB2List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DB2 `json:"items"`
}
